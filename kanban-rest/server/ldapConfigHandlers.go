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

// CreateLDAPConfig creates a new LDAP configuration
func CreateLDAPConfig(w http.ResponseWriter, r *http.Request) {
	var data model.LDAPConfig
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err == nil {
		jsonValue, err := json.Marshal(data)
		if err != nil {
			slog.Error("Error marshaling LDAPConfig data", "error", err.Error())
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to marshal LDAPConfig data")
			return
		}

		resp, err := http.Post(DBURL+"/create-ldap-config", "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			slog.Error("Error making POST request", "error", err.Error())
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to create LDAPConfig record")
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

// UpdateLDAPConfig updates an existing LDAP configuration
func UpdateLDAPConfig(w http.ResponseWriter, r *http.Request) {
	var data model.LDAPConfig
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err == nil {
		jsonValue, err := json.Marshal(data)
		if err != nil {
			slog.Error("Error marshaling LDAPConfig data", "error", err.Error())
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to marshal LDAPConfig data")
			return
		}

		resp, err := http.Post(DBURL+"/update-ldap-config", "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			slog.Error("Error making POST request", "error", err.Error())
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to update LDAPConfig record")
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

// GetDefaultLDAPConfig retrieves the default LDAP configuration
func GetDefaultLDAPConfig(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get(DBURL + "/get-default-ldap-config")
	if err != nil {
		slog.Error("Error making GET request", "error", err.Error())
		utils.SetResponse(w, http.StatusInternalServerError, "Failed to retrieve LDAPConfig")
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

// DeleteLDAPConfig deletes an LDAP configuration
func DeleteLDAPConfig(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	resp, err := http.Get(DBURL + "/delete-ldap-config?id=" + id)
	if err != nil {
		slog.Error("Error making GET request", "error", err.Error(), "id", id)
		utils.SetResponse(w, http.StatusInternalServerError, "Failed to delete LDAPConfig")
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
