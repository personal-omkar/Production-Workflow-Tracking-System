package dao

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
	m "irpl.com/kanban-commons/model"
	db "irpl.com/kanban-dao/db"
)

const (
	KBData_TABLE string = "kb_data" // Updated to match the table name
)

// CreateNewOrUpdateExistingKBData creates a new kbdata or updates an existing kbdata
func CreateNewOrUpdateExistingKBData(kbd *m.KbData) (ID int, err error) {
	var kbdataint int
	now := time.Now()
	if kbd.Id != 0 {
		kbd.ModifiedOn = now

		if err := db.GetDB().Table(KBData_TABLE).Omit("compound_id").Save(&kbd).Error; err != nil {
			return kbd.Id, err
		}
	} else {
		kbd.CreatedOn = now

		if err := db.GetDB().Table(KBData_TABLE).Create(&kbd).Error; err != nil {
			return kbd.Id, err
		}
	}
	kbdataint = kbd.Id
	return kbdataint, err
}

// GetAllKBData returns a all records present in kb_data table
func GetAllKBData() (kbe []m.KbData, err error) {
	result := db.GetDB().Table(KBData_TABLE).Find(&kbe)
	return kbe, result.Error
}

// GetKBDataByParam returns kb_data records based on a parameter
func GetKBDataByParam(key, value string) (kbe []m.KbData, err error) {
	query := db.GetDB().Table(KBData_TABLE)

	// Use placeholders to avoid SQL injection
	if key == "id" {
		id, convErr := strconv.Atoi(value)
		if convErr != nil {
			return nil, fmt.Errorf("invalid value for id: %w", convErr)
		}
		query = query.Where(key+" = ?", id).Find(&kbe)
	} else {
		query = query.Where(key+" = ?", value).Find(&kbe)
	}
	return kbe, query.Error
}

func GetCustomerOrderDetails(condition string) (ord []m.OrderDetails, err error) {
	var result *gorm.DB
	if len(condition) > 0 {
		result = db.GetDB().
			Table(KBData_TABLE).
			Select("users.username As CustomerName,kb_data.* , kb_extension.status ,  vendors.vendor_name, compounds.compound_name ").
			Joins(" JOIN users ON CAST(users.id AS TEXT)= kb_data.created_by ").
			Joins(" JOIN kb_extension ON kb_extension.id = kb_data.kb_extension_id ").
			Joins("JOIN vendors ON kb_extension.vendor_id = vendors.id").
			Joins("JOIN compounds  ON kb_data.compound_id = compounds.id").
			Where(condition).Order("kb_data.kb_extension_id DESC").Find(&ord)
	} else {
		result = db.GetDB().
			Table(KBData_TABLE).
			Select("users.username As CustomerName,kb_data.* , kb_extension.status ,  vendors.vendor_name, compounds.compound_name ").
			Joins(" JOIN users ON CAST(users.id AS TEXT)= kb_data.created_by ").
			Joins("JOIN vendors ON kb_extension.vendor_id = vendors.id").
			Joins("JOIN compounds  ON kb_data.compound_id = compounds.id;").
			Order("kb_data.kb_extension_id DESC").Find(&ord)
	}

	return ord, result.Error
}

func GetOrderDetails(condition string) (ord []m.OrderDetails, err error) {
	var result *gorm.DB
	if len(condition) > 0 {
		result = db.GetDB().
			Table(KBData_TABLE).
			Select(`  
				kb_data.*,
				kb_extension.status,
				vendors.vendor_name,
				compounds.compound_name,
				COALESCE(inventory.min_quantity, 0) AS min_quantity,
				COALESCE(inventory.available_quantity, 0) AS available_quantity
			`).
			Joins("JOIN kb_extension ON kb_extension.id = kb_data.kb_extension_id").
			Joins("JOIN vendors ON kb_extension.vendor_id = vendors.id").
			Joins("JOIN compounds ON kb_data.compound_id = compounds.id").
			Joins("LEFT JOIN inventory ON compounds.id = inventory.compound_id").
			Where(condition).
			Order("kb_data.kb_extension_id DESC").
			Find(&ord)
	} else {
		result = db.GetDB().
			Table(KBData_TABLE).
			Select(`
				kb_data.*,
				kb_extension.status,
				vendors.vendor_name,
				compounds.compound_name,
				COALESCE(inventory.min_quantity, 0) AS min_quantity,
				COALESCE(inventory.available_quantity, 0) AS available_quantity
			`).
			Joins("JOIN kb_extension ON kb_extension.id = kb_data.kb_extension_id").
			Joins("JOIN vendors ON kb_extension.vendor_id = vendors.id").
			Joins("JOIN compounds ON kb_data.compound_id = compounds.id").
			Joins("LEFT JOIN inventory ON compounds.id = inventory.compound_id").
			Order("kb_data.kb_extension_id DESC").
			Find(&ord)
	}
	return ord, result.Error
}

func DeleteKBDataByParam(key, value string) error {
	query := db.GetDB().Table(KBData_TABLE)

	// Use placeholders to avoid SQL injection
	if key == "id" {
		id, convErr := strconv.Atoi(value)
		if convErr != nil {
			return fmt.Errorf("invalid value for id: %w", convErr)
		}
		query = query.Where(key+" = ?", id).Delete(nil)
	} else {
		query = query.Where(key+" = ?", value).Delete(nil)
	}

	// Check if the query resulted in an error
	if query.Error != nil {
		return fmt.Errorf("failed to delete KB data: %w", query.Error)
	}

	// Check if any rows were actually deleted
	if query.RowsAffected == 0 {
		return fmt.Errorf("no record found for %s = %s", key, value)
	}

	return nil
}

func GetAllDetailsForOrder(kbDataId string) (m.OrderDetailsHistory, error) {
	var result m.OrderDetailsHistory
	id, err := strconv.Atoi(kbDataId)
	if err != nil {
		return m.OrderDetailsHistory{}, fmt.Errorf("invalid kbDataId: %v", err)
	}
	err = db.GetDB().
		Table("kb_data").
		Select(`
			kb_data.id,
			kb_data.no_of_lots,
			kb_data.cell_no,
			kb_data.created_on,
			kb_data.demand_date_time,
			kb_extension.modified_on,
			kb_extension.status,
			kb_extension.order_id,
			vendors.vendor_name,
			vendors.vendor_code,
			vendors.contact_info,
			users.username,
			users.email,
			compounds.compound_name
		`).
		Joins("JOIN kb_extension ON kb_extension.id = kb_data.kb_extension_id").
		Joins("JOIN compounds ON kb_data.compound_id = compounds.id").
		Joins("JOIN vendors ON kb_extension.vendor_id = vendors.id").
		Joins("JOIN users ON kb_data.created_by::integer = users.id").
		Where("kb_data.id = ?", id).
		Scan(&result).Error
	if err != nil {
		return m.OrderDetailsHistory{}, err
	}
	return result, nil
}

// DeleteKbDataByProdLine deletes kb_data records for the given production line ID.
func DeleteKbDataByProdLine(prodLineID int) error {
	return db.GetDB().
		Table(KBData_TABLE).
		Where("kb_data.kb_root_id in(select kb_root_id FROM kb_transaction WHERE prod_process_line_id IN (SELECT id FROM prod_process_line WHERE prod_line_id = ?)) ", prodLineID).
		Delete(&m.KbTransaction{}).Error //";"
}

func GetOrderStatusByKbData(id int) (status m.KbExtension, err error) {
	result := db.GetDB().
		Table("kb_data").
		Select("kb_extension.*").
		Joins("LEFT OUTER JOIN kb_extension on kb_data.kb_extension_id=kb_extension.id").
		Where("kb_data.id = ? ", id).
		Find(&status)
	return status, result.Error
}

func GetOrderDetailsBySearchAndPagination(pagination m.PaginationReq, conditions []string) (ord []*m.OrderDetails, paginationResp m.PaginationResp, err error) {
	// Base query
	dbQuery := db.GetDB().
		Table(KBData_TABLE).
		Select(`  
			kb_data.*,
			kb_extension.status,
			vendors.vendor_name,
			compounds.compound_name,
			COALESCE(inventory.min_quantity, 0) AS min_quantity,
			COALESCE(inventory.available_quantity, 0) AS available_quantity
		`).
		Joins("JOIN kb_extension ON kb_extension.id = kb_data.kb_extension_id").
		Joins("JOIN vendors ON kb_extension.vendor_id = vendors.id").
		Joins("JOIN compounds ON kb_data.compound_id = compounds.id").
		Joins("LEFT JOIN inventory ON compounds.id = inventory.compound_id").
		Where("kb_extension.status = ?", "pending")

	// Parse search conditions
	var parsedConditions []string
	for _, cond := range conditions {
		parts := strings.SplitN(cond, " ILIKE ", 2)
		if len(parts) < 2 {
			continue
		}

		field := strings.TrimSpace(parts[0])
		value := strings.Trim(parts[1], "'%")

		if value != "" {
			parsedConditions = append(parsedConditions, fmt.Sprintf("%s ILIKE '%%%s%%'", field, value))
		}
	}

	// Apply dynamic search filters
	if len(parsedConditions) > 0 {
		dbQuery = dbQuery.Where(strings.Join(parsedConditions, " AND "))
	}

	// Total count
	var totalRecords int64
	countQuery := db.GetDB().
		Table(KBData_TABLE).
		Joins("JOIN kb_extension ON kb_extension.id = kb_data.kb_extension_id").
		Joins("JOIN vendors ON kb_extension.vendor_id = vendors.id").
		Joins("JOIN compounds ON kb_data.compound_id = compounds.id").
		Where("kb_extension.status = ?", "pending")

	if len(parsedConditions) > 0 {
		countQuery = countQuery.Where(strings.Join(parsedConditions, " AND "))
	}

	if err := countQuery.Count(&totalRecords).Error; err != nil {
		return nil, m.PaginationResp{}, err
	}

	// Sorting
	orderBy := "kb_data.kb_extension_id DESC"
	if pagination.Order != "" {
		orderBy = pagination.Order
	}
	dbQuery = dbQuery.Order(orderBy)

	// Pagination logic
	limit, errLimit := strconv.Atoi(pagination.Limit)
	pageNo := pagination.PageNo
	if errLimit != nil || limit <= 0 {
		limit = 15
	}
	if pageNo <= 0 {
		pageNo = 1
	}
	offset := (pageNo - 1) * limit

	dbQuery = dbQuery.Limit(limit).Offset(offset)

	// Execute
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
