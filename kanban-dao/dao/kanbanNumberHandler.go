package dao

import (
	"errors"
	"fmt"
	"strconv"
	"time"

)

func GenerateKanbanNumber() (string, error) {
	now := time.Now()
	year := now.Year()
	month := int(now.Month())

	defaults, err := GetAllLogosByName("000001")
	if err != nil {
		return "", fmt.Errorf("unable to get kanban running number")
	}
	runningStr, ok := defaults["running_no"]
	if !ok {
		return "", errors.New("running_no not found in defaults")
	}

	runningNo, err := strconv.Atoi(runningStr)
	if err != nil {
		return "", fmt.Errorf("invalid running_no: %v", err)
	}
	kanban := fmt.Sprintf("RUB/%d/%02d/%04d", year, month, runningNo)
	err = UpdateKanbanRunningNumber()
	if err != nil {
		return "", err
	}
	return kanban, nil
}

func UpdateKanbanRunningNumber() error {
	err := UpdateKanbanRunningNumberByCode("000001")
	if err != nil {
		return fmt.Errorf("fail to update running_no: %v", err)
	}
	return nil

}
