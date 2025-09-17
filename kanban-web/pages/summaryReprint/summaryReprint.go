package summaryReprint

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"

	m "irpl.com/kanban-commons/model"
	"irpl.com/kanban-commons/utils"
	s "irpl.com/kanban-web/services"
)

type SummaryReprintPage struct {
	Username   string
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

func (k *SummaryReprintPage) Build() string {
	// var compoundlist []m.Compounds
	type TableConditions struct {
		Conditions []string `json:"Conditions"`
	}
	var con string
	var tablecondition TableConditions
	js := ""

	// Adding a condition
	if k.VendorCode != "" {
		con = utils.JoinStr(`vendors.vendor_code='`, k.VendorCode, `' AND  kb_extension.status !='creating'  AND TO_CHAR(kb_extension.created_on, 'YYYY-MM-DD') = TO_CHAR(NOW(), 'YYYY-MM-DD')`)
	} else {
		con = utils.JoinStr(`kb_extension.status !='creating' AND TO_CHAR(kb_extension.created_on, 'YYYY-MM-DD') = TO_CHAR(NOW(), 'YYYY-MM-DD') `)
	}

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
		// if o.UserID == i.CreatedBy {
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
		// }
	}
	var html strings.Builder
	html.WriteString(`
		<!--html-->
		
		<!--!html-->`)

	var kanbanreprint s.TableCard

	allCheckBox := `
		<input style="width: 20px; height: 20px;" class="form-check-input me-2 selected_item" type="checkbox" value="">`
	kanbanreprint.CardHeading = "Summary Reprint"
	kanbanreprint.CardHeadingActions.Component = []s.CardHeadActionComponent{
		{ComponentName: "button", ComponentType: s.ActionComponentElement{Button: s.ButtonAttributes{ID: "summaryPrint", Name: "summaryPrint", Type: "button", Text: "Summary Print", Disabled: true}}, Width: "col-4"},
	}

	kanbanreprint.BodyTables = s.CardTableBody{Columns: []s.CardTableBodyHeadCol{
		{
			Name:  "Tools",
			Lable: " ",
			Width: "col-1",
		},
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
			Lable: "Part Name",
			Name:  "Compound Name",
			ID:    "compound_code",
			Width: "col-1",
		},
		{
			Name:  "Customer Name",
			ID:    "Customer Name",
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
	},
		ColumnsWidth: []string{"col-1", "col-1", "col-1", "col-1", "col-1", "col-1", "col-1", "col-1", "col-1"},
		Data:         tabledata,
		Tools:        allCheckBox,
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
	html.WriteString(kanbanreprint.Build())
	js = `
		<script>
			// js

			document.addEventListener("DOMContentLoaded", function () {
				document.querySelectorAll(".selected_item").forEach(function (checkbox) {
					checkbox.checked = false;
				});
				function showNotification(status, msg, callback) {
					var notification = $('#notification');
					var message = "";
					if (status === "200") {
						message = "Success! " + msg + ".";
						notification.removeClass("alert-danger").addClass("alert-success");
					} else {
						message = "Fail! " + msg + ".";
						notification.removeClass("alert-success").addClass("alert-danger");
					}
					notification.html(message);
					notification.show();

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

				var urlParams = new URLSearchParams(window.location.search);
				var status = urlParams.get("status");
				var msg = urlParams.get("msg");
				if (status) {
					showNotification(status, msg, function () {
						removeQueryParams();
					});
				}
			});

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
						
					}
				}
			});

				$(document).on('click', '#kanbanPrint', function () {
					
					const selectedData = []; 
							$(".selected_item:checked").each(function () {
								const row = $(this).closest('tr');
								const datadata = row.attr('data-data');

								if (!datadata) return;

								try {
									const parsedData = JSON.parse(datadata);
									selectedData.push(parsedData.CellNo); // Add to array
								} catch (error) {
									console.error("Error parsing JSON:", error);
								}
							});

							
						for (let index = 0; index < selectedData.length; index++) {
							$.get("/cust/print-html-card?cellno=" + selectedData[index], function(data) {
								$("#report-div").html(data);
								printContent();
							});		
						}
				});

			$(document).on('click', '#summaryPrint', function () {
				var cellnos = [];

				$(".selected_item:checked").each(function () {
					const jsonStr = $(this).closest("tr").attr("data-data");
					const data = JSON.parse(jsonStr);
					cellnos.push(data.CellNo);  // <-- correct method to add to array
				});

				var requestData = {
					Type: "Summary",
					CellNo: cellnos
				};


				$.post("/print-kanban-summary", JSON.stringify(requestData), function(response) {
					$("#report-div").html(response)
					$("#report-dialog").modal("show")
				}).fail(function(xhr, status, error) {
					window.location.href = "/kanban-reprint?status=" + xhr.status + "&msg=" + xhr.responseText;
				});
			});
			

			$(document).on('change', '.selected_item', function () {
				if ($('.selected_item:checked').length > 0) {
					$('#kanbanPrint').prop('disabled', false); 
					$('#summaryPrint').prop('disabled', false); 
				} else {
					$('#summaryPrint').prop('disabled', true); 
					$('#kanbanPrint').prop('disabled', true); 
				}
			});

			$('#report-dialog').on('hidde.bs.modal', function () {
				window.location.href = "/kanban-reprint"
			});

			//!js
		</script>
		`
	html.WriteString(js)
	html.WriteString(`</div>`)
	html.WriteString(`<div id="report-div"></div>`)

	return html.String()

}
