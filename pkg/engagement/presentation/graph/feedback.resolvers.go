package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"time"

	"github.com/savannahghi/engagementcore/pkg/engagement/application/common/dto"
	"github.com/savannahghi/serverutils"
)

// RecordPatientFeedback is the resolver for the recordPatientFeedback field.
func (r *mutationResolver) RecordPatientFeedback(ctx context.Context, input dto.PatientFeedbackInput) (bool, error) {
	startTime := time.Now()

	r.checkPreconditions()
	r.CheckUserTokenInContext(ctx)

	response, err := r.infra.RecordPatientFeedback(ctx, input)
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
