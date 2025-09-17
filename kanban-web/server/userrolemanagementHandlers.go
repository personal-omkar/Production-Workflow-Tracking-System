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
	"irpl.com/kanban-web/pages/basepage"
	"irpl.com/kanban-web/pages/userrolemanagementpage"
	s "irpl.com/kanban-web/services"
)

func userRoleManagmentPage(w http.ResponseWriter, r *http.Request) {
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

	userRoleManagementPage := userrolemanagementpage.UserRoleManagement{
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
		Content:    userRoleManagementPage.Build(),
	}).Build()

	// Write the complete HTML page to the response
	w.Write([]byte(page))
}

// Create UserRole
func CreateUserRole(w http.ResponseWriter, r *http.Request) {
	var data m.UserRoles
	userID := r.Header.Get("X-Custom-Userid")
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err == nil {
		data.CreatedBy = userID
		jsonValue, err := json.Marshal(data)
		if err != nil {
			slog.Error("Error marshaling user role data", "error", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to marshal user role data")
			return
		}

		resp, err := http.Post(utils.RestURL+"/create-user-role", "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			slog.Error("Error making POST request", "error", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to create user role")
			return
		}
		defer resp.Body.Close()

		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			slog.Error("Error reading response body", "error", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to read response body")
			return
		}

		utils.SetResponse(w, resp.StatusCode, string(responseBody))
	} else {
		slog.Error("Record creation failed", "error", err.Error())
		utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to create record")
	}
}

// Update User Role
func UpdateUserRole(w http.ResponseWriter, r *http.Request) {
	var data m.UserRoles
	userID := r.Header.Get("X-Custom-Userid")

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err == nil {
		data.ModifiedBy = userID
		jsonValue, err := json.Marshal(data)
		if err != nil {
			slog.Error("Error marshaling user role data", "error", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to marshal user role data")
			return
		}

		resp, err := http.Post(utils.RestURL+"/update-user-role", "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			slog.Error("Error making POST request", "error", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to update user role")
			return
		}
		defer resp.Body.Close()

		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			slog.Error("Error reading response body", "error", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to read response body")
			return
		}

		utils.SetResponse(w, http.StatusOK, string(responseBody))
	} else {
		slog.Error("Record update failed", "error", err.Error())
		utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to update record")
	}
}

// Delete UserRole
func DeleteUserRole(w http.ResponseWriter, r *http.Request) {
	var ids []int64
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&ids); err == nil {
		jsonValue, err := json.Marshal(ids)
		if err != nil {
			slog.Error("Error marshaling delete request data", "error", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to marshal delete request data")
			return
		}

		resp, err := http.Post(utils.RestURL+"/delete-user-role", "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			slog.Error("Error making POST request for delete", "error", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to delete user role(s)")
			return
		}
		defer resp.Body.Close()

		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			slog.Error("Error reading response body for delete", "error", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to read response body for delete")
			return
		}

		utils.SetResponse(w, http.StatusOK, string(responseBody))
	} else {
		slog.Error("Record deletion failed", "error", err.Error())
		utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to delete record")
	}
}
func UserRoleSearchPagination(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Pagination m.PaginationReq `json:"pagination"`
		Conditions []string        `json:"conditions"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var Response struct {
		Pagination m.PaginationResp
		Data       []*m.UserRoles
	}

	jsonValue, err := json.Marshal(req)
	if err != nil {
		slog.Error("Error marshaling request", "error", err)
		return
	}

	resp, err := http.Post(RestURL+"/get-all-user-role-by-search-pagination", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		slog.Error("Error making POST request", "error", err)
		http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&Response); err != nil {
		slog.Error("Error decoding response", "error", err)
		http.Error(w, "Failed to decode data", http.StatusInternalServerError)
		return
	}

	var UserRoleTable s.TableCard
	UserRoleTable.CardHeading = "User Role Master"
	UserRoleTable.BodyTables = s.CardTableBody{
		Columns: []s.CardTableBodyHeadCol{
			{
				Name:         "Role Name",
				Type:         "input",
				DataField:    "role_name",
				IsSearchable: true,
				Width:        "col-2",
			},
			{
				Name:         "Description",
				Type:         "input",
				DataField:    "description",
				IsSearchable: true,
				Width:        "col-2",
			},
			{
				Name:  "Created On",
				Width: "col-2",
			},
			{
				Name:  "Tools",
				Width: "col-1",
			},
		},
		ID:    "advanced-search-table-body",
		Data:  Response.Data,
		Tools: `<button type="button" class="btn m-0 p-0" id="edit-User-Role-btn" data-bs-toggle="modal" data-bs-target="#EditUserRoleModel"><i class="fa fa-edit" style="color: #871a83;"></i></button>`,
	}

	tableBodyHTML := UserRoleTable.BodyTables.RenderBodyColumns()

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
