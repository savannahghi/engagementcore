package usecases

import (
	"github.com/savannahghi/engagementcore/pkg/engagement/infrastructure"
	"github.com/savannahghi/engagementcore/pkg/engagement/usecases/fcm"
	"github.com/savannahghi/engagementcore/pkg/engagement/usecases/feed"
	"github.com/savannahghi/engagementcore/pkg/engagement/usecases/library"
	"github.com/savannahghi/engagementcore/pkg/engagement/usecases/mail"
)

// Interactor is an implementation of the usecases interface
type Interactor struct {
	*feed.UseCaseImpl
	*feed.NotificationImpl
	*fcm.ImplFCM
	*library.ImplLibrary
	*mail.ImplMail
}

// NewUsecasesInteractor initializes a new usecases interactor
func NewUsecasesInteractor(infrastructure infrastructure.Interactor) Interactor {

	notification := feed.NewNotification(infrastructure)
	feed := feed.NewFeed(infrastructure)
	fcm := fcm.NewFCM(infrastructure)
	library := library.NewLibrary(infrastructure)
	mail := mail.NewMail(infrastructure)

	return Interactor{
		feed,
		notification,
		fcm,
		library,
		mail,
	}
}
