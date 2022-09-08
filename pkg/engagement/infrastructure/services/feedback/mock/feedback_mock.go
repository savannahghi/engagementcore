package mock

import (
	"context"

	"github.com/savannahghi/engagementcore/pkg/engagement/domain"
)

// FakeFeedback simulates the behavior of our Feedback response implementation
type FakeFeedback struct {
	RecordSurveyFeedbackResponseFn func(
		ctx context.Context,
		input domain.SurveyInput,
	) (bool, error)
}

// RecordSurveyFeedbackResponse is a mock of the RecordSurveyFeedbackResponse method
func (s *FakeFeedback) RecordSurveyFeedbackResponse(
	ctx context.Context,
	input domain.SurveyInput,
) (bool, error) {
	return s.RecordSurveyFeedbackResponseFn(ctx, input)
}
