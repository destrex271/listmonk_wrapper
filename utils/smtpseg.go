package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

// Create Campaigns
func CreateNewList(listEndpoint string, apiUsername string, accessToken string, membershipEndpoint string, list_title string, w http.ResponseWriter) (string, error) {
	url := listEndpoint
	auth := apiUsername + ":" + accessToken
	authHeader := base64.StdEncoding.EncodeToString([]byte(auth))

	body := &CreateListReq{
		Name:        list_title,
		Type:        "private",
		Optin:       "single",
		Tags:        []string{"newsletter"},
		Description: "Temp List",
	}

	content, err := json.Marshal(body)
	log.Println("Sending- > ", string(content))
	if err != nil {
		http.Error(w, "Unable to parse request", http.StatusInternalServerError)
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(content))
	if err != nil {
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return "", err
	}

	// Set appropriate headers
	log.Println(authHeader, auth)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Basic "+authHeader)

	// Send the request using the default HTTP client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Error sending request", http.StatusInternalServerError)
		return "", err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
	}
	defer resp.Body.Close()

	fmt.Println(string(bodyBytes))

	var result map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		log.Println("Error unmarshaling JSON:", err)
		return "", err
	}

	// Access the 'data' field
	data, ok := result["data"].(map[string]interface{})
	if !ok {
		return "", errors.New("No data found")
	}
	log.Println(data["id"])
	id := data["id"].(float64)
	return strconv.FormatFloat(id, 'f', -1, 64), nil
}

func DeleteList(listEndpoint string, apiUsername string, accessToken string, listId string) error {

	auth := apiUsername + ":" + accessToken
	authHeader := base64.StdEncoding.EncodeToString([]byte(auth))
	log.Println(authHeader, auth)
	req, err := http.NewRequest("DELETE", listEndpoint+"/"+listId, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Basic "+authHeader)
	if err != nil {
		return err
	}

	// Send the request using the default HTTP client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return err
	}
	log.Println(string(bodyBytes))
	defer resp.Body.Close()

	// Set appropriate headers
	return nil
}

// Update Recepients
func UpdateRecepients(apiUsername string, accessToken string, listEndpoint string, action string, recps []Recipient, listId string) error {
	ids := []string{}
	for _, rec := range recps {
		ids = append(ids, rec.UUID)
	}
	body := UpdateSubscribers{
		Ids:           ids,
		Action:        "add",
		TargetListIDs: []string{listId},
		Status:        "confirmed",
	}

	content, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", listEndpoint, bytes.NewBuffer(content))
	if err != nil {
		return err
	}

	// Set appropriate headers
	auth := apiUsername + ":" + accessToken
	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", authHeader)

	// Send the request using the default HTTP client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	fmt.Println(resp)
	defer resp.Body.Close()

	return nil
}

// Send Campaign
func SendCapmaign(campaignEndpoint string, apiUsername string, accessToken string, body Postback, recps []Recipient, w http.ResponseWriter) {
	url := campaignEndpoint
	auth := apiUsername + ":" + accessToken
	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))

	content, err := json.Marshal(body)
	if err != nil {
		http.Error(w, "Unable to parse request", http.StatusInternalServerError)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(content))
	if err != nil {
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}

	// Set appropriate headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", authHeader)

	fmt.Println("HII!!!!", req)

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
