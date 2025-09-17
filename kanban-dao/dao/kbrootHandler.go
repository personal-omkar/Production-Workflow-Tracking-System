package dao

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
	m "irpl.com/kanban-commons/model"
	db "irpl.com/kanban-dao/db"
)

const (
	KbRoot_TABLE string = "kb_root" // Updated to match the table name
)

// CreateNewOrUpdateExistingKbExtension creates a new kbdata or updates an existing kbdata
func CreateNewOrUpdateExistingKbRoot(kbr *m.KbRoot) (Id int, err error) {
	var rootID int
	now := time.Now()
	if kbr.Id != 0 {
		kbr.ModifiedOn = now

		if err := db.GetDB().Table(KbRoot_TABLE).Save(&kbr).Error; err != nil {
			return 0, err
		}
	} else {
		kbr.CreatedOn = now

		if err := db.GetDB().Table(KbRoot_TABLE).Create(&kbr).Error; err != nil {
			return 0, err
		}
	}
	rootID = kbr.Id
	return rootID, err
}

// GetAllKbRoot returns a all records present in kb_data table
func GetAllKbRoot() (kbr []m.KbRoot, err error) {
	result := db.GetDB().Table(KbRoot_TABLE).Find(&kbr)
	return kbr, result.Error
}

// GetKbRootByParam returns a kb_root records based on parameter
func GetKbRootByParam(key, value string) (kbr []m.KbRoot, err error) {
	query := db.GetDB().Table(KbRoot_TABLE)
	query = query.Where(key + " = " + value).Order("modified_on ASC").Find(&kbr)
	return kbr, query.Error
}

// GetKbRootByParam returns kb_root records based on a key and multiple values
func GetMultiKbRootByParam(key string, values []string) (kbr []m.KbRoot, err error) {
	query := db.GetDB().Table(KbRoot_TABLE)
	query = query.Where(key+" IN (?)", values).Order("modified_on ASC").Find(&kbr)
	return kbr, query.Error
}

// GetKbRootByParam returns a kb_root records based on parameter
func GetKbRootByParamAndStatus(key, value, status string) (kbr []m.KbRoot, err error) {
	query := db.GetDB().Table(KbRoot_TABLE)
	query = query.Where("status='" + status + "' AND " + key + " = " + value).Find(&kbr)
	return kbr, query.Error
}

// GetKbRootInitalNo returns a kb_root initial no based on existing records
func GetKbRootInitalNo() (kbr m.KbRoot, err error) {
	query := db.GetDB().Table(KbRoot_TABLE).Order("id DESC").Limit(1).Find(&kbr)
	return kbr, query.Error
}

// update running number in Table
func UpdateRunningNo(KbRoots []m.KbRoot) error {
	now := time.Now()
	for _, value := range KbRoots {
		if err := db.GetDB().Table(KbRoot_TABLE).Where("id = ?", value.Id).Updates(map[string]interface{}{
			"modified_on": now,
			"running_no":  value.RunningNo,
		}).Error; err != nil {
			return err
		}
	}
	return nil
}

// UpdateRunningNumber decrements the running number by 1 for each kb_root_id in the given slice.
func UpdateRunningNumber(kbRootIDs []int) error {
	// Validate input
	if len(kbRootIDs) == 0 {
		return fmt.Errorf("no kb_root_id provided")
	}

	// Start a transaction for batch updates
	tx := db.GetDB().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Iterate over the kb_root_ids and update their running_no
	for _, kbRootID := range kbRootIDs {
		var currentRunningNo int

		err := tx.Table(KbRoot_TABLE).Select("running_no").Where("id = ?", kbRootID).Scan(&currentRunningNo).Error
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to fetch running_no for kb_root_id %d: %v", kbRootID, err)
		}
		if currentRunningNo <= 0 {
			continue
		}
		err = tx.Table(KbRoot_TABLE).Where("id = ?", kbRootID).UpdateColumn("running_no", gorm.Expr("running_no - 1")).Error
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to update running_no for kb_root_id %d: %v", kbRootID, err)
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

// Define a struct to hold the results
type VendorAndCompound struct {
	VendorID   int `json:"vendor_id"`
	CompoundID int `json:"compound_id"`
}

// GetVendorAndCompoundByKRID fetches the vendorID and compoundID for the given KRID
func GetVendorAndCompoundByKRID(krID int) (int, int, error) {
	var result VendorAndCompound

	// Perform the query to fetch the vendor and compound IDs based on krID
	if err := db.GetDB().Table(KbRoot_TABLE+" kr").
		Joins("JOIN "+KBData_TABLE+" kd ON kr.kb_data_id = kd.id").
		Joins("JOIN "+KBExtension_TABLE+" ke ON kd.kb_extension_id = ke.id").
		Joins("JOIN "+VENDORS_TABLE+" v ON ke.vendor_id = v.id").
		Joins("JOIN "+COMPOUND_TABLE+" c ON kd.compound_id = c.id").
		Where("kr.id = ?", krID).
		Select("v.id AS Vendor_ID, c.id AS Compound_ID").
		Scan(&result).Error; err != nil {
		return 0, 0, fmt.Errorf("failed to fetch vendor and compound by KRID %d: %w", krID, err)
	}

	// Return the results
	return result.VendorID, result.CompoundID, nil
}

// update kb_root_status
func UpdateKbRootStatus(KBRootID int, status, UserID string) error {
	now := time.Now()
	if err := db.GetDB().Table(KbRoot_TABLE).Where("id = ?", KBRootID).Updates(map[string]interface{}{
		"modified_on": now,
		"modified_by": UserID,
		"status":      status,
	}).Error; err != nil {
		return err
	}
	return nil
}

func GetAllCompletedKBRootDetails() (ord []m.OrderDetails, err error) {
	result := db.GetDB().
		Table("kb_root AS kr").
		Select(`
                kr.id,
                vendors.vendor_name,
                compounds.compound_name,
                kb_data.cell_no,
                kb_data.demand_date_time,
                kr.status
            `).
		Joins("JOIN kb_data AS kb_data ON kr.kb_data_id = kb_data.id").
		Joins("JOIN kb_extension AS kb_extension ON kb_data.kb_extension_id = kb_extension.id").
		Joins("JOIN vendors AS vendors ON kb_extension.vendor_id = vendors.id").
		Joins("JOIN compounds AS compounds ON kb_data.compound_id = compounds.id").
		Joins("JOIN kb_transaction AS kb_transaction ON kr.id = kb_transaction.kb_root_id").
		Joins("JOIN prod_process_line AS prod_process_line ON kb_transaction.prod_process_line_id = prod_process_line.id").
		Where("prod_process_line.prod_process_id = ?", 2).
		Where("kb_extension.status IN (?) OR kb_extension.status LIKE ?", []string{"dispatch", "quality", "InProductionProcess"}, "dispatch%").
		Order("kr.id DESC").
		Find(&ord)
	return ord, result.Error
}

func GetCompletedKBRootDetailsBySearch(condition string) (ord []m.OrderDetails, err error) {
	result := db.GetDB().
		Table("kb_root AS kr").
		Select(`
                kr.id,
                vendors.vendor_name,
                compounds.compound_name,
                kb_data.cell_no,
                kb_data.demand_date_time,
                kr.status
            `).
		Joins("JOIN kb_data AS kb_data ON kr.kb_data_id = kb_data.id").
		Joins("JOIN kb_extension AS kb_extension ON kb_data.kb_extension_id = kb_extension.id").
		Joins("JOIN vendors AS vendors ON kb_extension.vendor_id = vendors.id").
		Joins("JOIN compounds AS compounds ON kb_data.compound_id = compounds.id").
		Joins("JOIN kb_transaction AS kb_transaction ON kr.id = kb_transaction.kb_root_id").
		Joins("JOIN prod_process_line AS prod_process_line ON kb_transaction.prod_process_line_id = prod_process_line.id").
		Where("prod_process_line.prod_process_id = ?", 2).
		Where("kb_extension.status IN (?) OR kb_extension.status LIKE ?", []string{"dispatch", "quality", "InProductionProcess"}, "dispatch%").
		Where(condition).
		Order("kr.id DESC").
		Find(&ord)
	return ord, result.Error
}

func GetDetailRootData(kbRootID int) (m.DetailRootData, error) {
	var result m.DetailRootData
	query := `
	SELECT 
		c.compound_name,
		kd.cell_no,
		kd.no_of_lots,
		pl.name,
		ke.order_id,
		ke.status,
		v.vendor_name,
		v.vendor_code,
		v.contact_info,
		pp.name,
		pp.expected_mean_time,
		kt.created_on,
		kt.completed_on,
		kt.operator AS operator, 
		kr.lot_no,
		kr.Status AS kanban_status,
		kr.kanban_no AS kanban_no,
		kr.quality_done_time AS quality_done_time,
		kr.dispatch_done_time AS dispatch_done_time,
		kr.quality_note AS quality_note,
		kr.dispatch_note AS dispatch_note,
		kr.quality_operator AS quality_operator,
		kr.packing_operator AS packing_operator
	FROM
		kb_root kr 
	JOIN 
		kb_transaction kt ON kr.id = kt.kb_root_id 
	JOIN 
		prod_process_line ppl ON kt.prod_process_line_id = ppl.id 
	JOIN 
		prod_process pp ON ppl.prod_process_id = pp.id 
	JOIN 
		prod_line pl ON ppl.prod_line_id = pl.id 
	JOIN 
		kb_data kd ON kr.kb_data_id = kd.id 
	JOIN 
		kb_extension ke ON kd.kb_extension_id = ke.id 
	JOIN 
		vendors v ON ke.vendor_id = v.id 
	JOIN 
		compounds c ON kd.compound_id = c.id 
	WHERE 
		kr.id = ?
	ORDER BY 
		ppl.order;
	`

	rows, err := db.GetDB().Raw(query, kbRootID).Rows()
	if err != nil {
		return m.DetailRootData{}, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	var kanban m.KanbanDetails

	for rows.Next() {
		var process m.ProductionProcess
		var compoundName, cellNo, prodLine, orderID, orderStatus, vendorName, vendorCode, lotNo string
		var noOfLots int
		var contactInfo, expectedMeanTime, KanbanNo, QualityNote, DispatchNote, KanbanStatus, operator, QualityOperator, PackingOperator sql.NullString
		var QualityDoneTime, DispatchDoneTime sql.NullTime
		var startedOn, completedOn string

		if err := rows.Scan(
			&compoundName, &cellNo, &noOfLots, &prodLine,
			&orderID, &orderStatus, &vendorName, &vendorCode,
			&contactInfo, &process.ProcessName, &expectedMeanTime,
			&startedOn, &completedOn, &operator, &lotNo, &KanbanStatus, &KanbanNo, &QualityDoneTime, &DispatchDoneTime, &QualityNote, &DispatchNote,
			&QualityOperator, &PackingOperator,
		); err != nil {
			return m.DetailRootData{}, fmt.Errorf("failed to scan row: %v", err)
		}

		// Set once only
		if len(result.KanbanDetails) == 0 {
			result.CompoundName = compoundName
			result.CellNo = cellNo
			result.NoOFLots = noOfLots
			result.OrderID = orderID
			result.Status = orderStatus
			result.VendorName = vendorName
			result.VendorCode = vendorCode
			result.LotNo = lotNo
			if contactInfo.Valid {
				result.ContactInfo = contactInfo.String
			}
			result.KanbanStatus = KanbanStatus.String
			result.KanbanNo = KanbanNo.String
			result.QualityDoneTime = QualityDoneTime.Time.String()
			result.DispatchDoneTime = DispatchDoneTime.Time.String()
			result.QualityNote = QualityNote.String
			result.DispatchNote = DispatchNote.String
			result.QualityOperator = QualityOperator.String
			result.PackingOperator = PackingOperator.String
			kanban.ProdLine = prodLine
		}

		if expectedMeanTime.Valid {
			process.ExpectedMeanTime = expectedMeanTime.String
		}
		process.StartedOn = startedOn
		process.CompletedOn = completedOn
		process.Operator = operator.String

		kanban.ProdProcesses = append(kanban.ProdProcesses, process)
	}

	if err := rows.Err(); err != nil {
		return m.DetailRootData{}, fmt.Errorf("row iteration error: %v", err)
	}

	if len(kanban.ProdProcesses) > 0 {
		result.KanbanDetails = append(result.KanbanDetails, kanban)
	}

	return result, nil
}

// DeleteKbExtensionByProdLine deletes kb_data records for the given production line ID.
func DeleteKbRootByProdLine(prodLineID int) error {
	return db.GetDB().
		Table(KbRoot_TABLE).
		Where("kb_root.id in(select kb_root_id FROM kb_transaction WHERE prod_process_line_id IN (SELECT id FROM prod_process_line WHERE prod_line_id = ?);) ", prodLineID).
		Delete(&m.KbTransaction{}).Error
}

// UpdateRunningAndInitialNumber sets the running and initial numbers to 0 and updates the modified timestamp.
func UpdateRunningAndInitialNumberToZero(KbRootId int) error {
	now := time.Now()

	// Update the database tableNoofLots
	if err := db.GetDB().Table(KbRoot_TABLE).Where("id = ?", KbRootId).Updates(map[string]interface{}{
		"modified_on": now,
		"running_no":  0,
		"initial_no":  0, // Corrected typo
	}).Error; err != nil {
		return err
	}
	return nil
}

// GetEntriesByMonth returns the number of entries in the kb_root table
// where lot_no is not blank for the given month.
func GetEntriesByMonth(currentMonth string) (int, error) {
	var countStr string

	// Perform the query to count the number of entries where lot_no is not blank
	err := db.GetDB().Table(KbRoot_TABLE+" kr").
		Where("EXTRACT(MONTH FROM kr.created_on) = ? AND kr.lot_no != ''", getMonthNumber(currentMonth)).
		Select("COUNT(*)").Scan(&countStr).Error
	if err != nil {
		return 0, fmt.Errorf("failed to fetch entries for the month %s: %w", currentMonth, err)
	}

	// Convert the result from string to int
	count, err := strconv.Atoi(countStr)
	if err != nil {
		return 0, fmt.Errorf("failed to convert count to integer: %w", err)
	}

	// Return the count of entries
	return count, nil
}

// Helper function to convert month name to month number (1-12)
func getMonthNumber(monthName string) int {
	monthTime, err := time.Parse("January", monthName)
	if err != nil {
		return 0
	}
	return int(monthTime.Month())
}

func GetinInventoryKbRootByCompoundID(compoundID string, limit int) ([]m.KbRoot, error) {
	var roots []m.KbRoot
	err := db.GetDB().
		Table(KbRoot_TABLE).
		Joins("LEFT OUTER JOIN kb_data ON kb_root.kb_data_id = kb_data.id").
		Where("in_inventory = true AND kb_root.status = ? AND kb_data.compound_id = ?", "3", compoundID).
		Order("kb_root.created_on ASC").
		Limit(limit).
		Find(&roots).Error

	return roots, err
}

func GetOrderStatusByKbroot(id int) (status m.KbExtension, err error) {
	result := db.GetDB().
		Table("kb_root").
		Select("kb_extension.*").
		Joins("LEFT OUTER JOIN kb_data on kb_data.id=kb_root.kb_data_id ").
		Joins("LEFT OUTER JOIN kb_extension on kb_data.kb_extension_id=kb_extension.id  ").
		Where("kb_root.id= ? ", id).
		Find(&status)
	return status, result.Error
}

func GetKbtransactionByKbroot(kbrootId int) (status m.KbTransaction, err error) {
	result := db.GetDB().
		Table("kb_root").
		Select("kb_transaction.*").
		Joins("LEFT OUTER JOIN kb_transaction on kb_transaction.kb_root_id =kb_root.id ").
		Where("kb_root.id= ? ", kbrootId).
		Find(&status)
	return status, result.Error
}

// DeleteKbRootByIDs deletes kb_root records based on a list of IDs
func DeleteKbRootByIDs(ids []string) error {
	query := db.GetDB().Table(KbRoot_TABLE).Where("id IN (?)", ids).Delete(&m.KbRoot{})
	if query.Error != nil {
		return query.Error // Return error if query fails
	}
	if query.RowsAffected == 0 {
		return fmt.Errorf("no records found for the given IDs")
	}
	return nil
}

func GetAllKanbanDetailsForReport(pagination m.PaginationReq, conditions []string) (ord []m.OrderDetails, paginationResp m.PaginationResp, err error) {
	dbQuery := db.GetDB().
		Table("kb_root AS kr").
		Select("DISTINCT ON (kr.id) kr.id, vendors.vendor_name, compounds.compound_name, kb_data.cell_no, kb_data.demand_date_time, kr.status, kr.lot_no").
		Joins("JOIN kb_data ON kr.kb_data_id = kb_data.id").
		Joins("JOIN kb_extension ON kb_data.kb_extension_id = kb_extension.id").
		Joins("JOIN vendors ON kb_extension.vendor_id = vendors.id").
		Joins("JOIN compounds ON kb_data.compound_id = compounds.id").
		Joins("LEFT JOIN kb_transaction ON kr.id = kb_transaction.kb_root_id").
		Joins("LEFT JOIN prod_process_line ON kb_transaction.prod_process_line_id = prod_process_line.id")

	// Apply search filters dynamically
	var parsedConditions []string
	var dateFrom, dateTo string

	for _, cond := range conditions {
		if strings.Contains(cond, "BETWEEN") {
			// Handle demand_date_time BETWEEN condition
			parts := strings.SplitN(cond, "BETWEEN", 2)
			if len(parts) == 2 {
				field := strings.TrimSpace(parts[0]) // Should be "demand_date_time"
				dateRange := strings.TrimSpace(parts[1])
				dateParts := strings.Split(dateRange, " AND ")

				if len(dateParts) == 2 {
					dateFrom = strings.TrimSpace(dateParts[0])
					dateTo = strings.TrimSpace(dateParts[1])
					parsedConditions = append(parsedConditions, fmt.Sprintf("kb_data.%s BETWEEN '%s' AND '%s'",
						field,
						strings.Trim(dateFrom, "'"),
						strings.Trim(dateTo, "'"),
					))
				}
			}
		} else {
			// Handle normal ILIKE conditions
			parts := strings.SplitN(cond, " ILIKE ", 2)
			if len(parts) < 2 {
				continue
			}

			field := strings.TrimSpace(parts[0])
			value := strings.Trim(parts[1], "'%")

			if value == "" {
				continue
			}

			if field == "status" {
				parsedConditions = append(parsedConditions, fmt.Sprintf("kr.status::INTEGER = %s", value))
			} else if field == "demand_date_time" {
				parsedConditions = append(parsedConditions, fmt.Sprintf("CAST(kb_data.%s AS TEXT) ILIKE '%%%s%%'", field, value))
			} else {
				parsedConditions = append(parsedConditions, fmt.Sprintf("%s ILIKE '%%%s%%'", field, value))
			}
		}
	}

	// Apply conditions if any
	if len(parsedConditions) > 0 {
		dbQuery = dbQuery.Where(strings.Join(parsedConditions, " AND "))
	}

	// Get total count of distinct kr.id before pagination
	var totalRecords int64
	countQuery := db.GetDB().
		Table("kb_root AS kr").
		Select("COUNT(DISTINCT kr.id)").
		Joins("JOIN kb_data ON kr.kb_data_id = kb_data.id").
		Joins("JOIN kb_extension ON kb_data.kb_extension_id = kb_extension.id").
		Joins("JOIN vendors ON kb_extension.vendor_id = vendors.id").
		Joins("JOIN compounds ON kb_data.compound_id = compounds.id").
		Joins("LEFT JOIN kb_transaction ON kr.id = kb_transaction.kb_root_id").
		Joins("LEFT JOIN prod_process_line ON kb_transaction.prod_process_line_id = prod_process_line.id")

	// Apply the same conditions to the count query
	if len(parsedConditions) > 0 {
		countQuery = countQuery.Where(strings.Join(parsedConditions, " AND "))
	}

	// Execute the count query
	if err := countQuery.Count(&totalRecords).Error; err != nil {
		return nil, m.PaginationResp{}, err
	}

	// Set sorting order
	orderBy := "kr.id DESC" // Default sorting
	if pagination.Order != "" {
		orderBy = pagination.Order
	}
	dbQuery = dbQuery.Order(orderBy)

	// Set pagination safely
	limit, errLimit := strconv.Atoi(pagination.Limit)
	pageNo := pagination.PageNo

	if errLimit != nil || limit <= 0 {
		limit = 15 // Default limit
	}
	if pageNo <= 0 {
		pageNo = 1 // Default page number
	}
	offset := (pageNo - 1) * limit

	dbQuery = dbQuery.Limit(limit).Offset(offset)

	// Execute query
	if err := dbQuery.Find(&ord).Error; err != nil {
		return nil, m.PaginationResp{}, err
	}

	// Build pagination response
	paginationResp = m.PaginationResp{
		TotalNo: int(totalRecords),
		Page:    pageNo,
		Offset:  offset,
	}
	return ord, paginationResp, nil
}
