package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
	m "irpl.com/kanban-commons/model"
	"irpl.com/kanban-commons/utils"
	websecure "irpl.com/kanban-commons/websecure"
	"irpl.com/kanban-web/pages/basepage"
	"irpl.com/kanban-web/pages/dashboard"
	kanbanhistorypage "irpl.com/kanban-web/pages/kanbanhistorypage"
	"irpl.com/kanban-web/pages/login"
	orderhistorypage "irpl.com/kanban-web/pages/orderhistorypage"
	"irpl.com/kanban-web/pages/orders"
	productionline "irpl.com/kanban-web/pages/productionLine"
	"irpl.com/kanban-web/pages/productionprocesses"
	"irpl.com/kanban-web/pages/register"
	"irpl.com/kanban-web/services"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	sideNav := basepage.SideNav{}
	sideNav.MenuItems = []basepage.SideNavItem{
		{
			Name:     "Dashboard",
			Icon:     "fas fa-chart-pie",
			Link:     "#",
			Style:    "font-size:1rem;",
			Selected: true,
		},
		{
			Name:  "Orders",
			Icon:  "fas fa-list-alt",
			Link:  "/vendor-orders",
			Style: "font-size:1rem;",
		},
		{
			Name:  "Orders",
			Icon:  "fas fa-list-alt",
			Link:  "/admin-orders",
			Style: "font-size:1rem;",
		},
		{
			Name:  "Heijunka Board",
			Icon:  "fas fa-calendar-alt",
			Link:  "/production-line",
			Style: "font-size:1rem;",
		},
		{
			Name:  "kanban",
			Icon:  "fas fa-calendar-day",
			Style: "font-size:1rem;",
			Items: []basepage.SideNavSubItem{
				{
					Name:  "Item1",
					URL:   "#",
					Style: "font-size:1rem;",
				},
				{
					Name: "Item2",
					URL:  "font-size:1rem;",
				},
			},
		},
		{
			Name:     "Comp Management",
			Icon:     "fas fa-vials",
			Link:     "/compounds-management",
			UserType: basepage.UserType{Admin: true, Operator: true},
			Style:    "font-size:1rem;",
		},
	}

	topNav := basepage.TopNav{}
	topNav.MenuItems = []basepage.TopNavItem{
		{
			ID:    "settings",
			Name:  "",
			Type:  "link",
			Icon:  "fa fa-cog",
			Link:  "#",
			Width: "col-1",
			Style: "bg-light bg-gradient",
		},
		{
			ID:    "notifications",
			Name:  "",
			Type:  "link",
			Icon:  "fa fa-bell",
			Link:  "#",
			Width: "col-1",
			Style: "bg-light bg-gradient",
		},
		{
			ID:    "username",
			Name:  "User Name",
			Type:  "button",
			Icon:  "",
			Width: "col-2",
			Style: "bg-light bg-gradient",
		},
		{
			ID:    "logout",
			Name:  "",
			Link:  "/logout",
			Type:  "link",
			Icon:  "fas fa-sign-out-alt",
			Width: "col-1",
			Style: "bg-light bg-gradient",
		},
	}

	page := (&basepage.BasePage{
		SideNavBar: sideNav,
		TopNavBar:  topNav,
	}).Build()
	w.Write([]byte(page))
}

func basePage(w http.ResponseWriter, r *http.Request) {

	http.Redirect(w, r, "/dashboard", http.StatusFound)
}

func loginPage(w http.ResponseWriter, r *http.Request) {

	var err = r.URL.Query().Get("error")
	page := (&login.LoginPage{
		Error: err,
	}).Build()
	w.Write([]byte(page))
}

func registerPage(w http.ResponseWriter, r *http.Request) {
	var err = r.URL.Query().Get("error")

	page := (&register.RegisterPage{
		Error: err,
	}).Build()
	w.Write([]byte(page))
}

// Status ss
func Status(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

// Web Configuration web interface
func Web() {

	r := mux.NewRouter()

	// Middleware to apply to all routes
	r.Use(websecure.CommonMiddleware)
	r.Use(websecure.CorsMiddleware)

	// Routes consist of a path and a handler function.
	r.HandleFunc("/", basePage)
	r.HandleFunc("/status", Status)
	r.HandleFunc("/home", homePage)
	r.HandleFunc("/dashboard", dashboard.RegisterPage)
	r.HandleFunc("/flowchart", flowChartPage)
	r.HandleFunc("/vendor-company", vendorCompanyPage)
	r.HandleFunc("/cold-storage", ColdStoragePage)
	r.HandleFunc("/vendor-orders", vendorOrderPage)
	r.HandleFunc("/admin-orders", adminOrderPage)
	r.HandleFunc("/vendor-management", vendorManagmentPage)
	r.HandleFunc("/kanban-history", kanabnHistoryPage)
	r.HandleFunc("/order-history", orderHistoryPage)
	r.HandleFunc("/user-management", usermanagmentPage)
	r.HandleFunc("/user-management-card", usermanagmentDialog).Methods("GET")
	r.HandleFunc("/user-role-management", userRoleManagmentPage)
	r.HandleFunc("/system-logs", LogsManagement)
	r.HandleFunc("/configuration-page", configurationPage)
	r.HandleFunc("/packing-dispatch-page", packingDispatchPage)
	r.HandleFunc("/recipe-management", recipeManagementPage)
	r.HandleFunc("/stage-management", stageManagementPage)
	r.HandleFunc("/login", loginPage).Methods("GET")   // Serves the login page
	r.HandleFunc("/do-login", doLogin).Methods("POST") // Handles login API call
	r.HandleFunc("/sign-in", SigneIn).Methods("POST")
	r.HandleFunc("/register", registerPage)
	r.HandleFunc("/do-register", doRegisterHandler)
	r.HandleFunc("/sign-up", SigneUp).Methods("POST")
	r.HandleFunc("/logout", LogoutHandler)
	r.HandleFunc("/get-user-by-email", GetUserByEmailHandler).Methods("GET")
	r.HandleFunc("/get-user-role-by-id", GetUserRoleByRoleIDHandler).Methods("GET")
	r.HandleFunc("/production-line", productionline.RenderProductionLinePage)
	r.HandleFunc("/get-production-line-items", GetProductionLinesHandler).Methods("GET")
	r.HandleFunc("/get-production-line-processes", GetProductionLineProcesses).Methods("POST")
	r.HandleFunc("/add-production-line", AddProdLine).Methods("POST")
	r.HandleFunc("/update-order-status", UpdateOrderStatus).Methods("POST")
	r.HandleFunc("/delete-production-line-cell", productionline.DeleteProductionLineCell).Methods("DELETE")
	r.HandleFunc("/update-coldstorage-quantity", UpdateColdStorageQuantity).Methods("POST")
	r.HandleFunc("/get-user-details", GetUserDetails).Methods("GET")
	r.HandleFunc("/user-search-pagination", UserSearchPagination).Methods("POST")
	// -- API: Compounds --
	// -- POST --
	r.HandleFunc("/add-compound-data-by-vendor", AddCompoundsForVendor).Methods("POST")
	r.HandleFunc("/add-compound-data-to-production-line", AddCompoundsInProductionLine).Methods("POST")
	r.HandleFunc("/update-compound-status-to-dispatch", UpdateCompoundStatusToDispatch).Methods("POST")
	r.HandleFunc("/update-compound-status-to-packing", UpdateCompoundStatusToPacking).Methods("POST")
	r.HandleFunc("/update-compound-quality-status-to-reject", UpdateCompoundQualityStatusToReject).Methods("POST")
	r.HandleFunc("/compounds-search-pagination", CompoundsSearchPagination).Methods("POST")

	r.HandleFunc("/get-customer-order-details", GetCustomerOrderDetails).Methods("GET")

	r.HandleFunc("/get-in-progress-count-by-line", GetInProgressCountByLine).Methods("POST")
	r.HandleFunc("/get-vendor-order-status-this-month", GetVendorOrderStatusThisMonth).Methods("POST")

	// -- API: Order --
	// -- POST --
	r.HandleFunc("/create-new-order-entry", CreateNewOrderEntry).Methods("POST")
	r.HandleFunc("/create-multi-new-order-entry", CreateMultiNewOrderEntry).Methods("POST")

	r.HandleFunc("/order-entry", orderentryPage)
	r.HandleFunc("/get-order-details", GetOrderDEtails).Methods("POST")
	r.HandleFunc("/search-order-details", SearchOrderDEtails).Methods("POST")
	r.HandleFunc("/get-all-kbRoot-details-by-search", GetCompletedKBRootDetailsBySearch).Methods("POST")
	r.HandleFunc("/search-customer-order-details", SearchCustomerOrderDetails).Methods("POST")

	// -- GET --
	r.HandleFunc("/cust/print-html-card", PrintHTMLCard).Methods("GET")
	r.HandleFunc("/get-compounds-list-by-parem", GetCompoundsListByParam).Methods("GET")

	// -- API: Vendors --
	// -- POST --
	r.HandleFunc("/create-new-vendor", AddNewVendor).Methods("POST")
	r.HandleFunc("/get-vendor-by-userid", GetVendorByUserID)
	r.HandleFunc("/create-new-KbTransaction", CreateNewKbTransaction).Methods("POST")
	r.HandleFunc("/get-all-orders-by-vendor-code", GetAllOrderByVendorCode).Methods("POST")
	r.HandleFunc("/get-vendor-details-by-vendor-code", GetVendorDetailsByVendorCode).Methods("POST")
	r.HandleFunc("/import-vendor-data", ImportVendorMasterData)
	r.HandleFunc("/vendor-search-pagination", VendorSearchPagination).Methods("POST")

	r.HandleFunc("/update-running-number", UpdateRunningNumberAfterTransactioin).Methods("POST")

	r.HandleFunc("/OrderDetailsForCustomer", OrderDetailsForCustomerHandler).Methods("POST")
	r.HandleFunc("/update-running-numbers", UpdateRunningNumbers).Methods("POST")
	r.HandleFunc("/update-flowchart", UpdateFlowchart).Methods("GET")
	r.HandleFunc("/delete-order-entry", DeleteOrder).Methods("DELETE")

	r.HandleFunc("/update-user-details", UpdateUserDetails).Methods("POST")
	r.HandleFunc("/create-new-user", CreateUser).Methods("POST")
	r.HandleFunc("/update-user-role", UpdateUserRole).Methods("POST")
	r.HandleFunc("/delete-user", DeleteUser).Methods("POST")
	r.HandleFunc("/create-user-role", CreateUserRole).Methods("POST")
	r.HandleFunc("/delete-user-role", DeleteUserRole).Methods("POST")
	r.HandleFunc("/build-comp-details-dialog", kanbanhistorypage.BuildCompDetailsDialog).Methods("POST")
	r.HandleFunc("/build-order-details-dialog", orderhistorypage.BuildOrderDetailsDialog).Methods("POST")
	r.HandleFunc("/user-role-search-pagination", UserRoleSearchPagination).Methods("POST")

	r.HandleFunc("/prod-line-management", productionline.ProdLineCRUDPage)
	r.HandleFunc("/edit-prod-line-dialog", productionline.EditProdLineDialog).Methods("POST")
	r.HandleFunc("/edit-prod-line", productionline.EditProdLine).Methods("POST")

	r.HandleFunc("/prod-processes-management", productionprocesses.ProdProcessesCRUDPage)
	r.HandleFunc("/add-production-process", productionprocesses.AddProdProcess).Methods("POST")
	r.HandleFunc("/edit-prod-process-dialog", services.EditProdProcessDialog).Methods("POST")
	r.HandleFunc("/edit-prod-process", productionprocesses.EditProdProcess).Methods("POST")

	// -- API: LDAP Configuration --
	// -- POST --
	r.HandleFunc("/secure/create-ldap-config", CreateLDAPConfiguration).Methods("POST")
	r.HandleFunc("/secure/update-ldap-config", UpdateLDAPConfiguration).Methods("POST")
	// -- GET --
	r.HandleFunc("/secure/get-ldap-config", GetDefaultLDAPConfig).Methods("GET")

	// -- API: SAMBA Configuration --
	// -- POST --
	r.HandleFunc("/secure/create-samba-config", CreateSAMBAConfiguration).Methods("POST")
	r.HandleFunc("/secure/update-samba-config", UpdateSAMBAConfiguration).Methods("POST")
	// -- GET --
	r.HandleFunc("/secure/get-samba-config", GetDefaultSAMBAConfig).Methods("GET")

	r.HandleFunc("/check-vendor-lot-limit", DailyAndMonthlyVendorLimit).Methods("POST")
	r.HandleFunc("/check-vendor-lot-limit-by-vendor-code", DailyAndMonthlyVendorLimitByVendorCode).Methods("POST")
	r.HandleFunc("/compounds-management", CompoundsManagement)
	r.HandleFunc("/add-update-compound", AddorUpdateCompound).Methods("POST")
	r.HandleFunc("/check-daily-lot-limit", DailyVendorLimit).Methods("POST")

	r.HandleFunc("/get-all-details-for-order", GetOrderDetailsForHistory).Methods("POST")
	r.HandleFunc("/get-user-details-by-email", GetUserDetailsByEmail).Methods("GET")

	r.HandleFunc("/sort-admin-order-table", orders.SortAdminOrderTable).Methods("POST")
	r.HandleFunc("/quality-testing", QualityTestingPage)
	r.HandleFunc("/delete-kbroot-by-ids", DeleteKbRootByIDsHandler).Methods("POST")

	r.HandleFunc("/kanban-report", kanabnReportPage)
	r.HandleFunc("/kanban-report-search-pagination", kanbanReportSearchPagination).Methods("POST")
	r.HandleFunc("/get-all-active-compounds", GetAllActiveCompounds).Methods("GET")
	r.HandleFunc("/import-compound-data", ImportCompoundData)

	// -- API:- Stage
	r.HandleFunc("/create-new-stage", CreateStage).Methods("POST")
	r.HandleFunc("/update-stage", UpdateExistingStage).Methods("POST")
	r.HandleFunc("/delete-stage", DeleteStageByID).Methods("POST")
	r.HandleFunc("/get-stages-by-header", GetStagesByHeader).Methods("POST")
	r.HandleFunc("/get-stage-by-param", GetStagesByParam).Methods("GET")
	r.HandleFunc("/get-all-stages", GetAllStage).Methods("GET")
	r.HandleFunc("/edit-stage-dialog", StageMangementEditDialog).Methods("GET")

	// -- API: Recipe
	// -- POST --
	r.HandleFunc("/create-new-recipe", CreateRecipe).Methods("POST")
	r.HandleFunc("/update-existing-recipe", UpdateRecipe).Methods("POST")
	r.HandleFunc("/delete-recipe-by-id", DeleteRecipe).Methods("POST")
	r.HandleFunc("/edit-recipe-dialog", EditRecipeDialog).Methods("POST")
	// -- GET --
	r.HandleFunc("/get-all-recipe", GetAllRecipe).Methods("GET")
	r.HandleFunc("/get-recipe-by-data-key", GetRecipeByDataKey).Methods("GET")
	r.HandleFunc("/get-recipe-by-data-value", GetRecipeByDataValue).Methods("GET")
	r.HandleFunc("/get-recipe-by-data-key-and-value", GetRecipeByDataKeyAndValue).Methods("GET")

	// Kanban-Sorting
	r.HandleFunc("/sort-kanban-asce", KanbanSortAsce).Methods("POST")
	r.HandleFunc("/sort-kanban-desc", KanbanSortDesc).Methods("POST")

	//All Kanban View
	r.HandleFunc("/all-kanban-view", AllKanbanView)
	r.HandleFunc("/search-kanban-view", SearchAllKanbanView).Methods("POST")

	// Quality test-search and sort
	r.HandleFunc("/quality-sort-kanban-desc", QualityKanbanSortDesc).Methods("POST")
	r.HandleFunc("/quality-sort-kanban-asce", QualityKanbanSortAsce).Methods("POST")

	// Packing-Kanban-Sorting-and-Searching
	r.HandleFunc("/packing-sort-kanban-asce", PackingKanbanSortAsce).Methods("POST")
	r.HandleFunc("/packing-sort-kanban-desc", PackingKanbanSortDesc).Methods("POST")

	// Kanban Reprint
	r.HandleFunc("/kanban-reprint", KanbanReprintPage)
	r.HandleFunc("/summary-reprint", SummaryReprintPage)
	r.HandleFunc("/print-kanban-summary", PrintKanbanSummary).Methods("POST")

	// Inventory Managmenet
	r.HandleFunc("/inventory-management", inventoryManagement).Methods("GET")
	r.HandleFunc("/create-new-or-update-existing-inventory", CreateOrUpdateInventory).Methods("POST")
	r.HandleFunc("/delete-inventory-by-id", DeleteInventoryById).Methods("POST")
	r.HandleFunc("/coldstorage-pagination-search", ColdStorageSearchPagination).Methods("POST")
	r.HandleFunc("/get-compound-data-by-parm", GetCompoundsByParm).Methods("GET")
	r.HandleFunc("/update-coldstorage-quantity-for-inventory-managment", UpdateColdStorageQuantityForInventoryManagement).Methods("POST")

	//--Operator--
	r.HandleFunc("/operator-management-card", operatormanagmentDialog).Methods("GET")
	r.HandleFunc("/operator-management", OperatorManagmentPage)
	r.HandleFunc("/create-new-or-update-existing-operator", CreateNewOrUpdateExistingOperator).Methods("POST")
	r.HandleFunc("/import-operator-data", ImportOperatorMasterData)
	r.HandleFunc("/operator-search-pagination", OperatorSearchPagination).Methods("POST")

	//--Raw Material--
	r.HandleFunc("/material-management-card", materialmanagmentDialog).Methods("GET")
	r.HandleFunc("/material-management", MaterialManagmentPage)
	r.HandleFunc("/create-new-or-update-existing-material", CreateNewOrUpdateExistingRawMaterial).Methods("POST")
	r.HandleFunc("/import-material-data", ImportMaterialMasterData)
	r.HandleFunc("/material-search-pagination", MaterialSearchPagination).Methods("POST")

	//--Chemicals--
	r.HandleFunc("/chemical-management-card", chemicalmanagmentDialog).Methods("GET")
	r.HandleFunc("/chemical-management", ChemicalManagmentPage)
	r.HandleFunc("/create-new-or-update-existing-chemical", CreateNewOrUpdateExistingChemical).Methods("POST")
	r.HandleFunc("/import-chemical-data", ImportChemicalMasterData)
	r.HandleFunc("/chemical-search-pagination", ChemicalSearchPagination).Methods("POST")

	// --ProdLine Master--
	r.HandleFunc("/prodline-search-pagination", ProdLineSearchPagination).Methods("POST")
	// --ProdProcess Master--
	r.HandleFunc("/prodprocess-search-pagination", ProdProcessSearchPagination).Methods("POST")
	// --Pending Order--
	r.HandleFunc("/pending-order-search-pagination", PendingOrderSearchPagination).Methods("POST")

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	r.PathPrefix("/RUBBER/").Handler(http.StripPrefix("/RUBBER/", http.FileServer(http.Dir("./RUBBER"))))

	// -----###### API endpoints used by the mobile application ######-----
	mobileAPI := r.PathPrefix("/api/v1").Subrouter()
	mobileAPI.HandleFunc("/get-pending-orders", GetAllPendingOrders)
	mobileAPI.HandleFunc("/reject-kanban-order", RejectOrder)
	mobileAPI.HandleFunc("/approve-kanban-order", ApproveOrder)
	mobileAPI.HandleFunc("/get-all-line-up-kanban", GetProductionLinesHandler).Methods("GET")
	mobileAPI.HandleFunc("/get-kanban-by-production-line-id", GetLinedUpProductionProcessByLineId).Methods("POST")
	mobileAPI.HandleFunc("/get-kanban-for-vendor", GetKanbanForVendor).Methods("POST")
	mobileAPI.HandleFunc("/get-kanban-for-all-vendor", GetKanbanForAllVendor).Methods("GET")
	mobileAPI.HandleFunc("/update-user-status", UpdateUserStatus)
	mobileAPI.HandleFunc("/get-all-quality-testing-kanban", GetAllQualityTestingKanban).Methods("GET")
	mobileAPI.HandleFunc("/get-quality-testing-kanban-for-vendor", GetQualityTestingKanbanForVendor).Methods("POST")
	mobileAPI.HandleFunc("/get-all-packing-kanban", GetAllPackingKanban).Methods("GET")
	mobileAPI.HandleFunc("/get-packing-kanban-for-vendor", GetPackingKanbanForVendor).Methods("POST")
	mobileAPI.HandleFunc("/get-order-details", GetOrderDetailsForHistoryMobileAPI).Methods("GET")

	// -----###### API endpoints used by the plc ######-----
	plcAPI := r.PathPrefix("/plc").Subrouter()
	// TODO- remove
	plcAPI.HandleFunc("/plc/1/2/test", TestPLCHandler).Methods("GET")
	plcAPI.HandleFunc("/plc//", TestPLCHandler).Methods("GET")

	port := os.Getenv("WEBSRV_PORT")
	if port == "" {
		port = "4300"
		log.Println("WEBSRV_PORT environment variable not set, defaulting to :", port)
	}
	log.Println("Server listening on:", utils.RestHost+":"+port)

	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Println("Failed to start server:", err)
	}
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	websecure.SetCookie("", "", false, time.Now().Add(-1000*24*time.Hour), w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func doRegisterHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve form values
	userType := r.FormValue("user_type")
	firstName := r.FormValue("FirstName")
	lastName := r.FormValue("LastName")
	email := r.FormValue("Email")
	password := r.FormValue("Password")
	confirmPassword := r.FormValue("ConfirmPassword")
	code := r.FormValue("Code")

	// Validate that the passwords match
	if !strings.EqualFold(strings.TrimSpace(password), strings.TrimSpace(confirmPassword)) {
		errMsg := "Passwords do not match!"
		http.Redirect(w, r, "/register?error="+errMsg, http.StatusFound)
		return
	}

	// Call the API to get role by name
	roleResp, err := http.Get("http://" + utils.RestHost + ":" + utils.RestPort + "/get-role-by-name?name=" + userType)
	if err != nil || roleResp.StatusCode != http.StatusOK {
		slog.Error("Failed to fetch role by name", "userType", userType, "error", err)
		errMsg := "Failed to fetch role information. Please try again."
		http.Redirect(w, r, "/register?error="+errMsg, http.StatusFound)
		return
	}
	defer roleResp.Body.Close()

	// Parse the role API response
	var role m.UserRoles
	if err := json.NewDecoder(roleResp.Body).Decode(&role); err != nil || role.ID == 0 {
		slog.Error("Error decoding role response", "error", err)
		errMsg := "Internal server error, please try again!"
		http.Redirect(w, r, "/register?error="+errMsg, http.StatusFound)
		return
	}

	// Create a userReq instance
	userReq := m.User{
		Username:    utils.JoinStr(firstName, " ", lastName),
		Email:       email,
		Password:    password,
		CreatedBy:   "",
		CreatedOn:   time.Now(),
		VendorsCode: code,
	}

	// Marshal the User into JSON
	jsonValue, err := json.Marshal(userReq)
	if err != nil {
		slog.Error("Error marshaling registration data", "error", err)
		http.Error(w, "Failed to process registration data", http.StatusInternalServerError)
		return
	}

	// Send the registration request to the backend API
	resp, err := http.Post("http://"+utils.RestHost+":"+utils.RestPort+"/register", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		slog.Error("Error making registration request", "error", err)
		http.Error(w, "Failed to connect to the registration service", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Parse the API response
	var apiResp m.ApiResp
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		slog.Error("Error decoding API response", "error", err)
	}

	// Handle response based on status code
	switch resp.StatusCode {
	case http.StatusOK:
		// Registration successful, set cookie and redirect to home
		websecure.SetCookie("ID", apiResp.User.Email, false, time.Now().Add(24*time.Hour), w)
		http.Redirect(w, r, "/dashboard?", http.StatusFound)
	case http.StatusForbidden:
		// Email already exists
		errMsg := "User with this email already exists"
		http.Redirect(w, r, "/register?error="+errMsg, http.StatusFound)
	case http.StatusNotAcceptable:
		//Invalid Vendor code
		errMsg := "Invalid Vendor Code"
		http.Redirect(w, r, "/register?error="+errMsg, http.StatusFound)
	default:
		// Other registration errors
		errMsg := "Registration failed. Please try again later."
		http.Redirect(w, r, "/register?error="+errMsg, http.StatusFound)
	}
}

func SigneUp(w http.ResponseWriter, r *http.Request) {
	var data m.RegisterRequest
	var ApiResp m.ApiRespMsg
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&data); err == nil {
		// Validate that the passwords match
		if !strings.EqualFold(strings.TrimSpace(data.Password), strings.TrimSpace(data.ConfirmPassword)) {
			errMsg := "Passwords do not match!"
			http.Redirect(w, r, "/register?error="+errMsg, http.StatusFound)
			return
		}

		if data.Email == "" || data.FirstName == "" || data.LastName == "" || data.UserType == "" {
			ApiResp.Code = 400
			ApiResp.Message = "Fields cannot be empty. Please fill out all required fields."
			body, err := json.Marshal(ApiResp)
			if err != nil {
				log.Printf("Error marshaling response: %v", err)
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Write(body)
			return
		}
		if data.UserType == "Customer" && data.Code == "" {
			ApiResp.Code = 406
			ApiResp.Message = "Invalid Vendor Code"
			body, err := json.Marshal(ApiResp)
			if err != nil {
				log.Printf("Error marshaling response: %v", err)
			}
			w.WriteHeader(http.StatusNotAcceptable)
			w.Write(body)
			return
		}
		// Call the API to get role by name
		roleResp, err := http.Get("http://" + utils.RestHost + ":" + utils.RestPort + "/get-role-by-name?name=" + data.UserType)
		if err != nil || roleResp.StatusCode != http.StatusOK {
			slog.Error("Failed to fetch role by name", "userType", data.UserType, "error", err)
			// errMsg := "Failed to fetch role information. Please try again."
			ApiResp.Code = 406
			ApiResp.Message = "Invalid User Type"
			body, err := json.Marshal(ApiResp)
			if err != nil {
				log.Printf("Error marshaling response: %v", err)
			}
			w.WriteHeader(http.StatusNotAcceptable)
			w.Write(body)

			// http.Redirect(w, r, "/register?error="+errMsg, http.StatusFound)
			return
		}
		defer roleResp.Body.Close()

		var role m.UserRoles
		if err := json.NewDecoder(roleResp.Body).Decode(&role); err != nil || role.ID == 0 {
			slog.Error("Error decoding role response", "error", err)
			errMsg := "Internal server error, please try again!"
			http.Redirect(w, r, "/register?error="+errMsg, http.StatusFound)
			return
		}

		// Create a userReq instance
		userReq := m.User{
			Username:    utils.JoinStr(data.FirstName, " ", data.LastName),
			RoleID:      uint(role.ID),
			Email:       data.Email,
			Password:    data.Password,
			CreatedBy:   "",
			CreatedOn:   time.Now(),
			VendorsCode: data.Code,
		}
		// Marshal the User into JSON
		jsonValue, err := json.Marshal(userReq)
		if err != nil {
			slog.Error("Error marshaling registration data", "error", err)
			http.Error(w, "Failed to process registration data", http.StatusInternalServerError)
			return
		}

		// Send the registration request to the backend API
		resp, err := http.Post("http://"+utils.RestHost+":"+utils.RestPort+"/register", "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			slog.Error("Error making registration request", "error", err)
			http.Error(w, "Failed to connect to the registration service", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
		// Parse the API response
		var apiResp m.ApiResp
		if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
			slog.Error("Error decoding API response", "error", err)
		}
		// Handle response based on status code
		switch resp.StatusCode {
		case http.StatusOK:
			// Registration successful, set cookie and redirect to home
			websecure.SetCookie("ID", apiResp.User.Email, false, time.Now().Add(24*time.Hour), w)
			ApiResp.Code = 200
			ApiResp.Message = "User Created Successfully"
			body, err := json.Marshal(apiResp)
			if err != nil {
				log.Printf("Error marshaling response: %v", err)
			}
			w.WriteHeader(http.StatusOK)
			w.Write(body)
		case http.StatusForbidden:
			// Email already exists
			ApiResp.Code = 403
			ApiResp.Message = "User with this email already exists"
			body, err := json.Marshal(ApiResp)
			if err != nil {
				log.Printf("Error marshaling response: %v", err)
			}
			w.WriteHeader(http.StatusForbidden)
			w.Write(body)
		case http.StatusNotAcceptable:
			ApiResp.Code = 406
			ApiResp.Message = "Invalid Vendor Code"
			body, err := json.Marshal(ApiResp)
			if err != nil {
				log.Printf("Error marshaling response: %v", err)
			}
			w.WriteHeader(http.StatusNotAcceptable)
			w.Write(body)
		default:
			// Other registration errors
			ApiResp.Code = 500
			ApiResp.Message = "Registration failed. Please try again later"
			body, err := json.Marshal(ApiResp)
			if err != nil {
				log.Printf("Error marshaling response: %v", err)
			}
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(body)

		}

	}
}

func doLogin(w http.ResponseWriter, r *http.Request) {
	var getresp m.ApiResp
	var user m.User
	// Extract username and password from the form data
	data := m.LoginReq{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}

	// Convert the data to JSON
	jsonValue, err := json.Marshal(data)
	if err != nil {
		slog.Error("Failed to marshal login request data", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Make the POST request to the login API
	resp, err := http.Post("http://"+utils.RestHost+":"+utils.RestPort+"/login", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		slog.Error("Error making login API request", "error", err)
		http.Error(w, "Unable to reach authentication service", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	// Check if the response status is 200 OK
	if resp.StatusCode != http.StatusOK {
		var errMsg string
		if resp.StatusCode == http.StatusForbidden {
			errMsg = "Your account is disabled please contact the administrator."
		} else {
			errMsg = "Invalid credentials, please try again."
		}
		http.Redirect(w, r, "/login?error="+errMsg, http.StatusFound)
		return
	}

	// Decode the response into getresp
	if err := json.NewDecoder(resp.Body).Decode(&getresp); err != nil {
		slog.Error("Error decoding response body", "error", err)
		http.Error(w, "Failed to process authentication response", http.StatusInternalServerError)
		return
	}

	user = getresp.User
	// Set a cookie with the user’s email and redirect to home
	websecure.SetCookie("ID", user.Email, false, time.Now().Add(24*time.Hour), w)
	http.Redirect(w, r, "/dashboard", http.StatusFound)
}

func SigneIn(w http.ResponseWriter, r *http.Request) {
	var getresp m.ApiResp
	var ApiResp m.ApiRespMsg
	var user m.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err == nil {
		// Convert the data to JSON
		jsonValue, err := json.Marshal(user)
		if err != nil {
			slog.Error("Failed to marshal login request data", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Make the POST request to the login API
		resp, err := http.Post("http://"+utils.RestHost+":"+utils.RestPort+"/login", "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			slog.Error("Error making login API request", "error", err)
			http.Error(w, "Unable to reach authentication service", http.StatusServiceUnavailable)
			return
		}
		defer resp.Body.Close()

		// Check if the response status is 200 OK
		if resp.StatusCode != http.StatusOK {
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
			w.WriteHeader(resp.StatusCode)
			w.Write(body)

			return
		}

		// Decode the response into getresp
		if err := json.NewDecoder(resp.Body).Decode(&getresp); err != nil {
			slog.Error("Error decoding response body", "error", err)
			http.Error(w, "Failed to process authentication response", http.StatusInternalServerError)
			return
		}

		user = getresp.User
		body, err := json.Marshal(getresp)
		if err != nil {
			log.Printf("Error marshaling response: %v", err)
		}

		// Set a cookie with the user’s email and redirect to home
		websecure.SetCookie("ID", user.Email, false, time.Now().Add(24*time.Hour), w)
		w.WriteHeader(http.StatusOK)
		w.Write(body)
		// http.Redirect(w, r, "/dashboard", http.StatusOK)
	}

}

func GetUserByEmailHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// log.Println("GetUserByEmailHandler: request received")

	email := r.URL.Query().Get("email")
	if len(strings.TrimSpace(email)) <= 0 {
		http.Error(w, "Mobile number is required", http.StatusBadRequest)
		return
	}

	// Use net/url to build the query parameters
	queryParams := url.Values{}
	queryParams.Add("email", email)

	// Make HTTP Get request
	resp, err := http.Get(utils.RestURL + "/get-user-by-email?" + queryParams.Encode())
	// resp, err := http.Get("http://" + DBHelperHost + ":" + DBHelperPort + "/get-user-by-email?" + queryParams.Encode())
	if err != nil {
		log.Printf("Error making POST request: %v", err)
		http.Error(w, "Failed to get user information", http.StatusForbidden)
		return
	}
	defer resp.Body.Close()

	// Process response
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		http.Error(w, "Failed to read response body", http.StatusInternalServerError)
		return
	}

	w.Write(responseBody)
}

func GetUserRoleByRoleIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.URL.Query().Get("id")
	if len(strings.TrimSpace(id)) <= 0 {
		http.Error(w, "Mobile number is required", http.StatusBadRequest)
		return
	}

	// Use net/url to build the query parameters
	queryParams := url.Values{}
	queryParams.Add("id", id)

	// Make HTTP Get request
	resp, err := http.Get(utils.RestURL + "/get-user-role-by-id?" + queryParams.Encode())
	// resp, err := http.Get("http://" + DBHelperHost + ":" + DBHelperPort + "/get-user-by-email?" + queryParams.Encode())
	if err != nil {
		log.Printf("Error making POST request: %v", err)
		http.Error(w, "Failed to get user information", http.StatusForbidden)
		return
	}
	defer resp.Body.Close()

	// Process response
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		http.Error(w, "Failed to read response body", http.StatusInternalServerError)
		return
	}

	w.Write(responseBody)
}

// Handelr func to add Production line
func AddProdLine(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID := r.Header.Get("X-Custom-Userid")
	// Decode the JSON payload
	var requestData m.AddLineStruct

	var APIresp m.ErrorResp

	json.NewDecoder(r.Body).Decode(&requestData)
	requestData.CreatedBy = userID
	marshalData, err := json.Marshal(requestData)
	if err != nil {
		APIresp.ErrCode = http.StatusInternalServerError
		APIresp.ErrMessage = "Error while marshaling data"
		body, jsonErr := json.Marshal(APIresp)
		if jsonErr != nil {
			http.Error(w, "{err_message:Internal Server Error}", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(body)
		return
	}

	if strings.TrimSpace(requestData.LineName) == "" {
		APIresp.ErrCode = http.StatusInternalServerError
		APIresp.ErrMessage = "Line name is empty"
		body, jsonErr := json.Marshal(APIresp)
		if jsonErr != nil {
			http.Error(w, "{err_message:Internal Server Error}", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(body)
		return
	}
	if len(requestData.ProcessOrders) == 0 {
		APIresp.ErrCode = http.StatusInternalServerError
		APIresp.ErrMessage = "No process is selected"
		body, jsonErr := json.Marshal(APIresp)
		if jsonErr != nil {
			http.Error(w, "{err_message:Internal Server Error}", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(body)
		return
	}
	// Construct the REST API URL
	url := utils.JoinStr("http://", utils.RestHost, ":", utils.RestPort, "/add-production-line")

	// Create the HTTP request
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(marshalData))
	if err != nil {
		APIresp.ErrCode = http.StatusInternalServerError
		APIresp.ErrMessage = "Error creating request to REST serviced"
		body, jsonErr := json.Marshal(APIresp)
		if jsonErr != nil {
			http.Error(w, "{err_message:Internal Server Error}", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(body)
		return
	}

	// Send the request using an HTTP client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		APIresp.ErrCode = http.StatusInternalServerError
		APIresp.ErrMessage = "Failed to send request to REST service"
		body, jsonErr := json.Marshal(APIresp)
		if jsonErr != nil {
			http.Error(w, "{err_message:Internal Server Error}", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(body)
		return
	}
	defer resp.Body.Close() // Ensure the response body is closed

	// Check for non-201 status codes
	if resp.StatusCode != http.StatusCreated {
		APIresp.ErrCode = resp.StatusCode
		APIresp.ErrMessage = "Fail to create New Line"
		body, jsonErr := json.Marshal(APIresp)
		if jsonErr != nil {
			http.Error(w, "{err_message:Internal Server Error}", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(resp.StatusCode)
		w.Write(body)
		return
	}
	APIresp.ErrCode = http.StatusOK
	APIresp.ErrMessage = "Production Line Created!"
	body, jsonErr := json.Marshal(APIresp)
	if jsonErr != nil {
		http.Error(w, "{err_message:Internal Server Error}", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func UpdateRunningNumbers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var updatesRunningNo []m.KbRoot
	err := json.NewDecoder(r.Body).Decode(&updatesRunningNo)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		log.Println("Error decoding JSON:", err)
		return
	}

	marshalData, err := json.Marshal(updatesRunningNo)
	if err != nil {
		http.Error(w, "Error while marshaling data", http.StatusInternalServerError)
		return
	}
	url := utils.JoinStr("http://", utils.RestHost, ":", utils.RestPort, "/update-running-numbers")

	// Create the HTTP request
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(marshalData))
	if err != nil {
		http.Error(w, "Error creating request to REST service", http.StatusInternalServerError)
		return
	}

	// Send the request using an HTTP client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to send request to REST service", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close() // Ensure the response body is closed

	// Check for non-201 status codes
	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Fail to create New Line", resp.StatusCode)
		return
	}
}

func GetProductionLinesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Build the URL
	url := utils.JoinStr("http://", utils.RestHost, ":", utils.RestPort, "/get-production-line-items")

	// Make the HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error making GET request: %v", err)
		http.Error(w, "{\"error\":\"Failed to fetch production line data\"}", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Decode JSON response
	var prodLineDetails []m.ProdLineDetails
	if err := json.NewDecoder(resp.Body).Decode(&prodLineDetails); err != nil {
		body, _ := io.ReadAll(resp.Body) // Optional: re-read the body for logging
		log.Printf("Error decoding JSON: %v, Body: %s", err, string(body))
		http.Error(w, "{\"error\":\"Failed to decode response\"}", http.StatusInternalServerError)
		return
	}

	// Encode and respond to the client
	if err := json.NewEncoder(w).Encode(prodLineDetails); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "{\"error\":\"Failed to encode response\"}", http.StatusInternalServerError)
		return
	}
}

func GetProductionLineProcesses(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Define the struct to hold the line number
	var LineNumber struct {
		LineNo int `json:"line_no"`
	}

	// Decode the incoming request body
	err := json.NewDecoder(r.Body).Decode(&LineNumber)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		log.Println("Error decoding JSON:", err)
		return
	}

	// Marshal the LineNumber struct to JSON to send in the POST request
	lineNumberJSON, err := json.Marshal(LineNumber)
	if err != nil {
		log.Printf("Error marshalling LineNumber: %v", err)
		http.Error(w, "Failed to prepare data", http.StatusInternalServerError)
		return
	}

	// Build the URL
	url := utils.JoinStr("http://", utils.RestHost, ":", utils.RestPort, "/get-production-line-status")

	// Make the HTTP POST request with the marshaled JSON as the body
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(lineNumberJSON))
	if err != nil {
		log.Printf("Error making POST request: %v", err)
		http.Error(w, "{\"error\":\"Failed to fetch production line data\"}", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Decode JSON response into prodLineDetails
	var prodLineDetails []m.ProdLineDetails
	if err := json.NewDecoder(resp.Body).Decode(&prodLineDetails); err != nil {
		body, _ := io.ReadAll(resp.Body) // Optional: re-read the body for logging
		log.Printf("Error decoding JSON: %v, Body: %s", err, string(body))
		http.Error(w, "{\"error\":\"Failed to decode response\"}", http.StatusInternalServerError)
		return
	}

	// Encode and send the response to the client
	if err := json.NewEncoder(w).Encode(prodLineDetails); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "{\"error\":\"Failed to encode response\"}", http.StatusInternalServerError)
		return
	}
}

func GetInProgressCountByLine(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ProdLineID int `json:"prod_line_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	rawQuery := m.RawQuery{
		Host:  utils.RestHost,
		Port:  utils.RestPort,
		Type:  "InProgressOrderCount",
		Query: dashboard.BuildInProgressOrderQuery(req.ProdLineID),
	}

	var result struct {
		ScheduledCount  int `json:"scheduled_count"`
		InProgressCount int `json:"in_progress_count"`
	}

	if err := rawQuery.RawQry(&result); err != nil {
		fmt.Println("Error in RawQry:", err)
		http.Error(w, "Query failed", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(result)
}
func GetVendorOrderStatusThisMonth(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Month string `json:"month"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON"})
		return
	}

	usertype := r.Header.Get("X-Custom-Role")

	rawQuery := m.RawQuery{
		Host:  utils.RestHost,
		Port:  utils.RestPort,
		Type:  "MonthlyVendorOrderStatus",
		Query: dashboard.BuildVendorOrderStatusQuery(payload.Month, usertype),
	}

	var result []m.VendorOrderStatus
	if err := rawQuery.RawQry(&result); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Query failed"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

// TODO- remove
// TestPLCHandler is a basic handler to verify Bearer token access
func TestPLCHandler(w http.ResponseWriter, r *http.Request) {
	utils.SetResponse(w, http.StatusOK, "PLC token authenticated successfully")
}

func RecipeDataForPLC(w http.ResponseWriter, r *http.Request) {

	var orderdetails []*m.OrderDetails
	var rawQuery m.RawQuery
	rawQuery.Host = utils.RestHost
	rawQuery.Port = utils.RestPort
	rawQuery.Type = "OrderDetails"
	rawQuery.Query = ""
	rawQuery.RawQry(&orderdetails)

	utils.SetResponse(w, http.StatusOK, "PLC token authenticated successfully")
}
