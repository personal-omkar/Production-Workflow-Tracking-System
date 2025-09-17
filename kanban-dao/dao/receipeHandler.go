package dao

import (
	"fmt"
	"strconv"
	"time"

	m "irpl.com/kanban-commons/model"
	db "irpl.com/kanban-dao/db"
)

const (
	Recipe_TABLE string = "recipe" // Updated to match the table name
)

// CreateNewOrUpdateExistingRecipe creates a new recipe or updates an existing recipe
func CreateNewOrUpdateExistingRecipe(rec *m.Recipe) error {

	now := time.Now()
	if rec.Id != 0 {
		rec.ModifiedOn = now
		if err := db.GetDB().Table(Recipe_TABLE).Omit("created_by, created_on").Save(&rec).Error; err != nil {
			return err
		}
	} else {
		rec.CreatedOn = now
		if err := db.GetDB().Table(Recipe_TABLE).Omit("modified_by, modified_on").Create(&rec).Error; err != nil {
			return err
		}
	}
	return nil
}

// GetAllRecipe returns a all records present in recipe table
func GetAllRecipe() (rec []*m.Recipe, err error) {
	result := db.GetDB().Table(Recipe_TABLE).Order("id").Find(&rec)
	return rec, result.Error
}

// GetRecipeByDataKey returns a all records present in recipe table
func GetRecipeByDataKey(key string) ([]*m.Recipe, error) {
	var rec []*m.Recipe

	query := `
	SELECT DISTINCT ON (recipe.id) recipe.*
	FROM recipe,
	     jsonb_each_text(data) AS kv(key, value)
	WHERE kv.key = ?
	`

	result := db.GetDB().Raw(query, key).Scan(&rec)
	return rec, result.Error
}

// GetRecipeByDataValue returns a all records present in recipe table
func GetRecipeByDataValue(value string) ([]*m.Recipe, error) {
	var rec []*m.Recipe

	query := `
	SELECT DISTINCT ON (recipe.id) recipe.*
	FROM recipe,
	     jsonb_each_text(data) AS kv(key, value)
	WHERE kv.value = ?
	`

	result := db.GetDB().Raw(query, value).Scan(&rec)
	return rec, result.Error
}

// GetRecipeByDataKeyAndValue returns a all records present in recipe table
func GetRecipeByDataKeyAndValue(key, value string) ([]*m.Recipe, error) {
	var rec []*m.Recipe

	query := `
	SELECT DISTINCT ON (recipe.id) recipe.*
	FROM recipe,
	     jsonb_each_text(data) AS kv(key, value)
	WHERE kv.key = ? AND  kv.value = ?
	`

	result := db.GetDB().Raw(query, key, value).Scan(&rec)
	return rec, result.Error
}

// DeleteUser deletes a User record by SrNo
func DeleteRecipeById(id int) error {
	return db.GetDB().Table(Recipe_TABLE).Where("id = ?", id).Delete(&m.Recipe{}).Error
}

// GetRecipeByParam returns a all records present in recipe table based on the condition
func GetRecipeByParam(key, value string) ([]m.Recipe, error) {
	var recipe []m.Recipe
	query := db.GetDB().Table(Recipe_TABLE)
	if key == "id" {
		idVal, err := strconv.Atoi(value)
		if err != nil {
			return nil, fmt.Errorf("invalid id: %v", err)
		}
		query = query.Where("id = ?", idVal)
	} else {
		query = query.Where(fmt.Sprintf("%s = ?", key), value)
	}
	err := query.Find(&recipe).Error
	return recipe, err
}
