package otp_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/savannahghi/converterandformatter"
	"github.com/savannahghi/engagementcore/pkg/engagement/application/common/dto"
	"github.com/savannahghi/engagementcore/pkg/engagement/infrastructure/services/mail"
	mailMock "github.com/savannahghi/engagementcore/pkg/engagement/infrastructure/services/mail/mock"
	"github.com/savannahghi/engagementcore/pkg/engagement/infrastructure/services/otp"
	otpMock "github.com/savannahghi/engagementcore/pkg/engagement/infrastructure/services/otp/mock"
	"github.com/savannahghi/engagementcore/pkg/engagement/infrastructure/services/sms"
	smsMock "github.com/savannahghi/engagementcore/pkg/engagement/infrastructure/services/sms/mock"
	"github.com/savannahghi/engagementcore/pkg/engagement/infrastructure/services/twilio"
	twilioMock "github.com/savannahghi/engagementcore/pkg/engagement/infrastructure/services/twilio/mock"
	"github.com/savannahghi/silcomms"

	"github.com/savannahghi/interserviceclient"
	"github.com/stretchr/testify/assert"
)

var fakeMail mailMock.FakeServiceMail
var mailSvc mail.ServiceMail = &fakeMail

var fakeSMS smsMock.FakeServiceSMS
var smsSvs sms.ServiceSMS = &fakeSMS

var fakeTwilio twilioMock.FakeServiceTwilio
var twilioSvc twilio.ServiceTwilio = &fakeTwilio

var fakeOTP otpMock.FakeServiceOTP
var OTPSvs otp.ServiceOTP = &fakeOTP

const (
	InternationalTestUserPhoneNumber = "+12028569601"
	ValidTestEmail                   = "automated.test.user.bewell-app-ci@healthcloud.co.ke"
)

func TestMain(m *testing.M) {
	os.Setenv("ROOT_COLLECTION_SUFFIX", "staging")
	os.Setenv("ENVIRONMENT", "staging")
	os.Exit(m.Run())
}

func TestNormalizeMSISDN(t *testing.T) {
	tests := map[string]struct {
		input                string
		expectError          bool
		expectedErrorMessage string
		expectedOutput       string
	}{
		"invalid_number": {
			input:                "1",
			expectError:          true,
			expectedErrorMessage: "invalid phone number: 1",
		},
		"international_with_plus": {
			input:          interserviceclient.TestUserPhoneNumber,
			expectError:    false,
			expectedOutput: interserviceclient.TestUserPhoneNumber,
		},
		"international_no_plus": {
			input:          interserviceclient.TestUserPhoneNumber[1:],
			expectError:    false,
			expectedOutput: interserviceclient.TestUserPhoneNumber,
		},
		"national_zero_prefix": {
			input:          "0" + interserviceclient.TestUserPhoneNumber[4:],
			expectError:    false,
			expectedOutput: interserviceclient.TestUserPhoneNumber,
		},
		"national_no_zero_prefix": {
			input:                interserviceclient.TestUserPhoneNumber[4:],
			expectError:          true,
			expectedErrorMessage: "invalid phone number: 711223344",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := converterandformatter.NormalizeMSISDN(tc.input)
			if tc.expectError {
				assert.NotNil(t, err)
				assert.Equal(t, tc.expectedErrorMessage, err.Error())
			}

			if !tc.expectError {
				assert.Nil(t, err)
				assert.Equal(t, tc.expectedOutput, *got)
			}
		})
	}
}

func TestService_GenerateAndSendOTP(t *testing.T) {
	ctx := context.Background()
	service := otp.NewService(mailSvc, smsSvs, twilioSvc)
	appID := uuid.New().String()
	type args struct {
		ctx    context.Context
		msisdn string
		appID  *string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Happy Case -> Successfully generate and send OTP",
			args: args{
				ctx:    ctx,
				msisdn: interserviceclient.TestUserPhoneNumber,
				appID:  &appID,
			},
			wantErr: false,
		},
		{
			name: "Sad Case -> Fail to send OTP",
			args: args{
				ctx:    ctx,
				msisdn: interserviceclient.TestUserPhoneNumber,
				appID:  &appID,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Happy Case -> Successfully generate and send OTP" {
				fakeOTP.GenerateOTPFn = func(ctx context.Context) (string, error) {
					return "1234", nil
				}

				fakeOTP.SaveOTPToFirestoreFn = func(otp dto.OTP) error {
					return nil
				}

				fakeSMS.SendFn = func(
					ctx context.Context,
					to, message string,
				) (*silcomms.BulkSMSResponse, error) {
					return &silcomms.BulkSMSResponse{}, nil
				}
			}

			if tt.name == "Sad Case -> Fail to send OTP" {
				fakeOTP.GenerateOTPFn = func(ctx context.Context) (string, error) {
					return "1234", nil
				}

				fakeOTP.SaveOTPToFirestoreFn = func(otp dto.OTP) error {
					return nil
				}

				fakeSMS.SendFn = func(
					ctx context.Context,
					to, message string,
				) (*silcomms.BulkSMSResponse, error) {
					return nil, fmt.Errorf("failed to send OTP")
				}
			}

			got, err := service.GenerateAndSendOTP(tt.args.ctx, tt.args.msisdn, tt.args.appID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GenerateAndSendOTP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == "" {
				t.Errorf("Service.GenerateAndSendOTP() = Expected an OTP to be returned")
			}
		})
	}
}

func TestService_SendOTPToEmail(t *testing.T) {
	ctx := context.Background()
	service := otp.NewService(mailSvc, smsSvs, twilioSvc)
	validEmail := ValidTestEmail
	phoneNumber := interserviceclient.TestUserPhoneNumber
	appID := uuid.New().String()
	type args struct {
		ctx    context.Context
		msisdn *string
		email  *string
		appID  *string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Happy Case -> Send OTP to Email",
			args: args{
				ctx:    ctx,
				msisdn: &phoneNumber,
				email:  &validEmail,
				appID:  &appID,
			},
			wantErr: false,
		},
		{
			name: "Sad Case -> Fail to Send OTP to Email",
			args: args{
				ctx:    ctx,
				msisdn: &phoneNumber,
				email:  &validEmail,
				appID:  &appID,
			},
			wantErr: true,
		},
		{
			name: "Sad Case -> Fail to generate and OTP",
			args: args{
				ctx:    ctx,
				msisdn: &phoneNumber,
				email:  &validEmail,
				appID:  &appID,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Happy Case -> Send OTP to Email" {
				fakeOTP.GenerateOTPFn = func(ctx context.Context) (string, error) {
					return "1234", nil
				}

				fakeOTP.SaveOTPToFirestoreFn = func(otp dto.OTP) error {
					return nil
				}

				fakeSMS.SendFn = func(
					ctx context.Context,
					to, message string,
				) (*silcomms.BulkSMSResponse, error) {
					return &silcomms.BulkSMSResponse{}, nil
				}

				fakeMail.SendEmailFn = func(
					ctx context.Context,
					subject, text string,
					body *string,
					to ...string,
				) (string, string, error) {
					return "", "", nil
				}
			}

			if tt.name == "Sad Case -> Fail to Send OTP to Email" {
				fakeOTP.GenerateOTPFn = func(ctx context.Context) (string, error) {
					return "1234", nil
				}

				fakeOTP.SaveOTPToFirestoreFn = func(otp dto.OTP) error {
					return nil
				}

				fakeSMS.SendFn = func(
					ctx context.Context,
					to, message string,
				) (*silcomms.BulkSMSResponse, error) {
					return &silcomms.BulkSMSResponse{}, nil
				}

				fakeMail.SendEmailFn = func(
					ctx context.Context,
					subject, text string,
					body *string,
					to ...string,
				) (string, string, error) {
					return "", "", fmt.Errorf("failed to send OTP via Email")
				}
			}

			if tt.name == "Sad Case -> Fail to generate and OTP" {
				fakeOTP.GenerateOTPFn = func(ctx context.Context) (string, error) {
					return "1234", nil
				}

				fakeOTP.SaveOTPToFirestoreFn = func(otp dto.OTP) error {
					return nil
				}

				fakeSMS.SendFn = func(
					ctx context.Context,
					to, message string,
				) (*silcomms.BulkSMSResponse, error) {
					return nil, fmt.Errorf("failed to generate and send OTP")
				}
			}

			got, err := service.SendOTPToEmail(tt.args.ctx, tt.args.msisdn, tt.args.email, tt.args.appID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.SendOTPToEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == "" {
				t.Errorf("Service.SendOTPToEmail() Expected an OTP to be returned")
			}
		})
	}
}

func TestService_GenerateRetryOtp(t *testing.T) {
	ctx := context.Background()
	service := otp.NewService(mailSvc, smsSvs, twilioSvc)
	phoneNumber := interserviceclient.TestUserPhoneNumber
	invalidPhoneNumber := "this is definitely not a number"

	appID := uuid.New().String()
	type args struct {
		ctx       context.Context
		msisdn    *string
		retryStep int
		appID     *string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Generate retry OTP whatsapp happy case",
			args: args{
				ctx:       ctx,
				msisdn:    &phoneNumber,
				retryStep: 1,
				appID:     &appID,
			},
			wantErr: false,
		},
		{
			name: "Generate retry OTP whatsapp sad case",
			args: args{
				ctx:       ctx,
				msisdn:    &phoneNumber,
				retryStep: 1,
				appID:     &appID,
			},
			wantErr: true,
		},
		{
			name: "Generate retry OTP twilio happy case",
			args: args{
				ctx:       ctx,
				msisdn:    &phoneNumber,
				retryStep: 2,
				appID:     &appID,
			},
			wantErr: false,
		},
		{
			name: "Generate retry OTP twilio sad case",
			args: args{
				ctx:       ctx,
				msisdn:    &phoneNumber,
				retryStep: 2,
				appID:     &appID,
			},
			wantErr: true,
		},
		{
			name: "Generate retry OTP sad case",
			args: args{
				ctx:    ctx,
				msisdn: &invalidPhoneNumber,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Generate retry OTP whatsapp happy case" {
				fakeOTP.GenerateOTPFn = func(ctx context.Context) (string, error) {
					return "1234", nil
				}

				fakeOTP.SaveOTPToFirestoreFn = func(otp dto.OTP) error {
					return nil
				}

				fakeTwilio.PhoneNumberVerificationCodeFn = func(
					ctx context.Context,
					to string,
					code string,
					marketingMessage string,
				) (bool, error) {
					return true, nil
				}
			}

			if tt.name == "Generate retry OTP whatsapp sad case" {
				fakeOTP.GenerateOTPFn = func(ctx context.Context) (string, error) {
					return "1234", nil
				}

				fakeOTP.SaveOTPToFirestoreFn = func(otp dto.OTP) error {
					return nil
				}

				fakeTwilio.PhoneNumberVerificationCodeFn = func(
					ctx context.Context,
					to string,
					code string,
					marketingMessage string,
				) (bool, error) {
					return false, fmt.Errorf("failed to generate OTP")
				}
			}

			if tt.name == "Generate retry OTP twilio happy case" {
				fakeOTP.GenerateOTPFn = func(ctx context.Context) (string, error) {
					return "1234", nil
				}

				fakeOTP.SaveOTPToFirestoreFn = func(otp dto.OTP) error {
					return nil
				}

				fakeTwilio.SendSMSFn = func(ctx context.Context, to string, msg string) error {
					return nil
				}
			}

			if tt.name == "Generate retry OTP twilio sad case" {
				fakeOTP.GenerateOTPFn = func(ctx context.Context) (string, error) {
					return "1234", nil
				}

				fakeOTP.SaveOTPToFirestoreFn = func(otp dto.OTP) error {
					return nil
				}

				fakeTwilio.SendSMSFn = func(ctx context.Context, to string, msg string) error {
					return fmt.Errorf("failed to generate OTP")
				}
			}

			otp, err := service.GenerateRetryOTP(tt.args.ctx, tt.args.msisdn, tt.args.retryStep, tt.args.appID)
			if err == nil {
				assert.NotNil(t, otp)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GenerateRetryOtp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestService_GenerateOTP(t *testing.T) {
	ctx := context.Background()
	service := otp.NewService(mailSvc, smsSvs, twilioSvc)

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case -> Generate OTP",
			args: args{
				ctx: ctx,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Happy Case -> Generate OTP" {
				fakeOTP.GenerateOTPFn = func(ctx context.Context) (string, error) {
					return "1234", nil
				}
			}

			if tt.name == "Sad Case -> Fail to Generate OTP" {
				fakeOTP.GenerateOTPFn = func(ctx context.Context) (string, error) {
					return "", fmt.Errorf("failed to generate OTP")
				}
			}

			otp, err := service.GenerateOTP(tt.args.ctx)
			if err == nil {
				assert.NotNil(t, otp)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GenerateOTP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestService_EmailVerificationOtp(t *testing.T) {
	email := ValidTestEmail
	invalidEmail := "not an email address"
	integrationTestEmail := otp.ITEmail

	service := otp.NewService(mailSvc, smsSvs, twilioSvc)
	type args struct {
		ctx   context.Context
		email *string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid email",
			args: args{
				ctx:   context.Background(),
				email: &email,
			},
			wantErr: false,
		},
		{
			name: "invalid email",
			args: args{
				ctx:   context.Background(),
				email: &invalidEmail,
			},
			wantErr: true,
		},
		{
			name: "valid I.T email",
			args: args{
				ctx:   context.Background(),
				email: &integrationTestEmail,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid email" {
				fakeOTP.GenerateOTPFn = func(ctx context.Context) (string, error) {
					return "1234", nil
				}

				fakeOTP.SaveOTPToFirestoreFn = func(otp dto.OTP) error {
					return nil
				}

				fakeMail.SendEmailFn = func(
					ctx context.Context,
					subject, text string,
					body *string,
					to ...string,
				) (string, string, error) {
					return "", "", nil
				}
			}

			if tt.name == "valid I.T email" {
				fakeOTP.GenerateOTPFn = func(ctx context.Context) (string, error) {
					return "1234", nil
				}

				fakeOTP.SaveOTPToFirestoreFn = func(otp dto.OTP) error {
					return nil
				}

				fakeMail.SendEmailFn = func(
					ctx context.Context,
					subject, text string,
					body *string,
					to ...string,
				) (string, string, error) {
					return "", "", nil
				}
			}

			if tt.name == "invalid email" {
				fakeOTP.GenerateOTPFn = func(ctx context.Context) (string, error) {
					return "1234", nil
				}

				fakeOTP.SaveOTPToFirestoreFn = func(otp dto.OTP) error {
					return nil
				}

				fakeMail.SendEmailFn = func(
					ctx context.Context,
					subject, text string,
					body *string,
					to ...string,
				) (string, string, error) {
					return "", "", fmt.Errorf("failed to send an email: invalid email provided")
				}
			}

			code, err := service.EmailVerificationOtp(tt.args.ctx, tt.args.email)
			if err == nil {
				assert.NotNil(t, code)
				if tt.args.email == &integrationTestEmail {
					assert.Equal(t, code, otp.ITCode)
				}
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.EmailVerificationOtp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestService_VerifyOtp(t *testing.T) {
	phoneNumber := interserviceclient.TestUserPhoneNumber
	invalidNumber := "1111"
	var srv otp.ServiceOTP = &fakeOTP
	assert.NotNil(t, srv, "service should not be bil")
	ctx := context.Background()

	otp_code := "123456"
	testPhone := otp.ITPhoneNumber
	testCode := otp.ITCode

	type args struct {
		ctx              context.Context
		msisdn           *string
		verificationCode *string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "verify otp happy case",
			args: args{
				ctx:              ctx,
				msisdn:           &phoneNumber,
				verificationCode: &otp_code,
			},
			wantErr: false,
			want:    true,
		},
		{
			name: "verify otp invalid phonenumber",
			args: args{
				ctx:              ctx,
				msisdn:           &invalidNumber,
				verificationCode: &otp_code,
			},
			wantErr: true,
			want:    false,
		},
		{
			name: "verify otp I.T case",
			args: args{
				ctx:              ctx,
				msisdn:           &testPhone,
				verificationCode: &testCode,
			},
			wantErr: false,
			want:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "verify otp happy case" {
				fakeOTP.VerifyOtpFn = func(
					ctx context.Context,
					email *string,
					verificationCode *string,
				) (bool, error) {
					return true, nil
				}
			}
			if tt.name == "verify otp invalid phonenumber" {
				fakeOTP.VerifyOtpFn = func(
					ctx context.Context,
					email *string,
					verificationCode *string,
				) (bool, error) {
					return false, fmt.Errorf("test error")
				}
			}
			if tt.name == "verify otp I.T case" {
				fakeOTP.VerifyOtpFn = func(
					ctx context.Context,
					email *string,
					verificationCode *string,
				) (bool, error) {
					return true, nil
				}
			}
			got, err := srv.VerifyOtp(tt.args.ctx, tt.args.msisdn, tt.args.verificationCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.VerifyOtp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Service.VerifyOtp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_VerifyEmailOtp(t *testing.T) {
	var s otp.ServiceOTP = &fakeOTP
	ctx := context.Background()
	email := ValidTestEmail
	testEmail := otp.ITEmail
	testCode := otp.ITCode
	randomCode := "random"
	otp := "123456"
	type args struct {
		ctx              context.Context
		email            *string
		verificationCode *string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    bool
	}{
		{
			name: "happy case",
			args: args{
				ctx:              ctx,
				email:            &email,
				verificationCode: &otp,
			},
			wantErr: false,
			want:    true,
		},
		{
			name: "happy case - integration tests",
			args: args{
				ctx:              ctx,
				email:            &testEmail,
				verificationCode: &testCode,
			},
			wantErr: false,
			want:    true,
		},
		{
			name: "sad case",
			args: args{
				ctx:              ctx,
				email:            &testEmail,
				verificationCode: &randomCode,
			},
			wantErr: true,
			want:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "happy case" {
				fakeOTP.VerifyEmailOtpFn = func(
					ctx context.Context,
					email *string,
					verificationCode *string,
				) (bool, error) {
					return true, nil
				}
			}

			if tt.name == "happy case - integration tests" {
				fakeOTP.VerifyEmailOtpFn = func(
					ctx context.Context,
					email *string,
					verificationCode *string,
				) (bool, error) {
					return true, nil
				}
			}

			if tt.name == "sad case" {
				fakeOTP.VerifyEmailOtpFn = func(
					ctx context.Context,
					email *string,
					verificationCode *string,
				) (bool, error) {
					return false, fmt.Errorf("test error")
				}
			}
			verify, err := s.VerifyEmailOtp(tt.args.ctx, tt.args.email, tt.args.verificationCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.VerifyEmailOtp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if verify != tt.want {
				t.Errorf("Service.VerifyEmailOtp() = %v, want %v", verify, tt.want)
			}
		})
	}
}

func Test_sendOtp(t *testing.T) {
	ctx := context.Background()
	service := otp.NewService(mailSvc, smsSvs, twilioSvc)
	appID := uuid.New().String()
	type args struct {
		ctx                   context.Context
		normalizedPhoneNumber string
		code                  string
		appID                 *string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Happy Case -> valid normalized Kenyan number",
			args: args{
				ctx:                   ctx,
				normalizedPhoneNumber: interserviceclient.TestUserPhoneNumber,
				code:                  "123456",
				appID:                 &appID,
			},
			want:    "123456",
			wantErr: false,
		},
		{
			name: "Sad Case -> fail to send otp to Kenyan number",
			args: args{
				ctx:                   ctx,
				normalizedPhoneNumber: interserviceclient.TestUserPhoneNumber,
				code:                  "123456",
				appID:                 &appID,
			},
			wantErr: true,
		},
		{
			name: "Happy Case -> valid normalized foreign number",
			args: args{
				ctx:                   ctx,
				normalizedPhoneNumber: "+1(202)856-9601",
				code:                  "123456",
				appID:                 &appID,
			},
			want:    "123456",
			wantErr: false,
		},
		{
			name: "Sad Case -> fail to send otp to international number",
			args: args{
				ctx:                   ctx,
				normalizedPhoneNumber: "+1(202)856-9601",
				code:                  "123456",
				appID:                 &appID,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Happy Case -> valid normalized Kenyan number" {
				fakeSMS.SendFn = func(
					ctx context.Context,
					to, message string,
				) (*silcomms.BulkSMSResponse, error) {
					return &silcomms.BulkSMSResponse{}, nil
				}
			}

			if tt.name == "Sad Case -> fail to send otp to Kenyan number" {
				fakeSMS.SendFn = func(
					ctx context.Context,
					to, message string,
				) (*silcomms.BulkSMSResponse, error) {
					return &silcomms.BulkSMSResponse{}, fmt.Errorf("failed to send OTP")
				}
			}

			if tt.name == "Happy Case -> valid normalized foreign number" {
				fakeTwilio.SendSMSFn = func(ctx context.Context, to string, msg string) error {
					return nil
				}
			}

			if tt.name == "Sad Case -> fail to send otp to international number" {
				fakeTwilio.SendSMSFn = func(ctx context.Context, to string, msg string) error {
					return fmt.Errorf("failed to send OTP")
				}
			}

			got, err := service.SendOTP(tt.args.ctx, tt.args.normalizedPhoneNumber, tt.args.code, tt.args.appID)
			if (err != nil) != tt.wantErr {
				t.Errorf("SendOTP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SendOTP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_SendTemporaryPIN(t *testing.T) {
	ctx := context.Background()
	service := otp.NewService(mailSvc, smsSvs, twilioSvc)
	phone := interserviceclient.TestUserPhoneNumber

	type args struct {
		ctx   context.Context
		input dto.TemporaryPIN
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "sad: error when sending via whatsapp",
			args: args{
				ctx: ctx,
				input: dto.TemporaryPIN{
					PhoneNumber: phone,
					FirstName:   "Test",
					PIN:         "1234",
					Channel:     1,
				},
			},
			wantErr: true,
		},
		{
			name: "sad: unable to send whatsapp",
			args: args{
				ctx: ctx,
				input: dto.TemporaryPIN{
					PhoneNumber: phone,
					FirstName:   "Test",
					PIN:         "1234",
					Channel:     1,
				},
			},
			wantErr: true,
		},
		{
			name: "happy: sent whatsapp",
			args: args{
				ctx: ctx,
				input: dto.TemporaryPIN{
					PhoneNumber: phone,
					FirstName:   "Test",
					PIN:         "1234",
					Channel:     1,
				},
			},
			wantErr: false,
		},
		{
			name: "sad: error when sending via sms",
			args: args{
				ctx: ctx,
				input: dto.TemporaryPIN{
					PhoneNumber: phone,
					FirstName:   "Test",
					PIN:         "1234",
					Channel:     2,
				},
			},
			wantErr: true,
		},
		{
			name: "happy: sent via sms",
			args: args{
				ctx: ctx,
				input: dto.TemporaryPIN{
					PhoneNumber: phone,
					FirstName:   "Test",
					PIN:         "1234",
					Channel:     2,
				},
			},
			wantErr: false,
		},
		{
			name: "sad: invalid retry step",
			args: args{
				ctx: ctx,
				input: dto.TemporaryPIN{
					PhoneNumber: phone,
					FirstName:   "Test",
					PIN:         "1234",
					Channel:     0,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "sad: error when sending via whatsapp" {
				fakeTwilio.TemporaryPINFn = func(ctx context.Context, to, message string) (bool, error) {
					return false, fmt.Errorf("unable to send whatsapp")
				}
			}
			if tt.name == "sad: unable to send whatsapp" {
				fakeTwilio.TemporaryPINFn = func(ctx context.Context, to, message string) (bool, error) {
					return false, nil
				}
			}
			if tt.name == "happy: sent whatsapp" {
				fakeTwilio.TemporaryPINFn = func(ctx context.Context, to, message string) (bool, error) {
					return true, nil
				}
			}
			if tt.name == "sad: error when sending via sms" {
				fakeTwilio.SendSMSFn = func(ctx context.Context, to, msg string) error {
					return fmt.Errorf("unable to send sms")
				}
			}
			if tt.name == "happy: sent via sms" {
				fakeTwilio.SendSMSFn = func(ctx context.Context, to, msg string) error {
					return nil
				}
			}
			if err := service.SendTemporaryPIN(tt.args.ctx, tt.args.input); (err != nil) != tt.wantErr {
				t.Errorf("Service.SendTemporaryPIN() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
