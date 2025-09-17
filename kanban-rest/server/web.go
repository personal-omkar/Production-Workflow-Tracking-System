// server\web.go

package server

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"

	// "time"
	m "irpl.com/kanban-commons/model"
	"irpl.com/kanban-commons/utils"
	restUtils "irpl.com/kanban-rest/services"

	"github.com/gorilla/mux"
)

const DefaultDBHelperHost string = "0.0.0.0" // Default port if not set in env
const DefaultDBHelperPort string = "4100"    // Default port if not set in env

var DBHelperHost string // Global variable to hold the DB helper host
var DBHelperPort string // Global variable to hold the DB helper port
var DBURL string        // Global variable to hold the DB url

func init() {
	// Initialize DBHelperHost and DBHelperPort with value from environment variable or fallback to default
	DBHelperHost = os.Getenv("DBHELPER_HOST")
	if strings.TrimSpace(DBHelperHost) == "" {
		DBHelperHost = DefaultDBHelperHost
	}

	DBHelperPort = os.Getenv("DBHELPER_PORT")
	if strings.TrimSpace(DBHelperPort) == "" {
		DBHelperPort = DefaultDBHelperPort
	}

	DBURL = utils.JoinStr("http://", DBHelperHost, ":", DBHelperPort)
}

// Status ss
func Status(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func Web() {
	r := mux.NewRouter()

	port := os.Getenv("RESTSRV_PORT")
	if port == "" {
		slog.Info("RESTSRV_PORT environment variable not set, defaulting to :4200")
		port = "4200"
	}

	// r.Use(websecure.CommonMiddleware)
	r.HandleFunc("/status", Status)
	r.HandleFunc("/register", RegisterHandler).Methods("POST")
	r.HandleFunc("/login", LoginHandler).Methods("POST")
	r.HandleFunc("/get-user-by-mobile", GetUserByMobileNoHandler).Methods("GET")
	r.HandleFunc("/get-user-by-email", GetUserByEmailHandler).Methods("GET")
	r.HandleFunc("/get-user-role-by-id", GetUserRoleByRoleIDHandler).Methods("GET")
	r.HandleFunc("/RawQuery", RawQuery).Methods("POST")

	// -- API: UserRoles --
	r.HandleFunc("/create-user-role", CreateUserRole).Methods("POST")
	r.HandleFunc("/update-user-role", UpdateUserRole).Methods("POST")
	r.HandleFunc("/delete-user-role", DeleteUserRole).Methods("POST")
	r.HandleFunc("/get-role-by-name", GetUserRoleByName).Methods("GET")
	r.HandleFunc("/get-role-by-id", GetUserRoleById).Methods("GET")
	r.HandleFunc("/get-all-user-roles", GetAllUserRoles).Methods("GET")
	// -- API: Get Porduction line specific data for prod line tab --
	// -- GET --
	r.HandleFunc("/get-production-line-items", GetProductionLinesHandler).Methods("GET")

	// -- API: Add Production Line --
	// -- POST --
	r.HandleFunc("/add-production-line", AddProdLine).Methods("POST")
	r.HandleFunc("/delete-production-line-by-id", DeleteProductionLineById).Methods("DELETE")
	r.HandleFunc("/delete-production-line-celldata-by-productionline-id", DeleteProductionLineCellDataByProductionLineId).Methods("DELETE")
	r.HandleFunc("/delete-production-line-by-prodline-id", DeleteProductionLineByProdLineId).Methods("DELETE")
	r.HandleFunc("/delete-production-line-cell", DeleteProductionLineCell).Methods("DELETE")
	// -- API: Get sepcific prod line stauts data --
	// -- POST --
	r.HandleFunc("/get-production-line-status", GetProductionLineStaus).Methods("POST")

	// -- API: Update kbroot running no--
	// -- POST --
	r.HandleFunc("/update-running-numbers", UpdateRunningNumbers).Methods("POST")

	// -- API: Vendors --
	// -- GET --
	r.HandleFunc("/get-kanban-for-all-vendor", GetKanbanForAllVendor).Methods("GET")
	// -- POST --
	r.HandleFunc("/create-new-vendor", CreateVendors).Methods("POST")
	r.HandleFunc("/get-all-vendors-data", GetVendorsAllrecords).Methods("POST")
	r.HandleFunc("/get-kanban-for-vendor", GetKanbanForVendor).Methods("POST")
	r.HandleFunc("/get-all-vendor-by-search-pagination", GetVendorSearchPaginationRecords).Methods("POST")
	// r.HandleFunc("/create-new-vendor", CreateVendors).Methods("POST")

	// -- API: Prod line --
	// -- GET --
	r.HandleFunc("/get-all-prod-line-data", GetProdLineAllrecords).Methods("GET")
	r.HandleFunc("/create-new-or-update-existing-production-line", CreateNewOrUpdateExistingProductionLine).Methods("POST")

	// -- API: Compounds --
	// -- GET --
	r.HandleFunc("/get-all-compound-data", GetCompountsAllrecords).Methods("GET")
	r.HandleFunc("/get-all-active-compounds", GetAllActiveCompounds).Methods("GET")
	r.HandleFunc("/get-compound-data-by-parm", GetCompoundsByParm).Methods("GET")

	r.HandleFunc("/get-compound-data-by-vendor", GetCompoundsByVendors).Methods("POST")
	r.HandleFunc("/get-packing-compound-data-by-vendor", GetPackingCompoundsByVendors).Methods("POST")
	r.HandleFunc("/get-quality-compound-data-by-vendor", GetQualityCompoundsByVendors).Methods("POST")
	// -- POST --
	r.HandleFunc("/add-compound-data-by-vendor", AddCompoundsForVendor).Methods("POST")
	r.HandleFunc("/add-compound-data-to-production-line", AddCompoundsInProductionLine).Methods("POST")
	r.HandleFunc("/update-compound-status-to-packing", UpdateCompoundStatusToPacking).Methods("POST")
	r.HandleFunc("/update-compound-quality-status-to-reject", UpdateCompoundQualityStatusToReject).Methods("POST")
	r.HandleFunc("/update-compound-status-to-dispatch", UpdateCompoundStatusToDispatch).Methods("POST")
	r.HandleFunc("/get-all-compounds-by-search-pagination", GetAllCompoundsBySearchAndPagination).Methods("POST")

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

	// -- API: Order --
	// -- POST --
	r.HandleFunc("/create-new-order-entry", CreateNewOrderEntry).Methods("POST")
	r.HandleFunc("/create-multi-new-order-entry", CreateMultiNewOrderEntry).Methods("POST")
	r.HandleFunc("/get-order-details", GetOrderDEtails).Methods("POST")
	r.HandleFunc("/update-order-status", UpdateOrderStatus).Methods("POST")
	r.HandleFunc("/get-customer-order-details", GetCustomerOrderDetails).Methods("POST")
	// API  to fetch orders details for order page
	r.HandleFunc("/OrderDetailsForCustomer", GetOrderDetaislForOrderPage).Methods("POST")
	r.HandleFunc("/get-all-orders-by-vendor-code", GetAllOrderByVendorCode).Methods("POST")
	r.HandleFunc("/get-vendor-details-by-vendor-code", GetVendorDetailsByVendorCode).Methods("POST")

	r.HandleFunc("/fetch-all-production-process-data", GetAllProductionProcessData).Methods("GET")
	r.HandleFunc("/get-all-production-process-for-line", GetAllProductionProcessEntries).Methods("GET")
	r.HandleFunc("/get-all-production-process", GetAllProductionProcess).Methods("GET")
	r.HandleFunc("/get-prod-line-by-param", GetProdLineByParam).Methods("POST")
	r.HandleFunc("/get-prod-process-by-param", GetProdProcessByParam).Methods("POST")
	r.HandleFunc("/edit-prod-process", EditProdProcess).Methods("POST")
	r.HandleFunc("/add-production-process", AddProdProcess).Methods("POST")
	r.HandleFunc("/edit-prod-line", EditProdLine).Methods("POST")
	// -- API: Inventory --
	// -- GET --
	r.HandleFunc("/get-all-inventory-data", GetAllInventoryrecords).Methods("GET")
	r.HandleFunc("/get-all-coldstorage-data", GetAllColdStoragerecords).Methods("GET")
	// r.HandleFunc("/get-inventory-by-param", GetInventoryByParam).Methods("GET")
	// r.HandleFunc("/delete-inventory-by-param", DeleteInventoryByParam).Methods("GET")
	// -- POST --
	r.HandleFunc("/update-coldstorage-quantity", UpdateColdStorageQuantity).Methods("POST")
	r.HandleFunc("/create-new-or-update-existing-inventory", CreateOrUpdateInventory).Methods("POST")
	r.HandleFunc("/delete-inventory-by-id", DeleteInventoryById).Methods("POST")
	r.HandleFunc("/get-inventory-by-param", GetInventoryByParam).Methods("GET")
	r.HandleFunc("/get-inventory-by-search", GetInventoryBySearch).Methods("POST")
	r.HandleFunc("/get-all-cold-storage-by-search-pagination", ColdStorageSearchPagination).Methods("POST")

	r.HandleFunc("/create-new-KbTransaction", CreateNewKbTransaction).Methods("POST")

	// -- API: UserToVemdor --
	// -- GET --
	r.HandleFunc("/get-vendor-by-userid", GetVendorByUserID).Methods("GET")

	r.HandleFunc("/update-running-number", UpdateRunningNumberAfterTransactioin).Methods("POST")

	r.HandleFunc("/delete-order-entry", DeleteOrder).Methods("DELETE")
	// -- API: User Details --
	// -- GET --
	r.HandleFunc("/get-user-details", GetUserDetails).Methods("GET")
	r.HandleFunc("/get-user-details-by-email", GetUserDetailsByEmail).Methods("GET")
	r.HandleFunc("/update-user-status", UpdateUserStatus).Methods("POST")
	r.HandleFunc("/update-user-details", UpdateUser).Methods("POST")
	r.HandleFunc("/create-new-user", CreateUser).Methods("POST")
	r.HandleFunc("/delete-user", DeleteUser).Methods("POST")
	r.HandleFunc("/get-all-kbRoot-details", GetAllCompletedKBRootDetails).Methods("POST")
	r.HandleFunc("/get-all-kbRoot-details-by-search", GetCompletedKBRootDetailsBySearch).Methods("POST")
	r.HandleFunc("/get-kbRoot-details", GetDetailRootData).Methods("POST")

	r.HandleFunc("/get-all-logos", GetAllLogosHandler).Methods("POST")

	// -- API: LDAPConfig --
	// -- POST --
	r.HandleFunc("/create-ldap-config", CreateLDAPConfig).Methods("POST")
	r.HandleFunc("/update-ldap-config", UpdateLDAPConfig).Methods("POST")
	r.HandleFunc("/delete-ldap-config", DeleteLDAPConfig).Methods("POST")
	// -- GET --
	r.HandleFunc("/get-default-ldap-config", GetDefaultLDAPConfig).Methods("GET")

	// -- API: SambaConfig --
	// -- POST --
	r.HandleFunc("/create-samba-config", CreateSambaConfig).Methods("POST")
	r.HandleFunc("/update-samba-config", UpdateSambaConfig).Methods("POST")
	r.HandleFunc("/delete-samba-config", DeleteSambaConfig).Methods("POST")
	// -- GET --
	r.HandleFunc("/get-default-samba-config", GetDefaultSambaConfig).Methods("GET")

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

	r.HandleFunc("/get-all-details-for-order", GetOrderDetailsForHistory)

	r.HandleFunc("/delete-kbroot-by-ids", DeleteKbRootByIDsHandler).Methods("POST")
	r.HandleFunc("/get-lineup-processes-by-lineid", GetLinedUpProductionProcessByLineId).Methods("POST")

	r.HandleFunc("/get-all-kanban-details-for-report", GetAllKanbanDetailsForReport).Methods("POST")

	// -- API:- Stage
	r.HandleFunc("/create-new-stage", CreateStage).Methods("POST")
	r.HandleFunc("/update-stage", UpdateExistingStage).Methods("POST")
	r.HandleFunc("/delete-stage", DeleteStageByID).Methods("POST")
	r.HandleFunc("/get-stages-by-header", GetStagesByHeader).Methods("POST")
	r.HandleFunc("/get-stage-by-param", GetStagesByParam).Methods("GET")
	r.HandleFunc("/get-all-stages", GetAllStage).Methods("GET")
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
	r.HandleFunc("/create-new-or-update-existing-material", CreateNewOrUpdateExistingRawMaterial).Methods("POST")
	r.HandleFunc("/get-all-material-by-search-pagination", GetAllMaterialBySearchAndPagination).Methods("POST")
	// --GET
	r.HandleFunc("/get-all-material", GetAllMaterial).Methods("GET")
	r.HandleFunc("/get-material-by-param", GetMaterialByParam).Methods("GET")

	//--Chemicals--
	// -- POST --
	r.HandleFunc("/create-new-or-update-existing-chemical", CreateNewOrUpdateExistingRawChemical).Methods("POST")
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
	r.HandleFunc("/add-api-key", AddOrUpdateAPIKey).Methods("POST")
	r.HandleFunc("/get-api-keys-by-search-pagination", GetAPIKeysWithSearchAndPagination).Methods("POST")

	slog.Info("Server listening on port: " + port)

	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// TODO- call db helper APIs using tokens

func RawQuery(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var data m.RawQuery
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Println("RawQuery: error- ", err.Error())
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	jsonValue, err := json.Marshal(data)
	if err != nil {
		log.Printf("RawQuery: Error marshaling user data: %v", err)
		http.Error(w, "Failed to marshal user data", http.StatusInternalServerError)
		return
	}

	resp, err := http.Post("http://"+DBHelperHost+":"+DBHelperPort+"/RawQuery", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Printf("RawQuery: Error making POST request: %v", err)
		http.Error(w, "Failed to execute raw query", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("RawQuery: Error reading response body: %v", err)
		http.Error(w, "Failed to read response body", http.StatusInternalServerError)
		return
	}

	w.Write(responseBody)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var User m.User
	err := json.NewDecoder(r.Body).Decode(&User)
	if err != nil {
		log.Println("RegisterHandler: error- ", err.Error())
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	jsonValue, err := json.Marshal(User)
	if err != nil {
		log.Printf("Error marshaling user data: %v", err)
		http.Error(w, "Failed to marshal user data", http.StatusInternalServerError)
		return
	}

	resp, err := http.Post("http://"+DBHelperHost+":"+DBHelperPort+"/create-new-user", "application/json", bytes.NewBuffer(jsonValue)) // Changed port to 3100
	if err != nil {
		log.Printf("Error making POST request: %v", err)
		http.Error(w, "Failed to create user", http.StatusForbidden)
		return
	}
	defer resp.Body.Close()

	restUtils.RespondFailure(resp, w)

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		http.Error(w, "Failed to read response body", http.StatusInternalServerError)
		return
	}

	w.Write(responseBody)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// log.Println("LoginHandler: request received")

	var credentials m.LoginReq
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	jsonValue, err := json.Marshal(credentials)
	if err != nil {
		log.Printf("Error marshaling login credentials: %v", err)
		http.Error(w, "Failed to marshal login credentials", http.StatusBadRequest)
		return
	}

	// log.Println("LoginHandler: API http://"+DBHelperHost+":"+DBHelperPort+"/login called")

	resp, err := http.Post("http://"+DBHelperHost+":"+DBHelperPort+"/login", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Printf("Error making POST request: %v", err)
		http.Error(w, "Failed to login", http.StatusForbidden)
		return
	}
	defer resp.Body.Close()

	// log.Println("LoginHandler: http://"+DBHelperHost+":"+DBHelperPort+"/login response received")

	restUtils.RespondFailure(resp, w)

	// Proceed with handling a successful response (200 OK or 201 Created)
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		http.Error(w, "Failed to read response body", http.StatusInternalServerError)
		return
	}

	// Write the successful response back to the client
	w.Write(responseBody)
}

func GetUserByMobileNoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// log.Println("GetUserByMobileNoHandler: request received")

	mobileNo := r.URL.Query().Get("mobile_no")
	if len(strings.TrimSpace(mobileNo)) <= 0 {
		http.Error(w, "Mobile number is required", http.StatusBadRequest)
		return
	}

	// log.Println("GetUserByMobileNoHandler: mobileNo", mobileNo)

	// log.Println("GetUserByMobileNoHandler: http://"+DBHelperHost+":"+DBHelperPort+"/get-user-by-mobile request sent", mobileNo)

	// Use net/url to build the query parameters
	queryParams := url.Values{}
	queryParams.Add("mobile_no", mobileNo)

	// Make HTTP POST request
	resp, err := http.Get("http://" + DBHelperHost + ":" + DBHelperPort + "/get-user-by-mobile?" + queryParams.Encode())
	if err != nil {
		log.Printf("Error making POST request: %v", err)
		http.Error(w, "Failed to get user information", http.StatusForbidden)
		return
	}
	defer resp.Body.Close()

	restUtils.RespondFailure(resp, w)

	// Process response
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		http.Error(w, "Failed to read response body", http.StatusInternalServerError)
		return
	}

	w.Write(responseBody)
}

func GetUserByEmailHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// log.Println("GetUserByEmailHandler: request received")

	email := r.URL.Query().Get("email")
	if len(strings.TrimSpace(email)) <= 0 {
		http.Error(w, "Mobile number is required", http.StatusBadRequest)
		return
	}

	// log.Println("GetUserByMobileNoHandler: email", email)

	// log.Println("GetUserByMobileNoHandler: http://"+DBHelperHost+":"+DBHelperPort+"/get-user-by-email request sent", email)

	// Use net/url to build the query parameters
	queryParams := url.Values{}
	queryParams.Add("email", email)

	// Make HTTP Get request
	resp, err := http.Get("http://" + DBHelperHost + ":" + DBHelperPort + "/get-user-by-email?" + queryParams.Encode())
	if err != nil {
		log.Printf("Error making POST request: %v", err)
		http.Error(w, "Failed to get user information", http.StatusForbidden)
		return
	}
	defer resp.Body.Close()

	restUtils.RespondFailure(resp, w)

	// log.Println("GetUserByMobileNoHandler: http://"+DBHelperHost+":"+DBHelperPort+"/get-user-by-mobile response", resp.Body)

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
	resp, err := http.Get("http://" + DBHelperHost + ":" + DBHelperPort + "/get-user-role-by-id?" + queryParams.Encode())
	if err != nil {
		log.Printf("Error making POST request: %v", err)
		http.Error(w, "Failed to get user information", http.StatusForbidden)
		return
	}
	defer resp.Body.Close()

	restUtils.RespondFailure(resp, w)

	// Process response
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		http.Error(w, "Failed to read response body", http.StatusInternalServerError)
		return
	}

	w.Write(responseBody)
}

func GetAllLogosHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	url := DBURL + "/get-all-logos"
	client := &http.Client{}
	req, err := http.NewRequest(r.Method, url, r.Body)
	if err != nil {
		http.Error(w, "Error creating request to DAO service", http.StatusInternalServerError)
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to send request to the target service", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Failed to process request on target service", http.StatusInternalServerError)
		return
	}
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("%s - error - %s", "Error reading response body", err)
		utils.SetResponse(w, http.StatusInternalServerError, "Failed to read response body")
		return
	}

	utils.SetResponse(w, http.StatusOK, string(responseBody))

}
