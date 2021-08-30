package sms

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/savannahghi/engagement/pkg/engagement/application/common/dto"
	"github.com/savannahghi/engagement/pkg/engagement/application/common/helpers"
	"github.com/savannahghi/engagement/pkg/engagement/infrastructure/database"
	"github.com/savannahghi/engagement/pkg/engagement/infrastructure/services/messaging"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/serverutils"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("github.com/savannahghi/engagement/pkg/engagement/services/sms")

// AIT environment variables names
const (
	APIKeyEnvVarName       = "AIT_API_KEY"
	APIUsernameEnvVarName  = "AIT_USERNAME"
	APISenderIDEnvVarName  = "AIT_SENDER_ID"
	AITEnvVarName          = "AIT_ENVIRONMENT"
	BeWellAITAPIKey        = "AIT_BEWELL_API_KEY"
	BeWellAITUsername      = "AIT_BEWELL_USERNAME"
	BeWellAITSenderID      = "AIT_BEWELL_SENDER_ID"
	AITAuthenticationError = "The supplied authentication is invalid"

	//AITCallbackCollectionName is the name of a Cloud Firestore collection into which AIT
	// callback data will be saved for future analysis
	AITCallbackCollectionName = "ait_callbacks"
)

// ServiceSMS defines the interactions with sms service
type ServiceSMS interface {
	SendToMany(
		ctx context.Context,
		message string,
		to []string,
		from enumutils.SenderID,
	) (*dto.SendMessageResponse, error)
	Send(
		ctx context.Context,
		to, message string,
		from enumutils.SenderID,
	) (*dto.SendMessageResponse, error)

	// TODO: Remove DB specific implementations
	SaveMarketingMessage(
		ctx context.Context,
		data dto.MarketingSMS,
	) (*dto.MarketingSMS, error)
	UpdateMarketingMessage(
		ctx context.Context,
		data *dto.MarketingSMS,
	) (*dto.MarketingSMS, error)
	GetMarketingSMSByPhone(
		ctx context.Context,
		phoneNumber string,
	) (*dto.MarketingSMS, error)
}

// Service defines a sms service struct
type Service struct {
	Env        string
	Repository database.Repository
	PubSub     messaging.NotificationService
}

// GetSmsURL is the sms endpoint
func GetSmsURL(env string) string {
	return GetAPIHost(env) + "/version1/messaging"
}

// GetAPIHost returns either sandbox or prod
func GetAPIHost(env string) string {
	return getHost(env, "api")
}

func getHost(env, service string) string {
	if env != "sandbox" {
		return fmt.Sprintf("https://%s.africastalking.com", service)
	}
	return fmt.Sprintf(
		"https://%s.sandbox.africastalking.com",
		service,
	)

}

// NewService returns a new service
func NewService(
	repository database.Repository,
	pubsub messaging.NotificationService,
) *Service {
	env := serverutils.MustGetEnvVar(AITEnvVarName)
	return &Service{env, repository, pubsub}
}

// SaveMarketingMessage saves the callback data for future analysis
func (s Service) SaveMarketingMessage(
	ctx context.Context,
	data dto.MarketingSMS,
) (*dto.MarketingSMS, error) {
	return s.Repository.SaveMarketingMessage(ctx, data)
}

// UpdateMarketingMessage adds a delivery report to an AIT SMS
func (s Service) UpdateMarketingMessage(
	ctx context.Context,
	data *dto.MarketingSMS,
) (*dto.MarketingSMS, error) {
	return s.Repository.UpdateMarketingMessage(ctx, data)
}

// GetMarketingSMSByPhone returns the latest message given a phone number
func (s Service) GetMarketingSMSByPhone(
	ctx context.Context,
	phoneNumber string,
) (*dto.MarketingSMS, error) {
	return s.Repository.GetMarketingSMSByPhone(ctx, phoneNumber)
}

// SendToMany is a utility method to send to many recipients at the same time
func (s Service) SendToMany(
	ctx context.Context,
	message string,
	to []string,
	from enumutils.SenderID,
) (*dto.SendMessageResponse, error) {
	recipients := strings.Join(to, ",")
	return s.Send(ctx, recipients, message, from)
}

// Send is a method used to send to a single recipient
func (s Service) Send(
	ctx context.Context,
	to, message string,
	from enumutils.SenderID,
) (*dto.SendMessageResponse, error) {

	switch from {
	case enumutils.SenderIDSLADE360:
		return s.SendSMS(
			ctx,
			to,
			message,
			serverutils.MustGetEnvVar(APISenderIDEnvVarName),
			serverutils.MustGetEnvVar(APIUsernameEnvVarName),
			serverutils.MustGetEnvVar(APIKeyEnvVarName),
		)
	case enumutils.SenderIDBewell:
		return s.SendSMS(
			ctx,
			to,
			message,
			serverutils.MustGetEnvVar(BeWellAITSenderID),
			serverutils.MustGetEnvVar(BeWellAITUsername),
			serverutils.MustGetEnvVar(BeWellAITAPIKey),
		)
	}
	return nil, fmt.Errorf("unknown AIT sender")
}

// SendSMS is a method used to send SMS
func (s Service) SendSMS(
	ctx context.Context,
	to, message string,
	from string,
	username string,
	key string,
) (*dto.SendMessageResponse, error) {
	ctx, span := tracer.Start(ctx, "SendSMS")
	defer span.End()
	values := url.Values{}
	values.Set("username", username)
	values.Set("to", to)
	values.Set("message", message)
	values.Set("from", from)

	smsURL := GetSmsURL(s.Env)
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"

	res, err := s.newPostRequest(ctx, smsURL, values, headers, key)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, err
	}
	defer func() {
		_ = res.Body.Close()
	}()

	// read the response body to a variable
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, err
	}

	smsMessageResponse := &dto.SendMessageResponse{}

	bodyAsString := string(bodyBytes)
	if strings.Contains(bodyAsString, AITAuthenticationError) {
		// return so that other processes don't break
		log.Println("AIT Authentication error encountered")
		return smsMessageResponse, nil
	}

	// reset the response body to the original unread state so that decode can
	// continue
	res.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	if err := json.NewDecoder(res.Body).Decode(smsMessageResponse); err != nil {
		return nil, errors.New(
			fmt.Errorf("SMS Error : unable to parse sms response ; %v", err).
				Error(),
		)
	}

	return smsMessageResponse, nil
}

func (s Service) newPostRequest(
	ctx context.Context,
	url string,
	values url.Values,
	headers map[string]string,
	key string,
) (*http.Response, error) {
	_, span := tracer.Start(ctx, "newPostRequest")
	defer span.End()
	reader := strings.NewReader(values.Encode())

	req, err := http.NewRequest(http.MethodPost, url, reader)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}
	req.Header.Set("Content-Length", strconv.Itoa(reader.Len()))
	req.Header.Set("apikey", key)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	return client.Do(req)
}
