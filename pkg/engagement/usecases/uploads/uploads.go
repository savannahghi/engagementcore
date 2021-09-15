package uploads

import (
	"context"

	"github.com/savannahghi/engagementcore/pkg/engagement/infrastructure"
	"github.com/savannahghi/profileutils"
)

// UsecaseUploads defines uploads service usecases interface
type UsecaseUploads interface {
	Upload(
		ctx context.Context,
		inp profileutils.UploadInput,
	) (*profileutils.Upload, error)

	FindUploadByID(
		ctx context.Context,
		id string,
	) (*profileutils.Upload, error)
}

// ImpUploads is the uploads service implementation
type ImpUploads struct {
	infrastructure infrastructure.Interactor
}

// NewUploads initializes an upload service instance
func NewUploads(infrastructure infrastructure.Interactor) *ImpUploads {
	return &ImpUploads{
		infrastructure: infrastructure,
	}
}

// Upload uploads the file to cloud storage
func (f *ImpUploads) Upload(
	ctx context.Context,
	inp profileutils.UploadInput,
) (*profileutils.Upload, error) {
	i := f.infrastructure.ServiceUploadImpl
	return i.Upload(
		ctx,
		inp,
	)
}

// FindUploadByID retrieves an upload by it's ID
func (f *ImpUploads) FindUploadByID(
	ctx context.Context,
	id string,
) (*profileutils.Upload, error) {
	i := f.infrastructure.ServiceUploadImpl
	return i.FindUploadByID(
		ctx,
		id,
	)
}
