package mock

import (
	"context"
	"time"

	"github.com/savannahghi/engagementcore/pkg/engagement/application/common/dto"
	"github.com/savannahghi/firebasetools"
)

// FakeServiceFcm simulates the behavior of our FCM push implementation
type FakeServiceFcm struct {
	SendNotificationFn func(
		ctx context.Context,
		registrationTokens []string,
		data map[string]string,
		notification *firebasetools.FirebaseSimpleNotificationInput,
		android *firebasetools.FirebaseAndroidConfigInput,
		ios *firebasetools.FirebaseAPNSConfigInput,
		web *firebasetools.FirebaseWebpushConfigInput,
	) (bool, error)

	NotificationsFn func(
		ctx context.Context,
		registrationToken string,
		newerThan time.Time,
		limit int,
	) ([]*dto.SavedNotification, error)

	SendFCMByPhoneOrEmailFn func(
		ctx context.Context,
		phoneNumber *string,
		email *string,
		data map[string]interface{},
		notification firebasetools.FirebaseSimpleNotificationInput,
		android *firebasetools.FirebaseAndroidConfigInput,
		ios *firebasetools.FirebaseAPNSConfigInput,
		web *firebasetools.FirebaseWebpushConfigInput,
	) (bool, error)
}

// SendNotification is a mock of the SendNotification method
func (f *FakeServiceFcm) SendNotification(
	ctx context.Context,
	registrationTokens []string,
	data map[string]string,
	notification *firebasetools.FirebaseSimpleNotificationInput,
	android *firebasetools.FirebaseAndroidConfigInput,
	ios *firebasetools.FirebaseAPNSConfigInput,
	web *firebasetools.FirebaseWebpushConfigInput,
) (bool, error) {
	return f.SendNotificationFn(
		ctx,
		registrationTokens,
		data,
		notification,
		android,
		ios,
		web,
	)
}

// Notifications is a mock of the Notifications method
func (f *FakeServiceFcm) Notifications(
	ctx context.Context,
	registrationToken string,
	newerThan time.Time,
	limit int,
) ([]*dto.SavedNotification, error) {
	return f.NotificationsFn(
		ctx,
		registrationToken,
		newerThan,
		limit,
	)
}

// SendFCMByPhoneOrEmail is a mock of the SendFCMByPhoneOrEmail method
func (f *FakeServiceFcm) SendFCMByPhoneOrEmail(
	ctx context.Context,
	phoneNumber *string,
	email *string,
	data map[string]interface{},
	notification firebasetools.FirebaseSimpleNotificationInput,
	android *firebasetools.FirebaseAndroidConfigInput,
	ios *firebasetools.FirebaseAPNSConfigInput,
	web *firebasetools.FirebaseWebpushConfigInput,
) (bool, error) {
	return f.SendFCMByPhoneOrEmailFn(
		ctx,
		phoneNumber,
		email,
		data,
		notification,
		android,
		ios,
		web,
	)
}
