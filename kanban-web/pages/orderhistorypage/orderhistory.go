package orderhistorypage

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"irpl.com/kanban-commons/utils"

	m "irpl.com/kanban-commons/model"
	s "irpl.com/kanban-web/services"
)

type OrderHistory struct {
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

func (o *OrderHistory) Build() string {

	// var compoundlist []m.Compounds
	type TableConditions struct {
		Conditions []string `json:"Conditions"`
	}
	var tablecondition TableConditions

	// Adding a condition
	con := "(kb_extension.status = 'dispatch' OR kb_extension.status LIKE 'dispatch%') " +
		"AND kb_extension.vendor_id != (" +
		"SELECT id FROM vendors " +
		"WHERE vendor_code LIKE 'I%' " +
		"LIMIT 1 )" //";"
	tablecondition.Conditions = append(tablecondition.Conditions, con)

	jsonValue, err := json.Marshal(tablecondition)
	if err != nil {
		slog.Error("%s - error - %s", "Error marshaling request payload", err)
	}
	url := RestURL + "/get-customer-order-details"
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		slog.Error("%s - error - %s", "Error making GET request", err)
	}

	defer resp.Body.Close()
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("%s - error - %s", "Error reading response body", err)
	}

	var ResponseData []*m.OrderDetails // Ensure correct field name in struct
	var tabledata []*m.CustomerOrderDetails

	if err := json.NewDecoder(bytes.NewReader(responseBody)).Decode(&ResponseData); err != nil {
		slog.Error("error decoding response body", "error", err)
	}
	for _, i := range ResponseData {
		var temp m.CustomerOrderDetails
		temp.Id = i.Id
		temp.KbRootId = i.Id
		temp.VendorName = i.VendorName
		temp.CompoundName = i.CompoundName
		temp.CellNo = i.CellNo
		temp.NoOFLots = i.NoOFLots
		temp.Status = strings.Title(i.Status)
		temp.DemandDateTime = i.DemandDateTime.String()
		temp.MFGDateTime = i.MFGDateTime.String()
		tabledata = append(tabledata, &temp)
	}
	var html strings.Builder
	html.WriteString(`
		<!--html-->
		
		<!--!html-->`)

	var vendorOrderTable s.TableCard
	tabletool := `
		<!--html-->
			<button type="button" class="btn m-0 p-0" id="viewAllOrderDetails" data-toggle="tooltip" data-placement="bottom" title="View Details"> 
 					<i class="fa fa-eye mx-2" style="color: #b250ad;"></i> 
			</button>
		<!--!html-->`
	vendorOrderTable.CardHeading = "Orders History"
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
		},
		{
			Lable:        "Kanban Summary",
			Name:         "Cell No",
			IsSearchable: true,
			Type:         "input",
			ID:           "CellNo",
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
			Name:  "Tools",
			Width: "col-1",
		},
	},
		ColumnsWidth: []string{"col-1", "col-1", "col-1", "col-1", "col-1", "col-1", "col-1"},
		Data:         tabledata,
		Tools:        tabletool,
	}

	html.WriteString(vendorOrderTable.Build())
	js := `
		<script>
		// js
			$(document).ready(function() {
				$('.search-input').on('input change', function() {
					
					var jsondata={
					"Criteria": {
							vendor_name: $('#search-by-VendorName').val(),
							compound_name: $('#search-by-CompoundName').val(),
							cell_no: $('#search-by-CellNo').val(),  
						}
					};

						$.post("/search-customer-order-details", JSON.stringify(jsondata), function (response) {
							$('#advanced-search-table-body').html(response);
						});
					
					});
			});
		let isOrderDialogLoading = false;
		let isKanbanDialogLoading = false;

		document.addEventListener("click", function (event) {
			// Handle Dialog 1 (viewOrderDetails)
			if (event.target.closest("#viewAllOrderDetails")) {
				if (isOrderDialogLoading) return;
				isOrderDialogLoading = true;

				const rowElement = event.target.closest("tr");
				if (rowElement && rowElement.dataset.data) {
					const rowData = JSON.parse(rowElement.dataset.data);
					const id = rowData.ID; 

					// Remove existing modals and backdrops
					document.querySelectorAll(".modal").forEach(modal => modal.remove());
					document.querySelectorAll(".modal-backdrop").forEach(backdrop => backdrop.remove());

					fetch("/build-order-details-dialog", {
						method: "POST",
						headers: { "Content-Type": "application/json" },
						body: JSON.stringify({ ID: id }),
					})
					.then(response => response.json())
					.then(data => {
						if (data.dialogHTML) {
							let tempDiv = document.createElement("div");
							tempDiv.innerHTML = data.dialogHTML;
							document.body.append(...tempDiv.children);

							const modalElement = document.getElementById("viewOrderDetails");

							if (modalElement) {
								const modal = new bootstrap.Modal(modalElement);
								modal.show();
								modalElement.removeAttribute("aria-hidden");
								modalElement.focus();
							} else {
								console.error("#viewOrderDetails not found.");
							}
						} else {
							console.error("Dialog 1 HTML not received.");
						}
					})
					.catch(error => console.error("Error fetching Dialog 1:", error))
					.finally(() => {
						isOrderDialogLoading = false;
					});
				}
			}

			// Handle Dialog 2 (viewKanbanDetails)
			if (event.target.closest("#viewKanbanDetails")) {
				if (isKanbanDialogLoading) return;
				isKanbanDialogLoading = true;

				const rowElement = event.target.closest("tr");
				if (rowElement && rowElement.dataset.data) {
					const rowData = JSON.parse(rowElement.dataset.data);
					const id = rowData.ID;

					const dialog1 = document.getElementById("viewOrderDetails");
					document.querySelectorAll(".modal").forEach(modal => modal.remove());

					if (dialog1) {
						const modal1 = bootstrap.Modal.getInstance(dialog1);
						if (modal1) {
							dialog1.blur();
							modal1.hide();
						}
					}

					fetch("/build-comp-details-dialog", {
						method: "POST",
						headers: { "Content-Type": "application/json" },
						body: JSON.stringify({ ID: id }),
					})
					.then(response => response.json())
					.then(data => {
						if (data.dialogHTML) {
							let tempDiv = document.createElement("div");
							tempDiv.innerHTML = data.dialogHTML;
							document.body.append(...tempDiv.children);

							const modalElement = document.getElementById("viewKanbanDetails");

							if (modalElement) {
								const modal = new bootstrap.Modal(modalElement);
								modal.show();

								modalElement.addEventListener("click", function(event) {
									event.stopPropagation();
								});
							} else {
								console.error("#viewKanbanDetails not found.");
							}
						} else {
							console.error("Dialog 2 HTML not received.");
						}
					})
					.catch(error => console.error("Error fetching Dialog 2:", error))
					.finally(() => {
						isKanbanDialogLoading = false;
					});
				}
			}
		});
		//!js
		</script>`
	html.WriteString(js)
	html.WriteString(`</div>`)
	return html.String()
}
