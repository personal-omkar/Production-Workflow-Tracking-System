package model

import (
	"encoding/json"
	"time"
)

type Stage struct {
	ID         uint            `json:"ID" gorm:"primaryKey;autoIncrement;column:id"`
	Name       string          `json:"Name" gorm:"column:name"`
	Headers    json.RawMessage `json:"Headers" gorm:"column:headers"`
	CreatedOn  time.Time       `json:"CreatedOn" gorm:"column:created_on"`
	CreatedBy  string          `json:"CreatedBy" gorm:"column:created_by"`
	ModifiedOn time.Time       `json:"ModifiedOn" gorm:"column:modified_on"`
	ModifiedBy string          `json:"ModifiedBy" gorm:"column:modified_by"`
	Active     bool            `json:"Active" gorm:"type:boolean;column:active"`
}
