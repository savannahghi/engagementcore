package fcm

import (
	"context"
	"time"

	"github.com/savannahghi/engagement/pkg/engagement/application/common/dto"
	"github.com/savannahghi/engagement/pkg/engagement/infrastructure"
	"github.com/savannahghi/firebasetools"
)

// Usecases ...
type Usecases interface {
	SendNotification(
		ctx context.Context,
		registrationTokens []string,
		data map[string]string,
		notification *firebasetools.FirebaseSimpleNotificationInput,
		android *firebasetools.FirebaseAndroidConfigInput,
		ios *firebasetools.FirebaseAPNSConfigInput,
		web *firebasetools.FirebaseWebpushConfigInput,
	) (bool, error)
	SendFCMByPhoneOrEmail(ctx context.Context, phoneNumber *string, email *string, data map[string]interface{}, notification firebasetools.FirebaseSimpleNotificationInput, android *firebasetools.FirebaseAndroidConfigInput, ios *firebasetools.FirebaseAPNSConfigInput, web *firebasetools.FirebaseWebpushConfigInput) (bool, error)
	Notifications(ctx context.Context, registrationToken string, newerThan time.Time, limit int) ([]*dto.SavedNotification, error)
}

// UsecasesImpl ...
type UsecasesImpl struct {
	infrastructure infrastructure.Infrastructure
}

// NewFCMUsecaseImpl ...
func NewFCMUsecaseImpl(infrastructure infrastructure.Infrastructure) *UsecasesImpl {
	return &UsecasesImpl{infrastructure: infrastructure}
}

// SendNotification ...
func (f *UsecasesImpl) SendNotification(
	ctx context.Context,
	registrationTokens []string,
	data map[string]string,
	notification *firebasetools.FirebaseSimpleNotificationInput,
	android *firebasetools.FirebaseAndroidConfigInput,
	ios *firebasetools.FirebaseAPNSConfigInput,
	web *firebasetools.FirebaseWebpushConfigInput,
) (bool, error) {
	return f.infrastructure.SendNotification(ctx, registrationTokens, data, notification, android, ios, web)
}

// Notifications is used to query a user's priorities
func (f *UsecasesImpl) Notifications(
	ctx context.Context,
	registrationToken string,
	newerThan time.Time,
	limit int,
) ([]*dto.SavedNotification, error) {
	return f.infrastructure.Notifications(ctx, registrationToken, newerThan, limit)
}

// SendFCMByPhoneOrEmail is used to send FCM notification by phone or email
func (f *UsecasesImpl) SendFCMByPhoneOrEmail(
	ctx context.Context,
	phoneNumber *string,
	email *string,
	data map[string]interface{},
	notification firebasetools.FirebaseSimpleNotificationInput,
	android *firebasetools.FirebaseAndroidConfigInput,
	ios *firebasetools.FirebaseAPNSConfigInput,
	web *firebasetools.FirebaseWebpushConfigInput,
) (bool, error) {
	return f.infrastructure.SendFCMByPhoneOrEmail(ctx, phoneNumber, email, data, notification, android, ios, web)
}
