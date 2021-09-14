package surveys_test

import (
	"context"
	"os"
	"testing"

	"github.com/savannahghi/engagementcore/pkg/engagement/application/common/dto"
	"github.com/savannahghi/engagementcore/pkg/engagement/infrastructure"
	"github.com/savannahghi/engagementcore/pkg/engagement/usecases/surveys"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/interserviceclient"
	"github.com/stretchr/testify/assert"
)

func InitializeTestNewSurveys(ctx context.Context) (*surveys.ImplSurveys, infrastructure.Interactor, error) {
	infra := infrastructure.NewInteractor()
	surveys := surveys.NewSurveys(infra)
	return surveys, infra, nil
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func onboardingISCClient(t *testing.T) *interserviceclient.InterServiceClient {
	deps, err := interserviceclient.LoadDepsFromYAML()
	if err != nil {
		t.Errorf("can't load inter-service config from YAML: %v", err)
		return nil
	}

	profileClient, err := interserviceclient.SetupISCclient(*deps, "profile")
	if err != nil {
		t.Errorf("can't set up profile interservice client: %v", err)
		return nil
	}

	return profileClient
}

func TestNewService(t *testing.T) {
	ctx := context.Background()
	s, _, err := InitializeTestNewSurveys(ctx)
	if err != nil {
		t.Errorf("failed to initialize test mail interractor: %v", err)
	}
	assert.NotNil(t, s)

	tests := []struct {
		name string
		want *surveys.ImplSurveys
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

func TestServiceSurveyImpl_RecordNPSResponse(t *testing.T) {
	onboardingClient := onboardingISCClient(t)
	ctx, token, err := interserviceclient.GetPhoneNumberAuthenticatedContextAndToken(t, onboardingClient)
	if err != nil {
		t.Errorf("cant get phone number authenticated context token: %v", err)
		return
	}
	_, err = firebasetools.GetAuthenticatedContextFromUID(ctx, token.UID)
	if err != nil {
		t.Errorf("cant get authenticated context from UID: %v", err)
		return
	}

	s, _, err := InitializeTestNewSurveys(ctx)
	if err != nil {
		t.Errorf("failed to initialize test mail interractor: %v", err)
	}

	type args struct {
		ctx   context.Context
		input dto.NPSInput
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
				got, err := s.RecordNPSResponse(tt.args.ctx, tt.args.input)
				if (err != nil) != tt.wantErr {
					t.Errorf("ServiceSurveyImpl.RecordNPSResponse() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got != tt.want {
					t.Errorf("ServiceSurveyImpl.RecordNPSResponse() = %v, want %v", got, tt.want)
				}
			}
			if tt.panics {
				fcRecordNPSResponse := func() { _, _ = s.RecordNPSResponse(tt.args.ctx, tt.args.input) }
				assert.Panics(t, fcRecordNPSResponse)
			}
		})
	}
}
