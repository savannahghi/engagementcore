package surveys

import (
	"context"

	"github.com/savannahghi/engagement/pkg/engagement/application/common/dto"
	"github.com/savannahghi/engagement/pkg/engagement/infrastructure"
)

// UseCases ...
type UseCases interface {
	RecordNPSResponse(ctx context.Context, input dto.NPSInput) (bool, error)
}

// UseCasesImpl ...
type UseCasesImpl struct {
	infrastructure infrastructure.Infrastructure
}

// NewSurveysImpl ...
func NewSurveysImpl(infrastructure infrastructure.Infrastructure) *UseCasesImpl {
	return &UseCasesImpl{infrastructure: infrastructure}
}

// RecordNPSResponse ...
func (s UseCasesImpl) RecordNPSResponse(ctx context.Context, input dto.NPSInput) (bool, error) {
	return s.infrastructure.RecordNPSResponse(ctx, input)
}
