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
	ID                int            `json:"id"`
	Name              string         `json:"name"`
	Subject           string         `json:"subject"`
	Type              string         `json:"type"` // "regular" or "A/B"
	EmailTemplateID   int            `json:"email_template_id"`
	EmailTemplateName string         `json:"email_template_name"`
	Content           string         `json:"content"`
	AltContent        string         `json:"alt_content"`
	CreatedAt         string         `json:"created_at"`
	Status            string         `json:"status"` // "draft", "scheduled", "running", "paused", "completed", etc.
	SentAt            *string        `json:"sent_at"`
	Source            string         `json:"source"`
	Tags              []string       `json:"tags"`
	Lists             []CampaignList `json:"lists"` // List IDs targeted
	FromEmail         string         `json:"from_email"`
	FromName          string         `json:"from_name"`
	SMTPHost          string         `json:"smtp_host"`
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

type SubscriberRequest struct {
	Email     string   `json:"email"`
	Name      string   `json:"name"`
	Lists     []string `json:"lists"`
	Frequency int      `json:"freq"`
}

type PostbackSubscriberRequest struct {
	Email   string         `json:"email"`
	Name    string         `json:"name"`
	Status  string         `json:"status"`
	Lists   []int          `json:"lists"`
	Attribs map[string]any `json:"attribs"`
	PreCon  bool           `json:"preconfirm_subscriptions"`
}

type AhaSendWebhook struct {
	Type string `json:"type"`
	Timestamp string `json:"timestamp"`
	WebhookID string `json:"webhook_id"`
	Data AhaSendWebhookBody `json:"data"`
}

type AhaSendWebhookBody struct{
	AccountID string `json:"account_id"`
	Event string `json:"event"`
	From string `json:"from"`
	Recepient string `json:"recepient"`
	Subject string `json:"subject"`
	MessageIDHeader string `json:"message_id_header"`
	ID string `json:"id"`
}

type ListMonkWebhook struct{
	Email string `json:"email"`
	CampaignUUID string `json:"campaign_uuid"`
	Source string `json:"source"`
	Type string `json:"type"`
	Meta string `json:"meta"`
}

// ZOHOMailAgentWebhook structs for new bounce payloads
type ZOHOMailAgentWebhook struct {
	EventName       []string          `json:"event_name"`
	EventMessage    []ZOHOEventMessage `json:"event_message"`
	MailagentKey    string            `json:"mailagent_key"`
	WebhookRequestID string           `json:"webhook_request_id"`
}

type ZOHOEventMessage struct {
	EmailInfo ZOHOEmailInfo   `json:"email_info"`
	EventData []ZOHOEventData `json:"event_data"`
	RequestID string      `json:"request_id"`
}

type ZOHOEmailInfo struct {
	CC               []ZOHOEmailAddress `json:"cc"`
	ClientReference  string         `json:"client_reference"`
	BCC              []ZOHOEmailAddress `json:"bcc"`
	IsSMTPTrigger    bool           `json:"is_smtp_trigger"`
	Subject          string         `json:"subject"`
	BounceAddress    string         `json:"bounce_address"`
	IsSynced         bool           `json:"is_synced"`
	EmailReference   string         `json:"email_reference"`
	ReplyTo          []ZOHOReplyTo    `json:"reply_to"`
	From             ZOHOFrom         `json:"from"`
	To               []ZOHOEmailAddress `json:"to"`
	Tag              string         `json:"tag"`
	ProcessedTime    string         `json:"processed_time"`
	Object           string         `json:"object"`
}

type ZOHOEmailAddress struct {
	EmailAddress struct {
		Address string `json:"address"`
		Name    string `json:"name"`
	} `json:"email_address"`
}

type ZOHOReplyTo struct {
	Address string `json:"address"`
	Name    string `json:"name"`
}

type ZOHOFrom struct {
	Address string `json:"address"`
	Name    string `json:"name"`
}

type ZOHOEventData struct {
	Details []ZOHODetails `json:"details"`
	Object  string    `json:"object"`
}

type ZOHODetails struct {
	Reason           string `json:"reason"`
	BouncedRecipient string `json:"bounced_recipient"`
	Time             string `json:"time"`
	DiagnosticMessage string `json:"diagnostic_message"`
}

