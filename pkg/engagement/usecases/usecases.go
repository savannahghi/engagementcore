package usecases

import (
	"context"
	"time"

	"github.com/savannahghi/engagement/pkg/engagement/application/common/dto"
	"github.com/savannahghi/engagement/pkg/engagement/application/common/helpers"
	"github.com/savannahghi/engagement/pkg/engagement/domain"
	"github.com/savannahghi/engagement/pkg/engagement/infrastructure"
	"github.com/savannahghi/engagement/pkg/engagement/usecases/fcm"
	"github.com/savannahghi/engagement/pkg/engagement/usecases/feed"
	"github.com/savannahghi/engagement/pkg/engagement/usecases/library"
	"github.com/savannahghi/engagement/pkg/engagement/usecases/mail"
	"github.com/savannahghi/engagement/pkg/engagement/usecases/otp"
	"github.com/savannahghi/engagement/pkg/engagement/usecases/sms"
	"github.com/savannahghi/engagement/pkg/engagement/usecases/surveys"
	tc "github.com/savannahghi/engagement/pkg/engagement/usecases/teleconsult"
	"github.com/savannahghi/engagement/pkg/engagement/usecases/uploads"
	"github.com/savannahghi/engagement/pkg/engagement/usecases/whatsapp"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/profileutils"
	"github.com/savannahghi/pubsubtools"
)

// Usecases is an interface that combines of all usescases
type Usecases interface {
	feed.FeedUseCases
	feed.NotificationUsecases
	fcm.FCMUsecases
	library.LibraryUsecases
	mail.MailUsecases
	whatsapp.WhatsappUsecases
	uploads.UploadsUsecases
	otp.OTPUsecases
	sms.SMSUsecases
	surveys.SurveysUsecases
	tc.TeleconsultUsecases
}

// Interactor is an implementation of the usecases interface
type Interactor struct {
	feed         *feed.FeedUseCaseImpl
	notification *feed.NotificationImpl
	fcm          *fcm.FCMUsecasesImpl
	lib          *library.LibraryUsecasesImpl
	mail         *mail.MailUsecasesImpl
	whatsapp     *whatsapp.WhatsappUsecasesImpl
	uploads      *uploads.UploadsImpl
	otp          *otp.OTPUsecasesImpl
	sms          *sms.SMSUsecasesImpl
	surveys      *surveys.SurveysImpl
	teleconsult  *tc.TeleconsultImpl
}

// NewUsecasesInteractor initializes a new usecases interactor
func NewUsecasesInteractor(infrastructure infrastructure.Infrastructure) Usecases {

	notification := feed.NewNotification(infrastructure)
	feed := feed.NewFeed(infrastructure)
	fcm := fcm.NewFCMUsecaseImpl(infrastructure)
	library := library.NewLibraryUsecasesImpl(infrastructure)
	mail := mail.NewMailUsecasesImpl(infrastructure)
	whatsapp := whatsapp.NewWhatsappUsecasesImpl(infrastructure)
	uploads := uploads.NewUploadsImpl(infrastructure)
	otp := otp.NewOTPUsecasesImpl(infrastructure)
	sms := sms.NewSMSUsecasesImpl(infrastructure)
	surveys := surveys.NewSurveysImpl(infrastructure)
	teleconsult := tc.NewTeleconsultImpl(infrastructure)

	impl := &Interactor{
		feed:         feed,
		notification: notification,
		fcm:          fcm,
		lib:          library,
		mail:         mail,
		whatsapp:     whatsapp,
		uploads:      uploads,
		otp:          otp,
		sms:          sms,
		surveys:      surveys,
		teleconsult:  teleconsult,
	}

	return impl
}

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
	return i.feed.GetFeed(ctx, uid, isAnonymous, flavour, persistent, status, visibility, expired, filterParams)
}

func (i *Interactor) GetThinFeed(
	ctx context.Context,
	uid *string,
	isAnonymous *bool,
	flavour feedlib.Flavour,
) (*domain.Feed, error) {
	return i.feed.GetThinFeed(ctx, uid, isAnonymous, flavour)
}

func (i *Interactor) GetFeedItem(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	itemID string,
) (*feedlib.Item, error) {
	return i.feed.GetFeedItem(ctx, uid, flavour, itemID)
}

func (i *Interactor) GetNudge(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	nudgeID string,
) (*feedlib.Nudge, error) {
	return i.feed.GetNudge(ctx, uid, flavour, nudgeID)
}

func (i *Interactor) GetAction(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	actionID string,
) (*feedlib.Action, error) {
	return i.feed.GetAction(ctx, uid, flavour, actionID)
}

func (i *Interactor) PublishFeedItem(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	item *feedlib.Item,
) (*feedlib.Item, error) {
	return i.feed.PublishFeedItem(ctx, uid, flavour, item)
}

func (i *Interactor) DeleteFeedItem(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	itemID string,
) error {
	return i.feed.DeleteFeedItem(ctx, uid, flavour, itemID)
}

func (i *Interactor) ResolveFeedItem(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	itemID string,
) (*feedlib.Item, error) {
	return i.feed.ResolveFeedItem(ctx, uid, flavour, itemID)
}

func (i *Interactor) PinFeedItem(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	itemID string,
) (*feedlib.Item, error) {
	return i.feed.PinFeedItem(ctx, uid, flavour, itemID)
}

func (i *Interactor) UnpinFeedItem(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	itemID string,
) (*feedlib.Item, error) {
	return i.feed.UnpinFeedItem(ctx, uid, flavour, itemID)
}

func (i *Interactor) UnresolveFeedItem(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	itemID string,
) (*feedlib.Item, error) {
	return i.feed.UnresolveFeedItem(ctx, uid, flavour, itemID)
}

func (i *Interactor) HideFeedItem(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	itemID string,
) (*feedlib.Item, error) {
	return i.feed.HideFeedItem(ctx, uid, flavour, itemID)
}

func (i *Interactor) ShowFeedItem(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	itemID string,
) (*feedlib.Item, error) {
	return i.feed.ShowFeedItem(ctx, uid, flavour, itemID)
}

func (i *Interactor) Labels(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
) ([]string, error) {
	return i.feed.Labels(ctx, uid, flavour)
}

func (i *Interactor) SaveLabel(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	label string,
) error {
	return i.feed.SaveLabel(ctx, uid, flavour, label)
}

func (i *Interactor) UnreadPersistentItems(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
) (int, error) {
	return i.feed.UnreadPersistentItems(ctx, uid, flavour)
}

func (i *Interactor) UpdateUnreadPersistentItemsCount(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
) error {
	return i.feed.UpdateUnreadPersistentItemsCount(ctx, uid, flavour)
}

func (i *Interactor) PublishNudge(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	nudge *feedlib.Nudge,
) (*feedlib.Nudge, error) {
	return i.feed.PublishNudge(ctx, uid, flavour, nudge)
}

func (i *Interactor) ResolveNudge(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	nudgeID string,
) (*feedlib.Nudge, error) {
	return i.feed.ResolveNudge(ctx, uid, flavour, nudgeID)
}

func (i *Interactor) UnresolveNudge(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	nudgeID string,
) (*feedlib.Nudge, error) {
	return i.feed.UnresolveNudge(ctx, uid, flavour, nudgeID)
}

func (i *Interactor) HideNudge(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	nudgeID string,
) (*feedlib.Nudge, error) {
	return i.feed.HideNudge(ctx, uid, flavour, nudgeID)
}

func (i *Interactor) ShowNudge(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	nudgeID string,
) (*feedlib.Nudge, error) {
	return i.feed.ResolveNudge(ctx, uid, flavour, nudgeID)
}

func (i *Interactor) GetDefaultNudgeByTitle(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	title string,
) (*feedlib.Nudge, error) {
	return i.feed.GetDefaultNudgeByTitle(ctx, uid, flavour, title)
}

func (i *Interactor) ProcessEvent(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	event *feedlib.Event,
) error {
	return i.feed.ProcessEvent(ctx, uid, flavour, event)
}

func (i *Interactor) DeleteMessage(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	itemID string,
	messageID string,
) error {
	return i.feed.DeleteMessage(ctx, uid, flavour, itemID, messageID)
}

func (i *Interactor) PostMessage(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	itemID string,
	message *feedlib.Message,
) (*feedlib.Message, error) {
	return i.feed.PostMessage(ctx, uid, flavour, itemID, message)
}

func (i *Interactor) DeleteAction(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	actionID string,
) error {
	return i.feed.DeleteAction(ctx, uid, flavour, actionID)
}

func (i *Interactor) PublishAction(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	action *feedlib.Action,
) (*feedlib.Action, error) {
	return i.feed.PublishAction(ctx, uid, flavour, action)
}

func (i *Interactor) DeleteNudge(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	nudgeID string,
) error {
	return i.feed.DeleteNudge(ctx, uid, flavour, nudgeID)
}

func (i *Interactor) HandleItemPublish(
	ctx context.Context,
	m *pubsubtools.PubSubPayload,
) error {
	return i.notification.HandleItemPublish(ctx, m)
}

func (i *Interactor) HandleItemDelete(
	ctx context.Context,
	m *pubsubtools.PubSubPayload,
) error {
	return i.notification.HandleItemDelete(ctx, m)
}

func (i *Interactor) HandleItemResolve(
	ctx context.Context,
	m *pubsubtools.PubSubPayload,
) error {
	return i.notification.HandleItemResolve(ctx, m)
}

func (i *Interactor) HandleItemUnresolve(
	ctx context.Context,
	m *pubsubtools.PubSubPayload,
) error {
	return i.notification.HandleItemUnresolve(ctx, m)
}

func (i *Interactor) HandleItemHide(
	ctx context.Context,
	m *pubsubtools.PubSubPayload,
) error {
	return i.notification.HandleItemHide(ctx, m)
}

func (i *Interactor) HandleItemShow(
	ctx context.Context,
	m *pubsubtools.PubSubPayload,
) error {
	return i.notification.HandleItemShow(ctx, m)
}

func (i *Interactor) HandleItemPin(
	ctx context.Context,
	m *pubsubtools.PubSubPayload,
) error {
	return i.notification.HandleItemPin(ctx, m)
}

func (i *Interactor) HandleItemUnpin(
	ctx context.Context,
	m *pubsubtools.PubSubPayload,
) error {
	return i.notification.HandleItemUnpin(ctx, m)
}

func (i *Interactor) HandleNudgePublish(
	ctx context.Context,
	m *pubsubtools.PubSubPayload,
) error {
	return i.notification.HandleNudgePublish(ctx, m)
}

func (i *Interactor) HandleNudgeDelete(
	ctx context.Context,
	m *pubsubtools.PubSubPayload,
) error {
	return i.notification.HandleNudgeDelete(ctx, m)
}

func (i *Interactor) HandleNudgeResolve(
	ctx context.Context,
	m *pubsubtools.PubSubPayload,
) error {
	return i.notification.HandleNudgeResolve(ctx, m)
}

func (i *Interactor) HandleNudgeUnresolve(
	ctx context.Context,
	m *pubsubtools.PubSubPayload,
) error {
	return i.notification.HandleNudgeUnresolve(ctx, m)
}

func (i *Interactor) HandleNudgeHide(
	ctx context.Context,
	m *pubsubtools.PubSubPayload,
) error {
	return i.notification.HandleNudgeHide(ctx, m)
}

func (i *Interactor) HandleNudgeShow(
	ctx context.Context,
	m *pubsubtools.PubSubPayload,
) error {
	return i.notification.HandleNudgeShow(ctx, m)
}

func (i *Interactor) HandleActionPublish(
	ctx context.Context,
	m *pubsubtools.PubSubPayload,
) error {
	return i.notification.HandleActionPublish(ctx, m)
}

func (i *Interactor) HandleActionDelete(
	ctx context.Context,
	m *pubsubtools.PubSubPayload,
) error {
	return i.notification.HandleActionDelete(ctx, m)
}

func (i *Interactor) HandleMessagePost(
	ctx context.Context,
	m *pubsubtools.PubSubPayload,
) error {
	return i.notification.HandleMessagePost(ctx, m)
}

func (i *Interactor) HandleMessageDelete(
	ctx context.Context,
	m *pubsubtools.PubSubPayload,
) error {
	return i.notification.HandleMessageDelete(ctx, m)
}

func (i *Interactor) HandleIncomingEvent(
	ctx context.Context,
	m *pubsubtools.PubSubPayload,
) error {
	return i.notification.HandleIncomingEvent(ctx, m)
}

func (i *Interactor) NotifyItemUpdate(
	ctx context.Context,
	sender string,
	includeNotification bool, // whether to show a tray notification
	m *pubsubtools.PubSubPayload,
) error {
	return i.notification.NotifyItemUpdate(ctx, sender, includeNotification, m)
}

func (i *Interactor) UpdateInbox(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
) error {
	return i.notification.UpdateInbox(ctx, uid, flavour)
}

func (i *Interactor) NotifyNudgeUpdate(
	ctx context.Context,
	sender string,
	m *pubsubtools.PubSubPayload,
) error {
	return i.notification.NotifyNudgeUpdate(ctx, sender, m)
}

func (i *Interactor) NotifyInboxCountUpdate(
	ctx context.Context,
	uid string,
	flavour feedlib.Flavour,
	count int,
) error {
	return i.notification.NotifyInboxCountUpdate(ctx, uid, flavour, count)
}

func (i *Interactor) GetUserTokens(
	ctx context.Context,
	uids []string,
) ([]string, error) {
	return i.notification.GetUserTokens(ctx, uids)
}

func (i *Interactor) SendNotificationViaFCM(
	ctx context.Context,
	uids []string,
	sender string,
	pl dto.NotificationEnvelope,
	notification *firebasetools.FirebaseSimpleNotificationInput,
) error {
	return i.notification.SendNotificationViaFCM(ctx, uids, sender, pl, notification)
}

func (i *Interactor) HandleSendNotification(
	ctx context.Context,
	m *pubsubtools.PubSubPayload,
) error {
	return i.notification.HandleSendNotification(ctx, m)
}

func (i *Interactor) SendEmail(
	ctx context.Context,
	m *pubsubtools.PubSubPayload,
) error {
	return i.notification.SendEmail(ctx, m)
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

func (i *Interactor) GetFaqsContent(ctx context.Context, flavour feedlib.Flavour) ([]*domain.GhostCMSPost, error) {
	return i.lib.GetFaqsContent(ctx, flavour)
}

func (i *Interactor) GetLibraryContent(ctx context.Context) ([]*domain.GhostCMSPost, error) {
	return i.lib.GetLibraryContent(ctx)
}

func (i *Interactor) SimpleEmail(
	ctx context.Context,
	subject, text string,
	body *string,
	to []string,
) (string, error) {
	return i.mail.SimpleEmail(ctx, subject, text, body, to...)
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

func (i *Interactor) FindUploadByID(ctx context.Context, id string) (*profileutils.Upload, error) {
	return i.uploads.FindUploadByID(ctx, id)
}

func (i *Interactor) Upload(ctx context.Context, input profileutils.UploadInput) (*profileutils.Upload, error) {
	return i.uploads.Upload(ctx, input)
}

func (i *Interactor) GenerateAndSendOTP(ctx context.Context, msisdn string, appID *string) (string, error) {
	return i.otp.GenerateAndSendOTP(ctx, msisdn, appID)
}

func (i *Interactor) SendOTPToEmail(ctx context.Context, msisdn string, email *string, appID *string) (string, error) {
	return i.otp.SendOTPToEmail(ctx, msisdn, email, appID)
}

func (i *Interactor) VerifyOtp(ctx context.Context, msisdn, verificationCode string) (bool, error) {
	return i.otp.VerifyOtp(ctx, msisdn, verificationCode)
}

func (i *Interactor) VerifyEmailOtp(ctx context.Context, email, verificationCode string) (bool, error) {
	return i.otp.VerifyEmailOtp(ctx, email, verificationCode)
}

func (i *Interactor) GenerateRetryOTP(ctx context.Context, msisdn string, retryStep int, appID *string) (string, error) {
	return i.otp.GenerateRetryOTP(ctx, msisdn, retryStep, appID)
}

func (i *Interactor) EmailVerificationOtp(ctx context.Context, email string) (string, error) {
	return i.otp.EmailVerificationOtp(ctx, email)
}

func (i *Interactor) Send(ctx context.Context, to string, message string) (*dto.SendMessageResponse, error) {
	return i.sms.Send(ctx, to, message)
}
func (i *Interactor) SendToMany(ctx context.Context, message string, to []string) (*dto.SendMessageResponse, error) {
	return i.sms.SendToMany(ctx, message, to)
}

func (i *Interactor) TwilioAccessToken(ctx context.Context) (*dto.AccessToken, error) {
	return i.teleconsult.TwilioAccessToken(ctx)
}

func (i *Interactor) RecordNPSResponse(ctx context.Context, input dto.NPSInput) (bool, error) {
	return i.surveys.RecordNPSResponse(ctx, input)
}
