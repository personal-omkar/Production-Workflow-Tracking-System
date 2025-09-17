package model

import (
	"time"
)

// LDAPConfig represents the ldap_config table in the database
type LDAPConfig struct {
	ID               int       `gorm:"column:id;primaryKey" json:"id"`
	LDAPURL          string    `gorm:"column:ldap_url" json:"ldap_url"`
	BindDN           string    `gorm:"column:bind_dn" json:"bind_dn"`
	BaseDN           string    `gorm:"column:base_dn" json:"base_dn"`
	Password         string    `gorm:"column:password" json:"password"`
	UniqueIdentifier string    `gorm:"column:unique_identifier" json:"unique_identifier"`
	TLSInsecure      bool      `gorm:"column:tls_insecure" json:"tls_insecure"`
	IsDefault        bool      `gorm:"column:is_default" json:"is_default"`
	CreatedOn        time.Time `gorm:"column:createdon" json:"createdon"`
	CreatedBy        string    `gorm:"column:createdby" json:"createdby"`
	ModifiedOn       time.Time `gorm:"column:modifiedon" json:"modifiedon"`
	ModifiedBy       string    `gorm:"column:modifiedby" json:"modifiedby"`
}
