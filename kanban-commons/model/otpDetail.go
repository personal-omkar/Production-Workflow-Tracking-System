package model

import "time"

type OtpDetail struct {
	ID         uint64     `json:"id"`
	MobileNo   uint64     `json:"mobileNo"`
	Otp        int        `json:"otp"`
	CreatedOn  *time.Time `json:"createdOn"`
	ModifiedOn *time.Time `json:"modifiedOn"`
}
