package onboarding

import (
	"context"

	"github.com/savannahghi/engagementcore/pkg/engagement/application/common/dto"
	"github.com/savannahghi/engagementcore/pkg/engagement/infrastructure"
	"github.com/savannahghi/engagementcore/pkg/engagement/infrastructure/services/onboarding"
	"github.com/savannahghi/profileutils"
)

// UsecaseOnboarding defines Onboarding service usecases interface
type UsecaseOnboarding interface {
	GetEmailAddresses(
		ctx context.Context,
		uids onboarding.UserUIDs,
	) (map[string][]string, error)
	GetPhoneNumbers(
		ctx context.Context,
		uids onboarding.UserUIDs,
	) (map[string][]string, error)
	GetDeviceTokens(
		ctx context.Context,
		uid onboarding.UserUIDs,
	) (map[string][]string, error)
	GetUserProfile(ctx context.Context, uid string) (*profileutils.UserProfile, error)
	GetUserProfileByPhoneOrEmail(ctx context.Context, payload *dto.RetrieveUserProfileInput) (*profileutils.UserProfile, error)
}

// ImplOnboarding is the Onboarding service implementation
type ImplOnboarding struct {
	infrastructure infrastructure.Interactor
}

// NewOnboarding initializes a Onboarding service instance
func NewOnboarding(infrastructure infrastructure.Interactor) *ImplOnboarding {
	return &ImplOnboarding{
		infrastructure: infrastructure,
	}
}

// GetEmailAddresses gets the specified users' email addresses from the
// staging / testing / prod profile service
func (f *ImplOnboarding) GetEmailAddresses(
	ctx context.Context,
	uids onboarding.UserUIDs,
) (map[string][]string, error) {
	i := f.infrastructure.ProfileService
	return i.GetEmailAddresses(
		ctx,
		uids,
	)
}

// GetPhoneNumbers gets the specified users' phone numbers from the
// staging / testing / prod profile service
func (f *ImplOnboarding) GetPhoneNumbers(
	ctx context.Context,
	uids onboarding.UserUIDs,
) (map[string][]string, error) {
	i := f.infrastructure.ProfileService
	return i.GetPhoneNumbers(
		ctx,
		uids,
	)
}

// GetDeviceTokens gets the specified users' FCM push tokens from the
// staging / testing / prod profile service
func (f *ImplOnboarding) GetDeviceTokens(
	ctx context.Context,
	uid onboarding.UserUIDs,
) (map[string][]string, error) {
	i := f.infrastructure.ProfileService
	return i.GetDeviceTokens(
		ctx,
		uid,
	)
}

// GetUserProfile gets the specified users' profile from the onboarding service
func (f *ImplOnboarding) GetUserProfile(
	ctx context.Context,
	uid string,
) (*profileutils.UserProfile, error) {
	i := f.infrastructure.ProfileService
	return i.GetUserProfile(
		ctx,
		uid,
	)
}

// GetUserProfileByPhoneOrEmail gets the specified users' profile from the onboarding service
func (f *ImplOnboarding) GetUserProfileByPhoneOrEmail(
	ctx context.Context,
	payload *dto.RetrieveUserProfileInput,
) (*profileutils.UserProfile, error) {
	i := f.infrastructure.ProfileService
	return i.GetUserProfileByPhoneOrEmail(
		ctx,
		payload,
	)
}
