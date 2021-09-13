package uploads_test

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/google/uuid"
	"github.com/savannahghi/engagementcore/pkg/engagement/infrastructure/services/uploads"
	uploadsMock "github.com/savannahghi/engagementcore/pkg/engagement/infrastructure/services/uploads/mock"
	"github.com/savannahghi/profileutils"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
)

var (
	fakeUploadsService uploadsMock.FakeServiceUploads
)

func TestUnit_Upload(t *testing.T) {
	ctx := context.Background()
	var s uploads.ServiceUploads = &fakeUploadsService

	bs, err := ioutil.ReadFile("testdata/gandalf.jpg")
	assert.Nil(t, err)
	sEnc := base64.StdEncoding.EncodeToString(bs)

	inp := profileutils.UploadInput{
		Title:       "Test file from automated tests",
		ContentType: "JPG",
		Language:    "en",
		Filename:    fmt.Sprintf("%s.jpg", uuid.New().String()),
		Base64data:  sEnc,
	}

	type args struct {
		ctx context.Context
		inp profileutils.UploadInput
	}
	tests := []struct {
		name      string
		args      args
		wantValue bool
		wantErr   bool
	}{
		{
			name: "valid case: successfully upload the file to cloud storage",
			args: args{
				ctx: ctx,
				inp: inp,
			},
			wantErr:   false,
			wantValue: true,
		},
		{
			name: "invalid case: upload the file to cloud storage, null upload input",
			args: args{
				ctx: ctx,
			},
			wantErr:   true,
			wantValue: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid case: successfully upload the file to cloud storage" {
				fakeUploadsService.UploadFn = func(
					ctx context.Context,
					inp profileutils.UploadInput,
				) (*profileutils.Upload, error) {
					return &profileutils.Upload{}, nil
				}
			}
			if tt.name == "invalid case: upload the file to cloud storage, null upload input" {
				fakeUploadsService.UploadFn = func(
					ctx context.Context,
					inp profileutils.UploadInput,
				) (*profileutils.Upload, error) {
					return nil, fmt.Errorf("test error")
				}
			}
			got, err := s.Upload(tt.args.ctx, tt.args.inp)
			if (err != nil) != tt.wantErr {
				t.Errorf("ServiceUploadImpl.Upload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantValue && got == nil {
				t.Errorf("ServiceUploadImpl.Upload() expected a value, got %v", got)
			}
		})
	}
}
func TestUnit_ServiceUploadImpl_FindUploadByID(t *testing.T) {
	ctx := context.Background()
	var s uploads.ServiceUploads = &fakeUploadsService

	data := profileutils.Upload{}
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name      string
		args      args
		wantValue bool
		wantErr   bool
	}{
		{
			name: "valid: successfully find upload by id",
			args: args{
				ctx: ctx,
				id:  data.ID,
			},
			wantValue: true,
			wantErr:   false,
		},
		{
			name: "invalid: find upload by id, invalid id",
			args: args{
				ctx: ctx,
				id:  ksuid.New().String(),
			},
			wantValue: false,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid: successfully find upload by id" {
				fakeUploadsService.FindUploadByIDFn = func(
					ctx context.Context,
					id string,
				) (*profileutils.Upload, error) {
					return &profileutils.Upload{}, nil
				}
			}
			if tt.name == "invalid: find upload by id, invalid id" {
				fakeUploadsService.FindUploadByIDFn = func(
					ctx context.Context,
					id string,
				) (*profileutils.Upload, error) {
					return nil, fmt.Errorf("unable to retrieve upload ")
				}
			}
			got, err := s.FindUploadByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("ServiceUploadImpl.FindUploadByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantValue && got == nil {
				t.Errorf("ServiceUploadImpl.FindUploadByID() expected a value, got %v", got)
			}
		})
	}
}
