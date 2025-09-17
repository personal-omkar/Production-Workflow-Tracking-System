package dao

import (
	"database/sql"
	"fmt"

	"gorm.io/gorm"
	db "irpl.com/kanban-dao/db"
)

const (
	SYSTEMDEFAULTS string = "systemdefaults"
)

// GetAllLogosByName retrieves a row by system_code and returns it as a map with column names as keys and values as the map values
func GetAllLogosByName(system_code string) (map[string]string, error) {
	result := make(map[string]string)
	rows, err := db.GetDB().Table(SYSTEMDEFAULTS).
		Where("system_code = ?", system_code).
		Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	values := make([]interface{}, len(columns))
	for i := range values {
		// Change the type to sql.NullString for each column
		values[i] = new(sql.NullString)
	}
	for rows.Next() {
		if err := rows.Scan(values...); err != nil {
			return nil, err
		}

		for i, col := range columns {
			// Type assertion to *sql.NullString
			ns, ok := values[i].(*sql.NullString)
			if !ok {
				return nil, fmt.Errorf("unexpected type assertion failure")
			}

			// Check if the value is valid (not NULL)
			if ns.Valid {
				result[col] = ns.String
			} else {
				// Handle the NULL value (set empty string or default value)
				result[col] = ""
			}
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

// UpdateKanbanRunningNumberByCode increments the running_no for a given system_code
func UpdateKanbanRunningNumberByCode(systemCode string) error {
	return db.GetDB().Table("systemdefaults").
		Where("system_code = ?", systemCode).
		UpdateColumn("running_no", gorm.Expr("CAST(CAST(running_no AS INTEGER) + 1 AS VARCHAR)")).Error
}
