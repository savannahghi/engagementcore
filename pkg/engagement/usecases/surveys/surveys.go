package surveys

import (
	"context"

	"github.com/savannahghi/engagement/pkg/engagement/application/common/dto"
	"github.com/savannahghi/engagement/pkg/engagement/infrastructure"
)

type SurveysUsecases interface {
	RecordNPSResponse(ctx context.Context, input dto.NPSInput) (bool, error)
}

type SurveysImpl struct {
	infrastructure infrastructure.Infrastructure
}

func NewSurveysImpl(infrastructure infrastructure.Infrastructure) *SurveysImpl {
	return &SurveysImpl{infrastructure: infrastructure}
}

func (s SurveysImpl) RecordNPSResponse(ctx context.Context, input dto.NPSInput) (bool, error) {
	return s.infrastructure.RecordNPSResponse(ctx, input)
}
