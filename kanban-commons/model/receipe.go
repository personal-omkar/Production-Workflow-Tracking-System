package model

import (
	"encoding/json"
	"time"
)

type Recipe struct {
	Id               int             `json:"ID" gorm:"id"`
	CompoundName     string          `json:"CompoundName" gorm:"compound_name"`
	CompoundCode     string          `json:"CompoundCode" gorm:"compound_code"`
	CreatedBy        string          `json:"CreatedBy" gorm:"created_by"`
	CreatedOn        time.Time       `json:"CreatedOn" gorm:"created_on"`
	ModifiedBy       string          `json:"ModifiedBy" gorm:"modified_by"`
	ModifiedOn       time.Time       `json:"ModifiedOn" gorm:"modified_on"`
	StageId          int             `json:"StageId" gorm:"-"`
	Data             json.RawMessage `json:"Data" gorm:"-"`
	ProdLineId       int             `json:"ProdLineId" gorm:"-"`
	ProdLineToRecipe int             `json:"ProdLineToRecipe" gorm:"-"`
	BaseQty          string          `json:"BaseQty" gorm:"base_qty"`
}
