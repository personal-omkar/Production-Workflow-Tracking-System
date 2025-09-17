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
	chemicalmanagement "irpl.com/kanban-web/pages/chemicalmanagement"
	s "irpl.com/kanban-web/services"
)

func ChemicalManagmentPage(w http.ResponseWriter, r *http.Request) {
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
				Selected: true,
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

	chemicalManagementPage := chemicalmanagement.ChemicalManagement{
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
		Content:    chemicalManagementPage.Build(),
	}).Build()

	// Write the complete HTML page to the response
	w.Write([]byte(page))
}

func CreateNewOrUpdateExistingChemical(w http.ResponseWriter, r *http.Request) {
	var data m.ChemicalTypes
	var apiresp m.ApiRespMsg
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Println("Invalid request body: error- ", err.Error())
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	jsonValue, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error: Error marshaling user data: %v", err)
		http.Error(w, "Failed to marshal data", http.StatusInternalServerError)
		return
	}
	resp, err := http.Post(utils.RestURL+"/create-new-or-update-existing-chemical", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Printf("Error making POST request: %v", err)
		http.Error(w, "Error making POST request:", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		http.Error(w, "Failed to read response body", http.StatusInternalServerError)
		return
	}
	apiresp.Code = resp.StatusCode
	apiresp.Message = string(responseBody)
	if err := json.NewEncoder(w).Encode(apiresp); err != nil {
		http.Error(w, "Failed to encode success message", http.StatusInternalServerError)
	}
}

func chemicalmanagmentDialog(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	chemicalManagementPage := chemicalmanagement.ChemicalManagement{
		ChemicalId: id,
	}
	card := chemicalManagementPage.CardBuild()
	if err := json.NewEncoder(w).Encode(card); err != nil {
		log.Println("Error encoding user details:", err)
		http.Error(w, "Error to retrive user details", http.StatusInternalServerError)
		return
	}
}

func ChemicalSearchPagination(w http.ResponseWriter, r *http.Request) {
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
		Data       []*m.ChemicalTypes
	}

	jsonValue, err := json.Marshal(req)
	if err != nil {
		slog.Error("Error marshaling table condition", "error", err)
	}

	// Send POST request
	resp, err := http.Post(RestURL+"/get-all-chemical-by-search-pagination", "application/json", bytes.NewBuffer(jsonValue))
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
	var ChemicalTable s.TableCard
	tableTools := `<button type="button" class="btn  m-0 p-0" id="edit-Chemical-btn" data-bs-toggle="modal"  >
	<i class="fa fa-edit " style="color: #871a83;"></i>
   </button>`
	ChemicalTable.CardHeading = "Chemical Type Master"
	ChemicalTable.CardHeadingActions.Component = []s.CardHeadActionComponent{
		{
			ComponentName: "modelbutton",
			ComponentType: s.ActionComponentElement{
				ModelButton: s.ModelButtonAttributes{Type: "button", Text: "Add Chemical Type", ModelID: "#AddChemicalModel"}},
			Width: "col-4",
		},
		{
			ComponentName: "button",
			ComponentType: s.ActionComponentElement{
				Button: s.ButtonAttributes{ID: "importChemical", Name: "importChemicalTypes", Type: "button", Text: `Import<input style="display:none" class="clickable btn btn-success" type="file" id="om-import-file" name="file" accept=".xlsx">`, Class: "bg-transparent", CSS: "color: #871a83; border: 1px solid #871a83;"},
			},
			Width: "col-3",
		},
		{
			ComponentName: "button",
			ComponentType: s.ActionComponentElement{
				Button: s.ButtonAttributes{ID: "downloadTmp", Name: "downloadTmp", Type: "button", Text: ` Template <i class="fa fa-download" aria-hidden="true"></i> `, Class: "bg-transparent", AdditionalAttr: `data-link="/files/sample.pdf"`, CSS: "color: #871a83; border: 1px solid #871a83;", IsDownloadBtn: true, DownloadLink: `/static/templates/import-chemical-master.xlsx`},
			},
			Width: "col-3",
		},
	}

	ChemicalTable.BodyTables = s.CardTableBody{Columns: []s.CardTableBodyHeadCol{
		{
			Lable:        "Chemical Type",
			Name:         `Type`,
			Type:         "input",
			DataField:    "type",
			IsSearchable: true,
			Width:        "col-2",
		},
		{
			Lable:        "Conv Code",
			Name:         "ConvCode",
			Type:         "input",
			DataField:    "conv_code",
			IsSearchable: true,
			Width:        "col-2",
		},
		{
			Name:  "Status",
			Type:  "input",
			ID:    "search-max-quantity",
			Width: "col-1",
		},
		{
			Name:  "Tools",
			Width: "col-1",
		},
	},
		ColumnsWidth: []string{"col-2", "col-2", "col-2", "col-2"},
		Data:         Response.Data,
		Tools:        tableTools,
		ID:           "Inventory-Table",
	}

	tableBodyHTML := ChemicalTable.BodyTables.RenderBodyColumns()

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
func ImportChemicalMasterData(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20) // 10MB limit

	userIDStr := r.Header.Get("X-Custom-Userid")

	file, _, err := r.FormFile("file")
	if err != nil {
		utils.CreateSystemLogInternal(m.SystemLog{
			Message:     "ImportChemicalMasterData: Invalid file uploaded! Failed to read file",
			MessageType: "ERROR",
			IsCritical:  true,
			CreatedBy:   "system_logger",
		})
		http.Error(w, "Unable to retrieve the uploaded file: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	records, err := ParseChemicalCSV(file)
	if err != nil {
		utils.CreateSystemLogInternal(m.SystemLog{
			Message:     "ImportChemicalMasterData: Invalid file uploaded! Failed to parse file",
			MessageType: "ERROR",
			IsCritical:  true,
			CreatedBy:   "system_logger",
		})
		slog.Error("Error parsing CSV", "error", err)
		http.Error(w, "Invalid CSV format", http.StatusBadRequest)
		return
	}

	createdCount := 0
	updatedCount := 0
	skippedCount := 0

	for _, rec := range records {
		if strings.TrimSpace(rec.Type) == "" && strings.TrimSpace(rec.ConvCode) == "" {
			utils.CreateSystemLogInternal(m.SystemLog{
				Message:     "ImportChemicalMasterData: Skipped record! Empty chemical code and description",
				MessageType: "INFO",
				IsCritical:  false,
				CreatedBy:   "system_logger",
			})
			skippedCount++
			continue
		}

		var chemicals []*m.ChemicalTypes
		rawQuery := m.RawQuery{
			Host:  utils.RestHost,
			Port:  utils.RestPort,
			Type:  "ChemicalTypes",
			Query: `SELECT * FROM chemical_types WHERE type = '` + rec.Type + `';`,
		}
		rawQuery.RawQry(&chemicals)

		var extChemicalMaster m.ChemicalTypes
		if len(chemicals) != 0 {
			extChemicalMaster = *chemicals[0]
			extChemicalMaster.ModifiedBy = userIDStr
			extChemicalMaster.ModifiedOn = time.Now()
			updatedCount++
		} else {
			extChemicalMaster.CreatedBy = userIDStr
			extChemicalMaster.CreatedOn = time.Now()
			createdCount++
		}

		extChemicalMaster.Type = rec.Type
		extChemicalMaster.ConvCode = rec.ConvCode
		extChemicalMaster.Status = true

		jsonData, err := json.Marshal(extChemicalMaster)
		if err != nil {
			slog.Error("Error marshalling chemical data", "error", err)
			continue
		}

		resp, err := http.Post(utils.RestURL+"/create-new-or-update-existing-chemical", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			utils.CreateSystemLogInternal(m.SystemLog{
				Message:     "ImportChemicalMasterData: Failed to save chemical: " + extChemicalMaster.Type,
				MessageType: "ERROR",
				IsCritical:  true,
				CreatedBy:   "system_logger",
			})
			slog.Error("Error POSTing chemical", "error", err)
			continue
		}
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}

	// Final success log
	utils.CreateSystemLogInternal(m.SystemLog{
		Message:     fmt.Sprintf("ImportChemicalMasterData: Successfully imported chemical details. Created: %d, Updated: %d, Skipped: %d", createdCount, updatedCount, skippedCount),
		MessageType: "SUCCESS",
		IsCritical:  false,
		CreatedBy:   "system_logger",
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":      "Chemical CSV imported successfully",
		"createdCount": createdCount,
		"updatedCount": updatedCount,
		"skippedCount": skippedCount,
	})
}

func ParseChemicalCSV(file io.Reader) ([]m.CSVChemical, error) {
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

	var records []m.CSVChemical

	rows = utils.NormalizeRowsToLength(rows, 3)
	for _, row := range rows {
		firstCol := strings.ToLower(strings.TrimSpace(row[0]))

		// If this is the header, set expected length
		if firstCol == "Chemical Type" {
			continue
		}
		// Pad rows shorter than expected
		record := m.CSVChemical{
			Type:     strings.TrimSpace(row[0]),
			ConvCode: strings.TrimSpace(row[1]),
			Status:   true,
		}

		records = append(records, record)
	}

	return records, nil
}
