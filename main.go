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

	// "time"

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
	asp_server            = os.Getenv("ASP_SERVER")
	asp_database 		  = os.Getenv("ASP_DATABASE")

	cronEnabled = os.Getenv("CRON_ENABLED")

	blockListDropTime, _ = strconv.Atoi(os.Getenv("BLOCKLIST_DROP_TIME_HOURS"))
	syncSubsTime, _      = strconv.Atoi(os.Getenv("SYNC_SUBS_TIME_HOURS"))
	mainWebsiteUnsubURL  = os.Getenv("MAIN_WEBSITE_UNSUB_LINK")
)

func markBlockListInSource() {
	//
	log.Print("Marking blocklisted subscribers in source database")
	conn, err := pgx.Connect(context.Background(), database_url)
	if err != nil {
		log.Printf("Unable to connect to database: %v\n", err)
	}
	defer conn.Close(context.Background())

	// Get blocklisted subscribers email
	query := "SELECT email FROM subscribers WHERE status = 'blocklisted';"
	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		log.Printf("Query failed: %v\n", err)
		return
	}
	defer rows.Close()

	var subscribers []string
	for rows.Next() {
		var email string
		if err := rows.Scan(&email); err != nil {
			fmt.Fprintf(os.Stderr, "Scan failed: %v\n", err)
			return
		}
		subscribers = append(subscribers, email)
	}

	// Mark blocklisted subsribers from ASP DB
	db, err := sql.Open("mssql", fmt.Sprintf("server=%s;database=%s;user id=%s;password=%s;Encrypt=True;TrustServerCertificate=True;", asp_server, asp_database, asp_username, asp_passwd))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	for _, email := range subscribers {
		_, err := db.Exec("UPDATE t_newsletter_subscriber SET activeyn='B' WHERE emailid=%s", email)
		if err != nil {
			log.Printf("err: %w\n", err)
		}else{
			log.Printf("%s marked as blocklisted", email)
		}
	}
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

func UnsubscribeHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Missing 'code' query parameter", http.StatusBadRequest)
		return
	}

	conn, err := pgx.Connect(context.Background(), database_url)
	if err != nil {
		log.Printf("Unable to connect to database for unsubscribe: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer conn.Close(context.Background())

	// 1. Get subscriber ID using the unsub_code in attribs
	var subscriberID int
	err = conn.QueryRow(context.Background(), "SELECT id FROM subscribers WHERE attribs->>'unsub_code' = $1", code).Scan(&subscriberID)
	if err == pgx.ErrNoRows {
		log.Printf("Subscriber with unsub_code '%s' not found", code)
		http.Redirect(w, r, mainWebsiteUnsubURL, http.StatusFound) // Redirect even if not found, to avoid information leakage
		return
	}
	if err != nil {
		log.Printf("Error getting subscriber by unsub_code '%s': %v", code, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// 2. Blocklist and unsubscribe the user from all lists
	_, err = conn.Exec(context.Background(), `
		WITH b AS (
			UPDATE subscribers SET status='blocklisted', updated_at=NOW()
			WHERE id = ANY($1::INT[])
		)
		UPDATE subscriber_lists SET status='unsubscribed', updated_at=NOW()
		WHERE subscriber_id = ANY($1::INT[]);
	`, []int{subscriberID})
	if err != nil {
		log.Printf("Error blocklisting subscriber ID %d: %v", subscriberID, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// 3. Redirect to the main website with the code as a query parameter
	redirectURL := fmt.Sprintf("%s?data=%s", mainWebsiteUnsubURL, code)
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	var webhookMessage AhaSendWebhook

	if err := json.NewDecoder(r.Body).Decode(&webhookMessage); err != nil {
		log.Printf("unable to parse webhook message! %v\n", err)
		http.Error(w, "Invalid webhook payload", http.StatusOK)
		return
	}

	fmt.Printf("%v\n", webhookMessage)

	// fetch latest campaign with the given subject
	conn, err := pgx.Connect(context.Background(), database_url)
	if err != nil {
		log.Printf("Unable to connect to database: %v\n", err)
		http.Error(w, "Internal server error", http.StatusOK)
		return
	}
	defer conn.Close(context.Background())

	query := "SELECT uuid FROM campaigns WHERE subject = $1 ORDER BY created_at DESC;"
	var campaignUUID string

	if err = conn.QueryRow(context.Background(), query, webhookMessage.Data.Subject).Scan(&campaignUUID); err != nil {
		log.Printf("unable to find campaign for subject '%s': %v\n", webhookMessage.Data.Subject, err)
		// We can't attribute the bounce without a campaign, but we should still acknowledge the webhook.
		json.NewEncoder(w).Encode(map[string]bool{"success": true})
		return
	}

	bounceReq := &ListMonkWebhook{
		Email:        webhookMessage.Data.Recepient,
		CampaignUUID: campaignUUID,
		Source:       "api",
		Type:         "soft", // Default
		Meta:         "{}",
	}

	// convert to listmonk type
	switch webhookMessage.Type {
	case "message.bounced":
		// handle as a soft bounce
		log.Printf("Soft bounce!")
		bounceReq.Type = "soft"
	case "message.suppressed", "suppression.created", "message.failed":
		log.Printf("Hard Bounce!")
		bounceReq.Type = "hard"
	}

	body, err := json.Marshal(bounceReq)
	if err != nil {
		log.Printf("error marshalling bounce request: %v", err)
		http.Error(w, "Internal server error", http.StatusOK)
		return
	}

	auth := apiUsername + ":" + accessToken
	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))

	log.Printf("Sending bounce report to listmonk webhook %v \t %v\n", *bounceReq, body)
	req, err := http.NewRequest("POST", listmonkURL+"/webhooks/bounce", bytes.NewBuffer(body))
	if err != nil {
		log.Printf("error creating bounce request: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("error sending bounce to listmonk: %v", err)
		http.Error(w, "Error sending bounce report", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Printf("listmonk bounce webhook returned non-200 status: %d. Body: %s", resp.StatusCode, string(bodyBytes))
		http.Error(w, "Error from listmonk", resp.StatusCode)
		return
	}

	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func webhookHandler_ZOHO(w http.ResponseWriter, r *http.Request) {
	var webhookMessage ZOHOMailAgentWebhook

	if err := json.NewDecoder(r.Body).Decode(&webhookMessage); err != nil {
		log.Printf("unable to parse ZOHO webhook message! %v\n", err)
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
		return
	}

	fmt.Printf("%+v\n", webhookMessage)

	conn, err := pgx.Connect(context.Background(), database_url)
	if err != nil {
		log.Printf("Unable to connect to database: %v\n", err)
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
		return
	}
	defer conn.Close(context.Background())

	auth := apiUsername + ":" + accessToken
	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))

	// Process each event message
	for _, msg := range webhookMessage.EventMessage {
		query := "SELECT uuid FROM campaigns WHERE subject = $1 ORDER BY created_at DESC;"
		var campaignUUID string

		if err := conn.QueryRow(context.Background(), query, msg.EmailInfo.Subject).Scan(&campaignUUID); err != nil {
			log.Printf("unable to find campaign for subject '%s': %v\n", msg.EmailInfo.Subject, err)
			// Still acknowledge the webhook even if campaign not found
			continue
		}

		bounceType := "soft" // Default
		if len(webhookMessage.EventName) > 0 {
			if webhookMessage.EventName[0] == "hardbounce" {
				bounceType = "hard"
			}
		}

		for _, eventData := range msg.EventData {
			for _, details := range eventData.Details {
				bounceReq := &ListMonkWebhook{
					Email:        details.BouncedRecipient,
					CampaignUUID: campaignUUID,
					Source:       "api",
					Type:         bounceType,
					Meta:         fmt.Sprintf(`{"reason": "%s", "diagnostic_message": "%s"}`, details.Reason, details.DiagnosticMessage),
				}

				body, err := json.Marshal(bounceReq)
				if err != nil {
					log.Printf("error marshalling bounce request: %v", err)
					continue
				}

				log.Printf("Sending bounce report to listmonk webhook: %s\n", string(body))
				req, err := http.NewRequest("POST", listmonkURL+"/webhooks/bounce", bytes.NewBuffer(body))
				if err != nil {
					log.Printf("error creating bounce request: %v", err)
					continue
				}

				req.Header.Set("Authorization", authHeader)
				req.Header.Set("Content-Type", "application/json")
				resp, err := http.DefaultClient.Do(req)
				if err != nil {
					log.Printf("error sending bounce to listmonk: %v", err)
					continue
				}
				defer resp.Body.Close()

				if resp.StatusCode != 200 {
					bodyBytes, _ := io.ReadAll(resp.Body)
					log.Printf("listmonk bounce webhook returned non-200 status: %d. Body: %s", resp.StatusCode, string(bodyBytes))
				}
			}
		}
	}


	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/proxy/send_campaign", withCORS(proxyHandler_SendCampaign))
	http.HandleFunc("/proxy/sync_subs", withCORS(proxyHandler_SyncSubs))
	http.HandleFunc("/proxy/unsubscribe", withCORS(UnsubscribeHandler))
	http.HandleFunc("/proxy/on_bounce", withCORS(webhookHandler))
	http.HandleFunc("/proxy/on_bounce_zoho", withCORS(webhookHandler_ZOHO))

	fmt.Println("Proxy server is running on port 8080")

	// go func() {
	// 	for true {
	// 		// Synchronize subscribers every 24 hours
	// 		syncSubscribers()
	// 		time.Sleep(time.Duration(syncSubsTime) * time.Hour)
	// 	}
	// }()

	// go func() {
	// 	for true {
	// 		dropBlocklist()
	// 		time.Sleep(time.Duration(blockListDropTime) * time.Hour)
	// 	}
	// }()

	// Go-Routine to run blocklist job every 24 hours.
	go func() {
		for true{
			if cronEnabled != "1"{
				continue
			}
			markBlockListInSource()
			time.Sleep(time.Duration(blockListDropTime) * time.Hour)
		}
	} ()

	log.Fatal(http.ListenAndServe(port, nil))
}
