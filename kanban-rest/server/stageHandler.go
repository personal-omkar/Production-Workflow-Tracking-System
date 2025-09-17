package server

import (
	"bytes"
	"io"
	"log"
	"log/slog"
	"net/http"

	"irpl.com/kanban-commons/utils"
)

// Create STAGE Entry
func CreateStage(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	url := DBURL + "/create-new-stage"
	res, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Failed to request: %v", err)
		http.Error(w, "Failed to Request", http.StatusForbidden)
		return
	}
	defer res.Body.Close()

	client := &http.Client{}
	resp, err := client.Do(res)
	if err != nil {
		log.Printf("Failed to request: %v", err)
		http.Error(w, "Failed to Request", http.StatusForbidden)
		return
	}
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("%s - error - %s", "Error reading response body", err)
		utils.SetResponse(w, http.StatusInternalServerError, "Failed to read response body")
		return
	}

	utils.SetResponse(w, resp.StatusCode, string(responseBody))
}

// Update Stage Entry
func UpdateExistingStage(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	url := DBURL + "/update-stage"
	res, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Failed to request: %v", err)
		http.Error(w, "Failed to Request", http.StatusForbidden)
		return
	}
	defer res.Body.Close()

	client := &http.Client{}
	resp, err := client.Do(res)
	if err != nil {
		log.Printf("Failed to request: %v", err)
		http.Error(w, "Failed to Request", http.StatusForbidden)
		return
	}
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("%s - error - %s", "Error reading response body", err)
		utils.SetResponse(w, http.StatusInternalServerError, "Failed to read response body")
		return
	}

	utils.SetResponse(w, resp.StatusCode, string(responseBody))
}

// Delete UserRole Entry
func DeleteStageByID(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	url := DBURL + "/delete-stage"
	res, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Failed to request: %v", err)
		http.Error(w, "Failed to Request", http.StatusForbidden)
		return
	}
	defer res.Body.Close()

	client := &http.Client{}
	resp, err := client.Do(res)
	if err != nil {
		log.Printf("Failed to request: %v", err)
		http.Error(w, "Failed to Request", http.StatusForbidden)
		return
	}
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("%s - error - %s", "Error reading response body", err)
		utils.SetResponse(w, http.StatusInternalServerError, "Failed to read response body")
		return
	}

	utils.SetResponse(w, resp.StatusCode, string(responseBody))
}

func GetAllStage(w http.ResponseWriter, r *http.Request) {
	url := DBURL + "/get-all-stages"

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		http.Error(w, "Failed to create request", http.StatusBadRequest)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to perform request: %v", err)
		http.Error(w, "Failed to perform request", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close() 

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		http.Error(w, "Failed to read response body", http.StatusInternalServerError)
		return
	}

	utils.SetResponse(w, resp.StatusCode, string(responseBody))
}

// Get Stages by Parameter
func GetStagesByParam(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")
	resp, err := http.Get(DBURL + "/get-stage-by-param?key=" + key + "&value=" + value)
	if err != nil {
		slog.Error("%s - error - %s", "Error making GET request", err)
		utils.SetResponse(w, http.StatusInternalServerError, "Failed to find record")
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

func GetStagesByHeader(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	url := DBURL + "/get-stages-by-header"
	res, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Failed to request: %v", err)
		http.Error(w, "Failed to Request", http.StatusForbidden)
		return
	}
	defer res.Body.Close()

	client := &http.Client{}
	resp, err := client.Do(res)
	if err != nil {
		log.Printf("Failed to request: %v", err)
		http.Error(w, "Failed to Request", http.StatusForbidden)
		return
	}
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("%s - error - %s", "Error reading response body", err)
		utils.SetResponse(w, http.StatusInternalServerError, "Failed to read response body")
		return
	}

	utils.SetResponse(w, resp.StatusCode, string(responseBody))
}
