package feedback

import (
	"context"

	"github.com/savannahghi/engagementcore/pkg/engagement/domain"
	"github.com/savannahghi/engagementcore/pkg/engagement/infrastructure"
)

// UsecaseFeedback defines feedback service usecases interface
type UsecaseFeedback interface {
	RecordSurveyFeedbackResponse(
		ctx context.Context,
		input domain.SurveyInput,
	) (bool, error)
}

// ImplFeedback is the Feedback service implementation
type ImplFeedback struct {
	infrastructure infrastructure.Interactor
}

// NewFeedback initializes a Surveys service instance
func NewFeedback(infrastructure infrastructure.Interactor) *ImplFeedback {
	return &ImplFeedback{
		infrastructure: infrastructure,
	}
}

// RecordSurveyFeedbackResponse records the Survey response
func (f *ImplFeedback) RecordSurveyFeedbackResponse(
	ctx context.Context,
	input domain.SurveyInput,
) (bool, error) {
	i := f.infrastructure.ServiceFeedbackImpl
	return i.RecordSurveyFeedbackResponse(
		ctx,
		input,
	)
}
