package dao

import (
	"errors"
	"strconv"
	"strings"
	"time"

	m "irpl.com/kanban-commons/model"
	u "irpl.com/kanban-commons/utils"
	db "irpl.com/kanban-dao/db"
)

const (
	API_KEYS_TABLE = "api_keys"
)

func CreateNewOrUpdateAPIKey(apiKey *m.APIKey) error {
	now := time.Now()

	if apiKey.Id != 0 {
		apiKey.ModifiedOn = now

		existing, _ := GetAPIKeyByParam("id", strconv.Itoa(apiKey.Id))
		if len(existing) == 0 {
			return errors.New("API key record not found for update")
		}

		apiKey.CreatedBy = existing[0].CreatedBy
		apiKey.CreatedOn = existing[0].CreatedOn

		if err := db.GetDB().Table(API_KEYS_TABLE).Save(&apiKey).Error; err != nil {
			return err
		}
	} else {
		// Generate unique key using ULID
		ulidGen := u.NewUlidGenerator()
		apiKey.Key = ulidGen.CreateID()
		apiKey.CreatedOn = now

		if err := db.GetDB().Table(API_KEYS_TABLE).Create(&apiKey).Error; err != nil {
			return err
		}
	}

	return nil
}

// GetAllAPIKeys returns all API key records
func GetAllAPIKeys() (entries []m.APIKey, err error) {
	result := db.GetDB().Table(API_KEYS_TABLE).Order("id").Find(&entries)
	return entries, result.Error
}

// GetAPIKeyByParam fetches API keys by a key-value parameter
func GetAPIKeyByParam(key string, value any) (entries []m.APIKey, err error) {
	query := db.GetDB().Table(API_KEYS_TABLE).Where(key+" = ?", value).Order("id").Find(&entries)
	return entries, query.Error
}

// GetAPIKeysWithPagination applies search and pagination on API keys
func GetAPIKeysWithPagination(pagination m.PaginationReq, conditions []string) (keys []*m.APIKey, paginationResp m.PaginationResp, err error) {
	dbQuery := db.GetDB().Table(API_KEYS_TABLE)

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
		parsedConditions = append(parsedConditions, field+" ILIKE '%"+value+"%'")
	}

	if len(parsedConditions) > 0 {
		dbQuery = dbQuery.Where(strings.Join(parsedConditions, " AND "))
	}

	var totalRecords int64
	countQuery := db.GetDB().Table(API_KEYS_TABLE)
	if len(parsedConditions) > 0 {
		countQuery = countQuery.Where(strings.Join(parsedConditions, " AND "))
	}
	if err := countQuery.Count(&totalRecords).Error; err != nil {
		return nil, m.PaginationResp{}, err
	}

	orderBy := "id DESC"
	if pagination.Order != "" {
		orderBy = pagination.Order
	}
	dbQuery = dbQuery.Order(orderBy)

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

	if err := dbQuery.Find(&keys).Error; err != nil {
		return nil, m.PaginationResp{}, err
	}

	paginationResp = m.PaginationResp{
		TotalNo: int(totalRecords),
		Page:    pageNo,
		Offset:  offset,
	}

	return keys, paginationResp, nil
}
