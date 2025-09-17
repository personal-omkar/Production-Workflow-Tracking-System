package allkanbanview

import (
	"os"
	"strconv"
	"strings"

	"irpl.com/kanban-commons/model"
	"irpl.com/kanban-commons/utils"
	s "irpl.com/kanban-web/services"
)

type AllKanbanView struct {
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

func (u *AllKanbanView) Build() string {
	var html strings.Builder

	var ProdLine []*model.ProdLine
	var details []*model.AllKanbanViewDetails
	var AllKanbanViewTable []*model.AllKanbanViewTable
	var rawQuery model.RawQuery
	rawQuery.Host = utils.RestHost
	rawQuery.Port = utils.RestPort

	rawQuery.Type = "ProdLine"
	rawQuery.Query = `select * from prod_line where status='true';`
	rawQuery.RawQry(&ProdLine)
	for i := range ProdLine {
		ProdLine[i].Name = strings.TrimSpace(strings.TrimSuffix(ProdLine[i].Name, "(Heijunka)"))
	}
	for _, v := range ProdLine {
		rawQuery.Type = "AllKanbanViewDetails"
		rawQuery.Query = `
		SELECT
			kr.id,
			c.compound_name AS compound_name,
			kr.kanban_no AS kanban_no,
			kd.note AS notes,
			kd.cell_no as cell_number
		FROM
			kb_transaction kt
		JOIN
			kb_root kr ON kt.kb_root_id = kr.id
		JOIN
			kb_data kd ON kr.kb_data_id = kd.id
		JOIN
			kb_extension ke ON kd.kb_extension_id = ke.id
		JOIN
			compounds c ON kd.compound_id = c.id
		JOIN
			prod_process_line ppl ON kt.prod_process_line_id = ppl.id
		JOIN
			prod_line pl ON ppl.prod_line_id = pl.id
		WHERE
			pl.status = true
		AND
			pl.id = ` + strconv.Itoa(v.Id) + `
		AND 
			ppl.prod_process_id IN (1)
		AND (
			SELECT COUNT(*)
			FROM kb_transaction sub_kt
			WHERE sub_kt.kb_root_id = kt.kb_root_id
			) <=(
			        SELECT COUNT(*)          
			        FROM prod_process_line
			        WHERE prod_line_id =` + strconv.Itoa(v.Id) + `
			      )
		AND NOT EXISTS (
			SELECT 1
			FROM kb_transaction kt2
			JOIN prod_process_line ppl2 ON kt2.prod_process_line_id = ppl2.id
			WHERE kt2.kb_root_id = kt.kb_root_id
			AND ppl2.prod_process_id = 2
		)
		GROUP BY
			kr.id, c.compound_name, kr.kanban_no, kd.note, kd.cell_no
		HAVING COUNT(kt.kb_root_id) = 1
		ORDER BY
			kr.running_no ASC;
		
		` //`;`

		rawQuery.RawQry(&details)

		var allcards string
		if len(details) > 0 {
			for _, v := range details {
				card := `
				<div class="d-flex border border-1 p-0 mx-2 my-2" style="background-color:` + utils.KanbanPriorityColors[strings.ToLower(v.CustomerNotes)]["bg-color"] + `; color:` + utils.KanbanPriorityColors[strings.ToLower(v.CustomerNotes)]["text-color"] + `; border-color: #ab71a2 !important; border-radius: 6px; user-select: none; cursor: pointer;">
						<span class="form-check p-0 m-0" style="cursor: pointer;">
							<label class="form-check-label m-1 px-1" for="290">
								` + v.CompoundName + `
							</label>
							<span class="mx">|</span>
							<label class="form-check-label m-0 px-1" for="290" data="Cell-Name">
							` + v.KanbanNo + `
							</label>	
					</span>
				</div>`
				allcards = allcards + card
			}
			v.Name = `<a href="/production-line">` + v.Name + `</a> <span style="margin-left:10%">   (` + strconv.Itoa(len(details)) + ") </span>"
		} else {
			v.Name = `<a href="/production-line">` + v.Name + `</a>`
		}

		var temptable model.AllKanbanViewTable
		temptable.MachineId = strconv.Itoa(v.Id)
		temptable.MachineName = v.Name
		temptable.PartNameorKanbanNo = allcards

		AllKanbanViewTable = append(AllKanbanViewTable, &temptable)
	}

	var AllKanbanView s.TableCard

	AllKanbanView.CardHeading = "All Kanban View"

	AllKanbanView.BodyTables = s.CardTableBody{Columns: []s.CardTableBodyHeadCol{
		{
			Lable:        "Machine Name",
			Name:         `Machine Name`,
			ID:           "machine-name",
			IsSearchable: true,
			Type:         "input",
			Width:        "col-2",
		},
		{
			Lable:        "Part Name or Kanban Number",
			Name:         "Part Name or Kanban No",
			Type:         "input",
			IsSearchable: true,
			ID:           "part-name-or-kanban-no",
			Width:        "col-10",
			Style:        "max-width:1300px",
		},
	},
		ColumnsWidth: []string{"col-2", "col-10 w-100 d-flex flex-wrap"},
		Data:         AllKanbanViewTable,
		ID:           "All Kanban View",
	}

	js := `
	<script>
	//js
   $(document).ready(function() {
		
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
			
		function refreshData() {
			var jsondata={
				"Criteria": {
							machine_name: $('#search-by-MachineName').val(),
							part_name_or_kanabn_no: $('#search-by-PartNameorKanbanNo').val(),  
							}
						};
			
						$.post('/search-kanban-view', JSON.stringify(jsondata), function (response) {
							$('#advanced-search-table-body').html(response);
						});

			}	
		
		setInterval(refreshData, 1000);


		$('.search-input').on('input change', function() {
			var jsondata={
				"Criteria": {
							machine_name: $('#search-by-MachineName').val(),
							part_name_or_kanabn_no: $('#search-by-PartNameorKanbanNo').val(),  
							}
						};
			
						$.post('/search-kanban-view', JSON.stringify(jsondata), function (response) {
							$('#advanced-search-table-body').html(response);
						});

		});
	});

	//!js
	</script>
	`

	helpDialog := utils.JoinStr(`
	<div class="modal fade" id="help-dialog" tabindex="-1" aria-labelledby="helpModalLabel" aria-hidden="true">
		<div class="modal-dialog">
			<div class="modal-content">
			<div class="modal-header">
				<h4 class="modal-title" id="helpModalLabel" style="color:#871a83;">Help</h4>
				<button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
			</div>
			<div class="modal-body">
				<div class="row ms-3">
					<h5 class="p-0">Color Codes - Priority</h5>
					<div class="row mt-1">
						<a class="col-2 btn btn-primary btn-sm" style="background:`, utils.KanbanPriorityColors["regular"]["bg-color"], `"></a><p class="col mb-0">- Regular</p>
					</div>
					<div class="row mt-1">
						<a class="col-2 btn btn-primary btn-sm" style="background:`, utils.KanbanPriorityColors["urgent"]["bg-color"], `"></a><p class="col mb-0">- Urgent</p>
					</div>
					<div class="row mt-1">
						<a class="col-2 btn btn-primary btn-sm" style="background:`, utils.KanbanPriorityColors["mosturgent"]["bg-color"], `"></a><p class="col mb-0">- Most Urgent</p>
					</div>
				</div>
			</div>
			<div class="modal-footer">
				<button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Close</button>
			</div>
			</div>
		</div>
	</div>
	`)

	AllKanbanView.CardHeadingActions.Component = []s.CardHeadActionComponent{{ComponentName: "modelbutton", ComponentType: s.ActionComponentElement{ModelButton: s.ModelButtonAttributes{ID: "help-modal", Name: "help", Type: "button", Text: " <i class=\"fas fa-question-circle\"></i>  Help", ModelID: "#help-dialog"}}, Width: "col-2"}}

	html.WriteString(AllKanbanView.Build())
	html.WriteString(helpDialog)
	html.WriteString(`</div>`)
	html.WriteString(js)
	return html.String()

}
