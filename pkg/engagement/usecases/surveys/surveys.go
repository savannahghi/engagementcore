package surveys

import (
	"context"

	"github.com/savannahghi/engagementcore/pkg/engagement/application/common/dto"
	"github.com/savannahghi/engagementcore/pkg/engagement/infrastructure"
)

// UsecaseSurveys defines surveys service usecases interface
type UsecaseSurveys interface {
	RecordNPSResponse(
		ctx context.Context,
		input dto.NPSInput,
	) (bool, error)
}

// ImplSurveys is the Surveys service implementation
type ImplSurveys struct {
	infrastructure infrastructure.Interactor
}

// NewSurveys initializes a Surveys service instance
func NewSurveys(infrastructure infrastructure.Interactor) *ImplSurveys {
	return &ImplSurveys{
		infrastructure: infrastructure,
	}
}

// RecordNPSResponse records the NPS response
func (f *ImplSurveys) RecordNPSResponse(
	ctx context.Context,
	input dto.NPSInput,
) (bool, error) {
	i := f.infrastructure.ServiceSurveyImpl
	return i.RecordNPSResponse(
		ctx,
		input,
	)
}
