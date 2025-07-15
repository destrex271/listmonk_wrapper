package utils

import (
	"os"
	"strconv"
)

var (
	apiUsername = os.Getenv("API_USER")
	// accessToken      = "7BXtarGYcQaCiCeS706G9M83DxC1ZJux"
	accessToken        = os.Getenv("API_TOKEN")
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
