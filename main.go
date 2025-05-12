package main

import (

	// "encoding/base64"

	"encoding/json"
	"fmt"

	// "io"
	"log"
	"net/http"

	// "net/url"
	"os"
	"strconv"

	. "github.com/destrex271/listmonk_proxy/utils"
)

var (
	apiUsername = os.Getenv("API_USER")
	// accessToken      = "7BXtarGYcQaCiCeS706G9M83DxC1ZJux"
	accessToken        = os.Getenv("API_TOKEN")
	internalEndpoint   = "http://" + os.Getenv("LISTMONK_URL") + "/api/subscribers/switch_list"
	campaignEndpoint   = "http://" + os.Getenv("LISTMONK_URL") + "/api/campaigns"
	listEndpoint       = "http://" + os.Getenv("LISTMONK_URL") + "/api/lists"
	membershipEndpoint = "http://" + os.Getenv("LISTMONK_URL") + "/api/subscribers/lists"
	hindiListId3m, _   = strconv.Atoi(os.Getenv("HINDI_LIST_3M"))
	englishListId3m, _ = strconv.Atoi(os.Getenv("ENGLISH_LIST_3M"))
	hindiListId1m, _   = strconv.Atoi(os.Getenv("HINDI_LIST_1M"))
	englishListId1m, _ = strconv.Atoi(os.Getenv("ENGLISH_LIST_1M"))
)

func proxyHandler_RoutingMessenger(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var data Postback
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid request from JSON", http.StatusBadRequest)
		return
	}

	var smtp_verified []Recipient
	var smtp_unverified []Recipient

	for _, sub := range data.Recipients {
		if sub.Attribs["verification_status"] == true {
			smtp_verified = append(smtp_verified, sub)
		} else {
			smtp_unverified = append(smtp_unverified, sub)
		}
	}

	campaign1, campaign2 := data, data
	campaign1.Recipients = smtp_verified
	campaign2.Recipients = smtp_unverified

	// re-route as new campaigns to listmonk
	// Send campaign requests to internal listmonk
	resp, err := CreateNewList(listEndpoint, apiUsername, accessToken, membershipEndpoint, "verified_1", w)
	fmt.Println("RESP", resp)
	fmt.Println(err)
	log.Println("Deleting now...")
	err = DeleteList(listEndpoint, apiUsername, accessToken, resp)
	log.Println(err)
	// sendCapmaign(campaign1, w)
	// sendCapmaign(campaign2, w)
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

	if data.Language == "hindi" {
		if data.Frequency == "1" {
			req.List1 = append(req.List1, hindiListId1m)
		} else if data.Frequency == "3" {
			req.List1 = append(req.List1, hindiListId3m)
		}
	} else if data.Language == "english" {
		if data.Frequency == "1" {
			req.List1 = append(req.List1, englishListId1m)
		} else if data.Frequency == "3" {
			req.List1 = append(req.List1, englishListId3m)
		}
	} else if data.Language == "both" {
		if data.Frequency == "1" {
			req.List1 = append(req.List1, englishListId1m)
			req.List1 = append(req.List1, hindiListId1m)
		} else if data.Frequency == "3" {
			req.List1 = append(req.List1, englishListId3m)
			req.List1 = append(req.List1, hindiListId3m)
		}
	}

	fmt.Println(req.List1)

	SendPreferenceChangeRequest(internalEndpoint, apiUsername, accessToken, req, w)
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/proxy/change_user_pref", proxyHandler_ChangeList)
	http.HandleFunc("/proxy/messenger", proxyHandler_RoutingMessenger)

	fmt.Println("Proxy server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
