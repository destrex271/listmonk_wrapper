package main

import (

	// "encoding/base64"

	"bytes"
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	// "io"
	"log"
	"net/http"

	// "net/url"
	"io"
	"os"
	"strconv"

	_ "github.com/denisenkom/go-mssqldb"
	. "github.com/destrex271/listmonk_proxy/utils"
	"github.com/jackc/pgx/v5"
)

var (
	apiUsername = os.Getenv("API_USER")
	port        = os.Getenv("PORT")
	// accessToken      = "7BXtarGYcQaCiCeS706G9M83DxC1ZJux"
	accessToken           = os.Getenv("API_TOKEN")
	listmonkURL           = "http://" + os.Getenv("LISTMONK_URL")
	internalEndpoint      = "http://" + os.Getenv("LISTMONK_URL") + "/api/subscribers/switch_list"
	campaignEndpoint      = "http://" + os.Getenv("LISTMONK_URL") + "/api/campaigns"
	listEndpoint          = "http://" + os.Getenv("LISTMONK_URL") + "/api/lists"
	membershipEndpoint    = "http://" + os.Getenv("LISTMONK_URL") + "/api/subscribers/lists"
	subsEndpoint          = "http://" + os.Getenv("LISTMONK_URL") + "/api/subscribers"
	hindiListId3m, _      = strconv.Atoi(os.Getenv("HINDI_LIST_3M"))
	englishListId3m, _    = strconv.Atoi(os.Getenv("ENGLISH_LIST_3M"))
	hindiListId1m, _      = strconv.Atoi(os.Getenv("HINDI_LIST_1M"))
	englishListId1m, _    = strconv.Atoi(os.Getenv("ENGLISH_LIST_1M"))
	database_url          = os.Getenv("DATABASE_URL")
	original_database_url = os.Getenv("ASP_DATABASE_URL")
	asp_username          = os.Getenv("ASP_USER_NAME")
	asp_passwd            = os.Getenv("ASP_PASSWD")

	blockListDropTime, _ = strconv.Atoi(os.Getenv("BLOCKLIST_DROP_TIME_HOURS"))
	syncSubsTime, _      = strconv.Atoi(os.Getenv("SYNC_SUBS_TIME_HOURS"))
)

func dropBlocklist() {
	//
	log.Print("Deleting blocklisted subscribers from databases")
	conn, err := pgx.Connect(context.Background(), database_url)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	// Get blocklisted subscribers email
	query := "SELECT email FROM subscribers WHERE status = 'blocklisted';"
	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		log.Fatalf("Query failed: %v\n", err)
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
		_, err := db.Exec("DELETE FROM subscribers WHERE emailid = @p1", email)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Delete blocklisted subscribers in listmonk
	log.Println(subscribers)
	_, err = conn.Exec(context.Background(), "DELETE FROM subscribers WHERE email = ANY($1)", subscribers)

	if err != nil {
		log.Fatal("unable to delete subscriber from db")
	}

	log.Println("Successfully removed: ", subscribers)

}

func syncSubscribers() {
	conn, err := pgx.Connect(context.Background(), database_url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return
	}
	defer conn.Close(context.Background())

	log.Println("Synchronizing subscriber lists...")

	move_verified := "SELECT mark_verified_on_view();"
	check_bounce_threshold := "SELECT check_bounce_threshold();"
	sync_subs := "SELECT sync_all_subscribers_to_verified_status();"

	// convert status of unverified to verified if viewed
	_, err = conn.Exec(context.Background(), move_verified)
	if err != nil {
		fmt.Println("%v", err)
		return
	}

	// convert status of verified to unverified if bounced
	_, err = conn.Exec(context.Background(), check_bounce_threshold)
	if err != nil {
		fmt.Println("%v", err)
		return
	}

	// TODO: For any bounced subs that need to be deleted, take
	// action on ASP Server DB

	// take action on verif and unverif by moving them around lists
	_, err = conn.Exec(context.Background(), sync_subs)
	if err != nil {
		fmt.Println("%v", err)
		return
	}

	log.Println("Successfully synced subsribers")

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

func proxyHandler_SyncSubs(w http.ResponseWriter, r *http.Request) {
	_, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	auth := apiUsername + ":" + accessToken
	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))

	req, _ := http.NewRequest("GET", listmonkURL+"/api/lists", nil)
	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != 200 {
		http.Error(w, "Unauthorized Request!\n"+err.Error(), http.StatusUnauthorized)
		return
	}

	syncSubscribers()

	// Send success msg
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})

}

func proxyHandler_SendCampaign(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	auth := apiUsername + ":" + accessToken
	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))

	// 2. Create the campaign
	reqCreate, _ := http.NewRequest("POST", listmonkURL+"/api/campaigns", bytes.NewBuffer(body))
	reqCreate.Header.Set("Authorization", authHeader)
	reqCreate.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(reqCreate)
	if err != nil || resp.StatusCode != 200 {
		http.Error(w, "Failed to create campaign "+err.Error(), http.StatusInternalServerError)
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
	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to decode campaign response "+err.Error(), http.StatusInternalServerError)
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
		http.Error(w, "Failed to start campaign "+err.Error(), http.StatusInternalServerError)
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
	http.HandleFunc("/proxy/sync_subs", withCORS(proxyHandler_SyncSubs))

	fmt.Println("Proxy server is running on port 8080")

	go func() {
		for true {
			// Synchronize subscribers every 24 hours
			syncSubscribers()
			time.Sleep(time.Duration(syncSubsTime) * time.Hour)
		}
	}()

	// go func() {
	// 	for true {
	// 		dropBlocklist()
	// 		time.Sleep(time.Duration(blockListDropTime) * time.Hour)
	// 	}
	// }()

	log.Fatal(http.ListenAndServe(port, nil))
}
