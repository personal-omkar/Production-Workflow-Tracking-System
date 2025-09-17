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

func CreateNewKbTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var transactionData m.KbTransaction
	// we are getting the root_id and prodProcessID for transaction tobe perform
	err := json.NewDecoder(r.Body).Decode(&transactionData)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		log.Println("Error decoding JSON:", err)
		return
	}
	// If this status is modified, make sure to update it in the frontend JavaScript as well.
	if transactionData.Status == "Packing" {
		kbtransactionData, err := dao.GetKbTransactionByParam("kb_root_id", strconv.Itoa(transactionData.KbRootId))
		if err != nil {
			http.Error(w, "Fail to get process ID for packing", http.StatusInternalServerError)
			slog.Error("ERROR: Fail to get process ID for packing", "error", err)
			return
		}

		// get Production_process_line id
		productionProcessLine, err := dao.GetProdProcessLineByParam("id", strconv.Itoa(kbtransactionData[0].ProdProcessLineID))
		if err != nil {
			http.Error(w, "Fail to get process ID for packing", http.StatusInternalServerError)
			slog.Error("ERROR: Fail to get process ID for packing", "error", err)
			return
		}

		// get Production_process_line id
		productionProcesses, err := dao.GetProdProcessLineByParam("prod_line_id", strconv.Itoa(productionProcessLine[0].ProdLineId))
		if err != nil {
			http.Error(w, "Fail to get process ID for packing", http.StatusInternalServerError)
			slog.Error("ERROR: Fail to get process ID for packing", "error", err)
			return
		}
		// Get prod_processes_line_id for packing stage
		for _, data := range productionProcesses {
			if data.ProdProcessID == 2 { //ProdProcesses 2 is by default set for packing
				transactionData.ProdProcessLineID = data.Id
				break
			}
		}
	}
	transactionData.StartedOn = time.Now()
	transactionData.CreatedOn = time.Now()

	TransactionExists, err := dao.CheckIfTransactionExists(transactionData.KbRootId, transactionData.ProdProcessLineID)
	if err != nil {
		http.Error(w, "Fail to check transaction", http.StatusInternalServerError)
		log.Println("Error while checking transaction:", err)
		return
	}
	if TransactionExists {
		log.Println("Transaction already exists for the process")
		http.Error(w, "Transaction already exists for the process", http.StatusInternalServerError)
		return
	}

	// Get last transaction of root we have got in transactionData
	LatestData, err := dao.GetLatestTransactionByKBrootID(transactionData.KbRootId)
	if err != nil {
		http.Error(w, "Fail to get current state", http.StatusInternalServerError)
		log.Println("Error while getting current state:", err)
		return
	}

	// Fetch data of ProdProcessLine using LatestData.ProdProcessLineID
	prodProcessLine, err := dao.GetProducProcesstLineByParam("id", strconv.Itoa(LatestData.ProdProcessLineID))
	if err != nil {
		http.Error(w, "Fail to fetch process line data", http.StatusInternalServerError)
		log.Println("Error while fetching process line data:", err)
		return
	}

	if !prodProcessLine[0].IsGroup {
		// If IsGroup is false, update completed_on date for the single entry
		LatestData.CompletedOn = time.Now()
		_, err = dao.CreateNewOrUpdateExistingKbTransactionData(&LatestData)
		if err != nil {
			log.Println("Error : Fail to update latest data : ", err)
			http.Error(w, "Fail to update latest data", http.StatusInternalServerError)
			return
		}
		kbroot, err := dao.GetKbRootByParam("id", strconv.Itoa(LatestData.KbRootId))
		if err != nil {
			log.Println("Error : Fail to get kanban data : ", err)
			http.Error(w, "Fail to get kanban data", http.StatusInternalServerError)
			return
		}
		kbdata, err := dao.GetKBDataByParam("id", strconv.Itoa(kbroot[0].KbDataId))
		if err != nil {
			log.Println("Error : Fail to get order data : ", err)
			http.Error(w, "Fail to get order data", http.StatusInternalServerError)
			return
		}
		prodprocessline, _ := dao.GetProducProcesstLineByParam("id", strconv.Itoa(LatestData.ProdProcessLineID))
		productionprocess, err := dao.GetProdProcessByParam("id", strconv.Itoa(prodprocessline[0].ProdProcessID))
		if err != nil {
			log.Println("Fail to get process for latest data:", err)
			http.Error(w, "Fail to get process for latest data", http.StatusInternalServerError)
			return
		} else {
			sysLog := m.SystemLog{
				Message:     utils.JoinStr(`Completed:  Cell Number: `, kbdata[0].CellNo, ` , Lot No: `, kbroot[0].LotNo, ` kanban completed `, productionprocess[0].Name),
				MessageType: "INFO",
				IsCritical:  false,
				CreatedBy:   transactionData.CreatedBy,
			}
			utils.CreateSystemLogInternal(sysLog)
		}

	} else {
		var processNames string
		// If IsGroup is true, get all grouped IDs where group is common
		groupedIDs, err := dao.GetGroupedProcessLineIDs(prodProcessLine[0].GroupName)
		if err != nil {
			http.Error(w, "Fail to fetch grouped IDs", http.StatusInternalServerError)
			log.Println("Error while fetching grouped process line IDs:", err)
			return
		}
		kbroot, err := dao.GetKbRootByParam("id", strconv.Itoa(LatestData.KbRootId))
		if err != nil {
			log.Println("Error : Fail to get kanban data : ", err)
			http.Error(w, "Fail to get kanban data", http.StatusInternalServerError)
			return
		}
		kbdata, err := dao.GetKBDataByParam("id", strconv.Itoa(kbroot[0].KbDataId))
		if err != nil {
			log.Println("Error : Fail to get order data : ", err)
			http.Error(w, "Fail to get order data", http.StatusInternalServerError)
			return
		}
		for _, v := range groupedIDs {
			prodprocessline, _ := dao.GetProducProcesstLineByParam("id", v)
			productionprocess, _ := dao.GetProdProcessByParam("id", strconv.Itoa(prodprocessline[0].ProdProcessID))
			processNames = processNames + " " + productionprocess[0].Name

		}

		// Update completed_on date for all grouped IDs
		err = dao.UpdateCompletedOnForGroupedIDs(groupedIDs, transactionData.KbRootId, time.Now())
		if err != nil {
			http.Error(w, "Fail to update grouped stages", http.StatusInternalServerError)
			log.Println("Error updating grouped transactions:", err)
			return
		} else {
			sysLog := m.SystemLog{
				Message:     utils.JoinStr(`Completed:  Cell Number: `, kbdata[0].CellNo, ` , Lot No: `, kbroot[0].LotNo, ` kanban for  `, prodProcessLine[0].GroupName, `  has been completed,  Consisting of`, processNames),
				MessageType: "INFO",
				IsCritical:  false,
				CreatedBy:   transactionData.CreatedBy,
			}
			utils.CreateSystemLogInternal(sysLog)
		}

	}

	// Check if transaction Exists
	Exists, _ := dao.CheckIfTransactionExists(transactionData.KbRootId, transactionData.ProdProcessLineID)
	if !Exists {
		var KbTransaction m.KbTransaction
		KbTransaction.Id = 0
		KbTransaction.KbRootId = transactionData.KbRootId
		KbTransaction.ProdProcessLineID = transactionData.ProdProcessLineID
		KbTransaction.Status = transactionData.Status
		KbTransaction.CreatedBy = transactionData.CreatedBy

		err := IfTransactionNOTExistsDoThis(KbTransaction)
		if err != nil {
			log.Println("Fail to create transaction")
			http.Error(w, "Fail to create transaction", http.StatusInternalServerError)
			return
		}
	} else if Exists {
		log.Println("Transaction already exists for the process")
		http.Error(w, "Transaction already exists for the process", http.StatusInternalServerError)
		return
	}

	err = dao.UpdateRunningAndInitialNumberToZero(transactionData.KbRootId)
	if err != nil {
		http.Error(w, "Fail to update Root Data", http.StatusInternalServerError)
		log.Println("Error updating root data:", err)
		return
	}

	if transactionData.Status == "Packing" {
		// // Get Vendor ID and compound ID by KbRootID
		// vendorID, compoundID, err := dao.GetVendorAndCompoundByKRID(transactionData.KbRootId)
		// if err != nil {
		// 	http.Error(w, "Can't get Vendor ID and compound ID", http.StatusInternalServerError)
		// 	log.Println("Error fetching Vendor/Compound data:", err)
		// 	return
		// }

		// if vendorID == 3 {
		// 	// update cold store avaliable quantity
		// 	err := dao.UpdateColdStoreAvailableQuantity(compoundID)
		// 	if err != nil {
		// 		http.Error(w, "Can't update cold store", http.StatusInternalServerError)
		// 		log.Println("Error updating cold store:", err)
		// 		return
		// 	}
		// }
		err = dao.UpdateKbRootStatus(transactionData.KbRootId, "2", transactionData.CreatedBy)
		if err != nil {
			http.Error(w, "Can't update root status", http.StatusInternalServerError)
			log.Println("Error updating root status:", err)
			return
		}
		err := UpdateCustomerOrderStatus([]string{strconv.Itoa(transactionData.KbRootId)}, "quality", transactionData.CreatedBy)
		if err != nil {
			http.Error(w, "Can't update customer order status", http.StatusInternalServerError)
			log.Println("Error updating customer order status:", err)
			return
		}
	} else if transactionData.Status == "Line-Up" {
		err := UpdateCustomerOrderStatus([]string{strconv.Itoa(transactionData.KbRootId)}, "InProductionProcess", transactionData.CreatedBy)
		if err != nil {
			http.Error(w, "Can't update customer order status", http.StatusInternalServerError)
			log.Println("Error updating customer order status:", err)
			return
		}
	}
	w.WriteHeader(http.StatusCreated)

}

// This function is used when we are sending Line-up process to production-Line Dynamically
func UpdateRunningNumberAfterTransactioin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var ProductionLine m.ProdLine
	err := json.NewDecoder(r.Body).Decode(&ProductionLine)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		log.Println("Error decoding JSON:", err)
		return
	}

	KbRootIds, err := dao.GetUniqueTransactionKBRoot(ProductionLine.Id)
	if err != nil {
		http.Error(w, "Fail to get Unique Transaction", http.StatusInternalServerError)
		return
	}

	if KbRootIds != nil {
		err = dao.UpdateRunningNumber(KbRootIds)
		if err != nil {
			http.Error(w, "Fail to update running numbers", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

// If transaction not exists
func IfTransactionNOTExistsDoThis(transactionData m.KbTransaction) error {

	var KbTransaction m.KbTransaction
	KbTransaction.Id = 0
	KbTransaction.KbRootId = transactionData.KbRootId
	KbTransaction.ProdProcessLineID = transactionData.ProdProcessLineID
	KbTransaction.CreatedBy = transactionData.CreatedBy

	//Get Prod Process Line Data
	ProdProcessLineData, err := dao.GetProducProcesstLineByParam("id", strconv.Itoa(transactionData.ProdProcessLineID))
	if err != nil {
		log.Println("Fail to Get Prod Process Line:", err)
		return err
	}
	// Check is group
	if ProdProcessLineData[0].IsGroup {
		prodProcessLineID, err := dao.GetProdProcessLineByGroup(ProdProcessLineData[0].GroupName, ProdProcessLineData[0].ProdLineId)
		if err != nil {
			log.Println("Fail to get production process Line")
			return err
		}
		for _, value := range prodProcessLineID {

			// get prodProcessLine order by its id
			ProdProcessLineData, err := dao.GetProdProcessLineByParam("id", strconv.Itoa(value.Id))
			if err != nil {
				log.Println("Fail to Get Prod Process Line data:", err)
				return err
			}

			KbTransaction.Id = 0
			KbTransaction.KbRootId = transactionData.KbRootId
			KbTransaction.ProdProcessLineID = value.Id
			KbTransaction.StartedOn = time.Now()
			KbTransaction.CreatedBy = transactionData.CreatedBy
			KbTransaction.Status = strconv.Itoa(ProdProcessLineData[0].Order)

			kbrrot, _ := dao.GetKbRootByParam("id", strconv.Itoa(transactionData.KbRootId))

			if kbrrot[0].LotNo == "" {
				kbt, _ := dao.GetKbTransactionByParam("kb_root_id", strconv.Itoa(transactionData.KbRootId))
				if len(kbt) == 1 {
					kbrrot[0].LotNo, _ = GenerateLotNumber(strconv.Itoa(value.ProdLineId))
					_, err = dao.CreateNewOrUpdateExistingKbRoot(&kbrrot[0])
					if err != nil {
						log.Println("Fail to create transaction")
						return err
					}
					prodline, _ := dao.GetProdLineByParam("id", strconv.Itoa(value.ProdLineId))
					prodline[0].RunningNumber = prodline[0].RunningNumber + 1
					err = dao.CreateNewOrUpdateExistingProductLine(&prodline[0])
					if err != nil {
						log.Println("Fail to update running number for production line")
						return err
					}
				}
			}
			_, err = dao.CreateNewOrUpdateExistingKbTransactionData(&KbTransaction)
			if err != nil {
				log.Println("Fail to create transaction")
				return err
			}
			err = UpdateCustomerOrderStatus([]string{strconv.Itoa(transactionData.KbRootId)}, "InProductionProcess", transactionData.CreatedBy)
			if err != nil {
				log.Println("Error updating customer order status:", err)
				return err
			}
		}
		kbrrot, _ := dao.GetKbRootByParam("id", strconv.Itoa(transactionData.KbRootId))
		kbdata, _ := dao.GetKBDataByParam("id", strconv.Itoa(kbrrot[0].KbDataId))
		sysLog := m.SystemLog{
			Message:     utils.JoinStr(`Moved:  Cell Number: `, kbdata[0].CellNo, ` , Lot No: `, kbrrot[0].LotNo, ` kanban move to   `, ProdProcessLineData[0].GroupName, ` `),
			MessageType: "INFO",
			IsCritical:  false,
			CreatedBy:   transactionData.CreatedBy,
		}
		utils.CreateSystemLogInternal(sysLog)
	} else {
		// get prodProcessLine order by its id
		ProdProcessLineData, err := dao.GetProdProcessLineByParam("id", strconv.Itoa(KbTransaction.ProdProcessLineID))
		if err != nil {
			log.Println("Fail to Get Prod Process Line data:", err)
			return err
		}

		KbTransaction.StartedOn = time.Now()
		KbTransaction.CreatedBy = transactionData.CreatedBy
		KbTransaction.Status = strconv.Itoa(ProdProcessLineData[0].Order)

		kbrrot, _ := dao.GetKbRootByParam("id", strconv.Itoa(transactionData.KbRootId))
		if kbrrot[0].LotNo == "" {
			kbt, _ := dao.GetKbTransactionByParam("kb_root_id", strconv.Itoa(transactionData.KbRootId))
			if len(kbt) == 1 {
				kbrrot[0].LotNo, _ = GenerateLotNumber(strconv.Itoa(ProdProcessLineData[0].ProdLineId))
				_, err = dao.CreateNewOrUpdateExistingKbRoot(&kbrrot[0])
				if err != nil {
					log.Println("Fail to create transaction")
					return err
				}
				prodline, _ := dao.GetProdLineByParam("id", strconv.Itoa(ProdProcessLineData[0].ProdLineId))
				prodline[0].RunningNumber = prodline[0].RunningNumber + 1
				err = dao.CreateNewOrUpdateExistingProductLine(&prodline[0])
				if err != nil {
					log.Println("Fail to update running number for production line")
					return err
				}
			}
		}

		_, err = dao.CreateNewOrUpdateExistingKbTransactionData(&KbTransaction)
		if err != nil {
			log.Println("Fail to create transaction")
			return err
		}
		err = UpdateCustomerOrderStatus([]string{strconv.Itoa(transactionData.KbRootId)}, "InProductionProcess", transactionData.CreatedBy)
		if err != nil {
			log.Println("Error updating customer order status:", err)
			return err
		} else {
			kbrrot, _ := dao.GetKbRootByParam("id", strconv.Itoa(transactionData.KbRootId))
			kbdata, _ := dao.GetKBDataByParam("id", strconv.Itoa(kbrrot[0].KbDataId))
			productionprocess, _ := dao.GetProdProcessByParam("id", strconv.Itoa(ProdProcessLineData[0].ProdProcessID))
			sysLog := m.SystemLog{
				Message:     utils.JoinStr(`Moved:  Cell Number: `, kbdata[0].CellNo, ` , Lot No: `, kbrrot[0].LotNo, ` kanban move to   `, productionprocess[0].Name, ` `),
				MessageType: "INFO",
				IsCritical:  false,
				CreatedBy:   transactionData.CreatedBy,
			}
			utils.CreateSystemLogInternal(sysLog)

		}

	}
	return nil
}
