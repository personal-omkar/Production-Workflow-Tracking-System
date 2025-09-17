package model

import "time"

type KbExtension struct {
	Id         int       `json:"ID" gorm:"id"`
	OrderID    int       `json:"OrderID" gorm:"order_id"`
	Code       string    `json:"Code" gorm:"code"`
	Status     string    `json:"Status" gorm:"status"`
	VendorID   int       `json:"VendorID" gorm:"vendor_id"`
	CreatedBy  string    `json:"CreatedBy" gorm:"created_by"`
	CreatedOn  time.Time `json:"CreatedOn" gorm:"created_on"`
	ModifiedBy string    `json:"ModifiedBy" gorm:"modified_by"`
	ModifiedOn time.Time `json:"ModifiedOn" gorm:"modified_on"`
}

// This struct is used to in update status function
type Status struct {
	ID         string `json:"ID"` // its and kb_data_id
	Status     string `json:"Status"`
	NoOFLots   int    `json:"NoOFLots" gorm:"no_of_lots"`
	Dispatch   int    `json:"DispatchQuantity"`
	Kanban     int    `json:"KanbanQuantity"`
	UserID     string `json:"userId"`
	CompoundID int    `json:"compoundId"`
}

type IsValidStatusUpdate struct {
	KbRootId int    `json:"KbRootId"`
	KbDataId int    `json:"KbDataId"`
	Status   string `json:"Status"`
}
