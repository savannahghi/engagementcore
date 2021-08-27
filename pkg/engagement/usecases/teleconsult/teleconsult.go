package tc

import (
	"context"

	"github.com/savannahghi/engagement/pkg/engagement/application/common/dto"
	"github.com/savannahghi/engagement/pkg/engagement/infrastructure"
)

type TeleconsultUsecases interface {
	TwilioAccessToken(ctx context.Context) (*dto.AccessToken, error)
}

type TeleconsultImpl struct {
	infrastructure infrastructure.Infrastructure
}

func NewTeleconsultImpl(infrastructure infrastructure.Infrastructure) *TeleconsultImpl {
	return &TeleconsultImpl{infrastructure: infrastructure}
}

func (t *TeleconsultImpl) TwilioAccessToken(ctx context.Context) (*dto.AccessToken, error) {
	return t.infrastructure.TwilioAccessToken(ctx)
}
