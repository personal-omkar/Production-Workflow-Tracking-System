package server

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"log/slog"
	"net/http"
	"sort"
	"strconv"
	"strings"

	m "irpl.com/kanban-commons/model"
	"irpl.com/kanban-commons/utils"
	"irpl.com/kanban-web/pages/basepage"
	"irpl.com/kanban-web/pages/vendorcompany"
	s "irpl.com/kanban-web/services"
)

func vendorCompanyPage(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("X-Custom-Username")
	usertype := r.Header.Get("X-Custom-Role")
	userID := r.Header.Get("X-Custom-Userid")
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
				Selected: true,
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

	// // Define the connections between steps (pairs of step indices)
	// connections := [][2]int{
	// 	{0, 1}, {1, 2}, {2, 3}, {3, 4},
	// 	{4, 5}, {5, 6}, {6, 7}, {7, 8},
	// 	{8, 9}, {9, 10},
	// }

	// Create the FlowchartPage with title, steps, and connections
	vendorcompanyPage := vendorcompany.VendorCompanyPage{}

	// Build the flowchart content for the BasePage

	// Build the complete BasePage with navigation and content
	page := (&basepage.BasePage{
		SideNavBar: sideNav,
		TopNavBar:  topNav,
		Username:   username,
		UserType:   usertype,
		Content:    vendorcompanyPage.Build(),
	}).Build()

	// Write the complete HTML page to the response
	w.Write([]byte(page))
}

// Create compound by vendor name
func AddCompoundsForVendor(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-Custom-Userid")
	var data m.AddCompoundsByVendor
	var ApiResp m.ApiRespMsg
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&data); err == nil {
		data.UserID = userID

		jsonValue, err := json.Marshal(data)
		if err != nil {
			slog.Error("%s - error - %s", "Error marshaling user data", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to marshal data")
			return
		}
		resp, err := http.Post(utils.RestURL+"/add-compound-data-by-vendor", "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			slog.Error("%s - error - %s", "Error making POST request", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to create Compound Entry")
			return
		}
		defer resp.Body.Close()
		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			slog.Error("%s - error - %s", "Error reading response body", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to read response body")
			return
		}
		ApiResp.Code = resp.StatusCode
		ApiResp.Message = string(responseBody)
		body, err := json.Marshal(ApiResp)
		if err != nil {
			log.Printf("Error marshaling response: %v", err)
		}
		utils.SetResponse(w, http.StatusOK, string(body))
		return
	} else {
		utils.SetResponse(w, http.StatusInternalServerError, "Failed to add the compound for the"+data.VendorName)
	}

}

// Add compound to production line
func AddCompoundsInProductionLine(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-Custom-Userid")
	type compoundData struct {
		LineID   string
		KbDataId []string
		UserID   string
	}
	var data compoundData
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err == nil {
		data.UserID = userID
		jsonValue, err := json.Marshal(data)
		if err != nil {
			slog.Error("%s - error - %s", "Error marshaling compound data", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to marshal compound data")
			return
		}

		resp, err := http.Post(utils.RestURL+"/add-compound-data-to-production-line", "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"status": "500", "msg": "Failed to add compound data in production line"})
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"status": "200", "msg": "Compounds added to the production line successfully"})
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"status": "500", "msg": "Failed to add compounds in production line"})
		}

	} else {

		utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to update record")
		slog.Error("%s - error - %s", "Failed to add compound data in production line", err.Error())
	}
}

func UpdateCompoundStatusToDispatch(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-Custom-Userid")
	var apiresp m.ApiRespMsg
	type compoundData struct {
		Notes    string
		KbDataId []string
		UserID   string
	}
	var data compoundData
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err == nil {
		data.UserID = userID

		jsonValue, err := json.Marshal(data)
		if err != nil {
			slog.Error("%s - error - %s", "Error marshaling compound data", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to marshal compound data")
			return
		}

		resp, err := http.Post(utils.RestURL+"/update-compound-status-to-dispatch", "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			w.WriteHeader(resp.StatusCode)
			body, _ := io.ReadAll(resp.Body) // Read response body
			apiresp.Code = resp.StatusCode
			apiresp.Message = string(body)
			http.Error(w, string(body), resp.StatusCode)

			return
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			body, _ := io.ReadAll(resp.Body) // Read response body
			w.WriteHeader(resp.StatusCode)
			w.Header().Set("Content-Type", "application/json")

			apiresp.Code = resp.StatusCode
			apiresp.Message = string(body)

			if err := json.NewEncoder(w).Encode(apiresp); err != nil {
				http.Error(w, "Failed to encode success message", http.StatusInternalServerError)
			}

		} else {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"status": "500", "msg": "Failed to add compounds in production line"})
		}

	} else {

		utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to update record")
		slog.Error("%s - error - %s", "Failed to add compound data in production line", err.Error())
	}
}

// Create new vendor
func AddNewVendor(w http.ResponseWriter, r *http.Request) {

	var data m.Vendors

	decoder := json.NewDecoder(r.Body)
	userID := r.Header.Get("X-Custom-Userid")
	if err := decoder.Decode(&data); err == nil {
		data.ModifiedBy = userID
		jsonValue, err := json.Marshal(data)
		if err != nil {
			slog.Error("%s - error - %s", "Error marshaling user data", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to marshal vendor data")
			return
		}

		resp, err := http.Post(utils.RestURL+"/create-new-vendor", "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			slog.Error("%s - error - %s", "Error making POST request", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to update vendor")
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
	} else {
		slog.Error("%s - error - %s", "Record updation failed", err.Error())
		utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to update record")
	}
}

func UpdateCompoundStatusToPacking(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-Custom-Userid")
	var apiresp m.ApiRespMsg
	type compoundData struct {
		Notes    string
		KbDataId []string
		UserID   string
	}
	var data compoundData
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err == nil {
		data.UserID = userID

		jsonValue, err := json.Marshal(data)
		if err != nil {
			slog.Error("%s - error - %s", "Error marshaling compound data", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to marshal compound data")
			return
		}

		resp, err := http.Post(utils.RestURL+"/update-compound-status-to-packing", "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			w.WriteHeader(resp.StatusCode)
			body, _ := io.ReadAll(resp.Body) // Read response body
			apiresp.Code = resp.StatusCode
			apiresp.Message = string(body)
			http.Error(w, string(body), resp.StatusCode)

			return
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			body, _ := io.ReadAll(resp.Body) // Read response body
			w.WriteHeader(resp.StatusCode)
			w.Header().Set("Content-Type", "application/json")

			apiresp.Code = resp.StatusCode
			apiresp.Message = string(body)

			if err := json.NewEncoder(w).Encode(apiresp); err != nil {
				http.Error(w, "Failed to encode success message", http.StatusInternalServerError)
			}

		} else {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"status": "500", "msg": "Failed to add compounds in production line"})
		}

	} else {

		utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to update record")
		slog.Error("%s - error - %s", "Failed to add compound data in production line", err.Error())
	}
}

func UpdateCompoundQualityStatusToReject(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-Custom-Userid")
	var apiresp m.ApiRespMsg
	type compoundData struct {
		Notes    string
		KbDataId []string
		UserID   string
	}
	var data compoundData
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err == nil {
		data.UserID = userID

		jsonValue, err := json.Marshal(data)
		if err != nil {
			slog.Error("%s - error - %s", "Error marshaling compound data", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to marshal compound data")
			return
		}

		resp, err := http.Post(utils.RestURL+"/update-compound-quality-status-to-reject", "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			w.WriteHeader(resp.StatusCode)
			body, _ := io.ReadAll(resp.Body) // Read response body
			apiresp.Code = resp.StatusCode
			apiresp.Message = string(body)
			http.Error(w, string(body), resp.StatusCode)

			return
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			body, _ := io.ReadAll(resp.Body) // Read response body
			w.WriteHeader(resp.StatusCode)
			w.Header().Set("Content-Type", "application/json")

			apiresp.Code = resp.StatusCode
			apiresp.Message = string(body)

			if err := json.NewEncoder(w).Encode(apiresp); err != nil {
				http.Error(w, "Failed to encode success message", http.StatusInternalServerError)
			}

		} else {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"status": "500", "msg": "Failed to add compounds in production line"})
		}

	} else {

		utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to update record")
		slog.Error("%s - error - %s", "Failed to add compound data in production line", err.Error())
	}
}

func DeleteKbRootByIDsHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-Custom-Userid")
	var apiresp m.ApiRespMsg
	type reqData struct {
		IDs    []string `json:"ids"`
		UserID string   `json:"UserID"`
	}
	var data reqData
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err == nil {
		data.UserID = userID

		jsonValue, err := json.Marshal(data)
		if err != nil {
			slog.Error("%s - error - %s", "Error marshaling compound data", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to marshal kanban data")
			return
		}

		resp, err := http.Post(utils.RestURL+"/delete-kbroot-by-ids", "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			w.WriteHeader(resp.StatusCode)
			body, _ := io.ReadAll(resp.Body) // Read response body
			apiresp.Code = resp.StatusCode
			apiresp.Message = string(body)
			http.Error(w, string(body), resp.StatusCode)

			return
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"status": "200", "msg": "Kanban deleted successfully"})
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"status": "500", "msg": "Failed to delete kanban"})
		}

	} else {

		utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to delete record")
		slog.Error("%s - error - %s", "Failed to delete kanban ", err.Error())
	}
}

// Get Kanban for vendor by id
func GetKanbanForVendor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	resp, err := http.Post(RestURL+"/get-kanban-for-vendor", "application/json", bytes.NewBuffer(body))
	if err != nil {
		slog.Error("%s - error - %s", "Error making POST request", err)
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

// Get Kanban for vendor by id
func GetKanbanForAllVendor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	resp, err := http.Get(RestURL + "/get-kanban-for-all-vendor")
	if err != nil {
		slog.Error("%s - error - %s", "Error making POST request", err)
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

func KanbanSortAsce(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Conditions []string `json:"Conditions"`
	}
	// Decode the request body
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	var vendor []*m.Vendors
	var vendortable []*m.VendorCompanyTable
	var compoundlist []m.Compounds
	var compound []m.CompoundsDataByVendor
	var prodline []m.ProdLine
	var action []s.DropDownOptions
	var root []m.KbRoot

	action = append(action, s.DropDownOptions{Text: "Select Action", Value: ""})

	kbRootResp, err := http.Get(RestURL + "/get-all-kb-root-data")
	if err != nil {
		slog.Error("Error fetching KbRoot", "error", err)
	}
	defer kbRootResp.Body.Close()

	if err := json.NewDecoder(kbRootResp.Body).Decode(&root); err != nil {
		slog.Error("Error decoding KbRoot", "error", err)
	}
	//fetching vendor records
	jsonReqData, err := json.Marshal(req)
	if err != nil {
		slog.Error("%s - error - %s", "Error marshaling user data", err)
		utils.SetResponse(w, http.StatusInternalServerError, "Failed to marshal data")
		return
	}
	vendorresp, err := http.Post(RestURL+"/get-all-vendors-data", "application/json", bytes.NewBuffer(jsonReqData))
	if err != nil {
		slog.Error("%s - error - %s", "Error making GET request", err)
	}

	defer vendorresp.Body.Close()

	if err := json.NewDecoder(vendorresp.Body).Decode(&vendor); err != nil {
		slog.Error("error decoding response body", "error", err)
	}

	//fetching prod line records
	prodlineresp, err := http.Get(RestURL + "/get-all-prod-line-data")
	if err != nil {
		slog.Error("%s - error - %s", "Error making GET request", err)
	}

	defer prodlineresp.Body.Close()

	responseData, err := io.ReadAll(prodlineresp.Body)
	if err != nil {
		slog.Error("Error reading response body", "error", err)
	}

	if err := json.Unmarshal(responseData, &prodline); err != nil {
		slog.Error("Error decoding response body", "error", err)
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

	for _, i := range vendor {
		var comp string
		// fetching component records by vendor
		coprresp, err := http.Post(RestURL+"/get-compound-data-by-vendor?key=id&value="+strconv.Itoa(i.ID), "application/json", bytes.NewBuffer(jsonReqData))
		if err != nil {
			slog.Error("%s - error - %s", "Error making GET request", err)
		}

		defer coprresp.Body.Close()

		if err := json.NewDecoder(coprresp.Body).Decode(&compound); err != nil {
			slog.Error("error decoding response body", "error", err)
		}

		// Sort by CreatedOn in ascending order
		sort.Slice(compound, func(i, j int) bool {
			return compound[i].CreatedOn.Before(compound[j].CreatedOn)
		})

		kbRootMap := make(map[int]string)
		for _, r := range root {
			kbRootMap[r.Id] = r.KanbanNo
		}

		for _, i := range compound {
			kanbanNo := kbRootMap[i.KbRootId]
			kanbanHTML := ""
			if strings.TrimSpace(kanbanNo) != "" {
				kanbanHTML = utils.JoinStr(
					`<span class="mx">|</span>
			<label class="form-check-label m-0 px-1" for="`, strconv.Itoa(i.KbRootId), `" data="Kanban-No">`,
					kanbanNo,
					`</label>`,
				)
			}
			opt := utils.JoinStr(`
					<div class="d-flex d-inline-flex  align-items-center pl-1 m-1" style="cursor: pointer !important;">
						<label class="form-check p-0 w-100" for="`, strconv.Itoa(i.KbRootId), `" style="cursor: pointer !important;">
						<div class="col-auto border border-1 p-0 mx-2" style="border-color: #ab71a2 !important; border-radius: 6px; user-select: none; cursor: pointer; background-color:`, utils.KanbanPriorityColors[i.CustomerNote]["bg-color"], `; color:`, utils.KanbanPriorityColors[i.CustomerNote]["text-color"], `;">
							<span class="form-check p-0 m-0" style="cursor: pointer;">
								<span class="border-end border-2 p-1 pl-1" style="border-color: #ab71a2 !important; cursor: pointer;">
									<input class="form-check-input m-1 mt-2 pl-1 component-code" type="checkbox" value="`, strconv.Itoa(i.KbRootId), `" id="`, strconv.Itoa(i.KbRootId), `">
								</span>
								<label class="form-check-label m-1 px-1" for="`, strconv.Itoa(i.KbRootId), `">
									`, i.CompoundName, ` 
								</label>
								<span class="mx">|</span>
								<label class="form-check-label m-0 px-1" for="`, strconv.Itoa(i.KbRootId), `" data="Cell-Name">
									`, i.CellNo, `
								</label>
								`, kanbanHTML, `
							
								<label class="form-check-label m-0 px-1 d-none" for="`, strconv.Itoa(i.KbRootId), `" data="Approved-Date">
									`, utils.FormatStringDate(i.DemandDate, "date"), `
								</label>
								
							</span>
						</div>
						</lable>
					</div>
			 `)

			comp += opt

		}
		var con string
		for _, v := range req.Conditions {
			con += con + v
		}

		var temptable m.VendorCompanyTable
		temptable.VendorCode = i.VendorCode
		temptable.VendorName = i.VendorName
		temptable.CompanyCodeAndNameString = comp
		if len(compound) > 0 {
			temptable.CreatedOn = compound[0].CreatedOn
		}
		if strings.Contains(con, "compound_name") {
			if len(compound) > 0 {
				vendortable = append(vendortable, &temptable)
			}
		} else {
			vendortable = append(vendortable, &temptable)
		}
	}

	// Now sort vendors by CreatedOn with vendors without CreatedOn moved to the end
	sort.Slice(vendortable, func(i, j int) bool {
		// If both have zero CreatedOn dates, they are equal
		if vendortable[i].CreatedOn.IsZero() && vendortable[j].CreatedOn.IsZero() {
			return false
		} else if vendortable[i].CreatedOn.IsZero() {
			// i has no CreatedOn, so it comes after j
			return false
		} else if vendortable[j].CreatedOn.IsZero() {
			// j has no CreatedOn, so i comes first
			return true
		}
		// Compare based on CreatedOn if both have valid dates
		return vendortable[i].CreatedOn.Before(vendortable[j].CreatedOn)
	})

	for _, i := range prodline {
		if i.Status {
			var tempopt s.DropDownOptions
			tempopt.Value = strconv.Itoa(i.Id)
			tempopt.Text = i.Name
			action = append(action, tempopt)
		}
	}

	var vendorcompanyTable s.TableCard
	tablebutton := `<button type="button" style="background-color:#871a83;color:white;" class="btn  m-0 p-2" id="del-License-btn" data-bs-toggle="modal" data-bs-target="#AddComp" > 
 						Add Comp
					</button>`
	allCheckBox := `<input style="width: 20px; height: 20px;" class="form-check-input me-2 select_all_compound" type="checkbox" value="" id="select_all_compound">`
	vendorcompanyTable.CardHeadingActions.Component = []s.CardHeadActionComponent{
		{ComponentName: "dropdown", ComponentType: s.ActionComponentElement{DropDown: s.DropdownAttributes{ID: "selectMenu-action-line", Name: "selectMenu-action-line", Options: action, Label: "Production Line", Disabled: true}}, Width: "col-9"},
		{ComponentName: "button", ComponentType: s.ActionComponentElement{Button: s.ButtonAttributes{ID: "deleteBtn", Disabled: true, Colour: "#c62f4a", Name: "deleteBtn", Type: "button", Text: "Delete"}}, Width: "col-3"},
	}
	vendorcompanyTable.CardHeadingActions.Style = "direction: ltr;"
	vendorcompanyTable.BodyTables = s.CardTableBody{Columns: []s.CardTableBodyHeadCol{
		{
			Name:  "Tools",
			Lable: " ",
		},
		{
			Name:         `Vendor Code`,
			IsSearchable: true,
			IsCheckbox:   true,
			ID:           "search-sn",
			Type:         "input",
			Width:        "col-1",
		},
		{
			Name:         "Vendor Name",
			IsSearchable: true,
			ID:           "search-MachinName",
			Type:         "input",
			Width:        "col-1",
		},
		{
			Name:  "Button",
			ID:    "search-MachinName",
			Width: "col-1",
		},
		{
			Name:             "Compound Code",
			IsSearchable:     true,
			SearchFieldWidth: "w-25",
			ID:               "search-MachinName",
			Type:             "input",
			Width:            "col9",
			Style:            "max-height: 40vh !important; overflow-y: auto;",
		},
		// {
		// 	Name:         "Vendor-Action",
		// 	IsSearchable: true,
		// 	ActionList:   action,
		// 	ID:           "search-MachinName",
		// 	Type:         "action",
		// 	Width:        "",
		// },
	},
		ColumnsWidth: []string{"", "col-1", "col-1", "col-1", "col- d-flex flex-wrap w-100 "},
		Data:         vendortable,
		Buttons:      tablebutton,
		Tools:        allCheckBox,
		ID:           "VendorComp",
	}

	tableBodyHTML := vendorcompanyTable.BodyTables.RenderBodyColumns()

	response := map[string]any{
		"tableBodyHTML": tableBodyHTML,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func KanbanSortDesc(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Conditions []string `json:"Conditions"`
	}
	// Decode the request body
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	var vendor []*m.Vendors
	var vendortable []*m.VendorCompanyTable
	var compoundlist []m.Compounds
	var compound []m.CompoundsDataByVendor
	var prodline []m.ProdLine
	var action []s.DropDownOptions
	var root []m.KbRoot

	action = append(action, s.DropDownOptions{Text: "Select Action", Value: ""})

	kbRootResp, err := http.Get(RestURL + "/get-all-kb-root-data")
	if err != nil {
		slog.Error("Error fetching KbRoot", "error", err)
	}
	defer kbRootResp.Body.Close()

	if err := json.NewDecoder(kbRootResp.Body).Decode(&root); err != nil {
		slog.Error("Error decoding KbRoot", "error", err)
	}
	//fetching vendor records
	jsonReqData, err := json.Marshal(req)
	if err != nil {
		slog.Error("%s - error - %s", "Error marshaling user data", err)
		utils.SetResponse(w, http.StatusInternalServerError, "Failed to marshal data")
		return
	}
	vendorresp, err := http.Post(RestURL+"/get-all-vendors-data", "application/json", bytes.NewBuffer(jsonReqData))
	if err != nil {
		slog.Error("%s - error - %s", "Error making GET request", err)
	}

	defer vendorresp.Body.Close()

	if err := json.NewDecoder(vendorresp.Body).Decode(&vendor); err != nil {
		slog.Error("error decoding response body", "error", err)
	}

	//fetching prod line records
	prodlineresp, err := http.Get(RestURL + "/get-all-prod-line-data")
	if err != nil {
		slog.Error("%s - error - %s", "Error making GET request", err)
	}

	defer prodlineresp.Body.Close()

	responseData, err := io.ReadAll(prodlineresp.Body)
	if err != nil {
		slog.Error("Error reading response body", "error", err)
	}

	if err := json.Unmarshal(responseData, &prodline); err != nil {
		slog.Error("Error decoding response body", "error", err)
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

	for _, i := range vendor {
		var comp string
		// fetching component records by vendor
		coprresp, err := http.Post(RestURL+"/get-compound-data-by-vendor?key=id&value="+strconv.Itoa(i.ID), "application/json", bytes.NewBuffer(jsonReqData))
		if err != nil {
			slog.Error("%s - error - %s", "Error making GET request", err)
		}

		defer coprresp.Body.Close()

		if err := json.NewDecoder(coprresp.Body).Decode(&compound); err != nil {
			slog.Error("error decoding response body", "error", err)
		}

		// Sort by CreatedOn in ascending order
		sort.Slice(compound, func(i, j int) bool {
			return compound[i].CreatedOn.After(compound[j].CreatedOn)
		})
		kbRootMap := make(map[int]string)
		for _, r := range root {
			kbRootMap[r.Id] = r.KanbanNo
		}

		// indianLocation := time.FixedZone("IST", 5*60*60+30*60)
		for _, i := range compound {
			kanbanNo := kbRootMap[i.KbRootId]
			kanbanHTML := ""
			if strings.TrimSpace(kanbanNo) != "" {
				kanbanHTML = utils.JoinStr(
					`<span class="mx">|</span>
			<label class="form-check-label m-0 px-1" for="`, strconv.Itoa(i.KbRootId), `" data="Kanban-No">`,
					kanbanNo,
					`</label>`,
				)
			}
			opt := utils.JoinStr(`
					<div class="d-flex d-inline-flex  align-items-center pl-1 m-1" style="cursor: pointer !important;">
						<label class="form-check p-0 w-100" for="`, strconv.Itoa(i.KbRootId), `" style="cursor: pointer !important;">
						<div class="col-auto border border-1 p-0 mx-2" style="border-color: #ab71a2 !important; border-radius: 6px; user-select: none; cursor: pointer; background-color:`, utils.KanbanPriorityColors[i.CustomerNote]["bg-color"], `; color:`, utils.KanbanPriorityColors[i.CustomerNote]["text-color"], `;">
							<span class="form-check p-0 m-0" style="cursor: pointer;">
								<span class="border-end border-2 p-1 pl-1" style="border-color: #ab71a2 !important; cursor: pointer;">
									<input class="form-check-input m-1 mt-2 pl-1 component-code" type="checkbox" value="`, strconv.Itoa(i.KbRootId), `" id="`, strconv.Itoa(i.KbRootId), `">
								</span>
								<label class="form-check-label m-1 px-1" for="`, strconv.Itoa(i.KbRootId), `">
									`, i.CompoundName, ` 
								</label>
								<span class="mx">|</span>
								<label class="form-check-label m-0 px-1" for="`, strconv.Itoa(i.KbRootId), `" data="Cell-Name">
									`, i.CellNo, `
								</label>
								`, kanbanHTML, `
							
								<label class="form-check-label m-0 px-1 d-none" for="`, strconv.Itoa(i.KbRootId), `" data="Approved-Date">
									`, utils.FormatStringDate(i.DemandDate, "date"), `
								</label>
								
							</span>
						</div>
						</lable>
					</div>
			 `)

			comp += opt

		}
		var con string
		for _, v := range req.Conditions {
			con += con + v
		}

		var temptable m.VendorCompanyTable
		temptable.VendorCode = i.VendorCode
		temptable.VendorName = i.VendorName
		temptable.CompanyCodeAndNameString = comp
		if len(compound) > 0 {
			temptable.CreatedOn = compound[0].CreatedOn
		}
		if strings.Contains(con, "compound_name") {
			if len(compound) > 0 {
				vendortable = append(vendortable, &temptable)
			}
		} else {
			vendortable = append(vendortable, &temptable)
		}
	}

	// Now sort vendors by CreatedOn with vendors without CreatedOn moved to the end
	sort.Slice(vendortable, func(i, j int) bool {
		// If both have zero CreatedOn dates, they are equal
		if vendortable[i].CreatedOn.IsZero() && vendortable[j].CreatedOn.IsZero() {
			return false
		} else if vendortable[i].CreatedOn.IsZero() {
			// i has no CreatedOn, so it comes after j
			return false
		} else if vendortable[j].CreatedOn.IsZero() {
			// j has no CreatedOn, so i comes first
			return true
		}
		// Compare based on CreatedOn if both have valid dates
		return vendortable[i].CreatedOn.After(vendortable[j].CreatedOn)
	})

	for _, i := range prodline {
		if i.Status {
			var tempopt s.DropDownOptions
			tempopt.Value = strconv.Itoa(i.Id)
			tempopt.Text = i.Name
			action = append(action, tempopt)
		}
	}

	var vendorcompanyTable s.TableCard
	tablebutton := `<button type="button" style="background-color:#871a83;color:white;" class="btn  m-0 p-2" id="del-License-btn" data-bs-toggle="modal" data-bs-target="#AddComp" > 
 						Add Comp
					</button>`
	allCheckBox := `<input style="width: 20px; height: 20px;" class="form-check-input me-2 select_all_compound" type="checkbox" value="" id="select_all_compound">`
	vendorcompanyTable.CardHeadingActions.Component = []s.CardHeadActionComponent{
		{ComponentName: "dropdown", ComponentType: s.ActionComponentElement{DropDown: s.DropdownAttributes{ID: "selectMenu-action-line", Name: "selectMenu-action-line", Options: action, Label: "Production Line", Disabled: true}}, Width: "col-9"},
		{ComponentName: "button", ComponentType: s.ActionComponentElement{Button: s.ButtonAttributes{ID: "deleteBtn", Disabled: true, Colour: "#c62f4a", Name: "deleteBtn", Type: "button", Text: "Delete"}}, Width: "col-3"},
	}
	vendorcompanyTable.CardHeadingActions.Style = "direction: ltr;"
	vendorcompanyTable.BodyTables = s.CardTableBody{Columns: []s.CardTableBodyHeadCol{
		{
			Name:  "Tools",
			Lable: " ",
		},
		{
			Name:         `Vendor Code`,
			IsSearchable: true,
			IsCheckbox:   true,
			ID:           "search-sn",
			Type:         "input",
			Width:        "col-1",
		},
		{
			Name:         "Vendor Name",
			IsSearchable: true,
			ID:           "search-MachinName",
			Type:         "input",
			Width:        "col-1",
		},
		{
			Name:  "Button",
			ID:    "search-MachinName",
			Width: "col-1",
		},
		{
			Name:             "Compound Code",
			IsSearchable:     true,
			SearchFieldWidth: "w-25",
			ID:               "search-MachinName",
			Type:             "input",
			Width:            "col9",
			Style:            "max-height: 40vh !important; overflow-y: auto;",
		},
		// {
		// 	Name:         "Vendor-Action",
		// 	IsSearchable: true,
		// 	ActionList:   action,
		// 	ID:           "search-MachinName",
		// 	Type:         "action",
		// 	Width:        "",
		// },
	},
		ColumnsWidth: []string{"", "col-1", "col-1", "col-1", "col- d-flex flex-wrap w-100 "},
		Data:         vendortable,
		Buttons:      tablebutton,
		Tools:        allCheckBox,
		ID:           "VendorComp",
	}

	tableBodyHTML := vendorcompanyTable.BodyTables.RenderBodyColumns()

	response := map[string]any{
		"tableBodyHTML": tableBodyHTML,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
