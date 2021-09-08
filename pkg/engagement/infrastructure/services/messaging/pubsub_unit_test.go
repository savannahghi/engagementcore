package messaging_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/savannahghi/engagementcore/pkg/engagement/infrastructure/services/messaging"
	messagingMock "github.com/savannahghi/engagementcore/pkg/engagement/infrastructure/services/messaging/mock"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/firebasetools"
	"github.com/segmentio/ksuid"
)

var (
	fakemessagingService messagingMock.FakeServiceMessaging
)

func getSendNotificationPayload() *firebasetools.SendNotificationPayload {

	img := "https://www.wxpr.org/sites/wxpr/files/styles/medium/public/202007/chipmunk-5401165_1920.jpg"
	key := uuid.New().String()
	fakeToken := uuid.New().String()
	pckg := "video"

	return &firebasetools.SendNotificationPayload{
		RegistrationTokens: []string{fakeToken},
		Data: map[string]string{
			"some": "data",
		},
		Notification: &firebasetools.FirebaseSimpleNotificationInput{
			Title:    "Test Notification",
			Body:     "From Integration Tests",
			ImageURL: &img,
			Data: map[string]interface{}{
				"more": "data",
			},
		},
		Android: &firebasetools.FirebaseAndroidConfigInput{
			Priority:              "high",
			CollapseKey:           &key,
			RestrictedPackageName: &pckg,
		},
	}
}

func TestUnit_Notify(t *testing.T) {
	ctx := context.Background()
	topicID := ksuid.New().String()
	uid := ksuid.New().String()
	flavour := feedlib.FlavourConsumer
	el := &feedlib.Message{
		ID:             ksuid.New().String(),
		SequenceNumber: 1,
		Text:           ksuid.New().String(),
		ReplyTo:        ksuid.New().String(),
		PostedByUID:    ksuid.New().String(),
		PostedByName:   ksuid.New().String(),
		Timestamp:      time.Now(),
	}
	metadata := map[string]interface{}{
		"test": "metadata",
	}

	var s messaging.NotificationService = &fakemessagingService

	type args struct {
		ctx      context.Context
		topicID  string
		uid      string
		flavour  feedlib.Flavour
		el       feedlib.Element
		metadata map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid: correct params passed",
			args: args{
				ctx:      ctx,
				topicID:  topicID,
				uid:      uid,
				flavour:  flavour,
				el:       el,
				metadata: metadata,
			},
			wantErr: false,
		},
		{
			name: "invalid: missing args",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		if tt.name == "valid: correct params passed" {
			fakemessagingService.NotifyFn = func(
				ctx context.Context,
				topicID string,
				uid string,
				flavour feedlib.Flavour,
				el feedlib.Element,
				metadata map[string]interface{},
			) error {
				return nil
			}
		}
		if tt.name == "invalid: missing args" {
			fakemessagingService.NotifyFn = func(
				ctx context.Context,
				topicID string,
				uid string,
				flavour feedlib.Flavour,
				el feedlib.Element,
				metadata map[string]interface{},
			) error {
				return fmt.Errorf("test error")
			}
		}
		t.Run(tt.name, func(t *testing.T) {
			if err := s.Notify(tt.args.ctx, tt.args.topicID, tt.args.uid, tt.args.flavour, tt.args.el, tt.args.metadata); (err != nil) != tt.wantErr {
				t.Errorf("PubSubNotificationService.Notify() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUnit_TopicIDs(t *testing.T) {

	var s messaging.NotificationService = &fakemessagingService

	tests := []struct {
		name string
		want []string
	}{
		{
			name: "default case",
			want: []string{"testIDs"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "default case" {
				fakemessagingService.TopicIDsFn = func() []string {
					return []string{"testIDs"}
				}
			}
			if got := s.TopicIDs(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PubSubNotificationService.TopicIDs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnit_SubscriptionIDs(t *testing.T) {

	var s messaging.NotificationService = &fakemessagingService

	tests := []struct {
		name string
		want map[string]string
	}{
		{
			name: "default case",
			want: map[string]string{"test": "id"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "default case" {
				fakemessagingService.SubscriptionIDsFn = func() map[string]string {
					return map[string]string{"test": "id"}
				}
			}
			if got := s.SubscriptionIDs(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PubSubNotificationService.SubscriptionIDs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnit_ReverseSubscriptionIDs(t *testing.T) {
	var s messaging.NotificationService = &fakemessagingService

	tests := []struct {
		name string
		want map[string]string
	}{
		{
			name: "default case",
			want: map[string]string{"test": "id"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "default case" {
				fakemessagingService.ReverseSubscriptionIDsFn = func() map[string]string {
					return map[string]string{"test": "id"}
				}
			}
			if got := s.ReverseSubscriptionIDs(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PubSubNotificationService.ReverseSubscriptionIDs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnit_Push(t *testing.T) {
	ctx := context.Background()
	sender := firebasetools.TestUserEmail
	notificationPayload := getSendNotificationPayload()

	var s messaging.NotificationService = &fakemessagingService

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
			name: "valid: correct params passed",
			args: args{
				ctx:                 ctx,
				sender:              sender,
				notificationPayload: *notificationPayload,
			},
			wantErr: false,
		},
		{
			name: "invalid: missing args passed",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid: correct params passed" {
				fakemessagingService.PushFn = func(
					ctx context.Context,
					sender string,
					payload firebasetools.SendNotificationPayload,
				) error {
					return nil
				}
			}
			if tt.name == "invalid: missing args passed" {
				fakemessagingService.PushFn = func(
					ctx context.Context,
					sender string,
					payload firebasetools.SendNotificationPayload,
				) error {
					return fmt.Errorf("test error")
				}
			}
			if err := s.Push(tt.args.ctx, tt.args.sender, tt.args.notificationPayload); (err != nil) != tt.wantErr {
				t.Errorf("PubSubNotificationService.Push() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
