package model

import "time"

type FinancialYear struct {
	SrNo          int       `gorm:"column:srno;primaryKey" json:"sr_no"`
	FinancialYear string    `gorm:"column:financialyear" json:"financial_year"`
	YearCode      string    `gorm:"column:yearcode" json:"year_code"`
	CreatedOn     time.Time `gorm:"column:createdon" json:"created_on"`
	CreatedBy     string    `gorm:"column:createdby" json:"created_by"`
	ModifiedOn    time.Time `gorm:"column:modifiedon" json:"modified_on"`
	ModifiedBy    string    `gorm:"column:modifiedby" json:"modified_by"`
}
