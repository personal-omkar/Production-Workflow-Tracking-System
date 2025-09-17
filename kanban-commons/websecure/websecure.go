package websecure

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	m "irpl.com/kanban-commons/model"
	"irpl.com/kanban-commons/utils"
)

func CommonMiddleware(next http.Handler) http.Handler {
	allowedURLs := []string{"/static/*", "/login", "/logout", "/do-login", "/register", "/do-register", "/get-user-by-email", "/sign-in", "/sign-up"}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if contains([]string{"/"}, path) || contains(allowedURLs, path) {
			// log.Println("Allowed URL accessed:", path)
			next.ServeHTTP(w, r)
			return
		}

		// If it's a PLC API, check for valid Bearer token
		if isPLCAPIRequest(path) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
				return
			}

			apiKey := strings.TrimPrefix(authHeader, "Bearer ")
			apiKeyRecord, err := getAPIKeyDetailsFromAPI(apiKey)
			if err != nil || apiKeyRecord == nil || !apiKeyRecord.IsActive {
				http.Error(w, "Invalid or inactive API key", http.StatusUnauthorized)
				return
			}

			// Authorized via API key
			next.ServeHTTP(w, r)
			return
		}

		homeURL := "/login"

		_, err := ValidateCookie(r)
		if err != nil {
			// log.Println("failed to ValidateCookie: ", err)
			http.Redirect(w, r, homeURL, http.StatusTemporaryRedirect)
			return
		}

		// log.Println("ValidateCookie successful for path: ", path)
		next.ServeHTTP(w, r)
	})
}

// CORS Middleware function
func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// Allow requests from any origin
		w.Header().Set("Access-Control-Allow-Origin", "*")
		// Allow GET, POST, OPTIONS methods
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		// Allow Content-Type header
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		// Continue with request
		if req.Method == "OPTIONS" {
			return
		}
		next.ServeHTTP(w, req)
	})
}

func isPLCAPIRequest(path string) bool {
	return strings.Contains(path, "/plc")
}

// checkApiToken extracts and validates JWT token from the Authorization header.
// func checkApiToken(r *http.Request) bool {

// 	// log.Println("checkApiToken: request received")

// 	authHeader := r.Header.Get("Authorization")
// 	if authHeader == "" {
// 		log.Println("checkApiToken: request rejected, no Authorization in header")
// 		return false
// 	}

// 	// Typically, Authorization header is in the format "Bearer <token>",
// 	// so we need to split by space and get the second part
// 	parts := strings.Split(authHeader, " ")
// 	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
// 		log.Println("checkApiToken: request rejected, no bearer in header")
// 		return false
// 	}

// 	token := parts[1]

// 	// log.Println("checkApiToken: token: ", token)

// 	return u.ValidateJwtToken(token)
// }

func contains(s []string, str string) bool {

	for _, v := range s {
		if strings.HasSuffix(v, "*") {
			if strings.Contains(str, strings.ReplaceAll(v, "*", "")) {
				return true
			}
		} else if v == str {
			return true
		}
	}

	return false
}

// getAPIKeyDetailsFromAPI calls the API to get an API key record using the key value.
func getAPIKeyDetailsFromAPI(key string) (*m.APIKey, error) {
	// webServURL := os.Getenv("HOME_URL")
	apiURL := utils.RestURL + "/get-api-key-by-param?key=key&value=" + key

	resp, err := http.Get(apiURL)
	if err != nil {
		log.Printf("getAPIKeyDetailsFromAPI - failed to call API, err: %s", err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("getAPIKeyDetailsFromAPI - API call returned status code %d", resp.StatusCode)
		return nil, errors.New("invalid response from API key service")
	}

	var apiKeys []m.APIKey
	err = json.NewDecoder(resp.Body).Decode(&apiKeys)
	if err != nil {
		log.Printf("getAPIKeyDetailsFromAPI - failed to decode API response, err: %s", err.Error())
		return nil, err
	}

	if len(apiKeys) == 0 {
		return nil, errors.New("no API key found for the given key")
	}

	return &apiKeys[0], nil
}
