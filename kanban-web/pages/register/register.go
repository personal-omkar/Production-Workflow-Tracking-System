package register

import (
	u "irpl.com/kanban-commons/utils"
)

type RegisterPage struct {
	Name     string
	Username string
	Config   map[string]string
	Error    string
}

const DefaultRestHost string = "0.0.0.0" // Default port if not set in env
const DefaultRestPort string = "4300"    // Default port if not set in env
var RestHost string
var RestPort string // Global variable to hold the DB helper port

func (d *RegisterPage) Build() (page string) {

	var invalidRegistration string

	if len(d.Error) > 0 {
		invalidRegistration = `
			<div class="d-flex justify-content-center" style="margin-top:1rem;">
				<h5 id="notification" class="text-danger">` + d.Error + `!</h5>
			<div>
		`
	}

	page = u.JoinStr(`
	<html>
	<!--html-->
		<head>
		<meta charset="utf-8">
		<meta http-equiv="X-UA-Compatible" content="IE=edge">
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<title>Registration</title>
		<link href="/static/assets/css/theme.css" rel="stylesheet"> 	
		<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bootstrap-icons@1.4.0/font/bootstrap-icons.css" />

		<!--    JavaScripts-->
 
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
				<form action="" method="post" data-group="registerForm">
             		 <div class="card-body p-0 flex-center">
						<div class="z-index-1 position-relative "><div class="d-flex align-items-center py-3 flex-center"><img class="mr-2 fs-5" src="`+u.DefaultsMap["registration_logo_url"]+`" alt="" width="180"></div>
						<div class="z-index-1 position-relative d-flex flex-center"><h4 class= "mb-3"  style="color:#CF7AC2">User Registration</h4></div>
					
						<div class="row g-0 h-100 d-flex justify-content-center">
					
							<div class="col-5" style="margin-right: 1rem;">
								<div class="mb-3">
									<label class="form-label  font-weight-bold text-dark" for="UserType">User Type *</label>
									<select id="UserType" name="UserType"  data-name="UserType"  data-type="text"  class="form-control" aria-label="Sizing example input" aria-describedby="inputGroup-sizing-default" tabindex="1">
										<option value="Customer" class="form-control" aria-label="Sizing example input" aria-describedby="inputGroup-sizing-default">Customer</option>
										<option value="Operator"  class="form-control" aria-label="Sizing example input" aria-describedby="inputGroup-sizing-default">Operator</option>
										<option value="Admin"  class="form-control" aria-label="Sizing example input" aria-describedby="inputGroup-sizing-default">Admin</option>
									</select>
								</div>	 
								<div class="mb-3">
									<label  class="form-label font-weight-bold text-dark" for="FirstName">First Name *</label>
									<input  class="form-control" data-type="text" aria-label="Sizing example input" aria-describedby="inputGroup-sizing-default" type="text" id="FirstName" name="FirstName" data-name="FirstName" placeholder="First Name" required tabindex="3">
								</div>
								<div class="mb-3">
									<label class="form-label  font-weight-bold text-dark" for="Password" >Password *</label>
									<input  class="form-control" data-type="text" aria-label="Sizing example input" aria-describedby="inputGroup-sizing-default" type="password" id="Password" name="Password"  data-name="Password" placeholder="Password" required tabindex="5">
								</div>
								<div class="mb-3">
									<label class="form-label  font-weight-bold text-dark" for="Code" >Code *</label>
									<input  class="form-control" data-type="text" aria-label="Sizing example input" aria-describedby="inputGroup-sizing-default" type="text" id="Code" name="Code" data-name="Code" placeholder="Code" required tabindex="5">
								</div>
							</div>
							<div class="col-5">
								<div class="mb-3">
									<label class="form-label  font-weight-bold text-dark" for="Email">Email *</label>
									<input  class="form-control" data-type="text" aria-label="Sizing example input" aria-describedby="inputGroup-sizing-default" type="email" id="Email" name="Email" data-name="Email" placeholder="Email" required tabindex="2">
								</div>
								<div class="mb-3">
									<label class="form-label  font-weight-bold text-dark" for="LastName">Last Name *</label>
									<input  class="form-control" data-type="text" aria-label="Sizing example input" aria-describedby="inputGroup-sizing-default" type="text" id="LastName" name="LastName" data-name="LastName" placeholder="Last Name" required tabindex="4">
								</div>
								<div class="mb-3">
									<label class="form-label  font-weight-bold text-dark" for="ConfirmPassword">Confirm Password *</label>
									<input  class="form-control" data-type="text" aria-label="Sizing example input" aria-describedby="inputGroup-sizing-default" type="Password" id="ConfirmPassword" name="ConfirmPassword" data-name="ConfirmPassword" placeholder="Confirm Password" required tabindex="6">
								</div>
							</div>
						</div>
					<div class="modal-footer">
						<div class="col-6">
							<a href="/login" style="color:#CF7AC2;margin-left: 3.5rem;"><b>Login?</b></a>
						</div>
						<div class="col-6 d-flex justify-content-end" style="padding-right:2.2rem">
							<button type="button" class="btn" data-submit="registerForm" id="register" style="background:#CF7AC2;color:white;margin-right:1rem"><b>Register</b></button>
							
						</div>
					</div>
				</form>
				`, invalidRegistration, `
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

		</main>
		<script>
			document.addEventListener("DOMContentLoaded", function () {

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

					$(document).on("click","#register",function(){
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

					
						$.post("/sign-up", JSON.stringify(result), function(xhr, status, error) {}, 'json')
							.done(function(response) {
								window.location.href = "/login?&error=User registered successfully and is waiting for approval"
							})
							.fail(function(xhr, status, error) {
							 	const data = JSON.parse(xhr.responseText);
								window.location.href = "/register?status=" + xhr.status + "&error=" +data.message;
							});						
					}
				})
				function removeQueryParams() {		
					var newUrl = window.location.protocol + "//" + window.location.host + window.location.pathname;
					window.history.replaceState({}, document.title, newUrl);
					
				}
		})

				const userType = document.getElementById("UserType");
				const codeDiv = document.getElementById("Code").closest(".mb-3");
				const codeInput = document.getElementById("Code");

				userType.addEventListener("change", function () {
					if (userType.value == "Customer" || userType.value == "Operator") {
						codeDiv.style.display = "block";
						codeInput.required = true;
					} else {
						codeDiv.style.display = "none"; 
						codeInput.required = false; 
				}
				});
				userType.dispatchEvent(new Event("change"));
		</script>

		<script src="/static/vendors/chart/chart.min.js"></script>
		<script src="/static/vendors/countup/countUp.umd.js"></script>
		<script src="/static/vendors/echarts/echarts.min.js"></script>
		<script src="/static/vendors/d3/d3.min.js"></script>
		<script src="/static/vendors/lodash/lodash.min.js"></script>
		
		<script src="/static/vendors/list.js/list.min.js"></script>
		<script src="/static/vendors/bootstrap/bootstrap.min.js"></script>	
	  </body>
	
	  <!--!html-->
	</html>
	`)
	return
}
