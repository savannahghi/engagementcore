package sms_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/savannahghi/engagement/pkg/engagement/application/common/dto"
	db "github.com/savannahghi/engagement/pkg/engagement/infrastructure/database/firestore"
	"github.com/savannahghi/engagement/pkg/engagement/infrastructure/services/messaging"
	"github.com/savannahghi/engagement/pkg/engagement/infrastructure/services/sms"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/serverutils"
)

func TestMain(m *testing.M) {
	os.Setenv("ROOT_COLLECTION_SUFFIX", "testing")
	os.Exit(m.Run())
}

func newTestSMSService() (*sms.ServiceSMSImpl, error) {
	ctx := context.Background()
	fr, err := db.NewFirebaseRepository(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't instantiate firebase repository in resolver: %w", err)
	}

	ps, err := messaging.NewPubSubNotificationService(
		ctx,
		serverutils.MustGetEnvVar(serverutils.GoogleCloudProjectIDEnvVarName),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"can't instantiate notification service in resolver: %w",
			err,
		)
	}
	return sms.NewService(fr, ps), nil
}

func TestSendToMany(t *testing.T) {
	ctx := context.Background()
	service, err := newTestSMSService()
	if err != nil {
		t.Errorf("unable to initialize test service with error %v", err)
		return
	}

	type args struct {
		message string
		to      []string
		sender  enumutils.SenderID
	}

	tests := []struct {
		name    string
		args    args
		want    *dto.SendMessageResponse
		wantErr bool
	}{
		{
			name: "valid:successfully send to many using BeWell",
			args: args{
				message: "This is a test",
				to:      []string{"+254711223344", "+254700990099"},
				sender:  enumutils.SenderIDBewell,
			},
			wantErr: false,
		},
		{
			name: "valid:successfully send to many using Slade260",
			args: args{
				message: "This is a test",
				to:      []string{"+254711223344", "+254700990099"},
				sender:  enumutils.SenderIDSLADE360,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.SendToMany(ctx, tt.args.message, tt.args.to, enumutils.SenderIDBewell)
			if (err != nil) != tt.wantErr {
				t.Errorf("SendToMany() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if err == nil {
					t.Errorf("error expected got %v", err)
					return
				}
			}
			if !tt.wantErr {
				if err != nil {
					t.Errorf("error not expected got %v", err)
					return
				}

				if got == nil {
					t.Errorf("nil response returned")
					return
				}
			}
		})
	}
}

func TestSend(t *testing.T) {
	ctx := context.Background()
	service, err := newTestSMSService()
	if err != nil {
		t.Errorf("unable to initialize test service with error %v", err)
		return
	}

	type args struct {
		to      string
		message string
		sender  enumutils.SenderID
	}

	tests := []struct {
		name    string
		args    args
		want    *dto.SendMessageResponse
		wantErr bool
	}{
		{
			name: "valid:successfully send",
			args: args{
				message: "This is a test",
				to:      "+254711223344",
				sender:  enumutils.SenderIDSLADE360,
			},
			wantErr: false,
		},
		{
			name: "invalid:fail to send",
			args: args{
				message: "",
				to:      "+",
				sender:  enumutils.SenderIDSLADE360,
			},
			wantErr: true,
		},
		{
			name: "send from an unknown sender",
			args: args{
				message: "This is a test",
				to:      "+254711223344",
				sender:  "na-kitambi-utaezana",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.Send(ctx, tt.args.to, tt.args.message, tt.args.sender)
			if (err != nil) != tt.wantErr {
				t.Errorf("Send error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if err == nil {
					t.Errorf("error expected got %v", err)
					return
				}
			}
			if !tt.wantErr {
				if err != nil {
					t.Errorf("error not expected got %v", err)
					return
				}

				if got == nil {
					t.Errorf("nil response returned")
					return
				}
			}
		})
	}
}
