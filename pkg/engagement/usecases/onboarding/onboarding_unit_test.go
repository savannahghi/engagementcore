package onboarding_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/savannahghi/engagementcore/pkg/engagement/application/common/dto"
	onboardingService "github.com/savannahghi/engagementcore/pkg/engagement/infrastructure/services/onboarding"
	onboardingMock "github.com/savannahghi/engagementcore/pkg/engagement/infrastructure/services/onboarding/mock"
	"github.com/savannahghi/engagementcore/pkg/engagement/usecases/onboarding"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/profileutils"
	"github.com/segmentio/ksuid"
)

var (
	fakeOnboardingService onboardingMock.FakeServiceOnboarding
)

func getTestProfile() profileutils.UserProfile {
	version := "1.0.0"
	testText := "testText"
	testTime := time.Now()
	testPhonenumber := interserviceclient.TestUserPhoneNumber
	testEmail := firebasetools.TestUserEmail
	return profileutils.UserProfile{
		ID:       ksuid.New().String(),
		UserName: &testText,
		VerifiedIdentifiers: []profileutils.VerifiedIdentifier{
			{
				UID:           ksuid.New().String(),
				Timestamp:     time.Now(),
				LoginProvider: profileutils.LoginProviderTypePhone,
			},
		},
		VerifiedUIDS:            []string{ksuid.New().String()},
		PrimaryPhone:            &testPhonenumber,
		PrimaryEmailAddress:     &testEmail,
		SecondaryPhoneNumbers:   []string{testPhonenumber},
		SecondaryEmailAddresses: []string{testEmail},
		PushTokens:              []string{ksuid.New().String()},
		Role:                    profileutils.RoleTypeAgent,
		Permissions:             []profileutils.PermissionType{"test"},
		FavNavActions:           []string{"test"},
		TermsAccepted:           true,
		Suspended:               false,
		PhotoUploadID:           ksuid.New().String(),
		UserBioData: profileutils.BioData{
			FirstName: &testText,
			LastName:  &testText,
			Gender:    enumutils.GenderMale,
		},
		HomeAddress: &profileutils.Address{
			Longitude:        "test",
			Latitude:         "test",
			Locality:         &testText,
			Name:             &testText,
			PlaceID:          &testText,
			FormattedAddress: &testText,
		},
		WorkAddress: &profileutils.Address{
			Longitude:        "test",
			Latitude:         "test",
			Locality:         &testText,
			Name:             &testText,
			PlaceID:          &testText,
			FormattedAddress: &testText,
		},
		CreatedByID:        &testText,
		Created:            &testTime,
		ConsumerAppVersion: &version,
		PROAppVersion:      &version,
	}
}

func TestUnit_GetEmailAddresses(t *testing.T) {
	var s onboarding.UsecaseOnboarding = &fakeOnboardingService

	ctx := context.Background()
	uids := onboardingService.UserUIDs{
		UIDs: []string{ksuid.New().String()},
	}
	type args struct {
		ctx  context.Context
		uids onboardingService.UserUIDs
	}
	tests := []struct {
		name    string
		args    args
		want    map[string][]string
		wantErr bool
	}{
		{
			name: "valid: correct params passed",
			args: args{
				ctx:  ctx,
				uids: uids,
			},
			want:    map[string][]string{"uids": {"testid"}},
			wantErr: false,
		},
		{
			name: "invalid: missing uids",
			args: args{
				ctx:  ctx,
				uids: uids,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid: correct params passed" {
				fakeOnboardingService.GetEmailAddressesFn = func(
					ctx context.Context,
					uids onboardingService.UserUIDs,
				) (map[string][]string, error) {
					return map[string][]string{"uids": {"testid"}}, nil
				}
			}
			if tt.name == "invalid: missing uids" {
				fakeOnboardingService.GetEmailAddressesFn = func(
					ctx context.Context,
					uids onboardingService.UserUIDs,
				) (map[string][]string, error) {
					return nil, fmt.Errorf("test error")
				}
			}
			got, err := s.GetEmailAddresses(tt.args.ctx, tt.args.uids)
			if (err != nil) != tt.wantErr {
				t.Errorf("RemoteProfileService.GetEmailAddresses() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RemoteProfileService.GetEmailAddresses() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnit_GetPhoneNumbers(t *testing.T) {
	var s onboarding.UsecaseOnboarding = &fakeOnboardingService

	ctx := context.Background()
	uids := onboardingService.UserUIDs{
		UIDs: []string{ksuid.New().String()},
	}

	type args struct {
		ctx  context.Context
		uids onboardingService.UserUIDs
	}
	tests := []struct {
		name    string
		args    args
		want    map[string][]string
		wantErr bool
	}{
		{
			name: "valid: correct params passed",
			args: args{
				ctx:  ctx,
				uids: uids,
			},
			want:    map[string][]string{"uids": {"testid"}},
			wantErr: false,
		},
		{
			name: "invalid: missing uids",
			args: args{
				ctx:  ctx,
				uids: uids,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid: correct params passed" {
				fakeOnboardingService.GetPhoneNumbersFn = func(
					ctx context.Context,
					uids onboardingService.UserUIDs,
				) (map[string][]string, error) {
					return map[string][]string{"uids": {"testid"}}, nil
				}
			}
			if tt.name == "invalid: missing uids" {
				fakeOnboardingService.GetPhoneNumbersFn = func(
					ctx context.Context,
					uids onboardingService.UserUIDs,
				) (map[string][]string, error) {
					return nil, fmt.Errorf("test error")
				}
			}
			got, err := s.GetPhoneNumbers(tt.args.ctx, tt.args.uids)
			if (err != nil) != tt.wantErr {
				t.Errorf("RemoteProfileService.GetPhoneNumbers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RemoteProfileService.GetPhoneNumbers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnit_GetDeviceTokens(t *testing.T) {
	var s onboarding.UsecaseOnboarding = &fakeOnboardingService

	ctx := context.Background()
	uids := onboardingService.UserUIDs{
		UIDs: []string{ksuid.New().String()},
	}

	type args struct {
		ctx  context.Context
		uids onboardingService.UserUIDs
	}
	tests := []struct {
		name    string
		args    args
		want    map[string][]string
		wantErr bool
	}{
		{
			name: "valid: correct params passed",
			args: args{
				ctx:  ctx,
				uids: uids,
			},
			want:    map[string][]string{"uids": {"testid"}},
			wantErr: false,
		},
		{
			name: "invalid: missing uids",
			args: args{
				ctx:  ctx,
				uids: uids,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid: correct params passed" {
				fakeOnboardingService.GetDeviceTokensFn = func(
					ctx context.Context,
					uids onboardingService.UserUIDs,
				) (map[string][]string, error) {
					return map[string][]string{"uids": {"testid"}}, nil
				}
			}
			if tt.name == "invalid: missing uids" {
				fakeOnboardingService.GetDeviceTokensFn = func(
					ctx context.Context,
					uids onboardingService.UserUIDs,
				) (map[string][]string, error) {
					return nil, fmt.Errorf("test error")
				}
			}
			got, err := s.GetDeviceTokens(tt.args.ctx, tt.args.uids)
			if (err != nil) != tt.wantErr {
				t.Errorf("RemoteProfileService.GetDeviceTokens() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RemoteProfileService.GetDeviceTokens() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnit_GetUserProfile(t *testing.T) {
	var s onboarding.UsecaseOnboarding = &fakeOnboardingService

	ctx := context.Background()
	uid := ksuid.New().String()

	profile := getTestProfile()

	type args struct {
		ctx context.Context
		uid string
	}
	tests := []struct {
		name    string
		args    args
		want    *profileutils.UserProfile
		wantErr bool
	}{
		{
			name: "valid: correct params passed",
			args: args{
				ctx: ctx,
				uid: uid,
			},
			want:    &profile,
			wantErr: false,
		},
		{
			name: "invalid: missing uid",
			args: args{
				ctx: ctx,
				uid: uid,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid: correct params passed" {
				fakeOnboardingService.GetUserProfileFn = func(
					ctx context.Context,
					uid string,
				) (*profileutils.UserProfile, error) {
					return &profile, nil
				}
			}
			if tt.name == "invalid: missing uid" {
				fakeOnboardingService.GetUserProfileFn = func(
					ctx context.Context,
					uid string,
				) (*profileutils.UserProfile, error) {
					return nil, fmt.Errorf("test error")
				}
			}
			got, err := s.GetUserProfile(tt.args.ctx, tt.args.uid)
			if (err != nil) != tt.wantErr {
				t.Errorf("RemoteProfileService.GetUserProfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RemoteProfileService.GetUserProfile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnit_GetUserProfileByPhoneOrEmail(t *testing.T) {

	var s onboarding.UsecaseOnboarding = &fakeOnboardingService

	testPhonenumber := interserviceclient.TestUserPhoneNumber
	testEmail := firebasetools.TestUserEmail

	ctx := context.Background()
	payload := dto.RetrieveUserProfileInput{
		PhoneNumber:  &testPhonenumber,
		EmailAddress: &testEmail,
	}

	profile := getTestProfile()

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
			name: "valid: correct params passed",
			args: args{
				ctx:     ctx,
				payload: &payload,
			},
			want:    &profile,
			wantErr: false,
		},
		{
			name: "invalid: missing payload",
			args: args{
				ctx: ctx,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid: correct params passed" {
				fakeOnboardingService.GetUserProfileByPhoneOrEmailFn = func(
					ctx context.Context,
					payload *dto.RetrieveUserProfileInput,
				) (*profileutils.UserProfile, error) {
					return &profile, nil
				}
			}
			if tt.name == "invalid: missing payload" {
				fakeOnboardingService.GetUserProfileByPhoneOrEmailFn = func(
					ctx context.Context,
					payload *dto.RetrieveUserProfileInput,
				) (*profileutils.UserProfile, error) {
					return nil, fmt.Errorf("test error")
				}
			}
			got, err := s.GetUserProfileByPhoneOrEmail(tt.args.ctx, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("RemoteProfileService.GetUserProfileByPhoneOrEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RemoteProfileService.GetUserProfileByPhoneOrEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}
