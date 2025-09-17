package server

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"log/slog"
	"net/http"

	m "irpl.com/kanban-commons/model"
	"irpl.com/kanban-commons/utils"
)

// Create User
func CreateUser(w http.ResponseWriter, r *http.Request) {

	var data m.User

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err == nil {

		jsonValue, err := json.Marshal(data)
		if err != nil {
			slog.Error("%s - error - %s", "Error marshaling user data", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to marshal user data")
			return
		}

		resp, err := http.Post(DBURL+"/create-new-user", "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			slog.Error("%s - error - %s", "Error making POST request", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to create user")
			return
		}
		defer resp.Body.Close()

		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			slog.Error("%s - error - %s", "Error reading response body", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to read response body")
			return
		}

		utils.SetResponse(w, resp.StatusCode, string(responseBody))
	} else {
		slog.Error("%s - error - %s", "Record creation failed", err.Error())
		utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to create record")
	}
}

// Update User
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	var data m.User

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err == nil {

		jsonValue, err := json.Marshal(data)
		if err != nil {
			slog.Error("%s - error - %s", "Error marshaling user data", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to marshal user data")
			return
		}

		resp, err := http.Post(DBURL+"/update-user-details", "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			slog.Error("%s - error - %s", "Error making POST request", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to update user")
			return
		}
		defer resp.Body.Close()

		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			slog.Error("%s - error - %s", "Error reading response body", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to read response body")
			return
		}

		utils.SetResponse(w, resp.StatusCode, string(responseBody))
	} else {
		slog.Error("%s - error - %s", "Record updation failed", err.Error())
		utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to update record")
	}
}

// Delete User
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	var ids []int64 // IDs to delete
	decoder := json.NewDecoder(r.Body)

	// Decode the request body to get the list of IDs
	if err := decoder.Decode(&ids); err == nil {
		// Marshal the IDs into JSON
		jsonValue, err := json.Marshal(ids)
		if err != nil {
			slog.Error("%s - error - %s", "Error marshaling delete request data", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to marshal delete request data")
			return
		}

		// Make a POST request to the DB service for deletion
		resp, err := http.Post(DBURL+"/delete-user", "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			slog.Error("%s - error - %s", "Error making POST request for delete", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to delete user(s)")
			return
		}
		defer resp.Body.Close()

		// Read the response body from the service
		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			slog.Error("%s - error - %s", "Error reading response body for delete", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to read response body for delete")
			return
		}

		// Respond to the client with the service's response
		utils.SetResponse(w, http.StatusOK, string(responseBody))
	} else {
		slog.Error("%s - error - %s", "Record deletion failed", err.Error())
		utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to delete record")
	}
}

// Get User by parameters
func GetUserByParam(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")

	// Create the payload for the DB service
	queryParams := map[string]string{
		"key":   key,
		"value": value,
	}
	jsonValue, err := json.Marshal(queryParams)
	if err != nil {
		slog.Error("%s - error - %s", "Error marshaling query parameters", err)
		utils.SetResponse(w, http.StatusInternalServerError, "Failed to marshal query parameters")
		return
	}

	// Make a POST request to the DB service to get user by parameters
	resp, err := http.Post(DBURL+"/get-user-by-param", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		slog.Error("%s - error - %s", "Error making POST request for query", err)
		utils.SetResponse(w, http.StatusInternalServerError, "Failed to query user")
		return
	}
	defer resp.Body.Close()

	// Read the response body from the service
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("%s - error - %s", "Error reading response body for query", err)
		utils.SetResponse(w, http.StatusInternalServerError, "Failed to read response body for query")
		return
	}

	// Respond to the client with the service's response
	utils.SetResponse(w, http.StatusOK, string(responseBody))
}

func GetUserDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var userManagement []m.UserManagement
	// Make HTTP GET request to DB
	resp, err := http.Get(DBURL + "/get-user-details")
	if err != nil {
		log.Printf("Error making GET request: %v", err)
		http.Error(w, "Failed to get user information i am in restserv", http.StatusForbidden)
		return
	}
	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&userManagement)

	if err := json.NewEncoder(w).Encode(userManagement); err != nil {
		log.Println("Error encoding user details:", err)
		http.Error(w, "Error to retrive user details", http.StatusInternalServerError)
		return
	}
}

func GetUserDetailsByEmail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	email := r.URL.Query().Get("email")
	var userManagement []m.UserManagement
	// Make HTTP GET request to DB
	resp, err := http.Get(DBURL + "/get-user-details-by-email?email=" + email)
	if err != nil {
		log.Printf("Error making GET request: %v", err)
		http.Error(w, "Failed to get user information i am in restserv", http.StatusForbidden)
		return
	}
	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&userManagement)

	if err := json.NewEncoder(w).Encode(userManagement); err != nil {
		log.Println("Error encoding user details:", err)
		http.Error(w, "Error to retrive user details", http.StatusInternalServerError)
		return
	}
}

// Update UserDetails
func UpdateUserDetails(w http.ResponseWriter, r *http.Request) {
	var data m.UserManagement

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err == nil {
		jsonValue, err := json.Marshal(data)
		if err != nil {
			slog.Error("Error marshaling user role data", "error", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to marshal user role data")
			return
		}

		resp, err := http.Post(DBURL+"/update-user-details", "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			slog.Error("Error making POST request", "error", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to update user role")
			return
		}
		defer resp.Body.Close()

		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			slog.Error("Error reading response body", "error", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to read response body")
			return
		}

		utils.SetResponse(w, http.StatusOK, string(responseBody))
	} else {
		slog.Error("Record updation failed", "error", err.Error())
		utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to update record")
	}
}

// Update User Status
func UpdateUserStatus(w http.ResponseWriter, r *http.Request) {
	var data m.User
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err == nil {

		jsonValue, err := json.Marshal(data)
		if err != nil {
			slog.Error("%s - error - %s", "Error marshaling user data", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to marshal user data")
			return
		}
		resp, err := http.Post(DBURL+"/update-user-status", "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			slog.Error("%s - error - %s", "Error making POST request", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to update user")
			return
		}
		defer resp.Body.Close()

		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			slog.Error("%s - error - %s", "Error reading response body", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to read response body")
			return
		}

		utils.SetResponse(w, resp.StatusCode, string(responseBody))
	} else {

		slog.Error("%s - error - %s", "Record updation failed", err.Error())
		utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to update record")
	}
}
func GetAllUserBySearchAndPagination(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	resp, err := http.Post(DBURL+"/get-all-user-by-search-pagination", "application/json", bytes.NewBuffer(body))
	if err != nil {
		slog.Error("%s - error - %s", "Error making POST request", err)
		utils.SetResponse(w, http.StatusInternalServerError, "Failed to fectch order details")
		return
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("%s - error - %s", "Error reading response body", err)
		utils.SetResponse(w, http.StatusInternalServerError, "Failed to read response body")
		return
	}

	utils.SetResponse(w, resp.StatusCode, string(responseBody))
}
