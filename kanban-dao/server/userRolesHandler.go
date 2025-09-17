package server

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"strconv"

	"irpl.com/kanban-commons/utils"
	"irpl.com/kanban-dao/dao"

	m "irpl.com/kanban-commons/model"
)

// Create UserRole Entry
func CreateUserRole(w http.ResponseWriter, r *http.Request) {
	var data m.UserRoles
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err == nil {
		err := dao.CreateUserRoles(data)
		if err != nil {
			slog.Error("Record creation failed - " + err.Error())
			utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to create record")
			return
		}
		utils.SetResponse(w, http.StatusOK, "Success: successfully created record")
	} else {
		slog.Error("Record creation failed - " + err.Error())
		utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to create record")
	}
}

// Update UserRole Entry
func UpdateUserRole(w http.ResponseWriter, r *http.Request) {
	var data m.UserRoles
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err == nil {
		err := dao.UpdateUserRoles(data)
		if err != nil {
			utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to update record")
			slog.Error("Record updation failed - " + err.Error())
			return
		}
		utils.SetResponse(w, http.StatusOK, "Success: successfully updated record")
	} else {
		utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to update record")
		slog.Error("Record updation failed - " + err.Error())
	}
}

// Delete UserRole Entry
func DeleteUserRole(w http.ResponseWriter, r *http.Request) {
	var ids []int64
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&ids); err == nil {
		for _, id := range ids {
			err := dao.DeleteUserRoles(id)
			if err != nil {
				utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to delete record")
				slog.Error("Record deletion failed", "id", id, "error", err.Error())
				return
			}
		}
		utils.SetResponse(w, http.StatusOK, "Success: successfully deleted record")
	} else {
		utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to delete record")
		slog.Error("Record deletion failed - " + err.Error())
	}
}

// Get UserRole Entry by parameters
func GetUserRoleByParam(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")
	m := map[string]interface{}{}
	m[key] = value

	userRoleEntries, _, err := dao.GetUserRolesByCriteria(1, 1, m)
	if err != nil {
		slog.Error("Record not found", "key", key, "value", value, "error", err.Error())
		utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to find record")
		return
	}

	bData, _ := json.Marshal(userRoleEntries)
	utils.SetResponse(w, http.StatusOK, string(bData))
}

// GetUserRoles returns a paginated list of UserRole records
func GetUserRoles(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, err := utils.ParseInt(pageStr, 1)
	if err != nil {
		page = 1
	}

	limit, err := utils.ParseInt(limitStr, 10)
	if err != nil {
		limit = 10
	}

	userRoles, totalRecords, err := dao.GetUserRoles(page, limit)
	if err != nil {
		slog.Error("Failed to retrieve UserRoles", "error", err.Error())
		utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to retrieve records")
		return
	}

	response := map[string]interface{}{
		"total_records": totalRecords,
		"entries":       userRoles,
	}

	bData, _ := json.Marshal(response)
	utils.SetResponse(w, http.StatusOK, string(bData))
}

// GetUserRoleByNameHandler retrieves a UserRole record by name (case-insensitive)
func GetUserRoleByNameHandler(w http.ResponseWriter, r *http.Request) {
	roleName := r.URL.Query().Get("name")
	if roleName == "" {
		utils.SetResponse(w, http.StatusBadRequest, "Fail: 'name' query parameter is required")
		return
	}

	userRole, err := dao.GetUserRolesByName(roleName)
	if err != nil {
		slog.Error("Failed to retrieve UserRole by name", "name", roleName, "error", err.Error())
		utils.SetResponse(w, http.StatusNotFound, "Fail: Role not found")
		return
	}

	bData, _ := json.Marshal(userRole)
	utils.SetResponse(w, http.StatusOK, string(bData))
}

// GetUserRoleByID retrieves a UserRole record by ID
func GetUserRoleByIDHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the 'id' parameter from the query
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		utils.SetResponse(w, http.StatusBadRequest, "Fail: 'id' query parameter is required")
		return
	}

	// Convert the 'id' string to int64
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		slog.Error("Invalid 'id' parameter", "id", idStr, "error", err.Error())
		utils.SetResponse(w, http.StatusBadRequest, "Fail: 'id' must be a valid integer")
		return
	}

	// Fetch the UserRole by ID using DAO function
	userRole, err := dao.GetUserRolesByID(id)
	if err != nil {
		slog.Error("Failed to retrieve UserRole by ID", "id", id, "error", err.Error())
		utils.SetResponse(w, http.StatusNotFound, "Fail: Role not found")
		return
	}

	// Serialize the result and send the response
	bData, _ := json.Marshal(userRole)
	utils.SetResponse(w, http.StatusOK, string(bData))
}

// GetAllUserRoles retrives all user roles
func GetAllUserRoles(w http.ResponseWriter, r *http.Request) {

	userroles, err := dao.GetAllUserRoles()
	if err != nil {
		slog.Error("Recordes not found", "error", err.Error())
		http.Error(w, "Failed to encode fetch records", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(userroles); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
func GetAllUserRoleBySearchAndPagination(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		Pagination m.PaginationReq `json:"pagination"`
		Conditions []string        `json:"conditions"`
	}
	// Decode JSON request
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	KbRoot, PaginationResp, err := dao.GetAllUserRoleBySearchAndPagination(requestData.Pagination, requestData.Conditions)
	if err != nil {
		slog.Error("Recordes not found", "error", err.Error())
		http.Error(w, "Failed to encode fetch records", http.StatusInternalServerError)
		return
	}

	var Response struct {
		Pagination m.PaginationResp
		Data       []*m.UserRoles
	}
	Response.Pagination = PaginationResp
	Response.Data = KbRoot

	if err := json.NewEncoder(w).Encode(Response); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
