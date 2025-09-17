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

// GetAllKBData returns a all records present in kb_data table
func GetAllKBData(w http.ResponseWriter, r *http.Request) {

	kbdata, err := dao.GetAllKBData()
	if err != nil {
		slog.Error("Recordes not found", "error", err.Error())
		http.Error(w, "Failed to encode fetch records", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(kbdata); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

}

// GetKBDataByParam returns a kb_data records based on parameter
func GetKBDataByParam(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")

	KbRoot, err := dao.GetKBDataByParam(key, value)
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

func GetOrderDEtails(w http.ResponseWriter, r *http.Request) {
	var condition string
	var kbroot []*model.KbRoot
	var rawQuery model.RawQuery
	rawQuery.Type = "KbRoot"
	rawQuery.Host = utils.RestHost
	rawQuery.Port = utils.RestPort
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

	ord, err := dao.GetOrderDetails(condition)
	if err != nil {
		slog.Error("Recordes not found", "error", err.Error())
		http.Error(w, "Failed to encode fetch records", http.StatusInternalServerError)
		return
	}
	for i, v := range ord {
		rawQuery.Query = utils.JoinStr(`				
		SELECT kb_root.* 
		FROM kb_root
		LEFT OUTER JOIN kb_data ON kb_data.id = kb_root.kb_data_id 
		LEFT OUTER JOIN kb_extension ON kb_extension.id = kb_data.kb_extension_id 
		WHERE kb_extension.vendor_id = (
			SELECT id FROM vendors 
			WHERE vendor_code LIKE 'I%' 
			ORDER BY id ASC 
			LIMIT 1
		) 
		AND kb_data.compound_id = `, strconv.Itoa(v.CompoundId), ` 
		AND kb_root.status != '3'
		AND kb_root.status != '-1';`)

		rawQuery.RawQry(&kbroot)
		ord[i].InventoryKanbanInProcessQty = len(kbroot)

	}

	if err := json.NewEncoder(w).Encode(ord); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func GetCustomerOrderDetails(w http.ResponseWriter, r *http.Request) {
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
	orddetails, err := dao.GetCustomerOrderDetails(con)
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

func GetAllDetailsForOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var kbData m.KbData

	if err := json.NewDecoder(r.Body).Decode(&kbData); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		log.Println("Error decoding JSON:", err)
		return
	}

	OrderDetails, err := dao.GetAllDetailsForOrder(strconv.Itoa(kbData.Id))
	if err != nil {
		http.Error(w, "Failed to fetch order details", http.StatusInternalServerError)
		log.Println("Error fetching order details:", err)
		return
	}

	KbRoot, err := dao.GetKbRootByParam("kb_data_id", strconv.Itoa(kbData.Id))
	if err != nil {
		http.Error(w, "Failed to fetch Kanban details", http.StatusInternalServerError)
		log.Println("Error fetching Kanban details:", err)
		return
	}
	for _, rootData := range KbRoot {
		FirstTxn, LastTnx, ProdLineName, err := dao.GetFirstAndLastTransactionWithProdLine(strconv.Itoa(rootData.Id))
		if err != nil {
			http.Error(w, "Failed to fetch Kanban details", http.StatusInternalServerError)
			log.Println("Error fetching Kanban details:", err)
			return
		}
		OrderDetails.KanbanDetails = append(OrderDetails.KanbanDetails, m.Kanban{
			ID:          rootData.Id,
			LotNo:       rootData.LotNo,
			CreatedOn:   FirstTxn.CreatedOn,
			CompletedOn: LastTnx.CreatedOn,
			ProdLine:    ProdLineName,
		})
	}

	// Send the response as JSON
	if err := json.NewEncoder(w).Encode(OrderDetails); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func GetPendingOrderDetailsBySearchAndPagination(w http.ResponseWriter, r *http.Request) {
	var condition string
	var kbroot []*model.KbRoot
	var rawQuery model.RawQuery
	rawQuery.Type = "KbRoot"
	rawQuery.Host = utils.RestHost
	rawQuery.Port = utils.RestPort
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

	ord, PaginationResp, err := dao.GetOrderDetailsBySearchAndPagination(data.PaginationReq, data.Conditions)
	if err != nil {
		slog.Error("Recordes not found", "error", err.Error())
		http.Error(w, "Failed to encode fetch records", http.StatusInternalServerError)
		return
	}
	for i, v := range ord {
		rawQuery.Query = utils.JoinStr(`				
		SELECT kb_root.* 
		FROM kb_root
		LEFT OUTER JOIN kb_data ON kb_data.id = kb_root.kb_data_id 
		LEFT OUTER JOIN kb_extension ON kb_extension.id = kb_data.kb_extension_id 
		WHERE kb_extension.vendor_id = (
			SELECT id FROM vendors 
			WHERE vendor_code LIKE 'I%' 
			ORDER BY id ASC 
			LIMIT 1
		) 
		AND kb_data.compound_id = `, strconv.Itoa(v.CompoundId), ` 
		AND kb_root.status != '3'
		AND kb_root.status != '-1';`)

		rawQuery.RawQry(&kbroot)
		ord[i].InventoryKanbanInProcessQty = len(kbroot)

	}

	var Response struct {
		Pagination m.PaginationResp
		Data       []*m.OrderDetails
	}
	Response.Pagination = PaginationResp
	Response.Data = ord

	if err := json.NewEncoder(w).Encode(Response); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
