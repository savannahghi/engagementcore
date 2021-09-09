package otp_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/savannahghi/engagementcore/pkg/engagement/application/common/dto"
	otpService "github.com/savannahghi/engagementcore/pkg/engagement/infrastructure/services/otp"
	otpMock "github.com/savannahghi/engagementcore/pkg/engagement/infrastructure/services/otp/mock"
	"github.com/savannahghi/engagementcore/pkg/engagement/usecases/otp"
	"github.com/savannahghi/interserviceclient"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
)

var (
	fakeOTP otpMock.FakeServiceOTP
)

const (
	InternationalTestUserPhoneNumber = "+12028569601"
	ValidTestEmail                   = "automated.test.user.bewell-app-ci@healthcloud.co.ke"
)

func TestMain(m *testing.M) {
	os.Setenv("ROOT_COLLECTION_SUFFIX", "staging")
	os.Setenv("ENVIRONMENT", "staging")
	os.Exit(m.Run())
}

func TestService_GenerateAndSendOTP(t *testing.T) {
	ctx := context.Background()
	var service otp.UsecaseOTP = &fakeOTP
	appID := ksuid.New().String()
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
				fakeOTP.GenerateAndSendOTPFn = func(
					ctx context.Context,
					msisdn string,
					appID *string,
				) (string, error) {
					return "otp", nil
				}
			}

			if tt.name == "Sad Case -> Fail to send OTP" {
				fakeOTP.GenerateAndSendOTPFn = func(
					ctx context.Context,
					msisdn string,
					appID *string,
				) (string, error) {
					return "", fmt.Errorf("test error")
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
	var service otp.UsecaseOTP = &fakeOTP
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
				fakeOTP.SendOTPToEmailFn = func(
					ctx context.Context,
					msisdn,
					email *string,
					appID *string,
				) (string, error) {
					return "otp", nil
				}
			}

			if tt.name == "Sad Case -> Fail to Send OTP to Email" {
				fakeOTP.SendOTPToEmailFn = func(
					ctx context.Context,
					msisdn,
					email *string,
					appID *string,
				) (string, error) {
					return "", fmt.Errorf("test error")
				}
			}

			if tt.name == "Sad Case -> Fail to generate and OTP" {
				fakeOTP.SendOTPToEmailFn = func(
					ctx context.Context,
					msisdn,
					email *string,
					appID *string,
				) (string, error) {
					return "", fmt.Errorf("test error")
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

func TestService_GenerateOTP(t *testing.T) {
	ctx := context.Background()
	var service otp.UsecaseOTP = &fakeOTP

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
	integrationTestEmail := otpService.ITEmail

	var service otp.UsecaseOTP = &fakeOTP
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
				fakeOTP.EmailVerificationOtpFn = func(
					ctx context.Context,
					email *string,
				) (string, error) {
					return "123456", nil
				}
			}

			if tt.name == "valid I.T email" {
				fakeOTP.EmailVerificationOtpFn = func(
					ctx context.Context,
					email *string,
				) (string, error) {
					return "123456", nil
				}
			}

			if tt.name == "invalid email" {
				fakeOTP.EmailVerificationOtpFn = func(
					ctx context.Context,
					email *string,
				) (string, error) {
					return "", fmt.Errorf("test error")
				}
			}

			code, err := service.EmailVerificationOtp(tt.args.ctx, tt.args.email)
			if err == nil {
				assert.NotNil(t, code)
				if tt.args.email == &integrationTestEmail {
					assert.Equal(t, code, otpService.ITCode)
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
	var srv otp.UsecaseOTP = &fakeOTP
	assert.NotNil(t, srv, "service should not be bil")
	ctx := context.Background()

	otp_code := "123456"
	testPhone := otpService.ITPhoneNumber
	testCode := otpService.ITCode

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
	var s otp.UsecaseOTP = &fakeOTP
	ctx := context.Background()
	email := ValidTestEmail
	testEmail := otpService.ITEmail
	testCode := otpService.ITCode
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

func TestService_SendTemporaryPIN(t *testing.T) {
	ctx := context.Background()
	var service otp.UsecaseOTP = &fakeOTP
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
			name: "happy case",
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
			name: "sad case",
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "happy case" {
				fakeOTP.SendTemporaryPINFn = func(
					ctx context.Context,
					input dto.TemporaryPIN,
				) error {
					return nil
				}
			}

			if tt.name == "sad case" {
				fakeOTP.SendTemporaryPINFn = func(
					ctx context.Context,
					input dto.TemporaryPIN,
				) error {
					return fmt.Errorf("test error")
				}
			}
			if err := service.SendTemporaryPIN(tt.args.ctx, tt.args.input); (err != nil) != tt.wantErr {
				t.Errorf("Service.SendTemporaryPIN() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestImplOTP_GenerateRetryOTP(t *testing.T) {

	ctx := context.Background()
	msisdn := InternationalTestUserPhoneNumber
	retryStep := 1
	appID := ksuid.New().String()

	var service otp.UsecaseOTP = &fakeOTP

	otp := "123456"

	type args struct {
		ctx       context.Context
		msisdn    *string
		retryStep int
		appID     *string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "happy case",
			args: args{
				ctx:       ctx,
				msisdn:    &msisdn,
				retryStep: retryStep,
				appID:     &appID,
			},
			want:    otp,
			wantErr: false,
		},
		{
			name: "sad case",
			args: args{
				ctx:       ctx,
				msisdn:    &msisdn,
				retryStep: retryStep,
				appID:     &appID,
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "happy case" {
				fakeOTP.GenerateRetryOTPFn = func(
					ctx context.Context,
					msisdn *string,
					retryStep int,
					appID *string,
				) (string, error) {
					return otp, nil
				}
			}

			if tt.name == "sad case" {
				fakeOTP.GenerateRetryOTPFn = func(
					ctx context.Context,
					msisdn *string,
					retryStep int,
					appID *string,
				) (string, error) {
					return "", fmt.Errorf("test error")
				}
			}
			got, err := service.GenerateRetryOTP(tt.args.ctx, tt.args.msisdn, tt.args.retryStep, tt.args.appID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ImplOTP.GenerateRetryOTP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ImplOTP.GenerateRetryOTP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestImplOTP_SaveOTPToFirestore(t *testing.T) {
	var service otp.UsecaseOTP = &fakeOTP

	otp := dto.OTP{
		MSISDN:            interserviceclient.TestUserPhoneNumber,
		Message:           "test message",
		AuthorizationCode: "test code",
		Timestamp:         time.Now(),
		IsValid:           true,
		Email:             "test@bewell.co.ke",
	}

	type args struct {
		otp dto.OTP
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case",
			args: args{
				otp: otp,
			},
		},
		{
			name:    "sad case",
			args:    args{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "happy case" {
				fakeOTP.SaveOTPToFirestoreFn = func(
					otp dto.OTP,
				) error {
					return nil
				}
			}

			if tt.name == "sad case" {
				fakeOTP.SaveOTPToFirestoreFn = func(
					otp dto.OTP,
				) error {
					return fmt.Errorf("test error")
				}
			}
			if err := service.SaveOTPToFirestore(tt.args.otp); (err != nil) != tt.wantErr {
				t.Errorf("ImplOTP.SaveOTPToFirestore() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
