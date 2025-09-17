package model

import (
	"time"

	"github.com/lib/pq"
)

type UserRoles struct {
	ID          int            `json:"ID" gorm:"id"`
	RoleName    string         `json:"RoleName" gorm:"role_name"`
	Description string         `json:"Description" gorm:"description"`
	Deny        pq.StringArray `json:"Deny" gorm:"column:deny;type:text[]"`
	CreatedBy   string         `json:"CreatedBy" gorm:"created_by"`
	CreatedOn   time.Time      `json:"CreatedOn" gorm:"created_on"`
	ModifiedBy  string         `json:"ModifiedBy" gorm:"modified_by"`
	ModifiedOn  time.Time      `json:"ModifiedOn" gorm:"modified_on"`
}
