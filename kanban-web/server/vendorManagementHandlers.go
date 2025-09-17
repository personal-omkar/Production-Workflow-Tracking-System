package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
	m "irpl.com/kanban-commons/model"
	"irpl.com/kanban-commons/utils"
	"irpl.com/kanban-web/pages/basepage"
	"irpl.com/kanban-web/pages/vendormanagement"
	s "irpl.com/kanban-web/services"
)

func vendorManagmentPage(w http.ResponseWriter, r *http.Request) {
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
			},
			{
				Name:     "Vendor Master",
				Icon:     "fa fa-briefcase",
				Link:     "/vendor-management",
				UserType: basepage.UserType{Admin: true},
				Style:    "font-size:1rem;",
				Selected: true,
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

	vendorManagementPage := vendormanagement.VendorManagement{
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
		Content:    vendorManagementPage.Build(),
	}).Build()

	// Write the complete HTML page to the response
	w.Write([]byte(page))
}

// GetVendorByUserID returns a vendor  records based on user id
func GetVendorByUserID(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")

	resp, err := http.Get(RestURL + "/get-vendor-by-userid?key=" + key + "&value=" + value)
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
	var ApiResp m.ApiRespMsg
	ApiResp.Code = resp.StatusCode
	ApiResp.Message = string(responseBody)
	body, err := json.Marshal(ApiResp)
	if err != nil {
		log.Printf("Error marshaling response: %v", err)
	}
	utils.SetResponse(w, resp.StatusCode, string(body))
}

func GetVendorDetailsByVendorCode(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	url := RestURL + "/get-vendor-details-by-vendor-code"
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

func ImportVendorMasterData(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20) // 10MB limit
	userIDStr := r.Header.Get("X-Custom-Userid")
	// Retrieve the uploaded file
	file, _, err := r.FormFile("file")
	if err != nil {
		sysLog := m.SystemLog{
			Message:     "ImportVendorMasterData: Invalid file uploaded! Failed to read file",
			MessageType: "ERROR",
			IsCritical:  true,
			CreatedBy:   "system_logger",
		}
		utils.CreateSystemLogInternal(sysLog)
		http.Error(w, "Unable to retrieve the uploaded file: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	records, err := ParseVendorCSV(file)
	if err != nil {
		sysLog := m.SystemLog{
			Message:     "ImportVendorMasterData: Invalid file uploaded! Failed to parse file",
			MessageType: "ERROR",
			IsCritical:  true,
			CreatedBy:   "system_logger",
		}
		utils.CreateSystemLogInternal(sysLog)
		log.Fatalf("Error parsing CSV: %v", err)
	}

	for _, rec := range records {
		if strings.TrimSpace(rec.VendorCode) == "" && strings.TrimSpace(rec.VendorName) == "" {
			sysLog := m.SystemLog{
				Message:     "ImportVendorMasterData: Skipped record! Empty vendor code and vendor name",
				MessageType: "INFO",
				IsCritical:  true,
				CreatedBy:   "system_logger",
			}
			utils.CreateSystemLogInternal(sysLog)
			continue
		}
		if strings.TrimSpace(rec.VendorCode) == "" {
			sysLog := m.SystemLog{
				Message:     "ImportVendorMasterData: Skipped record! Empty vendor code for vendor " + rec.VendorName,
				MessageType: "INFO",
				IsCritical:  true,
				CreatedBy:   "system_logger",
			}
			utils.CreateSystemLogInternal(sysLog)
			continue
		}
		if strings.TrimSpace(rec.VendorName) == "" {
			sysLog := m.SystemLog{
				Message:     "ImportVendorMasterData: Skipped record! Empty vendor name for vendor " + rec.VendorCode,
				MessageType: "INFO",
				IsCritical:  true,
				CreatedBy:   "system_logger",
			}
			utils.CreateSystemLogInternal(sysLog)
			continue
		}
		var Vendors []*m.Vendors
		var rawQuery m.RawQuery
		rawQuery.Host = utils.RestHost
		rawQuery.Port = utils.RestPort
		rawQuery.Type = "Vendors"
		rawQuery.Query = `SELECT * FROM vendors WHERE vendor_code = '` + rec.VendorCode + `';`
		rawQuery.RawQry(&Vendors)

		if len(Vendors) != 0 {
			var extVendorMaster = Vendors[0]
			extVendorMaster.VendorCode = rec.VendorCode
			extVendorMaster.VendorName = rec.VendorName
			extVendorMaster.ModifiedBy = userIDStr
			extVendorMaster.Address = rec.Address
			extVendorMaster.ContactInfo = rec.ContactInfo
			extVendorMaster.PerDayLotConfig, _ = strconv.Atoi(rec.PerDayLotConfig)
			extVendorMaster.PerMonthLotConfig, _ = strconv.Atoi(rec.PerMonthLotConfig)
			extVendorMaster.PerHourLotConfig, _ = strconv.Atoi(rec.PerHourLotConfig)
			extVendorMaster.ModifiedOn = time.Now()
			extVendorMaster.Isactive = true

			if extVendorMaster.PerHourLotConfig > extVendorMaster.PerDayLotConfig && extVendorMaster.PerDayLotConfig > extVendorMaster.PerMonthLotConfig {
				sysLog := m.SystemLog{
					Message:     "ImportVendorMasterData: Skipped record! Monthly Limit should be greater than Daily, and Daily greater than Hourly for vendor: " + extVendorMaster.VendorName,
					MessageType: "INFO",
					IsCritical:  true,
					CreatedBy:   "system_logger",
				}
				utils.CreateSystemLogInternal(sysLog)
				continue
			}
			jsonData, err := json.Marshal(extVendorMaster)
			if err != nil {
				slog.Error("%s - error - %s", "Error in marshalling data:", err)
			}
			resp, err := http.Post(utils.RestURL+"/create-new-vendor", "application/json", bytes.NewBuffer(jsonData))
			if err != nil {
				sysLog := m.SystemLog{
					Message:     "ImportVendorMasterData: Skipped record! failed to update vendor: " + extVendorMaster.VendorName,
					MessageType: "INFO",
					IsCritical:  true,
					CreatedBy:   "system_logger",
				}
				utils.CreateSystemLogInternal(sysLog)
				slog.Error("%s - error - %s", "Error making POST request", err)
			}
			defer resp.Body.Close()
		} else {
			var extVendorMaster m.Vendors
			extVendorMaster.VendorCode = rec.VendorCode
			extVendorMaster.VendorName = rec.VendorName
			extVendorMaster.CreatedBy = userIDStr
			extVendorMaster.Address = rec.Address
			extVendorMaster.ContactInfo = rec.ContactInfo
			extVendorMaster.PerDayLotConfig, _ = strconv.Atoi(rec.PerDayLotConfig)
			extVendorMaster.PerMonthLotConfig, _ = strconv.Atoi(rec.PerMonthLotConfig)
			extVendorMaster.PerHourLotConfig, _ = strconv.Atoi(rec.PerHourLotConfig)
			extVendorMaster.CreatedOn = time.Now()
			extVendorMaster.Isactive = true

			if extVendorMaster.PerHourLotConfig > extVendorMaster.PerDayLotConfig && extVendorMaster.PerDayLotConfig > extVendorMaster.PerMonthLotConfig {
				sysLog := m.SystemLog{
					Message:     "ImportVendorMasterData: Skipped record! Monthly Limit should be greater than Daily, and Daily greater than Hourly for vendor: " + extVendorMaster.VendorName,
					MessageType: "INFO",
					IsCritical:  true,
					CreatedBy:   "system_logger",
				}
				utils.CreateSystemLogInternal(sysLog)
				continue
			}

			jsonData, err := json.Marshal(extVendorMaster)
			if err != nil {
				slog.Error("%s - error - %s", "Error in marshalling data:", err)
			}
			resp, err := http.Post(utils.RestURL+"/create-new-vendor", "application/json", bytes.NewBuffer(jsonData))
			if err != nil {
				sysLog := m.SystemLog{
					Message:     "ImportVendorMasterData: Skipped record! failed to create vendor: " + rec.VendorName,
					MessageType: "ERROR",
					IsCritical:  true,
					CreatedBy:   "system_logger",
				}
				utils.CreateSystemLogInternal(sysLog)
				slog.Error("%s - error - %s", "Error making POST request", err)
			}
			defer resp.Body.Close()
		}

	}

	sysLog := m.SystemLog{
		Message:     "ImportVendorMasterData: Successfully Imported vendor details",
		MessageType: "SUCCESS",
		IsCritical:  false,
		CreatedBy:   "system_logger",
	}
	utils.CreateSystemLogInternal(sysLog)
	fmt.Fprintln(w, "CSV Imported Successfully")
}

func ParseVendorCSV(file io.Reader) ([]m.CSVVendor, error) {
	f, err := excelize.OpenReader(file)
	if err != nil {
		log.Println("ParseCSV: Failed to open xlsx file, error -", err)
	}
	defer f.Close()

	sheet := f.GetSheetName(0)
	rows, err := f.GetRows(sheet)
	if err != nil {
		log.Println("ParseCSV: Failed to read row from xlsx file, error -", err)
	}

	var records []m.CSVVendor
	rows = utils.NormalizeRowsToLength(rows, 7)
	// Skip the first 2 rows
	for _, row := range rows {
		// if len(row) < 2 {
		// 	continue // skip incomplete rows
		// }

		// Skip row if it's a header
		if strings.ToLower(strings.TrimSpace(row[0])) == "vendor code" {
			continue
		}

		record := m.CSVVendor{
			VendorCode:        strings.TrimSpace(row[0]),
			VendorName:        strings.TrimSpace(row[1]),
			ContactInfo:       strings.TrimSpace(row[2]),
			Address:           strings.TrimSpace(row[3]),
			PerDayLotConfig:   strings.TrimSpace(row[4]),
			PerMonthLotConfig: strings.TrimSpace(row[5]),
			PerHourLotConfig:  strings.TrimSpace(row[6]),
		}

		records = append(records, record)
	}

	return records, nil
}

func VendorSearchPagination(w http.ResponseWriter, r *http.Request) {
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
		Data       []*m.Vendors
	}

	jsonValue, err := json.Marshal(req)
	if err != nil {
		slog.Error("Error marshaling table condition", "error", err)
	}

	// Send POST request
	resp, err := http.Post(RestURL+"/get-all-vendor-by-search-pagination", "application/json", bytes.NewBuffer(jsonValue))
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
	var VendorTable s.TableCard
	tableTools := `<button  type="button" class="btn  m-0 p-0" id="edit-vendor-btn" data-bs-toggle="modal" data-bs-target="#EditVendorModel"> 
 					<i class="fa fa-edit " style="color: #871a83;"></i> 
					</button>`
	VendorTable.CardHeading = "Vendor Master"
	VendorTable.CardHeadingActions.Component = []s.CardHeadActionComponent{
		{
			ComponentName: "modelbutton",
			ComponentType: s.ActionComponentElement{
				ModelButton: s.ModelButtonAttributes{Type: "button", Text: "Add New Operator", ModelID: "#AddOperatorModel"}},
			Width: "col-4",
		},
		{
			ComponentName: "button",
			ComponentType: s.ActionComponentElement{
				Button: s.ButtonAttributes{ID: "importOperators", Name: "importOperators", Type: "button", Text: `Import<input style="display:none" class="clickable btn btn-success" type="file" id="om-import-file" name="file" accept=".xlsx">`, Class: "bg-transparent", CSS: "color: #871a83; border: 1px solid #871a83;"},
			},
			Width: "col-3",
		},
		{
			ComponentName: "button",
			ComponentType: s.ActionComponentElement{
				Button: s.ButtonAttributes{ID: "downloadTmp", Name: "downloadTmp", Type: "button", Text: ` Template <i class="fa fa-download" aria-hidden="true"></i> `, Class: "bg-transparent", AdditionalAttr: `data-link="/files/sample.pdf"`, CSS: "color: #871a83; border: 1px solid #871a83;", IsDownloadBtn: true, DownloadLink: `/static/templates/import-operator-master.xlsx`},
			},
			Width: "col-3",
		},
	}

	VendorTable.BodyTables = s.CardTableBody{Columns: []s.CardTableBodyHeadCol{
		{
			Name:         "Vendor Code",
			ID:           "vendorCode",
			Width:        "col-2",
			IsSearchable: true,
			DataField:    "vendor_code",
			Type:         "input",
		},
		{
			Name:         "Vendor Name",
			ID:           "vendorName",
			Width:        "col-2",
			IsSearchable: true,
			DataField:    "vendor_name",
			Type:         "input",
		},
		{
			Name:  "Contact Info",
			ID:    "contactInfo",
			Width: "col-2",
		},
		{
			Name:  "Created On",
			ID:    "createdOn",
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
		ID:           "VendorManagement",
	}

	tableBodyHTML := VendorTable.BodyTables.RenderBodyColumns()

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
