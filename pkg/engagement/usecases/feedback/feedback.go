package feedback

import (
	"context"

	"github.com/savannahghi/engagementcore/pkg/engagement/application/common/dto"
	"github.com/savannahghi/engagementcore/pkg/engagement/infrastructure"
)

type UsecaseFeedback interface {
	RecordPatientFeedback(
		ctx context.Context,
		input dto.FeedbackInput,
	) (bool, error)
}

type ImplFeedback struct {
	infrastructure infrastructure.Interactor
}

func NewPatientFeedback(infrastructure infrastructure.Interactor) *ImplFeedback {
	return &ImplFeedback{
		infrastructure: infrastructure,
	}
}

func (f *ImplFeedback) RecordPatientFeedback(
	ctx context.Context,
	input dto.PatientFeedbackInput,
) (bool, error) {
	i := f.infrastructure.ServiceFeedbackImpl
	return i.RecordPatientFeedback(
		ctx,
		input,
	)
}
