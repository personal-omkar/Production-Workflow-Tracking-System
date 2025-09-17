package services

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	m "irpl.com/kanban-commons/model"
	u "irpl.com/kanban-commons/utils"
)

const smSender string = "&sender=" + "VSSOTG"
const smApiKey string = "apiKey=" + "xuS3OdikGk4-hwcRhJh4W8lJjJrLSWQQUnsD9Lj9Hx"
const otpMessageTemplatePart1 string = "Your OTP for mobile number validation is "
const comapnyNameForFooter string = " - hi Systems"

// SendOTP- sends OTP to provided mobile no
// codeType flag- 1- OTP, 2- Reset Key and 3- Change password key
func SendOTP(mobileNo string, codeType int, user *m.User) bool {
	otpGeneratedInt, otpGeneratedStr, err := u.GenerateOtp(6)
	if err != nil {
		if codeType == 1 {
			log.Printf("Failed to generate OTP, error- %v", err)
		} else if codeType == 2 {
			log.Printf("Failed to generate Reset Key, error- %v", err)
		}
		return false
	}

	var numbers string = "&numbers=" + mobileNo
	var message string = "&message=" + otpMessageTemplatePart1 + otpGeneratedStr + comapnyNameForFooter

	var reqDataStr string = smApiKey + numbers + message + smSender

	reqData := strings.NewReader(reqDataStr)

	resp, err := http.Post("https://api.textlocal.in/send/?", "application/x-www-form-urlencoded", reqData)

	if err != nil {
		log.Printf("Failed to call SMS sending API, api call error- %v", err)
		return false
	}

	defer resp.Body.Close()
	bodyBytes, _ := io.ReadAll(resp.Body)

	// Convert response body to string
	responseInStr := string(bodyBytes)
	log.Printf("SMS API output: %v\n", responseInStr)

	// Assume `otpMessageTemplatePart1` is defined elsewhere in your code
	// Assume `mobileNo`, `otpGeneratedInt`, and `codeType` are defined based on your current context

	if strings.Contains(responseInStr, otpMessageTemplatePart1) {
		mobileNoInt, err := strconv.ParseUint(mobileNo, 10, 64)
		if err != nil {
			log.Printf("Failed to convert mobile number, error- %v", err)
			return false
		}

		// Define the API URL
		dbHelperHost := os.Getenv("DBHELPER_HOST")
		if dbHelperHost == "" {
			dbHelperHost = "0.0.0.0"
		}

		dbHelperPort := os.Getenv("DBHELPER_PORT")
		if dbHelperPort == "" {
			dbHelperPort = "7200"
		}
		apiUrl := "http://" + dbHelperHost + ":" + dbHelperPort + "/"

		if codeType == 1 {
			// Construct the request for finding OTP details
			findOtpDetailsUrl := apiUrl + "find-otp-details?mobileNo=" + strconv.FormatUint(mobileNoInt, 10)
			response, err := http.Get(findOtpDetailsUrl)
			if err != nil || response.StatusCode != http.StatusOK {
				// Assuming the record not found or an error occurred. Now, let's create a new one.
				var otpDetail = map[string]interface{}{
					"MobileNo": mobileNoInt,
					"Otp":      otpGeneratedInt,
				}
				otpDetailJson, _ := json.Marshal(otpDetail)
				response, err = http.Post(apiUrl+"save-otp-details", "application/json", bytes.NewBuffer(otpDetailJson))
				if err != nil || response.StatusCode != http.StatusCreated {
					log.Printf("Failed to save otpDetails record, error- %v", err)
					return false
				}
			}
			// If the record exists, you might need another API call to update the record, which depends on your API design.
		} else if codeType == 2 {
			// Construct the payload for updating the user
			var user = map[string]interface{}{
				"ResetKey":  otpGeneratedInt,
				"ResetDate": time.Now().Format(time.RFC3339),
				// Include other necessary fields from the `user` object
			}
			userJson, _ := json.Marshal(user)
			response, err := http.Post(apiUrl+"create-or-update-user", "application/json", bytes.NewBuffer(userJson))
			if err != nil || response.StatusCode != http.StatusOK {
				log.Printf("Failed to update reset key record, error- %v", err)
				return false
			}
		}
		return true
	}
	// failed to parse response of SMS sending API
	log.Printf("SMS sending API throws error- %v", responseInStr)
	return false
}
