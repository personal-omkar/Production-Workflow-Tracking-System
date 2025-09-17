package vendormanagement

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

type VendorManagement struct {
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
func (u *VendorManagement) Build() string {
	var Response struct {
		Pagination m.PaginationResp
		Data       []*m.Vendors
	}
	type TableConditions struct {
		Pagination m.PaginationReq `json:"pagination"`
		Conditions []string        `json:"conditions"`
	}
	var tablecondition TableConditions

	// Temporary pagination data
	tablecondition.Pagination = m.PaginationReq{
		Type:   "vendor",
		Limit:  "10",
		PageNo: 1,
	}

	// Convert to JSON
	jsonValue, err := json.Marshal(tablecondition)
	if err != nil {
		slog.Error("Error marshaling table condition", "error", err)
	}

	resp, err := http.Post(RestURL+"/get-all-vendor-by-search-pagination", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		slog.Error("%s - error - %s", "Error making GET request", err)
	}

	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&Response); err != nil {
		slog.Error("error decoding response body", "error", err)
	}

	var vendorstatus = []s.DropDownOptions{{Text: "False", Value: "false"}, {Text: "True", Value: "true"}}
	tableTools := `<button  type="button" class="btn  m-0 p-0" id="edit-vendor-btn" data-bs-toggle="modal" data-bs-target="#EditVendorModel"> 
 					<i class="fa fa-edit " style="color: #871a83;"></i> 
					</button>`
	var html strings.Builder
	html.WriteString(`
		<!--html-->
		<!--!html-->`)

	var userManagement s.TableCard

	userManagement.CardHeading = "Vendor Master"
	userManagement.CardHeadingActions.Component = []s.CardHeadActionComponent{
		{
			ComponentName: "modelbutton",
			ComponentType: s.ActionComponentElement{
				ModelButton: s.ModelButtonAttributes{ID: "addNewVendor", Name: "addNewVendor", Type: "button", Text: "Add New Vendor", ModelID: "#AddNewVendor"},
			},
			Width: "col-4",
		},
		{
			ComponentName: "button",
			ComponentType: s.ActionComponentElement{
				Button: s.ButtonAttributes{ID: "importVendors", Name: "importVendors", Type: "button", Text: `Import<input style="display:none" class="clickable btn btn-success" type="file" id="vm-import-file" name="file" accept=".xlsx">`, Class: "bg-transparent", CSS: "color: #871a83; border: 1px solid #871a83;"},
			},
			Width: "col-3",
		},
		{
			ComponentName: "button",
			ComponentType: s.ActionComponentElement{
				Button: s.ButtonAttributes{ID: "downloadTmp", Name: "downloadTmp", Type: "button", Text: ` Template <i class="fa fa-download" aria-hidden="true"></i> `, Class: "bg-transparent", AdditionalAttr: `data-link="/files/sample.pdf"`, CSS: "color: #871a83; border: 1px solid #871a83;", IsDownloadBtn: true, DownloadLink: `/static/templates/import-vendor-master.xlsx`},
			},
			Width: "col-3",
		},
	}
	userManagement.BodyTables = s.CardTableBody{Columns: []s.CardTableBodyHeadCol{

		{
			Name:         "Vendor Code",
			ID:           "vendorCode",
			Width:        "col-2",
			IsSearchable: true,
			DataField:    "vendor_code",
			Type:         "input",
		},
		{
			Name:         "Vendor Name",
			ID:           "vendorName",
			Width:        "col-2",
			IsSearchable: true,
			DataField:    "vendor_name",
			Type:         "input",
		},
		{
			Name:  "Contact Info",
			ID:    "contactInfo",
			Width: "col-2",
		},
		{
			Name:  "Created On",
			ID:    "createdOn",
			Width: "col-2",
		},
		{
			Name:  "Tools",
			Width: "col-1",
		},
	},
		ColumnsWidth: []string{"col-2", "col-2", "col-2", "col-2", "col-1"},
		Data:         Response.Data,
		Tools:        tableTools,
		ID:           "VendorManagement",
	}
	var Pagination s.Pagination
	Pagination.CurrentPage = 1
	Pagination.Offset = 0
	Pagination.TotalRecords = Response.Pagination.TotalNo
	Pagination.PerPage, _ = strconv.Atoi(tablecondition.Pagination.Limit)
	Pagination.PerPageOptions = []int{10, 25, 50, 100, 200}
	userManagement.CardFooter = Pagination.Build()

	//Add Dialogue
	AddVendorModel := s.ModelCard{
		ID:      "AddNewVendor",
		Type:    "modal-md",
		Heading: "Add Vendor",
		Form: s.ModelForm{FormID: "AddNewVendor",
			FormAction: "", Footer: s.Footer{CancelBtn: false, Buttons: []s.FooterButtons{{BtnType: "button", DataSubmitName: "AddNewVendor", BtnID: "add-new-vendor-submit", Text: "Add"}}},
			Inputfield: []s.InputAttributes{
				{Type: "text", Name: "VendorCode", ID: "addVendorCode", Label: `Vendor Code`, Width: "w-100"},
				{Type: "text", Name: "VendorName", ID: "addVendorName", Label: `Vendor Name`, Width: "w-100"},
				{Type: "text", Name: "ContactInfo", ID: "addContactInfo", Label: `Contact Info`, Width: "w-100"},
				{Type: "text", DataType: "int", Name: "PerDayLotConfig", ID: "addPerDayLotConfig", Label: `Lot Per Day`, Width: "w-50"},
				{Type: "text", DataType: "int", Name: "PerMonthLotConfig", ID: "addPerMonthLotConfig", Label: `Lot Per Month`, Width: "w-50"},
				{Type: "text", DataType: "int", Name: "PerHourLotConfig", ID: "addPerHourLotConfig", Label: `Lot Per Hour`, Width: "w-100"},
			},
			TextArea: []s.TextAreaAttributes{
				{Label: "Address", Name: "Address", ID: "addAddress"},
			},
		},
	}

	//Edit Dialogue
	EditVendorModel := s.ModelCard{
		ID:      "EditVendorModel",
		Type:    "modal-md",
		Heading: "Edit Vendor",
		Form: s.ModelForm{FormID: "EditVendorModel",
			FormAction: "", Footer: s.Footer{CancelBtn: true, Buttons: []s.FooterButtons{{BtnType: "button", DataSubmitName: "EditVendorModel", BtnID: "edit-vendor-submit", Text: "Save"}}},
			Inputfield: []s.InputAttributes{
				{Type: "hidden", DataType: "int", Name: "ID", ID: "showUserId", Label: `User ID`, Width: "w-100", Hidden: true},
				{Type: "text", Name: "VendorCode", ID: "showVendorCode", Label: `Vendor Code`, Width: "w-100", Disabled: true},
				{Type: "text", Name: "VendorName", ID: "showVendorName", Label: `Vendor Name`, Width: "w-100"},
				{Type: "text", Name: "ContactInfo", ID: "showContactInfo", Label: `Contact Info`, Width: "w-100"},
				{Type: "text", DataType: "int", Name: "PerDayLotConfig", ID: "showPerDayLotConfig", Label: `Lot Per Day`, Width: "w-50"},
				{Type: "text", DataType: "int", Name: "PerMonthLotConfig", ID: "showPerMonthLotConfig", Label: `Lot Per Month`, Width: "w-50"},
				{Type: "text", DataType: "int", Name: "PerHourLotConfig", ID: "showPerHourLotConfig", Label: `Lot Per Hour`, Width: "w-100"},
			},
			Dropdownfield: []s.DropdownAttributes{
				{Label: "Is Active", DataType: "text", Name: "Isactive", Options: vendorstatus, ID: "showIsActive", Width: "w-100"},
			},
			TextArea: []s.TextAreaAttributes{
				{Label: "Address", Name: "Address", ID: "showaddAddress"},
			},
		},
	}

	js :=
		`<script>
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

		fetch("/vendor-search-pagination", {
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
	});

	const urlParams = new URLSearchParams(window.location.search);
	const status = urlParams.get("status");
	const msg = urlParams.get("msg");

	if (status) {
		showNotification(status, msg, () => {
			removeQueryParams();
		});
	}

	$(document).on("click", "#edit-vendor-btn", function () {
		let data = JSON.parse($(this).closest("tr").attr("data-data"));
		$("#showUserId").val(data.ID);
		$("#showVendorCode").val(data.VendorCode);
		$("#showVendorName").val(data.VendorName);
		$("#showContactInfo").val(data.ContactInfo);
		$("#showaddAddress").val(data.Address);
		$("#showPerDayLotConfig").val(data.PerDayLotConfig);
		$("#showPerMonthLotConfig").val(data.PerMonthLotConfig);
		$('#showPerHourLotConfig').val(data.PerHourLotConfig);
		$("#showIsActive").val(data.Isactive.toString());
	});

	$("#showPerDayLotConfig,#showPerMonthLotConfig,#addPerDayLotConfig,#addPerMonthLotConfig,#addPerHourLotConfig,#showPerHourLotConfig").on("input", function () {
		this.value = this.value.replace(/[^1-9\d]/g, "").replace(/^0+/, "");
	});

	$(document).on("click", "#add-new-vendor-submit,#edit-vendor-submit", function () {
		let group = $(this).attr("data-submit");
		let result = {};
		let validated = true;

		$("[data-group='" + group + "']").find("[data-name]").each(function () {
			if (
				($(this).attr("data-validate") && $(this).val().trim().length === 0) ||
				($(this).attr("data-validate") && $(this).val() === "Nil")
			) {
				$(this).css("background-color", "rgba(128, 0, 128, 0.1)");
				let label = $(this).closest("label").length
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
						"margin-left": "0.5rem"
					}).insertAfter(label);
				}
				validated = false;
			} else {
				$(this).css("background-color", "rgb(255, 255, 255)");
				let label = $(this).closest("label").length
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

		if (validated) {
			$("[data-group='" + group + "']").find("[data-name]").each(function () {
				if ($(this).is("select")) {
					result[$(this).attr("data-name")] = $(this).find(":selected").val();
				} else if ($(this).attr("data-type") === "date") {
					let userDate = $(this).val();
					if (userDate.includes(":")) {
						result[$(this).attr("data-name")] = userDate;
					} else {
						let formattedDate = formatDateToYYYYMMDD(userDate);
						result[$(this).attr("data-name")] = new Date(formattedDate);
					}
				} else if ($(this).attr("data-type") === "int") {
					result[$(this).attr("data-name")] = parseInt($(this).val());
				} else {
					result[$(this).attr("data-name")] = $(this).val();
				}
			});

			if (group === "AddNewVendor") {
				result["Isactive"] = true;
			} else {
				result["Isactive"] = result["Isactive"] === "true";
			}

			$.post("/create-new-vendor", JSON.stringify(result), function () {}, "json").fail(function (xhr) {
				window.location.href = "/vendor-management?status=" + xhr.status + "&msg=" + xhr.responseText;
			});
		}
	});

	function showNotification(status, msg, callback) {
		const notification = $("#notification");
		let message = "";

		if (status === "200") {
			message = "<strong>Success!</strong> " + msg + ".";
			notification.removeClass("alert-danger").addClass("alert-success");
		} else {
			message = "<strong>Fail!</strong> " + msg + ".";
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
		let newUrl = window.location.protocol + "//" + window.location.host + window.location.pathname;
		window.history.replaceState({}, document.title, newUrl);
	}

	$(document).on("click", "#importVendors", function(event) {
				
				switch (event.target.id) {
					case 'importVendors':
					$('#vm-import-file').trigger('click');
						break;
				}
			});
			$(document).on("change", '#vm-import-file', function(e) {
						const url = "/import-vendor-data"
						if (url) {
							
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
									success: function(data) {
										window.location.href = "/vendor-management?status=200&msg=File uploaded successfully"
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
	html.WriteString(userManagement.Build())
	html.WriteString(AddVendorModel.Build())
	html.WriteString(EditVendorModel.Build())
	html.WriteString(js)

	html.WriteString(`</div>`)

	return html.String()

}
