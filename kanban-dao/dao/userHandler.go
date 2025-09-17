package dao

import (
	"errors"
	"fmt"
	"log"
	"log/slog"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
	m "irpl.com/kanban-commons/model"
	"irpl.com/kanban-commons/utils"
	u "irpl.com/kanban-commons/utils"
	db "irpl.com/kanban-dao/db"
)

const (
	USER_TABLE           string = "users" // Updated to match the table name
	PASSWORD_PLACEHOLDER string = "********"
)

const DefaultDbHelperHost string = "0.0.0.0" // Default port if not set in env
const DefaultDbHelperPort string = "4200"    // Default port if not set in env

var DBHelperHost string
var DBHelperPort string

func init() {
	DBHelperHost = os.Getenv("DBHELPER_HOST")
	if strings.TrimSpace(DBHelperHost) == "" {
		DBHelperHost = DefaultDbHelperHost
	}

	DBHelperPort = os.Getenv("DBHELPER_PORT")
	if strings.TrimSpace(DBHelperPort) == "" {
		DBHelperPort = DefaultDbHelperPort
	}
}

// GetUsersCount returns the count of entries in the users table
func GetUsersCount() (int64, error) {
	var count int64
	err := db.GetDB().Table(USER_TABLE).Count(&count)
	if err != nil {
		return count, err.Error
	}
	return count, nil
}

// GetAllUsersList returns all users from the users list
func GetAllUsersList() (users []m.User, err error) {
	userQueryErr := db.GetDB().Table(USER_TABLE).Find(&users)
	if userQueryErr != nil {
		err = userQueryErr.Error
	}
	for user := range users {
		users[user].Password = PASSWORD_PLACEHOLDER
	}
	return
}

// CheckIfUserExists checks if a user exists by email
func CheckIfUserExists(User *m.User, email string) error {
	// Attempt to find the first record matching the email.
	result := db.GetDB().Table(USER_TABLE).Where("email = ?", email).First(User)
	// Check if a record was found
	if result.Error == nil {
		// A record was found, return an error indicating the user exists.
		slog.Error("###user already exists with mail- " + email)
		return u.Error("user already exists")
	} else if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// No record was found, which is the expected outcome.
		slog.Error("###user does not exists with mail- " + email)
		return u.Error("record not found")
	}

	// An error occurred that wasn't due to the record not being found (e.g., DB connection issue).
	return result.Error
}

// CreateNewOrUpdateExistingUser creates a new user or updates an existing user
func CreateNewOrUpdateExistingUser(User *m.User) error {

	now := time.Now()
	if User.Email != "" {
		User.ModifiedOn = now
		userdetails, _ := GetRawUserWithEmail(User.Email)
		if userdetails.ID != 0 {

			layout := "2006-01-02 15:04:05"

			userdetails.Email = utils.IfElse(len(User.Email) > 0, User.Email, userdetails.Email)
			userdetails.Username = utils.IfElse(len(User.Username) > 0, User.Username, userdetails.Username)
			userdetails.Password = utils.IfElse(len(User.Password) > 0, User.Password, userdetails.Password)
			userdetails.ApprovedBy = utils.IfElse(len(User.ApprovedBy) > 0, User.ApprovedBy, userdetails.ApprovedBy)

			approvedOnStr := utils.IfElse(len(User.ApprovedOn.GoString()) > 0, User.ApprovedOn.GoString(), userdetails.ApprovedOn.GoString())
			approvedOnTime, _ := time.Parse(layout, approvedOnStr)
			userdetails.ApprovedOn = approvedOnTime

			userdetails.RejectedBy = utils.IfElse(len(User.RejectedBy) > 0, User.RejectedBy, userdetails.RejectedBy)

			rejectedOnStr := utils.IfElse(len(User.RejectedOn.GoString()) > 0, User.RejectedOn.GoString(), userdetails.RejectedOn.GoString())
			rejectedOnTime, _ := time.Parse(layout, rejectedOnStr)
			userdetails.RejectedOn = rejectedOnTime

			userdetails.ModifiedBy = utils.IfElse(len(User.ModifiedBy) > 0, User.ModifiedBy, userdetails.ModifiedBy)
			userdetails.ModifiedOn = time.Now()
			userdetails.Isactive = User.Isactive

			if userdetails.ID != 0 {
				if err := db.GetDB().Table(USER_TABLE).Omit("role_id , created_by, created_on").Save(&userdetails).Error; err != nil {
					return err
				}
			}
		} else {
			User.CreatedOn = now
			if err := db.GetDB().Table(USER_TABLE).Omit("role_id").Create(&User).Error; err != nil {
				return err
			}

		}
	}
	return nil
}

// UpdateExistingUser updates an existing user
func UpdateExistingUser(User *m.User) error {

	now := time.Now()
	if User.Email != "" {
		User.ModifiedOn = now
		userdetails, _ := GetRawUserWithEmail(User.Email)
		if userdetails.ID != 0 {

			layout := "2006-01-02 15:04:05"

			userdetails.Email = utils.IfElse(len(User.Email) > 0, User.Email, userdetails.Email)
			userdetails.Username = utils.IfElse(len(User.Username) > 0, User.Username, userdetails.Username)
			userdetails.Password = utils.IfElse(len(User.Password) > 0, User.Password, userdetails.Password)
			userdetails.ApprovedBy = utils.IfElse(len(User.ApprovedBy) > 0, User.ApprovedBy, userdetails.ApprovedBy)

			approvedOnStr := utils.IfElse(len(User.ApprovedOn.GoString()) > 0, User.ApprovedOn.GoString(), userdetails.ApprovedOn.GoString())
			approvedOnTime, _ := time.Parse(layout, approvedOnStr)
			userdetails.ApprovedOn = approvedOnTime

			userdetails.RejectedBy = utils.IfElse(len(User.RejectedBy) > 0, User.RejectedBy, userdetails.RejectedBy)

			rejectedOnStr := utils.IfElse(len(User.RejectedOn.GoString()) > 0, User.RejectedOn.GoString(), userdetails.RejectedOn.GoString())
			rejectedOnTime, _ := time.Parse(layout, rejectedOnStr)
			userdetails.RejectedOn = rejectedOnTime

			userdetails.ModifiedBy = utils.IfElse(len(User.ModifiedBy) > 0, User.ModifiedBy, userdetails.ModifiedBy)
			userdetails.ModifiedOn = time.Now()
			userdetails.Isactive = User.Isactive

			if userdetails.ID != 0 {
				if err := db.GetDB().Table(USER_TABLE).Omit("role_id , created_by, created_on").Save(&userdetails).Error; err != nil {
					return err
				}
			}
		} else {

			return fmt.Errorf("user with mail (%s) not exits", User.Email)

		}
	}
	return nil
}

// GetUserWithID returns a user associated with a given ID
func GetUserWithID(Id string) (user m.User, err error) {
	condition := u.JoinStr("id=", Id)
	result := db.GetDB().Table(USER_TABLE).Where(condition).First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return user, errors.New("user not found")
	}

	if result.Error != nil {
		return user, result.Error
	}

	user.Password = PASSWORD_PLACEHOLDER
	return user, nil
}

// DeleteUser deletes a User record by SrNo
func DeleteUser(srNo int64) error {
	return db.GetDB().Table(USER_TABLE).Where("sr_no = ?", srNo).Delete(&m.User{}).Error
}

// GetUserWithEmail returns a user associated with a given email
func GetUserWithEmail(email string) (user m.User, err error) {
	result := db.GetDB().Table(USER_TABLE).Where("email = ?", email).First(&user)
	if result.Error != nil {
		err = result.Error
	}
	user.Password = PASSWORD_PLACEHOLDER
	return
}

// GetRawUserWithEmail returns all user details associated with a given email
func GetRawUserWithEmail(email string) (user m.User, err error) {
	result := db.GetDB().Table(USER_TABLE).Where("email = ?", email).First(&user)
	if result.Error != nil {
		err = result.Error
	}
	return
}

// SearchUsersForString returns accumulated results from the user table matching the provided string
func SearchUsersForString(str string) (users []m.User, err error) {
	condtition := u.JoinStr(
		"user_name LIKE '%", str, "%' OR ",
		"email LIKE '%", str, "%' OR ",
		"pay_code LIKE '%", str, "%'")
	userQueryErr := db.GetDB().Table(USER_TABLE).Find(&users, condtition)
	if userQueryErr != nil {
		err = userQueryErr.Error
	}
	for user := range users {
		users[user].Password = PASSWORD_PLACEHOLDER
	}
	return
}

// GetUsersByCriteria returns a paginated list of User records based on criteria
func GetUsersByCriteria(page, limit int, criteria map[string]interface{}) (entries []m.User, totalRecords int, err error) {
	offset := (page - 1) * limit
	var totalRecords64 int64
	query := db.GetDB().Table(USER_TABLE)

	// Apply filters from criteria
	for key, value := range criteria {
		query = query.Where(key+" = ?", value)
	}

	// Get total record count with filters
	err = query.Count(&totalRecords64).Error
	if err != nil {
		return nil, 0, err
	}

	// Convert totalRecords64 to int
	totalRecords = int(totalRecords64)

	// Fetch records with pagination and filters
	err = query.Offset(offset).Limit(limit).Scan(&entries).Error
	if err != nil {
		return nil, 0, err
	}

	return entries, totalRecords, nil
}

// SearchUsersByUserRoleString returns accumulated results from the user table matching the provided string to the user role
func SearchUsersByUserRoleString(str string) (users []m.User, err error) {
	condtition := u.JoinStr("sr_no in (select (end_user_id) from `end_user_role` where role like '%", str, "%')")
	userQueryErr := db.GetDB().Table(USER_TABLE).Find(&users, condtition)
	if userQueryErr != nil {
		err = userQueryErr.Error
	}
	for user := range users {
		users[user].Password = PASSWORD_PLACEHOLDER
	}
	return
}

// GetPaginatedUsersDataDB returns classified results from the user table matching the provided string
func GetPaginatedUsersDataDB(str, limitstr, pagestr string) (users []m.User, err error) {
	var count int
	var limit int
	var page int

	limit, _ = strconv.Atoi(limitstr)
	page, _ = strconv.Atoi(pagestr)
	offset := (page - 1) * limit

	condition := u.JoinStr(
		"user_name LIKE '%", str, "%' OR ",
		"email LIKE '%", str, "%' OR pay_code LIKE '%", str, "%'")

	countErr := db.GetDB().Table(USER_TABLE).Select("COUNT(*)").Find(&count, condition)
	if countErr != nil {
		err = countErr.Error
	}

	res := float64(count) / float64(limit)
	totalPages := int(math.Round(res))
	if page > totalPages {
		offset = 0
	}

	userQueryErr := db.GetDB().Table(USER_TABLE).Limit(limit).Offset(offset).Find(&users, condition)
	if userQueryErr != nil {
		err = userQueryErr.Error
	}
	for user := range users {
		users[user].Password = PASSWORD_PLACEHOLDER
	}

	return
}

// GenerateToken creates a JWT token for the given user and logs any errors
func GenerateToken(User m.User) (string, time.Time, error) {

	token, expire, err := u.CreateJwtToken(User.Email)
	if err != nil {
		log.Printf("Error generating token: %v", err)
		return "", time.Now(), fmt.Errorf("failed to generate token: %v", err)
	}

	return token, expire, nil
}

// CheckUserCredentials checks if the username and password are correct
// On success, return token and user data
// On failure, return error
func CheckUserCredentials(email string, password string) (m.User, []*m.UserRoles, error) {
	var userFound m.User
	var userRoles []*m.UserRoles

	exists, findUserErr := CheckIfUserMasterExists(&userFound, email)
	if !exists {
		log.Printf("User does not exists locally hence check over ldap/active directory: %s, error : %v", email, findUserErr)
		// return userFound, userTypes, findUserErr
	}

	if findUserErr != nil {
		return userFound, userRoles, findUserErr
	}

	// do authentication over ldap or do locally based on the domain (used by user)
	ldapConfig, ldapConfigErr := GetDefaultLDAPConfig()

	if ldapConfigErr != nil {
		log.Printf("Error fetching default ldap config : %v", ldapConfigErr)
	}

	authResult, authErr := AuthenticateUserOverLdapOrLocal(ldapConfig, &userFound, email, password)

	if authErr != nil {
		log.Printf("Error performing authentication for user: %s, error : %v", email, authErr)
		return userFound, userRoles, authErr
	}

	// if user does not exists locally but exists on ldap/active directory, create new user
	if !exists && authResult {
		userName, _, _ := u.ExtractUsernameAndDomain(email)

		userFound.Email = email
		userFound.Password = password
		userFound.Username = userName
		//userFound.UserId = userName
		userFound.Isactive = true
		//userFound.CreatedBy = 1
		createNewUserErr := CreateNewOrUpdateExistingUser(&userFound)

		if createNewUserErr != nil {
			log.Printf("Error failed to create new user (it exists on ldap/active directory): %s, error : %v", email, createNewUserErr)
			return userFound, userRoles, createNewUserErr
		}
	}

	// if authentication is successful, get userRoles information
	if authResult {
		var rawQuery m.RawQuery
		rawQuery.Host = DBHelperHost
		rawQuery.Port = DBHelperPort
		rawQuery.Type = "UserRoles"
		rawQuery.Query = `SELECT * FROM user_roles WHERE id = (SELECT (userroleid) FROM usertorole WHERE userid = ` + fmt.Sprint(userFound.ID) + `)` //`;`
		rawQuery.RawQry(&userRoles)

		return userFound, userRoles, nil
	} else {
		// authentication failure
		return userFound, userRoles, errors.New("authentication failed")
	}
}

// GetUserRoleByRoleID returns a user role associated with a given id
func GetUserRoleByRoleID(roleID int) (user m.UserRoles, err error) {
	result := db.GetDB().Table("user_roles").Where("id = ?", roleID).First(&user)
	if result.Error != nil {
		err = result.Error
	}
	return
}

// CheckIfUserMasterExists checks if a user exists by email
func CheckIfUserMasterExists(User *m.User, email string) (bool, error) {
	// Attempt to find the first record matching the email.
	result := db.GetDB().Table(USER_TABLE).Where("email = ?", email).First(&User)

	if result.Error == nil {
		// A record was found, return true with no error.
		slog.Info("User exists with email: " + email)
		return true, nil
	} else if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// No record was found, return false with no error.
		slog.Info("User does not exist with email: " + email)
		return false, nil
	}

	// An error occurred that wasn't due to the record not being found (e.g., DB connection issue).
	slog.Error("Error occurred while checking user existence: " + result.Error.Error())
	return false, result.Error
}
func GetAllUserBySearchAndPagination(pagination m.PaginationReq, conditions []string) (op []*m.UserManagement, paginationResp m.PaginationResp, err error) {
	// Base query with double join to fetch role name
	dbQuery := db.GetDB().
		Table("users AS u").
		Select("u.*, r.role_name").
		Joins("LEFT JOIN usertorole ur ON u.id = ur.userid").
		Joins("LEFT JOIN user_roles r ON ur.userroleid = r.id")

	// Parse search conditions
	var parsedConditions []string
	for _, cond := range conditions {
		parts := strings.SplitN(cond, " ILIKE ", 2)
		if len(parts) < 2 {
			continue
		}
		field := strings.TrimSpace(parts[0])
		value := strings.Trim(parts[1], "'%")
		if value == "" {
			continue
		}
		parsedConditions = append(parsedConditions, fmt.Sprintf("%s ILIKE '%%%s%%'", field, value))
	}

	// Apply where clause
	if len(parsedConditions) > 0 {
		dbQuery = dbQuery.Where(strings.Join(parsedConditions, " AND "))
	}

	// Get total count
	var totalRecords int64
	countQuery := db.GetDB().Table("users AS u").
		Joins("LEFT JOIN usertorole ur ON u.id = ur.userid").
		Joins("LEFT JOIN user_roles r ON ur.userroleid = r.id")

	if len(parsedConditions) > 0 {
		countQuery = countQuery.Where(strings.Join(parsedConditions, " AND "))
	}
	if err := countQuery.Count(&totalRecords).Error; err != nil {
		return nil, m.PaginationResp{}, err
	}

	// Sorting
	orderBy := "u.id DESC"
	if pagination.Order != "" {
		orderBy = pagination.Order
	}
	dbQuery = dbQuery.Order(orderBy)

	// Pagination
	limit, errLimit := strconv.Atoi(pagination.Limit)
	pageNo := pagination.PageNo
	if errLimit != nil || limit <= 0 {
		limit = 15
	}
	if pageNo <= 0 {
		pageNo = 1
	}
	offset := (pageNo - 1) * limit
	dbQuery = dbQuery.Limit(limit).Offset(offset)

	// Final query execution
	if err := dbQuery.Find(&op).Error; err != nil {
		return nil, m.PaginationResp{}, err
	}

	// Prepare pagination response
	paginationResp = m.PaginationResp{
		TotalNo: int(totalRecords),
		Page:    pageNo,
		Offset:  offset,
	}

	return op, paginationResp, nil
}
