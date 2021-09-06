package mock

import (
	"context"
	"fmt"

	"github.com/savannahghi/engagementcore/pkg/engagement/domain"
	"github.com/savannahghi/feedlib"
)

// FakeLibrary defines a mock Libray service
type FakeLibrary struct {
	GetFeedContentFn func(
		ctx context.Context,
		flavour feedlib.Flavour,
	) ([]*domain.GhostCMSPost, error)
	GetFaqsContentFn func(
		ctx context.Context,
		flavour feedlib.Flavour,
	) ([]*domain.GhostCMSPost, error)
	GetLibraryContentFn func(
		ctx context.Context,
	) ([]*domain.GhostCMSPost, error)
}

// GetFeedContent mocks getting feed content from ghost cms
func (l *FakeLibrary) GetFeedContent(
	ctx context.Context,
	flavour feedlib.Flavour,
) ([]*domain.GhostCMSPost, error) {
	fmt.Println(">>>>@here")
	return l.GetFeedContentFn(ctx, flavour)
}

// GetFaqsContent mocks getting FAQ content from ghost cms
func (l *FakeLibrary) GetFaqsContent(
	ctx context.Context,
	flavour feedlib.Flavour,
) ([]*domain.GhostCMSPost, error) {
	return l.GetFaqsContentFn(ctx, flavour)

}

// GetLibraryContent mocks getting Library content from ghost cms
func (l *FakeLibrary) GetLibraryContent(
	ctx context.Context,
) ([]*domain.GhostCMSPost, error) {
	return l.GetLibraryContentFn(ctx)
}
