package recipemanagement

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"

	m "irpl.com/kanban-commons/model"
	"irpl.com/kanban-commons/utils"
	s "irpl.com/kanban-web/services"
)

type RecipeManagement struct {
	Username     string
	UserType     string
	UserID       string
	CompoundCode string
	CompoundName string
}

const DefaultRestHost string = "0.0.0.0" // Default port if not set in env
const DefaultRestPort string = "4200"    // Default port if not set in env

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
func (r *RecipeManagement) Build() string {
	var recipe []*m.Recipe
	var stage []*m.Stage
	var compounds []*m.Compounds
	var rawMaterials []*m.RawMaterial
	var chemTypes []*m.ChemicalTypes
	var prodLine []*m.ProdLine

	var rawQuery m.RawQuery
	rawQuery.Host = utils.RestHost
	rawQuery.Port = utils.RestPort
	rawQuery.Type = "Recipe"
	rawQuery.Query = `SELECT * FROM recipe  ;` //`;`
	rawQuery.RawQry(&recipe)

	rawQuery.Type = "Stage"
	rawQuery.Query = `SELECT * FROM Stage ;` //`;`
	rawQuery.RawQry(&stage)

	rawQuery.Type = "RawMaterial"
	rawQuery.Query = `SELECT * FROM Raw_Material ;` //`;`
	rawQuery.RawQry(&rawMaterials)

	rawQuery.Type = "ChemicalTypes"
	rawQuery.Query = `SELECT * FROM Chemical_Types ;` //`;`
	rawQuery.RawQry(&chemTypes)

	rawQuery.Type = "ProdLine"
	rawQuery.Query = `SELECT * FROM prod_line ;` //`;`
	rawQuery.RawQry(&prodLine)

	//get compound data
	compResp, err := http.Get(RestURL + "/get-all-compound-data")
	if err != nil {
		slog.Error("%s - error - %s", "Error making GET request", err)
	}

	defer compResp.Body.Close()

	if err := json.NewDecoder(compResp.Body).Decode(&compounds); err != nil {
		slog.Error("error decoding response body", "error", err)
	}

	tableTools := `<button  type="button" class="btn  m-0 p-0" id="edit-recipe-btn"> 
 					<i class="fa fa-edit " style="color: #871a83;"></i> 
					</button>
					`
	var html strings.Builder
	html.WriteString(`
		<!--html-->
		<!--!html-->`)

	var recipeManagement s.TableCard

	recipeManagement.CardHeading = "Recipe Master"
	recipeManagement.CardHeadingActions.Component = []s.CardHeadActionComponent{{ComponentName: "button", ComponentType: s.ActionComponentElement{Button: s.ButtonAttributes{ID: "redirectStagePage", Name: "redirectStagePage", Type: "button", Text: "Stages", Colour: "#871A83"}}, Width: "col-4"}, {ComponentName: "modelbutton", ComponentType: s.ActionComponentElement{ModelButton: s.ModelButtonAttributes{ID: "add-new-recipe", Name: "add-new-recipe", Type: "button", Text: "Add New Recipe", ModelID: "#AddRecipeModel"}}, Width: "col-4"}}
	recipeManagement.BodyTables = s.CardTableBody{Columns: []s.CardTableBodyHeadCol{

		{
			Lable: "Part Name",
			Name:  "Compound Name",
			ID:    "Compound Name",
			Width: "col-4",
		},
		{
			Name:  "Compound Code",
			ID:    "Compound Code",
			Width: "col-4",
		},
		{
			Name:  "Created On",
			ID:    "lot_no",
			Width: "col-2",
		},
		{
			Name:  "Tools",
			Width: "col-2",
		},
	},
		ColumnsWidth: []string{"col-4", "col-4", "col-2", "col-2"},
		Data:         recipe,
		Tools:        tableTools,
		ID:           "RecipeManagement",
	}

	//Add Dialogue
	var opt []s.DropDownOptions
	for _, v := range compounds {
		var tempopt s.DropDownOptions
		tempopt.Text = v.CompoundName
		tempopt.Value = v.CompoundName
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

	formBtn := utils.JoinStr(`
			<div class="-flex justify-content-end">
				<button type="button" class="btn " data-bs-dismiss="modal" style="background-color:#871a83 ;border:none;color:white">Cancel</button>
				<button id="add-recipe" data-submit="addrecipemodel" data-url="/create-new-recipe" class="btn ms-3" style="background-color:#871a83 ;border:none;color:white">Save</button>
			</div>
		`)

	var subformlist []s.FormSection
	for _, v := range stage {
		if v.Active {
			var fields []s.FormField
			var temp s.FormSection
			temp.ID = v.Name
			temp.SubFormName = v.Name
			temp.Data = `data-id="` + strconv.FormatUint(uint64(v.ID), 10) + `"`

			// Parse headers
			var headers []map[string]string
			if err := json.Unmarshal(v.Headers, &headers); err != nil {
				log.Printf("error decoding headers: %v", err)
			}

			// Add fields from headers
			for _, header := range headers {

				if header["type"] == "dropdown-RM" {

					fields = append(fields,
						s.FormField{Label: header["field"], ID: header["field"] + "-1", Width: "auto", Type: "select", DropDownOptions: rawOpts},
					)
				} else if header["type"] == "dropdown-CT" {

					fields = append(fields,
						s.FormField{Label: header["field"], ID: header["field"] + "-1", Width: "auto", Type: "select", DropDownOptions: chemTypeOpts},
					)
				} else {
					fields = append(fields,
						s.FormField{Label: header["field"], ID: header["field"] + "-1", Width: "auto", Type: "text", DataType: "string", Placeholder: "Enter " + header["field"]},
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

	}
	addRecipeModel := s.Model{
		ID:    "AddRecipeModel",
		CSS:   `<style> .modal-xxl { width: 1355px; max-width: 1355px} </style>`,
		Type:  "modal-xxl",
		Title: `Add Recipe`,
		Sections: []s.FormSection{
			{
				ID: "add-new-recipe",
				Fields: []s.FormField{
					{Label: "Part Name", ID: "CompoundName", Width: "25%", Type: "dropdownfield", DataType: "string", DropDownOptions: opt},
					{Label: "Compound Code", ID: "CompoundCode", Width: "25%", Type: "text", DataType: "string"},
					{Label: "Production Line", ID: "ProdLineId", Width: "25%", Type: "dropdownfield", DataType: "string", DropDownOptions: prodLineOpts},
					{Label: "Base Quantity", ID: "BaseQty", Width: "25%", Type: "text", DataType: "string"},
					{SubForms: subformlist, Type: "sectional-card"},
				},
			},
		},
		Buttons: formBtn,

		JS: `
<!--html-->
<script>
function addRecipeFields(button) {
	if (!$('#AddRecipeModel').is(':visible')) return;

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
	if (!$('#AddRecipeModel').is(':visible')) return;
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
`,
	}

	js :=
		`<script> 
	
			
				$(document).ready(function() {
					var selectedValue = $("#addRoleName").val(); 
					if (selectedValue==2){
						$("#addVendorCode").closest("div.w-100").removeAttr("style");
					}else{
						$("#addVendorCode").closest("div.w-100").attr("style", "display: none;");
					}

					
				});
				$(document).on("click", "#edit-User-btn, #del-User-btn", function(event) {  
					var data = JSON.parse($(this).closest("tr").attr("data-data"));
					if (data.RoleName==='Customer'){
						$("#showVendorCode").closest("div.w-100").removeAttr("style");
					}else{
						$("#showVendorCode").closest("div.w-100").attr("style", "display: none;");
					}
					// Get modal content for the clicked user
					$.get("/user-management-card?email=" + data.Email, function(response) {
						let modal = $("#EditUserModel");
						
						if (modal.length) {
							$("#EditUserModel").replaceWith(response);
						} else {
							$(".main-container").append(response);
						}

						// Show the modal after updating/appending
						$("#EditUserModel").modal("show");
					}, 'json');

					
				})

				$(document).on("click", "#redirectStagePage", function(event) {  
					window.location.href = "/stage-management"
				})

				$(document).on('click', '#edit-recipe-btn', function () {
					var data = JSON.parse($(this).closest("tr").attr("data-data"));
					$.post("/edit-recipe-dialog", JSON.stringify(data), function(response) {
						$('#dialog').html(response);    
						const modal = new bootstrap.Modal(document.getElementById('EditRecipeModel'));
						modal.show();
					})  
				});
				const urlParams = new URLSearchParams(window.location.search);
					const status = urlParams.get('status');
					const msg = urlParams.get('msg');
					
				if (status) {
					showNotification(status,msg, () => {
						removeQueryParams();
					});
				}
					$("#addRoleName ,#showRoleName").change(function() {
						var selectedValue = $(this).val();
						if (selectedValue==2){
							$("#addVendorCode ,#showVendorCode").closest("div.w-100").removeAttr("style");
						}else{
							$("#addVendorCode,#showVendorCode").closest("div.w-100").attr("style", "display: none;");
						}
					});

  			
				$(document).on("click", "#add-recipe, #update-recipe", function () {
					var group = $(this).attr("data-submit");
					var url = $(this).attr("data-url");
					var result = {};
					var validated = true;

					// Validate only fields with data-validate
					$("[data-group='" + group + "']:visible").find("[data-validate]").each(function () {
						var value = $(this).val() ? $(this).val().trim() : "";
						var parent = $(this).parent();
						var label = parent.find("label");

						if (value === "" || value === "Nil") {
							$(this).css("background-color", "rgba(128, 0, 128, 0.1)");
							parent.find(".required-label").remove();

							var labelRow = label.closest(".d-flex");
							if (!labelRow.length) {
								label.wrap('<div class="d-flex justify-content-between align-items-center "></div>');
								labelRow = label.closest(".d-flex");
							}

							$("<span class='required-label text-danger'>Required</span>").css({
								fontSize: "1em"
							}).appendTo(labelRow);

							validated = false;
						} else {
							$(this).css("background-color", "rgb(255, 255, 255)");
							parent.find(".required-label").remove();
						}
					});

					// Continue only if 
					
						if (validated) {
							var dataList = [];
							$("[data-group='" + group + "']:visible").find("[subform-name]").each(function () {
							const name = $(this).attr("subform-name");
							const dataId = $(this).attr("data-id");
							
								
								if (name) {
									var resultArray = [];

									$(this).find(".row.container-fluid").each(function () {
										var result = {};
										var hasValue = false; // Flag to detect at least one non-empty input

										$(this).find("[data-name]").each(function () {
											var key = $(this).attr("data-name");
											var type = $(this).attr("data-type");
											var value = $(this).val()?.trim();

											if ($(this).is("select")) {
												var selectedValue = $(this).find(":selected").val();
												if (selectedValue !== "" && selectedValue !== "Nil") {
													hasValue = true;
												}
												if (type === "int") {
													result[key] = parseInt(selectedValue, 10);
												} else if (type === "bool") {
													result[key] = selectedValue === "true";
												} else {
													result[key] = selectedValue;
												}
											} else if (type === "date") {
												if (value !== "") {
													hasValue = true;
													if (value.indexOf(":") !== -1) {
														result[key] = value;
													} else {
														var formattedDate = formatDateToYYYYMMDD(value);
														result[key] = new Date(formattedDate);
													}
												}
											} else if (type === "int") {
												if (value !== "" && !isNaN(value)) {
													hasValue = true;
													result[key] = parseInt(value, 10);
												}
											} else {
												if (value !== "") {
													hasValue = true;
												}
												result[key] = value;
											}
										});

										if (hasValue) {
											resultArray.push(result);
										}
									});

									var data = {
										ID:parseInt($("[data-group='" + group + "']:visible").find('#recipe-id').val()),
										ProdLineToRecipe:parseInt($("[data-group='" + group + "']:visible").find('#ProdLineToRecipe').val()),
										ProdLineId:parseInt($("[data-group='" + group + "']:visible").find('#ProdLineId').val()),
										CompoundName:  $("[data-group='" + group + "']:visible").find('#CompoundName').val(),
										CompoundCode: $("[data-group='" + group + "']:visible").find('#CompoundCode').val(),
										BaseQty:$("[data-group='" + group + "']:visible").find('#BaseQty').val(),
										Data:resultArray,
										StageId: parseInt(dataId)
									};
									dataList.push(data);
								}
							});
						}
						$.post(url, JSON.stringify(dataList), function() {}, 'json').fail(function(xhr, status, error) {
							 window.location.href = "/recipe-management?status="+xhr.status+"&msg="+xhr.responseText
						}); 	
					
				});

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
					notification.html(message);
					notification.show();

					setTimeout(() => {
						notification.fadeOut(() => {
							if (callback) callback();
						});
					}, 5000);
				}

				function removeQueryParams() {
					var newUrl = window.location.protocol + "//" + window.location.host + window.location.pathname;
					window.history.replaceState({}, document.title, newUrl);
				}
	
		</script>
		<!--!html-->`

	html.WriteString(recipeManagement.Build())
	html.WriteString(addRecipeModel.Build())
	html.WriteString(js)

	html.WriteString(`</div> </div></div></div></div></div></div></div>`)

	return html.String()

}
