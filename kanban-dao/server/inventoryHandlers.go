package server

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	m "irpl.com/kanban-commons/model"
	"irpl.com/kanban-commons/utils"
	"irpl.com/kanban-dao/dao"
)

// Create new inventory or update existing inventory
func CreateOrUpdateInventory(w http.ResponseWriter, r *http.Request) {
	var data m.Inventory
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err == nil {
		err := dao.CreateNewOrUpdateExistingInventory(&data)
		if err != nil {
			slog.Error("Record creation failed", "error", err.Error())
			utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to create record")
			return
		}
		utils.SetResponse(w, http.StatusOK, "Success: successfully created record")
	} else {
		slog.Error("Record creation failed", "error", err.Error())
		utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to create record")
	}
}

// Get All inventory
func GetAllInventoryrecords(w http.ResponseWriter, r *http.Request) {

	prodLineEntries, err := dao.GetAllInventory()
	if err != nil {
		slog.Error("Recordes not found", "error", err.Error())
		http.Error(w, "Failed to encode fetch records", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(prodLineEntries); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// GetInventoryByParam returns a inventory records based on parameter
func GetInventoryByParam(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")

	inventory, err := dao.GetInventoryByParam(key, value)
	if err != nil {
		slog.Error("Recordes not found", "error", err.Error())
		http.Error(w, "Failed to encode fetch records", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(inventory); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// DeleteInventoryByParam delete a inventory records based on parameter
func DeleteInventoryByParam(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")

	err := dao.DeleteInventoryByParam(key, value)
	if err != nil {
		utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to delete record")
		slog.Error("Record deletion failed for ", key, value, "error", err.Error())
		return
	}

	utils.SetResponse(w, http.StatusOK, "Success: successfully deleted record")
}

// GetAllColdStoragerecords All inventory records with the compound name
func GetAllColdStoragerecords(w http.ResponseWriter, r *http.Request) {
	var coldstorage []m.ColdStorage
	var tempcoldstorage m.ColdStorage

	inventory, err := dao.GetAllInventory()
	if err != nil {
		slog.Error("Recordes not found", "error", err.Error())
		http.Error(w, "Failed to encode fetch records", http.StatusInternalServerError)
		return
	}

	for _, i := range inventory {
		comp, _ := dao.GetCompoundDataByParam("id", strconv.Itoa(i.CompoundID))
		tempcoldstorage.CreatedBy = i.CreatedBy
		tempcoldstorage.CreatedOn = i.CreatedOn
		tempcoldstorage.ModifiedBy = i.ModifiedBy
		tempcoldstorage.ModifiedOn = i.ModifiedOn
		tempcoldstorage.MinQuantity = i.MinQuantity
		tempcoldstorage.MaxQuantity = i.MaxQuantity
		tempcoldstorage.Description = i.Description
		tempcoldstorage.Id = i.Id
		tempcoldstorage.CompoundID = i.CompoundID
		tempcoldstorage.AvailableQuantity = i.AvailableQuantity
		tempcoldstorage.CompoundName = comp[0].CompoundName

		coldstorage = append(coldstorage, tempcoldstorage)
	}
	if err := json.NewEncoder(w).Encode(coldstorage); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// UpdateColdStorageQuantity min and max quantity ofs existing inventory
func UpdateColdStorageQuantity(w http.ResponseWriter, r *http.Request) {
	var data m.Inventory
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&data); err == nil {
		inv, _ := dao.GetInventoryByParam("id", strconv.Itoa(data.Id))
		CompoundData, _ := dao.GetCompoundDataByParam("id", inv[0].CompoundID)
		inv[0].MinQuantity = data.MinQuantity
		inv[0].MaxQuantity = data.MaxQuantity
		inv[0].ModifiedOn = time.Now()
		err := dao.CreateNewOrUpdateExistingInventory(inv[0])
		if err != nil {
			// create log
			sysLog := m.SystemLog{
				Message:     "UpdateInventory: Failed to update inventory part : " + CompoundData[0].CompoundName,
				MessageType: "ERROR",
				IsCritical:  false,
				CreatedBy:   data.CreatedBy,
			}
			utils.CreateSystemLogInternal(sysLog)
			slog.Error("Record creation failed", "error", err.Error())
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to create record")
			return
		}
		// create log
		sysLog := m.SystemLog{
			Message:     "UpdateInventory: " + CompoundData[0].CompoundName + " Part updated",
			MessageType: "SUCCESS",
			IsCritical:  false,
			CreatedBy:   data.CreatedBy,
		}
		utils.CreateSystemLogInternal(sysLog)
		utils.SetResponse(w, http.StatusOK, "Successfully created record")
	} else {
		// create log
		sysLog := m.SystemLog{
			Message:     "UpdateInventory: Failed to update inventory",
			MessageType: "ERROR",
			IsCritical:  false,
			CreatedBy:   data.CreatedBy,
		}
		utils.CreateSystemLogInternal(sysLog)
		slog.Error("Record creation failed", "error", err.Error())
		utils.SetResponse(w, http.StatusInternalServerError, "Failed to create record")
	}
}

func GetInventoryBySearch(w http.ResponseWriter, r *http.Request) {
	var con string
	type TableConditions struct {
		Conditions []string `json:"Conditions"`
	}
	var tablecondition TableConditions

	err := json.NewDecoder(r.Body).Decode(&tablecondition)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		log.Println("Error decoding JSON:", err)
		return
	}

	for _, i := range tablecondition.Conditions {
		con += i
	}
	orddetails, err := dao.GetInventoryBySearch(con)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusInternalServerError)
		log.Println("Error while updating data:", err)
		return
	}

	if err := json.NewEncoder(w).Encode(orddetails); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func ColdStorageSearchPagination(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		Pagination m.PaginationReq `json:"pagination"`
		Conditions []string        `json:"conditions"`
	}
	// Decode JSON request
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	comp, PaginationResp, err := dao.GetInventoryBySearchPagination(requestData.Pagination, requestData.Conditions)
	if err != nil {
		slog.Error("Recordes not found", "error", err.Error())
		http.Error(w, "Failed to encode fetch records", http.StatusInternalServerError)
		return
	}

	var Response struct {
		Pagination m.PaginationResp
		Data       []*m.ColdStorage
	}
	Response.Pagination = PaginationResp
	Response.Data = comp

	if err := json.NewEncoder(w).Encode(Response); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
