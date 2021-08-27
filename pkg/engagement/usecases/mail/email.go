package mail

import (
	"context"

	"github.com/savannahghi/engagement/pkg/engagement/infrastructure"
)

type MailUsecases interface {
	SimpleEmail(ctx context.Context, subject string, text string, body *string, to []string) (string, error)
}

type MailUsecasesImpl struct {
	infrastructure infrastructure.Infrastructure
}

func NewMailUsecasesImpl(infrastructure infrastructure.Infrastructure) *MailUsecasesImpl {
	return &MailUsecasesImpl{infrastructure: infrastructure}
}

func (m *MailUsecasesImpl) SimpleEmail(
	ctx context.Context,
	subject, text string,
	body *string,
	to ...string,
) (string, error) {
	return m.infrastructure.SimpleEmail(ctx, subject, text, body, to...)
}
