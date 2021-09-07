package mail_test

import (
	"context"
	"fmt"
	"testing"

	mailMock "github.com/savannahghi/engagementcore/pkg/engagement/infrastructure/services/mail/mock"
	"github.com/savannahghi/engagementcore/pkg/engagement/usecases/mail"
	"github.com/savannahghi/firebasetools"
	"github.com/stretchr/testify/assert"
)

var (
	fakeFakeServiceMail mailMock.FakeServiceMail
)

func TestSendInBlue(t *testing.T) {
	ctx := context.Background()

	var s mail.UsecaseMail = &fakeFakeServiceMail

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

	var s mail.UsecaseMail = &fakeFakeServiceMail

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

	var s mail.UsecaseMail = &fakeFakeServiceMail
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

	var s mail.UsecaseMail = &fakeFakeServiceMail

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
