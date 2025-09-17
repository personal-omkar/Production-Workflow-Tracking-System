package server

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"irpl.com/kanban-commons/model"
	"irpl.com/kanban-commons/utils"
	"irpl.com/kanban-dao/dao"
)

// CreateLDAPConfig creates a new LDAP configuration entry
func CreateLDAPConfig(w http.ResponseWriter, r *http.Request) {
	var data model.LDAPConfig
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err == nil {
		err := dao.CreateLDAPConfig(data)
		if err != nil {
			slog.Error("LDAPConfig creation failed - " + err.Error())
			utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to create record")
			return
		}
		utils.SetResponse(w, http.StatusOK, "Success: successfully created record")
	} else {
		slog.Error("LDAPConfig creation failed - " + err.Error())
		utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to create record")
	}
}

// UpdateLDAPConfig updates an existing LDAP configuration entry
func UpdateLDAPConfig(w http.ResponseWriter, r *http.Request) {
	var data model.LDAPConfig
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err == nil {
		rec, _ := dao.GetLDAPConfigByID(data.ID)
		if rec.ID != 0 {
			err := dao.UpdateLDAPConfig(data)
			if err != nil {
				utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to update record")
				slog.Error("LDAPConfig update failed - " + err.Error())
				return
			}
			utils.SetResponse(w, http.StatusOK, "Success: successfully updated record")
		} else {
			data.IsDefault = true
			err := dao.CreateLDAPConfig(data)
			if err != nil {
				utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to update record")
				slog.Error("LDAPConfig create failed - " + err.Error())
				return
			}
			utils.SetResponse(w, http.StatusOK, "Success: successfully created record")
		}

	} else {
		utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to update record")
		slog.Error("LDAPConfig update failed - " + err.Error())
	}
}

// DeleteLDAPConfig deletes an LDAP configuration entry
func DeleteLDAPConfig(w http.ResponseWriter, r *http.Request) {
	var ids []int
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&ids); err == nil {
		for _, id := range ids {
			err := dao.DeleteLDAPConfig(id)
			if err != nil {
				utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to delete record")
				slog.Error("LDAPConfig deletion failed - " + err.Error())
				return
			}
		}
		utils.SetResponse(w, http.StatusOK, "Success: successfully deleted record")
	} else {
		utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to delete record")
		slog.Error("LDAPConfig deletion failed - " + err.Error())
	}
}

// GetDefaultLDAPConfig retrieves the default LDAP configuration
func GetDefaultLDAPConfig(w http.ResponseWriter, r *http.Request) {
	entry, err := dao.GetDefaultLDAPConfig()
	if err != nil {
		slog.Error("Failed to fetch default LDAPConfig", "error", err.Error())
		utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to fetch default record")
		return
	}

	bData, _ := json.Marshal(entry)
	utils.SetResponse(w, http.StatusOK, string(bData))
}
