package server

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	m "irpl.com/kanban-commons/model"
	"irpl.com/kanban-commons/utils"
	"irpl.com/kanban-web/pages/basepage"
	"irpl.com/kanban-web/pages/kanbanhistorypage"
	s "irpl.com/kanban-web/services"
)

func kanabnHistoryPage(w http.ResponseWriter, r *http.Request) {
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
				Selected: true,
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

	historyPagePage := kanbanhistorypage.Historypage{
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

func GetCompletedKBRootDetailsBySearch(w http.ResponseWriter, r *http.Request) {
	type SearchRequest struct {
		Criteria map[string]string `json:"Criteria"`
	}
	type TableRequest struct {
		Conditions []string `json:"conditions"`
	}
	var data SearchRequest
	var searchresult []*m.OrderDetails
	var orderdetails []*m.CustomerOrderDetails
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

		resp, err := http.Post(utils.RestURL+"/get-all-kbRoot-details-by-search", "application/json", bytes.NewBuffer(jsonValue))
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
			<button type="button" class="btn m-0 p-0" id="viewAllCompDetails" data-toggle="tooltip" data-placement="bottom" title="View Details"> 
 					<i class="fa fa-eye mx-2" style="color: #b250ad;"></i> 
			</button>
		<!--!html-->`

			for _, value := range searchresult {
				value.Status = strings.Title(value.Status)
			}

			for _, val := range searchresult {
				var temp m.CustomerOrderDetails
				if val.Status == "3" || val.Status == "4" || val.Status == "-1" {
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
					temp.Status = map[string]string{"-1": "Quality Fail", "3": "Packing", "4": "Dispatched"}[val.Status]
					temp.VendorName = val.VendorName
					temp.CompoundName = val.CompoundName
					temp.OrderId = val.OrderId
					temp.MinQuantity = val.MinQuantity
					temp.AvailableQuantity = val.AvailableQuantity
					orderdetails = append(orderdetails, &temp)
				}

			}
			odcard.BodyTables.Data = orderdetails
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
					Name:  "Status",
					Width: "col-1",
				},
				{
					Name:  "Action",
					Width: "col-1",
				},
			}
			odcard.BodyTables.ColumnsWidth = []string{"col-1", "col-1", "col-1", "col-1", "col-1", "col-1", "col-1"}
			odcard.BodyTables.Buttons = tablebutton
			tbody := odcard.BodyTables.CustomerOrderDetailsTable()
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
