package onboarding_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/savannahghi/engagementcore/pkg/engagement/application/common/dto"
	"github.com/savannahghi/engagementcore/pkg/engagement/infrastructure"
	onboardingService "github.com/savannahghi/engagementcore/pkg/engagement/infrastructure/services/onboarding"
	"github.com/savannahghi/engagementcore/pkg/engagement/usecases/onboarding"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/profileutils"
)

func InitializeTestNewOnboarding(ctx context.Context) (*onboarding.ImplOnboarding, infrastructure.Interactor, error) {
	infra := infrastructure.NewInteractor()
	onboarding := onboarding.NewOnboarding(infra)
	return onboarding, infra, nil
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

func TestNewOnboarding(t *testing.T) {
	ctx := context.Background()
	f, i, err := InitializeTestNewOnboarding(ctx)
	if err != nil {
		t.Errorf("failed to initialize new FCM: %v", err)
	}

	type args struct {
		infrastructure infrastructure.Interactor
	}
	tests := []struct {
		name string
		args args
		want *onboarding.ImplOnboarding
	}{
		{
			name: "default case",
			args: args{
				infrastructure: i,
			},
			want: f,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := onboarding.NewOnboarding(tt.args.infrastructure); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFCM() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRemoteProfileService_GetEmailAddresses(t *testing.T) {
	ctx, token, err := interserviceclient.GetPhoneNumberAuthenticatedContextAndToken(t, onboardingISCClient(t))
	if err != nil {
		t.Errorf("cant get phone number authenticated context and token: %v", err)
		return
	}

	f, _, err := InitializeTestNewOnboarding(ctx)
	if err != nil {
		t.Errorf("failed to initialize new FCM: %v", err)
	}

	type args struct {
		ctx  context.Context
		uids onboardingService.UserUIDs
	}
	tests := []struct {
		name    string
		args    args
		wantNil bool
		wantErr bool
	}{
		{
			name: "happy case: get email addresses inter-service API call to profile",
			args: args{
				ctx: ctx,
				uids: onboardingService.UserUIDs{
					UIDs: []string{token.UID},
				},
			},
			wantNil: false,
			wantErr: false,
		},
		{
			name: "sad case: get email addresses inter-service API call to profile",
			args: args{
				ctx: ctx,
				uids: onboardingService.UserUIDs{
					UIDs: []string{},
				},
			},
			wantNil: true,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := f.GetEmailAddresses(tt.args.ctx, tt.args.uids)
			if (err != nil) != tt.wantErr {
				t.Errorf("RemoteProfileService.GetEmailAddresses() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantNil {
				if got == nil {
					t.Errorf("got back nil contact data")
					return
				}
			}
		})
	}
}

func TestRemoteProfileService_GetPhoneNumbers(t *testing.T) {
	ctx, token, err := interserviceclient.GetPhoneNumberAuthenticatedContextAndToken(t, onboardingISCClient(t))
	if err != nil {
		t.Errorf("cant get phone number authenticated context and token: %v", err)
		return
	}

	f, _, err := InitializeTestNewOnboarding(ctx)
	if err != nil {
		t.Errorf("failed to initialize new FCM: %v", err)
	}

	type args struct {
		ctx  context.Context
		uids onboardingService.UserUIDs
	}
	tests := []struct {
		name    string
		args    args
		wantNil bool
		wantErr bool
	}{
		{
			name: "happy case: get phone numbers inter-service API call to profile",
			args: args{
				ctx: ctx,
				uids: onboardingService.UserUIDs{
					UIDs: []string{token.UID},
				},
			},
			wantNil: false,
			wantErr: false,
		},
		{
			name: "sad case: get phone numbers inter-service API call to profile",
			args: args{
				ctx: ctx,
				uids: onboardingService.UserUIDs{
					UIDs: []string{}, // empty UID list
				},
			},
			wantNil: true,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := f.GetPhoneNumbers(tt.args.ctx, tt.args.uids)
			if (err != nil) != tt.wantErr {
				t.Errorf("RemoteProfileService.GetPhoneNumbers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantNil {
				if got == nil {
					t.Errorf("got back nil contact data")
					return
				}
			}
		})
	}
}

func TestRemoteProfileService_GetDeviceTokens(t *testing.T) {
	ctx, token, err := interserviceclient.GetPhoneNumberAuthenticatedContextAndToken(t, onboardingISCClient(t))
	if err != nil {
		t.Errorf("cant get phone number authenticated context and token: %v", err)
		return
	}

	f, _, err := InitializeTestNewOnboarding(ctx)
	if err != nil {
		t.Errorf("failed to initialize new FCM: %v", err)
	}

	type args struct {
		ctx  context.Context
		uids onboardingService.UserUIDs
	}
	tests := []struct {
		name    string
		args    args
		wantNil bool
		wantErr bool
	}{
		{
			name: "happy case: get device tokens inter-service API call to profile",
			args: args{
				ctx: ctx,
				uids: onboardingService.UserUIDs{
					UIDs: []string{token.UID},
				},
			},
			wantNil: false,
			wantErr: false,
		},
		{
			name: "sad case: get device tokens inter-service API call to profile",
			args: args{
				ctx: ctx,
				uids: onboardingService.UserUIDs{
					UIDs: []string{},
				},
			},
			wantNil: true,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := f.GetDeviceTokens(tt.args.ctx, tt.args.uids)
			if (err != nil) != tt.wantErr {
				t.Errorf("RemoteProfileService.GetDeviceTokens() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantNil {
				if got == nil {
					t.Errorf("got back nil contact data")
					return
				}
			}
		})
	}
}

func TestRemoteProfileService_GetUserProfile(t *testing.T) {
	ctx, _, err := interserviceclient.GetPhoneNumberAuthenticatedContextAndToken(t, onboardingISCClient(t))
	if err != nil {
		t.Errorf("cant get phone number authenticated context and token: %v", err)
		return
	}

	f, _, err := InitializeTestNewOnboarding(ctx)
	if err != nil {
		t.Errorf("failed to initialize new FCM: %v", err)
	}
	type args struct {
		ctx context.Context
		uid string
	}

	UID, err := firebasetools.GetLoggedInUserUID(ctx)
	if err != nil {
		t.Errorf("can't get logged in user: %v", err)
		return
	}
	invalidUID := "9VwnREOH8GdSfaxH69J6MvCu1gp9"

	tests := []struct {
		name    string
		args    args
		wantNil bool
		wantErr bool
	}{
		{
			name:    "happy case: got user profile",
			args:    args{ctx: ctx, uid: UID},
			wantNil: false,
			wantErr: false,
		},
		{
			name:    "sad case: unable to get user profile",
			args:    args{ctx: ctx, uid: invalidUID},
			wantNil: true,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := f.GetUserProfile(tt.args.ctx, tt.args.uid)
			if (err != nil) != tt.wantErr {
				t.Errorf("RemoteProfileService.GetUserProfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantNil {
				if got == nil {
					t.Errorf("got back nil profile data")
					return
				}
			}
		})
	}
}

func TestRemoteProfileService_GetUserProfileByPhoneOrEmail(t *testing.T) {
	ctx, _, err := interserviceclient.GetPhoneNumberAuthenticatedContextAndToken(t, onboardingISCClient(t))
	if err != nil {
		t.Errorf("cant get phone number authenticated context and token: %v", err)
		return
	}
	f, _, err := InitializeTestNewOnboarding(ctx)
	if err != nil {
		t.Errorf("failed to initialize new FCM: %v", err)
	}

	validPhone := interserviceclient.TestUserPhoneNumber
	invalidPhone := "+2547+"
	invalidEmail := "test"
	type args struct {
		ctx     context.Context
		payload *dto.RetrieveUserProfileInput
	}
	tests := []struct {
		name    string
		args    args
		want    *profileutils.UserProfile
		wantErr bool
	}{
		{
			name: "Happy case:phone",
			args: args{
				ctx: ctx,
				payload: &dto.RetrieveUserProfileInput{
					PhoneNumber: &validPhone,
				},
			},
			wantErr: false,
		},

		{
			name: "Sad case:phone",
			args: args{
				ctx: context.Background(),
				payload: &dto.RetrieveUserProfileInput{
					PhoneNumber: &invalidPhone,
				},
			},
			want:    nil,
			wantErr: true,
		},

		{
			name: "Sad case:email",
			args: args{
				ctx: context.Background(),
				payload: &dto.RetrieveUserProfileInput{
					EmailAddress: &invalidEmail,
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := f.GetUserProfileByPhoneOrEmail(tt.args.ctx, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("RemoteProfileService.GetUserProfileByPhoneOrEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("RemoteProfileService.GetUserProfileByPhoneOrEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}
