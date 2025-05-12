package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Preference Change API
func SendPreferenceChangeRequest(internalEndpoint string, apiUsername string, accessToken string, body RequestBody, w http.ResponseWriter) {
	fmt.Println("Sending POST REQ to ", internalEndpoint)
	url := internalEndpoint
	auth := apiUsername + ":" + accessToken
	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
	// Marshal the body into JSON
	requestBody, err := json.Marshal(body)
	if err != nil {
		http.Error(w, "Error marshalling request body", http.StatusInternalServerError)
		return
	}

	// Create a new POST request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}

	// Set appropriate headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", authHeader)

	fmt.Println(req)

	// Send the request using the default HTTP client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Error sending request", http.StatusInternalServerError)
		return
	}
	fmt.Println(resp)
	defer resp.Body.Close()

	// Copy the response from the internal endpoint to the client
	w.WriteHeader(resp.StatusCode)
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		http.Error(w, "Error reading response", http.StatusInternalServerError)
	}
}
