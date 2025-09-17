package dao

import (
	"time"

	m "irpl.com/kanban-commons/model"
	db "irpl.com/kanban-dao/db"
)

// return total count of entries in attachment table
func GetOtpCount() (int64, error) {
	var count int64
	otpQueryErr := db.GetDB().Table("otp_details").Count(&count)
	if otpQueryErr != nil {
		return count, otpQueryErr.Error
	}
	return count, nil
}

func SaveOtpDetails(rec *m.OtpDetail) error {

	now := time.Now()
	if rec.ID != 0 {
		rec.ModifiedOn = &now
		if err := db.GetDB().Save(rec).Error; err != nil {
			return err
		}
	} else {
		rec.CreatedOn = &now
		if err := db.GetDB().Create(rec).Error; err != nil {
			return err
		}
	}

	return nil
}

func FindOtpDetailsForMobileNo(mobileNoFromReq uint64) (m.OtpDetail, error) {
	var otpDetail m.OtpDetail

	queryErr := db.GetDB().Where(&m.OtpDetail{MobileNo: mobileNoFromReq}).First(&otpDetail)

	if queryErr.Error != nil {
		return otpDetail, queryErr.Error
	}
	return otpDetail, nil
}
