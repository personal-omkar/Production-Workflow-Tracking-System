package model

import "time"

type APIKey struct {
	Id         int       `json:"ID" gorm:"column:id"`
	Name       string    `json:"Name" gorm:"column:name"`
	Key        string    `json:"Key" gorm:"column:key"`
	IsActive   bool      `json:"IsActive" gorm:"column:is_active"`
	CreatedBy  string    `json:"CreatedBy" gorm:"column:created_by"`
	CreatedOn  time.Time `json:"CreatedOn" gorm:"column:created_on"`
	ModifiedBy string    `json:"ModifiedBy" gorm:"column:modified_by"`
	ModifiedOn time.Time `json:"ModifiedOn" gorm:"column:modified_on"`
}
