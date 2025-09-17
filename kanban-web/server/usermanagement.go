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
	"irpl.com/kanban-web/pages/basepage"
	"irpl.com/kanban-web/pages/usermanagmentpage"
)

func usermanagmentPage(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("X-Custom-Username")
	userID := r.Header.Get("X-Custom-Userid")
	usertype := r.Header.Get("X-Custom-Role")
	links := r.Header.Get("X-Custom-Allowlist")

	var vendorName string
	var vendorRecord []m.Vendors
	if usertype != "Admin" {
		resp, err := http.Get(RestURL + "/get-vendor-by-userid?key=user_id&value=" + userID)
		if err != nil {
			slog.Error("%s - error - %s", "Error making GET request", err)
		}

		defer resp.Body.Close()

		if err := json.NewDecoder(resp.Body).Decode(&vendorRecord); err != nil {
			slog.Error("error decoding response body", "error", err)
		}
		if len(vendorRecord) != 0 {
			vendorName = vendorRecord[0].VendorName
		} else {
			vendorName = ""
		}

	} else {
		vendorName = ""
	}
	// Define the side navigation items
	sideNav := basepage.SideNav{
		MenuItems: []basepage.SideNavItem{
			{
				Name:     "Dashboard",
				Icon:     "fas fa-chart-pie",
				Link:     "/dashboard",
				UserType: basepage.UserType{Admin: true, Operator: true, Customer: true},
				Style:    "font-size:1rem;",
			},
			{
				Name:     "User Master",
				Icon:     "fas fa-users",
				Link:     "/user-management",
				UserType: basepage.UserType{Admin: true},
				Style:    "font-size:1rem;",
				Selected: true,
			},
			{
				Name:     "Vendor Master",
				Icon:     "fa fa-briefcase",
				Link:     "/vendor-management",
				UserType: basepage.UserType{Admin: true},
				Style:    "font-size:1rem;",
			},
			{
				Name:     "Operator Master",
				Icon:     "fa fa-user",
				Link:     "/operator-management",
				UserType: basepage.UserType{Admin: true},
				Style:    "font-size:1rem;",
			},
			{
				Name:     "Part Master",
				Icon:     "fas fa-vials",
				Link:     "/compounds-management",
				UserType: basepage.UserType{Admin: true, Operator: true},
				Style:    "font-size:1rem;",
			},
			{
				Name:     "Chemical Type Master",
				Icon:     "fas fa-vial",
				Link:     "/chemical-management",
				UserType: basepage.UserType{Admin: true},
				Style:    "font-size:1rem;",
			},
			{
				Name:     "Raw Material",
				Icon:     "fas fa-boxes",
				Link:     "/material-management",
				UserType: basepage.UserType{Admin: true},
				Style:    "font-size:1rem;",
			},
			{
				Name:     "Rubber Store Master",
				Icon:     "fas fa-memory",
				Link:     "/inventory-management",
				UserType: basepage.UserType{Admin: true},
				Style:    "font-size:1rem;",
			},
			{
				Name:     "Recipe Master",
				Icon:     "fa fa-clipboard-list",
				Link:     "/recipe-management",
				UserType: basepage.UserType{Admin: true, Operator: true},
				Style:    "font-size:1rem;",
			},
			{
				Name:     "Machine Master",
				Icon:     "fas fa-sliders-h",
				Link:     "/prod-line-management",
				UserType: basepage.UserType{Admin: true},
				Style:    "font-size:1rem;",
			},
			{
				Name:     "Process Master",
				Icon:     "fas fa-project-diagram",
				Link:     "/prod-processes-management",
				UserType: basepage.UserType{Admin: true},
				Style:    "font-size:1rem;",
			},
			{
				Name:     "Kanban Report",
				Icon:     "fas fa-list-alt",
				Link:     "/vendor-orders",
				UserType: basepage.UserType{Admin: true, Operator: true, Customer: true},
				Style:    "font-size:1rem;",
			},
			{
				Name:     utils.DefaultsMap["cold_store_menu"],
				Icon:     "fas fa-store",
				Link:     "/cold-storage",
				UserType: basepage.UserType{Admin: true},
				Style:    "font-size:1rem;",
			},
			{
				Name:     "Pending Orders",
				Icon:     "fas fa-list-alt",
				Link:     "/admin-orders",
				UserType: basepage.UserType{Admin: true, Operator: true, Customer: true},
				Style:    "font-size:1rem;",
			},
			{
				Name:     "All Kanban View",
				Icon:     "fa fa-th-list",
				Link:     "/all-kanban-view",
				Style:    "font-size:1rem;",
				UserType: basepage.UserType{Admin: true, Operator: true},
			},
			{
				Name:     "Kanban Board",
				Icon:     "fas fa-tasks",
				Link:     "/vendor-company",
				UserType: basepage.UserType{Admin: true, Operator: true},
				Style:    "font-size:1rem;",
			},
			{
				Name:     "Heijunka Board",
				Icon:     "fas fa-calendar-alt",
				Link:     "/production-line",
				UserType: basepage.UserType{Admin: true, Operator: true},
				Style:    "font-size:1rem;",
			},
			{
				Name:     "Kanban Entry",
				Icon:     "fas fa-plus",
				Link:     "/order-entry",
				UserType: basepage.UserType{Customer: true},
				Style:    "font-size:1rem;",
			},
			{
				Name:     "Quality Testing",
				Icon:     "fas fa-check-double",
				Link:     "/quality-testing",
				UserType: basepage.UserType{Admin: true},
				Style:    "font-size:1rem;",
			},
			{
				Name:     "Packing/Dispatch",
				Icon:     "fas fa-truck",
				Link:     "/packing-dispatch-page",
				UserType: basepage.UserType{Admin: true},
				Style:    "font-size:1rem;",
			},
			{
				Name:     "Kanban History",
				Icon:     "fas fa-history",
				Link:     "/kanban-history",
				UserType: basepage.UserType{Admin: true},
				Style:    "font-size:1rem;",
			},
			{
				Name:     "Order History",
				Icon:     "fas fa-file-alt",
				Link:     "/order-history",
				UserType: basepage.UserType{Admin: true},
				Style:    "font-size:1rem;",
			},
			{
				Name:     "Kanban Report",
				Icon:     "fas fa-scroll",
				Link:     "/kanban-report",
				UserType: basepage.UserType{Admin: true},
				Style:    "font-size:1rem;",
			},
			{
				Name:     "Kanban Reprint",
				Icon:     "fas fa-print",
				Link:     "/kanban-reprint",
				UserType: basepage.UserType{Admin: true},
				Style:    "font-size:1rem;",
			},
			{
				Name:     "Summary Reprint",
				Icon:     "fas fa-file",
				Link:     "/summary-reprint",
				UserType: basepage.UserType{Admin: true},
				Style:    "font-size:1rem;",
			},
		},
	}
	sideNav.MenuItems = basepage.CheckdisabledNavItems(sideNav.MenuItems, links, "|")

	// Define the top navigation items
	topNav := basepage.TopNav{VendorName: vendorName, UserType: usertype,
		MenuItems: []basepage.TopNavItem{
			{
				ID:    "settings",
				Name:  "",
				Title: "Settings",
				Type:  "link",
				Icon:  "fa fa-cog",
				Link:  "/configuration-page",
				Width: "col-1",
				Style: "bg-light bg-gradient",
			},
			{
				ID:    "notifications",
				Name:  "",
				Title: "Notifications",
				Type:  "link",
				Icon:  "fa fa-bell",
				Link:  "/system-logs",
				Width: "col-1",
				Style: "bg-light bg-gradient",
			},
			{
				ID:    "username",
				Title: "User Name",
				Name:  username,
				Type:  "button",
				Width: "col-2",
				Style: "bg-light bg-gradient",
			},
			{
				ID:    "logout",
				Title: "Log out",
				Name:  "",
				Link:  "/logout",
				Type:  "link",
				Icon:  "fas fa-sign-out-alt",
				Width: "col-1",
				Style: "bg-light bg-gradient",
			},
		},
	}

	userManagementPage := usermanagmentpage.UserManagement{
		Username: username,
		UserID:   userID,
		UserType: usertype,
	}

	// Build the complete BasePage with navigation and content
	page := (&basepage.BasePage{
		SideNavBar: sideNav,
		TopNavBar:  topNav,
		Username:   username,
		UserType:   usertype,
		Content:    userManagementPage.Build(),
	}).Build()

	// Write the complete HTML page to the response
	w.Write([]byte(page))
}

func UpdateUserDetails(w http.ResponseWriter, r *http.Request) {
	var data m.User
	var ApiResp m.ApiRespMsg
	decoder := json.NewDecoder(r.Body)
	userID := r.Header.Get("X-Custom-Userid")
	if err := decoder.Decode(&data); err == nil {
		if data.Email == "" {
			ApiResp.Code = 400
			ApiResp.Message = "Email can not be empty"
			body, err := json.Marshal(ApiResp)
			if err != nil {
				log.Printf("Error marshaling response: %v", err)
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Write(body)
			return
		}
		data.ModifiedBy = userID
		jsonValue, err := json.Marshal(data)
		if err != nil {
			slog.Error("%s - error - %s", "Error marshaling user data", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to marshal user data")
			return
		}

		resp, err := http.Post(utils.RestURL+"/update-user-details", "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			slog.Error("%s - error - %s", "Error making POST request", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to update user")
			return
		}
		defer resp.Body.Close()

		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			slog.Error("%s - error - %s", "Error reading response body", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to read response body")
			return
		}
		if resp.StatusCode != 200 {
			ApiResp.Code = resp.StatusCode
			ApiResp.Message = string(responseBody)
			body, err := json.Marshal(ApiResp)
			if err != nil {
				log.Printf("Error marshaling response: %v", err)
			}
			w.WriteHeader(resp.StatusCode)
			w.Write(body)
			return
		} else {
			ApiResp.Code = 200
			ApiResp.Message = string(responseBody)
			body, err := json.Marshal(ApiResp)
			if err != nil {
				log.Printf("Error marshaling response: %v", err)
			}
			w.WriteHeader(http.StatusOK)
			w.Write(body)
			return
		}

		// utils.SetResponse(w, http.StatusOK, string(responseBody))
	} else {
		slog.Error("%s - error - %s", "Record updation failed", err.Error())
		utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to update record")
	}
}

func GetUserDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var userManagement []m.UserManagement
	// Make HTTP GET request to DB
	resp, err := http.Get(utils.RestURL + "/get-user-details")
	if err != nil {
		log.Printf("Error making GET request: %v", err)
		http.Error(w, "Failed to get user information i am in restserv", http.StatusForbidden)
		return
	}
	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&userManagement)

	if err := json.NewEncoder(w).Encode(userManagement); err != nil {
		log.Println("Error encoding user details:", err)
		http.Error(w, "Error to retrive user details", http.StatusInternalServerError)
		return
	}
}

func GetUserDetailsByEmail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	email := r.URL.Query().Get("email")
	var userManagement []m.UserManagement
	// Make HTTP GET request to DB
	resp, err := http.Get(utils.RestURL + "/get-user-details-by-email?email=" + email)
	if err != nil {
		log.Printf("Error making GET request: %v", err)
		http.Error(w, "Failed to get user information i am in restserv", http.StatusForbidden)
		return
	}
	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&userManagement)

	if err := json.NewEncoder(w).Encode(userManagement); err != nil {
		log.Println("Error encoding user details:", err)
		http.Error(w, "Error to retrive user details", http.StatusInternalServerError)
		return
	}
}

func usermanagmentDialog(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	userManagementPage := usermanagmentpage.UserManagement{}
	card := userManagementPage.CardBuild(email)
	if err := json.NewEncoder(w).Encode(card); err != nil {
		log.Println("Error encoding user details:", err)
		http.Error(w, "Error to retrive user details", http.StatusInternalServerError)
		return
	}
}
