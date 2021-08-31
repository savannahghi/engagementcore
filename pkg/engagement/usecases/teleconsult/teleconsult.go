package tc

import (
	"context"

	"github.com/savannahghi/engagement/pkg/engagement/application/common/dto"
	"github.com/savannahghi/engagement/pkg/engagement/infrastructure"
)

// TeleconsultUsecases ...
type TeleconsultUsecases interface {
	TwilioAccessToken(ctx context.Context) (*dto.AccessToken, error)
}

// TeleconsultImpl ...
type TeleconsultImpl struct {
	infrastructure infrastructure.Infrastructure
}

// NewTeleconsultImpl ...
func NewTeleconsultImpl(infrastructure infrastructure.Infrastructure) *TeleconsultImpl {
	return &TeleconsultImpl{infrastructure: infrastructure}
}

// TwilioAccessToken ...
func (t *TeleconsultImpl) TwilioAccessToken(ctx context.Context) (*dto.AccessToken, error) {
	return t.infrastructure.TwilioAccessToken(ctx)
}
