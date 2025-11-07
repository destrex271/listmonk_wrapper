package main

import (

	// "encoding/base64"

	"bytes"
	"context"
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
	database_url       = os.Getenv("DATABASE_URL")
)

func syncSubscribers() {
	conn, err := pgx.Connect(context.Background(), database_url)
	if err != nil{
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())o

	move_verified := "SELECT mark_verified_on_view();"
	check_bounce_threshold := "SELECT check_bounce_threshold();"
	sync_subs := "SELECT sync_all_subscribers_to_verified_status();"

	// convert status of unverified to verified if viewed
	_, err = conn.Exec(context.Background(), move_verified)
	if err != nil{
		fmt.Println("%v", err)
		os.Exit(1)
	}

	// convert status of verified to unverified if bounced
	_, err = conn.Exec(context.Background(), check_bounce_threshold)
	if err != nil{
		fmt.Println("%v", err)
		os.Exit(1)
	}

	// TODO: For any bounced subs that need to be deleted, take
	// action on ASP Server DB

	// take action on verif and unverif by moving them around lists
	_, err = conn.Exec(context.Background(), sync_subs)
	if err != nil{
		fmt.Println("%v", err)
		os.Exit(1)
	}

}

func pollCampaignStatus(id str){
	// 4. Poll for completion
	for {
		time.Sleep(2 * time.Second)
		reqCheck, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/campaigns/%d", listmonkURL, id,created.Data.ID), nil)
		reqCheck.Header.Set("Authorization", authHeader)
		respCheck, err := http.DefaultClient.Do(reqCheck)
		if err != nil || respCheck.StatusCode != 200 {
			http.Error(w, "Failed to check campaign status " + err.Error(), http.StatusInternalServerError)
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
		http.Error(w, "Failed to create campaign " + err.Error(), http.StatusInternalServerError)
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
		http.Error(w, "Failed to decode campaign response " + err.Error(), http.StatusInternalServerError)
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
		http.Error(w, "Failed to start campaign " + err.Error(), http.StatusInternalServerError)
		return
	}
	defer respStatus.Body.Close()



	// 5. Respond success
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/proxy/send_campaign", withCORS(proxyHandler_SendCampaign))

	fmt.Println("Proxy server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
