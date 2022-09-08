package messaging

import (
	"context"

	"github.com/savannahghi/engagementcore/pkg/engagement/infrastructure"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/firebasetools"
)

// UsecaseNotification defines pupsub notification service usecases interface
type UsecaseNotification interface {
	// Send a message to a topic
	Notify(
		ctx context.Context,
		topicID string,
		uid string,
		flavour feedlib.Flavour,
		payload feedlib.Element,
		metadata map[string]interface{},
	) error

	// Ask the notification service about the topics that it knows about
	TopicIDs() []string

	SubscriptionIDs() map[string]string

	ReverseSubscriptionIDs() map[string]string

	Push(
		ctx context.Context,
		sender string,
		payload firebasetools.SendNotificationPayload,
	) error
}

// ImplNotification is the pubsub notification service implementation
type ImplNotification struct {
	infrastructure infrastructure.Interactor
}

// NewNotification initializes a pubsub notification service instance
func NewNotification(infrastructure infrastructure.Interactor) *ImplNotification {
	return &ImplNotification{
		infrastructure: infrastructure,
	}
}

// Notify Send a message to a topic
func (n *ImplNotification) Notify(
	ctx context.Context,
	topicID string,
	uid string,
	flavour feedlib.Flavour,
	payload feedlib.Element,
	metadata map[string]interface{},
) error {
	i := n.infrastructure.NotificationService
	return i.Notify(ctx, topicID, uid, flavour, payload, metadata)
}

// TopicIDs Ask the notification service about the topics that it knows about
func (n *ImplNotification) TopicIDs() []string {
	i := n.infrastructure.NotificationService
	return i.TopicIDs()
}

// SubscriptionIDs gets subscription IDs for the notification service
func (n *ImplNotification) SubscriptionIDs() map[string]string {
	i := n.infrastructure.NotificationService
	return i.SubscriptionIDs()
}

// ReverseSubscriptionIDs gets reverse subscription IDs for the notification service
func (n *ImplNotification) ReverseSubscriptionIDs() map[string]string {
	i := n.infrastructure.NotificationService
	return i.ReverseSubscriptionIDs()
}

// Push instructs a remote FCM service to send a push notification.
//
// This is done over Google Cloud Pub-Sub.
func (n *ImplNotification) Push(
	ctx context.Context,
	sender string,
	payload firebasetools.SendNotificationPayload,
) error {
	i := n.infrastructure.NotificationService
	return i.Push(ctx, sender, payload)
}
