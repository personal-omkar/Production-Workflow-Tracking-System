package model

import (
	"time"
)

type SystemLog struct {
	Id          int       `gorm:"column:log_id;primaryKey" json:"Id"`
	Timestamp   time.Time `gorm:"column:timestamp;default:now()" json:"Timestamp"`
	Message     string    `gorm:"column:message;not null" json:"Message"`
	MessageType string    `gorm:"column:message_type;not null" json:"MessageType"`
	IsCritical  bool      `gorm:"column:is_critical;default:false" json:"IsCritical"`
	Icon        string    `gorm:"column:icon" json:"Icon"`
	CreatedBy   string    `gorm:"column:created_by" json:"CreatedBy"`
	CreatedOn   time.Time `gorm:"column:created_on;default:now()" json:"CreatedOn"`
	ModifiedBy  string    `gorm:"column:modified_by" json:"ModifiedBy"`
	ModifiedOn  time.Time `gorm:"column:modified_on;default:now()" json:"ModifiedOn"`
}
