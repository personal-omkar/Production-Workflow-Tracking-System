package model

import (
	"time"
)

type User struct {
	ID          uint      `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Username    string    `json:"username" gorm:"size:50;unique;not null"`
	Email       string    `json:"email" gorm:"size:100;unique;not null"`
	Password    string    `json:"password" gorm:"column:password_hash;size:255;not null"`
	RoleID      uint      `json:"roleId"`
	ApprovedBy  string    `json:"approvedBy" gorm:"column:approved_by;size:50;unique;not null"`
	ApprovedOn  time.Time `json:"approvedOn" gorm:"column:approved_on"`
	RejectedBy  string    `json:"rejectedBy" gorm:"column:rejected_by;size:50;unique;not null"`
	RejectedOn  time.Time `json:"rejectedOn" gorm:"column:rejected_on"`
	CreatedBy   string    `json:"createdBy" gorm:"size:50"`
	CreatedOn   time.Time `json:"createdOn" gorm:"autoCreateTime"`
	ModifiedBy  string    `json:"modifiedBy" gorm:"size:50"`
	ModifiedOn  time.Time `json:"modifiedOn" gorm:"autoUpdateTime"`
	Isactive    bool      `json:"isactive" gorm:"isactive"`
	VendorsCode string    `json:"vendors" gorm:"-"`
}

type UserManagement struct {
	UserID     int       `json:"userid" gorm:"column:userid"`
	UserName   string    `json:"UserName" gorm:"column:username"`
	Email      string    `json:"Email" gorm:"column:email"`
	Password   string    `json:"Password" gorm:"column:password_hash"`
	RoleId     int       `json:"roleid" gorm:"column:roleid"`
	UserRoleId int       `json:"userroleid" gorm:"column:userroleid"`
	RoleName   string    `json:"RoleName" gorm:"column:role_name"`
	CreatedOn  time.Time `json:"CreatedOn" gorm:"column:created_on"`
	Isactive   bool      `json:"isactive" gorm:"isactive"`
	VendorCode string    `json:"VendorCode" gorm:"vendor_code"`
}

type RegisterRequest struct {
	UserType        string `json:"UserType"`
	FirstName       string `json:"FirstName"`
	LastName        string `json:"LastName"`
	Email           string `json:"Email"`
	Password        string `json:"Password"`
	ConfirmPassword string `json:"ConfirmPassword"`
	Code            string `json:"Code"`
}
