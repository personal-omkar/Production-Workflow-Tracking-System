package logs

import (
	"os"
	"sort"
	"strings"

	"irpl.com/kanban-commons/model"
	"irpl.com/kanban-commons/utils"
	"irpl.com/kanban-web/services"
)

type SystemLogPage struct {
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

func (s *SystemLogPage) Build() string {
	var html strings.Builder
	content := utils.JoinStr(`

	<div class="container-xxl mt-3" >
		<div class="card">
			<div class="card-header d-flex justify-content-between align-items-center" style="background-color:#F4F5FB">
				<h3 class="heading-text mb-0" style="color:#871a83;">System Logs</h3>
				<div>
				</div>
			</div>
			<div class="card-body scrollable p-0" id="system-logs" style="overflow-x: auto; max-height:75vh">
				<table id="system-logs-table" class="table table-sm table-striped fs--1 mb-0 overflow-auto scrollbar " style="white-space:nowrap;">
					`, s.systemLogsTable(), `
				</table>
			</div>
			<div class="card-footer text-muted d-flex justify-content-end  align-items-center">
			</div>
		</div>
	</div>

	`)

	html.WriteString(content)
	return html.String()
}

func (s *SystemLogPage) systemLogsTable() string {

	var systemLogTable services.TableCard
	var SysLogs []*model.SystemLog
	var rawQuery model.RawQuery
	rawQuery.Host = utils.RestHost
	rawQuery.Port = utils.RestPort
	rawQuery.Type = "SystemLog"
	// use this query if you need to show logs accordingly to user login
	// rawQuery.Query = `SELECT * FROM SystemLogs where CAST(created_by AS INTEGER) = ` + s.UserID + `;`
	rawQuery.Query = `SELECT * FROM SystemLogs;`
	rawQuery.RawQry(&SysLogs)

	// Sorting by time (most recent first)
	sort.Slice(SysLogs, func(i, j int) bool {
		return SysLogs[i].CreatedOn.After(SysLogs[j].Timestamp)
	})
	systemLogTable.BodyTables = services.CardTableBody{Columns: []services.CardTableBodyHeadCol{
		{
			Name:         `Created On`,
			IsSearchable: false,
			ID:           "search-createdon",
			Width:        "col-1",
		},
		{
			Name:         "Message Type",
			IsSearchable: false,
			ID:           "search-Type",
			Width:        "col-1",
		},
		{
			Name:         "Message",
			IsSearchable: false,
			ID:           "search-message",
			Width:        "col-10",
		},
	},
		Data:         SysLogs,
		ColumnsWidth: []string{"col-1", "col-1", "col-10"},
	}
	return systemLogTable.BodyTables.GenrateCaradTables()
}
