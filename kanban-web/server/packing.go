package server

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"log/slog"
	"net/http"
	"sort"
	"strconv"

	m "irpl.com/kanban-commons/model"
	"irpl.com/kanban-commons/utils"
	"irpl.com/kanban-web/pages/basepage"
	"irpl.com/kanban-web/pages/packing"
	s "irpl.com/kanban-web/services"
)

func packingDispatchPage(w http.ResponseWriter, r *http.Request) {
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
				Selected: true,
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

	// // Define the connections between steps (pairs of step indices)
	// connections := [][2]int{
	// 	{0, 1}, {1, 2}, {2, 3}, {3, 4},
	// 	{4, 5}, {5, 6}, {6, 7}, {7, 8},
	// 	{8, 9}, {9, 10},
	// }

	// Create the FlowchartPage with title, steps, and connections
	dispatchPackingPage := packing.PackingDispatchPage{}

	// Build the flowchart content for the BasePage

	// Build the complete BasePage with navigation and content
	page := (&basepage.BasePage{
		SideNavBar: sideNav,
		TopNavBar:  topNav,
		Username:   username,
		UserType:   usertype,
		Content:    dispatchPackingPage.Build(),
	}).Build()

	// Write the complete HTML page to the response
	w.Write([]byte(page))
}

func GetAllPackingKanban(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var vendors []*m.Vendors
	var vendorKanbans []m.VendorKanban

	// Define request body (optional filters)
	var req struct {
		Conditions []string `json:"Conditions"`
	}

	reqJSON, err := json.Marshal(req)
	if err != nil {
		log.Printf("Error marshaling user data: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// Fetch all vendors
	vendorResp, err := http.Post(RestURL+"/get-all-vendors-data?status=packing", "application/json", bytes.NewBuffer(reqJSON))
	if err != nil {
		slog.Error("Error fetching vendors", "error", err)
		http.Error(w, "failed to fetch vendor data", http.StatusInternalServerError)
		return
	}
	defer vendorResp.Body.Close()

	if err := json.NewDecoder(vendorResp.Body).Decode(&vendors); err != nil {
		slog.Error("Error decoding vendor response", "error", err)
		http.Error(w, "invalid vendor data", http.StatusInternalServerError)
		return
	}

	// Iterate over vendors and fetch their compounds
	for _, v := range vendors {
		var vendorKanban m.VendorKanban
		var compounds []m.CompoundsDataByVendor

		vendorKanban.Vendor = *v

		compoundResp, err := http.Post(RestURL+"/get-packing-compound-data-by-vendor?key=id&value="+strconv.Itoa(v.ID),
			"application/json", bytes.NewBuffer(reqJSON))
		if err != nil {
			slog.Error("Error fetching compounds for vendor", "vendorID", v.ID, "error", err)
			continue // Skip this vendor, proceed with next
		}

		func() {
			defer compoundResp.Body.Close()
			if err := json.NewDecoder(compoundResp.Body).Decode(&compounds); err != nil {
				slog.Error("Error decoding compound response", "vendorID", v.ID, "error", err)
				return
			}
			vendorKanban.Compounds = compounds
		}()

		vendorKanbans = append(vendorKanbans, vendorKanban)
	}

	// Return final response
	if err := json.NewEncoder(w).Encode(vendorKanbans); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func GetPackingKanbanForVendor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var vendor m.Vendors
	var req struct {
		Conditions []string `json:"Conditions"`
	}
	if err := json.NewDecoder(r.Body).Decode(&vendor); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if vendor.ID == 0 {
		utils.SetResponse(w, http.StatusInternalServerError, string("invalid vendor"))
		return
	}
	reqjsonValue, err := json.Marshal(req)
	if err != nil {
		log.Printf("Error: Error marshaling user data: %v", err)
	}
	resp, err := http.Post(RestURL+"/get-packing-compound-data-by-vendor?key=id&value="+strconv.Itoa(vendor.ID), "application/json", bytes.NewBuffer(reqjsonValue))
	if err != nil {
		slog.Error("%s - error - %s", "Error making GET request", err)
	}
	defer resp.Body.Close()
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		http.Error(w, "error reading response from service", http.StatusInternalServerError)
		return
	}
	utils.SetResponse(w, resp.StatusCode, string(responseBody))
}

func PackingKanbanSortAsce(w http.ResponseWriter, r *http.Request) {
	type SearchRequest struct {
		Criteria map[string]string `json:"Criteria"`
	}
	type TableRequest struct {
		Conditions []string `json:"conditions"`
	}
	var req struct {
		Conditions []string `json:"Conditions"`
	}
	var compstr string
	allCheckBox := `<input style="width: 20px; height: 20px;" class="form-check-input me-2 select_all_compound" type="checkbox" value="" id="select_all_compound">`
	var data SearchRequest
	var vendortable []*m.VendorCompanyTable
	var compound []m.CompoundsDataByVendor
	var searchresult []*m.Vendors
	// var orderdetails []*m.CustomerOrderDetails
	var tablereq TableRequest
	var subcondition []string
	var condition string
	var searcondition []string
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err == nil {
		for key, val := range data.Criteria {
			if key == "vendor_name" || key == "vendor_code" && val != "" {
				con := "vendors." + key + " iLIKE '%%%" + val + "%%' "
				subcondition = append(subcondition, con)
			} else if key == "compound_name" {
				compstr = val
			}
		}

		for i, v := range subcondition {
			if i < (len(subcondition) - 1) {
				condition = condition + v + " AND "
			} else if len(subcondition) == 1 {
				condition = condition + v
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

		resp, err := http.Post(utils.RestURL+"/get-all-vendors-data?status=packing", "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			slog.Error("%s - error - %s", "Error making POST request", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to fectch order details")
			return
		}
		defer resp.Body.Close()

		decoder := json.NewDecoder(resp.Body)
		if err := decoder.Decode(&searchresult); err == nil {

			for _, i := range searchresult {
				var comp string
				req.Conditions = append(req.Conditions, "compound_name iLIKE '%"+compstr+"%' ")
				reqjsonValue, _ := json.Marshal(req)
				// fetching component records by vendor
				coprresp, err := http.Post(RestURL+"/get-packing-compound-data-by-vendor?key=id&value="+strconv.Itoa(i.ID), "application/json", bytes.NewBuffer(reqjsonValue))
				if err != nil {
					slog.Error("%s - error - %s", "Error making GET request", err)
				}

				defer coprresp.Body.Close()

				if err := json.NewDecoder(coprresp.Body).Decode(&compound); err != nil {
					slog.Error("error decoding response body", "error", err)
				}

				// Sort by QualityDoneTime in ascending order
				sort.Slice(compound, func(i, j int) bool {
					return compound[i].QualityDoneTime.Before(compound[j].QualityDoneTime)
				})

				// indianLocation := time.FixedZone("IST", 5*60*60+30*60)
				for _, i := range compound {
					opt := utils.JoinStr(`
							<div class="d-flex d-inline-flex  align-items-center pl-1 m-1" style="cursor: pointer !important;">
								<label class="form-check p-0 w-100" for="`, strconv.Itoa(i.KbRootId), `" style="cursor: pointer !important;">
								<div class="col-auto border border-1 p-0 mx-2" style="border-color: #ab71a2 !important; border-radius: 6px; user-select: none; cursor: pointer; background-color:`, utils.KanbanPriorityColors[i.CustomerNote]["bg-color"], `; color:`, utils.KanbanPriorityColors[i.CustomerNote]["text-color"], `;">
									<span class="form-check p-0 m-0" style="cursor: pointer;">
										<span class="border-end border-2 p-1 pl-1" style="border-color: #ab71a2 !important; cursor: pointer;">
											<input class="form-check-input m-1 mt-2 pl-1 component-code" type="checkbox" value="`, strconv.Itoa(i.KbRootId), `" id="`, strconv.Itoa(i.KbRootId), `">
										</span>
										<label class="form-check-label m-1 px-1" for="`, strconv.Itoa(i.KbRootId), `">
											`, i.CompoundName, `
										</label>
										<span class="mx">|</span>
										<label class="form-check-label m-0 px-1" for="`, strconv.Itoa(i.KbRootId), `" data="Cell-Name">
											`, i.CellNo, `
										</label>
										<span class="mx">|</span>
										<label class="form-check-label m-0 px-1" for="`, i.KanbanNo, `" data="Kanba-Number">
											`, i.KanbanNo, `
										</label>
										
									</span>
								</div>
								</lable>
							</div>
					 `)

					comp += opt

				}
				var temptable m.VendorCompanyTable
				temptable.VendorCode = i.VendorCode
				temptable.VendorName = i.VendorName
				temptable.CompanyCodeAndNameString = comp
				if len(compound) > 0 {
					temptable.QualityDoneTime = compound[0].QualityDoneTime
				}
				if len(compound) > 0 && compstr != "" {
					vendortable = append(vendortable, &temptable)
				} else if compstr == "" {
					vendortable = append(vendortable, &temptable)
				}
			}

			sort.Slice(vendortable, func(i, j int) bool {
				if vendortable[i].QualityDoneTime.IsZero() && vendortable[j].QualityDoneTime.IsZero() {
					return false
				} else if vendortable[i].QualityDoneTime.IsZero() {
					return false
				} else if vendortable[j].QualityDoneTime.IsZero() {
					return true
				}
				return vendortable[i].QualityDoneTime.Before(vendortable[j].QualityDoneTime)
			})

			var qtyTesting s.TableCard
			tablebutton := `
		<!--html-->
			<button type="button" class="btn m-0 p-0" id="viewAllCompDetails" data-toggle="tooltip" data-placement="bottom" title="View Details"> 
 					<i class="fa fa-eye mx-2" style="color: #b250ad;"></i> 
			</button>
		<!--!html-->`

			qtyTesting.BodyTables.Data = vendortable
			qtyTesting.BodyTables.Columns = []s.CardTableBodyHeadCol{
				{
					Name:         `Vendor Code`,
					IsSearchable: true,
					IsCheckbox:   true,
					ID:           "VendorCode",
					Type:         "input",
					Width:        "col-1",
				},
				{
					Name:         "Vendor Name",
					IsSearchable: true,
					ID:           "VendorName",
					Type:         "input",
					Width:        "col-1",
				},

				{
					Lable:            "Part Name",
					Name:             "Compound Code",
					IsSearchable:     true,
					SearchFieldWidth: "w-25",
					ID:               "CompoundCode",
					Type:             "input",
					Width:            "col-10",
				},
			}
			qtyTesting.BodyTables.ColumnsWidth = []string{"col-1", "col-1", "col-9 d-flex flex-wrap w-100 "}
			qtyTesting.BodyTables.Buttons = tablebutton
			qtyTesting.BodyTables.Tools = allCheckBox
			tbody := qtyTesting.BodyTables.VendorCompanyTable()
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

func PackingKanbanSortDesc(w http.ResponseWriter, r *http.Request) {
	type SearchRequest struct {
		Criteria map[string]string `json:"Criteria"`
	}
	type TableRequest struct {
		Conditions []string `json:"conditions"`
	}
	var req struct {
		Conditions []string `json:"Conditions"`
	}
	var compstr string
	allCheckBox := `<input style="width: 20px; height: 20px;" class="form-check-input me-2 select_all_compound" type="checkbox" value="" id="select_all_compound">`
	var data SearchRequest
	var vendortable []*m.VendorCompanyTable
	var compound []m.CompoundsDataByVendor
	var searchresult []*m.Vendors
	// var orderdetails []*m.CustomerOrderDetails
	var tablereq TableRequest
	var subcondition []string
	var condition string
	var searcondition []string
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err == nil {
		for key, val := range data.Criteria {
			if key == "vendor_name" || key == "vendor_code" && val != "" {
				con := "vendors." + key + " iLIKE '%%%" + val + "%%' "
				subcondition = append(subcondition, con)
			} else if key == "compound_name" {
				compstr = val
			}
		}

		for i, v := range subcondition {
			if i < (len(subcondition) - 1) {
				condition = condition + v + " AND "
			} else if len(subcondition) == 1 {
				condition = condition + v
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

		resp, err := http.Post(utils.RestURL+"/get-all-vendors-data?status=packing", "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			slog.Error("%s - error - %s", "Error making POST request", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to fectch order details")
			return
		}
		defer resp.Body.Close()

		decoder := json.NewDecoder(resp.Body)
		if err := decoder.Decode(&searchresult); err == nil {

			for _, i := range searchresult {
				var comp string
				req.Conditions = append(req.Conditions, "compound_name iLIKE '%"+compstr+"%' ")
				reqjsonValue, _ := json.Marshal(req)
				// fetching component records by vendor
				coprresp, err := http.Post(RestURL+"/get-packing-compound-data-by-vendor?key=id&value="+strconv.Itoa(i.ID), "application/json", bytes.NewBuffer(reqjsonValue))
				if err != nil {
					slog.Error("%s - error - %s", "Error making GET request", err)
				}

				defer coprresp.Body.Close()

				if err := json.NewDecoder(coprresp.Body).Decode(&compound); err != nil {
					slog.Error("error decoding response body", "error", err)
				}

				// Sort by CreatedOn in ascending order
				sort.Slice(compound, func(i, j int) bool {
					return compound[i].QualityDoneTime.After(compound[j].QualityDoneTime)
				})

				// indianLocation := time.FixedZone("IST", 5*60*60+30*60)
				for _, i := range compound {
					opt := utils.JoinStr(`
							<div class="d-flex d-inline-flex  align-items-center pl-1 m-1" style="cursor: pointer !important;">
								<label class="form-check p-0 w-100" for="`, strconv.Itoa(i.KbRootId), `" style="cursor: pointer !important;">
								<div class="col-auto border border-1 p-0 mx-2" style="border-color: #ab71a2 !important; border-radius: 6px; user-select: none; cursor: pointer; background-color:`, utils.KanbanPriorityColors[i.CustomerNote]["bg-color"], `; color:`, utils.KanbanPriorityColors[i.CustomerNote]["text-color"], `;">
									<span class="form-check p-0 m-0" style="cursor: pointer;">
										<span class="border-end border-2 p-1 pl-1" style="border-color: #ab71a2 !important; cursor: pointer;">
											<input class="form-check-input m-1 mt-2 pl-1 component-code" type="checkbox" value="`, strconv.Itoa(i.KbRootId), `" id="`, strconv.Itoa(i.KbRootId), `">
										</span>
										<label class="form-check-label m-1 px-1" for="`, strconv.Itoa(i.KbRootId), `">
											`, i.CompoundName, `
										</label>
										<span class="mx">|</span>
										<label class="form-check-label m-0 px-1" for="`, strconv.Itoa(i.KbRootId), `" data="Cell-Name">
											`, i.CellNo, `
										</label>
										<span class="mx">|</span>
										<label class="form-check-label m-0 px-1" for="`, i.KanbanNo, `" data="Kanba-Number">
											`, i.KanbanNo, `
										</label>
										
									</span>
								</div>
								</lable>
							</div>
					 `)

					comp += opt

				}
				var temptable m.VendorCompanyTable
				temptable.VendorCode = i.VendorCode
				temptable.VendorName = i.VendorName
				temptable.CompanyCodeAndNameString = comp
				if len(compound) > 0 {
					temptable.QualityDoneTime = compound[0].QualityDoneTime
				}
				if len(compound) > 0 && compstr != "" {
					vendortable = append(vendortable, &temptable)
				} else if compstr == "" {
					vendortable = append(vendortable, &temptable)
				}
			}

			sort.Slice(vendortable, func(i, j int) bool {
				if vendortable[i].QualityDoneTime.IsZero() && vendortable[j].QualityDoneTime.IsZero() {
					return false
				} else if vendortable[i].QualityDoneTime.IsZero() {
					return false
				} else if vendortable[j].QualityDoneTime.IsZero() {
					return true
				}
				return vendortable[i].QualityDoneTime.After(vendortable[j].QualityDoneTime)
			})

			var qtyTesting s.TableCard
			tablebutton := `
		<!--html-->
			<button type="button" class="btn m-0 p-0" id="viewAllCompDetails" data-toggle="tooltip" data-placement="bottom" title="View Details"> 
 					<i class="fa fa-eye mx-2" style="color: #b250ad;"></i> 
			</button>
		<!--!html-->`

			qtyTesting.BodyTables.Data = vendortable
			qtyTesting.BodyTables.Columns = []s.CardTableBodyHeadCol{
				{
					Name:         `Vendor Code`,
					IsSearchable: true,
					IsCheckbox:   true,
					ID:           "VendorCode",
					Type:         "input",
					Width:        "col-1",
				},
				{
					Name:         "Vendor Name",
					IsSearchable: true,
					ID:           "VendorName",
					Type:         "input",
					Width:        "col-1",
				},

				{
					Lable:            "Part Name",
					Name:             "Compound Code",
					IsSearchable:     true,
					SearchFieldWidth: "w-25",
					ID:               "CompoundCode",
					Type:             "input",
					Width:            "col-10",
				},
			}
			qtyTesting.BodyTables.ColumnsWidth = []string{"col-1", "col-1", "col-9 d-flex flex-wrap w-100 "}
			qtyTesting.BodyTables.Buttons = tablebutton
			qtyTesting.BodyTables.Tools = allCheckBox
			tbody := qtyTesting.BodyTables.VendorCompanyTable()
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
