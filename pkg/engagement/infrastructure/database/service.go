package database

import (
	"context"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/savannahghi/engagement/pkg/engagement/application/common/dto"
	"github.com/savannahghi/engagement/pkg/engagement/application/common/helpers"
	"github.com/savannahghi/engagement/pkg/engagement/domain"
	fb "github.com/savannahghi/engagement/pkg/engagement/infrastructure/database/firestore"
	"github.com/savannahghi/feedlib"
)

// Repository is the interface to be implemented by the database(s)
// The method signatures are should be database independent
type Repository interface {
	// getting a feed...create a default feed if it does not exist
	// return: feed, matching count, total count, optional error
	GetFeed(
		ctx context.Context,
		uid *string,
		isAnonymous *bool,
		flavour feedlib.Flavour,
		persistent feedlib.BooleanFilter,
		status *feedlib.Status,
		visibility *feedlib.Visibility,
		expired *feedlib.BooleanFilter,
		filterParams *helpers.FilterParams,
	) (*domain.Feed, error)

	// getting a the LATEST VERSION of a feed item from a feed
	GetFeedItem(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
		itemID string,
	) (*feedlib.Item, error)

	// saving a new feed item
	SaveFeedItem(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
		item *feedlib.Item,
	) (*feedlib.Item, error)

	// updating an existing feed item
	UpdateFeedItem(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
		item *feedlib.Item,
	) (*feedlib.Item, error)

	// DeleteFeedItem permanently deletes a feed item and it's copies
	DeleteFeedItem(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
		itemID string,
	) error

	// getting THE LATEST VERSION OF a nudge from a feed
	GetNudge(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
		nudgeID string,
	) (*feedlib.Nudge, error)

	// saving a new modified nudge
	SaveNudge(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
		nudge *feedlib.Nudge,
	) (*feedlib.Nudge, error)

	// updating an existing nudge
	UpdateNudge(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
		nudge *feedlib.Nudge,
	) (*feedlib.Nudge, error)

	// DeleteNudge permanently deletes a nudge and it's copies
	DeleteNudge(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
		nudgeID string,
	) error

	// getting THE LATEST VERSION OF a single action
	GetAction(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
		actionID string,
	) (*feedlib.Action, error)

	// saving a new action
	SaveAction(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
		action *feedlib.Action,
	) (*feedlib.Action, error)

	// DeleteAction permanently deletes an action and it's copies
	DeleteAction(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
		actionID string,
	) error

	// PostMessage posts a message or a reply to a message/thread
	PostMessage(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
		itemID string,
		message *feedlib.Message,
	) (*feedlib.Message, error)

	// GetMessage retrieves THE LATEST VERSION OF a message
	GetMessage(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
		itemID string,
		messageID string,
	) (*feedlib.Message, error)

	// DeleteMessage deletes a message
	DeleteMessage(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
		itemID string,
		messageID string,
	) error

	// GetMessages retrieves a message
	GetMessages(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
		itemID string,
	) ([]feedlib.Message, error)

	SaveIncomingEvent(
		ctx context.Context,
		event *feedlib.Event,
	) error

	SaveOutgoingEvent(
		ctx context.Context,
		event *feedlib.Event,
	) error

	GetNudges(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
		status *feedlib.Status,
		visibility *feedlib.Visibility,
		expired *feedlib.BooleanFilter,
	) ([]feedlib.Nudge, error)

	GetActions(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
	) ([]feedlib.Action, error)

	GetItems(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
		persistent feedlib.BooleanFilter,
		status *feedlib.Status,
		visibility *feedlib.Visibility,
		expired *feedlib.BooleanFilter,
		filterParams *helpers.FilterParams,
	) ([]feedlib.Item, error)

	Labels(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
	) ([]string, error)

	SaveLabel(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
		label string,
	) error

	UnreadPersistentItems(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
	) (int, error)

	UpdateUnreadPersistentItemsCount(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
	) error

	GetDefaultNudgeByTitle(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
		title string,
	) (*feedlib.Nudge, error)

	SaveMarketingMessage(
		ctx context.Context,
		data dto.MarketingSMS,
	) (*dto.MarketingSMS, error)

	GetMarketingSMSByID(
		ctx context.Context,
		id string,
	) (*dto.MarketingSMS, error)

	GetMarketingSMSByPhone(
		ctx context.Context,
		phoneNumber string,
	) (*dto.MarketingSMS, error)

	UpdateMarketingMessage(
		ctx context.Context,
		data *dto.MarketingSMS,
	) (*dto.MarketingSMS, error)

	SaveTwilioResponse(
		ctx context.Context,
		data dto.Message,
	) error

	SaveNotification(
		ctx context.Context,
		firestoreClient *firestore.Client,
		notification dto.SavedNotification,
	) error

	RetrieveNotification(
		ctx context.Context,
		firestoreClient *firestore.Client,
		registrationToken string,
		newerThan time.Time,
		limit int,
	) ([]*dto.SavedNotification, error)

	SaveNPSResponse(
		ctx context.Context,
		response *dto.NPSResponse,
	) error

	SaveOutgoingEmails(
		ctx context.Context,
		payload *dto.OutgoingEmailsLog,
	) error
	UpdateMailgunDeliveryStatus(
		ctx context.Context,
		payload *dto.MailgunEvent,
	) (*dto.OutgoingEmailsLog, error)

	SaveTwilioVideoCallbackStatus(
		ctx context.Context,
		data dto.CallbackData,
	) error
}

// DbService is an implementation of the database repository
// It is implementation agnostic i.e logic should be handled using
// the preferred database
type DbService struct {
	firestore *fb.Repository
}

// NewDbService creates a new database service
func NewDbService() Repository {
	ctx := context.Background()

	firestore, err := fb.NewFirebaseRepository(ctx)
	if err != nil {
		log.Panicf("can't instantiate firebase repository in resolver: %v", err)
	}
	return &DbService{
		firestore: firestore,
	}
}

// CheckPreconditions ensures correct initialization
func (d DbService) CheckPreconditions() {
	if d.firestore == nil {
		log.Panicf("nil firestore service in database service")
	}
}

// GetFeed ...
func (d *DbService) GetFeed(
	ctx context.Context,
	uid *string,
	isAnonymous *bool,
	flavour feedlib.Flavour,
	persistent feedlib.BooleanFilter,
	status *feedlib.Status,
	visibility *feedlib.Visibility,
	expired *feedlib.BooleanFilter,
	filterParams *helpers.FilterParams,
) (*domain.Feed, error) {
	return d.firestore.GetFeed(ctx, uid, isAnonymous, flavour, persistent, status, visibility, expired, filterParams)
}

// GetFeedItem ...
func (d *DbService) GetFeedItem(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	itemID string,
) (*feedlib.Item, error) {
	return d.firestore.GetFeedItem(ctx, uid, flavour, itemID)
}

// SaveFeedItem ...
func (d *DbService) SaveFeedItem(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	item *feedlib.Item,
) (*feedlib.Item, error) {
	return d.firestore.SaveFeedItem(ctx, uid, flavour, item)
}

// UpdateFeedItem ...
func (d *DbService) UpdateFeedItem(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	item *feedlib.Item,
) (*feedlib.Item, error) {
	return d.firestore.UpdateFeedItem(ctx, uid, flavour, item)
}

// DeleteFeedItem permanently deletes a feed item and it's copies
func (d *DbService) DeleteFeedItem(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	itemID string,
) error {
	return d.firestore.DeleteFeedItem(ctx, uid, flavour, itemID)
}

// GetNudge gets THE LATEST VERSION OF a nudge from a feed
func (d *DbService) GetNudge(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	nudgeID string,
) (*feedlib.Nudge, error) {
	return d.firestore.GetNudge(ctx, uid, flavour, nudgeID)
}

// SaveNudge saves a new modified nudge
func (d *DbService) SaveNudge(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	nudge *feedlib.Nudge,
) (*feedlib.Nudge, error) {
	return d.firestore.SaveNudge(ctx, uid, flavour, nudge)
}

// UpdateNudge updates an existing nudge
func (d *DbService) UpdateNudge(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	nudge *feedlib.Nudge,
) (*feedlib.Nudge, error) {
	return d.firestore.UpdateNudge(ctx, uid, flavour, nudge)
}

// DeleteNudge permanently deletes a nudge and it's copies
func (d *DbService) DeleteNudge(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	nudgeID string,
) error {
	return d.firestore.DeleteNudge(ctx, uid, flavour, nudgeID)
}

// GetAction gets THE LATEST VERSION OF a single action
func (d *DbService) GetAction(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	actionID string,
) (*feedlib.Action, error) {
	return d.firestore.GetAction(ctx, uid, flavour, actionID)
}

// SaveAction saves a new action
func (d *DbService) SaveAction(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	action *feedlib.Action,
) (*feedlib.Action, error) {
	return d.firestore.SaveAction(ctx, uid, flavour, action)
}

// DeleteAction permanently deletes an action and it's copies
func (d *DbService) DeleteAction(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	actionID string,
) error {
	return d.firestore.DeleteAction(ctx, uid, flavour, actionID)
}

// PostMessage posts a message or a reply to a message/thread
func (d *DbService) PostMessage(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	itemID string,
	message *feedlib.Message,
) (*feedlib.Message, error) {
	return d.firestore.PostMessage(ctx, uid, flavour, itemID, message)
}

// GetMessage retrieves THE LATEST VERSION OF a message
func (d *DbService) GetMessage(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	itemID string,
	messageID string,
) (*feedlib.Message, error) {
	return d.firestore.GetMessage(ctx, uid, flavour, itemID, messageID)
}

// DeleteMessage deletes a message
func (d *DbService) DeleteMessage(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	itemID string,
	messageID string,
) error {
	return d.firestore.DeleteMessage(ctx, uid, flavour, itemID, messageID)
}

// GetMessages retrieves a message
func (d *DbService) GetMessages(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	itemID string,
) ([]feedlib.Message, error) {
	return d.firestore.GetMessages(ctx, uid, flavour, itemID)
}

// SaveIncomingEvent ...
func (d *DbService) SaveIncomingEvent(
	ctx context.Context,
	event *feedlib.Event,
) error {
	return d.firestore.SaveIncomingEvent(ctx, event)
}

// SaveOutgoingEvent ...
func (d *DbService) SaveOutgoingEvent(
	ctx context.Context,
	event *feedlib.Event,
) error {
	return d.firestore.SaveOutgoingEvent(ctx, event)
}

// GetNudges ...
func (d *DbService) GetNudges(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	status *feedlib.Status,
	visibility *feedlib.Visibility,
	expired *feedlib.BooleanFilter,
) ([]feedlib.Nudge, error) {
	return d.firestore.GetNudges(ctx, uid, flavour, status, visibility, expired)
}

// GetActions ...
func (d *DbService) GetActions(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
) ([]feedlib.Action, error) {
	return d.firestore.GetActions(ctx, uid, flavour)
}

// GetItems ...
func (d *DbService) GetItems(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	persistent feedlib.BooleanFilter,
	status *feedlib.Status,
	visibility *feedlib.Visibility,
	expired *feedlib.BooleanFilter,
	filterParams *helpers.FilterParams,
) ([]feedlib.Item, error) {
	return d.firestore.GetItems(ctx, uid, flavour, persistent, status, visibility, expired, filterParams)
}

// Labels ...
func (d *DbService) Labels(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
) ([]string, error) {
	return d.firestore.Labels(ctx, uid, flavour)
}

// SaveLabel ...
func (d *DbService) SaveLabel(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	label string,
) error {
	return d.firestore.SaveLabel(ctx, uid, flavour, label)
}

// UnreadPersistentItems ...
func (d *DbService) UnreadPersistentItems(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
) (int, error) {
	return d.firestore.UnreadPersistentItems(ctx, uid, flavour)
}

// UpdateUnreadPersistentItemsCount ...
func (d *DbService) UpdateUnreadPersistentItemsCount(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
) error {
	return d.firestore.UpdateUnreadPersistentItemsCount(ctx, uid, flavour)
}

// GetDefaultNudgeByTitle ...
func (d *DbService) GetDefaultNudgeByTitle(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	title string,
) (*feedlib.Nudge, error) {
	return d.firestore.GetDefaultNudgeByTitle(ctx, uid, flavour, title)
}

// SaveMarketingMessage saves the callback data for future analysis
func (d *DbService) SaveMarketingMessage(
	ctx context.Context,
	data dto.MarketingSMS,
) (*dto.MarketingSMS, error) {
	return d.firestore.SaveMarketingMessage(ctx, data)
}

// SaveTwilioResponse saves the callback data for future analysis
func (d *DbService) SaveTwilioResponse(
	ctx context.Context,
	data dto.Message,
) error {
	return d.firestore.SaveTwilioResponse(ctx, data)
}

// SaveNotification saves a notification
func (d *DbService) SaveNotification(
	ctx context.Context,
	firestoreClient *firestore.Client,
	notification dto.SavedNotification,
) error {
	return d.firestore.SaveNotification(ctx, firestoreClient, notification)
}

// RetrieveNotification retrieves a notification
func (d *DbService) RetrieveNotification(
	ctx context.Context,
	firestoreClient *firestore.Client,
	registrationToken string,
	newerThan time.Time,
	limit int,
) ([]*dto.SavedNotification, error) {
	return d.firestore.RetrieveNotification(ctx, firestoreClient, registrationToken, newerThan, limit)
}

// SaveNPSResponse saves a NPS response
func (d *DbService) SaveNPSResponse(
	ctx context.Context,
	response *dto.NPSResponse,
) error {
	return d.firestore.SaveNPSResponse(ctx, response)
}

// UpdateMarketingMessage ..
func (d *DbService) UpdateMarketingMessage(
	ctx context.Context,
	data *dto.MarketingSMS,
) (*dto.MarketingSMS, error) {
	return d.firestore.UpdateMarketingMessage(ctx, data)
}

// SaveOutgoingEmails ...
func (d *DbService) SaveOutgoingEmails(ctx context.Context, payload *dto.OutgoingEmailsLog) error {
	return d.firestore.SaveOutgoingEmails(ctx, payload)
}

// UpdateMailgunDeliveryStatus ...
func (d *DbService) UpdateMailgunDeliveryStatus(ctx context.Context, payload *dto.MailgunEvent) (*dto.OutgoingEmailsLog, error) {
	return d.firestore.UpdateMailgunDeliveryStatus(ctx, payload)
}

// GetMarketingSMSByPhone ..
func (d *DbService) GetMarketingSMSByPhone(ctx context.Context, phoneNumber string) (*dto.MarketingSMS, error) {
	return d.firestore.GetMarketingSMSByPhone(ctx, phoneNumber)
}

// GetMarketingSMSByID ..
func (d *DbService) GetMarketingSMSByID(
	ctx context.Context,
	id string,
) (*dto.MarketingSMS, error) {
	return d.firestore.GetMarketingSMSByID(ctx, id)
}

// SaveTwilioVideoCallbackStatus ..
func (d *DbService) SaveTwilioVideoCallbackStatus(
	ctx context.Context,
	data dto.CallbackData,
) error {
	return d.firestore.SaveTwilioVideoCallbackStatus(ctx, data)
}
