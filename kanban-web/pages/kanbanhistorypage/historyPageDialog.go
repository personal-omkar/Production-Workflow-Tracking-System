package kanbanhistorypage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	u "irpl.com/kanban-commons/utils"

	m "irpl.com/kanban-commons/model"
	s "irpl.com/kanban-web/services"
)

// BuildCompDetailsDialog handles the API request
func BuildCompDetailsDialog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	var RootData m.KbRoot
	err := json.NewDecoder(r.Body).Decode(&RootData)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.Println("Failed to decode request:", err)
		return
	}
	data, err := json.Marshal(RootData)
	if err != nil {
		http.Error(w, "Unable to marshal data", http.StatusBadRequest)
		log.Println("Failed to send request:", err)
		return
	}
	url := RestURL + "/get-kbRoot-details"
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		http.Error(w, "Error creating request to REST service", http.StatusInternalServerError)
		return
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to send request to the target service", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Failed to process request on target service", http.StatusInternalServerError)
		return
	}
	var RootDetails m.DetailRootData
	err = json.NewDecoder(resp.Body).Decode(&RootDetails)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.Println("Failed to decode request:", err)
		return
	}

	var List []s.AccordionAttributes
	for _, Linedata := range RootDetails.KanbanDetails {

		tabledata := []*m.ProdLine{
			{
				Name:      Linedata.ProdLine,
				CreatedBy: u.FormatStringDate(Linedata.ProdProcesses[0].StartedOn, "date-time"),
				ModifiedBy: func() string {
					lastProcess := Linedata.ProdProcesses[len(Linedata.ProdProcesses)-1]

					// Convert CompletedOn from string to time.Time
					completedOn, err := time.Parse(time.RFC3339, lastProcess.CompletedOn)
					zeroTime := time.Time{} // Default zero value of time.Time

					if err != nil || completedOn.Equal(zeroTime) || completedOn.Year() == 1 {
						return u.FormatStringDate(lastProcess.StartedOn, "date-time")
					}
					return u.FormatStringDate(completedOn, "date-time")
				}(),
				Operator: func() string {
					for _, p := range Linedata.ProdProcesses {
						if p.Operator != "" {
							return p.Operator
						}
					}
					return "-"
				}(),
			},
		}
		var ProductionLine s.TableCard
		ProductionLine.For = "Accordion"
		ProductionLine.BodyTables = s.CardTableBody{Columns: []s.CardTableBodyHeadCol{
			{
				Lable: "Process Line",
				Name:  "Name",
				Width: "col-1",
			},
			{
				Lable: "Operator",
				Name:  "Operator",
				Width: "col-1",
			},
			{
				Lable: "Started On",
				Name:  "CreatedBy",
				Width: "col-1",
			},
			{
				Lable: "Completed On",
				Name:  "ModifiedBy",
				Width: "col-1",
			},
		},
			ColumnsWidth: []string{"col-1", "col-1", "col-1", "col-1"},
			Data:         tabledata,
		}

		tabledata1 := []*m.ProdProcess{}
		for x, processdata := range Linedata.ProdProcesses {

			if x == 0 || x == len(Linedata.ProdProcesses)-1 {
				continue
			}

			expctedTime := ""
			if processdata.ExpectedMeanTime == "" {
				expctedTime = "N/A"
			} else {
				expctedTime = processdata.ExpectedMeanTime + " min"
			}

			startedTime, errStart := time.Parse(time.RFC3339Nano, processdata.StartedOn)
			completedTime, errCompleted := time.Parse(time.RFC3339Nano, processdata.CompletedOn)

			actualTime := "N/A"
			if errStart == nil && errCompleted == nil {
				// Truncate both times to seconds before calculating the difference
				startedTime = startedTime.Truncate(time.Second)
				completedTime = completedTime.Truncate(time.Second)
				if completedTime.Before(startedTime) {
					actualTime = "In process"
				} else {
					totalSeconds := int(completedTime.Sub(startedTime).Seconds())
					if totalSeconds == 0 {
						actualTime = "0.00 min"
					} else {
						minutes := totalSeconds / 60
						seconds := totalSeconds % 60
						actualTime = fmt.Sprintf("%d.%02d min", minutes, seconds)
					}
				}
			}

			ProcessTable := &m.ProdProcess{
				Name:             processdata.ProcessName,
				CreatedBy:        u.FormatStringDate(processdata.StartedOn, "date-time"),
				ModifiedBy:       u.FormatStringDate(processdata.CompletedOn, "date-time"),
				ExpectedMeanTime: expctedTime,
				Status:           actualTime,
			}

			tabledata1 = append(tabledata1, ProcessTable)
		}

		var ProdProcesses s.TableCard
		ProdProcesses.For = "Accordion"
		ProdProcesses.CardHeading = "Production Processes"
		ProdProcesses.BodyTables = s.CardTableBody{Columns: []s.CardTableBodyHeadCol{
			{
				Lable: "Process Name",
				Name:  "Name",
				Width: "col-1",
			},
			{
				Lable: "Started On",
				Name:  "CreatedBy",
				Width: "col-1",
			},
			{
				Lable: "Completed On",
				Name:  "ModifiedBy",
				Width: "col-1",
			},
			{
				Lable: "Expected Time",
				Name:  "expected_mean_time",
				Width: "col-1",
			},
			{
				Lable: "Actual Time",
				Name:  "Status",
				Width: "col-1",
			},
		},
			ColumnsWidth: []string{"col-1", "col-1", "col-1", "col-1", "col-1", "col-1", "col-1"},
			Data:         tabledata1,
		}

		var AccordionList s.AccordionAttributes
		AccordionList.Name = RootDetails.KanbanDetails[0].ProdLine
		AccordionList.Sections = []s.FormSection{
			{
				ID:         "Test",
				Fields:     []s.FormField{},
				ExtraField: ProductionLine.Build() + ProdProcesses.Build(),
			},
		}

		List = append(List, AccordionList)
	}

	KanbanLen := len(RootDetails.KanbanDetails) - 1
	kanbanProcessLen := len(RootDetails.KanbanDetails[KanbanLen].ProdProcesses) - 1

	// If the Quality test fail don't show Dispatch Section
	var dispatchDetails s.Details
	if RootDetails.KanbanStatus != "-1" {

		dispatchDetails = s.Details{
			Lable: "Dispatch Details",
			Style: "",
			Note:  RootDetails.DispatchNote,
			DetailsData: []s.DetailsData{
				{Heading: "",
					Style: "color:#ab71a2;",
					Data: []s.Data{
						{Lable: "Remark", Value: map[string]string{"-1": "Quality Fail", "2": "Quality Test Pending", "3": "Dispatch Pending", "4": "Dispatched"}[RootDetails.KanbanStatus]},
					},
				},
				{Heading: "",
					Data: []s.Data{
						{Lable: "Completed On", Value: map[string]string{"3": "Dispatch Pending", "4": u.FormatStringDate(RootDetails.DispatchDoneTime, "date-time")}[RootDetails.KanbanStatus]},
					},
				},
				{Heading: "",
					Data: []s.Data{
						{Lable: "Operator", Value: RootDetails.PackingOperator},
					},
				},
			},
		}
	}

	orderDetailAccordion := s.ModelCard{
		ID:      "viewKanbanDetails",
		Type:    "modal-xl",
		Heading: "Kanban Details",
		Form: s.ModelForm{
			Details: s.Details{
				Lable: "Kanban Details",
				Style: "",
				DetailsData: []s.DetailsData{
					{Heading: "Vendor Details",
						Style: "color:#ab71a2;",
						Data: []s.Data{
							{Lable: "Vendor Name", Value: RootDetails.VendorName},
							{Lable: "Vendor Code", Value: RootDetails.VendorCode},
							{Lable: "Contact Info", Value: RootDetails.ContactInfo},
						},
					},
					{Heading: "Order Details",
						Style: "color:#ab71a2;",
						Data: []s.Data{
							{Lable: "Part Name", Value: RootDetails.CompoundName},
							{Lable: "Cell Number", Value: RootDetails.CellNo},
							{Lable: "Number of Lots", Value: strconv.Itoa(RootDetails.NoOFLots)},
							{Lable: "Kanban Number", Value: RootDetails.KanbanNo},
						},
					},
					{Heading: "&nbsp;",
						Data: []s.Data{
							{Lable: "Status", Value: RootDetails.Status},
							{Lable: "Lot Number", Value: RootDetails.LotNo},
							{Lable: "Order ID", Value: RootDetails.OrderID},
						},
					},
				},
			},
			Accordion: []s.AccordionList{
				{
					Lable: "Kanban Process Details",
					List:  List,
				},
			},
			QualityDetails: s.Details{
				Lable: "Quality Testing Details",
				Style: "",
				Note:  RootDetails.QualityNote,
				DetailsData: []s.DetailsData{
					{Heading: "",
						Style: "color:#ab71a2;",
						Data: []s.Data{
							{Lable: "Remark", Value: map[string]string{"-1": "Quality Fail", "2": "Quality Test Pending", "3": "Quality Pass", "4": "Quality Pass"}[RootDetails.KanbanStatus]},
							{Lable: "Operator", Value: RootDetails.QualityOperator},
						},
					},
					{Heading: "",
						Style: "color:#ab71a2;",
						Data: []s.Data{
							{Lable: "Started On", Value: map[string]string{"-1": u.FormatStringDate(RootDetails.KanbanDetails[KanbanLen].ProdProcesses[kanbanProcessLen].StartedOn, "date-time"), "2": u.FormatStringDate(RootDetails.KanbanDetails[KanbanLen].ProdProcesses[kanbanProcessLen].StartedOn, "date-time"), "3": u.FormatStringDate(RootDetails.KanbanDetails[KanbanLen].ProdProcesses[kanbanProcessLen].StartedOn, "date-time"), "0": "", "1": "", "4": u.FormatStringDate(RootDetails.KanbanDetails[KanbanLen].ProdProcesses[kanbanProcessLen].StartedOn, "date-time")}[RootDetails.KanbanStatus]},
						},
					},
					{Heading: "",
						Data: []s.Data{
							{Lable: "Completed On", Value: u.FormatStringDate(RootDetails.QualityDoneTime, "date-time")},
						},
					},
				},
			},
			DispatchDetails: dispatchDetails,
		},
	}
	AccordionDialogBox := orderDetailAccordion.Build()
	response := map[string]string{
		"dialogHTML": AccordionDialogBox,
	}
	json.NewEncoder(w).Encode(response)
}
