package model

import "time"

type Compounds struct {
	Id           int       `json:"ID" gorm:"id"`
	CompoundName string    `json:"CompoundName" gorm:"compound_name"`
	Description  string    `json:"Description" gorm:"description"`
	CreatedBy    string    `json:"CreatedBy" gorm:"created_by"`
	CreatedOn    time.Time `json:"CreatedOn" gorm:"created_on"`
	ModifiedBy   string    `json:"ModifiedBy" gorm:"modified_by"`
	ModifiedOn   time.Time `json:"ModifiedOn" gorm:"modified_on"`
	Status       bool      `json:"Status" gorm:"status"`
	SAPCode      string    `json:"SAPCode" gorm:"column:sap_code"`
	SCADACode    string    `json:"SCADACode" gorm:"column:scada_code"`
}

type AddCompoundsByVendor struct {
	VendorCode   string `json:"VendorCode"`
	VendorName   string `json:"VendorName"`
	CompoundCode int    `json:"CompoundCode"`
	Quantity     int    `json:"Quantity"`
	UserID       string `json:"UserID"`
	Note         string `json:"Note"`
}
type CompoundsDataByVendor struct {
	CompoundCode    int       `json:"CompoundCode"`
	CompoundName    string    `json:"CompoundName"`
	CellNo          string    `json:"CellNo"`
	KbRootId        int       `json:"KbRootId"`
	CreatedOn       time.Time `json:"CreatedOn"`
	ModifiedOn      time.Time `json:"ModifiedOn"`
	DemandDate      time.Time `json:"DemandDate"`
	CustomerNote    string    `json:"customer_note"`
	KanbanNo        string    `json:"kanban_no"`
	QualityDoneTime time.Time `josn:"quality_done_time"`
}

type CSVCompounds struct {
	Id           int       `json:"ID" gorm:"id"`
	CompoundName string    `json:"CompoundName" gorm:"compound_name"`
	Description  string    `json:"Description" gorm:"description"`
	CreatedBy    string    `json:"CreatedBy" gorm:"created_by"`
	CreatedOn    time.Time `json:"CreatedOn" gorm:"created_on"`
	ModifiedBy   string    `json:"ModifiedBy" gorm:"modified_by"`
	ModifiedOn   time.Time `json:"ModifiedOn" gorm:"modified_on"`
	Status       bool      `json:"Status" gorm:"status"`
	SAPCode      string    `json:"SAPCode" gorm:"column:sap_code"`
	SCADACode    string    `json:"SCADACode" gorm:"column:scada_code"`
}
