package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"time"

	"github.com/savannahghi/serverutils"
)

func (r *mutationResolver) PhoneNumberVerificationCode(ctx context.Context, to string, code string, marketingMessage string) (bool, error) {
	startTime := time.Now()

	r.checkPreconditions()
	r.CheckUserTokenInContext(ctx)

	verificationCode, err := r.usecases.PhoneNumberVerificationCode(ctx, to, code, marketingMessage)
	if err != nil {
		return false, fmt.Errorf("failed to send a verification code: %v", err)
	}

	defer serverutils.RecordGraphqlResolverMetrics(
		ctx,
		startTime,
		"phoneNumberVerificationCode",
		err,
	)

	return verificationCode, nil
}
