package server

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"log/slog"
	"net/http"

	m "irpl.com/kanban-commons/model"
	"irpl.com/kanban-commons/utils"
)

// Get All Compounds
func GetCompountsAllrecords(w http.ResponseWriter, r *http.Request) {

	resp, err := http.Get(DBURL + "/get-all-compound-data")
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

	utils.SetResponse(w, http.StatusOK, string(responseBody))
}

// Get All Compounds
func GetAllActiveCompounds(w http.ResponseWriter, r *http.Request) {

	resp, err := http.Get(DBURL + "/get-all-active-compounds")
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

func GetCompoundsByVendors(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")
	w.Header().Set("Content-Type", "application/json")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	resp, err := http.Post(DBURL+"/get-compound-data-by-vendor?key="+key+"&value="+value, "application/json", bytes.NewBuffer(body))
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

	utils.SetResponse(w, http.StatusOK, string(responseBody))
}

func GetPackingCompoundsByVendors(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")
	w.Header().Set("Content-Type", "application/json")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	resp, err := http.Post(DBURL+"/get-packing-compound-data-by-vendor?key="+key+"&value="+value, "application/json", bytes.NewBuffer(body))
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

// Add Compounds data based on vendor records
func AddCompoundsForVendor(w http.ResponseWriter, r *http.Request) {
	var data m.AddCompoundsByVendor
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Println("Invalid request body: error- ", err.Error())
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	jsonValue, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error: Error marshaling user data: %v", err)
		http.Error(w, "Failed to marshal data", http.StatusInternalServerError)
		return
	}

	resp, err := http.Post("http://"+DBHelperHost+":"+DBHelperPort+"/add-compound-data-by-vendor", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Printf("Error making POST request: %v", err)
		http.Error(w, "Error making POST request:", http.StatusInternalServerError)
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

// Add compound to production line
func AddCompoundsInProductionLine(w http.ResponseWriter, r *http.Request) {
	type compoundData struct {
		LineID   string
		KbDataId []string
		UserID   string
	}
	var data compoundData
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Println("Invalid request body: error- ", err.Error())
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	jsonValue, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error: Error marshaling compound data: %v", err)
		http.Error(w, "Failed to marshal data", http.StatusInternalServerError)
		return
	}

	resp, err := http.Post("http://"+DBHelperHost+":"+DBHelperPort+"/add-compound-data-to-production-line", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Printf("Error making POST request: %v", err)
		http.Error(w, "Error making POST request:", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		http.Error(w, "Failed to read response body", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(resp.StatusCode)
	w.Write(responseBody)
}

func UpdateCompoundStatusToDispatch(w http.ResponseWriter, r *http.Request) {
	type compoundData struct {
		Notes    string
		KbDataId []string
		UserID   string
	}
	var data compoundData
	err := json.NewDecoder(r.Body).Decode(&data)

	if err != nil {
		log.Println("Invalid request body: error- ", err.Error())
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	jsonValue, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error: Error marshaling compound data: %v", err)
		http.Error(w, "Failed to marshal data", http.StatusInternalServerError)
		return
	}

	resp, err := http.Post("http://"+DBHelperHost+":"+DBHelperPort+"/update-compound-status-to-dispatch", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Printf("Error making POST request: %v", err)
		http.Error(w, "Error making POST request:", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		http.Error(w, "Failed to read response body", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(resp.StatusCode)
	w.Write(responseBody)
}

func AddorUpdateCompound(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	url := DBURL + "/add-update-compound"
	res, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		http.Error(w, "Error creating request to REST service", http.StatusInternalServerError)
		return
	}
	res.Header = r.Header
	client := &http.Client{}
	resp, err := client.Do(res)
	if err != nil {
		http.Error(w, "Failed to send request to REST service", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// If the response is not successful, return an error
	if resp.StatusCode != http.StatusCreated {
		http.Error(w, "Failed to read response body", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		http.Error(w, "Failed to forward response body", http.StatusInternalServerError)
	}
}

func GetQualityCompoundsByVendors(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")
	w.Header().Set("Content-Type", "application/json")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	resp, err := http.Post(DBURL+"/get-quality-compound-data-by-vendor?key="+key+"&value="+value, "application/json", bytes.NewBuffer(body))
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

func UpdateCompoundStatusToPacking(w http.ResponseWriter, r *http.Request) {
	type compoundData struct {
		Notes    string
		KbDataId []string
		UserID   string
	}
	var data compoundData
	err := json.NewDecoder(r.Body).Decode(&data)

	if err != nil {
		log.Println("Invalid request body: error- ", err.Error())
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	jsonValue, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error: Error marshaling compound data: %v", err)
		http.Error(w, "Failed to marshal data", http.StatusInternalServerError)
		return
	}

	resp, err := http.Post("http://"+DBHelperHost+":"+DBHelperPort+"/update-compound-status-to-packing", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Printf("Error making POST request: %v", err)
		http.Error(w, "Error making POST request:", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		http.Error(w, "Failed to read response body", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(resp.StatusCode)
	w.Write(responseBody)
}

func UpdateCompoundQualityStatusToReject(w http.ResponseWriter, r *http.Request) {
	type compoundData struct {
		Notes    string
		KbDataId []string
		UserID   string
	}
	var data compoundData
	err := json.NewDecoder(r.Body).Decode(&data)

	if err != nil {
		log.Println("Invalid request body: error- ", err.Error())
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	jsonValue, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error: Error marshaling compound data: %v", err)
		http.Error(w, "Failed to marshal data", http.StatusInternalServerError)
		return
	}

	resp, err := http.Post("http://"+DBHelperHost+":"+DBHelperPort+"/update-compound-quality-status-to-reject", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Printf("Error making POST request: %v", err)
		http.Error(w, "Error making POST request:", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		http.Error(w, "Failed to read response body", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(resp.StatusCode)
	w.Write(responseBody)
}

func GetCompoundsByParm(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")
	resp, err := http.Get(DBURL + "/get-compound-data-by-parm?key=" + key + "&value=" + value)
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

	utils.SetResponse(w, http.StatusOK, string(responseBody))
}

// GetAllCompoundsBySearchAndPagination returns a all records present in compounds table
func GetAllCompoundsBySearchAndPagination(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	resp, err := http.Post(DBURL+"/get-all-compounds-by-search-pagination", "application/json", bytes.NewBuffer(body))
	if err != nil {
		slog.Error("%s - error - %s", "Error making POST request", err)
		utils.SetResponse(w, http.StatusInternalServerError, "Failed to fectch order details")
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
