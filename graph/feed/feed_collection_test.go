package feed_test

import (
	"context"
	"testing"

	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/engagement/graph/feed"
	db "gitlab.slade360emr.com/go/engagement/graph/feed/infrastructure/database"
	"gitlab.slade360emr.com/go/engagement/graph/feed/infrastructure/messaging"
)

func getNotificationService(ctx context.Context, t *testing.T) feed.NotificationService {
	projectID, err := base.GetEnvVar(base.GoogleCloudProjectIDEnvVarName)
	if err != nil {
		t.Errorf("project ID not found in env var: %s", err)
		return nil
	}

	if projectID == "" {
		t.Errorf("nil project ID")
		return nil
	}

	ns, err := messaging.NewPubSubNotificationService(ctx, projectID)
	if err != nil {
		t.Errorf("can't instantiate notification service: %s", err)
		return nil
	}

	if ns == nil {
		t.Errorf("nil notification service")
		return nil
	}

	return ns
}

func getFeedAggregate(t *testing.T) *feed.Collection {
	ctx := context.Background()

	fr, err := db.NewFirebaseRepository(ctx)
	if err != nil {
		t.Errorf("project ID not found in env var: %s", err)
		return nil
	}

	if fr == nil {
		t.Errorf("nil firebase repository")
		return nil
	}

	ns := getNotificationService(ctx, t)
	if ns == nil {
		t.Errorf("nil notification service")
		return nil
	}

	agg, err := feed.NewCollection(fr, ns)
	assert.Nil(t, err)
	assert.NotNil(t, agg)

	return agg
}

func getTestFeedAggregate(t *testing.T) *feed.Collection {
	ctx := context.Background()
	repository, err := db.NewFirebaseRepository(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, repository)

	notificationService, err := messaging.NewMockNotificationService()
	assert.Nil(t, err)
	assert.NotNil(t, notificationService)

	agg, err := feed.NewCollection(
		repository,
		notificationService,
	)
	assert.NotNil(t, agg)
	assert.Nil(t, err)

	return agg
}

func TestNewAggregate(t *testing.T) {
	ctx := context.Background()

	fr, err := db.NewFirebaseRepository(ctx)
	if err != nil {
		t.Errorf("can't initialize firebase repository: %s", err)
		return
	}

	if fr == nil {
		t.Errorf("nil firebase repository")
		return
	}

	projectID, err := base.GetEnvVar(base.GoogleCloudProjectIDEnvVarName)
	if err != nil {
		t.Errorf("project ID not found in env var: %s", err)
		return
	}

	if projectID == "" {
		t.Errorf("nil project ID")
		return
	}

	ns, err := messaging.NewPubSubNotificationService(ctx, projectID)
	if err != nil {
		t.Errorf("can't initialize notification service: %s", err)
		return
	}
	if ns == nil {
		t.Errorf("nil notification service")
		return
	}

	type args struct {
		repository          feed.Repository
		notificationService feed.NotificationService
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case",
			args: args{
				repository:          fr,
				notificationService: ns,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := feed.NewCollection(
				tt.args.repository, tt.args.notificationService)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"NewAggregate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.NotNil(t, got)
			}
		})
	}
}

func TestAggregate_GetThinFeed(t *testing.T) {
	ctx := context.Background()

	fr, err := db.NewFirebaseRepository(ctx)
	if err != nil {
		t.Errorf("can't initialize firebase repository: %s", err)
		return
	}

	if fr == nil {
		t.Errorf("nil firebase repository")
		return
	}

	projectID, err := base.GetEnvVar(base.GoogleCloudProjectIDEnvVarName)
	if err != nil {
		t.Errorf("project ID not found in env var: %s", err)
		return
	}

	if projectID == "" {
		t.Errorf("nil project ID")
		return
	}

	ns, err := messaging.NewPubSubNotificationService(ctx, projectID)
	if err != nil {
		t.Errorf("can't initialize notification service: %s", err)
		return
	}
	if ns == nil {
		t.Errorf("nil notification service")
		return
	}

	agg, err := feed.NewCollection(fr, ns)
	if err != nil {
		t.Errorf("can't initialize aggregate: %s", err)
		return
	}

	if agg == nil {
		t.Errorf("nil feed collection/aggregate")
		return
	}

	type args struct {
		ctx         context.Context
		uid         string
		isAnonymous bool
		flavour     base.Flavour
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "default case - success",
			args: args{
				ctx:         ctx,
				uid:         ksuid.New().String(),
				isAnonymous: false,
				flavour:     base.FlavourConsumer,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := agg.GetThinFeed(
				tt.args.ctx,
				&tt.args.uid,
				&tt.args.isAnonymous,
				tt.args.flavour,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"Aggregate.GetThinFeed() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if !tt.wantErr {
				assert.NotNil(t, got)
			}
		})
	}
}

func TestAggregate_GetFeed(t *testing.T) {
	ctx := context.Background()
	fr, err := db.NewFirebaseRepository(ctx)
	if err != nil {
		t.Errorf("can't initialize firebase repository: %s", err)
		return
	}

	if fr == nil {
		t.Errorf("nil firebase repository")
		return
	}

	projectID, err := base.GetEnvVar(base.GoogleCloudProjectIDEnvVarName)
	if err != nil {
		t.Errorf("project ID not found in env var: %s", err)
		return
	}

	if projectID == "" {
		t.Errorf("nil project ID")
		return
	}

	ns, err := messaging.NewPubSubNotificationService(ctx, projectID)
	if err != nil {
		t.Errorf("can't initialize notification service: %s", err)
		return
	}
	if ns == nil {
		t.Errorf("nil notification service")
		return
	}

	agg, err := feed.NewCollection(fr, ns)
	if err != nil {
		t.Errorf("can't initialize aggregate: %s", err)
		return
	}

	if agg == nil {
		t.Errorf("nil feed collection/aggregate")
		return
	}

	uid := ksuid.New().String()
	flavour := base.FlavourConsumer
	persistent := base.BooleanFilterBoth
	status := base.StatusPending
	visibility := base.VisibilityHide
	expired := base.BooleanFilterFalse

	type args struct {
		ctx          context.Context
		uid          string
		isAnonymous  bool
		flavour      base.Flavour
		persistent   base.BooleanFilter
		status       base.Status
		visibility   base.Visibility
		expired      base.BooleanFilter
		filterParams *feed.FilterParams
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "no filter params",
			args: args{
				ctx:          ctx,
				uid:          uid,
				isAnonymous:  false,
				flavour:      flavour,
				persistent:   persistent,
				status:       status,
				visibility:   visibility,
				expired:      expired,
				filterParams: nil,
			},
			wantErr: false,
		},
		{
			name: "with filter params",
			args: args{
				ctx:         ctx,
				uid:         uid,
				isAnonymous: false,
				flavour:     flavour,
				persistent:  persistent,
				status:      status,
				visibility:  visibility,
				expired:     expired,
				filterParams: &feed.FilterParams{
					Labels: []string{ksuid.New().String()},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := agg.GetFeed(
				tt.args.ctx,
				&tt.args.uid,
				&tt.args.isAnonymous,
				tt.args.flavour,
				tt.args.persistent,
				&tt.args.status,
				&tt.args.visibility,
				&tt.args.expired,
				tt.args.filterParams,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"Aggregate.GetFeed() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if !tt.wantErr {
				assert.NotNil(t, got)
			}
		})
	}
}
