package dao

import (
	"fmt"
	"strconv"
	"time"

	m "irpl.com/kanban-commons/model"
	"irpl.com/kanban-commons/utils"
	db "irpl.com/kanban-dao/db"
)

const (
	KBExtension_TABLE string = "kb_extension" // Updated to match the table name
)

// CreateNewOrUpdateExistingKbExtension creates a new kbextesnion or updates an existing kbextension
func CreateNewOrUpdateExistingKbExtension(kbe *m.KbExtension) (Id int, err error) {
	var extensionID int
	now := time.Now()
	if kbe.Id != 0 {
		kbe.ModifiedOn = now

		if err := db.GetDB().Table(KBExtension_TABLE).Save(&kbe).Error; err != nil {
			return 0, err
		}
	} else {
		kbe.CreatedOn = now

		if err := db.GetDB().Table(KBExtension_TABLE).Create(&kbe).Error; err != nil {
			return 0, err
		}
	}
	extensionID = kbe.Id
	return extensionID, err
}

// GetAllKbExtensions returns a all records present in kb_extension table
func GetAllKbExtensions() (kbe []m.KbExtension, err error) {
	result := db.GetDB().Table(KBExtension_TABLE).Find(&kbe)
	return kbe, result.Error
}

// GetKbExtensionsByParam returns a kb_extension records based on parameter
func GetKbExtensionsByParam(key, value string) (kbe []m.KbExtension, err error) {
	query := db.GetDB().Table(KBExtension_TABLE)
	query = query.Where(key + " = " + value).Find(&kbe)
	return kbe, query.Error
}

// GetKbExtensionsBystatusandVendor returns a kb_extension records based on vendor id and there status
func GetKbExtensionsBystatusandVendor(status, vendor string) (kbe []m.KbExtension, err error) {
	query := db.GetDB().Table(KBExtension_TABLE)
	query = query.Where("status = ? AND vendor_id=?", status, vendor).Order("created_on ASC").Find(&kbe)
	return kbe, query.Error
}

// todo - make both update DB in sync if one fail cancle whole transaction and revert the updates
// Update only order stuats in kb_extension table
func UpdateOrderStatus(id, Dispatch, compoundID int, status, modifiedby string) error {
	// Define the SQL query
	modifiedOn := time.Now()
	query := `
		UPDATE ` + KBExtension_TABLE + ` ke
		SET status = ? , modified_by = ? , modified_on = ?
		FROM kb_data kd
		WHERE kd.kb_extension_id = ke.id
	 	AND kd.id = ?;
	`

	result := db.GetDB().Exec(query, status, modifiedby, modifiedOn, id)
	if result.Error != nil {
		return fmt.Errorf("failed to update order status: %w", result.Error)
	}

	// Check how many rows were affected
	if result.RowsAffected == 0 {
		return fmt.Errorf("no rows updated for id: %v", id)
	}
	// if status == "dispatch" {
	// 	err := UpdateAvailableQuantity(compoundID, Dispatch)
	// 	if err != nil {
	// 		return fmt.Errorf("failed to update order status: %w", result.Error)
	// 	}
	// }
	return nil
}

// Update order status for customer order
// This function works after Order is approved and flow through production line.
func UpdateCustomerOrderStatus(ID int, status, modified_by string) error {
	// Check Status of the order using kb_extension_id
	// Define the SQL query
	modifiedOn := time.Now()
	query := `
		UPDATE ` + KBExtension_TABLE + ` ke
		SET status = ? , modified_by = ? , modified_on = ?
		WHERE ke.id = ?;

	`
	result := db.GetDB().Exec(query, status, modified_by, modifiedOn, ID)
	if result.Error != nil {
		return fmt.Errorf("failed to update order status: %w", result.Error)
	}
	// Check how many rows were affected
	if result.RowsAffected == 0 {
		return fmt.Errorf("no rows updated for id: %v", ID)
	}
	return nil
}

func UpdateOrderStatusAndCreateKanban(id int, status, createdby string, kanban, dispatch, compoundID int, isInventory bool) error {
	modifiedOn := time.Now()

	query := `
        UPDATE ` + KBExtension_TABLE + ` ke
        SET status = ? , modified_by = ? , modified_on = ?
        FROM kb_data kd
        WHERE kd.kb_extension_id = ke.id
        AND kd.id = ?;
    `

	result := db.GetDB().Exec(query, status, createdby, modifiedOn, id)
	if result.Error != nil {
		return fmt.Errorf("failed to update order status: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("no rows updated for id: %v", id)
	}
	if dispatch > 0 {
		err := UpdateAvailableQuantity(compoundID, dispatch)
		if err != nil {
			return fmt.Errorf("failed to update order status: %w", result.Error)
		}
	}
	kbdata, _ := GetKBDataByParam("id", strconv.Itoa(id))
	for i := 0; i < kanban; i++ {
		kbroot := m.KbRoot{
			Status:    "0",
			RunningNo: -1,
			InitialNo: -1,
			CreatedBy: createdby,
			CreatedOn: time.Now(),
			KbDataId:  id,
		}
		if !isInventory {
			kbroot.KanbanNo = kbdata[0].KanbanNo[dispatch+i]
		}

		if isInventory {
			kbroot.KanbanNo = kbdata[0].KanbanNo[i]
			kbroot.InInventory = true
		}
		if _, err := CreateNewOrUpdateExistingKbRoot(&kbroot); err != nil {
			sysLog := m.SystemLog{
				Message:     "CreateKaban: Failed to created kanban for cell number : " + kbdata[0].CellNo,
				MessageType: "ERROR",
				IsCritical:  false,
				CreatedBy:   createdby,
			}
			utils.CreateSystemLogInternal(sysLog)
			return fmt.Errorf("failed to create or update KbRoot: %w", err)
		}
	}

	sysLog := m.SystemLog{
		Message:     "CreateKaban: Successfully created " + strconv.Itoa(kanban) + " kanban for cell number : " + kbdata[0].CellNo,
		MessageType: "SUCCESS",
		IsCritical:  false,
		CreatedBy:   createdby,
	}
	utils.CreateSystemLogInternal(sysLog)
	return nil
}

// DeleteKbExtensionsByParam deletes kb_extension records based on a given parameter
func DeleteKbExtensionsByParam(key, value string) error {
	query := db.GetDB().Table(KBExtension_TABLE)

	// Use placeholders to avoid SQL injection
	query = query.Where(key+" = ?", value).Delete(nil)

	// Check if the query resulted in an error
	if query.Error != nil {
		return fmt.Errorf("failed to delete KB extension: %w", query.Error)
	}

	// Check if any rows were actually deleted
	if query.RowsAffected == 0 {
		return fmt.Errorf("no records found for %s = %s", key, value)
	}

	return nil
}

// DeleteKbExtensionByProdLine deletes kb_data records for the given production line ID.
func DeleteKbExtensionByProdLine(prodLineID int) error {
	return db.GetDB().
		Table(KBExtension_TABLE).
		Where("kb_extension.kb_root_id in(select kb_root_id FROM kb_transaction WHERE prod_process_line_id IN (SELECT id FROM prod_process_line WHERE prod_line_id = ?)) ", prodLineID).
		Delete(&m.KbTransaction{}).Error
}

// DeleteKbExtensionByProdLine deletes kb_data records for the given production line ID.
func GetKbExtensiondataByProdLine(prodLineID int) (kbe []m.KbExtension, err error) {
	query := db.GetDB().Table(KBExtension_TABLE)
	query = query.Where("kb_extension.kb_root_id in(select kb_root_id FROM kb_transaction WHERE prod_process_line_id IN (SELECT id FROM prod_process_line WHERE prod_line_id = ?)) ", prodLineID).Find(&kbe)
	return kbe, query.Error
}

func GetOrderIdAndStatusByKDid(kbextensionid int) (status, orderID string, err error) {
	type KBExtension struct {
		Status  string
		OrderID string
	}

	var kbe KBExtension

	// Build the query
	err = db.GetDB().
		Table(KBExtension_TABLE).
		Select("status, order_id").
		Where("id = ? ", kbextensionid).
		First(&kbe).
		Error

	if err != nil {
		return "", "", err
	}

	return kbe.Status, kbe.OrderID, nil
}
