package model

import "time"

type MonthMaster struct {
	SrNo       int       `gorm:"column:srno;primaryKey" json:"sr_no"`
	Month      string    `gorm:"column:month" json:"month"`
	MonthCode  string    `gorm:"column:monthcode" json:"month_code"`
	CreatedOn  time.Time `gorm:"column:createdon" json:"created_on"`
	CreatedBy  string    `gorm:"column:createdby" json:"created_by"`
	ModifiedOn time.Time `gorm:"column:modifiedon" json:"modified_on"`
	ModifiedBy string    `gorm:"column:modifiedby" json:"modified_by"`
}
