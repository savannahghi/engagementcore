package mock

import (
	"context"

	"github.com/savannahghi/profileutils"
)

// FakeServiceUploads is an defines all interactions with the mock Uploads service
type FakeServiceUploads struct {
	UploadFn func(
		ctx context.Context,
		inp profileutils.UploadInput,
	) (*profileutils.Upload, error)

	FindUploadByIDFn func(
		ctx context.Context,
		id string,
	) (*profileutils.Upload, error)
}

// Upload is a mock of the Upload service
func (u *FakeServiceUploads) Upload(
	ctx context.Context,
	inp profileutils.UploadInput,
) (*profileutils.Upload, error) {
	return u.UploadFn(ctx, inp)
}

// FindUploadByID is a mock of the FindUploadByID service
func (u *FakeServiceUploads) FindUploadByID(
	ctx context.Context,
	id string,
) (*profileutils.Upload, error) {
	return u.FindUploadByIDFn(ctx, id)
}
