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
	"irpl.com/kanban-web/pages/compounds"
	s "irpl.com/kanban-web/services"
)

func CompoundsManagement(w http.ResponseWriter, r *http.Request) {
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

	var compoundlist []m.Compounds
	compoundsresp, err := http.Get(RestURL + "/get-all-compound-data")
	if err != nil {
		slog.Error("%s - error - %s", "Error making GET request", err)
	}

	defer compoundsresp.Body.Close()

	if err := json.NewDecoder(compoundsresp.Body).Decode(&compoundlist); err != nil {
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
				Selected: true,
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

	CompoundManagementPage := compounds.CompoundsManagement{
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
		Content:    CompoundManagementPage.Build(),
	}).Build()

	// Write the complete HTML page to the response
	w.Write([]byte(page))
}

func AddorUpdateCompound(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-Custom-Userid")
	var CompoundData m.Compounds

	// Decode the request body
	err := json.NewDecoder(r.Body).Decode(&CompoundData)
	if err != nil {
		utils.SetResponse(w, http.StatusInternalServerError, "Fail: Failed to decode data")
		slog.Error("%s - error - %s", "Failed to decode compound data", err.Error())
		return
	}

	// Set CreatedBy or ModifiedBy based on ID presence
	if CompoundData.Id == 0 {
		CompoundData.CreatedBy = userID
	} else {
		CompoundData.ModifiedBy = userID
	}

	// Marshal the compound data
	jsonValue, err := json.Marshal(CompoundData)
	if err != nil {
		slog.Error("%s - error - %s", "Error marshaling compound data", err)
		utils.SetResponse(w, http.StatusInternalServerError, "Failed to marshal compound data")
		return
	}

	// Make the POST request
	resp, err := http.Post(utils.RestURL+"/add-update-compound", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		slog.Error("%s - error - %s", "Error making POST request", err)
		utils.SetResponse(w, http.StatusInternalServerError, "Failed to update compound")
		return
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode == http.StatusCreated {
		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			slog.Error("%s - error - %s", "Error reading response body", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to read response body")
			return
		}
		utils.SetResponse(w, http.StatusOK, string(responseBody))
	} else {
		slog.Error("%s - error - %s", "Record updation failed", err)
		utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to update record")
	}
}

// Get All Compounds
func GetAllActiveCompounds(w http.ResponseWriter, r *http.Request) {

	resp, err := http.Get(RestURL + "/get-all-active-compounds")
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

	utils.SetResponse(w, resp.StatusCode, string(responseBody))
}

func ImportCompoundData(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20) // 10MB limit
	userIDStr := r.Header.Get("X-Custom-Userid")
	// Retrieve the uploaded file
	file, _, err := r.FormFile("file")
	if err != nil {
		sysLog := m.SystemLog{
			Message:     "ImportCompoundMasterData: Invalid file uploaded! Failed to read file",
			MessageType: "ERROR",
			IsCritical:  true,
			CreatedBy:   "system_logger",
		}
		utils.CreateSystemLogInternal(sysLog)
		http.Error(w, "Unable to retrieve the uploaded file: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	records, err := ParseCSV(file)
	if err != nil {
		sysLog := m.SystemLog{
			Message:     "ImportCompoundMasterData: Invalid file uploaded! Failed to parse file",
			MessageType: "ERROR",
			IsCritical:  true,
			CreatedBy:   "system_logger",
		}
		utils.CreateSystemLogInternal(sysLog)
		log.Fatalf("Error parsing CSV: %v", err)
	}

	for _, rec := range records {

		if strings.TrimSpace(rec.CompoundName) == "" {
			sysLog := m.SystemLog{
				Message:     "ImportCompoundMasterData: Skipping Compound!, empty compound name",
				MessageType: "INFO",
				IsCritical:  true,
				CreatedBy:   "system_logger",
			}
			utils.CreateSystemLogInternal(sysLog)
			continue
		}
		var PMs []*m.Compounds
		var rawQuery m.RawQuery
		rawQuery.Host = utils.RestHost
		rawQuery.Port = utils.RestPort
		rawQuery.Type = "Compounds"
		rawQuery.Query = `SELECT * FROM compounds WHERE compound_name = '` + rec.CompoundName + `';`
		rawQuery.RawQry(&PMs)

		if len(PMs) != 0 {
			var extPartMaster = PMs[0]
			extPartMaster.CompoundName = rec.CompoundName
			extPartMaster.Description = rec.Description
			extPartMaster.SCADACode = rec.SCADACode
			extPartMaster.SAPCode = rec.SAPCode
			extPartMaster.ModifiedBy = userIDStr
			extPartMaster.ModifiedOn = time.Now()

			jsonData, err := json.Marshal(extPartMaster)
			if err != nil {
				slog.Error("%s - error - %s", "Error in marshalling data:", err)
			}
			resp, err := http.Post(utils.RestURL+"/add-update-compound", "application/json", bytes.NewBuffer(jsonData))
			if err != nil {
				sysLog := m.SystemLog{
					Message:     "ImportCompoundMasterData: Failed to import! failed to update compound: " + extPartMaster.CompoundName,
					MessageType: "ERROR",
					IsCritical:  true,
					CreatedBy:   "system_logger",
				}
				utils.CreateSystemLogInternal(sysLog)
				slog.Error("%s - error - %s", "Error making POST request", err)
			}
			defer resp.Body.Close()
		} else {
			var extPartMaster m.Compounds
			extPartMaster.CompoundName = rec.CompoundName
			extPartMaster.Description = rec.Description
			extPartMaster.SCADACode = rec.SCADACode
			extPartMaster.SAPCode = rec.SAPCode
			extPartMaster.ModifiedBy = userIDStr
			extPartMaster.ModifiedOn = time.Now()
			extPartMaster.Status = true

			jsonData, err := json.Marshal(extPartMaster)
			if err != nil {
				slog.Error("%s - error - %s", "Error in marshalling data:", err)
			}
			resp, err := http.Post(utils.RestURL+"/add-update-compound", "application/json", bytes.NewBuffer(jsonData))
			if err != nil {
				sysLog := m.SystemLog{
					Message:     "ImportCompoundMasterData: Failed to import! failed to create compound: " + rec.CompoundName,
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
		Message:     "ImportPartData: Successfully Imported compound details",
		MessageType: "SUCCESS",
		IsCritical:  false,
		CreatedBy:   "system_logger",
	}
	utils.CreateSystemLogInternal(sysLog)
	fmt.Fprintln(w, "CSV Imported Successfully")
}

func ParseCSV(file io.Reader) ([]m.CSVCompounds, error) {
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

	var records []m.CSVCompounds
	rows = utils.NormalizeRowsToLength(rows, 2)
	// Skip the first 2 rows
	for _, row := range rows {
		// if len(row) < 2 {
		// 	continue // skip incomplete rows
		// }

		// Skip row if it's a header
		if strings.ToLower(strings.TrimSpace(row[0])) == "compound name" {
			continue
		}

		record := m.CSVCompounds{
			CompoundName: strings.TrimSpace(row[0]),
			Description:  strings.TrimSpace(row[1]),
			SCADACode:    strings.TrimSpace(row[2]),
			SAPCode:      strings.TrimSpace(row[3]),
		}

		records = append(records, record)
	}

	return records, nil
}

func CompoundsSearchPagination(w http.ResponseWriter, r *http.Request) {
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
		Data       []*m.Compounds
	}

	jsonValue, err := json.Marshal(req)
	if err != nil {
		slog.Error("Error marshaling table condition", "error", err)
	}

	// Send POST request
	resp, err := http.Post(RestURL+"/get-all-compounds-by-search-pagination", "application/json", bytes.NewBuffer(jsonValue))
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

	var html strings.Builder

	tablebutton := `
	<!--html-->
			<button type="button" class="btn m-0 p-0" id="viewCompoundDetails" data-toggle="tooltip" data-placement="bottom" title="View Details"> 
				<i class="fa fa-edit mx-2" style="color: #b250ad;"></i> 
			</button>
			<!--!html-->`

	html.WriteString(`
		<!--html-->
		<!--!html-->`)

	var CompoundTable s.TableCard
	CompoundTable.CardHeading = "Part Master"
	CompoundTable.CardHeadingActions.Component = []s.CardHeadActionComponent{
		{
			ComponentName: "modelbutton",
			ComponentType: s.ActionComponentElement{
				ModelButton: s.ModelButtonAttributes{
					ID:      "add-new-compound",
					Name:    "add-new-compound",
					Type:    "button",
					Text:    "Add New Compound",
					ModelID: "#AddCompoundModel",
				},
			},
			Width: "col-4",
		},
		{
			ComponentName: "button",
			ComponentType: s.ActionComponentElement{
				Button: s.ButtonAttributes{ID: "importParts", Name: "importParts", Type: "button", Text: `Import<input style="display:none" class="clickable btn btn-success" type="file" id="pm-import-file" name="file" accept=".xlsx">`, Class: "bg-transparent", CSS: "color: #871a83; border: 1px solid #871a83;"},
			},
			Width: "col-3",
		},
		{
			ComponentName: "button",
			ComponentType: s.ActionComponentElement{
				Button: s.ButtonAttributes{ID: "downloadTmp", Name: "downloadTmp", Type: "button", Text: ` Template <i class="fa fa-download" aria-hidden="true"></i> `, Class: "bg-transparent", AdditionalAttr: `data-link="/files/sample.pdf"`, CSS: "color: #871a83; border: 1px solid #871a83;", IsDownloadBtn: true, DownloadLink: `/static/templates/import-compound-master.xlsx`},
			},
			Width: "col-3",
		},
	}
	CompoundTable.BodyTables = s.CardTableBody{Columns: []s.CardTableBodyHeadCol{
		{
			Name:       `Sr. No.`,
			IsCheckbox: true,
			ID:         "search-sn",
			Width:      "col-1",
		},
		{
			Lable: "Part Name",
			Name:  "CompoundName",
			ID:    "Customer Name",
			Width: "col-1",
		},
		{
			Lable:        "SCADA Code",
			Name:         "SCADACode",
			ID:           "scada-code",
			Width:        "col-1",
			DataField:    "SCADACode",
			IsSearchable: true,
			Type:         "input",
		},
		{
			Lable:        "SAP Code",
			Name:         "SAPCode",
			ID:           "sap-code",
			Width:        "col-1",
			DataField:    "SAPCode",
			IsSearchable: true,
			Type:         "input",
		},
		{
			Name:  "Status",
			ID:    "compound_status",
			Width: "col-1",
		},
		{
			Name:  "Action",
			Width: "col-1",
		},
	},

		ColumnsWidth: []string{"col-1", "col-1", "col-1", "col-1", "col-1", "col-1"},
		Data:         Response.Data,
		Buttons:      tablebutton,
	}

	tableBodyHTML := CompoundTable.BodyTables.RenderBodyColumns()

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

// Create compound by vendor name
func GetCompoundsByParm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")
	var compounds []m.Compounds
	// Make HTTP GET request to DB
	resp, err := http.Get(utils.RestURL + "/get-compound-data-by-parm?key=" + key + "&value=" + value)
	if err != nil {
		log.Printf("Error making GET request: %v", err)
		http.Error(w, "Failed to get user information i am in restserv", http.StatusForbidden)
		return
	}
	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&compounds)

	if err := json.NewEncoder(w).Encode(compounds); err != nil {
		log.Println("Error encoding user details:", err)
		http.Error(w, "Error to retrive user details", http.StatusInternalServerError)
		return
	}
}
