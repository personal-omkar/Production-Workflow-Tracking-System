package services

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	m "irpl.com/kanban-commons/model"
	"irpl.com/kanban-commons/utils"
)

var visibilityStatus = []DropDownOptions{{Text: "False", Value: "false"}, {Text: "True", Value: "true"}}

var ProcessesStatus = []DropDownOptions{{Text: "In-Active", Value: "inActive"}, {Text: "Active", Value: "Active"}}

func AddProdProcessesModal(ID, Type, Heading string) string {
	// Construct the AddLineModel
	AddLineModel := ModelCard{
		ID:      ID,
		Type:    Type,
		Heading: Heading,
		Form: ModelForm{
			FormID:     "AddProdProcess",
			FormAction: "",
			Footer: Footer{
				CancelBtn: false,
				Buttons: []FooterButtons{
					{BtnType: "submit", BtnID: "SaveProcesses", Text: "Add", Disabled: true},
				},
			},
			Inputfield: []InputAttributes{
				{Type: "text", Name: "addProcessName", ID: "addProcessName", Label: `Process Name`, Width: "w-100", Required: true},
				{Type: "file", Name: "addProcessIcon", ID: "addProcessIcon", Label: `Icon `, Width: "w-100", AdditionalAttr: `accept="image/*"`},
				// {Type: "text", Name: "addProcessIcon", ID: "addProcessIcon", Label: `Icon Link`, Width: "w-100", Required: true},
				{Type: "text", Name: "addProcessTime", ID: "addProcessTime", Label: `Expected Mean Time`, Width: "w-100", Required: true},
			},
			TextArea: []TextAreaAttributes{
				{Label: "Process Description", DataType: "text", Name: "processDescription", ID: "processDescription"},
			},
			Dropdownfield: []DropdownAttributes{
				{Label: "Process Status", DataType: "text", Name: "isactive", Options: ProcessesStatus, ID: "processIsActive", Width: "w-100"},
				{Label: "Process Visisbility", DataType: "text", Name: "processIsVisibalse", Options: visibilityStatus, ID: "processIsVisibalse", Width: "w-100"},
			},
		},
	}

	return AddLineModel.Build()
}

func EditProdProcessDialog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		log.Println("Failed to read request body:", err)
		return
	}

	var prodProcessData m.ProdProcess
	err = json.Unmarshal(body, &prodProcessData)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.Println("Failed to decode request:", err)
		return
	}

	data, err := json.Marshal(prodProcessData)
	if err != nil {
		http.Error(w, "Unable to marshal data", http.StatusBadRequest)
		log.Println("Failed to send request:", err)
		return
	}

	url := utils.RestURL + "/get-prod-process-by-param"
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
	var ProdProcess []m.ProdProcess
	err = json.NewDecoder(resp.Body).Decode(&ProdProcess)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.Println("Failed to decode request:", err)
		return
	}

	line_visibility_status := ProdProcess[0].LineVisibility
	if line_visibility_status {
		visibilityStatus = []DropDownOptions{{Text: "False", Value: "false"}, {Text: "True", Value: "true", Selected: true}}
	} else {
		visibilityStatus = []DropDownOptions{{Text: "False", Value: "false", Selected: true}, {Text: "True", Value: "true"}}

	}

	prodprocessesStatus := ProdProcess[0].Status
	if prodprocessesStatus == "Active" {
		ProcessesStatus = []DropDownOptions{{Text: "In-Active", Value: "inActive"}, {Text: "Active", Value: "Active", Selected: true}}
	} else {
		ProcessesStatus = []DropDownOptions{{Text: "In-Active", Value: "inActive", Selected: true}, {Text: "Active", Value: "Active"}}
	}
	ReadOnly := false
	if (strings.TrimSpace(ProdProcess[0].Name) == "Line up" && ProdProcess[0].CreatedBy == "system") || (strings.TrimSpace(ProdProcess[0].Name) == "Packing" && ProdProcess[0].CreatedBy == "system") {
		ReadOnly = true
	}
	EditProcessModel := ModelCard{
		ID:      "EditProductionProcess",
		Type:    "",
		Heading: "Edit Production Process",
		Form: ModelForm{
			FormID:     "EditProdProcess",
			FormAction: "",
			Footer: Footer{
				CancelBtn: false,
				Buttons: []FooterButtons{
					{BtnType: "button", BtnID: "closeDialog", Text: "Close"},
					{BtnType: "submit", BtnID: "UpdateProcesses", Text: "Update"},
				},
			},
			Inputfield: []InputAttributes{
				{Type: "text", Name: "editProcessID", ID: "editProcessID", Label: `Process ID`, Width: "w-100", Required: true, Value: strconv.Itoa(ProdProcess[0].Id), Hidden: true, Readonly: true, Disabled: true},
				{Type: "text", Name: "editProcessName", ID: "editProcessName", Label: `Process Name`, Width: "w-100", Required: true, Value: ProdProcess[0].Name, Readonly: ReadOnly, Disabled: ReadOnly},
				{Type: "upload", Name: "editProcessIcon", ID: "editProcessIcon", Label: `Icon `, Width: "w-100", AdditionalAttr: `accept="image/*"`, Value: ProdProcess[0].Icon},
				// {Type: "text", Name: "editProcessIcon", ID: "editProcessIcon", Label: `Icon Link`, Width: "w-100", Required: true, Value: ProdProcess[0].Icon},
				{Type: "text", Name: "editProcessTime", ID: "editProcessTime", Label: `Expected Mean Time`, Width: "w-100", Required: true, Value: ProdProcess[0].ExpectedMeanTime},
			},
			TextArea: []TextAreaAttributes{
				{Label: "Process Description", DataType: "text", Name: "editprocessDescription", ID: "editprocessDescription", Value: ProdProcess[0].Description},
			},
			Dropdownfield: []DropdownAttributes{
				{Label: "Process Status", DataType: "text", Name: "editisactive", Options: ProcessesStatus, ID: "editisactive", Width: "w-100", Disabled: ReadOnly},
				{Label: "Process Visisbility", DataType: "text", Name: "editprocessIsVisibalse", Options: visibilityStatus, ID: "editprocessIsVisibalse", Width: "w-100", Disabled: ReadOnly},
			},
		},
	}

	DialogBox := EditProcessModel.Build()

	response := map[string]string{
		"dialogHTML": DialogBox,
	}
	json.NewEncoder(w).Encode(response)

}
