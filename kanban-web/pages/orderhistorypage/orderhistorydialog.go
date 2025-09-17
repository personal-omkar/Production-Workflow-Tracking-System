package orderhistorypage

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	m "irpl.com/kanban-commons/model"
	u "irpl.com/kanban-commons/utils"
	s "irpl.com/kanban-web/services"
)

// Dialog for Order history
func BuildOrderDetailsDialog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	var kbData m.KbData
	err := json.NewDecoder(r.Body).Decode(&kbData)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.Println("Failed to decode request:", err)
		return
	}
	data, err := json.Marshal(kbData)
	if err != nil {
		http.Error(w, "Unable to marshal data", http.StatusBadRequest)
		log.Println("Failed to send request:", err)
		return
	}

	url := u.RestURL + "/get-all-details-for-order"
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
	var orderdetails m.OrderDetailsHistory
	err = json.NewDecoder(resp.Body).Decode(&orderdetails)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.Println("Failed to decode request:", err)
		return
	}

	DispatchMsg := ""
	lotsDispatch := orderdetails.NoOFLots - len(orderdetails.KanbanDetails)
	if lotsDispatch == 0 {
		DispatchMsg = ""
	} else if lotsDispatch == 1 {
		DispatchMsg = "[" + strconv.Itoa(lotsDispatch) + " lot has been successfully dispatched directly from the inventory]"
	} else {
		DispatchMsg = "[" + strconv.Itoa(lotsDispatch) + " lots  has been successfully dispatched directly from the inventory]"
	}

	var tabledata1 []*m.CustomerOrderDetails
	var tabledata2 []*m.CustomerOrderDetails

	totalEntries := len(orderdetails.KanbanDetails)
	splitIndex := totalEntries / 2

	TableClass := ""

	for i, data := range orderdetails.KanbanDetails {
		ProcessTable := &m.CustomerOrderDetails{
			Id:        data.ID,
			Location:  data.LotNo,
			CreatedBy: data.ProdLine,
		}

		if totalEntries > 2 {
			if i < splitIndex {
				tabledata1 = append(tabledata1, ProcessTable)
			} else {
				tabledata2 = append(tabledata2, ProcessTable)
			}
			TableClass = "mx-1"
		} else {
			tabledata1 = append(tabledata1, ProcessTable)
			TableClass = ""
		}
	}

	Dialopgtool := `
	<!--html-->
		<button type="button" class="btn m-0 p-0" id="viewKanbanDetails" data-toggle="tooltip" data-placement="bottom" title="View Details"> 
				 <i class="fa fa-eye mx-2" style="color: #b250ad;"></i> 
		</button>
	<!--!html-->`
	// Define the first table
	var ProductionLine1 s.TableCard
	ProductionLine1.For = "Accordion"
	ProductionLine1.Class = "text-center " + TableClass
	ProductionLine1.BodyTables = s.CardTableBody{
		Columns: []s.CardTableBodyHeadCol{
			{Lable: "Lot Number", Name: "Location", Width: "col-1"},
			{Lable: "View Details", Name: "Tools", Width: "col-1"},
		},
		ColumnsWidth: []string{"col-1", "col-1"},
		Data:         tabledata1,
		Tools:        Dialopgtool,
	}

	// Define the second table (only if needed)
	var ProductionLine2 s.TableCard
	if totalEntries > 2 {
		ProductionLine2.For = "Accordion"
		ProductionLine2.Class = "text-center " + TableClass
		ProductionLine2.BodyTables = s.CardTableBody{
			Columns: []s.CardTableBodyHeadCol{
				{Lable: "Lot Number", Name: "Location", Width: "col-1"},
				{Lable: "View Details", Name: "Tools", Width: "col-1"},
			},
			ColumnsWidth: []string{"col-1", "col-1"},
			Data:         tabledata2,
			Tools:        Dialopgtool,
		}
	}

	// Create ExtractFields dynamically
	extractFields := []s.ExtractFields{
		{
			Heading:      "Lot Details " + DispatchMsg,
			HeadingStyle: "color: #ab71a2; font-size:1.25rem;",
			HTML:         ProductionLine1.Build(),
		},
	}

	if totalEntries > 2 {
		extractFields = append(extractFields, s.ExtractFields{
			Heading:      "Lot Details " + DispatchMsg,
			HeadingStyle: "color: #ab71a2; font-size:1.25rem;",
			HTML:         ProductionLine2.Build() + ProductionLine1.Build(),
		})
	}

	viewOrderDetailsModel := s.ModelCard{
		ID:      "viewOrderDetails",
		Type:    "modal-xl",
		Heading: "Order Details  [" + orderdetails.CellNo + "]",
		Form: s.ModelForm{
			FormID: "accept_order",
			Details: s.Details{
				Style: "",
				DetailsData: []s.DetailsData{
					{Heading: "Customer Details", Style: "color:#ab71a2;", Data: []s.Data{
						{Lable: "Name", Value: orderdetails.Username},
						{Lable: "Email", Value: orderdetails.Email},
						{Lable: "Part Name", Value: orderdetails.CompoundName},
					}},
					{Heading: "Order Details", Style: "color:#ab71a2;", Data: []s.Data{
						{Lable: "Order ID", Value: orderdetails.OrderId},
						{Lable: "No. of Lots", Value: strconv.Itoa(orderdetails.NoOFLots)},
						{Lable: "Status", Value: orderdetails.Status},
					}},
					{Heading: "&nbsp;", Data: []s.Data{
						{Lable: "Dispatch On", Value: u.FormatStringDate(orderdetails.DispatchDate, "date-time")},
						{Lable: "Ordered On", Value: u.FormatStringDate(orderdetails.OrderOn, "date")},
						{Lable: "Demand Date", Value: u.FormatStringDate(orderdetails.DemandDate, "date")},
					}},
					{Heading: "Vendor Details", Style: "color:#ab71a2;", Data: []s.Data{
						{Lable: "Vendor Name", Value: orderdetails.VendorName},
						{Lable: "Vendor Code", Value: orderdetails.VendorCode},
						{Lable: "Contact Info.", Value: orderdetails.ContactInfo},
					}},
				},
			},
			ExtractFields: extractFields,
		},
	}

	DialogBox := viewOrderDetailsModel.Build()

	response := map[string]string{
		"dialogHTML": DialogBox,
		"dialogJS": u.JoinStr(
			`<script>
			 </script>
		`),
	}
	json.NewEncoder(w).Encode(response)
}
