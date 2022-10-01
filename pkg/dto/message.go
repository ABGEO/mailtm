package dto

import "time"

type EmailAddress struct {
	Address string `json:"address"`
	Name    string `json:"name"`
}

type attachment struct {
	ID               string `json:"id"`
	Filename         string `json:"filename"`
	ContentType      string `json:"contentType"`
	Disposition      string `json:"disposition"`
	TransferEncoding string `json:"transferEncoding"`
	Related          bool   `json:"related"`
	Size             int    `json:"size"`
	DownloadURL      string `json:"downloadUrl"`
}

type baseMessage struct {
	ID             string         `json:"id"`
	AccountID      string         `json:"accountId"`
	MsgID          string         `json:"msgid"`
	From           EmailAddress   `json:"from"`
	To             []EmailAddress `json:"to"`
	Subject        string         `json:"subject"`
	Seen           bool           `json:"seen"`
	IsDeleted      bool           `json:"isDeleted"`
	HasAttachments bool           `json:"hasAttachments"`
	Size           int            `json:"size"`
	DownloadURL    string         `json:"downloadUrl"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
}

type Messages []struct {
	baseMessage
	Intro string `json:"intro"`
}

type Message struct {
	baseMessage
	Cc            []EmailAddress `json:"cc"`
	Bcc           []EmailAddress `json:"bcc"`
	Flagged       bool           `json:"flagged"`
	Verifications []string       `json:"verifications"`
	Retention     bool           `json:"retention"`
	RetentionDate time.Time      `json:"retentionDate"`
	Text          string         `json:"text"`
	HTML          []string       `json:"html"`
	Attachments   []attachment   `json:"attachments"`
}

type MessageWrite struct {
	Seen bool `json:"seen"`
}
