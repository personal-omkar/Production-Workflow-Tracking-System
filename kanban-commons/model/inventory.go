package model

import "time"

type Inventory struct {
	Id                int       `json:"ID" gorm:"id"`
	CompoundID        int       `json:"CompoundId" gorm:"compound_id"`
	MinQuantity       int       `json:"MinQuantity" gorm:"min_quantity"`
	MaxQuantity       int       `json:"MaxQuantity" gorm:"max_quantity"`
	AvailableQuantity int       `json:"AvailableQuantity" gorm:"available_quantity"`
	Description       string    `json:"Description" gorm:"description"`
	CreatedBy         string    `json:"CreatedBy" gorm:"created_by"`
	CreatedOn         time.Time `json:"CreatedOn" gorm:"created_on"`
	ModifiedBy        string    `json:"ModifiedBy" gorm:"modified_by"`
	ModifiedOn        time.Time `json:"ModifiedOn" gorm:"modified_on"`
}

type ColdStorage struct {
	Id                int       `json:"ID" gorm:"id"`
	CompoundID        int       `json:"CompoundID" gorm:"compound_id"`
	CompoundName      string    `json:"CompoundName" gorm:"compound_name"`
	MinQuantity       int       `json:"MinQuantity" gorm:"min_quantity"`
	MaxQuantity       int       `json:"MaxQuantity" gorm:"max_quantity"`
	AvailableQuantity int       `json:"AvailableQuantity" gorm:"available_quantity"`
	ProductType       string    `json:"ProductType" gorm:"product_type"`
	Description       string    `json:"Description" gorm:"description"`
	CreatedBy         string    `json:"CreatedBy" gorm:"created_by"`
	CreatedOn         time.Time `json:"CreatedOn" gorm:"created_on"`
	ModifiedBy        string    `json:"ModifiedBy" gorm:"modified_by"`
	ModifiedOn        time.Time `json:"ModifiedOn" gorm:"modified_on"`
}
