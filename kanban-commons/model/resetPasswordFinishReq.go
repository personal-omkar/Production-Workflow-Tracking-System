package model

type ResetPasswordFinishReq struct {
	MobileNo    string `json:"mobileNo"`
	ResetKey    int    `json:"resetKey"`
	NewPassword string `json:"newPassword"`
}
