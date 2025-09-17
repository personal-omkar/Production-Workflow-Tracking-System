package dao

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	m "irpl.com/kanban-commons/model"
	"irpl.com/kanban-commons/utils"
	db "irpl.com/kanban-dao/db"
)

const (
	PRODUCTPROCESSLINE_TABLE string = "prod_process_line" // Updated to match the table name
)

// CreateNewOrUpdateExistingProducProcesstLine creates a new product Prcess line or updates an existing vendor
func CreateNewOrUpdateExistingProducProcesstLine(ppl *m.ProdProcessLine) (ID int, err error) {

	now := time.Now()
	if ppl.Id != 0 {
		ppl.ModifiedOn = now

		if err := db.GetDB().Table(PRODUCTPROCESSLINE_TABLE).Save(&ppl).Error; err != nil {
			return 0, err
		}
	} else {
		ppl.CreatedOn = now

		if err := db.GetDB().Table(PRODUCTPROCESSLINE_TABLE).Create(&ppl).Error; err != nil {
			return 0, err
		}
	}
	return ppl.Id, nil
}

// GetAllProductionLineEntries returns a all records present in prod_line table
func GetAllProducProcesstLineEntries() (entries []m.ProdProcessLine, err error) {
	result := db.GetDB().Table(PRODUCTPROCESSLINE_TABLE).Find(&entries)
	return entries, result.Error
}

// GetKbExtensionsByParam returns a kb_extension records based on parameter
func GetProducProcesstLineByParam(key, value string) (ppl []m.ProdProcessLine, err error) {
	query := db.GetDB().Table(PRODUCTPROCESSLINE_TABLE)
	query = query.Where(key + " = " + value).Find(&ppl)
	return ppl, query.Error
}

// DeleteProductionProcessLineByProdLine delete the production line but there cell data is still present in the kb_data , kb_root, kb_extension
func DeleteProductionProcessLineByProdLine(id int) error {
	return db.GetDB().Table(PRODUCTPROCESSLINE_TABLE).Where("prod_line_id = ?", id).Delete(&m.ProdProcessLine{}).Error
}

func DeleteProductProcessLineByParam(key, value string) (err error) {
	query := db.GetDB().Table(PRODUCTPROCESSLINE_TABLE)
	query = query.Where(key + " = " + value).Delete(&m.ProdProcessLine{})
	return query.Error
}
func GetAllProdProcessesBYOrder(prodLineID int) ([]m.ProdProcessLine, error) {
	var prodProcessLine []m.ProdProcessLine
	result := db.GetDB().
		Table(PRODUCTPROCESSLINE_TABLE).
		Where("prod_line_id = ?", prodLineID).
		Order("\"order\" ASC").
		Find(&prodProcessLine)
	return prodProcessLine, result.Error
}

func CreateProdProcessLines(prodLineID int, processOrders []m.ProcessOrders) error {
	insertValues := ""
	args := []interface{}{}

	for i, processOrder := range processOrders {
		if i > 0 {
			insertValues += ", "
		}
		// Trim spaces from GroupName
		groupName := strings.TrimSpace(processOrder.GroupName)

		// Determine isGroup based on the trimmed GroupName
		isGroup := false
		if groupName != "" {
			isGroup = true
		}

		// Add placeholders for the query
		insertValues += "(?, ?, ?, ?, ?)"

		// Convert ProdProcessID to integer
		intProdProcessID, _ := strconv.Atoi(processOrder.ProdProcessID)

		// Append arguments for the placeholders
		args = append(args, intProdProcessID, prodLineID, processOrder.Order, isGroup, groupName)
	}

	// Include is_group and group_name in the query
	query := `INSERT INTO ` + PRODUCTPROCESSLINE_TABLE + ` (prod_process_id, prod_line_id, "order", isgroup, group_name) VALUES ` + insertValues + `;`

	// Execute the query with the prepared arguments
	if err := db.GetDB().Exec(query, args...).Error; err != nil {
		return err
	}

	return nil
}

func GetProducProcesstLineByParamAndOrder(key, value, order string) (ppl []m.ProdProcessLine, err error) {
	utils.JoinStr()
	query := db.GetDB().Table(PRODUCTPROCESSLINE_TABLE)
	query = query.Where(`"order"= ` + order + "AND " + key + " = " + value).Find(&ppl)
	return ppl, query.Error
}

func GetProducProcesstLineByParams(params map[string]interface{}) ([]m.ProdProcessLine, error) {
	var result []m.ProdProcessLine
	db := db.GetDB()

	// Start building the query
	query := db.Table("prod_process_line").Select("id, prod_line_id, group_name, isgroup, \"order\"")

	// Add WHERE conditions dynamically
	for key, value := range params {
		if key == "order" {
			key = "\"order\"" // Escape the reserved keyword
		}
		query = query.Where(fmt.Sprintf("%s = ?", key), value)
	}

	// Execute the query
	err := query.Find(&result).Error
	if err != nil {
		log.Printf("Error executing query with params %v: %v", params, err)
		return nil, err
	}

	return result, nil
}

func GetGroupedProcessLineIDs(groupID string) ([]string, error) {
	var groupedIDs []string
	err := db.GetDB().
		Table(PRODUCTPROCESSLINE_TABLE).
		Where("group_name= ?", groupID).
		Pluck("id", &groupedIDs).
		Error
	return groupedIDs, err
}

// GetProdProcessByParam returns productionprocess records based on a parameter
func GetProdProcessLineByParam(key, value string) (ppl []m.ProdProcessLine, err error) {
	query := db.GetDB().Table(PRODUCTPROCESSLINE_TABLE)
	if key == "id" {
		id, convErr := strconv.Atoi(value)
		if convErr != nil {
			return nil, fmt.Errorf("invalid value for id: %w", convErr)
		}
		query = query.Where(key+" = ?", id).Find(&ppl)
	} else {
		query = query.Where(key+" = ?", value).Find(&ppl)
	}
	return ppl, query.Error
}
