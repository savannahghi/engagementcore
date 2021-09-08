package library

import (
	"context"

	"github.com/savannahghi/engagementcore/pkg/engagement/domain"
	"github.com/savannahghi/engagementcore/pkg/engagement/infrastructure"
	"github.com/savannahghi/feedlib"
)

// UsecaseLibrary defines Libray service usecases interface
type UsecaseLibrary interface {
	GetFeedContent(
		ctx context.Context,
		flavour feedlib.Flavour,
	) ([]*domain.GhostCMSPost, error)
	GetFaqsContent(
		ctx context.Context,
		flavour feedlib.Flavour,
	) ([]*domain.GhostCMSPost, error)
	GetLibraryContent(
		ctx context.Context,
	) ([]*domain.GhostCMSPost, error)
}

// ImplLibrary is the library service implementation
type ImplLibrary struct {
	infrastructure infrastructure.Interactor
}

// NewLibrary initializes a Library service instance
func NewLibrary(infrastructure infrastructure.Interactor) *ImplLibrary {
	return &ImplLibrary{
		infrastructure: infrastructure,
	}
}

// GetFeedContent gets feed content from ghost cms
func (l *ImplLibrary) GetFeedContent(
	ctx context.Context,
	flavour feedlib.Flavour,
) ([]*domain.GhostCMSPost, error) {
	i := l.infrastructure.ServiceLibraryImpl
	return i.GetFeedContent(ctx, flavour)
}

// GetFaqsContent gets FAQ content from ghost cms
func (l *ImplLibrary) GetFaqsContent(
	ctx context.Context,
	flavour feedlib.Flavour,
) ([]*domain.GhostCMSPost, error) {
	i := l.infrastructure.ServiceLibraryImpl
	return i.GetFaqsContent(ctx, flavour)

}

// GetLibraryContent gets Library content from ghost cms
func (l *ImplLibrary) GetLibraryContent(
	ctx context.Context,
) ([]*domain.GhostCMSPost, error) {
	i := l.infrastructure.ServiceLibraryImpl
	return i.GetLibraryContent(ctx)
}
