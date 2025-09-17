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

// Add order entry
func CreateNewOrderEntry(w http.ResponseWriter, r *http.Request) {

	var data m.OrderEntry

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err == nil {

		jsonValue, err := json.Marshal(data)
		if err != nil {
			slog.Error("%s - error - %s", "Error marshaling compound data", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to marshal compound data")
			return
		}

		resp, err := http.Post(DBURL+"/create-new-order-entry", "application/json", bytes.NewBuffer(jsonValue))
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

// Add order entry
func CreateMultiNewOrderEntry(w http.ResponseWriter, r *http.Request) {

	var data []m.OrderEntry

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err == nil {

		jsonValue, err := json.Marshal(data)
		if err != nil {
			slog.Error("%s - error - %s", "Error marshaling compound data", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to marshal compound data")
			return
		}

		resp, err := http.Post(DBURL+"/create-multi-new-order-entry", "application/json", bytes.NewBuffer(jsonValue))
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

func UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	url := DBURL + "/update-order-status"
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
		body, _ := io.ReadAll(resp.Body) // Read response body
		http.Error(w, string(body), resp.StatusCode)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Delete order entry
func DeleteOrder(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	url := DBURL + "/delete-order-entry"
	res, err := http.NewRequest(http.MethodDelete, url, bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Failed to delete order: %v", err)
		http.Error(w, "Failed to delete order", http.StatusForbidden)
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

func DailyAndMonthlyVendorLimit(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	url := DBURL + "/check-vendor-lot-limit"
	res, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Failed to delete order: %v", err)
		http.Error(w, "Failed to delete order", http.StatusForbidden)
		return
	}
	defer res.Body.Close()

	client := &http.Client{}
	resp, err := client.Do(res)
	if err != nil {
		http.Error(w, "Failed to read response body", http.StatusInternalServerError)
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

func DailyAndMonthlyVendorLimitByVendorCode(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	url := DBURL + "/check-vendor-lot-limit-by-vendor-code"
	res, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Failed to delete order: %v", err)
		http.Error(w, "Failed to delete order", http.StatusForbidden)
		return
	}
	defer res.Body.Close()

	client := &http.Client{}
	resp, err := client.Do(res)
	if err != nil {
		http.Error(w, "Failed to read response body", http.StatusInternalServerError)
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

func GetAllOrderByVendorCode(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	url := DBURL + "/get-all-orders-by-vendor-code"
	res, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Failed to delete order: %v", err)
		http.Error(w, "Failed to delete order", http.StatusForbidden)
		return
	}
	defer res.Body.Close()
	client := &http.Client{}
	resp, err := client.Do(res)
	if err != nil {
		http.Error(w, "Failed to read response body", http.StatusInternalServerError)
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

func DailyVendorLimit(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	url := DBURL + "/check-daily-lot-limit"
	res, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Failed to delete order: %v", err)
		http.Error(w, "Failed to delete order", http.StatusForbidden)
		return
	}
	defer res.Body.Close()

	client := &http.Client{}
	resp, err := client.Do(res)
	if err != nil {
		http.Error(w, "Failed to read response body", http.StatusInternalServerError)
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
