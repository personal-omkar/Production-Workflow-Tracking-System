package services

import (
	"io"
	"log"
	"net/http"
)

// Respond,ailure if http status code is not equal to 200 or 201, then return failure
func RespondFailure(resp *http.Response, w http.ResponseWriter) {
	// Check the status code of the response
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		// log.Printf("Received non-success status code: %d", resp.StatusCode)

		// Read the response body from the unsuccessful request
		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading response body: %v", err)
			http.Error(w, "Failed to get requested data", resp.StatusCode)
			return
		}

		// Write the received status code and response body back to the client
		w.WriteHeader(resp.StatusCode)
		w.Write(responseBody)
		return
	}
}
