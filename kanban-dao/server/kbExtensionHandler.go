package server

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	m "irpl.com/kanban-commons/model"
	"irpl.com/kanban-dao/dao"
)

// GetAllKbExtensions returns a all records present in kb_extension table
func GetAllKbExtensions(w http.ResponseWriter, r *http.Request) {
	kbextension, err := dao.GetAllKbExtensions()
	if err != nil {
		slog.Error("Recordes not found", "error", err.Error())
		http.Error(w, "Failed to encode fetch records", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(kbextension); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// GetKbExtensionsByParam returns a kb_extension records based on parameter
func GetKbExtensionsByParam(w http.ResponseWriter, r *http.Request) {

	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")

	KbRoot, err := dao.GetKbExtensionsByParam(key, value)
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

func GetOrderDetaislForOrderPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var requestData struct {
		ID             int    `json:"ID"`
		CompoundName   string `json:"CompoundName"`
		CellNo         string `json:"CellNo"`
		DemandDateTime string `json:"DemandDateTime"`
		NoOFLots       int    `json:"NoOFLots"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestData)

	if err != nil {
		http.Error(w, "Failed to decode data", http.StatusInternalServerError)
		return
	}
	var OrderDetails m.OrderDetails
	id := strconv.Itoa(requestData.ID)
	kbdata, _ := dao.GetKBDataByParam("id", id)
	Status, orderID, err := dao.GetOrderIdAndStatusByKDid(kbdata[0].KbExtensionID)
	if err != nil {
		http.Error(w, "Failed to retive data", http.StatusInternalServerError)
		return
	}
	OrderDetails.OrderId = orderID
	OrderDetails.Status = Status
	if err := json.NewEncoder(w).Encode(OrderDetails); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

}

func UpdateCustomerOrderStatus(KbRootId []string, status, modified_by string) error {
	for _, KBRootID := range KbRootId {
		// Fetch kb_root data
		KRData, err := dao.GetKbRootByParam("id", KBRootID)
		if err != nil {
			return err
		}

		// Fetch kb_data from kb_root
		KDData, err := dao.GetKBDataByParam("id", strconv.Itoa(KRData[0].KbDataId))
		if err != nil {
			return err
		}

		// Fetch kb_extension from kb_data
		KEData, err := dao.GetKbExtensionsByParam("id", strconv.Itoa(KDData[0].KbExtensionID))
		if err != nil {
			return err
		}

		// Get current and new statuses
		currentStatus := KEData[0].Status
		newStatus := status

		// Define status priorities
		statusPriority := map[string]int{
			"approved":            0,
			"InProductionLine":    1,
			"InProductionProcess": 2,
			"quality":             3,
			"dispatch":            5,
		}
		getPriority := func(s string) int {
			if strings.HasPrefix(s, "dispatched(") {
				return 4 // Assign priority between "InProductionProcess" (2) and "dispatch" (4)
			}
			if val, exists := statusPriority[s]; exists {
				return val
			}
			return -1
		}
		if getPriority(newStatus) <= getPriority(currentStatus) {
			continue
		}

		// Update status logic here
		err = dao.UpdateCustomerOrderStatus(KEData[0].Id, newStatus, modified_by)
		if err != nil {
			return err
		}
	}
	return nil
}
