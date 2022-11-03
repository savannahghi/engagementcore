package sms_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	smsMock "github.com/savannahghi/engagementcore/pkg/engagement/infrastructure/services/sms/mock"
	"github.com/savannahghi/engagementcore/pkg/engagement/usecases/sms"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/silcomms"
)

var (
	fakeSMSService smsMock.FakeServiceSMS
)

func TestServiceSMSImpl_SendToMany(t *testing.T) {

	var s sms.UsecaseSMS = &fakeSMSService

	ctx := context.Background()
	message := "test message"
	to := []string{firebasetools.TestUserEmail}

	response := silcomms.BulkSMSResponse{
		Message:    "Hello",
		GUID:       "d818f20f-4258-4af2-bd28-4e4766440823",
		Recipients: to,
		SMS:        []string{"c3f3d2ac-8f0d-4e38-b772-d052857df6c2"},
		Updated:    "2022-11-03T13:07:10.563417+03:00",
		Created:    "2022-11-03T13:07:10.563417+03:00",
		State:      "QUEUED",
		Sender:     "BewellApp",
	}

	type args struct {
		ctx     context.Context
		message string
		to      []string
	}
	tests := []struct {
		name    string
		args    args
		want    *silcomms.BulkSMSResponse
		wantErr bool
	}{
		{
			name: "happy case",
			args: args{
				ctx:     ctx,
				message: message,
				to:      to,
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
				to []string,
				message string,
			) (*silcomms.BulkSMSResponse, error) {
				return &response, nil
			}
		}
		if tt.name == "sad case" {
			fakeSMSService.SendToManyFn = func(
				ctx context.Context,
				to []string,
				message string,
			) (*silcomms.BulkSMSResponse, error) {
				return nil, fmt.Errorf("test error")
			}
		}
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.SendToMany(tt.args.ctx, tt.args.to, tt.args.message)
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
	to := []string{firebasetools.TestUserEmail}

	response := silcomms.BulkSMSResponse{
		Message:    "Hello",
		GUID:       "d818f20f-4258-4af2-bd28-4e4766440823",
		Recipients: to,
		SMS:        []string{"c3f3d2ac-8f0d-4e38-b772-d052857df6c2"},
		Updated:    "2022-11-03T13:07:10.563417+03:00",
		Created:    "2022-11-03T13:07:10.563417+03:00",
		State:      "QUEUED",
		Sender:     "BewellApp",
	}

	type args struct {
		ctx     context.Context
		to      string
		message string
	}
	tests := []struct {
		name    string
		args    args
		want    *silcomms.BulkSMSResponse
		wantErr bool
	}{
		{
			name: "happy case",
			args: args{
				ctx:     ctx,
				message: message,
				to:      to[0],
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
					to string,
					message string,
				) (*silcomms.BulkSMSResponse, error) {
					return &response, nil
				}
			}
			if tt.name == "sad case" {
				fakeSMSService.SendFn = func(
					ctx context.Context,
					message string,
					to string,
				) (*silcomms.BulkSMSResponse, error) {
					return nil, fmt.Errorf("test error")
				}
			}
			got, err := s.Send(tt.args.ctx, tt.args.to, tt.args.message)
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
