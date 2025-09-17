package dashboard

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"

	m "irpl.com/kanban-commons/model"
	u "irpl.com/kanban-commons/utils"
	"irpl.com/kanban-web/pages/basepage"
	"irpl.com/kanban-web/services"
)

type DashboardPage struct {
	Username   string
	Role       string
	Config     map[string]string
	SideNavBar basepage.SideNav
	TopNavBar  basepage.TopNav
	Content    string
	BgColor    string
	UserRole   string
}

const DefaultRestHost string = "0.0.0.0" // Default port if not set in env
const DefaultRestPort string = "4300"    // Default port if not set in env

var RestHost string // Global variable to hold the Rest helper host
var RestPort string // Global variable to hold the Rest helper port
var RestURL string  // Global variable to hold the Rest URL

func init() {
	RestHost = os.Getenv("RESTSRV_HOST")
	if strings.TrimSpace(RestHost) == "" {
		RestHost = DefaultRestHost
	}

	RestPort = os.Getenv("RESTSRV_PORT")
	if strings.TrimSpace(RestPort) == "" {
		RestPort = DefaultRestPort
	}

	RestURL = u.JoinStr("http://", RestHost, ":", RestPort)
}

func RegisterPage(w http.ResponseWriter, r *http.Request) {
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

	sideNav := basepage.SideNav{}
	sideNav.MenuItems = []basepage.SideNavItem{
		{
			Name:     "Dashboard",
			Icon:     "fas fa-chart-pie",
			Link:     "/dashboard",
			UserType: basepage.UserType{Admin: true, Operator: true, Customer: true},
			Style:    "font-size:1rem;",
			Selected: true,
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
			Name:     u.DefaultsMap["cold_store_menu"],
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
		// {
		// 	Name:  "Production Line Status",
		// 	Icon:  "fas fa-calendar-day",
		// 	Link:  "/flowchart?line=1",
		// 	Style: "font-size:1rem;",
		// },
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
	}
	sideNav.MenuItems = basepage.CheckdisabledNavItems(sideNav.MenuItems, links, "|")
	topNav := basepage.TopNav{VendorName: vendorName, UserType: usertype}
	topNav.MenuItems = []basepage.TopNavItem{
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
	}

	var dashboardStats m.DashboardStatsResponse
	rawQuery := m.RawQuery{
		Host: u.RestHost,
		Port: u.RestPort,
		Type: "DashboardStats",
		Query: `
		SELECT
 
			COUNT(*) FILTER (
				WHERE DATE_TRUNC('month', submit_date_time) = DATE_TRUNC('month', CURRENT_DATE)
			) AS monthly_submitted,

			COUNT(*) FILTER (
				WHERE DATE(submit_date_time) = CURRENT_DATE
			) AS daily_submitted,

			(
				SELECT COUNT(*) 
				FROM kb_root kr
				JOIN kb_data kd ON kr.kb_data_id = kd.id
				WHERE 
				kr.status = '4'
				AND DATE_TRUNC('month', dispatch_done_time) = DATE_TRUNC('month', CURRENT_DATE)
				AND DATE_TRUNC('month', kd.submit_date_time) = DATE_TRUNC('month', CURRENT_DATE)
			) AS monthly_dispatched,

			(
				SELECT COUNT(*) 
				FROM kb_root kr
				JOIN kb_data kd ON kr.kb_data_id = kd.id
				WHERE 
				kr.status = '4'
				AND DATE(dispatch_done_time) = CURRENT_DATE
				AND DATE(kd.submit_date_time) = CURRENT_DATE
			) AS daily_dispatched,

			(
				SELECT COUNT(*) 
				FROM kb_extension 
				WHERE status = 'reject' 
				AND DATE_TRUNC('month', created_on) = DATE_TRUNC('month', CURRENT_DATE)
			) AS monthly_rejected

			FROM kb_data;
		`,
	}

	if err := rawQuery.RawQry(&dashboardStats); err != nil {
		slog.Error("Failed to fetch dashboard stats using RawQuery", "error", err)
	}

	dashboardStats.DailyPercentage = 0
	if dashboardStats.DailySubmitted > 0 {
		dashboardStats.DailyPercentage = int(float64(dashboardStats.DailyDispatched) / float64(dashboardStats.DailySubmitted) * 100)
	}
	dashboardStats.MonthlyPercentage = 0
	if dashboardStats.MonthlySubmitted > 0 {
		dashboardStats.MonthlyPercentage = int(float64(dashboardStats.MonthlyDispatched) / float64(dashboardStats.MonthlySubmitted) * 100)
	}

	dashboardStats.MonthlyRejPercentage = 0
	if dashboardStats.MonthlySubmitted > 0 {
		dashboardStats.MonthlyRejPercentage = int(float64(dashboardStats.MonthlyRejected) / float64(dashboardStats.MonthlySubmitted) * 100)
	}

	page := (&DashboardPage{
		Username:   username,
		SideNavBar: sideNav,
		TopNavBar:  topNav,
		UserRole:   usertype,
		Content:    BuildDashboardCards(dashboardStats),
	}).Build()
	w.Write([]byte(page))
}

func (d *DashboardPage) Build() string {
	out := (&basepage.BasePage{
		Username:   d.Username,
		SideNavBar: d.SideNavBar,
		TopNavBar:  d.TopNavBar,
		Content:    d.Content,
		UserType:   d.UserRole,
	}).Build()
	return out
}

func BuildDashboardCards(stats m.DashboardStatsResponse) string {
	dailySubmit := strconv.Itoa(stats.DailySubmitted)
	monthlySubmit := strconv.Itoa(stats.MonthlySubmitted)
	dailyPct := strconv.Itoa(stats.DailyPercentage)
	monthlyPct := strconv.Itoa(stats.MonthlyPercentage)
	monthlyRejPct := strconv.Itoa(stats.MonthlyRejPercentage)
	var rawQuery m.RawQuery
	rawQuery.Host = u.RestHost
	rawQuery.Port = u.RestPort
	var prodLines []*m.ProdLine
	rawQuery.Type = "ProdLine"
	rawQuery.Query = `SELECT * FROM prod_line WHERE status='true';`
	_ = rawQuery.RawQry(&prodLines)

	var machineOptions string
	//var defaultCount int
	for i, line := range prodLines {
		line.Name = strings.TrimSpace(strings.TrimSuffix(line.Name, "(Heijunka)"))
		machineOptions += `<option value="` + strconv.Itoa(line.Id) + `">` + line.Name + `</option>`

		if i == 0 {
			var result struct{ Count int }
			rawQuery.Type = "InProgressOrderCount"
			rawQuery.Query = BuildInProgressOrderQuery(line.Id)
			_ = rawQuery.RawQry(&result)
			//defaultCount = result.Count
		}
	}

	cards := []services.Card{
		{ //<!--html-->
			Body: u.JoinStr(`
			<div class="card-body">
						<div class="container row">
							<div class="col-3 rounded-circle d-flex justify-content-center align-items-center" style="background:#ab71a2;">
								<i class="fas fa-calendar-day" style="font-size:2rem; color:white"></i>
							</div>
							<div class="col-9">
								<div style="font-size: 2rem;font-weight: bold;">
									`, dailySubmit, `
								</div>
								<div>
									Todays Orders	
								</div>
							</div>
						</div>
						<div class="container">
							<div class="progress mt-3" style="height:5px">
								<div class="progress-bar" role="progressbar" style="width: `, dailyPct, `%" aria-valuenow="25" aria-valuemin="0" aria-valuemax="100">
								</div>
							</div>
						</div>
						<div class="container" style="font-size: 1.5rem;">
								`, dailyPct, `% Completed
						</div>
			     	</div>
		`),
			//<!--!html-->
			Width: "col-3",
		},
		{
			//<!--html-->
			Body: u.JoinStr(`
			<div class="card-body">
						<div class="container row">
							<div class="col-3 rounded-circle d-flex justify-content-center align-items-center" style="background:#ab71a2;">
								<i class="fas fa-calendar-check" style="font-size:2rem; color:white"></i>
							</div>
							<div class="col-9">
								<div style="font-size: 2rem;font-weight: bold;">
									`, monthlySubmit, `
								</div>
								<div>
									Monthly Orders	
								</div>
							</div>
						</div>
						<div class="container">
							<div class="progress mt-3" style="height:5px">
								<div class="progress-bar" role="progressbar" style="width: `, monthlyPct, `%" aria-valuenow="25" aria-valuemin="0" aria-valuemax="100">
								</div>
							</div>
						</div>
						<div class="container" style="font-size: 1.5rem;">
								`, monthlyPct, `% Completed
						</div>
			     	</div>
		`),
			//<!--!html-->
			Width: "col-3",
		},
		{ //<!--html-->
			Body: u.JoinStr(`
			<div class="card-body">
						<div class="container row">
							<div class="col-3 rounded-circle d-flex justify-content-center align-items-center" style="background:#ab71a2;">
								<i class="fas fa-ban" style="font-size:2rem; color:white"></i>
							</div>
							<div class="col-9">
								<div style="font-size: 2rem;font-weight: bold;">
									`, monthlySubmit, `
								</div>
								<div>
								Monthly Orders	
								</div>
							</div>
						</div>
						<div class="container">
							<div class="progress mt-3" style="height:5px">
								<div class="progress-bar" role="progressbar" style="width: `, monthlyRejPct, `%" aria-valuenow="25" aria-valuemin="0" aria-valuemax="100">
								</div>
							</div>
						</div>
						<div class="container" style="font-size: 1.5rem;">
								`, monthlyRejPct, `% Rejected
						</div>
			     	</div>
		`),
			//<!--!html-->
			Width: "col-3",
		},
		{
			Width: "col-3",
			//<!--html-->
			Body: u.JoinStr(`
			<div class="card-body position-relative">
					<div class="position-absolute top-0 end-0 mt-2 me-3" style="width: 140px;">
						<select id="machineSelect" class="form-select form-select-sm" onchange="updateMachineOrderCount()" style="font-size: 0.85rem;">
							`, machineOptions, `
						</select>
					</div>
				<div class="row" style="margin-left:1px;">
					<div class="col-3 rounded-circle d-flex justify-content-center align-items-center" style="background: #ab71a2; height: 68px; width: 68px;">
						<i class="fas fa-play" style="font-size: 2rem; color: white;"></i>
					</div>
					<div class="col-9">
						<div id="machineScheduledCount" style="font-size: 2rem; font-weight: bold;">0</div>
						<div>Scheduled Kanban's</div>
					</div>
				</div>
				<div class="mt-3">
					<div class="progress" style="height:5px;">
						<div class="progress-bar" role="progressbar" style="width: 100%;" aria-valuenow="25" aria-valuemin="0" aria-valuemax="100"></div>
					</div>
				</div>
				<div class="d-flex align-items-center mt-0" style="font-size: 1.5rem;">
					<div id="machineInProgressCount">0</div>
					<span class="ms-2">In-Progress Kanban's</span>
				</div>
			</div>
			`),
			//<!--!html-->
		},

		{
			//<!--html-->
			Width: "col-12",
			Style: "margin-top: 2rem;",

			Body: u.JoinStr(`
			<div class="card-header d-flex justify-content-between align-items-center">
				<h5 class="mb-0">Kanban Status Overview</h5>
				<div class="form-group mb-0" id="month-wrapper" style="cursor: pointer;">
					<input type="month" id="vendor-status-month" class="form-control form-control-lm" style="text-align: center;" />
				</div>
			</div>
			<div class="card-body">
				<div class="echart-horizontal-stacked-chart-example"
					style="width: 100%; height: 500px;"
					data-echart-responsive="true">
				</div>
			</div>
		`),
			//<!--!html-->
		},
	}

	var content string
	content += `<div class="row">`
	for _, card := range cards {
		content += card.Build()
	}
	content += `

		<script>
			function updateMachineOrderCount() {
				const dropdown = document.getElementById("machineSelect");
				const selectedMachineId = parseInt(dropdown.value);

				if (!selectedMachineId) {
					document.getElementById("machineScheduledCount").textContent = "0";
					document.getElementById("machineInProgressCount").textContent = "0";
					return;
				}

				fetch("/get-in-progress-count-by-line", {
					method: "POST",
					headers: { "Content-Type": "application/json" },
					body: JSON.stringify({ prod_line_id: selectedMachineId })
				})
				.then(response => {
					if (!response.ok) {
						throw new Error("HTTP error " + response.status);
					}
					return response.json();
				})
				.then(data => {
					document.getElementById("machineScheduledCount").textContent = data.scheduled_count || 0;
					document.getElementById("machineInProgressCount").textContent = data.in_progress_count || 0;
				})
				.catch(error => {
					console.error("Fetch error:", error);
					document.getElementById("machineScheduledCount").textContent = "0";
					document.getElementById("machineInProgressCount").textContent = "0";
				});
			}

			document.addEventListener("DOMContentLoaded", function () {
				updateMachineOrderCount();
				const dropdown = document.getElementById("machineSelect");
				const optionsToHide = ["quality", "packing"];
				for (let i = dropdown.options.length - 1; i >= 0; i--) {
					const optionText = dropdown.options[i].textContent.toLowerCase().trim();
					if (optionsToHide.includes(optionText)) {
						dropdown.remove(i);
					}
				}
			});
	</script>`

	content += `</div>`
	return content
}
