package fcm_test

import (
	"context"
	"github.com/pkg/errors"
	"reflect"
	"testing"
	"time"

	"github.com/savannahghi/engagementcore/pkg/engagement/application/common/dto"
	fcmMock "github.com/savannahghi/engagementcore/pkg/engagement/infrastructure/services/fcm/mock"
	"github.com/savannahghi/engagementcore/pkg/engagement/usecases/fcm"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/interserviceclient"
)

var (
	fakeFCMService fcmMock.FakeServiceFcm
)

func TestUnit_SendNotification(t *testing.T) {
	ctx := context.Background()
	data := map[string]string{
		"name": "user",
	}
	notification := firebasetools.FirebaseSimpleNotificationInput{}
	android := firebasetools.FirebaseAndroidConfigInput{}
	ios := firebasetools.FirebaseAPNSConfigInput{}
	web := firebasetools.FirebaseWebpushConfigInput{}

	var f fcm.UsecaseFCM = &fakeFCMService

	type args struct {
		ctx                context.Context
		registrationTokens []string
		data               map[string]string
		notification       *firebasetools.FirebaseSimpleNotificationInput
		android            *firebasetools.FirebaseAndroidConfigInput
		ios                *firebasetools.FirebaseAPNSConfigInput
		web                *firebasetools.FirebaseWebpushConfigInput
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "valid: correct params passed",
			args: args{
				ctx:                ctx,
				registrationTokens: []string{"tokens"},
				data:               data,
				notification:       &notification,
				android:            &android,
				ios:                &ios,
				web:                &web,
			},
			wantErr: false,
			want:    true,
		},
		{
			name: "invalid: missing args",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
			want:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid: correct params passed" {
				fakeFCMService.SendNotificationFn = func(
					ctx context.Context,
					registrationTokens []string,
					data map[string]string,
					notification *firebasetools.FirebaseSimpleNotificationInput,
					android *firebasetools.FirebaseAndroidConfigInput,
					ios *firebasetools.FirebaseAPNSConfigInput,
					web *firebasetools.FirebaseWebpushConfigInput,
				) (bool, error) {
					return true, nil
				}
			}
			if tt.name == "invalid: missing args" {
				fakeFCMService.SendNotificationFn = func(
					ctx context.Context,
					registrationTokens []string,
					data map[string]string,
					notification *firebasetools.FirebaseSimpleNotificationInput,
					android *firebasetools.FirebaseAndroidConfigInput,
					ios *firebasetools.FirebaseAPNSConfigInput,
					web *firebasetools.FirebaseWebpushConfigInput,
				) (bool, error) {
					return false, errors.New("tests error")
				}
			}
			got, err := f.SendNotification(tt.args.ctx, tt.args.registrationTokens, tt.args.data, tt.args.notification, tt.args.android, tt.args.ios, tt.args.web)
			if (err != nil) != tt.wantErr {
				t.Errorf("ImplFCM.SendNotification() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ImplFCM.SendNotification() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnit_Notifications(t *testing.T) {
	ctx := context.Background()
	registrationToken := "token"
	newerThan := time.Now()
	limit := 10

	notification := []*dto.SavedNotification{}

	var f fcm.UsecaseFCM = &fakeFCMService

	type args struct {
		ctx               context.Context
		registrationToken string
		newerThan         time.Time
		limit             int
	}
	tests := []struct {
		name    string
		args    args
		want    []*dto.SavedNotification
		wantErr bool
	}{
		{
			name: "valid: correct args passed",
			args: args{
				ctx:               ctx,
				registrationToken: registrationToken,
				newerThan:         newerThan,
				limit:             limit,
			},
			want:    notification,
			wantErr: false,
		},
		{
			name: "invalid: missing registrationToken, newerThan, limit args",
			args: args{
				ctx:               ctx,
				registrationToken: registrationToken,
				newerThan:         newerThan,
				limit:             limit,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid: correct args passed" {
				fakeFCMService.NotificationsFn = func(
					ctx context.Context,
					registrationToken string,
					newerThan time.Time,
					limit int,
				) ([]*dto.SavedNotification, error) {
					return notification, nil
				}
			}
			if tt.name == "invalid: missing registrationToken, newerThan, limit args" {
				fakeFCMService.NotificationsFn = func(
					ctx context.Context,
					registrationToken string,
					newerThan time.Time,
					limit int,
				) ([]*dto.SavedNotification, error) {
					return nil, errors.New("test error")
				}
			}
			got, err := f.Notifications(tt.args.ctx, tt.args.registrationToken, tt.args.newerThan, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("ImplFCM.Notifications() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ImplFCM.Notifications() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnit_SendFCMByPhoneOrEmail(t *testing.T) {
	ctx := context.Background()
	phoneNumber := interserviceclient.TestUserPhoneNumber
	data := map[string]interface{}{
		"name": "user",
	}
	notification := firebasetools.FirebaseSimpleNotificationInput{}
	android := firebasetools.FirebaseAndroidConfigInput{}
	ios := firebasetools.FirebaseAPNSConfigInput{}
	web := firebasetools.FirebaseWebpushConfigInput{}

	var f fcm.UsecaseFCM = &fakeFCMService

	type args struct {
		ctx          context.Context
		phoneNumber  *string
		email        *string
		data         map[string]interface{}
		notification firebasetools.FirebaseSimpleNotificationInput
		android      *firebasetools.FirebaseAndroidConfigInput
		ios          *firebasetools.FirebaseAPNSConfigInput
		web          *firebasetools.FirebaseWebpushConfigInput
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "valid: correct params passed",
			args: args{
				ctx:          ctx,
				phoneNumber:  &phoneNumber,
				data:         data,
				notification: notification,
				android:      &android,
				ios:          &ios,
				web:          &web,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "invalid: missing phoneNumber, data, notification,  android, ios, web params",
			args: args{
				ctx: ctx,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid: correct params passed" {
				fakeFCMService.SendFCMByPhoneOrEmailFn = func(
					ctx context.Context,
					phoneNumber *string,
					email *string,
					data map[string]interface{},
					notification firebasetools.FirebaseSimpleNotificationInput,
					android *firebasetools.FirebaseAndroidConfigInput,
					ios *firebasetools.FirebaseAPNSConfigInput,
					web *firebasetools.FirebaseWebpushConfigInput,
				) (bool, error) {
					return true, nil
				}
			}
			if tt.name == "invalid: missing phoneNumber, data, notification,  android, ios, web params" {
				fakeFCMService.SendFCMByPhoneOrEmailFn = func(
					ctx context.Context,
					phoneNumber *string,
					email *string,
					data map[string]interface{},
					notification firebasetools.FirebaseSimpleNotificationInput,
					android *firebasetools.FirebaseAndroidConfigInput,
					ios *firebasetools.FirebaseAPNSConfigInput,
					web *firebasetools.FirebaseWebpushConfigInput,
				) (bool, error) {
					return false, errors.New("test error")
				}
			}
			got, err := f.SendFCMByPhoneOrEmail(tt.args.ctx, tt.args.phoneNumber, tt.args.email, tt.args.data, tt.args.notification, tt.args.android, tt.args.ios, tt.args.web)
			if (err != nil) != tt.wantErr {
				t.Errorf("ImplFCM.SendFCMByPhoneOrEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ImplFCM.SendFCMByPhoneOrEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}
