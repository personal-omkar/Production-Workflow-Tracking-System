package dao

import (
	"strconv"
	"time"

	m "irpl.com/kanban-commons/model"
	db "irpl.com/kanban-dao/db"
)

const (
	ProdLine_To_Recipe_TABLE string = "prod_line_to_recipe" // Updated to match the table name
)

// CreateNewOrUpdateExistingProdLineToRecipe creates a new prod lien to recipe or updates an existing  prod lien to recipe
func CreateNewOrUpdateExistingProdLineToRecipe(rec *m.ProdLineToRecipe) error {

	now := time.Now()
	if rec.Id != 0 {
		rec.ModifiedOn = now
		if err := db.GetDB().Table(ProdLine_To_Recipe_TABLE).Omit("created_by, created_on").Save(&rec).Error; err != nil {
			return err
		}
	} else {
		rec.CreatedOn = now
		if err := db.GetDB().Table(ProdLine_To_Recipe_TABLE).Omit("modified_by, modified_on").Create(&rec).Error; err != nil {
			return err
		}
	}
	return nil
}

// GetProdLineToRecipeByRecipeAndProdId returns prodlinetorecipe records based on the prod and recipe id
func GetProdLineToRecipeByRecipeAndProdId(recipeId, pordid int) (r []m.ProdLineToRecipe, err error) {
	result := db.GetDB().Table(ProdLine_To_Recipe_TABLE).Where("recipe_id=" + strconv.Itoa(recipeId) + " AND prod_line_id=" + strconv.Itoa(pordid)).Find(&r)
	return r, result.Error
}

// GetProductionLineToRecipe returns a user associated with a given id
func GetProductionLineToRecipe(id int) (prodLineToRecipe m.ProdLineToRecipe, err error) {
	result := db.GetDB().Table(ProdLine_To_Recipe_TABLE).Where("id = ?", strconv.Itoa(id)).First(&prodLineToRecipe)
	if result.Error != nil {
		err = result.Error
	}
	return
}
