package chemicalmanagement

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"irpl.com/kanban-commons/model"
	s "irpl.com/kanban-web/services"
)

func (o *ChemicalManagement) CardBuild() string {

	var chemical []*model.ChemicalTypes

	var chemicallist = []s.DropDownOptions{{Text: "Enabled", Value: "true"}, {Text: "Disabled", Value: "false"}}
	// fetching compounds records
	resp, err := http.Get(RestURL + "/get-chemical-by-param?key=id&value=" + o.ChemicalId)
	if err != nil {
		slog.Error("%s - error - %s", "Error making GET request", err)
	}

	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&chemical); err != nil {
		slog.Error("error decoding response body", "error", err)
	}

	for i, v := range chemicallist {
		strValue := strconv.FormatBool(chemical[0].Status)
		if strValue == v.Value {
			chemicallist[i].Selected = true
		}
	}

	EditUserModel := s.ModelCard{
		ID:      "EditChemicalModel",
		Type:    "modal-m",
		Heading: "Edit Chemical Type",
		Form: s.ModelForm{FormID: "EditChemicalModel",
			FormAction: "", Footer: s.Footer{CancelBtn: true, Buttons: []s.FooterButtons{{BtnType: "button", DataSubmitName: "EditChemicalModel", BtnID: "edit-chemical-submit", Text: "Save"}}},
			Dropdownfield: []s.DropdownAttributes{
				{Label: "Status", ID: "EditStatus", Name: "Status", DataType: "bool", Options: chemicallist},
			},
			Inputfield: []s.InputAttributes{
				{Type: "text", Name: "ID", DataType: "int", ID: "EditID", Value: strconv.Itoa(chemical[0].Id), Width: "hidden", Required: true, Hidden: true},
				{Type: "text", Name: "Type", DataType: "string", Value: chemical[0].Type, ID: "EditChemicalType", Label: `Chemical Type`, Width: "w-100", Required: true},
				{Type: "text", Name: "ConvCode", DataType: "string", Value: chemical[0].ConvCode, ID: "EditChemicalConv", Label: `Conv Code`, Width: "w-100", Required: true},
			},
		},
	}

	return EditUserModel.Build()
}
