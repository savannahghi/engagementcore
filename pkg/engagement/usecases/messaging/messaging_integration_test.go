package messaging_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/rs/xid"
	"github.com/savannahghi/engagementcore/pkg/engagement/infrastructure"
	"github.com/savannahghi/engagementcore/pkg/engagement/usecases/messaging"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/firebasetools"
	"github.com/segmentio/ksuid"
)

func InitializeTestNewMessaging(ctx context.Context) (*messaging.ImplNotification, infrastructure.Interactor, error) {
	infra := infrastructure.NewInteractor()
	messaging := messaging.NewNotification(infra)
	return messaging, infra, nil
}

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
	f, i, err := InitializeTestNewMessaging(ctx)
	if err != nil {
		t.Errorf("failed to initialize new Messaging: %v", err)
	}

	type args struct {
		infrastructure infrastructure.Interactor
	}
	tests := []struct {
		name string
		args args
		want *messaging.ImplNotification
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
			if got := messaging.NewNotification(tt.args.infrastructure); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMessaging() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPubSubNotificationService_Notify(t *testing.T) {
	ctx := context.Background()
	_, i, err := InitializeTestNewMessaging(ctx)
	if err != nil {
		t.Errorf("failed to initialize new Messaging: %v", err)
	}

	s := messaging.NewNotification(i)
	if s == nil {
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
		pubsub  messaging.UsecaseNotification
		args    args
		wantErr bool
	}{
		{
			pubsub: s,
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
			pubsub: s,
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

func TestImplNotification_TopicIDs(t *testing.T) {
	ctx := context.Background()
	f, _, err := InitializeTestNewMessaging(ctx)
	if err != nil {
		t.Errorf("failed to initialize new Messaging: %v", err)
	}
	tests := []struct {
		name    string
		wantNil bool
	}{
		{
			name:    "default case",
			wantNil: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := f.TopicIDs()
			if !tt.wantNil && got == nil {
				t.Errorf("PubSubNotificationService.TopicIDs() did not expect nil, got %v", got)
			}
		})
	}
}

func TestPubSubNotificationService_SubscriptionIDs(t *testing.T) {
	ctx := context.Background()
	f, _, err := InitializeTestNewMessaging(ctx)
	if err != nil {
		t.Errorf("failed to initialize new Messaging: %v", err)
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
			if got := f.SubscriptionIDs(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PubSubNotificationService.ReverseSubscriptionIDs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPubSubNotificationService_ReverseSubscriptionIDs(t *testing.T) {
	ctx := context.Background()
	f, _, err := InitializeTestNewMessaging(ctx)
	if err != nil {
		t.Errorf("failed to initialize new Messaging: %v", err)
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
			if got := f.ReverseSubscriptionIDs(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PubSubNotificationService.ReverseSubscriptionIDs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRemotePushService_Push(t *testing.T) {
	ctx := context.Background()
	_, i, err := InitializeTestNewMessaging(ctx)
	if err != nil {
		t.Errorf("failed to initialize new Messaging: %v", err)
	}

	s := messaging.NewNotification(i)
	if s == nil {
		t.Errorf("nil remote Messaging push service")
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
			if err := s.Push(tt.args.ctx, tt.args.sender, tt.args.notificationPayload); (err != nil) != tt.wantErr {
				t.Errorf("RemotePushService.Push() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
