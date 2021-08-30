package mail

import (
	"context"

	"github.com/savannahghi/engagement/pkg/engagement/infrastructure"
)

// UseCases ...
type UseCases interface {
	SimpleEmail(ctx context.Context, subject string, text string, body *string, to []string) (string, error)
}

// UseCasesImpl ...
type UseCasesImpl struct {
	infrastructure infrastructure.Infrastructure
}

// NewMailUsecasesImpl ...
func NewMailUsecasesImpl(infrastructure infrastructure.Infrastructure) *UseCasesImpl {
	return &UseCasesImpl{infrastructure: infrastructure}
}

// SimpleEmail ...
func (m *UseCasesImpl) SimpleEmail(
	ctx context.Context,
	subject, text string,
	body *string,
	to ...string,
) (string, error) {
	return m.infrastructure.SimpleEmail(ctx, subject, text, body, to...)
}
