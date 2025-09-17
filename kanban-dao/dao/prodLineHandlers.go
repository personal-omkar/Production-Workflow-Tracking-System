package dao

import (
	"strconv"
	"strings"
	"time"

	"fmt"

	m "irpl.com/kanban-commons/model"
	db "irpl.com/kanban-dao/db"
)

const (
	PRODUCTLINE_TABLE string = "prod_line" // Updated to match the table name
)

// CreateNewOrUpdateExistingProductLine creates a new product line or updates an existing vendor
func CreateNewOrUpdateExistingProductLine(pl *m.ProdLine) error {

	now := time.Now()
	if pl.Id != 0 {
		pl.ModifiedOn = now

		if err := db.GetDB().Table(PRODUCTLINE_TABLE).Save(&pl).Error; err != nil {
			return err
		}
	} else {
		pl.CreatedOn = now

		if err := db.GetDB().Table(PRODUCTLINE_TABLE).Create(&pl).Error; err != nil {
			return err
		}
	}
	return nil
}

// GetAllProductionLineEntries returns a all records present in prod_line table
func GetAllProductionLineEntries() (entries []*m.ProdLine, err error) {
	result := db.GetDB().Table(PRODUCTLINE_TABLE).Order("id").Find(&entries)
	return entries, result.Error
}

// GetKBDataByParam returns kb_data records based on a parameter
func GetProdLineByParam(key, value string) (kbe []m.ProdLine, err error) {
	query := db.GetDB().Table(PRODUCTLINE_TABLE)
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

// Func to create a Prod line at prod_line and return the ID of the newly created record
func CreateProdLine(prod_line *m.ProdLine) (int, error) {
	// Save the new prod_line to the database
	if err := db.GetDB().Table(PRODUCTLINE_TABLE).Create(&prod_line).Error; err != nil {
		return 0, err
	}

	// Return the ID of the newly created prod_line
	return prod_line.Id, nil
}

// DeleteProductionLine delete the production line but there cell data is still present in the kb_data , kb_root, kb_extension
func DeleteProductionLine(id int) error {
	return db.GetDB().Table(PRODUCTLINE_TABLE).Where("id = ?", id).Delete(&m.ProdLine{}).Error
}

// GetProdLineDetails fetches production line details and returns them as a slice
func GetProdLineDetails() ([]m.ProdLineDetails, error) {
	var prodLines []m.ProdLineDetails

	// Query to fetch production line ID and name
	query := `
    SELECT 
        id AS prod_id,
        name AS prod_name
    FROM 
        prod_line
	WHERE
		prod_line.Status= TRUE;
    `

	rows, err := db.GetDB().Raw(query).Rows()
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	// Populate the slice with ProdLineDetails
	for rows.Next() {
		var prodID int
		var prodName string

		err := rows.Scan(&prodID, &prodName)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row data: %v", err)
		}

		prodLines = append(prodLines, m.ProdLineDetails{
			ProdLineID:   prodID,
			ProdLineName: prodName,
			Cells:        []m.Cell{}, // Start with an empty list of Cells
		})
	}

	return prodLines, nil
}

// GetProdLineCells fetches production line details along with cell data
func GetProdLinesWithCellsAndStatus() ([]m.ProdLineDetails, error) {
	// Fetch the initial prodLines slice
	prodLines, err := GetProdLineDetails()
	if err != nil {
		return nil, fmt.Errorf("failed to get prod lines: %v", err)
	}

	// Query to fetch the cell details based on the kb_root_id
	query := `
    SELECT 
        pl.id AS ProdLineID, 
        pl.Name AS prod_line_name,
        kt.kb_root_id,
        ppl.prod_process_id,
        kd.cell_no AS CellNumber,
        kr.running_no AS KBRunningNo,
        kr.initial_no AS KBInitialNo,
		kr.kanban_no AS KanbanNo,
        c.compound_name AS CompoundName,
        kd.mfg_date_time AS MfgDateTime,
        kd.demand_date_time AS DemandDateTime,
        kd.exp_date AS ExpDate,
        kd.no_of_lots AS NoOFLots,
        kd.location AS Location,
        ke.status AS Status,
		kr.lot_no as lot_no
    FROM 
        kb_transaction kt
    JOIN 
        kb_root kr ON kt.kb_root_id = kr.id
    JOIN 
        kb_data kd ON kr.kb_data_id = kd.id
    JOIN 
        kb_extension ke ON kd.kb_extension_id = ke.id
    JOIN 
        compounds c ON kd.compound_id = c.id
    JOIN
        prod_process_line ppl ON kt.prod_process_line_id = ppl.id
    JOIN
        prod_line pl ON ppl.prod_line_id = pl.id
	WHERE 
    	pl.Status = TRUE 
    	AND kt.kb_root_id IN (
			SELECT kb_root_id
            FROM kb_transaction
            GROUP BY kb_root_id
            HAVING COUNT(kb_root_id) = 1
        );
    `

	rows, err := db.GetDB().Raw(query).Rows()
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	// Iterate through the rows and append the cell data to the corresponding production line in the slice
	for rows.Next() {
		var prodLineID int
		var prodLineName string
		var cell m.Cell

		err := rows.Scan(
			&prodLineID,
			&prodLineName,
			&cell.KRId,
			&cell.ProdProcessID,
			&cell.CellNumber,
			&cell.KBRunningNo,
			&cell.KBInitialNo,
			&cell.KanbanNo,
			&cell.CompoundName,
			&cell.MfgDateTime,
			&cell.DemandDateTime,
			&cell.ExpDate,
			&cell.NoOFLots,
			&cell.Location,
			&cell.Status,
			&cell.LotNo,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row data: %v", err)
		}

		// Find the ProdLineDetails corresponding to the ProdLineID in the slice
		for i := range prodLines {
			if prodLines[i].ProdLineID == prodLineID {
				prodLines[i].Cells = append(prodLines[i].Cells, cell)
				break
			}
		}
	}
	return prodLines, nil
}

// Function to fetch production line details and the latest record for a given prod_id
func GetLatestRecordsForProdID(prodID int) ([]m.ProdLineDetails, error) {
	// Query to fetch production line details along with cell data based on the prod_id
	query := `
   WITH LatestRecords AS (
    SELECT 
        kt.kb_root_id, 
        MAX(kt.created_on) AS LatestCreatedOn
    FROM 
        kb_transaction kt
    GROUP BY 
        kt.kb_root_id
)
SELECT 
    pl.id AS ProdLineID,
    pl.Name AS ProdLineName,
    kt.kb_root_id, 
    ppl.prod_process_id, 
	ppl."order" as prod_process_line_order, 
    kt.created_on,
    kr.running_no AS KBRunningNo,
    kr.initial_no AS KBInitialNo,
	kr.kanban_no ,
    kd.cell_no AS CellNumber,
    c.compound_name AS CompoundName,
    kd.mfg_date_time AS MfgDateTime,
    kd.demand_date_time AS DemandDateTime,
    kd.exp_date AS ExpDate,
    kd.no_of_lots AS NoOFLots,
    kd.location AS Location,
    ke.status AS Status,
	kr.lot_no AS lot_no
FROM 
    kb_transaction kt
JOIN 
    LatestRecords lr 
    ON kt.kb_root_id = lr.kb_root_id AND kt.created_on = lr.LatestCreatedOn
JOIN 
    kb_root kr ON kt.kb_root_id = kr.id
JOIN 
    kb_data kd ON kr.kb_data_id = kd.id
JOIN 
    kb_extension ke ON kd.kb_extension_id = ke.id
JOIN 
    compounds c ON kd.compound_id = c.id
JOIN
    prod_process_line ppl ON kt.prod_process_line_id = ppl.id
JOIN
    prod_line pl ON ppl.prod_line_id = pl.id
WHERE 
    pl.id = ?
ORDER BY 
    kt.created_on DESC;
    `

	// Declare a slice to hold the result
	var prodLineDetails []m.ProdLineDetails

	// Execute the query
	rows, err := db.GetDB().Raw(query, prodID).Rows()
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	// Declare variables to hold ProdLineDetails and Cells
	var lastProdLineID int
	var currentProdLine m.ProdLineDetails

	// Iterate through the rows and populate the result
	for rows.Next() {
		var cell m.Cell
		var createdOn time.Time

		// Scan the row into the variables
		err := rows.Scan(
			&currentProdLine.ProdLineID,
			&currentProdLine.ProdLineName,
			&cell.KRId,
			&cell.ProdProcessID,
			&cell.ProductionProcessLineOrder,
			&createdOn,
			&cell.KBRunningNo,
			&cell.KBInitialNo,
			&cell.KanbanNo,
			&cell.CellNumber,
			&cell.CompoundName,
			&cell.MfgDateTime,
			&cell.DemandDateTime,
			&cell.ExpDate,
			&cell.NoOFLots,
			&cell.Location,
			&cell.Status,
			&cell.LotNo,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row data: %v", err)
		}

		// Format the createdOn date
		cell.MfgDateTime = createdOn.Format("2006-01-02 15:04:05")

		// If the current production line has changed or it's the first entry, save the current data
		if currentProdLine.ProdLineID != lastProdLineID {
			// If this isn't the first row, append the last ProdLineDetails to the slice
			if lastProdLineID != 0 {
				prodLineDetails = append(prodLineDetails, currentProdLine)
			}

			// Reset currentProdLine for the new ProdLineID
			currentProdLine = m.ProdLineDetails{
				ProdLineID:   currentProdLine.ProdLineID,
				ProdLineName: currentProdLine.ProdLineName,
				Cells:        []m.Cell{cell}, // Start with the first cell
			}
			lastProdLineID = currentProdLine.ProdLineID // Update lastProdLineID to the new one
		} else {
			// If it's the same ProdLineID, just append the cell data to the existing currentProdLine
			currentProdLine.Cells = append(currentProdLine.Cells, cell)
		}
	}

	// Don't forget to append the last ProdLineDetails to the slice
	if len(currentProdLine.Cells) > 0 {
		prodLineDetails = append(prodLineDetails, currentProdLine)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate over rows: %v", err)
	}

	// If no data found, fetch the basic production line info
	if len(prodLineDetails) == 0 {
		var prodLines []m.ProdLineDetails
		// Query to fetch production line ID and name
		query := `
		SELECT 
			id AS prod_id,
			name AS prod_name
		FROM 
			prod_line
		WHERE 
			id = ?;
		`
		rows, err := db.GetDB().Raw(query, prodID).Rows()
		if err != nil {
			return nil, fmt.Errorf("failed to execute query: %v", err)
		}
		defer rows.Close()
		// Populate the slice with ProdLineDetails
		for rows.Next() {
			var prodID int
			var prodName string
			err := rows.Scan(&prodID, &prodName)
			if err != nil {
				return nil, fmt.Errorf("failed to scan row data: %v", err)
			}
			prodLines = append(prodLines, m.ProdLineDetails{
				ProdLineID:   prodID,
				ProdLineName: prodName,
				Cells:        []m.Cell{}, // Start with an empty list of Cells
			})
		}

		return prodLines, nil
	}

	// Return the populated production line details
	return prodLineDetails, nil
}

// GetUniqueTransactionKBRoot fetches a slice of unique kb_root_id for a given production line ID.
func GetUniqueTransactionKBRoot(prodLineID int) ([]int, error) {
	// Validate input
	if prodLineID <= 0 {
		return nil, fmt.Errorf("invalid production line ID provided")
	}

	// Prepare the query to fetch unique kb_root_id
	query := `
    SELECT DISTINCT 
        kt.kb_root_id AS KBRootID
    FROM 
        kb_transaction kt
    JOIN 
        kb_root kr ON kt.kb_root_id = kr.id
    JOIN
        prod_process_line ppl ON kt.prod_process_line_id = ppl.id
    JOIN
        prod_line pl ON ppl.prod_line_id = pl.id
    WHERE 
        pl.id = ?
    GROUP BY 
        kt.kb_root_id
    HAVING 
        COUNT(kt.kb_root_id) = 1;
    `

	// Execute the query with the production line ID
	rows, err := db.GetDB().Raw(query, prodLineID).Rows()
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	// Slice to hold the kb_root_id results
	var kbRootIDs []int

	// Iterate through the rows
	for rows.Next() {
		var kbRootID int

		err := rows.Scan(&kbRootID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row data: %v", err)
		}

		// Append the kb_root_id to the slice
		kbRootIDs = append(kbRootIDs, kbRootID)
	}

	return kbRootIDs, nil
}

// Get Production Line ID By KbRootID
func GetProdLineIDByKbRootID(kbRootID int) (int, error) {
	var prodLineID int
	err := db.GetDB().
		Table(PRODUCTLINE_TABLE+" p").
		Select("p.id").
		Joins("JOIN "+KBTRANSACTIONDATA_TABLE+" t ON t.kb_root_id = ?", kbRootID).
		Joins("JOIN " + PRODUCTPROCESSLINE_TABLE + " ppl ON ppl.id = t.prod_process_line_id").
		Where("ppl.prod_line_id = p.id").
		Scan(&prodLineID).Error
	if err != nil {
		return 0, fmt.Errorf("error fetching prod_line_id for kb_root_id %d: %v", kbRootID, err)
	}
	return prodLineID, nil
}

func GetLinedUpKBRootsByProdLineID(prodLineID int) (KbRootData []m.KbRoot, err error) {
	query := `
		SELECT kr.id, kr.running_no, kr.initial_no, kr.created_by, kr.created_on, 
			kr.modified_by, kr.modified_on, kr.kb_data_id, kr.status, kr.lot_no
		FROM kb_root kr 
		JOIN kb_transaction kt ON kt.kb_root_id = kr.id
		JOIN prod_process_line ppl ON kt.prod_process_line_id = ppl.id
		JOIN prod_line pl ON ppl.prod_line_id = pl.id
		WHERE pl.id = ?
		GROUP BY kr.id
		HAVING COUNT(kt.kb_root_id) = 1;
	`
	// Execute the query
	rows, err := db.GetDB().Raw(query, prodLineID).Rows()
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var kbRoot m.KbRoot
		if err := rows.Scan(&kbRoot.Id, &kbRoot.RunningNo, &kbRoot.InitialNo,
			&kbRoot.CreatedBy, &kbRoot.CreatedOn, &kbRoot.ModifiedBy,
			&kbRoot.ModifiedOn, &kbRoot.KbDataId, &kbRoot.Status, &kbRoot.LotNo); err != nil {
			return nil, fmt.Errorf("failed to scan kb_root data: %v", err)
		}
		KbRootData = append(KbRootData, kbRoot)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %v", err)
	}

	if len(KbRootData) == 0 {
		return nil, fmt.Errorf("no kb_root records found for prodLineID %d", prodLineID)
	}

	return KbRootData, nil
}

// GetProdLinesWithCellsAndStatusByID fetches production line details along with cell data for a specific production line ID
func GetProdLinesWithCellsAndStatusByID(prodLineID int) (*m.ProdLineDetails, error) {
	var prodLine m.ProdLineDetails

	// Query to fetch production line details
	query := `
    SELECT 
        pl.id AS ProdLineID, 
        pl.Name AS prod_line_name,
        kt.kb_root_id,
        ppl.prod_process_id,
        kd.cell_no AS CellNumber,
        kr.running_no AS KBRunningNo,
        kr.initial_no AS KBInitialNo,
		kr.kanban_no AS KanbanNo,
        c.compound_name AS CompoundName,
        kd.mfg_date_time AS MfgDateTime,
        kd.demand_date_time AS DemandDateTime,
        kd.exp_date AS ExpDate,
        kd.no_of_lots AS NoOFLots,
        kd.location AS Location,
        ke.status AS Status,
		kr.lot_no AS LotNo
    FROM 
        kb_transaction kt
    JOIN 
        kb_root kr ON kt.kb_root_id = kr.id
    JOIN 
        kb_data kd ON kr.kb_data_id = kd.id
    JOIN 
        kb_extension ke ON kd.kb_extension_id = ke.id
    JOIN 
        compounds c ON kd.compound_id = c.id
    JOIN
        prod_process_line ppl ON kt.prod_process_line_id = ppl.id
    JOIN
        prod_line pl ON ppl.prod_line_id = pl.id
	WHERE 
        pl.id = ? 
        AND kt.kb_root_id IN (
            SELECT kb_root_id
            FROM kb_transaction
            GROUP BY kb_root_id
            HAVING COUNT(kb_root_id) = 1
        );
    `

	rows, err := db.GetDB().Raw(query, prodLineID).Rows()
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	// Initialize Cells slice
	prodLine.Cells = []m.Cell{}

	// Iterate through the rows and populate the prodLine struct
	for rows.Next() {
		var cell m.Cell

		err := rows.Scan(
			&prodLine.ProdLineID,
			&prodLine.ProdLineName,
			&cell.KRId,
			&cell.ProdProcessID,
			&cell.CellNumber,
			&cell.KBRunningNo,
			&cell.KBInitialNo,
			&cell.KanbanNo,
			&cell.CompoundName,
			&cell.MfgDateTime,
			&cell.DemandDateTime,
			&cell.ExpDate,
			&cell.NoOFLots,
			&cell.Location,
			&cell.Status,
			&cell.LotNo,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row data: %v", err)
		}

		// Append cell data
		prodLine.Cells = append(prodLine.Cells, cell)
	}

	return &prodLine, nil
}

func GetAllProdLineBySearchAndPagination(pagination m.PaginationReq, conditions []string) (op []*m.ProdLine, paginationResp m.PaginationResp, err error) {
	// Base query
	dbQuery := db.GetDB().Table(PRODUCTLINE_TABLE)

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
	countQuery := db.GetDB().Table(PRODUCTLINE_TABLE)
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
