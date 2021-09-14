package sms_test

import (
	"context"
	"os"
	"reflect"
	"testing"

	"github.com/savannahghi/engagementcore/pkg/engagement/application/common/dto"
	"github.com/savannahghi/engagementcore/pkg/engagement/infrastructure"
	"github.com/savannahghi/engagementcore/pkg/engagement/usecases/sms"
	"github.com/savannahghi/enumutils"
)

func InitializeTestNewSMS(ctx context.Context) (*sms.ImplSMS, infrastructure.Interactor, error) {
	infra := infrastructure.NewInteractor()
	sms := sms.NewSMS(infra)
	return sms, infra, nil
}

func TestMain(m *testing.M) {
	os.Setenv("ROOT_COLLECTION_SUFFIX", "testing")
	os.Exit(m.Run())
}

func TestNewSMSService(t *testing.T) {
	ctx := context.Background()
	f, i, err := InitializeTestNewSMS(ctx)
	if err != nil {
		t.Errorf("failed to initialize new sms: %v", err)
	}

	type args struct {
		infrastructure infrastructure.Interactor
	}

	tests := []struct {
		name string
		args args
		want *sms.ImplSMS
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
			if got := sms.NewSMS(tt.args.infrastructure); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSMS() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSendToMany(t *testing.T) {
	ctx := context.Background()

	service, _, err := InitializeTestNewSMS(ctx)
	if err != nil {
		t.Errorf("failed to initialize new sms: %v", err)
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
	service, _, err := InitializeTestNewSMS(ctx)
	if err != nil {
		t.Errorf("failed to initialize new sms: %v", err)
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
