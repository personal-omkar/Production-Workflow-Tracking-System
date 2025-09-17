package server

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"time"

	m "irpl.com/kanban-commons/model"
	re "irpl.com/kanban-commons/report"
	"irpl.com/kanban-commons/utils"
	"irpl.com/kanban-web/pages/basepage"
	summaryreprint "irpl.com/kanban-web/pages/summaryReprint"
)

func SummaryReprintPage(w http.ResponseWriter, r *http.Request) {
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
				Selected: true,
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
	kanbanreprintpage := summaryreprint.SummaryReprintPage{
		UserType: usertype,
		UserID:   userID,
	}

	if len(vendorRecord) > 0 {
		kanbanreprintpage.VendorCode = vendorRecord[0].VendorCode
	}

	// Build the complete BasePage with navigation and content
	page := (&basepage.BasePage{
		SideNavBar: sideNav,
		TopNavBar:  topNav,
		Username:   username,
		UserType:   usertype,
		Content:    kanbanreprintpage.Build(),
	}).Build()

	// Write the complete HTML page to the response
	w.Write([]byte(page))
}

func PrintKanbanSummary2(w http.ResponseWriter, r *http.Request) {
	var kanbanSummary struct {
		Type   string
		CellNo []string
	}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&kanbanSummary); err == nil {
		pdfOutDir := "/RUBBER/pdf/"
		reportName := "test"
		xlsxOut := fmt.Sprint("/RUBBER/xlsx/temp-summary-" + time.Now().Format("02-01-2006"))
		report := &re.Report{
			Template: "/RUBBER/template/rub-summary-template.xlsx",
			OutPath:  xlsxOut,
		}

		if len(kanbanSummary.CellNo) > 40 {
			utils.SetResponse(w, http.StatusInternalServerError, "To many values are selected")
		} else {
			var rowCount = 7
			var vendorName, vendorCode string

			for _, v := range kanbanSummary.CellNo {

				var rawQuery m.RawQuery
				// KBdata details
				rawQuery.Type = "KbData"
				rawQuery.Host = utils.RestHost
				rawQuery.Port = utils.RestPort
				rawQuery.Query = utils.JoinStr(`SELECT * FROM Kb_data WHERE cell_no = '`, v, `'`) // `;`
				var kbData []m.KbData
				rawQuery.RawQry(&kbData)

				// Compound Details
				rawQuery.Type = "Compounds"
				rawQuery.Host = utils.RestHost
				rawQuery.Port = utils.RestPort
				rawQuery.Query = utils.JoinStr(`SELECT * FROM Compounds WHERE id = `, fmt.Sprint(kbData[0].CompoundId)) // `;`
				var Comps []m.Compounds
				rawQuery.RawQry(&Comps)

				// KbExtension Details
				rawQuery.Type = "KbExtension"
				rawQuery.Host = utils.RestHost
				rawQuery.Port = utils.RestPort
				rawQuery.Query = utils.JoinStr(`SELECT * FROM kb_extension WHERE id = `, fmt.Sprint(kbData[0].KbExtensionID)) // `;`
				var KbExt []m.KbExtension
				rawQuery.RawQry(&KbExt)

				// Vendor Details
				rawQuery.Type = "Vendors"
				rawQuery.Host = utils.RestHost
				rawQuery.Port = utils.RestPort
				rawQuery.Query = utils.JoinStr(`SELECT * FROM Vendors WHERE id = `, fmt.Sprint(KbExt[0].VendorID)) // `;`
				var Vendors []m.Vendors
				rawQuery.RawQry(&Vendors)

				values := []string{
					fmt.Sprintf("${A%d}%v", rowCount, (rowCount - 6)), // SrNo: rowCount - excaped rows (6)
					fmt.Sprintf("${B%d}%v", rowCount, Comps[0].CompoundName),
					fmt.Sprintf("${C%d}%v", rowCount, kbData[0].CellNo),
					fmt.Sprintf("${D%d}%v", rowCount, kbData[0].DemandDateTime.Format("02.01.2006 15.04.03")),
					fmt.Sprintf("${E%d}%v", rowCount, kbData[0].NoOFLots),
				}
				report.Values = append(report.Values, values...)
				vendorName = Vendors[0].VendorName
				vendorCode = Vendors[0].VendorCode
				rowCount++
			}
			report.Values = append(report.Values, fmt.Sprintf("${A2}image:%v", "./static/icons/logo_black.png"))
			report.Values = append(report.Values, fmt.Sprintf("${B2}Vendor Name: %v", vendorName))
			report.Values = append(report.Values, fmt.Sprintf("${B3}Vendor Code: %v", vendorCode))
			report.Values = append(report.Values, fmt.Sprintf("${D3}%v", time.Now().Format("02.01.2006 15.04.03")))
			report.Values = append(report.Values, fmt.Sprintf("${D39}%v", vendorName))
			report.OutPath = fmt.Sprint("/RUBBER/xlsx/" + vendorCode + "-summary-" + time.Now().Format("02-01-2006_15.04.03") + ".xlsx")
			xlsxPath, _ := report.CreateXLSXPage()

			// Create temp xlsx file (with report sheet only)
			tempFilePath := fmt.Sprint("/RUBBER/" + vendorCode + "-summary-" + time.Now().Format("02-01-2006_15.04.03") + ".xlsx")
			// Convert the new Excel file to PDF
			pdfFilePath, err := re.ConvertExcelToPDF("report", xlsxPath, tempFilePath, pdfOutDir)
			if err != nil {
				log.Println(err)
			}
			re.HandleFileWithTimeout(pdfFilePath, tempFilePath)

			d := re.PdfDialog{
				RecordId:   reportName,
				PageSource: pdfFilePath,
			}

			utils.SetResponse(w, http.StatusOK, d.Build())
		}
	} else {
		slog.Error("PrintKanbanSummary: error decoding request body", "error", err)
	}
}
