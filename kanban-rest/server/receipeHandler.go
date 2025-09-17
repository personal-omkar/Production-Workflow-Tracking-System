package server

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/url"

	m "irpl.com/kanban-commons/model"
	"irpl.com/kanban-commons/utils"
)

// Create Recipe Entry
func CreateRecipe(w http.ResponseWriter, r *http.Request) {

	var data []m.Recipe

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err == nil {
		jsonValue, err := json.Marshal(data)
		if err != nil {
			slog.Error("%s - error - %s", "Error marshaling cold storage data", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to marshal cold storage data")
			return
		}

		resp, err := http.Post(DBURL+"/create-new-recipe", "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			slog.Error("%s - error - %s", "Error making POST request", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to update cold storage data in inventorytable")
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

		utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to update record")
		slog.Error("%s - error - %s", "Failed to update cold storage data in inventory table", err.Error())
	}
}

// Update Recipe Entry
func UpdateRecipe(w http.ResponseWriter, r *http.Request) {

	var data []m.Recipe

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err == nil {
		jsonValue, err := json.Marshal(data)
		if err != nil {
			slog.Error("%s - error - %s", "Error marshaling cold storage data", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to marshal cold storage data")
			return
		}

		resp, err := http.Post(DBURL+"/update-existing-recipe", "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			slog.Error("%s - error - %s", "Error making POST request", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to update cold storage data in inventorytable")
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

		utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to update record")
		slog.Error("%s - error - %s", "Failed to update cold storage data in inventory table", err.Error())
	}
}

// GetAllRecipe retrives all recipe
func GetAllRecipe(w http.ResponseWriter, r *http.Request) {

	resp, err := http.Get(DBURL + "/get-all-recipe")
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

// GetRecipeByDataKey retrives  recipe based on the datakey value
func GetRecipeByDataKey(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	resp, err := http.Get(DBURL + "/get-recipe-by-data-key?key=" + key)
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

// GetRecipeByDataKey retrives  recipe based on the datakey value
func GetRecipeByDataValue(w http.ResponseWriter, r *http.Request) {
	value := r.URL.Query().Get("value")
	resp, err := http.Get(DBURL + "get-recipe-by-data-value?value=" + value)
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

// GetRecipeByDataKey retrives  recipe based on the datakey value
func GetRecipeByDataKeyAndValue(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")
	resp, err := http.Get(DBURL + "/get-recipe-by-data-key-and-value?key=" + key + "&value=" + value)
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

// Delete recipe
func DeleteRecipe(w http.ResponseWriter, r *http.Request) {
	var recipe m.Recipe
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&recipe); err == nil {
		jsonValue, err := json.Marshal(recipe)
		if err != nil {
			slog.Error("%s - error - %s", "Error marshaling cold storage data", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to marshal cold storage data")
			return
		}

		resp, err := http.Post(DBURL+"/delete-recipe-by-id", "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			slog.Error("%s - error - %s", "Error making POST request", err)
			utils.SetResponse(w, http.StatusInternalServerError, "Failed to update cold storage data in inventorytable")
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
	} else {

		utils.SetResponse(w, http.StatusInternalServerError, "Fail: failed to update record")
		slog.Error("%s - error - %s", "Failed to update cold storage data in inventory table", err.Error())
	}
}

// Get Recipe by Parameter
func GetRecipeByParam(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")
	resp, err := http.Get(DBURL + "/get-recipe-by-param?key=" + key + "&value=" + url.QueryEscape(value))
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
