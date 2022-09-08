package feedback_test

import (
	"context"
	"os"
	"testing"

	"github.com/savannahghi/engagementcore/pkg/engagement/domain"
	"github.com/savannahghi/engagementcore/pkg/engagement/infrastructure"
	"github.com/savannahghi/engagementcore/pkg/engagement/usecases/feedback"
	"github.com/stretchr/testify/assert"
)

func InitializeTestNewFeedback(ctx context.Context) (*feedback.ImplFeedback, infrastructure.Interactor, error) {
	infra := infrastructure.NewInteractor()
	feedback := feedback.NewFeedback(infra)
	return feedback, infra, nil
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestNewService(t *testing.T) {
	ctx := context.Background()
	s, _, err := InitializeTestNewFeedback(ctx)
	if err != nil {
		t.Errorf("failed to initialize survey feedback response interractor: %v", err)
	}
	assert.NotNil(t, s)

	tests := []struct {
		name string
		want *feedback.ImplFeedback
	}{
		{
			name: "good case",
			want: s,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := s
			assert.NotNil(t, got)
		})
	}
}

func TestServiceFeedbackImpl_RecordSurveyFeedbackResponse(t *testing.T) {
	ctx := context.Background()
	s, _, err := InitializeTestNewFeedback(ctx)
	if err != nil {
		t.Errorf("failed to initialize survey feedback response: %v", err)
	}

	type args struct {
		ctx   context.Context
		input domain.SurveyInput
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
		panics  bool
	}{
		{
			name:   "invalid: missing input",
			args:   args{},
			panics: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.panics {
				got, err := s.RecordSurveyFeedbackResponse(tt.args.ctx, tt.args.input)
				if err != nil {
					t.Errorf("ServiceFeedbackImpl.RecordSurveyFeedbackResponse() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got != tt.want {
					t.Errorf("ServiceFeedbackImpl.RecordSurveyFeedbackResponse() = %v, want %v", got, tt.want)
				}
			}
			if tt.panics {
				fcRecordSurveyFeedbackResponse := func() { _, _ = s.RecordSurveyFeedbackResponse(tt.args.ctx, tt.args.input) }
				assert.Panics(t, fcRecordSurveyFeedbackResponse)
			}
		})
	}
}
