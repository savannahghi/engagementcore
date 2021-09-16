package twilio_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/savannahghi/engagementcore/pkg/engagement/application/common/dto"
	twilioService "github.com/savannahghi/engagementcore/pkg/engagement/infrastructure/services/twilio"
	twilioMock "github.com/savannahghi/engagementcore/pkg/engagement/infrastructure/services/twilio/mock"
	"github.com/savannahghi/firebasetools"
)

var (
	fakeTwilioService twilioMock.FakeServiceTwilio
)

func TestUnit_Room(t *testing.T) {
	var s twilioService.ServiceTwilio = &fakeTwilioService

	ctx := context.Background()

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    *dto.Room
		wantErr bool
	}{
		{
			name: "valid: valid context",
			args: args{
				ctx: ctx,
			},
			want:    &dto.Room{},
			wantErr: false,
		},
		{
			name: "invalid: invalid context",
			args: args{
				ctx: ctx,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid: valid context" {
				fakeTwilioService.RoomFn = func(
					ctx context.Context,
				) (*dto.Room, error) {
					return &dto.Room{}, nil
				}
			}
			if tt.name == "invalid: invalid context" {
				fakeTwilioService.RoomFn = func(
					ctx context.Context,
				) (*dto.Room, error) {
					return nil, fmt.Errorf("test error")
				}
			}
			got, err := s.Room(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("ServiceTwilioImpl.Room() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ServiceTwilioImpl.Room() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnit_TwilioAccessToken(t *testing.T) {
	var s twilioService.ServiceTwilio = &fakeTwilioService

	ctx := context.Background()

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    *dto.AccessToken
		wantErr bool
	}{
		{
			name: "valid: valid context",
			args: args{
				ctx: ctx,
			},
			want:    &dto.AccessToken{},
			wantErr: false,
		},
		{
			name: "invalid: invalid context",
			args: args{
				ctx: ctx,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid: valid context" {
				fakeTwilioService.TwilioAccessTokenFn = func(
					ctx context.Context,
				) (*dto.AccessToken, error) {
					return &dto.AccessToken{}, nil
				}
			}
			if tt.name == "invalid: invalid context" {
				fakeTwilioService.TwilioAccessTokenFn = func(
					ctx context.Context,
				) (*dto.AccessToken, error) {
					return nil, fmt.Errorf("test error")
				}
			}
			got, err := s.TwilioAccessToken(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("ServiceTwilioImpl.TwilioAccessToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ServiceTwilioImpl.TwilioAccessToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnit_SendSMS(t *testing.T) {
	var s twilioService.ServiceTwilio = &fakeTwilioService
	ctx := context.Background()
	to := firebasetools.TestUserEmail
	msg := "test message"

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
			name: "valid: correct params passed",
			args: args{
				ctx: ctx,
				to:  to,
				msg: msg,
			},
			wantErr: false,
		},
		{
			name: "invalid: missing params",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid: correct params passed" {
				fakeTwilioService.SendSMSFn = func(
					ctx context.Context,
					to string,
					msg string,
				) error {
					return nil
				}
			}
			if tt.name == "invalid: missing params" {
				fakeTwilioService.SendSMSFn = func(
					ctx context.Context,
					to string,
					msg string,
				) error {
					return fmt.Errorf("test error")
				}
			}
			if err := s.SendSMS(tt.args.ctx, tt.args.to, tt.args.msg); (err != nil) != tt.wantErr {
				t.Errorf("ServiceTwilioImpl.SendSMS() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUnit_SaveTwilioVideoCallbackStatus(t *testing.T) {
	var s twilioService.ServiceTwilio = &fakeTwilioService

	ctx := context.Background()

	type args struct {
		ctx  context.Context
		data dto.CallbackData
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid: correct params passed",
			args: args{
				ctx:  ctx,
				data: dto.CallbackData{},
			},
			wantErr: false,
		},
		{
			name: "invalid: missing params",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		if tt.name == "valid: correct params passed" {
			fakeTwilioService.SaveTwilioVideoCallbackStatusFn = func(
				ctx context.Context,
				data dto.CallbackData,
			) error {
				return nil
			}
		}
		if tt.name == "invalid: missing params" {
			fakeTwilioService.SaveTwilioVideoCallbackStatusFn = func(
				ctx context.Context,
				data dto.CallbackData,
			) error {
				return fmt.Errorf("test error")
			}
		}
		t.Run(tt.name, func(t *testing.T) {
			if err := s.SaveTwilioVideoCallbackStatus(tt.args.ctx, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("ServiceTwilioImpl.SaveTwilioVideoCallbackStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUnit_PhoneNumberVerificationCode(t *testing.T) {
	var s twilioService.ServiceTwilio = &fakeTwilioService

	ctx := context.Background()
	to := firebasetools.TestUserEmail
	code := "200"
	marketingMessage := "marketingMessage"

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
			name: "valid: correct params passed",
			args: args{
				ctx:              ctx,
				to:               to,
				code:             code,
				marketingMessage: marketingMessage,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "invalid: missing params",
			args: args{
				ctx: ctx,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid: correct params passed" {
				fakeTwilioService.PhoneNumberVerificationCodeFn = func(
					ctx context.Context,
					to string,
					code string,
					marketingMessage string,
				) (bool, error) {
					return true, nil
				}
			}
			if tt.name == "invalid: missing params" {
				fakeTwilioService.PhoneNumberVerificationCodeFn = func(
					ctx context.Context,
					to string,
					code string,
					marketingMessage string,
				) (bool, error) {
					return false, fmt.Errorf("test error")
				}
			}
			got, err := s.PhoneNumberVerificationCode(tt.args.ctx, tt.args.to, tt.args.code, tt.args.marketingMessage)
			if (err != nil) != tt.wantErr {
				t.Errorf("ServiceTwilioImpl.PhoneNumberVerificationCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ServiceTwilioImpl.PhoneNumberVerificationCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnit_SaveTwilioCallbackResponse(t *testing.T) {
	var s twilioService.ServiceTwilio = &fakeTwilioService

	ctx := context.Background()
	type args struct {
		ctx  context.Context
		data dto.Message
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid: correct params passed",
			args: args{
				ctx:  ctx,
				data: dto.Message{},
			},
			wantErr: false,
		},
		{
			name: "invalid: missing params",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid: correct params passed" {
				fakeTwilioService.SaveTwilioCallbackResponseFn = func(
					ctx context.Context,
					data dto.Message,
				) error {
					return nil
				}
			}
			if tt.name == "invalid: missing params" {
				fakeTwilioService.SaveTwilioCallbackResponseFn = func(
					ctx context.Context,
					data dto.Message,
				) error {
					return fmt.Errorf("test error")
				}
			}
			if err := s.SaveTwilioCallbackResponse(tt.args.ctx, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("ServiceTwilioImpl.SaveTwilioCallbackResponse() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUnit_TemporaryPIN(t *testing.T) {
	var s twilioService.ServiceTwilio = &fakeTwilioService

	to := firebasetools.TestUserEmail
	message := "test message"

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
			name: "valid: correct params passed",
			args: args{
				ctx:     ctx,
				to:      to,
				message: message,
			},
			wantErr: false,
			want:    true,
		},
		{
			name: "invalid: missing params",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
			want:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid: correct params passed" {
				fakeTwilioService.TemporaryPINFn = func(
					ctx context.Context,
					to string,
					message string,
				) (bool, error) {
					return true, nil
				}
			}
			if tt.name == "invalid: missing params" {
				fakeTwilioService.TemporaryPINFn = func(
					ctx context.Context,
					to string,
					message string,
				) (bool, error) {
					return false, fmt.Errorf("test error")
				}
			}
			got, err := s.TemporaryPIN(tt.args.ctx, tt.args.to, tt.args.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("ServiceTwilioImpl.TemporaryPIN() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ServiceTwilioImpl.TemporaryPIN() = %v, want %v", got, tt.want)
			}
		})
	}
}
