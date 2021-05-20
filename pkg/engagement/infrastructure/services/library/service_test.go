package library_test

import (
	"context"
	"testing"

	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/engagement/pkg/engagement/infrastructure/services/library"
)

func TestNewLibraryService(t *testing.T) {
	srv := library.NewLibraryService()
	if srv == nil {
		t.Errorf("nil library service")
	}
}

func TestService_GetFeedContent(t *testing.T) {
	s := library.NewLibraryService()
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name        string
		args        args
		wantNonZero bool
		wantErr     bool
	}{
		{
			name: "default case",
			args: args{
				ctx: context.Background(),
			},
			wantNonZero: true,
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.GetFeedContent(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetFeedContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantNonZero && len(got) < 1 {
				t.Errorf("expected a non zero count of posts")
				return
			}
		})
	}
}

func TestService_GetFaqsContent(t *testing.T) {
	s := library.NewLibraryService()
	type args struct {
		ctx     context.Context
		flavour base.Flavour
	}
	tests := []struct {
		name        string
		args        args
		wantNonZero bool
		wantErr     bool
	}{
		{
			name: "default case",
			args: args{
				ctx:     context.Background(),
				flavour: "PRO",
			},
			wantNonZero: true,
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.GetFaqsContent(tt.args.ctx, tt.args.flavour)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetFaqsContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantNonZero && len(got) < 1 {
				t.Errorf("expected a non zero count of posts")
				return
			}
		})
	}
}

func TestService_GetLibraryContent(t *testing.T) {
	s := library.NewLibraryService()
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name        string
		args        args
		wantNonZero bool
		wantErr     bool
	}{
		{
			name: "default case",
			args: args{
				ctx: context.Background(),
			},
			wantNonZero: true,
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.GetLibraryContent(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetLibraryContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantNonZero && len(got) < 1 {
				t.Errorf("expected a non zero count of posts")
				return
			}
		})
	}
}
