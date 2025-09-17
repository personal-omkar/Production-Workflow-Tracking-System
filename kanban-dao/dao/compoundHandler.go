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
	COMPOUND_TABLE string = "compounds" // Updated to match the table name
)

// CreateNewOrUpdateExistingCompound creates a new compound or updates an existing compound
func CreateNewOrUpdateExistingCompound(compound *m.Compounds) error {

	now := time.Now()
	if compound.Id != 0 {
		compound.ModifiedOn = now
		CompData, _ := GetCompoundDataByParam("id", strconv.Itoa(compound.Id))
		compound.CreatedBy = CompData[0].CreatedBy
		compound.CreatedOn = CompData[0].CreatedOn
		if err := db.GetDB().Table(COMPOUND_TABLE).Save(&compound).Error; err != nil {
			return err
		}
	} else {
		compound.CreatedOn = now

		if err := db.GetDB().Table(COMPOUND_TABLE).Create(&compound).Error; err != nil {
			return err
		}
	}
	return nil
}

// GetAllCompoundEntries returns a all records present in compound table
func GetAllCompoundEntries() (entries []m.Compounds, err error) {
	result := db.GetDB().Table(COMPOUND_TABLE).Order("id").Find(&entries)
	return entries, result.Error
}

// GetKBDataByParam returns a kb_data records based on parameter
func GetCompoundDataByParam(key string, value any) (comp []m.Compounds, err error) {
	query := db.GetDB().Table(COMPOUND_TABLE)
	// Use parameterized queries to prevent SQL injection
	query = query.Where(key+" = ?", value).Order("id").Find(&comp)

	return comp, query.Error
}

// GetCompoundDataByParamAndCondition returns compound records based on key-value pair and filtered conditions
func GetCompoundDataByParamAndCondition(key, value string, conditions []string) (comp []m.Compounds, err error) {

	query := db.GetDB().Table(COMPOUND_TABLE)
	// Apply key-value filter
	query = query.Where(key+" = ?", value)
	// Filter conditions that contain "compound"
	var compoundConditions []string
	for _, condition := range conditions {
		if strings.Contains(strings.ToLower(condition), "compound") {
			compoundConditions = append(compoundConditions, condition)
		}
	}
	// Apply only compound-related conditions
	if len(compoundConditions) > 0 {
		query = query.Where(strings.Join(compoundConditions, " AND "))
	}
	result := query.Find(&comp)
	return comp, result.Error
}

func GetAllCompoundsBySearchAndPagination(pagination m.PaginationReq, conditions []string) (comp []*m.Compounds, paginationResp m.PaginationResp, err error) {
	// Base query
	dbQuery := db.GetDB().Table(COMPOUND_TABLE)

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
	countQuery := db.GetDB().Table(COMPOUND_TABLE)
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
	if err := dbQuery.Find(&comp).Error; err != nil {
		return nil, m.PaginationResp{}, err
	}

	// Prepare pagination response
	paginationResp = m.PaginationResp{
		TotalNo: int(totalRecords),
		Page:    pageNo,
		Offset:  offset,
	}

	return comp, paginationResp, nil
}
