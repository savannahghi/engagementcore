package dto

import (
	"time"

	"github.com/savannahghi/enumutils"
)

// SendSMSPayload is used to serialise an SMS sent through the AIT service REST API
type SendSMSPayload struct {
	To      []string           `json:"to"`
	Message string             `json:"message"`
	Sender  enumutils.SenderID `json:"sender"`
	Segment *string            `json:"segment"`
}

// EMailMessage holds data required to send emails
type EMailMessage struct {
	Subject string   `json:"subject,omitempty"`
	Text    string   `json:"text,omitempty"`
	To      []string `json:"to,omitempty"`
}

// FeedbackInput is reason a user gave a certain NPS score
// Its stored as question answer in plain text
type FeedbackInput struct {
	Question string `json:"question" firestore:"question"`
	Answer   string `json:"answer" firestore:"answer"`
}

// Feedback is reason a user gave a certain NPS score
// Its stored as question answer in plain text
type Feedback struct {
	Question string `json:"question" firestore:"question"`
	Answer   string `json:"answer" firestore:"answer"`
}

// NPSInput is the input for a survey
type NPSInput struct {
	Name        string           `json:"name"`
	Score       int              `json:"score"`
	SladeCode   string           `json:"sladeCode"`
	Email       *string          `json:"email"`
	PhoneNumber *string          `json:"phoneNumber"`
	Feedback    []*FeedbackInput `json:"feedback"`
}

// UIDPayload is the user ID used in some inter-service requests
type UIDPayload struct {
	UID *string `json:"uid"`
}

// OutgoingEmailsLog contains the content of the sent email message sent via MailGun
type OutgoingEmailsLog struct {
	UUID    string   `json:"uuid" firestore:"uuid"`
	To      []string `json:"to" firestore:"to"`
	From    string   `json:"from" firestore:"from"`
	Subject string   `json:"subject" firestore:"subject"`
	Text    string   `json:"text" firestore:"text"`
	// MessageID is a unique identifier of mailgun's message
	MessageID   string              `json:"message-id" firestore:"messageID"`
	EmailSentOn time.Time           `json:"emailSentOn" firestore:"emailSentOn"`
	Event       *MailgunEventOutput `json:"mailgunEvent" firestore:"mailgunEvent"`
}

// MailgunEvent represents mailgun event input e.g delivered, rejected etc
type MailgunEvent struct {
	EventName   string `json:"event" firestore:"event"`
	DeliveredOn string `json:"timestamp" firestore:"deliveredOn"`
	// MessageID is a unique identifier of mailgun's message
	MessageID string `json:"message-id" firestore:"messageID"`
}

// RetrieveUserProfileInput used to retrieve user profile info using either email address or phone
type RetrieveUserProfileInput struct {
	PhoneNumber  *string `json:"phone"`
	EmailAddress *string `json:"email"`
}

//TemporaryPIN input used to send temporary PIN message
type TemporaryPIN struct {
	PhoneNumber string `json:"phoneNumber,omitempty"`
	FirstName   string `json:"firstName,omitempty"`
	PIN         string `json:"pin,omitempty"`
	Channel     int    `json:"channel,omitempty"`
}
