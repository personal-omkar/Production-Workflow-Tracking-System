package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"irpl.com/kanban-commons/utils"
)

func GetAllDefaultsHandler(ProjectName string) (map[string]string, error) {
	requestData := map[string]string{
		"project_code": ProjectName,
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		log.Printf("Error marshaling request data: %v", err)
		return nil, err
	}

	// Send the POST request
	url := utils.RestURL + "/get-all-logos"
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error sending request: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return nil, err
	}

	// Check for a non-200 response code
	if resp.StatusCode != http.StatusOK {
		log.Printf("Non-OK HTTP status: %v", resp.Status)
		return nil, fmt.Errorf("non-OK HTTP status: %v", resp.Status)
	}

	// Unmarshal the response body into a map
	var responseData map[string]string
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		log.Printf("Error unmarshaling response body: %v", err)
		return nil, err
	}

	// Return the fetched logos (or whatever the response data is)
	return responseData, nil
}
