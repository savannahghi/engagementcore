package usecases

import (
	"github.com/savannahghi/engagementcore/pkg/engagement/infrastructure"
	"github.com/savannahghi/engagementcore/pkg/engagement/usecases/fcm"
	"github.com/savannahghi/engagementcore/pkg/engagement/usecases/feed"
	"github.com/savannahghi/engagementcore/pkg/engagement/usecases/library"
	"github.com/savannahghi/engagementcore/pkg/engagement/usecases/mail"
	"github.com/savannahghi/engagementcore/pkg/engagement/usecases/messaging"
	"github.com/savannahghi/engagementcore/pkg/engagement/usecases/onboarding"
	"github.com/savannahghi/engagementcore/pkg/engagement/usecases/otp"
	"github.com/savannahghi/engagementcore/pkg/engagement/usecases/sms"
	"github.com/savannahghi/engagementcore/pkg/engagement/usecases/surveys"
	"github.com/savannahghi/engagementcore/pkg/engagement/usecases/uploads"
)

// Interactor is an implementation of the usecases interface
type Interactor struct {
	*feed.UseCaseImpl
	*feed.NotificationImpl
	*fcm.ImplFCM
	*library.ImplLibrary
	*mail.ImplMail
	*messaging.ImplNotification
	*onboarding.ImplOnboarding
	*otp.ImplOTP
	*sms.ImplSMS
	*surveys.ImplSurveys
	*uploads.ImpUploads
}

// NewUsecasesInteractor initializes a new usecases interactor
func NewUsecasesInteractor(infrastructure infrastructure.Interactor) Interactor {

	notification := feed.NewNotification(infrastructure)
	feed := feed.NewFeed(infrastructure)
	fcm := fcm.NewFCM(infrastructure)
	library := library.NewLibrary(infrastructure)
	mail := mail.NewMail(infrastructure)
	messaging := messaging.NewNotification(infrastructure)
	onboarding := onboarding.NewOnboarding(infrastructure)
	otp := otp.NewOTP(infrastructure)
	sms := sms.NewSMS(infrastructure)
	surveys := surveys.NewSurveys(infrastructure)
	uploads := uploads.NewUploads(infrastructure)

	return Interactor{
		feed,
		notification,
		fcm,
		library,
		mail,
		messaging,
		onboarding,
		otp,
		sms,
		surveys,
		uploads,
	}
}
