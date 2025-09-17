package orders

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

type Order struct {
	UserType   string
	UserID     string
	VendorCode string
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

func (o *Order) Build() string {

	// var compoundlist []m.Compounds
	type TableConditions struct {
		Conditions []string `json:"Conditions"`
	}
	var tablecondition TableConditions
	js := ""

	// type TableRequest struct {
	// 	PaginationReq m.PaginationReq `json:"pagination"`
	// 	Conditions []string `json:"conditions"`
	// }
	// PaginationReq := m.PaginationReq{
	// 	Limit:  "15",
	// 	PageNo: "1",
	// }

	// Adding a condition
	con := utils.JoinStr(`vendors.vendor_code='`, o.VendorCode, `' AND  kb_extension.status !='creating'`)
	tablecondition.Conditions = append(tablecondition.Conditions, con)

	// Create the combined request struct
	// requestPayload := TableRequest{
	// 	PaginationReq: PaginationReq,
	// 	Conditions:    tableConditions,
	// }

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
		temp.CustomerName = i.CustomerName
		tabledata = append(tabledata, &temp)
	}
	var html strings.Builder
	html.WriteString(`
		<!--html-->
		
		<!--!html-->`)

	var vendorOrderTable s.TableCard
	tabletool := `
		<!--html-->
			<button type="button" class="btn m-0 p-0" id="del-License-btn" data-toggle="tooltip" data-placement="bottom" title="View Details" > 
 					<i class="fa fa-eye mx-2" style="color: #b250ad;"></i> 
			</button>
		<!--!html-->`
	vendorOrderTable.CardHeading = "Kanban Requirement Report"
	vendorOrderTable.BodyTables = s.CardTableBody{Columns: []s.CardTableBodyHeadCol{
		{

			Name:       `Sr. No.`,
			IsCheckbox: true,
			ID:         "search-sn",
			Width:      "col-1",
		},
		{
			Name:  "Vendor Name",
			ID:    "Customer Name",
			Width: "col-1",
		},
		{
			Name:  "Customer Name",
			ID:    "Customer Name",
			Width: "col-1",
		},

		{
			Lable: "Part Name",
			Name:  "Compound Name",
			ID:    "compound_code",
			Width: "col-1",
		},
		{
			Lable: "Kanban Summary",
			Name:  "Cell No",
			ID:    "cell_no",
			Width: "col-1",
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
		ColumnsWidth: []string{"col-1", "col-1", "col-1", "col-1", "col-1", "col-1", "col-1", "col-1", "col-1"},
		Data:         tabledata,
		Tools:        tabletool,
	}
	// vendorOrderTable.BodyAction = s.CardActionBody{Component: []s.ActionComponent{
	// 	{
	// 		ComponentName: "Pagination",
	// 		ComponentType: s.ActionComponentElement{
	// 			PaginationResp: ResponseData.Pagination,
	// 		},
	// 	},
	// },
	// }
	html.WriteString(vendorOrderTable.Build())
	js = `
		<script>
			// js
			const tableBody = document.getElementById('advanced-search-table-body');
			tableBody.addEventListener('click', function (event) {
				const clickedButton = event.target.closest('button#del-License-btn'); 
				if (clickedButton) {
					const clickedRow = clickedButton.closest('tr'); 
					if (clickedRow) {
						const rowData = JSON.parse(clickedRow.getAttribute('data-data')); 
						
						const requestData = {
							ID: rowData.ID,
							KbRootId: rowData.KbRootId,
							CompoundName: rowData.CompoundName,
							CellNo: rowData.CellNo,
							DemandDateTime: rowData.DemandDateTime,
							NoOFLots: parseInt(rowData.NoOFLots, 10),
						};

						$.post("/OrderDetailsForCustomer", JSON.stringify(requestData), function(response) {
							// You can optionally handle inline response here if needed.
						}, 'json')
						.done(function(response, status, xhr) {
							try {
								// Parse the JSON response if needed
								let parsedData = typeof response === "string" ? JSON.parse(response) : response;

								// Get the content div
								const contentDiv = document.getElementById('Table-div');

								// Replace the inner HTML of the div with the response data
								contentDiv.innerHTML = parsedData.html;

								
							} catch (error) {
								console.error("Failed to parse response or update content:", error);
							}
						})

						// fetch('/OrderDetailsForCustomer', {
						// 	method: 'POST',
						// 	headers: {
						// 		'Content-Type': 'application/json',
						// 	},
						// 	body: JSON.stringify(requestData),
						// })
						// 	.then((response) => {
						// 		if (!response.ok) {
						// 			throw new Error('Network response was not ok');
						// 		}
						// 		return response.json();
						// 	})
						// 	.then((data) => {
						// 		const contentDiv = document.getElementById('Table-div');
						// 		contentDiv.innerHTML = data.html;
						// 	})
						// 	.catch((error) => {
						// 		console.error('Error fetching order details:', error); 
						// 	});
					}
				}
			});
			//!js
		</script>
		`
	html.WriteString(js)
	html.WriteString(`</div>`)
	return html.String()
}
