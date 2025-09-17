package dao

import (
	"time"

	m "irpl.com/kanban-commons/model"
	db "irpl.com/kanban-dao/db"
)

const sambaConfigTable = "sambaconfig"

// CreateSambaConfig creates a new SambaConfig record
func CreateSambaConfig(entry m.SambaConfig) error {
	now := time.Now()
	entry.CreatedOn = now
	return db.GetDB().Omit("id").Table(sambaConfigTable).Create(&entry).Error
}

// GetSambaConfigByID retrieves a SambaConfig record by ID
func GetSambaConfigByID(id int) (entry m.SambaConfig, err error) {
	result := db.GetDB().Table(sambaConfigTable).Where("id = ?", id).First(&entry)
	return entry, result.Error
}

// UpdateSambaConfig updates an existing SambaConfig record
func UpdateSambaConfig(entry m.SambaConfig) error {
	now := time.Now()
	entry.ModifiedOn = now
	return db.GetDB().Table(sambaConfigTable).Where("id = ?", entry.ID).Updates(&entry).Error
}

// DeleteSambaConfig deletes a SambaConfig record by ID
func DeleteSambaConfig(id int) error {
	return db.GetDB().Table(sambaConfigTable).Where("id = ?", id).Delete(&m.SambaConfig{}).Error
}

// GetDefaultSambaConfig retrieves the default SambaConfig record
func GetDefaultSambaConfig() (entry m.SambaConfig, err error) {
	result := db.GetDB().Table(sambaConfigTable).Where("is_default = ?", true).Limit(1).First(&entry)
	return entry, result.Error
}
