package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"

	m "irpl.com/kanban-commons/model"
	"irpl.com/kanban-commons/utils"
	"irpl.com/kanban-web/pages/basepage"
	stagemanagement "irpl.com/kanban-web/pages/stageManagement"
	s "irpl.com/kanban-web/services"
)

func stageManagementPage(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("X-Custom-Username")
	userID := r.Header.Get("X-Custom-Userid")
	usertype := r.Header.Get("X-Custom-Role")
	links := r.Header.Get("X-Custom-Allowlist")

	var vendorName string
	var vendorRecord []m.Vendors
	if usertype != "Admin" {
		resp, err := http.Get(RestURL + "/get-vendor-by-userid?key=user_id&value=" + userID)
		if err != nil {
			slog.Error("%s - error - %s", "Error making GET request", err)
		}

		defer resp.Body.Close()

		if err := json.NewDecoder(resp.Body).Decode(&vendorRecord); err != nil {
			slog.Error("error decoding response body", "error", err)
		}
		if len(vendorRecord) != 0 {
			vendorName = vendorRecord[0].VendorName
		} else {
			vendorName = ""
		}

	} else {
		vendorName = ""
	}
	// Define the side navigation items
	sideNav := basepage.SideNav{
		MenuItems: []basepage.SideNavItem{
			{
				Name:     "Dashboard",
				Icon:     "fas fa-chart-pie",
				Link:     "/dashboard",
				UserType: basepage.UserType{Admin: true, Operator: true, Customer: true},
				Style:    "font-size:1rem;",
			},
			{
				Name:     "User Master",
				Icon:     "fas fa-users",
				Link:     "/user-management",
				UserType: basepage.UserType{Admin: true},
				Style:    "font-size:1rem;",
			},
			{
				Name:     "Vendor Master",
				Icon:     "fa fa-briefcase",
				Link:     "/vendor-management",
				UserType: basepage.UserType{Admin: true},
				Style:    "font-size:1rem;",
			},
			{
				Name:     "Operator Master",
				Icon:     "fa fa-user",
				Link:     "/operator-management",
				UserType: basepage.UserType{Admin: true},
				Style:    "font-size:1rem;",
			},
			{
				Name:     "Part Master",
				Icon:     "fas fa-vials",
				Link:     "/compounds-management",
				UserType: basepage.UserType{Admin: true, Operator: true},
				Style:    "font-size:1rem;",
			},
			{
				Name:     "Chemical Type Master",
				Icon:     "fas fa-vial",
				Link:     "/chemical-management",
				UserType: basepage.UserType{Admin: true},
				Style:    "font-size:1rem;",
			},
			{
				Name:     "Raw Material",
				Icon:     "fas fa-boxes",
				Link:     "/material-management",
				UserType: basepage.UserType{Admin: true},
				Style:    "font-size:1rem;",
			},
			{
				Name:     "Rubber Store Master",
				Icon:     "fas fa-memory",
				Link:     "/inventory-management",
				UserType: basepage.UserType{Admin: true},
				Style:    "font-size:1rem;",
			},
			{
				Name:     "Recipe Master",
				Icon:     "fa fa-clipboard-list",
				Link:     "/recipe-management",
				UserType: basepage.UserType{Admin: true, Operator: true},
				Style:    "font-size:1rem;",
				Selected: true,
			},
			{
				Name:     "Machine Master",
				Icon:     "fas fa-sliders-h",
				Link:     "/prod-line-management",
				UserType: basepage.UserType{Admin: true},
				Style:    "font-size:1rem;",
			},
			{
				Name:     "Process Master",
				Icon:     "fas fa-project-diagram",
				Link:     "/prod-processes-management",
				UserType: basepage.UserType{Admin: true},
				Style:    "font-size:1rem;",
			},
			{
				Name:     "Kanban Report",
				Icon:     "fas fa-list-alt",
				Link:     "/vendor-orders",
				UserType: basepage.UserType{Admin: true, Operator: true, Customer: true},
				Style:    "font-size:1rem;",
			},
			{
				Name:     utils.DefaultsMap["cold_store_menu"],
				Icon:     "fas fa-store",
				Link:     "/cold-storage",
				UserType: basepage.UserType{Admin: true},
				Style:    "font-size:1rem;",
			},
			{
				Name:     "Pending Orders",
				Icon:     "fas fa-list-alt",
				Link:     "/admin-orders",
				UserType: basepage.UserType{Admin: true, Operator: true, Customer: true},
				Style:    "font-size:1rem;",
			},
			{
				Name:     "All Kanban View",
				Icon:     "fa fa-th-list",
				Link:     "/all-kanban-view",
				Style:    "font-size:1rem;",
				UserType: basepage.UserType{Admin: true, Operator: true},
			},
			{
				Name:     "Kanban Board",
				Icon:     "fas fa-tasks",
				Link:     "/vendor-company",
				UserType: basepage.UserType{Admin: true, Operator: true},
				Style:    "font-size:1rem;",
			},
			{
				Name:     "Heijunka Board",
				Icon:     "fas fa-calendar-alt",
				Link:     "/production-line",
				UserType: basepage.UserType{Admin: true, Operator: true},
				Style:    "font-size:1rem;",
			},
			{
				Name:     "Kanban Entry",
				Icon:     "fas fa-plus",
				Link:     "/order-entry",
				UserType: basepage.UserType{Customer: true},
				Style:    "font-size:1rem;",
			},
			{
				Name:     "Quality Testing",
				Icon:     "fas fa-check-double",
				Link:     "/quality-testing",
				UserType: basepage.UserType{Admin: true},
				Style:    "font-size:1rem;",
			},
			{
				Name:     "Packing/Dispatch",
				Icon:     "fas fa-truck",
				Link:     "/packing-dispatch-page",
				UserType: basepage.UserType{Admin: true},
				Style:    "font-size:1rem;",
			},
			{
				Name:     "Kanban History",
				Icon:     "fas fa-history",
				Link:     "/kanban-history",
				UserType: basepage.UserType{Admin: true},
				Style:    "font-size:1rem;",
			},
			{
				Name:     "Order History",
				Icon:     "fas fa-file-alt",
				Link:     "/order-history",
				UserType: basepage.UserType{Admin: true},
				Style:    "font-size:1rem;",
			},
			{
				Name:     "Kanban Report",
				Icon:     "fas fa-scroll",
				Link:     "/kanban-report",
				UserType: basepage.UserType{Admin: true},
				Style:    "font-size:1rem;",
			},
			{
				Name:     "Kanban Reprint",
				Icon:     "fas fa-print",
				Link:     "/kanban-reprint",
				UserType: basepage.UserType{Admin: true},
				Style:    "font-size:1rem;",
			},
			{
				Name:     "Summary Reprint",
				Icon:     "fas fa-file",
				Link:     "/summary-reprint",
				UserType: basepage.UserType{Admin: true},
				Style:    "font-size:1rem;",
			},
		},
	}
	sideNav.MenuItems = basepage.CheckdisabledNavItems(sideNav.MenuItems, links, "|")

	// Define the top navigation items
	topNav := basepage.TopNav{VendorName: vendorName, UserType: usertype,
		MenuItems: []basepage.TopNavItem{
			{
				ID:    "settings",
				Name:  "",
				Title: "Settings",
				Type:  "link",
				Icon:  "fa fa-cog",
				Link:  "/configuration-page",
				Width: "col-1",
				Style: "bg-light bg-gradient",
			},
			{
				ID:    "notifications",
				Name:  "",
				Title: "Notifications",
				Type:  "link",
				Icon:  "fa fa-bell",
				Link:  "/system-logs",
				Width: "col-1",
				Style: "bg-light bg-gradient",
			},
			{
				ID:    "username",
				Title: "User Name",
				Name:  username,
				Type:  "button",
				Width: "col-2",
				Style: "bg-light bg-gradient",
			},
			{
				ID:    "logout",
				Title: "Log out",
				Name:  "",
				Link:  "/logout",
				Type:  "link",
				Icon:  "fas fa-sign-out-alt",
				Width: "col-1",
				Style: "bg-light bg-gradient",
			},
		},
	}

	stageMgntPage := stagemanagement.StageManagement{
		Username: username,
		UserID:   userID,
		UserType: usertype,
	}

	// Build the complete BasePage with navigation and content
	page := (&basepage.BasePage{
		SideNavBar: sideNav,
		TopNavBar:  topNav,
		Username:   username,
		UserType:   usertype,
		Content:    stageMgntPage.Build(),
	}).Build()

	// Write the complete HTML page to the response
	w.Write([]byte(page))
}

// Create STAGE Entry
func CreateStage(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-Custom-Userid")
	var stage m.Stage
	err := json.NewDecoder(r.Body).Decode(&stage)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		log.Println("Failed to read request body:", err)
		return
	}
	v := []map[string]string{}
	json.Unmarshal(stage.Headers, &v)

	stage.CreatedBy = userID
	jsonValue, err := json.Marshal(stage)
	if err != nil {
		http.Error(w, "Failed to marshal the data", http.StatusBadRequest)
		log.Println("Failed to marshal the data:", err)
		return
	}
	url := utils.RestURL + "/create-new-stage"
	res, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Printf("Failed to request: %v", err)
		http.Error(w, "Failed to Request", http.StatusForbidden)
		return
	}
	defer res.Body.Close()

	client := &http.Client{}
	resp, err := client.Do(res)
	if err != nil {
		log.Printf("Failed to request: %v", err)
		http.Error(w, "Failed to Request", http.StatusForbidden)
		return
	}
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("%s - error - %s", "Error reading response body", err)
		utils.SetResponse(w, http.StatusInternalServerError, "Failed to read response body")
		return
	}

	utils.SetResponse(w, resp.StatusCode, string(responseBody))
}

// Update Stage Entry
func UpdateExistingStage(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-Custom-Userid")
	var stage m.Stage
	err := json.NewDecoder(r.Body).Decode(&stage)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		log.Println("Failed to read request body:", err)
		return
	}
	v := []map[string]string{}
	json.Unmarshal(stage.Headers, &v)

	stage.ModifiedBy = userID

	jsonValue, err := json.Marshal(stage)
	if err != nil {
		http.Error(w, "Failed to marshal the data", http.StatusBadRequest)
		log.Println("Failed to marshal the data:", err)
		return
	}

	url := utils.RestURL + "/update-stage"
	res, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Printf("Failed to request: %v", err)
		http.Error(w, "Failed to Request", http.StatusForbidden)
		return
	}
	defer res.Body.Close()

	client := &http.Client{}
	resp, err := client.Do(res)
	if err != nil {
		log.Printf("Failed to request: %v", err)
		http.Error(w, "Failed to Request", http.StatusForbidden)
		return
	}
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("%s - error - %s", "Error reading response body", err)
		utils.SetResponse(w, http.StatusInternalServerError, "Failed to read response body")
		return
	}

	utils.SetResponse(w, resp.StatusCode, string(responseBody))
}

// Delete UserRole Entry
func DeleteStageByID(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-Custom-Userid")
	var stage m.Stage
	err := json.NewDecoder(r.Body).Decode(&stage)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		log.Println("Failed to read request body:", err)
		return
	}
	stage.CreatedBy = userID
	jsonValue, err := json.Marshal(stage)
	if err != nil {
		http.Error(w, "Failed to marshal the data", http.StatusBadRequest)
		log.Println("Failed to marshal the data:", err)
		return
	}

	url := utils.RestURL + "/delete-stage"
	res, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Printf("Failed to request: %v", err)
		http.Error(w, "Failed to Request", http.StatusForbidden)
		return
	}
	defer res.Body.Close()

	client := &http.Client{}
	resp, err := client.Do(res)
	if err != nil {
		log.Printf("Failed to request: %v", err)
		http.Error(w, "Failed to Request", http.StatusForbidden)
		return
	}
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("%s - error - %s", "Error reading response body", err)
		utils.SetResponse(w, http.StatusInternalServerError, "Failed to read response body")
		return
	}

	utils.SetResponse(w, resp.StatusCode, string(responseBody))
}

func GetAllStage(w http.ResponseWriter, r *http.Request) {
	url := utils.RestURL + "/get-all-stages"

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		http.Error(w, "Failed to create request", http.StatusBadRequest)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to perform request: %v", err)
		http.Error(w, "Failed to perform request", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		http.Error(w, "Failed to read response body", http.StatusInternalServerError)
		return
	}

	utils.SetResponse(w, resp.StatusCode, string(responseBody))
}

// Get Stages by Parameter
func GetStagesByParam(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")
	resp, err := http.Get(utils.RestURL + "/get-stage-by-param?key=" + key + "&value=" + value)
	if err != nil {
		slog.Error("%s - error - %s", "Error making GET request", err)
		utils.SetResponse(w, http.StatusInternalServerError, "Failed to find record")
		return
	}
	defer resp.Body.Close()
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("%s - error - %s", "Error reading response body", err)
		utils.SetResponse(w, http.StatusInternalServerError, "Failed to read response body")
		return
	}
	utils.SetResponse(w, resp.StatusCode, string(responseBody))
}

func GetStagesByHeader(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	url := utils.RestURL + "/get-stages-by-header"
	res, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Failed to request: %v", err)
		http.Error(w, "Failed to Request", http.StatusForbidden)
		return
	}
	defer res.Body.Close()

	client := &http.Client{}
	resp, err := client.Do(res)
	if err != nil {
		log.Printf("Failed to request: %v", err)
		http.Error(w, "Failed to Request", http.StatusForbidden)
		return
	}
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("%s - error - %s", "Error reading response body", err)
		utils.SetResponse(w, http.StatusInternalServerError, "Failed to read response body")
		return
	}

	utils.SetResponse(w, resp.StatusCode, string(responseBody))
}

// Stage edit dialog
func StageMangementEditDialog(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	var stages []*m.Stage
	var rawQuery m.RawQuery
	rawQuery.Host = utils.DefaultRestHost
	rawQuery.Port = utils.DefaultRestPort
	rawQuery.Type = "Stage"
	rawQuery.Query = `SELECT * FROM stage WHERE id = ` + id //`;`
	rawQuery.RawQry(&stages)

	formBtn := utils.JoinStr(`
	<div class="col-md-6 d-flex justify-content-end">
		<button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Cancel</button>
		<button id="edit-stage" data-url="/update-stage"  data-submit="editstagemodel" class="btn btn-primary ms-3">Update</button>
	</div>
`)

	var fields []s.FormField

	if len(stages) > 0 {
		st := stages[0]

		// Static base fields
		fields = append(fields,
			s.FormField{Label: "ID", ID: "ID", Width: "hidden", Type: "text", DataType: "int", Value: fmt.Sprint(st.ID)},
			s.FormField{Label: "Name", ID: "Name", Width: "100%", Type: "text", DataType: "string", Value: st.Name, Placeholder: "Enter name", IsRequired: true},
		)

		// Parse headers
		var headers []map[string]string
		if err := json.Unmarshal(st.Headers, &headers); err != nil {
			log.Printf("error decoding headers: %v", err)
		}

		// Add fields from headers
		for i, header := range headers {
			// Determine selected type
			fieldType := header["type"]

			fieldTypeOptions := []s.DropDownOptions{
				{Value: "input", Text: "Input", Selected: fieldType == "input"},
				{Value: "dropdown-RM", Text: "RM Dropdown", Selected: fieldType == "dropdown-RM"},
				{Value: "dropdown-CT", Text: "CT Dropdown", Selected: fieldType == "dropdown-CT"},
			}

			fields = append(fields,
				s.FormField{Label: "Field Name", ID: "FieldName-" + fmt.Sprint(i+1), Width: "33%", Type: "text", DataType: "string", Placeholder: "Enter field name", Value: header["field"]},
				s.FormField{Label: "Field Type", ID: "FieldType-" + fmt.Sprint(i+1), Width: "33%", Type: "select", DropDownOptions: fieldTypeOptions},
				s.FormField{Label: "Delete", ID: "delete-header", Width: "25%", Type: "button", AdditionalAttr: "class='btn btn-danger w-100 mt-4 delete-btn' onclick='removeField(this)'"},
			)
		}

		// Final row: fresh "Add" field input
		fields = append(fields,
			s.FormField{Label: "Field Name", ID: "FieldName-" + fmt.Sprint(len(headers)+1), Width: "33%", Type: "text", DataType: "string", Placeholder: "Enter field name"},
			s.FormField{Label: "Field Type", ID: "FieldType-" + fmt.Sprint(len(headers)+1), Width: "33%", Type: "select", DropDownOptions: []s.DropDownOptions{
				{Value: "input", Text: "Input", Selected: true},
				{Value: "dropdown-RM", Text: "RM Dropdown"},
				{Value: "dropdown-CT", Text: "CT Dropdown"},
			}},
			s.FormField{Label: "Add", ID: "add-header", Width: "25%", Type: "button", AdditionalAttr: "class='btn btn-primary w-100 mt-4' onclick='addField(this)'"},
		)
	}

	//Add Dialogue
	editStageModel := s.Model{
		ID:    "EditStageModel",
		Type:  "modal-lg",
		Title: `<h5 class="modal-title text-primary fs-3" id="staticBackdropLabel"><b style="color:#871A83">Edit Stage</b></h5>`,
		Sections: []s.FormSection{
			{
				ID:     "edit-new-stage-section",
				Fields: fields,
			},
		},
		Buttons: formBtn,
		JS: `
		<script>

function addField(button) {
    if (!$('#EditStageModel').is(':visible')) return;

    // Convert current button to Delete
    button.textContent = 'Delete';
    button.classList.remove('btn-primary');
    button.classList.add('btn-danger', 'delete-btn');
    button.id = 'delete-header';
    button.setAttribute('onclick', 'removeField(this)');

    const row = document.querySelector('#EditStageModel #field-container-edit-new-stage-section .row');
    const existingFieldNames = row.querySelectorAll('input[id^="FieldName-"]');
    const fieldCount = existingFieldNames.length + 1;

    // Field Name
    const fieldNameCol = document.createElement('div');
    fieldNameCol.className = 'col-4 mb-0';
    fieldNameCol.innerHTML =
        '<label class="container-fluid p-0 my-1" for="FieldName-' + fieldCount + '">Field Name</label>' +
        '<input id="FieldName-' + fieldCount + '" data-name="FieldName-' + fieldCount + '" type="text" data-type="string" class="form-control p-2" placeholder="Enter field name">';

    // Field Type
    const fieldTypeCol = document.createElement('div');
    fieldTypeCol.className = 'col-4 mb-0';
    fieldTypeCol.innerHTML =
        '<label class="container-fluid p-0 my-1" for="FieldType-' + fieldCount + '">Field Type</label>' +
        '<select class="form-select" id="FieldType-' + fieldCount + '" data-name="FieldType-' + fieldCount + '" data-type="">' +
            '<option value="input" selected>Input</option>' +
            '<option value="dropdown-RM">RM Dropdown</option>' +
			'<option value="dropdown-CT">CT Dropdown</option>' +
        '</select>';

    // Add Button
    const buttonCol = document.createElement('div');
    buttonCol.className = 'col-3 mb-0';
    buttonCol.innerHTML =
        '<button id="add-header" data-name="add-header" name="add-header" class="btn btn-primary w-100 mt-4" onclick="addField(this)">Add</button>';

    row.appendChild(fieldNameCol);
    row.appendChild(fieldTypeCol);
    row.appendChild(buttonCol);
}

function removeField(button) {
    if (!$('#EditStageModel').is(':visible')) return;

    const buttonCol = button.closest('.col-3');
    const fieldTypeCol = buttonCol?.previousElementSibling;
    const fieldNameCol = fieldTypeCol?.previousElementSibling;

    if (buttonCol) buttonCol.remove();
    if (fieldTypeCol) fieldTypeCol.remove();
    if (fieldNameCol) fieldNameCol.remove();

    renumberFields();
}


function renumberFields() {
    const row = document.querySelector('#EditStageModel #field-container-edit-new-stage-section .row');
    if (!row) return;

    const fieldNameCols = row.querySelectorAll('input[id^="FieldName-"]');
    const fieldTypeCols = row.querySelectorAll('select[id^="FieldType-"]');
    const buttons = row.querySelectorAll('button');

    for (let i = 0; i < fieldNameCols.length; i++) {
        const index = i + 1;
        const input = fieldNameCols[i];
        input.id = 'FieldName-' + index;
        input.setAttribute('data-name', 'FieldName-' + index);

        const label = input.closest('.col-4')?.querySelector('label');
        if (label) label.setAttribute('for', 'FieldName-' + index);
    }

    for (let i = 0; i < fieldTypeCols.length; i++) {
        const index = i + 1;
        const select = fieldTypeCols[i];
        select.id = 'FieldType-' + index;
        select.setAttribute('data-name', 'FieldType-' + index);

        const label = select.closest('.col-4')?.querySelector('label');
        if (label) label.setAttribute('for', 'FieldType-' + index);
    }

    for (let i = 0; i < buttons.length; i++) {
        const btn = buttons[i];
        if (i === buttons.length - 1) {
            btn.textContent = 'Add';
            btn.className = 'btn btn-primary w-100 mt-4';
            btn.id = 'add-header';
            btn.setAttribute('onclick', 'addField(this)');
        } else {
            btn.textContent = 'Delete';
            btn.className = 'btn btn-danger delete-btn w-100 mt-4';
            btn.id = 'delete-header';
            btn.setAttribute('onclick', 'removeField(this)');
        }
    }
}

</script>`,
	}
	w.Write([]byte(editStageModel.Build()))
}
