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
	materialmanagement "irpl.com/kanban-web/pages/materialmanagement"
	s "irpl.com/kanban-web/services"
)

func MaterialManagmentPage(w http.ResponseWriter, r *http.Request) {
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
			},
			{
				Name:     "Raw Material",
				Icon:     "fas fa-boxes",
				Link:     "/material-management",
				UserType: basepage.UserType{Admin: true},
				Style:    "font-size:1rem;",
				Selected: true,
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

	materialManagementPage := materialmanagement.MaterialManagement{
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
		Content:    materialManagementPage.Build(),
	}).Build()

	// Write the complete HTML page to the response
	w.Write([]byte(page))
}

func CreateNewOrUpdateExistingRawMaterial(w http.ResponseWriter, r *http.Request) {
	var data m.RawMaterial
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
	resp, err := http.Post(utils.RestURL+"/create-new-or-update-existing-material", "application/json", bytes.NewBuffer(jsonValue))
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

func materialmanagmentDialog(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	materialManagementPage := materialmanagement.MaterialManagement{
		MaterialId: id,
	}
	card := materialManagementPage.CardBuild()
	if err := json.NewEncoder(w).Encode(card); err != nil {
		log.Println("Error encoding user details:", err)
		http.Error(w, "Error to retrive user details", http.StatusInternalServerError)
		return
	}
}

func MaterialSearchPagination(w http.ResponseWriter, r *http.Request) {
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
		Data       []*m.RawMaterial
	}

	jsonValue, err := json.Marshal(req)
	if err != nil {
		slog.Error("Error marshaling table condition", "error", err)
	}

	// Send POST request
	resp, err := http.Post(RestURL+"/get-all-material-by-search-pagination", "application/json", bytes.NewBuffer(jsonValue))
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
	var MaterialTable s.TableCard
	tableTools := `<button type="button" class="btn  m-0 p-0" id="edit-Material-btn" data-bs-toggle="modal"  >
	<i class="fa fa-edit " style="color: #871a83;"></i>
   </button>`
	MaterialTable.CardHeading = "Raw Material"
	MaterialTable.CardHeadingActions.Component = []s.CardHeadActionComponent{
		{
			ComponentName: "modelbutton",
			ComponentType: s.ActionComponentElement{
				ModelButton: s.ModelButtonAttributes{Type: "button", Text: "Add New Material", ModelID: "#AddMaterialModel"}},
			Width: "col-4",
		},
		{
			ComponentName: "button",
			ComponentType: s.ActionComponentElement{
				Button: s.ButtonAttributes{ID: "importMaterial", Name: "importMaterials", Type: "button", Text: `Import<input style="display:none" class="clickable btn btn-success" type="file" id="om-import-file" name="file" accept=".xlsx">`, Class: "bg-transparent", CSS: "color: #871a83; border: 1px solid #871a83;"},
			},
			Width: "col-3",
		},
		{
			ComponentName: "button",
			ComponentType: s.ActionComponentElement{
				Button: s.ButtonAttributes{ID: "downloadTmp", Name: "downloadTmp", Type: "button", Text: ` Template <i class="fa fa-download" aria-hidden="true"></i> `, Class: "bg-transparent", AdditionalAttr: `data-link="/files/sample.pdf"`, CSS: "color: #871a83; border: 1px solid #871a83;", IsDownloadBtn: true, DownloadLink: `/static/templates/import-material-master.xlsx`},
			},
			Width: "col-3",
		},
	}

	MaterialTable.BodyTables = s.CardTableBody{Columns: []s.CardTableBodyHeadCol{
		{
			Lable:        "Description",
			Name:         `Description`,
			Type:         "input",
			DataField:    "Material_desc",
			IsSearchable: true,
			Width:        "col-2",
		},
		{
			Lable:        "SCADA Code",
			Name:         "SCADACode",
			Type:         "input",
			DataField:    "scada_code",
			IsSearchable: true,
			Width:        "col-2",
		},
		{
			Name:      "SAP Code",
			Type:      "input",
			DataField: "sap_code",
			Width:     "col-2",
		},
		{

			Lable:     "Comment",
			Name:      "Comment",
			Type:      "input",
			DataField: "comment",
			Width:     "col-2",
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
		ColumnsWidth: []string{"col-2", "col-2", "col-2", "col-2", "col-1", "col-1"},
		Data:         Response.Data,
		Tools:        tableTools,
		ID:           "Inventory-Table",
	}

	tableBodyHTML := MaterialTable.BodyTables.RenderBodyColumns()

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
func ImportMaterialMasterData(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20) // 10MB limit

	userIDStr := r.Header.Get("X-Custom-Userid")

	file, _, err := r.FormFile("file")
	if err != nil {
		utils.CreateSystemLogInternal(m.SystemLog{
			Message:     "ImportMaterialMasterData: Invalid file uploaded! Failed to read file",
			MessageType: "ERROR",
			IsCritical:  true,
			CreatedBy:   "system_logger",
		})
		http.Error(w, "Unable to retrieve the uploaded file: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	records, err := ParseMaterialCSV(file)
	if err != nil {
		utils.CreateSystemLogInternal(m.SystemLog{
			Message:     "ImportMaterialMasterData: Invalid file uploaded! Failed to parse file",
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
		if strings.TrimSpace(rec.SCADACode) == "" && strings.TrimSpace(rec.Description) == "" {
			utils.CreateSystemLogInternal(m.SystemLog{
				Message:     "ImportMaterialMasterData: Skipped record! Empty material code and description",
				MessageType: "INFO",
				IsCritical:  false,
				CreatedBy:   "system_logger",
			})
			skippedCount++
			continue
		}

		if strings.TrimSpace(rec.SCADACode) == "" && strings.TrimSpace(rec.SAPCode) == "" {
			utils.CreateSystemLogInternal(m.SystemLog{
				Message:     "ImportMaterialMasterData: Skipped record! Missing SCADA and SAP codes for material code: " + rec.SCADACode,
				MessageType: "INFO",
				IsCritical:  false,
				CreatedBy:   "system_logger",
			})
			skippedCount++
			continue
		}

		var materials []*m.RawMaterial
		rawQuery := m.RawQuery{
			Host:  utils.RestHost,
			Port:  utils.RestPort,
			Type:  "RawMaterial",
			Query: `SELECT * FROM raw_material WHERE scada_code = '` + rec.SCADACode + `';`,
		}
		rawQuery.RawQry(&materials)

		var extMaterialMaster m.RawMaterial
		if len(materials) != 0 {
			extMaterialMaster = *materials[0]
			extMaterialMaster.ModifiedBy = userIDStr
			extMaterialMaster.ModifiedOn = time.Now()
			updatedCount++
		} else {
			extMaterialMaster.CreatedBy = userIDStr
			extMaterialMaster.CreatedOn = time.Now()
			createdCount++
		}

		extMaterialMaster.Description = rec.Description
		extMaterialMaster.SCADACode = rec.SCADACode
		extMaterialMaster.SAPCode = rec.SAPCode
		extMaterialMaster.Comment = rec.Comment
		extMaterialMaster.Status = true

		jsonData, err := json.Marshal(extMaterialMaster)
		if err != nil {
			slog.Error("Error marshalling material data", "error", err)
			continue
		}

		resp, err := http.Post(utils.RestURL+"/create-new-or-update-existing-material", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			utils.CreateSystemLogInternal(m.SystemLog{
				Message:     "ImportMaterialMasterData: Failed to save material: " + extMaterialMaster.Description,
				MessageType: "ERROR",
				IsCritical:  true,
				CreatedBy:   "system_logger",
			})
			slog.Error("Error POSTing material", "error", err)
			continue
		}
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}

	// Final success log
	utils.CreateSystemLogInternal(m.SystemLog{
		Message:     fmt.Sprintf("ImportMaterialMasterData: Successfully imported material details. Created: %d, Updated: %d, Skipped: %d", createdCount, updatedCount, skippedCount),
		MessageType: "SUCCESS",
		IsCritical:  false,
		CreatedBy:   "system_logger",
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":      "Material CSV imported successfully",
		"createdCount": createdCount,
		"updatedCount": updatedCount,
		"skippedCount": skippedCount,
	})
}

func ParseMaterialCSV(file io.Reader) ([]m.CSVMaterial, error) {
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

	var records []m.CSVMaterial

	rows = utils.NormalizeRowsToLength(rows, 3)
	for _, row := range rows {
		firstCol := strings.ToLower(strings.TrimSpace(row[0]))

		// If this is the header, set expected length
		if firstCol == "material name" {
			continue
		}
		// Pad rows shorter than expected
		record := m.CSVMaterial{
			Description: strings.TrimSpace(row[0]),
			SCADACode:   strings.TrimSpace(row[1]),
			SAPCode:     strings.TrimSpace(row[2]),
			Comment:     strings.TrimSpace(row[3]),
			//LineName:     strings.TrimSpace(row[2]),
		}

		records = append(records, record)
	}

	return records, nil
}
