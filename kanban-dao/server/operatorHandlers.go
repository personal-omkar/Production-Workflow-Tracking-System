package server

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"

	m "irpl.com/kanban-commons/model"
	"irpl.com/kanban-commons/utils"
	"irpl.com/kanban-dao/dao"
)

// CreateNewOrUpdateExistingOperator create a new operator entry or update existing operator
func CreateNewOrUpdateExistingOperator(w http.ResponseWriter, r *http.Request) {
	var data m.Operator
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err == nil {
		var msg string
		if data.Id != 0 {
			msg = "Success: successfully updated record"
		} else {
			msg = "Success: successfully created record"
		}
		err := dao.CreateNewOrUpdateExistingOperator(&data)
		if err != nil {
			slog.Error("Record creation failed - " + err.Error())
			utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to create record")
			return
		}

		utils.SetResponse(w, http.StatusOK, msg)
	} else {
		slog.Error("Record creation failed - " + err.Error())
		utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to create record")
	}
}

// GetAllOperator returns a all records present in operator table
func GetAllOperator(w http.ResponseWriter, r *http.Request) {
	KbRoot, err := dao.GetAllOperator()
	if err != nil {
		slog.Error("Recordes not found", "error", err.Error())
		http.Error(w, "Failed to encode fetch records", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(KbRoot); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// GetOperatorByParam returns a operator records based on parameter
func GetOperatorByParam(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")
	KbRoot, err := dao.GetOperatorByParam(key, value)
	if err != nil {
		slog.Error("Recordes not found", "error", err.Error())
		http.Error(w, "Failed to encode fetch records", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(KbRoot); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// GetAllOperator returns a all records present in operator table
func GetAllOperatorBySearchAndPagination(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		Pagination m.PaginationReq `json:"pagination"`
		Conditions []string        `json:"conditions"`
	}
	// Decode JSON request
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	KbRoot, PaginationResp, err := dao.GetAllOperatorBySearchAndPagination(requestData.Pagination, requestData.Conditions)
	if err != nil {
		slog.Error("Recordes not found", "error", err.Error())
		http.Error(w, "Failed to encode fetch records", http.StatusInternalServerError)
		return
	}

	var Response struct {
		Pagination m.PaginationResp
		Data       []*m.Operator
	}
	Response.Pagination = PaginationResp
	Response.Data = KbRoot

	if err := json.NewEncoder(w).Encode(Response); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
