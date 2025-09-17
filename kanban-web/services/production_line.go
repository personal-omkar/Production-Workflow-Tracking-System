package services

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"

	m "irpl.com/kanban-commons/model"
	u "irpl.com/kanban-commons/utils"
)

type ProductionLineService struct {
	restURL string
}

func NewProductionLineService() *ProductionLineService {
	restHost := os.Getenv("RESTSRV_HOST")
	restPort := os.Getenv("RESTSRV_PORT")
	if restHost == "" || restPort == "" {
		log.Println("RESTSRV_HOST or RESTSRV_PORT environment variables not set.")
		// handel this condition when host and port are not set
		restHost = "0.0.0.0"
		restPort = "4200"
	}
	restURL := fmt.Sprintf("http://%s:%s/get-production-line-items", restHost, restPort)
	return &ProductionLineService{restURL: restURL}
}

func (p *ProductionLineService) FetchProductionLines() ([]ProdLineDetailsCard, error) {
	resp, err := http.Get(p.restURL)
	if err != nil {
		log.Printf("Error fetching production lines: %v", err)
		return nil, errors.New("failed to fetch production line items")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to fetch production lines: status %s", resp.Status)
		return nil, errors.New("failed to fetch production line items")
	}

	var lines []ProdLineDetailsCard
	json.NewDecoder(resp.Body).Decode(&lines)
	return lines, nil
}

type ProdLineDetailsCard struct {
	ProdLineID   int                  `json:"prod_line_id"`
	ProdLineName string               `json:"prod_line_name"`
	Cells        []ProductionLineCell `json:"cells"`
	Style        string
}

type ProductionLineCell struct {
	CellNumber     string         `json:"cell_number"`
	KBRunningNo    string         `json:"kb_running_no"`
	KBInitialNo    string         `json:"kb_initial_no"`
	CompoundName   string         `json:"compound_name"`
	MfgDateTime    string         `json:"mfg_date_time"`
	DemandDateTime string         `json:"demand_date_time"`
	ExpDate        string         `json:"exp_date"`
	NoOFLots       int            `json:"NoOFLots"`
	Location       string         `json:"location"`
	Status         string         `json:"status"`
	KRId           string         `json:"krid"`
	ProdProcessID  string         `json:"prod_process_id"`
	LotNo          string         `json:"lot_no"`
	KanbanNo       sql.NullString `json:"kanban_no"`
	Style          string
}

var temp_no = 0

// Build generates the HTML for the ProductionLineCard
func (p *ProdLineDetailsCard) Build() string {
	sort.Slice(p.Cells, func(i, j int) bool {
		running1, _ := strconv.Atoi(p.Cells[i].KBRunningNo)
		running2, _ := strconv.Atoi(p.Cells[j].KBRunningNo)
		return running1 < running2
	})
	// Build the cells
	var cellsHTML strings.Builder
	temp_no = 0
	for _, cell := range p.Cells {
		cellsHTML.WriteString(cell.BuildProductionLine())
	}

	var c Card
	c.Body = u.JoinStr(`
		<div class="card-body p-0">
			<div class="production-line-card" id=`, strconv.Itoa(p.ProdLineID), `>
				<h5 class="line-title">`, p.ProdLineName, `</h5>
				<div style="padding: 8 8 8 8;padding-left: 8px;padding-right: 8px;">
					<button class="btn custom-button mb-2 w-100" onclick="window.location.href='/flowchart?line=`, strconv.Itoa(p.ProdLineID), `'">Live Status</button>
				</div>
				<div class="kanban-items-container" id="`, p.ProdLineName, `">`, cellsHTML.String(), `</div>
			</div>
		</div>
		<script src="/static/vendors/draggable/Sortable.min.js"></script>
		<script>
		//js
			document.addEventListener("DOMContentLoaded", function() {
				new Sortable(document.getElementById("`, p.ProdLineName, `"), {
					animation: 150,
					ghostClass: 'sortable-ghost',
					draggable: '.kanban-item',
					filter: '.non-draggable', // Prevent dragging these items
					onStart: function(evt) {
						if (evt.item.matches('.non-draggable')) {
							evt.preventDefault();
						}
					},
					onMove: function(evt) {
						if (evt.related && evt.related.matches('.non-draggable')) {
							return false;
						}
					}
				});
			});
		//!js
		</script>`)

	return c.Build()
}

// BuildProductionLine generates HTML for each cell inside the ProductionLineCard
func (c *ProductionLineCell) BuildProductionLine() string {
	if c.ProdProcessID == "1" {
		temp_no += 1
		kanbanNo := ""
		if c.KanbanNo.Valid {
			kanbanNo = c.KanbanNo.String
		}
		runningNo, _ := strconv.Atoi(c.KBRunningNo)
		NoDrag := ""
		CheckBox := ""
		if runningNo >= 1 && runningNo <= 2 {
			NoDrag = "non-draggable"
		} else {
			CheckBox = `<div class="icon-group d-flex cursor-pointer" style="padding-right:5px;">
						<input class="form-check-input" type="checkbox" value="" id="` + c.KRId + `">
						</div>`
		}
		return u.JoinStr(`
			<!--html-->
				<div class="kanban-item shadow-sm p-2 `, NoDrag, `" style="border-radius: 5px; background-color: #f9f9f9;">
					<div class="card-head">
						`, CheckBox, `
						<div class="cell-info" style="font-weight: 500;" running="`, c.KBRunningNo, `" temp_no="`, strconv.Itoa(temp_no), `">Cell Name.: `, c.CellNumber, `</div>
						<input type="hidden" class="KRid" value="`, c.KRId, `">
						<div class="icon-group d-flex cursor-pointer" style="padding-left: 5px;">
							<a class="expand text-dark" data-toggle="collapse" aria-expanded="false" aria-controls="collapseExample">
								<i class="fa fa-bars"></i>
							</a>
						</div>
					</div>
					<div class="collapse p-2" id="collapseExample" >
						<div class="row">
							<div class="col-6 row-label">Compound Code</div>
							<div class="col-6 row-label text-end">Kanban No.</div>
							<div class="col-6 row-data">`+c.CompoundName+`</div>
							<div class="col-6 row-data text-end">`+kanbanNo+`</div>
						</div>
						<div class="row">
							<div class="col-6 row-label">Demand Date/Time</div>
							<div class="col-6 row-label text-end">MFG. Date/Time</div>
							<div class="col-6 row-data">`+u.FormatStringDate(c.DemandDateTime, "date")+`</div>
							<div class="col-6 row-data text-end">`+u.FormatStringDate(c.MfgDateTime, "date-time")+`</div>
						</div>
					</div>
				</div>
			<!--!html-->
			`)
	}
	return ""
}

func AddProdLineModal() string {
	// Define the REST URL
	url := u.RestURL + "/fetch-all-production-process-data"

	// Create a new HTTP request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Print("Error while creating request: ", err)
		return "" // Exit the function if request creation fails
	}

	// Send the request using an HTTP client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Print("Error while sending request: ", err)
		return "" // Exit the function if the request fails
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to fetch data, status code: %d", resp.StatusCode)
		return ""
	}

	// Decode the response body into productionProcess slice
	var productionProcess []m.ProdProcess
	if err := json.NewDecoder(resp.Body).Decode(&productionProcess); err != nil {
		log.Print("Error while decoding response: ", err)
		return ""
	}

	var recipe []*m.Recipe

	var rawQuery m.RawQuery
	rawQuery.Host = u.RestHost
	rawQuery.Port = u.RestPort
	rawQuery.Type = "Recipe"
	rawQuery.Query = `SELECT * FROM recipe  ;` //`;`
	rawQuery.RawQry(&recipe)

	var recipelist []DropDownOptions

	for _, v := range recipe {
		var temp DropDownOptions
		temp.Text = v.CompoundName + " (" + v.CompoundCode + ")"
		temp.Value = strconv.Itoa(v.Id)
		recipelist = append(recipelist, temp)
	}
	// Construct the Checkbox slice dynamically from productionProcess
	var checkboxes []Checkbox
	for _, process := range productionProcess {
		checkboxes = append(checkboxes, Checkbox{
			ID:      fmt.Sprintf("%d", process.Id),
			Name:    process.Name,
			Width:   "w-100",
			Label:   process.Description,
			Visible: true,
		})
	}

	// Construct the AddLineModel
	AddLineModel := ModelCard{
		ID:      "AddLine",
		Type:    "modal-lg",
		Heading: "Add Line",
		Form: ModelForm{
			FormID:     "AddProdLine",
			FormAction: "",
			Footer: Footer{
				CancelBtn: false,
				Buttons: []FooterButtons{
					{BtnType: "submit", BtnID: "SaveLine", Text: "Add", Disabled: false},
				},
			},
			Inputfield: []InputAttributes{
				{Type: "text", Name: "addLineName", ID: "addLineName", Label: `Line Name`, Width: "w-100", Required: true},
				{Type: "text", Name: "addLineDescription", ID: "addLineDescription", Label: `Line Description`, Width: "w-100", Required: true},
				{Type: "file", Name: "addLineIcon", ID: "addLineIcon", Label: `Line Icon`, Width: "w-100", Hidden: true},
			},
			// Dropdownfield: []DropdownAttributes{
			// 	{Type: "text", Name: "addrecipe", ID: "addrecipe", Label: `Recipe`, Width: "w-100", Options: recipelist},
			// },
			Checkbox: checkboxes,
		},
	}

	return AddLineModel.Build()
}
