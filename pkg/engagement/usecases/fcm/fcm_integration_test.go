package fcm_test

// TODO: get register push tokens working and add happy case
import (
	"context"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/savannahghi/engagementcore/pkg/engagement/application/common/dto"
	"github.com/savannahghi/engagementcore/pkg/engagement/infrastructure"
	"github.com/savannahghi/engagementcore/pkg/engagement/usecases/fcm"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/interserviceclient"
)

const (
	onboardingService = "profile"
	intMax            = 9007199254740990
	registerPushToken = "testing/register_push_token"
)

func TestMain(m *testing.M) {
	os.Setenv("ROOT_COLLECTION_SUFFIX", "staging")
	os.Exit(m.Run())
}

func InitializeTestNewFCM(ctx context.Context) (*fcm.ImplFCM, infrastructure.Interactor, error) {
	infra := infrastructure.NewInteractor()
	fcm := fcm.NewFCM(infra)
	return fcm, infra, nil
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

func RegisterPushToken(
	ctx context.Context,
	t *testing.T,
	UID string,
	onboardingClient *interserviceclient.InterServiceClient,
) (bool, error) {
	token := "random"
	if onboardingClient == nil {
		return false, fmt.Errorf("nil ISC client")
	}

	payload := map[string]interface{}{
		"pushTokens": token,
		"uid":        UID,
	}
	resp, err := onboardingClient.MakeRequest(
		ctx,
		http.MethodPost,
		registerPushToken,
		payload,
	)
	if err != nil {
		return false, fmt.Errorf("unable to make a request to register push token: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("expected a StatusOK (200) status code but instead got %v", resp.StatusCode)
	}

	return true, nil
}

func TestNewFCM(t *testing.T) {
	ctx := context.Background()
	f, i, err := InitializeTestNewFCM(ctx)
	if err != nil {
		t.Errorf("failed to initialize new FCM: %v", err)
	}

	type args struct {
		infrastructure infrastructure.Interactor
	}
	tests := []struct {
		name string
		args args
		want *fcm.ImplFCM
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
			if got := fcm.NewFCM(tt.args.infrastructure); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFCM() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestImplFCM_SendNotification(t *testing.T) {
	ctx, token, err := interserviceclient.GetPhoneNumberAuthenticatedContextAndToken(
		t,
		onboardingISCClient(t),
	)
	if err != nil {
		t.Errorf("failed to create a test user: %v", err)
		return
	}

	ok, err := RegisterPushToken(ctx, t, token.UID, onboardingISCClient(t))
	if err != nil {
		t.Errorf("failed to get user push tokens: %v", err)
		return
	}
	if !ok {
		t.Errorf("failed to get user push tokens: %v", err)
		return
	}

	f, _, err := InitializeTestNewFCM(ctx)
	if err != nil {
		t.Errorf("failed to initialize new FCM: %v", err)
	}

	fakeToken := uuid.New().String()
	imgURL := "https://www.wxpr.org/sites/wxpr/files/styles/medium/public/202007/chipmunk-5401165_1920.jpg"

	data := map[string]string{
		"name": "user",
	}
	notification := firebasetools.FirebaseSimpleNotificationInput{
		Title:    "Test Notification",
		Body:     "From Integration Tests",
		ImageURL: &imgURL,
		Data: map[string]interface{}{
			"more": "data",
		},
	}
	android := firebasetools.FirebaseAndroidConfigInput{}
	ios := firebasetools.FirebaseAPNSConfigInput{}
	web := firebasetools.FirebaseWebpushConfigInput{}

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
			name: "invalid: non existent token - should fail gracefully",
			args: args{
				ctx:                ctx,
				registrationTokens: []string{fakeToken},
				data:               data,
				notification:       &notification,
				android:            &android,
				ios:                &ios,
				web:                &web,
			},
			wantErr: true,
			want:    false,
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

func TestImplFCM_Notifications(t *testing.T) {
	ctx := firebasetools.GetAuthenticatedContext(t)
	f, _, err := InitializeTestNewFCM(ctx)
	if err != nil {
		t.Errorf("failed to initialize new FCM: %v", err)
	}

	registrationToken := uuid.New().String()
	newerThan := time.Now()
	limit := 10

	notification := []*dto.SavedNotification{}

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
			name: "valid: returns no error with missing args",
			args: args{
				ctx: ctx,
			},
			want:    notification,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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

func TestImplFCM_SendFCMByPhoneOrEmail(t *testing.T) {
	ctx, token, err := interserviceclient.GetPhoneNumberAuthenticatedContextAndToken(
		t,
		onboardingISCClient(t),
	)
	if err != nil {
		t.Errorf("failed to create a test user: %v", err)
		return
	}

	ok, err := RegisterPushToken(ctx, t, token.UID, onboardingISCClient(t))
	if err != nil {
		t.Errorf("failed to get user push tokens: %v", err)
		return
	}
	if !ok {
		t.Errorf("failed to get user push tokens: %v", err)
		return
	}

	f, _, err := InitializeTestNewFCM(ctx)
	if err != nil {
		t.Errorf("failed to initialize new FCM: %v", err)
	}

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
