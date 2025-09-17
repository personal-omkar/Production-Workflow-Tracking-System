package model

type OtpVerificationReq struct {
	MobileNo string `json:"mobileNo"`
	OTP      string `json:"otp"`
}
