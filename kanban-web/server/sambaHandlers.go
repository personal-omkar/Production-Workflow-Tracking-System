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

// Create SAMBA Configuration
func CreateSAMBAConfiguration(w http.ResponseWriter, r *http.Request) {

	var data model.SambaConfig

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err == nil {
		userIDStr := r.Header.Get("X-Custom-Userid")
		data.CreatedBy = userIDStr
		data.ModifiedBy = userIDStr

		jsonValue, err := json.Marshal(data)
		if err != nil {
			slog.Error("%s - error - %s", "Error marshaling user data", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to marshal user data")
			return
		}

		resp, err := http.Post(utils.RestURL+"/create-samba-config", "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			slog.Error("%s - error - %s", "Error making POST request", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to create samba config record")
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

		slog.Error("%s - error - %s", "Record creation failed", err.Error())
		utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to create record")
	}
}

// Update SAMBA Configuration
func UpdateSAMBAConfiguration(w http.ResponseWriter, r *http.Request) {

	var data model.SambaConfig

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err == nil {
		userIDStr := r.Header.Get("X-Custom-Userid")
		data.CreatedBy = userIDStr
		data.ModifiedBy = userIDStr

		jsonValue, err := json.Marshal(data)
		if err != nil {
			slog.Error("%s - error - %s", "Error marshaling user data", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to marshal user data")
			return
		}

		resp, err := http.Post(utils.RestURL+"/update-samba-config", "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			slog.Error("%s - error - %s", "Error making POST request", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to update samba config record")
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

		slog.Error("%s - error - %s", "Record update failed", err.Error())
		utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to update record")
	}
}

// Get Default SAMBA Configuration
func GetDefaultSAMBAConfig(w http.ResponseWriter, r *http.Request) {

	resp, err := http.Get(utils.RestURL + "/get-default-samba-config")
	if err != nil {
		slog.Error("%s - error - %s", "Error making POST request", err)
		utils.SetResponse(w, http.StatusInternalServerError, "Failed to get default samba config record")
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
