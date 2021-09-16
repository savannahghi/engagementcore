package twilio_test

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/savannahghi/engagementcore/pkg/engagement/application/common/dto"
	"github.com/savannahghi/engagementcore/pkg/engagement/infrastructure"
	"github.com/savannahghi/engagementcore/pkg/engagement/infrastructure/services/otp"
	"github.com/savannahghi/engagementcore/pkg/engagement/infrastructure/services/twilio"
	twilioUsecase "github.com/savannahghi/engagementcore/pkg/engagement/usecases/twilio"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/serverutils"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	os.Setenv("ROOT_COLLECTION_SUFFIX", "testing")
	os.Setenv("ENVIRONMENT", "staging")
	os.Exit(m.Run())
}

func InitializeTestNewTwilio(ctx context.Context) (*twilioUsecase.ImplTwilio, infrastructure.Interactor, error) {
	infra := infrastructure.NewInteractor()
	twilio := twilioUsecase.NewImplTwilio(infra)
	return twilio, infra, nil
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

func TestNewImplTwilio(t *testing.T) {
	ctx := context.Background()
	f, i, err := InitializeTestNewTwilio(ctx)
	if err != nil {
		t.Errorf("failed to initialize test mail interractor: %v", err)
	}
	type args struct {
		infrastructure infrastructure.Interactor
	}
	tests := []struct {
		name string
		args args
		want *twilioUsecase.ImplTwilio
	}{
		{
			name: "default case",
			args: args{
				infrastructure: i,
			},
			want: f,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := twilioUsecase.NewImplTwilio(tt.args.infrastructure); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewImplTwilio() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestImplTwilio_Room(t *testing.T) {

	// A Room Can't be set up with test creds so for this test we make twilio creds live
	initialTwilioAuthToken, initialTwilioSID, err := setTwilioCredsToLive()
	if err != nil {
		t.Errorf("unable to set twilio credentials to live: %v", err)
		return
	}

	ctx := context.Background()
	f, _, err := InitializeTestNewTwilio(ctx)
	if err != nil {
		t.Errorf("failed to initialize test mail interractor: %v", err)
	}

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name      string
		args      args
		wantValue bool
		wantErr   bool
	}{
		{
			name: "default case",
			args: args{
				ctx: ctx,
			},
			wantErr:   false,
			wantValue: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := f.Room(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("ImplTwilio.Room() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantValue && got == nil {
				t.Errorf("ImplTwilio.Room(): expected to return a value but got nil")
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

func TestImplTwilio_TwilioAccessToken(t *testing.T) {

	// A Room Can't be set up with test creds so for this test we make twilio creds live
	initialTwilioAuthToken, initialTwilioSID, err := setTwilioCredsToLive()
	if err != nil {
		t.Errorf("unable to set twilio credentials to live: %v", err)
		return
	}

	ctx := firebasetools.GetAuthenticatedContext(t)
	f, _, err := InitializeTestNewTwilio(ctx)
	if err != nil {
		t.Errorf("failed to initialize test mail interractor: %v", err)
	}

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name      string
		args      args
		wantValue bool
		wantErr   bool
	}{
		{
			name: "valid: valid context",
			args: args{
				ctx: ctx,
			},
			wantErr:   false,
			wantValue: true,
		},
		{
			name: "invalid: invalid context",
			args: args{
				ctx: context.Background(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := f.TwilioAccessToken(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("ImplTwilio.TwilioAccessToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantValue && got == nil {
				t.Errorf("ImplTwilio.TwilioAccessToken(): expected to return a value but got nil")
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

func TestImplTwilio_SendSMS(t *testing.T) {
	// set test credentials
	initialSmsNumber := serverutils.MustGetEnvVar(twilio.TwilioSMSNumberEnvVarName)
	testSmsNumber := serverutils.MustGetEnvVar("TEST_TWILIO_SMS_NUMBER")
	os.Setenv(twilio.TwilioSMSNumberEnvVarName, testSmsNumber)

	ctx := context.Background()
	f, _, err := InitializeTestNewTwilio(ctx)
	if err != nil {
		t.Errorf("failed to initialize test mail interractor: %v", err)
	}
	type args struct {
		ctx context.Context
		to  string
		msg string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "good case",
			args: args{
				ctx: context.Background(),
				to:  testSmsNumber,
				msg: "Test message via Twilio",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := f.SendSMS(tt.args.ctx, tt.args.to, tt.args.msg); (err != nil) != tt.wantErr {
				t.Errorf("ImplTwilio.SendSMS() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	// restore twilio sms phone number
	err = os.Setenv(twilio.TwilioSMSNumberEnvVarName, initialSmsNumber)
	if err != nil {
		t.Errorf("unable to restore twilio sms number envar: %v", err)
	}
}

func TestImplTwilio_SaveTwilioVideoCallbackStatus(t *testing.T) {
	ctx := context.Background()
	f, _, err := InitializeTestNewTwilio(ctx)
	if err != nil {
		t.Errorf("failed to initialize test mail interractor: %v", err)
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
				if err := f.SaveTwilioVideoCallbackStatus(tt.args.ctx, tt.args.data); (err != nil) != tt.wantErr {
					t.Errorf("ImplTwilio.SaveTwilioVideoCallbackStatus() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
			if tt.panics {
				fcSaveTwilioVideoCallbackStatus := func() { _ = f.SaveTwilioVideoCallbackStatus(tt.args.ctx, tt.args.data) }
				assert.Panics(t, fcSaveTwilioVideoCallbackStatus)
			}
		})
	}
}

func TestImplTwilio_PhoneNumberVerificationCode(t *testing.T) {
	ctx := context.Background()
	f, _, err := InitializeTestNewTwilio(ctx)
	if err != nil {
		t.Errorf("failed to initialize test mail interractor: %v", err)
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
			got, err := f.PhoneNumberVerificationCode(tt.args.ctx, tt.args.to, tt.args.code, tt.args.marketingMessage)
			if (err != nil) != tt.wantErr {
				t.Errorf("ImplTwilio.PhoneNumberVerificationCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ImplTwilio.PhoneNumberVerificationCode() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestImplTwilio_SaveTwilioCallbackResponse(t *testing.T) {
	ctx := context.Background()
	f, _, err := InitializeTestNewTwilio(ctx)
	if err != nil {
		t.Errorf("failed to initialize test mail interractor: %v", err)
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
				if err := f.SaveTwilioCallbackResponse(tt.args.ctx, tt.args.data); (err != nil) != tt.wantErr {
					t.Errorf("ImplTwilio.SaveTwilioCallbackResponse() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
			if tt.panics {
				fcSaveTwilioCallbackResponse := func() { _ = f.SaveTwilioCallbackResponse(tt.args.ctx, tt.args.data) }
				assert.Panics(t, fcSaveTwilioCallbackResponse)
			}
		})
	}
}

func TestImplTwilio_TemporaryPIN(t *testing.T) {
	ctx := context.Background()
	f, _, err := InitializeTestNewTwilio(ctx)
	if err != nil {
		t.Errorf("failed to initialize test mail interractor: %v", err)
	}
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
			got, err := f.TemporaryPIN(tt.args.ctx, tt.args.to, tt.args.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("ImplTwilio.TemporaryPIN() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ImplTwilio.TemporaryPIN() = %v, want %v", got, tt.want)
			}
		})
	}
}
