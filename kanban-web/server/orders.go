package server

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"log/slog"
	"math"
	"net/http"
	"strconv"
	"strings"

	m "irpl.com/kanban-commons/model"
	"irpl.com/kanban-commons/utils"
	"irpl.com/kanban-web/pages/basepage"
	"irpl.com/kanban-web/pages/orders"
	s "irpl.com/kanban-web/services"
)

func vendorOrderPage(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("X-Custom-Username")
	usertype := r.Header.Get("X-Custom-Role")
	userID := r.Header.Get("X-Custom-Userid")
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
				Selected: true,
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
			// {
			// 	Name:  "Production Line Status",
			// 	Icon:  "fas fa-calendar-day",
			// 	Link:  "/flowchart?line=1",
			// 	Style: "font-size:1rem;",
			// },
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
	topNav := basepage.TopNav{VendorName: vendorName,
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

	// Create the FlowchartPage with title, steps, and connections
	vendorOrderPage := orders.Order{
		UserType:   usertype,
		UserID:     userID,
		VendorCode: vendorRecord[0].VendorCode,
	}

	// Build the complete BasePage with navigation and content
	page := (&basepage.BasePage{
		SideNavBar: sideNav,
		TopNavBar:  topNav,
		Username:   username,
		UserType:   usertype,
		Content:    vendorOrderPage.Build(),
	}).Build()

	// Write the complete HTML page to the response
	w.Write([]byte(page))
}

func UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-Custom-Userid")
	var status m.Status
	var kbroot []*m.KbRoot
	var apiresp m.ApiRespMsg
	err := json.NewDecoder(r.Body).Decode(&status)
	if err != nil {
		http.Error(w, "Invalid input format", http.StatusBadRequest)
		return
	}
	var orderdetails []*m.OrderDetails
	var rawQuery m.RawQuery
	rawQuery.Host = utils.RestHost
	rawQuery.Port = utils.RestPort
	rawQuery.Type = "OrderDetails"
	rawQuery.Query = `
				SELECT
				kb_data.*,
				kb_extension.status,
				vendors.vendor_name,
				compounds.compound_name,
				COALESCE(inventory.min_quantity, 0) AS min_quantity,
				COALESCE(inventory.available_quantity, 0) AS available_quantity
				FROM kb_data
				JOIN kb_extension ON kb_extension.id = kb_data.kb_extension_id
				JOIN vendors ON kb_extension.vendor_id = vendors.id
				JOIN compounds ON kb_data.compound_id = compounds.id
				LEFT JOIN inventory ON compounds.id = inventory.compound_id
				Where kb_data.id=` + status.ID + `;`
	rawQuery.RawQry(&orderdetails)

	rawQuery.Type = "KbRoot"
	rawQuery.Host = utils.RestHost
	rawQuery.Port = utils.RestPort
	rawQuery.Query = utils.JoinStr(`				
		SELECT kb_root.* 
		FROM kb_root
		LEFT OUTER JOIN kb_data ON kb_data.id = kb_root.kb_data_id 
		LEFT OUTER JOIN kb_extension ON kb_extension.id = kb_data.kb_extension_id 
		WHERE kb_extension.vendor_id =  (
			SELECT id FROM vendors
			WHERE vendor_code LIKE 'I%' 
			ORDER BY id ASC 
			LIMIT 1
		)
		AND kb_data.compound_id = `, strconv.Itoa(orderdetails[0].CompoundId), ` 
		AND kb_root.status != '3'
		AND kb_root.status != '-1';`)

	rawQuery.RawQry(&kbroot)

	orderdetails[0].InventoryKanbanInProcessQty = len(kbroot)
	if status.Status == "approved" || status.Status == "dispatch" {
		if orderdetails[0].InventoryKanbanInProcessQty > orderdetails[0].MinQuantity {
			if orderdetails[0].AvailableQuantity < orderdetails[0].NoOFLots {
				status.Status = "approved"
				status.Kanban = int(math.Abs(float64(orderdetails[0].MinQuantity - orderdetails[0].AvailableQuantity)))
				status.Dispatch = orderdetails[0].AvailableQuantity
			} else if orderdetails[0].AvailableQuantity >= orderdetails[0].NoOFLots {
				status.Status = "dispatch"
				status.Dispatch = orderdetails[0].NoOFLots
			}
		} else {
			if orderdetails[0].NoOFLots <= orderdetails[0].AvailableQuantity {
				if orderdetails[0].AvailableQuantity >= orderdetails[0].MinQuantity {
					status.Status = "dispatch"
					status.Dispatch = orderdetails[0].NoOFLots
				} else {
					if orderdetails[0].InventoryKanbanInProcessQty < (orderdetails[0].MinQuantity - (orderdetails[0].NoOFLots - orderdetails[0].AvailableQuantity)) {
						status.Status = "approved"
						status.Kanban = orderdetails[0].MinQuantity - (orderdetails[0].NoOFLots - orderdetails[0].AvailableQuantity) - orderdetails[0].InventoryKanbanInProcessQty
					} else {
						status.Status = "dispatch"
					}

					status.Dispatch = orderdetails[0].NoOFLots

				}
			} else if orderdetails[0].AvailableQuantity == 0 {
				status.Status = "approved"
				status.Kanban = (orderdetails[0].MinQuantity + orderdetails[0].NoOFLots) - orderdetails[0].InventoryKanbanInProcessQty
			} else {
				status.Status = "approved"
				status.Dispatch = orderdetails[0].AvailableQuantity
				status.Kanban = (orderdetails[0].MinQuantity + (orderdetails[0].NoOFLots - orderdetails[0].AvailableQuantity)) - orderdetails[0].InventoryKanbanInProcessQty
			}
			status.CompoundID = orderdetails[0].CompoundId
			status.NoOFLots = orderdetails[0].NoOFLots

		}
	} else {
		status.Status = "reject"
		status.CompoundID = orderdetails[0].CompoundId
		status.NoOFLots = orderdetails[0].NoOFLots
	}

	status.UserID = userID
	data, err := json.Marshal(status)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	url := utils.RestURL + "/update-order-status"
	res, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		http.Error(w, "Failed to update model license data", http.StatusForbidden)
		return
	}
	defer res.Body.Close()

	client := &http.Client{}
	resp, err := client.Do(res)
	if resp.StatusCode != http.StatusOK || err != nil {

		body, _ := io.ReadAll(resp.Body) // Read response body
		apiresp.Code = resp.StatusCode
		apiresp.Message = string(body)
		respbody, _ := json.Marshal(apiresp)
		http.Error(w, string(respbody), resp.StatusCode)
		return
	}
	// Write success status and message
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	apiresp.Code = 200
	apiresp.Message = "Order Status updated successfully"

	if err := json.NewEncoder(w).Encode(apiresp); err != nil {
		http.Error(w, "Failed to encode success message", http.StatusInternalServerError)
	}
}

func OrderDetailsForCustomerHandler(w http.ResponseWriter, r *http.Request) {
	// Validate the request method
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Initialize the OrderDetails structure to hold incoming JSON payload
	var OrderDetails orders.OrderDetails

	// Decode the JSON payload from the request body
	if err := json.NewDecoder(r.Body).Decode(&OrderDetails); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// Prepare the REST API URL for an external request (if necessary)
	url := orders.RestURL + "/OrderDetailsForCustomer"

	// Create a new HTTP POST request to the external service
	requestBody, err := json.Marshal(OrderDetails)
	if err != nil {
		http.Error(w, "Failed to encode request body", http.StatusInternalServerError)
		return
	}

	externalReq, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestBody))
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}
	externalReq.Header.Set("Content-Type", "application/json")

	// Execute the HTTP request using an HTTP client
	client := &http.Client{}
	externalResp, err := client.Do(externalReq)
	if err != nil {
		http.Error(w, "Failed to execute request", http.StatusInternalServerError)
		return
	}
	defer externalResp.Body.Close()
	var OrderData m.OrderDetails
	if err := json.NewDecoder(externalResp.Body).Decode(&OrderData); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}
	// Handle non-200 responses from the external service
	if externalResp.StatusCode != http.StatusOK {
		http.Error(w, "External service returned an error", http.StatusInternalServerError)
		return
	}

	var stage string = "Approved"
	if OrderData.Status == "pending" {
		stage = "Pending"
	} else if OrderData.Status == "reject" {
		stage = "Rejected"
	}
	OrderDetails.Status = OrderData.Status
	OrderDetails.OrderId = OrderData.OrderId
	// Add ProductionProcess with temporary test data
	OrderDetails.ProductionProcess = []m.ProdProcess{
		{
			Name:        "Created",
			Description: "You have created an order",
			Status:      "1",
		},
		{
			Name:        stage,
			Description: "Your order is " + stage,
			Status:      "2",
		},
		{
			Name:        "Line Up",
			Description: "Your order is Lined-up",
			Status:      "3",
		},
		{
			Name:        "In production Line",
			Description: "Your order is in Production Line",
			Status:      "4",
		},
		{
			Name:        "Quality Testing",
			Description: "Your order is in Quality Testing",
			Status:      "5",
		},
		{
			Name:        "Packing",
			Description: "Your order is in packing stage",
			Status:      "6",
		},
		{
			Name:        "Order Dispatch",
			Description: "Your order is Dispatched",
			Status:      "7",
		},
	}
	// Generate the HTML content
	Html := OrderDetails.Build() // Assuming this generates an HTML string
	// Send the HTML inside a JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{"html": Html}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// Admin order page
func adminOrderPage(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("X-Custom-Username")
	usertype := r.Header.Get("X-Custom-Role")
	userID := r.Header.Get("X-Custom-Userid")
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
				Selected: true,
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

	// Create the FlowchartPage with title, steps, and connections
	adminOrderPage := orders.AdminOrder{
		UserType: usertype,
		UserID:   userID,
	}

	// Build the complete BasePage with navigation and content
	page := (&basepage.BasePage{
		SideNavBar: sideNav,
		TopNavBar:  topNav,
		Username:   username,
		UserType:   usertype,
		Content:    adminOrderPage.Build(),
	}).Build()

	// Write the complete HTML page to the response
	w.Write([]byte(page))
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

		resp, err := http.Post(utils.RestURL+"/get-order-details", "application/json", bytes.NewBuffer(jsonValue))
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

func SearchOrderDEtails(w http.ResponseWriter, r *http.Request) {
	type SearchRequest struct {
		Criteria map[string]string `json:"Criteria"`
	}
	type TableRequest struct {
		Conditions []string `json:"conditions"`
	}
	var data SearchRequest
	var searchresult []*m.OrderDetails
	var tablereq TableRequest
	var subcondition []string
	var condition string
	var searcondition []string
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err == nil {
		for key, val := range data.Criteria {
			if key == "vendor_name" && val != "" {
				con := "vendors." + key + " iLIKE '%%%" + val + "%%' "
				subcondition = append(subcondition, con)
			} else if key == "compound_name" && val != "" {
				con := "compounds." + key + " iLIKE '%%%" + val + "%%' "
				subcondition = append(subcondition, con)
			} else if key == "cell_no" && val != "" {
				con := "kb_data." + key + " iLIKE '%%%" + val + "%%' "
				subcondition = append(subcondition, con)
			}
		}
		con := ` kb_extension.status='pending'`
		subcondition = append(subcondition, con)
		for i, v := range subcondition {
			if i < (len(subcondition) - 1) {
				condition = condition + v + " AND "
			} else {
				condition = condition + v
			}
		}

		searcondition = append(searcondition, condition)
		tablereq.Conditions = searcondition
		jsonValue, err := json.Marshal(tablereq)
		if err != nil {
			slog.Error("%s - error - %s", "Error marshaling table data", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to marshal table condition")
			return
		}

		resp, err := http.Post(utils.RestURL+"/get-order-details", "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			slog.Error("%s - error - %s", "Error making POST request", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to fectch order details")
			return
		}
		defer resp.Body.Close()

		decoder := json.NewDecoder(resp.Body)
		if err := decoder.Decode(&searchresult); err == nil {

			var odcard s.TableCard
			tablebutton := `
	<!--html-->
		<button type="button" class="btn m-0 p-0" id="viewOrderModel" data-toggle="tooltip" data-placement="bottom" title="View Details"  data-bs-toggle="modal" data-bs-target="#UpdateOrderStatus"> 
				 <i class="fa fa-eye mx-2" style="color: #b250ad;"></i> 
		</button>
	<!--!html-->`

			odcard.BodyTables.Data = searchresult
			odcard.BodyTables.Columns = []s.CardTableBodyHeadCol{
				{
					Name:       `Sr. No.`,
					IsCheckbox: true,
					ID:         "search-sn",
					Width:      "col-1",
				},
				{
					Name:         "Vendor Name",
					ID:           "VendorName",
					IsSearchable: true,
					Type:         "input",
					Width:        "col-1",
				},
				{
					Lable:        "Part Name",
					Name:         "Compound Name",
					ID:           "CompoundName",
					IsSearchable: true,
					Type:         "input",
					Width:        "col-1",
					// IsSortable:   true,
				},
				{
					Name:         "Cell No",
					ID:           "CellNo",
					IsSearchable: true,
					Type:         "input",
					Width:        "col-1",
				},
				{
					Name:  "Demand Date Time",
					ID:    "demand_date",
					Width: "col-1",
				},
				{
					Lable:  "No. of Lots",
					Name:   "NoOFLots",
					ID:     "lot_no",
					Width:  "col-1",
					GetSum: true,
				},
				{
					Name:  "Status",
					Width: "col-1",
				},
				{
					Name:  "Action",
					Width: "col-1",
				},
			}
			odcard.BodyTables.ColumnsWidth = []string{"col-1", "col-1", "col-1", "col-1", "col-1", "col-1", "col-1", "col-1"}
			odcard.BodyTables.Buttons = tablebutton
			tbody := odcard.BodyTables.OrderDetailsTable()
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(tbody))

		} else {
			slog.Error("Record creation failed - " + err.Error())
			utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to create Pagination")
		}

	} else {

		utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to fetch order record")
		slog.Error("%s - error - %s", "failed to fetch order record", err.Error())
	}

}

func GetAllOrderByVendorCode(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	url := RestURL + "/get-all-orders-by-vendor-code"
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
	var ApiResp m.ApiRespMsg
	ApiResp.Code = resp.StatusCode
	ApiResp.Message = string(responseBody)
	apiResp, err := json.Marshal(ApiResp)
	if err != nil {
		log.Printf("Error marshaling response: %v", err)
	}
	utils.SetResponse(w, resp.StatusCode, string(apiResp))
}

func GetAllPendingOrders(w http.ResponseWriter, r *http.Request) {
	type TableConditions struct {
		Conditions []string `json:"Conditions"`
	}

	var tablecondition TableConditions
	con := utils.JoinStr(`kb_extension.status='pending'`)
	tablecondition.Conditions = append(tablecondition.Conditions, con)

	jsonValue, err := json.Marshal(tablecondition)
	if err != nil {
		slog.Error("%s - error - %s", "Error marshaling table condition", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	resp, err := http.Post(RestURL+"/get-order-details", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		slog.Error("%s - error - %s", "Error making GET request", err)
		http.Error(w, "Failed to fetch data", http.StatusBadGateway)
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

func RejectOrder(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-Custom-Userid")
	var status m.Status
	var kbroot []*m.KbRoot
	var apiresp m.ApiRespMsg
	err := json.NewDecoder(r.Body).Decode(&status)
	if err != nil {
		http.Error(w, "Invalid input format", http.StatusBadRequest)
		return
	}
	if status.ID == "" {
		http.Error(w, "Missing or invalid status ID", http.StatusBadRequest)
		return
	}
	var orderdetails []*m.OrderDetails
	var rawQuery m.RawQuery
	rawQuery.Host = utils.RestHost
	rawQuery.Port = utils.RestPort
	rawQuery.Type = "OrderDetails"
	rawQuery.Query = `
				SELECT
				kb_data.*,
				kb_extension.status,
				vendors.vendor_name,
				compounds.compound_name,
				COALESCE(inventory.min_quantity, 0) AS min_quantity,
				COALESCE(inventory.available_quantity, 0) AS available_quantity
				FROM kb_data
				JOIN kb_extension ON kb_extension.id = kb_data.kb_extension_id
				JOIN vendors ON kb_extension.vendor_id = vendors.id
				JOIN compounds ON kb_data.compound_id = compounds.id
				LEFT JOIN inventory ON compounds.id = inventory.compound_id
				Where kb_data.id=` + status.ID + `;`
	rawQuery.RawQry(&orderdetails)
	if len(orderdetails) == 0 {
		http.Error(w, "No order details found", http.StatusNotFound)
		return
	}

	rawQuery.Type = "KbRoot"
	rawQuery.Host = utils.RestHost
	rawQuery.Port = utils.RestPort
	rawQuery.Query = utils.JoinStr(`				
		SELECT kb_root.* 
		FROM kb_root
		LEFT OUTER JOIN kb_data ON kb_data.id = kb_root.kb_data_id 
		LEFT OUTER JOIN kb_extension ON kb_extension.id = kb_data.kb_extension_id 
		WHERE kb_extension.vendor_id = (
			SELECT id FROM vendors
			WHERE vendor_code LIKE 'I%' 
			ORDER BY id ASC 
			LIMIT 1
		)
		AND kb_data.compound_id = `, strconv.Itoa(orderdetails[0].CompoundId), ` 
		AND kb_root.status != '3';`)

	rawQuery.RawQry(&kbroot)

	orderdetails[0].InventoryKanbanInProcessQty = len(kbroot)
	status.Status = "reject"
	status.CompoundID = orderdetails[0].CompoundId
	status.NoOFLots = orderdetails[0].NoOFLots

	status.UserID = userID
	data, err := json.Marshal(status)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	url := utils.RestURL + "/update-order-status"
	res, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		http.Error(w, "Failed to update model license data", http.StatusForbidden)
		return
	}
	defer res.Body.Close()

	client := &http.Client{}
	resp, err := client.Do(res)
	if resp.StatusCode != http.StatusOK || err != nil {

		body, _ := io.ReadAll(resp.Body) // Read response body
		apiresp.Code = resp.StatusCode
		apiresp.Message = string(body)
		respbody, _ := json.Marshal(apiresp)
		http.Error(w, string(respbody), resp.StatusCode)
		return
	}
	// Write success status and message
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	apiresp.Code = 200
	apiresp.Message = "Order Status updated successfully"

	if err := json.NewEncoder(w).Encode(apiresp); err != nil {
		http.Error(w, "Failed to encode success message", http.StatusInternalServerError)
	}
}

func ApproveOrder(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-Custom-Userid")
	var status m.Status
	var kbroot []*m.KbRoot
	var apiresp m.ApiRespMsg
	err := json.NewDecoder(r.Body).Decode(&status)
	if err != nil {
		http.Error(w, "Invalid input format", http.StatusBadRequest)
		return
	}
	if status.ID == "" {
		http.Error(w, "Missing or invalid status ID", http.StatusBadRequest)
		return
	}

	var orderdetails []*m.OrderDetails
	var rawQuery m.RawQuery
	rawQuery.Host = utils.RestHost
	rawQuery.Port = utils.RestPort
	rawQuery.Type = "OrderDetails"
	rawQuery.Query = `
					SELECT
					kb_data.*,
					kb_extension.status,
					vendors.vendor_name,
					compounds.compound_name,
					COALESCE(inventory.min_quantity, 0) AS min_quantity,
					COALESCE(inventory.available_quantity, 0) AS available_quantity
					FROM kb_data
					JOIN kb_extension ON kb_extension.id = kb_data.kb_extension_id
					JOIN vendors ON kb_extension.vendor_id = vendors.id
					JOIN compounds ON kb_data.compound_id = compounds.id
					LEFT JOIN inventory ON compounds.id = inventory.compound_id
					Where kb_data.id=` + status.ID + `;`
	rawQuery.RawQry(&orderdetails)
	if len(orderdetails) == 0 {
		http.Error(w, "No order details found", http.StatusNotFound)
		return
	}

	rawQuery.Type = "KbRoot"
	rawQuery.Host = utils.RestHost
	rawQuery.Port = utils.RestPort
	rawQuery.Query = utils.JoinStr(`				
			SELECT kb_root.* 
			FROM kb_root
			LEFT OUTER JOIN kb_data ON kb_data.id = kb_root.kb_data_id 
			LEFT OUTER JOIN kb_extension ON kb_extension.id = kb_data.kb_extension_id 
			WHERE kb_extension.vendor_id =(
				SELECT id FROM vendors
				WHERE vendor_code LIKE 'I%' 
				ORDER BY id ASC 
				LIMIT 1
			)
			AND kb_data.compound_id = `, strconv.Itoa(orderdetails[0].CompoundId), ` 
			AND kb_root.status != '3';`)

	rawQuery.RawQry(&kbroot)

	orderdetails[0].InventoryKanbanInProcessQty = len(kbroot)

	if orderdetails[0].InventoryKanbanInProcessQty > orderdetails[0].MinQuantity {
		if orderdetails[0].AvailableQuantity < orderdetails[0].NoOFLots {
			status.Status = "approved"
			status.Kanban = int(math.Abs(float64(orderdetails[0].MinQuantity - orderdetails[0].AvailableQuantity)))
			status.Dispatch = orderdetails[0].AvailableQuantity
		} else if orderdetails[0].AvailableQuantity >= orderdetails[0].NoOFLots {
			status.Status = "dispatch"
			status.Dispatch = orderdetails[0].NoOFLots
		}
	} else {
		if orderdetails[0].NoOFLots <= orderdetails[0].AvailableQuantity {
			if orderdetails[0].AvailableQuantity >= orderdetails[0].MinQuantity {
				status.Status = "dispatch"
				status.Dispatch = orderdetails[0].NoOFLots
			} else {
				if orderdetails[0].InventoryKanbanInProcessQty < (orderdetails[0].MinQuantity - (orderdetails[0].NoOFLots - orderdetails[0].AvailableQuantity)) {
					status.Status = "approved"
					status.Kanban = orderdetails[0].MinQuantity - (orderdetails[0].NoOFLots - orderdetails[0].AvailableQuantity) - orderdetails[0].InventoryKanbanInProcessQty
				} else {
					status.Status = "dispatch"
				}

				status.Dispatch = orderdetails[0].NoOFLots

			}
		} else if orderdetails[0].AvailableQuantity == 0 {
			status.Status = "approved"
			status.Kanban = (orderdetails[0].MinQuantity + orderdetails[0].NoOFLots) - orderdetails[0].InventoryKanbanInProcessQty
		} else {
			status.Status = "approved"
			status.Dispatch = orderdetails[0].AvailableQuantity
			status.Kanban = (orderdetails[0].MinQuantity + (orderdetails[0].NoOFLots - orderdetails[0].AvailableQuantity)) - orderdetails[0].InventoryKanbanInProcessQty
		}
		status.CompoundID = orderdetails[0].CompoundId
		status.NoOFLots = orderdetails[0].NoOFLots

	}
	status.UserID = userID
	data, err := json.Marshal(status)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	url := utils.RestURL + "/update-order-status"
	res, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		http.Error(w, "Failed to update model license data", http.StatusForbidden)
		return
	}
	defer res.Body.Close()

	client := &http.Client{}
	resp, err := client.Do(res)
	if resp.StatusCode != http.StatusOK || err != nil {

		body, _ := io.ReadAll(resp.Body) // Read response body
		apiresp.Code = resp.StatusCode
		apiresp.Message = string(body)
		respbody, _ := json.Marshal(apiresp)
		http.Error(w, string(respbody), resp.StatusCode)
		return
	}
	// Write success status and message
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	apiresp.Code = 200
	apiresp.Message = "Order Status updated successfully"

	if err := json.NewEncoder(w).Encode(apiresp); err != nil {
		http.Error(w, "Failed to encode success message", http.StatusInternalServerError)
	}
}

func PendingOrderSearchPagination(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Pagination m.PaginationReq `json:"pagination"`
		Conditions []string        `json:"Conditions"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var Response struct {
		Pagination m.PaginationResp
		Data       []*m.OrderDetails
	}

	jsonValue, err := json.Marshal(req)
	if err != nil {
		slog.Error("Error marshaling table condition", "error", err)
	}

	resp, err := http.Post(utils.RestURL+"/get-order-pending-details-by-search-pagination", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		slog.Error("%s - error - %s", "Error making GET request", err)
	}

	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&Response); err != nil {
		slog.Error("error decoding response body", "error", err)
	}

	// Extract DataField and Value
	searchFilters := make(map[string]string)
	for _, condition := range req.Conditions {
		parts := strings.Split(condition, " ")
		if len(parts) == 3 {
			dataField := parts[0]
			value := strings.Trim(parts[2], "'%") // Remove % from start and end
			searchFilters[dataField] = value
		}
	}

	var vendorOrderTable s.TableCard
	tablebutton := `
	<!--html-->
		<button type="button" class="btn m-0 p-0" id="viewOrderModel" data-toggle="tooltip" data-placement="bottom" title="View Details"  data-bs-toggle="modal" data-bs-target="#UpdateOrderStatus"> 
				 <i class="fa fa-eye mx-2" style="color: #b250ad;"></i> 
		</button>
	<!--!html-->`
	vendorOrderTable.CardHeading = "Customer Orders"
	vendorOrderTable.BodyTables = s.CardTableBody{Columns: []s.CardTableBodyHeadCol{
		{
			Name:       `Sr. No.`,
			IsCheckbox: true,
			ID:         "search-sn",
			Width:      "col-1",
		},
		{
			Name:         "Vendor Name",
			ID:           "VendorName",
			IsSearchable: true,
			Type:         "input",
			Width:        "col-1",
		},
		{
			Lable:        "Part Name",
			Name:         "Compound Name",
			ID:           "CompoundName",
			IsSearchable: true,
			Type:         "input",
			Width:        "col-1",
			IsSortable:   true,
		},
		{
			Lable:        "Kanban Summary",
			Name:         "Cell No",
			ID:           "CellNo",
			IsSearchable: true,
			Type:         "input",
			Width:        "col-1",
		},
		{
			Name:  "Demand Date Time",
			ID:    "demand_date",
			Width: "col-1",
		},
		{
			Lable:  "No. of Lots",
			Name:   "NoOFLots",
			ID:     "lot_no",
			Width:  "col-1",
			GetSum: true,
		},
		{
			Name:  "Status",
			Width: "col-1",
		},
		{
			Name:  "Action",
			Width: "col-1",
		},
	},
		ColumnsWidth: []string{"col-1", "col-1", "col-1", "col-1", "col-1", "col-1", "col-1", "col-1"},
		Data:         Response.Data,
		Buttons:      tablebutton,
	}

	tableBodyHTML := vendorOrderTable.BodyTables.RenderBodyColumns()

	var Pagination s.Pagination
	Pagination.TotalRecords = Response.Pagination.TotalNo
	Pagination.Offset = Response.Pagination.Offset
	Pagination.CurrentPage = Response.Pagination.Page
	Pagination.PerPage, _ = strconv.Atoi(req.Pagination.Limit)
	Pagination.PerPageOptions = []int{10, 25, 50, 100, 200}

	response := map[string]any{
		"tableBodyHTML":  tableBodyHTML,
		"paginationHTML": Pagination.Build(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}
