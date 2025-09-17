package server

import (
	"bytes"
	"io"
	"net/http"
)

func CreateNewKbTransaction(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	url := DBURL + "/create-new-KbTransaction"
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

func UpdateRunningNumberAfterTransactioin(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	url := DBURL + "/update-running-number"
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
	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Failed to read response body", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
