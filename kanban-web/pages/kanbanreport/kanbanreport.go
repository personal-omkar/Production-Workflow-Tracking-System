package kanbanreport

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

type KanbanReportPage struct {
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

func (h *KanbanReportPage) Build() string {
	var Response struct {
		Pagination m.PaginationResp
		Data       []*m.OrderDetails
	}
	var tabledata []*m.CustomerOrderDetails
	type TableConditions struct {
		Pagination m.PaginationReq `json:"pagination"`
		Conditions []string        `json:"conditions"`
	}

	var tablecondition TableConditions

	// Temporary pagination data
	tablecondition.Pagination = m.PaginationReq{
		Type:   "kanban",
		Limit:  "10",
		PageNo: 1,
	}

	// Temporary search conditions
	// tablecondition.Conditions = append(tablecondition.Conditions, "kr.status = 'Active'")
	// tablecondition.Conditions = append(tablecondition.Conditions, "vendors.vendor_name LIKE '%Test Vendor%'")

	// Convert to JSON
	jsonValue, err := json.Marshal(tablecondition)
	if err != nil {
		slog.Error("Error marshaling table condition", "error", err)
	}

	// Send POST request
	resp, err := http.Post(RestURL+"/get-all-kanban-details-for-report", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		slog.Error("Error making POST request", "error", err)
	}

	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&Response); err != nil {
		slog.Error("error decoding response body", "error", err)
	}
	for _, value := range Response.Data {
		value.Status = strings.Title(value.Status)
	}

	for _, val := range Response.Data {
		var temp m.CustomerOrderDetails
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
		temp.Status = map[string]string{"0": "Kanban", "1": "In Process", "2": "Quality Test", "3": "Packing", "4": "Dispatched", "-1": "Quality Test Fail"}[val.Status]
		temp.VendorName = val.VendorName
		temp.CompoundName = val.CompoundName
		temp.OrderId = val.OrderId
		temp.MinQuantity = val.MinQuantity
		temp.AvailableQuantity = val.AvailableQuantity
		temp.LotNo = val.LotNo
		tabledata = append(tabledata, &temp)
	}

	var html strings.Builder
	html.WriteString(`
		<!--html-->
		<!--!html-->`)

	var vendorOrderTable s.TableCard

	vendorOrderTable.CardHeading = "Kanban Report"
	vendorOrderTable.CardHeadingActions = s.CardHeadActionBody{
		Style: "direction: ltr;",
	}
	vendorOrderTable.CardHeadingActions.Component = []s.CardHeadActionComponent{
		{ComponentName: "input", ComponentType: s.ActionComponentElement{Input: s.InputAttributes{ID: "FromDate", Name: "FromDate", Type: "date", Icon: "From Date"}}, Width: "col-5"},
		{ComponentName: "input", ComponentType: s.ActionComponentElement{Input: s.InputAttributes{ID: "ToDate", Name: "ToDate", Type: "date", Icon: "To Date"}}, Width: "col-5"},
		{ComponentName: "button", ComponentType: s.ActionComponentElement{Button: s.ButtonAttributes{ID: "searchBtn", Colour: "#007BFF", Name: "searchBtn", Type: "button", Text: "Search"}}, Width: "col-2"},
	}
	vendorOrderTable.BodyTables = s.CardTableBody{Columns: []s.CardTableBodyHeadCol{
		{
			Name:       `Sr. No.`,
			IsCheckbox: true,
			ID:         "search-sn",
		},
		{
			Name:         "Vendor Name",
			ID:           "Customer Name",
			Width:        "col-2",
			IsSearchable: true,
			DataField:    "vendor_name",
			Type:         "input",
		},
		{
			Lable:        "Part Name",
			Name:         "Compound Name",
			ID:           "compound_code",
			Width:        "col-2",
			IsSearchable: true,
			DataField:    "compound_name",
			Type:         "input",
		},
		{
			Lable:        "Kanban Summary",
			Name:         "Cell No",
			ID:           "cell_no",
			Width:        "col-2",
			IsSearchable: true,
			DataField:    "cell_no",
			Type:         "input",
		},
		{
			Lable:        "Lot Number",
			Name:         "LotNo",
			ID:           "lot_no",
			Width:        "col-1",
			IsSearchable: true,
			DataField:    "lot_no",
			Type:         "input",
		},
		{
			Name:         "Demand Date Time",
			ID:           "demand_date",
			Width:        "col-2",
			IsSearchable: true,
			DataField:    "demand_date_time",
			Type:         "input",
		},
		{
			Name:         "Status",
			Width:        "col-2",
			IsSearchable: true,
			DataField:    "status",
			Type:         "select",
			ActionList: []s.DropDownOptions{
				{Value: "", Text: "Select Status"},
				{Value: "0", Text: "Kanban"},
				{Value: "1", Text: "In Process"},
				{Value: "2", Text: "Quality Test"},
				{Value: "3", Text: "Packing"},
				{Value: "4", Text: "Dispatched"},
				{Value: "-1", Text: "Quality Test Fail"},
			},
		},
	},
		ColumnsWidth: []string{"col-1", "col-1", "col-1", "col-1", "col-1", "col-1", "col-1"},
		Data:         tabledata,
	}

	var Pagination s.Pagination
	Pagination.CurrentPage = 1
	Pagination.Offset = 0
	Pagination.TotalRecords = Response.Pagination.TotalNo
	Pagination.PerPage, _ = strconv.Atoi(tablecondition.Pagination.Limit)
	Pagination.PerPageOptions = []int{10, 25, 50, 100, 200}
	vendorOrderTable.CardFooter = Pagination.Build()
	html.WriteString(vendorOrderTable.Build())

	html.WriteString(`</div>`)
	js := `
		<script>
			// js
		document.addEventListener("DOMContentLoaded", function () {
			document.querySelectorAll("[id^='search-by']").forEach(function (input) {
				input.addEventListener("input", function () {
					pagination(1);
				});
			});

			attachPerPageListener();

			let searchBtn = document.getElementById("searchBtn");
			if (searchBtn) {
				searchBtn.addEventListener("click", function () {
					pagination(1, true); // Always pass the flag on click
				});
			}
		});

		function pagination(pageNo, isSearchClicked = false) {
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

			// Apply date range filter only when search button is clicked
			if (isSearchClicked) {
				let fromDate = document.getElementById("FromDate")?.value;
				let toDate = document.getElementById("ToDate")?.value;

				if (fromDate && toDate) {
					searchCriteria.push("demand_date_time BETWEEN '" + fromDate + "' AND '" + toDate + "'");
				}
			}

			let requestData = {
				pagination: {
					Limit: limit,
					Pageno: pageNo
				},
				Conditions: searchCriteria
			};

			fetch("/kanban-report-search-pagination", {
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
		//!js
		</script>`
	html.WriteString(js)

	return html.String()

}
