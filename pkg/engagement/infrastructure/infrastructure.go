package infrastructure

import (
	"context"

	"github.com/labstack/gommon/log"
	"github.com/savannahghi/engagementcore/pkg/engagement/infrastructure/database"
	"github.com/savannahghi/engagementcore/pkg/engagement/infrastructure/services/fcm"
	"github.com/savannahghi/engagementcore/pkg/engagement/infrastructure/services/library"
	"github.com/savannahghi/engagementcore/pkg/engagement/infrastructure/services/mail"
	"github.com/savannahghi/engagementcore/pkg/engagement/infrastructure/services/messaging"
	"github.com/savannahghi/engagementcore/pkg/engagement/infrastructure/services/onboarding"
	"github.com/savannahghi/engagementcore/pkg/engagement/infrastructure/services/otp"
	"github.com/savannahghi/engagementcore/pkg/engagement/infrastructure/services/sms"
	"github.com/savannahghi/engagementcore/pkg/engagement/infrastructure/services/surveys"
	"github.com/savannahghi/engagementcore/pkg/engagement/infrastructure/services/twilio"
	"github.com/savannahghi/engagementcore/pkg/engagement/infrastructure/services/uploads"
	"github.com/savannahghi/engagementcore/pkg/engagement/infrastructure/services/whatsapp"
	"github.com/savannahghi/engagementcore/pkg/engagement/repository"
	"github.com/savannahghi/serverutils"
)

// Interactor is an implementation of the infrastructure interface
// It combines each individual service implementation
type Interactor struct {
	repository.Repository
	*fcm.ServiceFCMImpl
	*fcm.RemotePushService
	*library.ServiceLibraryImpl
	*mail.ServiceMailImpl
	messaging.NotificationService
	onboarding.ProfileService
	*otp.ServiceOTPImpl
	*sms.ServiceSMSImpl
	*surveys.ServiceSurveyImpl
	*twilio.ServiceTwilioImpl
	*uploads.ServiceUploadImpl
	*whatsapp.ServiceWhatsappImpl
}

// NewInteractor initializes a new infrastructure interactor
func NewInteractor() Interactor {
	ctx := context.Background()

	projectID, err := serverutils.GetEnvVar(serverutils.GoogleCloudProjectIDEnvVarName)
	if err != nil {
		log.Fatal(err)
	}

	db := database.NewDbService()

	onboarding := onboarding.NewRemoteProfileService(onboarding.NewOnboardingClient())

	fcmOne := fcm.NewService(db, onboarding)
	push, err := fcm.NewRemotePushService(ctx)
	if err != nil {
		log.Fatal(err)
	}

	lib := library.NewLibraryService(onboarding)

	mail := mail.NewService(db)

	pubsub, err := messaging.NewPubSubNotificationService(ctx, projectID)
	if err != nil {
		log.Fatal(err)
	}

	whatsapp := whatsapp.NewService()
	sms := sms.NewService(db, pubsub)
	twilio := twilio.NewService(*sms, db)

	uploads := uploads.NewUploadsService()

	otp := otp.NewService(*whatsapp, *mail, *sms, twilio)

	surveys := surveys.NewService(db)

	return Interactor{
		db,
		fcmOne,
		push,
		lib,
		mail,
		pubsub,
		onboarding,
		otp,
		sms,
		surveys,
		twilio,
		uploads,
		whatsapp,
	}
}
