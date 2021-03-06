package mail_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/savannahghi/engagementcore/pkg/engagement/application/common/dto"
	"github.com/savannahghi/engagementcore/pkg/engagement/infrastructure/services/mail"
	mailMock "github.com/savannahghi/engagementcore/pkg/engagement/infrastructure/services/mail/mock"
	"github.com/savannahghi/firebasetools"
	"github.com/stretchr/testify/assert"

	db "github.com/savannahghi/engagementcore/pkg/engagement/infrastructure/database/firestore"
)

var (
	fakeFakeServiceMail mailMock.FakeServiceMail
)

func TestSendInBlue(t *testing.T) {
	ctx := context.Background()

	var s mail.ServiceMail = &fakeFakeServiceMail

	type args struct {
		ctx     context.Context
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
				ctx:     ctx,
				subject: "Test Email",
				text:    "This is a test email",
				to:      []string{firebasetools.TestUserEmail},
			},
			wantStatus: "ok",
			wantErr:    false,
		},
		{
			name: "sad case: missing recipient",
			args: args{
				ctx:     ctx,
				subject: "Test Email",
				text:    "Test text",
				to:      []string{},
			},
			wantErr:    true,
			wantStatus: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "happy case" {
				fakeFakeServiceMail.SendInBlueFn = func(
					ctx context.Context,
					subject string,
					text string,
					to ...string,
				) (string, string, error) {
					return "ok", "id", nil
				}
			}
			if tt.name == "sad case: missing recipient" {
				fakeFakeServiceMail.SendInBlueFn = func(
					ctx context.Context,
					subject string,
					text string,
					to ...string,
				) (string, string, error) {
					return "", "", fmt.Errorf("test error")
				}
			}
		})
		got, _, err := s.SendInBlue(tt.args.ctx,
			tt.args.subject,
			tt.args.text,
			tt.args.to...,
		)
		if tt.wantStatus == "ok" {
			assert.NotEmpty(t, got)
		}
		if !tt.wantErr && err != nil {
			t.Errorf("error not expected, got %v", err)
		}
	}
}

func TestSendMailgun(t *testing.T) {
	ctx := context.Background()

	var s mail.ServiceMail = &fakeFakeServiceMail

	testBody := "This is a test email"

	type args struct {
		ctx           context.Context
		subject, text string
		body          *string
		to            []string
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
				ctx:     ctx,
				subject: "Test Email",
				text:    "Test text",
				body:    &testBody,
				to:      []string{firebasetools.TestUserEmail},
			},
			wantStatus: "ok",
			wantErr:    false,
		},
		{
			name: "sad case: missing recipient",
			args: args{
				ctx:     ctx,
				subject: "Test Email",
				text:    "Test text",
				body:    &testBody,
				to:      []string{},
			},
			wantStatus: "",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "happy case" {
				fakeFakeServiceMail.SendMailgunFn = func(
					ctx context.Context,
					subject string,
					text string,
					body *string,
					to ...string,
				) (string, string, error) {
					return "ok", "", nil
				}
			}
			if tt.name == "sad case: missing recipient" {
				fakeFakeServiceMail.SendMailgunFn = func(
					ctx context.Context,
					subject string,
					text string,
					body *string,
					to ...string,
				) (string, string, error) {
					return "", "", fmt.Errorf("test error")
				}
			}
			got, _, err := s.SendMailgun(tt.args.ctx,
				tt.args.subject,
				tt.args.text,
				tt.args.body,
				tt.args.to...)
			if tt.wantStatus == "ok" {
				assert.NotEmpty(t, got)
			}
			if !tt.wantErr && err != nil {
				t.Errorf("error not expected, got %v", err)
			}
		})
	}
}

func TestSendEmail(t *testing.T) {
	ctx := context.Background()

	var s mail.ServiceMail = &fakeFakeServiceMail
	testBody := "This is a test email"

	type args struct {
		ctx           context.Context
		subject, text string
		body          *string
		to            []string
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
				ctx:     ctx,
				subject: "Test Email",
				text:    "Test text",
				body:    &testBody,
				to:      []string{firebasetools.TestUserEmail},
			},
			wantStatus: "ok",
			wantErr:    false,
		},
		{
			name: "sad case: missing recipient",
			args: args{
				ctx:     ctx,
				subject: "Test Email",
				text:    "Test text",
				body:    &testBody,
				to:      []string{},
			},
			wantStatus: "",
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "happy case" {
				fakeFakeServiceMail.SendEmailFn = func(
					ctx context.Context,
					subject string,
					text string,
					body *string,
					to ...string,
				) (string, string, error) {
					return "ok", "id", nil
				}
			}
			if tt.name == "sad case: missing recipient" {
				fakeFakeServiceMail.SendEmailFn = func(
					ctx context.Context,
					subject string,
					text string,
					body *string,
					to ...string,
				) (string, string, error) {
					return "", "", fmt.Errorf("test error")
				}
			}
			got, _, err := s.SendEmail(tt.args.ctx,
				tt.args.subject,
				tt.args.text,
				tt.args.body,
				tt.args.to...,
			)
			if !tt.wantErr && err != nil {
				t.Errorf("error not expected, got %v", err)
			}
			if tt.wantStatus == "ok" {
				assert.NotEmpty(t, got)
			}
		})
	}
}

func TestSimpleEmail(t *testing.T) {
	ctx := context.Background()

	var s mail.ServiceMail = &fakeFakeServiceMail

	testBody := "This is a test email"

	type args struct {
		ctx           context.Context
		subject, text string
		body          *string
		to            []string
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
				ctx:     ctx,
				subject: "Test Email",
				text:    "Test text",
				body:    &testBody,
				to:      []string{firebasetools.TestUserEmail},
			},
			wantStatus: "ok",
			wantErr:    false,
		},
		{
			name: "sad case: missing recipient",
			args: args{
				ctx:     ctx,
				subject: "Test Email",
				text:    "Test text",
				body:    &testBody,
				to:      []string{},
			},
			wantStatus: "",
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "happy case" {
				fakeFakeServiceMail.SimpleEmailFn = func(
					ctx context.Context,
					subject string,
					text string,
					body *string,
					to ...string,
				) (string, error) {
					return "ok", nil
				}
			}
			if tt.name == "sad case: missing recipient" {
				fakeFakeServiceMail.SimpleEmailFn = func(
					ctx context.Context,
					subject string,
					text string,
					body *string,
					to ...string,
				) (string, error) {
					return "", fmt.Errorf("test error")
				}
			}
		})
		got, err := s.SimpleEmail(tt.args.ctx,
			tt.args.subject,
			tt.args.text,
			tt.args.body,
			tt.args.to...,
		)

		if !tt.wantErr && err != nil {
			t.Errorf("error not expected, got %v", err)
		}
		if tt.wantStatus == "ok" {
			assert.NotEmpty(t, got)
		}

	}
}

func TestUpdateMailgunDeliveryStatus(t *testing.T) {
	ctx := context.Background()
	repo, err := db.NewFirebaseRepository(ctx)
	if err != nil {
		t.Errorf("error initializing new firebase repo:%v", err)
		return
	}
	s := mail.NewService(repo)

	// save new email logs with status pending
	emailPayload := dto.OutgoingEmailsLog{
		UUID:    uuid.New().String(),
		To:      []string{firebasetools.TestUserEmail},
		From:    "test@email.com",
		Subject: "test subject",
		Text:    "test text",
		// MessageID is a unique identifier of mailgun's message
		MessageID:   "1234",
		EmailSentOn: time.Now(),
		Event:       &dto.MailgunEventOutput{EventName: "rejected", DeliveredOn: time.Now()},
	}
	err = s.SaveOutgoingEmails(ctx, &emailPayload)
	if err != nil {
		t.Errorf("failed to to save outgoing emails log: %v", err)
	}

	updateStatusPayload := dto.MailgunEvent{
		MessageID:   "1234",
		EventName:   "delivered",
		DeliveredOn: time.Now().Local().String(),
	}

	invalidMessageID := dto.MailgunEvent{
		MessageID:   "invalid",
		EventName:   "delivered",
		DeliveredOn: time.Now().Local().String(),
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
			name: "happy case",
			args: args{
				ctx:     ctx,
				payload: &updateStatusPayload,
			},
			wantErr: true,
		},
		{
			name: "sad case: invalid messageID",
			args: args{
				ctx:     ctx,
				payload: &invalidMessageID,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := s.UpdateMailgunDeliveryStatus(
				tt.args.ctx,
				tt.args.payload)
			if !tt.wantErr && err != nil {
				t.Errorf("error not expected, got %v", err)
			}
			if tt.name == "sad case: missing recipient" {
				fakeFakeServiceMail.UpdateMailgunDeliveryStatusFn = func(
					ctx context.Context,
					payload *dto.MailgunEvent,
				) (*dto.OutgoingEmailsLog, error) {
					return nil, fmt.Errorf("test error")
				}
				got, err := fakeFakeServiceMail.UpdateMailgunDeliveryStatusFn(
					tt.args.ctx,
					tt.args.payload,
				)
				assert.NotNil(t, err)
				assert.NotNil(t, got)
			}
		})

	}
}

func TestSaveOutgoingEmails(t *testing.T) {
	ctx := context.Background()

	repo, err := db.NewFirebaseRepository(ctx)
	if err != nil {
		t.Errorf("error initializing new firebase repo:%v", err)
		return
	}

	s := mail.NewService(repo)

	emptyPayload := dto.OutgoingEmailsLog{}
	emailPayload := dto.OutgoingEmailsLog{
		UUID:    uuid.New().String(),
		To:      []string{firebasetools.TestUserEmail},
		From:    "test@email.com",
		Subject: "test subject",
		Text:    "test text",
		// MessageID is a unique identifier of mailgun's message
		MessageID:   "1234",
		EmailSentOn: time.Now(),
		Event:       &dto.MailgunEventOutput{EventName: "rejected", DeliveredOn: time.Now()},
	}

	type args struct {
		ctx     context.Context
		payload *dto.OutgoingEmailsLog
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
				ctx:     ctx,
				payload: &emailPayload,
			},
			wantErr: true,
		},
		{
			name: "sad case: empty payload",
			args: args{
				ctx:     ctx,
				payload: &emptyPayload,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "happy case" {
				fakeFakeServiceMail.SaveOutgoingEmailsFn = func(
					ctx context.Context,
					payload *dto.OutgoingEmailsLog,
				) error {
					return nil
				}
			}
			if tt.name == "sad case: missing recipient" {
				fakeFakeServiceMail.SaveOutgoingEmailsFn = func(
					ctx context.Context,
					payload *dto.OutgoingEmailsLog,
				) error {
					return fmt.Errorf("test error")
				}
			}
			err := s.SaveOutgoingEmails(
				tt.args.ctx,
				tt.args.payload,
			)
			if !tt.wantErr && err != nil {
				t.Errorf("error not expected, got %v", err)
			}
		})
	}
}

func TestGenerateEmailTemplate(t *testing.T) {
	ctx := context.Background()
	repo, err := db.NewFirebaseRepository(ctx)
	if err != nil {
		t.Errorf("error initializing new firebase repo:%v", err)
		return
	}
	s := mail.NewService(repo)

	type args struct {
		name         string
		templateName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case",
			args: args{
				name:         "test name",
				templateName: "<div>test template<div>",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "happy case" {
				fakeFakeServiceMail.GenerateEmailTemplateFn = func(
					name string, templateName string,
				) string {
					return "some template"
				}
			}
			got := s.GenerateEmailTemplate(
				tt.args.name,
				tt.args.templateName,
			)

			if !tt.wantErr {
				assert.NotNil(t, got)
			}
		})

	}
}
