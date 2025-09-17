package server

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"strconv"

	"irpl.com/kanban-commons/model"
	m "irpl.com/kanban-commons/model"
	"irpl.com/kanban-commons/utils"
	"irpl.com/kanban-dao/dao"
)

// Get All Prod Process
func GetAllProductionProcessData(w http.ResponseWriter, r *http.Request) {

	prodProcessEntries, err := dao.GetAllProductionProcessEntries()
	if err != nil {
		slog.Error("Recordes not found", "error", err.Error())
		http.Error(w, "Failed to encode fetch records", http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(prodProcessEntries); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// Get All Prod Process
func GetAllProductionProcessEntries(w http.ResponseWriter, r *http.Request) {
	var data []model.ProdProcessCardData
	queryParams := r.URL.Query()
	line_no := queryParams.Get("line")
	if line_no == "" {
		line_no = "1"
	}

	prodProcessEntries, err := dao.GetAllProductionProcessForLine(line_no)
	if err != nil {
		slog.Error("Recordes not found", "error", err.Error())
		http.Error(w, "Failed to encode fetch records", http.StatusInternalServerError)
		return
	}

	tempdata, _ := dao.GetProductionProcessCardData(line_no)
	for _, i := range prodProcessEntries {
		var temp model.ProdProcessCardData
		temp.Id = i.Id
		temp.Name = i.Name
		temp.Link = i.Link
		temp.Icon = i.Icon
		temp.Description = i.Description
		temp.Status = i.Status
		temp.ProdProcessLineId = i.ProdProcessLineId
		temp.Order = i.Order
		temp.IsGroup = i.IsGroup
		temp.GroupName = i.GroupName
		for _, v := range tempdata {
			if i.ProdProcessLineId == v.ProdProcessLineId && i.ProdProcessId == v.ProdProcessId {
				temp.CompoundName = v.CompoundName
				temp.CellNo = v.CellNo
				temp.KbRootId = v.KbRootId
				temp.MFGDateTime = v.MFGDateTime
				temp.CreatedOn = v.CreatedOn
			}
		}
		data = append(data, temp)
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// Get All Prod Process
func GetAllProductionProcess(w http.ResponseWriter, r *http.Request) {

	prodProcessEntries, err := dao.GetAllProductionProcess()
	if err != nil {
		slog.Error("Recordes not found", "error", err.Error())
		http.Error(w, "Failed to encode fetch records", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(prodProcessEntries); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
func AddProdProcess(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		utils.SetResponse(w, http.StatusMethodNotAllowed, "Invalid request method")
		return
	}

	var payload m.ProdProcess
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.SetResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	err := dao.CreateNewOrUpdateExistingProdProcess(&payload)
	if err != nil {
		slog.Error("Records not found", "error", err.Error())
		utils.SetResponse(w, http.StatusInternalServerError, "Failed to add process records")
		return
	}
	utils.SetResponse(w, http.StatusOK, "Successfully added new the production process")
}

func GetProdProcessByParam(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var payload m.ProdProcess
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	ProdLineData, err := dao.GetProdProcessByParam("id", strconv.Itoa(payload.Id))
	if err != nil {
		slog.Error("Failed to get production line data.")
		utils.SetResponse(w, http.StatusInternalServerError, "Fail: Failed to get production line data.")
		return
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(ProdLineData)
	if err != nil {
		http.Error(w, "Error while encoding data", http.StatusInternalServerError)
		return
	}
}

func EditProdProcess(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var payload m.ProdProcess
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		// create log
		sysLog := m.SystemLog{
			Message:     "UpdateProductionProcess : Fail to update Production process",
			MessageType: "ERROR",
			IsCritical:  false,
			CreatedBy:   payload.CreatedBy,
		}
		utils.CreateSystemLogInternal(sysLog)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	ProdProcessData, err := dao.GetProdProcessByParam("id", strconv.Itoa(payload.Id))
	if err != nil {
		// create log
		sysLog := m.SystemLog{
			Message:     "UpdateProductionProcess : Fail to update Production process " + payload.Name,
			MessageType: "ERROR",
			IsCritical:  false,
			CreatedBy:   payload.CreatedBy,
		}
		utils.CreateSystemLogInternal(sysLog)
		slog.Error("Failed to retrieve production process data: ", err)
		utils.SetResponse(w, http.StatusInternalServerError, "Fail: Failed to retrieve production process data.")
		return
	}

	if len(ProdProcessData) == 0 {
		// create log
		sysLog := m.SystemLog{
			Message:     "UpdateProductionProcess : Fail to update Production process " + payload.Name,
			MessageType: "ERROR",
			IsCritical:  false,
			CreatedBy:   payload.CreatedBy,
		}
		utils.CreateSystemLogInternal(sysLog)
		slog.Error("Production process not found for ID: ", payload.Id)
		utils.SetResponse(w, http.StatusNotFound, "Fail: Production process not found.")
		return
	}

	payload.CreatedBy = ProdProcessData[0].CreatedBy
	payload.CreatedOn = ProdProcessData[0].CreatedOn // This should match CreatedOn, not ModifiedOn

	err = dao.CreateNewOrUpdateExistingProdProcess(&payload)
	if err != nil {
		// create log
		sysLog := m.SystemLog{
			Message:     "UpdateProductionProcess : Fail to update Production process " + payload.Name,
			MessageType: "ERROR",
			IsCritical:  false,
			CreatedBy:   payload.CreatedBy,
		}
		utils.CreateSystemLogInternal(sysLog)
		slog.Error("Failed to update production process: ", err)
		utils.SetResponse(w, http.StatusInternalServerError, "Fail: Failed to update production process.")
		return
	}

	// create log
	sysLog := m.SystemLog{
		Message:     "UpdateProductionProcess : Successfully updated Production process " + payload.Name,
		MessageType: "SUCCESS",
		IsCritical:  false,
		CreatedBy:   payload.CreatedBy,
	}
	utils.CreateSystemLogInternal(sysLog)
	utils.SetResponse(w, http.StatusOK, "Successfully updated the production process")
}

// GetAllProcessBySearchAndPagination returns a all records present in prod_process table
func GetAllProcessBySearchAndPagination(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		Pagination m.PaginationReq `json:"pagination"`
		Conditions []string        `json:"conditions"`
	}
	// Decode JSON request
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	prodProcess, PaginationResp, err := dao.GetAllProcessBySearchAndPagination(requestData.Pagination, requestData.Conditions)
	if err != nil {
		slog.Error("Recordes not found", "error", err.Error())
		http.Error(w, "Failed to encode fetch records", http.StatusInternalServerError)
		return
	}

	var Response struct {
		Pagination m.PaginationResp
		Data       []*m.ProdProcess
	}
	Response.Pagination = PaginationResp
	Response.Data = prodProcess

	if err := json.NewEncoder(w).Encode(Response); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
