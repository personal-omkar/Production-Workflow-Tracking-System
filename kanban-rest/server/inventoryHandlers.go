package server

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	m "irpl.com/kanban-commons/model"
	"irpl.com/kanban-commons/utils"
)

// GetAllInventoryrecords returns a all records present in inventory table
func GetAllInventoryrecords(w http.ResponseWriter, r *http.Request) {

	resp, err := http.Get(DBURL + "/get-all-inventory-data")
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

// GetAllColdStoragerecords returns a all records present in inventory table with compound name
func GetAllColdStoragerecords(w http.ResponseWriter, r *http.Request) {

	resp, err := http.Get(DBURL + "/get-all-coldstorage-data")
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

// Add Quantity min and max quantity of cold storage
func UpdateColdStorageQuantity(w http.ResponseWriter, r *http.Request) {

	var data m.Inventory

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err == nil {
		jsonValue, err := json.Marshal(data)
		if err != nil {
			slog.Error("%s - error - %s", "Error marshaling cold storage data", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to marshal cold storage data")
			return
		}

		resp, err := http.Post(DBURL+"/update-coldstorage-quantity", "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			slog.Error("%s - error - %s", "Error making POST request", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to update cold storage data in inventorytable")
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
	} else {

		utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to update record")
		slog.Error("%s - error - %s", "Failed to update cold storage data in inventory table", err.Error())
	}
}

func CreateOrUpdateInventory(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("Unable to read request body", "error", err)
		utils.SetResponse(w, http.StatusInternalServerError, "Failed to read response body")
		return
	}
	defer r.Body.Close()

	url := DBURL + "/create-new-or-update-existing-inventory"
	res, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		slog.Error("Error creating request to REST service", "error", err)
		utils.SetResponse(w, http.StatusInternalServerError, "Failed to read response body")
		return
	}
	res.Header = r.Header
	client := &http.Client{}
	resp, err := client.Do(res)
	if err != nil {
		slog.Error("Failed to send request to REST service", "error", err)
		utils.SetResponse(w, http.StatusInternalServerError, "Failed to read response body")
		return
	}
	defer resp.Body.Close()

	// If the response is not successful, return an error
	if resp.StatusCode != http.StatusOK {
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

func DeleteInventoryById(w http.ResponseWriter, r *http.Request) {
	var inventory m.Inventory
	json.NewDecoder(r.Body).Decode(&inventory)
	url := DBURL + "/delete-inventory-by-param?key=id&value=" + strconv.Itoa(inventory.Id)
	res, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		slog.Error("Error creating request to REST service", "error", err)
		utils.SetResponse(w, http.StatusInternalServerError, "Failed to read response body")
		return
	}
	res.Header = r.Header
	client := &http.Client{}
	resp, err := client.Do(res)
	if err != nil {
		slog.Error("Failed to send request to REST service", "error", err)
		utils.SetResponse(w, http.StatusInternalServerError, "Failed to read response body")
		return
	}
	defer resp.Body.Close()

	// If the response is not successful, return an error
	if resp.StatusCode != http.StatusOK {
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

// GetOperatorByParam returns a operator records based on parameter
func GetInventoryByParam(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")
	resp, err := http.Get(DBURL + "/get-inventory-by-param?key=" + key + "&value=" + value)
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

func GetInventoryBySearch(w http.ResponseWriter, r *http.Request) {
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

		resp, err := http.Post(DBURL+"/get-inventory-by-search", "application/json", bytes.NewBuffer(jsonValue))
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

		utils.SetResponse(w, http.StatusOK, string(responseBody))
	} else {

		utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to update record")
		slog.Error("%s - error - %s", "Failed to add order data in KbData table", err.Error())
	}
}

func ColdStorageSearchPagination(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	resp, err := http.Post(DBURL+"/get-all-cold-storage-by-search-pagination", "application/json", bytes.NewBuffer(body))
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
