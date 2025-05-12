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
	Ids           []string `json:"ids"`
	Action        string   `json:"action"`
	TargetListIDs []string `json:"target_list_ids"`
	Status        string   `json:"status"`
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
