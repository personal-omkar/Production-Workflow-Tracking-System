package materialmanagement

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"irpl.com/kanban-commons/model"
	s "irpl.com/kanban-web/services"
)

func (o *MaterialManagement) CardBuild() string {

	var material []*model.RawMaterial

	var materiallist = []s.DropDownOptions{{Text: "Enabled", Value: "true"}, {Text: "Disabled", Value: "false"}}
	// fetching compounds records
	resp, err := http.Get(RestURL + "/get-material-by-param?key=id&value=" + o.MaterialId)
	if err != nil {
		slog.Error("%s - error - %s", "Error making GET request", err)
	}

	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&material); err != nil {
		slog.Error("error decoding response body", "error", err)
	}

	for i, v := range materiallist {
		strValue := strconv.FormatBool(material[0].Status)
		if strValue == v.Value {
			materiallist[i].Selected = true
		}
	}

	EditUserModel := s.ModelCard{
		ID:      "EditMaterialModel",
		Type:    "modal-m",
		Heading: "Edit Raw Material",
		Form: s.ModelForm{FormID: "EditMaterialModel",
			FormAction: "", Footer: s.Footer{CancelBtn: true, Buttons: []s.FooterButtons{{BtnType: "button", DataSubmitName: "EditMaterialModel", BtnID: "edit-material-submit", Text: "Save"}}},
			Dropdownfield: []s.DropdownAttributes{
				{Label: "Status", ID: "EditStatus", Name: "Status", DataType: "bool", Options: materiallist},
			},
			Inputfield: []s.InputAttributes{
				{Type: "text", Name: "ID", DataType: "int", ID: "EditID", Value: strconv.Itoa(material[0].Id), Width: "hidden", Required: true, Hidden: true},
				{Type: "text", Name: "Description", DataType: "string", Value: material[0].Description, ID: "EditMaterialDesc", Label: `Description`, Width: "w-100", Required: true},
				{Type: "text", Name: "SCADACode", DataType: "string", Value: material[0].SCADACode, ID: "EditMaterialSCADACode", Label: `SCADA Code`, Width: "w-100", Required: false},
				{Type: "text", Name: "SAPCode", DataType: "string", Value: material[0].SAPCode, ID: "EditMaterialSAPCode", Label: `SAP Code`, Width: "w-100", Required: false},
				{Type: "text", Name: "Comment", DataType: "string", Value: material[0].Comment, ID: "EditMaterialComment", Label: `Comment`, Width: "w-100", Required: false},
			},
		},
	}

	return EditUserModel.Build()
}
