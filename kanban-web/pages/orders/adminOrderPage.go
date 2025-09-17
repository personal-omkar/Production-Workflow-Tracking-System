package orders

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	m "irpl.com/kanban-commons/model"
	"irpl.com/kanban-commons/utils"
	s "irpl.com/kanban-web/services"
)

type AdminOrder struct {
	UserType string
	UserID   string
}

func (o *AdminOrder) Build() string {
	var Response struct {
		Pagination m.PaginationResp
		Data       []*m.OrderDetails
	}
	type TableConditions struct {
		Pagination m.PaginationReq `json:"pagination"`
		Conditions []string        `json:"conditions"`
	}

	var tablecondition TableConditions

	// Temporary pagination data
	tablecondition.Pagination = m.PaginationReq{
		Type:   "PendingOrders",
		Limit:  "10",
		PageNo: 1,
	}

	// Convert to JSON
	jsonValue, err := json.Marshal(tablecondition)
	if err != nil {
		slog.Error("Error marshaling table condition", "error", err)
	}

	resp, err := http.Post(utils.RestURL+"/get-order-pending-details-by-search-pagination", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		slog.Error("%s - error - %s", "Error making GET request", err)
	}

	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&Response); err != nil {
		slog.Error("error decoding response body", "error", err)
	}
	for _, value := range Response.Data {
		value.Status = strings.Title(value.Status)
		formattedDate := value.DemandDateTime.Format("02.01.2006")
		value.DemandDateTime, _ = time.Parse("02.01.2006", formattedDate)
	}

	var html strings.Builder
	html.WriteString(`
	<!--html-->
	<!--!html-->`)

	var vendorOrderTable s.TableCard
	tablebutton := `
	<!--html-->
		<button type="button" class="btn m-0 p-0" id="viewOrderModel" data-toggle="tooltip" data-placement="bottom" title="View Details"  data-bs-toggle="modal" data-bs-target="#UpdateOrderStatus"> 
				 <i class="fa fa-eye mx-2" style="color: #b250ad;"></i> 
		</button>
	<!--!html-->`
	vendorOrderTable.CardHeading = "Customer Orders"
	vendorOrderTable.BodyTables = s.CardTableBody{Columns: []s.CardTableBodyHeadCol{
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
			DataField:    "vendor_name",
			Width:        "col-1",
		},
		{
			Lable:        `<span id="search-span">Part Name</span>`,
			Name:         "Compound Name",
			ID:           "CompoundName",
			IsSearchable: true,
			Type:         "input",
			Width:        "col-1",
			DataField:    "compound_name",
			IsSortable:   true,
		},
		{
			Lable:        "Kanban Summary",
			Name:         "Cell No",
			ID:           "CellNo",
			IsSearchable: true,
			Type:         "input",
			DataField:    "cell_no",
			Width:        "col-1",
		},
		{
			Name:  "Demand Date Time",
			ID:    "demand_date",
			Width: "col-1",
		},
		{
			Lable:  "No. of Lots",
			Name:   "NoOFLots",
			ID:     "lot_no",
			Width:  "col-1",
			GetSum: true,
		},
		{
			Name:  "Status",
			Width: "col-1",
		},
		{
			Name:  "Action",
			Width: "col-1",
		},
	},
		ColumnsWidth: []string{"col-1", "col-1", "col-1", "col-1", "col-1", "col-1", "col-1", "col-1"},
		Data:         Response.Data,
		Buttons:      tablebutton,
	}
	var Pagination s.Pagination
	Pagination.CurrentPage = 1
	Pagination.Offset = 0
	Pagination.TotalRecords = Response.Pagination.TotalNo
	Pagination.PerPage, _ = strconv.Atoi(tablecondition.Pagination.Limit)
	Pagination.PerPageOptions = []int{10, 25, 50, 100, 200}
	vendorOrderTable.CardFooter = Pagination.Build()

	js := `
		<!--html-->
		<script>
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

			fetch("/pending-order-search-pagination", {
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
					attachEventListeners();
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

		function attachEventListeners() {
				const viewLinks = document.querySelectorAll('#viewOrderModel');
				viewLinks.forEach(link => {
					link.addEventListener('click', function (event) {
						// Get the parent row <tr> element
						const row = event.target.closest('tr');

					// Retrieve the data from the 'data-data' attribute
					const data = JSON.parse(row.getAttribute('data-data'));
					
					// Populate the modal form fields with the data
					document.getElementById('showCompoundCode').value = data.CompoundName;
					let date = new Date(data.DemandDateTime);
					let formattedDate = ("0" + date.getDate()).slice(-2) + "." + ("0" + (date.getMonth() + 1)).slice(-2) + "." + date.getFullYear();
					document.getElementById('showDemandDate').value = formattedDate;
					document.getElementById('showLotNo').value = data.NoOFLots;
					if (data.min_quantity=="0" && data.available_quantity=="0"){
						document.getElementById('ColdStroreInfo').innerText ="Compound is not available in Inventory!";
					}else{
						document.getElementById('ColdStroreInfo').innerText =
						"Cold Store :\n Minimum Quantity : " + data.min_quantity + ",\n Available Quantity : " + data.available_quantity +  ", \nIn-Process Quantity : "+ data.InventoryKanbanInProcessQty; 
					}

					// Convert strings to integers
					const availableQuantity = parseInt(data.available_quantity, 10);
					const NoOFLots = parseInt(data.NoOFLots, 10);
					const minQuantity = parseInt(data.min_quantity, 10);
					var inprocesskanban = parseInt(data.InventoryKanbanInProcessQty,10)
					
					let approveOrderText = "";
					let status = "";
					let dispatchQuantity = 0;
					let kanbanQuantity = 0;
					var kanbantxt=""

					if (inprocesskanban>minQuantity){
						if(availableQuantity < NoOFLots){
							if (availableQuantity==0){
								kanbantxt="Kanban (" +Math.abs((minQuantity  - availableQuantity)) +")"
								status = "approved";
								kanbanQuantity = Math.abs((minQuantity - availableQuantity))
								approveOrderText = kanbantxt;
								dispatchQuantity = availableQuantity;
							}
							else{
								kanbantxt="Kanban (" +Math.abs((minQuantity  - availableQuantity)) +")"
								status = "approved";
								kanbanQuantity =Math.abs((minQuantity - availableQuantity))
								approveOrderText =
									"Dispatch (" + availableQuantity + "), "+ kanbantxt;
								dispatchQuantity = availableQuantity;
							}
						}
						else if(availableQuantity >= NoOFLots){
							approveOrderText =
								"Dispatch (" + NoOFLots + ")";
							status = "dispatch";
							dispatchQuantity = NoOFLots;			
						}
					}else{
						if (NoOFLots <= availableQuantity) {
							if (availableQuantity >= minQuantity) {
							
								if (Math.abs(NoOFLots - availableQuantity) >= minQuantity) {
									approveOrderText = "Dispatch (" + NoOFLots + ")";
									status = "dispatch";
									dispatchQuantity = NoOFLots;
								} else {
									if (inprocesskanban < (minQuantity - Math.abs(NoOFLots - availableQuantity))){

									kanbantxt=", Kanban (" +Math.abs( (minQuantity - Math.abs(NoOFLots - availableQuantity)) - inprocesskanban) +")"
									status = "approved";
									kanbanQuantity = Math.abs( (minQuantity - Math.abs(NoOFLots - availableQuantity)) - inprocesskanban)
									}else{
									status = "dispatch";
									}
									approveOrderText = "Dispatch (" + NoOFLots + ") "+ kanbantxt;
									dispatchQuantity = NoOFLots;
									
								}
							} else {
						
								approveOrderText = "Dispatch (" + NoOFLots + "), Kanban (" + Math.abs( (minQuantity - Math.abs(NoOFLots - availableQuantity))-inprocesskanban) + ")";
								status = "approved";
								dispatchQuantity = NoOFLots;
								kanbanQuantity = (minQuantity - Math.abs(NoOFLots - availableQuantity))-inprocesskanban;
							}
						} else if (availableQuantity === 0) {
							
							approveOrderText = "Kanban (" + (minQuantity + NoOFLots - inprocesskanban) + ")";
							status = "approved";
						
							kanbanQuantity =(minQuantity + NoOFLots)- inprocesskanban;
							
						} else {
							
							if (inprocesskanban < (minQuantity + Math.abs(NoOFLots - availableQuantity))){
								kanbantxt=", Kanban (" +Math.abs((minQuantity + Math.abs(NoOFLots - availableQuantity)) - inprocesskanban) +")"
								status = "approved";
								kanbanQuantity =Math.abs((minQuantity + Math.abs(NoOFLots - availableQuantity)) - inprocesskanban)
							}
							approveOrderText =
								"Dispatch (" + availableQuantity + ")"+ kanbantxt;
								if (status == "approved"){
									status = "approved";
								}
								else{
									status = "dispatch";
								}
							dispatchQuantity = availableQuantity;
						}
					}
						document.getElementById("Approve_Order").innerText = approveOrderText;

						// Set data attributes on buttons for later use
						document.getElementById('Approve_Order').setAttribute('data-id', data.ID);
						document.getElementById('Approve_Order').setAttribute('compound-id', data.CompoundId);
						document.getElementById('Approve_Order').setAttribute('data-status', status);
						document.getElementById('Approve_Order').setAttribute('data-dispatch', dispatchQuantity);
						document.getElementById('Approve_Order').setAttribute('data-kanban', kanbanQuantity);
						document.getElementById('Reject_Order').setAttribute('data-id', data.ID);
						document.getElementById('Reject_Order').setAttribute('data-status', 'reject');
					});
				});

				// Add click event for 'Approve' (Approve Order) button
				document.getElementById('Approve_Order').addEventListener('click', function (event) {
					event.preventDefault();

					const id = this.getAttribute('data-id');
					const status = this.getAttribute('data-status');
					const dispatchQuantity = this.getAttribute('data-dispatch');
					const kanbanQuantity = this.getAttribute('data-kanban');
					const compoundId = this.getAttribute('compound-id'); 
					const NoOFLots = document.getElementById('showLotNo').value;

					var payload = { 
						ID: id, 
						Status: status, 
						NoOFLots : parseInt(NoOFLots,10), 
						DispatchQuantity: parseInt(dispatchQuantity, 10), 
						KanbanQuantity: parseInt(kanbanQuantity, 10), 
						compoundId: parseInt(compoundId ,10)
					};

					const openModals = document.querySelectorAll('.modal.show');
					openModals.forEach(modal => {
						const modalInstance = bootstrap.Modal.getInstance(modal);
						if (modalInstance) {
							modalInstance.hide();
						}
					});
					document.querySelectorAll('.modal-backdrop').forEach(backdrop => {
						backdrop.remove();
					});

					const confirmModal = new bootstrap.Modal(document.getElementById('Approve_Order_confirmation'));
					confirmModal.show();

					document.getElementById('Close_Approve_Confirmation_Modal').addEventListener('click', function () {
						const modalElement = document.getElementById('Approve_Order_confirmation');
						const modalInstance = bootstrap.Modal.getInstance(modalElement);
						
						if (modalInstance) {
							modalInstance.hide();
						}
						payload = {};
					});
					const submitButton = document.getElementById('Submit_Approve_Confirmation_Modal');
        			submitButton.replaceWith(submitButton.cloneNode(true));
					document.getElementById('Submit_Approve_Confirmation_Modal').addEventListener('click', function () {
						fetch("/update-order-status", {
							method: "POST",
							headers: {
								"Content-Type": "application/json"
							},
							body: JSON.stringify(payload)
						}).then(response => {
							if (!response.ok || response.ok) {
								return response.text().then(msg => {
									var data = JSON.parse(msg)
									 window.location.href = "/admin-orders?status=" + data.code + "&msg=" +data.message;
									throw new Error(msg);
								});
							}
							else{
								window.location.reload();
							}
						}).catch(error => {
							console.error("Error:", error);
						});
					});
				});

				// Add click event for 'Reject' (Reject Order) button
				document.getElementById('Reject_Order').addEventListener('click', function (event) {
					event.preventDefault();

					const id = this.getAttribute('data-id');
					const status = this.getAttribute('data-status');

					var payload = { ID: id, Status: status };

					const openModals = document.querySelectorAll('.modal.show');
					openModals.forEach(modal => {
						const modalInstance = bootstrap.Modal.getInstance(modal);
						if (modalInstance) {
							modalInstance.hide();
						}
					});
					document.querySelectorAll('.modal-backdrop').forEach(backdrop => {
						backdrop.remove();
					});

					const confirmModal = new bootstrap.Modal(document.getElementById('Reject_Order_confirmation'));
					confirmModal.show();


					document.getElementById('Close_Reject_Confirmation_Modal').addEventListener('click', function () {
						const modalElement = document.getElementById('Reject_Order_confirmation');
						const modalInstance = bootstrap.Modal.getInstance(modalElement);
						
						if (modalInstance) {
							modalInstance.hide();
						}
						payload ={};
					});
					const submitButton = document.getElementById('Submit_Reject_Confirmation_Modal');
        			submitButton.replaceWith(submitButton.cloneNode(true));
					document.getElementById('Submit_Reject_Confirmation_Modal').addEventListener('click', function () {
						fetch("/update-order-status", {
							method: "POST",
							headers: {
								"Content-Type": "application/json"
							},
							body: JSON.stringify(payload)
						}).then(response => {
							if (!response.ok || response.ok) {
								return response.text().then(msg => {
									var data = JSON.parse(msg)
									 window.location.href = "/admin-orders?status=" + data.code + "&msg=" +data.message;
									throw new Error(msg);
								});
							}
							else{
								window.location.reload();
							}
						}).catch(error => {
							console.error("Error:", error);
						});
					});
				});
			}

			function bindSorting() {
				document.querySelectorAll('th[data-sort]').forEach(function (header) {
					header.addEventListener('click', function (event) {
						// Prevent sorting if the click target is an input or inside one
						if (event.target.closest('input')) {
							return;
						}
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
								Pageno: 1
							},
							Conditions: searchCriteria
						};

						fetch("/sort-admin-order-table", {
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
								attachEventListeners();
								attachPerPageListener();
							} else {
								console.error("Error: Table body element not found.");
							}
						})
						.catch(error => console.error("Error fetching paginated results:", error));
					});
				});
			}


		document.addEventListener('DOMContentLoaded', function () {
			
			document.querySelectorAll("[id^='search-by']").forEach(function (input) {
				input.addEventListener("input", function () {
					pagination(1);
				});
			});
			
			attachPerPageListener();

			const urlParams = new URLSearchParams(window.location.search);
			const status = urlParams.get('status');
			const msg = urlParams.get('msg');
			if (status) {
				showNotification(status,msg, () => {
					removeQueryParams();
				});
			}


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
			attachEventListeners();
			bindSorting();
		});
		</script>

		<!--!html-->
		`
	viewOrderDetailsModel := s.ModelCard{
		ID:      "UpdateOrderStatus",
		Type:    "modal-md",
		Heading: "Order Details",
		Form: s.ModelForm{FormID: "accept_order",
			FormAction: "", Footer: s.Footer{CancelBtn: false, Buttons: []s.FooterButtons{{BtnType: "submit", BtnID: "Reject_Order", Text: "Reject", Style: "background-color:#c62f4a ;border:none"}, {BtnType: "submit", BtnID: "Approve_Order", Text: "Add"}}},
			Inputfield: []s.InputAttributes{
				{Type: "text", Name: "showCompoundCode", ID: "showCompoundCode", Label: `Compound Code`, Readonly: true, Width: "w-100", Disabled: true},
				{Type: "text", Name: "showDemandDate", ID: "showDemandDate", Label: `Demand Date/Time `, Readonly: true, Width: "w-100", Disabled: true},
				{Type: "text", Name: "showLotNo", ID: "showLotNo", Label: `Number of Lots`, Readonly: true, Width: "w-100", Disabled: true},
			},
			Info: s.InfoLine{IsVisible: true, ID: "ColdStroreInfo", Lable: "Cold Store"},
		},
	}

	confirmReject := s.ConfirmationModal{
		For:   "Delete",
		ID:    "Reject_Order_confirmation",
		Title: "Attention!",
		Body: []string{
			"Are you sure you want to reject this order?",
			"This action cannot be undone.",
		},
		Footer: s.Footer{CancelBtn: false, Buttons: []s.FooterButtons{{BtnType: "submit", BtnID: "Close_Reject_Confirmation_Modal", Text: "Close", Style: "background-color:#636E7E; border:none;"}, {BtnType: "submit", BtnID: "Submit_Reject_Confirmation_Modal", Text: "Reject"}}},
	}

	confirmApprove := s.ConfirmationModal{
		For:   "Submit",
		ID:    "Approve_Order_confirmation",
		Title: "Confirm!",
		Body: []string{
			"Are you sure you wants to approve/dispatch the order?",
		},
		Footer: s.Footer{CancelBtn: false, Buttons: []s.FooterButtons{{BtnType: "submit", BtnID: "Close_Approve_Confirmation_Modal", Text: "Close", Style: "background-color:#636E7E; border:none;"}, {BtnType: "submit", BtnID: "Submit_Approve_Confirmation_Modal", Text: "Submit"}}},
	}

	html.WriteString(vendorOrderTable.Build())
	html.WriteString(js)
	html.WriteString(`</div>`)
	html.WriteString(viewOrderDetailsModel.Build())
	html.WriteString(confirmReject.Build())
	html.WriteString(confirmApprove.Build())
	return html.String()
}

type SortRequest struct {
	Column string `json:"column"`
}

var isSort bool = true

func SortAdminOrderTable(w http.ResponseWriter, r *http.Request) {

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
		Data       []*m.OrderDetails
	}

	jsonValue, err := json.Marshal(req)
	if err != nil {
		slog.Error("Error marshaling table condition", "error", err)
	}

	resp, err := http.Post(utils.RestURL+"/get-order-pending-details-by-search-pagination", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		slog.Error("%s - error - %s", "Error making GET request", err)
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

	if isSort {
		sort.Slice(Response.Data, func(i, j int) bool {
			return Response.Data[i].CompoundName < Response.Data[j].CompoundName
		})
		isSort = false
	} else {
		isSort = true
	}

	for _, value := range Response.Data {
		value.Status = strings.Title(value.Status)
	}

	var html strings.Builder
	html.WriteString(`
	<!--html-->
	<!--!html-->`)

	var vendorOrderTable s.TableCard
	tablebutton := `
	<!--html-->
		<button type="button" class="btn m-0 p-0" id="viewOrderModel" data-toggle="tooltip" data-placement="bottom" title="View Details"  data-bs-toggle="modal" data-bs-target="#UpdateOrderStatus"> 
				 <i class="fa fa-eye mx-2" style="color: #b250ad;"></i> 
		</button>
	<!--!html-->`
	vendorOrderTable.CardHeading = "Customer Orders"
	vendorOrderTable.BodyTables = s.CardTableBody{Columns: []s.CardTableBodyHeadCol{
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
			IsSortable:   true,
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
			Lable:  "No. of Lots",
			Name:   "NoOFLots",
			ID:     "lot_no",
			Width:  "col-1",
			GetSum: true,
		},
		{
			Name:  "Status",
			Width: "col-1",
		},
		{
			Name:  "Action",
			Width: "col-1",
		},
	},
		ColumnsWidth: []string{"col-1", "col-1", "col-1", "col-1", "col-1", "col-1", "col-1", "col-1"},
		Data:         Response.Data,
		Buttons:      tablebutton,
	}

	tableBodyHTML := vendorOrderTable.BodyTables.RenderBodyColumns()

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
