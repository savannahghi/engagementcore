package whatsapp

import (
	"context"

	"github.com/savannahghi/engagement/pkg/engagement/infrastructure"
)

// UseCases ...
type UseCases interface {
	PhoneNumberVerificationCode(
		ctx context.Context,
		to string,
		code string,
		marketingMessage string,
	) (bool, error)
}

// UseCasesImpl ...
type UseCasesImpl struct {
	infrastructure infrastructure.Infrastructure
}

// NewWhatsappUsecasesImpl ...
func NewWhatsappUsecasesImpl(infrastructure infrastructure.Infrastructure) *UseCasesImpl {
	return &UseCasesImpl{infrastructure: infrastructure}
}

// PhoneNumberVerificationCode ...
func (w *UseCasesImpl) PhoneNumberVerificationCode(
	ctx context.Context,
	to string,
	code string,
	marketingMessage string,
) (bool, error) {
	return w.infrastructure.PhoneNumberVerificationCode(ctx, to, code, marketingMessage)
}
