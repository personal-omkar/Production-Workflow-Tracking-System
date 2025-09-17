package model

import (
	"time"
)

type KbRoot struct {
	Id               int       `json:"ID" gorm:"id"`
	RunningNo        int       `json:"RunningNo" gorm:"running_no"`
	InitialNo        int       `json:"InitialNo" gorm:"initial_no"`
	CreatedBy        string    `json:"CreatedBy" gorm:"created_by"`
	CreatedOn        time.Time `json:"CreatedOn" gorm:"created_on"`
	ModifiedBy       string    `json:"ModifiedBy" gorm:"modified_by"`
	ModifiedOn       time.Time `json:"ModifiedOn" gorm:"modified_on"`
	KbDataId         int       `json:"kbDataId" gorm:"kb_data_id"`
	Status           string    `json:"Status" gorm:"status"` // -->> 0=Kanban state , 1=production process , 2=Qualitytest (-1= quality fail), 3=Packing , 4=Completed
	LotNo            string    `json:"LotNo" gorm:"lot_no"`
	InInventory      bool      `json:"InInventory" gorm:"inInventory"`
	Comment          string    `json:"Comment" gorm:"comment"`
	DispatchNote     string    `json:"DispatchNote" gorm:"dispatch_note"`
	Remark           string    `json:"Remark" gorm:"remark"`
	QualityNote      string    `josn:"QualityNote" gorm:"quality_note"`
	KanbanNo         string    `josn:"KanbanNo" gorm:"kanban_no"`
	QualityDoneTime  time.Time `josn:"quality_done_time" gorm:"quality_done_time"`
	DispatchDoneTime time.Time `josn:"dispatch_done_time" gorm:"dispatch_done_time"`
	QualityOperator  string    `josn:"QualityOperator" gorm:"quality_operator"`
	PackingOperator  string    `josn:"PackingOperator" gorm:"packing_operator"`
}

type ProductionProcess struct {
	ProcessName      string `json:"process_name"`
	StartedOn        string `json:"started_on"`
	CompletedOn      string `json:"completed_on"`
	ExpectedMeanTime string `json:"expected_mean_time"`
	Operator         string `json:"Operator"`
}

type DetailRootData struct {
	VendorName       string          `json:"vendor_name"`
	VendorCode       string          `json:"vendor_code"`
	ContactInfo      string          `json:"contact_info"`
	CompoundName     string          `json:"compound_name"`
	CellNo           string          `json:"cell_no"`
	NoOFLots         int             `json:"NoOFLots" gorm:"no_of_lots"`
	LotNo            string          `json:"lot_no"`
	Status           string          `json:"status"`
	OrderID          string          `json:"order_id"`
	KanbanNo         string          `json:"kanban_no"`
	KanbanStatus     string          `json:"kanban_status"`
	DispatchNote     string          `json:"dispatch_note"`
	QualityNote      string          `josn:"quality_note"`
	QualityDoneTime  string          `josn:"quality_done_time"`
	DispatchDoneTime string          `josn:"dispatch_done_time"`
	QualityOperator  string          `josn:"QualityOperator" gorm:"quality_operator"`
	PackingOperator  string          `josn:"PackingOperator" gorm:"packing_operator"`
	KanbanDetails    []KanbanDetails `json:"kanban_details "`
}

type KanbanDetails struct {
	KanbanName    string              `json:"kanban_name"`
	ProdLine      string              `json:"prod_line"`
	ProdProcesses []ProductionProcess `json:"prod_proces"`
}
