package server

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"irpl.com/kanban-commons/utils"
)

// GetAllKBData returns a all records present in kb_data table
func GetAllKBData(w http.ResponseWriter, r *http.Request) {

	resp, err := http.Get(DBURL + "/get-all-kb-data")
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

// GetKBDataByParam returns a kb_data records based on parameter
func GetKBDataByParam(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")

	resp, err := http.Get(DBURL + "/get-kb-data-by-param?key=" + key + "&value=" + value)
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

func GetOrderDEtails(w http.ResponseWriter, r *http.Request) {
	type TableRequest struct {
		Conditions []string `json:"conditions"`
	}
	var data TableRequest
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err == nil {
		jsonValue, err := json.Marshal(data)
		if err != nil {
			slog.Error("%s - error - %s", "Error marshaling table data", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to marshal table condition")
			return
		}

		resp, err := http.Post(DBURL+"/get-order-details", "application/json", bytes.NewBuffer(jsonValue))
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
	} else {

		utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to fetch order record")
		slog.Error("%s - error - %s", "failed to fetch order record", err.Error())
	}

}

func GetCustomerOrderDetails(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	type TableConditions struct {
		Conditions []string `json:"Conditions"`
	}
	var data TableConditions

	if err := decoder.Decode(&data); err == nil {

		jsonValue, err := json.Marshal(data)
		if err != nil {
			slog.Error("%s - error - %s", "Error marshaling compound data", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to marshal compound data")
			return
		}

		resp, err := http.Post(DBURL+"/get-customer-order-details", "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			slog.Error("%s - error - %s", "Error making POST request", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to add compound data in production line")
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

		utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to update record")
		slog.Error("%s - error - %s", "Failed to add order data in KbData table", err.Error())
	}
}

func GetOrderDetailsForHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	url := DBURL + "/get-all-details-for-order"
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

	utils.SetResponse(w, http.StatusOK, string(responseBody))
}

func GetPendingOrderDetailsBySearchAndPagination(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	resp, err := http.Post(DBURL+"/get-order-pending-details-by-search-pagination", "application/json", bytes.NewBuffer(body))
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
