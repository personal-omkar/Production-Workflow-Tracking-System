package server

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"irpl.com/kanban-commons/model"
	"irpl.com/kanban-commons/utils"
	"irpl.com/kanban-dao/dao"
)

// CreateSambaConfig creates a new Samba configuration entry
func CreateSambaConfig(w http.ResponseWriter, r *http.Request) {
	var data model.SambaConfig
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err == nil {
		err := dao.CreateSambaConfig(data)
		if err != nil {
			slog.Error("SambaConfig creation failed - " + err.Error())
			utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to create record")
			return
		}
		utils.SetResponse(w, http.StatusOK, "Success: successfully created record")
	} else {
		slog.Error("SambaConfig creation failed - " + err.Error())
		utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to create record")
	}
}

// UpdateSambaConfig updates an existing Samba configuration entry
func UpdateSambaConfig(w http.ResponseWriter, r *http.Request) {
	var data model.SambaConfig
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err == nil {
		rec, _ := dao.GetSambaConfigByID(data.ID)
		if rec.ID != 0 {
			err := dao.UpdateSambaConfig(data)
			if err != nil {
				utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to update record")
				slog.Error("SambaConfig update failed - " + err.Error())
				return
			}
			utils.SetResponse(w, http.StatusOK, "Success: successfully updated record")
		} else {
			data.IsDefault = true
			err := dao.CreateSambaConfig(data)
			if err != nil {
				utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to update record")
				slog.Error("SambaConfig create failed - " + err.Error())
				return
			}
			utils.SetResponse(w, http.StatusOK, "Success: successfully created record")
		}
	} else {
		utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to update record")
		slog.Error("SambaConfig update failed - " + err.Error())
	}
}

// DeleteSambaConfig deletes a Samba configuration entry
func DeleteSambaConfig(w http.ResponseWriter, r *http.Request) {
	var ids []int
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&ids); err == nil {
		for _, id := range ids {
			err := dao.DeleteSambaConfig(id)
			if err != nil {
				utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to delete record")
				slog.Error("SambaConfig deletion failed - " + err.Error())
				return
			}
		}
		utils.SetResponse(w, http.StatusOK, "Success: successfully deleted record")
	} else {
		utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to delete record")
		slog.Error("SambaConfig deletion failed - " + err.Error())
	}
}

// GetDefaultSambaConfig retrieves the default Samba configuration
func GetDefaultSambaConfig(w http.ResponseWriter, r *http.Request) {
	entry, err := dao.GetDefaultSambaConfig()
	if err != nil {
		slog.Error("Failed to fetch default SambaConfig", "error", err.Error())
		utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to fetch default record")
		return
	}

	bData, _ := json.Marshal(entry)
	utils.SetResponse(w, http.StatusOK, string(bData))
}
