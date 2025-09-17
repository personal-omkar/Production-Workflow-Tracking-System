package login

import (
	"irpl.com/kanban-commons/utils"
)

type LoginPage struct {
	Name     string
	Username string
	Error    string
	Config   map[string]string
}

const DefaultRestHost string = "0.0.0.0" // Default port if not set in env
const DefaultRestPort string = "4300"    // Default port if not set in env
var RestHost string
var RestPort string // Global variable to hold the DB helper port

func (l *LoginPage) Build() (page string) {

	var invalidLogin string

	if len(l.Error) > 0 {
		invalidLogin = `
			<div  class="d-flex justify-content-center" style="margin-top:1rem;">
				<h5 id="notification" class="text-danger">` + l.Error + `!</h5>
			<div>
		`
	}
	loginForm := LoginForm{
		FormID:     "loginForm",
		FormAction: "",
		// Dropdownfield: []Dropdownfield{{ID: "LoginType", Name: "Login Type", Options: Dropdownoptions{Option: []string{"Customer", "Operator", "Admin"}}}},
		Inputfield: []Inputfield{
			{Type: "text", DataType: "text", Name: "username", Required: true, ID: "username", Label: "Email", Visible: true},
			{Type: "password", DataType: "text", Name: "password", Required: true, ID: "password", Label: "Password", Visible: true},
		},
		Buttons: FooterButtons{BtnType: "button", BtnID: "login", BtnSubmitGroup: "loginForm", Text: "Login"},
	}

	page = utils.JoinStr(
		`
		<!DOCTYPE html>
		<!--html-->
		<head>
		<meta charset="utf-8">
		<meta http-equiv="X-UA-Compatible" content="IE=edge">
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<title>`, utils.DefaultsMap["name"], ` - Login</title>
		<link href="/static/assets/css/theme.css" rel="stylesheet"> 	
		<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bootstrap-icons@1.4.0/font/bootstrap-icons.css"/>

		<!--JavaScripts-->
 
		<script src="/static/vendors/popper/popper.min.js"></script>
		<script src="/static/vendors/anchorjs/anchor.min.js"></script>
		<script src="/static/vendors/is/is.min.js"></script>
		<script src="/static/vendors/fontawesome/all.min.js"></script>
		<script src="/static/vendors/jquery/jquery.min.js"></script>
	 
	  </head>
	  <body>
		<main class="main" id="top">
			<div class="container-fluid">
				<div class="row min-vh-100 flex-center g-0">
					<div class="col-lg-8 col-xxl-5 py-3 position-relative">
						<div class="card overflow-hidden z-index-1">
							<div class="card-body p-0">
								<div class="row g-0 h-100">
									<div class="col-md-5 text-center bg-card-dark" style="background-color:#260020">
										<div class="position-relative p-4 pt-md-5 pb-md-7 light">
											<div class="bg-holder bg-auth-card-shape" style="background-image:url('/static/assets/img/icons/spot-illustrations/half-circle.png');">
											</div>
										<!--/.bg-holder-->
											<div class="z-index-1 position-relative"><a class="link-light mb-4 font-sans-serif fs-4 d-inline-block fw-bolder"><div class="d-flex align-items-center py-3"><img class="mr-2 fs-5" src="`+utils.DefaultsMap["login_logo_url"]+`" alt="" width="180"/></div></a>
											</div>
										</div>
										<div class="light text-center" style="margin-top:7rem">
											<h3 style="color:#ab71a2;padding-bottom:3rem;">`, utils.DefaultsMap["name"], `</h3>
										</div>
									</div>
									<div class="col-md-7 d-flex flex-center">
										<div class="p-4 p-md-5 flex-grow-1">
											<div class="row flex-between-center mb-3">
												<div class="col-auto">
													<h3 style="color:#CF7AC2">Account Login</h3>
												</div>
											</div>	
											`, loginForm.Build(), invalidLogin, `

										</div>
									</div>
								</div>
							</div>
						</div>
					</div>
				</div>
			</div>

		</main>

		<script src="/static/vendors/chart/chart.min.js"></script>
		<script src="/static/vendors/countup/countUp.umd.js"></script>
		<script src="/static/vendors/echarts/echarts.min.js"></script>
		<script src="/static/vendors/d3/d3.min.js"></script>
		<script src="/static/vendors/lodash/lodash.min.js"></script>
		
		<script src="/static/vendors/list.js/list.min.js"></script>
		<script src="/static/vendors/bootstrap/bootstrap.min.js"></script>
		
		<script>
		$(function(){
		
		const urlParams = new URLSearchParams(window.location.search);
					const status = urlParams.get('status');
					const msg = urlParams.get('msg');
					const error=urlParams.get('error');
					
					
				if (error) {
					
						removeQueryParams();
					
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
					}, 50);
		}

					$(document).on("click","#login",function(){
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

						$.post("/sign-in", JSON.stringify(result), function() {}, 'json')
							.done(function(response) {
								window.location.href = "/dashboard";
							})
							.fail(function(xhr, status, error) {
							 	const data = JSON.parse(xhr.responseText);
								window.location.href = "/login?status=" + xhr.status + "&error=" +data.message;
							});						
					}
				})
				function removeQueryParams() {		
					var newUrl = window.location.protocol + "//" + window.location.host + window.location.pathname;
					window.history.replaceState({}, document.title, newUrl);
					
				}
	})
		</script>
	  </body>
	<!--!html-->
	</html> `)
	return
}
