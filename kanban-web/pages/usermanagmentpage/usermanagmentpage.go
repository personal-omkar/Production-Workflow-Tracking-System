package usermanagmentpage

import (
	"bytes"
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

type UserManagement struct {
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
func (u *UserManagement) Build() string {
	var user []*m.UserManagement
	var userRoleslist []s.DropDownOptions
	var userRoles []m.UserRoles
	var userstatus = []s.DropDownOptions{{Text: "False", Value: "false"}, {Text: "True", Value: "true"}}
	//fetching user records
	// Inside func (u *UserManagement) Build() string {
	type TableConditions struct {
		Pagination m.PaginationReq `json:"pagination"`
		Conditions []string        `json:"conditions"`
	}

	var Response struct {
		Pagination m.PaginationResp    `json:"pagination"`
		Data       []*m.UserManagement `json:"data"`
	}

	tablecondition := TableConditions{
		Pagination: m.PaginationReq{
			Type:   "user",
			Limit:  "10",
			PageNo: 1,
		},
	}

	// Marshal and call paginated API
	jsonValue, err := json.Marshal(tablecondition)
	if err != nil {
		slog.Error("Error marshaling table condition", "error", err)
	}

	resp, err := http.Post(RestURL+"/get-all-user-by-search-pagination", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		slog.Error("Error making POST request", "error", err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&Response); err != nil {
		slog.Error("Error decoding response body", "error", err)
	}
	user = Response.Data

	log.Println("Fetched users", "count", len(Response.Data))

	//fetching compounds records
	userroleresp, err := http.Get(RestURL + "/get-all-user-roles")
	if err != nil {
		slog.Error("%s - error - %s", "Error making GET request", err)
	}

	defer userroleresp.Body.Close()

	if err := json.NewDecoder(userroleresp.Body).Decode(&userRoles); err != nil {
		slog.Error("error decoding response body", "error", err)
	}
	for _, i := range userRoles {
		var userroles s.DropDownOptions
		userroles.Text = i.RoleName
		userroles.Value = strconv.Itoa(i.ID)

		userRoleslist = append(userRoleslist, userroles)
	}
	tableTools := `<button  type="button" class="btn  m-0 p-0" id="edit-User-btn"> 
 					<i class="fa fa-edit " style="color: #871a83;"></i> 
					</button>
					`
	var html strings.Builder
	html.WriteString(`
		<!--html-->
		<!--!html-->`)

	var userManagement s.TableCard

	userManagement.CardHeading = "User Master"
	userManagement.CardHeadingActions.Component = []s.CardHeadActionComponent{{ComponentName: "button", ComponentType: s.ActionComponentElement{Button: s.ButtonAttributes{ID: "redirectuserRolePage", Name: "redirectuserRolePage", Type: "button", Text: "User Roles", Colour: "#871A83"}}, Width: "col-4"}, {ComponentName: "modelbutton", ComponentType: s.ActionComponentElement{ModelButton: s.ModelButtonAttributes{ID: "add-new-user-role", Name: "add-new-user-role", Type: "button", Text: "Add New User", ModelID: "#AddUserModel"}}, Width: "col-4"}}
	userManagement.BodyTables = s.CardTableBody{Columns: []s.CardTableBodyHeadCol{

		{
			Name:         "User Name",
			ID:           "search-by-username",
			DataField:    "username",
			IsSearchable: true,
			Type:         "input",
			Width:        "col-2",
		},

		{
			Name:  "Email",
			ID:    "Email",
			Width: "col-2",
		},
		{
			Name:  "Role Name",
			ID:    "Role Name",
			Width: "col-2",
		},
		{
			Name:  "Created On",
			ID:    "lot_no",
			Width: "col-2",
		},
		{
			Name:  "Tools",
			Width: "col-1",
		},
	},
		ColumnsWidth: []string{"col-2", "col-2", "col-2", "col-2", "col-1"},
		Data:         user,
		Tools:        tableTools,
		ID:           "UserManagement",
	}
	//var Pagination s.Pagination
	var Pagination s.Pagination
	Pagination.CurrentPage = 1
	Pagination.Offset = 0
	Pagination.TotalRecords = Response.Pagination.TotalNo
	Pagination.PerPage, _ = strconv.Atoi(tablecondition.Pagination.Limit)
	Pagination.PerPageOptions = []int{10, 25, 50, 100, 200}
	userManagement.CardFooter = Pagination.Build()

	//Add Dialogue
	AddUserModel := s.ModelCard{
		ID:      "AddUserModel",
		Type:    "modal-md",
		Heading: "Add New User",
		Form: s.ModelForm{FormID: "AddUserModel",
			FormAction: "", Footer: s.Footer{CancelBtn: true, Buttons: []s.FooterButtons{{BtnType: "button", DataSubmitName: "AddUserModel", BtnID: "add-user-submit", Text: "Save"}}},
			Inputfield: []s.InputAttributes{
				{Type: "text", DataType: "text", Name: "username", ID: "addUserName", Label: `User Name`, Width: "w-100", Required: true},
				{Type: "text", DataType: "text", Name: "email", ID: "addEmail", Label: `Email`, Width: "w-100", Required: true},
				{Type: "password", DataType: "text", Name: "password", ID: "addPassword", Label: `Password`, Width: "w-100", Required: true},
				{Type: "text", DataType: "text", Name: "vendors", ID: "addVendorCode", Label: `Vendor Code`, Width: "w-100", Hidden: true},
			},
			Dropdownfield: []s.DropdownAttributes{
				{Label: "Role Name", DataType: "int", Name: "roleId", ID: "addRoleName", Options: userRoleslist, Width: "w-100"},
				{Label: "Is Active", DataType: "text", Name: "isactive", Options: userstatus, ID: "addIsActive", Width: "w-100"},
			},
		},
	}

	//Delete Dialogue
	DeleteUserModel := s.ModelCard{
		ID:      "DeleteUserModel",
		Type:    "modal-md",
		Heading: "Delete User",
		Form: s.ModelForm{FormID: "DeleteUserModel",
			FormAction: "", Footer: s.Footer{CancelBtn: true, Buttons: []s.FooterButtons{{BtnType: "button", DataSubmitName: "DeleteUserModel", BtnID: "delete-user-submit", Text: "Delete", Style: "background-color:#dc3545; ;border:none"}}},
			Inputfield: []s.InputAttributes{
				{Type: "hidden", Name: "id", ID: "deleteUserId", Label: `User ID`, Width: "w-100", Hidden: true},
				{Type: "text", Name: "username", ID: "deleteUserName", Label: `User Name`, Width: "w-100", Readonly: true, Disabled: true},
				{Type: "text", Name: "email", ID: "deleteEmail", Label: `Email`, Width: "w-100", Readonly: true, Disabled: true},
				{Type: "password", Name: "password", ID: "deletePassword", Label: `Password`, Width: "w-100", Readonly: true, Disabled: true},
			},
			Dropdownfield: []s.DropdownAttributes{
				{Label: "Role Name", Name: "roleId", ID: "deleteRoleName", Options: userRoleslist, Disabled: true, Width: "w-100"},
				{Label: "Is Active", DataType: "text", Name: "isactive", Options: userstatus, Disabled: true, ID: "deleteIsActive", Width: "w-100"},
			},
		},
	}
	js :=
		`<script> 
		//js
				$(document).ready(function () {
		// Attach event listeners to all searchable inputs
		document.querySelectorAll("[id^='search-by']").forEach(function (input) {
			input.addEventListener("input", function () {
				pagination(1);
			});
		});

		attachPerPageListener(); // On page load
	});

	// Pagination and search fetch handler
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
				Pageno: pageNo
			},
			Conditions: searchCriteria
		};

		fetch("/user-search-pagination", {
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
				attachPerPageListener(); // rebind dropdown listener
			}
		})
		.catch(err => console.error("Pagination error:", err));
	}

	// Re-bind per-page select on page update
	function attachPerPageListener() {
		let perPageSelect = document.getElementById("perPageSelect");
		if (perPageSelect) {
			perPageSelect.addEventListener("change", function () {
				pagination(1);
			});
		}
	}

				$(document).ready(function() {
					var selectedValue = $("#addRoleName option:selected").text(); 
					if (selectedValue=="Customer" || selectedValue=="Operator"){
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
						$("#additional-content").html(response);
						// Show the modal after updating/appending
						$("#EditUserModel").modal("show");
					}, 'json');

					
				})

				$(document).on("click", "#redirectuserRolePage", function(event) {  
					window.location.href = "/user-role-management"
				})
					
				const urlParams = new URLSearchParams(window.location.search);
					const status = urlParams.get('status');
					const msg = urlParams.get('msg');
					
				if (status) {
					showNotification(status,msg, () => {
						removeQueryParams();
					});
				}
					$(document).on("change", "#addRoleName", function() {
						var selectedValue = $("#addRoleName option:selected").text(); 
						if (selectedValue=="Customer" || selectedValue=="Operator"){
							$("#addVendorCode").closest("div.w-100").removeAttr("style");
						}else{
							$("#addVendorCode").closest("div.w-100").attr("style", "display: none;");
						}
					});

					$(document).on("change", "#showRoleName", function() {
						var selectedValue = $("#showRoleName option:selected").text(); 
						if (selectedValue=="Customer" || selectedValue=="Operator"){
							$("#showVendorCode").closest("div.w-100").removeAttr("style");
						}else{
							$("#showVendorCode").closest("div.w-100").attr("style", "display: none;");
						}
					});

  			
				$(document).on("click","#add-user-submit,#edit-user-submit,#delete-user-submit",function(){
					
					var group = $(this).attr("data-submit");
					var result = {}
					var validated = true;
				
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
						result['roleId']=parseInt(result['roleId'], 10)
						if (result['isactive']==="true"){
							 result['isactive']=true
						}else{
						 	result['isactive']=false
						}	
							var data 
							if (group === "EditUserModel") {
							
								$.post("/update-user-details", JSON.stringify(result), function() {}, 'json')
									.done(function(response, textStatus, jqXHR) {
										 data = JSON.parse(jqXHR.responseText);
										window.location.href = "/user-management?status=" +data.code+"&msg="  +data.message;
									})
									.fail(function(xhr, status, error) {
										 data = JSON.parse(xhr.responseText);
										window.location.href = "/user-management?status=" +data.code+"&msg="  +data.message;
									});
							} else if (group === "AddUserModel") {
								$.post("/create-new-user", JSON.stringify(result), function() {}, 'json')
									.done(function(response, textStatus, jqXHR) {	
									 	data = JSON.parse(jqXHR.responseText);
										 window.location.href = "/user-management?status=" +data.code+"&msg="  +data.message;
									})
									.fail(function(xhr, status, error) {
										data = JSON.parse(xhr.responseText);
										window.location.href ="/user-management?status=" +data.code+"&msg="  +data.message;
									});
							} else if (group === "DeleteUserModel") {
								let ids = [];
								ids.push(result["id"]);
								$.post("/delete-user", JSON.stringify(ids), function() {}, 'json')
									.done(function(response) {
										window.location.href = "/user-management?status=200&msg=" + encodeURIComponent("User deleted successfully");
									})
									.fail(function(xhr, status, error) {
										window.location.href = "/user-management?status=" + xhr.status + "&msg=" + encodeURIComponent(xhr.responseText);
									});
							}

								
						
					}
				})

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
		//!js
		</script>`
	html.WriteString(userManagement.Build())
	html.WriteString(AddUserModel.Build())
	// html.WriteString(EditUserModel.Build())
	html.WriteString(DeleteUserModel.Build())
	html.WriteString(js)

	html.WriteString(`</div>`)

	return html.String()

}
