package configuration

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"

	mo "irpl.com/kanban-commons/model"
	"irpl.com/kanban-commons/utils"
	s "irpl.com/kanban-web/services"
)

type ConfigurationPage struct {
}

var Pagination mo.PaginationResp
var MachinMasterTable s.TableCard

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

func (m *ConfigurationPage) Build() string {

	content := utils.JoinStr(`

	<div class="container-xxl mt-3" >
		<div class="card">
			<div class="card-header d-flex justify-content-between align-items-center" style="background-color:#F4F5FB">
	
				
				<div class="row d-flex align-items-center">
					<div class="col-auto">
						<h3 class="heading-text mb-0">System Configuration</h3>
					</div>
					<div class="col-auto ms-3">
						<span id="notification" class="alert alert-success alert-dismissible fade show mb-0 mt-1 p-0" role="alert" style="display: none;">
							
						</span>
					</div>
				</div>

			</div>
			<div class="card-body" id="system-config-card">
				`, GetLDAPConfigCard(), GetSAMBAConfigCard(), `
			</div>
				<div class="card-footer text-muted d-flex justify-content-end  align-items-center">
			</div>
		</div>
	</div>

	 <div class="modal" id="uploadModal" tabindex="-1">
		<div class="modal-dialog">
			<div class="modal-content">
				<div class="modal-header">
					<h5 class="modal-title">Uploading File...</h5>
				</div>
				<div class="modal-body" style="display: flex; flex-direction: column; align-items: center; justify-content: center; height: 200px;">
					<div class="spinner-border text-primary" role="status">
						<span class="visually-hidden">Loading...</span>
					</div>
					<p id="uploadStatus" style="text-align: center; margin-top: 20px;">Please wait while the file is being uploaded.</p>
				</div>
			</div>
		</div>
	</div>

	`)

	return content
}

func GetLDAPConfigCard() string {
	var ldapConfig mo.LDAPConfig
	resp, err := http.Get(RestURL + "/get-default-ldap-config")
	if err != nil {
		slog.Error("%s - error - %s", "Error making POST request", err)

	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&ldapConfig); err != nil {
		slog.Error("error decoding response body", "error", err)
	}

	// <!--html-->
	ret := utils.JoinStr(`
		<div class="col-sm-12">
			<div class="card">
				<div class="card-body">
					<h5 class="card-title">LDAP Configuration</h5>
					<div class="row mt-3" data-group="ldap-config">

						<div class="col-md-12 d-none">
							<div class="input-group mb-2">
								<span class="input-group-text">ID</span>
								<input id="ID" data-name="ID"  type='hidden' value="`, strconv.Itoa(ldapConfig.ID), `" data-type="int" class="form-control">
							</div>
						</div>	
						<div class="col-md-12">
							<div class="input-group mb-2">
								<span class="input-group-text">Server URL</span>
								<input id="LDAP_URL" data-name="LDAP_URL"  type='text' value="`, ldapConfig.LDAPURL, `" data-type="string" class="form-control">
							</div>	
						</div>	
						<div class="col-md-12">
							<div class="input-group mb-2">
								<span class="input-group-text">Bind DN</span>
								<input id="Bind_DN" data-name="Bind_DN"  type='text' value="`, ldapConfig.BindDN, `" data-type="string" class="form-control">
							</div>	
						</div>	
						<div class="col-md-12">
							<div class="input-group mb-2">
								<span class="input-group-text">Password</span>
								<input id="Password" data-name="Password"  type='password' value="`, ldapConfig.Password, `" data-type="string" class="form-control">
							</div>	
						</div>	
						<div class="col-md-12">
							<div class="input-group mb-2">
								<span class="input-group-text">Domain (Base DN)</span>
								<input id="Base_DN" data-name="Base_DN"  type='text' value="`, ldapConfig.BaseDN, `" data-type="string" class="form-control">
							</div>
						</div>
						<div class="col-md-12">
							<div class="input-group mb-2">
								<span class="input-group-text">Unique Identifier (e.g. uid)</span>
								<input id="unique_identifier" data-name="unique_identifier"  type='text' value="`, ldapConfig.UniqueIdentifier, `" data-type="string" class="form-control">
							</div>
						</div>
					</div>
					<div class="mt-2 d-flex justify-content-end  align-items-end">	
						<button type="button" class="btn btn-secondary">Clear</button>
						<button id="save-ldap-config" data-submit="ldap-config" data-url="/secure/update-ldap-config" type="button" class="btn btn-primary ms-2">Update</button>
					</div>
				</div>
			</div>
		</div>
		`)
	// <!--!html-->

	return ret
}

func GetSAMBAConfigCard() string {
	var SAMBA mo.SambaConfig
	sambaresp, err := http.Get(RestURL + "/get-default-samba-config")
	if err != nil {
		slog.Error("%s - error - %s", "Error making POST request", err)

	}
	defer sambaresp.Body.Close()

	if err := json.NewDecoder(sambaresp.Body).Decode(&SAMBA); err != nil {
		slog.Error("error decoding response body", "error", err)
	}

	// <!--html-->
	ret := utils.JoinStr(`
		<div class="col-sm-12 mt-5">
			<div class="card">
				<div class="card-body">
					<h5 class="card-title">SAMBA Configuration</h5>
					<div class="row mt-3" data-group="samba-config">

						<div class="col-md-12">
							<div class="input-group mb-2 d-none">
								<span class="input-group-text">ID</span>
								<input id="ID" data-name="ID"  type='hidden' value="`, strconv.Itoa(SAMBA.ID), `" data-type="int" class="form-control">
							</div>	
						</div>	
						<div class="col-md-12">
							<div class="input-group mb-2">
								<span class="input-group-text">Username</span>
								<input id="Server_String" data-name="Server_String"  type='text' value="`, SAMBA.ServerString, `"  data-type="string" class="form-control">
							</div>	
						</div>	
						<div class="col-md-12">
							<div class="input-group mb-2"> 
								<span class="input-group-text">Password</span>
								<input id="WorkGroup" data-name="WorkGroup"  type='password' value="`, SAMBA.Workgroup, `"   data-type="string" class="form-control">
							</div>	
						</div>	
						<div class="col-md-12">
							<div class="input-group mb-2">
								<span class="input-group-text">SAP Export Location</span>
								<input id="Security" data-name="Security"  type='text'  value="`, SAMBA.Security, `"  data-type="string" class="form-control">
							</div>	
						</div>	
					</div>
					<div class="mt-2 d-flex justify-content-end  align-items-end">	
						<button type="button" class="btn btn-secondary">Clear</button>
						<button id="save-samba-config" data-submit="samba-config" data-url="/secure/update-samba-config" type="button" class="btn btn-primary ms-2">Update</button>
					</div>
				</div>
			</div>
		</div>
		`)
	// <!--!html-->
	js := `

<script>
	$(document).ready(function() {

		const urlParams = new URLSearchParams(window.location.search);
		const status = urlParams.get('status');
		const msg = urlParams.get('msg');

		if (status) {
		
			showNotification(status,msg, () => {
				removeQueryParams();
			});
		}

		$(document).on("click","#save-ldap-config, #save-samba-config",function(){

			var group = $(this).attr("data-submit");
			var url = $(this).attr("data-url");
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

			$.post(url, JSON.stringify(result), function(data) {
		
			}, 'json')
				.done(function(data, textStatus, jqXHR) {
					window.location.href = "/configuration-page?status=" + jqXHR.status + "&msg=" + encodeURIComponent(jqXHR.responseText);
				})
				.fail(function(jqXHR, textStatus, errorThrown) {
					window.location.href = "/configuration-page?status=" + jqXHR.status + "&msg=" + encodeURIComponent(jqXHR.responseText);
				});
			
		}
	})


	function showNotification(status, msg, callback) {
		const notification = $('#notification');
		var message = '';
		if (status === "200") {
			message = '<strong>Success!</strong> ' + "Record updated successfully" + '.';
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
	})
</script>
`
	return ret + js
}
