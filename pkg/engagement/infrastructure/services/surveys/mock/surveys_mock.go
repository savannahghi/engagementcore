package mock

import (
	"context"

	"github.com/savannahghi/engagementcore/pkg/engagement/application/common/dto"
)

// FakeSUrveys simulates the behavior of our Surveys response implementation
type FakeSUrveys struct {
	RecordNPSResponseFn func(
		ctx context.Context,
		input dto.NPSInput,
	) (bool, error)
}

// RecordNPSResponse is a mock of the RecordNPSResponse method
func (s *FakeSUrveys) RecordNPSResponse(
	ctx context.Context,
	input dto.NPSInput,
) (bool, error) {
	return s.RecordNPSResponseFn(ctx, input)
}
