package productionprocesses

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	m "irpl.com/kanban-commons/model"
	"irpl.com/kanban-commons/utils"
	"irpl.com/kanban-web/pages/basepage"
	s "irpl.com/kanban-web/services"
)

type ProdProcessManagement struct {
	Username string
	UserType string
	UserID   string
}

const (
	uploadDir = "./static/uploaded-icons/uploadedIcons"
)

func ProdProcessesCRUDPage(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("X-Custom-Username")
	usertype := r.Header.Get("X-Custom-Role")
	userID := r.Header.Get("X-Custom-Userid")
	links := r.Header.Get("X-Custom-Allowlist")
	var vendorName string
	var vendorRecord []m.Vendors
	if usertype != "Admin" {
		resp, err := http.Get(utils.RestURL + "/get-vendor-by-userid?key=user_id&value=" + userID)
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
				Selected: true,
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
			// {
			// 	Name:  "Production Line Status",
			// 	Icon:  "fas fa-calendar-day",
			// 	Link:  "/flowchart?line=1",
			// 	Style: "font-size:1rem;",
			// },
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
	topNav := basepage.TopNav{VendorName: vendorName,
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

	ProdProcessManagement := ProdProcessManagement{
		UserType: usertype,
		UserID:   userID,
	}

	// Build the complete BasePage with navigation and content
	page := (&basepage.BasePage{
		SideNavBar: sideNav,
		TopNavBar:  topNav,
		Username:   username,
		UserType:   usertype,
		Content:    ProdProcessManagement.Build(),
	}).Build()

	// Write the complete HTML page to the response
	w.Write([]byte(page))
}

func (c *ProdProcessManagement) Build() string {

	var Response struct {
		Pagination m.PaginationResp
		Data       []*m.ProdProcess
	}
	type TableConditions struct {
		Pagination m.PaginationReq `json:"pagination"`
		Conditions []string        `json:"conditions"`
	}

	var tablecondition TableConditions

	// Temporary pagination data
	tablecondition.Pagination = m.PaginationReq{
		Type:   "ProdProcess",
		Limit:  "10",
		PageNo: 1,
	}

	// Convert to JSON
	jsonValue, err := json.Marshal(tablecondition)
	if err != nil {
		slog.Error("Error marshaling table condition", "error", err)
	}

	resp, err := http.Post(utils.RestURL+"/get-all-production-process-by-search-paginations", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		slog.Error("%s - error - %s", "Error making GET request", err)
	}

	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&Response); err != nil {
		slog.Error("error decoding response body", "error", err)
	}

	var html strings.Builder

	tablebutton := `
	<!--html-->
			<button type="button" class="btn m-0 p-0" id="ViewProdProcessesDetails" data-toggle="tooltip" data-placement="bottom" title="View Details"> 
				<i class="fa fa-edit mx-2" style="color: #b250ad;"></i> 
			</button>
			<!--!html-->`

	html.WriteString(`
		<!--html-->
		<!--!html-->`)

	var ProdProcessesTable s.TableCard

	ProdProcessesTable.CardHeading = "Production Process Master"
	ProdProcessesTable.CardHeadingActions.Component = []s.CardHeadActionComponent{{ComponentName: "modelbutton", ComponentType: s.ActionComponentElement{ModelButton: s.ModelButtonAttributes{ID: "add-new-prod-process", Name: "add-new-processes", Type: "button", Text: "Add New Process", ModelID: "#AddNewProcesses"}}, Width: "col-4"}}
	ProdProcessesTable.BodyTables = s.CardTableBody{Columns: []s.CardTableBodyHeadCol{
		{
			Name:       `Sr. No.`,
			IsCheckbox: true,
			ID:         "search-sn",
			Width:      "col-1",
		},
		{
			Lable:        "Processes Name",
			Name:         "Name",
			ID:           "prod-line-name",
			Width:        "col-2",
			Type:         "input",
			DataField:    "name",
			IsSearchable: true,
		},
		{
			Lable: "Status",
			Name:  "Status",
			ID:    "status",
			Width: "col-1",
		},
		{
			Lable: "Line Visibility",
			Name:  "line_visibility",
			ID:    "LineVisibility",
			Width: "col-1",
		},
		{
			Lable: "Estimated Average Time",
			Name:  "expected_mean_time",
			ID:    "ExpectedAverageTime",
			Width: "col-1",
		},
		{
			Name:  "Action",
			Width: "col-1",
		},
	},
		ColumnsWidth: []string{"col-1", "col-2", "col-1", "col-1", "col-1", "col-1"},
		Data:         Response.Data,
		Buttons:      tablebutton,
	}
	var Pagination s.Pagination
	Pagination.CurrentPage = 1
	Pagination.Offset = 0
	Pagination.TotalRecords = Response.Pagination.TotalNo
	Pagination.PerPage, _ = strconv.Atoi(tablecondition.Pagination.Limit)
	Pagination.PerPageOptions = []int{10, 25, 50, 100, 200}
	ProdProcessesTable.CardFooter = Pagination.Build()

	js := `
	<script>
		//js
		function pagination(pageNo) {
			let searchCriteria = [];

			document.querySelectorAll("[id^='search-by']").forEach(function (input) {
				let field = input.dataset.field;
				let value = input.value.trim();
				if (value) {
					searchCriteria.push(field + " ILIKE '%" + value + "%'");
				}
			});

			// Get selected limit
			let limit = document.getElementById("perPageSelect")?.value || "15";
			let requestData = {
				pagination: {
					Limit: limit,
					Pageno: pageNo
				},
				Conditions: searchCriteria
			};

			fetch("/prodprocess-search-pagination", {
				method: "POST",
				headers: {
					"Content-Type": "application/json"
				},
				body: JSON.stringify(requestData)
			})
			.then(response => {
				if (!response.ok) {
					throw new Error("HTTP error! Status: " + response.status);
				}
				return response.json();
			})
			.then(data => {
				let tableBody = document.getElementById("advanced-search-table-body");
				let cardFooter = document.getElementsByClassName("card-footer")[0];

				if (tableBody) {
					tableBody.innerHTML = "";
					tableBody.insertAdjacentHTML("beforeend", data.tableBodyHTML);
					cardFooter.innerHTML = "";
					cardFooter.insertAdjacentHTML("beforeend", data.paginationHTML);

					attachPerPageListener();
				} else {
					console.error("Error: Table body element not found.");
				}
			})
			.catch(error => console.error("Error fetching paginated results:", error));
		}

		// Function to attach event listener for perPageSelect dropdown
		function attachPerPageListener() {
			let perPageSelect = document.getElementById("perPageSelect");
			if (perPageSelect) {
				perPageSelect.addEventListener("change", function () {
					pagination(1);
				});
			}
		}
	document.addEventListener("DOMContentLoaded", function () {
		document.querySelectorAll("[id^='search-by']").forEach(function (input) {
			input.addEventListener("input", function () {
				pagination(1);
			});
		});
		
		attachPerPageListener();
		// Notification Functionality
		function showNotification(status, msg, callback) {
			const notification = $('#notification');
			let message = "";

			if (status === "200") {
				message =msg;
				notification.removeClass("alert-danger").addClass("alert-success");
			} else {
				message =msg;
				notification.removeClass("alert-success").addClass("alert-danger");
			}

			notification.html(message).show();
			setTimeout(function () {
				notification.fadeOut(function () {
					if (callback) callback();
				});
			}, 3000);
		}

		function removeQueryParams() {
			const newUrl = window.location.protocol + "//" + window.location.host + window.location.pathname;
			window.history.replaceState({}, document.title, newUrl);
		}

		const urlParams = new URLSearchParams(window.location.search);
		const status = urlParams.get("status");
		const msg = urlParams.get("msg");
		
		if (status) {
			showNotification(status, msg, function () {
				removeQueryParams();
			});
		}
		// Notification End

		// Form Validation
		const form = document.getElementById("AddProdProcess");
		const addButton = document.getElementById("SaveProcesses");
		const requiredFields = form.querySelectorAll("[data-validate='required']");

		function checkFormValidity() {
			let isValid = true;
			requiredFields.forEach(function (field) {
				if (!field.value.trim()) {
					isValid = false;
				}
			});
			addButton.disabled = !isValid;
		}

		requiredFields.forEach(function (field) {
			field.addEventListener("input", checkFormValidity);
		});
		checkFormValidity();

		// Form Submission
		form.addEventListener("submit", function (event) {
			event.preventDefault();

			const payload = new FormData();

			const iconFile = $("#addProcessIcon")[0].files[0];
			if (iconFile) {
				payload.append("addProcessIcon", iconFile); // matches Go's r.FormFile key
			}

			payload.append("Name", document.getElementById("addProcessName").value.trim());
			payload.append("Description", document.getElementById("processDescription").value.trim());
			payload.append("expected_mean_time", document.getElementById("addProcessTime").value.trim());
			payload.append("Status", document.getElementById("processIsActive").value);
			payload.append("line_visibility", document.getElementById("processIsVisibalse").value.trim() );

			$.ajax({
				url: "/add-production-process",
				type: "POST",
				data: payload,
				processData: false,
				contentType: false,
				success: function (response) {
					window.location.href = "/prod-processes-management?status=200"+ "&msg=" + response;
				},
				error: function (xhr) {
					window.location.href = "/prod-processes-management?status=" + xhr.status + "&msg=" + xhr.responseText;
				}
			});
		});

		// Handle "ViewProdProcessesDetails" click event
		document.addEventListener("click", function (event) {
			const viewButton = event.target.closest("#ViewProdProcessesDetails");
			if (viewButton) {
				const rowElement = viewButton.closest("tr");
				if (rowElement && rowElement.dataset.data) {
					const rowData = JSON.parse(rowElement.dataset.data);
					const id = rowData.ID;

					const xhr = new XMLHttpRequest();
					xhr.open("POST", "/edit-prod-process-dialog");
					xhr.setRequestHeader("Content-Type", "application/json");
					xhr.onload = function () {
						if (xhr.status === 200) {
							const data = JSON.parse(xhr.responseText);
							if (data.dialogHTML) {
								const existingModal = document.getElementById("EditProductionProcess");
								if (existingModal) {
									existingModal.remove();
								}

								document.body.insertAdjacentHTML("beforeend", data.dialogHTML);

								const modal = new bootstrap.Modal(document.getElementById("EditProductionProcess"));
								modal.show();

								// Dialog close Code
								const closeDialogButton = document.getElementById("closeDialog");
								if (closeDialogButton) {
									closeDialogButton.addEventListener("click", function (event) {
										event.preventDefault();
										const modalElement = closeDialogButton.closest('.modal');
										if (modalElement) {
											const modalInstance = bootstrap.Modal.getInstance(modalElement);
											if (modalInstance) {
												modalInstance.hide(); // Hide the modal
											}
											modalElement.parentNode.removeChild(modalElement);
											const modalBackdrop = document.querySelector('.modal-backdrop');
											if (modalBackdrop) {
												modalBackdrop.parentNode.removeChild(modalBackdrop);
											}
											document.body.classList.remove('modal-open');
											document.body.style.paddingRight = '';
										}
									});
								}

								const updateLineButton = document.getElementById("UpdateProcesses");
								if (updateLineButton) {
									updateLineButton.addEventListener("click", function (event) {
										event.preventDefault();
										updateProcess();
									});
								}
							} else {
								console.error("Dialog HTML not received");
							}
						}
					};
					xhr.onerror = function () {
						console.error("Error fetching dialog box");
					};
					xhr.send(JSON.stringify({ ID: id }));
				} else {
					console.error("Row data not found or invalid");
				}
			}
		});
		document.addEventListener("click", function (event) {
			// Check if the clicked element is a file input
			if (event.target && event.target.id === "editProcessIcon") {
				event.target.onchange = function () {
					const file = this.files[0];
					if (file) {
						document.getElementById("editProcessIcon-value").value = file.name;
					}
				};
			}
		});

		// Function to handle process update
		function updateProcess() {

			const payload = new FormData();

			const iconFile = $("#editProcessIcon")[0].files[0];
			if (iconFile) {
				payload.append("editProcessIcon", iconFile); // matches Go's r.FormFile key
			}

			payload.append("ID", document.getElementById("editProcessID").value.trim());
			payload.append("Name", document.getElementById("editProcessName").value.trim());
			payload.append("Description", document.getElementById("editprocessDescription").value.trim());
			payload.append("expected_mean_time", document.getElementById("editProcessTime").value.trim());
			payload.append("Status", document.getElementById("editisactive").value);
			payload.append("line_visibility", document.getElementById("editprocessIsVisibalse").value.trim() );
			payload.append("editProcessIconView", document.getElementById("editProcessIcon-value").value.trim() );
			$.ajax({
				url: "/edit-prod-process",
				type: "POST",
				data: payload,
				processData: false,
				contentType: false,
				success: function (response) {
					window.location.href = "/prod-processes-management?status=200"+ "&msg=" + response;
				},
				error: function (xhr) {
					window.location.href = "/prod-processes-management?status=" + xhr.status + "&msg=" + xhr.responseText;
				}
			});
		}
	});
	//!js
	</script>`

	html.WriteString(ProdProcessesTable.Build())
	html.WriteString(js)
	html.WriteString(s.AddProdProcessesModal("AddNewProcesses", "", "Add Process"))
	return html.String()
}

func AddProdProcess(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form with a max memory of 10 MB
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Can't parse multipart form", http.StatusBadRequest)
		return
	}

	userID := r.Header.Get("X-Custom-Userid")
	var ProdProcess m.ProdProcess

	ProdProcess.Name = strings.TrimSpace(r.FormValue("Name"))
	ProdProcess.ExpectedMeanTime = strings.TrimSpace(r.FormValue("expected_mean_time"))
	ProdProcess.Description = strings.TrimSpace(r.FormValue("Description"))
	ProdProcess.Status = strings.TrimSpace(r.FormValue("Status"))

	// Get uploaded file
	file, fileHeader, err := r.FormFile("addProcessIcon")
	if err != nil {
		http.Error(w, "Missing icon file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Ensure unique folder
	folder := ProdProcess.Name
	targetDir := filepath.Join(uploadDir, folder)
	for i := 1; ; i++ {
		if _, err := os.Stat(targetDir); os.IsNotExist(err) {
			break
		}
		folder = fmt.Sprintf("%s_%d", ProdProcess.Name, i)
		targetDir = filepath.Join(uploadDir, folder)
	}

	// Create folder
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		http.Error(w, "Server error: cannot create folder", http.StatusInternalServerError)
		return
	}

	// Save file
	savePath := filepath.Join(targetDir, fileHeader.Filename)
	out, err := os.Create(savePath)
	if err != nil {
		http.Error(w, "Server error: cannot save image", http.StatusInternalServerError)
		return
	}
	defer out.Close()
	if _, err := io.Copy(out, file); err != nil {
		http.Error(w, "Server error: failed to write image", http.StatusInternalServerError)
		return
	}
	ProdProcess.Icon = filepath.ToSlash(filepath.Join("uploaded-icons/uploadedIcons", folder))
	if strings.TrimSpace(r.FormValue("line_visibility")) == "true" {
		ProdProcess.LineVisibility = true
	} else {
		ProdProcess.LineVisibility = false
	}

	ProdProcess.CreatedBy = userID

	jsonValue, err := json.Marshal(ProdProcess)
	if err != nil {
		http.Error(w, "Failed to marshal the data", http.StatusBadRequest)
		log.Println("Failed to marshal the data:", err)
		return
	}

	resp, err := http.Post(utils.RestURL+"/add-production-process", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		slog.Error("%s - error - %s", "Error making GET request", err)
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

func EditProdProcess(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form with a max memory of 10 MB
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Can't parse multipart form", http.StatusBadRequest)
		return
	}
	userID := r.Header.Get("X-Custom-Userid")

	// Decode the request body into ProdLine struct
	var prodProcessData m.ProdProcess
	prodProcessData.Id, _ = strconv.Atoi(strings.TrimSpace(r.FormValue("ID")))
	prodProcessData.Name = strings.TrimSpace(r.FormValue("Name"))
	prodProcessData.ExpectedMeanTime = strings.TrimSpace(r.FormValue("expected_mean_time"))
	prodProcessData.Description = strings.TrimSpace(r.FormValue("Description"))
	prodProcessData.Status = strings.TrimSpace(r.FormValue("Status"))

	// Get uploaded file
	file, fileHeader, err := r.FormFile("editProcessIcon")
	if err != nil {
		file = nil
		fileHeader = nil
	} else {
		defer file.Close()
	}
	if file != nil {
		// Ensure unique folder
		folder := prodProcessData.Name
		targetDir := filepath.Join(uploadDir, folder)
		for i := 1; ; i++ {
			if _, err := os.Stat(targetDir); os.IsNotExist(err) {
				break
			}
			folder = fmt.Sprintf("%s_%d", prodProcessData.Name, i)
			targetDir = filepath.Join(uploadDir, folder)
		}

		// Create folder
		if err := os.MkdirAll(targetDir, 0755); err != nil {
			http.Error(w, "Server error: cannot create folder", http.StatusInternalServerError)
			return
		}

		// Save file
		savePath := filepath.Join(targetDir, fileHeader.Filename)
		out, err := os.Create(savePath)
		if err != nil {
			http.Error(w, "Server error: cannot save image", http.StatusInternalServerError)
			return
		}
		defer out.Close()
		if _, err := io.Copy(out, file); err != nil {
			http.Error(w, "Server error: failed to write image", http.StatusInternalServerError)
			return
		}
		prodProcessData.Icon = filepath.ToSlash(filepath.Join("uploaded-icons/uploadedIcons", folder))
	} else {
		prodProcessData.Icon = strings.TrimSpace(r.FormValue("editProcessIconView"))
	}

	if strings.TrimSpace(r.FormValue("line_visibility")) == "true" {
		prodProcessData.LineVisibility = true
	} else {
		prodProcessData.LineVisibility = false
	}

	prodProcessData.CreatedBy = userID

	jsonValue, err := json.Marshal(prodProcessData)
	if err != nil {
		http.Error(w, "Failed to marshal the data", http.StatusBadRequest)
		log.Println("Failed to marshal the data:", err)
		return
	}

	url := utils.RestURL + "/edit-prod-process"
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonValue))
	if err != nil {
		http.Error(w, "Error creating request to REST service", http.StatusInternalServerError)
		return
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to send request to the target service", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("%s - error - %s", "Error reading response body", err)
		utils.SetResponse(w, http.StatusInternalServerError, "Failed to read response body")
		return
	}

	utils.SetResponse(w, http.StatusOK, string(responseBody))

}
