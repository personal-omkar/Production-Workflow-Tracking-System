package server

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"time"

	"irpl.com/kanban-commons/model"
	"irpl.com/kanban-commons/utils"
	"irpl.com/kanban-dao/dao"
)

// Get all API keys
func GetAllAPIKeys(w http.ResponseWriter, r *http.Request) {
	entries, err := dao.GetAllAPIKeys()
	if err != nil {
		slog.Error("Records not found", "error", err.Error())
		http.Error(w, "Failed to fetch API keys", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(entries); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// Get API keys by parameter (key=value)
func GetAPIKeyByParam(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")

	if key == "" || value == "" {
		http.Error(w, "Missing key or value in query params", http.StatusBadRequest)
		return
	}

	entries, err := dao.GetAPIKeyByParam(key, value)
	if err != nil {
		slog.Error("Failed to fetch API keys by param", "error", err.Error())
		http.Error(w, "Failed to fetch records", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(entries); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// Create or update API key
func AddOrUpdateAPIKey(w http.ResponseWriter, r *http.Request) {
	var data model.APIKey
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		utils.SetResponse(w, http.StatusBadRequest, "Fail: Decode data")
		slog.Error("Failed to decode API key data", "error", err.Error())
		return
	}

	if data.Id == 0 {
		data.CreatedOn = time.Now()
		utils.CreateSystemLogInternal(model.SystemLog{
			Message:     "APIKey: Added new key " + data.Name,
			MessageType: "SUCCESS",
			IsCritical:  false,
			CreatedBy:   data.CreatedBy,
		})
	} else {
		utils.CreateSystemLogInternal(model.SystemLog{
			Message:     "APIKey: Updated key " + data.Name,
			MessageType: "SUCCESS",
			IsCritical:  false,
			CreatedBy:   data.ModifiedBy,
		})
		data.ModifiedOn = time.Now()
	}

	if err := dao.CreateNewOrUpdateAPIKey(&data); err != nil {
		utils.SetResponse(w, http.StatusInternalServerError, "Fail: Update data")
		slog.Error("Failed to update API key", "error", err.Error())
		return
	}

	utils.SetResponse(w, http.StatusCreated, "API Key saved successfully")
}

// Get API keys with pagination and filters
func GetAPIKeysWithSearchAndPagination(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Pagination model.PaginationReq `json:"pagination"`
		Conditions []string            `json:"conditions"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	data, pagination, err := dao.GetAPIKeysWithPagination(req.Pagination, req.Conditions)
	if err != nil {
		slog.Error("Failed to fetch API keys", "error", err.Error())
		http.Error(w, "Failed to fetch records", http.StatusInternalServerError)
		return
	}

	res := struct {
		Pagination model.PaginationResp `json:"pagination"`
		Data       []*model.APIKey      `json:"data"`
	}{
		Pagination: pagination,
		Data:       data,
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
