package server

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	m "irpl.com/kanban-commons/model"
	"irpl.com/kanban-commons/utils"
	s "irpl.com/kanban-web/services"
)

func ProdProcessSearchPagination(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Pagination m.PaginationReq `json:"pagination"`
		Conditions []string        `json:"Conditions"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var Response struct {
		Pagination m.PaginationResp
		Data       []*m.ProdProcess
	}

	jsonValue, err := json.Marshal(req)
	if err != nil {
		slog.Error("Error marshaling table condition", "error", err)
	}

	// Send POST request
	resp, err := http.Post(utils.RestURL+"/get-all-production-process-by-search-paginations", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		slog.Error("Error making POST request", "error", err)
	}

	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&Response); err != nil {
		slog.Error("error decoding response body", "error", err)
	}

	// Extract DataField and Value
	searchFilters := make(map[string]string)
	for _, condition := range req.Conditions {
		parts := strings.Split(condition, " ")
		if len(parts) == 3 {
			dataField := parts[0]
			value := strings.Trim(parts[2], "'%") // Remove % from start and end
			searchFilters[dataField] = value
		}
	}

	var ProdProcessesTable s.TableCard
	tablebutton := `
	<!--html-->
			<button type="button" class="btn m-0 p-0" id="ViewProdProcessesDetails" data-toggle="tooltip" data-placement="bottom" title="View Details"> 
				<i class="fa fa-edit mx-2" style="color: #b250ad;"></i> 
			</button>
			<!--!html-->`
	ProdProcessesTable.CardHeading = "Production Process Master"
	ProdProcessesTable.CardHeadingActions.Component = []s.CardHeadActionComponent{{ComponentName: "modelbutton", ComponentType: s.ActionComponentElement{ModelButton: s.ModelButtonAttributes{ID: "add-new-prod-process", Name: "add-new-processes", Type: "button", Text: "Add New Process", ModelID: "#AddNewProcesses"}}, Width: "col-4"}}
	ProdProcessesTable.BodyTables = s.CardTableBody{Columns: []s.CardTableBodyHeadCol{
		{
			Name:       `Sr. No.`,
			IsCheckbox: true,
			ID:         "search-sn",
			Width:      "col-1",
		},
		{
			Lable:        "Processes Name",
			Name:         "Name",
			ID:           "prod-line-name",
			Width:        "col-2",
			Type:         "input",
			DataField:    "name",
			IsSearchable: true,
		},
		{
			Lable: "Status",
			Name:  "Status",
			ID:    "status",
			Width: "col-1",
		},
		{
			Lable: "Line Visibility",
			Name:  "line_visibility",
			ID:    "LineVisibility",
			Width: "col-1",
		},
		{
			Lable: "Estimated Average Time",
			Name:  "expected_mean_time",
			ID:    "ExpectedAverageTime",
			Width: "col-1",
		},
		{
			Name:  "Action",
			Width: "col-1",
		},
	},
		ColumnsWidth: []string{"col-1", "col-2", "col-1", "col-1", "col-1", "col-1"},
		Data:         Response.Data,
		Buttons:      tablebutton,
	}

	tableBodyHTML := ProdProcessesTable.BodyTables.RenderBodyColumns()

	var Pagination s.Pagination
	Pagination.TotalRecords = Response.Pagination.TotalNo
	Pagination.Offset = Response.Pagination.Offset
	Pagination.CurrentPage = Response.Pagination.Page
	Pagination.PerPage, _ = strconv.Atoi(req.Pagination.Limit)
	Pagination.PerPageOptions = []int{10, 25, 50, 100, 200}

	response := map[string]any{
		"tableBodyHTML":  tableBodyHTML,
		"paginationHTML": Pagination.Build(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}
