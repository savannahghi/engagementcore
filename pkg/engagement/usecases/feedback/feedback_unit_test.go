package feedback_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/savannahghi/engagementcore/pkg/engagement/domain"
	feedbackMock "github.com/savannahghi/engagementcore/pkg/engagement/infrastructure/services/feedback/mock"
	"github.com/savannahghi/engagementcore/pkg/engagement/usecases/feedback"
)

var (
	fakeFeedbackService feedbackMock.FakeFeedback
)

func TestUnit_RecordSurveyFeedbackResponse(t *testing.T) {
	var f feedback.UsecaseFeedback = &fakeFeedbackService

	feedback := []*domain.SurveyFeedbackInput{
		{
			Question: "test",
			Answer:   "Test",
		},
	}

	input := domain.SurveyInput{
		Feedback:      feedback,
		ExtraFeedback: "test",
	}

	type args struct {
		ctx   context.Context
		input domain.SurveyInput
	}
	tests := []struct {
		name string

		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "happy case",
			args: args{
				input: input,
			},
			wantErr: false,
			want:    true,
		}, {
			name:    "sad case",
			args:    args{},
			wantErr: true,
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "happy case" {
				fakeFeedbackService.RecordSurveyFeedbackResponseFn = func(
					ctx context.Context,
					input domain.SurveyInput,
				) (bool, error) {
					return true, nil
				}
			}
			if tt.name == "sad case" {
				fakeFeedbackService.RecordSurveyFeedbackResponseFn = func(
					ctx context.Context,
					input domain.SurveyInput,
				) (bool, error) {
					return false, fmt.Errorf("test error")
				}
			}

			got, err := f.RecordSurveyFeedbackResponse(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ServiceFeedbackImpl.RecordSurveyFeedbackResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("ServiceFeedbackImpl.RecordSurveyFeedbackResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}
