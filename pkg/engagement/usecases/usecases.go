package usecases

import (
	"github.com/savannahghi/engagement/pkg/engagement/infrastructure"
	"github.com/savannahghi/engagement/pkg/engagement/usecases/feed"
)

// Interactor is an implementation of the usecases interface
type Interactor struct {
	*feed.UseCaseImpl
	*feed.NotificationImpl
}

// NewInteractor initializes a new usecases interactor
func NewInteractor(infrastructure infrastructure.Interactor) Interactor {

	notification := feed.NewNotification(infrastructure)
	feed := feed.NewFeed(infrastructure)

	return Interactor{
		feed,
		notification,
	}
}
