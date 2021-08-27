package sms

import (
	"context"

	"github.com/savannahghi/engagement/pkg/engagement/application/common/dto"
	"github.com/savannahghi/engagement/pkg/engagement/infrastructure"
	"github.com/savannahghi/enumutils"
)

type SMSUsecases interface {
	Send(ctx context.Context, to string, message string) (*dto.SendMessageResponse, error)
	SendToMany(ctx context.Context, message string, to []string) (*dto.SendMessageResponse, error)
}

type SMSUsecasesImpl struct {
	infrastructure infrastructure.Infrastructure
}

func NewSMSUsecasesImpl(infrastructure infrastructure.Infrastructure) *SMSUsecasesImpl {
	return &SMSUsecasesImpl{infrastructure: infrastructure}
}

func (s *SMSUsecasesImpl) Send(ctx context.Context, to string, message string) (*dto.SendMessageResponse, error) {
	return s.infrastructure.Send(ctx, to, message, enumutils.SenderIDBewell)
}
func (s *SMSUsecasesImpl) SendToMany(ctx context.Context, message string, to []string) (*dto.SendMessageResponse, error) {
	return s.infrastructure.SendToMany(ctx, message, to, enumutils.SenderIDBewell)
}
