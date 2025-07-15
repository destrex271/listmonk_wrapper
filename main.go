package main

import (

	// "encoding/base64"

	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"time"

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
	listmonkURL		   = "http://" + os.Getenv("LISTMONK_URL")
	internalEndpoint   = "http://" + os.Getenv("LISTMONK_URL") + "/api/subscribers/switch_list"
	campaignEndpoint   = "http://" + os.Getenv("LISTMONK_URL") + "/api/campaigns"
	listEndpoint       = "http://" + os.Getenv("LISTMONK_URL") + "/api/lists"
	membershipEndpoint = "http://" + os.Getenv("LISTMONK_URL") + "/api/subscribers/lists"
	subsEndpoint       = "http://" + os.Getenv("LISTMONK_URL") + "/api/subscribers"
	hindiListId3m, _   = strconv.Atoi(os.Getenv("HINDI_LIST_3M"))
	englishListId3m, _ = strconv.Atoi(os.Getenv("ENGLISH_LIST_3M"))
	hindiListId1m, _   = strconv.Atoi(os.Getenv("HINDI_LIST_1M"))
	englishListId1m, _ = strconv.Atoi(os.Getenv("ENGLISH_LIST_1M"))
)

func operate(listName string, campaign Postback, messenger string) error {
	// Create Temporary List
	listId, err := CreateNewList(listEndpoint, apiUsername, accessToken, membershipEndpoint, listName)
	if err != nil {
		return err
	}
	log.Println("Created list with ID", listId)
	// Get Recepients IDs
	ids, err := FetchIDsFromUUIDs(apiUsername, accessToken, subsEndpoint, campaign.Recipients)
	if err != nil {
		return err
	}
	log.Println("Got recepient ids as: ", ids)
	// Add Recepients to List
	err = UpdateRecepients(apiUsername, accessToken, membershipEndpoint, ids, listId, "add")
	if err != nil {
		return err
	}
	log.Println("Added recepeients to List")

	// Create Campaign with new List
	lId, err := strconv.Atoi(listId)
	if err != nil {
		return err
	}
	campaign.Campaign.Name += "_" + listName
	cpId, err := SendCapmaign(campaignEndpoint, apiUsername, accessToken, campaign, lId, messenger)
	if err != nil {
		return err
	}

	log.Println("Sent campaign", cpId)
	// log.Println("Polling status....")

	// for !GetCampaignStatus(cpId, campaignEndpoint, apiUsername, accessToken) {
	// 	log.Println("Campaign still running...")
	// 	time.Sleep(time.Second * 10)
	// }
	// log.Println("Campaign sent.., time for cleanup")

	// // Remove Recepients from List
	// err = UpdateRecepients(apiUsername, accessToken, membershipEndpoint, ids, listId, "remove")
	// if err != nil {
	// 	return err
	// }

	// // Wait till status is compeleted

	// // Delete the List
	// err = DeleteList(listEndpoint, apiUsername, accessToken, listId)
	// if err != nil {
	// 	return err
	// }

	return nil
}

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

	log.Println("VERIFIED: ", smtp_verified)
	log.Println("UNVERIFIED: ", smtp_unverified)

	if len(campaign1.Recipients) > 0 {
		log.Println("Sending for verified....")
		log.Println("Recps are: ", campaign1.Recipients)
		fmt.Scanln()
		err := operate("verified_"+campaign1.Campaign.Name, campaign1, "email-verified")
		if err != nil {
			log.Println("HELLO -.", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	// if len(campaign2.Recipients) > 0 {
	// 	log.Println("Sending for unverified....")
	// 	log.Println("Sending for verified....")
	// 	err := operate("unverified_"+campaign2.Campaign.Name, campaign2, "email-unverified")
	// 	if err != nil {
	// 		log.Println("HELLO -.", err.Error())
	// 		http.Error(w, err.Error(), http.StatusInternalServerError)
	// 		return
	// 	}
	// }
	http.Error(w, "Successfully sent data", http.StatusOK)
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

func withCORS(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Required CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // Or a specific origin like http://localhost:5500
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Preflight request
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Call the actual handler
		handler(w, r)
	}
}

func proxyHandler_SendCampaign(w http.ResponseWriter, r *http.Request){
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	auth := apiUsername + ":" + accessToken
	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))

	// 2. Create the campaign
	reqCreate, _ := http.NewRequest("POST", listmonkURL +"/api/campaigns", bytes.NewBuffer(body))
	reqCreate.Header.Set("Authorization", authHeader)
	reqCreate.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(reqCreate)
	if err != nil || resp.StatusCode != 200 {
		http.Error(w, "Failed to create campaign", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// bodyBytes, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	log.Fatal("Failed to read response body:", err)
	// }
	// fmt.Println("Raw response body:", string(bodyBytes))

	var created struct {
		Data CRCampaign `json:"data"`
	}
	fmt.Println("hi->", resp.Body)
	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to decode campaign response", http.StatusInternalServerError)
		return
	}

	// 3. Set campaign status to running
	statusPayload := map[string]string{"status": "running"}
	payloadBytes, _ := json.Marshal(statusPayload)
	reqStatus, _ := http.NewRequest("PUT", fmt.Sprintf("%s/api/campaigns/%d/status", listmonkURL, created.Data.ID), bytes.NewBuffer(payloadBytes))
	reqStatus.Header.Set("Authorization", authHeader)
	reqStatus.Header.Set("Content-Type", "application/json")
	respStatus, err := http.DefaultClient.Do(reqStatus)
	if err != nil || respStatus.StatusCode != 200 {
		bodyText, _ := io.ReadAll(respStatus.Body)
		log.Printf("Failed to set campaign status: %s", string(bodyText))
		http.Error(w, "Failed to start campaign", http.StatusInternalServerError)
		return
	}
	defer respStatus.Body.Close()

	// 4. Poll for completion
	for {
		time.Sleep(2 * time.Second)
		reqCheck, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/campaigns/%d", listmonkURL, created.Data.ID), nil)
		reqCheck.Header.Set("Authorization", authHeader)
		respCheck, err := http.DefaultClient.Do(reqCheck)
		if err != nil || respCheck.StatusCode != 200 {
			http.Error(w, "Failed to check campaign status", http.StatusInternalServerError)
			return
		}
		var statusResp struct {
			Data CRCampaign `json:"data"`
		}
		json.NewDecoder(respCheck.Body).Decode(&statusResp)
		respCheck.Body.Close()

		if statusResp.Data.Status == "finished" {
			break
		}
	}

	// 5. Respond success
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/proxy/change_user_pref", proxyHandler_ChangeList)
	http.HandleFunc("/proxy/send_campaign", withCORS(proxyHandler_SendCampaign))

	fmt.Println("Proxy server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
