package productionline

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

	"irpl.com/kanban-commons/model"

	m "irpl.com/kanban-commons/model"
	"irpl.com/kanban-commons/utils"

	"irpl.com/kanban-web/pages/basepage"
	"irpl.com/kanban-web/services"
	s "irpl.com/kanban-web/services"
)

type ProdLineManagement struct {
	Username string
	UserType string
	UserID   string
}

func ProdLineCRUDPage(w http.ResponseWriter, r *http.Request) {
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
				Selected: true,
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

	ProdLineManagement := ProdLineManagement{
		UserType: usertype,
		UserID:   userID,
	}

	// Build the complete BasePage with navigation and content
	page := (&basepage.BasePage{
		SideNavBar: sideNav,
		TopNavBar:  topNav,
		Username:   username,
		UserType:   usertype,
		Content:    ProdLineManagement.Build(),
	}).Build()

	// Write the complete HTML page to the response
	w.Write([]byte(page))
}

func (c *ProdLineManagement) Build() string {

	var Response struct {
		Pagination m.PaginationResp
		Data       []*m.ProdLine
	}
	type TableConditions struct {
		Pagination m.PaginationReq `json:"pagination"`
		Conditions []string        `json:"conditions"`
	}

	var tablecondition TableConditions

	// Temporary pagination data
	tablecondition.Pagination = m.PaginationReq{
		Type:   "ProdLine",
		Limit:  "10",
		PageNo: 1,
	}

	// Convert to JSON
	jsonValue, err := json.Marshal(tablecondition)
	if err != nil {
		slog.Error("Error marshaling table condition", "error", err)
	}

	resp, err := http.Post(utils.RestURL+"/get-all-prod-line-data-by-search-paginations", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		slog.Error("%s - error - %s", "Error making GET request", err)
	}

	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&Response); err != nil {
		slog.Error("error decoding response body", "error", err)
	}

	// Fetch operator list
	var operators []m.Operator
	respOp, err := http.Get(utils.RestURL + "/get-all-operator")
	if err != nil {
		slog.Error("Failed to fetch operators", "error", err)
	} else {
		defer respOp.Body.Close()
		if err := json.NewDecoder(respOp.Body).Decode(&operators); err != nil {
			slog.Error("Error decoding operators", "error", err)
		}
	}

	// Build lookup map
	opMap := make(map[string]string)
	for _, op := range operators {
		opMap[op.OperatorCode] = op.OperatorName
	}

	for _, line := range Response.Data {
		if line.Operator != "" {
			if strings.Contains(line.Operator, "(") && strings.Contains(line.Operator, ")") {
				open := strings.Index(line.Operator, "(")
				close := strings.Index(line.Operator, ")")
				line.OperatorDisplay = strings.TrimSpace(line.Operator[:open])
				line.OperatorCode = line.Operator[open+1 : close]
			} else {
				line.OperatorDisplay = opMap[line.Operator]
				line.OperatorCode = line.Operator
			}
		} else {
			line.OperatorDisplay = ""
			line.OperatorCode = ""
		}

	}

	var html strings.Builder

	tablebutton := `
	<!--html-->
			<button type="button" class="btn m-0 p-0" id="ViewProdLineDetails" data-toggle="tooltip" data-placement="bottom" title="View Details"> 
				<i class="fa fa-edit mx-2" style="color: #b250ad;"></i> 
			</button>
			<!--!html-->`

	html.WriteString(`
		<!--html-->
		<!--!html-->`)

	var ProdLineTable s.TableCard

	ProdLineTable.CardHeading = "Machine Master"
	ProdLineTable.CardHeadingActions.Component = []s.CardHeadActionComponent{{ComponentName: "modelbutton", ComponentType: s.ActionComponentElement{ModelButton: s.ModelButtonAttributes{ID: "add-new-prod-line", Name: "add-new-compound", Type: "button", Text: "Add New Line", ModelID: "#AddLine"}}, Width: "col-3"}}
	ProdLineTable.BodyTables = s.CardTableBody{Columns: []s.CardTableBodyHeadCol{
		{
			Name:       `Sr. No.`,
			IsCheckbox: true,
			ID:         "search-sn",
			Width:      "col-1",
		},
		{
			Lable:        "Line Name",
			Name:         "Name",
			ID:           "prod-line-name",
			Width:        "col-2",
			Type:         "input",
			DataField:    "name",
			IsSearchable: true,
		},
		{
			Lable:        "Operator Name",
			Name:         "OperatorDisplay",
			ID:           "op-name",
			Width:        "col-2",
			Type:         "input",
			DataField:    "operator",
			IsSearchable: true,
		},
		{
			Lable:        "Operator Code",
			Name:         "OperatorCode",
			ID:           "op-code",
			Width:        "col-2",
			Type:         "input",
			DataField:    "operator",
			IsSearchable: true,
		},

		{
			Lable: "Status",
			Name:  "Status",
			ID:    "status",
			Width: "col-2",
		},
		{
			Name:  "Action",
			Width: "col-1",
		},
	},
		ColumnsWidth: []string{"col-1", "col-2", "col-2", "col-2", "col-2", "col-1"},

		Data:    Response.Data,
		Buttons: tablebutton,
	}
	var Pagination s.Pagination
	Pagination.CurrentPage = 1
	Pagination.Offset = 0
	Pagination.TotalRecords = Response.Pagination.TotalNo
	Pagination.PerPage, _ = strconv.Atoi(tablecondition.Pagination.Limit)
	Pagination.PerPageOptions = []int{10, 25, 50, 100, 200}
	ProdLineTable.CardFooter = Pagination.Build()

	js := `
		<script>
		//js
		$(function(){
			const urlParams = new URLSearchParams(window.location.search);
			const status = urlParams.get('status');
			const msg = urlParams.get('msg');
			
			if (status) {
				showNotification(status,msg, () => {
					removeQueryParams();
				});
			}
		});
		function showNotification(status, msg, callback) {
			const notification = $('#notification');
			var message = '';
			if (status === "200") {
				message = '<strong>Success!</strong> ' + msg + '.';
				notification.removeClass("alert-danger").addClass("alert-success");
			} else {
				message = '<strong>Fail!</strong> ' + msg + '.';
				notification.removeClass("alert-success").addClass("alert-danger");
			}
			notification.html(message);
			notification.show();

			setTimeout(() => {
				notification.fadeOut(() => {
					if (callback) callback();
				});
			}, 5000);
		}

		function removeQueryParams() {
			var newUrl = window.location.protocol + "//" + window.location.host + window.location.pathname;
			window.history.replaceState({}, document.title, newUrl);
		}
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

			fetch("/prodline-search-pagination", {
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

		document.addEventListener('DOMContentLoaded', () => {

		document.querySelectorAll("[id^='search-by']").forEach(function (input) {
			input.addEventListener("input", function () {
				pagination(1);
			});
		});
		
		attachPerPageListener();


			let selectedProcesses = [];
			
			const checkboxes = document.querySelectorAll('.form-check-input');
			const clearBtn = document.getElementById("clear_check_box");
			const saveButton = document.getElementById("SaveLine");
			
			function validateForm() {
				const lineName = document.getElementById("addLineName").value.trim();
				const lineDesc = document.getElementById("addLineDescription").value.trim();
				const isValid = lineName && lineDesc && selectedProcesses.length > 0;
				saveButton.disabled = !isValid;
			}

			function updateClearButtonState() {
				clearBtn.disabled = selectedProcesses.length === 0;
			}

			function updateInputValues() {
				selectedProcesses.forEach((processId, index) => {
					const inputField = $('#' + processId + '.processNo');
					inputField.val(index + 1); 
				});
			}

			checkboxes.forEach((checkbox) => {
				checkbox.addEventListener('change', (event) => {
					const checkboxId = event.target.id;
					const inputField = $('#' + checkboxId + '.processNo');
					const groupNameInput = $('#' + checkboxId + '.GroupName');

					if (event.target.checked) {
						selectedProcesses.push(checkboxId);
						groupNameInput.prop('disabled', false);
						inputField.val(selectedProcesses.length);
					} else {
						const index = selectedProcesses.indexOf(checkboxId);
						if (index > -1) {
							selectedProcesses.splice(index, 1);
						}
						inputField.val("");
						groupNameInput.val("");
						groupNameInput.prop('disabled', true);
						updateInputValues();
					}

					validateForm();
					updateClearButtonState();
				});
			});

			document.getElementById("addLineName").addEventListener('input', validateForm);
			document.getElementById("addLineDescription").addEventListener('input', validateForm);

			clearBtn.addEventListener('click', () => {
				checkboxes.forEach((checkbox) => {
					checkbox.checked = false;
					$('#' + checkbox.id + '.processNo').val("");
					$('#' + checkbox.id + '.GroupName').val("").prop('disabled', true);
				});

				selectedProcesses = [];
				validateForm();
				updateClearButtonState();
			});

			document.getElementById("SaveLine").addEventListener("click", (event) => {
				event.preventDefault();

				let LineName = document.getElementById("addLineName").value;
				let LineDesc = document.getElementById("addLineDescription").value;

				const processOrders = selectedProcesses.map((processId, index) => {
					const groupName = $('#' + processId + '.GroupName').val();
					return {
						prod_process_id: processId,
						order: index + 1,
						group_name: groupName
					};
				});

				const payload = {
					line_name: LineName,
					line_description: LineDesc,
					process_orders: processOrders
				};
				$.post("/add-production-line", JSON.stringify(payload), function(response) {

					if (response.ErrCode === 200) {
						window.location.href = "/prod-line-management?status=" + response.ErrCode + "&msg=" + encodeURIComponent(response.ErrMessage);
					}
				})
			});

			// Initialize form validation and button states
			validateForm();
			updateClearButtonState();
		});

		function CheckDropDown() {  
			$("#processIsActive").change(function () {  
				if ($(this).val() === "false") {  
					console.log("We selected false");  
					$(".d-none").removeClass("d-none"); // Make hidden fields visible
				} else {  
					$("#moveToLineDropDown").closest(".w-100").addClass("d-none"); 
					$("#moveToLineInfo").closest(".w-100").addClass("d-none"); 

				}  
			});  
		}  


		document.addEventListener('DOMContentLoaded', () => {
			document.addEventListener("click", function (event) {	
				if (event.target.closest("#ViewProdLineDetails")) {
					const rowElement = event.target.closest("tr");
					if (rowElement && rowElement.dataset.data) {
						const rowData = JSON.parse(rowElement.dataset.data);
						const id = rowData.ID; 

						fetch("/edit-prod-line-dialog", {
							method: "POST",
							headers: {
								"Content-Type": "application/json",
							},
							body: JSON.stringify({ ID: id }),
						})
							.then((response) => response.json())
							.then((data) => {
								if (data.dialogHTML) {
									const existingModal = document.getElementById("EditLine");
									if (existingModal) {
										existingModal.remove();
									}

									document.body.insertAdjacentHTML("beforeend", data.dialogHTML);

									const modal = new bootstrap.Modal(document.getElementById("EditLine"));
									modal.show();
									CheckDropDown();
									// Dialoge close Code
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


									const updateLineButton = document.getElementById("updateLine");
									if (updateLineButton) {
										updateLineButton.addEventListener("click", function (event) {
											event.preventDefault(); // Prevent the default form submission

											const lineID = parseInt(document.getElementById("editLineID").value,10);
											const lineName = document.getElementById("editLineName").value;
											const lineDescription = document.getElementById("editLineDescription").value;
										
											const lineStatus = document.getElementById("processIsActive").value === "true"; // Convert text to boolean
											const moveToLineElement = document.getElementById("moveToLineDropDown");
											const MoveToLineID = moveToLineElement ? parseInt(moveToLineElement.value, 10) || null : null;
											const operatorCode = document.getElementById("EditOperator").value;

											const data = {
												id: lineID,
												name: lineName,
												description: lineDescription,
												status: lineStatus,
												MoveToLineID: MoveToLineID,
												operator: operatorCode,
											};

											$.post("/edit-prod-line", JSON.stringify(data), function(response) {

												if (response.ErrCode === 200) {
													window.location.href = "/prod-line-management?status=" + response.ErrCode + "&msg=" + encodeURIComponent(response.ErrMessage);
												}
											})
										});
									}							} else {
									console.error("Dialog HTML not received");
								}
							})
							.catch((error) => {
								console.error("Error fetching dialog box:", error);
							});
					} else {
						console.error("Row data not found or invalid");
					}
				}
			}); 
		});
		//!js
		</script>
	`

	html.WriteString(ProdLineTable.Build())
	html.WriteString(js)
	html.WriteString(s.AddProdLineModal())
	return html.String()
}

func EditProdLineDialog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		log.Println("Failed to read request body:", err)
		return
	}

	// Decode the request body into ProdLine struct
	var prodLineData m.ProdLine
	err = json.Unmarshal(body, &prodLineData)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.Println("Failed to decode request:", err)
		return
	}

	data, err := json.Marshal(prodLineData)
	if err != nil {
		http.Error(w, "Unable to marshal data", http.StatusBadRequest)
		log.Println("Failed to send request:", err)
		return
	}

	url := utils.RestURL + "/get-prod-line-by-param"
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
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
	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Failed to process request on target service", http.StatusInternalServerError)
		return
	}
	var ProdLine []m.ProdLine
	err = json.NewDecoder(resp.Body).Decode(&ProdLine)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.Println("--Failed to decode request:", err)
		return
	}

	var ProdProcessData *m.ProdLineDetails
	respTemp, err := http.Post(utils.RestURL+"/get-lineup-processes-by-lineid", "application/json", bytes.NewBuffer(data))
	if err != nil {
		slog.Error("%s - error - %s", "Error making GET request", err)
	}

	defer resp.Body.Close()

	if err := json.NewDecoder(respTemp.Body).Decode(&ProdProcessData); err != nil {
		slog.Error("error decoding response body", "error", err)
	}
	var operators []m.Operator
	var operatoroptions []s.DropDownOptions
	operatorresp, err := http.Get(utils.RestURL + "/get-all-operator")
	if err != nil {
		slog.Error("%s - error - %s", "Error making GET request", err)
	}

	defer operatorresp.Body.Close()

	err = json.NewDecoder(operatorresp.Body).Decode(&operators)
	if err != nil {
		slog.Error("%s - error - %s", "Error decoding request body", err)
	}

	for _, v := range operators {
		//	display := v.OperatorName + " (" + v.OperatorCode + ")"
		operatoroptions = append(operatoroptions, s.DropDownOptions{
			Text:  fmt.Sprintf("%s (%s)", v.OperatorName, v.OperatorCode), // Displayed to the user
			Value: fmt.Sprintf("%s (%s)", v.OperatorName, v.OperatorCode), // ✅ Submitted to backend
		})

	}

	for i, v := range operatoroptions {
		if v.Value == ProdLine[0].Operator {
			operatoroptions[i].Selected = true
		}
	}

	var ProdLineStatus = []s.DropDownOptions{{Text: "False", Value: "false"}, {Text: "True", Value: "true"}}
	Status := ProdLine[0].Status
	if Status {
		ProdLineStatus = []s.DropDownOptions{{Text: "False", Value: "false"}, {Text: "True", Value: "true", Selected: true}}
	} else {
		ProdLineStatus = []s.DropDownOptions{{Text: "False", Value: "false", Selected: true}, {Text: "True", Value: "true"}}

	}

	var recipe []*m.Recipe

	var rawQuery m.RawQuery
	rawQuery.Host = utils.RestHost
	rawQuery.Port = utils.RestPort
	rawQuery.Type = "Recipe"
	rawQuery.Query = `SELECT * FROM recipe  ;` //`;`
	rawQuery.RawQry(&recipe)

	var recipelist []services.DropDownOptions

	for _, v := range recipe {
		var temp services.DropDownOptions
		temp.Text = v.CompoundName + " (" + v.CompoundCode + ")"
		temp.Value = strconv.Itoa(v.Id)
		recipelist = append(recipelist, temp)
	}

	for i, v := range recipelist {
		if v.Value == strconv.Itoa(ProdLine[0].RecipeId) {
			recipelist[i].Selected = true
		}
	}

	var TempDropDown s.ModelCard // Declare outside to avoid scope issues
	var TempInfo s.InfoLine
	if len(ProdProcessData.Cells) != 0 {

		TempInfo = s.InfoLine{
			IsVisible: false,
			ID:        "moveToLineInfo",
			Lable:     "To disable a production line, you must first reassign the lined-up Kanban to an alternative production line.",
		}

		var ProdLineData []*m.ProdLine
		resp, err = http.Get(utils.RestURL + "/get-all-prod-line-data")
		if err != nil {
			slog.Error("%s - error - %s", "Error making GET request", err)
		}

		defer resp.Body.Close()

		if err := json.NewDecoder(resp.Body).Decode(&ProdLineData); err != nil {
			slog.Error("error decoding response body", "error", err)
		}

		var ProdLines []s.DropDownOptions
		for _, data := range ProdLineData {
			if data.Status && data.Id != prodLineData.Id {
				ProdLines = append(ProdLines, s.DropDownOptions{
					Text:  data.Name,
					Value: strconv.Itoa(data.Id),
				})
			}
		}

		TempDropDown = s.ModelCard{
			Form: s.ModelForm{
				Dropdownfield: []s.DropdownAttributes{
					{
						Label:    "Move Kanban to",
						DataType: "text",
						Name:     "moveToLineDropDown",
						Options:  ProdLines,
						ID:       "moveToLineDropDown",
						Width:    "w-100",
						Hidden:   true, // Ensures it is initially hidden
					},
				},
			},
		}
	}

	var defaultline bool
	if strings.ToLower(ProdLine[0].Name) == "quality" || strings.ToLower(ProdLine[0].Name) == "packing" {
		defaultline = true
	} else {
		defaultline = false
	}
	// Now correctly add `TempDropDown` to the form
	EditLineModel := s.ModelCard{
		ID:      "EditLine",
		Type:    "",
		Heading: "Edit Line",
		Form: s.ModelForm{
			FormID:     "EditProdLine",
			FormAction: "",
			Footer: s.Footer{
				CancelBtn: false,
				Buttons: []s.FooterButtons{
					{BtnType: "submit", BtnID: "closeDialog", Text: "Close", Disabled: false},
					{BtnType: "submit", BtnID: "updateLine", Text: "Update", Disabled: false},
				},
			},
			Inputfield: []s.InputAttributes{
				{Type: "text", Name: "editLineID", ID: "editLineID", Label: `Line Name`, Width: "w-100", Required: true, Value: strconv.Itoa(ProdLine[0].Id), Hidden: true, Readonly: true},
				{Type: "text", Name: "editLineName", ID: "editLineName", Label: `Line Name`, Width: "w-100", Required: true, Value: ProdLine[0].Name, Disabled: defaultline, Readonly: defaultline},
			},
			TextArea: []s.TextAreaAttributes{
				{Label: "Line Description", DataType: "text", Name: "editLineDescription", ID: "editLineDescription", Value: ProdLine[0].Description},
			},
			Dropdownfield: []s.DropdownAttributes{
				{Label: "Operator", DataType: "text", Name: "operator", Options: operatoroptions, ID: "EditOperator", Width: "w-100"},
				{Label: "Line Status", DataType: "text", Name: "isactive", Options: ProdLineStatus, ID: "processIsActive", Width: "w-100", Disabled: defaultline},
				// {Label: "Recipe", DataType: "text", Name: "editRecipt", Options: recipelist, ID: "editRecipt", Width: "w-100", Disabled: defaultline},
			},
			Info: TempInfo,
			ExtractFields: []s.ExtractFields{
				{HTML: TempDropDown.Form.DropdownRender()},
			},
		},
	}

	DialogBox := EditLineModel.Build()

	response := map[string]string{
		"dialogHTML": DialogBox,
	}
	json.NewEncoder(w).Encode(response)

}

func EditProdLine(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-Custom-Userid")

	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		log.Println("Failed to read request body:", err)
		return
	}

	// Decode the request body into ProdLine struct
	var prodLineData m.ProdLine
	err = json.Unmarshal(body, &prodLineData)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.Println("Failed to decode request:", err)
		return
	}

	prodLineData.ModifiedBy = userID

	data, err := json.Marshal(prodLineData)
	if err != nil {
		http.Error(w, "Unable to marshal data", http.StatusBadRequest)
		log.Println("Failed to send request:", err)
		return
	}

	url := utils.RestURL + "/edit-prod-line"
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		http.Error(w, "Error creating request to REST service", http.StatusInternalServerError)
		return
	}
	client := &http.Client{}
	APIresp := model.ErrorResp{}
	resp, err := client.Do(req)
	if err != nil {
		APIresp.ErrCode = http.StatusInternalServerError
		APIresp.ErrMessage = "Failed to send request to REST service"
		body, jsonErr := json.Marshal(APIresp)
		if jsonErr != nil {
			http.Error(w, "{err_message:Internal Server Error}", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(body)
		return
	}
	defer resp.Body.Close()

	APIresp.ErrCode = http.StatusOK
	APIresp.ErrMessage = "Production Line Updated!"
	body, jsonErr := json.Marshal(APIresp)
	if jsonErr != nil {
		http.Error(w, "{err_message:Internal Server Error}", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(body)

}
