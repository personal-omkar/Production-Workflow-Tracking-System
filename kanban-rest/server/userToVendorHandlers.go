package server

import (
	"io"
	"log/slog"
	"net/http"

	"irpl.com/kanban-commons/utils"
)

// GetVendorByUserID returns a vendor  records based on user id
func GetVendorByUserID(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")

	resp, err := http.Get(DBURL + "/get-vendor-by-userid?key=" + key + "&value=" + value)
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
