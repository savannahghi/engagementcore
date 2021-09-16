package fcm

import (
	"context"
	"time"

	"github.com/savannahghi/engagementcore/pkg/engagement/application/common/dto"
	"github.com/savannahghi/engagementcore/pkg/engagement/infrastructure"
	"github.com/savannahghi/firebasetools"
)

// UsecaseFCM defines FCM service usecases interface
type UsecaseFCM interface {
	SendNotification(
		ctx context.Context,
		registrationTokens []string,
		data map[string]string,
		notification *firebasetools.FirebaseSimpleNotificationInput,
		android *firebasetools.FirebaseAndroidConfigInput,
		ios *firebasetools.FirebaseAPNSConfigInput,
		web *firebasetools.FirebaseWebpushConfigInput,
	) (bool, error)

	Notifications(
		ctx context.Context,
		registrationToken string,
		newerThan time.Time,
		limit int,
	) ([]*dto.SavedNotification, error)

	SendFCMByPhoneOrEmail(
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

// ImplFCM is the FCM service implementation
type ImplFCM struct {
	infrastructure infrastructure.Interactor
}

// NewFCM initializes a FCM service instance
func NewFCM(infrastructure infrastructure.Interactor) *ImplFCM {
	return &ImplFCM{
		infrastructure: infrastructure,
	}
}

// SendNotification sends a notification to ios, android, web users
func (f *ImplFCM) SendNotification(
	ctx context.Context,
	registrationTokens []string,
	data map[string]string,
	notification *firebasetools.FirebaseSimpleNotificationInput,
	android *firebasetools.FirebaseAndroidConfigInput,
	ios *firebasetools.FirebaseAPNSConfigInput,
	web *firebasetools.FirebaseWebpushConfigInput,
) (bool, error) {
	i := f.infrastructure.ServiceFCMImpl
	return i.SendNotification(
		ctx,
		registrationTokens,
		data,
		notification,
		android,
		ios,
		web,
	)
}

// Notifications fetches notifications with the defined limit and set date
func (f *ImplFCM) Notifications(
	ctx context.Context,
	registrationToken string,
	newerThan time.Time,
	limit int,
) ([]*dto.SavedNotification, error) {
	i := f.infrastructure.ServiceFCMImpl
	return i.Notifications(
		ctx,
		registrationToken,
		newerThan,
		limit,
	)
}

// SendFCMByPhoneOrEmail sends fcm by phone or email
func (f *ImplFCM) SendFCMByPhoneOrEmail(
	ctx context.Context,
	phoneNumber *string,
	email *string,
	data map[string]interface{},
	notification firebasetools.FirebaseSimpleNotificationInput,
	android *firebasetools.FirebaseAndroidConfigInput,
	ios *firebasetools.FirebaseAPNSConfigInput,
	web *firebasetools.FirebaseWebpushConfigInput,
) (bool, error) {
	i := f.infrastructure.ServiceFCMImpl
	return i.SendFCMByPhoneOrEmail(
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
