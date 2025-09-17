package server

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	m "irpl.com/kanban-commons/model"
	"irpl.com/kanban-commons/utils"
	"irpl.com/kanban-web/pages/basepage"
	"irpl.com/kanban-web/pages/coldstorage"
	s "irpl.com/kanban-web/services"
)

var RestHost string // Global variable to hold the Rest helper host
var RestPort string // Global variable to hold the Rest helper port
var RestURL string  // Global variable to hold the Rest URL

func init() {
	RestHost = os.Getenv("RESTSRV_HOST")
	if strings.TrimSpace(RestHost) == "" {
		RestHost = utils.DefaultRestHost
	}

	RestPort = os.Getenv("RESTSRV_PORT")
	if strings.TrimSpace(RestPort) == "" {
		RestPort = utils.DefaultRestPort
	}

	RestURL = utils.JoinStr("http://", RestHost, ":", RestPort)
}
func ColdStoragePage(w http.ResponseWriter, r *http.Request) {
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
				Selected: true,
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

	// Create the FlowchartPage with title, steps, and connections
	coldstorage := coldstorage.ColdStoragePage{}

	// Build the flowchart content for the BasePage

	// Build the complete BasePage with navigation and content
	page := (&basepage.BasePage{
		SideNavBar: sideNav,
		TopNavBar:  topNav,
		Username:   username,
		UserType:   usertype,
		Content:    coldstorage.Build(),
	}).Build()

	// Write the complete HTML page to the response
	w.Write([]byte(page))
}

// Create compound by vendor name
func UpdateColdStorageQuantity(w http.ResponseWriter, r *http.Request) {
	var data m.Inventory
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	values, err := url.ParseQuery(string(body))
	if err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	compoundName := values.Get("editcompoundname")
	data.Id, _ = strconv.Atoi(values.Get("editinventoryid"))
	data.MinQuantity, _ = strconv.Atoi(values.Get("editminqty"))
	data.MaxQuantity, _ = strconv.Atoi(values.Get("editmaxqty"))
	jsonValue, err := json.Marshal(data)

	if err != nil {
		slog.Error("%s - error - %s", "Error marshaling user data", err)
		utils.SetResponse(w, http.StatusInternalServerError, "Failed to marshal data")
		return
	}
	resp, err := http.Post(utils.RestURL+"/update-coldstorage-quantity", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		slog.Error("%s - error - %s", "Error making POST request", err)
		utils.SetResponse(w, http.StatusInternalServerError, "Failed to create Compound Entry")
		return
	}

	if resp.StatusCode == http.StatusOK {
		http.Redirect(w, r, "/cold-storage?status=200&msg="+compoundName+" Quantity updated successfully", http.StatusFound)

	} else {
		http.Redirect(w, r, "/cold-storage?status=500&msg=Faild to update quantity of "+compoundName, http.StatusFound)
	}
}

func ColdStorageSearchPagination(w http.ResponseWriter, r *http.Request) {
	// Parse query param
	showDelete := r.URL.Query().Get("showDelete") == "true"
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
		Data       []*m.ColdStorage
	}

	jsonValue, err := json.Marshal(req)
	if err != nil {
		slog.Error("Error marshaling table condition", "error", err)
	}

	// Send POST request
	resp, err := http.Post(RestURL+"/get-all-cold-storage-by-search-pagination", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		slog.Error("Error making POST request", "error", err)
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

	var coldstorageTable s.TableCard
	tableTools := `<button type="button" class="btn m-0 p-0" data-bs-toggle="modal" data-bs-target="#EditCompound">
		<i class="fa fa-edit" style="color: #CF7AC2;"></i>
	</button>`

	if showDelete {
		tableTools += `<button type="button" class="btn mx-1 p-0" id="del-part">
			<i class="fa fa-trash" style="color:rgb(207, 97, 97);"></i>
		</button>`
	}

	coldstorageTable.CardHeading = "Rubber Store Master"
	coldstorageTable.CardHeadingActions.Component = []s.CardHeadActionComponent{
		{ComponentName: "modelbutton", ComponentType: s.ActionComponentElement{ModelButton: s.ModelButtonAttributes{Type: "button", Text: "Add Compound to Inventory", ModelID: "#AddParToInventoryModel"}}, Width: "col-4"}}

	coldstorageTable.BodyTables = s.CardTableBody{Columns: []s.CardTableBodyHeadCol{
		{
			Lable:        "Part Name",
			Name:         `Compound Name`,
			ID:           "search-compound-name",
			IsSearchable: true,
			Type:         "input",
			Width:        "col-2",
		},
		{
			Name:  "Max Quantity",
			ID:    "search-min-quantity",
			Type:  "input",
			Width: "col-2",
		},
		{
			Name:  "Min Quantity",
			Type:  "input",
			ID:    "search-max-quantity",
			Width: "col-2",
		},
		{
			Name:  "Available Quantity",
			Type:  "input",
			ID:    "search-max-quantity",
			Width: "col-2",
		},

		{
			Name:  "Tools",
			Width: "col-1",
		},
	},
		ColumnsWidth: []string{"col-2", "col-2", "col-2", "col-2", "col-1"},
		Data:         Response.Data,
		Tools:        tableTools,
		ID:           "Inventory-Table",
	}

	tableBodyHTML := coldstorageTable.BodyTables.RenderBodyColumns()

	var Pagination s.Pagination
	Pagination.TotalRecords = Response.Pagination.TotalNo
	Pagination.CurrentPage = req.Pagination.PageNo
	Pagination.PerPage, _ = strconv.Atoi(req.Pagination.Limit)
	Pagination.Offset = (Pagination.CurrentPage - 1) * Pagination.PerPage
	Pagination.PerPage, _ = strconv.Atoi(req.Pagination.Limit)
	Pagination.PerPageOptions = []int{10, 25, 50, 100, 200}

	response := map[string]any{
		"tableBodyHTML":  tableBodyHTML,
		"paginationHTML": Pagination.Build(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

// Create compound by vendor name
func UpdateColdStorageQuantityForInventoryManagement(w http.ResponseWriter, r *http.Request) {
	var data m.Inventory
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	values, err := url.ParseQuery(string(body))
	if err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	compoundName := values.Get("editcompoundname")
	data.Id, _ = strconv.Atoi(values.Get("editinventoryid"))
	data.MinQuantity, _ = strconv.Atoi(values.Get("editminqty"))
	data.MaxQuantity, _ = strconv.Atoi(values.Get("editmaxqty"))

	if data.MinQuantity > data.MaxQuantity {
		slog.Error("Max Quantity is less than Min Quantity")
		http.Redirect(w, r, "/inventory-management?status=500&msg=Max Quantity is less than Min Quantity", http.StatusFound)
		return
	}

	jsonValue, err := json.Marshal(data)

	if err != nil {
		slog.Error("%s - error - %s", "Error marshaling user data", err)
		utils.SetResponse(w, http.StatusInternalServerError, "Failed to marshal data")
		return
	}
	resp, err := http.Post(utils.RestURL+"/update-coldstorage-quantity", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		slog.Error("%s - error - %s", "Error making POST request", err)
		utils.SetResponse(w, http.StatusInternalServerError, "Failed to create Compound Entry")
		return
	}

	if resp.StatusCode == http.StatusOK {
		http.Redirect(w, r, "/inventory-management?status=200&msg="+compoundName+" Quantity updated successfully", http.StatusFound)

	} else {
		http.Redirect(w, r, "/inventory-management?status=500&msg=Faild to update quantity of "+compoundName, http.StatusFound)
	}
}
