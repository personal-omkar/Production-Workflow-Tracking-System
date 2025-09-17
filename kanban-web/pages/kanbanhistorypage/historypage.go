package kanbanhistorypage

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"strings"

	m "irpl.com/kanban-commons/model"
	"irpl.com/kanban-commons/utils"
	s "irpl.com/kanban-web/services"
)

type Historypage struct {
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

func (h *Historypage) Build() string {
	var orderdetails []*m.OrderDetails
	var tabledata []*m.CustomerOrderDetails
	type TableConditions struct {
		PaginationReq m.PaginationReq `json:"pagination"`
		Conditions    []string        `json:"Conditions"`
	}
	var tablecondition TableConditions
	con := utils.JoinStr(``)
	tablecondition.Conditions = append(tablecondition.Conditions, con)

	jsonValue, err := json.Marshal(tablecondition)
	if err != nil {
		slog.Error("%s - error - %s", "Error marshaling table condition", err)
	}

	// GetOrderDetails retrieves all KB_root details based on the provided condition
	resp, err := http.Post(RestURL+"/get-all-kbRoot-details", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		slog.Error("%s - error - %s", "Error making GET request", err)
	}

	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&orderdetails); err != nil {
		slog.Error("error decoding response body", "error", err)
	}
	for _, value := range orderdetails {
		value.Status = strings.Title(value.Status)
	}

	for _, val := range orderdetails {
		var temp m.CustomerOrderDetails
		if val.Status == "3" || val.Status == "4" || val.Status == "-1" {
			temp.Id = val.Id
			temp.CompoundId = val.CompoundId
			temp.MFGDateTime = val.MFGDateTime.String()
			temp.DemandDateTime = val.DemandDateTime.String()
			temp.ExpDate = val.ExpDate
			temp.CellNo = val.CellNo
			temp.NoOFLots = val.NoOFLots
			temp.Location = val.Location
			temp.KbRootId = val.KbRootId
			temp.CreatedBy = val.CreatedBy
			temp.CreatedOn = val.CreatedOn
			temp.ModifiedBy = val.ModifiedBy
			temp.ModifiedOn = val.ModifiedOn
			temp.Status = map[string]string{"-1": "Quality Fail", "3": "Packing", "4": "Dispatched"}[val.Status]
			temp.VendorName = val.VendorName
			temp.CompoundName = val.CompoundName
			temp.OrderId = val.OrderId
			temp.MinQuantity = val.MinQuantity
			temp.AvailableQuantity = val.AvailableQuantity
			tabledata = append(tabledata, &temp)
		}

	}

	tablebutton := `
		<!--html-->
			<button type="button" class="btn m-0 p-0" id="viewAllCompDetails" data-toggle="tooltip" data-placement="bottom" title="View Details"> 
 					<i class="fa fa-eye mx-2" style="color: #b250ad;"></i> 
			</button>
		<!--!html-->`

	var html strings.Builder
	html.WriteString(`
		<!--html-->
		<!--!html-->`)

	var vendorOrderTable s.TableCard

	vendorOrderTable.CardHeading = "Kanban History"
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
			Name:  "Status",
			Width: "col-1",
		},
		{
			Name:  "Action",
			Width: "col-1",
		},
	},
		ColumnsWidth: []string{"col-1", "col-1", "col-1", "col-1", "col-1", "col-1", "col-1"},
		Data:         tabledata,
		Buttons:      tablebutton,
	}

	html.WriteString(vendorOrderTable.Build())

	html.WriteString(`</div>`)
	js := `
			<script>
			$(document).ready(function() {
				$('.search-input').on('input change', function() {
					
					var jsondata={
					"Criteria": {
							vendor_name: $('#search-by-VendorName').val(),
							compound_name: $('#search-by-CompoundName').val(),
							cell_no: $('#search-by-CellNo').val(),  
						}
					};

						$.post("/get-all-kbRoot-details-by-search", JSON.stringify(jsondata), function (response) {
							$('#advanced-search-table-body').html(response);
						});
					
					});
			});
		

			// js
			let isKanbanDialogLoading = false;
			document.addEventListener("click", function (event) {
			// Handle show modal functionality
			if (event.target.closest("#viewAllCompDetails")) {
				if (isKanbanDialogLoading) return;
				isKanbanDialogLoading = true;
				// Get the closest <tr> element to extract the data-data attribute
				const rowElement = event.target.closest("tr");
				if (rowElement && rowElement.dataset.data) {
					const rowData = JSON.parse(rowElement.dataset.data);
					const id = rowData.ID; // Extract the ID from the row's data

					// Make an API call to fetch the modal HTML
					fetch("/build-comp-details-dialog", {
						method: "POST",
						headers: {
							"Content-Type": "application/json",
						},
						body: JSON.stringify({ ID: id }),
					})
						.then((response) => response.json())
						.then((data) => {
							if (data.dialogHTML) {
								// Remove any existing modal with the same ID to avoid duplication
								const existingModal = document.getElementById("viewKanbanDetails");
								if (existingModal) {
									existingModal.remove();
								}

								// Append the modal directly to the body
								document.body.insertAdjacentHTML("beforeend", data.dialogHTML);

								// Show the modal (Bootstrap example)
								const modal = new bootstrap.Modal(document.getElementById("viewKanbanDetails"));
								modal.show();
							} else {
								console.error("Dialog HTML not received");
							}
						})
						.catch((error) => {
							console.error("Error fetching dialog box:", error);
						})
						.finally(() => {
							isKanbanDialogLoading = false;
						});
				} else {
					console.error("Row data not found or invalid");
				}
			}

			// Handle close modal functionality
			const modalElement = document.getElementById("viewKanbanDetails");


			// Handle custom close button
			if (event.target.id === "Close") {
				if (modalElement) {
					event.preventDefault();
					event.stopPropagation();

					// Use Bootstrap's modal method to hide the modal
					const modalInstance = bootstrap.Modal.getInstance(modalElement);
					if (modalInstance) {
						modalInstance.hide();
					}

					// Remove modal from the DOM after it's hidden
					modalElement.addEventListener("hidden.bs.modal", function () {
						modalElement.remove();
					});
				}
			}
		});
		//!js
		</script>`
	html.WriteString(js)

	return html.String()

}
