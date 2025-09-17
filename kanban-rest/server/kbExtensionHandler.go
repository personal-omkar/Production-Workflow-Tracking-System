package server

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"

	"irpl.com/kanban-commons/utils"
)

// GetAllKbExtensions returns a all records present in kb_extension table
func GetAllKbExtensions(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get(DBURL + "/get-all-kb-extension-data")
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

// GetKbExtensionsByParam returns a kb_extension records based on parameter
func GetKbExtensionsByParam(w http.ResponseWriter, r *http.Request) {

	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")

	resp, err := http.Get(DBURL + "/get-kb-extension-by-param?key=" + key + "&value=" + value)
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

func GetOrderDetaislForOrderPage(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	url := DBURL + "/OrderDetailsForCustomer"
	res, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		http.Error(w, "Failed to make request", http.StatusForbidden)
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
