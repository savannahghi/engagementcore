package mock

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/savannahghi/engagementcore/pkg/engagement/application/common/dto"
	"github.com/savannahghi/engagementcore/pkg/engagement/application/common/helpers"
	"github.com/savannahghi/engagementcore/pkg/engagement/domain"
	"github.com/savannahghi/engagementcore/pkg/engagement/infrastructure/services/onboarding"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/profileutils"
)

// FakeInfrastructure mocks the infrastructure interface
type FakeInfrastructure struct {
	GetFeedFn func(
		ctx context.Context,
		uid *string,
		isAnonymous *bool,
		flavour feedlib.Flavour,
		playMP4 bool,
		persistent feedlib.BooleanFilter,
		status *feedlib.Status,
		visibility *feedlib.Visibility,
		expired *feedlib.BooleanFilter,
		filterParams *helpers.FilterParams,
	) (*domain.Feed, error)

	// getting a the LATEST VERSION of a feed item from a feed
	GetFeedItemFn func(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
		itemID string,
	) (*feedlib.Item, error)

	// saving a new feed item
	SaveFeedItemFn func(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
		item *feedlib.Item,
	) (*feedlib.Item, error)

	// updating an existing feed item
	UpdateFeedItemFn func(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
		item *feedlib.Item,
	) (*feedlib.Item, error)

	// DeleteFeedItem permanently deletes a feed item and it's copies
	DeleteFeedItemFn func(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
		itemID string,
	) error

	// getting THE LATEST VERSION OF a nudge from a feed
	GetNudgeFn func(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
		nudgeID string,
	) (*feedlib.Nudge, error)

	// saving a new modified nudge
	SaveNudgeFn func(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
		nudge *feedlib.Nudge,
	) (*feedlib.Nudge, error)

	// updating an existing nudge
	UpdateNudgeFn func(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
		nudge *feedlib.Nudge,
	) (*feedlib.Nudge, error)

	// DeleteNudge permanently deletes a nudge and it's copies
	DeleteNudgeFn func(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
		nudgeID string,
	) error

	// getting THE LATEST VERSION OF a single action
	GetActionFn func(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
		actionID string,
	) (*feedlib.Action, error)

	// saving a new action
	SaveActionFn func(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
		action *feedlib.Action,
	) (*feedlib.Action, error)

	// DeleteAction permanently deletes an action and it's copies
	DeleteActionFn func(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
		actionID string,
	) error

	// PostMessage posts a message or a reply to a message/thread
	PostMessageFn func(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
		itemID string,
		message *feedlib.Message,
	) (*feedlib.Message, error)

	// GetMessage retrieves THE LATEST VERSION OF a message
	GetMessageFn func(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
		itemID string,
		messageID string,
	) (*feedlib.Message, error)

	// DeleteMessage deletes a message
	DeleteMessageFn func(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
		itemID string,
		messageID string,
	) error

	// GetMessages retrieves a message
	GetMessagesFn func(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
		itemID string,
	) ([]feedlib.Message, error)

	SaveIncomingEventFn func(
		ctx context.Context,
		event *feedlib.Event,
	) error

	SaveOutgoingEventFn func(
		ctx context.Context,
		event *feedlib.Event,
	) error

	GetNudgesFn func(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
		status *feedlib.Status,
		visibility *feedlib.Visibility,
		expired *feedlib.BooleanFilter,
	) ([]feedlib.Nudge, error)

	GetActionsFn func(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
	) ([]feedlib.Action, error)

	GetItemsFn func(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
		persistent feedlib.BooleanFilter,
		status *feedlib.Status,
		visibility *feedlib.Visibility,
		expired *feedlib.BooleanFilter,
		filterParams *helpers.FilterParams,
	) ([]feedlib.Item, error)

	LabelsFn func(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
	) ([]string, error)

	SaveLabelFn func(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
		label string,
	) error

	UnreadPersistentItemsFn func(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
	) (int, error)

	UpdateUnreadPersistentItemsCountFn func(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
	) error

	GetDefaultNudgeByTitleFn func(
		ctx context.Context,
		uid string,
		flavour feedlib.Flavour,
		title string,
	) (*feedlib.Nudge, error)

	SaveTwilioResponseFn func(
		ctx context.Context,
		data dto.Message,
	) error

	SaveNotificationFn func(
		ctx context.Context,
		firestoreClient *firestore.Client,
		notification dto.SavedNotification,
	) error

	RetrieveNotificationFn func(
		ctx context.Context,
		firestoreClient *firestore.Client,
		registrationToken string,
		newerThan time.Time,
		limit int,
	) ([]*dto.SavedNotification, error)

	SaveNPSResponseFn func(
		ctx context.Context,
		response *dto.NPSResponse,
	) error

	// save Patients Feedback
	SavePatientFeedbackFn func(
		ctx context.Context,
		response *dto.PatientFeedbackResponse,
	) error

	SaveTwilioVideoCallbackStatusFn func(
		ctx context.Context,
		data dto.CallbackData,
	) error

	SendNotificationFn func(
		ctx context.Context,
		registrationTokens []string,
		data map[string]string,
		notification *firebasetools.FirebaseSimpleNotificationInput,
		android *firebasetools.FirebaseAndroidConfigInput,
		ios *firebasetools.FirebaseAPNSConfigInput,
		web *firebasetools.FirebaseWebpushConfigInput,
	) (bool, error)

	NotificationsFn func(
		ctx context.Context,
		registrationToken string,
		newerThan time.Time,
		limit int,
	) ([]*dto.SavedNotification, error)

	SendFCMByPhoneOrEmailFn func(
		ctx context.Context,
		phoneNumber *string,
		email *string,
		data map[string]interface{},
		notification firebasetools.FirebaseSimpleNotificationInput,
		android *firebasetools.FirebaseAndroidConfigInput,
		ios *firebasetools.FirebaseAPNSConfigInput,
		web *firebasetools.FirebaseWebpushConfigInput,
	) (bool, error)

	GetFeedContentFn    func(ctx context.Context, flavour feedlib.Flavour) ([]*domain.GhostCMSPost, error)
	GetFaqsContentFn    func(ctx context.Context, flavour feedlib.Flavour) ([]*domain.GhostCMSPost, error)
	GetLibraryContentFn func(ctx context.Context) ([]*domain.GhostCMSPost, error)

	SendInBlueFn  func(ctx context.Context, subject, text string, to ...string) (string, string, error)
	SendMailgunFn func(
		ctx context.Context,
		subject, text string,
		body *string,
		to ...string,
	) (string, string, error)
	SendEmailFn func(
		ctx context.Context,
		subject, text string,
		body *string,
		to ...string,
	) (string, string, error)
	SimpleEmailFn func(
		ctx context.Context,
		subject, text string,
		body *string,
		to ...string,
	) (string, error)
	SaveOutgoingEmailsFn          func(ctx context.Context, payload *dto.OutgoingEmailsLog) error
	UpdateMailgunDeliveryStatusFn func(
		ctx context.Context,
		payload *dto.MailgunEvent,
	) (*dto.OutgoingEmailsLog, error)
	GenerateEmailTemplateFn func(name string, templateName string) string

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

	GetEmailAddressesFn            func(ctx context.Context, uids onboarding.UserUIDs) (map[string][]string, error)
	GetPhoneNumbersFn              func(ctx context.Context, uids onboarding.UserUIDs) (map[string][]string, error)
	GetDeviceTokensFn              func(ctx context.Context, uid onboarding.UserUIDs) (map[string][]string, error)
	GetUserProfileFn               func(ctx context.Context, uid string) (*profileutils.UserProfile, error)
	GetUserProfileByPhoneOrEmailFn func(ctx context.Context, payload *dto.RetrieveUserProfileInput) (*profileutils.UserProfile, error)

	GenerateAndSendOTPFn   func(ctx context.Context, msisdn string, appID *string) (string, error)
	SendOTPToEmailFn       func(ctx context.Context, msisdn, email *string, appID *string) (string, error)
	SaveOTPToFirestoreFn   func(otp dto.OTP) error
	VerifyOtpFn            func(ctx context.Context, msisdn, verificationCode *string) (bool, error)
	VerifyEmailOtpFn       func(ctx context.Context, email, verificationCode *string) (bool, error)
	GenerateRetryOTPFn     func(ctx context.Context, msisdn *string, retryStep int, appID *string) (string, error)
	EmailVerificationOtpFn func(ctx context.Context, email *string) (string, error)
	GenerateOTPFn          func(ctx context.Context) (string, error)

	SendToManyFn func(
		ctx context.Context,
		message string,
		to []string,
		from enumutils.SenderID,
	) (*dto.SendMessageResponse, error)

	SendFn func(
		ctx context.Context,
		to, message string,
		from enumutils.SenderID,
	) (*dto.SendMessageResponse, error)

	RecordNPSResponseFn func(ctx context.Context, input dto.NPSInput) (bool, error)

	RecordPatientFeedbackFn func(ctx context.Context, input dto.PatientFeedbackInput) (bool, error)

	RoomFn func(ctx context.Context) (*dto.Room, error)

	TwilioAccessTokenFn func(ctx context.Context) (*dto.AccessToken, error)

	SendSMSFn func(ctx context.Context, to string, msg string) error

	UploadFn func(
		ctx context.Context,
		inp profileutils.UploadInput,
	) (*profileutils.Upload, error)

	FindUploadByIDFn func(
		ctx context.Context,
		id string,
	) (*profileutils.Upload, error)

	PhoneNumberVerificationCodeFn func(
		ctx context.Context,
		to string,
		code string,
		marketingMessage string,
	) (bool, error)

	WellnessCardActivationDependantFn func(
		ctx context.Context,
		to string,
		memberName string,
		cardName string,
		marketingMessage string,
	) (bool, error)

	WellnessCardActivationPrincipalFn func(
		ctx context.Context,
		to string,
		memberName string,
		cardName string,
		minorAgeThreshold string,
		marketingMessage string,
	) (bool, error)

	BillNotificationFn func(
		ctx context.Context,
		to string,
		productName string,
		billingPeriod string,
		billAmount string,
		paymentInstruction string,
		marketingMessage string,
	) (bool, error)

	VirtualCardsFn func(
		ctx context.Context,
		to string,
		wellnessCardFamily string,
		virtualCardLink string,
		marketingMessage string,
	) (bool, error)

	VisitStartFn func(
		ctx context.Context,
		to string,
		memberName string,
		benefitName string,
		locationName string,
		startTime string,
		balance string,
		marketingMessage string,
	) (bool, error)

	ClaimNotificationFn func(
		ctx context.Context,
		to string,
		claimReference string,
		claimTypeParenthesized string,
		provider string,
		visitType string,
		claimTime string,
		marketingMessage string,
	) (bool, error)

	PreauthApprovalFn func(
		ctx context.Context,
		to string,
		currency string,
		amount string,
		benefit string,
		provider string,
		member string,
		careContact string,
		marketingMessage string,
	) (bool, error)

	PreauthRequestFn func(
		ctx context.Context,
		to string,
		currency string,
		amount string,
		benefit string,
		provider string,
		requestTime string,
		member string,
		careContact string,
		marketingMessage string,
	) (bool, error)

	SladeOTPFn func(
		ctx context.Context,
		to string,
		name string,
		otp string,
		marketingMessage string,
	) (bool, error)
}

// GetFeed ...
func (f *FakeInfrastructure) GetFeed(
	ctx context.Context,
	uid *string,
	isAnonymous *bool,
	flavour feedlib.Flavour,
	playMP4 bool,
	persistent feedlib.BooleanFilter,
	status *feedlib.Status,
	visibility *feedlib.Visibility,
	expired *feedlib.BooleanFilter,
	filterParams *helpers.FilterParams,
) (*domain.Feed, error) {
	return f.GetFeedFn(ctx, uid, isAnonymous, flavour, playMP4, persistent, status, visibility, expired, filterParams)
}

// GetFeedItem ...
func (f *FakeInfrastructure) GetFeedItem(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	itemID string,
) (*feedlib.Item, error) {
	return f.GetFeedItemFn(ctx, uid, flavour, itemID)
}

// SaveFeedItem ...
func (f *FakeInfrastructure) SaveFeedItem(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	item *feedlib.Item,
) (*feedlib.Item, error) {
	return f.SaveFeedItemFn(ctx, uid, flavour, item)
}

// UpdateFeedItem ...
func (f *FakeInfrastructure) UpdateFeedItem(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	item *feedlib.Item,
) (*feedlib.Item, error) {
	return f.UpdateFeedItemFn(ctx, uid, flavour, item)
}

// DeleteFeedItem permanently deletes a feed item and it's copies
func (f *FakeInfrastructure) DeleteFeedItem(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	itemID string,
) error {
	return f.DeleteFeedItemFn(ctx, uid, flavour, itemID)
}

// GetNudge gets THE LATEST VERSION OF a nudge from a feed
func (f *FakeInfrastructure) GetNudge(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	nudgeID string,
) (*feedlib.Nudge, error) {
	return f.GetNudgeFn(ctx, uid, flavour, nudgeID)
}

// SaveNudge saves a new modified nudge
func (f *FakeInfrastructure) SaveNudge(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	nudge *feedlib.Nudge,
) (*feedlib.Nudge, error) {
	return f.SaveNudgeFn(ctx, uid, flavour, nudge)
}

// UpdateNudge updates an existing nudge
func (f *FakeInfrastructure) UpdateNudge(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	nudge *feedlib.Nudge,
) (*feedlib.Nudge, error) {
	return f.UpdateNudgeFn(ctx, uid, flavour, nudge)
}

// DeleteNudge permanently deletes a nudge and it's copies
func (f *FakeInfrastructure) DeleteNudge(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	nudgeID string,
) error {
	return f.DeleteNudgeFn(ctx, uid, flavour, nudgeID)
}

// GetAction gets THE LATEST VERSION OF a single action
func (f *FakeInfrastructure) GetAction(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	actionID string,
) (*feedlib.Action, error) {
	return f.GetActionFn(ctx, uid, flavour, actionID)
}

// SaveAction saves a new action
func (f *FakeInfrastructure) SaveAction(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	action *feedlib.Action,
) (*feedlib.Action, error) {
	return f.SaveActionFn(ctx, uid, flavour, action)
}

// DeleteAction permanently deletes an action and it's copies
func (f *FakeInfrastructure) DeleteAction(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	actionID string,
) error {
	return f.DeleteActionFn(ctx, uid, flavour, actionID)
}

// PostMessage posts a message or a reply to a message/thread
func (f *FakeInfrastructure) PostMessage(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	itemID string,
	message *feedlib.Message,
) (*feedlib.Message, error) {
	return f.PostMessageFn(ctx, uid, flavour, itemID, message)
}

// GetMessage retrieves THE LATEST VERSION OF a message
func (f *FakeInfrastructure) GetMessage(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	itemID string,
	messageID string,
) (*feedlib.Message, error) {
	return f.GetMessageFn(ctx, uid, flavour, itemID, messageID)
}

// DeleteMessage deletes a message
func (f *FakeInfrastructure) DeleteMessage(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	itemID string,
	messageID string,
) error {
	return f.DeleteMessageFn(ctx, uid, flavour, itemID, messageID)
}

// GetMessages retrieves a message
func (f *FakeInfrastructure) GetMessages(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	itemID string,
) ([]feedlib.Message, error) {
	return f.GetMessagesFn(ctx, uid, flavour, itemID)
}

// SaveIncomingEvent ...
func (f *FakeInfrastructure) SaveIncomingEvent(
	ctx context.Context,
	event *feedlib.Event,
) error {
	return f.SaveIncomingEventFn(ctx, event)
}

// SaveOutgoingEvent ...
func (f *FakeInfrastructure) SaveOutgoingEvent(
	ctx context.Context,
	event *feedlib.Event,
) error {
	return f.SaveOutgoingEventFn(ctx, event)
}

// GetNudges ...
func (f *FakeInfrastructure) GetNudges(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	status *feedlib.Status,
	visibility *feedlib.Visibility,
	expired *feedlib.BooleanFilter,
) ([]feedlib.Nudge, error) {
	return f.GetNudgesFn(ctx, uid, flavour, status, visibility, expired)
}

// GetActions ...
func (f *FakeInfrastructure) GetActions(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
) ([]feedlib.Action, error) {
	return f.GetActionsFn(ctx, uid, flavour)
}

// GetItems ...
func (f *FakeInfrastructure) GetItems(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	persistent feedlib.BooleanFilter,
	status *feedlib.Status,
	visibility *feedlib.Visibility,
	expired *feedlib.BooleanFilter,
	filterParams *helpers.FilterParams,
) ([]feedlib.Item, error) {
	return f.GetItemsFn(ctx, uid, flavour, persistent, status, visibility, expired, filterParams)
}

// Labels ...
func (f *FakeInfrastructure) Labels(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
) ([]string, error) {
	return f.LabelsFn(ctx, uid, flavour)
}

// SaveLabel ...
func (f *FakeInfrastructure) SaveLabel(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	label string,
) error {
	return f.SaveLabelFn(ctx, uid, flavour, label)
}

// UnreadPersistentItems ...
func (f *FakeInfrastructure) UnreadPersistentItems(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
) (int, error) {
	return f.UnreadPersistentItemsFn(ctx, uid, flavour)
}

// UpdateUnreadPersistentItemsCount ...
func (f *FakeInfrastructure) UpdateUnreadPersistentItemsCount(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
) error {
	return f.UpdateUnreadPersistentItemsCountFn(ctx, uid, flavour)
}

// GetDefaultNudgeByTitle ...
func (f *FakeInfrastructure) GetDefaultNudgeByTitle(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	title string,
) (*feedlib.Nudge, error) {
	return f.GetDefaultNudgeByTitleFn(ctx, uid, flavour, title)
}

// SaveTwilioResponse saves the callback data for future analysis
func (f *FakeInfrastructure) SaveTwilioResponse(
	ctx context.Context,
	data dto.Message,
) error {
	return f.SaveTwilioResponseFn(ctx, data)
}

// SaveNotification saves a notification
func (f *FakeInfrastructure) SaveNotification(
	ctx context.Context,
	firestoreClient *firestore.Client,
	notification dto.SavedNotification,
) error {
	return f.SaveNotificationFn(ctx, firestoreClient, notification)
}

// RetrieveNotification retrieves a notification
func (f *FakeInfrastructure) RetrieveNotification(
	ctx context.Context,
	firestoreClient *firestore.Client,
	registrationToken string,
	newerThan time.Time,
	limit int,
) ([]*dto.SavedNotification, error) {
	return f.RetrieveNotificationFn(ctx, firestoreClient, registrationToken, newerThan, limit)
}

// SaveNPSResponse saves a NPS response
func (f *FakeInfrastructure) SaveNPSResponse(
	ctx context.Context,
	response *dto.NPSResponse,
) error {
	return f.SaveNPSResponseFn(ctx, response)
}

// SavePatientFeedback  --> TODO: Will work on this tomorrow
// func (f *FakeInfrastructure) SavePatientFeedback(
// 	ctx context.Context,
// 	response *dto.PatientFeedbackResponse,
// ) error {
// 	return f.SavePatientFeedback(ctx, response)
// }

// SaveTwilioVideoCallbackStatus ..
func (f *FakeInfrastructure) SaveTwilioVideoCallbackStatus(
	ctx context.Context,
	data dto.CallbackData,
) error {
	return f.SaveTwilioVideoCallbackStatusFn(ctx, data)
}

// SendInBlue ...
func (f *FakeInfrastructure) SendInBlue(ctx context.Context, subject, text string, to ...string) (string, string, error) {
	return f.SendInBlueFn(ctx, subject, text, to...)
}

// SendMailgun ...
func (f *FakeInfrastructure) SendMailgun(
	ctx context.Context,
	subject, text string,
	body *string,
	to ...string,
) (string, string, error) {
	return f.SendMailgunFn(ctx, subject, text, body, to...)
}

// SendEmail ...
func (f *FakeInfrastructure) SendEmail(
	ctx context.Context,
	subject, text string,
	body *string,
	to ...string,
) (string, string, error) {
	return f.SendEmailFn(ctx, subject, text, body, to...)
}

// SimpleEmail ...
func (f *FakeInfrastructure) SimpleEmail(
	ctx context.Context,
	subject, text string,
	body *string,
	to ...string,
) (string, error) {
	return f.SimpleEmailFn(ctx, subject, text, body, to...)
}

// SaveOutgoingEmails ...
func (f *FakeInfrastructure) SaveOutgoingEmails(ctx context.Context, payload *dto.OutgoingEmailsLog) error {
	return f.SaveOutgoingEmailsFn(ctx, payload)
}

// UpdateMailgunDeliveryStatus ...
func (f *FakeInfrastructure) UpdateMailgunDeliveryStatus(
	ctx context.Context,
	payload *dto.MailgunEvent,
) (*dto.OutgoingEmailsLog, error) {
	return f.UpdateMailgunDeliveryStatusFn(ctx, payload)
}

// GenerateEmailTemplate ...
func (f *FakeInfrastructure) GenerateEmailTemplate(name string, templateName string) string {
	return f.GenerateEmailTemplateFn(name, templateName)
}

// Notify Sends a message to a topic
func (f *FakeInfrastructure) Notify(
	ctx context.Context,
	topicID string,
	uid string,
	flavour feedlib.Flavour,
	payload feedlib.Element,
	metadata map[string]interface{},
) error {
	return f.NotifyFn(ctx, topicID, uid, flavour, payload, metadata)
}

// TopicIDs ...
func (f *FakeInfrastructure) TopicIDs() []string {
	return f.TopicIDsFn()
}

// SubscriptionIDs ...
func (f *FakeInfrastructure) SubscriptionIDs() map[string]string {
	return f.SubscriptionIDsFn()
}

// ReverseSubscriptionIDs ...
func (f *FakeInfrastructure) ReverseSubscriptionIDs() map[string]string {
	return f.ReverseSubscriptionIDsFn()
}

// GetEmailAddresses ...
func (f *FakeInfrastructure) GetEmailAddresses(ctx context.Context, uids onboarding.UserUIDs) (map[string][]string, error) {
	return f.GetEmailAddressesFn(ctx, uids)
}

// GetPhoneNumbers ...
func (f *FakeInfrastructure) GetPhoneNumbers(ctx context.Context, uids onboarding.UserUIDs) (map[string][]string, error) {
	return f.GetPhoneNumbersFn(ctx, uids)
}

// GetDeviceTokens ...
func (f *FakeInfrastructure) GetDeviceTokens(ctx context.Context, uid onboarding.UserUIDs) (map[string][]string, error) {
	return f.GetDeviceTokensFn(ctx, uid)
}

// GetUserProfile ...
func (f *FakeInfrastructure) GetUserProfile(ctx context.Context, uid string) (*profileutils.UserProfile, error) {
	return f.GetUserProfileFn(ctx, uid)
}

// GetUserProfileByPhoneOrEmail ...
func (f *FakeInfrastructure) GetUserProfileByPhoneOrEmail(ctx context.Context, payload *dto.RetrieveUserProfileInput) (*profileutils.UserProfile, error) {
	return f.GetUserProfileByPhoneOrEmailFn(ctx, payload)
}

// GenerateAndSendOTP ...
func (f *FakeInfrastructure) GenerateAndSendOTP(ctx context.Context, msisdn string, appID *string) (string, error) {
	return f.GenerateAndSendOTPFn(ctx, msisdn, appID)
}

// SendOTPToEmail ...
func (f *FakeInfrastructure) SendOTPToEmail(ctx context.Context, msisdn, email *string, appID *string) (string, error) {
	return f.SendOTPToEmailFn(ctx, msisdn, email, appID)
}

// SaveOTPToFirestore ...
func (f *FakeInfrastructure) SaveOTPToFirestore(otp dto.OTP) error {
	return f.SaveOTPToFirestoreFn(otp)
}

// VerifyOtp ...
func (f *FakeInfrastructure) VerifyOtp(ctx context.Context, msisdn, verificationCode *string) (bool, error) {
	return f.VerifyOtpFn(ctx, msisdn, verificationCode)
}

// VerifyEmailOtp ...
func (f *FakeInfrastructure) VerifyEmailOtp(ctx context.Context, email, verificationCode *string) (bool, error) {
	return f.VerifyEmailOtpFn(ctx, email, verificationCode)
}

// GenerateRetryOTP ...
func (f *FakeInfrastructure) GenerateRetryOTP(ctx context.Context, msisdn *string, retryStep int, appID *string) (string, error) {
	return f.GenerateRetryOTPFn(ctx, msisdn, retryStep, appID)
}

// EmailVerificationOtp ...
func (f *FakeInfrastructure) EmailVerificationOtp(ctx context.Context, email *string) (string, error) {
	return f.EmailVerificationOtpFn(ctx, email)
}

// GenerateOTP ...
func (f *FakeInfrastructure) GenerateOTP(ctx context.Context) (string, error) {
	return f.GenerateOTPFn(ctx)
}

// SendToMany ...
func (f *FakeInfrastructure) SendToMany(
	ctx context.Context,
	message string,
	to []string,
	from enumutils.SenderID,
) (*dto.SendMessageResponse, error) {
	return f.SendToManyFn(ctx, message, to, from)
}

// Send ...
func (f *FakeInfrastructure) Send(
	ctx context.Context,
	to, message string,
	from enumutils.SenderID,
) (*dto.SendMessageResponse, error) {
	return f.SendFn(ctx, to, message, from)
}

// RecordNPSResponse ...
func (f *FakeInfrastructure) RecordNPSResponse(ctx context.Context, input dto.NPSInput) (bool, error) {
	return f.RecordNPSResponseFn(ctx, input)
}

// RecordPatientFeedback
func (f *FakeInfrastructure) RecordPatientFeedback(ctx context.Context, input dto.PatientFeedbackInput) (bool, error) {
	return f.RecordPatientFeedbackFn(ctx, input)
}

// Room ...
func (f *FakeInfrastructure) Room(ctx context.Context) (*dto.Room, error) {
	return f.RoomFn(ctx)
}

// TwilioAccessToken ...
func (f *FakeInfrastructure) TwilioAccessToken(ctx context.Context) (*dto.AccessToken, error) {
	return f.TwilioAccessTokenFn(ctx)
}

// SendSMS ...
func (f *FakeInfrastructure) SendSMS(ctx context.Context, to string, msg string) error {
	return f.SendSMSFn(ctx, to, msg)
}

// Upload ...
func (f *FakeInfrastructure) Upload(
	ctx context.Context,
	inp profileutils.UploadInput,
) (*profileutils.Upload, error) {
	return f.UploadFn(ctx, inp)
}

// FindUploadByID ...
func (f *FakeInfrastructure) FindUploadByID(
	ctx context.Context,
	id string,
) (*profileutils.Upload, error) {
	return f.FindUploadByIDFn(ctx, id)
}

// PhoneNumberVerificationCode ...
func (f *FakeInfrastructure) PhoneNumberVerificationCode(
	ctx context.Context,
	to string,
	code string,
	marketingMessage string,
) (bool, error) {
	return f.PhoneNumberVerificationCodeFn(
		ctx,
		to,
		code,
		marketingMessage,
	)
}

// WellnessCardActivationDependant ...
func (f *FakeInfrastructure) WellnessCardActivationDependant(
	ctx context.Context,
	to string,
	memberName string,
	cardName string,
	marketingMessage string,
) (bool, error) {
	return f.WellnessCardActivationDependantFn(
		ctx,
		to,
		memberName,
		cardName,
		marketingMessage,
	)
}

// WellnessCardActivationPrincipal ...
func (f *FakeInfrastructure) WellnessCardActivationPrincipal(
	ctx context.Context,
	to string,
	memberName string,
	cardName string,
	minorAgeThreshold string,
	marketingMessage string,
) (bool, error) {
	return f.WellnessCardActivationPrincipalFn(
		ctx,
		to,
		memberName,
		cardName,
		minorAgeThreshold,
		marketingMessage,
	)
}

// BillNotification ...
func (f *FakeInfrastructure) BillNotification(
	ctx context.Context,
	to string,
	productName string,
	billingPeriod string,
	billAmount string,
	paymentInstruction string,
	marketingMessage string,
) (bool, error) {
	return f.BillNotificationFn(
		ctx,
		to,
		productName,
		billingPeriod,
		billAmount,
		paymentInstruction,
		marketingMessage,
	)
}

// VirtualCards ...
func (f *FakeInfrastructure) VirtualCards(
	ctx context.Context,
	to string,
	wellnessCardFamily string,
	virtualCardLink string,
	marketingMessage string,
) (bool, error) {
	return f.VirtualCardsFn(
		ctx,
		to,
		wellnessCardFamily,
		virtualCardLink,
		marketingMessage,
	)
}

// VisitStart ...
func (f *FakeInfrastructure) VisitStart(
	ctx context.Context,
	to string,
	memberName string,
	benefitName string,
	locationName string,
	startTime string,
	balance string,
	marketingMessage string,
) (bool, error) {
	return f.VisitStartFn(
		ctx,
		to,
		memberName,
		benefitName,
		locationName,
		startTime,
		balance,
		marketingMessage,
	)
}

// ClaimNotification ...
func (f *FakeInfrastructure) ClaimNotification(
	ctx context.Context,
	to string,
	claimReference string,
	claimTypeParenthesized string,
	provider string,
	visitType string,
	claimTime string,
	marketingMessage string,
) (bool, error) {
	return f.ClaimNotificationFn(
		ctx,
		to,
		claimReference,
		claimTypeParenthesized,
		provider,
		visitType,
		claimTime,
		marketingMessage,
	)
}

// PreauthApproval ...
func (f *FakeInfrastructure) PreauthApproval(
	ctx context.Context,
	to string,
	currency string,
	amount string,
	benefit string,
	provider string,
	member string,
	careContact string,
	marketingMessage string,
) (bool, error) {
	return f.PreauthApprovalFn(
		ctx,
		to,
		currency,
		amount,
		benefit,
		provider,
		member,
		careContact,
		marketingMessage,
	)
}

// PreauthRequest ...
func (f *FakeInfrastructure) PreauthRequest(
	ctx context.Context,
	to string,
	currency string,
	amount string,
	benefit string,
	provider string,
	requestTime string,
	member string,
	careContact string,
	marketingMessage string,
) (bool, error) {
	return f.PreauthRequestFn(
		ctx,
		to,
		currency,
		amount,
		benefit,
		provider,
		requestTime,
		member,
		careContact,
		marketingMessage,
	)
}

// SladeOTP ...
func (f *FakeInfrastructure) SladeOTP(
	ctx context.Context,
	to string,
	name string,
	otp string,
	marketingMessage string,
) (bool, error) {
	return f.SladeOTPFn(
		ctx,
		to,
		name,
		otp,
		marketingMessage,
	)
}

// SendNotification ...
func (f *FakeInfrastructure) SendNotification(
	ctx context.Context,
	registrationTokens []string,
	data map[string]string,
	notification *firebasetools.FirebaseSimpleNotificationInput,
	android *firebasetools.FirebaseAndroidConfigInput,
	ios *firebasetools.FirebaseAPNSConfigInput,
	web *firebasetools.FirebaseWebpushConfigInput,
) (bool, error) {
	return f.SendNotificationFn(
		ctx,
		registrationTokens,
		data,
		notification,
		android,
		ios,
		web,
	)
}

// Notifications ...
func (f *FakeInfrastructure) Notifications(
	ctx context.Context,
	registrationToken string,
	newerThan time.Time,
	limit int,
) ([]*dto.SavedNotification, error) {
	return f.NotificationsFn(
		ctx,
		registrationToken,
		newerThan,
		limit,
	)
}

// SendFCMByPhoneOrEmail ...
func (f *FakeInfrastructure) SendFCMByPhoneOrEmail(
	ctx context.Context,
	phoneNumber *string,
	email *string,
	data map[string]interface{},
	notification firebasetools.FirebaseSimpleNotificationInput,
	android *firebasetools.FirebaseAndroidConfigInput,
	ios *firebasetools.FirebaseAPNSConfigInput,
	web *firebasetools.FirebaseWebpushConfigInput,
) (bool, error) {
	return f.SendFCMByPhoneOrEmailFn(
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

// GetFeedContent ...
func (f *FakeInfrastructure) GetFeedContent(ctx context.Context, flavour feedlib.Flavour) ([]*domain.GhostCMSPost, error) {
	return f.GetFeedContentFn(ctx, flavour)
}

// GetFaqsContent ...
func (f *FakeInfrastructure) GetFaqsContent(ctx context.Context, flavour feedlib.Flavour) ([]*domain.GhostCMSPost, error) {
	return f.GetFaqsContentFn(ctx, flavour)
}

// GetLibraryContent ...
func (f *FakeInfrastructure) GetLibraryContent(ctx context.Context) ([]*domain.GhostCMSPost, error) {
	return f.GetLibraryContentFn(ctx)
}
