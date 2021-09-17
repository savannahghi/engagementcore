package library_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/savannahghi/engagementcore/pkg/engagement/domain"
	"github.com/savannahghi/engagementcore/pkg/engagement/infrastructure/services/library"
	libraryMock "github.com/savannahghi/engagementcore/pkg/engagement/infrastructure/services/library/mock"
	"github.com/savannahghi/feedlib"
	"github.com/segmentio/ksuid"
)

var (
	fakeLibraryService libraryMock.FakeLibrary
)

func getGhostCMSPost() []*domain.GhostCMSPost {
	var (
		standardTime = time.Date(2021, 2, 3, 4, 4, 4, 5, time.Local)
		testText     = "test text"
	)
	return []*domain.GhostCMSPost{
		{
			ID:           ksuid.New().String(),
			UUID:         ksuid.New().String(),
			Slug:         "slug",
			Title:        "title",
			HTML:         "<html></html>",
			Excerpt:      "test",
			URL:          "https://example.com",
			FeatureImage: "https://example.com/test.jepeg",
			Featured:     true,
			Visibility:   "published",
			ReadingTime:  30,
			CreatedAt:    standardTime,
			UpdatedAt:    standardTime,
			PublishedAt:  standardTime,
			CommentID:    ksuid.New().String(),
			Tags: []domain.GhostCMSTag{
				{
					ID:          ksuid.New().String(),
					Name:        "test",
					Slug:        "slug",
					Description: &testText,
					Visibility:  "published",
					URL:         "https://example.com",
				},
			},
			Authors: []domain.GhostCMSAuthor{
				{
					ID:           ksuid.New().String(),
					Name:         "test",
					Slug:         "slug",
					ProfileImage: "https://example.com/test.jepeg",
					Website:      "https://example.com",
					Location:     "test",
					Facebook:     "https://facebook.com",
					Twitter:      "https://twitter.com",
					URL:          "https://example.com",
				},
			},
			PrimaryAuthor: domain.GhostCMSAuthor{
				ID:           ksuid.New().String(),
				Name:         "test",
				Slug:         "slug",
				ProfileImage: "https://example.com/test.jepeg",
				Website:      "https://example.com",
				Location:     "test",
				Facebook:     "https://facebook.com",
				Twitter:      "https://twitter.com",
				URL:          "https://example.com",
			},
			PrimaryTag: domain.GhostCMSTag{
				ID:   ksuid.New().String(),
				Name: "test",
				Slug: "slug",
				URL:  "https://example.com",
			},
		},
	}

}

func TestUnit_GetFeedContent(t *testing.T) {
	var s library.ServiceLibrary = &fakeLibraryService

	ctx := context.Background()
	flavor := feedlib.FlavourPro

	post := getGhostCMSPost()

	type args struct {
		ctx     context.Context
		flavour feedlib.Flavour
	}
	tests := []struct {
		name    string
		args    args
		want    []*domain.GhostCMSPost
		wantErr bool
	}{
		{
			name: "valid: correct params passed",
			args: args{
				ctx:     ctx,
				flavour: flavor,
			},
			want:    post,
			wantErr: false,
		},
		{
			name: "invalid: no flavor passed",
			args: args{
				ctx: ctx,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid: correct params passed" {
				fakeLibraryService.GetFeedContentFn = func(
					ctx context.Context,
					flavour feedlib.Flavour,
				) ([]*domain.GhostCMSPost, error) {
					return post, nil
				}
			}
			if tt.name == "invalid: no flavor passed" {
				fakeLibraryService.GetFeedContentFn = func(
					ctx context.Context,
					flavour feedlib.Flavour,
				) ([]*domain.GhostCMSPost, error) {
					return nil, fmt.Errorf("test error")
				}
			}
			got, err := s.GetFeedContent(tt.args.ctx, tt.args.flavour)
			if (err != nil) != tt.wantErr {
				t.Errorf("ServiceLibraryImpl.GetFeedContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ServiceLibraryImpl.GetFeedContent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnit_GetFaqsContent(t *testing.T) {
	var s library.ServiceLibrary = &fakeLibraryService

	ctx := context.Background()
	flavour := feedlib.FlavourPro

	post := getGhostCMSPost()

	type args struct {
		ctx     context.Context
		flavour feedlib.Flavour
	}
	tests := []struct {
		name    string
		args    args
		want    []*domain.GhostCMSPost
		wantErr bool
	}{
		{
			name: "valid: correct params passed",
			args: args{
				ctx:     ctx,
				flavour: flavour,
			},
			want:    post,
			wantErr: false,
		},
		{
			name: "invalid: no flavor passed",
			args: args{
				ctx: ctx,
			},
			want:    nil,
			wantErr: true,
		},

		{
			name: "invalid: bad flavor passed",
			args: args{
				ctx:     ctx,
				flavour: "BAD FLAVOUR",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid: correct params passed" {
				fakeLibraryService.GetFaqsContentFn = func(
					ctx context.Context,
					flavour feedlib.Flavour,
				) ([]*domain.GhostCMSPost, error) {
					return post, nil
				}
			}
			if tt.name == "invalid: no flavor passed" {
				fakeLibraryService.GetFaqsContentFn = func(
					ctx context.Context,
					flavour feedlib.Flavour,
				) ([]*domain.GhostCMSPost, error) {
					return nil, fmt.Errorf("test error")
				}
			}
			got, err := s.GetFaqsContent(tt.args.ctx, tt.args.flavour)
			if (err != nil) != tt.wantErr {
				t.Errorf("ServiceLibraryImpl.GetFaqsContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ServiceLibraryImpl.GetFaqsContent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnit_GetLibraryContent(t *testing.T) {
	var s library.ServiceLibrary = &fakeLibraryService

	ctx := context.Background()

	post := getGhostCMSPost()

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    []*domain.GhostCMSPost
		wantErr bool
	}{
		{
			name: "valid: correct params passed",
			args: args{
				ctx: ctx,
			},
			want:    post,
			wantErr: false,
		},
		{
			name: "invalid: empty args",
			args: args{
				ctx: ctx,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid: correct params passed" {
				fakeLibraryService.GetLibraryContentFn = func(
					ctx context.Context,
				) ([]*domain.GhostCMSPost, error) {
					return post, nil
				}
			}
			if tt.name == "invalid: empty args" {
				fakeLibraryService.GetLibraryContentFn = func(
					ctx context.Context,
				) ([]*domain.GhostCMSPost, error) {
					return nil, fmt.Errorf("test error")
				}
			}
			got, err := s.GetLibraryContent(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("ServiceLibraryImpl.GetLibraryContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ServiceLibraryImpl.GetLibraryContent() = %v, want %v", got, tt.want)
			}
		})
	}
}
