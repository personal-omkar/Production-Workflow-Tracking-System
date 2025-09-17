package model

import "time"

type Vendors struct {
	ID                int       `json:"ID" gorm:"id"`
	VendorCode        string    `json:"VendorCode" gorm:"vendor_code"`
	VendorName        string    `json:"VendorName" gorm:"vendor_name"`
	ContactInfo       string    `json:"ContactInfo" gorm:"contact_info"`
	Address           string    `json:"Address" gorm:"address"`
	Isactive          bool      `json:"Isactive" gorm:"isactive"`
	PerDayLotConfig   int       `json:"PerDayLotConfig" gorm:"per_day_lot_config"`
	PerMonthLotConfig int       `json:"PerMonthLotConfig" gorm:"per_month_lot_config"`
	PerHourLotConfig  int       `json:"PerHourLotConfig" gorm:"per_hour_lot_config"`
	CreatedBy         string    `json:"CreatedBy" gorm:"created_by"`
	CreatedOn         time.Time `json:"CreatedOn" gorm:"created_on"`
	ModifiedBy        string    `json:"ModifiedBy" gorm:"modified_by"`
	ModifiedOn        time.Time `json:"ModifiedOn" gorm:"modified_on"`
}

type VendorCompanyTable struct {
	VendorCode               string                 `json:"VendorCode"`
	VendorName               string                 `json:"VendorName"`
	Button                   string                 `json:""`
	CompanyCodeAndName       []VendorCompanyDetails `json:"CompanyCodeAndNameStruct"`
	CompanyCodeAndNameString string                 `json:"CompoundCode"`
	CreatedOn                time.Time              `json:""`
	QualityDoneTime          time.Time
	ModifiedOn               time.Time
}
type VendorCompanyDetails struct {
	CompanyCode string
	CompanyName string
}

type VendorKanban struct {
	Vendor    Vendors                 `json:"vendor"`
	Compounds []CompoundsDataByVendor `json:"compounds"`
}

type CSVVendor struct {
	ID                int       `json:"ID" gorm:"id"`
	VendorCode        string    `json:"VendorCode" gorm:"vendor_code"`
	VendorName        string    `json:"VendorName" gorm:"vendor_name"`
	ContactInfo       string    `json:"ContactInfo" gorm:"contact_info"`
	Address           string    `json:"Address" gorm:"address"`
	Isactive          bool      `json:"Isactive" gorm:"column:is_active"`
	PerDayLotConfig   string    `json:"PerDayLotConfig" gorm:"per_day_lot_config"`
	PerMonthLotConfig string    `json:"PerMonthLotConfig" gorm:"per_month_lot_config"`
	PerHourLotConfig  string    `json:"PerHourLotConfig" gorm:"per_hour_lot_config"`
	CreatedBy         string    `json:"CreatedBy" gorm:"created_by"`
	CreatedOn         time.Time `json:"CreatedOn" gorm:"created_on"`
	ModifiedBy        string    `json:"ModifiedBy" gorm:"modified_by"`
	ModifiedOn        time.Time `json:"ModifiedOn" gorm:"modified_on"`
}
type VendorOrderStatus struct {
	VendorName string `json:"vendor_name"`
	Created    int    `json:"created"`
	Approved   int    `json:"approved"`
	InProgress int    `json:"in_progress"`
	Quality    int    `json:"quality"`
	Dispatch   int    `json:"dispatch"`
	Rejected   int    `json:"rejected"`
	Packing    int    `json:"packing"`
	Completed  int    `json:"completed"`
	Pending    int    `json:"pending"`
}
