package server

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	m "irpl.com/kanban-commons/model"
	"irpl.com/kanban-commons/utils"
	s "irpl.com/kanban-web/services"
)

// Get LinedUp Production Process By LineId
func GetLinedUpProductionProcessByLineId(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	resp, err := http.Post(RestURL+"/get-lineup-processes-by-lineid", "application/json", bytes.NewBuffer(body))
	if err != nil {
		slog.Error("%s - error - %s", "Error making POST request", err)
		utils.SetResponse(w, http.StatusInternalServerError, "Failed to find record")
		return
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("%s - error - %s", "Error reading response body", err)
		utils.SetResponse(w, http.StatusInternalServerError, "Failed to read response body")
		return
	}

	utils.SetResponse(w, resp.StatusCode, string(responseBody))
}

func ProdLineSearchPagination(w http.ResponseWriter, r *http.Request) {
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
		Data       []*m.ProdLine
	}

	jsonValue, err := json.Marshal(req)
	if err != nil {
		slog.Error("Error marshaling table condition", "error", err)
	}

	// Send POST request
	resp, err := http.Post(utils.RestURL+"/get-all-prod-line-data-by-search-paginations", "application/json", bytes.NewBuffer(jsonValue))
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

	// Fetch operator list
	var operators []m.Operator
	respOp, err := http.Get(utils.RestURL + "/get-all-operator")
	if err != nil {
		slog.Error("Failed to fetch operators", "error", err)
	} else {
		defer respOp.Body.Close()
		if err := json.NewDecoder(respOp.Body).Decode(&operators); err != nil {
			slog.Error("Error decoding operators", "error", err)
		}
	}

	// Build lookup map
	opMap := make(map[string]string)
	for _, op := range operators {
		opMap[op.OperatorCode] = op.OperatorName
	}

	for _, line := range Response.Data {
		if line.Operator != "" {
			if strings.Contains(line.Operator, "(") && strings.Contains(line.Operator, ")") {
				open := strings.Index(line.Operator, "(")
				close := strings.Index(line.Operator, ")")
				line.OperatorDisplay = strings.TrimSpace(line.Operator[:open])
				line.OperatorCode = line.Operator[open+1 : close]
			} else {
				line.OperatorDisplay = opMap[line.Operator]
				line.OperatorCode = line.Operator
			}
		} else {
			line.OperatorDisplay = ""
			line.OperatorCode = ""
		}

	}

	var ProdLineTable s.TableCard
	tablebutton := `
	<!--html-->
			<button type="button" class="btn m-0 p-0" id="ViewProdLineDetails" data-toggle="tooltip" data-placement="bottom" title="View Details"> 
				<i class="fa fa-edit mx-2" style="color: #b250ad;"></i> 
			</button>
			<!--!html-->`
	ProdLineTable.CardHeading = "Machine Master"
	ProdLineTable.CardHeadingActions.Component = []s.CardHeadActionComponent{{ComponentName: "modelbutton", ComponentType: s.ActionComponentElement{ModelButton: s.ModelButtonAttributes{ID: "add-new-prod-line", Name: "add-new-compound", Type: "button", Text: "Add New Line", ModelID: "#AddLine"}}, Width: "col-3"}}
	ProdLineTable.BodyTables = s.CardTableBody{Columns: []s.CardTableBodyHeadCol{
		{
			Name:       `Sr. No.`,
			IsCheckbox: true,
			ID:         "search-sn",
			Width:      "col-1",
		},
		{
			Lable:        "Line Name",
			Name:         "Name",
			ID:           "prod-line-name",
			Width:        "col-2",
			Type:         "input",
			DataField:    "name",
			IsSearchable: true,
		},
		{
			Lable:        "Operator Name",
			Name:         "OperatorDisplay",
			ID:           "op-name",
			Width:        "col-2",
			Type:         "input",
			DataField:    "operator",
			IsSearchable: true,
		},
		{
			Lable:        "Operator Code",
			Name:         "OperatorCode",
			ID:           "op-code",
			Width:        "col-2",
			Type:         "input",
			DataField:    "operator",
			IsSearchable: true,
		},

		{
			Lable: "Status",
			Name:  "Status",
			ID:    "status",
			Width: "col-2",
		},
		{
			Name:  "Action",
			Width: "col-1",
		},
	},
		ColumnsWidth: []string{"col-1", "col-2", "col-2", "col-2", "col-2", "col-1"},

		Data:    Response.Data,
		Buttons: tablebutton,
	}

	tableBodyHTML := ProdLineTable.BodyTables.RenderBodyColumns()

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
