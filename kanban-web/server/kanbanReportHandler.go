package server

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"irpl.com/kanban-commons/model"
	"irpl.com/kanban-commons/utils"
	"irpl.com/kanban-web/pages/basepage"
	"irpl.com/kanban-web/pages/kanbanreport"
	s "irpl.com/kanban-web/services"
)

func kanabnReportPage(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("X-Custom-Username")
	userID := r.Header.Get("X-Custom-Userid")
	usertype := r.Header.Get("X-Custom-Role")
	links := r.Header.Get("X-Custom-Allowlist")
	var vendorName string
	var vendorRecord []model.Vendors
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
				Selected: true,
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

	historyPagePage := kanbanreport.KanbanReportPage{
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
		Content:    historyPagePage.Build(),
	})

	page.AddStyleCode(`
	/*css*/
		[class^="col-"] {
			flex-grow: 0;
			flex-shrink: 0;
		}

		.col-0_5 { flex-basis: 4.1667%; max-width: 4.1667%; } 
		.col-1_5 { flex-basis: 12.5%; max-width: 12.5%; }   
		.col-2_5 { flex-basis: 20.8333%; max-width: 20.8333%; } 
		.col-3_5 { flex-basis: 29.1667%; max-width: 29.1667%; }
		.col-4_5 { flex-basis: 37.5%; max-width: 37.5%; }
		.col-5_5 { flex-basis: 45.8333%; max-width: 45.8333%; }
		.col-6_5 { flex-basis: 54.1667%; max-width: 54.1667%; }
		.col-7_5 { flex-basis: 62.5%; max-width: 62.5%; }
		.col-8_5 { flex-basis: 70.8333%; max-width: 70.8333%; }
		.col-9_5 { flex-basis: 79.1667%; max-width: 79.1667%; }
		.col-10_5 { flex-basis: 87.5%; max-width: 87.5%; }
		.col-11_5 { flex-basis: 95.8333%; max-width: 95.8333%; }

	/*!css*/
	`)

	// Write the complete HTML page to the response
	w.Write([]byte(page.Build()))
}

func kanbanReportSearchPagination(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Pagination model.PaginationReq `json:"pagination"`
		Conditions []string            `json:"Conditions"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var tabledata []*model.CustomerOrderDetails
	var Response struct {
		Pagination model.PaginationResp
		Data       []*model.OrderDetails
	}

	jsonValue, err := json.Marshal(req)
	if err != nil {
		slog.Error("Error marshaling table condition", "error", err)
	}

	// Send POST request
	resp, err := http.Post(RestURL+"/get-all-kanban-details-for-report", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		slog.Error("Error making POST request", "error", err)
	}

	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&Response); err != nil {
		slog.Error("error decoding response body", "error", err)
	}
	for _, value := range Response.Data {
		value.Status = strings.Title(value.Status)
	}

	for _, val := range Response.Data {
		var temp model.CustomerOrderDetails
		temp.Id = val.Id
		temp.CompoundId = val.CompoundId
		temp.MFGDateTime = val.MFGDateTime.String()
		temp.DemandDateTime = val.DemandDateTime.String()
		temp.ExpDate = val.ExpDate
		temp.CellNo = val.CellNo
		temp.NoOFLots = val.NoOFLots
		temp.Location = val.Location
		temp.KbRootId = val.KbRootId
		temp.CreatedBy = val.CreatedBy
		temp.CreatedOn = val.CreatedOn
		temp.ModifiedBy = val.ModifiedBy
		temp.ModifiedOn = val.ModifiedOn
		temp.Status = map[string]string{"0": "Kanban", "1": "In Process", "2": "Quality Test", "3": "Packing", "4": "Dispatched", "-1": "Quality Test Fail"}[val.Status]
		temp.VendorName = val.VendorName
		temp.CompoundName = val.CompoundName
		temp.OrderId = val.OrderId
		temp.MinQuantity = val.MinQuantity
		temp.AvailableQuantity = val.AvailableQuantity
		temp.LotNo = val.LotNo
		tabledata = append(tabledata, &temp)
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

	vendorOrderTable.CardHeading = "Kanban Report"
	vendorOrderTable.CardHeadingActions = s.CardHeadActionBody{
		Style: "direction: ltr;",
	}
	vendorOrderTable.CardHeadingActions.Component = []s.CardHeadActionComponent{
		{ComponentName: "input", ComponentType: s.ActionComponentElement{Input: s.InputAttributes{ID: "FromDate", Name: "FromDate", Type: "date", Icon: "From Date"}}, Width: "col-5"},
		{ComponentName: "input", ComponentType: s.ActionComponentElement{Input: s.InputAttributes{ID: "ToDate", Name: "ToDate", Type: "date", Icon: "To Date"}}, Width: "col-5"},
		{ComponentName: "button", ComponentType: s.ActionComponentElement{Button: s.ButtonAttributes{ID: "searchBtn", Disabled: true, Colour: "#007BFF", Name: "searchBtn", Type: "button", Text: "Search"}}, Width: "col-2"},
	}
	vendorOrderTable.BodyTables = s.CardTableBody{Columns: []s.CardTableBodyHeadCol{
		{
			Name:       `Sr. No.`,
			IsCheckbox: true,
			ID:         "search-sn",
		},
		{
			Name:         "Vendor Name",
			ID:           "Customer Name",
			Width:        "col-2",
			IsSearchable: true,
			DataField:    "vendor_name",
			Type:         "input",
			Value:        searchFilters["vendor_name"],
		},
		{
			Lable:        "Part Name",
			Name:         "Compound Name",
			ID:           "compound_code",
			Width:        "col-2",
			IsSearchable: true,
			DataField:    "compound_name",
			Type:         "input",
			Value:        searchFilters["compound_name"],
		},
		{
			Lable:        "Kanban Summary",
			Name:         "Cell No",
			ID:           "cell_no",
			Width:        "col-2",
			IsSearchable: true,
			DataField:    "cell_no",
			Type:         "input",
			Value:        searchFilters["cell_no"],
		},
		{
			Lable:        "Lot Number",
			Name:         "LotNo",
			ID:           "lot_no",
			Width:        "col-1",
			IsSearchable: true,
			DataField:    "lot_no",
			Type:         "input",
		},
		{
			Name:         "Demand Date Time",
			ID:           "demand_date",
			Width:        "col-2",
			IsSearchable: true,
			DataField:    "demand_date_time",
			Type:         "input",
			Value:        searchFilters["demand_date_time"],
		},
		{
			Name:         "Status",
			Width:        "col-2",
			IsSearchable: true,
			DataField:    "status",
			Type:         "input",
			Value:        searchFilters["status"],
		},
	},
		ColumnsWidth: []string{"col-1", "col-1", "col-1", "col-1", "col-1", "col-1", "col-1"},
		Data:         tabledata,
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
