package services

import (
	u "irpl.com/kanban-commons/utils"
)

type Card struct {
	Header     string
	Body       string
	Footer     string
	Width      string
	Style      string
	CustomAttr string
}

func (c *Card) Build() string {
	header := ""
	if len(c.Header) > 0 {
		header = c.Header
	}

	footer := ""
	if len(footer) > 0 {
		footer = c.Footer
	}

	card := u.JoinStr(`
	<div class="`, c.Width, `" style="`, c.Style, `">
		<div class="card">
			`, header, `
			`, c.Body, `
			`, footer, `
		</div>
	</div>	
	`)
	return card
}

// func (mf *CardForm) InputFieldRender() string {

// 	var inputFields string
// 	var readonly string

// 	for _, input := range mf.Inputfield {
// 		// var col string
// 		if input.Readonly {
// 			readonly = "readonly"
// 		} else {
// 			readonly = ""
// 		}
// 		temp := u.JoinStr(`
// 		<div class="`, input.Width, `">
// 		<div class="input-group  mb-3">
// 				<label for="`, input.ID, `" class="col-form-label w-100"> `, input.Label, ` </label>
// 				<input type="`, input.Type, `" id="`, input.ID, `"  class="form-control p-2" name=`, input.Name, `  value="`, input.Value, `" placeholder="`, input.Placeholder, `" aria-label="Username" aria-describedby="basic-addon1" `, readonly, `>
// 			</div>
// 		</div>
// 		`)
// 		inputFields = inputFields + temp
// 	}

// 	return inputFields
// }

// func (mf *CardForm) DropdownRender() string {
// 	var options string
// 	var dropDownFields string

// 	for _, dropdown := range mf.Dropdownfield {
// 		for _, option := range dropdown.Options {
// 			opt := u.JoinStr(`
// 			<option value="`, option.Value, `" class="form-control" aria-label="Sizing example input" aria-describedby="inputGroup-sizing-default">`, option.Text, `</option>
// 			`)
// 			options = options + opt
// 		}
// 		dropDownField := u.JoinStr(`
// 			<div class="`, dropdown.Width, ` mb-3">
// 				<div class="input-group">
// 					<label for="`, dropdown.ID, `" class="col-form-label w-100">`, dropdown.Label, `</label>
// 					<select id="`, dropdown.ID, `" name="`, dropdown.Name, `" class="form-control  p-2" aria-label="Sizing example input" aria-describedby="inputGroup-sizing-default">
// 					<option value="select" class="form-control" aria-label="Sizing example input" aria-describedby="inputGroup-sizing-default">Select</option>
// 					`, options, `
// 					</select>
// 				</div>
// 			</div>

// 		`)
// 		dropDownFields = dropDownFields + dropDownField
// 	}
// 	return dropDownFields

// }
