package compounds

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

type CompoundsManagement struct {
	Username string
	UserType string
	UserID   string
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
func (c *CompoundsManagement) Build() string {

	var Response struct {
		Pagination m.PaginationResp
		Data       []*m.Compounds
	}
	type TableConditions struct {
		Pagination m.PaginationReq `json:"pagination"`
		Conditions []string        `json:"conditions"`
	}
	var tablecondition TableConditions

	// Temporary pagination data
	tablecondition.Pagination = m.PaginationReq{
		Type:   "operator",
		Limit:  "10",
		PageNo: 1,
	}
	tablecondition.Conditions = []string{"compound_name ILIKE '%%'"}
	jsonValue, err := json.Marshal(tablecondition)
	if err != nil {
		slog.Error("Error marshaling table condition", "error", err)
	}

	resp, err := http.Post(RestURL+"/get-all-compounds-by-search-pagination", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		slog.Error("%s - error - %s", "Error making GET request", err)
	}

	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&Response); err != nil {
		slog.Error("error decoding response body", "error", err)
	}
	var html strings.Builder

	var compoundstatus = []s.DropDownOptions{{Text: "False", Value: "false"}, {Text: "True", Value: "true"}}
	tablebutton := `
	<!--html-->
			<button type="button" class="btn m-0 p-0" id="viewCompoundDetails" data-toggle="tooltip" data-placement="bottom" title="View Details"> 
				<i class="fa fa-edit mx-2" style="color: #b250ad;"></i> 
			</button>
			<!--!html-->`

	html.WriteString(`
		<!--html-->
		<!--!html-->`)

	var CompoundTable s.TableCard

	CompoundTable.CardHeading = "Part Master"
	CompoundTable.CardHeadingActions.Component = []s.CardHeadActionComponent{
		{
			ComponentName: "modelbutton",
			ComponentType: s.ActionComponentElement{
				ModelButton: s.ModelButtonAttributes{
					ID:      "add-new-compound",
					Name:    "add-new-compound",
					Type:    "button",
					Text:    "Add New Part",
					ModelID: "#AddCompoundModel",
				},
			},
			Width: "col-4",
		},
		{
			ComponentName: "button",
			ComponentType: s.ActionComponentElement{
				Button: s.ButtonAttributes{ID: "importParts", Name: "importParts", Type: "button", Text: `Import<input style="display:none" class="clickable btn btn-success" type="file" id="pm-import-file" name="file" accept=".xlsx">`, Class: "bg-transparent", CSS: "color: #871a83; border: 1px solid #871a83;"},
			},
			Width: "col-3",
		},
		{
			ComponentName: "button",
			ComponentType: s.ActionComponentElement{
				Button: s.ButtonAttributes{ID: "downloadTmp", Name: "downloadTmp", Type: "button", Text: ` Template <i class="fa fa-download" aria-hidden="true"></i> `, Class: "bg-transparent", AdditionalAttr: `data-link="/files/sample.pdf"`, CSS: "color: #871a83; border: 1px solid #871a83;", IsDownloadBtn: true, DownloadLink: `/static/templates/import-compound-master.xlsx`},
			},
			Width: "col-3",
		},
	}
	CompoundTable.BodyTables = s.CardTableBody{Columns: []s.CardTableBodyHeadCol{
		{
			Name:       `Sr. No.`,
			IsCheckbox: true,
			ID:         "search-sn",
			Width:      "col-1",
		},
		{
			Lable:        "Part Name",
			Name:         "CompoundName",
			ID:           "Customer Name",
			Width:        "col-1",
			DataField:    "compound_name",
			IsSearchable: true,
			Type:         "input",
		},
		{
			Lable:        "SCADA Code",
			Name:         "SCADACode",
			ID:           "scadacode",
			Width:        "col-1",
			DataField:    "scada_code",
			IsSearchable: true,
			Type:         "input",
		},
		{
			Lable:        "SAP Code",
			Name:         "SAPCode",
			ID:           "sapcode",
			Width:        "col-1",
			DataField:    "sap_code",
			IsSearchable: true,
			Type:         "input",
		},
		{
			Name:  "Status",
			ID:    "compound_status",
			Width: "col-1",
		},
		{
			Name:  "Action",
			Width: "col-1",
		},
	},
		ColumnsWidth: []string{"col-1", "col-1", "col-1", "col-1", "col-1", "col-1"},

		Data:    Response.Data,
		Buttons: tablebutton,
	}

	var Pagination s.Pagination
	Pagination.CurrentPage = 1
	Pagination.Offset = 0
	Pagination.TotalRecords = Response.Pagination.TotalNo
	Pagination.PerPage, _ = strconv.Atoi(tablecondition.Pagination.Limit)
	Pagination.PerPageOptions = []int{10, 25, 50, 100, 200}
	CompoundTable.CardFooter = Pagination.Build()

	//Add Dialogue
	AddCompoundModel := s.ModelCard{
		ID:      "AddCompoundModel",
		Type:    "modal-md",
		Heading: "Add New Part",
		Form: s.ModelForm{FormID: "AddUserModel",
			FormAction: "", Footer: s.Footer{CancelBtn: true, Buttons: []s.FooterButtons{{BtnType: "button", DataSubmitName: "AddCompoundModel", BtnID: "add-Compound-submit", Text: "Save"}}},
			Inputfield: []s.InputAttributes{
				{Type: "text", DataType: "text", Name: "compoundName", ID: "compoundName", Label: `Part Name`, Width: "w-100", Required: true},
				{Type: "text", DataType: "text", Name: "scadaCode", ID: "scadaCode", Label: `SCADA Code`, Width: "w-100", Required: true},
				{Type: "text", DataType: "text", Name: "sapCode", ID: "sapCode", Label: `SAP Code`, Width: "w-100", Required: true},
			},

			TextArea: []s.TextAreaAttributes{
				{Label: "Part Description", DataType: "text", Name: "compoundDescription", ID: "compoundDescription"},
			},
			Dropdownfield: []s.DropdownAttributes{
				{Label: "Is Active", DataType: "text", Name: "isactive", Options: compoundstatus, ID: "compIsActive", Width: "w-100"},
			},
		},
	}

	//Edit Dialogue
	EditCompoundModel := s.ModelCard{
		ID:      "EditCompoundModel",
		Type:    "modal-md",
		Heading: "Edit Part",
		Form: s.ModelForm{FormID: "AddUserModel",
			FormAction: "", Footer: s.Footer{CancelBtn: true, Buttons: []s.FooterButtons{{BtnType: "button", DataSubmitName: "AddCompModel", BtnID: "update-comp-submit", Text: "Update"}}},
			Inputfield: []s.InputAttributes{
				{Type: "hidden", Name: "EditCompoundId", ID: "EditCompoundId", Label: `User ID`, Width: "w-100", Hidden: true},
				{Type: "text", Name: "EditcompoundName", ID: "EditcompoundName", Label: `Part Name`, Width: "w-100", Required: true},
				{Type: "text", Name: "EditScadaCode", ID: "EditScadaCode", Label: `SCADA Code`, Width: "w-100", Required: true},
				{Type: "text", Name: "EditSapCode", ID: "EditSapCode", Label: `SAP Code`, Width: "w-100", Required: true},
			},
			TextArea: []s.TextAreaAttributes{
				{Label: "Part Description", DataType: "text", Name: "EditcompoundDescription", ID: "EditcompoundDescription"},
			},
			Dropdownfield: []s.DropdownAttributes{
				{Label: "Is Active", DataType: "text", Name: "isactive", Options: compoundstatus, ID: "EditcompIsActive", Width: "w-100"},
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

				let limit = document.getElementById("perPageSelect")?.value || "15";
				let requestData = {
					pagination: {
						Limit: limit,
						Pageno: pageNo
					},
					Conditions: searchCriteria
				};

				fetch("/compounds-search-pagination", {
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
						tableBody.innerHTML = data.tableBodyHTML;
						cardFooter.innerHTML = data.paginationHTML;

						attachPerPageListener();
						reattachAllEvents(); // Important to re-attach events after DOM update
					} else {
						console.error("Error: Table body element not found.");
					}
				})
				.catch(error => console.error("Error fetching paginated results:", error));
			}

			function attachPerPageListener() {
				let perPageSelect = document.getElementById("perPageSelect");
				if (perPageSelect) {
					perPageSelect.addEventListener("change", function () {
						pagination(1);
					});
				}
			}

			function reattachAllEvents() {
				// Edit button click handler
				document.querySelectorAll("#viewCompoundDetails").forEach(function (button) {
					button.addEventListener("click", function () {
						var row = button.closest("tr");
						var rowData = JSON.parse(row.getAttribute("data-data"));
						document.getElementById("EditCompoundId").value = rowData.ID;
						document.getElementById("EditScadaCode").value = rowData.SCADACode || "";
						document.getElementById("EditSapCode").value = rowData.SAPCode || "";
						document.getElementById("EditcompoundName").value = rowData.CompoundName;
						document.getElementById("EditcompoundDescription").value = rowData.Description;
						document.getElementById("EditcompIsActive").value = rowData.Status ? "true" : "false";

						var modal = new bootstrap.Modal(document.getElementById("EditCompoundModel"));
						modal.show();
					});
				});

				// Clear validation message on input
				document.getElementById("compoundName")?.addEventListener("input", function () {
					this.setCustomValidity("");
				});
				document.getElementById("EditcompoundName")?.addEventListener("input", function () {
					this.setCustomValidity("");
				});
			}

			$(document).ready(function () {
				// Live search inputs
				document.querySelectorAll("[id^='search-by']").forEach(function (input) {
					input.addEventListener("input", function () {
						pagination(1);
					});
				});

				attachPerPageListener();
				reattachAllEvents();
			});

			document.addEventListener("DOMContentLoaded", function () {
				function showNotification(status, msg, callback) {
					var notification = $('#notification');
					var message = status === "200"
						? "Success! " + msg + "."
						: "Fail! " + msg + ".";

					notification
						.removeClass("alert-success alert-danger")
						.addClass(status === "200" ? "alert-success" : "alert-danger")
						.html(message)
						.show();

					setTimeout(function () {
						notification.fadeOut(function () {
							if (callback) callback();
						});
					}, 5000);
				}

				function removeQueryParams() {
					var newUrl = window.location.protocol + "//" + window.location.host + window.location.pathname;
					window.history.replaceState({}, document.title, newUrl);
				}

				// Import handler
				$(document).on("click", "#importParts", function (event) {
					if (event.target.id === "importParts") {
						$('#pm-import-file').trigger('click');
					}
				});

				$(document).on("change", '#pm-import-file', function (e) {
					const url = "/import-compound-data";
					const file = e.target.files[0];
					if (file) {
						const formData = new FormData();
						formData.append("file", file);

						$.ajax({
							url: url,
							type: 'POST',
							data: formData,
							processData: false,
							contentType: false,
							success: function () {
								window.location.href = "/compounds-management?status=200&msg=File uploaded successfully";
							},
							error: function (xhr) {
								window.location.href = "/compounds-management?status=500&msg=Failed to upload file";
							}
						});
					}
				});

				// Show notification on page load if params present
				let urlParams = new URLSearchParams(window.location.search);
				let status = urlParams.get("status");
				let msg = urlParams.get("msg");
				if (status) {
					showNotification(status, msg, removeQueryParams);
				}

				// Add Compound
				document.getElementById("add-Compound-submit")?.addEventListener("click", function () {
					var compoundName = String(document.getElementById("compoundName").value);
					var compoundData = {
						ID: 0,
						CompoundName: compoundName,
						SCADACode: document.getElementById("scadaCode").value,
						SAPCode: document.getElementById("sapCode").value,
						Description: document.getElementById("compoundDescription").value,
						Status: document.getElementById("compIsActive").value === "true",
					};

					if (compoundName.trim() === "") {
						var field = document.getElementById("compoundName");
						field.setCustomValidity("Part Name is required");
						field.reportValidity();
					} else {
						document.getElementById("compoundName").setCustomValidity("");
						$.post("/add-update-compound", JSON.stringify(compoundData), function () {}, "json").fail(function (xhr) {
							window.location.href = "/compounds-management?status=" + xhr.status + "&msg=" + xhr.responseText;
						});
					}
				});

				// Update Compound
				document.getElementById("update-comp-submit")?.addEventListener("click", function () {
					var compoundName = String(document.getElementById("EditcompoundName").value);
					var updatedData = {
						ID: parseInt(document.getElementById("EditCompoundId").value, 10),
						CompoundName: compoundName,
						SCADACode: document.getElementById("EditScadaCode").value,
						SAPCode: document.getElementById("EditSapCode").value,
						Description: document.getElementById("EditcompoundDescription").value,
						Status: document.getElementById("EditcompIsActive").value === "true",
					};

					if (compoundName.trim() === "") {
						var field = document.getElementById("EditcompoundName");
						field.setCustomValidity("Part Name is required");
						field.reportValidity();
					} else {
						document.getElementById("EditcompoundName").setCustomValidity("");
						$.post("/add-update-compound", JSON.stringify(updatedData), function () {}, "json").fail(function (xhr) {
							window.location.href = "/compounds-management?status=" + xhr.status + "&msg=" + xhr.responseText;
						});
					}
				});
			});
		//!js
		</script>
	`
	html.WriteString(AddCompoundModel.Build())
	html.WriteString(EditCompoundModel.Build())
	html.WriteString(CompoundTable.Build())
	html.WriteString(js)

	return html.String()
}
