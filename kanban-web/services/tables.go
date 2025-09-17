package services

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"

	"irpl.com/kanban-commons/model"
	u "irpl.com/kanban-commons/utils"
)

type TableCard struct {
	CardHeading        string
	CardHeadingActions CardHeadActionBody
	BodyAction         CardActionBody
	BodyTables         CardTableBody
	CardFooter         string
	Alert              bool
	For                string
	Class              string
}

type CardActionBody struct {
	Component []ActionComponent
}

type CardHeadActionBody struct {
	Component []CardHeadActionComponent
	Style     string
}

type ActionComponent struct {
	ComponentName string
	ComponentType ActionComponentElement
	Width         string
}

type CardHeadActionComponent struct {
	ComponentName string
	ComponentType ActionComponentElement
	Width         string
	Style         string
}
type ActionComponentElement struct {
	ModelButton    ModelButtonAttributes
	Input          InputAttributes
	Button         ButtonAttributes
	Browse         BrowseAttributes
	TextArea       TextAreaAttributes
	DropDown       DropdownAttributes
	Span           SpanAttributes
	PaginationReq  model.PaginationReq
	PaginationResp model.PaginationResp
}

type TextAreaAttributes struct {
	Name        string
	ID          string
	Label       string
	Placeholder string
	Value       string
	Readonly    bool
	Visibility  bool
	DataType    string
	Disabled    bool
}

type SpanAttributes struct {
	Text  string
	Class string
}
type InputAttributes struct {
	Name           string
	ID             string
	Type           string
	Label          string
	Icon           string
	Placeholder    string
	Value          string
	Width          string
	AdditionalAttr string
	Readonly       bool
	Disabled       bool
	Required       bool
	Hidden         bool
	DataType       string
}
type DropdownAttributes struct {
	Type           string
	Name           string
	ID             string
	Options        []DropDownOptions
	Label          string
	Width          string
	Disabled       bool
	AdditionalAttr string
	DataType       string
	Selected       bool
	Hidden         bool
}
type DropDownOptions struct {
	Value    string
	Text     string
	Selected bool
	Disabled bool
}
type ButtonAttributes struct {
	Name           string
	ID             string
	Text           string
	Type           string
	Disabled       bool
	Colour         string
	Class          string
	CSS            string
	AdditionalAttr string
	IsDownloadBtn  bool
	DownloadLink   string
}
type ModelButtonAttributes struct {
	Name     string
	ID       string
	Text     string
	Type     string
	ModelID  string
	Disabled bool
}

type BrowseAttributes struct {
	InputName   string
	InputID     string
	BrowseLabel string
	ButtonsName string
	ButtoneID   string
	ButtoneText string
	ButtoneType string
}
type CardTableBodyHeadCol struct {
	Lable            string
	Name             string
	ID               string
	IsSearchable     bool
	SearchFieldWidth string
	IsCheckbox       bool
	ActionList       []DropDownOptions
	Type             string
	Value            string
	Width            string
	GetSum           bool
	IsSortable       bool
	Style            string
	DataField        string
}

type CardTableBody struct {
	TableHeading      CardTableHeading
	ID                string
	Username          string
	Columns           []CardTableBodyHeadCol
	ColumnsWidth      []string
	Data              interface{}
	Tools             string
	Buttons           string
	Conditional_Tools map[bool]string
}

type CardTableHeading struct {
	Heading      string
	HeadingWidth string
	Component    []Component
}
type Subcomponents struct {
	Component string
	Width     string
	Item      string
	AlignItem string
}
type Component struct {
	Component     string
	Width         string
	Subcomponents []Subcomponents
}

var currentPage, rowsPerPage int

func (t *TableCard) Build() string {
	if strings.ToLower(t.For) == "accordion" {
		CardHeading := ``

		if t.CardHeading != "" {
			CardHeading = u.JoinStr(`<div class="row container">
						<span class="col-6 p-2"><h5 style="color:#AB71A2; font-size:0.9rem;">`, t.CardHeading, `</h5>  </span>
					</div>`)
		}
		tableCard := u.JoinStr(`
		<div class="container-xxl m-0 p-0 mt-0" id="Table-div">
			<div class="card overflow-auto container-fluid p-0 m-0"  style="box-shadow: none !important;">
				<div class="card-body container-fluid p-0 m-0">
					`, CardHeading, `
					<div class="card-body p-0 mt-0"  style="box-shadow: none !important;">
						<div class="d-flex" style="overflow-y: auto;overflow-x: auto; max-height:50vh box-shadow: none !important;" >
							<table id="advanced-search-table" class="table table-sm table-striped fs--1 mb-0 `, t.Class, `" style="white-space:nowrap;">
								`, t.BodyTables.genrateCaradTables(), `
							</table>
						</div>
					</div>		
				</div>
			</div>
		</div>
		`)
		return tableCard
	}
	tableCard := u.JoinStr(`
	<div class="container-xxl m-0 p-0" id="Table-div">
		<div class="card overflow-auto">
			<div class="card-body py-0">
				<div class="row">
					<span class="col-6"><h1 class="mt-3" style="color:#871a83;">`, t.CardHeading, `</h1>  </span>
					`, t.CardHeadingActions.genrateCaradHeadingActions(), `
				</div>
				<div id="notification" class="alert alert-success alert-dismissible fade show mb-0 mt-1 p-1 w-100" role="alert" style="display: none;"></div>
				`, t.BodyAction.genrateCaradActions(), ` 
				<div class="card-body p-1 mt-1 ">
					<div style="overflow-y: auto;overflow-x: auto; max-height:70vh;" >
					<table id="advanced-search-table" class="table table-sm table-striped fs--1 mb-0 `, t.Class, `" style="white-space:nowrap;">
						`, t.BodyTables.genrateCaradTables(), `
					</table>
					</div>
			  	</div>
			</div>		
	`)

	tableCard = u.JoinStr(tableCard, `
			<div class="card-footer text-muted d-flex justify-content-end  align-items-center">
				`, t.CardFooter, `
			</div>
		</div>
	</div>
	`)
	return tableCard
}

func (t CardActionBody) genrateCaradActions() string {
	var components string
	var element string

	for _, i := range t.Component {
		if i.ComponentName == "button" {
			element = u.JoinStr(`
			<button type="`, i.ComponentType.Button.Type, `"  id="`, i.ComponentType.Button.ID, `" name="`, i.ComponentType.Button.Name, `" class="btn btn-primary p-2" style="background-color:#871a83 ;border:none">`, i.ComponentType.Button.Text, `</button>`)
		} else if i.ComponentName == "modelbutton" {
			element = u.JoinStr(`
			<button type="`, i.ComponentType.ModelButton.Type, `" id="`, i.ComponentType.ModelButton.ID, `" name="`, i.ComponentType.ModelButton.Name, `" class="btn btn-primary  p-2" data-bs-toggle="modal" data-bs-target="`, i.ComponentType.ModelButton.ModelID, `" style="background-color:#871a83 ;border:none;">`, i.ComponentType.ModelButton.Text, `</button>`)
		} else if i.ComponentName == "input" {
			element = u.JoinStr(`
			<div class="input-group mb-3">
				<span class="input-group-text" id="basic-addon1" > `, i.ComponentType.Input.Icon, `</span>
				<input type="`, i.ComponentType.Input.Type, `"  class="form-control" name=`, i.ComponentType.Input.Name, ` placeholder="`, i.ComponentType.Input.Placeholder, `" aria-label="Username" aria-describedby="basic-addon1">
			</div>`)

		} else if i.ComponentName == "browse" {
			element = u.JoinStr(`
			<div class="input-group ">
				<span class="input-group-text" id="basic-addon1" > `, i.ComponentType.Browse.BrowseLabel, `</span>
				<input type="`, i.ComponentType.Browse.InputID, `"  class="form-control" name=`, i.ComponentType.Browse.InputName, ` " aria-label="Username" aria-describedby="basic-addon1">
				<button type="`, i.ComponentType.Browse.ButtoneType, `"  id="`, i.ComponentType.Browse.ButtoneID, `" name="`, i.ComponentType.Browse.ButtonsName, `" class="btn btn-primary btn-lg" style="background-color:#871a83  ;border:none">`, i.ComponentType.Browse.ButtoneText, `</button>
				</div>`)

		} else if i.ComponentName == "dropdown" {
			var options string
			for _, opt := range i.ComponentType.DropDown.Options {
				tempopt := u.JoinStr(`<option selected value="`, opt.Value, `">`, opt.Text, `</option>`)
				options = options + tempopt
			}
			if i.ComponentType.DropDown.Label == "" {
				element = u.JoinStr(`
				<div class="col-12 d-flex p-0  align-items-center">
					<select id="search-by- " class="form-select p-2 " id=`, i.ComponentType.DropDown.ID, ` name=`, i.ComponentType.DropDown.Name, ` aria-label=".form-select-sm example">
						`, options, `
					</select>
				</div>`)
			} else if i.ComponentType.DropDown.Label != "" {
				element = u.JoinStr(`
				<div class="col-12 d-flex p-0 w-50 align-items-center">
					<label for="cheese" class="col-6 d-flex align-items-end flex-column p-2">`, i.ComponentType.DropDown.Label, `</label>
					<select id="search-by- " class="form-select p-2 " id=`, i.ComponentType.DropDown.ID, ` name=`, i.ComponentType.DropDown.Name, ` aria-label=".form-select-sm example">
						`, options, `
					</select>
				</div>`)

			}
		} else if i.ComponentName == "date" {
			var options string
			if i.ComponentType.Input.Placeholder == "H:i" {
				options = `{"enableTime":true,"noCalendar":true,"dateFormat":"H:i","disableMobile":true}`
			} else if i.ComponentType.Input.Placeholder == "dd/mm/yyyy" {
				options = `{"disableMobile":true,"dateFormat":"d-m-Y"}`
			} else {
				options = `{"disableMobile":true}`
			}
			element = u.JoinStr(`
					<label for="cheese" class="col-6 d-flex p-0 flex-column">`, i.ComponentType.Input.Label, `</label>
					<div class="input-group mb-2">
						<span class="input-group-text p-0" onclick="$('#`, i.ComponentType.Input.ID, `').click()"></span>
						<input class="form-control datetimepicker" id="`, i.ComponentType.Input.ID, `" data-name="`, i.ComponentType.Input.ID, `" type="text" data-type="`, i.ComponentType.Input.Type, `" placeholder="`, i.ComponentType.Input.Placeholder, `" data-options='`, options, `' / `, `> 
					</div> `)
		} else if i.ComponentName == "Pagination" {
			element = u.JoinStr(`
			<!--html-->
			<style>
				.active .page-link{
					background-color: #b250ad;
					border : none;
					outline : none;
				}
				.page-link:hover{
					background-color: #b250ad;
					border : none;
					outline : none;
					color: white;
				}
				.pipe-symbol {
					line-height: 100%; /* Ensure the symbol takes the full height */
				}
			</style>
			<div class="d-flex align-item-center justify-content-end" style="gap: 20px;">
			<!-- Pagination Navigation -->
			 `, BuildPagination(i.ComponentType.PaginationResp.TotalNo, i.ComponentType.PaginationResp.Page, i.ComponentType.PaginationResp.Offset), `
			</div>
			<!--!html-->
			`)
			currentPage = i.ComponentType.PaginationResp.Page
			rowsPerPage = i.ComponentType.PaginationResp.Offset
		}

		tempcomponent := u.JoinStr(`
		<div class="`, i.Width, ` d-flex flex-column ">
			`, element, `
		</div>`)
		components = components + tempcomponent
	}
	cardActions := u.JoinStr(`
	        <div class="row" >
            `, components, ` 
        	</div>
	`)
	return cardActions
}

func (t CardHeadActionBody) genrateCaradHeadingActions() string {
	var components string
	var element string
	var color string
	var disable string
	for _, i := range t.Component {
		if i.ComponentName == "button" {

			if i.ComponentType.Button.Colour != "" {
				color = i.ComponentType.Button.Colour
			} else {
				color = "#871a83"
			}
			if i.ComponentType.Button.Disabled {
				disable = "disabled"
			} else {
				disable = ""
			}

			if i.ComponentType.Button.CSS == "" {
				i.ComponentType.Button.CSS = "border:none;"
			}
			if i.ComponentType.Button.IsDownloadBtn {
				element = u.JoinStr(`
				 <a href="`, i.ComponentType.Button.DownloadLink, `" download> <button type="`, i.ComponentType.Button.Type, `"  id="`, i.ComponentType.Button.ID, `" name="`, i.ComponentType.Button.Name, `" class="btn btn-primary p-2 `, i.ComponentType.Button.Class, `" style="`, i.ComponentType.Button.CSS, ` background-color:`, color, ` ;" `, disable, ` `, i.ComponentType.Button.AdditionalAttr, `>`, i.ComponentType.Button.Text, `</button> </a>`)
			} else {
				element = u.JoinStr(`
				<button type="`, i.ComponentType.Button.Type, `"  id="`, i.ComponentType.Button.ID, `" name="`, i.ComponentType.Button.Name, `" class="btn btn-primary p-2 `, i.ComponentType.Button.Class, `" style="`, i.ComponentType.Button.CSS, ` background-color:`, color, ` ;" `, disable, ` `, i.ComponentType.Button.AdditionalAttr, `>`, i.ComponentType.Button.Text, `</button>`)
			}

		} else if i.ComponentName == "modelbutton" {
			if i.ComponentType.ModelButton.Disabled {
				disable = "disabled"
			} else {
				disable = ""
			}
			element = u.JoinStr(`
			<button type="`, i.ComponentType.ModelButton.Type, `" id="`, i.ComponentType.ModelButton.ID, `" name="`, i.ComponentType.ModelButton.Name, `" class="btn btn-primary  p-2" data-bs-toggle="modal" data-bs-target="`, i.ComponentType.ModelButton.ModelID, `" style="background-color:#871a83 ;border:none;"  `, disable, `>`, i.ComponentType.ModelButton.Text, `</button>`)
		} else if i.ComponentName == "input" {
			element = u.JoinStr(`
			<div class="input-group mb-3" style="min-height:75%">
				<span class="input-group-text" id="basic-addon1" > `, i.ComponentType.Input.Icon, `</span>
				<input id=`, i.ComponentType.Input.ID, ` type="`, i.ComponentType.Input.Type, `"  class="form-control" name=`, i.ComponentType.Input.Name, ` placeholder="`, i.ComponentType.Input.Placeholder, `" aria-label="Username" aria-describedby="basic-addon1">
			</div>`)

		} else if i.ComponentName == "browse" {
			element = u.JoinStr(`
			<div class="input-group ">
				<span class="input-group-text" id="basic-addon1" > `, i.ComponentType.Browse.BrowseLabel, `</span>
				<input type="`, i.ComponentType.Browse.InputID, `"  class="form-control" name=`, i.ComponentType.Browse.InputName, ` " aria-label="Username" aria-describedby="basic-addon1">
				<button type="`, i.ComponentType.Browse.ButtoneType, `"  id="`, i.ComponentType.Browse.ButtoneID, `" name="`, i.ComponentType.Browse.ButtonsName, `" class="btn btn-primary btn-lg" style="background-color:#871a83  ;border:none">`, i.ComponentType.Browse.ButtoneText, `</button>
				</div>`)

		} else if i.ComponentName == "span" {
			element = u.JoinStr(`
									<span class="`, i.ComponentType.Span.Class, `" id="basic-addon1" > `, i.ComponentType.Span.Text, `</span>
								`)
		} else if i.ComponentName == "dropdown" {

			if i.ComponentType.DropDown.Disabled {
				disable = "disabled"
			} else {
				disable = ""
			}

			var options string
			if !i.ComponentType.DropDown.Selected {
				for i, opt := range i.ComponentType.DropDown.Options {
					selected := ""
					if i == 0 {
						selected = "selected"
					}
					tempopt := u.JoinStr(`<option `, selected, ` value="`, opt.Value, `">`, opt.Text, `</option>`)
					options = options + tempopt
				}
			} else {
				for _, opt := range i.ComponentType.DropDown.Options {
					selected := ""
					if opt.Selected {
						selected = "selected"
					}
					tempopt := u.JoinStr(`<option `, selected, ` value="`, opt.Value, `">`, opt.Text, `</option>`)
					options = options + tempopt
				}
			}
			if i.ComponentType.DropDown.Label == "" {
				element = u.JoinStr(`
				<div class="col-12 d-flex p-0  align-items-center">
					<select id="search-by- " class="form-select p-2 " id=`, i.ComponentType.DropDown.ID, ` name=`, i.ComponentType.DropDown.Name, ` aria-label=".form-select-sm example">
						`, options, `
					</select>
				</div>`)
			} else if i.ComponentType.DropDown.Label != "" {
				element = u.JoinStr(`
				<div class="input-group mb-3 col-12 d-flex justify-content-end"">
					<div class="input-group-prepend">
						<label class="input-group-text h-100" for="`, i.ComponentType.DropDown.ID, `">`, i.ComponentType.DropDown.Label, `</label>
					</div>
					<select class="custom-select col-7" id="`, i.ComponentType.DropDown.ID, `" `, disable, `>
						`, options, `
					</select>
				</div>`)

			}
		} else if i.ComponentName == "date" {
			var options string
			if i.ComponentType.Input.Placeholder == "H:i" {
				options = `{"enableTime":true,"noCalendar":true,"dateFormat":"H:i","disableMobile":true}`
			} else if i.ComponentType.Input.Placeholder == "dd/mm/yyyy" {
				options = `{"disableMobile":true,"dateFormat":"d-m-Y"}`
			} else {
				options = `{"disableMobile":true}`
			}
			element = u.JoinStr(`
					<label for="cheese" class="col-6 d-flex p-0 flex-column">`, i.ComponentType.Input.Label, `</label>
					<div class="input-group mb-2">
						<span class="input-group-text p-0" onclick="$('#`, i.ComponentType.Input.ID, `').click()"></span>
						<input class="form-control datetimepicker" id="`, i.ComponentType.Input.ID, `" data-name="`, i.ComponentType.Input.ID, `" type="text" data-type="`, i.ComponentType.Input.Type, `" placeholder="`, i.ComponentType.Input.Placeholder, `" data-options='`, options, `' / `, `> 
					</div> `)
		}

		tempcomponent := u.JoinStr(`
		<div class="`, i.Width, ` d-flex flex-column ">
			`, element, `
		</div>`)
		components = components + tempcomponent
	}
	if t.Style == "" {
		cardActions := u.JoinStr(`
				<span class="row py-2 col-6  mt-3" style="direction: rtl;">
				`, components, ` 
				</span>
				
		`)
		return cardActions
	} else {
		cardActions := u.JoinStr(`
				<span class="row py-2 col-6  mt-3" style="`, t.Style, `">
				`, components, ` 
				</span>
		`)
		return cardActions
	}

}

func (t CardTableBody) genrateCaradTables() string {
	var tHead, tBoday string

	tHead = `<thead class="text-900" style="white-space: nowrap; position: sticky; top: 0; "><tr>`

	for _, col := range t.Columns {
		tHead = u.JoinStr(tHead, col.renderHeadColumns(t.Data))
	}
	tHead = u.JoinStr(tHead, `</tr></thead>`)

	tBoday = u.JoinStr(`<tbody class="list" id="advanced-search-table-body">`, t.RenderBodyColumns(), `</tbody>`)
	table := u.JoinStr(tHead, tBoday)
	// table := tHead

	return table
}

func (tableColumn *CardTableBodyHeadCol) renderHeadColumns(data interface{}) string {
	field := tableColumn.renderSearchField()
	if tableColumn.Name == "Columnse" {
		return u.JoinStr(`<th class="`, tableColumn.Width, ` p-2 fs-10 pe-1 align-top " style="background-color:#f1f1f1;color:#F4F5FB ;white-space: nowrap;  margin-right: 10px; "> `, tableColumn.Name, " ", field, `</th>`)
	} else if tableColumn.Name == "Button" {
		return u.JoinStr(`<th class="`, tableColumn.Width, ` p-2 fs-10 pe-1 align-top " style="background-color:#f1f1f1;color:#f1f1f1 ;white-space: nowrap;  margin-right: 10px; "> `, tableColumn.Name, " ", field, `</th>`)
	} else if tableColumn.Name == "Vendor-Action" {
		return u.JoinStr(`<th class="`, tableColumn.Width, ` p-2 fs-10 pe-1 align-top " style="background-color:#f1f1f1;color:#f1f1f1 ;white-space: nowrap;  margin-right: 10px; "> `, tableColumn.Name, " ", field, `</th>`)
	} else {
		sortData := ""
		if tableColumn.IsSortable {
			sortData = `data-Sort=Data-Sort-` + strings.ReplaceAll(tableColumn.Name, " ", "")
		}
		if tableColumn.GetSum {
			DataSum := getSumOfColData(tableColumn.Name, data)
			if tableColumn.Lable == "" {
				return u.JoinStr(`<th class="`, tableColumn.Width, ` p-2 fs-10 pe-1 align-top " `, sortData, `style="background-color:#f1f1f1;color:#5B5B5B; white-space: nowrap;  margin-right: 10px; "> `, tableColumn.Name, ` (`, DataSum, `)  `, field, `</th>`)
			}
			return u.JoinStr(`<th class="`, tableColumn.Width, ` p-2 fs-10 pe-1 align-top " `, sortData, ` style="background-color:#f1f1f1;color:#5B5B5B; white-space: nowrap;  margin-right: 10px; "> `, tableColumn.Lable, ` (`, DataSum, `)  `, field, `</th>`)
		} else {
			if tableColumn.Lable == "" {
				return u.JoinStr(`<th class="`, tableColumn.Width, ` p-2 fs-10 pe-1 align-top " `, sortData, ` style="background-color:#f1f1f1;color:#5B5B5B; white-space: nowrap;  margin-right: 10px; "> `, tableColumn.Name, " ", field, `</th>`)
			}
			return u.JoinStr(`<th class="`, tableColumn.Width, ` p-2 fs-10 pe-1 align-top " `, sortData, ` style="background-color:#f1f1f1;color:#5B5B5B; white-space: nowrap;  margin-right: 10px; "> `, tableColumn.Lable, " ", field, `</th>`)
		}
	}
}

func (tableColumn *CardTableBodyHeadCol) renderSearchField() string {
	var field string
	var searchFieldWidth string

	//Search Icon
	if tableColumn.SearchFieldWidth != "" {
		searchFieldWidth = tableColumn.SearchFieldWidth
	} else {
		searchFieldWidth = "w-100"
	}

	if tableColumn.IsSearchable {
		switch tableColumn.Type {
		case "input":
			field += u.JoinStr(`
					<div class="d-flex align-items-center">
						<div class="content-search-box mt-2 border border-2 `, searchFieldWidth, `" style="border-radius: 0rem">
							<form class="position-relative" data-toggle="search" data-display="static">
								<div class="input-group border-0">
									<input id="search-by-`, strings.ReplaceAll(tableColumn.Name, " ", ""), `" class="form-control search-input ps-1 border-0" type="search" placeholder="" data-field="`, tableColumn.DataField, `" aria-label="Search" value="`, tableColumn.Value, `" />		
								</div>
							</form>
						</div>
					</div>`)
		case "select":
			var options string
			for _, i := range tableColumn.ActionList {
				op := u.JoinStr(`<option id="select-action-opt" value="`, i.Value, `">`, i.Text, `</option>`)
				options += op
			}
			field += u.JoinStr(`
					<div class="col d-md-block d-none p-0 mt-2 `, searchFieldWidth, `">
						<select data-field="`, tableColumn.DataField, `" id="search-by-`, strings.ReplaceAll(tableColumn.Name, " ", ""), `" class="form-select form-select col-7" aria-label=".form-select-sm example">
							`, options, `
						</select>
					</div>`)
		case "action":
			var options string
			for _, i := range tableColumn.ActionList {
				op := u.JoinStr(`<option id="select-action-opt" value="`, i.Value, `">`, i.Text, `</option>`)
				options += op
			}
			field += u.JoinStr(`
					<div class="d-flex align-items-center" id="action-line" >
					<label for="selectMenu" class="me-2" style="color:#5B5B5B;"> <b>Action:</b> </label>
					<select id="selectMenu-action-line" class="form-select" disabled>
						`, options, `
					</select>
				</div>
				`)
		case "date":
			field += u.JoinStr(`
				<div class="d-flex align-items-center">
					<div class="content-search-box mt-2 border border-2 `, searchFieldWidth, `" style="border-radius: 0rem">
						<div class="input-group border-0">
							<input id="search-by-`, strings.ReplaceAll(tableColumn.Name, " ", "_"), `" 
								class="form-control search-input ps-1 border-0" 
								type="date" 
								data-field="`, tableColumn.DataField, `" 
								aria-label="Search" 
								value="`, tableColumn.Value, `" 
								placeholder="Select a date" />
						</div>
					</div>
				</div>`)

		}
	} else {
		field += ""
	}
	return field
}

func (table *CardTableBody) RenderBodyColumns() string {
	var res string
	switch v := table.Data.(type) {
	case []*model.VendorCompanyTable:
		res = table.VendorCompanyTable()
	case []*model.OrderDetails:
		res = table.OrderDetailsTable()
	case []*model.CustomerOrderDetails:
		res = table.CustomerOrderDetailsTable()
	case []*model.ColdStorage:
		res = table.ColdStorageTable()
	case []*model.UserManagement:
		res = table.UserManagementTable()
	case []*model.UserRoles:
		res = table.UserRoleManagementTable()
	case []*model.Compounds:
		res = table.CompoundsManagementTable()
	case []*model.Vendors:
		res = table.VendorManagementTable()
	case []*model.SystemLog:
		res = table.SystemLogsTable()
	case []*model.ProdLine:
		res = table.ProdLineManagementTable()
	case []*model.ProdProcess:
		res = table.ProdProcessesManagementTable()
	case []*model.Recipe:
		res = table.RecipeManagementTable()
	case []*model.Stage:
		res = table.StageManagementTable()
	case []*model.AllKanbanViewTable:
		res = table.AllKanbanViewTable()
	case []*model.Operator:
		res = table.OperatorManagementTable()
	case []*model.RawMaterial:
		res = table.MaterialManagementTable()
	case []*model.ChemicalTypes:
		res = table.ChemicalManagementTable()
	default:
		fmt.Printf("The type is string with value %s\n", v)
	}
	return res
}

func (table *CardTableBody) loadColumsData(Table map[string]interface{}, rowIndex int) string {
	var sb strings.Builder
	textColor, _ := Table["TextColor"].(string)

	for i, col := range table.Columns {
		colKey := strings.ReplaceAll(col.Name, " ", "")
		var colWidth string
		if i < len(table.ColumnsWidth) {
			colWidth = table.ColumnsWidth[i]
		} else {
			colWidth = ""
		}
		baseClass := `class="` + colWidth + ` p-2 pe-1 align-middle"`

		switch col.Name {
		case "Button":
			sb.WriteString(u.JoinStr(`<td `, baseClass, `>`, table.Buttons, `</td>`))

		case "Tools":
			sb.WriteString(`<td class="p-2 pe-1 align-middle">` + table.Tools + `</td>`)

		case "Sr. No.":
			serialNo := (currentPage-1)*rowsPerPage + rowIndex + 1
			sb.WriteString(u.JoinStr(`<td `, baseClass, `>`, strconv.Itoa(serialNo), `</td>`))

		case "Message Type":
			sb.WriteString(u.JoinStr(`<td class="`, colWidth, ` `, textColor, ` p-3 pe-1 align-middle">`, fmt.Sprint(Table[colKey]), `</td>`))

		case "Action":
			sb.WriteString(u.JoinStr(`<td `, baseClass, `>`, table.Buttons, table.ID, `</td>`))

		case "Conditional_Tools":
			tool := table.Conditional_Tools[isToday(Table["DemandDateTime"])]
			sb.WriteString(`<td class="p-2 pe-1 align-middle">` + tool + `</td>`)

		default:
			value := Table[colKey]
			switch val := value.(type) {
			case float64:
				sb.WriteString(u.JoinStr(`<td class="align-middle p-2 pe-1 `, colWidth, `">`, fmt.Sprintf("%.f", val), `</td>`))

			default:
				// Special formatting for date fields
				if col.Name == "Created On" || col.Name == "Modified On" || col.Name == "Demand Date Time" || col.Name == "MFGDateTime" {
					format := "date-time"
					if col.Name == "MFGDateTime" {
						format = "date"
					}
					if col.Name == "Demand Date Time" {
						format = "short-date-time"
					}
					formatted := u.FormatStringDate(fmt.Sprint(val), format)
					sb.WriteString(u.JoinStr(`<td `, baseClass, `>`, formatted, `</td>`))

				} else if col.Name == "Columnse" {
					sb.WriteString(u.JoinStr(`<td class="`, colWidth, ` p-3 pe-1 align-middle" style="background-color:#F4F5FB;color:#5B5B5B">`, fmt.Sprint(val), `</td>`))

				} else if col.Name == "Status" {
					statusStr := fmt.Sprint(val)
					if statusVal, ok := u.StatusMap[statusStr]; ok {
						sb.WriteString(u.JoinStr(`<td `, baseClass, `>`, statusVal, `</td>`))
					} else {
						sb.WriteString(u.JoinStr(`<td `, baseClass, `>`, statusStr, `</td>`))
					}

				} else {
					if valStr := fmt.Sprint(val); valStr != "" {
						sb.WriteString(u.JoinStr(`<td class="`, colWidth, ` p-2 pe-1 align-middle" style="`, col.Style, `">`, valStr, `</td>`))
					} else {
						sb.WriteString(`<td class="p-2 pe-1 align-middle"></td>`)
					}
				}
			}
		}
	}

	return sb.String()
}

func (table *CardTableBody) VendorCompanyTable() string {
	var sb strings.Builder

	var TableMap map[string]interface{}

	for i, column := range table.Data.([]*model.VendorCompanyTable) {

		var bgColor = ""

		datadata, _ := json.Marshal(column)

		rowData, _ := json.Marshal(struct {
			VendorCode string
			VendorName string
		}{
			VendorCode: column.VendorCode,
			VendorName: column.VendorName,
		})

		json.Unmarshal(datadata, &TableMap)

		s := strings.ReplaceAll(string(rowData), "'", "&#39;")
		json.Unmarshal(datadata, &TableMap)

		sb.WriteString(u.JoinStr(`<tr id='`, table.ID, `-r-`, strconv.Itoa(i), `' data-data='`, s, `' class="tb-row p-3 btn-reveal-trigger fw-semi-bold" style="background:`, bgColor, `">`, table.loadColumsData(TableMap, i), `</tr>`))

	}
	return sb.String()
}

func (table *CardTableBody) KbDataTable() string {
	var sb strings.Builder

	var TableMap map[string]interface{}

	for i, column := range table.Data.([]*model.KbData) {

		var bgColor = ""

		datadata, _ := json.Marshal(column)

		json.Unmarshal(datadata, &TableMap)

		s := strings.ReplaceAll(string(datadata), "'", "&#39;")

		sb.WriteString(u.JoinStr(`<tr id='`, table.ID, `-r-`, strconv.Itoa(i), `' data-data='`, s, `' class="tb-row p-3 btn-reveal-trigger fw-semi-bold" style="background:`, bgColor, `">`, table.loadColumsData(TableMap, i), `</tr>`))

	}
	return sb.String()
}

func (table *CardTableBody) OrderDetailsTable() string {
	var sb strings.Builder

	var TableMap map[string]interface{}

	for i, column := range table.Data.([]*model.OrderDetails) {
		var bgColor = ""

		datadata, _ := json.Marshal(column)

		json.Unmarshal(datadata, &TableMap)

		s := strings.ReplaceAll(string(datadata), "'", "&#39;")

		sb.WriteString(u.JoinStr(`<tr id='`, table.ID, `-r-`, strconv.Itoa(i), `' data-data='`, s, `' class="tb-row p-3 btn-reveal-trigger fw-semi-bold" style="background:`, bgColor, `">`, table.loadColumsData(TableMap, i), `</tr>`))

	}
	return sb.String()
}

func (table *CardTableBody) CustomerOrderDetailsTable() string {
	var sb strings.Builder

	var TableMap map[string]interface{}

	for i, column := range table.Data.([]*model.CustomerOrderDetails) {
		var bgColor = ""

		datadata, _ := json.Marshal(column)

		json.Unmarshal(datadata, &TableMap)

		s := strings.ReplaceAll(string(datadata), "'", "&#39;")

		sb.WriteString(u.JoinStr(`<tr id='`, table.ID, `-r-`, strconv.Itoa(i), `' data-data='`, s, `' class="tb-row p-3 btn-reveal-trigger fw-semi-bold" style="background:`, bgColor, `">`, table.loadColumsData(TableMap, i), `</tr>`))

	}
	return sb.String()
}

func (table *CardTableBody) ColdStorageTable() string {
	var sb strings.Builder

	var TableMap map[string]interface{}

	for i, column := range table.Data.([]*model.ColdStorage) {

		var bgColor = ""

		datadata, _ := json.Marshal(column)

		json.Unmarshal(datadata, &TableMap)

		s := strings.ReplaceAll(string(datadata), "'", "&#39;")

		sb.WriteString(u.JoinStr(`<tr id='`, table.ID, `-r-`, strconv.Itoa(i), `' data-data='`, s, `' class="tb-row p-3 btn-reveal-trigger fw-semi-bold" style="background:`, bgColor, `">`, table.loadColumsData(TableMap, i), `</tr>`))

	}
	return sb.String()
}

func (table *CardTableBody) UserManagementTable() string {
	var sb strings.Builder

	var TableMap map[string]interface{}

	for i, column := range table.Data.([]*model.UserManagement) {
		var bgColor = ""

		datadata, _ := json.Marshal(column)
		rowData, _ := json.Marshal(struct {
			Email    string
			RoleName string
		}{
			Email:    column.Email,
			RoleName: column.RoleName,
		})
		json.Unmarshal(datadata, &TableMap)

		s := strings.ReplaceAll(string(rowData), "'", "&#39;")

		sb.WriteString(u.JoinStr(`<tr id='`, table.ID, `-r-`, strconv.Itoa(i), `' data-data='`, s, `' class="tb-row p-3 btn-reveal-trigger fw-semi-bold" style="background:`, bgColor, `">`, table.loadColumsData(TableMap, i), `</tr>`))

	}
	return sb.String()
}

func BuildPagination(totalRecordsNo, currentPage, offset int) string {
	var html strings.Builder
	totalPages := int(math.Ceil(float64(totalRecordsNo) / float64(offset)))
	// Start the pagination HTML
	html.WriteString(`
		<!--html-->
		<div class="d-flex align-items-center">
				<label for="n" class="m-0" style="font-size:1rem; color:#b250ad;"><b>` + strconv.Itoa(totalRecordsNo) + ` results</b></label>
		</div>
		<nav aria-label="Page navigation example" class="d-flex align-items-center">
			<ul class="pagination align-items-center p-0 m-0">
		<!--!html-->
	`)

	// Handle Previous button (disable if on the first page)
	if currentPage > 1 {
		html.WriteString(`<li class="page-item" onclick=pagination(` + strconv.Itoa(currentPage-1) + `)><a class="page-link" href="#" aria-label="Previous"><span aria-hidden="true" class="">&lsaquo;</span><span class="sr-only">Previous</span></a></li>`)
	} else {
		html.WriteString(`<li class="page-item disabled"><a class="page-link" href="#" aria-label="Previous"><span aria-hidden="true" class="">&lsaquo;</span><span class="sr-only">Previous</span></a></li>`)
	}

	// Determine the range of page numbers to display
	startPage := currentPage - 1
	endPage := currentPage + 1

	// Adjust for edge cases at the beginning or end of the page range
	if currentPage == 1 {
		startPage = 1
		endPage = 3
	} else if currentPage == totalPages {
		startPage = totalPages - 2
		endPage = totalPages
	}

	// Ensure page numbers stay within valid range
	if startPage < 1 {
		startPage = 1
	}
	if endPage > totalPages {
		endPage = totalPages
	}

	// Generate the page numbers dynamically based on the range
	for i := startPage; i <= endPage; i++ {
		if i == currentPage {
			html.WriteString(`<li class="page-item active" onclick=pagination(` + strconv.Itoa(i) + `)><a class="page-link">` + strconv.Itoa(i) + `</a></li>`)
		} else {
			html.WriteString(`<li class="page-item" onclick=pagination(` + strconv.Itoa(i) + `)><a class="page-link" >` + strconv.Itoa(i) + `</a></li>`)
		}
	}

	// Handle Next button (disable if on the last page)
	if currentPage < totalPages {
		html.WriteString(`<li class="page-item" onclick=pagination(` + strconv.Itoa(currentPage+1) + `)><a class="page-link" href="#" aria-label="Next"><span aria-hidden="true" class="">&rsaquo;</span><span class="sr-only">Next</span></a></li>`)
	} else {
		html.WriteString(`<li class="page-item disabled"><a class="page-link" href="#" aria-label="Next"><span aria-hidden="true" class="">&rsaquo;</span><span class="sr-only">Next</span></a></li>`)
	}

	// Close the pagination HTML
	html.WriteString(`
		<!--html-->
		</ul>
		</nav>
		<div class="align-items-center" style="position: relative; height: 90%; border: 1px solid gray;"></div>
		<!-- Number Input Field -->
		<div class="d-flex align-items-center">
			<label for="n" class="m-0" style="font-size:1rem;">Show&nbsp;</label>
			<input type="number" min="1" max="15" step="1" id="n" value="15" oninput="(validity.valid)||(value='');" class="form-control" style="width: 100px; background-color:#F1F1F1;">
		</div>
		<!--!html-->
	`)

	// Return the constructed HTML string
	return html.String()
}

func BuildPaginationJS(url string, condition []string) string {
	var html strings.Builder
	html.WriteString(`
		<script>
		//js
		function pagination(i) {
			var url = "` + url + `";
			Offset = $("#search-by-SerialNumber").val(),
			
			// Data to be sent to the backend (marshaled as JSON)
			var requestData = {
				PageNumber: i,  // The page number to request
				Limit: Offset,       // The number of records per page
			};
			// Make the AJAX request
			$.ajax({
				url: url,
				type: 'POST',                    // Sending a POST request
				contentType: 'application/json', // We are sending JSON data
				data: JSON.stringify(requestData), // Marshal requestData to JSON
				success: function (response) {
					// On success, clear the existing table body
					$('#Table-div').empty();
					$('#Table-div').html(response.tableBody);
				},
				error: function (xhr, status, error) {
					console.error('Error fetching data:', error);
				}
			});
		}
		//!js
		</script>
	`)
	return html.String()
}

func (table *CardTableBody) UserRoleManagementTable() string {
	var sb strings.Builder

	var TableMap map[string]interface{}

	for i, column := range table.Data.([]*model.UserRoles) {
		var bgColor = ""

		datadata, _ := json.Marshal(column)

		json.Unmarshal(datadata, &TableMap)

		s := strings.ReplaceAll(string(datadata), "'", "&#39;")

		sb.WriteString(u.JoinStr(`<tr id='`, table.ID, `-r-`, strconv.Itoa(i), `' data-data='`, s, `' class="tb-row p-3 btn-reveal-trigger fw-semi-bold" style="background:`, bgColor, `">`, table.loadColumsData(TableMap, i), `</tr>`))

	}
	return sb.String()
}

func (table *CardTableBody) VendorManagementTable() string {
	var sb strings.Builder

	var TableMap map[string]interface{}

	for i, column := range table.Data.([]*model.Vendors) {
		var bgColor = ""

		datadata, _ := json.Marshal(column)

		json.Unmarshal(datadata, &TableMap)

		s := strings.ReplaceAll(string(datadata), "'", "&#39;")

		sb.WriteString(u.JoinStr(`<tr id='`, table.ID, `-r-`, strconv.Itoa(i), `' data-data='`, s, `' class="tb-row p-3 btn-reveal-trigger fw-semi-bold" style="background:`, bgColor, `">`, table.loadColumsData(TableMap, i), `</tr>`))

	}
	return sb.String()
}

func (table *CardTableBody) CompoundsManagementTable() string {
	var compounDataTable strings.Builder

	var TableMap map[string]interface{}

	for i, column := range table.Data.([]*model.Compounds) {
		var bgColor = ""

		datadata, _ := json.Marshal(column)

		json.Unmarshal(datadata, &TableMap)

		s := strings.ReplaceAll(string(datadata), "'", "&#39;")

		compounDataTable.WriteString(u.JoinStr(`<tr id='`, table.ID, `-r-`, strconv.Itoa(i), `' data-data='`, s, `' class="tb-row p-3 btn-reveal-trigger fw-semi-bold" style="background:`, bgColor, `">`, table.loadColumsData(TableMap, i), `</tr>`))

	}
	return compounDataTable.String()
}

func (t CardTableBody) GenrateCaradTables() string {
	var tHead, tBoday string

	tHead = `<thead class="text-900" style="background-color:#F4F5FB; white-space: nowrap; position: sticky; top: 0;  "><tr>`

	for _, col := range t.Columns {
		tHead = u.JoinStr(tHead, col.renderHeadColumns(t.Data))
	}
	tHead = u.JoinStr(tHead, `</tr></thead>`)

	tBoday = u.JoinStr(`<tbody class="list" id="advanced-search-table-body">`, t.RenderBodyColumns(), `</tbody>`)
	table := u.JoinStr(tHead, tBoday)
	// table := tHead

	return table
}

func (table *CardTableBody) SystemLogsTable() string {
	var compounDataTable strings.Builder

	var TableMap map[string]interface{}

	for i, column := range table.Data.([]*model.SystemLog) {
		var bgColor = ""

		if i%2 == 0 {
			bgColor = ""
		}

		datadata, _ := json.Marshal(column)

		json.Unmarshal(datadata, &TableMap)

		s := strings.ReplaceAll(string(datadata), "'", "&#39;")

		switch column.MessageType {
		case "SUCCESS":
			TableMap["TextColor"] = "text-success"
		case "INFO":
			TableMap["TextColor"] = "text-warning"
		case "ERROR":
			TableMap["TextColor"] = "text-danger"
		case "DELETE":
			TableMap["TextColor"] = "text-danger"
		default:
		}

		compounDataTable.WriteString(u.JoinStr(`<tr id='`, table.ID, `-r-`, strconv.Itoa(i), `' data-data='`, s, `' class="tb-row p-3 btn-reveal-trigger fw-semi-bold" style="background:`, bgColor, `">`, table.loadColumsData(TableMap, i), `</tr>`))

	}
	return compounDataTable.String()
}

func getSumOfColData(ColName string, data interface{}) string {
	sum := 0
	dataType := reflect.ValueOf(data)
	if dataType.Kind() == reflect.Slice {
		for i := 0; i < dataType.Len(); i++ {
			item := dataType.Index(i)
			if item.Kind() == reflect.Ptr && item.Elem().Kind() == reflect.Struct {
				field := item.Elem().FieldByName(strings.TrimSpace(ColName))
				if field.IsValid() {
					if field.Kind() == reflect.Int {
						sum += int(field.Int())
					}
				} else {
					log.Printf("Field %s not found in struct %v", ColName, item.Elem().Type())
				}
			}
		}
	} else {
		log.Println("Error: Data is not a slice")
	}
	return strconv.Itoa(sum)
}

func isToday(DemandDateTime interface{}) bool {
	demandDate, ok := DemandDateTime.(string)
	if !ok {
		log.Println("Invalid date format")
		return false
	}
	parsedDate, err := time.Parse("2006-01-02 15:04:05 +0000 UTC", demandDate)
	if err != nil {
		log.Println("Error parsing date:", err)
		return false
	}
	today := time.Now().Truncate(24 * time.Hour)
	return parsedDate.Truncate(24 * time.Hour).Equal(today)
}

func (table *CardTableBody) ProdLineManagementTable() string {
	var compounDataTable strings.Builder

	var TableMap map[string]interface{}

	for i, column := range table.Data.([]*model.ProdLine) {
		var bgColor = ""

		datadata, _ := json.Marshal(column)

		json.Unmarshal(datadata, &TableMap)

		s := strings.ReplaceAll(string(datadata), "'", "&#39;")

		compounDataTable.WriteString(u.JoinStr(`<tr id='`, table.ID, `-r-`, strconv.Itoa(i), `' data-data='`, s, `' class="tb-row p-3 btn-reveal-trigger fw-semi-bold" style="background:`, bgColor, `">`, table.loadColumsData(TableMap, i), `</tr>`))

	}
	return compounDataTable.String()
}

func (table *CardTableBody) ProdProcessesManagementTable() string {
	var compounDataTable strings.Builder

	var TableMap map[string]interface{}

	for i, column := range table.Data.([]*model.ProdProcess) {
		var bgColor = ""

		datadata, _ := json.Marshal(column)

		json.Unmarshal(datadata, &TableMap)

		s := strings.ReplaceAll(string(datadata), "'", "&#39;")

		compounDataTable.WriteString(u.JoinStr(`<tr id='`, table.ID, `-r-`, strconv.Itoa(i), `' data-data='`, s, `' class="tb-row p-3 btn-reveal-trigger fw-semi-bold" style="background:`, bgColor, `">`, table.loadColumsData(TableMap, i), `</tr>`))

	}
	return compounDataTable.String()
}

func (table *CardTableBody) RecipeManagementTable() string {
	var compounDataTable strings.Builder

	var TableMap map[string]interface{}

	for i, column := range table.Data.([]*model.Recipe) {
		var bgColor = ""

		datadata, _ := json.Marshal(column)
		rowData, _ := json.Marshal(struct {
			CompoundName string
			CompoundCode string
		}{
			CompoundName: column.CompoundName,
			CompoundCode: column.CompoundCode,
		})
		json.Unmarshal(datadata, &TableMap)

		s := strings.ReplaceAll(string(rowData), "'", "&#39;")

		compounDataTable.WriteString(u.JoinStr(`<tr id='`, table.ID, `-r-`, strconv.Itoa(i), `' data-data='`, s, `' class="tb-row p-3 btn-reveal-trigger fw-semi-bold" style="background:`, bgColor, `">`, table.loadColumsData(TableMap, i), `</tr>`))

	}
	return compounDataTable.String()
}

func (table *CardTableBody) StageManagementTable() string {
	var compounDataTable strings.Builder

	var TableMap map[string]interface{}

	for i, column := range table.Data.([]*model.Stage) {
		var bgColor = ""

		datadata, _ := json.Marshal(column)

		json.Unmarshal(datadata, &TableMap)

		s := strings.ReplaceAll(string(datadata), "'", "&#39;")

		compounDataTable.WriteString(u.JoinStr(`<tr id='`, table.ID, `-r-`, strconv.Itoa(i), `' data-data='`, s, `' class="tb-row p-3 btn-reveal-trigger fw-semi-bold" style="background:`, bgColor, `">`, table.loadColumsData(TableMap, i), `</tr>`))

	}
	return compounDataTable.String()
}

func (table *CardTableBody) AllKanbanViewTable() string {
	var sb strings.Builder

	var TableMap map[string]interface{}

	for i, column := range table.Data.([]*model.AllKanbanViewTable) {

		var bgColor = ""

		datadata, _ := json.Marshal(column)

		json.Unmarshal(datadata, &TableMap)

		json.Unmarshal(datadata, &TableMap)

		sb.WriteString(u.JoinStr(`<tr id='`, table.ID, `-r-`, strconv.Itoa(i), `' class="tb-row p-3 btn-reveal-trigger fw-semi-bold" style="background:`, bgColor, `">`, table.loadColumsData(TableMap, i), `</tr>`))

	}
	return sb.String()
}

func (table *CardTableBody) OperatorManagementTable() string {
	var compounDataTable strings.Builder

	var TableMap map[string]interface{}

	for i, column := range table.Data.([]*model.Operator) {
		var bgColor = ""

		datadata, _ := json.Marshal(column)

		json.Unmarshal(datadata, &TableMap)

		s := strings.ReplaceAll(string(datadata), "'", "&#39;")

		compounDataTable.WriteString(u.JoinStr(`<tr id='`, table.ID, `-r-`, strconv.Itoa(i), `' data-data='`, s, `' class="tb-row p-3 btn-reveal-trigger fw-semi-bold" style="background:`, bgColor, `">`, table.loadColumsData(TableMap, i), `</tr>`))

	}
	return compounDataTable.String()
}
func (table *CardTableBody) MaterialManagementTable() string {
	var compounDataTable strings.Builder

	var TableMap map[string]interface{}

	for i, column := range table.Data.([]*model.RawMaterial) {
		var bgColor = ""

		datadata, _ := json.Marshal(column)

		json.Unmarshal(datadata, &TableMap)

		s := strings.ReplaceAll(string(datadata), "'", "&#39;")

		compounDataTable.WriteString(u.JoinStr(`<tr id='`, table.ID, `-r-`, strconv.Itoa(i), `' data-data='`, s, `' class="tb-row p-3 btn-reveal-trigger fw-semi-bold" style="background:`, bgColor, `">`, table.loadColumsData(TableMap, i), `</tr>`))

	}
	return compounDataTable.String()
}
func (table *CardTableBody) ChemicalManagementTable() string {
	var compounDataTable strings.Builder

	var TableMap map[string]interface{}

	for i, column := range table.Data.([]*model.ChemicalTypes) {
		var bgColor = ""

		datadata, _ := json.Marshal(column)

		json.Unmarshal(datadata, &TableMap)

		s := strings.ReplaceAll(string(datadata), "'", "&#39;")

		compounDataTable.WriteString(u.JoinStr(`<tr id='`, table.ID, `-r-`, strconv.Itoa(i), `' data-data='`, s, `' class="tb-row p-3 btn-reveal-trigger fw-semi-bold" style="background:`, bgColor, `">`, table.loadColumsData(TableMap, i), `</tr>`))

	}
	return compounDataTable.String()
}
