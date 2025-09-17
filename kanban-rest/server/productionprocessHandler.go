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

func GetAllProductionProcessData(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	// Make HTTP GET request to DB
	resp, err := http.Get(DBURL + "/fetch-all-production-process-data")
	if err != nil {
		log.Printf("Error making POST request: %v", err)
		http.Error(w, "Failed to get user information i am in restserv", http.StatusForbidden)
		return
	}
	defer resp.Body.Close()

	var productionProcess []m.ProdProcess
	json.NewDecoder(resp.Body).Decode(&productionProcess)

	if err := json.NewEncoder(w).Encode(productionProcess); err != nil {
		log.Println("Error encoding production lines:", err)
		http.Error(w, "Error generating production lines data", http.StatusInternalServerError)
		return
	}
}

func GetAllProductionProcessEntries(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	line_no := queryParams.Get("line")
	if line_no == "" {
		line_no = "1"
	}
	w.Header().Set("Content-Type", "application/json")

	// Make HTTP GET request to DB
	resp, err := http.Get(DBURL + "/get-all-production-process-for-line?line=" + line_no)
	if err != nil {
		log.Printf("Error making POST request: %v", err)
		http.Error(w, "Failed to get user information i am in restserv", http.StatusForbidden)
		return
	}
	defer resp.Body.Close()

	var data []m.ProdProcessCardData
	json.NewDecoder(resp.Body).Decode(&data)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Println("Error encoding production lines:", err)
		http.Error(w, "Error generating production lines data", http.StatusInternalServerError)
		return
	}
}

func GetAllProductionProcess(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	// Make HTTP GET request to DB
	resp, err := http.Get(DBURL + "/get-all-production-process")
	if err != nil {
		log.Printf("Error making POST request: %v", err)
		http.Error(w, "Failed to get user information i am in restserv", http.StatusForbidden)
		return
	}
	defer resp.Body.Close()

	var productionProcess []m.ProdProcess
	json.NewDecoder(resp.Body).Decode(&productionProcess)

	if err := json.NewEncoder(w).Encode(productionProcess); err != nil {
		log.Println("Error encoding production lines:", err)
		http.Error(w, "Error generating production lines data", http.StatusInternalServerError)
		return
	}
}

func AddProdProcess(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	resp, err := http.Post(DBURL+"/add-production-process", "application/json", bytes.NewBuffer(body))
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

	if resp.StatusCode != http.StatusOK {
		utils.SetResponse(w, http.StatusInternalServerError, string(responseBody))
		return
	}
	utils.SetResponse(w, resp.StatusCode, string(responseBody))

}

func GetProdProcessByParam(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	resp, err := http.Post(DBURL+"/get-prod-process-by-param", "application/json", bytes.NewBuffer(body))
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

	if resp.StatusCode != http.StatusOK {
		utils.SetResponse(w, http.StatusInternalServerError, string(responseBody))
		return
	}
	utils.SetResponse(w, http.StatusOK, string(responseBody))

}

func EditProdProcess(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	resp, err := http.Post(DBURL+"/edit-prod-process", "application/json", bytes.NewBuffer(body))
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

	if resp.StatusCode != http.StatusOK {
		utils.SetResponse(w, http.StatusInternalServerError, string(responseBody))
		return
	}
	utils.SetResponse(w, http.StatusOK, string(responseBody))

}

func GetAllProcessBySearchAndPagination(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	resp, err := http.Post(DBURL+"/get-all-production-process-by-search-paginations", "application/json", bytes.NewBuffer(body))
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
