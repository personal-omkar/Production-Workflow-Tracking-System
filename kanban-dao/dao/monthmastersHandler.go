package dao

import (
	"time"

	m "irpl.com/kanban-commons/model"
	db "irpl.com/kanban-dao/db"
)

const (
	MONTHMASTER string = "monthmasters"
)

// CreateNewOrUpdatefinancialyears creates a new financialyears or updates an existing financialyears
func CreateNewOrUpdatemonthmasters(mm *m.MonthMaster) (ID int, err error) {
	var FinancialYearID int
	now := time.Now()
	if mm.SrNo != 0 {
		mm.ModifiedOn = now

		if err := db.GetDB().Table(MONTHMASTER).Save(&mm).Error; err != nil {
			return mm.SrNo, err
		}
	} else {
		mm.CreatedOn = now

		if err := db.GetDB().Table(MONTHMASTER).Create(&mm).Error; err != nil {
			return mm.SrNo, err
		}
	}
	FinancialYearID = mm.SrNo
	return FinancialYearID, err
}

// GetMonthMasterByMonth retrieves a MonthMaster record by month name.
func GetMonthMasterByMonth(month string) (m.MonthMaster, error) {
	var monthMaster m.MonthMaster
	result := db.GetDB().Table(MONTHMASTER).Where("month = ?", month).First(&monthMaster)
	return monthMaster, result.Error
}
