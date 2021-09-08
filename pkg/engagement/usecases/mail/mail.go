package mail

import (
	"context"

	"github.com/savannahghi/engagementcore/pkg/engagement/application/common/dto"
	"github.com/savannahghi/engagementcore/pkg/engagement/infrastructure"
)

// UsecaseMail defines mail service usecases interface
type UsecaseMail interface {
	SendInBlue(
		ctx context.Context,
		subject,
		text string,
		to ...string,
	) (string, string, error)

	SendMailgun(
		ctx context.Context,
		subject,
		text string,
		body *string,
		to ...string,
	) (string, string, error)

	SendEmail(
		ctx context.Context,
		subject,
		text string,
		body *string,
		to ...string,
	) (string, string, error)

	SimpleEmail(
		ctx context.Context,
		subject,
		text string,
		body *string,
		to ...string,
	) (string, error)

	UpdateMailgunDeliveryStatus(
		ctx context.Context,
		payload *dto.MailgunEvent,
	) (*dto.OutgoingEmailsLog, error)

	GenerateEmailTemplate(
		name string,
		templateName string,
	) string
}

// ImplMail is the mail service implementation
type ImplMail struct {
	infrastructure infrastructure.Interactor
}

// NewMail initializes a mail service instance
func NewMail(infrastructure infrastructure.Interactor) *ImplMail {
	return &ImplMail{
		infrastructure: infrastructure,
	}
}

// SendInBlue sends email using sendinblue service
func (m *ImplMail) SendInBlue(
	ctx context.Context,
	subject,
	text string,
	to ...string,
) (string, string, error) {
	i := m.infrastructure.ServiceMailImpl
	return i.SendInBlue(
		ctx,
		subject,
		text,
		to...,
	)
}

// SendMailgun sends email using mailgun service
func (m *ImplMail) SendMailgun(
	ctx context.Context,
	subject,
	text string,
	body *string,
	to ...string,
) (string, string, error) {
	i := m.infrastructure.ServiceMailImpl
	return i.SendMailgun(
		ctx,
		subject,
		text,
		body,
		to...,
	)
}

// SendEmail sends email using mailgun service
func (m *ImplMail) SendEmail(
	ctx context.Context,
	subject,
	text string,
	body *string,
	to ...string,
) (string, string, error) {
	i := m.infrastructure.ServiceMailImpl
	return i.SendEmail(
		ctx,
		subject,
		text,
		body,
		to...,
	)
}

// SimpleEmail sends email using mailgun service
func (m *ImplMail) SimpleEmail(
	ctx context.Context,
	subject,
	text string,
	body *string,
	to ...string,
) (string, error) {
	i := m.infrastructure.ServiceMailImpl
	return i.SimpleEmail(
		ctx,
		subject,
		text,
		body,
		to...,
	)
}

// UpdateMailgunDeliveryStatus updates mailgun delivery status
func (m *ImplMail) UpdateMailgunDeliveryStatus(
	ctx context.Context,
	payload *dto.MailgunEvent,
) (*dto.OutgoingEmailsLog, error) {
	i := m.infrastructure.ServiceMailImpl
	return i.UpdateMailgunDeliveryStatus(
		ctx,
		payload,
	)
}

// GenerateEmailTemplate generates templates for email
func (m *ImplMail) GenerateEmailTemplate(
	name string,
	templateName string,
) string {
	i := m.infrastructure.ServiceMailImpl
	return i.GenerateEmailTemplate(
		name,
		templateName,
	)
}
