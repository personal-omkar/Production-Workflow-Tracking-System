package recipemanagement

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"

	m "irpl.com/kanban-commons/model"
	"irpl.com/kanban-commons/utils"
	s "irpl.com/kanban-web/services"
)

func (r *RecipeManagement) EditDialogBuild() string {

	var recipe []*m.Recipe
	var stage []*m.Stage
	var recipeToStage []*m.RecipeToStage
	var compounds []*m.Compounds
	var rawMaterials []*m.RawMaterial
	var chemTypes []*m.ChemicalTypes
	var prodLine []*m.ProdLine
	var prodLineToRecipe []*m.ProdLineToRecipe
	var ProdLineToReciptId int

	//get compound data
	compResp, err := http.Get(RestURL + "/get-all-compound-data")
	if err != nil {
		slog.Error("%s - error - %s", "Error making GET request", err)
	}

	defer compResp.Body.Close()

	if err := json.NewDecoder(compResp.Body).Decode(&compounds); err != nil {
		slog.Error("error decoding response body", "error", err)
	}

	//get stages
	stageResp, err := http.Get(RestURL + "/get-all-stages")
	if err != nil {
		slog.Error("%s - error - %s", "Error making GET request", err)
	}

	defer stageResp.Body.Close()

	if err := json.NewDecoder(stageResp.Body).Decode(&stage); err != nil {
		slog.Error("error decoding response body", "error", err)
	}

	//get recipes
	recipeResp, err := http.Get(RestURL + "/get-recipe-by-param?key=compound_code&value=" + url.QueryEscape(r.CompoundCode))
	if err != nil {
		slog.Error("%s - error - %s", "Error making GET request", err)
	}

	defer recipeResp.Body.Close()

	if err := json.NewDecoder(recipeResp.Body).Decode(&recipe); err != nil {
		slog.Error("error decoding response body", "error", err)
	}

	var rawQuery m.RawQuery
	rawQuery.Host = utils.RestHost
	rawQuery.Port = utils.RestPort
	rawQuery.Type = "RecipeToStage"
	rawQuery.Query = `SELECT * FROM RecipeToStage where recipe_id = ` + fmt.Sprint(recipe[0].Id) + ` ;` //`;`
	rawQuery.RawQry(&recipeToStage)

	rawQuery.Type = "RawMaterial"
	rawQuery.Query = `SELECT * FROM Raw_Material ;` //`;`
	rawQuery.RawQry(&rawMaterials)

	rawQuery.Type = "ChemicalTypes"
	rawQuery.Query = `SELECT * FROM Chemical_Types ;` //`;`
	rawQuery.RawQry(&chemTypes)

	rawQuery.Type = "ProdLine"
	rawQuery.Query = `SELECT * FROM prod_line ;` //`;`
	rawQuery.RawQry(&prodLine)

	rawQuery.Type = "ProdLineToRecipe"
	rawQuery.Query = `SELECT * FROM prod_line_to_recipe where recipe_id = ` + fmt.Sprint(recipe[0].Id) + ` ;` //`;`
	rawQuery.RawQry(&prodLineToRecipe)

	//Edit Dialogue
	var opt []s.DropDownOptions
	for _, v := range compounds {
		var tempopt s.DropDownOptions
		tempopt.Text = v.CompoundName
		tempopt.Value = v.CompoundName
		if v.CompoundName == r.CompoundName {
			tempopt.Selected = true
		}
		opt = append(opt, tempopt)
	}

	var rawOpts []s.DropDownOptions
	for _, v := range rawMaterials {
		var tempopt s.DropDownOptions
		tempopt.Text = v.SCADACode
		tempopt.Value = v.SCADACode
		rawOpts = append(rawOpts, tempopt)
	}

	var chemTypeOpts []s.DropDownOptions
	for _, v := range chemTypes {
		var tempopt s.DropDownOptions
		tempopt.Text = v.Type
		tempopt.Value = v.Type
		chemTypeOpts = append(chemTypeOpts, tempopt)
	}

	var prodLineOpts []s.DropDownOptions
	for _, v := range prodLine {
		var tempopt s.DropDownOptions
		tempopt.Text = v.Name
		tempopt.Value = strconv.Itoa(v.Id)
		prodLineOpts = append(prodLineOpts, tempopt)
	}

	if len(prodLineToRecipe) > 0 {
		for i, v := range prodLineOpts {
			if v.Value == strconv.Itoa(prodLineToRecipe[0].ProdLineId) {
				prodLineOpts[i].Selected = true
			}
		}
		ProdLineToReciptId = prodLineToRecipe[0].Id
	}

	formBtn := utils.JoinStr(`
		<div class="col-md-12 d-flex justify-content-end">
			<button type="button" class="btn " data-bs-dismiss="modal" style="background-color:#871a83 ;border:none;color:white">Cancel</button>
			<button id="update-recipe" data-submit="editrecipemodel" data-url="/update-existing-recipe" class="btn ms-3" style="background-color:#871a83 ;border:none;color:white">Save</button>
		</div>
	`)

	var subformlist []s.FormSection
	for _, v := range stage {
		var fields []s.FormField
		var temp s.FormSection
		temp.ID = v.Name + "edit"
		temp.SubFormName = v.Name
		temp.Data = `data-id="` + strconv.FormatUint(uint64(v.ID), 10) + `"`

		// Parse headers
		var headers []map[string]string
		if err := json.Unmarshal(v.Headers, &headers); err != nil {
			log.Printf("error decoding headers: %v", err)
		}

		var rts *m.RecipeToStage
		for _, each := range recipeToStage {
			if each.StageId == int(v.ID) {
				rts = each
				break
			}
		}

		jsonArray := []map[string]string{}
		json.Unmarshal(rts.Data, &jsonArray)
		for i, data := range jsonArray {

			// Add fields from headers
			for _, header := range headers {

				if header["type"] == "dropdown-RM" {
					var custOpt []s.DropDownOptions
					for _, opt := range rawOpts {
						opt.Selected = (opt.Value == data[header["field"]+"-"+fmt.Sprint(i+1)])
						custOpt = append(custOpt, opt)
					}

					fields = append(fields,
						s.FormField{Label: header["field"], ID: header["field"] + "-" + fmt.Sprint(i+1), Width: "auto", Type: "select", DropDownOptions: custOpt},
					)
				} else if header["type"] == "dropdown-CT" {
					var custOpt []s.DropDownOptions
					for _, opt := range chemTypeOpts {
						opt.Selected = (opt.Value == data[header["field"]+"-"+fmt.Sprint(i+1)])
						custOpt = append(custOpt, opt)
					}

					fields = append(fields,
						s.FormField{Label: header["field"], ID: header["field"] + "-" + fmt.Sprint(i+1), Width: "auto", Type: "select", DropDownOptions: custOpt},
					)
				} else {
					fields = append(fields,
						s.FormField{Label: header["field"], ID: header["field"] + "-" + fmt.Sprint(i+1), Width: "auto", Type: "text", DataType: "string", Value: data[header["field"]+"-"+fmt.Sprint(i+1)]},
					)
				}

			}

			if len(fields) > 0 {
				fields = append(fields,
					s.FormField{Label: "Delete", ID: "delete-stg-fields", Width: "15%", Type: "button", AdditionalAttr: "class='btn btn-danger delete-btn w-100 mt-4' onclick='removeRecipeField(this)'"},
				)
			}
		}

		for _, header := range headers {

			if header["type"] == "dropdown" {

				fields = append(fields,
					s.FormField{Label: header["field"], ID: header["field"] + "-" + fmt.Sprint(len(jsonArray)+1), Width: "auto", Type: "select", DropDownOptions: rawOpts},
				)
			} else {
				fields = append(fields,
					s.FormField{Label: header["field"], ID: header["field"] + "-" + fmt.Sprint(len(jsonArray)+1), Width: "auto", Type: "text", DataType: "string", Placeholder: "Enter " + header["field"]},
				)
			}

		}

		if len(fields) > 0 {
			fields = append(fields,
				s.FormField{Label: "Add", ID: "add-stg-fields", Width: "15%", Type: "button", AdditionalAttr: "class='btn btn-primary w-100 mt-4' onclick='addRecipeFields(this)'"},
			)
		}

		temp.Fields = fields
		subformlist = append(subformlist, temp)
	}

	EditRecipeModel := s.Model{
		ID:    "EditRecipeModel",
		CSS:   `<style> .modal-xxl { width: 1355px; max-width: 1355px} </style>`,
		Type:  "modal-xxl",
		Title: `Edit Recipe`,
		Sections: []s.FormSection{
			{
				ID: "edit-recipe",
				Fields: []s.FormField{
					{Label: "Id", ID: "ProdLineToRecipe", Width: "hidden", Type: "text", DataType: "text", Value: strconv.Itoa(ProdLineToReciptId), Hidden: true},
					{Label: "Id", ID: "recipe-id", Width: "hidden", Type: "text", DataType: "int", Value: strconv.Itoa(recipe[0].Id), Hidden: true},
					{Label: "Part Name", ID: "CompoundName", Width: "25%", Type: "dropdownfield", DataType: "string", DropDownOptions: opt, Readonly: true},
					{Label: "Compound Code", ID: "CompoundCode", Width: "25%", Type: "text", DataType: "string", Value: recipe[0].CompoundCode},
					{Label: "Production Line", ID: "ProdLineId", Width: "25%", Type: "dropdownfield", DataType: "int", DropDownOptions: prodLineOpts},
					{Label: "Base Quantity", ID: "BaseQty", Width: "25%", Value: recipe[0].BaseQty, Type: "text", DataType: "string"},
					{SubForms: subformlist, Type: "sectional-card"},
				},
			},
		},
		Buttons: formBtn,
		JS: `
<!--html-->
<script>

function addRecipeFields(button) {
	if (!$('#EditRecipeModel').is(':visible')) return;

	// Convert current button to Delete
	button.textContent = 'Delete';
	button.classList.remove('btn-primary');
	button.classList.add('btn-danger', 'delete-btn');
	button.setAttribute('onclick', 'removeRecipeField(this)');
	button.setAttribute('id', 'delete-stg-fields');

	// Get the container and clone last row
	const container = button.closest('[id^="field-container-"]');
	if (!container) return;

	const lastRow = container.querySelector('.row.container-fluid:last-of-type');
	if (!lastRow) return;

	const newRow = lastRow.cloneNode(true);

	// Clear values
	const inputs = newRow.querySelectorAll('input, select');
	inputs.forEach(input => {
		if (input.tagName.toLowerCase() === 'select') {
			input.selectedIndex = 0;
		} else {
			input.value = '';
		}
	});

	container.appendChild(newRow);

	renumberRecipeFields(container);
}


function removeRecipeField(button) {
	if (!$('#EditRecipeModel').is(':visible')) return;
  const row = button.closest('.row.container-fluid');
	if (!row) return;

	const container = button.closest('[id^="field-container-"]');
	if (!container) return;

	row.remove();
	renumberRecipeFields(container);
}

function renumberRecipeFields(container) {
	if (!container) return;

	const rows = container.querySelectorAll('.row.container-fluid');

	rows.forEach((row, index) => {
		const indexNum = index + 1;
		const cols = row.querySelectorAll('.col.mb-0');

		cols.forEach(col => {
			const input = col.querySelector('input, select');
			const label = col.querySelector('label');

			if (!input || !label) return;

			const oldName = input.getAttribute('data-name') || '';
			const base = oldName.replace(/-\d+$/, '');
			const newId = base + '-' + indexNum;

			input.setAttribute('id', newId);
			input.setAttribute('data-name', newId);
			label.setAttribute('for', newId);
		});

		// Now handle the button in the current row
		const btn = row.querySelector('button');
		if (btn) {
			const isLastRow = (index === rows.length - 1);

			if (isLastRow) {
				// Set as Add button
				btn.textContent = 'Add';
				btn.setAttribute('id', 'add-stg-fields');
				btn.setAttribute('onclick', 'addRecipeFields(this)');
				btn.classList.remove('btn-danger', 'delete-btn');
				btn.classList.add('btn-primary');
			} else {
				// Set as Delete button
				btn.textContent = 'Delete';
				btn.setAttribute('id', 'delete-stg-fields');
				btn.setAttribute('onclick', 'removeRecipeField(this)');
				btn.classList.remove('btn-primary');
				btn.classList.add('btn-danger', 'delete-btn');
			}
		}
	});
}
</script>
<!--!html-->
		`,
	}

	return EditRecipeModel.Build()
}
