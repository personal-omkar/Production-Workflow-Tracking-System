package server

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"irpl.com/kanban-commons/utils"
)

// GetAllKbRoot returns a all records present in kb_data table
func GetAllKbRoot(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get(DBURL + "/get-all-kb-root-data")
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

// GetKbRootByParam returns a kb_root records based on parameter
func GetKbRootByParam(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")
	resp, err := http.Get(DBURL + "/get-kb-root-by-param?key=" + key + "&value=" + value)
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

func UpdateRunningNumbers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	url := DBURL + "/update-running-numbers"
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
	w.WriteHeader(http.StatusOK)
}

func GetAllCompletedKBRootDetails(w http.ResponseWriter, r *http.Request) {
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

		resp, err := http.Post(DBURL+"/get-all-kbRoot-details", "application/json", bytes.NewBuffer(jsonValue))
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

		utils.SetResponse(w, http.StatusOK, string(responseBody))
	} else {

		utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to fetch order record")
		slog.Error("%s - error - %s", "failed to fetch order record", err.Error())
	}

}

func GetCompletedKBRootDetailsBySearch(w http.ResponseWriter, r *http.Request) {
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

		resp, err := http.Post(DBURL+"/get-all-kbRoot-details-by-search", "application/json", bytes.NewBuffer(jsonValue))
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

		utils.SetResponse(w, http.StatusOK, string(responseBody))
	} else {

		utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to fetch order record")
		slog.Error("%s - error - %s", "failed to fetch order record", err.Error())
	}

}

func GetDetailRootData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	url := DBURL + "/get-kbRoot-details"
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

func DeleteKbRootByIDsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("Failed to read request body: %v", err)
		utils.SetResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	url := DBURL + "/delete-kbroot-by-ids"
	req, err := http.NewRequest(http.MethodDelete, url, bytes.NewReader(body))
	if err != nil {
		slog.Error("Error creating request to REST service: %v", err)
		utils.SetResponse(w, http.StatusInternalServerError, "Error creating request to REST service")
		return
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("Failed to send request to REST service: %v", err)
		utils.SetResponse(w, http.StatusInternalServerError, "Failed to send request to the target service")
		return
	}
	defer resp.Body.Close()
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Error reading response body: %v", err)
		utils.SetResponse(w, http.StatusInternalServerError, "Failed to read response body")
		return
	}
	utils.SetResponse(w, resp.StatusCode, string(responseBody))
}

func GetAllKanbanDetailsForReport(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	resp, err := http.Post(DBURL+"/get-all-kanban-details-for-report", "application/json", bytes.NewBuffer(body))
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
