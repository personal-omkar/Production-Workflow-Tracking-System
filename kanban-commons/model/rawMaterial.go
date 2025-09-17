package model

import "time"

type RawMaterial struct {
	Id          int       `json:"ID" gorm:"column:id;primaryKey;autoIncrement"`
	Status      bool      `json:"Status" gorm:"column:status"`
	CreatedBy   string    `json:"CreatedBy" gorm:"column:created_by"`
	CreatedOn   time.Time `json:"CreatedOn" gorm:"column:created_on"`
	ModifiedBy  string    `json:"ModifiedBy" gorm:"column:modified_by"`
	ModifiedOn  time.Time `json:"ModifiedOn" gorm:"column:modified_on"`
	Description string    `json:"Description" gorm:"column:material_desc"`
	SCADACode   string    `json:"SCADACode" gorm:"column:scada_code"`
	SAPCode     string    `json:"SAPCode" gorm:"column:sap_code"`
	Comment     string    `json:"Comment" gorm:"column:comment"`
}

type CSVMaterial struct {
	Id          int       `json:"ID" gorm:"column:id;primaryKey;autoIncrement"`
	Status      bool      `json:"Status" gorm:"column:status"`
	CreatedBy   string    `json:"CreatedBy" gorm:"column:created_by"`
	CreatedOn   time.Time `json:"CreatedOn" gorm:"column:created_on"`
	ModifiedBy  string    `json:"ModifiedBy" gorm:"column:modified_by"`
	ModifiedOn  time.Time `json:"ModifiedOn" gorm:"column:modified_on"`
	Description string    `json:"Description" gorm:"column:material_desc"`
	SCADACode   string    `json:"SCADACode" gorm:"column:scada_code"`
	SAPCode     string    `json:"SAPCode" gorm:"column:sap_code"`
	Comment     string    `json:"Comment" gorm:"column:comment"`
}
