package server

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"irpl.com/kanban-commons/model"
	"irpl.com/kanban-commons/utils"
	"irpl.com/kanban-dao/dao"
)

// Get All Compounds
func GetCompountsAllrecords(w http.ResponseWriter, r *http.Request) {

	prodLineEntries, err := dao.GetAllCompoundEntries()
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

func GetKanbanByVendors(w http.ResponseWriter, r *http.Request) {
	var requestData model.Vendors
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		slog.Error("Failed to decode request data", "error", err.Error())
		utils.SetResponse(w, http.StatusBadRequest, "Fail: Unable to decode request data")
		return
	}
	var compounData []model.CompoundsDataByVendor
	vendor, err := dao.GetVendorByParam("id", strconv.Itoa(requestData.ID))
	if len(vendor) == 0 || err != nil {
		log.Printf("Error while getting vendor: %v", err)
		http.Error(w, "Vendor not found", http.StatusInternalServerError)
		return
	}
	kbextensions, _ := dao.GetKbExtensionsByParam("vendor_id", strconv.Itoa(vendor[0].ID))

	for _, kbextensionsData := range kbextensions {
		kbdata, _ := dao.GetKBDataByParam("kb_extension_id", strconv.Itoa(kbextensionsData.Id))
		if len(kbdata) > 0 {
			kbroot, _ := dao.GetKbRootByParamAndStatus("kb_data_id", strconv.Itoa(kbdata[0].Id), "0")
			if len(kbroot) > 0 {
				for i := 1; i <= len(kbroot); i++ {
					compound, err := dao.GetCompoundDataByParamAndCondition("id", strconv.Itoa(kbdata[0].CompoundId), []string{})
					if err != nil || len(compound) == 0 {
						continue
					}
					var tempComp model.CompoundsDataByVendor
					tempComp.CellNo = kbdata[0].CellNo
					tempComp.CompoundCode = compound[0].Id
					tempComp.CompoundName = compound[0].CompoundName
					tempComp.KbRootId = kbroot[i-1].Id
					tempComp.ModifiedOn = kbextensionsData.ModifiedOn
					tempComp.DemandDate = kbdata[0].DemandDateTime
					tempComp.KanbanNo = kbroot[i-1].KanbanNo
					tempComp.CustomerNote = kbdata[0].Note
					compounData = append(compounData, tempComp)
				}
			}
		}
	}

	if err := json.NewEncoder(w).Encode(compounData); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func GetKanbanForAllVendors(w http.ResponseWriter, r *http.Request) {
	var kanbanList []model.VendorKanban

	vendors, _ := dao.GetVendorByParam("isactive", true)

	for _, v := range vendors {
		var compoundList []model.CompoundsDataByVendor

		kbextensions, _ := dao.GetKbExtensionsByParam("vendor_id", strconv.Itoa(v.ID))
		for _, kbext := range kbextensions {
			kbdata, _ := dao.GetKBDataByParam("kb_extension_id", strconv.Itoa(kbext.Id))
			if len(kbdata) == 0 {
				continue
			}

			kbroot, _ := dao.GetKbRootByParamAndStatus("kb_data_id", strconv.Itoa(kbdata[0].Id), "0")
			if len(kbroot) == 0 {
				continue
			}

			for i := 1; i <= len(kbroot); i++ {
				compound, err := dao.GetCompoundDataByParamAndCondition("id", strconv.Itoa(kbdata[0].CompoundId), []string{})
				if err != nil || len(compound) == 0 {
					continue
				}

				compoundList = append(compoundList, model.CompoundsDataByVendor{
					CellNo:       kbdata[0].CellNo,
					CompoundCode: compound[0].Id,
					CompoundName: compound[0].CompoundName,
					KbRootId:     kbroot[i-1].Id,
					ModifiedOn:   kbext.ModifiedOn,
					DemandDate:   kbdata[0].DemandDateTime,
					KanbanNo:     kbroot[i-1].KanbanNo,
					CustomerNote: kbdata[0].Note,
				})
			}
		}

		kanbanList = append(kanbanList, model.VendorKanban{
			Vendor:    v,
			Compounds: compoundList,
		})
	}

	if err := json.NewEncoder(w).Encode(kanbanList); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

}

// Get Compounds data based on vendor records
func GetCompoundsByVendors(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		Conditions []string `json:"Conditions"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		slog.Error("Failed to decode request data", "error", err.Error())
		utils.SetResponse(w, http.StatusBadRequest, "Fail: Unable to decode request data")
		return
	}
	var compounData []model.CompoundsDataByVendor
	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")
	if key == "" || value == "" {
		log.Printf("Insufficient data")
		http.Error(w, "Insufficient data", http.StatusInternalServerError)
		return
	}
	vendor, err := dao.GetVendorByParam(key, value)
	if len(vendor) == 0 || err != nil {
		log.Printf("Vendor not found: %v", err)
		http.Error(w, "Vendor not found", http.StatusInternalServerError)
		return
	}
	kbextensions, _ := dao.GetKbExtensionsByParam("vendor_id", strconv.Itoa(vendor[0].ID))

	for _, kbextensionsData := range kbextensions {
		kbdata, _ := dao.GetKBDataByParam("kb_extension_id", strconv.Itoa(kbextensionsData.Id))
		if len(kbdata) > 0 {
			kbroot, _ := dao.GetKbRootByParamAndStatus("kb_data_id", strconv.Itoa(kbdata[0].Id), "0")
			if len(kbroot) > 0 {
				for i := 1; i <= len(kbroot); i++ {
					compound, err := dao.GetCompoundDataByParamAndCondition("id", strconv.Itoa(kbdata[0].CompoundId), requestData.Conditions)
					if err != nil || len(compound) == 0 {
						continue
					}
					var tempComp model.CompoundsDataByVendor
					tempComp.CellNo = kbdata[0].CellNo
					tempComp.CompoundCode = compound[0].Id
					tempComp.CompoundName = compound[0].CompoundName
					tempComp.KbRootId = kbroot[i-1].Id
					tempComp.ModifiedOn = kbextensionsData.ModifiedOn
					tempComp.DemandDate = kbdata[0].DemandDateTime
					tempComp.CreatedOn = kbroot[i-1].CreatedOn
					tempComp.CustomerNote = kbdata[0].Note
					tempComp.KanbanNo = kbroot[i-1].KanbanNo
					compounData = append(compounData, tempComp)
				}
			}
		}
	}

	if err := json.NewEncoder(w).Encode(compounData); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// Get packing Compounds data based on vendor records
func GetPackingCompoundsByVendors(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		Conditions []string `json:"Conditions"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		slog.Error("Failed to decode request data", "error", err.Error())
		utils.SetResponse(w, http.StatusBadRequest, "Fail: Unable to decode request data")
		return
	}
	var compounData []model.CompoundsDataByVendor
	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")

	vendor, _ := dao.GetVendorByParam(key, value)
	if len(vendor) == 0 {
		http.Error(w, "Vendor not found", http.StatusInternalServerError)
		return
	}
	kbextensions, _ := dao.GetKbExtensionsByParam("vendor_id", strconv.Itoa(vendor[0].ID))

	for _, kbextensionsData := range kbextensions {
		kbdata, _ := dao.GetKBDataByParam("kb_extension_id", strconv.Itoa(kbextensionsData.Id))
		if len(kbdata) > 0 {
			kbroot, _ := dao.GetKbRootByParamAndStatus("kb_data_id", strconv.Itoa(kbdata[0].Id), "3")
			if len(kbroot) > 0 {
				for i := range kbroot {
					compound, err := dao.GetCompoundDataByParamAndCondition("id", strconv.Itoa(kbdata[0].CompoundId), requestData.Conditions)
					if err != nil || len(compound) == 0 {
						continue
					}
					var tempComp model.CompoundsDataByVendor
					tempComp.CellNo = kbdata[0].CellNo
					tempComp.CompoundCode = compound[0].Id
					tempComp.CompoundName = compound[0].CompoundName
					tempComp.KbRootId = kbroot[i].Id
					tempComp.ModifiedOn = kbextensionsData.ModifiedOn
					tempComp.DemandDate = kbdata[0].DemandDateTime
					tempComp.CustomerNote = kbdata[0].Note
					tempComp.KanbanNo = kbroot[i].KanbanNo
					tempComp.QualityDoneTime = kbroot[i].QualityDoneTime
					compounData = append(compounData, tempComp)
				}
			}
		}
	}

	if err := json.NewEncoder(w).Encode(compounData); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// Get Quality testing Compounds data based on vendor records
func GetQualityCompoundsByVendors(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		Conditions []string `json:"Conditions"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		slog.Error("Failed to decode request data", "error", err.Error())
		utils.SetResponse(w, http.StatusBadRequest, "Fail: Unable to decode request data")
		return
	}
	var compounData []model.CompoundsDataByVendor
	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")

	vendor, _ := dao.GetVendorByParam(key, value)
	if len(vendor) == 0 {
		http.Error(w, "Vendor not found", http.StatusInternalServerError)
		return
	}
	kbextensions, _ := dao.GetKbExtensionsByParam("vendor_id", strconv.Itoa(vendor[0].ID))

	for _, kbextensionsData := range kbextensions {
		kbdata, _ := dao.GetKBDataByParam("kb_extension_id", strconv.Itoa(kbextensionsData.Id))
		if len(kbdata) > 0 {
			kbroot, _ := dao.GetKbRootByParamAndStatus("kb_data_id", strconv.Itoa(kbdata[0].Id), "2")
			if len(kbroot) > 0 {
				for i := range kbroot {
					if kbroot[i].Status == "2" && kbroot[i].Comment == "" {

						compound, err := dao.GetCompoundDataByParamAndCondition("id", strconv.Itoa(kbdata[0].CompoundId), requestData.Conditions)
						if err != nil || len(compound) == 0 {
							continue
						}
						var tempComp model.CompoundsDataByVendor
						tempComp.CellNo = kbdata[0].CellNo
						tempComp.CompoundCode = compound[0].Id
						tempComp.CompoundName = compound[0].CompoundName
						tempComp.KbRootId = kbroot[i].Id
						tempComp.ModifiedOn = kbextensionsData.ModifiedOn
						tempComp.DemandDate = kbdata[0].DemandDateTime
						tempComp.CustomerNote = kbdata[0].Note
						tempComp.KanbanNo = kbroot[i].KanbanNo
						tempComp.ModifiedOn = kbroot[i].ModifiedOn
						compounData = append(compounData, tempComp)
					}
				}
			}
		}
	}

	if err := json.NewEncoder(w).Encode(compounData); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// Add Compounds data based on vendor records
func AddCompoundsForVendor(w http.ResponseWriter, r *http.Request) {
	var cellNo string
	var data model.AddCompoundsByVendor
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&data); err == nil {
		// for i := 1; i <= data.Quantity; i++ {
		if len(data.VendorCode) > 0 && data.VendorCode[0] == 'I' {
			// Fetch inventory data
			inventoryList, _ := dao.GetAllInventory()

			// Convert inventory slice to a map for fast lookups
			inventoryMap := make(map[int]struct{}, len(inventoryList))
			for _, inventoryData := range inventoryList {
				inventoryMap[inventoryData.CompoundID] = struct{}{}
			}
			// Check if the compound exists in the inventory
			if _, exists := inventoryMap[data.CompoundCode]; !exists {
				slog.Error("This part is not available for inventory")
				http.Error(w, "This part is not available for inventory", http.StatusInternalServerError)
				return
			}
		}
		compoundData, _ := dao.GetCompoundDataByParam("id", strconv.Itoa(data.CompoundCode))
		//Fetch vendor data using vendor code to get vendor id
		vendorRec, err := dao.GetVendorDetailsByVendorCode(data.VendorCode)
		if err != nil {
			slog.Error("Failed to Fetch vendor data", "error", err.Error())
			http.Error(w, "Failed to fetch vendor data", http.StatusInternalServerError)
			return
		}
		kbExtension := model.KbExtension{
			VendorID:  vendorRec.ID,
			CreatedBy: data.UserID,
			Status:    "approved",
			CreatedOn: time.Now(),
		}
		//Creating kb_extensions records
		ExID, err := dao.CreateNewOrUpdateExistingKbExtension(&kbExtension)
		if err != nil {
			slog.Error("Failed to create kbextension records", "error", err.Error())
			http.Error(w, "Failed to create kbextension records", http.StatusInternalServerError)
			return
		} else {
			var kanbansRunningNumber []string
			// creating cell based on the + vendor code + Id
			cellNo = data.VendorCode + "/" + strconv.Itoa(ExID)
			for range data.Quantity {
				kanbanNo, err := dao.GenerateKanbanNumber()
				if err != nil {
					slog.Error("Failed to generate Kanban Number", "error", err.Error())
					http.Error(w, "Failed to generate Kanban Number", http.StatusInternalServerError)
					return
				}
				kanbansRunningNumber = append(kanbansRunningNumber, kanbanNo)
			}
			kbData := model.KbData{
				KbExtensionID:  ExID,
				CompoundId:     compoundData[0].Id,
				CellNo:         cellNo,
				CreatedBy:      data.UserID,
				CreatedOn:      time.Now(),
				NoOFLots:       data.Quantity,
				MFGDateTime:    time.Now(),
				DemandDateTime: time.Now(),
				ModifiedOn:     time.Now(),
				KanbanNo:       kanbansRunningNumber,
				Note:           data.Note,
			}
			//Creating kb_data records
			ID, err := dao.CreateNewOrUpdateExistingKBData(&kbData)
			if err != nil {
				slog.Error("Failed to create kbData records", "error", err.Error())
				http.Error(w, "Failed to create kbData records", http.StatusInternalServerError)
				return
			}
			for i := 1; i <= data.Quantity; i++ {
				kbRoot := model.KbRoot{
					RunningNo:  0,
					InitialNo:  0,
					CreatedBy:  data.UserID,
					CreatedOn:  time.Now(),
					ModifiedOn: time.Now(),
					KbDataId:   ID,
					Status:     "0",
					KanbanNo:   kanbansRunningNumber[i-1],
				}
				if len(vendorRec.VendorCode) > 0 && vendorRec.VendorCode[0] == 'I' {
					kbRoot.InInventory = true
				}
				_, err = dao.CreateNewOrUpdateExistingKbRoot(&kbRoot)
				if err != nil {
					http.Error(w, "Failed to create part records", http.StatusInternalServerError)
					return
				}
			}
		}
		// }
		utils.SetResponse(w, http.StatusOK, "Successfully created "+strconv.Itoa(data.Quantity)+" parts for vendor "+data.VendorName)
	} else {

		utils.SetResponse(w, http.StatusInternalServerError, "Fail: Decode data")
		slog.Error("%s - error - %s", "Failed to decode data", err.Error())
	}
}

// Add compound to production line
func AddCompoundsInProductionLine(w http.ResponseWriter, r *http.Request) {
	type compoundData struct {
		LineID   string
		KbDataId []string
		UserID   string
	}
	var data compoundData
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&data); err == nil {
		for _, i := range data.KbDataId {
			kbrootid, _ := strconv.Atoi(i)
			ppl, _ := dao.GetProducProcesstLineByParamAndOrder("prod_line_id", data.LineID, "1 ")

			//  Get kb_base_info_id id by production line ID
			LineData, err := dao.GetProdLineByParam("id", data.LineID)
			if len(LineData) == 0 || err != nil {
				slog.Error("Failed to get production line", "error", err.Error())
				http.Error(w, "Failed to get production line", http.StatusInternalServerError)
				return
			}

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
				kbTransaction := model.KbTransaction{
					ProdProcessId:     1,
					Status:            strconv.Itoa(ppl[0].Order),
					KbRootId:          kbrootid,
					ProdProcessLineID: ppl[0].Id,
					StartedOn:         time.Now(),
					CreatedOn:         time.Now(),
					CreatedBy:         data.UserID,
					Operator:          LineData[0].Operator,
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
					if strconv.Itoa(i.ProdLineID) == data.LineID {
						ppl = len(i.Cells)
					}
				}

				kbrootdt, _ := dao.GetKbRootByParam("id", i)
				kbdata, _ := dao.GetKBDataByParam("id", strconv.Itoa(kbrootdt[0].KbDataId))
				// LotNumber, err := GenerateLotNumber(data.LineID)
				if err != nil {
					slog.Error("Failed to create Lot Number", "error", err.Error())
					http.Error(w, "Failed to create Lot Number", http.StatusInternalServerError)
					return
				}
				kbroot := model.KbRoot{
					Id:          kbrootid,
					InInventory: kbrootdt[0].InInventory,
					CreatedOn:   kbrootdt[0].CreatedOn,
					CreatedBy:   kbrootdt[0].CreatedBy,
					Status:      "1",
					KbDataId:    kbdata[0].Id,
					RunningNo:   ppl,
					InitialNo:   ppl,
					ModifiedBy:  data.UserID,
					ModifiedOn:  time.Now(),
					KanbanNo:    kbrootdt[0].KanbanNo,
				}
				_, err = dao.CreateNewOrUpdateExistingKbRoot(&kbroot)
				if err != nil {
					slog.Error("Failed to update kbroot inital and running number records", "error", err.Error())
					http.Error(w, "Failed to update kbroot inital and running number records", http.StatusInternalServerError)
					return
				}

				err = UpdateCustomerOrderStatus(data.KbDataId, "InProductionLine", data.UserID)
				if err != nil {
					slog.Error("Failed to update kbextension status of records", "error", err.Error())
					http.Error(w, "Failed to update kbextension status of records", http.StatusInternalServerError)
					return
				}

				// Create Log
				sysLog := model.SystemLog{
					Message:     "Moved: Cell Number: " + kbdata[0].CellNo + " kanban move to Line Up",
					MessageType: "INFO",
					IsCritical:  false,
					CreatedBy:   data.UserID,
				}
				utils.CreateSystemLogInternal(sysLog)
			} else {
				slog.Error("Failed to create ProdProcessLine records", "error", err.Error())
				http.Error(w, "Failed to create ProdProcessLine records", http.StatusInternalServerError)
				return
			}
		}
	} else {
		utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to update record")
		slog.Error("%s - error - %s", "Failed to Add  Compound data in production line", err.Error())
	}
}

// Get Compounds data based on vendor records
func GetCompoundsByParm(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")

	compounds, _ := dao.GetCompoundDataByParam(key, value)

	if err := json.NewEncoder(w).Encode(compounds); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// Add or Update Compound
func AddorUpdateCompound(w http.ResponseWriter, r *http.Request) {
	var CompoundData model.Compounds
	decoder := json.NewDecoder(r.Body)
	var message string
	var sysLog model.SystemLog
	if err := decoder.Decode(&CompoundData); err == nil {
		if CompoundData.Id == 0 {
			message = "Part added"
			// Create Log
			sysLog = model.SystemLog{
				Message:     "AddPart: Added new part " + CompoundData.CompoundName,
				MessageType: "SUCCESS",
				IsCritical:  false,
				CreatedBy:   CompoundData.CreatedBy,
			}
		} else {
			// Create Log
			sysLog = model.SystemLog{
				Message:     "UpdatePart: Successfully updated part " + CompoundData.CompoundName,
				MessageType: "SUCCESS",
				IsCritical:  false,
				CreatedBy:   CompoundData.CreatedBy,
			}
			message = "Part updated"
		}
		err := dao.CreateNewOrUpdateExistingCompound(&CompoundData)
		if err != nil {
			utils.SetResponse(w, http.StatusInternalServerError, "Fail: Update data")
			slog.Error("%s - error - %s", "Failed to Update data", err.Error())
		}
		// Respond with HTTP 201 Created on success
		utils.CreateSystemLogInternal(sysLog)
		utils.SetResponse(w, http.StatusCreated, message)
	} else {
		utils.SetResponse(w, http.StatusInternalServerError, "Fail: Decode data")
		slog.Error("%s - error - %s", "Failed to decode data", err.Error())
	}
}

func UpdateCompoundStatusToDispatch(w http.ResponseWriter, r *http.Request) {
	type compoundData struct {
		Notes    string
		KbDataId []string //-- It's Kb_root id
		UserID   string
	}
	var data compoundData
	var kbextensiontatus string
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&data); err == nil {
		packingLine, err := dao.GetProdLineByParam("name", "Packing")
		if err != nil {
			slog.Error("Failed to get quality line", "error", err.Error())
		}
		var packingOperator string
		if len(packingLine) != 0 {
			packingOperator = packingLine[0].Operator
		}
		for _, i := range data.KbDataId {
			kbroot, _ := dao.GetKbRootByParam("id", i)
			kbroot[0].Status = "4"
			kbroot[0].DispatchNote = data.Notes
			kbroot[0].DispatchDoneTime = time.Now()
			kbroot[0].PackingOperator = packingOperator
			_, err := dao.CreateNewOrUpdateExistingKbRoot(&kbroot[0])
			if err != nil {
				// Create Logs
				sysLog := model.SystemLog{
					Message:     "Dispatch/Packing: Error while dispatching part",
					MessageType: "ERROR",
					IsCritical:  false,
					CreatedBy:   data.UserID,
				}
				utils.CreateSystemLogInternal(sysLog)
				slog.Error("Failed to update Status", "error", err.Error())
				http.Error(w, "Failed to update Status", http.StatusInternalServerError)
				return
			}
			// Update the packing date and time for respactive kb-root
			// get packing dataStatus from kb_transaction table by
			rootId, _ := strconv.Atoi(i)
			kb_Transaction, _ := dao.GetPackingtTransactionByKBrootID(rootId)
			kb_Transaction.CompletedOn = time.Now()
			kb_Transaction.ModifiedOn = time.Now()
			// just update packing date i.e. completed on date to current date.
			_, err = dao.CreateNewOrUpdateExistingKbTransactionData(&kb_Transaction)
			if err != nil {
				// Create Logs
				sysLog := model.SystemLog{
					Message:     "Dispatch/Packing: Error while dispatching part",
					MessageType: "ERROR",
					IsCritical:  false,
					CreatedBy:   data.UserID,
				}
				utils.CreateSystemLogInternal(sysLog)
				slog.Error("Failed to update packing time", "error", err.Error())
				http.Error(w, "Failed to update packing time", http.StatusInternalServerError)
				return
			}
		}

		for _, i := range data.KbDataId {
			kbroot, _ := dao.GetKbRootByParam("id", i)
			kbdata, _ := dao.GetKBDataByParam("id", strconv.Itoa(kbroot[0].KbDataId))
			kbroots, _ := dao.GetKbRootByParamAndStatus("kb_data_id", strconv.Itoa(kbdata[0].Id), "4")
			if len(kbroots) == kbdata[0].NoOFLots {
				kbextensiontatus = "dispatch"
			} else {
				kbextensiontatus = "dispatched(" + strconv.Itoa(len(kbroots)) + "/" + strconv.Itoa(kbdata[0].NoOFLots) + ")"
			}

			kbext, _ := dao.GetKbExtensionsByParam("id", strconv.Itoa(kbdata[0].KbExtensionID))
			kbext[0].Status = kbextensiontatus
			_, err := dao.CreateNewOrUpdateExistingKbExtension(&kbext[0])
			if err != nil {
				// Create Logs
				sysLog := model.SystemLog{
					Message:     "Dispatch/Packing: Error while dispatching part",
					MessageType: "ERROR",
					IsCritical:  false,
					CreatedBy:   data.UserID,
				}
				utils.CreateSystemLogInternal(sysLog)
				slog.Error("Failed to update Status", "error", err.Error())
				http.Error(w, "Failed to update Status", http.StatusInternalServerError)
				return
			}
		}
		// Create Logs
		sysLog := model.SystemLog{
			Message:     "Dispatch/Packing: Part dispatch successfully",
			MessageType: "INFO",
			IsCritical:  false,
			CreatedBy:   data.UserID,
		}
		utils.CreateSystemLogInternal(sysLog)
		utils.SetResponse(w, http.StatusOK, "Kanban dispatched successfully")
	} else {
		// Create Logs
		sysLog := model.SystemLog{
			Message:     "Dispatch/Packing: Error while dispatching part",
			MessageType: "ERROR",
			IsCritical:  false,
			CreatedBy:   data.UserID,
		}
		utils.CreateSystemLogInternal(sysLog)
		utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to update record")
		slog.Error("%s - error - %s", "Failed to Add  Compound data in production line", err.Error())
	}
}

func UpdateCompoundStatusToPacking(w http.ResponseWriter, r *http.Request) {
	type compoundData struct {
		Notes    string
		KbDataId []string //-- it's an root id
		UserID   string
	}
	var data compoundData
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&data); err == nil {
		qualityLine, err := dao.GetProdLineByParam("name", "Quality")
		if err != nil {
			slog.Error("Failed to get quality line", "error", err.Error())
		}
		var qualityOperator string
		if len(qualityLine) != 0 {
			qualityOperator = qualityLine[0].Operator
		}
		for _, i := range data.KbDataId {
			kbroot, _ := dao.GetKbRootByParam("id", i)
			kbroot[0].Status = "3"
			kbroot[0].QualityNote = data.Notes
			kbroot[0].QualityDoneTime = time.Now()
			kbroot[0].QualityOperator = qualityOperator
			_, err := dao.CreateNewOrUpdateExistingKbRoot(&kbroot[0])
			if err != nil {
				// Create Logs
				sysLog := model.SystemLog{
					Message:     "QualityTest: Error while quality testing ",
					MessageType: "ERROR",
					IsCritical:  false,
					CreatedBy:   data.UserID,
				}
				utils.CreateSystemLogInternal(sysLog)
				slog.Error("Failed to update Status", "error", err.Error())
				http.Error(w, "Failed to update Status", http.StatusInternalServerError)
				return
			}
			// Get Vendor ID and compound ID by KbRootID
			kbRootId, _ := strconv.Atoi(i)
			vendorID, compoundID, err := dao.GetVendorAndCompoundByKRID(kbRootId)
			if err != nil {
				// Create Logs
				sysLog := model.SystemLog{
					Message:     "QualityTest: Error while quality testing ",
					MessageType: "ERROR",
					IsCritical:  false,
					CreatedBy:   data.UserID,
				}
				utils.CreateSystemLogInternal(sysLog)
				http.Error(w, "Can't get Vendor ID and part ID", http.StatusInternalServerError)
				log.Println("Error fetching Vendor/Part data:", err)
				return
			}

			vendorrec, err := dao.GetVendorByParam("id", strconv.Itoa(vendorID))
			if err != nil {
				http.Error(w, "Can't get Vendor", http.StatusInternalServerError)
				log.Println("Error fetching Vendor", err)
				return
			}
			if string(vendorrec[0].VendorCode[0]) == "I" {
				// update cold store avaliable quantity
				err := dao.UpdateColdStoreAvailableQuantity(compoundID)
				if err != nil {
					// Create Logs
					sysLog := model.SystemLog{
						Message:     "QualityTest: Error while quality testing ",
						MessageType: "ERROR",
						IsCritical:  false,
						CreatedBy:   data.UserID,
					}
					utils.CreateSystemLogInternal(sysLog)
					http.Error(w, "Can't update cold store", http.StatusInternalServerError)
					log.Println("Error updating cold store:", err)
					return
				}
			}
		}
		// Create Logs
		sysLog := model.SystemLog{
			Message:     "QualityTest: Part successfully passed quality testing ",
			MessageType: "INFO",
			IsCritical:  false,
			CreatedBy:   data.UserID,
		}
		utils.CreateSystemLogInternal(sysLog)
		utils.SetResponse(w, http.StatusOK, "Kanban Approved successfully")
	} else {
		// Create Logs
		sysLog := model.SystemLog{
			Message:     "QualityTest: Error while quality testing ",
			MessageType: "ERROR",
			IsCritical:  false,
			CreatedBy:   data.UserID,
		}
		utils.CreateSystemLogInternal(sysLog)
		utils.SetResponse(w, http.StatusInternalServerError, "Fail to update record")
		slog.Error("%s - error - %s", "Failed to Add  Compound data in production line", err.Error())
	}
}

func UpdateCompoundQualityStatusToReject(w http.ResponseWriter, r *http.Request) {
	type compoundData struct {
		Notes    string
		KbDataId []string
		UserID   string
	}
	var data compoundData
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&data); err == nil {
		qualityLine, err := dao.GetProdLineByParam("name", "Quality")
		if err != nil {
			slog.Error("Failed to get quality line", "error", err.Error())
		}
		var qualityOperator string
		if len(qualityLine) != 0 {
			qualityOperator = qualityLine[0].Operator
		}
		for _, i := range data.KbDataId {
			kbroot, _ := dao.GetKbRootByParam("id", i)
			kbroot[0].Status = "-1"
			kbroot[0].QualityNote = data.Notes
			kbroot[0].InInventory = false
			kbroot[0].QualityDoneTime = time.Now()
			kbroot[0].QualityOperator = qualityOperator
			_, err := dao.CreateNewOrUpdateExistingKbRoot(&kbroot[0])
			if err != nil {
				// Create Log
				sysLog := model.SystemLog{
					Message:     "QualityTest: Error while quality testing",
					MessageType: "INFO",
					IsCritical:  false,
					CreatedBy:   data.UserID,
				}
				utils.CreateSystemLogInternal(sysLog)
				slog.Error("Failed to update Status", "error", err.Error())
				http.Error(w, "Failed to update Status", http.StatusInternalServerError)
				return
			}
		}
		// Create Log
		sysLog := model.SystemLog{
			Message:     "QualityTest: Part failed quality testing",
			MessageType: "INFO",
			IsCritical:  false,
			CreatedBy:   data.UserID,
		}
		utils.CreateSystemLogInternal(sysLog)
		utils.SetResponse(w, http.StatusOK, "Kanban Rejected successfully")
	} else {
		// Create Log
		sysLog := model.SystemLog{
			Message:     "QualityTest: Error while quality testing",
			MessageType: "INFO",
			IsCritical:  false,
			CreatedBy:   data.UserID,
		}
		utils.CreateSystemLogInternal(sysLog)
		utils.SetResponse(w, http.StatusInternalServerError, "Fail to update record")
		slog.Error("%s - error - %s", "Failed to Add  Compound data in production line", err.Error())
	}
}

// Get All Active Compounds
func GetAllActiveCompounds(w http.ResponseWriter, r *http.Request) {

	compounds, err := dao.GetCompoundDataByParam("status", true)
	if err != nil {
		slog.Error("Recordes not found", "error", err.Error())
		http.Error(w, "Failed to encode fetch records", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(compounds); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func GetAllOrderByVendorCode(w http.ResponseWriter, r *http.Request) {
	var data model.Vendors
	var OrderDetails []model.OrderEntry
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&data); err == nil {
		vendorData, err := dao.GetVendorByParam("vendor_code", data.VendorCode)
		if err != nil || len(vendorData) == 0 {
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to get vendor record")
			return
		}
		kbExtensionData, _ := dao.GetKbExtensionsByParam("vendor_id", strconv.Itoa(vendorData[0].ID))
		for _, extensionData := range kbExtensionData {
			kbDataData, _ := dao.GetKBDataByParam("kb_extension_id", strconv.Itoa(extensionData.Id))
			if len(kbDataData) == 0 {
				continue
			}
			compoundData, _ := dao.GetCompoundDataByParam("id", kbDataData[0].CompoundId)
			var order model.OrderEntry
			order.CompoundCode = compoundData[0].CompoundName
			order.DemandDateTime = kbDataData[0].DemandDateTime
			order.NoOFLots = kbDataData[0].NoOFLots
			order.UserID = extensionData.CreatedBy
			order.Status = extensionData.Status
			order.Location = kbDataData[0].Location
			order.CellNo = kbDataData[0].CellNo
			order.MFGDateTime = kbDataData[0].MFGDateTime
			OrderDetails = append(OrderDetails, order)
		}
		// Send the response as JSON
		if err := json.NewEncoder(w).Encode(OrderDetails); err != nil {
			slog.Error("%s - error - %s", "Failed to add order data in KbData table", err.Error())
			utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to update record")
			return
		}
	} else {
		utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to update record")
		slog.Error("%s - error - %s", "Failed to add order data in KbData table", err.Error())
	}
}

// GetAllOperator returns a all records present in operator table
func GetAllCompoundsBySearchAndPagination(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		Pagination model.PaginationReq `json:"pagination"`
		Conditions []string            `json:"conditions"`
	}
	// Decode JSON request
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	comp, PaginationResp, err := dao.GetAllCompoundsBySearchAndPagination(requestData.Pagination, requestData.Conditions)
	if err != nil {
		slog.Error("Recordes not found", "error", err.Error())
		http.Error(w, "Failed to encode fetch records", http.StatusInternalServerError)
		return
	}

	var Response struct {
		Pagination model.PaginationResp
		Data       []*model.Compounds
	}
	Response.Pagination = PaginationResp
	Response.Data = comp

	if err := json.NewEncoder(w).Encode(Response); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
