package vendorcompany

import (
	"bytes"
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"

	m "irpl.com/kanban-commons/model"
	"irpl.com/kanban-commons/utils"
	s "irpl.com/kanban-web/services"
)

type VendorCompanyPage struct {
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
func (v VendorCompanyPage) Build() string {
	var vendor []*m.Vendors
	var vendortable []*m.VendorCompanyTable
	var compoundlist []m.Compounds
	var compound []*m.CompoundsDataByVendor
	var prodline []m.ProdLine
	var action []s.DropDownOptions
	var compoundnames []s.DropDownOptions
	var html strings.Builder
	action = append(action, s.DropDownOptions{Text: "Select Action", Value: ""})
	//fetching vendor records
	var req struct {
		Conditions []string `json:"Conditions"`
	}
	reqjsonValue, err := json.Marshal(req)
	if err != nil {
		log.Printf("Error: Error marshaling user data: %v", err)
	}
	vendorresp, err := http.Post(RestURL+"/get-all-vendors-data", "application/json", bytes.NewBuffer(reqjsonValue))
	if err != nil {
		slog.Error("%s - error - %s", "Error making GET request", err)
	}

	defer vendorresp.Body.Close()

	if err := json.NewDecoder(vendorresp.Body).Decode(&vendor); err != nil {
		slog.Error("error decoding response body", "error", err)
	}

	//fetching prod line records
	prodlineresp, err := http.Get(RestURL + "/get-all-prod-line-data")
	if err != nil {
		slog.Error("%s - error - %s", "Error making GET request", err)
	}

	defer prodlineresp.Body.Close()

	if err := json.NewDecoder(prodlineresp.Body).Decode(&prodline); err != nil {
		slog.Error("error decoding response body", "error", err)
	}

	//fetching compounds records
	compoundsresp, err := http.Get(RestURL + "/get-all-compound-data")
	if err != nil {
		slog.Error("%s - error - %s", "Error making GET request", err)
	}

	defer compoundsresp.Body.Close()

	if err := json.NewDecoder(compoundsresp.Body).Decode(&compoundlist); err != nil {
		slog.Error("error decoding response body", "error", err)
	}

	for _, i := range vendor {
		var comp string
		// fetching component records by vendor
		coprresp, err := http.Post(RestURL+"/get-compound-data-by-vendor?key=id&value="+strconv.Itoa(i.ID), "application/json", bytes.NewBuffer(reqjsonValue))
		if err != nil {
			slog.Error("%s - error - %s", "Error making GET request", err)
		}

		defer coprresp.Body.Close()

		if err := json.NewDecoder(coprresp.Body).Decode(&compound); err != nil {
			slog.Error("error decoding response body", "error", err)
		}
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
								<label class="form-check-label m-0 px-1" for="`, strconv.Itoa(i.KbRootId), `" data="Kanban-No">
									`, i.KanbanNo, `
								</label>									
								<label class="form-check-label m-0 px-1 d-none" for="`, strconv.Itoa(i.KbRootId), `" data="Approved-Date">
									`, utils.FormatStringDate(i.DemandDate, "date"), `
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

		vendortable = append(vendortable, &temptable)
	}
	for _, i := range prodline {
		if i.Status {
			var tempopt s.DropDownOptions
			tempopt.Value = strconv.Itoa(i.Id)
			tempopt.Text = i.Name
			action = append(action, tempopt)
		}
	}

	for _, i := range compoundlist {
		if i.Status {

			var comp s.DropDownOptions
			comp.Text = i.CompoundName
			comp.Value = strconv.Itoa(i.Id)

			compoundnames = append(compoundnames, comp)
		}
	}

	var vendorcompanyTable s.TableCard
	tablebutton := `<button type="button" style="background-color:#871a83;color:white;" class="btn  m-0 p-2" id="del-License-btn" data-bs-toggle="modal" data-bs-target="#AddComp" > 
 						Add Comp
					</button>`
	allCheckBox := `<input style="width: 20px; height: 20px;" class="form-check-input me-2 select_all_compound" type="checkbox" value="" id="select_all_compound">`
	vendorcompanyTable.CardHeading = "Kanban Board"
	vendorcompanyTable.CardHeadingActions.Component = []s.CardHeadActionComponent{
		{ComponentName: "modelbutton", ComponentType: s.ActionComponentElement{ModelButton: s.ModelButtonAttributes{ID: "help-modal", Name: "help", Type: "button", Text: " <i class=\"fas fa-question-circle\"></i>  Help", ModelID: "#help-dialog"}}, Width: "col-2"},
		{ComponentName: "dropdown", ComponentType: s.ActionComponentElement{DropDown: s.DropdownAttributes{ID: "selectMenu-action-line", Name: "selectMenu-action-line", Options: action, Label: "Production Line", Disabled: true}}, Width: "col-8"},
		// {ComponentName: "button", ComponentType: s.ActionComponentElement{Button: s.ButtonAttributes{ID: "deleteBtn", Disabled: true, Colour: "#c62f4a", Name: "deleteBtn", Type: "button", Text: "Delete"}}, Width: "col-2"},
		{ComponentName: "button", ComponentType: s.ActionComponentElement{Button: s.ButtonAttributes{ID: "sortBtn", Colour: "#007BFF", Name: "sortBtn", Type: "button", Text: "<i class=\"fas fa-sort\"></i>  Sort"}}, Width: "col-2"},
	}
	vendorcompanyTable.CardHeadingActions.Style = "direction: ltr;"
	vendorcompanyTable.BodyTables = s.CardTableBody{Columns: []s.CardTableBodyHeadCol{
		{
			Name:  "Tools",
			Lable: " ",
		},
		{
			Name:         `Vendor Code`,
			IsSearchable: true,
			DataField:    "vendor_code",
			IsCheckbox:   true,
			ID:           "search-sn",
			Type:         "input",
			Width:        "col-1",
		},
		{
			Name:         "Vendor Name",
			IsSearchable: true,
			DataField:    "vendor_name",
			ID:           "search-MachinName",
			Type:         "input",
			Width:        "col-1",
		},
		{
			Name:  "Button",
			ID:    "search-MachinName",
			Width: "col-1",
		},
		{Lable: "Part Name",
			Name:             "Compound Code",
			IsSearchable:     true,
			DataField:        "compound_name",
			SearchFieldWidth: "w-25",
			ID:               "search-MachinName",
			Type:             "input",
			Width:            "col-9",
			Style:            "max-height: 25vh !important; overflow-y: auto;",
		},
		// {
		// 	Name:         "Vendor-Action",
		// 	IsSearchable: true,
		// 	ActionList:   action,
		// 	ID:           "search-MachinName",
		// 	Type:         "action",
		// 	Width:        "",
		// },
	},
		ColumnsWidth: []string{"", "col-1", "col-1", "col-1", "col-9 d-flex flex-wrap w-100 "},
		Data:         vendortable,
		Buttons:      tablebutton,
		Tools:        allCheckBox,
		ID:           "VendorComp",
	}

	var customerNoteOptions = []s.DropDownOptions{
		{Value: "regular", Text: "Regular", Selected: true},
		{Value: "urgent", Text: "Urgent"},
		{Value: "mosturgent", Text: "Most Urgent"},
	}

	AddCompModel := s.ModelCard{
		ID:      "AddComp",
		Type:    "modal-md",
		Heading: "Add Comp",
		Form: s.ModelForm{FormID: "AddComp",
			FormAction: "/add-compound-data-by-vendor", Footer: s.Footer{CancelBtn: false, Buttons: []s.FooterButtons{{BtnType: "button", DataSubmitName: "AddComp", BtnID: "Add-Comp-Submit", Text: "Load"}}},
			Inputfield: []s.InputAttributes{
				{Type: "text", DataType: "text", Name: "VendorCode", ID: "showVendorCode", Label: `Vendor Code`, Readonly: true, Width: "w-100", Required: true},
				{Type: "text", DataType: "text", Name: "VendorName", ID: "showVendorName", Label: `Vendor Name`, Readonly: true, Width: "w-100"},
			},
			Dropdownfield: []s.DropdownAttributes{
				{Label: "Note", DataType: "text", Name: "Note", ID: "Note", Options: customerNoteOptions, Width: "w-100"},
				{Label: "Comp Code", DataType: "int", Name: "CompoundCode", ID: "CompCode", Options: compoundnames, Width: "w-50"},
			},
			IncrementalButton: s.IncrementalButtonAttributes{Label: "", DataType: "int", IsVisible: true, Width: "w-50", MinValue: 0, MaxValue: 100},
		},
	}

	AddVendorModel := s.ModelCard{
		ID:      "AddNewVendor",
		Type:    "modal-md",
		Heading: "Add Vendor",
		Form: s.ModelForm{FormID: "Add-Vendor",
			FormAction: "/create-new-vendor", Footer: s.Footer{CancelBtn: false, Buttons: []s.FooterButtons{{BtnType: "submit", BtnID: "Add-New-Vendor", Text: "Add"}}},
			Inputfield: []s.InputAttributes{
				{Type: "text", Name: "addVendorCode", ID: "addVendorCode", Label: `Vendor Code`, Width: "w-100"},
				{Type: "text", Name: "addVendorName", ID: "addVendorName", Label: `Vendor Name`, Width: "w-100"},
				{Type: "text", Name: "addContactInfo", ID: "addContactInfo", Label: `Contact Info`, Width: "w-100"},
			},
			TextArea: []s.TextAreaAttributes{
				{Label: "Address", Name: "addAddress", ID: "addAddress"},
			},
		},
	}
	DeleteKanban := s.InfoModal{
		ID:        "ConfirmDeleteKanban",
		ModelSize: "modal-lg",
		Title:     "Confirm Deletion",
		Body: []string{
			`<div>
				<div class="input-group-prepend">
					<Label style="font-size:20px">Once the Kanban is deleted, it cannot be undone. Are you sure you want to delete the Kanban?</Label>
				</div>
			</div>
			`,
		},
		Footer: s.Footer{CancelBtn: false, Buttons: []s.FooterButtons{{BtnType: "button", BtnID: "deleteKanban", Text: "Delete"}}},
	}

	helpDialog := utils.JoinStr(`
	<div class="modal fade" id="help-dialog" tabindex="-1" aria-labelledby="helpModalLabel" aria-hidden="true">
		<div class="modal-dialog">
			<div class="modal-content">
			<div class="modal-header">
				<h4 class="modal-title" id="helpModalLabel" style="color:#871a83;">Help</h4>
				<button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
			</div>
			<div class="modal-body">
				<div class="row ms-3">
					<h5 class="p-0">Color Codes - Priority</h5>
					<div class="row mt-1">
						<a class="col-2 btn btn-primary btn-sm" style="background:`, utils.KanbanPriorityColors["regular"]["bg-color"], `"></a><p class="col mb-0">- Regular</p>
					</div>
					<div class="row mt-1">
						<a class="col-2 btn btn-primary btn-sm" style="background:`, utils.KanbanPriorityColors["urgent"]["bg-color"], `"></a><p class="col mb-0">- Urgent</p>
					</div>
					<div class="row mt-1">
						<a class="col-2 btn btn-primary btn-sm" style="background:`, utils.KanbanPriorityColors["mosturgent"]["bg-color"], `"></a><p class="col mb-0">- Most Urgent</p>
					</div>
				</div>
			</div>
			<div class="modal-footer">
				<button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Close</button>
			</div>
			</div>
		</div>
	</div>
	`)

	html.WriteString(helpDialog)

	js := ` <script>
	//js
		document.addEventListener('DOMContentLoaded', () => {
			checkBoxEvents();
			var asce = true;

			// Handle search 
			document.querySelectorAll("[id^='search-by']").forEach(function (input) {
				input.addEventListener("input", function () {
					performSortedSearch();
				});
			});

			// For sort button
			$('#sortBtn').on('click', function () {
				asce = !asce;
				document.getElementById('selectMenu-action-line').disabled = true;
				performSortedSearch();
			});


			function performSortedSearch() {
				var url = asce ? "/sort-kanban-asce" : "/sort-kanban-desc";

				// Update icon
				var sortIconSvg = document.querySelector('#sortBtn svg');
				if (sortIconSvg) {
					if (asce) {
						sortIconSvg.setAttribute('data-icon', 'sort-up');
						sortIconSvg.setAttribute('viewBox', '0 0 320 512');
						sortIconSvg.innerHTML = '<path fill="currentColor" d="M279 224H41c-21.4 0-32.1-25.9-17-41l119-119c9.4-9.4 24.6-9.4 33.9 0l119 119c15.1 15.1 4.4 41-17 41z"></path>';
						sortIconSvg.style.fill = 'currentColor';
						sortIconSvg.style.stroke = 'none';
					} else {
						sortIconSvg.setAttribute('data-icon', 'sort-down');
						sortIconSvg.setAttribute('viewBox', '0 0 320 512');
						sortIconSvg.innerHTML = '<path d="M41 288h238c21.4 0 32.1 25.9 17 41l-119 119c-9.4 9.4-24.6 9.4-33.9 0l-119-119c-15.1-15.1-4.4-41 17-41z"></path>';
						sortIconSvg.style.fill = 'none';
						sortIconSvg.style.stroke = 'currentColor';
						sortIconSvg.style.strokeWidth = '20';
					}
				}

				// Get current flow from URL
				var urlParams = new URLSearchParams(window.location.search);
				var flow = urlParams.get("flow") || "";

				// Collect current search input filters
				var searchCriteria = [];
				document.querySelectorAll("[id^='search-by']").forEach(function (input) {
					var field = input.dataset.field;
					var value = input.value.trim();
					if (value) {
						searchCriteria.push(field + " ILIKE '%" + value + "%'");
					}
				});

				var requestData = {
					Conditions: searchCriteria,
					Flow: flow
				};

				fetch(url, {
					method: "POST",
					headers: {
						"Content-Type": "application/json"
					},
					body: JSON.stringify(requestData)
				})
				.then(function (response) {
					if (!response.ok) throw new Error("HTTP error! Status: " + response.status);
					return response.json();
				})
				.then(function (data) {
					var tableBody = document.getElementById("advanced-search-table-body");
					if (tableBody) {
						tableBody.innerHTML = "";
						tableBody.insertAdjacentHTML("beforeend", data.tableBodyHTML);

						$(document).trigger('reset-selection');
						checkBoxEvents();
					}
				})
				.catch(function (error) {
					console.error("Error in sorted search:", error);
				});
			}
			
			function checkBoxEvents(){
				// Handle Select All checkbox functionality
				document.querySelectorAll('#select_all_compound').forEach(selectAll => {
					selectAll.addEventListener('change', function() {
						const row = this.closest('tr');
						const checkboxes = row ? row.querySelectorAll('.component-code') : document.querySelectorAll('.component-code');
						checkboxes.forEach(checkbox => checkbox.checked = this.checked);
						updateSelectAllState();
						toggleActionMenu();
					});
				});
				// Handle individual checkbox changes
				document.querySelectorAll('.component-code').forEach(checkbox => {
					checkbox.addEventListener('change', () => {
						updateSelectAllState();
						toggleActionMenu();
					});
				});
			}
			// Function to update "Select All" checkbox state based on individual checkboxes
			function updateSelectAllState() {
				document.querySelectorAll('#select_all_compound').forEach(selectAll => {
					const row = selectAll.closest('tr');
					const checkboxes = row ? row.querySelectorAll('.component-code') : document.querySelectorAll('.component-code');
					selectAll.checked = checkboxes.length > 0 && [...checkboxes].every(cb => cb.checked);
					selectAll.indeterminate = checkboxes.length > 0 && [...checkboxes].some(cb => cb.checked) && !selectAll.checked;
				});
			}

			// Enable/disable action dropdown based on checked checkboxes
			function toggleActionMenu() {
				const anyChecked = document.querySelectorAll('.component-code:checked').length > 0;
				document.getElementById('selectMenu-action-line').disabled = !anyChecked;
			}
		});
		$(document).ready(function () {
			function checkVendorLotLimit(quantity) {
			const vendorCode = $('#showVendorCode').val();
			if (!vendorCode || quantity <= 0) {
				$('#Add-Comp-Submit').prop('disabled', true);
				$('#limit-warning').hide().text('');
				return;
			}

			const lotData = {
				VendorCode: vendorCode,
				NoOFLots: quantity,
				DemandDateTime: new Date().toISOString()
			};

			$.ajax({
				url: '/check-vendor-lot-limit-by-vendor-code',
				type: 'POST',
				contentType: 'application/json',
				data: JSON.stringify(lotData),
				success: function () {
					$('#Add-Comp-Submit').prop('disabled', false);
					$('#limit-warning').hide().text('');
				},
				error: function (xhr) {
					$('#Add-Comp-Submit').prop('disabled', true);

					let resp;
					try {
						resp = JSON.parse(xhr.responseText);
					} catch (e) {
						resp = { message: 'Limit exceeded', exceed_by: '?' };
					}

					var msg = 'Lot limit exceeded by ' + resp.exceed_by + ': ' + resp.message;
					$('#limit-warning').remove();
					var warningHtml = '<div id="limit-warning" class="text-danger mt-2 small me-auto">' + msg + '</div>';
					$('#AddComp .modal-footer').prepend(warningHtml);
					}
				});
			 }
						
			$(document).on('click', '#btn-plus', function () {
			let $quantity = $('#quantity');
			let current = parseInt($quantity.val()) || 0;
			current = current + 1;
			$quantity.val(current);
			checkVendorLotLimit(current);
			});

			$(document).on('click', '#btn-minus', function () {
				let $quantity = $('#quantity');
				let current = parseInt($quantity.val()) || 0;
				if (current > 0) {
					current = current - 1;
					$quantity.val(current);
					checkVendorLotLimit(current);
				}
			});

			$(document).on('input', '#quantity', function () {
				const val = parseInt($(this).val()) || 0;
				checkVendorLotLimit(val);
			});

			$('#AddComp').on('shown.bs.modal', function () {
				checkVendorLotLimit(parseInt($('#quantity').val()) || 0);
			});
		});

		$(function() {
			const urlParams = new URLSearchParams(window.location.search);
			const status = urlParams.get('status');
			const msg = urlParams.get('msg');
			const flow = urlParams.get('flow');
			

			if (status) {
				showNotification(status, msg, removeQueryParams);
			}

			// Show modal and populate fields
			$("#AddComp").on("show.bs.modal", function(event) {
				var data = JSON.parse($(event.relatedTarget).closest("tr").attr("data-data"));
				$('#showVendorName').val(data.VendorName);
				$('#showVendorCode').val(data.VendorCode);
				$('#Add-Comp-Submit').prop('disabled', true);
				$('#limit-warning').remove();
			});

			let selectedItems = [];

			$(document).on('reset-selection', function () {
				selectedItems = [];
			});

			// Handle individual checkbox changes
			$(document).on('change', '.component-code', function () {
				const value = $(this).val();
				if ($(this).is(':checked')) {
					if (!selectedItems.includes(value)) {
						selectedItems.push(value);
					}
				} else {
					selectedItems = selectedItems.filter(item => item !== value);
				}
				updateSelectAllStateForRow($(this).closest('tr')); // scope to row
			});

			// Handle "Select All" per row
			$(document).on('change', '.select_all_compound', function () {
				const $row = $(this).closest('tr');
				const $checkboxes = $row.find('.component-code');

				if ($(this).is(':checked')) {
					$checkboxes.each(function () {
						const value = $(this).val();
						this.checked = true;
						if (!selectedItems.includes(value)) {
							selectedItems.push(value);
						}
					});
				} else {
					$checkboxes.each(function () {
						const value = $(this).val();
						this.checked = false;
						selectedItems = selectedItems.filter(item => item !== value);
					});
				}
			});

			// Helper to update select_all_compound state for a row
			function updateSelectAllStateForRow($row) {
				const $checkboxes = $row.find('.component-code');
				const $selectAll = $row.find('.select_all_compound');

				const total = $checkboxes.length;
				const checked = $checkboxes.filter(':checked').length;

				if (checked === total) {
					$selectAll.prop('checked', true).prop('indeterminate', false);
				} else if (checked > 0) {
					$selectAll.prop('checked', false).prop('indeterminate', true);
				} else {
					$selectAll.prop('checked', false).prop('indeterminate', false);
				}
			}

			// Handle production line assignment
			$('#selectMenu-action-line').on('change', function() {	
				const lineId = $(this).val();
				var jsondata = {
					LineID: lineId,
					KbDataId: selectedItems,
				};
				
				$.post("/add-compound-data-to-production-line", JSON.stringify(jsondata), function(response) {
					if (response.status === "200") {
						window.location.href = "/vendor-company?flow="+flow+"&status=" + response.status + "&msg=" + encodeURIComponent(response.msg);
					}
				}, 'json').fail(function(xhr) {
					var redirectURL = "/vendor-company?flow="+flow+"&status=" + xhr.status + "&msg=" + encodeURIComponent("Failed to process the request");
					window.location.href = redirectURL;
				});
			});
				$(document).on("click","#Add-Comp-Submit",function(){
					
					var group = $(this).attr("data-submit");
					var result = {}
					var validated = true;
				
					$("[data-group='" + group + "']").find("[data-name]").each(function () {
						
						if ($(this).attr("data-validate") && $(this).val().trim().length === 0 ||$(this).attr("data-validate") && $(this).val()==="Nil" ) {
							$(this).css("background-color", "rgba(128, 0, 128, 0.1)");
							const label = $(this).closest("label").length 
										? $(this).closest("label") 
										: $(this).siblings("label").length 
										? $(this).siblings("label") 
										: $(this).parent().siblings("label");
					
							if (label.length) {
								label.find(".required-label").remove();
								label.siblings(".required-label").remove();

								$("<span class='required-label'>Required</span>").css({
									color: "red",
									fontSize: "1em",
									"margin-left": "0.5rem",
								}).insertAfter(label);
							}
					
							validated = false;
						} else {
							$(this).css("background-color", "rgb(255, 255, 255)");
							
							const label = $(this).closest("label").length 
										? $(this).closest("label")
										: $(this).siblings("label").length
										? $(this).siblings("label")
										: $(this).parent().siblings("label");

							if (label.length) {
								label.siblings(".required-label").remove();
								label.find(".required-label").remove();
							}
						}
					});
					
					if (validated){
						
						$("[data-group='"+group+"'").find("[data-name]").each(function(){
							if ($(this).is("select")){
								var selectedValue = $(this).find(":selected").val();
								if ($(this).attr("data-type") === "int") {
									result[$(this).attr("data-name")] = parseInt(selectedValue, 10); // Convert to integer
								} else {
									result[$(this).attr("data-name")] = selectedValue; // Keep as string
								}
							} else if ($(this).attr("data-type") == "date"){
								var userDate = $(this).val();                 
								if (userDate.includes(":")) {
									result[$(this).attr("data-name")] = userDate
								}else{
									var formattedDate = formatDateToYYYYMMDD(userDate);
									result[$(this).attr("data-name")] = new Date(formattedDate)      
								}
							} else if ($(this).attr("data-type") == "int") {
								result[$(this).attr("data-name")] = parseInt($(this).val());
							} else { 
								result[$(this).attr("data-name")] = $(this).val();
							}						
						})
						
							const lotData = {
									VendorCode:result.VendorCode,
									NoOFLots: result.quantity
								};
								
							$.post("/check-vendor-lot-limit-by-vendor-code", JSON.stringify(lotData), function (xhr, status, error) {
									$.post("/add-compound-data-by-vendor", JSON.stringify(result), function (xhr, status, error)  {	
										var resp = JSON.parse(xhr)
										 window.location.href = "/vendor-company?status=" + resp.code + "&msg=" + resp.message;
									});
								}, 'json').fail(function (xhr, status, error) {
								 	 var resp=JSON.parse(xhr.responseText)
									 window.location.href = "/vendor-company?status=" + resp.status + "&msg=" + resp.message +" Exceeded by : " +resp.exceed_by;	
									
								});													
					}
				})
			// Show notifications
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

				notification.html(message).show();
				setTimeout(() => {
					notification.fadeOut(callback);
				}, 5000);
			}

			// Remove query parameters from URL
			function removeQueryParams() {
				var newUrl = window.location.protocol + "//" + window.location.host + window.location.pathname;
				window.history.replaceState({}, document.title, newUrl);
			}
		});

		$('#AddComp').on('shown.bs.modal', function () {
			const $select = $('#CompCode');

			if ($select.hasClass('select2-hidden-accessible')) {
				$select.select2('destroy');
			}

			$select.select2({
				theme: 'bootstrap-5',
				placeholder: 'Search Compound',
				allowClear: true,
				dropdownParent: $('#AddComp'),
				ajax: {
				transport: function (params, success, failure) {
					$.ajax({
					url: '/get-compounds-list-by-parem',
					data: { partName: params.data.q },
					dataType: 'json',
					success: function (data) {
						success({
						results: $.map(data, function (item) {
							return { id: item.Id, text: item.Text };
						})
						});
					},
					error: failure
					});
				},
				minimumInputLength: 1
				}
			}).on('select2:open', function () {
				const input = document.querySelector('.select2-search__field');
				if (input) input.focus();
			});
			});

		$('#AddComp').on('show.bs.modal', function () {
			$('#quantity').val(0);
			$('#Note').val('regular');
			$('#compound-dropdown').val('');
			$('#Add-Comp-Submit').prop('disabled', true);
			$('#limit-warning').remove();
			const $select = $('#CompCode');
			if ($select.hasClass('select2-hidden-accessible')) {
				$select.val(null).trigger('change');
			}
		});
	//!js
	</script>
	`
	html.WriteString(vendorcompanyTable.Build())

	html.WriteString(AddCompModel.Build())
	html.WriteString(AddVendorModel.Build())
	html.WriteString(DeleteKanban.Build())
	html.WriteString(js)

	return html.String()
}
