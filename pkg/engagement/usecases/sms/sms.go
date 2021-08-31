package sms

import (
	"context"

	"github.com/savannahghi/engagement/pkg/engagement/application/common/dto"
	"github.com/savannahghi/engagement/pkg/engagement/infrastructure"
	"github.com/savannahghi/enumutils"
)

// UseCases ...
type UseCases interface {
	Send(ctx context.Context, to string, message string) (*dto.SendMessageResponse, error)
	SendToMany(ctx context.Context, message string, to []string) (*dto.SendMessageResponse, error)
}

// UseCasesImpl ...
type UseCasesImpl struct {
	infrastructure infrastructure.Infrastructure
}

// NewSMSUsecasesImpl ...
func NewSMSUsecasesImpl(infrastructure infrastructure.Infrastructure) *UseCasesImpl {
	return &UseCasesImpl{infrastructure: infrastructure}
}

// Send ...
func (s *UseCasesImpl) Send(ctx context.Context, to string, message string) (*dto.SendMessageResponse, error) {
	return s.infrastructure.Send(ctx, to, message, enumutils.SenderIDBewell)
}

// SendToMany ...
func (s *UseCasesImpl) SendToMany(ctx context.Context, message string, to []string) (*dto.SendMessageResponse, error) {
	return s.infrastructure.SendToMany(ctx, message, to, enumutils.SenderIDBewell)
}
