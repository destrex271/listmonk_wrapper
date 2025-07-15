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
	"net/url"
	"strconv"
)

// func GetCampaginIDFromUUID(uuid string) (string, err){
// }

// Create Campaigns
func CreateNewList(listEndpoint string, apiUsername string, accessToken string, membershipEndpoint string, list_title string) (string, error) {
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
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(content))
	if err != nil {
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
		return "", err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return "", err
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
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return err
	}
	log.Println(string(bodyBytes))
	// Set appropriate headers
	return nil
}

func FetchIDsFromUUIDs(apiUsername string, accessToken string, subsEndpoint string, recps []Recipient) ([]int, error) {
	queryString := "subscribers.uuid IN ("
	for i, recp := range recps {
		queryString += ("'" + recp.UUID + "'")
		if i < len(recps)-1 {
			queryString += ","
		}
	}
	queryString += ")"
	fmt.Println(queryString)
	queryString = url.QueryEscape(queryString)
	fmt.Println(queryString)
	req, err := http.NewRequest("GET", subsEndpoint+"?query="+queryString, nil)
	if err != nil {
		return nil, err
	}

	auth := apiUsername + ":" + accessToken
	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", authHeader)

	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	fmt.Println("I am here 5....")
	if err != nil {
		return nil, err
	}
	fmt.Println("I am here 6....")
	fmt.Println("OLA", resp)
	bodyBytes, _ := io.ReadAll(resp.Body)
	log.Println("HERE!!!")
	log.Println(string(bodyBytes))

	var response ResponseSubQuery
	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		return nil, err
	}

	fmt.Println("RES->>", response)

	var idStrings []int
	for _, sub := range response.Data.Results {
		idStrings = append(idStrings, sub.ID)
	}

	return idStrings, nil
}

// Update Recepients - ADD or REMOVE from list
func UpdateRecepients(apiUsername string, accessToken string, membershipEndpoint string, ids []int, listId string, action string) error {
	fmt.Println("I am here...", action, listId)
	id, _ := strconv.Atoi(listId)
	body := UpdateSubscribers{
		Ids:           ids,
		Action:        action, // add or remove
		TargetListIDs: []int{id},
	}
	if action == "add" {
		body.Status = "confirmed"
	}

	content, err := json.Marshal(body)
	fmt.Println("I am here 2", content)
	if err != nil {
		log.Println("HERE -1!!!")
		log.Println(err.Error())
		return err
	}

	fmt.Println("I am here 3....", membershipEndpoint)
	req, err := http.NewRequest("PUT", membershipEndpoint, bytes.NewBuffer(content))
	if err != nil {
		log.Println("HERE 1!!!")
		log.Println(err.Error())
		return err
	}

	// Set appropriate headers
	fmt.Println("I am here 4....")
	auth := apiUsername + ":" + accessToken
	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", authHeader)

	// Send the request using the default HTTP client
	client := &http.Client{}
	resp, err := client.Do(req)
	fmt.Println("I am here 5....")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	fmt.Println("I am here 6....")
	fmt.Println("OLA", resp)
	bodyBytes, _ := io.ReadAll(resp.Body)
	log.Println("HERE!!!")
	log.Println(string(bodyBytes))

	fmt.Println("I am here 7....")
	return nil
}

// Send Campaign
func SendCapmaign(campaignEndpoint string, apiUsername string, accessToken string, body Postback, listId int, messenger string) (int, error) {
	url := campaignEndpoint
	auth := apiUsername + ":" + accessToken
	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))

	campaignBody := CreateCampaignRequest{
		FromEmail:   body.FromEmail,
		Name:        body.Campaign.Name,
		Subject:     body.Subject,
		Lists:       []int{listId},
		Type:        "regular",
		ContentType: body.ContentType,
		Body:        body.Body,
		Messenger:   messenger,
	}
	content, err := json.Marshal(&campaignBody)
	if err != nil {
		return -1, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(content))
	if err != nil {
		return -1, err
	}

	// Set appropriate headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", authHeader)

	fmt.Println("HII!!!!", req)

	// Send the request using the default HTTP client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return -1, err
	}
	fmt.Println(resp)
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return -1, fmt.Errorf("failed to read response body: %w", err)
	}

	var result CampaignResponse
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return -1, fmt.Errorf("failed to parse JSON: %w", err)
	}

	payload, err := json.Marshal(map[string]string{
		"status": "running",
	})
	if err != nil {
		return -1, err
	}

	req, err = http.NewRequest("PUT", campaignEndpoint+"/"+strconv.Itoa(result.Data.ID)+"/status", bytes.NewBuffer(payload))
	if err != nil {
		return -1, err
	}

	// Set appropriate headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", authHeader)

	fmt.Println("HII!!!!", req)

	// Send the request using the default HTTP client
	// client = &http.Client{}
	// resp, err = client.Do(req)
	// if err != nil {
	// 	return -1, err
	// }
	// fmt.Println(resp)
	// defer resp.Body.Close()
	// bodyBytes, err = io.ReadAll(resp.Body)
	// if err != nil {
	// 	return -1, fmt.Errorf("failed to read response body: %w", err)
	// }

	return result.Data.ID, nil
}

func GetCampaignStatus(cpId int, campaignEndpoint string, apiUsername string, accessToken string) bool {
	req, err := http.NewRequest("GET", campaignEndpoint+"/"+strconv.Itoa(cpId), nil)
	if err != nil {
		log.Println(err)
		return false
	}

	auth := apiUsername + ":" + accessToken
	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))

	// Set appropriate headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", authHeader)

	// Send the request using the default HTTP client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return false
	}
	fmt.Println(resp)
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("failed to read response body: %w\n", err)
		return false
	}

	data := make(map[string]interface{})
	err = json.Unmarshal(bodyBytes, &data)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	cpData := data["data"].(map[string]interface{})
	cpSent := cpData["status"].(string)

	return cpSent != "running"
}
