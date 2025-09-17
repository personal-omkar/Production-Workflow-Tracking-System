package server

import (
	"encoding/json"
	"log"
	"log/slog"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	m "irpl.com/kanban-commons/model"
	"irpl.com/kanban-commons/utils"
	"irpl.com/kanban-dao/dao"
)

// Add order entry
func CreateNewOrderEntry(w http.ResponseWriter, r *http.Request) {

	var data m.OrderEntry
	var isValidStatusUpdate m.IsValidStatusUpdate

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err == nil {
		// get vendor by used id ,check it status
		userToVendor, _ := dao.GetUserToVendorByParam("user_id", data.UserID)
		vendor, _ := dao.GetVendorByParam("id", strconv.Itoa(userToVendor[0].VendorId))
		if vendor[0].Isactive {
			if data.CellNo != "" {
				kbdata, _ := dao.GetKBDataByParam("cell_no", data.CellNo)
				isValidStatusUpdate.KbDataId = kbdata[0].Id
				isValidStatusUpdate.Status = data.Status
				KanbanNosLen := len(kbdata[0].KanbanNo)
				kanbansRunningNumber := make([]string, 0, data.NoOFLots)
				if KanbanNosLen != data.NoOFLots {
					// Append extra Kanban Nos.
					if data.NoOFLots > KanbanNosLen {
						for range data.NoOFLots - KanbanNosLen {
							kanban, err := dao.GenerateKanbanNumber()
							if err != nil {
								slog.Error("Failed to generate Kanban Number", "error", err.Error())
								http.Error(w, "Failed to generate Kanban Number", http.StatusInternalServerError)
								return
							}
							kanbansRunningNumber = append(kanbansRunningNumber, kanban)
						}
					}
					// Remove extra Kanban Nos.
					if data.NoOFLots < KanbanNosLen {
						for range data.NoOFLots {
							newLength := data.NoOFLots
							if newLength < len(kbdata[0].KanbanNo) {
								kbdata[0].KanbanNo = kbdata[0].KanbanNo[:newLength]
							}
						}
					}
				}
				if IsValidStatusUpdate(isValidStatusUpdate) {
					kbdata[0].DemandDateTime = data.DemandDateTime
					kbdata[0].Location = data.Location
					kbdata[0].NoOFLots = data.NoOFLots
					kbdata[0].CompoundId, _ = strconv.Atoi(data.CompoundCode)
					kbdata[0].ModifiedOn = time.Now()
					kbdata[0].Note = data.CustomerNote
					kbdata[0].SubmitDateTime = time.Now()
					if KanbanNosLen != data.NoOFLots {
						kbdata[0].KanbanNo = append(kbdata[0].KanbanNo, kanbansRunningNumber...)
					}

					_, err = dao.CreateNewOrUpdateExistingKBData(&kbdata[0])
					if err != nil {
						slog.Error("Failed to Update KbData record", "error", err.Error())
						http.Error(w, "Failed to Update KbData record", http.StatusInternalServerError)
						return
					}

					kbext, _ := dao.GetKbExtensionsByParam("id", strconv.Itoa(kbdata[0].KbExtensionID))
					kbext[0].Status = data.Status
					kbext[0].ModifiedOn = time.Now()

					_, err = dao.CreateNewOrUpdateExistingKbExtension(&kbext[0])
					if err != nil {
						slog.Error("Failed to Update KbExtension record", "error", err.Error())
						http.Error(w, "Failed to Update KbExtension record", http.StatusInternalServerError)
						return
					}
					if data.Status == "creating" {
						sysLog := m.SystemLog{
							Message:     "Updated: Order for Cell No :" + kbdata[0].CellNo + " has been updated successfully",
							MessageType: "SUCCESS",
							IsCritical:  false,
							CreatedBy:   kbdata[0].CreatedBy,
						}
						utils.CreateSystemLogInternal(sysLog)
						utils.SetResponse(w, http.StatusOK, "Order updated successfully")
					} else {
						sysLog := m.SystemLog{
							Message:     "Submitted: Order for Cell No :" + kbdata[0].CellNo + " consisting of  " + strconv.Itoa(kbdata[0].NoOFLots) + " Lots has been submitted  successfully",
							MessageType: "SUCCESS",
							IsCritical:  false,
							CreatedBy:   kbdata[0].CreatedBy,
						}
						utils.CreateSystemLogInternal(sysLog)
						utils.SetResponse(w, http.StatusOK, "Order submitted successfully for CellNo: "+kbdata[0].CellNo+";")
					}
				}
			} else {
				compoundID, err := strconv.Atoi(data.CompoundCode)
				// compoundData, err := dao.GetCompoundDataByParam("id", data.CompoundCode)
				if err != nil {
					slog.Error("Failed to fetch compound record to create Kb data ", "error", err.Error())
					http.Error(w, "Failed to fetch compound record to create Kb data", http.StatusInternalServerError)
					return
				}
				// get Vendor ID by user ID
				vendorData, err := dao.GetUserToVendorByParam("user_id", data.UserID)
				if err != nil {
					slog.Error("Failed to fetch Vendor ID", "error", err.Error())
					http.Error(w, "Failed to fetch Vendor ID", http.StatusInternalServerError)
					return
				}
				KbExtension := m.KbExtension{
					Status:    data.Status,
					VendorID:  vendorData[0].VendorId,
					CreatedBy: data.UserID,
					CreatedOn: time.Now(),
				}
				ID, err := dao.CreateNewOrUpdateExistingKbExtension(&KbExtension)
				if err != nil {
					slog.Error("Failed to create KbExtension record", "error", err.Error())
					http.Error(w, "Failed to create KbExtension record", http.StatusInternalServerError)
					return
				}
				vendorId, _ := dao.GetUserToVendorByParam("user_id", data.UserID)
				vendorCodeDetails, _ := dao.GetVendorByParam("id", strconv.Itoa(vendorId[0].VendorId))
				// creating cell based on the + vendor code + Id
				cellNo := vendorCodeDetails[0].VendorCode + "/" + strconv.Itoa(ID)
				kanbansRunningNumber := make([]string, 0, data.NoOFLots)
				for range data.NoOFLots {
					kanban, err := dao.GenerateKanbanNumber()
					if err != nil {
						slog.Error("Failed to generate Kanban Number", "error", err.Error())
						http.Error(w, "Failed to generate Kanban Number", http.StatusInternalServerError)
						return
					}
					kanbansRunningNumber = append(kanbansRunningNumber, kanban)
				}
				kbgroup := m.KbData{
					CompoundId:     compoundID,
					DemandDateTime: data.DemandDateTime,
					MFGDateTime:    data.MFGDateTime,
					CellNo:         cellNo,
					NoOFLots:       data.NoOFLots,
					KbExtensionID:  ID,
					Location:       data.Location,
					CreatedBy:      data.UserID,
					CreatedOn:      time.Now(),
					KanbanNo:       kanbansRunningNumber,
					Note:           data.CustomerNote,
				}
				_, err = dao.CreateNewOrUpdateExistingKBData(&kbgroup)
				if err != nil {
					slog.Error("Failed to create KbData record", "error", err.Error())
					http.Error(w, "Failed to create KbData record", http.StatusInternalServerError)
					return
				} else {
					sysLog := m.SystemLog{
						Message:     "Saved: Order for Cell No :" + kbgroup.CellNo + " consisting of  " + strconv.Itoa(kbgroup.NoOFLots) + " Lots has been saved successfully",
						MessageType: "SUCCESS",
						IsCritical:  false,
						CreatedBy:   kbgroup.CreatedBy,
					}
					utils.CreateSystemLogInternal(sysLog)
				}
				utils.SetResponse(w, http.StatusOK, "Order created successfully CellNo:"+kbgroup.CellNo+";")
			}
		} else {
			utils.SetResponse(w, http.StatusInternalServerError, "You can't place order, vendor is disabled")
		}

	}
}

// Add order entry
func CreateMultiNewOrderEntry(w http.ResponseWriter, r *http.Request) {
	var data []m.OrderEntry
	var isValidStatusUpdate m.IsValidStatusUpdate
	var msg string
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err == nil {
		for _, data := range data {
			if data.CellNo != "" {
				kbdata, _ := dao.GetKBDataByParam("cell_no", data.CellNo)
				isValidStatusUpdate.KbDataId = kbdata[0].Id
				isValidStatusUpdate.Status = data.Status

				if IsValidStatusUpdate(isValidStatusUpdate) {
					kbdata[0].SubmitDateTime = time.Now()
					kbdata[0].DemandDateTime = data.DemandDateTime
					kbdata[0].Location = data.Location
					kbdata[0].NoOFLots = data.NoOFLots
					kbdata[0].CompoundId, _ = strconv.Atoi(data.CompoundCode)
					kbdata[0].ModifiedOn = time.Now()

					_, err = dao.CreateNewOrUpdateExistingKBData(&kbdata[0])
					if err != nil {
						slog.Error("Failed to Update KbData record", "error", err.Error())
						http.Error(w, "Failed to Update KbData record", http.StatusInternalServerError)
						return
					}

					kbext, _ := dao.GetKbExtensionsByParam("id", strconv.Itoa(kbdata[0].KbExtensionID))
					kbext[0].Status = data.Status
					kbext[0].ModifiedOn = time.Now()

					_, err = dao.CreateNewOrUpdateExistingKbExtension(&kbext[0])
					if err != nil {
						slog.Error("Failed to Update KbExtension record", "error", err.Error())
						http.Error(w, "Failed to Update KbExtension record", http.StatusInternalServerError)
						return
					}
					if data.Status == "creating" {
						sysLog := m.SystemLog{
							Message:     "Updated: Order for Cell No :" + kbdata[0].CellNo + " has been updated successfully",
							MessageType: "SUCCESS",
							IsCritical:  false,
							CreatedBy:   kbdata[0].CreatedBy,
						}
						utils.CreateSystemLogInternal(sysLog)
						// utils.SetResponse(w, http.StatusOK, "Order updated successfully")
					} else {
						sysLog := m.SystemLog{
							Message:     "Submitted: Order for Cell No :" + kbdata[0].CellNo + " consisting of  " + strconv.Itoa(kbdata[0].NoOFLots) + " Lots has been submitted  successfully",
							MessageType: "SUCCESS",
							IsCritical:  false,
							CreatedBy:   kbdata[0].CreatedBy,
						}
						utils.CreateSystemLogInternal(sysLog)
					}
				}
			}
			msg = msg + "Successfully submitted order for CellNo:" + data.CellNo + ";"
		}
		utils.SetResponse(w, http.StatusOK, msg)
	}
}

func UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var status m.Status
	var isValidStatusUpdate m.IsValidStatusUpdate
	err := json.NewDecoder(r.Body).Decode(&status)
	if err != nil {
		http.Error(w, "Invalid input format", http.StatusBadRequest)
		return
	}
	id, _ := strconv.Atoi(status.ID)
	isValidStatusUpdate.KbDataId = id
	isValidStatusUpdate.Status = status.Status
	if IsValidStatusUpdate(isValidStatusUpdate) {

		if (status.Status == "dispatch" || status.Status == "reject") && status.Kanban == 0 {
			if status.Status == "dispatch" {
				kbData, err := dao.GetKBDataByParam("id", status.ID)
				if err != nil {
					http.Error(w, "Fail to get order details", http.StatusBadRequest)
					return
				}
				kbroot, err := dao.GetinInventoryKbRootByCompoundID(strconv.Itoa(status.CompoundID), status.Dispatch)
				if err != nil {
					http.Error(w, "Fail to get inventory kanban", http.StatusBadRequest)
					return
				}
				for i := range kbroot {
					kbroot[i].InInventory = false
					kbroot[i].Comment = strconv.Itoa(kbroot[i].KbDataId)
					kbroot[i].KbDataId = id
					kbroot[i].KanbanNo = kbroot[i].KanbanNo + " | " + kbData[0].KanbanNo[i]
					dao.CreateNewOrUpdateExistingKbRoot(&kbroot[i])
					dao.UpdateAvailablePartByCompoundId(status.CompoundID)
				}
			}
			if status.Status == "reject" {
				err = dao.UpdateOrderStatus(id, status.Dispatch, status.CompoundID, status.Status, status.UserID)
			} else {
				err = dao.UpdateOrderStatus(id, status.Dispatch, status.CompoundID, utils.JoinStr("dispatched(0/", strconv.Itoa(status.Dispatch), ")"), status.UserID)
			}
			if err != nil {
				http.Error(w, "Fail to update data", http.StatusBadRequest)
				return
			} else {
				var status string
				kbdata, _ := dao.GetKBDataByParam("id", strconv.Itoa(id))
				kbext, _ := dao.GetKbExtensionsByParam("id", strconv.Itoa(kbdata[0].KbExtensionID))
				if kbext[0].Status == "reject" {
					status = utils.IF(kbext[0].Status == "reject", "Rejected")
				} else {
					status = utils.IF(strings.HasPrefix(kbext[0].Status, "dispatch"), "Dispatched")
				}
				sysLog := m.SystemLog{
					Message:     status + ": Order for Cell No :" + kbdata[0].CellNo + " consisting of  " + strconv.Itoa(kbdata[0].NoOFLots) + " Lots has been " + kbext[0].Status,
					MessageType: "SUCCESS",
					IsCritical:  false,
					CreatedBy:   kbdata[0].CreatedBy,
				}
				utils.CreateSystemLogInternal(sysLog)
			}
		} else {
			if status.NoOFLots == status.Dispatch {
				kbData, err := dao.GetKBDataByParam("id", status.ID)
				if err != nil {
					http.Error(w, "Fail to get order details", http.StatusBadRequest)
					return
				}
				kbroot, err := dao.GetinInventoryKbRootByCompoundID(strconv.Itoa(status.CompoundID), status.Dispatch)
				if err != nil {
					http.Error(w, "Fail to get inventory kanban", http.StatusBadRequest)
					return
				}
				for i := range kbroot {
					kbroot[i].InInventory = false
					kbroot[i].Comment = strconv.Itoa(kbroot[i].KbDataId)
					kbroot[i].KbDataId = id
					kbroot[i].KanbanNo = kbroot[i].KanbanNo + " | " + kbData[0].KanbanNo[i]
					dao.CreateNewOrUpdateExistingKbRoot(&kbroot[i])
					dao.UpdateAvailablePartByCompoundId(status.CompoundID)
				}

				err = dao.UpdateOrderStatus(id, status.Dispatch, status.CompoundID, utils.JoinStr("dispatched(0/", strconv.Itoa(status.Dispatch), ")"), status.UserID)
				if err != nil {
					http.Error(w, "Fail to update data", http.StatusBadRequest)
					return
				} else {
					var status string
					kbdata, _ := dao.GetKBDataByParam("id", strconv.Itoa(id))
					kbext, _ := dao.GetKbExtensionsByParam("id", strconv.Itoa(kbdata[0].KbExtensionID))
					status = utils.IF(strings.HasPrefix(kbext[0].Status, "dispatch"), "Dispatched")
					sysLog := m.SystemLog{
						Message:     status + ": Order for Cell No :" + kbdata[0].CellNo + " consisting of  " + strconv.Itoa(kbdata[0].NoOFLots) + " Lots has been " + kbext[0].Status,
						MessageType: "SUCCESS",
						IsCritical:  false,
						CreatedBy:   kbdata[0].CreatedBy,
					}
					utils.CreateSystemLogInternal(sysLog)
				}
			}
			if math.Abs(float64(status.NoOFLots-status.Dispatch)) > 0 {
				kbData, err := dao.GetKBDataByParam("id", status.ID)
				if err != nil {
					http.Error(w, "Fail to get order details", http.StatusBadRequest)
					return
				}
				// Order Kanban
				kbroot, err := dao.GetinInventoryKbRootByCompoundID(strconv.Itoa(status.CompoundID), status.Dispatch)
				if err != nil {
					http.Error(w, "Fail to get inventory kanban", http.StatusBadRequest)
					return
				}
				for i := range kbroot {
					kbroot[i].Comment = strconv.Itoa(kbroot[i].KbDataId)
					kbroot[i].InInventory = false
					kbroot[i].KbDataId = id
					kbroot[i].KanbanNo = kbroot[i].KanbanNo + " | " + kbData[0].KanbanNo[i]
					dao.CreateNewOrUpdateExistingKbRoot(&kbroot[i])
					// dao.UpdateAvailablePartByCompoundId(status.CompoundID)
				}
				err = dao.UpdateOrderStatusAndCreateKanban(id, status.Status, status.UserID, int(math.Abs(float64(status.NoOFLots-status.Dispatch))), status.Dispatch, status.CompoundID, false)
				if err != nil {
					http.Error(w, "Fail to update data", http.StatusBadRequest)
					return
				} else {
					var status string
					kbdata, _ := dao.GetKBDataByParam("id", strconv.Itoa(id))
					kbext, _ := dao.GetKbExtensionsByParam("id", strconv.Itoa(kbdata[0].KbExtensionID))
					status = utils.IF(kbext[0].Status == "approved", "Approved")
					sysLog := m.SystemLog{
						Message:     status + ": Order for Cell No :" + kbdata[0].CellNo + " consisting of  " + strconv.Itoa(kbdata[0].NoOFLots) + " Lots has been " + kbext[0].Status,
						MessageType: "SUCCESS",
						IsCritical:  false,
						CreatedBy:   kbdata[0].CreatedBy,
					}
					utils.CreateSystemLogInternal(sysLog)
				}
			}
			ColdStorageKamban := math.Abs(float64(status.NoOFLots - (status.Kanban + status.Dispatch)))
			if ColdStorageKamban > 0 {
				vendor, _ := dao.GetVendorByParamStartsWith("vendor_code", "I")
				// Code Storage Kanban
				// Create new KBextension with vendor id - 3(it's for inventory), return Kb Extension ID
				var KBextension m.KbExtension
				KBextension.OrderID = 0
				KBextension.Status = "approved"
				KBextension.VendorID = vendor[0].ID
				KBextension.CreatedBy = status.UserID
				Id, err := dao.CreateNewOrUpdateExistingKbExtension(&KBextension)
				if err != nil {
					http.Error(w, "Fail to update data", http.StatusBadRequest)
					return
				}

				// Create New  KBdata with compound KBextension ID and Vendor id received By creating new Extension.\
				var KBdata m.KbData
				KBdata.CompoundId = status.CompoundID
				KBdata.NoOFLots = int(ColdStorageKamban)
				KBdata.KbExtensionID = Id
				KBdata.DemandDateTime = time.Now()
				KBdata.MFGDateTime = time.Now()
				kanbansRunningNumber := make([]string, 0, KBdata.NoOFLots)
				for range KBdata.NoOFLots {
					kanban, err := dao.GenerateKanbanNumber()
					if err != nil {
						slog.Error("Failed to generate Kanban Number", "error", err.Error())
						http.Error(w, "Failed to generate Kanban Number", http.StatusInternalServerError)
						return
					}
					kanbansRunningNumber = append(kanbansRunningNumber, kanban)
				}
				KBdata.KanbanNo = kanbansRunningNumber
				ID, err := dao.CreateNewOrUpdateExistingKBData(&KBdata)
				if err != nil {
					http.Error(w, "Fail to update data", http.StatusBadRequest)
					return
				}
				KBdata.Id = ID

				// Get Vendor Code By vendor ID and Compound Name by compound id
				vendor, err = dao.GetVendorByParam("id", strconv.Itoa(KBextension.VendorID))
				if err != nil {
					http.Error(w, "Fail to update data", http.StatusBadRequest)
					return
				}
				compound, err := dao.GetCompoundDataByParam("id", strconv.Itoa(KBdata.CompoundId))
				if err != nil {
					http.Error(w, "Fail to update data", http.StatusBadRequest)
					return
				}

				KBdata.CreatedBy = status.UserID
				KBdata.CellNo = strconv.Itoa(ID) + vendor[0].VendorCode + compound[0].CompoundName
				KBdata.CellNo = vendor[0].VendorCode + "/" + strconv.Itoa(KBdata.KbExtensionID)
				_, err = dao.CreateNewOrUpdateExistingKBData(&KBdata)
				if err != nil {
					http.Error(w, "Fail to update data", http.StatusBadRequest)
					return
				}

				err = dao.UpdateOrderStatusAndCreateKanban(KBdata.Id, "approved", status.UserID, int(ColdStorageKamban), 0, status.CompoundID, true)
				if err != nil {
					http.Error(w, "Fail to update data", http.StatusBadRequest)
					return
				}
			}
		}
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Invalide Operation", http.StatusBadRequest)
		return
	}

}

// Delete order entry
func DeleteOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var kbdata m.KbData
	err := json.NewDecoder(r.Body).Decode(&kbdata)
	if err != nil {
		http.Error(w, "Invalid input format", http.StatusBadRequest)
		return
	}
	//get ke_id
	KDdata, err := dao.GetKBDataByParam("id", strconv.Itoa(kbdata.Id))
	if err != nil {
		http.Error(w, "Order not found", http.StatusBadRequest)
		return
	}

	//delete kb_data by id
	err = dao.DeleteKBDataByParam("id", strconv.Itoa(kbdata.Id))
	if err != nil {
		http.Error(w, "Fail to delete order", http.StatusBadRequest)
		return
	}

	//delete ke_data by id
	err = dao.DeleteKbExtensionsByParam("id", strconv.Itoa(KDdata[0].KbExtensionID))
	if err != nil {
		http.Error(w, "Fail to delete order", http.StatusBadRequest)
		return
	}

	utils.SetResponse(w, http.StatusOK, "Order deleted successfully")
}

func DailyAndMonthlyVendorLimit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var orderData []m.KbData
	var todayLots int
	var monthLots int
	var hourLots int

	err := json.NewDecoder(r.Body).Decode(&orderData)
	if err != nil {
		http.Error(w, "Invalid input format", http.StatusBadRequest)
		return
	}

	userVendorData, err := dao.GetUserToVendorByParam("user_id", strconv.Itoa(orderData[0].Id))
	if err != nil || len(userVendorData) == 0 {
		http.Error(w, "Vendor not found for the user", http.StatusNotFound)
		return
	}

	vendorData, err := dao.GetVendorByParam("id", strconv.Itoa(userVendorData[0].VendorId))
	if err != nil || len(vendorData) == 0 {
		http.Error(w, "Vendor configuration not found", http.StatusNotFound)
		return
	}

	kbExtensionData, err := dao.GetKbExtensionsByParam("vendor_id", strconv.Itoa(vendorData[0].ID))
	if err != nil {
		http.Error(w, "Failed to retrieve KB extensions for the vendor", http.StatusInternalServerError)
		return
	}

	var kbData []m.KbData
	for _, extension := range kbExtensionData {
		if extension.Status == "creating" {
			continue
		}
		data, err := dao.GetKBDataByParam("kb_extension_id", strconv.Itoa(extension.Id))
		if err == nil && len(data) > 0 {
			kbData = append(kbData, data[0])
		}
	}

	var filteredkbData []m.KbData

	for _, v := range orderData {

		todayLots = todayLots + v.NoOFLots
		monthLots = monthLots + v.NoOFLots
		hourLots = hourLots + v.NoOFLots

		for _, kbgroup := range kbData {
			if v.CellNo != kbgroup.CellNo && !IsCellExist(filteredkbData, kbgroup.CellNo) {
				filteredkbData = append(filteredkbData, kbgroup)
			}
		}

	}

	currentDate := orderData[0].DemandDateTime.Format("2006-01-02")
	currentMonth := orderData[0].DemandDateTime.Format("2006-01")

	endHourTime := time.Now()

	startHourTime := endHourTime.Add(-1 * time.Hour)

	for _, record := range filteredkbData {

		if record.DemandDateTime.Format("2006-01-02") == currentDate {
			todayLots += record.NoOFLots
		}

		if record.DemandDateTime.Format("2006-01") == currentMonth {
			monthLots += record.NoOFLots
		}

		if !record.SubmitDateTime.Before(startHourTime) && record.SubmitDateTime.Before(endHourTime) {
			hourLots += record.NoOFLots
		}

	}

	response := map[string]interface{}{
		"vendor_id":     vendorData[0].ID,
		"daily_limit":   vendorData[0].PerDayLotConfig,
		"monthly_limit": vendorData[0].PerMonthLotConfig,
		"hourly_limit":  vendorData[0].PerHourLotConfig,
	}

	// Check if the limits are exceeded
	if hourLots > vendorData[0].PerHourLotConfig {
		response["message"] = "Hourly lot limit exceeded!"
		response["exceed_by"] = math.Abs(float64(vendorData[0].PerHourLotConfig - hourLots))
		w.WriteHeader(http.StatusInternalServerError)
	} else if todayLots > vendorData[0].PerDayLotConfig {
		response["message"] = "Daily lot limit exceeded!"
		response["exceed_by"] = math.Abs(float64(vendorData[0].PerDayLotConfig - todayLots))
		w.WriteHeader(http.StatusInternalServerError)
	} else if monthLots > vendorData[0].PerMonthLotConfig {
		response["message"] = "Monthly lot limit exceeded!"
		response["exceed_by"] = math.Abs(float64(vendorData[0].PerMonthLotConfig - monthLots))
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		response["message"] = "Lots are within the limit"
		w.WriteHeader(http.StatusOK)
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func DailyAndMonthlyVendorLimitByVendorCode(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var orderData m.OrderDetails
	err := json.NewDecoder(r.Body).Decode(&orderData)
	if err != nil {
		http.Error(w, "Invalid input format", http.StatusBadRequest)
		return
	}
	vendorData, err := dao.GetVendorByParam("vendor_code", orderData.VendorCode)
	if err != nil || len(vendorData) == 0 {
		http.Error(w, "Vendor configuration not found", http.StatusNotFound)
		return
	}

	kbExtensionData, err := dao.GetKbExtensionsByParam("vendor_id", strconv.Itoa(vendorData[0].ID))
	if err != nil {
		http.Error(w, "Failed to retrieve KB extensions for the vendor", http.StatusInternalServerError)
		return
	}

	var kbData []m.KbData
	for _, extension := range kbExtensionData {
		if extension.Status == "reject" {
			continue
		}
		data, err := dao.GetKBDataByParam("kb_extension_id", strconv.Itoa(extension.Id))
		if err == nil && len(data) > 0 {
			kbData = append(kbData, data[0])
		}
	}

	currentDate := orderData.DemandDateTime.Format("2006-01-02")
	currentMonth := orderData.DemandDateTime.Format("2006-01")
	currentHour := orderData.DemandDateTime.Format("2006-01-02 15")

	todayLots := orderData.NoOFLots
	monthLots := orderData.NoOFLots
	hourLots := orderData.NoOFLots
	for _, record := range kbData {
		if record.DemandDateTime.Format("2006-01-02") == currentDate && record.CellNo != orderData.CellNo {
			todayLots += record.NoOFLots
		}
		if record.DemandDateTime.Format("2006-01") == currentMonth && record.CellNo != orderData.CellNo {
			monthLots += record.NoOFLots
		}
		if record.DemandDateTime.Format("2006-01-02 15") == currentHour {
			hourLots += record.NoOFLots
		}
	}

	response := map[string]interface{}{
		"vendor_id":     vendorData[0].ID,
		"daily_limit":   vendorData[0].PerDayLotConfig,
		"monthly_limit": vendorData[0].PerMonthLotConfig,
		"hourly_limit":  vendorData[0].PerHourLotConfig,
	}

	// Check if the limits are exceeded
	if hourLots > vendorData[0].PerHourLotConfig {
		response["message"] = "Hourly lot limit exceeded!"
		response["exceed_by"] = math.Abs(float64(vendorData[0].PerHourLotConfig - hourLots))
		w.WriteHeader(http.StatusInternalServerError)
	} else if todayLots > vendorData[0].PerDayLotConfig {
		response["message"] = "Daily lot limit exceeded!"
		response["exceed_by"] = math.Abs(float64(vendorData[0].PerDayLotConfig - todayLots))
		w.WriteHeader(http.StatusInternalServerError)
	} else if monthLots > vendorData[0].PerMonthLotConfig {
		response["message"] = "Monthly lot limit exceeded!"
		response["exceed_by"] = math.Abs(float64(vendorData[0].PerMonthLotConfig - monthLots))
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		response["message"] = "Lots are within the limit"
		w.WriteHeader(http.StatusOK)
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to encode response")
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func DailyVendorLimit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var orderData m.KbData
	err := json.NewDecoder(r.Body).Decode(&orderData)
	if err != nil {
		http.Error(w, "Invalid input format", http.StatusBadRequest)
		return
	}
	userVendorData, err := dao.GetUserToVendorByParam("user_id", strconv.Itoa(orderData.Id))
	if err != nil || len(userVendorData) == 0 {
		http.Error(w, "Vendor not found for the user", http.StatusNotFound)
		return
	}

	vendorData, err := dao.GetVendorByParam("id", strconv.Itoa(userVendorData[0].VendorId))
	if err != nil || len(vendorData) == 0 {
		http.Error(w, "Vendor configuration not found", http.StatusNotFound)
		return
	}

	kbExtensionData, err := dao.GetKbExtensionsByParam("vendor_id", strconv.Itoa(vendorData[0].ID))
	if err != nil {
		http.Error(w, "Failed to retrieve KB extensions for the vendor", http.StatusInternalServerError)
		return
	}

	var kbData []m.KbData
	for _, extension := range kbExtensionData {
		if extension.Status == "creating" {
			continue
		}
		data, err := dao.GetKBDataByParam("kb_extension_id", strconv.Itoa(extension.Id))
		if err == nil && len(data) > 0 {
			kbData = append(kbData, data[0])
		}
	}

	currentDate := orderData.DemandDateTime.Format("2006-01-02")
	// currentMonth := orderData.DemandDateTime.Format("2006-01")
	// currentHour := orderData.DemandDateTime.Format("2006-01-02 15")

	todayLots := orderData.NoOFLots
	// monthLots := orderData.NoOFLots
	// hourLots := orderData.NoOFLots
	for _, record := range kbData {
		if record.DemandDateTime.Format("2006-01-02") == currentDate && record.CellNo != orderData.CellNo {
			todayLots += record.NoOFLots
		}
	}

	response := map[string]interface{}{
		"vendor_id":   vendorData[0].ID,
		"daily_limit": vendorData[0].PerDayLotConfig,
	}

	// Check if the limits are exceeded
	if todayLots > vendorData[0].PerDayLotConfig {
		response["message"] = "Daily lot limit exceeded!"
		response["exceed_by"] = math.Abs(float64(vendorData[0].PerDayLotConfig - todayLots))
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		response["message"] = "Lots are within the limit"
		w.WriteHeader(http.StatusOK)
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func IsCellExist(kbdata []m.KbData, cellNo string) bool {
	for _, g := range kbdata {
		if g.CellNo == cellNo {
			return true
		}
	}
	return false
}
