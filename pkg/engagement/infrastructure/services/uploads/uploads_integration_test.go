package uploads_test

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/savannahghi/engagementcore/pkg/engagement/infrastructure/services/uploads"
	"github.com/savannahghi/profileutils"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	os.Setenv("ROOT_COLLECTION_SUFFIX", "testing")
	m.Run()
}

func TestNewUploadsService(t *testing.T) {
	tests := []struct {
		name      string
		wantValue bool
	}{
		{
			name:      "default case",
			wantValue: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := uploads.NewUploadsService()
			if tt.wantValue && got == nil {
				t.Errorf("NewUploadsService() = expected a response")
			}
		})
	}
}

func TestUpload(t *testing.T) {
	ctx := context.Background()

	bs, err := ioutil.ReadFile("testdata/gandalf.jpg")
	assert.Nil(t, err)
	sEnc := base64.StdEncoding.EncodeToString(bs)

	service := uploads.NewUploadsService()
	tests := map[string]struct {
		inp                  profileutils.UploadInput
		expectError          bool
		expectedErrorMessage string
	}{
		"simple_case": {
			inp: profileutils.UploadInput{
				Title:       "Test file from automated tests",
				ContentType: "JPG",
				Language:    "en",
				Filename:    fmt.Sprintf("%s.jpg", uuid.New().String()),
				Base64data:  sEnc,
			},
			expectError: false,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			upload, err := service.Upload(ctx, tc.inp)
			if tc.expectError {
				assert.NotNil(t, err)
				assert.Equal(t, tc.expectedErrorMessage, err.Error())
			}
			if !tc.expectError {
				assert.NotNil(t, upload)
				assert.Nil(t, err)

				assert.NotZero(t, upload.ID)
				assert.NotZero(t, upload.URL)
				assert.NotZero(t, upload.Size)
				assert.NotZero(t, upload.Hash)
				assert.NotZero(t, upload.Creation)
				assert.NotZero(t, upload.Title)
				assert.NotZero(t, upload.ContentType)
				assert.NotZero(t, upload.Language)
				assert.NotZero(t, upload.Base64data)
			}
		})
	}
}

func TestServiceUploadImpl_FindUploadByID(t *testing.T) {
	ctx := context.Background()
	s := uploads.NewUploadsService()

	bs, err := ioutil.ReadFile("testdata/gandalf.jpg")
	assert.Nil(t, err)
	sEnc := base64.StdEncoding.EncodeToString(bs)

	UploadInput := profileutils.UploadInput{
		Title:       "Test file from automated tests",
		ContentType: "JPG",
		Language:    "en",
		Filename:    fmt.Sprintf("%s.jpg", uuid.New().String()),
		Base64data:  sEnc,
	}
	data, err := s.Upload(ctx, UploadInput)
	if err != nil {
		t.Errorf("Upload failed: %v", err)
	}

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
			name: "valid: correct id passed",
			args: args{
				ctx: ctx,
				id:  data.ID,
			},
			wantValue: true,
			wantErr:   false,
		},
		{
			name: "invalid: incorrect id passed",
			args: args{
				ctx: ctx,
				id:  ksuid.New().String(),
			},
			wantValue: false,
			wantErr:   true,
		},
		{
			name: "invalid: missing params",
			args: args{
				ctx: ctx,
			},
			wantValue: false,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
