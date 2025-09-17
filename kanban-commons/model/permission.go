package model

import "time"

type Permissions struct {
	Id             int       `json:"ID" gorm:"id"`
	PermissionName string    `json:"PermissionName" gorm:"permission_name"`
	Description    string    `json:"Description" gorm:"description"`
	CreatedBy      string    `json:"CreatedBy" gorm:"created_by"`
	CreatedOn      time.Time `json:"CreatedOn" gorm:"created_on"`
	ModifiedBy     string    `json:"ModifiedBy" gorm:"modified_by"`
	ModifiedOn     time.Time `json:"ModifiedOn" gorm:"modified_on"`
}
