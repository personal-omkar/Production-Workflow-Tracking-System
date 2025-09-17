package server

import (
	"encoding/json"
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
		err := dao.CreateNewOrUpdateExistingUserToRole(&data)
		if err != nil {
			slog.Error("Record creation failed", "error", err.Error())
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to add user role")
			return
		}
		utils.SetResponse(w, http.StatusOK, "Successfully Added user role")
	} else {
		slog.Error("Record creation failed", "error", err.Error())
		utils.SetResponse(w, http.StatusInternalServerError, "Failed to add user role")
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
