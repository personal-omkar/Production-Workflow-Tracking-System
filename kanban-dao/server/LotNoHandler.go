package server

import (
	"fmt"
	"strings"
	"time"

	"irpl.com/kanban-dao/dao"
)

// GenerateLotNumber generates a LotNo based on -> current year + production line initial + month code + the number of entries for the current month.
func GenerateLotNumber(ProdLineId string) (string, error) {
	var LotNo string

	// Get current year
	currentYear := time.Now().Year()

	// Get the production line initial
	prodLineData, err := dao.GetProdLineByParam("id", ProdLineId)
	if err != nil || len(prodLineData) == 0 {
		return "", fmt.Errorf("error fetching production line data: %v", err)
	}
	prodLineNameInitial := strings.ToUpper(string(prodLineData[0].Name[0]))

	// Get current month code
	currentMonth := time.Now().Format("January")
	monthMaster, err := dao.GetMonthMasterByMonth(currentMonth)
	if err != nil {
		return "", fmt.Errorf("error fetching month master: %v", err)
	}
	monthCode := monthMaster.MonthCode

	// Get the number of entries for the current month
	// Entries, err := dao.GetEntriesByMonth(currentMonth)
	// if err != nil {
	// 	return "", fmt.Errorf("error fetching entries for the month %s: %v", currentMonth, err)
	// }

	numericCount := prodLineData[0].RunningNumber + 1

	// Format the LotNo without spaces
	LotNo = fmt.Sprintf("%04d%s%s%04d", currentYear, prodLineNameInitial, monthCode, numericCount)
	LotNo = strings.ReplaceAll(strings.TrimSpace(LotNo), " ", "")
	return LotNo, nil
}

// GenerateInventoryLotNumber generates a inventoryLotNo based on -> current year + D (dispatched) + month code + the number of entries for the current month.
func GenerateInventoryLotNumber(num int) (string, error) {
	var LotNo string

	// Get current year
	currentYear := time.Now().Year()
	// Get current month code
	currentMonth := time.Now().Format("January")
	monthMaster, err := dao.GetMonthMasterByMonth(currentMonth)
	if err != nil {
		return "", fmt.Errorf("error fetching month master: %v", err)
	}
	monthCode := monthMaster.MonthCode

	// Get the number of entries for the current month
	Entries, err := dao.GetEntriesByMonth(currentMonth)
	if err != nil {
		return "", fmt.Errorf("error fetching entries for the month %s: %v", currentMonth, err)
	}
	numericCount := Entries + num

	// Format the LotNo without spaces
	LotNo = fmt.Sprintf("%04d%s%s%04d", currentYear, "D", monthCode, numericCount)
	LotNo = strings.ReplaceAll(strings.TrimSpace(LotNo), " ", "")
	return LotNo, nil
}

// Helper function to get the current financial year
func GetCurrentYearName() string {
	now := time.Now()
	if now.Month() >= time.April {
		return fmt.Sprintf("%d-%d", now.Year(), now.Year()+1)
	}
	return fmt.Sprintf("%d-%d", now.Year()-1, now.Year())
}

// Update Lot Number while changing the line
func UpdateLotNumber(lineID, lotNumber string) string {
	yearPart := lotNumber[:4]
	lotPart := lotNumber[4:]
	lineInitial := GetLineInitial(lineID)
	if lineInitial == "" {
		return ""
	}
	// Replace first character of lotPart with the new line initial
	newLotPart := lineInitial + lotPart[1:]
	return strings.TrimSpace(yearPart + newLotPart)
}

func GetLineInitial(LineID string) string {
	LineData, _ := dao.GetProdLineByParam("id", LineID)
	return strings.ToUpper(string(LineData[0].Name[0]))
}
