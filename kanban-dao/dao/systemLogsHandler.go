package dao

import (
	"time"

	m "irpl.com/kanban-commons/model"
	db "irpl.com/kanban-dao/db"
)

const systemLogTable = "systemlogs"

// CreateSystemLog creates a new SystemLog record
func CreateSystemLog(entry m.SystemLog) error {
	now := time.Now()
	entry.Timestamp = now
	entry.CreatedOn = now
	entry.ModifiedOn = now

	return db.GetDB().Table(systemLogTable).Create(&entry).Error
}

// GetSystemLogByID retrieves a SystemLog record by ID
func GetSystemLogByID(id int) (entry m.SystemLog, err error) {
	result := db.GetDB().Table(systemLogTable).Where("log_id = ?", id).First(&entry)
	return entry, result.Error
}

// UpdateSystemLog updates an existing SystemLog record
func UpdateSystemLog(entry m.SystemLog) error {
	now := time.Now()
	entry.ModifiedOn = now
	return db.GetDB().Table(systemLogTable).Where("log_id = ?", entry.Id).
		Select("Message", "MessageType", "IsCritical", "Icon", "ModifiedBy", "ModifiedOn").
		Updates(&entry).Error
}

// DeleteSystemLog deletes a SystemLog record by ID
func DeleteSystemLog(id int) error {
	return db.GetDB().Table(systemLogTable).Where("log_id = ?", id).Delete(&m.SystemLog{}).Error
}

// GetSystemLogs returns a paginated list of SystemLog records
func GetSystemLogs(page, limit int) (entries []m.SystemLog, totalRecords int, err error) {
	offset := (page - 1) * limit
	var totalRecords64 int64

	// Get total record count
	err = db.GetDB().Table(systemLogTable).Count(&totalRecords64).Error
	if err != nil {
		return nil, 0, err
	}

	// Convert totalRecords64 to int
	totalRecords = int(totalRecords64)

	// Fetch records with pagination
	err = db.GetDB().Table(systemLogTable).Offset(offset).Limit(limit).Find(&entries).Error
	if err != nil {
		return nil, 0, err
	}

	return entries, totalRecords, nil
}

// GetSystemLogsByCriteria returns a paginated list of SystemLog records based on criteria
func GetSystemLogsByCriteria(page, limit int, criteria map[string]interface{}) (entries []m.SystemLog, totalRecords int, err error) {
	offset := (page - 1) * limit
	var totalRecords64 int64
	query := db.GetDB().Table(systemLogTable)

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

// GetAllSystemLogs returns all records present in the systemlogs table
func GetAllSystemLogs() (entries []m.SystemLog, err error) {
	result := db.GetDB().Table(systemLogTable).Find(&entries)
	return entries, result.Error
}
