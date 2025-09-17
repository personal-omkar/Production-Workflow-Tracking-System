package services

import (
	"fmt"
	"strconv"
	"strings"

	"irpl.com/kanban-commons/utils"
	u "irpl.com/kanban-commons/utils"
)

type ModelCard struct {
	ID      string
	Type    string
	Form    ModelForm
	Heading string
}

type ModelForm struct {
	FormID            string
	FormAction        string
	Accordion         []AccordionList
	Inputfield        []InputAttributes
	Dropdownfield     []DropdownAttributes
	TextArea          []TextAreaAttributes
	IncrementalButton IncrementalButtonAttributes
	Footer            Footer
	Checkbox          []Checkbox
	CheckBoxList      CheckBox
	Info              InfoLine
	InputGroup        []InputGroups
	Details           Details
	QualityDetails    Details
	DispatchDetails   Details
	ExtractFields     []ExtractFields
}

type ExtractFields struct {
	Heading      string
	HTML         string
	HeadingStyle string
	Style        string
	Class        string
}
type Details struct {
	Lable       string
	Style       string
	DetailsData []DetailsData
	Note        string
}

type DetailsData struct {
	Heading string
	Data    []Data
	Style   string
}

type Data struct {
	Lable string
	Value string
}

type InputGroups struct {
	Title                string
	InputGroupAttributes []InputGroupAttributes
}

type InputGroupAttributes struct {
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
}

type InfoLine struct {
	IsVisible bool
	ID        string
	Lable     string
	Data      []string
}

type Checkbox struct {
	ID       string
	Name     string
	Width    string
	Height   string
	Label    string
	Visible  bool
	CheckBox TextAreaAttributes
}

type CheckBox struct {
	Label        string
	Width        string
	CheckBoxList []CheckboxAttribut
}

type CheckboxAttribut struct {
	PageLink  string
	PageName  string
	IsChecked bool
	Readonly  bool
}
type AccordionAttributes struct {
	Name     string
	Sections []FormSection
}
type IncrementalButtonAttributes struct {
	DataType  string
	Label     string
	IsVisible bool
	Width     string
	MinValue  int
	MaxValue  int
}
type AccordionList struct {
	Lable string
	List  []AccordionAttributes
}

type AccordionrowAttribute struct {
	Inputfield    []InputAttributes
	Dropdownfield []DropdownAttributes
}
type Footer struct {
	Buttons   []FooterButtons
	CancelBtn bool
}
type FooterButtons struct {
	BtnType        string
	BtnID          string
	Text           string
	Disabled       bool
	Style          string
	DataSubmitName string
}

type ConfirmationModal struct {
	For    string
	ID     string
	Title  string
	Body   []string
	Footer Footer
}

type InfoModal struct {
	ID        string
	Title     string
	ModelSize string
	Body      []string
	Footer    Footer
}

func (mf *ModelCard) Build() string {

	model := utils.JoinStr(`
	<div class="container">
		<div class="modal fade" id="`, mf.ID, `" data-bs-backdrop="static" data-bs-keyboard="false" tabindex="-1" aria-labelledby="staticBackdropLabel" aria-hidden="true">
			<div class="modal-dialog modal-dialog-centered `, mf.Type, `">
				<div class="modal-content">
					<div class="modal-header">
						<h5 class="modal-title" id="staticBackdropLabel" ><b style="color:#871A83">`, mf.Heading, `</b></h5>
						<button type="button" class="btn-close" data-bs-dismiss="modal"  aria-label="Close" ></button>
					</div>
					<div class="modal-body">
						`, mf.Form.FormRender(), `
					</div>
				</div>
			</div>
		</div>
	</div>
	`)
	return model
}

func (mf *ModelForm) FormRender() string {
	var accordionStr string
	for _, accordion := range mf.Accordion {
		accordionStr += accordion.AccordionRender()
	}

	form := utils.JoinStr(`
			<form action="`, mf.FormAction, `" method="post" data-group="`, mf.FormID, `"  id="`, mf.FormID, `"  class="p-0">
			<div class="row">
				`, mf.InputFieldRender(), `	
				`, mf.TextAreaRender(), `	
				`, mf.DropdownRender(), `
				`, mf.IncrementalButtonRender(), `
				`, mf.InputGroupFieldRender(), `
				
				`, mf.CheckBoxs(), `
				`, mf.InfoRendender(), `
				`, mf.CheckBoxListRender(), `
				`, mf.Details.Build(), `
				`, mf.ExtractFieldsrender(), `

						
			</div>
			`, accordionStr, `
			`, mf.QualityDetails.Build(), `
			`, mf.DispatchDetails.Build(), `
			<div id="modelnotification" class="alert alert-success alert-dismissible fade show mb-0 mt-1 p-1 w-100" role="alert" style="display: none;"></div>
			`, mf.Footer.FormButtons(), `
			</form>
			
			`)
	return form

}

func (mf *ModelForm) InputFieldRender() string {

	var inputFields string
	var readonly string
	var hidden string
	var required string
	var datavalidate string
	var disabled string
	for _, input := range mf.Inputfield {
		if input.Hidden {
			hidden = `style="display:none"`
		} else {
			hidden = ``
		}
		// var col string
		if input.Readonly {
			readonly = "readonly"
		} else {
			readonly = ""
		}
		//required
		if input.Required {
			required = "required"
			datavalidate = "data-validate='required'"

		} else {
			required = ""
			datavalidate = ""
		}

		if input.Disabled {
			disabled = "disabled"
		} else {
			disabled = ""
		}
		// if input.Type == "text" {
		// 	if t, ok := input.Value.; ok {
		// 	}
		// }
		if input.Type != "upload" {
			temp := utils.JoinStr(`
			<div class="`, input.Width, ` " `, hidden, ` >
			<div class="input-group  mb-3">	
					<label for="`, input.ID, `" class="col-form-label w-100"> `, input.Label, ` </label>
					<input type="`, input.Type, `"  data-name="`, input.Name, `" data-type="`, input.DataType, `" id="`, input.ID, `"  class="form-control p-2" name=`, input.Name, ` `, datavalidate, `  value="`, input.Value, `" placeholder="`, input.Placeholder, `" aria-label="Username" aria-describedby="basic-addon1" `, required, ` `, readonly, ` `, input.AdditionalAttr, disabled, `>
				</div>
			</div>	
			`)

			inputFields = inputFields + temp
		} else {
			temp := u.JoinStr(`
			<div class="input-group mb-2">
				<span class="input-group-text p-0" onclick="$('#`, input.ID, `').click()">`, input.Icon, `</span>
				<input type="`, input.Type, `"  data-name="`, input.Name, `-value" data-type="`, input.DataType, `" id="`, input.ID, `-value"  class="form-control p-2" name=`, input.Name, ` `, datavalidate, `  value="`, input.Value, `" placeholder="`, input.Placeholder, `" aria-label="Username" aria-describedby="basic-addon1" `, required, ` `, readonly, ` `, input.AdditionalAttr, disabled, ` style="width:65%">
				<input type="file" class="form-control file-button-only"  id="`, input.ID, `" >
			</div>`)

			style := `
			<style>	
				.file-button-only {
				 	font-size: 0; 
					 color: transparent;
					 
				}

				.file-button-only::file-selector-button {
					font-size: 14px; /* Show the button text */
					padding: 10px 12px;
				
					
					border: none;
					
					cursor: pointer;
				}
								
			</style>	
					
			`
			inputFields = inputFields + temp + style
		}

	}
	return inputFields
}

func (mf *ModelForm) DropdownRender() string {
	var options string
	var dropDownFields string
	var selected string

	for _, dropdown := range mf.Dropdownfield {
		// Determine if the dropdown should be hidden
		hiddenClass := ""
		if dropdown.Hidden {
			hiddenClass = "d-none" // Bootstrap's `d-none` hides the element
		}

		// Check if the dropdown should be disabled
		disabled := ""
		if dropdown.Disabled {
			disabled = "disabled"
		}

		for _, option := range dropdown.Options {
			// Check if the option should be selected
			if option.Selected {
				selected = "selected"
			} else {
				selected = ""
			}

			// Create the option element
			opt := utils.JoinStr(`
                <option value="`, option.Value, `" class="form-control" 
                        aria-label="Sizing example input" 
                        aria-describedby="inputGroup-sizing-default" 
                        `, selected, `>`, option.Text, `</option>
            `)
			options = options + opt
		}

		// Create the dropdown field with the `d-none` class if hidden
		dropDownField := utils.JoinStr(`
            <div class="`, dropdown.Width, ` mb-3 `, hiddenClass, `">
                <label for="`, dropdown.ID, `" class="col-form-label">`, dropdown.Label, `</label>
                <div class="input-group">
                    <select id="`, dropdown.ID, `" 
                            data-name="`, dropdown.Name, `" 
                            data-type="`, dropdown.DataType, `" 
                            name="`, dropdown.Name, `" 
                            data-validate='required' 
                            class="form-control" 
                            aria-label="Sizing example input" 
                            aria-describedby="inputGroup-sizing-default" 
                            `, disabled, ` 
                            `, dropdown.AdditionalAttr, `>
                    `, options, `
                    </select>
                </div>
            </div>
        `)

		// Add the dropdown field to the result
		dropDownFields = dropDownFields + dropDownField

		// Reset options for the next dropdown
		options = ""
	}

	return dropDownFields
}

func (f *Footer) FormButtons() string {
	var button string
	var disabled string
	button = button + `<div class="modal-footer">`
	for _, data := range f.Buttons {
		if data.Disabled {
			disabled = "disabled"
		} else {
			disabled = ""
		}
		if f.CancelBtn {
			btn := utils.JoinStr(`<button type="button" class="btn btn-secondary btn-lg" data-bs-dismiss="modal" style="background-color:#c62f4a ;border:none" >Cancel</button>`)
			button = button + btn
		}
		var btn = ""
		if len(data.Style) > 0 {
			btn = utils.JoinStr(`
				<button type="`, data.BtnType, `" data-submit="`, data.DataSubmitName, `" id="`, data.BtnID, `"class="btn btn-primary btn-lg" style="`, data.Style, `" `, disabled, `>`, data.Text, `</button>`)
			button = button + btn
		} else {
			btn = utils.JoinStr(`
				<button type="`, data.BtnType, `" data-submit="`, data.DataSubmitName, `" id="`, data.BtnID, `"class="btn btn-primary btn-lg" style="background-color:#871A83 ;border:none" `, disabled, `>`, data.Text, `</button>`)
			button = button + btn
		}
	}
	button = button + `<div id="limit-warning" class="text-danger small ms-2 mt-2" style="display: none;"></div>`
	button = button + `</div>`
	return button
}
func DigitToWord(digit int) string {
	words := map[int]string{
		0: "zero", 1: "One", 2: "Two", 3: "Three", 4: "Four",
		5: "Five", 6: "Six", 7: "Seven", 8: "Eight", 9: "Nine",
	}
	return words[digit]
}

func (accList AccordionList) AccordionRender() string {
	var rows string

	if len(accList.List) != 0 {
		for i, acc := range accList.List {

			element := utils.JoinStr(`
		
		<div class="accordion-item">
			<h2 class="accordion-header" id="heading`, DigitToWord(i+1), `"  style="background-color:#f1f1f1">
			  <button class="accordion-button collapsed p-2" type="button" data-bs-toggle="collapse" data-bs-target="#collapse`, DigitToWord(i+1), `" aria-expanded="false" aria-controls="collapse`, DigitToWord(i+1), `">
					<span class="text-secondary">`, acc.Name, `</span>
			  </button>
			</h2>
			<div id="collapse`, DigitToWord(i+1), `" class="accordion-collapse collapse" aria-labelledby="heading`, DigitToWord(i+1), `" data-bs-parent="#accordionExample">
				<div class="accordion-body p-0 m-0 w-100">
					
	`)
			for _, field := range acc.Sections {
				element = utils.JoinStr(element, field.GenerateSection())
			}

			element = utils.JoinStr(element,
				`
				</div>
			</div>
		</div>
	`)
			rows = rows + element

		}
		row := utils.JoinStr(`
		<div class="w-100">
			<lable ><h4 class="fw-bold" style="color: #ab71a2;">`, accList.Lable, ` :</h4></lable>
		</div>
		<div class="accordion pb-2" id="accordionExample">
			`, rows, `
		</div>`)

		return row
	}
	return ""
}

func (mf *ModelForm) TextAreaRender() string {
	var readonly string
	var result string
	var disabled string
	for _, text := range mf.TextArea {
		if text.Readonly {
			readonly = "readonly"
		}
		if text.Disabled {
			disabled = "disabled"
		}
		field := utils.JoinStr(`<div class="mb-3 mt-2">
		<div class="mb-3">
		<label for="`, text.ID, `" class="col-form-label">`, text.Label, `</label>
		<textarea class="form-control" data-name="`, text.Name, `" id="`, text.ID, `" name="`, text.Name, `"  data-type="`, text.DataType, `" placeholder="`, text.Placeholder, `" data-validate='required' rows="4" `, readonly, ` `, disabled, ` ,">`, text.Value, `</textarea>
		</div>
		</div>`)
		result = result + field
	}
	return result

}

func (mf *ModelForm) IncrementalButtonRender() string {
	var result string
	if mf.IncrementalButton.IsVisible {
		result = utils.JoinStr(
			`<div class="container `, mf.IncrementalButton.Width, `">
				<div class="row">
					<div class="col-12 mt-4">
						<h6 class="form-label"></h6>
						<div class="d-flex align-items-center">
							<button type="button" id="btn-minus" class="quantity-left-minus btn btn-number me-2" 
								style="background-color:#871a83;color:white;width: 40px;">
								-
							</button>
							<input type="text" id="quantity" name="quantity" data-name="quantity" data-type="`, mf.IncrementalButton.DataType, `" 
								class="form-control text-center me-2" 
								value="0" min="`, strconv.Itoa(mf.IncrementalButton.MinValue), `" 
								max="`, strconv.Itoa(mf.IncrementalButton.MaxValue), `" 
								style="width: 80px;" oninput="validateQuantity(this)">
							<button type="button" id="btn-plus" class="quantity-right-plus btn btn-number" 
								style="background-color:#871a83;color:white;width: 40px;">
								+
							</button>
						</div>
					</div>
				</div>
				<script>
					function validateQuantity(input) {
						let min = parseInt(input.min, 10);
						let max = parseInt(input.max, 10);
						let value = parseInt(input.value, 10);

						if (isNaN(value) || value < min) {
							input.value = min;
						} else if (value > max) {
							input.value = max;
						}
					}
				</script>
			</div>`)
	}
	return result
}

func (mf *ModelForm) CheckBoxs() string {
	var html string
	if mf.Checkbox != nil {
		html = `
		<div class="input-group  mb-3">	
			<div class="d-flex justify-content-between w-100">
				<label for="form-check" class="col-form-label">Production Processes</label>
  				<button id="clear_check_box" class="btn btn-secondary me-3" disabled>Clear</button>
			</div>
			<div class=" scroll-div rounded overflow-auto w-100 bg-secondary bg-gradient align-items-center p-3" style="max-height: 40vh;">
			<div class="process-row d-flex justify-content-between">
            	<div class="col-6"><b>Processes Name</b></div> 
            	<div class="col-4 d-flex justify-content-start"><b>Group Name</b></div> 
        		<div class="col-2 d-flex justify-content-end"><b>Order</b></div>
			</div>
		`
		for _, text := range mf.Checkbox {
			if !text.Visible {
				html = ""
				return html
			}
			html += utils.JoinStr(`
				<!-- Checkbox and Label -->
				<div class="process-row d-flex justify-content-between">
					<div class="form-check d-flex align-items-center col-6">
						<input style="width: 20px; height: 20px;" class="form-check-input me-2" type="checkbox" value="" id="` + text.ID + `">
						<label class="form-check-label text-white m-0 pt-2 h-5" for="temporaryCheckbox" style="font-size:15px;">
							` + text.Name + `
						</label>
					</div>
					<div class="d-flex align-items-center col-4">
						<input style="height: 25px; border:none; outline:none;" class="me-2 GroupName" type="Text" value="" id="` + text.ID + `" disabled="">					
					</div>
					<!-- Square Input Field -->
					<div class="d-flex  align-items-center justify-content-end col-2">
						<input class="me-2 text-center processNo" type="text" value="" id="` + text.ID + `" style="width: 35px; height: 25px;" min="0" max="999" step="1" disabled="">
					</div>
				</div>
			
		`)
		}
		html += `</div> </div> 
`

	}
	return html
}

func (mf *ModelForm) InfoRendender() string {
	hiddenClass := ""
	if !mf.Info.IsVisible {
		hiddenClass = " d-none"
	}

	return `<div class="w-100 mb-3` + hiddenClass + `">
    	<span class="d-flex text-info w-100 justify-content-start p-0" id="` + mf.Info.ID + `" style="font-size:0.9rem;font-weight: bold;">` + mf.Info.Lable + `</span>
	</div>`
}

func (cm *ConfirmationModal) Build() string {
	style := ""
	if cm.For == "Submit" || cm.For == "submit" {
		style = "color:green;"
	} else if cm.For == "Delete" || cm.For == "delete" {
		style = "color:red;"
	}

	var bodyHTML string
	for _, line := range cm.Body {
		bodyHTML += line + "<br><br>"
	}

	html := fmt.Sprintln(`
    <!-- html -->
    <div class="modal fade" id="` + cm.ID + `" data-bs-backdrop="static" data-bs-keyboard="false" tabindex="-1" aria-labelledby="staticBackdropLabel" aria-hidden="true">
        <div class="modal-dialog modal-lg">
            <div class="modal-content">
                <div class="modal-header">
                    <h5 class="modal-title" id="staticBackdropLabel" style="` + style + `"><b>` + cm.Title + `</b></h5>
                    <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
                </div>
                <div class="modal-body" id="modal-body">
                    ` + bodyHTML + `
                </div>
            	` + cm.Footer.FormButtons() + `
            </div>
        </div>
    </div>
    <!-- !html -->`)

	return html
}

func (cm *InfoModal) Build() string {
	modelsize := ""
	if cm.ModelSize != "" {
		modelsize = cm.ModelSize
	} else {
		modelsize = "modal-lg"
	}
	var bodyHTML string
	for _, line := range cm.Body {
		bodyHTML += line
	}

	html := fmt.Sprintln(`
    <!-- html -->
    <div class="modal fade" id="` + cm.ID + `" data-bs-backdrop="static" data-bs-keyboard="false" tabindex="-1" aria-labelledby="staticBackdropLabel" aria-hidden="true">
        <div class="modal-dialog ` + modelsize + `">
            <div class="modal-content ">
                <div class="modal-header">
                    <h5 class="modal-title" id="staticBackdropLabel" style="color:#871a83"><b>` + cm.Title + `</b></h5>
                    <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
                </div>
                <div class="modal-body" id="modal-body">
                    ` + bodyHTML + `
                </div>
            	` + cm.Footer.FormButtons() + `
            </div>
        </div>
    </div>
    <!-- !html -->`)

	return html
}
func (mf *ModelForm) InputGroupFieldRender() string {
	var formContent string
	var readonly string
	var hidden string
	var required string
	var disabled string

	// Loop through each input group
	for _, group := range mf.InputGroup {
		// Render the group title in bold and red
		groupTitle := `
		<div class="section-title" style="font-weight: bold; color: #AB71A2;">` + group.Title + `</div>
		`

		// Initialize a variable for the input fields in this group
		var inputFields string

		// Loop through each input field in the group
		for _, input := range group.InputGroupAttributes {
			// Handle hidden, readonly, required, and disabled attributes for each input
			if input.Hidden {
				hidden = `style="display:none"`
			} else {
				hidden = ``
			}
			if input.Readonly {
				readonly = "readonly"
			} else {
				readonly = ""
			}
			if input.Required {
				required = "required"
			} else {
				required = ""
			}
			if input.Disabled {
				disabled = "disabled"
			} else {
				disabled = ""
			}

			// Create the HTML structure for each input field
			inputField := `
			<div class="` + input.Width + `">
				<div class="input-group mb-3" ` + hidden + `>
					<div class="input-group-prepend">
						<span class="input-group-text">` + input.Label + `</span>
					</div>
					<input type="` + input.Type + `" id="` + input.ID + `" class="form-control" name="` + input.Name + `" value="` + input.Value + `" placeholder="` + input.Placeholder + `" ` + required + ` ` + readonly + ` ` + input.AdditionalAttr + ` ` + disabled + `>
				</div>
			</div>
			`

			// Append the input field HTML to the group's input fields
			inputFields = inputFields + inputField
		}

		// Combine the group title and the input fields
		formContent = formContent + groupTitle + inputFields + `<br>`
	}

	return formContent
}

func (mf *ModelForm) CheckBoxListRender() string {
	var result string

	for _, checkbox := range mf.CheckBoxList.CheckBoxList {
		// var checkedJS string
		var readonly string
		var checked string
		if checkbox.Readonly {
			readonly = "disabled"
		}
		if checkbox.IsChecked {
			checked = "checked"
		}
		opt := utils.JoinStr(`
		
			<span class="d-flex d-inline-flex  align-items-center pl-1 m-1">
				<div class="col-auto border border-1 p-0 mx-2" style="border-color: #ab71a2 !important; border-radius:6px;">
					<span class="form-check p-0">
						<span class="border-end border-2 pl-1" style="border-color: #ab71a2 !important;">
							<input class="form-check-input m-1 mt-2 pl-1 checkbox-element" type="checkbox" `, readonly, ` `, checked, ` value="`, checkbox.PageLink, `" id="`, strings.ReplaceAll(checkbox.PageName, " ", ""), `">
						</span>
						<label class="form-check-label m-1 px-1" for="`, strings.ReplaceAll(checkbox.PageName, " ", ""), `">
							`, checkbox.PageName, `
						</label>
					</span>
				</div>
				
			</span>
		
	 `)

		result += opt

	}
	result = utils.JoinStr(` <span class="`, mf.CheckBoxList.Width, `" >`, `<label for="`, mf.CheckBoxList.Label, `" class="form-label w-100">`, mf.CheckBoxList.Label, `</label>`, result, `</span>`)
	return result
}

func (d *Details) Build() string {
	var result string
	var columHTML string
	var note string
	if d.Note != "" {
		note += `<div class="mx-1" style="white-space: nowrap; overflow: hidden; text-overflow: ellipsis;">Note : ` + d.Note + `</div>`
	}
	for _, data := range d.DetailsData {
		columHTML += ` <div class="p-1" style="flex: 1; margin: 0;">
						  	<h5 class="fw-bold" style="padding: 0; margin: 0; ` + data.Style + `">` + data.Heading + `</h5>`
		for _, data := range data.Data {
			columHTML += `
			<div class="d-flex" style="width:100%;">
				<div class="col-4" style="font-size:14px;">` + data.Lable + `</div> 
				<div class="col-1"> : </div> 
    			<div class="col-7" style="font-size:14px;">` + data.Value + `</div>
			</div>`
		}
		columHTML += `
						</div>`
	}
	if len(d.DetailsData) != 0 {

		Heading := ""
		if d.Lable != "" {
			Heading = utils.JoinStr(`<lable ><h4 class="fw-bold" style="color: #ab71a2;">` + d.Lable + ` :</h4></b></lable>`)
		}

		result += `
		<div class="row">
			<div class="w-100">
				` + Heading + `
			</div>
		`
		result += `
			<div class="d-flex" style="` + d.Style + `"> 
				` + columHTML + `
			</div>
		` + note + `
		</div>
		<hr>
		`
	}
	return result
}

func (mf *ModelForm) ExtractFieldsrender() string {
	var html string
	if len(mf.ExtractFields) != 0 {
		for _, data := range mf.ExtractFields {
			heading := utils.JoinStr(`<lable><h4 class="fw-bold" style="`, data.HeadingStyle, `">`, data.Heading, ` </h4></lable>`)
			html = utils.JoinStr(`
			<div class="row">
				<div class="w-100">
					`, heading, `
				</div>
				<div class="d-flex `, data.Class, `" style="`, data.Style, `"> 
					`, data.HTML, `
				</div>
			</div>
			`)
		}
	}
	return html
}
