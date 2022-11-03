package sms

import (
	"context"

	"github.com/savannahghi/engagementcore/pkg/engagement/infrastructure"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/silcomms"
)

// UsecaseSMS defines SMS service usecases interface
type UsecaseSMS interface {
	SendToMany(
		ctx context.Context,
		to []string,
		message string,
	) (*silcomms.BulkSMSResponse, error)
	Send(
		ctx context.Context,
		to, message string,
	) (*silcomms.BulkSMSResponse, error)
}

// ImplSMS is the SMS service implementation
type ImplSMS struct {
	infrastructure infrastructure.Interactor
}

// NewSMS initializes an SMS service instance
func NewSMS(infrastructure infrastructure.Interactor) *ImplSMS {
	return &ImplSMS{
		infrastructure: infrastructure,
	}
}

// SendToMany sends sms to many recipients
func (s *ImplSMS) SendToMany(
	ctx context.Context,
	message string,
	to []string,
	from enumutils.SenderID,
) (*silcomms.BulkSMSResponse, error) {
	i := s.infrastructure.ServiceSMSImpl
	return i.SendToMany(
		ctx,
		to,
		message,
	)
}

// Send is a method used to send to a single recipient
func (s *ImplSMS) Send(
	ctx context.Context,
	to, message string,
	from enumutils.SenderID,
) (*silcomms.BulkSMSResponse, error) {
	i := s.infrastructure.ServiceSMSImpl
	return i.Send(
		ctx,
		to,
		message,
	)
}
