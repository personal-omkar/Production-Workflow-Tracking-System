package coldstorage

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"

	m "irpl.com/kanban-commons/model"
	"irpl.com/kanban-commons/utils"
	s "irpl.com/kanban-web/services"
)

type ColdStoragePage struct {
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

	RestURL = utils.JoinStr("http://", RestHost, ":", RestPort)
}
func (v ColdStoragePage) Build() string {
	var coldstorage []*m.ColdStorage

	// var compoundnames []s.DropDownOptions
	var html strings.Builder

	var Response struct {
		Pagination m.PaginationResp
		Data       []*m.ColdStorage
	}
	type TableConditions struct {
		Pagination m.PaginationReq `json:"pagination"`
		Conditions []string        `json:"conditions"`
	}

	var tablecondition TableConditions

	// Set default pagination
	tablecondition.Pagination = m.PaginationReq{
		Type:   "operator",
		Limit:  "10", // Default; overridden later by JS
		PageNo: 1,
	}

	// Fetch paginated data
	jsonValue, err := json.Marshal(tablecondition)
	if err != nil {
		slog.Error("Error marshaling table condition", "error", err)
	}

	resp, err := http.Post(RestURL+"/get-all-cold-storage-by-search-pagination", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		slog.Error("Error making POST request", "error", err)
	}

	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&Response); err != nil {
		slog.Error("error decoding paginated response body", "error", err)
	}
	var coldstorageTable s.TableCard
	tableTools := `<button type="button" class="btn  m-0 p-0" id="del-License-btn" data-bs-toggle="modal" data-bs-target="#EditCompound" > 
	<i class="fa fa-edit " style="color: #CF7AC2;"></i> 
   </button>`
	coldstorageTable.CardHeading = utils.DefaultsMap["cold_store_board"]
	coldstorageTable.BodyTables = s.CardTableBody{Columns: []s.CardTableBodyHeadCol{
		{
			Lable:        "Part Name",
			Name:         `Compound Name`,
			IsSearchable: true,
			ID:           "search-by-compound-name",
			DataField:    "compound_name",
			Type:         "input",
			Width:        "col-2",
		},
		{
			Name:  "Max Quantity",
			ID:    "search-min-quantity",
			Type:  "input",
			Width: "col-2",
		},
		{
			Name:  "Min Quantity",
			Type:  "input",
			ID:    "search-max-quantity",
			Width: "col-2",
		},
		{
			Name:  "Available Quantity",
			Type:  "input",
			ID:    "search-max-quantity",
			Width: "col-2",
		},

		{
			Name:  "Tools",
			Width: "col-1",
		},
	},

		ColumnsWidth: []string{"col-2", "col-2", "col-2", "col-2", "col-1"},
		Data:         coldstorage,
		Tools:        tableTools,
		ID:           "Inventory-Table",
	}
	coldstorageTable.BodyTables.Data = Response.Data
	var Pagination s.Pagination
	Pagination.CurrentPage = 1
	Pagination.Offset = 0
	Pagination.TotalRecords = Response.Pagination.TotalNo
	Pagination.PerPage, _ = strconv.Atoi(tablecondition.Pagination.Limit)
	Pagination.PerPageOptions = []int{10, 25, 50, 100, 200}
	coldstorageTable.CardFooter = Pagination.Build()

	AddCompModel := s.ModelCard{
		ID:      "EditCompound",
		Type:    "modal-md",
		Heading: `Edit Record  ( <span id="element-txt"></span> ) `,
		Form: s.ModelForm{FormID: "Inventory-Compound",
			FormAction: "/update-coldstorage-quantity", Footer: s.Footer{CancelBtn: true, Buttons: []s.FooterButtons{{BtnType: "submit", BtnID: "edit-compound", Text: "Update"}}},
			Inputfield: []s.InputAttributes{
				{Type: "text", Name: "editinventoryid", ID: "editinventoryid", Readonly: true, Hidden: true},
				{Type: "text", Name: "showcompoundname", ID: "showcompoundname", Label: `Part Name`, Disabled: true, Width: "w-100"},
				{Type: "hidden", Name: "editcompoundname", ID: "editcompoundname", Label: `Part Name`, Hidden: true, Width: "w-100"},
				{Type: "text", Name: "editavailableqty", ID: "editavailableqty", Label: `Available Quantity`, Disabled: true, Width: "w-100"},
				{Type: "text", Name: "editmaxqty", ID: "editmaxqty", Label: `Max Quantity`, Width: "w-100"},
				{Type: "text", Name: "editminqty", ID: "editminqty", Label: `Min Quantity`, Width: "w-100"},
			},
		},
	}
	js :=

		` <script>	
		//js
			$(document).ready(function () {
					// Attach event listeners to all searchable inputs
					document.querySelectorAll("[id^='search-by']").forEach(function (input) {
						input.addEventListener("input", function () {
							pagination(1);
						});
					});

					attachPerPageListener(); // On page load
				});
				// Pagination and search fetch handler
				function pagination(pageNo) {
					let searchCriteria = [];

					document.querySelectorAll("[id^='search-by']").forEach(function (input) {
						let field = input.dataset.field;
						let value = input.value.trim();
						if (value) {
							searchCriteria.push(field + " ILIKE '%" + value + "%'");
						}
					});

					let limit = document.getElementById("perPageSelect")?.value || "10";

					let requestData = {
						pagination: {
							Limit: limit,
							Pageno: pageNo
						},
						Conditions: searchCriteria
					};

					fetch("/coldstorage-pagination-search", {
						method: "POST",
						headers: {
							"Content-Type": "application/json"
						},
						body: JSON.stringify(requestData)
					})
					.then(response => response.json())
					.then(data => {
						let tableBody = document.getElementById("advanced-search-table-body");
						let cardFooter = document.querySelector(".card-footer");

						if (tableBody) {
							tableBody.innerHTML = data.tableBodyHTML;
						}
						if (cardFooter) {
							cardFooter.innerHTML = data.paginationHTML;
							attachPerPageListener(); // rebind dropdown listener
						}
					})
					.catch(err => console.error("Pagination error:", err));
				}

				// Re-bind per-page select on page update
				function attachPerPageListener() {
					let perPageSelect = document.getElementById("perPageSelect");
					if (perPageSelect) {
						perPageSelect.addEventListener("change", function () {
							pagination(1);
						});
					}
				}

				$(function(){
						const urlParams = new URLSearchParams(window.location.search);
    					const status = urlParams.get('status');
    					const msg = urlParams.get('msg');

						if (status) {
							showNotification(status,msg, () => {
								   removeQueryParams();
							});
						}

					$("#EditCompound").on("show.bs.modal" ,function(event){   
					
						var data = JSON.parse($(event.relatedTarget).closest("tr").attr("data-data"))		
						document.getElementById('editinventoryid').value = data.ID;
						document.getElementById('showcompoundname').value = data.CompoundName;
						document.getElementById('editcompoundname').value = data.CompoundName;
						document.getElementById('editminqty').value = data.MinQuantity; 
						document.getElementById('editmaxqty').value = data.MaxQuantity;
						document.getElementById('editavailableqty').value = data.AvailableQuantity; 
						$('#element-txt').text( data.CompoundName) 

						
						
					})

					$('#editminqty , #editmaxqty').on('input', function () {
						this.value = this.value.replace(/[^0-9]/g, '');
					});

					$('#editminqty').on('input change', function () {
						if(parseInt($('#editmaxqty').val()) < parseInt( this.value))	{
							var msg , status
							status=400
							msg="Minimum quantity cannot exceed the maximum quantity"
							showModelNotification(status,msg);
							$("#edit-compound").prop("disabled",true)
						}else{
							status="remove"
							msg=""
							showModelNotification(status,msg);
							$("#edit-compound").prop("disabled",false)
						}
					});

					$('#editmaxqty').on('input change', function () { 
						if(parseInt($('#editminqty').val()) > parseInt( this.value))	{
							var msg , status
							status=400
							msg="Maximum quantity cannot be less than the minimum quantity"
							showModelNotification(status,msg);
							$("#edit-compound").prop("disabled",true)
						}else{
							status="remove"
							msg=""
							showModelNotification(status,msg);
							$("#edit-compound").prop("disabled",false)
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

					function showModelNotification(status, msg, callback) {
						const notification = $('#modelnotification');
						var message = '';
						if (status === "200") {
							notification.removeClass("d-none")
							message = msg + '.';
							notification.removeClass("alert-danger").addClass("alert-success");
							
						} else if  (status === "remove") {
							message =  msg ;
							notification.removeClass("alert-success");
							notification.removeClass("alert-danger");
							notification.addClass("d-none")
						}else {
							notification.removeClass("d-none")
							message =  msg + '.';
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
					$('#search-by-compound-name').on('input', function () {
					pagination(1);
				});	
			})
		//!js
	</script>
`
	html.WriteString(coldstorageTable.Build())

	html.WriteString(AddCompModel.Build())
	html.WriteString(js)

	return html.String()
}
