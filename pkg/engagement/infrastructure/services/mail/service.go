// Package mail implements an email sending API that uses MailGun and SendInBlue.
package mail

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/mailgun/mailgun-go/v4"
	"github.com/savannahghi/serverutils"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/application/common/helpers"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("gitlab.slade360emr.com/go/engagement/pkg/engagement/services/mail")

// Mail configuration constants
const (
	MailGunAPIKeyEnvVarName     = "MAILGUN_API_KEY"
	MailGunAPIBaseURLEnvVarName = "MAILGUN_API_BASE_URL"
	MailGunDomainEnvVarName     = "MAILGUN_DOMAIN"
	MailGunFromEnvVarName       = "MAILGUN_FROM"
	MailGunTimeoutSeconds       = 15

	SendInBlueAPIKeyEnvVarName  = "SEND_IN_BLUE_API_KEY"
	SendInBlueEnabledEnvVarName = "SEND_IN_BLUE_ENABLED"

	appName           = "Slade 360 HealthCloud"
	defaultUser       = "HealthCloud User"
	sendInBlueBaseURL = "https://api.sendinblue.com/v3/smtp/email"
)

// ServiceMail defines a Mail service interface
type ServiceMail interface {
	SendInBlue(ctx context.Context, subject, text string, to ...string) (string, string, error)
	SendMailgun(ctx context.Context, subject, text string, body *string, to ...string) (string, string, error)
	SendEmail(ctx context.Context, subject, text string, body *string, to ...string) (string, string, error)
	SimpleEmail(ctx context.Context, subject, text string, body *string, to ...string) (string, error)
}

// NewService initializes a new MailGun service
func NewService() *Service {
	apiKey := serverutils.MustGetEnvVar(MailGunAPIKeyEnvVarName)
	domain := serverutils.MustGetEnvVar(MailGunDomainEnvVarName)
	from := serverutils.MustGetEnvVar(MailGunFromEnvVarName)
	mg := mailgun.NewMailgun(domain, apiKey)
	mg.SetAPIBase(mailgun.APIBaseEU)

	// special case for sandbox
	if strings.Contains(domain, "sandbox") {
		mg.SetAPIBase(mailgun.APIBase)
	}

	sendInBlueEnabled := false
	if serverutils.MustGetEnvVar(SendInBlueEnabledEnvVarName) == "true" {
		sendInBlueEnabled = true
	}

	return &Service{
		Mg:                mg,
		From:              from,
		SendInBlueAPIKey:  serverutils.MustGetEnvVar(SendInBlueAPIKeyEnvVarName),
		SendInBlueEnabled: sendInBlueEnabled,
	}
}

// Service is an email sending service
type Service struct {
	Mg                *mailgun.MailgunImpl
	From              string
	SendInBlueEnabled bool
	SendInBlueAPIKey  string
}

// CheckPreconditions checks that all the required preconditions are satisfied
func (s Service) CheckPreconditions() {
	if s.Mg == nil {
		log.Panicf("uninitialized MailGun")
	}
	if s.From == "" {
		log.Panicf("uninitialized email from")
	}
	if s.SendInBlueAPIKey == "" {
		log.Panicf("uninitialized sendInBlueAPIKey")
	}
}

// MakeSendInBlueRequest makes a request to SendInBlue
func (s Service) MakeSendInBlueRequest(ctx context.Context, data map[string]interface{}, target interface{}) error {
	_, span := tracer.Start(ctx, "MakeSendInBlueRequest")
	defer span.End()
	bs, err := json.Marshal(data)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return fmt.Errorf("makeSendInBlueRequest: can't marshal data [%#v] to JSON: %w", data, err)
	}
	r := bytes.NewBuffer(bs)
	req, reqErr := http.NewRequest(http.MethodPost, sendInBlueBaseURL, r)
	if reqErr != nil {
		helpers.RecordSpanError(span, err)
		return reqErr
	}

	sendInBlueAPIKey := serverutils.MustGetEnvVar(SendInBlueAPIKeyEnvVarName)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", sendInBlueAPIKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return fmt.Errorf("SendInBlue API error: %w", err)
	}

	respBs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return fmt.Errorf("SendInBlue content error: %w", err)
	}

	if resp.StatusCode > 201 {
		return fmt.Errorf("SendInBlue API Error: %s", string(respBs))
	}

	err = json.Unmarshal(respBs, target)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return fmt.Errorf("unable to unmarshal SendInBlue resp: %w", err)
	}

	return nil
}

// SendInBlue sends email via the SendInBlue service
func (s Service) SendInBlue(ctx context.Context, subject, text string, to ...string) (string, string, error) {
	ctx, span := tracer.Start(ctx, "SendInBlue")
	defer span.End()
	s.CheckPreconditions()
	sender := map[string]string{
		"name":  appName,
		"email": s.From,
	}
	addresses := []map[string]string{}
	for _, address := range to {
		addresses = append(addresses, map[string]string{
			"email": address,
		})
	}
	data := map[string]interface{}{
		"sender":      sender,
		"to":          addresses,
		"subject":     subject,
		"htmlContent": text,
	}
	result := map[string]interface{}{}
	err := s.MakeSendInBlueRequest(ctx, data, &result)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return "", "", fmt.Errorf("unable to send email via sendInBlue: %w", err)
	}
	messageID, ok := result["messageId"].(string)
	if !ok {
		return "", "", fmt.Errorf("string messageID not found in SendInBlue response %#v", result)
	}
	return "ok", messageID, nil
}

// SendMailgun sends email via MailGun
func (s Service) SendMailgun(ctx context.Context, subject, text string, body *string, to ...string) (string, string, error) {
	_, span := tracer.Start(ctx, "SendMailgun")
	defer span.End()
	s.CheckPreconditions()
	message := s.Mg.NewMessage(s.From, subject, text, to...)
	if body != nil {
		message.SetHtml(*body)
	} else {
		message.SetHtml(text)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*MailGunTimeoutSeconds)
	defer cancel()

	resp, id, err := s.Mg.Send(ctx, message)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return resp, id, fmt.Errorf("mailgun email sending error: %s", err)
	}
	return resp, id, err
}

// SendEmail sends the specified email to the recipient(s) specified in `to`
// and returns the status
func (s Service) SendEmail(ctx context.Context, subject, text string, body *string, to ...string) (string, string, error) {
	s.CheckPreconditions()
	if s.SendInBlueEnabled {
		return s.SendInBlue(ctx, subject, text, to...)
	}
	return s.SendMailgun(ctx, subject, text, body, to...)
}

// SimpleEmail is a simplified API to send email.
// It returns only a status or error.
func (s Service) SimpleEmail(ctx context.Context, subject, text string, body *string, to ...string) (string, error) {
	s.CheckPreconditions()
	status, _, err := s.SendEmail(ctx, subject, text, nil, to...)
	return status, err
}
