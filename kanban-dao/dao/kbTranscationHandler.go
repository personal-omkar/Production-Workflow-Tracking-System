package dao

import (
	"fmt"
	"time"

	m "irpl.com/kanban-commons/model"
	db "irpl.com/kanban-dao/db"
)

const (
	KBTRANSACTIONDATA_TABLE string = "kb_transaction" // Updated to match the table name
)

// CreateNewOrUpdateExistingKbTransactionData creates a new kb_data  or updates an existing kb_data
func CreateNewOrUpdateExistingKbTransactionData(kbt *m.KbTransaction) (ID int, err error) {

	now := time.Now()
	if kbt.Id != 0 {
		kbt.ModifiedOn = now

		if err := db.GetDB().Table(KBTRANSACTIONDATA_TABLE).Save(&kbt).Error; err != nil {
			return 0, err
		}
	} else {
		kbt.CreatedOn = now

		if err := db.GetDB().Table(KBTRANSACTIONDATA_TABLE).Create(&kbt).Error; err != nil {
			return 0, err
		}
	}
	return kbt.Id, nil
}

// DeleteKBTransactionsByRootID deletes all kb_transaction records where kb_root_id matches the given ID.
func DeleteKBTransactionsByRootID(kbRootID string) error {
	if err := db.GetDB().Table(KBTRANSACTIONDATA_TABLE).Where("kb_root_id = ?", kbRootID).Delete(&m.KbTransaction{}).Error; err != nil {
		return err
	}
	return nil
}

// GetKbTransactionByGroup retrieves all ProdProcessLine records by GroupName and ProdLineId.
func GetProdProcessLineByGroup(GroupName string, ProdLineId int) ([]m.ProdProcessLine, error) {
	var prodProcessLines []m.ProdProcessLine

	// Prepare the query
	query := `SELECT * 
              FROM prod_process_line 
              WHERE group_name = ? AND prod_line_id = ?;`

	// Execute the query and fetch the results into the slice
	if err := db.GetDB().Raw(query, GroupName, ProdLineId).Scan(&prodProcessLines).Error; err != nil {
		return nil, fmt.Errorf("could not execute query: %v", err)
	}

	return prodProcessLines, nil
}

// GetLatestTransactionByKBrootID fetches the latest kb_transaction record by joining with prod_process_line table and returns the data with the highest order.
func GetLatestTransactionByKBrootID(KBrootID int) (KBTransactionData m.KbTransaction, err error) {

	// Query to get the latest transaction based on the created_on column
	err = db.GetDB().
		Table(KBTRANSACTIONDATA_TABLE+" as kbt").
		Select("kbt.*, ppl.order as process_order").
		Joins("JOIN prod_process_line as ppl ON kbt.prod_process_line_id = ppl.id").
		Where("kbt.kb_root_id = ?", KBrootID).
		Order("kbt.created_on DESC").
		First(&KBTransactionData).
		Error

	if err != nil {
		return KBTransactionData, err
	}

	return KBTransactionData, nil
}

func GetPackingtTransactionByKBrootID(KBrootID int) (KBTransactionData m.KbTransaction, err error) {
	err = db.GetDB().
		Table(KBTRANSACTIONDATA_TABLE+" as kbt").
		Select("kbt.*").
		Joins("JOIN prod_process_line as ppl ON kbt.prod_process_line_id = ppl.id").
		Where("kbt.kb_root_id = ? AND ppl.prod_process_id = ?", KBrootID, 2).
		First(&KBTransactionData).
		Error

	if err != nil {
		return KBTransactionData, err
	}

	return KBTransactionData, nil
}

func CheckIfTransactionExists(KbRoot, prodProcessLineId int) (bool, error) {
	var exists bool
	query := `
	SELECT EXISTS (
			SELECT 1
			FROM kb_transaction
			WHERE kb_root_id = ? AND prod_process_line_id = ?
		);
	`
	err := db.GetDB().Raw(query, KbRoot, prodProcessLineId).Scan(&exists).Error
	if err != nil {
		return false, fmt.Errorf("could not check if transaction exists: %v", err)
	}
	return exists, nil
}

func UpdateCompletedOnForGroupedIDs(groupedIDs []string, kbRootID int, completedOn time.Time) error {
	err := db.GetDB().
		Table(KBTRANSACTIONDATA_TABLE).
		Where("prod_process_line_id IN (?) AND kb_root_id = ?", groupedIDs, kbRootID).
		Update("completed_on", completedOn).
		Error
	return err
}

// GetFirstAndLastTransactionWithProdLine retrieves the first and last transactions for a given kb_root_id, along with the prod_line name.
func GetFirstAndLastTransactionWithProdLine(kbRootID string) (firstTxn, lastTxn m.KbTransaction, prodLineName string, err error) {
	dbConn := db.GetDB()

	// Fetch first transaction (earliest created_at)
	err = dbConn.Table(KBTRANSACTIONDATA_TABLE).
		Where("kb_root_id = ?", kbRootID).
		Order("created_on ASC").
		First(&firstTxn).Error
	if err != nil {
		return
	}

	// Fetch last transaction (latest created_at)
	err = dbConn.Table(KBTRANSACTIONDATA_TABLE).
		Where("kb_root_id = ?", kbRootID).
		Order("created_on DESC").
		First(&lastTxn).Error
	if err != nil {
		return
	}

	// Fetch prod_line name by joining with prod_line table
	err = dbConn.Table(PRODUCTLINE_TABLE).
		Select("name").
		Joins("JOIN prod_process_line ON prod_process_line.prod_line_id = prod_line.id").
		Where("prod_process_line.id = ?", firstTxn.ProdProcessLineID).
		Pluck("name", &prodLineName).Error
	if err != nil {
		return
	}

	return firstTxn, lastTxn, prodLineName, nil
}

func StatusKbRootExists(status string, kbRootID int) (bool, error) {
	var count int64
	err := db.GetDB().
		Table(KBTRANSACTIONDATA_TABLE).
		Where("status = ? AND kb_root_id = ?", status, kbRootID).
		Count(&count).
		Error

	if err != nil {
		return false, err
	}

	return count > 0, nil // Returns true if record exists, false otherwise
}

// DeleteKbTransactionByProdLine deletes kb_transaction records for the given production line ID.
func DeleteKbTransactionByProdLine(prodLineID int) error {
	return db.GetDB().
		Table(KBTRANSACTIONDATA_TABLE).
		Where("prod_process_line_id IN (SELECT id FROM prod_process_line WHERE prod_line_id = ?)", prodLineID).
		Delete(&m.KbTransaction{}).Error
}

// DeleteKbTransactionByProdLine deletes kb_transaction records for the given production line ID.
func DeleteProductionLineDataInKbTransactionByProdLine(prodLineID int) error {
	return db.GetDB().
		Table(KBTRANSACTIONDATA_TABLE).
		Where("kb_transaction.kb_root_id in(select kb_root_id FROM kb_transaction WHERE prod_process_line_id IN (SELECT id FROM prod_process_line WHERE prod_line_id = ?)) ", prodLineID).
		Delete(&m.KbTransaction{}).Error
}

func DeleteKbTransactionByParam(key, value string) (err error) {
	query := db.GetDB().Table(KBTRANSACTIONDATA_TABLE)
	query = query.Where(key + " = " + value).Delete(&m.KbTransaction{})
	return query.Error
}

// GetKbTransactionByParam returns a kb_transition records based on parameter
func GetKbTransactionByParam(key, value string) (kbt []m.KbTransaction, err error) {
	query := db.GetDB().Table(KBTRANSACTIONDATA_TABLE)
	query = query.Where(key + " = " + value).Find(&kbt)
	return kbt, query.Error
}

// GetStatusFromKbTransaction retrieves the status from the kb_transaction_data table using the provided kbRootID.
func GetStatusFromKbTransaction(kbRootID int) (status string, err error) {
	// Execute the query to fetch only the status field
	err = db.GetDB().
		Table(KBTRANSACTIONDATA_TABLE).
		Select("status").
		Where("kb_root_id = ?", kbRootID).
		Scan(&status). // Use Scan to directly map the single field
		Error

	// Return the status and error
	return status, err
}
