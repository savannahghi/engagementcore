package messaging_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/rs/xid"
	"github.com/savannahghi/engagementcore/pkg/engagement/infrastructure/services/messaging"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/serverutils"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
)

func getNotificationPayload(t *testing.T) firebasetools.SendNotificationPayload {
	imgURL := "https://example.com/img.png"
	return firebasetools.SendNotificationPayload{
		RegistrationTokens: []string{xid.New().String(), xid.New().String()},
		Data: map[string]string{
			xid.New().String(): xid.New().String(),
			xid.New().String(): xid.New().String(),
		},
		Notification: &firebasetools.FirebaseSimpleNotificationInput{
			Title:    xid.New().String(),
			Body:     xid.New().String(),
			ImageURL: &imgURL,
			Data: map[string]interface{}{
				xid.New().String(): xid.New().String(),
				xid.New().String(): xid.New().String(),
			},
		},
	}
}

func TestNewPubSubNotificationService(t *testing.T) {
	ctx := context.Background()
	projectID := serverutils.MustGetEnvVar(serverutils.GoogleCloudProjectIDEnvVarName)

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name: "default case",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := messaging.NewPubSubNotificationService(ctx, projectID)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"NewPubSubNotificationService() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if !tt.wantErr {
				assert.NotNil(t, got)
			}
		})
	}
}

func TestPubSubNotificationService_Notify(t *testing.T) {
	ctx := context.Background()
	projectID := serverutils.MustGetEnvVar(serverutils.GoogleCloudProjectIDEnvVarName)
	srv, err := messaging.NewPubSubNotificationService(ctx, projectID)
	if err != nil {
		t.Errorf("can't initialize pubsub notification service: %s", err)
		return
	}

	if srv == nil {
		t.Errorf("nil pubsub notification service")
		return
	}

	type args struct {
		channel  string
		uid      string
		flavour  feedlib.Flavour
		el       feedlib.Element
		metadata map[string]interface{}
	}
	tests := []struct {
		name    string
		pubsub  messaging.NotificationService
		args    args
		wantErr bool
	}{
		{
			pubsub: srv,
			args: args{
				channel: "message.post",
				el: &feedlib.Message{
					ID:             ksuid.New().String(),
					SequenceNumber: 1,
					Text:           ksuid.New().String(),
					ReplyTo:        ksuid.New().String(),
					PostedByUID:    ksuid.New().String(),
					PostedByName:   ksuid.New().String(),
					Timestamp:      time.Now(),
				},
				uid:      ksuid.New().String(),
				flavour:  feedlib.FlavourConsumer,
				metadata: map[string]interface{}{},
			},
			wantErr: false,
		},
		{
			name:   "invalid message, missing posted by info",
			pubsub: srv,
			args: args{
				channel: "message.post",
				el: &feedlib.Message{
					ID:        ksuid.New().String(),
					Text:      ksuid.New().String(),
					ReplyTo:   ksuid.New().String(),
					Timestamp: time.Now(),
				},
				uid:      ksuid.New().String(),
				flavour:  feedlib.FlavourPro,
				metadata: map[string]interface{}{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.pubsub.Notify(
				context.Background(),
				tt.args.channel,
				tt.args.uid,
				tt.args.flavour,
				tt.args.el,
				tt.args.metadata,
			); (err != nil) != tt.wantErr {
				t.Errorf(
					"PubSubNotificationService.Notify() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
			}
		})
	}
}

func TestPubSubNotificationService_SubscriptionIDs(t *testing.T) {
	ctx := context.Background()
	projectID := serverutils.MustGetEnvVar(serverutils.GoogleCloudProjectIDEnvVarName)
	srv, err := messaging.NewPubSubNotificationService(ctx, projectID)
	if err != nil {
		t.Errorf("can't initialize pubsub notification service: %s", err)
		return
	}

	if srv == nil {
		t.Errorf("nil pubsub notification service")
		return
	}

	tests := []struct {
		name string
		want map[string]string
	}{
		{
			name: "default case",
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := srv.SubscriptionIDs(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PubSubNotificationService.ReverseSubscriptionIDs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPubSubNotificationService_ReverseSubscriptionIDs(t *testing.T) {
	ctx := context.Background()
	projectID := serverutils.MustGetEnvVar(serverutils.GoogleCloudProjectIDEnvVarName)
	srv, err := messaging.NewPubSubNotificationService(ctx, projectID)
	if err != nil {
		t.Errorf("can't initialize pubsub notification service: %s", err)
		return
	}
	tests := []struct {
		name string
		want map[string]string
	}{
		{
			name: "default case",
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := srv.ReverseSubscriptionIDs(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PubSubNotificationService.ReverseSubscriptionIDs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRemotePushService_Push(t *testing.T) {
	ctx := context.Background()
	projectID := serverutils.MustGetEnvVar(serverutils.GoogleCloudProjectIDEnvVarName)
	srv, err := messaging.NewPubSubNotificationService(ctx, projectID)
	if err != nil {
		t.Errorf("can't initialize pubsub notification service: %s", err)
		return
	}
	if srv == nil {
		t.Errorf("nil remote FCM push service")
		return
	}

	type args struct {
		ctx                 context.Context
		sender              string
		notificationPayload firebasetools.SendNotificationPayload
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid send - won't actually push but won't error",
			args: args{
				ctx:                 ctx,
				sender:              "test",
				notificationPayload: getNotificationPayload(t),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := srv.Push(tt.args.ctx, tt.args.sender, tt.args.notificationPayload); (err != nil) != tt.wantErr {
				t.Errorf("RemotePushService.Push() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
