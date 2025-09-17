package orderentrypage

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	m "irpl.com/kanban-commons/model"
	"irpl.com/kanban-commons/utils"

	s "irpl.com/kanban-web/services"
)

type OrderEntrypage struct {
	Username   string
	UserType   string
	UserID     string
	VendorCode string
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
func (v OrderEntrypage) Build() string {

	var orderdetails []*m.OrderDetails
	var tabledata []*m.CustomerOrderDetails
	var compoundlist []m.Compounds
	var compoundnames []s.DropDownOptions
	// var todaysOrders string

	type TableConditions struct {
		Conditions []string `json:"Conditions"`
	}

	var vendor []*m.Vendors
	var rawQuery m.RawQuery
	rawQuery.Host = utils.RestHost
	rawQuery.Port = utils.RestPort
	rawQuery.Type = "Vendors"
	rawQuery.Query = `SELECT * FROM Vendors where vendor_code='` + v.VendorCode + `';`
	rawQuery.RawQry(&vendor)

	var kbdata []*m.KbData
	rawQuery.Host = utils.RestHost
	rawQuery.Port = utils.RestPort
	rawQuery.Type = "KbData"
	rawQuery.Query = `SELECT kb_data.* FROM kb_extension LEFT JOIN kb_data ON kb_data.kb_extension_id = kb_extension.id where vendor_id='` + strconv.Itoa(vendor[0].ID) + `' AND  kb_extension.status != 'creating' 
	  AND DATE(kb_data.demand_date_time) = CURRENT_DATE;`
	rawQuery.RawQry(&kbdata)

	var todayorders int
	for _, v := range kbdata {

		todayorders = todayorders + v.NoOFLots
	}
	var tablecondition TableConditions
	con := utils.JoinStr(`kb_data.created_by='`, v.UserID, `' AND kb_extension.status='creating'`)
	tablecondition.Conditions = append(tablecondition.Conditions, con)

	jsonValue, err := json.Marshal(tablecondition)
	if err != nil {
		slog.Error("%s - error - %s", "Error marshaling table condition", err)
	}
	var html strings.Builder

	resp, err := http.Post(RestURL+"/get-customer-order-details", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		slog.Error("%s - error - %s", "Error making GET request", err)
	}

	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&orderdetails); err != nil {
		slog.Error("error decoding response body", "error", err)
	}

	for _, i := range orderdetails {
		var temp m.CustomerOrderDetails
		temp.CompoundId = i.CompoundId
		temp.Id = i.Id
		temp.Location = i.Location
		temp.KbRootId = i.Id
		temp.VendorName = i.VendorName
		temp.CompoundName = i.CompoundName
		temp.CellNo = i.CellNo
		temp.NoOFLots = i.NoOFLots
		temp.Status = i.Status
		temp.CompoundName = i.CompoundName
		temp.DemandDateTime = i.DemandDateTime.String()
		temp.MFGDateTime = i.MFGDateTime.String()
		temp.CustomerNote = i.CustomerNote
		tabledata = append(tabledata, &temp)
	}
	//fetching compounds records
	compoundsresp, err := http.Get(RestURL + "/get-all-compound-data")
	if err != nil {
		slog.Error("%s - error - %s", "Error making GET request", err)
	}

	defer compoundsresp.Body.Close()

	if err := json.NewDecoder(compoundsresp.Body).Decode(&compoundlist); err != nil {
		slog.Error("error decoding response body", "error", err)
	}

	for _, i := range compoundlist {
		if i.Status {
			var comp s.DropDownOptions
			comp.Text = i.CompoundName
			comp.Value = strconv.Itoa(i.Id)
			compoundnames = append(compoundnames, comp)
		}
	}

	var orderEntryTable s.TableCard
	allCheckBox := `<input style="width: 20px; height: 20px;" class="form-check-input me-2" type="checkbox" value="" id="select_all_compound">`
	Conditional_Tools := map[bool]string{
		true: `<button  type="button" class="btn  m-0 p-0 edit-order-entry" > 
 					<i class="fa fa-edit " style="color: #871a83;"></i> 
					</button>
					
					<button   id="order-entry-delete" type="button" class="btn  m-0 p-0" data-toggle="tooltip" data-placement="bottom" title="Delete Order" data-bs-toggle="modal" data-bs-target="#DeleteOrderStatus" > 
 					<i class="fa fa-trash" aria-hidden="true"  style="color: #FF5C5C;"></i>
					</button>
					
					<button id="order-entry-submit" type="button" class="btn  m-0 p-0"  data-toggle="tooltip" data-placement="bottom" title="Submit Order" data-bs-toggle="modal" data-bs-target="#SubmitOrderStatus"> 
 					<i class="fa fa-arrow-alt-circle-right" aria-hidden="true"  style="color:#00D27A;"></i>
					</button>`,
		false: `<button  type="button" class="btn  m-0 p-0 edit-order-entry" > 
 					<i class="fa fa-edit " style="color: #871a83;"></i> 
					</button>
					
					<button   id="order-entry-delete" type="button" class="btn  m-0 p-0" data-toggle="tooltip" data-placement="bottom" title="Delete Order" data-bs-toggle="modal" data-bs-target="#DeleteOrderStatus" > 
 					<i class="fa fa-trash" aria-hidden="true"  style="color: #FF5C5C;"></i>
					</button>
				`,
	}

	orderEntryTable.CardHeading = `<button class="btn btn-secondary ms-auto" id="submitAllBtn" type="button" style="background-color:#871a83; color: white; border:none; cursor: pointer;" disabled>
									Submit
								</button>`

	orderEntryTable.CardHeadingActions.Component = []s.CardHeadActionComponent{
		{ComponentName: "span", ComponentType: s.ActionComponentElement{Span: s.SpanAttributes{Text: "Today Orders :" + strconv.Itoa(todayorders), Class: "text-danger fw-bold"}}, Width: "col-4"},
		{ComponentName: "span", ComponentType: s.ActionComponentElement{Span: s.SpanAttributes{Text: "Daily Limt :" + strconv.Itoa(vendor[0].PerDayLotConfig), Class: "text-danger fw-bold"}}, Width: "col-4"},
		{ComponentName: "span", ComponentType: s.ActionComponentElement{Span: s.SpanAttributes{Text: "Hour Limt :" + strconv.Itoa(vendor[0].PerHourLotConfig), Class: "text-danger fw-bold"}}, Width: "col-4"},
	}
	orderEntryTable.BodyTables = s.CardTableBody{Columns: []s.CardTableBodyHeadCol{
		{
			Name:  "Tools",
			Lable: " ",
		},
		{
			Lable: "Kanban Summary",
			Name:  `Cell No`,
			Width: "col-2",
		},
		{
			Lable: "Part Name",
			Name:  `Compound Name`,
			Width: "col-2",
		},
		{
			Lable:  "No. of Lots",
			Name:   "NoOFLots ",
			Width:  "col-2",
			GetSum: true,
		},
		// {
		// 	Name:  "Location",
		// 	Width: "col-1",
		// },
		{
			Name:  "Status",
			Width: "col-2",
		},
		{
			Name:  "Demand Date Time",
			Width: "col-2",
		},
		{
			Name:  "MFGDateTime",
			Width: "col-2",
		},
		{
			Name:  "Conditional_Tools",
			Lable: "Tools",
			Width: "col-2",
		},
	},
		ColumnsWidth:      []string{"col-1", "col-1", "col-2", "col-2", "col-2", "col-2", "col-2", "col-2"},
		Data:              tabledata,
		Tools:             allCheckBox,
		Conditional_Tools: Conditional_Tools,
		ID:                "VendorComp",
	}
	formBtn := utils.JoinStr(`<div class="col-md-4 d-flex justify-content-end">
	<button id="order-entry-clear" class="btn btn-secondary" style="width:100%;">Clear</button>
	<button id="order-entry-save" data-submit="kanbanentry" class="btn" style="margin-left:1rem;width:100%;  border:none; color:#FFFFFF; background-color:#871A83;" disabled>Create</button>
	
</div>
`)

	var customerNoteOptions = []s.DropDownOptions{
		{Value: "regular", Text: "Regular", Selected: true},
		{Value: "urgent", Text: "Urgent"},
		{Value: "mosturgent", Text: "Most Urgent"},
	}

	//Form Date
	form := s.Form{
		Title: "Kanban Entry",
		Style: `style="padding-left:0px;padding-right:5%;"`,
		Sections: []s.FormSection{
			{
				ID: "order-entry",
				Fields: []s.FormField{
					{Label: "Compound Code", ID: "CompoundCode", Width: "100%", Type: "dropdown", DropDownOptions: compoundnames},
					{Label: "Demand Date/Time", ID: "DemandDateTime", Width: "100%", Type: "text", DataType: "date", Placeholder: "Y-m-d H:i", DateOpts: `"minDate": "` + time.Now().String() + `", "defaultDate": "` + time.Now().Format("2006-01-02 15.04.03") + `", "maxDate": "` + time.Now().AddDate(0, 0, 7).String() + `"`, IsRequired: true},
					// {Label: "Location", ID: "Location", Width: "100%", Type: "Hidden", DataType: "text"},
					{Label: "No. of Lots", ID: "NoOFLots", Width: "100%", Type: "number", DataType: "number", Min: "1", Max: "999", Step: "1", IsRequired: true},
					{Label: "Cell Name", ID: "CellNo", Width: "hidden", Type: "hidden", DataType: "text"},
					// {Label: "Notes", ID: "Notes", Width: "100%", Type: "textarea", DataType: "text"},
					{Label: "Note", ID: "Note", Width: "100%", Type: "dropdownfield", DropDownOptions: customerNoteOptions},
				},
			},
		},
		Buttons: formBtn,
		Alert:   true,
	}

	confirmSubmit := s.ConfirmationModal{
		For:   "Submit",
		ID:    "SubmitOrderStatus",
		Title: "Confirm!",
		Body: []string{
			"You are about to submit this order [Cell Name: __], Once it is submitted you will no longer be able to edit/modify this record.",
			"For any modification, you need to contact admin.",
			"Are you sure you want to Submit this order?",
			"<b>Note: This order will be listed in your Ordered List where you can check it's status.</b>",
		},
		Footer: s.Footer{CancelBtn: false, Buttons: []s.FooterButtons{{BtnType: "submit", BtnID: "Close_Modal", Text: "Close", Style: "background-color:#636E7E; border:none;"}, {BtnType: "submit", BtnID: "Submit_Order", Text: "Submit"}}},
	}
	confirmDelete := s.ConfirmationModal{
		For:   "Delete",
		ID:    "DeleteOrderStatus",
		Title: "Attention!",
		Body: []string{
			"You are about to delete [Cell Name: __] this order, are you sure you want to perform this action?",
		},
		Footer: s.Footer{CancelBtn: false, Buttons: []s.FooterButtons{{BtnType: "submit", BtnID: "Close_Modal", Text: "Close", Style: "background-color:#636E7E; border:none;"}, {BtnType: "submit", BtnID: "Delete_Order", Text: "Delete"}}},
	}

	confirmMultiSubmit := s.ConfirmationModal{
		For:   "MultiSubmit",
		ID:    "SubmitMultipleOrders",
		Title: "Confirm Multiple Submissions!",
		Body: []string{
			"You are about to submit multiple Kanban entries for Cell(s): [CELL_LIST].",
			"For any modification, you need to contact admin.",
			"Are you sure you want to Submit this orders?",
			"<b>Note: This orders will be listed in your Ordered List where you can check it's status.</b>",
		},
		Footer: s.Footer{CancelBtn: false, Buttons: []s.FooterButtons{{BtnType: "button", BtnID: "Close_Modal", Text: "Close", Style: "background-color:#636E7E; border:none;"},
			{BtnType: "submit", BtnID: "Confirm_Multi_Submit", Text: "Submit All"}}},
	}

	html.WriteString(confirmSubmit.Build())
	html.WriteString(confirmDelete.Build())
	html.WriteString(confirmMultiSubmit.Build())

	js := ` <script>
		//js 
			function updateSubmitButtonState() {
				const today = new Date().toISOString().split('T')[0];
				const selectedCheckboxes = document.querySelectorAll('#select_all_compound:checked');

				if (selectedCheckboxes.length === 0) {
					document.getElementById('submitAllBtn').disabled = true;
					return;
				}

				const allSelectedAreToday = [...selectedCheckboxes].every(checkbox => {
					const row = checkbox.closest('tr');
					const datadata = row.getAttribute('data-data');
					if (!datadata) return false;
					try {
						const parsedData = JSON.parse(datadata);
						const demandDate = parsedData.DemandDateTime.split(' ')[0];
						return demandDate === today;
					} catch (e) {
						return false;
					}
				});

				document.getElementById('submitAllBtn').disabled = !allSelectedAreToday;
			}
					document.addEventListener("DOMContentLoaded", function () {
			});

			function toggleActionMenu() {
				const anyChecked = document.querySelectorAll('.component-code:checked').length > 0;
				document.getElementById('selectMenu-action-line').disabled = !anyChecked;

				// NEW: update Submit button
				updateSubmitButtonState();
			}

					$(function(){
						function convertToISO(dateTime) {
						
						 const [datePart, timePart] = dateTime.split(" ");
   	 						return datePart + 'T' + timePart + ':00Z';
						}
							
						document.querySelectorAll('#select_all_compound').forEach(selectAll => {
							selectAll.addEventListener('change', function() {
								const today = new Date().toISOString().split('T')[0]; // Get today's date (YYYY-MM-DD)
								const selectedCheckboxes = document.querySelectorAll('#select_all_compound:checked');
								if (selectedCheckboxes.length === 0) {
									document.getElementById('submitAllBtn').disabled = true;
									return;
								}
								const allSelectedAreToday = [...selectedCheckboxes].every(checkbox => {
									const row = checkbox.closest('tr');
									const datadata = row.getAttribute('data-data');
									if (!datadata) return false; 
									try {
										const parsedData = JSON.parse(datadata);
										const demandDateTimeStr = parsedData.DemandDateTime; 
										const demandDate = demandDateTimeStr.split(' ')[0]; 
										return demandDate === today;
									} catch (error) {
										console.error("Error parsing JSON:", error);
										return false;
									}
								});
								document.getElementById('submitAllBtn').disabled = !allSelectedAreToday;
							});
						});
						$("#submitAllBtn").click(function () {
							$(this).prop("disabled", true); // Disable button after click
							const selectedData = []; // Array to store selected items
							$("#select_all_compound:checked").each(function () {
								const row = $(this).closest('tr');
								const datadata = row.attr('data-data');
								if (!datadata) return;
								try {
									const parsedData = JSON.parse(datadata);
									const item = {
										CompoundCode: String(parsedData.CompoundId),
										DemandDateTime: convertToISO(parsedData.DemandDateTime).replace(/:00Z$/, 'Z'),
										NoOFLots: parsedData.NoOFLots,
										Status: "pending",
										CellNo: parsedData.CellNo,
										MFGDateTime: convertToISO(parsedData.MFGDateTime).replace(/:00Z$/, 'Z'),
									};
									selectedData.push(item); // Add to array
								} catch (error) {
									console.error("Error parsing JSON:", error);
								}
							});
							// Store selection globally
						window.multiSubmitData = selectedData;

						// Inject CellNos into modal
						const cellNos = selectedData.map(item => item.CellNo).join(", ");
					const modalBody = document.querySelector("#SubmitMultipleOrders #modal-body");
					if (modalBody) {
					modalBody.innerHTML = 
						'<p>You are about to submit multiple Kanban entries for Cell(s): <b>' + cellNos + '</b>.</p>' +
						'<p>For any modification, you need to contact admin.</p>' +
						'<p>Are you sure you want to Submit these orders?</p>' +
						'<b>Note: These orders will be listed in your Ordered List where you can check their status.</b>';
					}

						// Show modal
						$('#SubmitMultipleOrders').modal('show');
					});

			// Confirm Submit (final submission)
			$(document).on("click", "#Confirm_Multi_Submit", function () {
				if (window.multiSubmitData && window.multiSubmitData.length > 0) {
				
								$.post("/check-vendor-lot-limit", JSON.stringify(window.multiSubmitData), function (xhr, status, error) {
									$.post("/create-multi-new-order-entry", JSON.stringify(window.multiSubmitData), function (data) {
										var resp = JSON.parse(data);

										// ✅ Clear selection after successful submit
										window.multiSubmitData = [];
										$('#SubmitMultipleOrders').modal('hide');

										// Optional: uncheck checkboxes
										document.querySelectorAll('#select_all_compound').forEach(cb => cb.checked = false);
										document.getElementById('submitAllBtn').disabled = true;

										window.location.href = "/order-entry?status=" + resp.code + "&msg=" + resp.message;
									});
								}, 'json').fail(function (xhr, status, error) {
								 	const data = JSON.parse(xhr.responseText);
									$('body').append(data.dialog);
									$('#LotSizeExceed').modal('show');
									return { success: false };
									
								});	
			}
			});
			
			// Close modal (cancel warning) — DO NOT clear selection or disable submit button
			$(document).on("click", "#Close_Modal", function () {
				$('#SubmitMultipleOrders').modal('hide');
				updateSubmitButtonState();
			});
				const urlParams = new URLSearchParams(window.location.search);
				const status = urlParams.get('status');
				const msg = urlParams.get('msg');
				if (status) {
					var messages = msg.split(";")
					for (let index = 0; index < messages.length; index++) {
						const element = messages[index];
						if (element.includes("submit")) {

							var match = element.match(/CellNo:\s*([\w\/-]+)/);
							if (match) {
								var cellno = match[1];

								$.get("/cust/print-html-card?cellno=" + cellno, function(data) {
									$("#report-div").html(data);
									printContent();
								});
							}
						}
					}
					showNotification(status,msg, () => {
							removeQueryParams();
					});
				}


						$(".edit-order-entry").on("click" ,function(event){   
							$('#CompoundCode').prop('disabled', true); 
							var data = JSON.parse($(this).closest("tr").attr("data-data"));
							var option = new Option(data.CompoundName, data.CompoundId, true, true);
							$('#CompoundCode').append(option).trigger('change');
							var formatedatestring = data.DemandDateTime.replace(/\./g, '-');
							var date = new Date(formatedatestring); 
							var formattedDate = date.toISOString().replace("T", " ").replace(":00.000Z",""); 
							$('#DemandDateTime').val(formattedDate);
							$('#Location').val(data.Location);
							$('#NoOFLots').val(data.NoOFLots);
							$('#CellNo').val(data.CellNo);
							$('#Note').val(data.Note);
								
						})

						$("#order-entry-clear").on("click" ,function(event){   
							$('#CompoundCode').prop('disabled', false); 
							$('#CompoundCode').val("1");
							$('#DemandDateTime').val("");
							$('#NoOFLots').val("");
							$('#Location').val("");
							$('#CellNo').val("");
						})

						function formatDateMMDDYYYYToYYYYMMDD(inputDate) {
							// Split the input date into day, month, and year
							var dateParts = inputDate.split('-');
							
							// Ensure the date format is correct
							if (dateParts.length === 3) {
								var day = dateParts[0];
								var month = dateParts[1];
								var year = dateParts[2];

								// Return the formatted date in yyyy-MM-dd format
								return year + '-' + month + '-' + day;
							} else {
								// If the input date is not in the correct format
								return 'Invalid date format';
							}
						}

						$('#DemandDateTime, #NoOFLots, #Note,#CompoundCode').on('change', function () {
							 let allFilled = true;

							$('#DemandDateTime, #NoOFLots, #Note,#CompoundCode').each(function () {
								if ($(this).val().trim() === '') {
									allFilled = false;
									return false; // Exit the loop early
								}
							});

							$('#order-entry-save').prop('disabled', !allFilled); 
							
						});

						function formatDate(dateString) {
							var formatedatestring =dateString.replace(/\./g, '-')
							var parts = formatedatestring.split('-');
							if (parts.length !== 3) {
								throw new Error("Invalid date format. Expected YYYY-MM-DD");
							}

							var year = parts[0];
							var month = parts[1].padStart(2, '0');
							var day = parts[2].padStart(2, '0');

							// Validate the extracted parts
							if (isNaN(year) || isNaN(month) || isNaN(day)) {
								throw new Error("Invalid date parts. Ensure they are numbers.");
							}

							// Return the date in ISO 8601 format with time and timezone
							return '${year}-${month}-${day}T00:00:00Z';
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

						function removeQueryParams() {
							var newUrl = window.location.protocol + "//" + window.location.host + window.location.pathname;
							window.history.replaceState({}, document.title, newUrl);
						}

						
						$(document).on("click", "#order-entry-save", function () {
							$('#CompoundCode').prop('disabled', false); 
							var btnname=this.id;
							var group = $(this).attr("data-submit");
							var result = {};
							var validated = true;

							$("[data-group='" + group + "']").find("[data-name]").each(function () {
								if ($(this).attr("data-validate") && ($(this).val().trim().length === 0 || $(this).val() === "Nil")) {
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
							
							if (validated) {
								$("[data-group='" + group + "']").find("[data-name]").each(function () {
									if ($(this).is("select")) {
										result[$(this).attr("data-name")] = $(this).find(":selected").val();
									} else if ($(this).attr("data-type") == "date") {
										var userDate = $(this).val();
										if (userDate.includes(":")) {
											result[$(this).attr("data-name")] = userDate;
										} else {
											var formattedDate = formatDate(userDate);
											result[$(this).attr("data-name")] = new Date(userDate);
										}
									} else if ($(this).attr("data-type") == "int") {
										result[$(this).attr("data-name")] = parseInt($(this).val());
									} else if($(this).attr("data-type") == "number") {
										result[$(this).attr("data-name")] = parseInt($(this).val());
									}else {
										result[$(this).attr("data-name")] = $(this).val();
									}
								});

								if (btnname == "order-entry-save") {
									result["status"] = "creating";
								}
						
								const lotData = {
									CompoundCode:result.CompoundCode,
									CellNo:result.CellNo,
									DemandDateTime:convertToISO(result.DemandDateTime),
									NoOFLots: result.NoOFLots
								};
								
								result.DemandDateTime=convertToISO(result.DemandDateTime)
								
								$.post("/check-daily-lot-limit", JSON.stringify(lotData), function () {
									$.post("/create-new-order-entry", JSON.stringify(result), function (data)  {	
										var resp = JSON.parse(data)
									
										window.location.href = "/order-entry?status=" + resp.code + "&msg=" + resp.message;
									});
								}, 'json').fail(function (xhr, status, error) {
									console.log(xhr.responseText)
									var data=JSON.parse(xhr.responseText)

									 $('body').append(data.dialog);
									 $('#LotSizeExceed').modal('show');
									
								});
								
							}
						});
						
						function formatDateToISO(dateString) {
							// Assuming input format is "DD.MM.YYYY"
							const [day, month, year] = dateString.split(".");
							return new Date(year + "-" + month + "-" + day).toISOString(); 
						}

						$(document).on("click", "#Close_Modal, .btn-close", function() {
							window.location.reload();
						})

						$(document).on("click", "#order-entry-submit", function () {					
							const rowData = $(this).closest("tr").data("data");
							if (rowData) {
								const parsedData = typeof rowData === "string" ? JSON.parse(rowData) : rowData;
								const CompoundId = parsedData.CompoundId;
								const compoundCode = parsedData.CompoundName;
								const demandDateTime =parsedData.DemandDateTime;
								const NoOFLots = parsedData.NoOFLots;
								const status = "pending";
								const location = parsedData.Location;
								const cellNo = parsedData.CellNo;
								const mfgDateTime =parsedData.MFGDateTime;
								// Populate the data
								var modalBody = document.querySelector("#SubmitOrderStatus #modal-body");
								if (modalBody) {
									var bodyContent = modalBody.innerHTML;
									// Replace placeholders like [Cell Name] with actual cellNo
									modalBody.innerHTML = bodyContent.replace(/\[Cell Name:\s*__\]/, "<b>[Cell Name: " + cellNo + "]</b>");
								}
								const dataList = [];
								const data = {
									CompoundCode : String(CompoundId),
									DemandDateTime :convertToISO(demandDateTime).replace(/:00Z$/, 'Z'),
									NoOFLots :NoOFLots,
									Status : status,
									Location : location,
									CellNo :cellNo,
									MFGDateTime : convertToISO(mfgDateTime).replace(/:00Z$/, 'Z'),
									Note : parsedData.Note,
								}
								// Forward Data
								dataList.push(data)
								$(document).on("click", "#Submit_Order", function () {
									$.post("/check-vendor-lot-limit", JSON.stringify(dataList), function () {
										$.post("/create-new-order-entry", JSON.stringify(data), function(msg) {
											$('#SubmitOrderStatus').modal('hide')
											var data = JSON.parse(msg)
											window.location.href = "/order-entry?status=" + data.code + "&msg=" + data.message;
										})
									}, 'json').fail(function (xhr, status, error) {
									var data=JSON.parse(xhr.responseText)

									 $('body').append(data.dialog);
									 $('#LotSizeExceed').modal('show');
									
								});
								})
								// Close_Modal
								$(document).on("click", "#Close_Modal", function () {
									$('#SubmitOrderStatus').modal('hide')
								})
							}
						});

						// order-entry-delete
						$(document).on("click", "#order-entry-delete", function () {					
							const rowData = $(this).closest("tr").data("data");
							if (rowData) {
								const parsedData = typeof rowData === "string" ? JSON.parse(rowData) : rowData;
								const KDID = parsedData.ID;
								const cellNo = parsedData.CellNo; 

								// Populate the data
								var modalBody = document.querySelector("#DeleteOrderStatus #modal-body");
								if (modalBody) {
									var bodyContent = modalBody.innerHTML;
									// Replace the placeholder with the actual cell number
									modalBody.innerHTML = bodyContent.replace(/\[Cell Name:\s*__\]/, "<b>[Cell Name: " + cellNo + "]</b>");
								}


								const data = {
									id : KDID,
									CellNo : cellNo, 
								}
								// Forward Data
								$(document).on("click", "#Delete_Order", function () {
									fetch("/delete-order-entry",{
										method: 'DELETE',
										headers: {
											'Content-Type': 'application/json',
										},
										body: JSON.stringify(data)
									})
									.then(response => {
										if (!response.ok || response.ok) {
											return response.text().then(msg => {
												$('#SubmitOrderStatus').modal('hide')
												var data = JSON.parse(msg)
												
												window.location.href = "/order-entry?status=" + data.code + "&msg=" +  data.message;
												throw new Error(data.message);
											});
										}
									})
									.catch(error => {
										console.error("Error occurred:", error);
									});
								})
								// Close_Modal
								$(document).on("click", "#Close_Modal", function () {
									$('#DeleteOrderStatus').modal('hide')
								})
							} else {
								console.error("No data-data attribute found for the clicked row.");
							}
						});

					})
					
			$(document).ready(function() {
				$('#CompoundCode').select2({
					theme: 'bootstrap-5',
					placeholder: 'Search Compound',
					allowClear: true,
					ajax: {
						transport: function(params, success, failure) {
							var partName = params.data.q;
							$.ajax({
								url: '/get-compounds-list-by-parem',
								data: {
									partName: partName
								},
								dataType: 'json',
								success: function(data) {
									var transformedData = $.map(data, function(item) {
										return {
											id: item.Id,     
											text: item.Text
										};
									});

									success({ results: transformedData });
								},
								error: function() {
									failure();
								}
							});
						},
						minimumInputLength: 1
					}
				}).on('select2:open', function() {
					const select2SearchField = document.querySelector('.select2-container--bootstrap-5 .select2-search__field');
					if (select2SearchField) {
						select2SearchField.focus();
					}
				});
			});

				//!js
		   </script>`

	html.WriteString(`
	<div class="d-flex">
		<div class="col-4">` + form.GenerateForm() + `</div>` +
		`<div class="col-8"  style="max-height: 90vh; overflow-y: auto;">` + orderEntryTable.Build() + `</div> 
	</div>`)

	// html.WriteString(AddEntryCard.Build())
	html.WriteString(js)

	html.WriteString(`<div id="report-div"></div>`)

	return html.String()
}
