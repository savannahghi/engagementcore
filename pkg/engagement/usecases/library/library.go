package library

import (
	"context"

	"github.com/savannahghi/engagement/pkg/engagement/domain"
	"github.com/savannahghi/engagement/pkg/engagement/infrastructure"
	"github.com/savannahghi/feedlib"
)

// UseCases ...
type UseCases interface {
	GetLibraryContent(ctx context.Context) ([]*domain.GhostCMSPost, error)
	GetFaqsContent(ctx context.Context, flavour feedlib.Flavour) ([]*domain.GhostCMSPost, error)
}

// UseCasesImpl ...
type UseCasesImpl struct {
	infrastructure infrastructure.Infrastructure
}

// NewLibraryUsecasesImpl ...
func NewLibraryUsecasesImpl(infrastructure infrastructure.Infrastructure) *UseCasesImpl {
	return &UseCasesImpl{infrastructure: infrastructure}
}

// GetFaqsContent ...
func (l *UseCasesImpl) GetFaqsContent(ctx context.Context, flavour feedlib.Flavour) ([]*domain.GhostCMSPost, error) {
	return l.infrastructure.GetFaqsContent(ctx, flavour)
}

// GetLibraryContent ...
func (l *UseCasesImpl) GetLibraryContent(ctx context.Context) ([]*domain.GhostCMSPost, error) {
	return l.infrastructure.GetLibraryContent(ctx)
}
