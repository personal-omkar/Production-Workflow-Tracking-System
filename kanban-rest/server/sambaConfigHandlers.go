package server

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"irpl.com/kanban-commons/model"
	"irpl.com/kanban-commons/utils"
)

// CreateSambaConfig creates a new Samba configuration
func CreateSambaConfig(w http.ResponseWriter, r *http.Request) {
	var data model.SambaConfig
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err == nil {
		jsonValue, err := json.Marshal(data)
		if err != nil {
			slog.Error("Error marshaling SambaConfig data", "error", err.Error())
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to marshal SambaConfig data")
			return
		}

		resp, err := http.Post(DBURL+"/create-samba-config", "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			slog.Error("Error making POST request", "error", err.Error())
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to create SambaConfig record")
			return
		}
		defer resp.Body.Close()

		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			slog.Error("Error reading response body", "error", err.Error())
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to read response body")
			return
		}

		utils.SetResponse(w, resp.StatusCode, string(responseBody))
	} else {
		slog.Error("Record creation failed", "error", err.Error())
		utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to create record")
	}
}

// UpdateSambaConfig updates an existing Samba configuration
func UpdateSambaConfig(w http.ResponseWriter, r *http.Request) {
	var data model.SambaConfig
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err == nil {
		jsonValue, err := json.Marshal(data)
		if err != nil {
			slog.Error("Error marshaling SambaConfig data", "error", err.Error())
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to marshal SambaConfig data")
			return
		}

		resp, err := http.Post(DBURL+"/update-samba-config", "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			slog.Error("Error making POST request", "error", err.Error())
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to update SambaConfig record")
			return
		}
		defer resp.Body.Close()

		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			slog.Error("Error reading response body", "error", err.Error())
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to read response body")
			return
		}

		utils.SetResponse(w, resp.StatusCode, string(responseBody))
	} else {
		slog.Error("Record update failed", "error", err.Error())
		utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to update record")
	}
}

// GetDefaultSambaConfig retrieves the default Samba configuration
func GetDefaultSambaConfig(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get(DBURL + "/get-default-samba-config")
	if err != nil {
		slog.Error("Error making GET request", "error", err.Error())
		utils.SetResponse(w, http.StatusInternalServerError, "Failed to retrieve SambaConfig")
		return
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Error reading response body", "error", err.Error())
		utils.SetResponse(w, http.StatusInternalServerError, "Failed to read response body")
		return
	}

	utils.SetResponse(w, http.StatusOK, string(responseBody))
}

// DeleteSambaConfig deletes a Samba configuration
func DeleteSambaConfig(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	resp, err := http.Get(DBURL + "/delete-samba-config?id=" + id)
	if err != nil {
		slog.Error("Error making GET request", "error", err.Error(), "id", id)
		utils.SetResponse(w, http.StatusInternalServerError, "Failed to delete SambaConfig")
		return
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Error reading response body", "error", err.Error(), "id", id)
		utils.SetResponse(w, http.StatusInternalServerError, "Failed to read response body")
		return
	}

	utils.SetResponse(w, http.StatusOK, string(responseBody))
}
