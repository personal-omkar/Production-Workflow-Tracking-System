package model

import "time"

type UserToVendor struct {
	Id         int       `json:"ID" gorm:"id"`
	UserId     int       `json:"UserId" gorm:"user_id"`
	VendorId   int       `json:"VendorId" gorm:"vendor_id"`
	ModifiedOn time.Time `json:"ModifiedOn" gorm:"modified_on"`
	CreatedOn  time.Time `json:"CreatedOn" gorm:"created_on"`
}
