package dao

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	m "irpl.com/kanban-commons/model"
	db "irpl.com/kanban-dao/db"
)

const (
	INVENTORY_TABLE string = "inventory" // Updated to match the table name
)

// CreateNewOrUpdateExistingInventory creates a new inventory or updates an existing inventory
func CreateNewOrUpdateExistingInventory(inv *m.Inventory) error {

	now := time.Now()
	if inv.Id != 0 {
		inv.ModifiedOn = now

		if err := db.GetDB().Table(INVENTORY_TABLE).Save(&inv).Error; err != nil {
			return err
		}
	} else {
		inv.CreatedOn = now

		if err := db.GetDB().Table(INVENTORY_TABLE).Create(&inv).Error; err != nil {
			return err
		}
	}
	return nil
}

// GetAllInventory returns a all records present in inventory table
func GetAllInventory() (inv []*m.Inventory, err error) {
	result := db.GetDB().Table(INVENTORY_TABLE).Find(&inv)
	return inv, result.Error
}

// GetInventoryByParam returns a inventory records based on parameter
func GetInventoryByParam(key, value string) (inv []*m.Inventory, err error) {
	query := db.GetDB().Table(INVENTORY_TABLE)
	query = query.Where(key + " = " + value).Find(&inv)
	return inv, query.Error
}

// DeleteInventoryByParam deletes inventory records for the given parameter
func DeleteInventoryByParam(key, value string) error {
	return db.GetDB().
		Table(INVENTORY_TABLE).
		Where(key + " = " + value).
		Delete(&m.Inventory{}).Error
}

// UpdateColdStoreAvailableQuantity increases the available quantity by 1 for the given compoundID
func UpdateColdStoreAvailableQuantity(compoundID int) error {
	// Get the current available quantity for the given compoundID
	var currentAvailableQuantity int
	query := "SELECT available_quantity FROM " + INVENTORY_TABLE + " WHERE compound_id = ?;"
	err := db.GetDB().Raw(query, compoundID).Scan(&currentAvailableQuantity).Error
	if err != nil {
		return fmt.Errorf("failed to get current available quantity for compound_id %d: %w", compoundID, err)
	}

	// Increment the available quantity by 1
	newAvailableQuantity := currentAvailableQuantity + 1

	// Update the available quantity in the database
	updateQuery := `
    UPDATE ` + INVENTORY_TABLE + ` 
    SET available_quantity = ?, modified_by = ?, modified_on = ?
    WHERE compound_id = ?;`
	modifiedOn := time.Now()
	modifiedBy := "System" // Or use the actual user performing the update

	result := db.GetDB().Exec(updateQuery, newAvailableQuantity, modifiedBy, modifiedOn, compoundID)
	if result.Error != nil {
		return fmt.Errorf("failed to update available quantity for compound_id %d: %w", compoundID, result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("no rows updated for compound_id: %v", compoundID)
	}
	return nil
}

func UpdateAvailableQuantity(Id, DispatchAmount int) error {
	var currentAvailableQuantity int
	query := `
	SELECT available_quantity FROM ` + INVENTORY_TABLE + ` WHERE compound_id = ?;
	`
	err := db.GetDB().Raw(query, Id).Scan(&currentAvailableQuantity).Error
	if err != nil {
		return fmt.Errorf("failed to get current available quantity: %w", err)
	}
	newAvailableQuantity := currentAvailableQuantity - DispatchAmount
	if newAvailableQuantity < 0 {
		return fmt.Errorf("insufficient available quantity for dispatch: %d", currentAvailableQuantity)
	}

	updateQuery := `
		UPDATE ` + INVENTORY_TABLE + ` 
		SET available_quantity = ?, modified_by = ?, modified_on = ?
		WHERE compound_id = ?;`
	modifiedOn := time.Now()
	modifiedBy := "System" // Or use the actual user performing the update

	result := db.GetDB().Exec(updateQuery, newAvailableQuantity, modifiedBy, modifiedOn, Id)
	if result.Error != nil {
		return fmt.Errorf("failed to update available quantity: %w", result.Error)
	}
	// if result.RowsAffected == 0 {
	// 	return fmt.Errorf("no rows updated for compound_id: %v", Id)
	// }

	return nil
}

func UpdateAvailablePartByCompoundId(compoundId int) error {
	db := db.GetDB()

	var currentAvailableQuantity int
	query := "SELECT available_quantity FROM " + INVENTORY_TABLE + " WHERE compound_id = ?;"
	err := db.Raw(query, compoundId).Scan(&currentAvailableQuantity).Error
	if err != nil {
		log.Printf("Failed to get available quantity for compound_id %d: %v", compoundId, err)
		return fmt.Errorf("failed to get available quantity: %w", err)
	}

	if currentAvailableQuantity <= 0 {
		log.Printf("Insufficient available quantity for compound_id %d", compoundId)
		return fmt.Errorf("insufficient available quantity for compound_id: %d", compoundId)
	}

	updateQuery := `
		UPDATE ` + INVENTORY_TABLE + ` 
		SET available_quantity = ?, modified_by = ?, modified_on = ?
		WHERE compound_id = ?;`
	modifiedOn := time.Now()
	modifiedBy := "System"

	result := db.Exec(updateQuery, currentAvailableQuantity-1, modifiedBy, modifiedOn, compoundId)
	if result.Error != nil {
		log.Printf("Failed to update available quantity for compound_id %d: %v", compoundId, result.Error)
		return fmt.Errorf("failed to update available quantity: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		log.Printf("No rows updated for compound_id %d", compoundId)
		return fmt.Errorf("no rows updated for compound_id: %d", compoundId)
	}

	return nil
}

// GetInventoryBySearch returns a inventory records based on condition
func GetInventoryBySearch(con string) (inv []*m.ColdStorage, err error) {
	query := db.GetDB().Table(INVENTORY_TABLE).Select("inventory.*,compounds.compound_name").Joins("join compounds on compounds.id =inventory.compound_id ")
	query = query.Where(con).Find(&inv)
	return inv, query.Error
}

// GetInventoryBySearchPagination returns inventory records based on search and pagination
func GetInventoryBySearchPagination(pagination m.PaginationReq, conditions []string) (inv []*m.ColdStorage, paginationResp m.PaginationResp, err error) {
	// Base query with join
	dbQuery := db.GetDB().Table(INVENTORY_TABLE).Select("inventory.*, compounds.compound_name").
		Joins("JOIN compounds ON compounds.id = inventory.compound_id")

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
	countQuery := db.GetDB().Table(INVENTORY_TABLE).Joins("JOIN compounds ON compounds.id = inventory.compound_id")
	if len(parsedConditions) > 0 {
		countQuery = countQuery.Where(strings.Join(parsedConditions, " AND "))
	}
	if err := countQuery.Count(&totalRecords).Error; err != nil {
		return nil, m.PaginationResp{}, err
	}

	// Set sorting order
	orderBy := "inventory.id DESC"
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
	if err := dbQuery.Find(&inv).Error; err != nil {
		return nil, m.PaginationResp{}, err
	}

	// Prepare pagination response
	paginationResp = m.PaginationResp{
		TotalNo: int(totalRecords),
		Page:    pageNo,
		Offset:  offset,
	}

	return inv, paginationResp, nil
}
