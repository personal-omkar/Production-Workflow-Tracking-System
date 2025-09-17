package operatormanagement

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"irpl.com/kanban-commons/model"
	s "irpl.com/kanban-web/services"
)

func (o *OperatorManagement) CardBuild() string {

	var operator []*model.Operator

	var operatorlist = []s.DropDownOptions{{Text: "Enabled", Value: "true"}, {Text: "Disabled", Value: "false"}}
	// fetching compounds records
	resp, err := http.Get(RestURL + "/get-operator-by-param?key=id&value=" + o.OperatorId)
	if err != nil {
		slog.Error("%s - error - %s", "Error making GET request", err)
	}

	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&operator); err != nil {
		slog.Error("error decoding response body", "error", err)
	}

	for i, v := range operatorlist {
		strValue := strconv.FormatBool(operator[0].Status)
		if strValue == v.Value {
			operatorlist[i].Selected = true
		}
	}

	EditUserModel := s.ModelCard{
		ID:      "EditOperatorModel",
		Type:    "modal-m",
		Heading: "Edit Operator",
		Form: s.ModelForm{FormID: "EditOperatorModel",
			FormAction: "", Footer: s.Footer{CancelBtn: true, Buttons: []s.FooterButtons{{BtnType: "button", DataSubmitName: "EditOperatorModel", BtnID: "edit-operator-submit", Text: "Save"}}},
			Dropdownfield: []s.DropdownAttributes{
				{Label: "Status", ID: "EditStatus", Name: "Status", DataType: "bool", Options: operatorlist},
			},
			Inputfield: []s.InputAttributes{
				{Type: "text", Name: "ID", DataType: "int", ID: "EditID", Value: strconv.Itoa(operator[0].Id), Width: "hidden", Required: true, Hidden: true},
				{Type: "text", Name: "OperatorName", DataType: "string", Value: operator[0].OperatorName, ID: "EditOperatorName", Label: `Operator Name`, Width: "w-100", Required: true},
				{Type: "text", Name: "OperatorCode", DataType: "string", Value: operator[0].OperatorCode, ID: "EditOperatorCode", Label: `Operator Code`, Width: "w-100", Required: true},
			},
		},
	}

	return EditUserModel.Build()
}
