package server

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"log/slog"
	"net/http"

	"irpl.com/kanban-commons/model"
	"irpl.com/kanban-commons/utils"
)

// Create SystemLog
func CreateSystemLog(w http.ResponseWriter, r *http.Request) {
	var data model.SystemLog
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err == nil {
		jsonValue, err := json.Marshal(data)
		if err != nil {
			slog.Error("Error marshaling SystemLog data", "error", err.Error())
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to marshal SystemLog data")
			return
		}

		resp, err := http.Post(DBURL+"/create-system-log", "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			log.Print("In the rest : error to requesrt DB ")
			slog.Error("Error making POST request", "error", err.Error())
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to create SystemLog record")
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

// Update SystemLog
func UpdateSystemLog(w http.ResponseWriter, r *http.Request) {
	var data model.SystemLog
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err == nil {
		jsonValue, err := json.Marshal(data)
		if err != nil {
			slog.Error("Error marshaling SystemLog data", "error", err.Error())
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to marshal SystemLog data")
			return
		}

		resp, err := http.Post(DBURL+"/update-system-log", "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			slog.Error("Error making POST request", "error", err.Error())
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to update SystemLog record")
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

// Get All SystemLogs
func GetAllSystemLogs(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get(DBURL + "/get-all-entries-systemlogs")
	if err != nil {
		slog.Error("Error making GET request", "error", err.Error())
		utils.SetResponse(w, http.StatusInternalServerError, "Failed to retrieve logs")
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

// Get SystemLog by ID
func GetSystemLogByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	resp, err := http.Get(DBURL + "/get-system-log-by-id?id=" + id)
	if err != nil {
		slog.Error("Error making GET request", "error", err.Error(), "id", id)
		utils.SetResponse(w, http.StatusInternalServerError, "Failed to retrieve log")
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

// Delete SystemLog by ID
func DeleteSystemLog(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	resp, err := http.Get(DBURL + "/delete-system-log?id=" + id)
	if err != nil {
		slog.Error("Error making GET request", "error", err.Error(), "id", id)
		utils.SetResponse(w, http.StatusInternalServerError, "Failed to delete log")
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
