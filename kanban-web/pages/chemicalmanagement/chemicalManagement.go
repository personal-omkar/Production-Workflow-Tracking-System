package chemicalmanagement

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

type ChemicalManagement struct {
	ChemicalId string
	Username   string
	UserType   string
	UserID     string
}

const DefaultRestHost string = "0.0.0.0" // Default port if not set in env
const DefaultRestPort string = "7200"    // Default port if not set in env

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

func (u *ChemicalManagement) Build() string {
	var html strings.Builder

	var Response struct {
		Pagination m.PaginationResp
		Data       []*m.ChemicalTypes
	}
	type TableConditions struct {
		Pagination m.PaginationReq `json:"pagination"`
		Conditions []string        `json:"conditions"`
	}

	var tablecondition TableConditions

	// Temporary pagination data
	tablecondition.Pagination = m.PaginationReq{
		Type:   "",
		Limit:  "10",
		PageNo: 1,
	}

	// Convert to JSON
	jsonValue, err := json.Marshal(tablecondition)
	if err != nil {
		slog.Error("Error marshaling table condition", "error", err)
	}

	resp, err := http.Post(RestURL+"/get-all-chemical-by-search-pagination", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		slog.Error("%s - error - %s", "Error making GET request", err)
	}

	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&Response); err != nil {
		slog.Error("error decoding response body", "error", err)
	}

	var ChemicalTable s.TableCard
	tableTools := `<button type="button" class="btn  m-0 p-0" id="edit-Chemical-btn" data-bs-toggle="modal"  > 
	<i class="fa fa-edit " style="color: #871a83;"></i> 
   </button>`
	ChemicalTable.CardHeading = "Chemical Type Master"
	ChemicalTable.CardHeadingActions.Component = []s.CardHeadActionComponent{
		{
			ComponentName: "modelbutton",
			ComponentType: s.ActionComponentElement{
				ModelButton: s.ModelButtonAttributes{Type: "button", Text: "Add Chemical Type", ModelID: "#AddChemicalModel"}},
			Width: "col-4",
		}, {
			ComponentName: "button",
			ComponentType: s.ActionComponentElement{
				Button: s.ButtonAttributes{ID: "importChemicals", Name: "importChemicals", Type: "button", Text: `Import<input style="display:none" class="clickable btn btn-success" type="file" id="om-import-file" name="file" accept=".xlsx">`, Class: "bg-transparent", CSS: "color: #871a83; border: 1px solid #871a83;"},
			},
			Width: "col-3",
		},
		{
			ComponentName: "button",
			ComponentType: s.ActionComponentElement{
				Button: s.ButtonAttributes{ID: "downloadTmp", Name: "downloadTmp", Type: "button", Text: ` Template <i class="fa fa-download" aria-hidden="true"></i> `, Class: "bg-transparent", AdditionalAttr: `data-link="/files/sample.pdf"`, CSS: "color: #871a83; border: 1px solid #871a83;", IsDownloadBtn: true, DownloadLink: `/static/templates/import-chemical-master.xlsx`},
			},
			Width: "col-3",
		},
	}

	ChemicalTable.BodyTables = s.CardTableBody{Columns: []s.CardTableBodyHeadCol{
		{
			Lable:        "Chemical Type",
			Name:         "Type",
			Type:         "input",
			DataField:    "type",
			IsSearchable: true,
			Width:        "col-2",
		},
		{
			Lable:        "Conv Code",
			Name:         "ConvCode",
			Type:         "input",
			DataField:    "conv_code",
			IsSearchable: true,
			Width:        "col-2",
		},
		{
			Name:  "Status",
			Type:  "input",
			ID:    "search-max-quantity",
			Width: "col-1",
		},
		{
			Name:  "Tools",
			Width: "col-1",
		},
	},
		ColumnsWidth: []string{"col-2", "col-2", "col-2", "col-2"},
		Data:         Response.Data,
		Tools:        tableTools,
		ID:           "Inventory-Table",
	}
	var Pagination s.Pagination
	Pagination.CurrentPage = 1
	Pagination.Offset = 0
	Pagination.TotalRecords = Response.Pagination.TotalNo
	Pagination.PerPage, _ = strconv.Atoi(tablecondition.Pagination.Limit)
	Pagination.PerPageOptions = []int{10, 25, 50, 100, 200}
	ChemicalTable.CardFooter = Pagination.Build()

	var statusopt []s.DropDownOptions
	statusopt = append(statusopt, s.DropDownOptions{Text: "Enabled", Value: "true"}, s.DropDownOptions{Text: "Disabled", Value: "false"})
	// Create Dialogue
	AddChemicalModel := s.ModelCard{
		ID:      "AddChemicalModel",
		Type:    "modal-m",
		Heading: "Add Chemical Type",
		Form: s.ModelForm{FormID: "AddChemicalModel",
			FormAction: "", Footer: s.Footer{CancelBtn: true, Buttons: []s.FooterButtons{{BtnType: "button", DataSubmitName: "AddChemicalModel", BtnID: "add-chemical-submit", Text: "Save"}}},
			Dropdownfield: []s.DropdownAttributes{
				{Label: "Status", Name: "Status", DataType: "bool", Options: statusopt},
			},
			Inputfield: []s.InputAttributes{
				{Type: "text", Name: "Type", DataType: "string", ID: "ChemicalType", Label: `Chemical Type`, Width: "w-100", Required: true},
				{Type: "text", Name: "ConvCode", DataType: "string", ID: "ChemicalConvCode", Label: `Conv Code`, Width: "w-100", Required: true},
			},
		},
	}

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

				fetch("/chemical-search-pagination", {
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
			
		$(document).ready(function() {
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
				showNotification(status, msg, removeQueryParams);
			}
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
			
						$(document).on("click", "#edit-Chemical-btn", function(event) {  
						var data = JSON.parse($(this).closest("tr").attr("data-data"));
						// Get modal content for the clicked user
						$.get("chemical-management-card?id=" + data.ID, function(response) {
							let modal = $("#EditChemicalModel");
							
							$("#additional-content").html(response)
							// Show the modal after updating/appending
							$("#EditChemicalModel").modal("show");
						}, 'json');

						
					})



				$(document).on("click","#add-chemical-submit,#edit-chemical-submit",function(){
				
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
							} else if ($(this).attr("data-type") == "bool") {
								if ($(this).val()=="true"){
									result[$(this).attr("data-name")] =true;
								}else{
									result[$(this).attr("data-name")] = false;
								}
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
						}else if ($(this).attr("data-type") == "bool"){
						
							if ($(this).val()=="true"){
								result[$(this).attr("data-name")] = true;
							}else{
								result[$(this).attr("data-name")] = false;
							}
						} else { 
							result[$(this).attr("data-name")] = $(this).val();
						}						
					})
					
						
				$.post("/create-new-or-update-existing-chemical", JSON.stringify(result), function() {}, 'json')
					.done(function(response, textStatus, jqXHR) {
							
							window.location.href = "/chemical-management?status=" +response.code+"&msg="  +response.message;
					})
					.fail(function(xhr, status, error) {
							data = JSON.parse(xhr.responseText);
						window.location.href = "/chemical-management?status=" +data.code+"&msg="  +data.message;
					});		
					
				}
			})	
		});
		$(document).on("click", "#importChemicals", function(event) {
				
				switch (event.target.id) {
					case 'importChemicals':
					$('#om-import-file').trigger('click');
						break;
				}
			});
			$(document).on("change", '#om-import-file', function(e) {
						const url = "/import-chemical-data"
						if (url) {
							
							const file = e.target.files[0];
							console.log(file)
							if (file) {
								const formData = new FormData();
								formData.append("file", file);
						
								$.ajax({
									url: url,
									type: 'POST',
									data: formData,
									processData: false,
									contentType: false,
									success: function(data) {
										window.location.href = "/chemical-management?status=200&msg=File uploaded successfully"
									},
									error: function(xhr, status, error) {
										console.error('Failed to upload file:', error);
										// window.location.href = "/operator-management?status=500&msg=Failed to upload file"
									}
								});
							}
						} else {
							console.error('No URL mapped for the file input ID:', fileInputId);
						}
			});

	//!js
	</script>
	`
	slog.Info("Fetched Chemicals", "count", len(Response.Data))

	html.WriteString(ChemicalTable.Build())
	html.WriteString(AddChemicalModel.Build())
	html.WriteString(`</div>`)
	html.WriteString(js)
	return html.String()

}
