package server

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"strconv"

	m "irpl.com/kanban-commons/model"
	"irpl.com/kanban-commons/utils"
	"irpl.com/kanban-dao/dao"
)

// Get All Vendors
func GetVendorsAllrecords(w http.ResponseWriter, r *http.Request) {
	var vendorEntries []m.Vendors
	var requestData struct {
		Conditions []string `json:"Conditions"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		slog.Error("Fail to decode request data", "error", err.Error())
		utils.SetResponse(w, http.StatusInternalServerError, "Fail: Fail to decode request data")
	}

	packing := r.URL.Query().Get("status")

	if packing == "packing" {
		requestData.Conditions = append(requestData.Conditions, "vendor_code != 'Inventory001'")
		vendorEntries, err = dao.GetAllVendorEntriesByCondition(requestData.Conditions)
	} else {
		vendorEntries, err = dao.GetAllVendorEntriesByCondition(requestData.Conditions)
	}

	if err != nil {
		slog.Error("Recordes not found", "error", err.Error())
		http.Error(w, "Failed to encode fetch records", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(vendorEntries); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// GetKbRootByParam returns a kb_root records based on parameter
func GetVendorsByParam(w http.ResponseWriter, r *http.Request) {
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

// Create Vendors
func CreateVendors(w http.ResponseWriter, r *http.Request) {
	var data m.Vendors
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err == nil {
		var OldData m.Vendors = data
		err := dao.CreateNewOrUpdateExistingVendors(&data)
		if err != nil {
			// Create Log
			if OldData.ID == 0 {
				sysLog := m.SystemLog{
					Message:     "AddVendor: Fail to add vendor " + data.VendorCode,
					MessageType: "ERROR",
					IsCritical:  false,
					CreatedBy:   data.CreatedBy,
				}
				utils.CreateSystemLogInternal(sysLog)
			} else {
				sysLog := m.SystemLog{
					Message:     "UpdateVendor: Fail to update vendor " + data.VendorCode,
					MessageType: "ERROR",
					IsCritical:  false,
					CreatedBy:   data.CreatedBy,
				}
				utils.CreateSystemLogInternal(sysLog)
			}
			slog.Error("Record creation failed", "error", err.Error())
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to create record")
			return
		}
		// Create Log
		if OldData.ID == 0 {
			sysLog := m.SystemLog{
				Message:     "AddVendor: Vendor " + data.VendorCode + " added successfully",
				MessageType: "SUCCESS",
				IsCritical:  false,
				CreatedBy:   data.CreatedBy,
			}
			utils.CreateSystemLogInternal(sysLog)
			utils.SetResponse(w, http.StatusOK, "Successfully created vendor")
			return
		} else {
			sysLog := m.SystemLog{
				Message:     "UpdateVendor: Vendor " + data.VendorCode + " updated successfully",
				MessageType: "SUCCESS",
				IsCritical:  false,
				CreatedBy:   data.CreatedBy,
			}
			utils.CreateSystemLogInternal(sysLog)
			utils.SetResponse(w, http.StatusOK, "Successfully updated vendor")
			return
		}
	} else {
		// Create Log
		sysLog := m.SystemLog{
			Message:     "AddVendor: Fail to add vendor ",
			MessageType: "ERROR",
			IsCritical:  false,
			CreatedBy:   data.CreatedBy,
		}
		utils.CreateSystemLogInternal(sysLog)
		slog.Error("Record creation failed", "error", err.Error())
		utils.SetResponse(w, http.StatusInternalServerError, "Failed to create record")
	}
}

// Delete Vendor records
func DeleteVendorsById(w http.ResponseWriter, r *http.Request) {
	value := r.URL.Query().Get("id")

	id, _ := strconv.Atoi(value)

	err := dao.DeleteVendor(id)
	if err != nil {
		utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to delete record")
		slog.Error("Record deletion failed", "id", id, "error", err.Error())
		return
	}
	utils.SetResponse(w, http.StatusOK, "Success: successfully deleted record")
}

func GetVendorDetailsByVendorCode(w http.ResponseWriter, r *http.Request) {
	var vendor m.Vendors
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&vendor); err == nil {
		vendor, err := dao.GetVendorByParam("vendor_code", vendor.VendorCode)
		if err != nil || len(vendor) == 0 {
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to get vendor record")
			return
		}
		if err := json.NewEncoder(w).Encode(vendor[0]); err != nil {
			slog.Error("%s - error - %s", "Failed to add order data in KbData table", err.Error())
			utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to update record")
			return
		}
	} else {
		slog.Error("Record creation failed", "error", err.Error())
		utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to create record")
	}
}

// GetVendorSearchPaginationRecords handles paginated vendor search requests
func GetVendorSearchPaginationRecords(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		Pagination m.PaginationReq `json:"pagination"`
		Conditions []string        `json:"conditions"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		slog.Error("Failed to decode request data", "error", err.Error())
		utils.SetResponse(w, http.StatusBadRequest, "Fail: Unable to decode request")
		return
	}

	// Optional packing-specific condition
	if r.URL.Query().Get("status") == "packing" {
		requestData.Conditions = append(requestData.Conditions, "vendor_code != 'Inventory001'")
	}

	vendorEntries, paginationData, err := dao.GetAllVendorBySearchAndPagination(requestData.Pagination, requestData.Conditions)
	if err != nil {
		slog.Error("Failed to fetch vendor records", "error", err.Error())
		http.Error(w, "Failed to fetch records", http.StatusInternalServerError)
		return
	}

	response := struct {
		Pagination m.PaginationResp `json:"pagination"`
		Data       []*m.Vendors     `json:"data"`
	}{
		Pagination: paginationData,
		Data:       vendorEntries,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		slog.Error("Error encoding response", "error", err.Error())
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
