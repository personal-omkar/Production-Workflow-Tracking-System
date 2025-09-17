package dao

import (
	m "irpl.com/kanban-commons/model" // Adjust import path for your models
	db "irpl.com/kanban-dao/db"
)

// GetUserDetails returns a user data
func GetUserDetails() (entries []m.UserManagement, err error) {
	if err := db.GetDB().Select("users.id as UserID,users.*, user_roles.role_name, user_roles.id as RoleId, usertorole.id as userRoleId,vendors.vendor_code").Table("users").
		Joins("LEFT OUTER JOIN usertorole  on usertorole.userid =users.id").
		Joins("LEFT OUTER JOIN user_roles  on user_roles.id =usertorole.userroleid").
		Joins("LEFT OUTER JOIN user_to_vendor on user_to_vendor.user_id =users.id").
		Joins("LEFT OUTER JOIN vendors on user_to_vendor.vendor_id =vendors.id").
		Scan(&entries).Error; err != nil {
		return nil, err
	}
	// if err := db.GetDB().Select("users.*").Table("users").Scan(&entries).Error; err != nil {
	// 	return nil, err
	// }
	return entries, nil
}

// GetUserDetails returns a user data
func GetUserDetailsByEmail(email string) (entries []m.UserManagement, err error) {
	if err := db.GetDB().Select("users.id as UserID,users.*, user_roles.role_name, user_roles.id as RoleId, usertorole.id as userRoleId,vendors.vendor_code").Table("users").
		Joins("LEFT OUTER JOIN usertorole  on usertorole.userid =users.id").
		Joins("LEFT OUTER JOIN user_roles  on user_roles.id =usertorole.userroleid").
		Joins("LEFT OUTER JOIN user_to_vendor on user_to_vendor.user_id =users.id").
		Joins("LEFT OUTER JOIN vendors on user_to_vendor.vendor_id =vendors.id").
		Where("users.email=?", email).Scan(&entries).Error; err != nil {
		return nil, err
	}
	// if err := db.GetDB().Select("users.*").Table("users").Scan(&entries).Error; err != nil {
	// 	return nil, err
	// }
	return entries, nil
}
