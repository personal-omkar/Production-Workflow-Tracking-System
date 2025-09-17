package services

import (
	"strings"

	u "irpl.com/kanban-commons/utils"
)

type FormField struct {
	Label           string
	ID              string
	Width           string
	Icon            string
	Type            string
	DataType        string
	Options         []string
	DateOpts        string
	Readonly        bool
	IsRequired      bool
	Placeholder     string
	Style           string
	Value           string
	AdditionalAttr  string
	DropDownOptions []DropDownOptions
	Min             string
	Step            string
	Max             string
	Hidden          bool
	SubForms        []FormSection
}

type FormSection struct {
	ID          string
	Fields      []FormField
	ExtraField  string
	SubFormName string
	Data        string
}

type Form struct {
	Title    string
	Sections []FormSection
	Buttons  string
	Alert    bool
	Style    string
}
type DropDownOpts struct {
	Value    string
	Text     string
	Selected bool
}

func (f *Form) GenerateForm() string {
	var result string

	var alert = u.JoinStr(`
	<div id="form-notification" class="alert alert-success alert-dismissible fade show mb-0 mt-1 p-0" role="alert" style="display: none;"></div>
	`)
	// Start the page container
	result = u.JoinStr(`
	<div class="container mt-0" `, f.Style, `>
		<div class="card">
			<div class="card-header d-flex justify-content-start align-items-start">
				<h3 class="heading-text mb-0" style="margin-right:1rem;">`, f.Title, `</h3>
				`, u.IF(f.Alert, alert), `
			</div>
			<div class="card-body" data-group="`, strings.ToLower(strings.ReplaceAll(f.Title, " ", "")), `">
	`)

	for _, section := range f.Sections {
		result = u.JoinStr(result, section.GenerateSection())
	}

	result = u.JoinStr(result, `
			</div>
			<div class="card-footer text-muted d-flex justify-content-end  align-items-center">
				`, f.Buttons, `
			</div>
		</div>
	</div>
	`)

	return result
}

func (section *FormSection) GenerateSection() string {
	var result string

	result = u.JoinStr(`
	<div id="`, section.ID, `" subform-name="`, section.SubFormName, `" class="container-fluid mb-0 p-0" style="" `, section.Data, `>
		<div class="container-fluid p-0 m-0">
			<div id="field-container-`, section.ID, `">
				<div class="row container-fluid m-0 p-0">    
	`)

	// Generate all fields in the section
	for _, field := range section.Fields {
		result = u.JoinStr(result, field.GenerateField())
	}

	// Close the row and section container
	result = u.JoinStr(result, section.ExtraField, `
				</div>
			</div>
		</div>
	</div>`)

	return result
}

func (field *FormField) GenerateField() string {
	var result string
	var col string
	switch field.Width {
	case "hidden":
		col = "d-none"
	case "auto":
		col = "col"
	case "7%":
		col = "col-1"
	case "15%":
		col = "col-2"
	case "25%":
		col = "col-3"
	case "33%":
		col = "col-4"
	case "50%":
		col = "col-6"
	case "75%":
		col = "col-9"
	case "100%":
		col = "col-12"
	}

	// Start the field container
	result = u.JoinStr(`<div class="`, col, ` mb-0">`)

	if !(field.Type == "button") {
		if field.Label != "" {
			// Label with required indicator
			result = u.JoinStr(result, `
			<div class="col mb-0">
			<label class="container-fluid p-0 my-1" for="`, field.ID, `">`, field.Label)
			if field.IsRequired {
				result = u.JoinStr(result, `<span style="color: red;">*</span>`)
			}
			result = u.JoinStr(result, `</label> </div>`)
		} else {
			result = u.JoinStr(result, `
				<div class="col-4 mb-0">
				<label class="container-fluid p-0 my-1"></label></div>`)
		}
	}

	// Generate the input/select element
	switch field.Type {
	case "select":
		disabled := ""
		if field.Readonly {
			disabled = "disabled"
		}

		result = u.JoinStr(result, `<select class="form-select" id="`, field.ID, `" data-name="`, field.ID, `" data-type="`, field.DataType, `"  style="`, field.Style, `" `, u.IF(field.IsRequired, "data-validate='required'"), " ", disabled, " ", field.AdditionalAttr, `>`)
		for _, option := range field.DropDownOptions {
			selected := ""
			if option.Selected {
				selected = "selected"
			}
			result = u.JoinStr(result, `<option value="`, option.Value, `" `, selected, `>`, option.Text, `</option>`)
		}
		result = u.JoinStr(result, `</select>`)

	case "upload":
		readonly := ""
		if field.Readonly {
			readonly = "disabled"
		}

		result = u.JoinStr(result, `
				<div class="input-group mb-2">
					<span class="input-group-text p-0" onclick="$('#`, field.ID, `').click()">`, field.Icon, `</span>
					<input class="form-control datetimepicker" id="`, field.ID, `-value" data-name="`, field.ID, `" type="text" data-type="`, field.DataType, `" placeholder="`, field.Placeholder, `" style="`, field.Style, `" value="`, field.Value, `" `, u.IF(field.IsRequired, "data-validate='required'"), `  `, readonly, " ", field.AdditionalAttr, `> 
					<input type="file" class="form-control" id="`, field.ID, `" >
				</div>`)

	case "text":

		readonly := ""
		if field.Readonly {
			readonly = "disabled"
		}

		if field.DataType == "date" {
			var options string
			dateOptions := map[string]string{
				"H:i":        `{"enableTime":true,"noCalendar":true,"dateFormat":"H:i","disableMobile":true, "time_24hr": true`,
				"h:i K":      `{"enableTime":true,"noCalendar":true,"dateFormat":"h:i K","disableMobile":true`,
				"H:i:S":      `{"enableTime":true,"noCalendar":true,"dateFormat":"H:i:S","disableMobile":true, "time_24hr": true, "enableSeconds": true`,
				"d-m-y H:i":  `{"enableTime":true,"dateFormat":"d-m-y H:i","disableMobile":true`,
				"d-m-Y H:i":  `{"enableTime":true,"dateFormat":"d-m-Y H:i","disableMobile":true`,
				"y-m-d H:i":  `{"enableTime":true,"dateFormat":"y-m-d H:i","disableMobile":true`,
				"Y-m-d H:i":  `{"enableTime":true,"dateFormat":"Y-m-d H:i","disableMobile":true`,
				"dd/mm/yyyy": `{"disableMobile":true`,
				"default":    `{"disableMobile":true`,
			}

			placeholderKey := field.Placeholder
			if _, ok := dateOptions[placeholderKey]; !ok {
				placeholderKey = "default"
			}

			options = u.JoinStr(dateOptions[placeholderKey], u.IF(len(field.DateOpts) > 0, ",", field.DateOpts), `}`)

			result = u.JoinStr(result, `
				<div class="input-group mb-2">
					<span class="input-group-text p-0" onclick="$('#`, field.ID, `').click()">`, field.Icon, `</span>
					<input class="form-control datetimepicker" id="`, field.ID, `" data-name="`, field.ID, `" type="text" data-type="`, field.DataType, `" placeholder="`, field.Placeholder, `" style="`, field.Style, `" value="`, field.Value, `" `, u.IF(field.IsRequired, "data-validate='required'"), ` data-options='`, options, `' `, readonly, " ", field.AdditionalAttr, `> 
				</div>`)
		} else {
			result = u.JoinStr(result, `<input id="`, field.ID, `" data-name="`, field.ID, `" type="text"  data-type="`, field.DataType, `" class="form-control p-2" placeholder="`, field.Placeholder, `" style="`, field.Style, `" value="`, field.Value, `" `, u.IF(field.IsRequired, "data-validate='required'"), " ", readonly, " ", field.AdditionalAttr, `>`)
		}

	case "dropdown":
		disabled := ""
		if field.Readonly {
			disabled = "disabled"
		}

		result = u.JoinStr(result, `<select class="form-select selectpicker" id="`, field.ID, `" data-name="`, field.ID, `"  name="`, field.ID, `" placeholder="Search `, field.Label, `" data-type="`, field.DataType, `" style="`, field.Style, `" `, u.IF(field.IsRequired, "data-validate='required'"), " ", disabled, " ", field.AdditionalAttr, `>`)
		for _, option := range field.Options {
			result = u.JoinStr(result, `<option value="`, option, `">`, option, `</option>`)
		}
		result = u.JoinStr(result, `</select>`)
	case "dropdownfield":
		disabled := ""
		if field.Readonly {
			disabled = "disabled"
		}

		result = u.JoinStr(result, `<select class="form-select p-2" id="`, field.ID, `" data-name="`, field.ID, `" data-type="`, field.DataType, `"  style="`, field.Style, `" `, u.IF(field.IsRequired, "data-validate='required'"), " ", disabled, " ", field.AdditionalAttr, `>`)
		for _, option := range field.DropDownOptions {
			selected := ""
			if option.Selected {
				selected = "selected"
			}
			result = u.JoinStr(result, `<option value="`, option.Value, `" `, selected, `>`, option.Text, `</option>`)
		}
		result = u.JoinStr(result, `</select>`)

	case "number":
		readonly := ""
		if field.Readonly {
			readonly = "disabled"
		}
		result = u.JoinStr(result, `<input id="`, field.ID, `" data-name="`, field.ID, `" type="number"  data-type="`, field.DataType, `" class="form-control" placeholder="`, field.Placeholder, `" style="`, field.Style, `" value="`, field.Value, `"  min="`, field.Min, `"max="`, field.Max, `" step="`, field.Step, `"`, u.IF(field.IsRequired, "data-validate='required'"), " ", readonly, " ", field.AdditionalAttr, ` oninput="this.value = this.value === '' ? '' : Math.min(Math.max(this.value, `, field.Min, `), `, field.Max, `);">`)

	case "hidden":
		readonly := ""
		if field.Readonly {
			readonly = "disabled"
		}
		result = u.JoinStr(result, `
				<div class="input-group mb-2" style="display:none">
					<span class="input-group-text p-0" onclick="$('#`, field.ID, `').click()">`, field.Icon, `</span>
					<input class="form-control" id="`, field.ID, `" data-name="`, field.ID, `" type="text" data-type="`, field.DataType, `" placeholder="`, field.Placeholder, `" style="`, field.Style, `" value="`, field.Value, `" `, u.IF(field.IsRequired, "data-validate='required'"), readonly, " ", field.AdditionalAttr, `> 
				</div>`)
	case "email":
		readonly := ""
		if field.Readonly {
			readonly = "disabled"
		}
		result = u.JoinStr(result, `<input id="`, field.ID, `" data-name="`, field.ID, `" type="email"  data-type="`, field.DataType, `" class="form-control p-2" placeholder="`, field.Placeholder, `" style="`, field.Style, `" value="`, field.Value, `" `, u.IF(field.IsRequired, "data-validate='required'"), " ", readonly, " ", field.AdditionalAttr, `>`)

	case "textarea":
		disabled := ""
		if field.Readonly {
			disabled = "disabled"
		}
		result = u.JoinStr(result, `<textarea class="form-control" rows="5" id="`, field.ID, `" data-name="`, field.ID, `"  name="`, field.ID, `"  placeholder="Search `, field.Label, `" data-type="`, field.DataType, `" style="`, field.Style, `" `, u.IF(field.IsRequired, "data-validate='required'"), " ", disabled, " ", field.AdditionalAttr, `>`, field.Value)
		result = u.JoinStr(result, `</textarea>`)

	case "button":
		disabled := ""
		if field.Readonly {
			disabled = "disabled"
		}
		result = u.JoinStr(result, `<button id="`, field.ID, `" data-name="`, field.ID, `"  name="`, field.ID, `"  style="`, field.Style, `" `, disabled, " ", field.AdditionalAttr, `>`, field.Label, `</button>`)

	case "div":
		disabled := ""
		if field.Readonly {
			disabled = "disabled"
		}
		result = u.JoinStr(result, `<div id="`, field.ID, `" data-name="`, field.ID, `"  name="`, field.ID, `"  style="`, field.Style, `" `, disabled, " ", field.AdditionalAttr, `>`, field.Value, `</div>`)

	case "sectional-card":

		var Sections string
		var tabs string
		var active string
		var selected string
		for i, v := range field.SubForms {
			if i == 0 {
				active = "active"
				selected = "true"
			} else {
				active = ""
				selected = "false"
			}
			sec := u.JoinStr(`
				<div class="card-body p-0">
					<div class="tab-content">
						<div class="tab-pane `, active, `" id="`, v.ID, `" role="tabpanel" aria-labelledby="`, v.ID, `">
							<div class="z-index-1" id="`, v.ID, `" >
								<div class="px-0 py-0">
								
									<div class="" style="max-height: 600px; overflow-y: auto;overflow-x: hidden;">
										`, v.GenerateGroupedSectionWithCondition(func(field FormField) bool {
				return field.Type == "button" &&
					(strings.HasSuffix(field.ID, "add-stg-fields") || strings.HasSuffix(field.ID, "delete-stg-fields"))
			}), `
									</div>
								</div>
							</div>
						</div>
					</div>
			`)
			Sections = Sections + sec
			tab := u.JoinStr(`			
			<li class="nav-item `, active, `" role="presentation">
				<a class="nav-link  mb-0 `, active, `" role="tab" id="`, v.ID, `-tab" data-bs-toggle="tab" href="#`, v.ID, `" aria-controls="`, v.ID, `" aria-selected="`, selected, `">
					<div class="d-flex gap-1 py-1 pe-3"> 
						<div class="ms-2">
							<h5 class="mb-0 lh-1">`, v.SubFormName, `</h5>
						</div>
					</div>
			 	</a>
			</li>`)
			tabs = tabs + tab
		}
		result = u.JoinStr(result, `         
		
		<div class="card-header p-0 ">
		  <ul class="nav nav-tabs border-0 top-courses-tab flex-nowrap" role="tablist">
			`, tabs, `	
		  </ul>
		</div>
		`, Sections, `
			
		</div>`)
	}
	// Close the field container
	result = u.JoinStr(result, `</div>`)

	return result
}

func (section *FormSection) GenerateGroupedSectionWithCondition(condition func(FormField) bool) string {
	var result strings.Builder

	result.WriteString(`<div id="` + section.ID + `" subform-name="` + section.SubFormName + `" class="container-fluid mb-0 p-0" style="" ` + section.Data + `>
		<div class="container-fluid p-0 m-0">
			<div id="field-container-` + section.ID + `">
	`)

	var group []FormField
	for _, field := range section.Fields {
		group = append(group, field)

		if condition(field) {
			result.WriteString(`<div class="row container-fluid m-0 p-0">`)
			for _, f := range group {
				result.WriteString(f.GenerateField())
			}
			result.WriteString(`</div>`)
			group = []FormField{}
		}
	}

	if len(group) > 0 {
		result.WriteString(`<div class="row container-fluid m-0 p-0">`)
		for _, f := range group {
			result.WriteString(f.GenerateField())
		}
		result.WriteString(`</div>`)
	}

	result.WriteString(`</div></div></div>`)
	return result.String()
}
