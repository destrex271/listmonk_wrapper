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
	"database/sql"

	// "io"
	"log"
	"net/http"

	// "net/url"
	"os"
	"strconv"

	. "github.com/destrex271/listmonk_proxy/utils"
	"github.com/jackc/pgx/v5"

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
	original_database_url = os.Getenv("ASP_DATABASE_URL")
	asp_username = os.Getenv("ASP_USER_NAME")
	asp_passwd = os.Getenv("ASP_PASSWD")
)

func dropBlocklist() {
	//
	conn, err := pgx.Connect(context.Background(), database_url)
	if err != nil{
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	// Get blocklisted subscribers email
	query := "SELECT email FROM subscribers WHERE status = 'blocklisted';"
	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Query failed: %v\n", err)
		os.Exit(1)
	}
	defer rows.Close()

	var subscribers []string
	for rows.Next() {
		var email string
		if err := rows.Scan(&email); err != nil {
			fmt.Fprintf(os.Stderr, "Scan failed: %v\n", err)
			os.Exit(1)
		}
		subscribers = append(subscribers, email)
	}


	// delete blocklisted subsribers from ASP DB
	db, err := sql.Open("sqlserver", original_database_url)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	for _, email := range subscribers {
		_, err := db.Exec("DELETE FROM subscribers WHERE email = @p1", email)
		if err != nil {
			log.Fatal(err)
		}
	}


	// Delete blocklisted subscribers in listmonk
	conn.Exec(context.Background(), "DELETE FROM subscribers WHERE email IN $1", subscribers)
}

func syncSubscribers() {
	conn, err := pgx.Connect(context.Background(), database_url)
	if err != nil{
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return
	}
	defer conn.Close(context.Background())

	move_verified := "SELECT mark_verified_on_view();"
	check_bounce_threshold := "SELECT check_bounce_threshold();"
	sync_subs := "SELECT sync_all_subscribers_to_verified_status();"

	// convert status of unverified to verified if viewed
	_, err = conn.Exec(context.Background(), move_verified)
	if err != nil{
		fmt.Println("%v", err)
		return
	}

	// convert status of verified to unverified if bounced
	_, err = conn.Exec(context.Background(), check_bounce_threshold)
	if err != nil{
		fmt.Println("%v", err)
		return
	}

	// TODO: For any bounced subs that need to be deleted, take
	// action on ASP Server DB

	// take action on verif and unverif by moving them around lists
	_, err = conn.Exec(context.Background(), sync_subs)
	if err != nil{
		fmt.Println("%v", err)
		return
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

	go func() {
		for true{
			// Synchronize subscribers every 24 hours
			syncSubscribers()
			time.Sleep(24 * time.Hour) // Every hour
		}
	}()

	go func() {
		for true{
			dropBlocklist()
			time.Sleep(48 * time.Hour) // Every two days
		}
	}()

	log.Fatal(http.ListenAndServe(":8080", nil))
}
