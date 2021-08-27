package fcm

import (
	"context"
	"time"

	"github.com/savannahghi/engagement/pkg/engagement/application/common/dto"
	"github.com/savannahghi/engagement/pkg/engagement/infrastructure"
	"github.com/savannahghi/firebasetools"
)

type FCMUsecases interface {
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

type FCMUsecasesImpl struct {
	infrastructure infrastructure.Infrastructure
}

func NewFCMUsecaseImpl(infrastructure infrastructure.Infrastructure) *FCMUsecasesImpl {
	return &FCMUsecasesImpl{infrastructure: infrastructure}
}

func (f *FCMUsecasesImpl) SendNotification(
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
func (f *FCMUsecasesImpl) Notifications(
	ctx context.Context,
	registrationToken string,
	newerThan time.Time,
	limit int,
) ([]*dto.SavedNotification, error) {
	return f.infrastructure.Notifications(ctx, registrationToken, newerThan, limit)
}

// SendFCMByPhoneOrEmail is used to send FCM notification by phone or email
func (f *FCMUsecasesImpl) SendFCMByPhoneOrEmail(
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
