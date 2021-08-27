package library

import (
	"context"

	"github.com/savannahghi/engagement/pkg/engagement/domain"
	"github.com/savannahghi/engagement/pkg/engagement/infrastructure"
	"github.com/savannahghi/feedlib"
)

type LibraryUsecases interface {
	GetLibraryContent(ctx context.Context) ([]*domain.GhostCMSPost, error)
	GetFaqsContent(ctx context.Context, flavour feedlib.Flavour) ([]*domain.GhostCMSPost, error)
}

type LibraryUsecasesImpl struct {
	infrastructure infrastructure.Infrastructure
}

func NewLibraryUsecasesImpl(infrastructure infrastructure.Infrastructure) *LibraryUsecasesImpl {
	return &LibraryUsecasesImpl{infrastructure: infrastructure}
}

func (l *LibraryUsecasesImpl) GetFaqsContent(ctx context.Context, flavour feedlib.Flavour) ([]*domain.GhostCMSPost, error) {
	return l.infrastructure.GetFaqsContent(ctx, flavour)
}

func (l *LibraryUsecasesImpl) GetLibraryContent(ctx context.Context) ([]*domain.GhostCMSPost, error) {
	return l.infrastructure.GetLibraryContent(ctx)
}
