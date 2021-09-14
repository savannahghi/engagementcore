package sms_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/savannahghi/engagementcore/pkg/engagement/application/common/dto"
	smsMock "github.com/savannahghi/engagementcore/pkg/engagement/infrastructure/services/sms/mock"
	"github.com/savannahghi/engagementcore/pkg/engagement/usecases/sms"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/firebasetools"
	"github.com/segmentio/ksuid"
)

var (
	fakeSMSService smsMock.FakeServiceSMS
)

func TestServiceSMSImpl_SendToMany(t *testing.T) {

	var s sms.UsecaseSMS = &fakeSMSService

	ctx := context.Background()
	message := "test message"
	to := []string{firebasetools.TestUserEmail}
	from := enumutils.SenderIDBewell

	recipients := []dto.Recipient{
		{
			Number:    "2",
			Cost:      "0.7",
			Status:    "ok",
			MessageID: ksuid.New().String(),
		},
	}
	msg := dto.SMS{
		Recipients: recipients,
	}
	response := dto.SendMessageResponse{
		SMSMessageData: &msg,
	}

	type args struct {
		ctx     context.Context
		message string
		to      []string
		from    enumutils.SenderID
	}
	tests := []struct {
		name    string
		args    args
		want    *dto.SendMessageResponse
		wantErr bool
	}{
		{
			name: "happy case",
			args: args{
				ctx:     ctx,
				message: message,
				to:      to,
				from:    from,
			},
			want:    &response,
			wantErr: false,
		},
		{
			name: "sad case",
			args: args{
				ctx: ctx,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		if tt.name == "happy case" {
			fakeSMSService.SendToManyFn = func(
				ctx context.Context,
				message string,
				to []string,
				from enumutils.SenderID,
			) (*dto.SendMessageResponse, error) {
				return &response, nil
			}
		}
		if tt.name == "sad case" {
			fakeSMSService.SendToManyFn = func(
				ctx context.Context,
				message string,
				to []string,
				from enumutils.SenderID,
			) (*dto.SendMessageResponse, error) {
				return nil, fmt.Errorf("test error")
			}
		}
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.SendToMany(tt.args.ctx, tt.args.message, tt.args.to, tt.args.from)
			if (err != nil) != tt.wantErr {
				t.Errorf("ServiceSMSImpl.SendToMany() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ServiceSMSImpl.SendToMany() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServiceSMSImpl_Send(t *testing.T) {
	var s sms.UsecaseSMS = &fakeSMSService

	ctx := context.Background()
	message := "test message"
	to := firebasetools.TestUserEmail
	from := enumutils.SenderIDBewell

	recipients := []dto.Recipient{
		{
			Number:    "2",
			Cost:      "0.7",
			Status:    "ok",
			MessageID: ksuid.New().String(),
		},
	}
	msg := dto.SMS{
		Recipients: recipients,
	}
	response := dto.SendMessageResponse{
		SMSMessageData: &msg,
	}

	type args struct {
		ctx     context.Context
		to      string
		message string
		from    enumutils.SenderID
	}
	tests := []struct {
		name    string
		args    args
		want    *dto.SendMessageResponse
		wantErr bool
	}{
		{
			name: "happy case",
			args: args{
				ctx:     ctx,
				message: message,
				to:      to,
				from:    from,
			},
			want:    &response,
			wantErr: false,
		},
		{
			name: "sad case",
			args: args{
				ctx: ctx,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "happy case" {
				fakeSMSService.SendFn = func(
					ctx context.Context,
					message string,
					to string,
					from enumutils.SenderID,
				) (*dto.SendMessageResponse, error) {
					return &response, nil
				}
			}
			if tt.name == "sad case" {
				fakeSMSService.SendFn = func(
					ctx context.Context,
					message string,
					to string,
					from enumutils.SenderID,
				) (*dto.SendMessageResponse, error) {
					return nil, fmt.Errorf("test error")
				}
			}
			got, err := s.Send(tt.args.ctx, tt.args.to, tt.args.message, tt.args.from)
			if (err != nil) != tt.wantErr {
				t.Errorf("ServiceSMSImpl.Send() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ServiceSMSImpl.Send() = %v, want %v", got, tt.want)
			}
		})
	}
}
