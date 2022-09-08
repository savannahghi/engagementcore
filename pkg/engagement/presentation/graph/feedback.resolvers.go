package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"time"

	"github.com/savannahghi/engagementcore/pkg/engagement/domain"
	"github.com/savannahghi/serverutils"
)

func (r *mutationResolver) RecordSurveyFeedbackResponse(ctx context.Context, input *domain.SurveyInput) (bool, error) {
	startTime := time.Now()

	r.checkPreconditions()
	r.CheckUserTokenInContext(ctx)

	response, err := r.usecases.RecordSurveyFeedbackResponse(ctx, *input)
	if err != nil {
		return false, fmt.Errorf("failed to record patient's feedback")
	}

	defer serverutils.RecordGraphqlResolverMetrics(
		ctx,
		startTime,
		"RecordPatientFeedback",
		err,
	)
	return response, nil
}
