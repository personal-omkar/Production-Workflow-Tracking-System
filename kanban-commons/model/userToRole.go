package model

import "time"

type UserToRole struct {
	ID         int       `gorm:"column:id;primaryKey" json:"ID"`
	UserId     int       `gorm:"column:userid" json:"UserID"`
	UserRoleID int       `gorm:"column:userroleid" json:"UserRoleId"`
	CreatedOn  time.Time `gorm:"column:createdon" json:"created_on"`
	CreatedBy  string    `gorm:"column:createdby" json:"created_by"`
	ModifiedOn time.Time `gorm:"column:modifiedon" json:"modified_on"`
	ModifiedBy string    `gorm:"column:modifiedby" json:"modified_by"`
}
