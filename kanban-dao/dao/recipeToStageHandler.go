package dao

import (
	"strconv"
	"time"

	m "irpl.com/kanban-commons/model"
	db "irpl.com/kanban-dao/db"
)

const (
	RecipeToStage_TABLE string = "recipetostage" // Updated to match the table name
)

// CreateRecipeToStage creates a new recipe to stage record
func CreateRecipeToStage(rec *m.RecipeToStage) error {

	now := time.Now()

	rec.CreatedOn = now
	if err := db.GetDB().Table(RecipeToStage_TABLE).Omit("modified_by, modified_on").Create(&rec).Error; err != nil {
		return err
	}

	return nil
}

// UpdateRecipeToStage updates a new recipe to stage record
func UpdateRecipeToStage(rec *m.RecipeToStage) error {

	now := time.Now()

	rec.ModifiedOn = now
	if err := db.GetDB().Table(RecipeToStage_TABLE).Omit("created_by, created_on").Save(&rec).Error; err != nil {
		return err
	}

	return nil
}

// GetAllRecipe returns a all records present in recipe table
func GetAllRecipeToStage() (rec []*m.Recipe, err error) {
	result := db.GetDB().Table(RecipeToStage_TABLE).Order("id").Find(&rec)
	return rec, result.Error
}

// GetRecipeToStageByRecipeAndStageId returns recipesTostage records based on the recipe and stage id
func GetRecipeToStageByRecipeAndStageId(recipeId, stageId int) (r []m.RecipeToStage, err error) {
	result := db.GetDB().Table(RecipeToStage_TABLE).Where("recipe_id=" + strconv.Itoa(recipeId) + " AND stage_id=" + strconv.Itoa(stageId)).Find(&r)
	return r, result.Error
}
