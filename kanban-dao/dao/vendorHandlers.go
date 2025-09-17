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
	VENDORS_TABLE string = "vendors" // Updated to match the table name
)

// CreateNewOrUpdateExistingVendors creates a new vendor or updates an existing vendor
func CreateNewOrUpdateExistingVendors(vendor *m.Vendors) error {
	now := time.Now()
	if vendor.ID != 0 {
		vendor.ModifiedOn = now

		if err := db.GetDB().Table(VENDORS_TABLE).Omit("created_on").Save(&vendor).Error; err != nil {
			return err
		}
	} else {
		vendor.CreatedOn = now

		if err := db.GetDB().Table(VENDORS_TABLE).Create(&vendor).Error; err != nil {
			return err
		}
	}
	return nil
}

// GetAllVendorEntries returns a all records present in Vendors table
func GetAllVendorEntries() (entries []m.Vendors, err error) {
	result := db.GetDB().Table(VENDORS_TABLE).Order("id ASC").Find(&entries)
	return entries, result.Error
}

// GetAllKbRootByParam returns a all records present in kb_data table
func GetAllVendorByParam(code string) (v []m.Vendors, err error) {
	result := db.GetDB().Table(VENDORS_TABLE).Where("vendor_code" + "!=" + code).Find(&v)
	return v, result.Error
}

// GetAllVendorEntriesByCondition returns vendor records based on vendor-related conditions
func GetAllVendorEntriesByCondition(conditions []string) (entries []m.Vendors, err error) {
	query := db.GetDB().Table(VENDORS_TABLE).Order("id ASC")

	// Filter conditions that contain "vendor"
	var vendorConditions []string
	for _, condition := range conditions {
		if strings.Contains(strings.ToLower(condition), "vendor") {
			vendorConditions = append(vendorConditions, condition)
		}
	}

	// Apply only vendor-related conditions
	if len(vendorConditions) > 0 {
		query = query.Where(strings.Join(vendorConditions, " AND "))
	}

	result := query.Find(&entries)
	return entries, result.Error
}

// GetVendorByParam returns a vendor records based on parameter
func GetVendorByParam(key, value any) (kbe []m.Vendors, err error) {
	query := db.GetDB().Table(VENDORS_TABLE)
	query = query.Where(fmt.Sprintf("%s = ?", key), value)
	// Execute the query and return the results
	query = query.Find(&kbe)
	return kbe, query.Error
}

// GetAllVendorDetails returns a all vendor records based on vendor code
func GetVendorDetailsByVendorCode(vendorcode string) (entries m.Vendors, err error) {
	result := db.GetDB().Table(VENDORS_TABLE).Where("vendor_code = ?", vendorcode).Find(&entries)
	return entries, result.Error
}

// DeleteVendor deletes a vendor record by id
func DeleteVendor(id int) error {
	return db.GetDB().Table(VENDORS_TABLE).Where("id = ?", id).Delete(&m.Vendors{}).Error
}

// GetVendorByParamStartsWith returns vendor records where the specified key starts with the given value
func GetVendorByParamStartsWith(key, value string) (kbe []m.Vendors, err error) {
	query := db.GetDB().Table(VENDORS_TABLE)
	if key == "id" {
		return nil, fmt.Errorf("starts with search is not supported for id")
	}
	// Use LIKE to match values that start with the input
	query = query.Where(fmt.Sprintf("%s LIKE ?", key), value+"%")
	query = query.Find(&kbe)
	return kbe, query.Error
}

func GetAllVendorBySearchAndPagination(pagination m.PaginationReq, conditions []string) (op []*m.Vendors, paginationResp m.PaginationResp, err error) {
	// Base query
	dbQuery := db.GetDB().Table(VENDORS_TABLE)

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
	countQuery := db.GetDB().Table(VENDORS_TABLE)
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
