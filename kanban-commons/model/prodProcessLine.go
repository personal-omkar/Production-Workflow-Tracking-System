package model

import "time"

type ProdProcessLine struct {
	Id            int       `json:"ID" gorm:"id"`
	ProdProcessID int       `json:"ProdProcessID" gorm:"prod_process_id"`
	ProdLineId    int       `json:"ProdLineId" gorm:"prod_line_id"`
	Order         int       `json:"Order" gorm:"order"`
	CreatedBy     string    `json:"CreatedBy" gorm:"created_by"`
	CreatedOn     time.Time `json:"CreatedOn" gorm:"created_on"`
	ModifiedBy    string    `json:"ModifiedBy" gorm:"modified_by"`
	ModifiedOn    time.Time `json:"ModifiedOn" gorm:"modified_on"`
	IsGroup       bool      `json:"is_group" gorm:"column:isgroup"`
	GroupName     string    `json:"GroupName" gorm:"group_name"`
}
