package server

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	m "irpl.com/kanban-commons/model"
	"irpl.com/kanban-commons/utils"
	allkanbanview "irpl.com/kanban-web/pages/allKanbanView"
	"irpl.com/kanban-web/pages/basepage"
	"irpl.com/kanban-web/services"
)

func AllKanbanView(w http.ResponseWriter, r *http.Request) {
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
				Selected: true,
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

	allKanbanView := allkanbanview.AllKanbanView{
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
		Content:    allKanbanView.Build(),
	})
	page.AddScriptLink("/static/js/partmanagement.js")

	// Write the complete HTML page to the response
	w.Write([]byte(page.Build()))
}

func SearchAllKanbanView(w http.ResponseWriter, r *http.Request) {
	type SearchRequest struct {
		Criteria map[string]string `json:"Criteria"`
	}

	var ProdLine []*m.ProdLine
	var AllKanbanViewTable []*m.AllKanbanViewTable

	var data SearchRequest
	var MachineAndKanban []*m.AllKanbanViewTable
	var details []*m.AllKanbanViewDetails
	// var tablereq TableRequest
	var subcondition []string
	var condition string
	var partcondition string
	// var searcondition []string
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err == nil {

		for key, val := range data.Criteria {
			if key == "machine_name" && val != "" {
				con := "status='true' AND name" + " iLIKE '%%%" + val + "%%'   "
				subcondition = append(subcondition, con)
			} else if key == "machine_name" && val == "" {
				con := "status='true'"
				subcondition = append(subcondition, con)
			} else if key == "part_name_or_kanabn_no" && val != "" {
				con := " ( c.compound_name" + " iLIKE '%%%" + val + "%%' " + " OR " + "  kr.kanban_no" + " iLIKE '%%%" + val + "%%' )"
				partcondition = con
			}
		}

		for i, v := range subcondition {
			if i < (len(subcondition) - 1) {
				condition = condition + v + " AND "
			} else {
				condition = condition + v
			}
		}
		if condition != "" {
			condition = " where " + condition
		}
		if partcondition != "" {
			partcondition = " AND " + partcondition
		}

		var rawQuery m.RawQuery
		rawQuery.Host = utils.RestHost
		rawQuery.Port = utils.RestPort
		rawQuery.Type = "ProdLine"
		rawQuery.Query = `select * from prod_line` + condition
		rawQuery.RawQry(&ProdLine)

		for i := range ProdLine {
			ProdLine[i].Name = strings.TrimSpace(strings.TrimSuffix(ProdLine[i].Name, "(Heijunka)"))
		}
		for _, v := range ProdLine {
			rawQuery.Type = "AllKanbanViewDetails"
			rawQuery.Query = `
			SELECT
				kr.id,
				c.compound_name AS compound_name,
				kr.kanban_no AS kanban_no,
				kd.note AS notes,
				kd.cell_no as cell_number
			FROM
				kb_transaction kt
			JOIN
				kb_root kr ON kt.kb_root_id = kr.id
			JOIN
				kb_data kd ON kr.kb_data_id = kd.id
			JOIN
				kb_extension ke ON kd.kb_extension_id = ke.id
			JOIN
				compounds c ON kd.compound_id = c.id
			JOIN
				prod_process_line ppl ON kt.prod_process_line_id = ppl.id
			JOIN
				prod_line pl ON ppl.prod_line_id = pl.id
			WHERE
				pl.status = true
			AND
				pl.id = ` + strconv.Itoa(v.Id) + `
			AND 
				ppl.prod_process_id IN (1)
			AND (
				SELECT COUNT(*)
				FROM kb_transaction sub_kt
				WHERE sub_kt.kb_root_id = kt.kb_root_id
				) <=(
			        SELECT COUNT(*)          
			        FROM prod_process_line
			        WHERE prod_line_id =` + strconv.Itoa(v.Id) + `
			      )
			AND NOT EXISTS (
				SELECT 1
				FROM kb_transaction kt2
				JOIN prod_process_line ppl2 ON kt2.prod_process_line_id = ppl2.id
				WHERE kt2.kb_root_id = kt.kb_root_id
				AND ppl2.prod_process_id = 2
			)
			` + partcondition + ` 
			GROUP BY
			kr.id, c.compound_name, kr.kanban_no, kd.note, kd.cell_no
			HAVING COUNT(kt.kb_root_id) <=2
			ORDER BY
				kr.running_no ASC;	
			` //`;`
			log.Printf("query %v", rawQuery.Query)
			rawQuery.RawQry(&details)

			var allcards string
			if len(details) > 0 {
				for _, v := range details {
					card := `
					<div class="d-flex border border-1 p-0 mx-2 my-2" style="background-color:` + utils.KanbanPriorityColors[strings.ToLower(v.CustomerNotes)]["bg-color"] + `; color:` + utils.KanbanPriorityColors[strings.ToLower(v.CustomerNotes)]["text-color"] + `; border-color: #ab71a2 !important; border-radius: 6px; user-select: none; cursor: pointer;">
							<span class="form-check p-0 m-0" style="cursor: pointer;">
								<label class="form-check-label m-1 px-1" for="290">
									` + v.CompoundName + `
								</label>
								<span class="mx">|</span>
								<label class="form-check-label m-0 px-1" for="290" data="Cell-Name">
								` + v.KanbanNo + `
								</label>	
						</span>
					</div>`
					allcards = allcards + card
				}
				v.Name = `<a href="/production-line">` + v.Name + `</a> <span style="margin-left:10%">   (` + strconv.Itoa(len(details)) + ") </span>"
			} else {
				v.Name = `<a href="/production-line">` + v.Name + `</a>`
			}

			var temptable m.AllKanbanViewTable
			temptable.MachineId = strconv.Itoa(v.Id)
			temptable.MachineName = v.Name
			temptable.PartNameorKanbanNo = allcards

			AllKanbanViewTable = append(AllKanbanViewTable, &temptable)
		}

		var allkanbanview services.TableCard
		if partcondition != "" {
			for _, v := range AllKanbanViewTable {
				if len(v.PartNameorKanbanNo) != 0 {
					MachineAndKanban = append(MachineAndKanban, v)
				}
			}
			allkanbanview.BodyTables.Data = MachineAndKanban
		} else {
			allkanbanview.BodyTables.Data = AllKanbanViewTable
		}

		allkanbanview.BodyTables.Columns = []services.CardTableBodyHeadCol{
			{
				Lable:        "Machine Name",
				Name:         `Machine Name`,
				ID:           "machine-name",
				IsSearchable: true,
				Type:         "input",
				Width:        "col-2",
			},
			{
				Lable:        "Part Name or Kanban No",
				Name:         "Part Name or Kanban No",
				Type:         "input",
				IsSearchable: true,
				ID:           "part-name-or-kanban-no",
				Width:        "col-10",
				Style:        "max-width: 1300px;",
			},
		}
		allkanbanview.BodyTables.ColumnsWidth = []string{"col-2", "col-10 w-100 d-flex flex-wrap"}
		tbody := allkanbanview.BodyTables.AllKanbanViewTable()
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(tbody))

	} else {
		slog.Error("Record creation failed - " + err.Error())
		utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to create Pagination")
	}

}
