package model

import "time"

type ProdProcess struct {
	Id               int       `json:"ID" gorm:"id"`
	Name             string    `json:"Name" gorm:"name"`
	Link             string    `json:"Link" gorm:"link"`
	Icon             string    `json:"Icon" gorm:"icon"`
	Description      string    `json:"Description" gorm:"description"`
	Status           string    `json:"Status" gorm:"status"`
	LineVisibility   bool      `json:"line_visibility" gorm:"line_visibility"`
	ExpectedMeanTime string    `json:"expected_mean_time" gorm:"expected_mean_time"`
	CreatedBy        string    `json:"CreatedBy" gorm:"created_by"`
	CreatedOn        time.Time `json:"CreatedOn" gorm:"created_on"`
	ModifiedBy       string    `json:"ModifiedBy" gorm:"modified_by"`
	ModifiedOn       time.Time `json:"ModifiedOn" gorm:"modified_on"`
}

type ProdProcessCardData struct {
	Id                int       `json:"ID" gorm:"id"`
	Name              string    `json:"Name" gorm:"name"`
	Link              string    `json:"Link" gorm:"link"`
	Icon              string    `json:"Icon" gorm:"icon"`
	Description       string    `json:"Description" gorm:"description"`
	Status            string    `json:"Status" gorm:"status"`
	CompoundName      string    `json:"CompoundName" gorm:"compound_name"`
	CellNo            string    `json:"CellNo" gorm:"cell_no"`
	KbRootId          int       `json:"KbRootId" gorm:"kb_root_id"`
	Order             int       `json:"Order" gorm:"order"`
	ProdProcessLineId int       `json:"ProdProcessLineId" gorm:"prod_process_line_id"`
	ProdProcessId     int       `json:"ProdProcessId" gorm:"prod_process_id"`
	MFGDateTime       time.Time `json:"MfgDateTime" gorm:"mfg_date_time"`
	IsGroup           bool      `json:"is_group" gorm:"column:isgroup"`
	GroupName         string    `json:"GroupName" gorm:"group_name"`
	CreatedOn         string    `json:"created_on" gorm:"created_on"`
}

type ProdProcessCard struct {
	Id          int                 `json:"ID" gorm:"id"`
	Name        string              `json:"Name" gorm:"name"`
	Link        string              `json:"Link" gorm:"link"`
	Icon        string              `json:"Icon" gorm:"icon"`
	Description string              `json:"Description" gorm:"description"`
	Status      string              `json:"Status" gorm:"status"`
	CardData    ProdProcessCardData `json:"CardData" gorm:"status"`
}

type ProcessOrders struct {
	ProdProcessID string `json:"prod_process_id"`
	Order         int    `json:"order"`
	GroupName     string `json:"group_name"`
}
