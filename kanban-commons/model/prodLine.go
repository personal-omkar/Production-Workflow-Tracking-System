package model

import (
	"database/sql"
	"time"
)

type ProdLine struct {
	Id              int       `json:"ID" gorm:"id"`
	Name            string    `json:"Name" gorm:"name"`
	Icon            string    `json:"Icon" gorm:"icon"`
	Description     string    `json:"Description" gorm:"description"`
	Status          bool      `json:"Status" gorm:"status"`
	OperatorDisplay string    `json:"OperatorDisplay" gorm:"-"`
	Operator        string    `json:"Operator"`
	CreatedBy       string    `json:"CreatedBy" gorm:"created_by"`
	CreatedOn       time.Time `json:"CreatedOn" gorm:"created_on"`
	ModifiedBy      string    `json:"ModifiedBy" gorm:"modified_by"`
	ModifiedOn      time.Time `json:"ModifiedOn" gorm:"modified_on"`
	OperatorCode    string    `json:"OperatorCode" gorm:"-"`
	MoveToLineID    int       `json:"MoveToLineID" gorm:"-"`
	RunningNumber   int       `json:"RunningNumber" gorm:"running_number"`
	RecipeId        int       `json:"RecipeId" gorm:"recipe_id"`
}

type ProdLineDetails struct {
	ProdLineID   int    `json:"prod_line_id"`
	ProdLineName string `json:"prod_line_name"`
	Cells        []Cell `json:"cells"`
}

type Cell struct {
	CellNumber                 string         `json:"cell_number"`
	KBRunningNo                string         `json:"kb_running_no"`
	KBInitialNo                string         `json:"kb_initial_no"`
	CompoundName               string         `json:"compound_name"`
	CreatedOn                  string         `json:"created_on"`
	MfgDateTime                string         `json:"mfg_date_time"`
	DemandDateTime             string         `json:"demand_date_time"`
	ExpDate                    string         `json:"exp_date"`
	NoOFLots                   int            `json:"NoOFLots" gorm:"no_of_lots"`
	Location                   string         `json:"location"`
	Status                     string         `json:"status"`
	KRId                       string         `json:"krid"`
	KanbanNo                   sql.NullString `json:"kanban_no"`
	ProdProcessID              string         `json:"prod_process_id"`
	LotNo                      string         `json:"lot_no"`
	ProductionProcessLineOrder int            `json:"prod_process_line_order"`
}

type ProductionLineStatus struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	CellData []Cell `json:"cell"`
}

type AddLineStruct struct {
	LineName        string          `json:"line_name"`
	LineDescription string          `json:"line_description"`
	CreatedBy       string          `json:"created_by"`
	ProcessOrders   []ProcessOrders `json:"process_orders"`
}
