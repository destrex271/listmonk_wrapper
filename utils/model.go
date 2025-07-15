package utils

import "net/textproto"

type JSON map[string]any
type Headers []map[string]string

type Postback struct {
	Subject     string       `json:"subject"`
	FromEmail   string       `json:"from_email"`
	ContentType string       `json:"content_type"`
	Body        string       `json:"body"`
	Recipients  []Recipient  `json:"recipients"`
	Campaign    *Campaign    `json:"campaign"`
	Attachments []Attachment `json:"attachments"`
}

type Campaign struct {
	FromEmail string   `json:"from_email"`
	UUID      string   `json:"uuid"`
	Name      string   `json:"name"`
	Headers   Headers  `json:"headers"`
	Tags      []string `json:"tags"`
}

type CampaignList struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type CRCampaign struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Subject      string    `json:"subject"`
	Type         string    `json:"type"` // "regular" or "A/B"
	EmailTemplateID int    `json:"email_template_id"`
	EmailTemplateName string `json:"email_template_name"`
	Content      string    `json:"content"`
	AltContent   string    `json:"alt_content"`
	CreatedAt    string    `json:"created_at"`
	Status       string    `json:"status"` // "draft", "scheduled", "running", "paused", "completed", etc.
	SentAt       *string   `json:"sent_at"`
	Source       string    `json:"source"`
	Tags         []string  `json:"tags"`
	Lists        []CampaignList `json:"lists"` // List IDs targeted
	FromEmail    string    `json:"from_email"`
	FromName     string    `json:"from_name"`
	SMTPHost     string    `json:"smtp_host"`
}

type CreateCampaignRequest struct {
	FromEmail   string `json:"from_email"`
	Name        string `json:"name"`
	Subject     string `json:"subject"`
	Lists       []int  `json:"lists"`
	Type        string `json:"type"`
	ContentType string `json:"content_type"`
	Body        string `json:"body"`
	Messenger   string `json:"messenger"`
}

type CampaignResponse struct {
	Data struct {
		ID int `json:"id"`
	} `json:"data"`
}

type Recipient struct {
	UUID    string `json:"uuid"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Attribs JSON   `json:"attribs"`
	Status  string `json:"status"`
}

type Attachment struct {
	Name    string               `json:"name"`
	Header  textproto.MIMEHeader `json:"header"`
	Content []byte               `json:"content"`
}

type CreateListReq struct {
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	Optin       string   `json:"optin"`
	Tags        []string `json:"tags"`
	Description string   `json:"description"`
}

type UpdateSubscribers struct {
	Ids           []int  `json:"ids"`
	Action        string `json:"action"`
	TargetListIDs []int  `json:"target_list_ids"`
	Status        string `json:"status"`
}

// Preference selection
type RequestBody struct {
	Email string `json:"email"`
	List1 []int  `json:"lista"`
}

type RequestData struct {
	Email     string `json:"email"`
	Language  string `json:"language"`
	Frequency string `json:"frequency"`
}

type Subscriber struct {
	ID int `json:"id"`
}

type ResponseSubQuery struct {
	Data struct {
		Results []Subscriber `json:"results"`
	} `json:"data"`
}



// ---------------------------------------------------------


type SubscriberRequest struct{
	Email	string	`json:"email"`
	Name	string 	`json:"name"`
	Lists	[]string	`json:"lists"`
	Frequency int	`json:"freq"`
}

type PostbackSubscriberRequest struct{
	Email	string 	`json:"email"`
	Name	string	`json:"name"`
	Status	string	`json:"status"`
	Lists	[]int	`json:"lists"`
	Attribs	map[string]any	`json:"attribs"`
	PreCon	bool	`json:"preconfirm_subscriptions"`
}


