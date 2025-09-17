package websecure

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/securecookie"

	m "irpl.com/kanban-commons/model"
	"irpl.com/kanban-commons/utils"
)

const DefaultRestHost string = "0.0.0.0" // Default port if not set in env
const DefaultRestPort string = "4200"    // Default port if not set in env

// TODO- move following hash key and block key to the common conf
var (
	CookieHandler = securecookie.New([]byte("44BBBFF4126C6B6E8CD7E1D97C2E4330"),
		[]byte("443C4F976ABF^%)B42ABA92B868443E3"))
)

// ValidateCookie checks if the cookie with the given name is present.
func ValidateCookie(request *http.Request) (cookie *http.Cookie, err error) {
	cookieName := os.Getenv("COOKIE_ID")
	cookie, err = request.Cookie(cookieName)
	if err != nil {
		log.Printf("ValidateCookie - failed to find cookie with name- %s in request, err- %s", cookieName, err.Error())
		return nil, err
	}

	email := GetDataFromCookie(cookie, "ID")
	if len(strings.TrimSpace(email)) <= 0 {
		log.Printf("ValidateCookie - failed to find email from cookie in request, cookiename- %s", cookieName)
		return nil, err
	}
	user, userRoles, err := getUserDetailsFromAPI(email)
	if err != nil {
		log.Printf("ValidateCookie - failed to get user details from API, err- %s", err.Error())
		// return nil, nil, err
	}

	for _, denyUrl := range userRoles.Deny {
		if strings.Contains(denyUrl, request.URL.Path) {
			return nil, fmt.Errorf("unauthorized access")
		}
	}
	// userole, err := getUserRoleFromAPI(int(user.RoleID))
	// if err != nil {
	// 	log.Printf("ValidateCookie - failed to get user details from API, err- %s", err.Error())
	// 	// return nil, nil, err
	// }

	if len(user.Email) > 0 {

		var urole, denyList string
		if userRoles != nil {
			urole = userRoles.RoleName
			denyList = strings.Join(userRoles.Deny, "|")
		}

		// update request headers
		request.Header.Set("X-Custom-Role", urole)
		request.Header.Set("X-Custom-Username", user.Username)
		request.Header.Set("X-Custom-Email", user.Email)
		request.Header.Set("X-Custom-Allowlist", denyList)
		request.Header.Set("X-Custom-Userid", strconv.FormatUint(uint64(user.ID), 10))
	}

	// log.Printf("ValidateCookie - found cookie with name- %s, value- %s", cookieName, cookie.Value)
	return cookie, nil
}

// GetDataFromCookie retrieves data using the key from a validated cookie.
func GetDataFromCookie(cookie *http.Cookie, key string) (cookieId string) {

	cookieValue := make(map[string]string)
	err := CookieHandler.Decode(os.Getenv("COOKIE_ID"), cookie.Value, &cookieValue)
	if err != nil {
		log.Printf("GetDataFromCookie - failed to decode cookie, err- %s", err.Error())
		return
	}

	cookieId = cookieValue[key]
	log.Printf("GetDataFromCookie - data retrieved for key- %s, value- %s", key, cookieId)
	return
}

// Create cookie using key-value pair (string), having provided expiry time and set cookie having SERVERS.SECURITY.COOKIE in response
func SetCookie(key string, value string, isSecure bool, cookieAndSessionExpiryTime time.Time, response http.ResponseWriter) (err error) {

	keyVal := map[string]string{
		key: value,
	}

	// TODO - error handling
	cookieName := os.Getenv("COOKIE_ID")
	domainFromEnv := os.Getenv("DOMAIN")

	slog.Info("cookiename found from env: " + cookieName)
	slog.Info("domain found from env: " + domainFromEnv)

	if encoded, encodeErr := CookieHandler.Encode(cookieName, keyVal); encodeErr == nil {
		cookie := &http.Cookie{
			Name:       cookieName,
			Domain:     "." + domainFromEnv, //The leading dot ensures the cookie is accessible for domain and its subdomains
			Value:      encoded,
			Path:       "/",
			Expires:    cookieAndSessionExpiryTime,
			RawExpires: "",
			MaxAge:     0,
			Secure:     isSecure,
			HttpOnly:   true,
			SameSite:   0,
			Raw:        "",
			Unparsed:   []string{},
		}
		http.SetCookie(response, cookie)
		//log.Printf("SetCookie - created cookie successfully: name- %s and value- %s", cookie.Name, cookie.Value)
	} else {
		// failed to set cookie
		log.Printf("SetCookie - failed to set cookie with name- %s and value- %s", cookieName, value)
		err = encodeErr
	}
	return
}

func SetCookieValues(values map[string]string, response http.ResponseWriter) (err error) {

	keyVal := values

	// TODO - error handling
	cookieName := os.Getenv("COOKIE_ID")
	domainFromEnv := os.Getenv("DOMAIN")

	if encoded, encodeErr := CookieHandler.Encode(cookieName, keyVal); encodeErr == nil {
		cookie := &http.Cookie{
			Name:       cookieName,
			Domain:     "." + domainFromEnv, //The leading dot ensures the cookie is accessible for domain and its subdomains
			Value:      encoded,
			Path:       "/",
			Expires:    time.Now().Add(24 * time.Hour),
			RawExpires: "",
			MaxAge:     0,
			Secure:     false,
			HttpOnly:   true,
			SameSite:   0,
			Raw:        "",
			Unparsed:   []string{},
		}
		http.SetCookie(response, cookie)
		//log.Printf("SetCookie - created cookie successfully: name- %s and value- %s", cookie.Name, cookie.Value)
	} else {
		// failed to set cookie
		err = encodeErr
	}
	return
}

// Clear cookie having SERVERS.SECURITY.COOKIE
func ClearCookie(response http.ResponseWriter) {
	// TODO - error handling
	cookieName := os.Getenv("COOKIE_ID")
	cookie := &http.Cookie{
		Name:       cookieName,
		Value:      "",
		Path:       "/",
		Expires:    time.Now().Add(-10000 * time.Hour),
		RawExpires: "",
		MaxAge:     -1,
		HttpOnly:   true,
		SameSite:   0,
		Raw:        "",
		Unparsed:   []string{},
	}
	http.SetCookie(response, cookie)
}

// post data to specified url using content type as app/json and return the result
func HttpPostReq(reqData []byte, url string) (result []byte, httpCodeResp int, err error) {
	resp, postErr := http.Post(url, "application/json", bytes.NewBuffer(reqData))
	if postErr != nil {
		// failed to execute post req
		log.Printf("HttpPostReq - failed to post request to url- %s, reqData- %s and error- %s", url, string(reqData), postErr.Error())
	}
	defer resp.Body.Close()
	if postErr == nil {
		result, err = io.ReadAll(resp.Body)
		if err != nil {
			// failed to read data from post response body
			log.Printf("HttpPostReq - failed to parse response body, url- %s, reqData- %s and error- %s", url, string(reqData), err.Error())
		} else {
			// success
			httpCodeResp = resp.StatusCode
		}

	} else {
		err = postErr
	}
	return
}

// getUserDetailsFromAPI calls the API to get user details by email.
func getUserDetailsFromAPI(email string) (*m.User, *m.UserRoles, error) {

	webServURL := os.Getenv("HOME_URL")
	apiURL := webServURL + "/get-user-by-email?email=" + email
	resp, err := http.Get(apiURL)
	if err != nil {
		log.Printf("getUserDetailsFromAPI - failed to call API, err- %s", err.Error())
		return nil, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("getUserDetailsFromAPI - API call returned status code %d", resp.StatusCode)
		return nil, nil, err
	}

	var user m.User
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		log.Printf("getUserDetailsFromAPI - failed to decode API response, err- %s", err.Error())
		return nil, nil, err
	}
	var URs []*m.UserRoles
	var userrole *m.UserRoles
	var rawQuery m.RawQuery
	rawQuery.Host = utils.RestHost
	rawQuery.Port = utils.RestPort
	rawQuery.Type = "UserRoles"
	rawQuery.Query = `SELECT * FROM user_roles WHERE id = (SELECT (userroleid) FROM usertorole WHERE userid = ` + fmt.Sprint(user.ID) + `)` //`;`
	rawQuery.RawQry(&URs)

	if len(URs) > 0 {
		userrole = URs[0]
	}
	return &user, userrole, err
}

// getRoleDetailsFromAPI calls the API to get role details by ID.
func getRoleDetailsFromAPI(roleID uint) (*m.UserRoles, error) {
	// Build the API URL
	restServURL := utils.JoinStr("http://", os.Getenv("RESTSRV_HOST"), ":", os.Getenv("RESTSRV_PORT"))
	apiURL := utils.JoinStr(restServURL, "/get-role-by-id?id=", strconv.FormatUint(uint64(roleID), 10))
	// Make the GET request
	resp, err := http.Get(apiURL)
	if err != nil {
		log.Printf("getRoleDetailsFromAPI - failed to call API, err- %s", err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	// Check if the response status code indicates success
	if resp.StatusCode != http.StatusOK {
		log.Printf("getRoleDetailsFromAPI - API call returned status code %d", resp.StatusCode)
		return nil, fmt.Errorf("API call failed with status code %d", resp.StatusCode)
	}

	// Decode the API response into a UserRoles struct
	var role m.UserRoles
	err = json.NewDecoder(resp.Body).Decode(&role)
	if err != nil {
		log.Printf("getRoleDetailsFromAPI - failed to decode API response, err- %s", err.Error())
		return nil, err
	}

	return &role, nil
}
