package uploads

import (
	"context"

	"github.com/savannahghi/engagement/pkg/engagement/infrastructure"
	"github.com/savannahghi/profileutils"
)

// UseCases ...
type UseCases interface {
	FindUploadByID(ctx context.Context, id string) (*profileutils.Upload, error)
	Upload(ctx context.Context, input profileutils.UploadInput) (*profileutils.Upload, error)
}

// UseCasesImpl ...
type UseCasesImpl struct {
	infrastructure infrastructure.Infrastructure
}

// NewUploadsImpl ...
func NewUploadsImpl(infrastructure infrastructure.Infrastructure) *UseCasesImpl {
	return &UseCasesImpl{infrastructure: infrastructure}
}

// FindUploadByID ...
func (u *UseCasesImpl) FindUploadByID(ctx context.Context, id string) (*profileutils.Upload, error) {
	return u.infrastructure.FindUploadByID(ctx, id)
}

// Upload ...
func (u *UseCasesImpl) Upload(ctx context.Context, input profileutils.UploadInput) (*profileutils.Upload, error) {
	return u.infrastructure.Upload(ctx, input)
}
