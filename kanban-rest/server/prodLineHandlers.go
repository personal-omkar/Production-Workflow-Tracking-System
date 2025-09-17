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

// Get All Prod Line records
func GetProdLineAllrecords(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get(DBURL + "/get-all-prod-line-data")
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

	utils.SetResponse(w, http.StatusOK, string(responseBody))
}

func AddProdLine(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	url := DBURL + "/add-production-line"
	res, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		http.Error(w, "Error creating request to REST service", http.StatusInternalServerError)
		return
	}
	res.Header = r.Header
	client := &http.Client{}
	resp, err := client.Do(res)
	if err != nil {
		http.Error(w, "Failed to send request to REST service", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// If the response is not successful, return an error
	if resp.StatusCode != http.StatusCreated {
		http.Error(w, "Failed to read response body", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		http.Error(w, "Failed to forward response body", http.StatusInternalServerError)
	}
}

func GetProductionLinesHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	// Make HTTP GET request to DB
	resp, err := http.Get(DBURL + "/get-production-line-items")
	if err != nil {
		log.Printf("Error making POST request: %v", err)
		http.Error(w, "Failed to get user information i am in restserv", http.StatusForbidden)
		return
	}
	defer resp.Body.Close()

	var prodLineDetails []m.ProdLineDetails
	json.NewDecoder(resp.Body).Decode(&prodLineDetails)

	if err := json.NewEncoder(w).Encode(prodLineDetails); err != nil {
		log.Println("Error encoding production lines:", err)
		http.Error(w, "Error generating production lines data", http.StatusInternalServerError)
		return
	}
}

func GetProductionLineStaus(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	url := DBURL + "/get-production-line-status"
	res, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Failed to update model license data: %v", err)
		http.Error(w, "Failed to update model license data", http.StatusForbidden)
		return
	}
	defer res.Body.Close()

	client := &http.Client{}
	resp, err := client.Do(res)
	if resp.StatusCode != http.StatusOK || err != nil {
		http.Error(w, "Failed to read response body", http.StatusInternalServerError)
		return
	}
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("%s - error - %s", "Error reading response body", err)
		utils.SetResponse(w, http.StatusInternalServerError, "Failed to read response body")
		return
	}

	utils.SetResponse(w, http.StatusOK, string(responseBody))

}

// DeleteProductionLineByProdLineId Delete production line and sent there data to previous stage (kanban page)
func DeleteProductionLineByProdLineId(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")
	resp, err := http.Get(DBURL + "/delete-production-line-by-prodline-id?id=" + id)
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

	utils.SetResponse(w, http.StatusOK, string(responseBody))
}

// DeleteProductionLineCell delet the cell from production line
func DeleteProductionLineCell(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	resp, err := http.Post(DBURL+"/delete-production-line-cell", "application/json", bytes.NewBuffer(body))
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

// DeleteProductionLineCellDataByProductionLineId Delet cell data production line from data base
func DeleteProductionLineCellDataByProductionLineId(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")
	resp, err := http.Get(DBURL + "/delete-production-line-celldata-by-productionline-id?id=" + id)
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

	utils.SetResponse(w, http.StatusOK, string(responseBody))
}

// DeleteProductionLineById Delet production line with there data
func DeleteProductionLineById(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")
	resp, err := http.Get(DBURL + "/delete-production-line-by-id?id=" + id)
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

	utils.SetResponse(w, http.StatusOK, string(responseBody))
}

func GetProdLineByParam(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	resp, err := http.Post(DBURL+"/get-prod-line-by-param", "application/json", bytes.NewBuffer(body))
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

func EditProdLine(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	resp, err := http.Post(DBURL+"/edit-prod-line", "application/json", bytes.NewBuffer(body))
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

// Get LinedUp Production Process By LineId
func GetLinedUpProductionProcessByLineId(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	resp, err := http.Post(DBURL+"/get-lineup-processes-by-lineid", "application/json", bytes.NewBuffer(body))
	if err != nil {
		slog.Error("%s - error - %s", "Error making POST request", err)
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

// CreateNewOrUpdateExistingProductionLine creates a new production line or updates an existing production line
func CreateNewOrUpdateExistingProductionLine(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	url := DBURL + "/create-new-or-update-existing-production-line"
	client := &http.Client{}

	req, err := http.NewRequest(r.Method, url, r.Body)
	if err != nil {
		http.Error(w, "Error creating request to REST service", http.StatusInternalServerError)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to send request to the target service", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Failed to process request on target service", http.StatusInternalServerError)
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

func GetAllProdLineBySearchAndPagination(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	resp, err := http.Post(DBURL+"/get-all-prod-line-data-by-search-paginations", "application/json", bytes.NewBuffer(body))
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
