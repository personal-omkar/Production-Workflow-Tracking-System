package dao

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	m "irpl.com/kanban-commons/model" // Adjust import path for your models
	db "irpl.com/kanban-dao/db"       // Adjust import path for your DB connection
)

const UserRolesTable = "user_roles" // Define the table name for UserRoless

// CreateUserRoles creates a new UserRoles record
func CreateUserRoles(entry m.UserRoles) error {
	now := time.Now()
	entry.CreatedOn = now
	return db.GetDB().Table(UserRolesTable).Create(&entry).Error
}

// GetUserRolesByID retrieves a UserRoles record by ID
func GetUserRolesByID(id int64) (entry m.UserRoles, err error) {
	result := db.GetDB().Table(UserRolesTable).Where("id = ?", id).First(&entry)
	return entry, result.Error
}

// UpdateUserRoles updates an existing UserRoles record
func UpdateUserRoles(entry m.UserRoles) error {
	now := time.Now()
	entry.ModifiedOn = now
	return db.GetDB().Table(UserRolesTable).Where("id = ?", entry.ID).Updates(&entry).Error
}

// DeleteUserRoles deletes a UserRoles record by ID
func DeleteUserRoles(id int64) error {
	return db.GetDB().Table(UserRolesTable).Where("id = ?", id).Delete(&m.UserRoles{}).Error
}

// GetUserRoles returns a paginated list of UserRoles records
func GetUserRoles(page, limit int) (entries []m.UserRoles, totalRecords int, err error) {
	offset := (page - 1) * limit
	var totalRecords64 int64

	// Get total record count
	err = db.GetDB().Table(UserRolesTable).Count(&totalRecords64).Error
	if err != nil {
		return nil, 0, err
	}

	// Convert totalRecords64 to int
	totalRecords = int(totalRecords64)

	// Fetch records with pagination
	err = db.GetDB().Table(UserRolesTable).Offset(offset).Limit(limit).Find(&entries).Error
	if err != nil {
		return nil, 0, err
	}

	return entries, totalRecords, nil
}

// GetUserRolesByCriteria returns a paginated list of UserRoles records based on criteria
func GetUserRolesByCriteria(page, limit int, criteria map[string]interface{}) (entries []m.UserRoles, totalRecords int, err error) {
	offset := (page - 1) * limit
	var totalRecords64 int64
	query := db.GetDB().Table(UserRolesTable)

	// Apply filters from criteria
	for key, value := range criteria {
		query = query.Where(key+" = ?", value)
	}

	// Get total record count with filters
	err = query.Count(&totalRecords64).Error
	if err != nil {
		return nil, 0, err
	}

	// Convert totalRecords64 to int
	totalRecords = int(totalRecords64)

	// Fetch records with pagination and filters
	err = query.Offset(offset).Limit(limit).Scan(&entries).Error
	if err != nil {
		return nil, 0, err
	}

	return entries, totalRecords, nil
}

// GetUserRolesByName retrieves a UserRoles record by name (case-insensitive)
func GetUserRolesByName(roleName string) (entry m.UserRoles, err error) {
	// Use the ILIKE operator for case-insensitive comparison
	result := db.GetDB().Table(UserRolesTable).Where("LOWER(role_name) = LOWER(?)", roleName).First(&entry)
	return entry, result.Error
}

// GetAllUserRoles returns a all records present in user roles table
func GetAllUserRoles() (entries []*m.UserRoles, err error) {
	result := db.GetDB().Table(UserRolesTable).Find(&entries)
	return entries, result.Error
}
func GetAllUserRoleBySearchAndPagination(pagination m.PaginationReq, conditions []string) (op []*m.UserRoles, paginationResp m.PaginationResp, err error) {
	// Base query
	dbQuery := db.GetDB().Table(UserRolesTable)

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
	countQuery := db.GetDB().Table(UserRolesTable)
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
