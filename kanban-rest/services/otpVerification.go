package services

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	m "irpl.com/kanban-commons/model"
	dao "irpl.com/kanban-dao/dao"
)

func VerifyOtp(w http.ResponseWriter, r *http.Request) {
	var otpVerificationReq m.OtpVerificationReq

	// Decode JSON request body
	err := json.NewDecoder(r.Body).Decode(&otpVerificationReq)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Basic OTP validation
	if len(otpVerificationReq.OTP) == 0 || len(otpVerificationReq.OTP) < 6 {
		http.Error(w, "Invalid OTP format", http.StatusBadRequest)
		return
	}

	// Convert mobile number to uint64
	mobileNoInt, err := strconv.ParseUint(otpVerificationReq.MobileNo, 10, 64)
	if err != nil {
		http.Error(w, "Invalid mobile number", http.StatusBadRequest)
		return
	}

	// Check OTP details
	otpDetail, err := dao.FindOtpDetailsForMobileNo(mobileNoInt)
	if err != nil {
		http.Error(w, "OTP details not found", http.StatusFailedDependency)
		return
	}

	// Convert the OTP to int and validate
	otpInIntFormat, err := strconv.Atoi(otpVerificationReq.OTP)
	if err != nil || otpDetail.Otp != otpInIntFormat {
		http.Error(w, "Invalid OTP", http.StatusBadRequest)
		return
	}

	// OTP verified successfully
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, `"OTP verified successfully"`)
}

func main() {
	http.HandleFunc("/verifyotp", VerifyOtp)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
