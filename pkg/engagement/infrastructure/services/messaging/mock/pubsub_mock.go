package mock

import (
	"context"

	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/firebasetools"
)

// FakeServiceMessaging is a mock implementation of the "messaging" service
type FakeServiceMessaging struct {
	NotifyFn func(
		ctx context.Context,
		topicID string,
		uid string,
		flavour feedlib.Flavour,
		payload feedlib.Element,
		metadata map[string]interface{},
	) error

	// Ask the notification service about the topics that it knows about
	TopicIDsFn func() []string

	SubscriptionIDsFn func() map[string]string

	ReverseSubscriptionIDsFn func() map[string]string
	PushFn                   func(
		ctx context.Context,
		sender string,
		payload firebasetools.SendNotificationPayload,
	) error
}

// Notify Sends a message to a topic
func (f *FakeServiceMessaging) Notify(
	ctx context.Context,
	topicID string,
	uid string,
	flavour feedlib.Flavour,
	payload feedlib.Element,
	metadata map[string]interface{},
) error {
	return f.NotifyFn(ctx, topicID, uid, flavour, payload, metadata)
}

// TopicIDs gets topic IDs
func (f *FakeServiceMessaging) TopicIDs() []string {
	return f.TopicIDsFn()
}

// SubscriptionIDs gets subscription IDs
func (f *FakeServiceMessaging) SubscriptionIDs() map[string]string {
	return f.SubscriptionIDsFn()
}

// ReverseSubscriptionIDs get Reverse Subscription IDs
func (f *FakeServiceMessaging) ReverseSubscriptionIDs() map[string]string {
	return f.ReverseSubscriptionIDsFn()
}

// Push sends push notifications
func (f *FakeServiceMessaging) Push(
	ctx context.Context,
	sender string,
	payload firebasetools.SendNotificationPayload,
) error {
	return f.PushFn(ctx, sender, payload)
}
