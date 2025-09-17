package model

import "time"

type ProdLineToRecipe struct {
	Id         int       `json:"ID" gorm:"id"`
	RecipeId   int       `json:"RecipeId" gorm:"recipe_id"`
	ProdLineId int       `json:"ProdLineId" gorm:"prod_line_id"`
	Status     bool      `json:"Status" gorm:"status"`
	CreatedBy  int       `json:"CreatedBy" gorm:"created_by"`
	CreatedOn  time.Time `json:"CreatedOn" gorm:"created_on"`
	ModifiedBy int       `json:"ModifiedBy" gorm:"modified_by"`
	ModifiedOn time.Time `json:"ModifiedOn" gorm:"modified_on"`
}
