package model

import (
	"time"
)

type Operator struct {
	Id           int       `json:"ID" gorm:"id"`
	OperatorName string    `json:"OperatorName" gorm:"operator_name"`
	OperatorCode string    `json:"OperatorCode" gorm:"operator_code"`
	Status       bool      `json:"Status" gorm:"status"`
	CreatedBy    string    `json:"CreatedBy" gorm:"created_by"`
	CreatedOn    time.Time `json:"CreatedOn" gorm:"created_on"`
	ModifiedBy   string    `json:"ModifiedBy" gorm:"modified_by"`
	ModifiedOn   time.Time `json:"ModifiedOn" gorm:"modified_on"`
}

type CSVOperator struct {
	ID           int       `json:"ID" gorm:"id"`
	OperatorName string    `json:"OperatorName" gorm:"operator_name"`
	OperatorCode string    `json:"OperatorCode" gorm:"operator_code"`
	LineName     string    `json:"LineName" gorm:"line_name"`
	KanbanName   string    `json:"KanbanName" gorm:"kanban_name"`
	Status       bool      `json:"Status" gorm:"column:status"`
	CreatedBy    string    `json:"CreatedBy" gorm:"created_by"`
	CreatedOn    time.Time `json:"CreatedOn" gorm:"created_on"`
	ModifiedBy   string    `json:"ModifiedBy" gorm:"modified_by"`
	ModifiedOn   time.Time `json:"ModifiedOn" gorm:"modified_on"`
}
