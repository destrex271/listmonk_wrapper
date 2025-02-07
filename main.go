package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

const (
	apiUsername      = "test"
	accessToken      = "77ql14YHnr1eAKgQzRrOJU8O3mcWacXe"
    internalEndpoint = "http://localhost:9000/api/subscribers/switch_to_hindi"
)

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	// Create the URL for the internal endpoint with the API key
	proxyURL, err := url.Parse(internalEndpoint)
	if err != nil {
		http.Error(w, "Invalid internal endpoint URL", http.StatusInternalServerError)
		return
	}

	query := proxyURL.Query()
	query.Set("email", email)
	proxyURL.RawQuery = query.Encode()

	// Create the request to the internal endpoint
	req, err := http.NewRequest("GET", proxyURL.String(), nil)
	if err != nil {
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}

	// Add the Authorization header
	auth := apiUsername + ":" + accessToken
	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Add("Authorization", authHeader)

	// Forward the request to the internal endpoint
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Error forwarding request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Copy the response from the internal endpoint to the client
	w.WriteHeader(resp.StatusCode)
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		http.Error(w, "Error reading response", http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/proxy/switch_to_hindi", proxyHandler)
	fmt.Println("Proxy server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
