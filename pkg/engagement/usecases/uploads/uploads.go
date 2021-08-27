package uploads

import (
	"context"

	"github.com/savannahghi/engagement/pkg/engagement/infrastructure"
	"github.com/savannahghi/profileutils"
)

type UploadsUsecases interface {
	FindUploadByID(ctx context.Context, id string) (*profileutils.Upload, error)
	Upload(ctx context.Context, input profileutils.UploadInput) (*profileutils.Upload, error)
}

type UploadsImpl struct {
	infrastructure infrastructure.Infrastructure
}

func NewUploadsImpl(infrastructure infrastructure.Infrastructure) *UploadsImpl {
	return &UploadsImpl{infrastructure: infrastructure}
}

func (u *UploadsImpl) FindUploadByID(ctx context.Context, id string) (*profileutils.Upload, error) {
	return u.infrastructure.FindUploadByID(ctx, id)
}

func (u *UploadsImpl) Upload(ctx context.Context, input profileutils.UploadInput) (*profileutils.Upload, error) {
	return u.infrastructure.Upload(ctx, input)
}
