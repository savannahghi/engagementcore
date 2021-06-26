package library_test

import (
	"context"
	"testing"

	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/application/common/helpers"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/infrastructure/services/library"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/infrastructure/services/onboarding"
)

const (
	onboardingService = "profile"
)

func TestNewLibraryService(t *testing.T) {
	onboardingClient := helpers.InitializeInterServiceClient(onboardingService)
	onboarding := onboarding.NewRemoteProfileService(onboardingClient)
	srv := library.NewLibraryService(onboarding)
	if srv == nil {
		t.Errorf("nil library service")
	}
}

func TestService_GetFeedContent(t *testing.T) {
	onboardingClient := helpers.InitializeInterServiceClient(onboardingService)
	onboarding := onboarding.NewRemoteProfileService(onboardingClient)
	s := library.NewLibraryService(onboarding)
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name        string
		args        args
		wantNonZero bool
		wantErr     bool
	}{
		{
			name: "default case",
			args: args{
				ctx: context.Background(),
			},
			wantNonZero: true,
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.GetFeedContent(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetFeedContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantNonZero && len(got) < 1 {
				t.Errorf("expected a non zero count of posts")
				return
			}
		})
	}
}

func TestService_GetFaqsContent(t *testing.T) {
	onboardingClient := helpers.InitializeInterServiceClient(onboardingService)
	onboarding := onboarding.NewRemoteProfileService(onboardingClient)
	s := library.NewLibraryService(onboarding)
	type args struct {
		ctx     context.Context
		flavour base.Flavour
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid:retrieved_user_rbac_faq",
			args: args{
				ctx:     context.Background(),
				flavour: "PRO",
			},
			wantErr: false,
		},
		{
			name: "valid:retrieved_consumer_faq",
			args: args{
				ctx:     context.Background(),
				flavour: "CONSUMER",
			},
			wantErr: false,
		},
		{
			name: "invalid:failed_to_retrieved_consumer_faq",
			args: args{
				ctx:     context.Background(),
				flavour: "CONSUMER",
			},
			wantErr: true,
		},
		{
			name: "invalid:failed_to_get_logged_in_user",
			args: args{
				ctx:     context.Background(),
				flavour: "PRO",
			},
			wantErr: true,
		},
		{
			name: "ivalid:failed_to_get_user_profile",
			args: args{
				ctx:     context.Background(),
				flavour: "PRO",
			},
			wantErr: true,
		},
		{
			name: "ivalid:failed_to_get_user_faqs",
			args: args{
				ctx:     context.Background(),
				flavour: "PRO",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := s.GetFaqsContent(tt.args.ctx, tt.args.flavour)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetFaqsContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestService_GetLibraryContent(t *testing.T) {
	onboardingClient := helpers.InitializeInterServiceClient(onboardingService)
	onboarding := onboarding.NewRemoteProfileService(onboardingClient)
	s := library.NewLibraryService(onboarding)
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name        string
		args        args
		wantNonZero bool
		wantErr     bool
	}{
		{
			name: "default case",
			args: args{
				ctx: context.Background(),
			},
			wantNonZero: true,
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.GetLibraryContent(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetLibraryContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantNonZero && len(got) < 1 {
				t.Errorf("expected a non zero count of posts")
				return
			}
		})
	}
}
