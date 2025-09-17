package dao

import (
	"crypto/tls"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/go-ldap/ldap/v3"
	m "irpl.com/kanban-commons/model"
	u "irpl.com/kanban-commons/utils"
	db "irpl.com/kanban-dao/db"
)

const ldapConfigTable = "ldapconfig"

// CreateLDAPConfig creates a new LDAPConfig record
func CreateLDAPConfig(entry m.LDAPConfig) error {
	now := time.Now()
	entry.CreatedOn = now
	return db.GetDB().Table(ldapConfigTable).Create(&entry).Error
}

// GetLDAPConfigByID retrieves an LDAPConfig record by ID
func GetLDAPConfigByID(id int) (entry m.LDAPConfig, err error) {
	result := db.GetDB().Table(ldapConfigTable).Where("id = ?", id).First(&entry)
	return entry, result.Error
}

// UpdateLDAPConfig updates an existing LDAPConfig record
func UpdateLDAPConfig(entry m.LDAPConfig) error {
	now := time.Now()
	entry.ModifiedOn = now
	return db.GetDB().Table(ldapConfigTable).Where("id = ?", entry.ID).Updates(&entry).Error
}

// DeleteLDAPConfig deletes an LDAPConfig record by ID
func DeleteLDAPConfig(id int) error {
	return db.GetDB().Table(ldapConfigTable).Where("id = ?", id).Delete(&m.LDAPConfig{}).Error
}

// GetDefaultLDAPConfig retrieves the default LDAPConfig record
func GetDefaultLDAPConfig() (entry m.LDAPConfig, err error) {
	result := db.GetDB().Table(ldapConfigTable).Where("is_default = ?", true).Limit(1).First(&entry)
	return entry, result.Error
}

// AuthenticateUserOverLdapOrLocal authenticates a user against the LDAP server using the `LDAPConfig` model or do locally
func AuthenticateUserOverLdapOrLocal(config m.LDAPConfig, userFound *m.User, email, inputPassword string) (bool, error) {
	// Extract userName (username) and domain from email
	userName, domain, err := u.ExtractUsernameAndDomain(email)
	if err != nil {
		log.Printf("failed to extract username from email : %v", err)
		return false, fmt.Errorf("failed to extract username from email: %v", err)
	}

	// Convert input domain to BaseDN format
	var inputDomainInBaseDNFormat string
	if strings.Contains(domain, ".") {
		domainParts := strings.Split(domain, ".")
		for _, part := range domainParts {
			if inputDomainInBaseDNFormat != "" {
				inputDomainInBaseDNFormat += ","
			}
			inputDomainInBaseDNFormat += fmt.Sprintf("dc=%s", part)
		}
	} else {
		// Handle single-word domains
		inputDomainInBaseDNFormat = fmt.Sprintf("dc=%s", domain)
	}

	// Check if the domain matches the BaseDN
	if !strings.EqualFold(inputDomainInBaseDNFormat, config.BaseDN) {

		// do the authentication locally, assuming administrative users here
		if inputPassword == userFound.Password {
			return true, nil
		}

		log.Printf("baseDN does not match, baseDN (created using input mail)- %s# and baseDN (from config) : %s#", inputDomainInBaseDNFormat, config.BaseDN)

		return false, fmt.Errorf("failed to authenticate user over active directory and locally: %s", userFound.Email)
	}
	log.Printf("domains/baseDNs matched, baseDN (created using input mail)- %s# and baseDN (from config) : %s#", inputDomainInBaseDNFormat, config.BaseDN)

	// Connect to the LDAP server
	ldapConn, err := ldap.DialURL(config.LDAPURL, ldap.DialWithTLSConfig(&tls.Config{InsecureSkipVerify: config.TLSInsecure}))
	if err != nil {
		return false, fmt.Errorf("failed to connect to LDAP server: %v", err)
	}
	defer ldapConn.Close()

	// Bind with the given Bind DN and passw 	ord
	err = ldapConn.Bind(config.BindDN, config.Password)
	if err != nil {
		return false, fmt.Errorf("failed to bind to LDAP server: %v", err)
	} else {
		log.Printf("Bind successful with given bindDN- %s#", config.BindDN)
	}

	// Search for the user by uid
	filter := fmt.Sprintf("(%s=%s)", config.UniqueIdentifier, email)

	// for openldap use this
	if strings.EqualFold(config.UniqueIdentifier, "uid") {
		filter = fmt.Sprintf("(%s=%s)", config.UniqueIdentifier, email)
	}

	// print filter created
	log.Printf("filter created- %s#", filter)

	searchRequest := ldap.NewSearchRequest(
		config.BaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		filter,
		[]string{"dn"}, // We only need the DN
		nil,
	)

	searchResult, err := ldapConn.Search(searchRequest)
	if err != nil {
		return false, fmt.Errorf("failed to search for user: %v", err)
	} else {
		log.Printf("search result found- %v#", searchResult)
	}

	if len(searchResult.Entries) == 0 {
		return false, fmt.Errorf("user not found or multiple entries returned for uid: %s", userName)
	} else {
		log.Printf("search result found- %v#", searchRequest)
	}

	// Attempt to bind as the user with the provided password using given uniqueIdentifier
	//userPrincipalName := searchResult.Entries[0].GetAttributeValue(config.UniqueIdentifier)
	// for OpenLDAP use DN
	//if strings.EqualFold(config.UniqueIdentifier, "uid") {
	userPrincipalName := searchResult.Entries[0].DN
	//}
	//if userPrincipalName == "" {
	//	return false, fmt.Errorf("%s not found for user %s", config.UniqueIdentifier, userFound.Email)
	//}

	log.Printf("user's unique identifier found- %s#", userPrincipalName)

	err = ldapConn.Bind(userPrincipalName, inputPassword)
	if err != nil {
		log.Printf("bind failed for user, uniqueName- %s# and password used- %s", userPrincipalName, inputPassword)
		return false, fmt.Errorf("authentication failed for user: %s#, userPrincipalName: %s: %v", userPrincipalName, userFound.Email, err)
	}

	return true, nil
}
