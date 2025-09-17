package model

import (
	"time"
)

type SambaOtherConfig struct {
	DomainName  string `gorm:"column:domain_name" json:"domain_name"`
	Username    string `gorm:"column:username" json:"username"`
	Password    string `gorm:"column:password" json:"password"`
	TLSInsecure bool   `gorm:"column:tls_insecure" json:"tls_insecure"`
}

// SambaConfig represents the sambaconfig table in the database
type SambaConfig struct {
	ID           int       `gorm:"column:id;primaryKey" json:"id"`
	Workgroup    string    `gorm:"column:workgroup" json:"workgroup"`
	ServerString string    `gorm:"column:server_string" json:"server_string"`
	Security     string    `gorm:"column:security" json:"security"`
	IsDefault    bool      `gorm:"column:is_default" json:"is_default"`
	CreatedOn    time.Time `gorm:"column:createdon" json:"createdon"`
	CreatedBy    string    `gorm:"column:createdby" json:"createdby"`
	ModifiedOn   time.Time `gorm:"column:modifiedon" json:"modifiedon"`
	ModifiedBy   string    `gorm:"column:modifiedby" json:"modifiedby"`
}
