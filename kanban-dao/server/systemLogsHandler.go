package server

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"irpl.com/kanban-commons/model"
	"irpl.com/kanban-commons/utils"
	"irpl.com/kanban-dao/dao"
)

// Create SystemLog Entry
func CreateSystemLog(w http.ResponseWriter, r *http.Request) {
	var data model.SystemLog
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&data); err == nil {
		err := dao.CreateSystemLog(data)
		if err != nil {
			slog.Error("Log creation failed - " + err.Error())
			utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to create log")
			return
		}
		utils.SetResponse(w, http.StatusOK, "Success: successfully created log")
	} else {
		slog.Error("Log creation failed - " + err.Error())
		utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to create log")
	}
}

// Update SystemLog Entry
func UpdateSystemLog(w http.ResponseWriter, r *http.Request) {
	var data model.SystemLog
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err == nil {
		err := dao.UpdateSystemLog(data)
		if err != nil {
			utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to update log")
			slog.Error("Log update failed - " + err.Error())
			return
		}
		utils.SetResponse(w, http.StatusOK, "Success: successfully updated log")
	} else {
		utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to update log")
		slog.Error("Log update failed - " + err.Error())
	}
}

// Delete SystemLog Entry
func DeleteSystemLog(w http.ResponseWriter, r *http.Request) {
	var ids []int
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&ids); err == nil {
		for _, id := range ids {
			err := dao.DeleteSystemLog(id)
			if err != nil {
				utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to delete log")
				slog.Error("Log deletion failed - " + err.Error())
				return
			}
		}
		utils.SetResponse(w, http.StatusOK, "Success: successfully deleted log")
	} else {
		utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to delete log")
		slog.Error("Log deletion failed  - " + err.Error())
	}
}

// Get SystemLog by ID
func GetSystemLogByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := utils.ParseInt(idStr, 0)
	if err != nil || id <= 0 {
		utils.SetResponse(w, http.StatusBadRequest, "Fail: invalid ID")
		slog.Error("Invalid ID - " + err.Error())
		return
	}

	logEntry, err := dao.GetSystemLogByID(id)
	if err != nil {
		utils.SetResponse(w, http.StatusNotFound, "Fail: failed to find log")
		slog.Error("Log retrieval failed - " + err.Error())
		return
	}

	bData, _ := json.Marshal(logEntry)
	utils.SetResponse(w, http.StatusOK, string(bData))
}

// Get SystemLogs (Paginated)
func GetSystemLogs(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, err := utils.ParseInt(pageStr, 1)
	if err != nil {
		page = 1
	}

	limit, err := utils.ParseInt(limitStr, 10)
	if err != nil {
		limit = 10
	}

	logEntries, totalRecords, err := dao.GetSystemLogs(page, limit)
	if err != nil {
		slog.Error("Failed to retrieve logs", "error", err.Error())
		utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to retrieve logs")
		return
	}

	response := map[string]interface{}{
		"total_records": totalRecords,
		"entries":       logEntries,
	}

	bData, _ := json.Marshal(response)
	utils.SetResponse(w, http.StatusOK, string(bData))
}

// Get All SystemLogs
func GetAllSystemLogs(w http.ResponseWriter, r *http.Request) {
	logEntries, err := dao.GetAllSystemLogs()
	if err != nil {
		slog.Error("Logs not found", "error", err.Error())
		utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to fetch logs")
		return
	}

	bData, _ := json.Marshal(logEntries)
	utils.SetResponse(w, http.StatusOK, string(bData))
}

// Search SystemLogs by Criteria
// func SearchSystemLogs(w http.ResponseWriter, r *http.Request) {
// 	var data model.PaginationReq
// 	decoder := json.NewDecoder(r.Body)

// 	if err := decoder.Decode(&data); err == nil {
// 		entries, err := dao.GetPaginatedData(data)
// 		if err != nil {
// 			utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to fetch logs")
// 			slog.Error("Log search failed - " + err.Error())
// 			return
// 		}
// 		resultData, _ := json.Marshal(entries)
// 		utils.SetResponse(w, http.StatusOK, string(resultData))
// 	} else {
// 		slog.Error("Invalid request - " + err.Error())
// 		utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to process search request")
// 	}
// }
