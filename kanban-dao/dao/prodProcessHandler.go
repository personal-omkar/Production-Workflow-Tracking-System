package dao

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	m "irpl.com/kanban-commons/model"
	db "irpl.com/kanban-dao/db"
)

const (
	PRODUCTPROCESS_TABLE string = "prod_process" // Updated to match the table name
)

// CreateNewOrUpdateExistingProductLine creates a new product line or updates an existing vendor
func CreateNewOrUpdateExistingProdProcess(pp *m.ProdProcess) error {

	now := time.Now()
	if pp.Id != 0 {
		pp.ModifiedOn = now

		if err := db.GetDB().Table(PRODUCTPROCESS_TABLE).Save(&pp).Error; err != nil {
			return err
		}
	} else {
		pp.CreatedOn = now

		if err := db.GetDB().Table(PRODUCTPROCESS_TABLE).Create(&pp).Error; err != nil {
			return err
		}
	}
	return nil
}

// GetAllProductionProcessEntries returns all records from the prod_process table
// with status = 'Active' and line_visibility = true.
func GetAllProductionProcessEntries() (entries []m.ProdProcess, err error) {
	result := db.GetDB().
		Table(PRODUCTPROCESS_TABLE).
		Where("status = ? AND line_visibility = ?", "Active", true).
		Find(&entries)
	return entries, result.Error
}

func GetAllProductionProcess() (entries []m.ProdProcess, err error) {
	result := db.GetDB().
		Table(PRODUCTPROCESS_TABLE).
		Order("id").
		Find(&entries)
	return entries, result.Error
}

func GetAllProductionProcessForLine(lineid string) (entries []m.ProdProcessCardData, err error) {
	result := db.GetDB().
		Select("prod_process_line.* , prod_process.* , prod_process_line.id As prod_process_line_id ").Table("prod_process_line").
		Joins("LEFT OUTER JOIN prod_process ON prod_process.id=prod_process_line.prod_process_id ").
		Where("prod_process.line_visibility = ? AND prod_process_line.prod_process_id NOT IN (?) AND prod_process_line.prod_line_id=?", true, []int{1, 2}, lineid).
		Find(&entries)
	return entries, result.Error
}

func GetProductionProcessCardData(productionlineid string) (entries []m.ProdProcessCardData, err error) {
	// Get Distinct ProdProcessLineIDs based on the production line
	ProdProcessLineID, err := GetDistinctProdProcessLineIDs(productionlineid)
	if err != nil {
		return nil, err
	}

	// Get filtered ProdProcessLineIDs based on the ProdProcessLineIDs obtained above
	ProdProcessLineIDs, err := GetFilteredProdProcessLineIDs(ProdProcessLineID)
	if err != nil {
		return nil, err
	}

	// Check if ProdProcessLineIDs is not empty
	if len(ProdProcessLineIDs) == 0 {
		return nil, fmt.Errorf("no valid prod_process_line_ids found")
	}

	// Convert the slice of ProdProcessLineIDs to integers
	var ids []int
	for _, id := range ProdProcessLineIDs {
		intID, err := strconv.Atoi(id)
		if err != nil {
			return nil, fmt.Errorf("invalid prod_process_line_id: %s", id)
		}
		ids = append(ids, intID)
	}

	// Run the query with the corrected IN clause
	result := db.GetDB().
		Select("DISTINCT ON (kb_transaction.kb_root_id) kb_transaction.* ,kb_root.* ,kb_data.* ,prod_process_line.* ,compounds.compound_name,prod_process.* ,kb_transaction.created_on").Table("kb_transaction").
		Joins("LEFT OUTER JOIN kb_root on kb_root.id = kb_transaction.kb_root_id").
		Joins("LEFT OUTER JOIN kb_data on kb_data.id = kb_root.kb_data_id").
		Joins("LEFT OUTER JOIN prod_process_line on prod_process_line.id = kb_transaction.prod_process_line_id").
		Joins("LEFT OUTER JOIN compounds  on kb_data.compound_id =compounds.id").
		Joins("LEFT OUTER JOIN prod_process on prod_process.id = prod_process_line.prod_process_id ").
		Where(" kb_transaction.prod_process_line_id IN (SELECT prod_process_line.id FROM prod_process_line WHERE prod_process_line.prod_line_id = ? )", productionlineid).
		Order("kb_transaction.kb_root_id, kb_transaction.started_on DESC, kb_transaction.id DESC;").Find(&entries)
	return entries, result.Error
}

func GetDistinctProdProcessLineIDs(productionLineID string) ([]string, error) {
	// Define a slice to store the result
	var prodProcessLineIDs []string

	// Execute the query
	result := db.GetDB().
		Raw(`
            SELECT DISTINCT ON (kb_transaction.kb_root_id)
                kb_transaction.prod_process_line_id
            FROM 
                kb_transaction
            WHERE 
                kb_transaction.prod_process_line_id IN (
                    SELECT 
                        prod_process_line.id 
                    FROM 
                        prod_process_line 
                    WHERE 
                        prod_process_line.prod_line_id = ?
                )
            ORDER BY 
                kb_transaction.kb_root_id, 
                kb_transaction.started_on DESC, 
                kb_transaction.id DESC;
        `, productionLineID).
		Scan(&prodProcessLineIDs)

	// Check for errors and return results
	if result.Error != nil {
		return nil, result.Error
	}
	return prodProcessLineIDs, nil
}

func GetFilteredProdProcessLineIDs(prodProcessLineIDs []string) ([]string, error) {
	var result []string

	for _, prodProcessLineID := range prodProcessLineIDs {
		var prodProcessLine m.ProdProcessLine

		query := `SELECT * FROM "prod_process_line" WHERE id = ? LIMIT 1;`
		if err := db.GetDB().Raw(query, prodProcessLineID).Scan(&prodProcessLine).Error; err != nil {
			return nil, err
		}

		if prodProcessLine.IsGroup {
			groupName := prodProcessLine.GroupName
			prodLineID := prodProcessLine.ProdLineId

			var prodProcessLines []m.ProdProcessLine
			groupQuery := `SELECT * FROM "prod_process_line" WHERE group_name = ? AND prod_line_id = ?;`
			if err := db.GetDB().Raw(groupQuery, groupName, prodLineID).Scan(&prodProcessLines).Error; err != nil {
				return nil, err
			}

			for _, line := range prodProcessLines {
				result = append(result, strconv.Itoa(line.Id))
			}
		} else {
			result = append(result, prodProcessLineID)
		}
	}

	return result, nil
}

// GetProdProcessByParam returns productionprocess records based on a parameter
func GetProdProcessByParam(key, value string) (kbe []m.ProdProcess, err error) {
	query := db.GetDB().Table(PRODUCTPROCESS_TABLE)
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

func GetAllProcessBySearchAndPagination(pagination m.PaginationReq, conditions []string) (op []*m.ProdProcess, paginationResp m.PaginationResp, err error) {
	// Base query
	dbQuery := db.GetDB().Table(PRODUCTPROCESS_TABLE)

	// Parse search conditions
	var parsedConditions []string
	for _, cond := range conditions {
		parts := strings.SplitN(cond, " ILIKE ", 2)
		if len(parts) < 2 {
			continue
		}

		field := strings.TrimSpace(parts[0])
		value := strings.Trim(parts[1], "'%")

		if value == "" {
			continue
		}

		parsedConditions = append(parsedConditions, fmt.Sprintf("%s ILIKE '%%%s%%'", field, value))
	}

	// Apply where clause
	if len(parsedConditions) > 0 {
		dbQuery = dbQuery.Where(strings.Join(parsedConditions, " AND "))
	}

	// Get total count
	var totalRecords int64
	countQuery := db.GetDB().Table(PRODUCTPROCESS_TABLE)
	if len(parsedConditions) > 0 {
		countQuery = countQuery.Where(strings.Join(parsedConditions, " AND "))
	}
	if err := countQuery.Count(&totalRecords).Error; err != nil {
		return nil, m.PaginationResp{}, err
	}

	// Set sorting order
	orderBy := "id DESC"
	if pagination.Order != "" {
		orderBy = pagination.Order
	}
	dbQuery = dbQuery.Order(orderBy)

	// Pagination
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

	// Execute the query
	if err := dbQuery.Find(&op).Error; err != nil {
		return nil, m.PaginationResp{}, err
	}

	// Prepare pagination response
	paginationResp = m.PaginationResp{
		TotalNo: int(totalRecords),
		Page:    pageNo,
		Offset:  offset,
	}

	return op, paginationResp, nil
}
