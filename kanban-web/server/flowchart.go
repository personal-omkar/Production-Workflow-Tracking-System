package server

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"log/slog"
	"net/http"
	"strconv"

	m "irpl.com/kanban-commons/model"
	u "irpl.com/kanban-commons/utils"
	"irpl.com/kanban-web/pages/basepage"
	"irpl.com/kanban-web/pages/flowchart"
	"irpl.com/kanban-web/services"
)

func flowChartPage(w http.ResponseWriter, r *http.Request) {
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

	queryParams := r.URL.Query()
	line_no := queryParams.Get("line")
	if line_no == "" {
		line_no = "1"
	}
	lineNo, _ := strconv.Atoi(line_no)
	payload := map[string]int{
		"line_no": lineNo,
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling payload: %v", err)
		return
	}

	// Construct the REST API URL
	url := u.JoinStr("http://", u.RestHost, ":", u.RestPort, "/get-production-line-status")
	// Create the HTTP request
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(jsonPayload))
	if err != nil {
		http.Error(w, "Error creating request to REST service", http.StatusInternalServerError)
		return
	}
	// Send the request using an HTTP client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to send request to REST service", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close() // Ensure the response body is closed
	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Fail to get data", resp.StatusCode)
		return
	}

	var prodLineStatus []m.ProdLineDetails
	json.NewDecoder(resp.Body).Decode(&prodLineStatus)

	var steps []services.StepCard
	RestURL := u.JoinStr("http://", u.RestHost, ":", u.RestPort)
	processresp, err := http.Get(RestURL + "/get-all-production-process-for-line?line=" + line_no)
	if err != nil {
		slog.Error("%s - error - %s", "Error making GET request", err)
	}

	defer processresp.Body.Close()

	if err := json.NewDecoder(processresp.Body).Decode(&steps); err != nil {
		slog.Error("error decoding response body", "error", err)
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
				Name:     u.DefaultsMap["cold_store_menu"],
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
				Selected: true,
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
				Type:  "link",
				Icon:  "fa fa-cog",
				Link:  "/configuration-page",
				Width: "col-1",
				Style: "bg-light bg-gradient",
			},
			{
				ID:    "notifications",
				Name:  "",
				Type:  "link",
				Icon:  "fa fa-bell",
				Link:  "/system-logs",
				Width: "col-1",
				Style: "bg-light bg-gradient",
			},
			{
				ID:    "username",
				Name:  username,
				Type:  "button",
				Width: "col-2",
				Style: "bg-light bg-gradient",
			},
			{
				ID:    "logout",
				Name:  "",
				Link:  "/logout",
				Type:  "link",
				Icon:  "fas fa-sign-out-alt",
				Width: "col-1",
				Style: "bg-light bg-gradient",
			},
		},
	}

	// Define the steps for the flowchart as StepCards
	// steps := []services.StepCard{
	// 	{IconFilename: "chemical_batching.png", Text: "Chemical Batching", Description: "Description for Chemical Batching."},
	// 	{IconFilename: "rubber_batching.png", Text: "Rubber Batching", Description: "Description for Rubber Batching."},
	// 	{IconFilename: "accelerator_batching.png", Text: "Accelerator Batching", Description: "Description for Accelerator Batching."},
	// 	{IconFilename: "carbon_batching.png", Text: "Carbon Batching", Description: "Description for Carbon Batching."},
	// 	{IconFilename: "oil_batching.png", Text: "Oil Batching", Description: "Description for Oil Batching."},
	// 	{IconFilename: "intermix.png", Text: "Intermix", Description: "Description for Intermix."},
	// 	{IconFilename: "mill1.png", Text: "Mill 1", Description: "Description for Mill 1."},
	// 	{IconFilename: "carbon_batching.png", Text: "Carbon Batching", Description: "Description for Carbon Batching."},
	// }

	// // Define the connections between steps (pairs of step indices)
	// connections := [][2]int{
	// 	{0, 1}, {1, 2}, {2, 3}, {3, 4},
	// 	{4, 5}, {5, 6}, {6, 7}, {7, 8},
	// 	{8, 9}, {9, 10},
	// }

	// Create the FlowchartPage with title, steps, and connections
	var Name string
	var CellData []m.Cell
	for _, data := range prodLineStatus {
		Name = data.ProdLineName
		CellData = data.Cells
	}
	flowchartPage := flowchart.NewFlowchartPage(Name, steps, CellData)
	// Build the flowchart content for the BasePage
	flowchartContent := flowchartPage.Build()

	// Build the complete BasePage with navigation and content
	page := (&basepage.BasePage{
		SideNavBar: sideNav,
		TopNavBar:  topNav,
		Username:   username,
		UserType:   usertype,
		Content:    flowchartContent,
	}).Build()

	// Write the complete HTML page to the response
	w.Write([]byte(page))
}

// Update Flowchart page for each 5 sec
func UpdateFlowchart(w http.ResponseWriter, r *http.Request) {

	queryParams := r.URL.Query()
	line_no := queryParams.Get("line")
	if line_no == "" {
		line_no = "1"
	}

	lineNo, _ := strconv.Atoi(line_no)
	payload := map[string]int{
		"line_no": lineNo,
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling payload: %v", err)
		return
	}

	// Construct the REST API URL
	url := u.JoinStr("http://", u.RestHost, ":", u.RestPort, "/get-production-line-status")
	// Create the HTTP request
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(jsonPayload))
	if err != nil {
		http.Error(w, "Error creating request to REST service", http.StatusInternalServerError)
		return
	}
	// Send the request using an HTTP client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to send request to REST service", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close() // Ensure the response body is closed
	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Fail to get data", resp.StatusCode)
		return
	}

	var prodLineStatus []m.ProdLineDetails
	json.NewDecoder(resp.Body).Decode(&prodLineStatus)

	var steps []services.StepCard
	RestURL := u.JoinStr("http://", u.RestHost, ":", u.RestPort)
	processresp, err := http.Get(RestURL + "/get-all-production-process-for-line?line=" + line_no)
	if err != nil {
		slog.Error("%s - error - %s", "Error making GET request", err)
	}

	defer processresp.Body.Close()

	if err := json.NewDecoder(processresp.Body).Decode(&steps); err != nil {
		slog.Error("error decoding response body", "error", err)
	}
	// Define the steps for the flowchart as StepCards
	// steps := []services.StepCard{
	// 	{IconFilename: "chemical_batching.png", Text: "Chemical Batching", Description: "Description for Chemical Batching."},
	// 	{IconFilename: "rubber_batching.png", Text: "Rubber Batching", Description: "Description for Rubber Batching."},
	// 	{IconFilename: "accelerator_batching.png", Text: "Accelerator Batching", Description: "Description for Accelerator Batching."},
	// 	{IconFilename: "carbon_batching.png", Text: "Carbon Batching", Description: "Description for Carbon Batching."},
	// 	{IconFilename: "oil_batching.png", Text: "Oil Batching", Description: "Description for Oil Batching."},
	// 	{IconFilename: "intermix.png", Text: "Intermix", Description: "Description for Intermix."},
	// 	{IconFilename: "mill1.png", Text: "Mill 1", Description: "Description for Mill 1."},
	// 	{IconFilename: "carbon_batching.png", Text: "Carbon Batching", Description: "Description for Carbon Batching."},
	// }

	// Create the FlowchartPage with title, steps, and connections
	var Name string
	var CellData []m.Cell
	for _, data := range prodLineStatus {
		Name = data.ProdLineName
		CellData = data.Cells
	}
	flowchartPage := flowchart.NewFlowchartPage(Name, steps, CellData)
	// Build the flowchart content for the BasePage
	flowchartContent := flowchartPage.Build()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"content": flowchartContent,
	}); err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}
}

func CreateNewKbTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID := r.Header.Get("X-Custom-Userid")
	var transactionData m.KbTransaction
	err := json.NewDecoder(r.Body).Decode(&transactionData)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	transactionData.CreatedBy = userID

	marshalData, err := json.Marshal(transactionData)
	if err != nil {
		http.Error(w, "Error while marshaling data", http.StatusInternalServerError)
		return
	}

	// Construct the REST API URL
	url := u.JoinStr("http://", u.RestHost, ":", u.RestPort, "/create-new-KbTransaction")

	// Create the HTTP request
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(marshalData))
	if err != nil {
		http.Error(w, "Error creating request to REST service", http.StatusInternalServerError)
		return
	}

	// Send the request using an HTTP client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to send request to REST service", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close() // Ensure the response body is closed

	// Check for non-201 status codes
	if resp.StatusCode != http.StatusCreated {
		http.Error(w, "Fail to create New Line", resp.StatusCode)
		return
	}
	// Parse the response from the REST service

	// Return success response to the client
	w.WriteHeader(http.StatusOK)
	// json.NewEncoder(w).Encode(responsePayload)
}

func UpdateRunningNumberAfterTransactioin(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	url := u.RestURL + "/update-running-number"
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
	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Failed to read response body", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
