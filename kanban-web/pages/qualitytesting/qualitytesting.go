package qualitytesting

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

type QualityTesting struct {
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

func (v QualityTesting) Build() string {
	var vendor []*m.Vendors
	var vendortable []*m.VendorCompanyTable

	var compound []m.CompoundsDataByVendor

	var html strings.Builder

	// fetching vendor records
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

	for _, i := range vendor {
		var comp string
		// fetching component records by vendor
		coprresp, err := http.Post(RestURL+"/get-quality-compound-data-by-vendor?key=id&value="+strconv.Itoa(i.ID), "application/json", bytes.NewBuffer(reqjsonValue))
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
								<label class="form-check-label m-0 px-1" for="`, i.KanbanNo, `" data="Kanba-Number">
									`, i.KanbanNo, `
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

	var vendorcompanyTable s.TableCard
	allCheckBox := `<input style="width: 20px; height: 20px;" class="form-check-input me-2 select_all_compound" type="checkbox" value="" id="select_all_compound">`
	vendorcompanyTable.CardHeading = "Quality Testing"
	vendorcompanyTable.CardHeadingActions.Component = []s.CardHeadActionComponent{
		{ComponentName: "button", ComponentType: s.ActionComponentElement{Button: s.ButtonAttributes{ID: "sortBtn", Colour: "#007BFF", Name: "sortBtn", Type: "button", Text: " <i class=\"fas fa-sort\"></i>  Sort"}}, Width: "col-3"},
		{ComponentName: "button", ComponentType: s.ActionComponentElement{Button: s.ButtonAttributes{ID: "approveBtn", Disabled: true, Name: "approveBtn", Type: "button", Text: "Approve"}}, Width: "col-3"},
		{ComponentName: "button", ComponentType: s.ActionComponentElement{Button: s.ButtonAttributes{ID: "rejectBtn", Disabled: true, Colour: "#c62f4a", Name: "rejectBtn", Type: "button", Text: "Reject"}}, Width: "col-3"},
		{ComponentName: "modelbutton", ComponentType: s.ActionComponentElement{ModelButton: s.ModelButtonAttributes{ID: "help-modal", Name: "help", Type: "button", Text: " <i class=\"fas fa-question-circle\"></i>  Help", ModelID: "#help-dialog"}}, Width: "col-2"},
	}

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
			Type:         "input",
			Width:        "col-1",
		},
		{
			Name:         "Vendor Name",
			DataField:    "vendor_name",
			IsSearchable: true,
			ID:           "search-MachinName",
			Type:         "input",
			Width:        "col-1",
		},

		{
			Name:             "Compound Code",
			IsSearchable:     true,
			DataField:        "compound_name",
			SearchFieldWidth: "w-25",
			ID:               "search-MachinName",
			Type:             "input",
			Width:            "col-10",
		},
	},
		ColumnsWidth: []string{"col-1", "col-1", "col-1", "col-9 d-flex flex-wrap w-100 "},
		Data:         vendortable,
		Tools:        allCheckBox,
	}
	addNotes := s.InfoModal{
		ID:        "TestKanban",
		ModelSize: "modal-lg",
		Title:     "Confirm Approval",
		Body: []string{
			`<div>
				<div class="input-group-prepend">
					<Label for="approvekbNotes" style="font-size:20px">Kanban approval is final and cannot be undone. Are you sure you want to proceed?<br><br>Note</Label>
					<textarea class="form-control" rows="5" id="approvekbNotes" data-name="approvekbNotes"  name="approvekbNotes"  placeholder="Note" data-type="string"></textarea>
				</div>
			</div>
			`,
		},
		Footer: s.Footer{CancelBtn: false, Buttons: []s.FooterButtons{{BtnType: "button", BtnID: "approveKanban", Text: "Approve"}}},
	}

	RejectKanban := s.InfoModal{
		ID:        "TestKanbanReject",
		ModelSize: "modal-lg",
		Title:     "Confirm Rejection",
		Body: []string{
			`<div>
				<div class="input-group-prepend">
					<Label style="font-size:20px" for="rejectkbNotes">Once the Kanban is rejected, it cannot be undone. Are you sure you want to reject the Kanban? <br><br>Note</Label>
					<textarea class="form-control" rows="5" id="rejectkbNotes" data-name="rejectkbNotes"  name="rejectkbNotes"  placeholder="Note" data-type="string"></textarea>
				</div>
			</div>
			`,
		},
		Footer: s.Footer{CancelBtn: false, Buttons: []s.FooterButtons{{BtnType: "button", BtnID: "rejectKanban", Text: "Reject"}}},
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
			
			$(document).ready(function() {
				var asce = true;				
				// For sort button
				$('#sortBtn').on('click', function () {
					asce = !asce;
					performSortedSearch();
				});
				
				function performSortedSearch() {
					var url = asce ? "/quality-sort-kanban-asce" : "/quality-sort-kanban-desc";
				
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
					var jsondata={
					"Criteria": {
							vendor_name: $('#search-by-VendorName').val(),
							compound_name: $('#search-by-CompoundCode').val(),
							vendor_code: $('#search-by-VendorCode').val(),  
						}
					};

					$.post(url, JSON.stringify(jsondata), function (response) {
						$('#advanced-search-table-body').html(response);
						$(document).trigger('reset-selection');
					});				
				}
				
				$('.search-input').on('input change', function() {
					performSortedSearch();	
				});
			});

			$(function(){

				// Handle "Select All" checkbox using event delegation
				document.addEventListener('change', function (e) {
					if (e.target.matches('#select_all_compound')) {
						const selectAll = e.target;
						const row = selectAll.closest('tr');
						const checkboxes = row ? row.querySelectorAll('.component-code') : document.querySelectorAll('.component-code');
						checkboxes.forEach(cb => cb.checked = selectAll.checked);
						updateSelectAllState();
						toggleActionMenu();
					}
				});

				// Handle individual component checkbox changes using event delegation
				document.addEventListener('change', function (e) {
					if (e.target.matches('.component-code')) {
						updateSelectAllState();
						toggleActionMenu();
					}
				});

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
					document.getElementById('approveBtn').disabled = !anyChecked;
					document.getElementById('rejectBtn').disabled = !anyChecked;
				}


				const urlParams = new URLSearchParams(window.location.search);
				const status = urlParams.get('status');
				const msg = urlParams.get('msg');
				
				if (status) {
					showNotification(status,msg, () => {
						removeQueryParams();
					});
				}

			

				$('.component-code').on('change', () => {
					if ($('.component-code:checked').length > 0) {
						$('#approveBtn').prop('disabled', false); 
						$('#rejectBtn').prop('disabled', false); 
					} else {
						$('#approveBtn').prop('disabled', true); 
						$('#rejectBtn').prop('disabled', true); 
					}
				});

				 

				
				$('#approveBtn').on('click', function () {
					selectedItems = $('.component-code:checked').map(function () {
				return $(this).val();
				}).get();
					$("#TestKanban").modal("show");
				});

				$('#rejectBtn').on('click', function () {
					selectedItems = $('.component-code:checked').map(function () {
				return $(this).val();
				}).get();
					$("#TestKanbanReject").modal("show");
				});

				let selectedItems = [];
				$(document).on('reset-selection', function () {
					selectedItems = [];
					// Disable buttons
					$('#approveBtn').prop('disabled', true);
					$('#rejectBtn').prop('disabled', true);
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
				$('#approveKanban').on('click', function () {
					var jsondata = {
						Notes: $.trim($('#approvekbNotes').val()), 
						KbDataId: selectedItems,
					};
					$.post("/update-compound-status-to-packing", JSON.stringify(jsondata), function(response) {
						if (response.code === 200) {
							window.location.href = "/quality-testing?status=" + response.code + "&msg=" + encodeURIComponent(response.message);
						}
					}, 'json').fail(function(xhr) {
							window.location.href = "/quality-testing?status=" + xhr.status + "&msg=" + xhr.responseText;
						});
					});

				$('#rejectKanban').on('click', function () {
					var jsondata = {
						Notes: $.trim($('#rejectkbNotes').val()), 
						KbDataId: selectedItems,
					};
					$.post("/update-compound-quality-status-to-reject", JSON.stringify(jsondata), function(response) {
						if (response.code === 200) {
							window.location.href = "/quality-testing?status=" + response.code + "&msg=" + encodeURIComponent(response.message);
						}
					}, 'json').fail(function(xhr) {
						window.location.href = "/quality-testing?status=" + xhr.status + "&msg=" + xhr.responseText;
					});
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
			})
				

			//!js
		</script>
`

	html.WriteString(RejectKanban.Build())
	html.WriteString(addNotes.Build())
	html.WriteString(vendorcompanyTable.Build())

	html.WriteString(js)

	return html.String()
}
