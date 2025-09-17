package dao

import (
	"time"

	m "irpl.com/kanban-commons/model" // Adjust import path for your models
	db "irpl.com/kanban-dao/db"       // Adjust import path for your DB connection
)

const (
	USERTOVENDOR_TABLE string = "user_to_vendor" // Updated to match the table name
)

// GetUserToVendorByParam returns a usertovendor records based on parameter
func GetUserToVendorByParam(key, value string) (UsrToVen []m.UserToVendor, err error) {
	query := db.GetDB().Table(USERTOVENDOR_TABLE)
	query = query.Where(key + " = " + value).Find(&UsrToVen)
	return UsrToVen, query.Error
}

// CreateNewOrUpdateVendorToUser creates a new user or updates an existing user
func CreateNewOrUpdateVendorToUser(UserToVendor *m.UserToVendor) error {
	now := time.Now()

	if UserToVendor.Id != 0 {
		// Update record
		UserToVendor.ModifiedOn = now // Update the modified timestamp
		if err := db.GetDB().Table(USERTOVENDOR_TABLE).
			Omit("created_on").
			Save(&UserToVendor).Error; err != nil {
			return err
		}
	} else {
		// Create new record
		UserToVendor.CreatedOn = now
		// Exclude ModifiedOn if it's not set
		if err := db.GetDB().Table(USERTOVENDOR_TABLE).
			Omit("modified_on").
			Create(&UserToVendor).Error; err != nil {
			return err
		}
	}
	return nil
}

func GetVendorDetailsByUserID(id int) (vendor m.Vendors, err error) {
	result := db.GetDB().Table("users").Select("vendors.*").Joins("join user_to_vendor  on user_to_vendor.user_id =users.id").Joins("join vendors  on vendors.id =user_to_vendor.vendor_id").Where("users.id = ?", id).First(&vendor)
	return vendor, result.Error
}
