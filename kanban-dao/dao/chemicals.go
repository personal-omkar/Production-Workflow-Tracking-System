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
	CHEMICALS_TABLE string = "chemical_types"
)

// CreateNewOrUpdateExistingChemical creates a new chemical or updates an existing chemical
func CreateNewOrUpdateExistingChemical(op *m.ChemicalTypes) error {

	now := time.Now()
	if op.Id != 0 {
		op.ModifiedOn = now

		if err := db.GetDB().Table(CHEMICALS_TABLE).Save(&op).Error; err != nil {
			return err
		}
	} else {
		op.CreatedOn = now

		if err := db.GetDB().Table(CHEMICALS_TABLE).Create(&op).Error; err != nil {
			return err
		}
	}
	return nil
}

// GetAllChemical returns a all records present in chemical table
func GetAllChemical() (op []*m.ChemicalTypes, err error) {
	result := db.GetDB().Table(CHEMICALS_TABLE).Order("id").Find(&op)
	return op, result.Error
}

// GetChemicalByParam returns a chemical records based on parameter
func GetChemicalByParam(key, value string) (op []*m.ChemicalTypes, err error) {
	query := db.GetDB().Table(CHEMICALS_TABLE)
	query = query.Where(key + " = " + value).Find(&op)
	return op, query.Error
}

// DeleteChemicalByParam deletes chemical records for the given parameter
func DeleteChemicalByParam(key, value string) error {
	return db.GetDB().
		Table(CHEMICALS_TABLE).
		Where(key + " = " + value).
		Delete(&m.ChemicalTypes{}).Error
}

func GetAllChemicalBySearchAndPagination(pagination m.PaginationReq, conditions []string) (op []*m.ChemicalTypes, paginationResp m.PaginationResp, err error) {
	// Base query
	dbQuery := db.GetDB().Table(CHEMICALS_TABLE)

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
	countQuery := db.GetDB().Table(CHEMICALS_TABLE)
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
