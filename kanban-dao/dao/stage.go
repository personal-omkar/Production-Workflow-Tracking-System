package dao

import (
	"fmt"
	"strconv"
	"time"

	"github.com/lib/pq"
	m "irpl.com/kanban-commons/model"
	db "irpl.com/kanban-dao/db"
)

const STAGETABLE = "stage"

// Create Stage
func CreateStage(stage m.Stage) error {
	now := time.Now()
	stage.CreatedOn = now
	stage.Active = true
	return db.GetDB().Table(STAGETABLE).Create(&stage).Error
}

func UpdateExistingStage(stage *m.Stage) error {
	now := time.Now()
	stage.ModifiedOn = now
	result := db.GetDB().Table(STAGETABLE).
		Where("id = ?", stage.ID).
		Omit("created_by", "created_on").
		Updates(&stage)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("no record found with id %d", stage.ID)
	}
	return nil
}

// Delete Stage By ID
func DeleteStageByID(id uint) error {
	result := db.GetDB().Table(STAGETABLE).Where("id = ?", id).Delete(&m.Stage{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("no stage found with id %d", id)
	}
	return nil
}

// GetStageByID retrieves a stage record by ID
func GetStageByID(id uint) (stage m.Stage, err error) {
	result := db.GetDB().Table(STAGETABLE).Where("id = ?", id).First(&stage)
	return stage, result.Error
}

// GetAllStage returns a all records present in stage table
func GetAllStage() (stage []m.Stage, err error) {
	result := db.GetDB().Table(STAGETABLE).Order("id").Find(&stage)
	return stage, result.Error
}

func GetStagesByParam(key, value string) ([]m.Stage, error) {
	var stages []m.Stage
	query := db.GetDB().Table(STAGETABLE)
	if key == "id" {
		idVal, err := strconv.Atoi(value)
		if err != nil {
			return nil, fmt.Errorf("invalid id: %v", err)
		}
		query = query.Where("id = ?", idVal)
	} else {
		query = query.Where(fmt.Sprintf("%s = ?", key), value)
	}
	err := query.Find(&stages).Error
	return stages, err
}

func GetStagesByHeader(headers []string) ([]m.Stage, error) {
	var stages []m.Stage

	err := db.GetDB().Table(STAGETABLE).
		Where("headers && ?", pq.Array(headers)). // Postgres array overlap operator
		Find(&stages).Error

	return stages, err
}
