package server

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	m "irpl.com/kanban-commons/model"
	"irpl.com/kanban-commons/utils"
)

// Create UserRole
func CreateUserRole(w http.ResponseWriter, r *http.Request) {
	var data m.UserRoles

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err == nil {

		jsonValue, err := json.Marshal(data)
		if err != nil {
			slog.Error("Error marshaling user role data", "error", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to marshal user role data")
			return
		}

		resp, err := http.Post(DBURL+"/create-user-role", "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			slog.Error("Error making POST request", "error", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to create user role")
			return
		}
		defer resp.Body.Close()

		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			slog.Error("Error reading response body", "error", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to read response body")
			return
		}

		utils.SetResponse(w, resp.StatusCode, string(responseBody))
	} else {
		slog.Error("Record creation failed", "error", err.Error())
		utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to create record")
	}
}

// Update UserRole
func UpdateUserRole(w http.ResponseWriter, r *http.Request) {
	var data m.UserRoles

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err == nil {
		jsonValue, err := json.Marshal(data)
		if err != nil {
			slog.Error("Error marshaling user role data", "error", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to marshal user role data")
			return
		}

		resp, err := http.Post(DBURL+"/update-user-role", "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			slog.Error("Error making POST request", "error", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to update user role")
			return
		}
		defer resp.Body.Close()

		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			slog.Error("Error reading response body", "error", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to read response body")
			return
		}

		utils.SetResponse(w, http.StatusOK, string(responseBody))
	} else {
		slog.Error("Record updation failed", "error", err.Error())
		utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to update record")
	}
}

// Delete UserRole
func DeleteUserRole(w http.ResponseWriter, r *http.Request) {
	var ids []int64
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&ids); err == nil {
		jsonValue, err := json.Marshal(ids)
		if err != nil {
			slog.Error("Error marshaling delete request data", "error", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to marshal delete request data")
			return
		}

		resp, err := http.Post(DBURL+"/delete-user-role", "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			slog.Error("Error making POST request for delete", "error", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to delete user role(s)")
			return
		}
		defer resp.Body.Close()

		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			slog.Error("Error reading response body for delete", "error", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to read response body for delete")
			return
		}

		utils.SetResponse(w, http.StatusOK, string(responseBody))
	} else {
		slog.Error("Record deletion failed", "error", err.Error())
		utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to delete record")
	}
}

// Get UserRole by Name
func GetUserRoleByName(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")

	resp, err := http.Get(DBURL + "/get-role-by-name?name=" + name)
	if err != nil {
		slog.Error("Error making POST request for query", "error", err)
		utils.SetResponse(w, http.StatusInternalServerError, "Failed to query user role")
		return
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Error reading response body for query", "error", err)
		utils.SetResponse(w, http.StatusInternalServerError, "Failed to read response body for query")
		return
	}

	utils.SetResponse(w, http.StatusOK, string(responseBody))
}

// Get UserRole by ID
func GetUserRoleById(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	resp, err := http.Get(DBURL + "/get-role-by-id?id=" + id)
	if err != nil {
		slog.Error("Error making POST request for query", "error", err)
		utils.SetResponse(w, http.StatusInternalServerError, "Failed to query user role")
		return
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Error reading response body for query", "error", err)
		utils.SetResponse(w, http.StatusInternalServerError, "Failed to read response body for query")
		return
	}

	utils.SetResponse(w, http.StatusOK, string(responseBody))
}

// GetAllUserRoles returns a all records present in user roles table
func GetAllUserRoles(w http.ResponseWriter, r *http.Request) {

	resp, err := http.Get(DBURL + "/get-all-user-roles")
	if err != nil {
		slog.Error("%s - error - %s", "Error making GET request", err)
		utils.SetResponse(w, http.StatusInternalServerError, "Failed to find record")
		return
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("%s - error - %s", "Error reading response body", err)
		utils.SetResponse(w, http.StatusInternalServerError, "Failed to read response body")
		return
	}

	utils.SetResponse(w, http.StatusOK, string(responseBody))
}
func GetAllUserRoleBySearchAndPagination(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	resp, err := http.Post(DBURL+"/get-all-user-role-by-search-pagination", "application/json", bytes.NewBuffer(body))
	if err != nil {
		slog.Error("%s - error - %s", "Error making POST request", err)
		utils.SetResponse(w, http.StatusInternalServerError, "Failed to fectch order details")
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
}
