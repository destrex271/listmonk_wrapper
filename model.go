package main

import "net/textproto"

type JSON map[string]any
type Headers []map[string]string

type postback struct {
	Subject     string       `json:"subject"`
	FromEmail   string       `json:"from_email"`
	ContentType string       `json:"content_type"`
	Body        string       `json:"body"`
	Recipients  []recipient  `json:"recipients"`
	Campaign    *campaign    `json:"campaign"`
	Attachments []attachment `json:"attachments"`
}

type campaign struct {
	FromEmail string   `json:"from_email"`
	UUID      string   `json:"uuid"`
	Name      string   `json:"name"`
	Headers   Headers  `json:"headers"`
	Tags      []string `json:"tags"`
}

type recipient struct {
	UUID    string `json:"uuid"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Attribs JSON   `json:"attribs"`
	Status  string `json:"status"`
}

type attachment struct {
	Name    string               `json:"name"`
	Header  textproto.MIMEHeader `json:"header"`
	Content []byte               `json:"content"`
}
