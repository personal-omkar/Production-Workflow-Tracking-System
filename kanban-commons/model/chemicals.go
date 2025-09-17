package model

import "time"

type ChemicalTypes struct {
	Id         int       `json:"ID" gorm:"column:id;primaryKey;autoIncrement"`
	Type       string    `json:"Type" gorm:"column:type"`
	ConvCode   string    `json:"ConvCode" gorm:"column:conv_code"`
	Status     bool      `json:"Status" gorm:"status"`
	CreatedBy  string    `json:"CreatedBy" gorm:"column:created_by"`
	CreatedOn  time.Time `json:"CreatedOn" gorm:"column:created_on"`
	ModifiedBy string    `json:"ModifiedBy" gorm:"column:modified_by"`
	ModifiedOn time.Time `json:"ModifiedOn" gorm:"column:modified_on"`
}

type CSVChemical struct {
	Id         int       `json:"ID" gorm:"column:id;primaryKey;autoIncrement"`
	Type       string    `json:"Type" gorm:"column:type"`
	ConvCode   string    `json:"ConvCode" gorm:"column:conv_code"`
	Status     bool      `json:"Status" gorm:"status"`
	CreatedBy  string    `json:"CreatedBy" gorm:"column:created_by"`
	CreatedOn  time.Time `json:"CreatedOn" gorm:"column:created_on"`
	ModifiedBy string    `json:"ModifiedBy" gorm:"column:modified_by"`
	ModifiedOn time.Time `json:"ModifiedOn" gorm:"column:modified_on"`
}
