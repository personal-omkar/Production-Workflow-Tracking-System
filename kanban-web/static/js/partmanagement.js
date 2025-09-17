$(function(){  
    $(document).on("click","#add-part-submit,#edit-part-submit",function(){
					
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
                    var selectedValue = $(this).find(":selected").val();
                    if ($(this).attr("data-type") === "int") {
                        result[$(this).attr("data-name")] = parseInt(selectedValue, 10); // Convert to integer
                    } else if ($(this).attr("data-type") == "bool") {
                        if ($(this).val()=="true"){
                            result[$(this).attr("data-name")] =true;
                        }else{
                            result[$(this).attr("data-name")] = false;
                        }
                        
                    } else {
                        result[$(this).attr("data-name")] = selectedValue; // Keep as string
                    }
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
                }else { 
                    result[$(this).attr("data-name")] = $(this).val();
                }						
            })
              
                    
                $.post("/add-compound", JSON.stringify(result), function (xhr, status, error) {
                    window.location.href = "/part-management?status=" + xhr.code + "&msg=" + xhr.message;
                    }, 'json').fail(function (xhr, status, error) {
                         window.location.href = "/part-management?status=" + xhr.code + "&msg=" + xhr.message;	
                        
                    });													
        }
    })
            const urlParams = new URLSearchParams(window.location.search);
			const status = urlParams.get('status');
			const msg = urlParams.get('msg');

			if (status) {
				showNotification(status, msg, removeQueryParams);
			}

		
			// Show notifications
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

				notification.html(message).show();
				setTimeout(() => {
					notification.fadeOut(callback);
				}, 5000);
			}

			// Remove query parameters from URL
			function removeQueryParams() {
				var newUrl = window.location.protocol + "//" + window.location.host + window.location.pathname;
				window.history.replaceState({}, document.title, newUrl);
			}

            $(document).on("click", "#edit-part-btn", function(event) {  
                var data = JSON.parse($(this).closest("tr").attr("data-data"));

                // Get modal content for the clicked user
                $.get("/part-management-card?key=id&value=" + data.ID, function(response) {
                    let modal = $("#EditPartModel");
                    
                    if (modal.length) {
                        $("#EditPartModel").replaceWith(response);
                    } else {
                        $(".main-container").append(response);
                    }

                    // Show the modal after updating/appending
                    $("#EditPartModel").modal("show");
                }, 'json');

                
            })
});