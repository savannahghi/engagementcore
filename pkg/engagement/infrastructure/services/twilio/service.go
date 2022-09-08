package twilio

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/kevinburke/twilio-go"
	"github.com/kevinburke/twilio-go/token"
	"github.com/savannahghi/converterandformatter"
	"github.com/savannahghi/engagementcore/pkg/engagement/application/common/dto"
	"github.com/savannahghi/engagementcore/pkg/engagement/application/common/helpers"
	"github.com/savannahghi/engagementcore/pkg/engagement/infrastructure/database"
	"github.com/savannahghi/engagementcore/pkg/engagement/infrastructure/services/sms"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/serverutils"
	"go.opentelemetry.io/otel"
	"moul.io/http2curl"
)

var tracer = otel.Tracer("github.com/savannahghi/engagementcore/pkg/engagement/services/twilio")

/* #nosec */
// DefaultTwilioRegion is set to global low latency auto-selection
const (
	TwilioRegionEnvVarName            = "TWILIO_REGION"
	TwilioVideoAPIURLEnvVarName       = "TWILIO_VIDEO_API_URL"
	TwilioVideoAPIKeySIDEnvVarName    = "TWILIO_VIDEO_SID"
	TwilioVideoAPIKeySecretEnvVarName = "TWILIO_VIDEO_SECRET"
	TwilioAccountSIDEnvVarName        = "TWILIO_ACCOUNT_SID"
	TwilioAccountAuthTokenEnvVarName  = "TWILIO_ACCOUNT_AUTH_TOKEN"
	TwilioSMSNumberEnvVarName         = "TWILIO_SMS_NUMBER"
	ServerPublicDomainEnvVarName      = "SERVER_PUBLIC_DOMAIN"
	TwilioCallbackPath                = "/twilio_callback"
	TwilioHTTPClientTimeoutSeconds    = 10
	TwilioPeerToPeerMaxParticipants   = 3
	TwilioAccessTokenTTL              = 14400

	TwilioWhatsappSIDEnvVarName = "TWILIO_WHATSAPP_SID"

	// gosec false positive
	TwilioWhatsappAuthTokenEnvVarName = "TWILIO_WHATSAPP_AUTH_TOKEN" /* #nosec */

	TwilioWhatsappSenderEnvVarName = "TWILIO_WHATSAPP_SENDER"

	twilioWhatsappBaseURL = "https://api.twilio.com/2010-04-01/Accounts/"
)

// ServiceTwilio defines the interaction with the twilio service
type ServiceTwilio interface {
	Room(ctx context.Context) (*dto.Room, error)

	TwilioAccessToken(ctx context.Context) (*dto.AccessToken, error)

	SendSMS(ctx context.Context, to string, msg string) error

	// TODO: Remove db call
	SaveTwilioVideoCallbackStatus(
		ctx context.Context,
		data dto.CallbackData,
	) error

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

	MakeTwilioRequest(
		method string,
		urlPath string,
		content url.Values,
		target interface{},
	) error

	MakeWhatsappTwilioRequest(
		ctx context.Context,
		method string,
		urlPath string,
		content url.Values,
		target interface{},
	) error
}

// NewService initializes a service to interact with Twilio
func NewService(sms sms.ServiceSMS, repo database.Repository) *ServiceTwilioImpl {
	region := serverutils.MustGetEnvVar(TwilioRegionEnvVarName)
	videoBaseURL := serverutils.MustGetEnvVar(TwilioVideoAPIURLEnvVarName)
	videoAPIKeySID := serverutils.MustGetEnvVar(TwilioVideoAPIKeySIDEnvVarName)
	videoAPIKeySecret := serverutils.MustGetEnvVar(TwilioVideoAPIKeySecretEnvVarName)
	accountSID := serverutils.MustGetEnvVar(TwilioAccountSIDEnvVarName)
	accountAuthToken := serverutils.MustGetEnvVar(TwilioAccountAuthTokenEnvVarName)
	httpClient := &http.Client{
		Timeout: time.Second * TwilioHTTPClientTimeoutSeconds,
	}
	publicDomain := serverutils.MustGetEnvVar(ServerPublicDomainEnvVarName)
	callbackURL := publicDomain + TwilioCallbackPath
	smsNumber := serverutils.MustGetEnvVar(TwilioSMSNumberEnvVarName)

	sid := serverutils.MustGetEnvVar(TwilioWhatsappSIDEnvVarName)
	authToken := serverutils.MustGetEnvVar(TwilioWhatsappAuthTokenEnvVarName)
	sender := serverutils.MustGetEnvVar(TwilioWhatsappSenderEnvVarName)

	srv := &ServiceTwilioImpl{
		region:             region,
		videoBaseURL:       videoBaseURL,
		videoAPIKeySID:     videoAPIKeySID,
		videoAPIKeySecret:  videoAPIKeySecret,
		accountSID:         accountSID,
		accountAuthToken:   accountAuthToken,
		twilioClient:       twilio.NewClient(accountSID, accountAuthToken, httpClient),
		callbackURL:        callbackURL,
		smsNumber:          smsNumber,
		sms:                sms,
		BaseURL:            twilioWhatsappBaseURL,
		WhatsappAccountSID: sid,
		AccountAuthToken:   authToken,
		Sender:             sender,
		HTTPClient:         httpClient,
	}
	srv.CheckPreconditions()
	return srv
}

// ServiceTwilioImpl organizes methods needed to interact with Twilio for video, voice
// and text
type ServiceTwilioImpl struct {
	region             string
	videoBaseURL       string
	videoAPIKeySID     string
	videoAPIKeySecret  string
	accountSID         string
	accountAuthToken   string
	twilioClient       *twilio.Client
	callbackURL        string
	smsNumber          string
	sms                sms.ServiceSMS
	BaseURL            string
	WhatsappAccountSID string
	AccountAuthToken   string
	Sender             string
	HTTPClient         *http.Client
	Repository         database.Repository
}

// CheckPreconditions checks preconditions for the twilio service
func (s ServiceTwilioImpl) CheckPreconditions() {
	if s.region == "" {
		log.Panicf("Twilio region not set")
	}

	if s.videoBaseURL == "" {
		log.Panicf("Twilio video base URL not set")
	}

	if !govalidator.IsURL(s.videoBaseURL) {
		log.Panicf("Twilio Video base URL (%s) is not a valid URL", s.videoBaseURL)
	}

	if s.videoAPIKeySID == "" {
		log.Panicf("Twilio Video API Key SID not set")
	}

	if s.videoAPIKeySecret == "" {
		log.Panicf("Twilio Video API Key secret not set")
	}

	if s.accountSID == "" {
		log.Panicf("Twilio Video account SID not set")
	}

	if s.accountAuthToken == "" {
		log.Panicf("Twilio Video account auth token not set")
	}

	if s.twilioClient == nil {
		log.Panicf("nil Twilio client in Twilio service")
	}

	if s.callbackURL == "" {
		log.Panicf("empty Twilio callback URL")
	}
	if s.HTTPClient == nil {
		log.Panicf("nil http client in Twilio WhatsApp service")
	}

	if s.BaseURL == "" {
		log.Panicf("blank base URL in Twilio WhatsApp service")
	}

	if s.WhatsappAccountSID == "" {
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
func (s ServiceTwilioImpl) MakeTwilioRequest(
	method string,
	urlPath string,
	content url.Values,
	target interface{},
) error {
	s.CheckPreconditions()
	if serverutils.IsDebug() {
		log.Printf("Twilio request data: \n%s\n", content)
	}

	r := strings.NewReader(content.Encode())
	req, reqErr := http.NewRequest(method, s.videoBaseURL+urlPath, r)
	if reqErr != nil {
		return reqErr
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(s.accountSID, s.accountAuthToken)

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("twilio API error: %w", err)
	}

	respBs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
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
		return fmt.Errorf("unable to unmarshal Twilio resp: %w", err)
	}

	command, _ := http2curl.GetCurlCommand(req)
	fmt.Println(command)

	return nil
}

// MakeWhatsappTwilioRequest makes a twilio request
func (s ServiceTwilioImpl) MakeWhatsappTwilioRequest(
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
	req.SetBasicAuth(s.WhatsappAccountSID, s.AccountAuthToken)

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

// Room represents a real-time audio, data, video, and/or screen-share session,
// and is the basic building block for a Programmable Video application.
//
// In a Peer-to-peer Room, media flows directly between participants. This
// supports up to 10 participants in a mesh topology.
//
// In a Group Room, media is routed through Twilio's Media Servers. This
// supports up to 50 participants.
//
// Participants represent client applications that are connected to a Room and
// sharing audio, data, and/or video media with one another.
//
// Tracks represent the individual audio, data, and video media streams that
// are shared within a Room.
//
// LocalTracks represent the audio, data, and video captured from the local
// client's media sources (for example, microphone and camera).
//
// RemoteTracks represent the audio, data, and video tracks from other
// participants connected to the Room.
//
// Room names must be unique within an account.
//
// Rooms created via the REST API exist for five minutes to allow the first
// Participant to connect. If no Participants join within five minutes,
// the Room times out and a new Room must be created.
//
// Because of confidentiality issues in healthcare, we do not enable recording
// for these meetings.
func (s ServiceTwilioImpl) Room(ctx context.Context) (*dto.Room, error) {
	_, span := tracer.Start(ctx, "Room")
	defer span.End()
	s.CheckPreconditions()

	roomReqData := url.Values{}
	roomReqData.Set("Type", "peer-to-peer")
	roomReqData.Set("MaxParticipants", strconv.Itoa(TwilioPeerToPeerMaxParticipants))
	roomReqData.Set("StatusCallbackMethod", "POST")
	roomReqData.Set("StatusCallback", s.callbackURL)
	roomReqData.Set("EnableTurn", strconv.FormatBool(true))

	var room dto.Room
	err := s.MakeTwilioRequest("POST", "/v1/Rooms", roomReqData, &room)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf("twilio room API call error: %w", err)
	}
	return &room, nil
}

// TwilioAccessToken is used to generate short-lived credentials used to authenticate
// the client-side application to Twilio.
//
// An access token should be generated for every user of the application.
//
// An access token can optionally encode a room name, which would allow the user
// to connect only to the room specified in the token.
//
// Access tokens are JSON Web Tokens (JWTs).
func (s ServiceTwilioImpl) TwilioAccessToken(ctx context.Context) (*dto.AccessToken, error) {
	ctx, span := tracer.Start(ctx, "TwilioAccessToken")
	defer span.End()
	s.CheckPreconditions()

	uid, err := firebasetools.GetLoggedInUserUID(ctx)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf("unable to get logged in user uid: %w", err)
	}

	room, err := s.Room(ctx)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf("unable to get room to issue a grant to: %w", err)
	}

	ttl := time.Second * TwilioAccessTokenTTL
	accessToken := token.New(s.accountSID, s.videoAPIKeySID, s.videoAPIKeySecret, uid, ttl)
	videoGrant := token.NewVideoGrant(room.SID)
	accessToken.AddGrant(videoGrant)

	jwt, err := accessToken.JWT()
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf("unable to generate JWT for Twilio access token: %w", err)
	}
	payload := dto.AccessToken{
		JWT:             jwt,
		UniqueName:      room.UniqueName,
		SID:             room.SID,
		DateUpdated:     room.DateUpdated,
		Status:          room.Status,
		Type:            room.Type,
		MaxParticipants: room.MaxParticipants,
		Duration:        room.Duration,
	}
	return &payload, nil
}

// SendSMS sends a text message through Twilio's programmable SMS
func (s ServiceTwilioImpl) SendSMS(ctx context.Context, to string, msg string) error {
	_, span := tracer.Start(ctx, "SendSMS")
	defer span.End()
	s.CheckPreconditions()

	t, err := s.twilioClient.Messages.SendMessage(s.smsNumber, to, msg, nil)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return fmt.Errorf("twilio SMS API error: %w", err)
	}

	if t.ErrorCode != 0 {
		return fmt.Errorf("sms could not be sent: %v", t.ErrorMessage)
	}

	fmt.Printf("Raw Twilio SMS response: %v", t)
	return nil
}

// SaveTwilioVideoCallbackStatus saves status callback data
func (s ServiceTwilioImpl) SaveTwilioVideoCallbackStatus(
	ctx context.Context,
	data dto.CallbackData,
) error {
	return s.Repository.SaveTwilioVideoCallbackStatus(ctx, data)
}

// PhoneNumberVerificationCode sends Phone Number verification codes via WhatsApp
func (s ServiceTwilioImpl) PhoneNumberVerificationCode(
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
	path := fmt.Sprintf("%s/Messages.json", s.WhatsappAccountSID)
	err = s.MakeWhatsappTwilioRequest(
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
func (s ServiceTwilioImpl) SaveTwilioCallbackResponse(
	ctx context.Context,
	data dto.Message,
) error {
	return s.Repository.SaveTwilioResponse(ctx, data)
}

// TemporaryPIN send PIN via whatsapp to user
func (s ServiceTwilioImpl) TemporaryPIN(
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
	path := fmt.Sprintf("%s/Messages.json", s.WhatsappAccountSID)

	err = s.MakeWhatsappTwilioRequest(
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
	return true, nil
}
