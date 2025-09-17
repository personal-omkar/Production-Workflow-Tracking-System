package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	m "irpl.com/kanban-commons/model"
	"irpl.com/kanban-commons/report"
	"irpl.com/kanban-commons/utils"
	"irpl.com/kanban-web/pages/basepage"
	"irpl.com/kanban-web/pages/orderentrypage"
	"irpl.com/kanban-web/services"
	s "irpl.com/kanban-web/services"
)

func orderentryPage(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("X-Custom-Username")
	userID := r.Header.Get("X-Custom-Userid")
	usertype := r.Header.Get("X-Custom-Role")
	links := r.Header.Get("X-Custom-Allowlist")
	var vendorName string
	var vendorCode string
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
			vendorCode = vendorRecord[0].VendorCode
		} else {
			vendorName = ""
			vendorCode = ""
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
				Selected: true,
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

	orderentryPage := orderentrypage.OrderEntrypage{
		Username:   username,
		UserID:     userID,
		UserType:   usertype,
		VendorCode: vendorCode,
	}

	// Build the complete BasePage with navigation and content
	page := (&basepage.BasePage{
		SideNavBar: sideNav,
		TopNavBar:  topNav,
		Username:   username,
		UserType:   usertype,
		Content:    orderentryPage.Build(),
	}).Build()

	// Write the complete HTML page to the response
	w.Write([]byte(page))
}

// Add order entry
func CreateNewOrderEntry(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-Custom-Userid")

	var data m.OrderEntry
	var ApiResp m.ApiRespMsg
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err == nil {
		data.UserID = userID
		data.MFGDateTime = time.Now()
		jsonValue, err := json.Marshal(data)
		if err != nil {
			slog.Error("%s - error - %s", "Error marshaling compound data", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to marshal order data")
			return
		}

		resp, err := http.Post(utils.RestURL+"/create-new-order-entry", "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			slog.Error("%s - error - %s", "Error making POST request", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to add marshal data in kanban")
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
	} else {
		ApiResp.Code = 500
		ApiResp.Message = "Fail: failed to create record"
		body, err := json.Marshal(ApiResp)
		if err != nil {
			log.Printf("Error marshaling response: %v", err)
		}
		utils.SetResponse(w, http.StatusInternalServerError, string(body))
		slog.Error("%s - error - %s", "Failed to add order data in table", err.Error())
	}
}

func CreateMultiNewOrderEntry(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-Custom-Userid")

	var data []m.OrderEntry
	var ApiResp m.ApiRespMsg
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err == nil {
		for index := range data {
			data[index].UserID = userID
			data[index].MFGDateTime = time.Now()
		}
		jsonValue, err := json.Marshal(data)
		if err != nil {
			slog.Error("%s - error - %s", "Error marshaling compound data", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to marshal order data")
			return
		}

		resp, err := http.Post(utils.RestURL+"/create-multi-new-order-entry", "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			slog.Error("%s - error - %s", "Error making POST request", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to add marshal data in kanban")
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

		utils.SetResponse(w, resp.StatusCode, string(body))
	} else {
		ApiResp.Code = 500
		ApiResp.Message = "Fail: failed to create record"
		body, err := json.Marshal(ApiResp)
		if err != nil {
			log.Printf("Error marshaling response: %v", err)
		}
		utils.SetResponse(w, http.StatusInternalServerError, string(body))
		slog.Error("%s - error - %s", "Failed to add order data in table", err.Error())
	}
}

func DeleteOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var data m.KbData
	var ApiResp m.ApiRespMsg
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&data); err == nil {
		jsonValue, err := json.Marshal(data)
		if err != nil {
			slog.Error("%s - error - %s", "Error marshaling compound data", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to marshal compound data")
			return
		}
		url := utils.RestURL + "/delete-order-entry"
		req, err := http.NewRequest(http.MethodDelete, url, bytes.NewBuffer(jsonValue))
		if err != nil {
			slog.Error("%s - error - %s", "Error creating Delete request", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to create Delete request")
			return
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			slog.Error("%s - error - %s", "Error Delete order", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to delete order")
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
	} else {
		ApiResp.Code = 500
		ApiResp.Message = "Fail: Failed to delete order"
		body, err := json.Marshal(ApiResp)
		if err != nil {
			log.Printf("Error marshaling response: %v", err)
		}
		utils.SetResponse(w, http.StatusInternalServerError, string(body))
	}
}

type VendorLotLimitResponse struct {
	DailyLimit   int    `json:"daily_limit"`
	Message      string `json:"message"`
	MonthlyLimit int    `json:"monthly_limit"`
	HourlyLimit  int    `json:"hourly_limit"`
	VendorID     int    `json:"vendor_id"`
	Dialog       string `json:"dialog"`
	StatusCode   int    `json:"status"`
	ExcedBy      int    `json:"exceed_by"`
}

func DailyAndMonthlyVendorLimit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID := r.Header.Get("X-Custom-Userid")
	var data []m.KbData
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&data); err != nil {
		slog.Error("%s - error - %s", "Failed to decode request body", err)
		utils.SetResponse(w, http.StatusBadRequest, "Invalid input data")
		return
	}
	for i := range data {
		data[i].Id, _ = strconv.Atoi(userID)
	}
	jsonValue, err := json.Marshal(data)
	if err != nil {
		slog.Error("%s - error - %s", "Error marshaling compound data", err)
		utils.SetResponse(w, http.StatusInternalServerError, "Failed to marshal compound data")
		return
	}

	url := utils.RestURL + "/check-vendor-lot-limit"
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonValue))
	if err != nil {
		slog.Error("%s - error - %s", "Error creating POST request", err)
		utils.SetResponse(w, http.StatusInternalServerError, "Failed to create POST request")
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("%s - error - %s", "Error sending POST request", err)
		utils.SetResponse(w, http.StatusInternalServerError, "Failed to send request")
		return
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("%s - error - %s", "Error reading response body", err)
		utils.SetResponse(w, http.StatusInternalServerError, "Failed to read response body")
		return
	}

	var lotLimitResponse VendorLotLimitResponse
	err = json.Unmarshal(responseBody, &lotLimitResponse)
	if err != nil {
		slog.Error("%s - error - %s", "Error unmarshaling response JSON", err)
		utils.SetResponse(w, http.StatusInternalServerError, "Failed to parse response data")
		return
	}

	if resp.StatusCode != http.StatusOK {
		confirmSubmit := s.ConfirmationModal{
			For:   "Delete",
			ID:    "LotSizeExceed",
			Title: "Alert!",
			Body: []string{
				lotLimitResponse.Message,
				"<b>Please note: The Hourly lot limit is " + strconv.Itoa(lotLimitResponse.HourlyLimit) + ", the daily lot limit is " + strconv.Itoa(lotLimitResponse.DailyLimit) + " and the monthly lot limit is " + strconv.Itoa(lotLimitResponse.MonthlyLimit) + ".</b>",
				"<b>Exceed By : " + strconv.Itoa(lotLimitResponse.ExcedBy) + "<b>",
			},
			Footer: s.Footer{CancelBtn: false, Buttons: []s.FooterButtons{{BtnType: "submit", BtnID: "Close_Modal", Text: "Ok", Style: "background-color:#636E7E; border:none;"}}},
		}

		Dialog := confirmSubmit.Build()
		lotLimitResponse.Dialog = Dialog
		lotLimitResponse.StatusCode = http.StatusInternalServerError
		jsonResponse, err := json.Marshal(lotLimitResponse)
		if err != nil {
			slog.Error("%s - error - %s", "Error marshaling response data", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to marshal response data")
			return
		}
		utils.SetResponse(w, resp.StatusCode, string(jsonResponse))
	} else {
		lotLimitResponse.StatusCode = http.StatusOK
		jsonResponse, err := json.Marshal(lotLimitResponse)
		if err != nil {
			slog.Error("%s - error - %s", "Error marshaling response data", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to marshal response data")
			return
		}
		utils.SetResponse(w, http.StatusOK, string(jsonResponse))
	}
}

func DailyAndMonthlyVendorLimitByVendorCode(w http.ResponseWriter, r *http.Request) {
	var data m.OrderDetails
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&data); err == nil {
		if data.DemandDateTime.IsZero() {
			data.DemandDateTime = time.Now()
		}

		jsonValue, err := json.Marshal(data)
		if err != nil {
			slog.Error("%s - error - %s", "Error marshaling compound data", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to marshal compound data")
			return
		}

		url := utils.RestURL + "/check-vendor-lot-limit-by-vendor-code"
		res, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonValue))
		if err != nil {
			log.Printf("Failed to delete order: %v", err)
			http.Error(w, "Failed to delete order", http.StatusForbidden)
			return
		}
		defer res.Body.Close()

		client := &http.Client{}
		resp, err := client.Do(res)
		if err != nil {
			http.Error(w, "Failed to read response body", http.StatusInternalServerError)
			return
		}

		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			slog.Error("%s - error - %s", "Error reading response body", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to read response body")
			return
		}

		utils.SetResponse(w, resp.StatusCode, string(responseBody))
	} else {
		log.Printf("Error %v", err)
	}
}

func GetCustomerOrderDetails(w http.ResponseWriter, r *http.Request) {
	Id := r.URL.Query().Get("id")
	status := r.URL.Query().Get("status")

	type TableConditions struct {
		Conditions []string `json:"Conditions"`
	}

	var tablecondition TableConditions
	con := utils.JoinStr(`kb_data.created_by='`, Id, `' AND kb_extension.status='`, status, `'`)

	if status == "submit" {
		con = utils.JoinStr(`kb_data.created_by='`, Id, `' AND kb_extension.status!='creating'`)
	}

	tablecondition.Conditions = append(tablecondition.Conditions, con)

	jsonValue, err := json.Marshal(tablecondition)
	if err != nil {
		slog.Error("%s - error - %s", "Error marshaling table condition", err)
	}

	resp, err := http.Post(RestURL+"/get-customer-order-details", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		slog.Error("%s - error - %s", "Error making GET request", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("%s - error - %s", "Error reading response body", err)
		utils.SetResponse(w, http.StatusInternalServerError, "Failed to read response body")
		return
	}

	utils.SetResponse(w, http.StatusOK, string(responseBody))
}

func PrintHTMLCard(w http.ResponseWriter, r *http.Request) {

	cellno := r.URL.Query().Get("cellno")

	var rawQuery m.RawQuery
	// KBdata details
	rawQuery.Type = "KbData"
	rawQuery.Host = utils.RestHost
	rawQuery.Port = utils.RestPort
	rawQuery.Query = utils.JoinStr(`SELECT * FROM kb_data WHERE cell_no = '`, cellno, `'`) // `;`
	var kbData []m.KbData
	rawQuery.RawQry(&kbData)

	// Compound Details
	rawQuery.Type = "Compounds"
	rawQuery.Host = utils.RestHost
	rawQuery.Port = utils.RestPort
	rawQuery.Query = utils.JoinStr(`SELECT * FROM Compounds WHERE id = `, fmt.Sprint(kbData[0].CompoundId)) // `;`
	var Comps []m.Compounds
	rawQuery.RawQry(&Comps)

	// Vendors Details
	rawQuery.Type = "Vendors"
	rawQuery.Host = utils.RestHost
	rawQuery.Port = utils.RestPort
	rawQuery.Query = utils.JoinStr(`SELECT * FROM Vendors WHERE id in (select (vendor_id) from kb_extension where id = `, fmt.Sprint(kbData[0].KbExtensionID), `)`) // `;`
	var Vendors []*m.Vendors
	rawQuery.RawQry(&Vendors)

	var multiReport string

	for _, v := range kbData[0].KanbanNo {

		multiReport = utils.JoinStr(multiReport, `
  		<div style="page-break-after : always;  width: 374px;  margin-top: 7px;">
			<div>
				<table class="printTable" cellpadding="0" cellspacing="0" style="border-collapse: collapse;  height: 50px;">
					<tbody>
					<tr style="background: #32cefd;">
						<td width="374px" style="text-align: center;">
							<p style="font-size: 26px; font-weight: bold;">Rubber Kanban Card</p>
						</td>
					</tr>
				</tbody>
				</table>
			</div>
			<div style="margin-top: -1px;">
				<table class="printTable" cellpadding="0" cellspacing="0" style="border-collapse: collapse;  height: 170px;">
					<tbody>
						<tr>
							<td style="background: #32cefd; width: 140px;">
							<p style="text-align: start;font-size: 13px; height: 20px;">CELL NO / VENDOR:-</p>
							</td>
							<td style="text-align: center;width: 234px;"> <span style="height: 20px; display: block; font-weight: bold;"> `, Vendors[0].VendorName, ` </span> </td>
						</tr>
						<tr>
							<td style="background: #32cefd; width: 140px;">
							<p style="text-align: start;font-size: 13px; height: 20px;">KANBAN NO:-</p>
							</td>
							<td style="text-align: center;width: 234px;"> <span style="height: 20px; display: block; font-weight: bold;"> `, v, ` </span> </td>
						</tr>
						<tr>
							<td style="background: #32cefd; width: 140px;">
							<p style="text-align: start;font-size: 13px; height: 20px;">COMPOUND CODE</p>
							</td>
							<td style="text-align: center;width: 234px;"> <span style="height: 20px; display: block; font-weight: bold;font-size: 0.6rem;">`, Comps[0].CompoundName, `</span> </td>
						</tr>
						<tr>
							<td style="background: #32cefd; width: 140px;">
							<p style="text-align: start;font-size: 13px; height: 20px;">DEMAND DATE/TIME</p>
							</td>
							<td style="text-align: center; width: 234px;"> <span style="height: 20px; display: block; font-weight: bold;">`, kbData[0].DemandDateTime.Format("02.01.2006 15.04.05"), `</span> </td>
						</tr>
						<tr>
							<td style="background: #32cefd; width: 140px;">
							<p style="text-align: start;font-size: 13px; height: 20px;">MFG. DATE/TIME</p>
							</td>
							<td style="text-align: center; width: 234px;"> <span style="height: 20px; display: block; font-weight: bold;"></span> </td>
						</tr>
						<tr>
							<td style="background: #32cefd; width: 140px;">
							<p style="text-align: start;font-size: 13px; height: 20px;">EXP. DATE/TIME</p>
							</td>
							<td style="text-align: center; width: 234px;"> <span style="height: 20px; display: block; font-weight: bold;"></span> </td>
						</tr>
						<tr>
							<td style="background: #32cefd; width: 140px;">
							<p style="text-align: start;font-size: 13px; height: 20px;">LOCATION</p>
							</td>
							<td style="text-align: center; width: 234px;"> <span style="height: 20px; display: block; font-weight: bold;">`, kbData[0].Location, `</span> </td>
						</tr>
					</tbody>
				</table>
			</div>
		</div>
	
		`)
	}

	d := report.HtmlCardDialog{
		RecordId:   "demand-card",
		PageSource: multiReport,
	}

	utils.SetResponse(w, http.StatusOK, d.Build())
}

// GetCompoundListByParam get parts list by parameter
func GetCompoundsListByParam(w http.ResponseWriter, r *http.Request) {
	var Parts []*m.Compounds
	pn := r.URL.Query().Get("partName")

	var rawQuery m.RawQuery
	rawQuery.Host = utils.RestHost
	rawQuery.Port = utils.RestPort
	rawQuery.Type = "Compounds"
	rawQuery.Query = `SELECT * FROM Compounds WHERE (Compound_Name iLIKE '%` + pn + `%') AND status IS TRUE` //`;`
	rawQuery.RawQry(&Parts)

	var val = []struct {
		Id   string
		Text string
	}{}
	for _, pm := range Parts {
		var ret = struct {
			Id   string
			Text string
		}{}
		ret.Id = fmt.Sprint(pm.Id)
		ret.Text = pm.CompoundName
		val = append(val, ret)
	}

	bData, _ := json.Marshal(val)
	utils.SetResponse(w, http.StatusOK, string(bData))
}

func DailyVendorLimit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID := r.Header.Get("X-Custom-Userid")
	var data m.KbData
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&data); err != nil {
		slog.Error("%s - error - %s", "Failed to decode request body", err)
		utils.SetResponse(w, http.StatusBadRequest, "Invalid input data")
		return
	}
	data.Id, _ = strconv.Atoi(userID)
	jsonValue, err := json.Marshal(data)
	if err != nil {
		slog.Error("%s - error - %s", "Error marshaling compound data", err)
		utils.SetResponse(w, http.StatusInternalServerError, "Failed to marshal compound data")
		return
	}

	url := utils.RestURL + "/check-daily-lot-limit"
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonValue))
	if err != nil {
		slog.Error("%s - error - %s", "Error creating POST request", err)
		utils.SetResponse(w, http.StatusInternalServerError, "Failed to create POST request")
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("%s - error - %s", "Error sending POST request", err)
		utils.SetResponse(w, http.StatusInternalServerError, "Failed to send request")
		return
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("%s - error - %s", "Error reading response body", err)
		utils.SetResponse(w, http.StatusInternalServerError, "Failed to read response body")
		return
	}

	var lotLimitResponse VendorLotLimitResponse
	err = json.Unmarshal(responseBody, &lotLimitResponse)
	if err != nil {
		slog.Error("%s - error - %s", "Error unmarshaling response JSON", err)
		utils.SetResponse(w, http.StatusInternalServerError, "Failed to parse response data")
		return
	}

	if resp.StatusCode != http.StatusOK {
		confirmSubmit := services.ConfirmationModal{
			For:   "Delete",
			ID:    "LotSizeExceed",
			Title: "Alert!",
			Body: []string{
				lotLimitResponse.Message,
				"<b>Please note:  The Daily lot limit is " + strconv.Itoa(lotLimitResponse.DailyLimit) +
					"<b>Exceed By : " + strconv.Itoa(lotLimitResponse.ExcedBy) + "<b>",
			},
			Footer: services.Footer{CancelBtn: false, Buttons: []services.FooterButtons{{BtnType: "submit", BtnID: "Close_Modal", Text: "Ok", Style: "background-color:#636E7E; border:none;"}}},
		}

		Dialog := confirmSubmit.Build()
		lotLimitResponse.Dialog = Dialog
		lotLimitResponse.StatusCode = http.StatusInternalServerError
		jsonResponse, err := json.Marshal(lotLimitResponse)
		if err != nil {
			slog.Error("%s - error - %s", "Error marshaling response data", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to marshal response data")
			return
		}
		utils.SetResponse(w, resp.StatusCode, string(jsonResponse))
	} else {
		lotLimitResponse.StatusCode = http.StatusOK
		jsonResponse, err := json.Marshal(lotLimitResponse)
		if err != nil {
			slog.Error("%s - error - %s", "Error marshaling response data", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to marshal response data")
			return
		}
		utils.SetResponse(w, http.StatusOK, string(jsonResponse))
	}
}
