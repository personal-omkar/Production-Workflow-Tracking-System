package userrolemanagementpage

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"

	m "irpl.com/kanban-commons/model"
	"irpl.com/kanban-commons/utils"
	s "irpl.com/kanban-web/services"
)

type UserRoleManagement struct {
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

func (u *UserRoleManagement) Build() string {
	var Response struct {
		Pagination m.PaginationResp
		Data       []*m.UserRoles
	}
	type TableConditions struct {
		Pagination m.PaginationReq `json:"pagination"`
		Conditions []string        `json:"conditions"`
	}

	var tablecondition TableConditions

	// Temporary pagination data
	tablecondition.Pagination = m.PaginationReq{
		Type:   "userRole",
		Limit:  "10",
		PageNo: 1,
	}

	// Convert to JSON
	jsonValue, err := json.Marshal(tablecondition)
	if err != nil {
		slog.Error("Error marshaling table condition", "error", err)
	}

	resp, err := http.Post(RestURL+"/get-all-user-role-by-search-pagination", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		slog.Error("%s - error - %s", "Error making GET request", err)
	}

	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&Response); err != nil {
		slog.Error("error decoding response body", "error", err)
	}

	tableTools := `<button  type="button" class="btn  m-0 p-0" id="edit-User-Role-btn" data-bs-toggle="modal" data-bs-target="#EditUserRoleModel"> 
 					<i class="fa fa-edit " style="color: #871a83;"></i> 
					</button>
					`
	var html strings.Builder
	html.WriteString(`
		<!--html-->
		<!--!html-->`)

	var userRoleManagement s.TableCard

	userRoleManagement.CardHeading = "User Role Master"
	userRoleManagement.CardHeadingActions.Component = []s.CardHeadActionComponent{
		{ComponentName: "button", ComponentType: s.ActionComponentElement{Button: s.ButtonAttributes{ID: "redirectuserManagementPage", Name: "redirectuserManagementPage", Type: "button", Text: "User Master", Colour: "#871A83"}}, Width: "col-4"},
		{ComponentName: "modelbutton", ComponentType: s.ActionComponentElement{ModelButton: s.ModelButtonAttributes{ID: "add-new-user-role", Name: "add-new-user-role", Type: "button", Text: "Add New User Role", ModelID: "#AddUserRoleModel"}}, Width: "col-4"}}
	userRoleManagement.BodyTables = s.CardTableBody{Columns: []s.CardTableBodyHeadCol{

		{
			Name:         "Role Name",
			ID:           "search-by-role-name",
			DataField:    "role_name",
			IsSearchable: true,
			Type:         "input",
			Width:        "col-2",
		},
		{
			Name:  "Description",
			ID:    "description",
			Width: "col-2",
		},
		{
			Name:  "Created On",
			ID:    "createdOn",
			Width: "col-2",
		},
		{
			Name:  "Tools",
			Width: "col-1",
		},
	},
		ColumnsWidth: []string{"col-2", "col-2", "col-2", "col-2", "col-1"},
		Data:         Response.Data,
		Tools:        tableTools,
		ID:           "UserRoleManagement",
	}

	CheckBoxList := s.CheckBox{Label: "Allow List", Width: "w-100 mb-2", CheckBoxList: []s.CheckboxAttribut{
		{PageName: "Dashboard", PageLink: "/dashboard", IsChecked: true, Readonly: true},
		{PageName: "Kanban Report", PageLink: "/vendor-orders"},
		{PageName: "Kanban Entry", PageLink: "/order-entry"},
		{PageName: "Pending Orders", PageLink: "/admin-orders"},
		{PageName: "Heijunka Board", PageLink: "/production-line"},
		{PageName: "Kanban Board", PageLink: "/vendor-company"},
		{PageName: utils.DefaultsMap["cold_store_menu"], PageLink: "/cold-storage"},
		{PageName: "User Master", PageLink: "/user-management"},
		{PageName: "Vendor Master", PageLink: "/vendor-management"},
		{PageName: "Part Master", PageLink: "/compounds-management"},
		{PageName: "Chemical Type Master", PageLink: "/chemical-management"},
		{PageName: "Raw Material", PageLink: "/material-management"},
		{PageName: "Kanban History", PageLink: "/kanban-history"},
		{PageName: "Order History", PageLink: "/order-history"},
		{PageName: "Packing/Dispatch", PageLink: "/packing-dispatch-page"},
		{PageName: "Quality Tetsing", PageLink: "/quality-testing"},
		{PageName: "Machine Master", PageLink: "/prod-line-management"},
		{PageName: "Process Master", PageLink: "/prod-processes-management"},
		{PageName: "Kanban Report", PageLink: "/kanban-report"},
		{PageName: "Recipe Master", PageLink: "/recipe-management"},
		{PageName: "All Kanban View", PageLink: "/all-kanban-view"},
		{PageName: "Operator Master", PageLink: "/operator-management"},
		{PageName: "Kanban Reprint", PageLink: "/kanban-reprint"},
		{PageName: "Summary Reprint", PageLink: "/summary-reprint"},
		{PageName: "Rubber Store Master", PageLink: "/inventory-management"},
	}}

	var Pagination s.Pagination
	Pagination.CurrentPage = 1
	Pagination.Offset = 0
	Pagination.TotalRecords = Response.Pagination.TotalNo
	Pagination.PerPage, _ = strconv.Atoi(tablecondition.Pagination.Limit)
	Pagination.PerPageOptions = []int{10, 25, 50, 100, 200}
	userRoleManagement.CardFooter = Pagination.Build()

	// Create Dialogue
	CreateUserRoleModel := s.ModelCard{
		ID:      "AddUserRoleModel",
		Type:    "modal-lg",
		Heading: "Add User Role",
		Form: s.ModelForm{FormID: "AddUserRoleModel",
			FormAction: "", Footer: s.Footer{CancelBtn: true, Buttons: []s.FooterButtons{{BtnType: "button", DataSubmitName: "AddUserRoleModel", BtnID: "add-user-role-submit", Text: "Save"}}},
			Inputfield: []s.InputAttributes{
				{Type: "text", Name: "RoleName", DataType: "text", ID: "addRoleName", Label: `Role Name`, Width: "w-100"},
			},
			TextArea: []s.TextAreaAttributes{
				{Label: "Description", DataType: "text", Name: "Description", ID: "adddescription"},
			},
			CheckBoxList: CheckBoxList,
		},
	}

	// Edit Dialogue
	EditUserRoleModel := s.ModelCard{
		ID:      "EditUserRoleModel",
		Type:    "modal-lg",
		Heading: "Edit User Role",
		Form: s.ModelForm{FormID: "EditUserRoleModel",
			FormAction: "/update-user-role", Footer: s.Footer{CancelBtn: true, Buttons: []s.FooterButtons{{BtnType: "button", DataSubmitName: "EditUserRoleModel", BtnID: "edit-user-role-submit", Text: "Update"}}},
			Inputfield: []s.InputAttributes{
				{Type: "hidden", DataType: "int", Name: "ID", ID: "editid", Label: `User ID`, Width: "w-100", Hidden: true, Readonly: true},
				{Type: "text", DataType: "text", Name: "RoleName", ID: "editRoleName", Label: `Role Name`, Width: "w-100"},
			},
			TextArea: []s.TextAreaAttributes{
				{Label: "Description", DataType: "text", Name: "Description", ID: "editdescription"},
			},
			CheckBoxList: CheckBoxList,
		},
	}

	//Delete Dialogue
	DeleteUserRoleModel := s.ModelCard{
		ID:      "DeleteUserRoleModel",
		Type:    "modal-lg",
		Heading: "Delete User Role",
		Form: s.ModelForm{FormID: "DeleteUserRoleModel",
			FormAction: "", Footer: s.Footer{CancelBtn: true, Buttons: []s.FooterButtons{{BtnType: "button", BtnID: "delete-user-role-submit", DataSubmitName: "DeleteUserRoleModel", Text: "Delete", Style: "background-color:#dc3545; ;border:none"}}},
			Inputfield: []s.InputAttributes{
				{Type: "text", DataType: "int", Name: "ID", ID: "deleteid", Label: `User ID`, Width: "w-100", Hidden: true, Readonly: true, Disabled: true},
				{Type: "text", DataType: "text", Name: "RoleName", ID: "deleteRoleName", Label: `Role Name`, Width: "w-100", Readonly: true, Disabled: true},
			},
			TextArea: []s.TextAreaAttributes{
				{Label: "Description", DataType: "text", Name: "Description", ID: "deletedescription", Disabled: true, Readonly: true},
			},
			CheckBoxList: CheckBoxList,
		},
	}
	js := utils.JoinStr(`
	<script> 
	//js
			$(document).ready(function () {
				document.querySelectorAll("[id^='search-by']").forEach(function (input) {
					input.addEventListener("input", function () {
						pagination(1);
					});
				});
				attachPerPageListener(); 
			});
	
			function pagination(pageNo) {
				let searchCriteria = [];

				document.querySelectorAll("[id^='search-by']").forEach(function (input) {
					let field = input.dataset.field;
					let value = input.value.trim();
					if (value) {
						searchCriteria.push(field + " ILIKE '%" + value + "%'");
					}
				});

				let limit = document.getElementById("perPageSelect")?.value || "10";

				let requestData = {
					pagination: {
						Limit: limit,
						PageNo: pageNo
					},
					Conditions: searchCriteria
				};

				fetch("/user-role-search-pagination", {
					method: "POST",
					headers: {
						"Content-Type": "application/json"
					},
					body: JSON.stringify(requestData)
				})
				.then(response => response.json())
				.then(data => {
					let tableBody = document.getElementById("advanced-search-table-body");
					let cardFooter = document.querySelector(".card-footer");

					if (tableBody) {
						tableBody.innerHTML = data.tableBodyHTML;
					}
					if (cardFooter) {
						cardFooter.innerHTML = data.paginationHTML;
						attachPerPageListener(); 
					}
				})
				.catch(err => console.error("Pagination error:", err));
			}

				function attachPerPageListener() {
					let perPageSelect = document.getElementById("perPageSelect");
					if (perPageSelect) {
						perPageSelect.addEventListener("change", function () {
							pagination(1);
						});
					}
				}

				$(document).on("click", "#edit-User-Role-btn, #del-User-Role-btn", function(event) {  
					var data = JSON.parse($(this).closest("tr").attr("data-data"));
					
					//edit dialogue data
					$('#editid').val(data.ID);
					$('#editRoleName').val(data.RoleName);
					$('#editdescription').val(data.Description);
					
					//delete dialogue data
					$('#deleteid').val(data.ID);
					$('#deleteRoleName').val(data.RoleName);
					$('#deletedescription').val(data.Description);

					$('.checkbox-element').each(function () {
						const currentValue = $(this).val();
						// If the current value is not in the list, check the checkbox
						if (!data.Deny.includes(currentValue)) {
							$(this).prop('checked', true);
						}else{
							$(this).prop('checked', false); 
						}
					});

				})

				$("#EditUserRoleModel").on("hide.bs.modal", function(event) {
						$('.checkbox-element').prop('checked', false);  
						$("#Dashboard").prop("checked", true);
				});

			
				const urlParams = new URLSearchParams(window.location.search);
					const status = urlParams.get('status');
					const msg = urlParams.get('msg');
					
				if (status) {
					showNotification(status,msg, () => {
						removeQueryParams();
					});
				}
				


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

				$(document).on("click", "#redirectuserManagementPage", function(event) {  
					window.location.href = "/user-management"
				})

				$(document).on("click","#add-user-role-submit,#edit-user-role-submit,#delete-user-role-submit",function(){
					
					var group = $(this).attr("data-submit");
					var result = {}
					var validated = true;
				
					const uncheckedValues = $("[data-group='" + group + "'] .checkbox-element:not(:checked)").map(function() {
						return $(this).val();
					}).get();

					$("[data-group='" + group + "']").find("[data-name]").each(function () {
						if ($(this).attr("data-validate") && $(this).val().trim().length === 0 ||$(this).attr("data-validate") && $(this).val()==="Nil" ) {
							$(this).css("background-color", "rgba(128, 0, 128, 0.1)");
							const label = $(this).closest("label").length 
										? $(this).closest("label") 
										: $(this).siblings("label").length 
										? $(this).siblings("label") 
										: $(this).parent().siblings("label");
					
							if (label.length) {
								label.find(".required-label").remove();
								label.siblings(".required-label").remove();

								$("<span class='required-label'>Required</span>").css({
									color: "red",
									fontSize: "1em",
									"margin-left": "0.5rem",
								}).insertAfter(label);
							}
					
							validated = false;
						} else {
							$(this).css("background-color", "rgb(255, 255, 255)");
							
							const label = $(this).closest("label").length 
										? $(this).closest("label")
										: $(this).siblings("label").length
										? $(this).siblings("label")
										: $(this).parent().siblings("label");

							if (label.length) {
								label.siblings(".required-label").remove();
								label.find(".required-label").remove();
							}
						}
					});
					
					if (validated){
					
						$("[data-group='"+group+"'").find("[data-name]").each(function(){
							if ($(this).is("select")){
								result[$(this).attr("data-name")] = $(this).find(":selected").val();
							} else if ($(this).attr("data-type") == "date"){
								var userDate = $(this).val();                 
								if (userDate.includes(":")) {
									result[$(this).attr("data-name")] = userDate
								}else{
									var formattedDate = formatDateToYYYYMMDD(userDate);
									result[$(this).attr("data-name")] = new Date(formattedDate)      
								}
							} else if ($(this).attr("data-type") == "int") {
								result[$(this).attr("data-name")] = parseInt($(this).val());
							} else { 
								result[$(this).attr("data-name")] = $(this).val();
							}						
						})

						result['Deny'] = uncheckedValues; 
						
						if (group==="EditUserRoleModel"){
							$.post("/update-user-role", JSON.stringify(result), function() {}, 'json').fail(function(xhr, status, error) {
								window.location.href = "/user-role-management?status="+xhr.status+"&msg="+xhr.responseText
							}); 
						}else if (group==="AddUserRoleModel"){
							$.post("/create-user-role", JSON.stringify(result), function() {}, 'json').fail(function(xhr, status, error) {
								window.location.href = "/user-role-management?status="+xhr.status+"&msg="+xhr.responseText
							}); 
						}else if (group==="DeleteUserRoleModel"){
							let ids = [];
							ids.push(result["ID"]) ;
							$.post("/delete-user-role", JSON.stringify(ids), function() {}, 'json').fail(function(xhr, status, error) {
								window.location.href = "/user-role-management?status="+xhr.status+"&msg="+xhr.responseText
							}); 
						}
						
					}
				})

				function removeQueryParams() {
					var newUrl = window.location.protocol + "//" + window.location.host + window.location.pathname;
					window.history.replaceState({}, document.title, newUrl);
				}
		 //!js
		</script>
	`)

	html.WriteString(userRoleManagement.Build())
	html.WriteString(CreateUserRoleModel.Build())
	html.WriteString(EditUserRoleModel.Build())
	html.WriteString(DeleteUserRoleModel.Build())
	html.WriteString(js)

	html.WriteString(`</div>`)

	return html.String()

}
