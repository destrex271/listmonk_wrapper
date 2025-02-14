package main

import (
	"bytes"
	"io"
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
}

func sendPostRequest(body RequestBody, w http.ResponseWriter) {
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
            req.List1 = append(req.List1, hindiListId3m)
        }
    }

    fmt.Println(req.List1)

   sendPostRequest(req, w)
}

func main() {
    http.Handle("/", http.FileServer(http.Dir("./static")))
    http.HandleFunc("/proxy/change_user_pref", proxyHandler_ChangeList)

	fmt.Println("Proxy server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
