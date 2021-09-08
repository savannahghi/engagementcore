package mail_test

import (
	"context"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/savannahghi/engagementcore/pkg/engagement/application/common/dto"
	"github.com/savannahghi/engagementcore/pkg/engagement/infrastructure"
	"github.com/savannahghi/engagementcore/pkg/engagement/usecases/mail"
	"github.com/savannahghi/firebasetools"
)

func InitializeTestNewMail(ctx context.Context) (*mail.ImplMail, infrastructure.Interactor, error) {
	infra := infrastructure.NewInteractor()
	mail := mail.NewMail(infra)
	return mail, infra, nil
}

func TestMain(m *testing.M) {
	os.Setenv("ROOT_COLLECTION_SUFFIX", "testing")
	os.Exit(m.Run())
}

func TestNewService(t *testing.T) {
	ctx := context.Background()
	f, i, err := InitializeTestNewMail(ctx)
	if err != nil {
		t.Errorf("failed to initialize test mail interractor: %v", err)
	}
	type args struct {
		infrastructure infrastructure.Interactor
	}

	tests := []struct {
		name string
		args args
		want *mail.ImplMail
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
			if got := mail.NewMail(tt.args.infrastructure); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewLibrary() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_SendInBlue(t *testing.T) {
	ctx := context.Background()

	f, _, err := InitializeTestNewMail(ctx)
	if err != nil {
		t.Errorf("failed to initialize test mail interractor: %v", err)
	}

	type args struct {
		subject string
		text    string
		to      []string
	}
	tests := []struct {
		name       string
		args       args
		wantStatus string
		wantErr    bool
	}{
		{
			name: "happy case",
			args: args{
				subject: "Test Email",
				text:    "This is a test email",
				to:      []string{firebasetools.TestUserEmail},
			},
			wantStatus: "ok",
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _, err := f.SendInBlue(ctx, tt.args.subject, tt.args.text, tt.args.to...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.SendInBlue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.wantStatus {
				t.Errorf("Service.SendInBlue() got = %v, want %v", got, tt.wantStatus)
			}
		})
	}
}

func TestService_SendMailgun(t *testing.T) {
	testUserMail := "test@bewell.co.ke"
	ctx := context.Background()

	f, _, err := InitializeTestNewMail(ctx)
	if err != nil {
		t.Errorf("failed to initialize test mail interractor: %v", err)
	}

	type args struct {
		subject string
		text    string
		to      []string
	}
	tests := []struct {
		name string

		args    args
		wantErr bool
	}{
		{
			name: "valid email",
			args: args{
				subject: "Test Email",
				text:    "Test Email",
				to:      []string{testUserMail},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg, id, err := f.SendMailgun(ctx, tt.args.subject, tt.args.text, nil, tt.args.to...)
			if tt.wantErr {
				if err == nil {
					t.Errorf("an error was expected")
					return
				}
				if msg != "" && id != "" {
					t.Errorf("expected no message and message ID")
					return
				}
			}
			if !tt.wantErr {
				if err != nil {
					t.Errorf("an error was not expected")
					return
				}
				if msg == "" && id == "" {
					t.Errorf("expected a message and message ID")
					return
				}
			}
		})
	}
}

func TestService_SendEmail(t *testing.T) {
	testUserMail := "test@bewell.co.ke"
	ctx := context.Background()

	f, _, err := InitializeTestNewMail(ctx)
	if err != nil {
		t.Errorf("failed to initialize test mail interractor: %v", err)
	}

	tests := []struct {
		name    string
		subject string
		text    string
		to      []string

		expectMsg bool
		expectID  bool
		expectErr bool
	}{
		{
			name:    "valid email",
			subject: "Test Email",
			text:    "Test Email",
			to:      []string{testUserMail},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg, id, err := f.SendEmail(ctx, tt.subject, tt.text, nil, tt.to...)
			if tt.expectErr {
				if err == nil {
					t.Errorf("an error was expected")
					return
				}
				if msg != "" && id != "" {
					t.Errorf("expected no message and message ID")
					return
				}
			}
			if !tt.expectErr {
				if err != nil {
					t.Errorf("an error was not expected")
					return
				}
				if msg == "" && id == "" {
					t.Errorf("expected a message and message ID")
					return
				}
			}
		})
	}
}

func TestService_SimpleEmail(t *testing.T) {
	testUserMail := "test@bewell.co.ke"
	testBody := "This is a test email"
	ctx := context.Background()

	f, _, err := InitializeTestNewMail(ctx)
	if err != nil {
		t.Errorf("failed to initialize test mail interractor: %v", err)
	}

	type args struct {
		subject string
		text    string
		body    *string
		to      []string
	}
	tests := []struct {
		name string

		args    args
		wantErr bool
	}{
		{
			name: "valid email",
			args: args{
				subject: "Test Email",
				text:    "Test Email",
				body:    &testBody,
				to:      []string{testUserMail},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg, err := f.SimpleEmail(ctx, tt.args.subject, tt.args.text, nil, tt.args.to...)
			if tt.wantErr {
				if err == nil {
					t.Errorf("an error was expected")
					return
				}
				if msg != "" {
					t.Errorf("expected no message")
					return
				}
			}
			if !tt.wantErr {
				if err != nil {
					t.Errorf("an error was not expected")
					return
				}
				if msg == "" {
					t.Errorf("expected a message")
					return
				}
			}
		})
	}
}

func TestService_UpdateMailgunDeliveryStatus(t *testing.T) {
	ctx := context.Background()

	f, _, err := InitializeTestNewMail(ctx)
	if err != nil {
		t.Errorf("failed to initialize test mail interractor: %v", err)
	}

	payload := dto.MailgunEvent{
		EventName:   "delivered",
		DeliveredOn: time.Now().String(),
		MessageID:   "20210715172955.1.63EC29EF167F09B9@sandboxb30d61fba25641a9983c3b3a3c84abde.mailgun.org",
	}
	invalidPayload := dto.MailgunEvent{
		EventName:   "delivered",
		DeliveredOn: time.Now().String(),
		MessageID:   "invalidmessageid",
	}
	type args struct {
		ctx     context.Context
		payload *dto.MailgunEvent
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid: correct payload passed",
			args: args{
				ctx:     ctx,
				payload: &payload,
			},
			wantErr: false,
		},
		{
			name: "invalid: invalid message id",
			args: args{
				ctx:     ctx,
				payload: &invalidPayload,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := f.UpdateMailgunDeliveryStatus(ctx, tt.args.payload)
			if tt.wantErr {
				if err == nil {
					t.Errorf("an error was expected")
					return
				}
				if got != nil {
					t.Errorf("expected no log")
					return
				}
			}
			if !tt.wantErr {
				if err != nil {
					t.Errorf("an error was not expected")
					return
				}
				if got == nil {
					t.Errorf("expected an outgoing email log")
					return
				}
			}
		})
	}
}

func TestService_GenerateEmailTemplate(t *testing.T) {
	ctx := context.Background()

	f, _, err := InitializeTestNewMail(ctx)
	if err != nil {
		t.Errorf("failed to initialize test mail interractor: %v", err)
	}

	name := "test name"
	templateName := "testTemplateName"

	type args struct {
		name         string
		templateName string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "valid: correct parameters passed",
			args: args{
				name:         name,
				templateName: templateName,
			},
			want: "some template",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := f.GenerateEmailTemplate(tt.args.name, tt.args.templateName)
			if tt.want != "" {
				if got == "" {
					t.Errorf("expected a template to be generated, got %q", got)
				}
			}
		})
	}
}
