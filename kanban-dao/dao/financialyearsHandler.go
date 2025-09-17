package dao

import (
	"time"

	m "irpl.com/kanban-commons/model"
	db "irpl.com/kanban-dao/db"
)

const (
	FINANCIALYEAR string = "financialyears"
)

// CreateNewOrUpdatefinancialyears creates a new financialyears or updates an existing financialyears
func CreateNewOrUpdatefinancialyears(fy *m.FinancialYear) (ID int, err error) {
	var FinancialYearID int
	now := time.Now()
	if fy.SrNo != 0 {
		fy.ModifiedOn = now

		if err := db.GetDB().Table(FINANCIALYEAR).Save(&fy).Error; err != nil {
			return fy.SrNo, err
		}
	} else {
		fy.CreatedOn = now

		if err := db.GetDB().Table(FINANCIALYEAR).Create(&fy).Error; err != nil {
			return fy.SrNo, err
		}
	}
	FinancialYearID = fy.SrNo
	return FinancialYearID, err
}

// GetFinancialYearByName retrieves a FinancialYear record by financial year name.
func GetFinancialYearByName(yearName string) (m.FinancialYear, error) {
	var financialYear m.FinancialYear
	result := db.GetDB().Table(FINANCIALYEAR).Where("financialyear = ?", yearName).First(&financialYear)
	return financialYear, result.Error
}
