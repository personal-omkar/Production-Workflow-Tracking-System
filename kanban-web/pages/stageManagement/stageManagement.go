package stagemanagement

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"strings"

	m "irpl.com/kanban-commons/model"
	"irpl.com/kanban-commons/utils"

	s "irpl.com/kanban-web/services"
)

type StageManagement struct {
	Username string
	UserType string
	UserID   string
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
func (sm *StageManagement) Build() string {
	var stages []*m.Stage

	stagesResp, err := http.Get(RestURL + "/get-all-stages")
	if err != nil {
		slog.Error("%s - error - %s", "Error making GET request", err)
	}

	defer stagesResp.Body.Close()

	if err := json.NewDecoder(stagesResp.Body).Decode(&stages); err != nil {
		slog.Error("error decoding response body", "error", err)
	}

	tableTools := `<button  type="button" class="btn  m-0 p-0 edit-stage-btn" id="edit-stage-btn"> 
 					<i class="fa fa-edit " style="color: #871a83;"></i> 
					</button>
					`
	var html strings.Builder

	var stageManagement s.TableCard

	stageManagement.CardHeading = "Stage Master"
	stageManagement.CardHeadingActions.Component = []s.CardHeadActionComponent{{ComponentName: "button", ComponentType: s.ActionComponentElement{Button: s.ButtonAttributes{ID: "redirectRecipePage", Name: "redirectRecipePage", Type: "button", Text: "Recipe", Colour: "#871A83"}}, Width: "col-4"}, {ComponentName: "modelbutton", ComponentType: s.ActionComponentElement{ModelButton: s.ModelButtonAttributes{ID: "add-new-stage-dialog", Name: "add-new-stage-dialog", Type: "button", Text: "Add New Stage", ModelID: "#AddStageModel"}}, Width: "col-4"}}
	stageManagement.BodyTables = s.CardTableBody{Columns: []s.CardTableBodyHeadCol{

		{
			Name:  "Name",
			ID:    "Name",
			Width: "col-4",
		},
		{
			Name:  "Created On",
			ID:    "CreatedOn",
			Width: "col-4",
		},
		{
			Name:  "Active",
			ID:    "Active",
			Width: "col-2",
		},
		{
			Name:  "Tools",
			Width: "col-2",
		},
	},
		ColumnsWidth: []string{"col-4", "col-4", "col-2", "col-2"},
		Data:         stages,
		Tools:        tableTools,
		ID:           "StageManagement",
	}

	formBtn := utils.JoinStr(`
			<div class="col-md-6 d-flex justify-content-end">
				<button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Cancel</button>
				<button id="add-new-stage" data-submit="addstagemodel" data-url="/create-new-stage" class="btn btn-primary ms-3" style="background-color:#871A83 ;border:none">Save</button>
			</div>
		`)

	var fieldTypeOptions = []s.DropDownOptions{
		{Value: "input", Text: "Input", Selected: true},
		{Value: "dropdown-RM", Text: "RM Dropdown"},
		{Value: "dropdown-CT", Text: "CT Dropdown"},
	}

	//Add Dialogue
	addStageModel := s.Model{
		ID:    "AddStageModel",
		Type:  "modal-lg",
		Title: `<h5 class="modal-title text-primary fs-3" id="staticBackdropLabel"><b style="color:#871A83">Add Stage</b></h5>`,
		Sections: []s.FormSection{
			{
				ID: "add-new-stage-section",
				Fields: []s.FormField{
					{Label: "Name", ID: "Name", Width: "100%", Type: "text", DataType: "string", Placeholder: "Enter name", IsRequired: true},

					{Label: "Field Name", ID: "FieldName-1", Width: "33%", Type: "text", DataType: "string", Placeholder: "Enter field name"},

					{Label: "Field Type", ID: "FieldType-1", Width: "33%", Type: "select", DropDownOptions: fieldTypeOptions},

					{Label: "Add", ID: "add-header", Width: "25%", Type: "button", AdditionalAttr: "class='btn btn-primary w-100 mt-4' onclick='addField(this)'"},
				},
			},
		},
		Buttons: formBtn,
		JS: `
<!--html-->
<script>
let fieldCount = 1;

function addField(button) {
    if (!$('#AddStageModel').is(':visible')) return;

    // Change current Add button to Delete
    button.textContent = 'Delete';
    button.classList.remove('btn-primary');
    button.classList.add('btn-danger', 'delete-btn');
    button.id = 'delete-header';
    button.setAttribute('onclick', 'removeField(this)');

    fieldCount++;

    const container = document.querySelector('#AddStageModel #field-container-add-new-stage-section');

    const row = document.createElement('div');
    row.className = 'row container-fluid m-0 p-0 mt-2';

    row.innerHTML =
        '<div class="col-4 mb-0">' +
            '<div class="col mb-0">' +
                '<label class="container-fluid p-0 my-1" for="FieldName-' + fieldCount + '">Field Name</label>' +
            '</div>' +
            '<input id="FieldName-' + fieldCount + '" data-name="FieldName-' + fieldCount + '" type="text" data-type="string" class="form-control p-2" placeholder="Enter field name" value="">' +
        '</div>' +
        '<div class="col-4 mb-0">' +
            '<div class="col mb-0">' +
                '<label class="container-fluid p-0 my-1" for="FieldType-' + fieldCount + '">Field Type</label>' +
            '</div>' +
            '<select class="form-select" id="FieldType-' + fieldCount + '" data-name="FieldType-' + fieldCount + '" data-type="">' +
                '<option value="input" selected>Input</option>' +
                '<option value="dropdown-RM">RM Dropdown</option>' +
				'<option value="dropdown-CT">CT Dropdown</option>' +
            '</select>' +
        '</div>' +
        '<div class="col-3 mb-0">' +
            '<button id="add-header" data-name="add-header" name="add-header" class="btn btn-primary w-100 mt-4" onclick="addField(this)">Add</button>' +
        '</div>';

    container.appendChild(row);
}

function removeField(button) {
    if (!$('#AddStageModel').is(':visible')) return;
    const row = button.closest('.row');
    if (row) row.remove();
    renumberFields();
}

function renumberFields() {
    const container = document.querySelector('#AddStageModel #field-container-add-new-stage-section');
    const rows = container.querySelectorAll('.row');
    fieldCount = rows.length;

    rows.forEach((row, index) => {
        const i = index + 1;

        const input = row.querySelector('input');
        const inputLabel = row.querySelector('label[for^="FieldName"]');
        if (input) {
            input.id = 'FieldName-' + i;
            input.setAttribute('data-name', 'FieldName-' + i);
        }
        if (inputLabel) {
            inputLabel.setAttribute('for', 'FieldName-' + i);
        }

        const select = row.querySelector('select');
        const selectLabel = row.querySelector('label[for^="FieldType"]');
        if (select) {
            select.id = 'FieldType-' + i;
            select.setAttribute('data-name', 'FieldType-' + i);
        }
        if (selectLabel) {
            selectLabel.setAttribute('for', 'FieldType-' + i);
        }

        const button = row.querySelector('button');
        if (i === rows.length) {
            // Last row should be "Add"
            button.textContent = 'Add';
            button.className = 'btn btn-primary w-100 mt-4';
            button.id = 'add-header';
            button.setAttribute('onclick', 'addField(this)');
        } else {
            // All others should be "Delete"
            button.textContent = 'Delete';
            button.className = 'btn btn-danger delete-btn w-100 mt-4';
            button.id = 'delete-header';
            button.setAttribute('onclick', 'removeField(this)');
        }
    });
}

</script>
`,
	}
	js :=
		`
		<script> 

				$(document).on("click", "#redirectRecipePage", function(event) {  
					window.location.href = "/recipe-management"
				})
			
				$(document).on("shown.bs.modal", "#AddStageModel , #EditStageModel", function(){
					var thisModel = $(this)

					$(thisModel).on("input", "#HeaderName", function(){
						if ($(this).val() !== "") {
							thisModel.find("#add-header").prop('disabled', false);
						} else {
							thisModel.find("#add-header").prop('disabled', true);
						}
					
					})
					
					$(document).on("click", "#add-header", function(){
						
						var headerName = thisModel.find('#HeaderName').val()

						var toastHtml = '<div id="toast-' + headerName + '" class="toast show text-white bg-primary border-0" ' +
							'style="width: 230px; height: 28px;" role="alert" data-bs-autohide="false">' +
							'<div class="d-flex align-items-center justify-content-between" style="height: 28px; overflow: hidden;">'+
							'<p class="mb-0 text-truncate" style="padding-left: 10px; max-width: calc(100% - 32px);" data-name="' + headerName + '" data-value="' + headerName + '">' +
							headerName + '</p>' +
							'<button type="button" class="btn-close btn-close-white me-2 m-auto" data-bs-dismiss="toast" aria-label="Close"></button>' +
							'</div></div>';

							
						thisModel.find('#header-list-toast').append(toastHtml);
						thisModel.find("#add-header").prop('disabled', true);
						thisModel.find("#HeaderName").val('');
					})
				});
				
				$('#AddStageModel').on('hidden.bs.modal', function () {
					$(this).find('input, textarea').val('');
					$("#header-list-toast").html('')
				});

				const urlParams = new URLSearchParams(window.location.search);
					const status = urlParams.get('status');
					const msg = urlParams.get('msg');
					
				if (status) {
					showNotification(status,msg, () => {
						removeQueryParams();
					});
				}

				$(document).on("click", "#redirectRecipePage", function(event) {  
					window.location.href = "/recipe-management"
				})

				$(document).on("click", ".edit-stage-btn", function(event) {  
					var data = JSON.parse($(this).closest("tr").attr("data-data"));
					
					// Get modal content for the clicked user
					$.get("/edit-stage-dialog?id=" + data.ID, function(response) {
						let modal = $("#EditStageModel");
						if (modal.length) {
							$("#EditStageModel").replaceWith(response);
						} else {
							$(".main-container").append(response);
						}
						// Show the modal after updating/appending
						$("#EditStageModel").modal("show");
					});

					
				})

				$(document).on("click", "#add-new-stage, #edit-stage", function () {
					var group = $(this).attr("data-submit");
					var url = $(this).attr("data-url");
					var result = {};
					var validated = true;
					var headerlist = [];

					// Determine container based on clicked button
    				var containerId = $(this).attr("id") === "add-new-stage" 
                      ? "#field-container-add-new-stage-section" 
                      : "#field-container-edit-new-stage-section";

					$(containerId).find('input[type="text"]').each(function () {
						var dataName = $(this).data('name');
						
						// Only process inputs with data-name containing 'FieldName'
						if (!dataName || !dataName.includes('FieldName')) return;

						var fieldName = $(this).val().trim();

						// Get the corresponding field type (select element) only if it has valid data-name
						var fieldTypeSelect = $(this).closest('.col-4').next('.col-4').find('select');
						var fieldTypeDataName = fieldTypeSelect.data('name');
						
						if (!fieldTypeDataName || !fieldTypeDataName.includes('FieldType')) return;

						var fieldType = fieldTypeSelect.val();

						if (fieldName && fieldType) {
							headerlist.push({
								field: fieldName,
								type: fieldType
							});
						}
					});
					
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
								label.wrap('<div class="d-flex justify-content-between align-items-center container-fluid p-0 my-1"></div>');
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

					// Continue only if validated
					if (validated) {
						$("[data-group='" + group + "']:visible").find("[data-name]").each(function () {
							var key = $(this).attr("data-name");
							var type = $(this).attr("data-type");

							if ($(this).is("select")) {
								var selectedValue = $(this).find(":selected").val();
								if (type === "int") {
									result[key] = parseInt(selectedValue, 10);
								} else if (type === "bool") {
									result[key] = selectedValue === "true";
								} else {
									result[key] = selectedValue;
								}
							} else if (type === "date") {
								var userDate = $(this).val();
								if (userDate.indexOf(":") !== -1) {
									result[key] = userDate;
								} else {
									var formattedDate = formatDateToYYYYMMDD(userDate);
									result[key] = new Date(formattedDate);
								}
							} else if (type === "int") {
								result[key] = parseInt($(this).val(), 10);
							} else {
								result[key] = $(this).val();
							}
						});

						result["Headers"] = headerlist

						$.post(url, JSON.stringify(result), function() {}, 'json').fail(function(xhr, status, error) {
							window.location.href = "/stage-management?status="+xhr.status+"&msg="+xhr.responseText
						}); 	
					}
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
	html.WriteString(stageManagement.Build())
	html.WriteString(addStageModel.Build())
	html.WriteString(js)

	html.WriteString(`</div>`)

	return html.String()

}
