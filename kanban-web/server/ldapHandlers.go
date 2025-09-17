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

// Create LDAP Configuration
func CreateLDAPConfiguration(w http.ResponseWriter, r *http.Request) {

	var data model.LDAPConfig

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

		resp, err := http.Post(utils.RestURL+"/create-ldap-config", "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			slog.Error("%s - error - %s", "Error making POST request", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to create cuttube entry record")
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

// Update LDAP Configuration
func UpdateLDAPConfiguration(w http.ResponseWriter, r *http.Request) {

	var data model.LDAPConfig

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

		resp, err := http.Post(utils.RestURL+"/update-ldap-config", "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			slog.Error("%s - error - %s", "Error making POST request", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to update ldap record")
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

// Get Default LDAP Configuration
func GetDefaultLDAPConfig(w http.ResponseWriter, r *http.Request) {

	resp, err := http.Get(utils.RestURL + "/get-default-ldap-config")
	if err != nil {
		slog.Error("%s - error - %s", "Error making POST request", err)
		utils.SetResponse(w, http.StatusInternalServerError, "Failed to get default ldap config record")
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
