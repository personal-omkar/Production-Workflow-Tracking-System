package server

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"strconv"

	"irpl.com/kanban-dao/dao"
)

// GetUserToVendorByParam returns a usertovendor records based on parameter
func GetUserToVendorByParam(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")

	KbRoot, err := dao.GetUserToVendorByParam(key, value)
	if err != nil {
		slog.Error("Recordes not found", "error", err.Error())
		http.Error(w, "Failed to encode fetch records", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(KbRoot); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// GetVendorByUserID returns a vendor  records based on user id
func GetVendorByUserID(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")

	UserToVendor, err := dao.GetUserToVendorByParam(key, value)
	if err != nil {
		slog.Error("Recordes not found", "error", err.Error())
		http.Error(w, "Failed to encode fetch records", http.StatusInternalServerError)
		return
	}

	if len(UserToVendor) != 0 {
		vendorDetails, err := dao.GetVendorByParam("id", strconv.Itoa(UserToVendor[0].VendorId))
		if err != nil {
			slog.Error("Recordes not found", "error", err.Error())
			http.Error(w, "Failed to encode fetch records", http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(vendorDetails); err != nil {
			log.Printf("Error encoding response: %v", err)
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	} else {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to get vendor record", http.StatusInternalServerError)
		return
	}

}
