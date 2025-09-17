package server

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"strconv"

	m "irpl.com/kanban-commons/model"
	u "irpl.com/kanban-commons/utils"
	"irpl.com/kanban-dao/dao"
)

// GetAllKbRoot returns a all records present in kb_data table
func GetAllKbRoot(w http.ResponseWriter, r *http.Request) {
	KbRoot, err := dao.GetAllKbRoot()
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

// GetKbRootByParam returns a kb_root records based on parameter
func GetKbRootByParam(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")
	KbRoot, err := dao.GetKbRootByParam(key, value)
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

// func to udpate running number:
func UpdateRunningNumbers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var updatesRunningNo []m.KbRoot
	err := json.NewDecoder(r.Body).Decode(&updatesRunningNo)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		log.Println("Error decoding JSON:", err)
		return
	}
	err = dao.UpdateRunningNo(updatesRunningNo)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusInternalServerError)
		log.Println("Error while updating data:", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func GetAllCompletedKBRootDetails(w http.ResponseWriter, r *http.Request) {
	var condition string
	type TableRequest struct {
		PaginationReq m.PaginationReq `json:"pagination"`
		Conditions    []string        `json:"conditions"`
	}
	var data TableRequest

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&data); err == nil {
		for _, i := range data.Conditions {
			condition += i
		}
	}
	ord, err := dao.GetAllCompletedKBRootDetails()
	if err != nil {
		slog.Error("Recordes not found", "error", err.Error())
		http.Error(w, "Failed to encode fetch records", http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(ord); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func GetCompletedKBRootDetailsBySearch(w http.ResponseWriter, r *http.Request) {
	var condition string
	type TableRequest struct {
		PaginationReq m.PaginationReq `json:"pagination"`
		Conditions    []string        `json:"conditions"`
	}
	var data TableRequest

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&data); err == nil {
		for _, i := range data.Conditions {
			condition += i
		}
	}
	ord, err := dao.GetCompletedKBRootDetailsBySearch(condition)
	if err != nil {
		slog.Error("Recordes not found", "error", err.Error())
		http.Error(w, "Failed to encode fetch records", http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(ord); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func GetDetailRootData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var krData m.KbRoot
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&krData); err != nil {
		log.Printf("Error to get data: %v", err)
		http.Error(w, "Failed to decode data", http.StatusInternalServerError)
		return
	}
	Data, err := dao.GetDetailRootData(krData.Id)
	if err != nil {
		slog.Error("Recordes not found", "error", err.Error())
		http.Error(w, "Failed to encode fetch records", http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(Data); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func IsValidStatusUpdate(ids m.IsValidStatusUpdate) bool {
	var currentstatus, updatingstatus int
	var kbe m.KbExtension
	StatusMap := map[string]int{
		"creating":            1,
		"pending":             2,
		"reject":              3,
		"approved":            3,
		"InProductionLine":    4,
		"InProductionProcess": 5,
		"dispatch":            6,
	}
	if ids.KbDataId != 0 {
		kbe, _ = dao.GetOrderStatusByKbData(ids.KbDataId)
	} else if ids.KbRootId != 0 {
		kbe, _ = dao.GetOrderStatusByKbroot(ids.KbRootId)
	}

	for key, val := range StatusMap {
		if key == kbe.Status {
			currentstatus = val
		}
		if key == ids.Status {
			updatingstatus = val
		}
	}
	if updatingstatus > currentstatus && updatingstatus == currentstatus+1 {
		return true
	} else if updatingstatus > currentstatus && currentstatus == 2 && updatingstatus == 6 {
		return true
	} else if updatingstatus < currentstatus && currentstatus == 4 && updatingstatus == currentstatus-1 {
		return true
	} else if currentstatus == 1 && updatingstatus == 1 { //Added check, cuz while updating order we dont change status
		return true
	} else {
		return false
	}

}

// DeleteKbRootByIDsHandler handles the deletion of kb_root records based on a list of IDs
func DeleteKbRootByIDsHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		IDs    []string `json:"ids"`
		UserID string   `json:"UserID"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		u.SetResponse(w, http.StatusBadRequest, "Fail to get kanban details to delete kanban")
		return
	}
	if len(req.IDs) == 0 {
		u.SetResponse(w, http.StatusBadRequest, "Fail to get kanban details to delete kanban")
		return
	}
	if err := dao.DeleteKbRootByIDs(req.IDs); err != nil {
		u.SetResponse(w, http.StatusInternalServerError, "Fail to delete kanban")
		return
	}

	sysLog := m.SystemLog{
		Message:     strconv.Itoa(len(req.IDs)) + " Kanban got deleted successfully from Kanban",
		MessageType: "DELETE",
		IsCritical:  false,
		CreatedBy:   req.UserID,
	}
	u.CreateSystemLogInternal(sysLog)

	u.SetResponse(w, http.StatusOK, "Records deleted successfully")
}

func GetAllKanbanDetailsForReportHandler(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		Pagination m.PaginationReq `json:"pagination"`
		Conditions []string        `json:"conditions"`
	}

	// Decode JSON request
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Fetch paginated data with filters
	ord, PaginationResp, err := dao.GetAllKanbanDetailsForReport(requestData.Pagination, requestData.Conditions)
	if err != nil {
		slog.Error("Records not found", "error", err.Error())
		http.Error(w, "Failed to fetch records", http.StatusInternalServerError)
		return
	}

	var Reponse struct {
		Pagination m.PaginationResp
		Data       []m.OrderDetails
	}
	Reponse.Pagination = PaginationResp
	Reponse.Data = ord
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(Reponse); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
