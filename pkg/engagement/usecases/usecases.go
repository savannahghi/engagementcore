package usecases

import (
	"github.com/savannahghi/engagementcore/pkg/engagement/infrastructure"
	"github.com/savannahghi/engagementcore/pkg/engagement/usecases/fcm"
	"github.com/savannahghi/engagementcore/pkg/engagement/usecases/feed"
)

// Interactor is an implementation of the usecases interface
type Interactor struct {
	*feed.UseCaseImpl
	*feed.NotificationImpl
	*fcm.ImplFCM
}

// NewUsecasesInteractor initializes a new usecases interactor
func NewUsecasesInteractor(infrastructure infrastructure.Interactor) Interactor {

	notification := feed.NewNotification(infrastructure)
	feed := feed.NewFeed(infrastructure)
	fcm := fcm.NewFCM(infrastructure)

	return Interactor{
		feed,
		notification,
		fcm,
	}
}
