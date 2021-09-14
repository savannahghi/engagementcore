package surveys_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/savannahghi/engagementcore/pkg/engagement/application/common/dto"
	surveysMock "github.com/savannahghi/engagementcore/pkg/engagement/infrastructure/services/surveys/mock"
	"github.com/savannahghi/engagementcore/pkg/engagement/usecases/surveys"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/interserviceclient"
	"github.com/segmentio/ksuid"
)

var (
	fakeSurveysService surveysMock.FakeSUrveys
)

func TestUnit_RecordNPSResponse(t *testing.T) {

	var s surveys.UsecaseSurveys = &fakeSurveysService

	email := firebasetools.TestUserEmail
	phoneNumber := interserviceclient.TestUserPhoneNumber

	feedback := []*dto.FeedbackInput{
		{
			Question: "test?",
			Answer:   "test",
		},
	}

	input := dto.NPSInput{
		Name:        "test",
		Score:       8,
		SladeCode:   ksuid.New().String(),
		Email:       &email,
		PhoneNumber: &phoneNumber,
		Feedback:    feedback,
	}

	type args struct {
		ctx   context.Context
		input dto.NPSInput
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
		},
		{
			name:    "sad case",
			args:    args{},
			wantErr: true,
			want:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "happy case" {
				fakeSurveysService.RecordNPSResponseFn = func(
					ctx context.Context,
					input dto.NPSInput,
				) (bool, error) {
					return true, nil
				}
			}
			if tt.name == "sad case" {
				fakeSurveysService.RecordNPSResponseFn = func(
					ctx context.Context,
					input dto.NPSInput,
				) (bool, error) {
					return false, fmt.Errorf("test error")
				}
			}

			got, err := s.RecordNPSResponse(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ServiceSurveyImpl.RecordNPSResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ServiceSurveyImpl.RecordNPSResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}
