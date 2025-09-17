package model

import (
	"time"
)

type VendorOrderTable struct {
	ID           string `json:"id"`
	SrNo         string `json:"srno"`
	CustomerName string `json:"customername"`
	CompoundCode string `json:"compoundcode"`
	DemandDate   string `json:"demanddate"`
	LotNumber    string `json:"lotnumber"`
	Status       string `json:"status"`
}

type OrderEntry struct {
	CompoundCode   string    `json:"CompoundCode"`
	DemandDateTime time.Time `json:"DemandDateTime"`
	NoOFLots       int       `json:"NoOFLots" gorm:"no_of_lots"`
	UserID         string    `json:"UserID" `
	Status         string    `json:"Status" `
	Location       string    `json:"Location" `
	CellNo         string    `json:"CellNo"`
	MFGDateTime    time.Time `json:"MFGDateTime"`
	CustomerNote   string    `json:"Note"`
}

type OrderDetails struct {
	Id                          int       `json:"ID" gorm:"id"`
	CompoundId                  int       `json:"CompoundId" gorm:"compound_id"`
	MFGDateTime                 time.Time `json:"MFGDateTime" gorm:"mfg_date_time"`
	DemandDateTime              time.Time `json:"DemandDateTime" gorm:"demand_date_time"`
	ExpDate                     time.Time `json:"ExpDate" gorm:"exp_date"`
	CellNo                      string    `json:"CellNo" gorm:"cell_no"`
	NoOFLots                    int       `json:"NoOFLots" gorm:"no_of_lots"`
	Location                    string    `json:"Location" gorm:"location"`
	KbRootId                    int       `json:"KbRootId" gorm:"kb_root_id"`
	CreatedBy                   string    `json:"CreatedBy" gorm:"created_by"`
	CreatedOn                   time.Time `json:"CreatedOn" gorm:"created_on"`
	ModifiedBy                  string    `json:"ModifiedBy" gorm:"modified_by"`
	ModifiedOn                  time.Time `json:"ModifiedOn" gorm:"modified_on"`
	Status                      string    `json:"Status" gorm:"status"`
	VendorName                  string    `json:"VendorName" gorm:"vendor_name"`
	VendorCode                  string    `json:"VendorCode" gorm:"vendor_code"`
	CompoundName                string    `json:"CompoundName" gorm:"compound_name"`
	OrderId                     string    `josn:"orderid" gorm:"order_id"`
	MinQuantity                 int       `json:"min_quantity" gorm:"column:min_quantity"`
	AvailableQuantity           int       `json:"available_quantity" gorm:"column:available_quantity"`
	InventoryKanbanInProcessQty int       `json:"InventoryKanbanInProcessQty"`
	CustomerName                string    `json:"CustomerName" gorm:"column:customername"`
	LotNo                       string    `json:"LotNo" gorm:"column:lot_no"`
	CustomerNote                string    `json:"Note" gorm:"column:note"`
	KanbanNo                    string    `json:"kanban_no" gorm:"kanban_no"`
}

type CustomerOrderDetails struct {
	Id                int       `json:"ID" gorm:"id"`
	CompoundId        int       `json:"CompoundId" gorm:"compound_id"`
	MFGDateTime       string    `json:"MFGDateTime" gorm:"mfg_date_time"`
	DemandDateTime    string    `json:"DemandDateTime" gorm:"demand_date_time"`
	ExpDate           time.Time `json:"ExpDate" gorm:"exp_date"`
	CellNo            string    `json:"CellNo" gorm:"cell_no"`
	NoOFLots          int       `json:"NoOFLots" gorm:"no_of_lots"`
	Location          string    `json:"Location" gorm:"location"`
	KbRootId          int       `json:"KbRootId" gorm:"kb_root_id"`
	CreatedBy         string    `json:"CreatedBy" gorm:"created_by"`
	CreatedOn         time.Time `json:"CreatedOn" gorm:"created_on"`
	ModifiedBy        string    `json:"ModifiedBy" gorm:"modified_by"`
	ModifiedOn        time.Time `json:"ModifiedOn" gorm:"modified_on"`
	Status            string    `json:"Status" gorm:"status"`
	VendorName        string    `json:"VendorName" gorm:"vendor_name"`
	CompoundName      string    `json:"CompoundName" gorm:"compound_name"`
	OrderId           string    `josn:"orderid" gorm:"order_id"`
	MinQuantity       int       `json:"min_quantity" gorm:"column:min_quantity"`
	AvailableQuantity int       `json:"available_quantity" gorm:"column:available_quantity"`
	CustomerName      string    `json:"CustomerName" `
	LotNo             string    `json:"LotNo" gorm:"column:lot_no"`
	CustomerNote      string    `json:"Note" gorm:"column:note"`
	KanbanNo          string    `json:"kanban_no" gorm:"kanban_no"`
}

type OrderDetailsHistory struct {
	ID            int       `json:"kb_date_id" gorm:"column:id"`
	NoOFLots      int       `json:"no_of_lots" gorm:"column:no_of_lots"`
	Status        string    `json:"status" gorm:"column:status"`
	CellNo        string    `json:"cell_no" gorm:"column:cell_no"`
	VendorName    string    `json:"vendor_name" gorm:"column:vendor_name"`
	VendorCode    string    `json:"vendor_code" gorm:"column:vendor_code"`
	ContactInfo   string    `json:"contact_info" gorm:"column:contact_info"`
	Username      string    `json:"username" gorm:"column:username"`
	Email         string    `json:"email" gorm:"column:email"`
	CompoundName  string    `json:"compound_name" gorm:"column:compound_name"`
	OrderId       string    `json:"order_id" gorm:"column:order_id"`
	OrderOn       time.Time `json:"order_on" gorm:"column:created_on"`
	DemandDate    time.Time `json:"demand_date" gorm:"column:demand_date_time"`
	DispatchDate  time.Time `json:"dispatch_date_time" gorm:"column:modified_on"`
	KanbanDetails []Kanban  `json:"Kanban_details" gorm:"-"`
}

type Kanban struct {
	ID          int       `json:"kb_root_id" gorm:"column:id"`
	LotNo       string    `json:"LotNo" gorm:"column:lot_no"`
	ProdLine    string    `json:"prod_line" gorm:"column:name"`
	CreatedOn   time.Time `json:"StartedOn" gorm:"created_on"`
	CompletedOn time.Time `json:"CompletedOn" gorm:"modified_on"`
	KanbanNo    string    `json:"kanban_no" gorm:"kanban_no"`
}
