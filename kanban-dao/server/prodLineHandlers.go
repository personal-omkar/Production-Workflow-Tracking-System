package server

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"sort"
	"strconv"
	"time"

	m "irpl.com/kanban-commons/model"
	"irpl.com/kanban-commons/utils"
	"irpl.com/kanban-dao/dao"
)

// Get All Prod line
func GetProdLineAllrecords(w http.ResponseWriter, r *http.Request) {

	prodLineEntries, err := dao.GetAllProductionLineEntries()
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

func GetAllProductionLineRecords(w http.ResponseWriter, r *http.Request) {

	prodLineEntries, err := dao.GetProdLinesWithCellsAndStatus()
	if err != nil {
		slog.Error("Recordes not found", "error", err.Error())
		http.Error(w, "Failed to encode fetch records", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(prodLineEntries)
	if err != nil {
		http.Error(w, "Error whille encoding data", http.StatusInternalServerError)
		// w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func CreateProdLineWithProcesses(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Decode the JSON payload
	var requestData m.AddLineStruct

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		// Create Log
		sysLog := m.SystemLog{
			Message:     "CreateProductionLine : Fail to create new production line",
			MessageType: "ERROR",
			IsCritical:  false,
			CreatedBy:   requestData.CreatedBy,
		}
		utils.CreateSystemLogInternal(sysLog)
		http.Error(w, "Invalid input format", http.StatusBadRequest)
		return
	}

	// Set system-generated fields for prodLine
	currentTime := time.Now().UTC()
	prodLine := m.ProdLine{
		Name:        requestData.LineName + " (Heijunka)",
		Description: requestData.LineDescription,
		CreatedBy:   requestData.CreatedBy,
		CreatedOn:   currentTime,
		ModifiedOn:  currentTime,
		Status:      true,
	}

	// Step 1: Create prod_line
	id, err := dao.CreateProdLine(&prodLine)
	if err != nil {
		// Create Log
		sysLog := m.SystemLog{
			Message:     "CreateProductionLine : Fail to create new production line " + requestData.LineName,
			MessageType: "ERROR",
			IsCritical:  false,
			CreatedBy:   requestData.CreatedBy,
		}
		utils.CreateSystemLogInternal(sysLog)
		log.Println("Error saving production line:", err)
		http.Error(w, "Failed to save production line", http.StatusInternalServerError)
		return
	}

	for i := range requestData.ProcessOrders {
		requestData.ProcessOrders[i].Order += 1
	}
	defaultStart := struct {
		ProdProcessID string `json:"prod_process_id"`
		Order         int    `json:"order"`
		GroupName     string `json:"group_name"`
	}{
		ProdProcessID: "1",
		Order:         1,
	}

	defaultEnd := struct {
		ProdProcessID string `json:"prod_process_id"`
		Order         int    `json:"order"`
		GroupName     string `json:"group_name"`
	}{
		ProdProcessID: "2",
		Order:         len(requestData.ProcessOrders) + 2,
	}
	requestData.ProcessOrders = append([]m.ProcessOrders{defaultStart}, requestData.ProcessOrders...)
	requestData.ProcessOrders = append(requestData.ProcessOrders, defaultEnd)
	err = dao.CreateProdProcessLines(id, requestData.ProcessOrders)
	if err != nil {
		// Create Log
		sysLog := m.SystemLog{
			Message:     "CreateProductionLine : Fail to create new production line " + requestData.LineName,
			MessageType: "ERROR",
			IsCritical:  false,
			CreatedBy:   requestData.CreatedBy,
		}
		utils.CreateSystemLogInternal(sysLog)
		log.Println("Error saving production process lines:", err)
		http.Error(w, "Failed to save production process lines", http.StatusInternalServerError)
		return
	}
	// Create Log
	sysLog := m.SystemLog{
		Message:     "CreateProductionLine : Successfully created new production line " + requestData.LineName,
		MessageType: "SUCCESS",
		IsCritical:  false,
		CreatedBy:   requestData.CreatedBy,
	}
	utils.CreateSystemLogInternal(sysLog)
	utils.SetResponse(w, http.StatusCreated, "Production line and processes created successfully")
}

func GetProductionLineStatus(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		LineNo int `json:"line_no"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	prodLineStatus, err := dao.GetLatestRecordsForProdID(requestBody.LineNo)
	if err != nil {
		http.Error(w, "Failed to fetch records", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(prodLineStatus)
	if err != nil {
		http.Error(w, "Error while encoding data", http.StatusInternalServerError)
		return
	}
}

// DeleteProductionLineById Delet production line with there data
func DeleteProductionLineById(w http.ResponseWriter, r *http.Request) {
	idstr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idstr)
	var err error

	err = dao.DeleteKbTransactionByProdLine(id)
	if err != nil {
		slog.Error("Recordes not found", "error", err.Error())
		http.Error(w, "Failed to delete records from kb_transaction", http.StatusInternalServerError)
		return
	} else {
		err = dao.DeleteProductionProcessLineByProdLine(id)
		if err != nil {
			slog.Error("Recordes not found", "error", err.Error())
			http.Error(w, "Failed to delete records from prod_process_line", http.StatusInternalServerError)
			return
		} else {
			err = dao.DeleteProductionLine(id)
			if err != nil {
				slog.Error("Recordes not found", "error", err.Error())
				http.Error(w, "Failed to delete records from prod_line", http.StatusInternalServerError)
				return
			} else {
				utils.SetResponse(w, http.StatusOK, "Success: successfully deleted record")
			}
		}
	}

}

func DeleteProductionLineCellDataByProductionLineId(w http.ResponseWriter, r *http.Request) {
	idstr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idstr)
	var err error

	err = dao.DeleteKbDataByProdLine(id)
	if err != nil {
		slog.Error("Recordes not found", "error", err.Error())
		http.Error(w, "Failed to delete records from kb_data", http.StatusInternalServerError)
		return
	} else {
		err = dao.DeleteKbExtensionByProdLine(id)
		if err != nil {
			slog.Error("Recordes not found", "error", err.Error())
			http.Error(w, "Failed to delete records from kb_extension", http.StatusInternalServerError)
			return
		} else {
			err = dao.DeleteProductionLineDataInKbTransactionByProdLine(id)
			if err != nil {
				slog.Error("Recordes not found", "error", err.Error())
				http.Error(w, "Failed to delete records from kb_transcation", http.StatusInternalServerError)
				return
			} else {
				err = dao.DeleteKbRootByProdLine(id)
				if err != nil {
					slog.Error("Recordes not found", "error", err.Error())
					http.Error(w, "Failed to delete records from kb_root", http.StatusInternalServerError)
					return
				} else {
					utils.SetResponse(w, http.StatusOK, "Success: successfully deleted record")
				}
			}

		}
	}

}

// Created an set to store all unique prod-line id
type Set map[int]struct{}

func DeleteProductionLineCell(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	var payload struct {
		KRid   []string `json:"KRid"`
		UserId string   `json:"userID"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	prodLineSet := Set{}

	// Get whole root data  by kbRootID
	kbroot, _ := dao.GetMultiKbRootByParam("id", payload.KRid)

	// Delete Kanban from production line
	for i := range kbroot {

		// get production line id from which kanban is deleted
		prodLineId, err := dao.GetProdLineIDByKbRootID(kbroot[i].Id)
		if err != nil {
			slog.Error("Prod Line Recordes not found", "error", err.Error())
			utils.SetResponse(w, http.StatusInternalServerError, "Fail: Failed to update cell status")
			return
		}
		// apepnd prodLine Id to set
		prodLineSet[prodLineId] = struct{}{}

		// Chnage its Running no. , initial no. n clear lot No.
		kbroot[i].RunningNo = -1
		kbroot[i].InitialNo = -1
		kbroot[i].Status = "0"
		kbroot[i].ModifiedOn = time.Now()
		kbroot[i].LotNo = ""
		_, err = dao.CreateNewOrUpdateExistingKbRoot(&kbroot[i])
		if err != nil {
			slog.Error("Recordes not found", "error", err.Error())
			utils.SetResponse(w, http.StatusInternalServerError, "Fail: Failed to update cell status")
			return
		}

		// Get transaction for that root id
		kbt, _ := dao.GetKbTransactionByParam("kb_root_id", strconv.Itoa(kbroot[i].Id))
		if len(kbt) > 0 {
			if len(kbt) == 1 {
				// delete transaction if it is only one(one transaction means -> only lined up transaction is created)
				err = dao.DeleteKbTransactionByParam("id", strconv.Itoa(kbt[0].Id))
				if err != nil {
					slog.Error("Recordes not found", "error", err.Error())
					utils.SetResponse(w, http.StatusInternalServerError, "Fail: Failed to update cell status")
					return
				}
			} else {
				slog.Error("Failed to Delete kanban")
				utils.SetResponse(w, http.StatusInternalServerError, "Fail: Failed to Delete kanban")
				return
			}
		} else {
			http.Error(w, "Failed to update cell status", http.StatusInternalServerError)
			utils.SetResponse(w, http.StatusInternalServerError, "Fail: Failed to update cell status")
			return
		}

	}

	//Update Running Number (we have prod line whose kanban was deleted)
	for key := range prodLineSet {
		// Get all lined-up roots in that line
		KbRoots, _ := dao.GetLinedUpKBRootsByProdLineID(key)
		if len(KbRoots) > 0 {
			sort.Slice(KbRoots, func(i, j int) bool {
				return KbRoots[i].RunningNo < KbRoots[j].RunningNo
			})
			for i := 1; i <= len(KbRoots); i++ {
				KbRoots[i-1].RunningNo = i
				KbRoots[i-1].ModifiedBy = payload.UserId
				KbRoots[i-1].ModifiedOn = time.Now()
				// update all running number in line
				_, err := dao.CreateNewOrUpdateExistingKbRoot(&KbRoots[i-1])
				if err != nil {
					slog.Error("Failed to update running no.")
					utils.SetResponse(w, http.StatusInternalServerError, "Fail: Failed to update running no.")
					return
				}
			}
		}
	}

	// Now lets check for status Updation

	// we have data which was updated
	for _, value := range kbroot {

		// now get kb_data id
		allStatus0 := true
		KRdata, err := dao.GetKbRootByParam("kb_data_id", strconv.Itoa(value.KbDataId))
		if err != nil {
			slog.Error("Failed to update order status.")
			utils.SetResponse(w, http.StatusInternalServerError, "Fail: Failed to update order status.")
			return
		}
		for _, RootValue := range KRdata {
			if RootValue.Status != "0" {
				allStatus0 = false
				break
			}
		}
		if allStatus0 {
			// get whole kbdata
			kbData, _ := dao.GetKBDataByParam("id", strconv.Itoa(value.KbDataId))
			// get extension
			kbe, _ := dao.GetKbExtensionsByParam("id", strconv.Itoa(kbData[0].KbExtensionID))
			kbe[0].Status = "approved"
			kbe[0].ModifiedOn = time.Now()
			// Update extension status i.e. order status
			_, err = dao.CreateNewOrUpdateExistingKbExtension(&kbe[0])
			if err != nil {
				slog.Error("Recordes not found", "error", err.Error())
				utils.SetResponse(w, http.StatusInternalServerError, "Fail: Failed to update cell status")
			}
		}
	}
}

func GetProdLineByParam(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var payload m.ProdLine
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	ProdLineData, err := dao.GetProdLineByParam("id", strconv.Itoa(payload.Id))
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

func EditProdLine(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var payload m.ProdLine
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if payload.MoveToLineID == 0 {
		ProdLineData, err := dao.GetProdLineByParam("id", strconv.Itoa(payload.Id))
		if err != nil {
			slog.Error("Failed to get production line data.")
			utils.SetResponse(w, http.StatusInternalServerError, "Fail: Failed to get production line data.")
			return
		}

		payload.CreatedBy = ProdLineData[0].CreatedBy
		payload.CreatedOn = ProdLineData[0].CreatedOn
		err = dao.CreateNewOrUpdateExistingProductLine(&payload)
		if err != nil {
			slog.Error("Failed to update production line data.")
			utils.SetResponse(w, http.StatusInternalServerError, "Fail: Failed to update production line data.")
			return
		}
		// Create Log
		sysLog := m.SystemLog{
			Message:     "UpdateProductionLine : Successfully updated production line " + payload.Name,
			MessageType: "SUCCESS",
			IsCritical:  false,
			CreatedBy:   payload.CreatedBy,
		}
		utils.CreateSystemLogInternal(sysLog)
		utils.SetResponse(w, http.StatusOK, "Successfully updated the production line")
	} else {
		ProdLineData, _ := dao.GetProdLinesWithCellsAndStatusByID(payload.Id)

		// Delete existing transaction from KBtransaction Table
		for _, j := range ProdLineData.Cells {
			err := dao.DeleteKBTransactionsByRootID(j.KRId)
			if err != nil {
				slog.Error("Failed to change Kanban production line", "error", err.Error())
				http.Error(w, "Failed to change Kanban production line", http.StatusInternalServerError)
				return
			}
		}
		for _, i := range ProdLineData.Cells {
			kbrootid, _ := strconv.Atoi(i.KRId)
			ppl, _ := dao.GetProducProcesstLineByParamAndOrder("prod_line_id", strconv.Itoa(payload.MoveToLineID), "1 ")

			// check for unique_status_KbRootId CONSTRAINT,  get status and kb_root_id and check if pair exists in kb_transaction or not
			exists, err := dao.StatusKbRootExists(strconv.Itoa(ppl[0].Order), kbrootid)
			if err != nil {
				slog.Error("Failed to check Status", "error", err.Error())
				http.Error(w, "Failed to check Status", http.StatusInternalServerError)
				return
			}
			if exists {
				continue
			}
			if err == nil {
				kbTransaction := m.KbTransaction{
					ProdProcessId:     1,
					Status:            strconv.Itoa(ppl[0].Order),
					KbRootId:          kbrootid,
					ProdProcessLineID: ppl[0].Id,
					StartedOn:         time.Now(),
					CreatedOn:         time.Now(),
					CreatedBy:         payload.ModifiedBy,
				}
				_, err := dao.CreateNewOrUpdateExistingKbTransactionData(&kbTransaction)
				if err != nil {
					slog.Error("Failed to create kbtransaction records", "error", err.Error())
					http.Error(w, "Failed to create kbtransaction records", http.StatusInternalServerError)
					return
				}
				var ppl int
				ppldt, _ := dao.GetProdLinesWithCellsAndStatus()
				for _, i := range ppldt {
					if i.ProdLineID == payload.MoveToLineID {
						ppl = len(i.Cells)
					}
				}

				kbrootdt, _ := dao.GetKbRootByParam("id", i.KRId)
				kbdata, _ := dao.GetKBDataByParam("id", strconv.Itoa(kbrootdt[0].KbDataId))
				LotNumber := UpdateLotNumber(strconv.Itoa(payload.MoveToLineID), kbrootdt[0].LotNo)
				kbroot := m.KbRoot{
					Id:          kbrootid,
					InInventory: kbrootdt[0].InInventory,
					CreatedOn:   kbrootdt[0].CreatedOn,
					CreatedBy:   kbrootdt[0].CreatedBy,
					Status:      "1",
					KbDataId:    kbdata[0].Id,
					RunningNo:   ppl,
					InitialNo:   ppl,
					ModifiedBy:  payload.ModifiedBy,
					ModifiedOn:  time.Now(),
					LotNo:       LotNumber,
					Remark:      "This root has been transferred from one line to another.",
				}
				_, err = dao.CreateNewOrUpdateExistingKbRoot(&kbroot)
				if err != nil {
					slog.Error("Failed to update kbroot inital and running number records", "error", err.Error())
					http.Error(w, "Failed to update kbroot inital and running number records", http.StatusInternalServerError)
					return
				}
			} else {
				slog.Error("Failed to create ProdProcessLine records", "error", err.Error())
				http.Error(w, "Failed to create ProdProcessLine records", http.StatusInternalServerError)
				return
			}
		}

		ProdLine, err := dao.GetProdLineByParam("id", strconv.Itoa(payload.Id))
		if err != nil {
			slog.Error("Failed to get production line data.")
			utils.SetResponse(w, http.StatusInternalServerError, "Fail: Failed to get production line data.")
			return
		}

		payload.CreatedBy = ProdLine[0].CreatedBy
		payload.CreatedOn = ProdLine[0].CreatedOn
		err = dao.CreateNewOrUpdateExistingProductLine(&payload)
		if err != nil {
			slog.Error("Failed to update production line data.")
			utils.SetResponse(w, http.StatusInternalServerError, "Fail: Failed to update production line data.")
			return
		}
		// Create Log
		sysLog := m.SystemLog{
			Message:     "UpdateProductionLine : Successfully updated production line " + payload.Name,
			MessageType: "SUCCESS",
			IsCritical:  false,
			CreatedBy:   payload.CreatedBy,
		}
		utils.CreateSystemLogInternal(sysLog)
		utils.SetResponse(w, http.StatusOK, "Successfully updated the production line")
	}
}

func GetLineUpProcessesByLineId(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var payload m.ProdLine
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	ProdLineData, err := dao.GetProdLinesWithCellsAndStatusByID(payload.Id)
	if err != nil {
		slog.Error("Failed to get production line data.")
		utils.SetResponse(w, http.StatusInternalServerError, "Fail: Failed to get production line data.")
		return
	}

	if ProdLineData == nil || len(ProdLineData.Cells) == 0 {
		ProdLine, err := dao.GetProdLineByParam("id", strconv.Itoa(payload.Id))
		if len(ProdLine) == 0 || err != nil {
			slog.Error("Failed to get production line data.")
			utils.SetResponse(w, http.StatusInternalServerError, "Fail: Failed to get production line data.")
			return
		}
		ProdLineData.ProdLineID = ProdLine[0].Id
		ProdLineData.ProdLineName = ProdLine[0].Name
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ProdLineData)
}

// CreateNewOrUpdateExistingProductionLine create a new production line entry or update existing production line
func CreateNewOrUpdateExistingProductionLine(w http.ResponseWriter, r *http.Request) {
	var data m.ProdLine
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err == nil {
		var msg string
		if data.Id != 0 {
			msg = "Success: successfully updated record"
		} else {
			msg = "Success: successfully created record"
		}
		err := dao.CreateNewOrUpdateExistingProductLine(&data)
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

// GetAllProdLineBySearchAndPagination returns a all records present in prodline table
func GetAllProdLineBySearchAndPagination(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		Pagination m.PaginationReq `json:"pagination"`
		Conditions []string        `json:"conditions"`
	}
	// Decode JSON request
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	prodLine, PaginationResp, err := dao.GetAllProdLineBySearchAndPagination(requestData.Pagination, requestData.Conditions)
	if err != nil {
		slog.Error("Recordes not found", "error", err.Error())
		http.Error(w, "Failed to encode fetch records", http.StatusInternalServerError)
		return
	}

	var Response struct {
		Pagination m.PaginationResp
		Data       []*m.ProdLine
	}
	Response.Pagination = PaginationResp
	Response.Data = prodLine

	if err := json.NewEncoder(w).Encode(Response); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
