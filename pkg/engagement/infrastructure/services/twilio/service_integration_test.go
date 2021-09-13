package twilio_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/savannahghi/engagementcore/pkg/engagement/application/common/dto"
	"github.com/savannahghi/engagementcore/pkg/engagement/infrastructure/database"
	"github.com/savannahghi/engagementcore/pkg/engagement/infrastructure/services/messaging"
	"github.com/savannahghi/engagementcore/pkg/engagement/infrastructure/services/otp"
	"github.com/savannahghi/engagementcore/pkg/engagement/infrastructure/services/sms"
	"github.com/savannahghi/engagementcore/pkg/engagement/infrastructure/services/twilio"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/serverutils"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	os.Setenv("ROOT_COLLECTION_SUFFIX", "testing")
	os.Exit(m.Run())
}

func newTwilioService(ctx context.Context) (*twilio.ServiceTwilioImpl, error) {
	var repo database.Repository
	projectID := serverutils.MustGetEnvVar(serverutils.GoogleCloudProjectIDEnvVarName)
	ns, err := messaging.NewPubSubNotificationService(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf(
			"can't instantiate notification service: %w",
			err,
		)
	}

	sms := sms.NewService(repo, ns)

	return twilio.NewService(*sms, repo), nil
}
func onboardingISCClient(t *testing.T) *interserviceclient.InterServiceClient {
	deps, err := interserviceclient.LoadDepsFromYAML()
	if err != nil {
		t.Errorf("can't load inter-service config from YAML: %v", err)
		return nil
	}

	profileClient, err := interserviceclient.SetupISCclient(*deps, "profile")
	if err != nil {
		t.Errorf("can't set up profile interservice client: %v", err)
		return nil
	}

	return profileClient
}

func TestNewService(t *testing.T) {
	srv, err := newTwilioService(context.Background())
	if err != nil {
		t.Errorf("failed to initialize new twilio test service: %v", err)
		return
	}
	assert.NotNil(t, srv)
	if srv == nil {
		t.Errorf("nil twilio service")
		return
	}
}

func setTwilioCredsToLive() (string, string, error) {
	initialTwilioAuthToken := serverutils.MustGetEnvVar("TWILIO_ACCOUNT_AUTH_TOKEN")
	initialTwilioSID := serverutils.MustGetEnvVar("TWILIO_ACCOUNT_SID")

	liveTwilioAuthToken := serverutils.MustGetEnvVar("TESTING_TWILIO_ACCOUNT_AUTH_TOKEN")
	liveTwilioSID := serverutils.MustGetEnvVar("TESTING_TWILIO_ACCOUNT_SID")

	err := os.Setenv("TWILIO_ACCOUNT_AUTH_TOKEN", liveTwilioAuthToken)
	if err != nil {
		return "", "", fmt.Errorf("unable to set twilio auth token to live: %v", err)
	}
	err = os.Setenv("TWILIO_ACCOUNT_SID", liveTwilioSID)
	if err != nil {
		return "", "", fmt.Errorf("unable to set test twilio auth token to live: %v", err)
	}

	return initialTwilioAuthToken, initialTwilioSID, nil
}

func restoreTwilioCreds(initialTwilioAuthToken, initialTwilioSID string) error {
	err := os.Setenv("TWILIO_ACCOUNT_AUTH_TOKEN", initialTwilioAuthToken)
	if err != nil {
		return fmt.Errorf("unable to restore twilio auth token: %v", err)
	}
	err = os.Setenv("TWILIO_ACCOUNT_SID", initialTwilioSID)
	if err != nil {
		return fmt.Errorf("unable to restore twilio sid: %v", err)
	}
	return nil
}

func TestService_Room(t *testing.T) {

	// A Room Can't be set up with test creds so for this test we make twilio creds live
	initialTwilioAuthToken, initialTwilioSID, err := setTwilioCredsToLive()
	if err != nil {
		t.Errorf("unable to set twilio credentials to live: %v", err)
		return
	}

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid test case",
			args: args{
				ctx: context.Background(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := newTwilioService(context.Background())
			if err != nil {
				t.Errorf("failed to initialize new twilio test service: %v", err)
				return
			}
			room, err := s.Room(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.Room() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if room == nil {
				t.Errorf("nil room")
				return
			}

			if tt.wantErr == false {
				if room.Type != "peer-to-peer" {
					t.Errorf("room.Type is not peer to peer")
					return
				}
			}
		})
	}

	// Restore envs after test
	err = restoreTwilioCreds(initialTwilioAuthToken, initialTwilioSID)
	if err != nil {
		t.Errorf("unable to restore twilio credentials: %v", err)
		return
	}
}

func TestService_AccessToken(t *testing.T) {

	// A Room Can't be set up with test creds so for this test we make twilio creds live
	initialTwilioAuthToken, initialTwilioSID, err := setTwilioCredsToLive()
	if err != nil {
		t.Errorf("unable to set twilio credentials to live: %v", err)
		return
	}

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid case",
			args: args{
				ctx: firebasetools.GetAuthenticatedContext(t),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := newTwilioService(context.Background())
			if err != nil {
				t.Errorf("failed to initialize new twilio test service: %v", err)
				return
			}
			got, err := s.TwilioAccessToken(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.AccessToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Errorf("nil AccessToken value got")
				return
			}
			if got.JWT == "" {
				t.Errorf("empty access token JWT value got")
				return
			}
			if got.UniqueName == "" {
				t.Errorf("empty access token Unique Name value got")
				return
			}
			if got.SID == "" {
				t.Errorf("empty access token SID value got")
				return
			}
			if got.DateUpdated.IsZero() {
				t.Errorf("empty access token Date Updated value got")
				return
			}
			if got.Status == "" {
				t.Errorf("empty access token Status value got")
				return
			}
			if got.Type == "" {
				t.Errorf("empty access token Type value got")
				return
			}
			if got.MaxParticipants == 0 {
				t.Errorf("empty access token Max Participants value got")
				return
			}
		})
	}

	// Restore envs after test
	err = restoreTwilioCreds(initialTwilioAuthToken, initialTwilioSID)
	if err != nil {
		t.Errorf("unable to restore twilio credentials: %v", err)
		return
	}

}

func TestService_SendSMS(t *testing.T) {

	// set test credentials
	initialSmsNumber := serverutils.MustGetEnvVar(twilio.TwilioSMSNumberEnvVarName)
	testSmsNumber := serverutils.MustGetEnvVar("TEST_TWILIO_SMS_NUMBER")
	os.Setenv(twilio.TwilioSMSNumberEnvVarName, testSmsNumber)

	type args struct {
		ctx                              context.Context
		normalizedDestinationPhoneNumber string
		msg                              string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "good case",
			args: args{
				ctx:                              context.Background(),
				normalizedDestinationPhoneNumber: testSmsNumber,
				msg:                              "Test message via Twilio",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := newTwilioService(context.Background())
			if err != nil {
				t.Errorf("failed to initialize new twilio test service: %v", err)
				return
			}
			if err := s.SendSMS(tt.args.ctx, tt.args.normalizedDestinationPhoneNumber, tt.args.msg); (err != nil) != tt.wantErr {
				t.Errorf("Service.SendSMS() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	// restore twilio sms phone number
	err := os.Setenv(twilio.TwilioSMSNumberEnvVarName, initialSmsNumber)
	if err != nil {
		t.Errorf("unable to restore twilio sms number envar: %v", err)
	}
}

func TestImplTwilio_SaveTwilioVideoCallbackStatus(t *testing.T) {
	onboardingClient := onboardingISCClient(t)
	ctx, token, err := interserviceclient.GetPhoneNumberAuthenticatedContextAndToken(t, onboardingClient)
	if err != nil {
		t.Errorf("cant get phone number authenticated context token: %v", err)
		return
	}
	_, err = firebasetools.GetAuthenticatedContextFromUID(ctx, token.UID)
	if err != nil {
		t.Errorf("cant get authenticated context from UID: %v", err)
		return
	}
	s, err := newTwilioService(ctx)
	if err != nil {
		t.Errorf("failed to initialize new twilio test service: %v", err)
		return
	}

	type args struct {
		ctx  context.Context
		data dto.CallbackData
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		panics  bool
	}{
		{
			name: "invalid: data not passed",
			args: args{
				ctx: ctx,
			},
			panics: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.panics {
				if err := s.SaveTwilioVideoCallbackStatus(tt.args.ctx, tt.args.data); (err != nil) != tt.wantErr {
					t.Errorf("ImplTwilio.SaveTwilioVideoCallbackStatus() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
			if tt.panics {
				fcSaveTwilioVideoCallbackStatus := func() { _ = s.SaveTwilioVideoCallbackStatus(tt.args.ctx, tt.args.data) }
				assert.Panics(t, fcSaveTwilioVideoCallbackStatus)
			}
		})
	}
}

func TestService_PhoneNumberVerificationCode(t *testing.T) {
	s, err := newTwilioService(context.Background())
	if err != nil {
		t.Errorf("failed to initialize new twilio test service: %v", err)
		return
	}
	type args struct {
		ctx              context.Context
		to               string
		code             string
		marketingMessage string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "invalid number",
			args: args{
				ctx: context.Background(),
				to:  "this is not a valid number",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "valid number",
			args: args{
				ctx:              context.Background(),
				to:               "+25423002959",
				code:             "345",
				marketingMessage: "This is a test",
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.PhoneNumberVerificationCode(tt.args.ctx, tt.args.to, tt.args.code, tt.args.marketingMessage)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.PhoneNumberVerificationCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Service.PhoneNumberVerificationCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestImplTwilio_SaveTwilioCallbackResponse(t *testing.T) {
	onboardingClient := onboardingISCClient(t)
	ctx, token, err := interserviceclient.GetPhoneNumberAuthenticatedContextAndToken(t, onboardingClient)
	if err != nil {
		t.Errorf("cant get phone number authenticated context token: %v", err)
		return
	}
	_, err = firebasetools.GetAuthenticatedContextFromUID(ctx, token.UID)
	if err != nil {
		t.Errorf("cant get authenticated context from UID: %v", err)
		return
	}
	s, err := newTwilioService(ctx)
	if err != nil {
		t.Errorf("failed to initialize new twilio test service: %v", err)
		return
	}
	type args struct {
		ctx  context.Context
		data dto.Message
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		panics  bool
	}{
		{
			name: "invalid: data not passed",
			args: args{
				ctx: ctx,
			},
			panics: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.panics {
				if err := s.SaveTwilioCallbackResponse(tt.args.ctx, tt.args.data); (err != nil) != tt.wantErr {
					t.Errorf("ImplTwilio.SaveTwilioCallbackResponse() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
			if tt.panics {
				fcSaveTwilioCallbackResponse := func() { _ = s.SaveTwilioCallbackResponse(tt.args.ctx, tt.args.data) }
				assert.Panics(t, fcSaveTwilioCallbackResponse)
			}
		})
	}
}

func TestService_TemporaryPIN(t *testing.T) {
	s, err := newTwilioService(context.Background())
	if err != nil {
		t.Errorf("failed to initialize new twilio test service: %v", err)
		return
	}
	ctx := context.Background()
	type args struct {
		ctx     context.Context
		to      string
		message string
	}

	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "sad invalid number",
			args: args{
				ctx: ctx,
				to:  "12345",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "happy sent temporary pin message",
			args: args{
				ctx:     ctx,
				to:      "+25423002959",
				message: fmt.Sprintf(otp.PINWhatsApp, "Test", "1234"),
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.TemporaryPIN(tt.args.ctx, tt.args.to, tt.args.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.TemporaryPIN() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Service.TemporaryPIN() = %v, want %v", got, tt.want)
			}
		})
	}
}
