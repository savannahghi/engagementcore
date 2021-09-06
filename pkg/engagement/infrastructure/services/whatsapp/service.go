package whatsapp

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/savannahghi/converterandformatter"
	"github.com/savannahghi/engagementcore/pkg/engagement/application/common/dto"
	"github.com/savannahghi/engagementcore/pkg/engagement/application/common/helpers"
	"github.com/savannahghi/engagementcore/pkg/engagement/infrastructure/database"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/serverutils"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("github.com/savannahghi/engagementcore/pkg/engagement/services/whatsapp")

// Twilio Whatsapp API contants
const (
	// TwilioHTTPClientTimeoutSeconds determines how long to wait (in seconds) before giving up on a
	// request to the Twilio API
	TwilioHTTPClientTimeoutSeconds = 30
	TwilioWhatsappSIDEnvVarName    = "TWILIO_WHATSAPP_SID"

	// gosec false positive
	TwilioWhatsappAuthTokenEnvVarName = "TWILIO_WHATSAPP_AUTH_TOKEN" /* #nosec */

	TwilioWhatsappSenderEnvVarName = "TWILIO_WHATSAPP_SENDER"

	twilioWhatsappBaseURL = "https://api.twilio.com/2010-04-01/Accounts/"
)

// NewService initializes a properly set up WhatsApp service
func NewService() *ServiceWhatsappImpl {
	sid := serverutils.MustGetEnvVar(TwilioWhatsappSIDEnvVarName)
	authToken := serverutils.MustGetEnvVar(TwilioWhatsappAuthTokenEnvVarName)
	sender := serverutils.MustGetEnvVar(TwilioWhatsappSenderEnvVarName)
	httpClient := &http.Client{
		Timeout: time.Second * TwilioHTTPClientTimeoutSeconds,
	}
	return &ServiceWhatsappImpl{
		BaseURL:          twilioWhatsappBaseURL,
		AccountSID:       sid,
		AccountAuthToken: authToken,
		Sender:           sender,
		HTTPClient:       httpClient,
	}
}

// ServiceWhatsapp defines the interactions with the whatsapp service
type ServiceWhatsapp interface {
	PhoneNumberVerificationCode(
		ctx context.Context,
		to string,
		code string,
		marketingMessage string,
	) (bool, error)

	// TODO: Remove db implementation
	SaveTwilioCallbackResponse(
		ctx context.Context,
		data dto.Message,
	) error

	TemporaryPIN(
		ctx context.Context,
		to string,
		message string,
	) (bool, error)
}

// ServiceWhatsappImpl is a WhatsApp service. The receivers implement the query and mutation resolvers.
type ServiceWhatsappImpl struct {
	BaseURL          string
	AccountSID       string
	AccountAuthToken string
	Sender           string
	HTTPClient       *http.Client
	Repository       database.Repository
}

// CheckPreconditions ...
func (s ServiceWhatsappImpl) CheckPreconditions() {
	if s.HTTPClient == nil {
		log.Panicf("nil http client in Twilio WhatsApp service")
	}

	if s.BaseURL == "" {
		log.Panicf("blank base URL in Twilio WhatsApp service")
	}

	if s.AccountSID == "" {
		log.Panicf("blank accountSID in Twilio WhatsApp service")
	}

	if s.AccountAuthToken == "" {
		log.Panicf("blank account auth token in Twilio WhatsApp service")
	}

	if s.Sender == "" {
		log.Panicf("blank sender in Twilio WhatsApp service")
	}
}

// MakeTwilioRequest makes a twilio request
func (s ServiceWhatsappImpl) MakeTwilioRequest(
	ctx context.Context,
	method string,
	urlPath string,
	content url.Values,
	target interface{},
) error {
	_, span := tracer.Start(ctx, "MakeTwilioRequest")
	defer span.End()
	s.CheckPreconditions()

	if serverutils.IsDebug() {
		log.Printf("Twilio request data: \n%s\n", content)
	}

	r := strings.NewReader(content.Encode())
	req, err := http.NewRequest(method, s.BaseURL+urlPath, r)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(s.AccountSID, s.AccountAuthToken)

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return fmt.Errorf("twilio API error: %w", err)
	}

	respBs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return fmt.Errorf("twilio room content error: %w", err)
	}

	if resp.StatusCode > 201 {
		return fmt.Errorf("twilio API Error: %s", string(respBs))
	}

	if serverutils.IsDebug() {
		log.Printf("Twilio response: \n%s\n", string(respBs))
	}
	err = json.Unmarshal(respBs, target)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return fmt.Errorf("unable to unmarshal Twilio resp: %w", err)
	}

	return nil
}

// PhoneNumberVerificationCode sends Phone Number verification codes via WhatsApp
func (s ServiceWhatsappImpl) PhoneNumberVerificationCode(
	ctx context.Context,
	to string,
	code string,
	marketingMessage string,
) (bool, error) {
	ctx, span := tracer.Start(ctx, "PhoneNumberVerificationCode")
	defer span.End()
	s.CheckPreconditions()

	normalizedPhoneNo, err := converterandformatter.NormalizeMSISDN(to)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return false, fmt.Errorf("%s is not a valid E164 phone number: %w", to, err)
	}

	msgFrom := fmt.Sprintf("whatsapp:%s", s.Sender)
	msgTo := fmt.Sprintf("whatsapp:%s", *normalizedPhoneNo)
	msg := fmt.Sprintf("Your phone number verification code is %s", code)

	payload := url.Values{}
	payload.Add("From", msgFrom)
	payload.Add("Body", msg)
	payload.Add("To", msgTo)

	target := dto.Message{}
	path := fmt.Sprintf("%s/Messages.json", s.AccountSID)
	err = s.MakeTwilioRequest(
		ctx,
		http.MethodPost,
		path,
		payload,
		&target,
	)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return false, fmt.Errorf("error from Twilio: %w", err)
	}

	// save Twilio response for audit purposes
	_, _, err = firebasetools.CreateNode(ctx, &target)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return false, fmt.Errorf("unable to save Twilio response: %w", err)
	}
	// TODO Find out why /ide is not working (401s)
	// TODO deploy UAT, deploy prod, tag (semver)
	return true, nil
}

// SaveTwilioCallbackResponse saves the twilio callback response for future
// analysis
func (s ServiceWhatsappImpl) SaveTwilioCallbackResponse(
	ctx context.Context,
	data dto.Message,
) error {
	return s.Repository.SaveTwilioResponse(ctx, data)
}

//TemporaryPIN send PIN via whatsapp to user
func (s ServiceWhatsappImpl) TemporaryPIN(
	ctx context.Context,
	to string,
	message string,
) (bool, error) {
	ctx, span := tracer.Start(ctx, "TemporaryPIN")
	defer span.End()

	s.CheckPreconditions()

	normalizedPhoneNo, err := converterandformatter.NormalizeMSISDN(to)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return false, fmt.Errorf("%s is not a valid E164 phone number: %w", to, err)
	}

	msgFrom := fmt.Sprintf("whatsapp:%s", s.Sender)
	msgTo := fmt.Sprintf("whatsapp:%s", *normalizedPhoneNo)

	payload := url.Values{}
	payload.Add("From", msgFrom)
	payload.Add("Body", message)
	payload.Add("To", msgTo)

	target := dto.Message{}
	path := fmt.Sprintf("%s/Messages.json", s.AccountSID)

	err = s.MakeTwilioRequest(ctx, http.MethodPost, path, payload, &target)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return false, fmt.Errorf("error from Twilio: %w", err)
	}

	// save Twilio response for audit purposes
	_, _, err = firebasetools.CreateNode(ctx, &target)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return false, fmt.Errorf("unable to save Twilio response: %w", err)
	}
	return true, nil
}
