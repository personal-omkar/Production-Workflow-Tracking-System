package inventorymanagement

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

type InventoryManagement struct {
	Username string
	UserType string
	UserID   string
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

func (u *InventoryManagement) Build() string {
	var html strings.Builder

	var compounds []*m.Compounds
	var rawQuery m.RawQuery
	rawQuery.Host = utils.RestHost
	rawQuery.Port = utils.RestPort
	rawQuery.Type = "Compounds"
	rawQuery.Query = `SELECT * FROM compounds ;`
	rawQuery.RawQry(&compounds)

	var compoundlist []*m.Compounds
	compoundsresp, err := http.Get(RestURL + "/get-all-compound-data")
	if err != nil {
		slog.Error("%s - error - %s", "Error making GET request", err)
	}

	defer compoundsresp.Body.Close()

	if err := json.NewDecoder(compoundsresp.Body).Decode(&compoundlist); err != nil {
		slog.Error("error decoding response body", "error", err)
	}

	var Response struct {
		Pagination m.PaginationResp
		Data       []*m.ColdStorage
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

	// Convert to JSON
	jsonValue, err := json.Marshal(tablecondition)
	if err != nil {
		slog.Error("Error marshaling table condition", "error", err)
	}

	resp, err := http.Post(RestURL+"/get-all-cold-storage-by-search-pagination", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		slog.Error("%s - error - %s", "Error making GET request", err)
	}

	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&Response); err != nil {
		slog.Error("error decoding response body", "error", err)
	}

	var processflowList []s.DropDownOptions
	for _, v := range compounds {
		shouldSkip := false
		if !v.Status {
			continue
		}
		for _, cs := range Response.Data {
			if v.Id == cs.CompoundID {
				shouldSkip = true
				break
			}
		}
		if shouldSkip {
			continue
		}

		var temp s.DropDownOptions
		temp.Value = strconv.Itoa(v.Id)
		temp.Text = v.CompoundName
		processflowList = append(processflowList, temp)
	}

	var coldstorageTable s.TableCard
	tableTools := `<button type="button" class="btn  m-0 p-0" id="del-License-btn" data-bs-toggle="modal" data-bs-target="#EditCompound" > 
	<i class="fa fa-edit " style="color: #CF7AC2;"></i> 
   </button>
   <button type="button" class="btn  mx-1 p-0" id="del-part" > 
		<i class="fa fa-trash" style="color:rgb(207, 97, 97);" </i> 
   </button>`
	coldstorageTable.CardHeading = "Rubber Store Master"
	coldstorageTable.CardHeadingActions.Component = []s.CardHeadActionComponent{
		{ComponentName: "modelbutton", ComponentType: s.ActionComponentElement{ModelButton: s.ModelButtonAttributes{Type: "button", Text: "Add Part to Inventory", ModelID: "#AddParToInventoryModel"}}, Width: "col-4"}}

	coldstorageTable.BodyTables = s.CardTableBody{Columns: []s.CardTableBodyHeadCol{
		{
			Lable:        "Part Name",
			Name:         `Compound Name`,
			ID:           "search-compound-name",
			DataField:    "compound_name",
			IsSearchable: true,
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
	coldstorageTable.CardFooter = Pagination.Build()

	// Create Dialogue
	AddPartModel := s.ModelCard{
		ID:      "AddParToInventoryModel",
		Type:    "modal-m",
		Heading: "Update Inventory (Add Part)",
		Form: s.ModelForm{FormID: "AddParToInventoryModel",
			FormAction: "", Footer: s.Footer{CancelBtn: true, Buttons: []s.FooterButtons{{BtnType: "button", DataSubmitName: "AddPartModel", BtnID: "add-part-submit", Text: "Save"}}},
			Dropdownfield: []s.DropdownAttributes{
				{Label: "Part", Name: "CompoundId", DataType: "int", Options: processflowList},
			},
			Inputfield: []s.InputAttributes{
				{Type: "number", Name: "MaxQuantity", DataType: "int", ID: "maxQuantity", Label: `Maximum Quantity`, Width: "w-100", Required: true},
				{Type: "number", Name: "MinQuantity", DataType: "int", ID: "minQuantity", Label: `Minimum Quantity`, Width: "w-100", Required: true},
			},
		},
	}

	EditPartModel := s.ModelCard{
		ID:      "EditCompound",
		Type:    "modal-md",
		Heading: `Edit Record  ( <span id="element-txt"></span> ) `,
		Form: s.ModelForm{FormID: "Inventory-Compound",
			FormAction: "/update-coldstorage-quantity-for-inventory-managment", Footer: s.Footer{CancelBtn: true, Buttons: []s.FooterButtons{{BtnType: "submit", BtnID: "edit-compound", Text: "Update"}}},
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

			fetch("/coldstorage-pagination-search?showDelete=true", {
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
    document.addEventListener("DOMContentLoaded", function () {

		document.querySelectorAll("[id^='search-by']").forEach(function (input) {
			input.addEventListener("input", function () {
				pagination(1);
			});
		});
		
		attachPerPageListener();

		const form = document.getElementById("AddParToInventoryModel");
		const saveButton = document.getElementById("add-part-submit");
		const maxQuantityInput = document.getElementById("maxQuantity");
		const minQuantityInput = document.getElementById("minQuantity"); 
		const partSelect = form.querySelector('select[name="CompoundId"]');

		// Function to check form validity
		function validateForm() {
			const maxQty = maxQuantityInput.value.trim();
			const minQty = minQuantityInput.value.trim();
			const partValue = partSelect.value.trim();

			if (maxQty !== "" && minQty !== "" && partValue !== "") {
				saveButton.removeAttribute("disabled");
			} else {
				saveButton.setAttribute("disabled", "true");
			}
		}

		// Add event listeners for input validation
		[maxQuantityInput, minQuantityInput, partSelect].forEach(input => {
			input.addEventListener("input", validateForm);
		});

		// Handle form submission
		saveButton.addEventListener("click", function () {
			if (saveButton.disabled) return; // Prevent click if still disabled

			const formData = {
				MaxQuantity: parseInt(maxQuantityInput.value, 10),
				MinQuantity: parseInt(minQuantityInput.value, 10),
				CompoundId: parseInt(partSelect.value, 10)
			};


			$.post("/create-new-or-update-existing-inventory", JSON.stringify(formData), function () {}, "json").fail(function (xhr) {
				window.location.href = "/inventory-management?status=" + xhr.status + "&msg=" + xhr.responseText;
			});
		});

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
		$("#EditCompound").on("show.bs.modal" ,function(event){   		
			
			var data = JSON.parse($(event.relatedTarget).closest("tr").attr("data-data"))	
			$.get("/get-compound-data-by-parm?key=id&value=" + data.CompoundID, function(compdata) {
			document.getElementById('editinventoryid').value = data.ID;
			document.getElementById('showcompoundname').value = data.CompoundName;
			document.getElementById('editcompoundname').value = data.CompoundName;
			document.getElementById('editminqty').value = data.MinQuantity; 
			document.getElementById('editmaxqty').value = data.MaxQuantity;
			document.getElementById('editavailableqty').value = data.AvailableQuantity; 
			$('#element-txt').text( data.CompoundName) 			
			});
		});
		
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


		$(document).on("click", "#del-part", function () {
			var rowData = $(this).closest("tr").attr("data-data");
			if (!rowData) {
				console.error("No data found in the closest row.");
				return;
			}
			var data = JSON.parse(rowData);
			var partId = data.ID;
			if (!partId) {
				console.error("No valid ID found for deletion.");
				return;
			}
			$.post("/delete-inventory-by-id", JSON.stringify({ ID: partId }), function (response) {
				window.location.href = "/inventory-management?status=200&msg=Part deleted successfully";
			}, "json").fail(function (xhr) {
				window.location.href = "/inventory-management?status=" + xhr.status + "&msg=" + xhr.responseText;
			});
		});

		// Initial validation check
		validateForm();
	});

	//!js
	</script>
	`

	html.WriteString(coldstorageTable.Build())
	html.WriteString(AddPartModel.Build())
	html.WriteString(EditPartModel.Build())
	html.WriteString(`</div>`)
	html.WriteString(js)
	return html.String()

}
