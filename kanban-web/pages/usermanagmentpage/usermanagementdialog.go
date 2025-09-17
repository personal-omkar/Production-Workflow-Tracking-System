package usermanagmentpage

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	m "irpl.com/kanban-commons/model"
	s "irpl.com/kanban-web/services"
)

func (u *UserManagement) CardBuild(email string) string {
	var userstatus = []s.DropDownOptions{{Text: "False", Value: "false"}, {Text: "True", Value: "true"}}
	var userRoleslist []s.DropDownOptions
	var userRoles []m.UserRoles
	var user []m.UserManagement
	var userrole string
	//fetching compounds records
	userroleresp, err := http.Get(RestURL + "/get-all-user-roles")
	if err != nil {
		slog.Error("%s - error - %s", "Error making GET request", err)
	}

	defer userroleresp.Body.Close()

	if err := json.NewDecoder(userroleresp.Body).Decode(&userRoles); err != nil {
		slog.Error("error decoding response body", "error", err)
	}
	//fetching compounds records
	userresp, err := http.Get(RestURL + "/get-user-details-by-email?email=" + email)
	if err != nil {
		slog.Error("%s - error - %s", "Error making GET request", err)
	}

	defer userresp.Body.Close()

	if err := json.NewDecoder(userresp.Body).Decode(&user); err != nil {
		slog.Error("error decoding response body", "error", err)
	}
	for _, i := range userRoles {
		var userroles s.DropDownOptions
		userroles.Text = i.RoleName
		userroles.Value = strconv.Itoa(i.ID)

		userRoleslist = append(userRoleslist, userroles)
	}

	for i, v := range userRoleslist {
		if strconv.Itoa(user[0].RoleId) == v.Value {
			userRoleslist[i].Selected = true
			userrole = v.Text
		}
	}
	for i, v := range userstatus {
		strValue := strconv.FormatBool(user[0].Isactive)
		if strValue == v.Value {
			userstatus[i].Selected = true
		}
	}
	EditUserModel := s.ModelCard{
		ID:      "EditUserModel",
		Type:    "modal-md",
		Heading: "Edit User",
		Form: s.ModelForm{FormID: "EditUserModel",
			FormAction: "", Footer: s.Footer{CancelBtn: true, Buttons: []s.FooterButtons{{BtnType: "button", DataSubmitName: "EditUserModel", BtnID: "edit-user-submit", Text: "Save"}}},
			Inputfield: []s.InputAttributes{
				{Type: "hidden", DataType: "int", Name: "id", Value: strconv.Itoa(user[0].UserID), ID: "showUserId", Label: `User ID`, Width: "w-100", Hidden: true},
				{Type: "text", DataType: "text", Name: "username", Value: user[0].UserName, ID: "showUserName", Label: `User Name`, Width: "w-100", Required: true},
				{Type: "text", DataType: "text", Name: "email", Value: user[0].Email, ID: "showEmail", Label: `Email`, Width: "w-100", Disabled: true, Required: true},
				{Type: "password", DataType: "text", Name: "password", Value: user[0].Password, ID: "showPassword", Label: `Password`, Width: "w-100", Required: true},
				// {Type: "text", DataType: "text", Name: "vendors", Value: user[0].VendorCode, ID: "showVendorCode", Label: `Vendor Code`, Width: "w-100", Hidden: true},
			},
			Dropdownfield: []s.DropdownAttributes{
				{Label: "Role Name", DataType: "int", Name: "roleId", ID: "showRoleName", Options: userRoleslist, Width: "w-100"},
				{Label: "Is Active", DataType: "text", Name: "isactive", Options: userstatus, ID: "showIsActive", Width: "w-100"},
			},
		},
	}
	if userrole == "Customer" || userrole == "Operator" {
		EditUserModel.Form.Inputfield = append(EditUserModel.Form.Inputfield, s.InputAttributes{Type: "text", DataType: "text", Name: "vendors", Value: user[0].VendorCode, ID: "showVendorCode", Label: `Vendor Code`, Width: "w-100"})
	} else {
		EditUserModel.Form.Inputfield = append(EditUserModel.Form.Inputfield, s.InputAttributes{Type: "text", DataType: "text", Name: "vendors", Value: user[0].VendorCode, ID: "showVendorCode", Label: `Vendor Code`, Width: "w-100", Hidden: true})
	}
	return EditUserModel.Build()
}
