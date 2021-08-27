package infrastructure

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/labstack/gommon/log"
	"github.com/savannahghi/engagement/pkg/engagement/application/common/dto"
	"github.com/savannahghi/engagement/pkg/engagement/application/common/helpers"
	"github.com/savannahghi/engagement/pkg/engagement/domain"
	"github.com/savannahghi/engagement/pkg/engagement/infrastructure/database"
	"github.com/savannahghi/engagement/pkg/engagement/infrastructure/services/fcm"
	"github.com/savannahghi/engagement/pkg/engagement/infrastructure/services/library"
	"github.com/savannahghi/engagement/pkg/engagement/infrastructure/services/mail"
	"github.com/savannahghi/engagement/pkg/engagement/infrastructure/services/messaging"
	"github.com/savannahghi/engagement/pkg/engagement/infrastructure/services/onboarding"
	"github.com/savannahghi/engagement/pkg/engagement/infrastructure/services/otp"
	"github.com/savannahghi/engagement/pkg/engagement/infrastructure/services/sms"
	"github.com/savannahghi/engagement/pkg/engagement/infrastructure/services/surveys"
	"github.com/savannahghi/engagement/pkg/engagement/infrastructure/services/twilio"
	"github.com/savannahghi/engagement/pkg/engagement/infrastructure/services/uploads"
	"github.com/savannahghi/engagement/pkg/engagement/infrastructure/services/whatsapp"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/profileutils"
	"github.com/savannahghi/serverutils"
	crmDomain "gitlab.slade360emr.com/go/commontools/crm/pkg/domain"
)

// Infrastructure defines the contract provided by the infrastructure layer
// It's a combination of interactions with external services/dependencies
//
// Enables easy mocking during tests
type Infrastructure interface {
	database.Repository
	fcm.ServiceFCM
	fcm.PushService
	library.ServiceLibrary
	mail.ServiceMail
	messaging.NotificationService
	onboarding.ProfileService
	otp.ServiceOTP
	sms.ServiceSMS
	surveys.ServiceSurveys
	twilio.ServiceTwilio
	uploads.ServiceUploads
	whatsapp.ServiceWhatsapp
}

// Interactor is an implementation of the infrastructure interface
// It combines each individual service implementation
type Interactor struct {
	db         database.Repository
	fcm        *fcm.Service
	fcmTwo     *fcm.RemotePushService
	lib        *library.Service
	mail       *mail.Service
	messaging  messaging.NotificationService
	onboarding onboarding.ProfileService
	otp        *otp.Service
	sms        *sms.Service
	surveys    *surveys.Service
	twilio     *twilio.Service
	uploads    *uploads.Service
	whatsapp   *whatsapp.Service
}

// NewInfrastructureInteractor initializes a new infrastructure interactor
func NewInfrastructureInteractor() Infrastructure {
	ctx := context.Background()

	projectID, err := serverutils.GetEnvVar(serverutils.GoogleCloudProjectIDEnvVarName)
	if err != nil {
		// TODO
		log.Error(err)
	}

	db := database.NewDbService()

	onboarding := onboarding.NewRemoteProfileService(onboarding.NewOnboardingClient())

	fcmOne := fcm.NewService(db, onboarding)
	fcmTwo, err := fcm.NewRemotePushService(ctx)
	if err != nil {
		// TODO
		log.Error(err)
	}

	lib := library.NewLibraryService(onboarding)

	mail := mail.NewService(db)

	pubsub, err := messaging.NewPubSubNotificationService(ctx, projectID)
	if err != nil {
		// TODO
		log.Error(err)
	}

	whatsapp := whatsapp.NewService()
	sms := sms.NewService(db, pubsub)
	twilio := twilio.NewService(sms, db)

	uploads := uploads.NewUploadsService()

	otp := otp.NewService(whatsapp, mail, sms, twilio)

	surveys := surveys.NewService(db)

	return &Interactor{
		db:         db,
		fcm:        fcmOne,
		fcmTwo:     fcmTwo,
		lib:        lib,
		mail:       mail,
		messaging:  pubsub,
		onboarding: onboarding,
		otp:        otp,
		sms:        sms,
		surveys:    surveys,
		twilio:     twilio,
		uploads:    uploads,
		whatsapp:   whatsapp,
	}
}

// CheckPreconditions ensures correct initialization
func (i Interactor) CheckPreconditions() {}

func (i *Interactor) GetFeed(
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
	return i.db.GetFeed(ctx, uid, isAnonymous, flavour, persistent, status, visibility, expired, filterParams)
}

func (i *Interactor) GetFeedItem(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	itemID string,
) (*feedlib.Item, error) {
	return i.db.GetFeedItem(ctx, uid, flavour, itemID)
}

func (i *Interactor) SaveFeedItem(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	item *feedlib.Item,
) (*feedlib.Item, error) {
	return i.db.SaveFeedItem(ctx, uid, flavour, item)
}

func (i *Interactor) UpdateFeedItem(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	item *feedlib.Item,
) (*feedlib.Item, error) {
	return i.db.UpdateFeedItem(ctx, uid, flavour, item)
}

// DeleteFeedItem permanently deletes a feed item and it's copies
func (i *Interactor) DeleteFeedItem(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	itemID string,
) error {
	return i.db.DeleteFeedItem(ctx, uid, flavour, itemID)
}

// GetNudge gets THE LATEST VERSION OF a nudge from a feed
func (i *Interactor) GetNudge(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	nudgeID string,
) (*feedlib.Nudge, error) {
	return i.db.GetNudge(ctx, uid, flavour, nudgeID)
}

// SaveNudge saves a new modified nudge
func (i *Interactor) SaveNudge(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	nudge *feedlib.Nudge,
) (*feedlib.Nudge, error) {
	return i.db.SaveNudge(ctx, uid, flavour, nudge)
}

// UpdateNudge updates an existing nudge
func (i *Interactor) UpdateNudge(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	nudge *feedlib.Nudge,
) (*feedlib.Nudge, error) {
	return i.db.UpdateNudge(ctx, uid, flavour, nudge)
}

// DeleteNudge permanently deletes a nudge and it's copies
func (i *Interactor) DeleteNudge(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	nudgeID string,
) error {
	return i.db.DeleteNudge(ctx, uid, flavour, nudgeID)
}

// GetAction gets THE LATEST VERSION OF a single action
func (i *Interactor) GetAction(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	actionID string,
) (*feedlib.Action, error) {
	return i.db.GetAction(ctx, uid, flavour, actionID)
}

// SaveAction saves a new action
func (i *Interactor) SaveAction(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	action *feedlib.Action,
) (*feedlib.Action, error) {
	return i.db.SaveAction(ctx, uid, flavour, action)
}

// DeleteAction permanently deletes an action and it's copies
func (i *Interactor) DeleteAction(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	actionID string,
) error {
	return i.db.DeleteAction(ctx, uid, flavour, actionID)
}

// PostMessage posts a message or a reply to a message/thread
func (i *Interactor) PostMessage(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	itemID string,
	message *feedlib.Message,
) (*feedlib.Message, error) {
	return i.db.PostMessage(ctx, uid, flavour, itemID, message)
}

// GetMessage retrieves THE LATEST VERSION OF a message
func (i *Interactor) GetMessage(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	itemID string,
	messageID string,
) (*feedlib.Message, error) {
	return i.db.GetMessage(ctx, uid, flavour, itemID, messageID)
}

// DeleteMessage deletes a message
func (i *Interactor) DeleteMessage(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	itemID string,
	messageID string,
) error {
	return i.db.DeleteMessage(ctx, uid, flavour, itemID, messageID)
}

// GetMessages retrieves a message
func (i *Interactor) GetMessages(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	itemID string,
) ([]feedlib.Message, error) {
	return i.db.GetMessages(ctx, uid, flavour, itemID)
}

func (i *Interactor) SaveIncomingEvent(
	ctx context.Context,
	event *feedlib.Event,
) error {
	return i.db.SaveIncomingEvent(ctx, event)
}

func (i *Interactor) SaveOutgoingEvent(
	ctx context.Context,
	event *feedlib.Event,
) error {
	return i.db.SaveOutgoingEvent(ctx, event)
}

func (i *Interactor) GetNudges(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	status *feedlib.Status,
	visibility *feedlib.Visibility,
	expired *feedlib.BooleanFilter,
) ([]feedlib.Nudge, error) {
	return i.db.GetNudges(ctx, uid, flavour, status, visibility, expired)
}

func (i *Interactor) GetActions(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
) ([]feedlib.Action, error) {
	return i.db.GetActions(ctx, uid, flavour)
}

func (i *Interactor) GetItems(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	persistent feedlib.BooleanFilter,
	status *feedlib.Status,
	visibility *feedlib.Visibility,
	expired *feedlib.BooleanFilter,
	filterParams *helpers.FilterParams,
) ([]feedlib.Item, error) {
	return i.db.GetItems(ctx, uid, flavour, persistent, status, visibility, expired, filterParams)
}

func (i *Interactor) Labels(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
) ([]string, error) {
	return i.db.Labels(ctx, uid, flavour)
}

func (i *Interactor) SaveLabel(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	label string,
) error {
	return i.db.SaveLabel(ctx, uid, flavour, label)
}

func (i *Interactor) UnreadPersistentItems(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
) (int, error) {
	return i.db.UnreadPersistentItems(ctx, uid, flavour)
}

func (i *Interactor) UpdateUnreadPersistentItemsCount(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
) error {
	return i.db.UpdateUnreadPersistentItemsCount(ctx, uid, flavour)
}

func (i *Interactor) GetDefaultNudgeByTitle(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	title string,
) (*feedlib.Nudge, error) {
	return i.db.GetDefaultNudgeByTitle(ctx, uid, flavour, title)
}

// SaveMarketingMessage saves the callback data for future analysis
func (i *Interactor) SaveMarketingMessage(
	ctx context.Context,
	data dto.MarketingSMS,
) (*dto.MarketingSMS, error) {
	return i.db.SaveMarketingMessage(ctx, data)
}

// SaveTwilioResponse saves the callback data for future analysis
func (i *Interactor) SaveTwilioResponse(
	ctx context.Context,
	data dto.Message,
) error {
	return i.db.SaveTwilioResponse(ctx, data)
}

// SaveNotification saves a notification
func (i *Interactor) SaveNotification(
	ctx context.Context,
	firestoreClient *firestore.Client,
	notification dto.SavedNotification,
) error {
	return i.db.SaveNotification(ctx, firestoreClient, notification)
}

// RetrieveNotification retrieves a notification
func (i *Interactor) RetrieveNotification(
	ctx context.Context,
	firestoreClient *firestore.Client,
	registrationToken string,
	newerThan time.Time,
	limit int,
) ([]*dto.SavedNotification, error) {
	return i.db.RetrieveNotification(ctx, firestoreClient, registrationToken, newerThan, limit)
}

// SaveNPSResponse saves a NPS response
func (i *Interactor) SaveNPSResponse(
	ctx context.Context,
	response *dto.NPSResponse,
) error {
	return i.db.SaveNPSResponse(ctx, response)
}

// UpdateMarketingMessage ..
func (i *Interactor) UpdateMarketingMessage(
	ctx context.Context,
	data *dto.MarketingSMS,
) (*dto.MarketingSMS, error) {
	return i.db.UpdateMarketingMessage(ctx, data)
}

func (i *Interactor) SaveOutgoingEmails(ctx context.Context, payload *dto.OutgoingEmailsLog) error {
	return i.db.SaveOutgoingEmails(ctx, payload)
}

func (i *Interactor) UpdateMailgunDeliveryStatus(ctx context.Context, payload *dto.MailgunEvent) (*dto.OutgoingEmailsLog, error) {
	return i.db.UpdateMailgunDeliveryStatus(ctx, payload)
}

// GetMarketingSMSByPhone ..
func (i *Interactor) GetMarketingSMSByPhone(ctx context.Context, phoneNumber string) (*dto.MarketingSMS, error) {
	return i.db.GetMarketingSMSByPhone(ctx, phoneNumber)
}

// GetMarketingSMSByID ..
func (i *Interactor) GetMarketingSMSByID(
	ctx context.Context,
	id string,
) (*dto.MarketingSMS, error) {
	return i.db.GetMarketingSMSByID(ctx, id)
}

// SaveTwilioVideoCallbackStatus ..
func (i *Interactor) SaveTwilioVideoCallbackStatus(
	ctx context.Context,
	data dto.CallbackData,
) error {
	return i.db.SaveTwilioVideoCallbackStatus(ctx, data)
}

func (i *Interactor) SendNotification(
	ctx context.Context,
	registrationTokens []string,
	data map[string]string,
	notification *firebasetools.FirebaseSimpleNotificationInput,
	android *firebasetools.FirebaseAndroidConfigInput,
	ios *firebasetools.FirebaseAPNSConfigInput,
	web *firebasetools.FirebaseWebpushConfigInput,
) (bool, error) {
	return i.fcm.SendNotification(ctx, registrationTokens, data, notification, android, ios, web)
}

// Notifications is used to query a user's priorities
func (i *Interactor) Notifications(
	ctx context.Context,
	registrationToken string,
	newerThan time.Time,
	limit int,
) ([]*dto.SavedNotification, error) {
	return i.fcm.Notifications(ctx, registrationToken, newerThan, limit)
}

// SendFCMByPhoneOrEmail is used to send FCM notification by phone or email
func (i *Interactor) SendFCMByPhoneOrEmail(
	ctx context.Context,
	phoneNumber *string,
	email *string,
	data map[string]interface{},
	notification firebasetools.FirebaseSimpleNotificationInput,
	android *firebasetools.FirebaseAndroidConfigInput,
	ios *firebasetools.FirebaseAPNSConfigInput,
	web *firebasetools.FirebaseWebpushConfigInput,
) (bool, error) {
	return i.fcm.SendFCMByPhoneOrEmail(ctx, phoneNumber, email, data, notification, android, ios, web)
}

func (i *Interactor) GetFeedContent(ctx context.Context, flavour feedlib.Flavour) ([]*domain.GhostCMSPost, error) {
	return i.lib.GetFeedContent(ctx, flavour)
}

func (i *Interactor) GetFaqsContent(ctx context.Context, flavour feedlib.Flavour) ([]*domain.GhostCMSPost, error) {
	return i.lib.GetFaqsContent(ctx, flavour)
}

func (i *Interactor) GetLibraryContent(ctx context.Context) ([]*domain.GhostCMSPost, error) {
	return i.lib.GetLibraryContent(ctx)
}

func (i *Interactor) SendInBlue(
	ctx context.Context,
	subject, text string,
	to ...string,
) (string, string, error) {
	return i.mail.SendInBlue(ctx, subject, text, to...)
}

func (i *Interactor) SendMailgun(
	ctx context.Context,
	subject, text string,
	body *string,
	to ...string,
) (string, string, error) {
	return i.mail.SendMailgun(ctx, subject, text, body, to...)
}

func (i *Interactor) SendEmail(
	ctx context.Context,
	subject, text string,
	body *string,
	to ...string,
) (string, string, error) {
	return i.mail.SendMailgun(ctx, subject, text, body, to...)
}

func (i *Interactor) SimpleEmail(
	ctx context.Context,
	subject, text string,
	body *string,
	to ...string,
) (string, error) {
	return i.mail.SimpleEmail(ctx, subject, text, body, to...)
}

func (i *Interactor) GenerateEmailTemplate(name string, templateName string) string {
	return i.mail.GenerateEmailTemplate(name, templateName)
}

func (i *Interactor) Notify(
	ctx context.Context,
	topicID string,
	uid string,
	flavour feedlib.Flavour,
	payload feedlib.Element,
	metadata map[string]interface{},
) error {
	return i.messaging.Notify(ctx, topicID, uid, flavour, payload, metadata)
}

func (i *Interactor) NotifyEngagementCreate(
	ctx context.Context,
	phone string,
	messageID string,
	engagementType crmDomain.EngagementType,
	metadata map[string]interface{},
	topicID string,
) error {
	return i.messaging.NotifyEngagementCreate(ctx, phone, messageID, engagementType, metadata, topicID)
}

func (i *Interactor) TopicIDs() []string {
	return i.messaging.TopicIDs()
}

func (i *Interactor) SubscriptionIDs() map[string]string {
	return i.messaging.SubscriptionIDs()
}

func (i *Interactor) ReverseSubscriptionIDs() map[string]string {
	return i.messaging.ReverseSubscriptionIDs()
}

func (i *Interactor) GetEmailAddresses(
	ctx context.Context,
	uids onboarding.UserUIDs,
) (map[string][]string, error) {
	return i.onboarding.GetEmailAddresses(ctx, uids)
}

func (i *Interactor) GetPhoneNumbers(
	ctx context.Context,
	uids onboarding.UserUIDs,
) (map[string][]string, error) {
	return i.onboarding.GetPhoneNumbers(ctx, uids)
}

func (i *Interactor) GetDeviceTokens(
	ctx context.Context,
	uid onboarding.UserUIDs,
) (map[string][]string, error) {
	return i.onboarding.GetDeviceTokens(ctx, uid)
}

func (i *Interactor) GetUserProfile(ctx context.Context, uid string) (*profileutils.UserProfile, error) {
	return i.onboarding.GetUserProfile(ctx, uid)
}

func (i *Interactor) GetUserProfileByPhoneOrEmail(ctx context.Context, payload *dto.RetrieveUserProfileInput) (*profileutils.UserProfile, error) {
	return i.onboarding.GetUserProfileByPhoneOrEmail(ctx, payload)
}

func (i *Interactor) GenerateAndSendOTP(ctx context.Context, msisdn string, appID *string) (string, error) {
	return i.otp.GenerateAndSendOTP(ctx, msisdn, appID)
}

func (i *Interactor) SendOTPToEmail(ctx context.Context, msisdn, email *string, appID *string) (string, error) {
	return i.otp.SendOTPToEmail(ctx, msisdn, email, appID)
}
func (i *Interactor) SaveOTPToFirestore(otp dto.OTP) error {
	return i.otp.SaveOTPToFirestore(otp)
}
func (i *Interactor) VerifyOtp(ctx context.Context, msisdn, verificationCode *string) (bool, error) {
	return i.otp.VerifyOtp(ctx, msisdn, verificationCode)
}

func (i *Interactor) VerifyEmailOtp(ctx context.Context, email, verificationCode *string) (bool, error) {
	return i.otp.VerifyEmailOtp(ctx, email, verificationCode)
}

func (i *Interactor) GenerateRetryOTP(ctx context.Context, msisdn *string, retryStep int, appID *string) (string, error) {
	return i.otp.GenerateRetryOTP(ctx, msisdn, retryStep, appID)
}

func (i *Interactor) EmailVerificationOtp(ctx context.Context, email *string) (string, error) {
	return i.otp.EmailVerificationOtp(ctx, email)
}

func (i *Interactor) GenerateOTP(ctx context.Context) (string, error) {
	return i.otp.GenerateOTP(ctx)
}

func (i *Interactor) SendToMany(
	ctx context.Context,
	message string,
	to []string,
	from enumutils.SenderID,
) (*dto.SendMessageResponse, error) {
	return i.sms.SendToMany(ctx, message, to, from)
}
func (i *Interactor) Send(
	ctx context.Context,
	to, message string,
	from enumutils.SenderID,
) (*dto.SendMessageResponse, error) {
	return i.sms.Send(ctx, to, message, from)
}

func (i *Interactor) RecordNPSResponse(ctx context.Context, input dto.NPSInput) (bool, error) {
	return i.surveys.RecordNPSResponse(ctx, input)
}

func (i *Interactor) Room(ctx context.Context) (*dto.Room, error) {
	return i.twilio.Room(ctx)
}
func (i *Interactor) TwilioAccessToken(ctx context.Context) (*dto.AccessToken, error) {
	return i.twilio.TwilioAccessToken(ctx)
}
func (i *Interactor) SendSMS(ctx context.Context, to string, msg string) error {
	return i.twilio.SendSMS(ctx, to, msg)
}

func (i *Interactor) Upload(
	ctx context.Context,
	inp profileutils.UploadInput,
) (*profileutils.Upload, error) {
	return i.uploads.Upload(ctx, inp)
}

func (i *Interactor) FindUploadByID(
	ctx context.Context,
	id string,
) (*profileutils.Upload, error) {
	return i.uploads.FindUploadByID(ctx, id)
}

func (i *Interactor) PhoneNumberVerificationCode(
	ctx context.Context,
	to string,
	code string,
	marketingMessage string,
) (bool, error) {
	return i.whatsapp.PhoneNumberVerificationCode(ctx, to, code, marketingMessage)
}

func (i *Interactor) WellnessCardActivationDependant(
	ctx context.Context,
	to string,
	memberName string,
	cardName string,
	marketingMessage string,
) (bool, error) {
	return i.whatsapp.WellnessCardActivationDependant(ctx, to, memberName, cardName, marketingMessage)
}

func (i *Interactor) WellnessCardActivationPrincipal(
	ctx context.Context,
	to string,
	memberName string,
	cardName string,
	minorAgeThreshold string,
	marketingMessage string,
) (bool, error) {
	return i.whatsapp.WellnessCardActivationPrincipal(ctx, to, memberName, cardName, minorAgeThreshold, marketingMessage)
}

func (i *Interactor) BillNotification(
	ctx context.Context,
	to string,
	productName string,
	billingPeriod string,
	billAmount string,
	paymentInstruction string,
	marketingMessage string,
) (bool, error) {
	return i.whatsapp.BillNotification(ctx, to, productName, billingPeriod, billAmount, paymentInstruction, marketingMessage)
}

func (i *Interactor) VirtualCards(
	ctx context.Context,
	to string,
	wellnessCardFamily string,
	virtualCardLink string,
	marketingMessage string,
) (bool, error) {
	return i.whatsapp.VirtualCards(ctx, to, wellnessCardFamily, virtualCardLink, marketingMessage)
}

func (i *Interactor) VisitStart(
	ctx context.Context,
	to string,
	memberName string,
	benefitName string,
	locationName string,
	startTime string,
	balance string,
	marketingMessage string,
) (bool, error) {
	return i.whatsapp.VisitStart(ctx, to, memberName, benefitName, locationName, startTime, balance, marketingMessage)
}

func (i *Interactor) ClaimNotification(
	ctx context.Context,
	to string,
	claimReference string,
	claimTypeParenthesized string,
	provider string,
	visitType string,
	claimTime string,
	marketingMessage string,
) (bool, error) {
	return i.whatsapp.ClaimNotification(ctx, to, claimReference, claimTypeParenthesized, provider, visitType, claimTime, marketingMessage)
}

func (i *Interactor) PreauthApproval(
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
	return i.whatsapp.PreauthApproval(ctx, to, currency, amount, benefit, provider, member, careContact, marketingMessage)
}

func (i *Interactor) PreauthRequest(
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
	return i.whatsapp.PreauthRequest(ctx, to, currency, amount, benefit, provider, requestTime, member, careContact, marketingMessage)
}

func (i *Interactor) SladeOTP(
	ctx context.Context,
	to string,
	name string,
	otp string,
	marketingMessage string,
) (bool, error) {
	return i.whatsapp.SladeOTP(ctx, to, name, otp, marketingMessage)
}

func (i *Interactor) SaveTwilioCallbackResponse(
	ctx context.Context,
	data dto.Message,
) error {
	return i.whatsapp.SaveTwilioCallbackResponse(ctx, data)
}

func (i *Interactor) Push(
	ctx context.Context,
	sender string,
	payload firebasetools.SendNotificationPayload,
) error {
	return i.fcmTwo.Push(ctx, sender, payload)
}
