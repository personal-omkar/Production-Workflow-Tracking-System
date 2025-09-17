package dao

import (
	"time"

	m "irpl.com/kanban-commons/model" // Adjust import path for your models
	db "irpl.com/kanban-dao/db"       // Adjust import path for your DB connection
)

const UserToRolesTable = "usertorole" // Define the table name for UserRoless

// CreateUserRoles creates a new UserRoles record
func CreateUserToRole(entry m.UserToRole) error {
	now := time.Now()
	entry.CreatedOn = now
	return db.GetDB().Table(UserToRolesTable).Create(&entry).Error
}

// CreateNewOrUpdateExistingUserToRole creates a new or update existing UserToRoles record
func CreateNewOrUpdateExistingUserToRole(userTorole *m.UserToRole) error {

	now := time.Now()
	if userTorole.ID != 0 {
		userTorole.ModifiedOn = now

		if err := db.GetDB().Table(UserToRolesTable).Save(&userTorole).Error; err != nil {
			return err
		}
	} else {
		userTorole.CreatedOn = now

		if err := db.GetDB().Table(UserToRolesTable).Create(&userTorole).Error; err != nil {
			return err
		}
	}
	return nil
}

// GetUserToRolesByID retrieves a UserToRoles record by ID
func GetUserToRolesByID(id int64) (entry m.UserToRole, err error) {
	result := db.GetDB().Table(UserToRolesTable).Where("id = ?", id).First(&entry)
	return entry, result.Error
}

// UpdateUserToRoles updates an existing UserToRoles record
func UpdateUserToRoles(entry m.UserToRole) error {
	now := time.Now()
	entry.ModifiedOn = now
	return db.GetDB().Table(UserToRolesTable).Where("id = ?", entry.ID).Updates(&entry).Error
}

// GetKbRootByParam returns a kb_root records based on parameter
func GetUserToRolesByParam(key, value string) (usr []m.UserToRole, err error) {
	query := db.GetDB().Table(UserToRolesTable)
	query = query.Where(key + " = " + value).Find(&usr)
	return usr, query.Error
}
