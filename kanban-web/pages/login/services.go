package login

import "irpl.com/kanban-commons/utils"

type LoginForm struct {
	FormID        string
	FormAction    string
	Inputfield    []Inputfield
	Dropdownfield []Dropdownfield
	Buttons       FooterButtons
}

type Inputfield struct {
	Visible  bool
	Label    string
	Type     string
	Name     string
	Required bool
	Value    string
	ID       string
	Function string
	DataType string
}

type Dropdownfield struct {
	ID       string
	Name     string
	Options  Dropdownoptions
	DataType string
}

type Dropdownoptions struct {
	Option []string
}

type FooterButtons struct {
	BtnType        string
	BtnSubmitGroup string
	BtnID          string
	Text           string
}

func (mf *LoginForm) Build() string {
	model := utils.JoinStr(`
	`, mf.FormRender(), `
	`)
	return model
}

func (mf *LoginForm) FormRender() string {
	form := utils.JoinStr(`
			<form action="`, mf.FormAction, `" method="post" data-group="`, mf.FormID, `" id="`, mf.FormID, `">
				`, mf.DropdownRender(), `
				`, mf.InputFieldRender(), `	
				`, mf.FormButtons(), `
			</form>
			`)
	return form

}

func (mf *LoginForm) InputFieldRender() string {
	var inputFields string
	var required string
	var imp string

	for _, input := range mf.Inputfield {
		if input.Required {
			required = "required"
			imp = "*"
		} else {
			required = ""
			imp = ""
		}

		inputField := utils.JoinStr(`

		<div class="mb-3">
            <label class="form-label" for="`, input.ID, `">`, input.Label, " ", imp, `</label>
            <input class="form-control" id="`, input.ID, `" data-type="`, input.DataType, `"  data-name="`, input.Name, `" type=`, input.Type, `  `, required, ` class="form-control" id="`, input.ID, `" name="`, input.Name, `" />
        </div>
		
		`)
		inputFields = inputFields + inputField
	}
	return inputFields
}

func (mf *LoginForm) DropdownRender() string {
	var options string
	var dropDownFields string

	for _, dropdown := range mf.Dropdownfield {
		for _, option := range dropdown.Options.Option {
			opt := utils.JoinStr(`
			<option value="`, option, `" class="form-control" aria-label="Sizing example input" aria-describedby="inputGroup-sizing-default">`, option, `</option>
			`)
			options = options + opt
		}
		dropDownField := utils.JoinStr(`
		<div class="col-12 mb-3">
		 	<label class="form-label" for="`, dropdown.ID, `">`, dropdown.Name, ` *</label>
			<div class="input-group">
				<select id="`, dropdown.ID, `" name="`, dropdown.Name, `" class="form-control" aria-label="Sizing example input" aria-describedby="inputGroup-sizing-default">
				`, options, `
				</select>
			</div>
		</div>
		`)
		dropDownFields = dropDownFields + dropDownField
	}
	return dropDownFields

}
func (f *LoginForm) FormButtons() string {
	var button string
	if f.Buttons.Text == "Login" {
		btn := utils.JoinStr(`
		<div class="row flex-between-center">
			<div class="col-auto">
				<div class="form-check mb-0">
					<input class="form-check-input" type="checkbox" id="card-checkbox">
					<label class="form-check-label mb-0" for="card-checkbox">Remember me</label>
				</div>
			</div>
			<div class="col-auto">
				<a class="fs--1" href="#">Forgot Password?</a>
			</div>
		</div>
		 <div class="" style="margin-top: 2rem;">
            <button class="btn d-block w-100 mt-3" id="`, f.Buttons.BtnID, `" type=`, f.Buttons.BtnType, ` data-submit="`, f.FormID, `" name=`, f.Buttons.BtnID, ` style="background:#CF7AC2;color:white">`, f.Buttons.Text, `</button>
        </div>
		<div class="mt-3">
			<a href="/register" style="color:#CF7AC2"><b>Register?</b></a>
		</div>
		
		`)
		button = button + btn
	} else {
		btn := utils.JoinStr(`
		<div class="modal-footer">
			<button type="`, f.Buttons.BtnType, `"  id="`, f.Buttons.BtnID, `"class="btn btn-primary">`, f.Buttons.Text, `</button>
		</div>
		`)
		button = button + btn
	}

	return button
}
