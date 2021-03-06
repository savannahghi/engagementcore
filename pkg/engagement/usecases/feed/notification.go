package feed

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/savannahghi/converterandformatter"
	"github.com/savannahghi/engagementcore/pkg/engagement/application/common"
	"github.com/savannahghi/engagementcore/pkg/engagement/infrastructure"
	"github.com/savannahghi/engagementcore/pkg/engagement/infrastructure/services/onboarding"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/firebasetools"

	"github.com/savannahghi/engagementcore/pkg/engagement/application/common/dto"
	"github.com/savannahghi/engagementcore/pkg/engagement/application/common/helpers"
	"github.com/savannahghi/pubsubtools"
)

const (
	itemPublishSender   = "ITEM_PUBLISHED"
	itemDeleteSender    = "ITEM_DELETED"
	itemResolveSender   = "ITEM_RESOLVED"
	itemUnresolveSender = "ITEM_UNRESOLVED"
	itemHideSender      = "ITEM_HIDE"
	itemShowSender      = "ITEM_SHOW"
	itemPinSender       = "ITEM_PIN"
	itemUnpinSender     = "ITEM_UNPIN"

	nudgePublishSender   = "NUDGE_PUBLISHED"
	nudgeDeleteSender    = "NUDGE_DELETED"
	nudgeResolveSender   = "NUDGE_RESOLVED"
	nudgeUnresolveSender = "NUDGE_UNRESOLVED"
	nudgeShowSender      = "NUDGE_SHOW"
	nudgeHideSender      = "NUDGE_HIDE"

	feedUpdate       = "FEED_UPDATE"
	inboxCountUpdate = "INBOX_COUNT_CHANGED"
)

// NotificationUsecases represent logic required to make notification
type NotificationUsecases interface {
	HandleItemPublish(
		ctx context.Context,
		m *pubsubtools.PubSubPayload,
	) error

	HandleItemDelete(
		ctx context.Context,
		m *pubsubtools.PubSubPayload,
	) error

	HandleItemResolve(
		ctx context.Context,
		m *pubsubtools.PubSubPayload,
	) error

	HandleItemUnresolve(
		ctx context.Context,
		m *pubsubtools.PubSubPayload,
	) error

	HandleItemHide(
		ctx context.Context,
		m *pubsubtools.PubSubPayload,
	) error

	HandleItemShow(
		ctx context.Context,
		m *pubsubtools.PubSubPayload,
	) error

	HandleItemPin(
		ctx context.Context,
		m *pubsubtools.PubSubPayload,
	) error

	HandleItemUnpin(
		ctx context.Context,
		m *pubsubtools.PubSubPayload,
	) error

	HandleNudgePublish(
		ctx context.Context,
		m *pubsubtools.PubSubPayload,
	) error

	HandleNudgeDelete(
		ctx context.Context,
		m *pubsubtools.PubSubPayload,
	) error

	HandleNudgeResolve(
		ctx context.Context,
		m *pubsubtools.PubSubPayload,
	) error

	HandleNudgeUnresolve(
		ctx context.Context,
		m *pubsubtools.PubSubPayload,
	) error

	HandleNudgeHide(
		ctx context.Context,
		m *pubsubtools.PubSubPayload,
	) error

	HandleNudgeShow(
		ctx context.Context,
		m *pubsubtools.PubSubPayload,
	) error

	HandleActionPublish(
		ctx context.Context,
		m *pubsubtools.PubSubPayload,
	) error

	HandleActionDelete(
		ctx context.Context,
		m *pubsubtools.PubSubPayload,
	) error

	HandleMessagePost(
		ctx context.Context,
		m *pubsubtools.PubSubPayload,
	) error

	HandleMessageDelete(
		ctx context.Context,
		m *pubsubtools.PubSubPayload,
	) error

	HandleIncomingEvent(
		ctx context.Context,
		m *pubsubtools.PubSubPayload,
	) error

	NotifyItemUpdate(
		ctx context.Context,
		sender string,
		includeNotification bool, // whether to show a tray notification
		m *pubsubtools.PubSubPayload,
	) error

	UpdateInbox(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
	) error

	NotifyNudgeUpdate(
		ctx context.Context,
		sender string,
		m *pubsubtools.PubSubPayload,
	) error

	NotifyInboxCountUpdate(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
		count int,
	) error

	GetUserTokens(
		ctx context.Context,
		uids []string,
	) ([]string, error)

	SendNotificationViaFCM(
		ctx context.Context,
		uids []string,
		sender string,
		pl dto.NotificationEnvelope,
		notification *firebasetools.FirebaseSimpleNotificationInput,
	) error

	HandleSendNotification(
		ctx context.Context,
		m *pubsubtools.PubSubPayload,
	) error

	SendNotificationEmail(
		ctx context.Context,
		m *pubsubtools.PubSubPayload,
	) error
}

// HandlePubsubPayload defines the signature of a function that handles
// payloads received from Google Cloud Pubsub
type HandlePubsubPayload func(ctx context.Context, m *pubsubtools.PubSubPayload) error

// NotificationImpl represents the notification usecase implementation
type NotificationImpl struct {
	infrastructure infrastructure.Interactor
}

// NewNotification initializes a notification usecase
func NewNotification(infrastructure infrastructure.Interactor) *NotificationImpl {
	return &NotificationImpl{
		infrastructure: infrastructure,
	}
}

// HandleItemPublish responds to item publish messages
func (n NotificationImpl) HandleItemPublish(
	ctx context.Context,
	m *pubsubtools.PubSubPayload,
) error {
	ctx, span := tracer.Start(ctx, "HandleItemPublish")
	defer span.End()
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}
	err := n.NotifyItemUpdate(ctx, itemPublishSender, true, m)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return fmt.Errorf("can't notify item update over FCM: %w", err)
	}

	return nil
}

// HandleItemDelete responds to item delete messages
func (n NotificationImpl) HandleItemDelete(
	ctx context.Context,
	m *pubsubtools.PubSubPayload,
) error {
	ctx, span := tracer.Start(ctx, "HandleItemDelete")
	defer span.End()
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}

	err := n.NotifyItemUpdate(ctx, itemDeleteSender, false, m)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return fmt.Errorf("can't notify item update over FCM: %w", err)
	}

	return nil
}

// HandleItemResolve responds to item resolve messages
func (n NotificationImpl) HandleItemResolve(
	ctx context.Context,
	m *pubsubtools.PubSubPayload,
) error {
	ctx, span := tracer.Start(ctx, "HandleItemResolve")
	defer span.End()
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}

	err := n.NotifyItemUpdate(ctx, itemResolveSender, false, m)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return fmt.Errorf("can't notify item update over FCM: %w", err)
	}

	return nil
}

// HandleItemUnresolve responds to item unresolve messages
func (n NotificationImpl) HandleItemUnresolve(
	ctx context.Context,
	m *pubsubtools.PubSubPayload,
) error {
	ctx, span := tracer.Start(ctx, "HandleItemUnresolve")
	defer span.End()
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}

	err := n.NotifyItemUpdate(ctx, itemUnresolveSender, false, m)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return fmt.Errorf("can't notify item update over FCM: %w", err)
	}

	return nil
}

// HandleItemHide responds to item hide messages
func (n NotificationImpl) HandleItemHide(
	ctx context.Context,
	m *pubsubtools.PubSubPayload,
) error {
	ctx, span := tracer.Start(ctx, "HandleItemHide")
	defer span.End()

	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}

	err := n.NotifyItemUpdate(ctx, itemHideSender, false, m)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return fmt.Errorf("can't notify item update over FCM: %w", err)
	}

	return nil
}

// HandleItemShow responds to item show messages
func (n NotificationImpl) HandleItemShow(
	ctx context.Context,
	m *pubsubtools.PubSubPayload,
) error {
	ctx, span := tracer.Start(ctx, "HandleItemShow")
	defer span.End()
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}

	err := n.NotifyItemUpdate(ctx, itemShowSender, false, m)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return fmt.Errorf("can't notify item update over FCM: %w", err)
	}

	return nil
}

// HandleItemPin responds to item pin messages
func (n NotificationImpl) HandleItemPin(
	ctx context.Context,
	m *pubsubtools.PubSubPayload,
) error {
	ctx, span := tracer.Start(ctx, "HandleItemPin")
	defer span.End()
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}

	err := n.NotifyItemUpdate(ctx, itemPinSender, false, m)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return fmt.Errorf("can't notify item update over FCM: %w", err)
	}

	return nil
}

// HandleItemUnpin responds to item unpin messages
func (n NotificationImpl) HandleItemUnpin(
	ctx context.Context,
	m *pubsubtools.PubSubPayload,
) error {
	ctx, span := tracer.Start(ctx, "HandleItemUnpin")
	defer span.End()
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}

	err := n.NotifyItemUpdate(ctx, itemUnpinSender, false, m)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return fmt.Errorf("can't notify item update over FCM: %w", err)
	}

	return nil
}

// HandleNudgePublish responds to nudge publish messages
func (n NotificationImpl) HandleNudgePublish(
	ctx context.Context,
	m *pubsubtools.PubSubPayload,
) error {
	ctx, span := tracer.Start(ctx, "HandleNudgePublish")
	defer span.End()
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}

	err := n.NotifyNudgeUpdate(ctx, nudgePublishSender, m)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return fmt.Errorf("can't notify nudge update over FCM: %w", err)
	}

	return nil
}

// HandleNudgeDelete responds to nudge delete messages
func (n NotificationImpl) HandleNudgeDelete(
	ctx context.Context,
	m *pubsubtools.PubSubPayload,
) error {
	ctx, span := tracer.Start(ctx, "HandleNudgeDelete")
	defer span.End()
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}

	err := n.NotifyNudgeUpdate(ctx, nudgeDeleteSender, m)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return fmt.Errorf("can't notify nudge update over FCM: %w", err)
	}

	return nil
}

// HandleNudgeResolve responds to nudge resolve messages
func (n NotificationImpl) HandleNudgeResolve(
	ctx context.Context,
	m *pubsubtools.PubSubPayload,
) error {
	ctx, span := tracer.Start(ctx, "HandleNudgeResolve")
	defer span.End()

	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}

	err := n.NotifyNudgeUpdate(ctx, nudgeResolveSender, m)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return fmt.Errorf("can't notify nudge update over FCM: %w", err)
	}

	return nil
}

// HandleNudgeUnresolve responds to nudge unresolve messages
func (n NotificationImpl) HandleNudgeUnresolve(
	ctx context.Context,
	m *pubsubtools.PubSubPayload,
) error {
	ctx, span := tracer.Start(ctx, "HandleNudgeUnresolve")
	defer span.End()
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}

	err := n.NotifyNudgeUpdate(ctx, nudgeUnresolveSender, m)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return fmt.Errorf("can't notify nudge update over FCM: %w", err)
	}

	return nil
}

// HandleNudgeHide responds to nudge hide messages
func (n NotificationImpl) HandleNudgeHide(
	ctx context.Context,
	m *pubsubtools.PubSubPayload,
) error {
	ctx, span := tracer.Start(ctx, "HandleNudgeHide")
	defer span.End()
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}

	err := n.NotifyNudgeUpdate(ctx, nudgeHideSender, m)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return fmt.Errorf("can't notify nudge update over FCM: %w", err)
	}

	return nil
}

// HandleNudgeShow responds to nudge hide messages
func (n NotificationImpl) HandleNudgeShow(
	ctx context.Context,
	m *pubsubtools.PubSubPayload,
) error {
	ctx, span := tracer.Start(ctx, "HandleNudgeShow")
	defer span.End()
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}

	err := n.NotifyNudgeUpdate(ctx, nudgeShowSender, m)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return fmt.Errorf("can't notify nudge update over FCM: %w", err)
	}

	return nil
}

// HandleActionPublish responds to action publish messages
func (n NotificationImpl) HandleActionPublish(
	ctx context.Context,
	m *pubsubtools.PubSubPayload,
) error {
	_, span := tracer.Start(ctx, "HandleActionPublish")
	defer span.End()
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}

	// TODO Notify action publish

	return nil
}

// HandleActionDelete responds to action publish messages
func (n NotificationImpl) HandleActionDelete(
	ctx context.Context,
	m *pubsubtools.PubSubPayload,
) error {
	_, span := tracer.Start(ctx, "HandleActionDelete")
	defer span.End()

	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}

	// TODO Notify action delete

	return nil
}

// HandleMessagePost responds to message post pubsub messages
func (n NotificationImpl) HandleMessagePost(
	ctx context.Context,
	m *pubsubtools.PubSubPayload,
) error {
	_, span := tracer.Start(ctx, "HandleMessagePost")
	defer span.End()
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}

	// TODO Notify the message and it's context i.e item and feed flavour

	return nil
}

// HandleMessageDelete responds to message delete pubsub messages
func (n NotificationImpl) HandleMessageDelete(
	ctx context.Context,
	m *pubsubtools.PubSubPayload,
) error {
	_, span := tracer.Start(ctx, "HandleMessageDelete")
	defer span.End()
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}

	// TODO Notify the message delete and it's context i.e item and feed flavour

	return nil
}

// HandleIncomingEvent responds to message delete pubsub messages
func (n NotificationImpl) HandleIncomingEvent(
	ctx context.Context,
	m *pubsubtools.PubSubPayload,
) error {
	_, span := tracer.Start(ctx, "HandleIncomingEvent")
	defer span.End()

	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}

	log.Printf("incoming event data: \n%s\n", string(m.Message.Data))
	log.Printf("incoming event subscription: %s", m.Subscription)
	log.Printf("incoming event message ID: %s", m.Message.MessageID)
	log.Printf("incoming event message attributes: %#v", m.Message.Attributes)

	return nil
}

// HandleSendNotification responds to send notification messages
func (n NotificationImpl) HandleSendNotification(
	ctx context.Context,
	m *pubsubtools.PubSubPayload,
) error {
	ctx, span := tracer.Start(ctx, "HandleSendNotification")
	defer span.End()
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}

	payload := &firebasetools.SendNotificationPayload{}
	err := json.Unmarshal(m.Message.Data, payload)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return fmt.Errorf(
			"can't unmarshal notification notification from pubsub data: %w",
			err,
		)
	}

	_, err = n.infrastructure.SendNotification(
		ctx,
		payload.RegistrationTokens,
		payload.Data,
		payload.Notification,
		payload.Android,
		payload.Ios,
		payload.Web,
	)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return fmt.Errorf("can't send notification: %v", err)
	}

	return nil
}

// NotifyItemUpdate sends a Firebase Cloud Messaging notification
func (n NotificationImpl) NotifyItemUpdate(
	ctx context.Context,
	sender string,
	includeNotification bool, // whether to show a tray notification
	m *pubsubtools.PubSubPayload,
) error {
	ctx, span := tracer.Start(ctx, "NotifyItemUpdate")
	defer span.End()
	var envelope dto.NotificationEnvelope
	err := json.Unmarshal(m.Message.Data, &envelope)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return fmt.Errorf(
			"can't unmarshal notification envelope from pubsub data: %w", err)
	}

	var item feedlib.Item
	err = json.Unmarshal(envelope.Payload, &item)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return fmt.Errorf("can't unmarshal item from pubsub data: %w", err)
	}
	// include notifications for persistent items
	var notification *firebasetools.FirebaseSimpleNotificationInput
	iconURL := common.DefaultIconPath
	if item.Persistent && includeNotification {
		// also include a notification
		notification = &firebasetools.FirebaseSimpleNotificationInput{
			Title:    item.Tagline,
			Body:     item.Summary,
			ImageURL: &iconURL,
		}

		err = n.SendNotificationViaFCM(
			ctx,
			item.Users,
			sender,
			envelope,
			notification,
		)
		if err != nil {
			helpers.RecordSpanError(span, err)
			return fmt.Errorf("unable to notify item: %w", err)
		}
	}

	// TODO Send email notifications
	// TODO For urgent (tray), consider whatsapp and sms notifications

	switch sender {
	case itemPublishSender:
		existingLabels, err := n.infrastructure.Labels(
			ctx,
			envelope.UID,
			envelope.Flavour,
		)
		if err != nil {
			helpers.RecordSpanError(span, err)
			return fmt.Errorf("can't fetch existing labels: %w", err)
		}

		if !converterandformatter.StringSliceContains(
			existingLabels,
			item.Label,
		) {
			err = n.infrastructure.SaveLabel(
				ctx,
				envelope.UID,
				envelope.Flavour,
				item.Label,
			)
			if err != nil {
				helpers.RecordSpanError(span, err)
				return fmt.Errorf("can't save label: %w", err)
			}
		}
	case itemDeleteSender,
		itemResolveSender,
		itemUnresolveSender,
		itemHideSender,
		itemShowSender,
		itemPinSender,
		itemUnpinSender:
		// do nothing...inbox update code will run in the outer scope
	default:
		return fmt.Errorf("unexpected item publish sender: %s", sender)
	}

	err = n.UpdateInbox(ctx, envelope.UID, envelope.Flavour)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return fmt.Errorf("unable to update inbox count: %w", err)
	}

	return nil
}

// UpdateInbox recalculates the inbox count and notifies the client over FCM
func (n NotificationImpl) UpdateInbox(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
) error {
	ctx, span := tracer.Start(ctx, "UpdateInbox")
	defer span.End()
	err := n.infrastructure.UpdateUnreadPersistentItemsCount(ctx, uid, flavour)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return fmt.Errorf("can't update inbox count: %w", err)
	}

	_, err = n.infrastructure.UnreadPersistentItems(ctx, uid, flavour)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return fmt.Errorf("can't get inbox count: %w", err)
	}

	// The inbox has been descoped for this milestone
	// Does not make sense to send notification updates to our users
	// TODO: Restore after the milestone @mathenge
	// err = n.NotifyInboxCountUpdate(ctx, uid, flavour, unread)
	// if err != nil {
	// 	return fmt.Errorf("can't notify inbox count: %w", err)
	// }

	return nil
}

// NotifyNudgeUpdate sends a nudge update notification via FCM
func (n NotificationImpl) NotifyNudgeUpdate(
	ctx context.Context,
	sender string,
	m *pubsubtools.PubSubPayload,
) error {
	ctx, span := tracer.Start(ctx, "NotifyNudgeUpdate")
	defer span.End()
	var envelope dto.NotificationEnvelope
	err := json.Unmarshal(m.Message.Data, &envelope)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return fmt.Errorf(
			"can't unmarshal notification envelope from pubsub data: %w",
			err,
		)
	}

	var nudge feedlib.Nudge
	err = json.Unmarshal(envelope.Payload, &nudge)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return fmt.Errorf("can't unmarshal nudge from pubsub data: %w", err)
	}

	links := nudge.Links
	var imageURL string
	for _, link := range links {
		imageURL = link.Thumbnail
	}

	var notification *firebasetools.FirebaseSimpleNotificationInput

	switch sender {
	case nudgePublishSender:
		notification = &firebasetools.FirebaseSimpleNotificationInput{
			Title:    nudge.Title,
			Body:     nudge.NotificationBody.PublishMessage,
			ImageURL: &imageURL,
		}

	case nudgeResolveSender:
		notification = &firebasetools.FirebaseSimpleNotificationInput{
			Title:    nudge.Title,
			Body:     nudge.NotificationBody.ResolveMessage,
			ImageURL: &imageURL,
		}

	case nudgeDeleteSender,
		nudgeUnresolveSender,
		nudgeShowSender,
		nudgeHideSender:
		// Do nothing..our scope for nudges does not contain these
		return nil

	default:
		return fmt.Errorf("unexpected nudge sender: %s", sender)
	}

	err = n.SendNotificationViaFCM(
		ctx,
		nudge.Users,
		sender,
		envelope,
		notification,
	)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return fmt.Errorf("unable to notify nudge: %w", err)
	}

	return nil
}

// NotifyInboxCountUpdate sends a message notifying of an update to inbox
// item counts.
func (n NotificationImpl) NotifyInboxCountUpdate(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	count int,
) error {
	ctx, span := tracer.Start(ctx, "NotifyInboxCountUpdate")
	defer span.End()
	notificationEnvelope := dto.NotificationEnvelope{
		UID:     uid,
		Flavour: flavour,
		Payload: []byte(fmt.Sprintf("%d", count)),
		Metadata: map[string]interface{}{
			"sender": inboxCountUpdate,
			"count":  count,
		},
	}

	notification := &firebasetools.FirebaseSimpleNotificationInput{
		Title: "Be.Well Inbox",
		Body:  fmt.Sprintf("You have %v unread notification(s).", count),
	}

	notifyUIDs := []string{uid}
	err := n.SendNotificationViaFCM(
		ctx, notifyUIDs, feedUpdate, notificationEnvelope, notification)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return fmt.Errorf("unable to notify thin feed: %w", err)
	}

	return nil
}

// GetUserTokens retrieves the user tokens corresponding to the supplied UIDs
func (n NotificationImpl) GetUserTokens(
	ctx context.Context,
	uids []string,
) ([]string, error) {
	ctx, span := tracer.Start(ctx, "GetUserTokens")
	defer span.End()
	userTokens, err := n.infrastructure.GetDeviceTokens(ctx, onboarding.UserUIDs{
		UIDs: uids,
	})
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf("can't get push tokens: %w", err)
	}
	tokens := []string{}
	for _, toks := range userTokens {
		tokens = append(tokens, toks...)
	}
	return tokens, nil
}

// SendNotificationViaFCM publishes an FCM notification
func (n NotificationImpl) SendNotificationViaFCM(
	ctx context.Context,
	uids []string,
	sender string,
	pl dto.NotificationEnvelope,
	notification *firebasetools.FirebaseSimpleNotificationInput,
) error {
	ctx, span := tracer.Start(ctx, "SendNotificationViaFCM")
	defer span.End()
	if notification == nil {
		return fmt.Errorf("nil notification")
	}

	tokens, err := n.GetUserTokens(ctx, uids)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return fmt.Errorf("can't get user tokens: %w", err)
	}
	if len(tokens) == 0 {
		return nil
	}
	marshalled, err := json.Marshal(pl)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return fmt.Errorf(
			"can't send element that failed validation over FCM: %w", err)
	}
	payload := firebasetools.SendNotificationPayload{
		RegistrationTokens: tokens,
		Data: map[string]string{
			sender: string(marshalled),
		},
		Notification: &firebasetools.FirebaseSimpleNotificationInput{
			Title:    notification.Title,
			Body:     notification.Body,
			Data:     notification.Data,
			ImageURL: notification.ImageURL,
		},
	}

	err = n.infrastructure.Push(ctx, sender, payload)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return fmt.Errorf("can't send element over FCM: %w", err)
	}
	return nil
}

// SendNotificationEmail sends an email
func (n NotificationImpl) SendNotificationEmail(
	ctx context.Context,
	m *pubsubtools.PubSubPayload,
) error {
	ctx, span := tracer.Start(ctx, "SendNotificationEmail")
	defer span.End()
	if m == nil {
		return fmt.Errorf("nil pub sub payload")
	}

	payload := &dto.EMailMessage{}
	err := json.Unmarshal(m.Message.Data, &payload)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return fmt.Errorf("failed to unmarshal data: %v", err)
	}

	_, _, err = n.infrastructure.SendEmail(
		ctx,
		payload.Subject,
		payload.Text,
		nil,
		payload.To...,
	)

	if err != nil {
		helpers.RecordSpanError(span, err)
		return fmt.Errorf("unable to send email: %v", err)
	}
	return nil
}
