package packing

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

type PackingDispatchPage struct {
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

func (v PackingDispatchPage) Build() string {
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
	vendorresp, err := http.Post(RestURL+"/get-all-vendors-data?status=packing", "application/json", bytes.NewBuffer(reqjsonValue))
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
		coprresp, err := http.Post(RestURL+"/get-packing-compound-data-by-vendor?key=id&value="+strconv.Itoa(i.ID), "application/json", bytes.NewBuffer(reqjsonValue))
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

	vendorcompanyTable.CardHeading = "Packing/Dispatch"
	vendorcompanyTable.CardHeadingActions.Component = []s.CardHeadActionComponent{
		{ComponentName: "button", ComponentType: s.ActionComponentElement{Button: s.ButtonAttributes{ID: "sortBtn", Colour: "#007BFF", Name: "sortBtn", Type: "button", Text: " <i class=\"fas fa-sort\"></i>  Sort"}}, Width: "col-3"},
		{ComponentName: "button", ComponentType: s.ActionComponentElement{Button: s.ButtonAttributes{ID: "dispatch", Disabled: true, Name: "dispatch", Type: "button", Text: "Dispatch"}}, Width: "col-4"},
		{ComponentName: "modelbutton", ComponentType: s.ActionComponentElement{ModelButton: s.ModelButtonAttributes{ID: "help-modal", Name: "help", Type: "button", Text: " <i class=\"fas fa-question-circle\"></i>  Help", ModelID: "#help-dialog"}}, Width: "col-2"},
	}

	vendorcompanyTable.BodyTables = s.CardTableBody{Columns: []s.CardTableBodyHeadCol{
		{
			Name:         `Vendor Code`,
			IsSearchable: true,
			IsCheckbox:   true,
			ID:           "search-sn",
			Type:         "input",
			Width:        "col-1",
			DataField:    "vendor_code",
		},
		{
			Name:         "Vendor Name",
			IsSearchable: true,
			ID:           "search-MachinName",
			Type:         "input",
			Width:        "col-1",
			DataField:    "vendor_name",
		},

		{
			Name:             "Compound Code",
			IsSearchable:     true,
			SearchFieldWidth: "w-25",
			ID:               "search-MachinName",
			Type:             "input",
			Width:            "col-7",
			DataField:        "compound_name",
		},
	},
		ColumnsWidth: []string{"col-1", "col-1", "col-9 d-flex flex-wrap w-100 "},
		Data:         vendortable,
		ID:           "VendorComp",
	}
	addNotes := s.InfoModal{
		ID:        "SubmitKanban",
		ModelSize: "modal-lg",
		Title:     "Confirm Dispatch",
		Body: []string{
			`<div>
				<div class="input-group-prepend">
					<Label style="font-size:20px">Notes</Label>
				</div>
				<textarea class="form-control" rows="5" id="kbNotes" data-name="kbNotes"  name="kbNotes"  placeholder="Notes" data-type="string"></textarea>
			</div>
			`,
		},
		Footer: s.Footer{CancelBtn: false, Buttons: []s.FooterButtons{{BtnType: "button", BtnID: "dispatchOrder", Text: "Dispatch"}}},
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
					var url = asce ? "/packing-sort-kanban-asce" : "/packing-sort-kanban-desc";
				
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
				const urlParams = new URLSearchParams(window.location.search);
				const status = urlParams.get('status');
				const msg = urlParams.get('msg');
				
				if (status) {
					showNotification(status,msg, () => {
						removeQueryParams();
					});
				}

			

			
				$(document).on('change', '.component-code', function () {
					if ($('.component-code:checked').length > 0) {
						$('#dispatch').prop('disabled', false); 
					} else {
						$('#dispatch').prop('disabled', true); 
					}
				});

				 

				
				$('#dispatch').on('click', function () {
					$("#SubmitKanban").modal("show");
				});

				let selectedItems = [];
				$(document).on('reset-selection', function () {
					selectedItems = [];
					// Disable buttons
					$('#dispatch').prop('disabled', true);
				});
				let selectedallItems=[];

				$(document).on('change', '.component-code', function () {
					const value = $(this).val();
					if ($(this).is(':checked')) {
						if (!selectedItems.includes(value)) {
							selectedItems.push(value);
						}
					} else {
						selectedItems = selectedItems.filter(item => item !== value);
					}
				});

			
				$(document).on('change', '.select_all_compound', function () {
					if ($(this).is(':checked')) {
						$('.component-code').each(function () {
							if ($(this).is(':checked')) {
								const value = $(this).val();
								if (!selectedallItems.includes(value)) {
									selectedallItems.push(value);
								}
							}
						});
					} else {
						selectedallItems = []; 
					}
				});

				// Detect the selected value when the selection changes
				$('#dispatchOrder').on('click', function () {
					if (selectedallItems.length > 0) {
					selectedallItems.forEach(item => {
						if (!selectedItems.includes(item)) {
							selectedItems.push(item);
						}
					});
				}
					
					notes=$("#kbNotes").val()
					var jsondata = {
						Notes:notes,
						KbDataId: selectedItems,
					};
					
					
						
					$.post("/update-compound-status-to-dispatch", JSON.stringify(jsondata), function(response) {
						
						if (response.code === 200) {
						
							window.location.href = "/packing-dispatch-page?status=" + response.code + "&msg=" + encodeURIComponent(response.message);
						}
					}, 'json').fail(function(xhr) {
						var redirectURL = "/packing-dispatch-page?status=" + xhr.status + "&msg=" + encodeURIComponent("Failed to process the request");
						window.location.href = redirectURL;
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

	html.WriteString(addNotes.Build())
	html.WriteString(vendorcompanyTable.Build())

	html.WriteString(js)

	return html.String()
}
