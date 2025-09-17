package server

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	m "irpl.com/kanban-commons/model"
	"irpl.com/kanban-commons/utils"
	"irpl.com/kanban-dao/dao"
)

// Create New or Update UserToRole
func CreateNewUserToRoel(w http.ResponseWriter, r *http.Request) {
	var data m.UserToRole
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

// UpdateUserToRoel  UserToRole
func UpdateUserToRoel(w http.ResponseWriter, r *http.Request) {
	var data m.UserToRole
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err == nil {
		err := dao.CreateNewOrUpdateExistingUserToRole(&data)
		if err != nil {
			slog.Error("Record creation failed", "error", err.Error())
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to update user role")
			return
		}
		utils.SetResponse(w, http.StatusOK, "Successfully update user record")
	} else {
		slog.Error("Record creation failed", "error", err.Error())
		utils.SetResponse(w, http.StatusInternalServerError, "Failed to update user role ")
	}
}
