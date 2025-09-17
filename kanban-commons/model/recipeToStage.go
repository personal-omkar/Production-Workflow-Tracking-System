package model

import (
	"encoding/json"
	"time"
)

type RecipeToStage struct {
	Id         int             `json:"ID" gorm:"id"`
	RecipeId   int             `json:"RecipeId" gorm:"recipe_id"`
	StageId    int             `json:"StageId" gorm:"stage_id"`
	Data       json.RawMessage `json:"Data" gorm:"data"`
	CreatedBy  string          `json:"CreatedBy" gorm:"created_by"`
	CreatedOn  time.Time       `json:"CreatedOn" gorm:"created_on"`
	ModifiedBy string          `json:"ModifiedBy" gorm:"modified_by"`
	ModifiedOn time.Time       `json:"ModifiedOn" gorm:"modified_on"`
}
