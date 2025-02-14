package main

import (
	"bytes"
	// "encoding/base64"
	"encoding/base64"
	"encoding/json"
	"fmt"

	// "io"
	"log"
	"net/http"

	// "net/url"
	"os"
	"strconv"
)

var (
	apiUsername      = os.Getenv("API_USER")
	// accessToken      = "7BXtarGYcQaCiCeS706G9M83DxC1ZJux"
	accessToken      = os.Getenv("API_TOKEN")
    internalEndpoint = "http://" + os.Getenv("LISTMONK_URL") + "/api/subscribers/switch_list"
    hindiListId3m, _ = strconv.Atoi(os.Getenv("HINDI_LIST_3M"))
    englishListId3m, _ = strconv.Atoi(os.Getenv("ENGLISH_LIST_3M"))
    hindiListId1m, _    = strconv.Atoi(os.Getenv("HINDI_LIST_1M"))
    englishListId1m, _  = strconv.Atoi(os.Getenv("ENGLISH_LIST_1M"))
)

type RequestBody struct {
    Email string `json:"email"`
    List1 []int  `json:"lista"`
    List2 []int  `json:"listr"`
}

func sendPostRequest(body RequestBody) (*http.Response, error) {
    url := internalEndpoint
    auth := apiUsername + ":" + accessToken
    authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
    // Marshal the body into JSON
    requestBody, err := json.Marshal(body)
    if err != nil {
        return nil, fmt.Errorf("error marshalling request body: %v", err)
    }

    // Create a new POST request
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
    if err != nil {
        return nil, fmt.Errorf("error creating request: %v", err)
    }

    // Set appropriate headers
    req.Header.Set("Content-Type", "application/json")
    req.Header.Add("Authorization", authHeader)


    // Send the request using the default HTTP client
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("error sending request: %v", err)
    }

    fmt.Println("Sent request: ", req, resp)

    return resp, nil
}

type RequestData struct {
    Email     string `json:"email"`
	Language  string `json:"language"`
	Frequency string `json:"frequency"`
}

func proxyHandler_ChangeList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	var data RequestData
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

    var req RequestBody
    req.Email = data.Email

    if data.Language == "hindi"{
        if data.Frequency == "1"{
            req.List1 = append(req.List1, hindiListId1m)
        }else if data.Frequency == "3"{
            req.List1 = append(req.List1, hindiListId3m)
        }
    }else if data.Language == "english"{
        if data.Frequency == "1"{
            req.List1 = append(req.List1, englishListId1m)
        }else if data.Frequency == "3"{
            req.List1 = append(req.List1, englishListId3m)
        }
    }else if data.Language == "both"{
        if data.Frequency == "1"{
            req.List1 = append(req.List1, englishListId1m)
            req.List1 = append(req.List1, hindiListId1m)
        }else if data.Frequency == "3"{
            req.List1 = append(req.List1, englishListId3m)
            req.List1 = append(req.List2, hindiListId3m)
        }
    }

    fmt.Println(req.List1)

    sendPostRequest(req)

	// Simple response
	response := map[string]string{
		"message":   "Received data successfully",
		"language":  data.Language,
		"frequency": data.Frequency,
	}

	// Send JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

//
// func proxyHandler_HIN_TO_ENG(w http.ResponseWriter, r *http.Request) {
//     fmt.Println(r.URL.Query())
// 	email := r.URL.Query().Get("email")
//     fmt.Println(email + " from HIN to ENG")
//
// 	if email == "" {
// 		http.Error(w, "Email is required", http.StatusBadRequest)
// 		return
// 	}
//
// 	// Create the URL for the internal endpoint with the API key
// 	proxyURL, err := url.Parse(internalEndpoint)
// 	if err != nil {
// 		http.Error(w, "Invalid internal endpoint URL", http.StatusInternalServerError)
// 		return
// 	}
//
// 	query := proxyURL.Query()
// 	query.Set("email", email)
//     query.Set("lista", englishListId)
//     query.Set("listr", hindiListId)
//     query.Set("addBoth", "0")
// 	proxyURL.RawQuery = query.Encode()
//
//     sendRequest(proxyURL, w)
// }
//
// func proxyHandler_ENG_TO_HIN(w http.ResponseWriter, r *http.Request) {
//     fmt.Println(r.URL.Query())
// 	email := r.URL.Query().Get("email")
//     fmt.Println(email + " from ENG to HIN")
//
// 	if email == "" {
// 		http.Error(w, "Email is required", http.StatusBadRequest)
// 		return
// 	}
//
// 	// Create the URL for the internal endpoint with the API key
// 	proxyURL, err := url.Parse(internalEndpoint)
// 	if err != nil {
// 		http.Error(w, "Invalid internal endpoint URL", http.StatusInternalServerError)
// 		return
// 	}
//
// 	query := proxyURL.Query()
// 	query.Set("email", email)
//     query.Set("lista", hindiListId)
//     query.Set("listr", englishListId)
//     query.Set("addBoth", "0")
// 	proxyURL.RawQuery = query.Encode()
//
//     sendRequest(proxyURL, w)
// }
//
// func proxyHandler_BOTH(w http.ResponseWriter, r *http.Request) {
//     fmt.Println(r.URL.Query())
// 	email := r.URL.Query().Get("email")
//     fmt.Println(email + " from BOTH ENG & HIN")
//
// 	if email == "" {
// 		http.Error(w, "Email is required", http.StatusBadRequest)
// 		return
// 	}
//
// 	// Create the URL for the internal endpoint with the API key
// 	proxyURL, err := url.Parse(internalEndpoint)
// 	if err != nil {
// 		http.Error(w, "Invalid internal endpoint URL", http.StatusInternalServerError)
// 		return
// 	}
//
// 	query := proxyURL.Query()
// 	query.Set("email", email)
//     query.Set("lista", hindiListId)
//     query.Set("listr", englishListId)
//     query.Set("addBoth", "1")
// 	proxyURL.RawQuery = query.Encode()
//
//     sendRequest(proxyURL, w)
// }


func main() {
    http.Handle("/", http.FileServer(http.Dir("./static")))
    http.HandleFunc("/proxy/change_user_pref", proxyHandler_ChangeList)
	// http.HandleFunc("/proxy/switch_to_hindi", proxyHandler_ENG_TO_HIN)
	// http.HandleFunc("/proxy/switch_to_english", proxyHandler_HIN_TO_ENG)
	// http.HandleFunc("/proxy/use_both", proxyHandler_BOTH)

	fmt.Println("Proxy server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
