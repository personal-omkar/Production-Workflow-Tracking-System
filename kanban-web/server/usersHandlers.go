package server

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	m "irpl.com/kanban-commons/model"
	"irpl.com/kanban-commons/utils"
	s "irpl.com/kanban-web/services"
)

// Create User
func CreateUser(w http.ResponseWriter, r *http.Request) {

	var data m.User
	var ApiResp m.ApiRespMsg
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err == nil {

		jsonValue, err := json.Marshal(data)
		if err != nil {
			slog.Error("%s - error - %s", "Error marshaling user data", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to marshal user data")
			return
		}

		resp, err := http.Post(utils.RestURL+"/create-new-user", "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			slog.Error("%s - error - %s", "Error making POST request", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to create user")
			return
		}
		defer resp.Body.Close()

		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			slog.Error("%s - error - %s", "Error reading response body", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to read response body")
			return
		}
		ApiResp.Code = resp.StatusCode
		ApiResp.Message = string(responseBody)
		body, err := json.Marshal(ApiResp)
		if err != nil {
			log.Printf("Error marshaling response: %v", err)
		}

		utils.SetResponse(w, resp.StatusCode, string(body))
	} else {

		slog.Error("%s - error - %s", "Record creation failed", err.Error())
		utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to create record")
	}
}

// Update User
func UpdateUser(w http.ResponseWriter, r *http.Request) {

	var data m.User

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err == nil {

		jsonValue, err := json.Marshal(data)
		if err != nil {
			slog.Error("%s - error - %s", "Error marshaling user data", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to marshal user data")
			return
		}

		resp, err := http.Post(utils.RestURL+"/update-user", "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			slog.Error("%s - error - %s", "Error making POST request", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to update user")
			return
		}
		defer resp.Body.Close()

		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			slog.Error("%s - error - %s", "Error reading response body", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to read response body")
			return
		}

		utils.SetResponse(w, resp.StatusCode, string(responseBody))
	} else {

		slog.Error("%s - error - %s", "Record updation failed", err.Error())
		utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to update record")
	}
}

// Delete User
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	var ids []int64 // IDs to delete
	decoder := json.NewDecoder(r.Body)

	// Decode the request body to get the list of IDs
	if err := decoder.Decode(&ids); err == nil {
		// Marshal the IDs into JSON
		jsonValue, err := json.Marshal(ids)
		if err != nil {
			slog.Error("%s - error - %s", "Error marshaling delete request data", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to marshal delete request data")
			return
		}

		// Make a POST request to the DB service for deletion
		resp, err := http.Post(utils.RestHost+"/delete-user", "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			slog.Error("%s - error - %s", "Error making POST request for delete", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to delete user(s)")
			return
		}
		defer resp.Body.Close()

		// Read the response body from the service
		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			slog.Error("%s - error - %s", "Error reading response body for delete", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to read response body for delete")
			return
		}

		// Respond to the client with the service's response
		utils.SetResponse(w, http.StatusOK, string(responseBody))
	} else {
		slog.Error("%s - error - %s", "Record deletion failed", err.Error())
		utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to delete record")
	}

}

// Get User by parameters
func GetUserByParam(w http.ResponseWriter, r *http.Request) {

}

// Update User Status
func UpdateUserStatus(w http.ResponseWriter, r *http.Request) {
	var data m.User
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err == nil {

		jsonValue, err := json.Marshal(data)
		if err != nil {
			slog.Error("%s - error - %s", "Error marshaling user data", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to marshal user data")
			return
		}
		resp, err := http.Post(utils.RestURL+"/update-user-status", "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			slog.Error("%s - error - %s", "Error making POST request", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to update user")
			return
		}
		defer resp.Body.Close()

		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			slog.Error("%s - error - %s", "Error reading response body", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to read response body")
			return
		}

		utils.SetResponse(w, resp.StatusCode, string(responseBody))
	} else {

		slog.Error("%s - error - %s", "Record updation failed", err.Error())
		utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to update record")
	}
}
func UserSearchPagination(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Pagination m.PaginationReq `json:"pagination"`
		Conditions []string        `json:"Conditions"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var Response struct {
		Pagination m.PaginationResp
		Data       []*m.UserManagement
	}

	jsonValue, err := json.Marshal(req)
	if err != nil {
		slog.Error("Error marshaling table condition", "error", err)
	}

	// Send POST request to REST service
	resp, err := http.Post(RestURL+"/get-all-user-by-search-pagination", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		slog.Error("Error making POST request", "error", err)
	}

	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&Response); err != nil {
		slog.Error("error decoding response body", "error", err)
	}

	// Build searchable filter map (optional use)
	searchFilters := make(map[string]string)
	for _, condition := range req.Conditions {
		parts := strings.Split(condition, " ")
		if len(parts) == 3 {
			dataField := parts[0]
			value := strings.Trim(parts[2], "'%")
			searchFilters[dataField] = value
		}
	}

	var UserTable s.TableCard
	UserTable.CardHeading = "User Master"

	UserTable.CardHeadingActions.Component = []s.CardHeadActionComponent{
		{
			ComponentName: "modelbutton",
			ComponentType: s.ActionComponentElement{
				ModelButton: s.ModelButtonAttributes{
					Type: "button", Text: "Add New User", ModelID: "#AddUserModel",
				},
			},
			Width: "col-4",
		},
		{
			ComponentName: "button",
			ComponentType: s.ActionComponentElement{
				Button: s.ButtonAttributes{
					ID: "redirectuserRolePage", Name: "redirectuserRolePage", Type: "button",
					Text: "User Roles", Colour: "#871A83",
				},
			},
			Width: "col-4",
		},
	}

	tableTools := `<button type="button" class="btn m-0 p-0" id="edit-User-btn"> 
	<i class="fa fa-edit" style="color: #871a83;"></i> 
	</button>`

	UserTable.BodyTables = s.CardTableBody{
		Columns: []s.CardTableBodyHeadCol{
			{
				Lable:        "User Name",
				Name:         `User Name`,
				Type:         "input",
				DataField:    "username",
				IsSearchable: true,
				Width:        "col-2",
			},
			{
				Name:         "Email",
				Type:         "input",
				DataField:    "email",
				IsSearchable: true,
				Width:        "col-2",
			},
			{
				Name:         "Role Name",
				Type:         "input",
				DataField:    "role_name",
				IsSearchable: true,
				Width:        "col-2",
			},
			{
				Name:  "Created On",
				Type:  "input",
				Width: "col-2",
			},
			{
				Name:  "Tools",
				Width: "col-1",
			},
		},
		ColumnsWidth: []string{"col-2", "col-2", "col-2", "col-2", "col-1"},
		Data:         Response.Data,
		Tools:        tableTools,
		ID:           "UserManagement",
	}

	tableBodyHTML := UserTable.BodyTables.RenderBodyColumns()

	var Pagination s.Pagination
	Pagination.TotalRecords = Response.Pagination.TotalNo
	Pagination.Offset = Response.Pagination.Offset
	Pagination.CurrentPage = Response.Pagination.Page
	Pagination.PerPage, _ = strconv.Atoi(req.Pagination.Limit)
	Pagination.PerPageOptions = []int{10, 25, 50, 100, 200}

	response := map[string]any{
		"tableBodyHTML":  tableBodyHTML,
		"paginationHTML": Pagination.Build(),
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}
