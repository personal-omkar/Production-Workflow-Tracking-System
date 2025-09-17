package server

import (
	"encoding/json"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	m "irpl.com/kanban-commons/model"
	dao "irpl.com/kanban-dao/dao"

	"github.com/gorilla/mux"
)

var InventoryVendorCode string = "Inventory001"

// Status ss
func Status(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func Web() {

	r := mux.NewRouter()

	port := os.Getenv("DBHELPER_PORT")
	if port == "" {
		slog.Info("DBHELPER_PORT environment variable not set, defaulting to :4100")
		port = "4100"
	}

	// TODO- apply middleware
	r.HandleFunc("/status", Status)
	r.HandleFunc("/create-user", RegisterHandler).Methods("POST")
	r.HandleFunc("/login", LoginHandler).Methods("POST")
	r.HandleFunc("/get-user-by-email", GetUserByEmailHandler).Methods("GET")
	// r.HandleFunc("/get-user-role-by-id", GetUserRoleByRoleIDHandler).Methods("GET")
	r.HandleFunc("/get-otp-count", GetOtpCountHandler).Methods("GET")
	r.HandleFunc("/save-otp-details", SaveOtpDetailsHandler).Methods("POST")
	r.HandleFunc("/find-otp-details", FindOtpDetailsHandler).Methods("GET") // Assuming mobileNo is passed as query param
	r.HandleFunc("/RawQuery", RawQuery).Methods("POST")
	// -- API: Vendors --
	// -- GET --
	r.HandleFunc("/get-all-vendors-data", GetVendorsAllrecords).Methods("POST")
	r.HandleFunc("/get-all-vendor-by-search-pagination", GetVendorSearchPaginationRecords).Methods("POST")
	r.HandleFunc("/delete-vendors-by-id", DeleteVendorsById).Methods("GET")
	// -- POST --
	r.HandleFunc("/create-new-vendor", CreateVendors).Methods("POST")

	// -- API: Prod line --
	// -- GET --
	r.HandleFunc("/get-all-prod-line-data", GetProdLineAllrecords).Methods("GET")
	r.HandleFunc("/delete-production-line-by-id", DeleteProductionLineById).Methods("GET")
	r.HandleFunc("/delete-production-line-celldata-by-productionline-id", DeleteProductionLineCellDataByProductionLineId).Methods("GET")
	// r.HandleFunc("/delete-production-line-by-prodline-id", DeleteProductionLineByProdLineId).Methods("GET")
	r.HandleFunc("/delete-production-line-cell", DeleteProductionLineCell).Methods("POST")
	r.HandleFunc("/create-new-or-update-existing-production-line", CreateNewOrUpdateExistingProductionLine).Methods("POST")

	// -- API: Compounds --
	// -- GET --
	r.HandleFunc("/get-all-compound-data", GetCompountsAllrecords).Methods("GET")
	r.HandleFunc("/get-all-active-compounds", GetAllActiveCompounds).Methods("GET")
	r.HandleFunc("/get-compound-data-by-vendor", GetCompoundsByVendors).Methods("POST")
	r.HandleFunc("/get-kanban-for-vendor", GetKanbanByVendors).Methods("POST")
	r.HandleFunc("/get-kanban-for-all-vendor", GetKanbanForAllVendors).Methods("GET")
	r.HandleFunc("/get-packing-compound-data-by-vendor", GetPackingCompoundsByVendors).Methods("POST")
	r.HandleFunc("/get-quality-compound-data-by-vendor", GetQualityCompoundsByVendors).Methods("POST")
	r.HandleFunc("/get-compound-data-by-parm", GetCompoundsByParm).Methods("GET")
	r.HandleFunc("/get-all-coldstorage-data", GetAllColdStoragerecords).Methods("GET")
	r.HandleFunc("/update-compound-status-to-packing", UpdateCompoundStatusToPacking).Methods("POST")
	r.HandleFunc("/update-compound-quality-status-to-reject", UpdateCompoundQualityStatusToReject).Methods("POST")
	r.HandleFunc("/update-compound-status-to-dispatch", UpdateCompoundStatusToDispatch).Methods("POST")
	r.HandleFunc("/get-all-compounds-by-search-pagination", GetAllCompoundsBySearchAndPagination).Methods("POST")

	// -- POST --
	r.HandleFunc("/add-compound-data-by-vendor", AddCompoundsForVendor).Methods("POST")
	r.HandleFunc("/add-compound-data-to-production-line", AddCompoundsInProductionLine).Methods("POST")
	r.HandleFunc("/update-coldstorage-quantity", UpdateColdStorageQuantity).Methods("POST")

	// -- API: KbRoot --
	// -- GET --
	r.HandleFunc("/get-all-kb-root-data", GetAllKbRoot).Methods("GET")
	r.HandleFunc("/get-kb-root-by-param", GetKbRootByParam).Methods("GET")

	// -- API: KbExtension --
	// -- GET --
	r.HandleFunc("/get-all-kb-extension-data", GetAllKbExtensions).Methods("GET")
	r.HandleFunc("/get-kb-extension-by-param", GetKbExtensionsByParam).Methods("GET")

	// -- API: KbData --
	// -- GET --
	r.HandleFunc("/get-all-kb-data", GetAllKBData).Methods("GET")
	r.HandleFunc("/get-kb-data-by-param", GetKBDataByParam).Methods("GET")
	// -- API: Prod line with cell data --
	// -- GET --
	r.HandleFunc("/get-production-line-items", GetAllProductionLineRecords).Methods("GET")

	r.HandleFunc("/add-production-line", CreateProdLineWithProcesses).Methods("POST")

	// -- API: Get sepcific prod line stauts data --
	// -- POST --
	r.HandleFunc("/get-production-line-status", GetProductionLineStatus).Methods("POST")

	// -- API: Update kbroot running no--
	// -- POST --
	r.HandleFunc("/update-running-numbers", UpdateRunningNumbers).Methods("POST")

	// -- API: Order --
	// -- POST --
	r.HandleFunc("/create-new-order-entry", CreateNewOrderEntry).Methods("POST")
	r.HandleFunc("/create-multi-new-order-entry", CreateMultiNewOrderEntry).Methods("POST")

	r.HandleFunc("/get-customer-order-details", GetCustomerOrderDetails).Methods("POST")
	// -- API: Order --
	// -- POST --

	r.HandleFunc("/get-order-details", GetOrderDEtails).Methods("POST")

	r.HandleFunc("/update-order-status", UpdateOrderStatus).Methods("POST")

	// -- API: Users --
	// -- POST --
	r.HandleFunc("/create-new-user", CreateUser).Methods("POST")
	r.HandleFunc("/delete-user", DeleteUser).Methods("POST")
	// -- GET --
	r.HandleFunc("/get-user-by-param", GetUserByParam).Methods("GET")

	// -- API: UserRoles --
	r.HandleFunc("/get-role-by-name", GetUserRoleByNameHandler).Methods("GET")
	r.HandleFunc("/get-role-by-id", GetUserRoleByIDHandler).Methods("GET")
	r.HandleFunc("/get-all-user-roles", GetAllUserRoles).Methods("GET")
	r.HandleFunc("/delete-user-role", DeleteUserRole).Methods("POST")

	// API  to fetch orders details for order page
	r.HandleFunc("/OrderDetailsForCustomer", GetOrderDetaislForOrderPage).Methods("POST")
	r.HandleFunc("/update-user-role", UpdateUserRole).Methods("POST")
	r.HandleFunc("/create-user-role", CreateUserRole).Methods("POST")

	// -- API: Inventory --
	// -- GET --
	r.HandleFunc("/get-all-inventory-data", GetAllInventoryrecords).Methods("GET")
	r.HandleFunc("/get-inventory-by-param", GetInventoryByParam).Methods("GET")
	r.HandleFunc("/delete-inventory-by-param", DeleteInventoryByParam).Methods("GET")
	// -- POST --
	r.HandleFunc("/create-new-or-update-existing-inventory", CreateOrUpdateInventory).Methods("POST")
	r.HandleFunc("/get-inventory-by-search", GetInventoryBySearch).Methods("POST")
	r.HandleFunc("/get-all-cold-storage-by-search-pagination", ColdStorageSearchPagination).Methods("POST")
	r.HandleFunc("/coldstorage-pagination-search", ColdStorageSearchPagination).Methods("POST")

	r.HandleFunc("/fetch-all-production-process-data", GetAllProductionProcessData).Methods("GET")
	r.HandleFunc("/get-all-production-process-for-line", GetAllProductionProcessEntries).Methods("GET")
	r.HandleFunc("/get-all-production-process", GetAllProductionProcess).Methods("GET")
	r.HandleFunc("/get-prod-line-by-param", GetProdLineByParam).Methods("POST")
	r.HandleFunc("/add-production-process", AddProdProcess).Methods("POST")
	r.HandleFunc("/edit-prod-line", EditProdLine).Methods("POST")
	r.HandleFunc("/get-prod-process-by-param", GetProdProcessByParam).Methods("POST")
	r.HandleFunc("/edit-prod-process", EditProdProcess).Methods("POST")

	r.HandleFunc("/create-new-KbTransaction", CreateNewKbTransaction).Methods("POST")

	// -- API: UserToVemdor --
	// -- GET --
	r.HandleFunc("/get-usertovendor-by-param", GetUserToVendorByParam).Methods("GET")
	r.HandleFunc("/get-vendor-by-userid", GetVendorByUserID).Methods("GET")
	r.HandleFunc("/get-all-orders-by-vendor-code", GetAllOrderByVendorCode).Methods("POST")
	r.HandleFunc("/get-vendor-details-by-vendor-code", GetVendorDetailsByVendorCode).Methods("POST")

	r.HandleFunc("/update-running-number", UpdateRunningNumberAfterTransactioin).Methods("POST")

	r.HandleFunc("/delete-order-entry", DeleteOrder).Methods("DELETE")

	r.HandleFunc("/get-all-kbRoot-details", GetAllCompletedKBRootDetails).Methods("POST")
	r.HandleFunc("/get-all-kbRoot-details-by-search", GetCompletedKBRootDetailsBySearch).Methods("POST")
	r.HandleFunc("/get-kbRoot-details", GetDetailRootData).Methods("POST")

	// -- API: User Details --
	// -- GET --
	r.HandleFunc("/get-user-details", GetUserDetails).Methods("GET")
	r.HandleFunc("/get-user-details-by-email", GetUserDetailsByEmail).Methods("GET")
	r.HandleFunc("/update-user-details", UpdateUser).Methods("POST")
	r.HandleFunc("/update-user-status", UpdateUserStatus).Methods("POST")

	// -- API: LDAP Config
	// -- POST --
	r.HandleFunc("/create-ldap-config", CreateLDAPConfig).Methods("POST")

	r.HandleFunc("/update-ldap-config", UpdateLDAPConfig).Methods("POST")

	// -- GET --
	r.HandleFunc("/get-default-ldap-config", GetDefaultLDAPConfig).Methods("GET")

	r.HandleFunc("/delete-ldap-config", DeleteLDAPConfig).Methods("POST")

	// -- API: Samba Config
	// -- POST --
	r.HandleFunc("/create-samba-config", CreateSambaConfig).Methods("POST")

	r.HandleFunc("/update-samba-config", UpdateSambaConfig).Methods("POST")

	// -- GET --
	r.HandleFunc("/get-default-samba-config", GetDefaultSambaConfig).Methods("GET")

	r.HandleFunc("/delete-samba-config", DeleteSambaConfig).Methods("POST")

	r.HandleFunc("/get-all-logos", GetAllLogosHandler).Methods("POST")

	r.HandleFunc("/add-update-compound", AddorUpdateCompound).Methods("POST")

	// -- API: SystemLogs --
	// -- POST --
	r.HandleFunc("/create-system-log", CreateSystemLog).Methods("POST")
	r.HandleFunc("/update-system-log", UpdateSystemLog).Methods("POST")

	// -- API: User To Role
	// -- POST --
	r.HandleFunc("/create-user-to-role", CreateNewUserToRoel).Methods("POST")
	r.HandleFunc("/update-user-to-role", UpdateUserToRoel).Methods("POST")
	r.HandleFunc("/get-all-user-by-search-pagination", GetAllUserBySearchAndPagination).Methods("POST")
	r.HandleFunc("/get-all-user-role-by-search-pagination", GetAllUserRoleBySearchAndPagination).Methods("POST")

	r.HandleFunc("/check-vendor-lot-limit", DailyAndMonthlyVendorLimit).Methods("POST")
	r.HandleFunc("/check-vendor-lot-limit-by-vendor-code", DailyAndMonthlyVendorLimitByVendorCode).Methods("POST")
	r.HandleFunc("/check-daily-lot-limit", DailyVendorLimit).Methods("POST")

	r.HandleFunc("/delete-kbroot-by-ids", DeleteKbRootByIDsHandler).Methods("DELETE")
	r.HandleFunc("/get-lineup-processes-by-lineid", GetLineUpProcessesByLineId).Methods("POST")

	r.HandleFunc("/get-all-details-for-order", GetAllDetailsForOrder).Methods("POST")

	r.HandleFunc("/get-all-kanban-details-for-report", GetAllKanbanDetailsForReportHandler).Methods("POST")

	// -- API:- Stage
	r.HandleFunc("/create-new-stage", CreateStage).Methods("POST")
	r.HandleFunc("/get-all-stages", GetAllStage).Methods("GET")
	r.HandleFunc("/update-stage", UpdateExistingStage).Methods("POST")
	r.HandleFunc("/delete-stage", DeleteStageByID).Methods("POST")
	// r.HandleFunc("/get-stages-by-header", GetStagesByHeader).Methods("POST")
	r.HandleFunc("/get-stage-by-param", GetStagesByParam).Methods("GET")
	// -- API: Recipe
	// -- POST --
	r.HandleFunc("/create-new-recipe", CreateRecipe).Methods("POST")
	r.HandleFunc("/update-existing-recipe", UpdateRecipe).Methods("POST")
	r.HandleFunc("/delete-recipe-by-id", DeleteRecipe).Methods("POST")
	// -- GET --
	r.HandleFunc("/get-all-recipe", GetAllRecipe).Methods("GET")
	r.HandleFunc("/get-recipe-by-data-key", GetRecipeByDataKey).Methods("GET")
	r.HandleFunc("/get-recipe-by-data-value", GetRecipeByDataValue).Methods("GET")
	r.HandleFunc("/get-recipe-by-data-key-and-value", GetRecipeByDataKeyAndValue).Methods("GET")
	r.HandleFunc("/get-recipe-by-param", GetRecipeByParam).Methods("GET")

	//--Operator--
	// -- POST --
	r.HandleFunc("/create-new-or-update-existing-operator", CreateNewOrUpdateExistingOperator).Methods("POST")
	r.HandleFunc("/get-all-operator-by-search-pagination", GetAllOperatorBySearchAndPagination).Methods("POST")
	// --GET
	r.HandleFunc("/get-all-operator", GetAllOperator).Methods("GET")
	r.HandleFunc("/get-operator-by-param", GetOperatorByParam).Methods("GET")

	//--Material--
	// -- POST --
	r.HandleFunc("/create-new-or-update-existing-material", CreateNewOrUpdateExistingMaterial).Methods("POST")
	r.HandleFunc("/get-all-material-by-search-pagination", GetAllMaterialBySearchAndPagination).Methods("POST")
	// --GET
	r.HandleFunc("/get-all-material", GetAllMaterial).Methods("GET")
	r.HandleFunc("/get-material-by-param", GetMaterialByParam).Methods("GET")

	//--Chemicals--
	// -- POST --
	r.HandleFunc("/create-new-or-update-existing-chemical", CreateNewOrUpdateExistingChemical).Methods("POST")
	r.HandleFunc("/get-all-chemical-by-search-pagination", GetAllChemicalBySearchAndPagination).Methods("POST")
	// --GET
	r.HandleFunc("/get-all-chemical", GetAllChemical).Methods("GET")
	r.HandleFunc("/get-chemical-by-param", GetChemicalByParam).Methods("GET")

	//--Prod Line Master Pagination--
	r.HandleFunc("/get-all-prod-line-data-by-search-paginations", GetAllProdLineBySearchAndPagination).Methods("POST")
	// --Process Master Pagination Search--
	r.HandleFunc("/get-all-production-process-by-search-paginations", GetAllProcessBySearchAndPagination).Methods("POST")
	// --Pending Orders--
	r.HandleFunc("/get-order-pending-details-by-search-pagination", GetPendingOrderDetailsBySearchAndPagination).Methods("POST")

	// -- API: API Keys --
	// -- GET --
	r.HandleFunc("/get-all-api-keys", GetAllAPIKeys).Methods("GET")
	r.HandleFunc("/get-api-key-by-param", GetAPIKeyByParam).Methods("GET")

	// -- POST --
	r.HandleFunc("/add-or-update-api-key", AddOrUpdateAPIKey).Methods("POST")
	r.HandleFunc("/get-api-keys-by-search-pagination", GetAPIKeysWithSearchAndPagination).Methods("POST")

	slog.Info("Server listening on port: " + port)
	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var User m.User
	err := json.NewDecoder(r.Body).Decode(&User)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	userExistsErr := dao.CheckIfUserExists(&User, User.Email)

	if userExistsErr != nil && strings.Contains(userExistsErr.Error(), "user already exists") {
		log.Printf("Error creating end user, already exists: %v", err)
		http.Error(w, "Failed to create end user (101)", http.StatusForbidden)
		return
	}

	if User.VendorsCode != "" {
		// get vendor data by ID
		vendorData, err := dao.GetVendorByParam("vendor_code", User.VendorsCode)
		if err != nil || len(vendorData) == 0 {
			log.Printf("Error Getting vendor: %v", err)
			http.Error(w, "Failed to get vendor data (406)", http.StatusNotAcceptable)
			return
		}

		if err := dao.CreateNewOrUpdateExistingUser(&User); err != nil {
			log.Printf("Error creating end user: %v", err)
			http.Error(w, "Failed to create end user (102)", http.StatusInternalServerError)
			return
		}

		// get user data
		userData, err := dao.GetUserWithEmail(User.Email)
		if err != nil {
			log.Printf("Error creating end user: %v", err)
			http.Error(w, "Failed to create end user (102)", http.StatusInternalServerError)
			return
		}

		// map user to vendor
		var VendorToUser m.UserToVendor
		VendorToUser.UserId = int(userData.ID)
		VendorToUser.VendorId = vendorData[0].ID

		err = dao.CreateNewOrUpdateVendorToUser(&VendorToUser)
		if err != nil {
			log.Printf("Error creating end user: %v", err)
			http.Error(w, "Failed to create end user (102)", http.StatusInternalServerError)
			return
		}
	} else {
		role, _ := dao.GetUserRoleByRoleID(int(User.RoleID))
		if role.RoleName != "Customer" {
			if err := dao.CreateNewOrUpdateExistingUser(&User); err != nil {
				log.Printf("Error creating end user: %v", err)
				http.Error(w, "Failed to create end user (102)", http.StatusInternalServerError)
				return
			}
			// get user data
			userData, err := dao.GetUserWithEmail(User.Email)
			if err != nil {
				log.Printf("Error creating end user: %v", err)
				http.Error(w, "Failed to create end user (102)", http.StatusInternalServerError)
				return
			}
			//get vendor id
			vendorData, _ := dao.GetVendorDetailsByVendorCode(InventoryVendorCode)
			// map user to vendor
			var VendorToUser m.UserToVendor
			VendorToUser.UserId = int(userData.ID)
			VendorToUser.VendorId = vendorData.ID

			err = dao.CreateNewOrUpdateVendorToUser(&VendorToUser)
			if err != nil {
				log.Printf("Error creating end user: %v", err)
				http.Error(w, "Failed to create end user (102)", http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, "Failed to get vendor data (406)", http.StatusNotAcceptable)
			return
		}

	}
	token, expire, err := dao.GenerateToken(User) // with expiration time
	if err != nil {
		log.Printf("Error generating token: %v", err)
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	resp := m.ApiResp{
		Code:   http.StatusOK,
		Token:  token,
		Expire: expire.Format(time.RFC3339),
		User:   User,
		// Assuming you handle EndUserRoles and EndUsers appropriately here
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var credentials m.LoginReq
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate login credentials
	User, userRoles, validateErr := dao.CheckUserCredentials(credentials.Username, credentials.Password)
	if validateErr != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	if userRoles == nil {
		http.Error(w, "Your account is waiting for approval ,please contact the administrator", http.StatusUnauthorized)
		return
	}

	if !User.Isactive {
		http.Error(w, "Your account is disabled please contact the administrator", http.StatusForbidden)
		return
	}
	token, expire, err := dao.GenerateToken(User) // with expiration time
	if err != nil {
		log.Printf("Error generating token: %v", err)
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	resp := m.ApiResp{
		Code:      http.StatusOK,
		Token:     token,
		Expire:    expire.Format(time.RFC3339),
		User:      User,
		UserRoles: userRoles,
		// Assuming you handle EndUserRoles and EndUsers appropriately here
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// GetUserByEmailHandler handles the API request to get user info by email
func GetUserByEmailHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract the email from query params
	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	// Fetch user information based on email
	User, err := dao.GetUserWithEmail(email)
	if err != nil {
		log.Printf("Error fetching user with email %s: %v", email, err)
		http.Error(w, "Failed to get user information", http.StatusNotFound)
		return
	}
	// Assuming (m.User{}) is the zero value when no user is found; this may need adjustment
	if (m.User{}) == User {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(User); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func GetOtpCountHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	count, err := dao.GetOtpCount()
	if err != nil {
		http.Error(w, "Failed to get OTP count", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(map[string]int64{"count": count}); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func SaveOtpDetailsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var rec m.OtpDetail
	if err := json.NewDecoder(r.Body).Decode(&rec); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := dao.SaveOtpDetails(&rec); err != nil {
		log.Printf("Error saving OTP details: %v", err)
		http.Error(w, "Failed to save OTP details", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "OTP details saved successfully"})
}

func FindOtpDetailsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	mobileNo, ok := r.URL.Query()["mobileNo"]
	if !ok || len(mobileNo[0]) < 1 {
		http.Error(w, "Mobile number is required", http.StatusBadRequest)
		return
	}

	mobileNoFromReq, err := strconv.ParseUint(mobileNo[0], 10, 64)
	if err != nil {
		http.Error(w, "Invalid mobile number format", http.StatusBadRequest)
		return
	}

	otpDetail, err := dao.FindOtpDetailsForMobileNo(mobileNoFromReq)
	if err != nil {
		log.Printf("Error finding OTP details: %v", err)
		http.Error(w, "Failed to find OTP details", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(otpDetail); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// GetAllLogosHandler handles the incoming API request to fetch logos
func GetAllLogosHandler(w http.ResponseWriter, r *http.Request) {
	// Read the incoming request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// RequestData represents the incoming request data structure
	type RequestData struct {
		ProjectCode string `json:"project_code"`
	}

	// Unmarshal the JSON data into the RequestData structure
	var requestData RequestData
	err = json.Unmarshal(body, &requestData)
	if err != nil {
		log.Printf("Error unmarshaling request body: %v", err)
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	data, err := dao.GetAllLogosByName(requestData.ProjectCode)
	if err != nil {
		log.Printf("Error while getting data: %v", err)
		http.Error(w, "Invalid name", http.StatusBadRequest)
		return
	}
	response, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshaling API response: %v", err)
		http.Error(w, "Error processing response", http.StatusInternalServerError)
		return
	}

	// Send the response back to the client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
