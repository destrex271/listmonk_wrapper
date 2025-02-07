package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

var (
	apiUsername      = os.Getenv("API_USER")
	// accessToken      = "7BXtarGYcQaCiCeS706G9M83DxC1ZJux"
	accessToken      = os.Getenv("API_TOKEN")
    internalEndpoint = "http://0.0.0.0/api/subscribers/switch_list"
    hindiListId      = os.Getenv("HINDI_LIST")
    englishListId    = os.Getenv("ENGLISH_LIST")
)

func sendRequest(proxyURL *url.URL, w http.ResponseWriter) {
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

func proxyHandler_HIN_TO_ENG(w http.ResponseWriter, r *http.Request) {
    fmt.Println(r.URL.Query())
	email := r.URL.Query().Get("email")
    fmt.Println(email + " from HIN to ENG")

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
    query.Set("lista", hindiListId)
    query.Set("listr", englishListId)
    query.Set("addBoth", "0")
	proxyURL.RawQuery = query.Encode()

    sendRequest(proxyURL, w)
}

func proxyHandler_ENG_TO_HIN(w http.ResponseWriter, r *http.Request) {
    fmt.Println(r.URL.Query())
	email := r.URL.Query().Get("email")
    fmt.Println(email + " from ENG to HIN")

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
    query.Set("lista", englishListId)
    query.Set("listr", hindiListId)
    query.Set("addBoth", "0")
	proxyURL.RawQuery = query.Encode()

    sendRequest(proxyURL, w)
}

func proxyHandler_BOTH(w http.ResponseWriter, r *http.Request) {
    fmt.Println(r.URL.Query())
	email := r.URL.Query().Get("email")
    fmt.Println(email + " from BOTH ENG & HIN")

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
    query.Set("lista", hindiListId)
    query.Set("listr", englishListId)
    query.Set("addBoth", "1")
	proxyURL.RawQuery = query.Encode()

    sendRequest(proxyURL, w)
}


func main() {
	http.HandleFunc("/proxy/switch_to_hindi", proxyHandler_ENG_TO_HIN)
	http.HandleFunc("/proxy/switch_to_english", proxyHandler_HIN_TO_ENG)
	http.HandleFunc("/proxy/use_both", proxyHandler_BOTH)

	fmt.Println("Proxy server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
