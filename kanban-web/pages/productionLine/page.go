package productionline

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"strings"

	m "irpl.com/kanban-commons/model"
	"irpl.com/kanban-commons/utils"
	"irpl.com/kanban-web/pages/basepage"
	"irpl.com/kanban-web/services"
)

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

	RestURL = utils.JoinStr("http://", RestHost, ":", RestPort)
}
func RenderProductionLinePage(w http.ResponseWriter, r *http.Request) {
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

	productionLineService := services.NewProductionLineService()
	lines, err := productionLineService.FetchProductionLines()
	if err != nil {
		http.Error(w, "Failed to load production lines", http.StatusInternalServerError)
		return
	}
	var contentHTML string
	for _, line := range lines {
		// skip build for quality and packing lines
		if strings.ToLower(line.ProdLineName) == "quality" || strings.ToLower(line.ProdLineName) == "packing" {
			continue
		}
		contentHTML += line.Build()
	}
	contentHTML +=
		`<!--html-->
		<script>
    		document.addEventListener("DOMContentLoaded", () => {
      			document.getElementById('AddProdLine').addEventListener('submit', function(event) {
        			event.preventDefault(); 
        			this.submit();
      			});
    		});
  		</script>
		<!--!html-->
		` + services.AddProdLineModal() + `
		<!--html-->
		<script>
			// JS to update running number
			$(document).ready(function () {
				$(document).on('dragstart', '.production-line-card .kanban-item', function () {
					$(this).data('initial-position', $(this).index() + 1);
				});
				$(document).on('dragend', '.production-line-card .kanban-item', function () {
					let dataToSend = [];
					$(this).closest('.production-line-card').find('.kanban-item').each(function (index) {
						const krId = parseInt($(this).find('.KRid').val(), 10);
						const finalPosition = index + 1;
						const initialNo = parseInt($(this).find('.cell-info').attr('temp_no'), 10);

						if (initialNo !== finalPosition) {
							dataToSend.push({
								ID: krId,
								RunningNo: finalPosition
							});

							$(this).find('.cell-info').attr('temp_no', finalPosition);
						}
					});

					if (dataToSend.length > 0) {
						fetch('/update-running-numbers', {
							method: 'POST', // HTTP method
							headers: {
								'Content-Type': 'application/json' // Specify content type as JSON
							},
							body: JSON.stringify(dataToSend) // Convert the data to JSON string
						})
					}
				});

				$(document).on('click', '#delete-kanban-btn', function () {
					let selectedCheckboxes = document.querySelectorAll(".form-check-input:checked"); 
					let selectedIds = [];

					selectedCheckboxes.forEach(checkbox => {
						let hiddenInput = checkbox.closest(".kanban-item").querySelector(".KRid");
						if (hiddenInput) {
							selectedIds.push(hiddenInput.value);
						}
					});	
					const requestData = JSON.stringify({ KRid: selectedIds });
					// Send the DELETE request
					fetch('/delete-production-line-cell', {
						method: 'DELETE',
						headers: {
							'Content-Type': 'application/json',
						},
						body: requestData,
					})
					.then((response) => {
						if (!response.ok) {
							throw new Error('Failed to delete kanban');
						}
						location.reload();
					})
					.catch((error) => {
						alert("Error: Unable to delete the production line cell");
					});
				});

				$(document).on('click', '.expand', function () {
					const button = $(this);
					const collapseElement = button.closest('.kanban-item').find('.collapse');
					const iconElement = button.find('svg');
					button.prop('disabled', true); // Disable button during animation
					collapseElement
						.css('font-size', '12px')
						.find('.row')
						.each(function (index) {
							$(this).css('background-color', index % 2 === 0 ? '#e0f7ff' : '#ffffff');
						});
					collapseElement.collapse('toggle');
					collapseElement.on('shown.bs.collapse hidden.bs.collapse', function () {
						button.prop('disabled', false); // Enable button after animation
					});
					if (iconElement.hasClass('rotated')) {
						iconElement.removeClass('rotated');
						iconElement.css({
							animation: 'rotateTo0 0.2s ease-out forwards',
						});
					} else {
						iconElement.addClass('rotated');
						iconElement.css({
							animation: 'rotateTo90 0.2s ease-in forwards',
						});
					}
				});
			});
			document.addEventListener('DOMContentLoaded', () => {
				const wrapper = document.querySelector('.wrapper-production-line');
				const scrollableItems = wrapper.querySelectorAll('.kanban-items-container');

				// Handle horizontal scroll for wrapper
				wrapper.addEventListener('wheel', (event) => {
					if (event.target.closest('.kanban-items-container')) return; // Skip if over a vertical scrollable item

					event.preventDefault(); // Prevent vertical scrolling
					wrapper.scrollBy({
						left: event.deltaY * 3, // Horizontal scroll
						behavior: 'smooth',
					});
				});

				// Handle vertical scroll for inner elements
				scrollableItems.forEach((item) => {
					item.addEventListener('wheel', (event) => {
						const canScrollVertically =
							item.scrollHeight > item.clientHeight;

						if (canScrollVertically) {
							event.stopPropagation(); // Prevent scroll propagation to the wrapper
						}
					});
				});
			});
		</script>

		<!-- This is js for Dialog box -->
		<script>
		//js
		document.addEventListener('DOMContentLoaded', () => {
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
				let RecipeId =  parseInt(document.getElementById("addrecipe").value,10);
				
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
					process_orders: processOrders,
					recipe_id:RecipeId
				};

				fetch('/add-production-line', {
					method: 'POST',
					headers: { 'Content-Type': 'application/json' },
					body: JSON.stringify(payload)
				})
				.then(response => {
					if (response.status === 200) {
						// location.reload();
					} else {
						console.log("Failed to create the production line. Status:", response.status);
					}
					return response.json();
				});
			});

			// Initialize form validation and button states
			validateForm();
			updateClearButtonState();
		});

		//!js
		</script>
		<script>
			document.addEventListener('DOMContentLoaded', () => {
				const scrollDiv = document.querySelector('.scroll-div');

				scrollDiv.addEventListener('wheel', (event) => {
					event.preventDefault();

					// Smooth scroll by using scrollBy with smooth behavior
					scrollDiv.scrollBy({
						top: event.deltaY * 3, // Adjust scroll speed
						behavior: 'smooth'     // Enable smooth scrolling
					});
				});
			});
		</script>
	
		<!--!html-->`

	css := `
		<style>
		/*css*/
		@keyframes rotateTo90 {
			from {
				transform: rotate(0deg);
			}
			to {
				transform: rotate(90deg);
			}
		}
		@keyframes rotateTo0 {
			from {
				transform: rotate(90deg);
			}
			to {
				transform: rotate(0deg);
			}
		}
			.add-line-btn-container {
				margin-top: 0; 
				padding-bottom: 10px;
			}

			.add-line-btn {
				background-color: #871a83;
				color: white;
				border: none;
				border-radius: 5px;
				padding: 10px;
				cursor: pointer;
			}
			#delete-kanban-btn {
				background-color: #E63757;
				color: white;
				border: none;
				border-radius: 5px;
				padding: 10px;
				cursor: pointer;
			}
			
			#delete-kanban-btn i {
				color: white !important;  
			}
			#delete-kanban-btn svg {
				fill: white !important; 
				color:#ffffff;
			}
			.wrapper-production-line {
				overflow-y: hidden;
				box-sizing: border-box;
				display: flex;
				align-items: flex-start;
				width: 100%;
				height:85vh !important;
				padding-bottom : 10px;
				overflow-x: auto;
				&::-webkit-scrollbar-thumb:hover {
					background: #a8a8a8 !important;
				}
				&::-webkit-scrollbar-track {
					-webkit-box-shadow: inset 0 0 6px rgba(83, 83, 83, 0.07);
					background-color: #f1f1f1;
				}
				&::-webkit-scrollbar {
					height: 7px;
					background-color: #f1f1f1;
				}
				&::-webkit-scrollbar-thumb {
					background-color: #c1c1c1;
				}
			}
			.kanban-items-container{
				&::-webkit-scrollbar-thumb:hover {
					background: #a8a8a8 !important;
				}
				&::-webkit-scrollbar-track {
					-webkit-box-shadow: inset 0 0 6px rgba(83, 83, 83, 0.07);
					background-color: #f1f1f1;
				}
				&::-webkit-scrollbar {
					width: 7px;
					background-color: #f1f1f1;
				}
				&::-webkit-scrollbar-thumb {
					background-color: #c1c1c1;
				}
			}
			.container-production-line {
				display: flex;
				gap: 20px;
			}
			.production-line-card {
				width: 300px;
				border: 1px solid #ddd;
				border-radius: 8px;
				background-color: white;
				box-shadow: 0px 4px 8px rgba(0, 0, 0, 0.1);
				overflow: hidden;
			}
			.line-title {
				background-color: white;
				color: black;
				padding: 10px;
				text-align: center;
				font-weight: bold;
			}
			.kanban-item {
				padding: 8px;
				margin: 5px 0;
				background-color: #871a83;
				border-radius: 5px;
			}
			/* .non-draggable{
				cursor: default;
			} */
			.card-head{
				display: flex;
				justify-content: space-between;
				align-items: center;	
			}
			.cell-info {
				font-weight: 500;
			}
			.icon-group {
				display: flex;
				gap: 8px;
				align-items: center;
			}
			.custom-button {
				background-color: #871a83;
				color: white;
				width: 100%;
				padding: 8px;
				border: none;
				border-radius: 5px;
				margin: 10px 0;
			}
			/*!css*/
		</style>
	`
	sideNav := basepage.SideNav{
		MenuItems: []basepage.SideNavItem{
			{
				Name:     "Dashboard",
				Icon:     "fas fa-chart-pie",
				Link:     "/dashboard",
				Style:    "font-size:1rem;",
				UserType: basepage.UserType{Admin: true, Operator: true, Customer: true},
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
				Style:    "font-size:1rem;",
				UserType: basepage.UserType{Admin: true, Operator: true, Customer: true},
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
				Style:    "font-size:1rem;",
				UserType: basepage.UserType{Admin: true, Operator: true},
			},
			{
				Name:     "Heijunka Board",
				Icon:     "fas fa-calendar-alt",
				Link:     "/production-line",
				Style:    "font-size:1rem;",
				UserType: basepage.UserType{Admin: true, Operator: true},
				Selected: true,
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
				Style:    "font-size:1rem;",
				UserType: basepage.UserType{Customer: true},
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

	finalContent := `
		<!--html-->
		<div class="add-line-btn-container d-flex justify-content-between w-100">
			<button class="add-line-btn" id="add-Line-btn" data-bs-toggle="modal" data-bs-target="#AddLine">+ Add Line</button>
			<button  id="delete-kanban-btn" data-bs-toggle="modal"><span class="fa fa-trash"></span> Delete</button>
		</div>
		<div class="wrapper-production-line">
			<div class="container-production-line">` + contentHTML + `</div>
		</div>
		<!--!html-->
		`
	page := &basepage.BasePage{
		ExtraHeaders: css,
		SideNavBar:   sideNav,
		TopNavBar:    topNav,
		Username:     username,
		UserType:     usertype,
		Content:      string(template.HTML(finalContent)),
	}

	page.AddStyleCode(`
	/*css*/
	#collapseExample{
		font-size: 11px !important;
	}
	/*!css*/
	`)
	w.Write([]byte(page.Build()))
}

func DeleteProductionLineCell(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	userID := r.Header.Get("X-Custom-Userid")

	var payload struct {
		KRid []string `json:"KRid"`
		user string   `json:"userID"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if len(payload.KRid) == 0 {
		http.Error(w, "Please select kanban!", http.StatusBadRequest)
		return
	}

	payload.user = userID

	marshalData, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "Error while marshaling data", http.StatusInternalServerError)
		return
	}

	url := utils.RestURL + "/delete-production-line-cell"
	req, err := http.NewRequest(http.MethodDelete, url, bytes.NewReader(marshalData))
	if err != nil {
		http.Error(w, "Error creating request to REST service", http.StatusInternalServerError)
		return
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to send request to REST service", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Fail to create New Line", resp.StatusCode)
		return
	}
	w.WriteHeader(http.StatusOK)
}
