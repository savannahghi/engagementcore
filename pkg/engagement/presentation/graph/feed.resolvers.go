package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"time"

	"github.com/savannahghi/engagementcore/pkg/engagement/application/common/helpers"
	"github.com/savannahghi/engagementcore/pkg/engagement/domain"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/serverutils"
)

// ResolveFeedItem is the resolver for the resolveFeedItem field.
func (r *mutationResolver) ResolveFeedItem(ctx context.Context, flavour feedlib.Flavour, itemID string) (*feedlib.Item, error) {
	startTime := time.Now()

	uid, err := r.getLoggedInUserUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't get logged in user UID")
	}
	item, err := r.usecases.ResolveFeedItem(ctx, uid, flavour, itemID)
	if err != nil {
		return nil, fmt.Errorf("unable to resolve a Feed item: %v", err)
	}

	defer serverutils.RecordGraphqlResolverMetrics(ctx, startTime, "resolveFeedItem", err)

	return item, nil
}

// UnresolveFeedItem is the resolver for the unresolveFeedItem field.
func (r *mutationResolver) UnresolveFeedItem(ctx context.Context, flavour feedlib.Flavour, itemID string) (*feedlib.Item, error) {
	startTime := time.Now()

	uid, err := r.getLoggedInUserUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't get logged in user UID")
	}
	item, err := r.usecases.UnresolveFeedItem(ctx, uid, flavour, itemID)
	if err != nil {
		return nil, fmt.Errorf("unable to unresolve Feed item: %v", err)
	}

	defer serverutils.RecordGraphqlResolverMetrics(ctx, startTime, "unresolveFeedItem", err)

	return item, nil
}

// PinFeedItem is the resolver for the pinFeedItem field.
func (r *mutationResolver) PinFeedItem(ctx context.Context, flavour feedlib.Flavour, itemID string) (*feedlib.Item, error) {
	startTime := time.Now()

	uid, err := r.getLoggedInUserUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't get logged in user UID")
	}
	item, err := r.usecases.PinFeedItem(ctx, uid, flavour, itemID)
	if err != nil {
		return nil, fmt.Errorf("unable to pin Feed item: %v", err)
	}

	defer serverutils.RecordGraphqlResolverMetrics(ctx, startTime, "pinFeedItem", err)

	return item, nil
}

// UnpinFeedItem is the resolver for the unpinFeedItem field.
func (r *mutationResolver) UnpinFeedItem(ctx context.Context, flavour feedlib.Flavour, itemID string) (*feedlib.Item, error) {
	startTime := time.Now()

	uid, err := r.getLoggedInUserUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't get logged in user UID")
	}
	item, err := r.usecases.UnpinFeedItem(ctx, uid, flavour, itemID)
	if err != nil {
		return nil, fmt.Errorf("unable to unpin Feed item: %v", err)
	}

	defer serverutils.RecordGraphqlResolverMetrics(ctx, startTime, "unpinFeedItem", err)

	return item, nil
}

// HideFeedItem is the resolver for the hideFeedItem field.
func (r *mutationResolver) HideFeedItem(ctx context.Context, flavour feedlib.Flavour, itemID string) (*feedlib.Item, error) {
	startTime := time.Now()

	uid, err := r.getLoggedInUserUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't get logged in user UID")
	}
	item, err := r.usecases.HideFeedItem(ctx, uid, flavour, itemID)
	if err != nil {
		return nil, fmt.Errorf("unable to hide Feed item: %v", err)
	}

	defer serverutils.RecordGraphqlResolverMetrics(ctx, startTime, "hideFeedItem", err)

	return item, nil
}

// ShowFeedItem is the resolver for the showFeedItem field.
func (r *mutationResolver) ShowFeedItem(ctx context.Context, flavour feedlib.Flavour, itemID string) (*feedlib.Item, error) {
	startTime := time.Now()

	uid, err := r.getLoggedInUserUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't get logged in user UID")
	}
	item, err := r.usecases.ShowFeedItem(ctx, uid, flavour, itemID)
	if err != nil {
		return nil, fmt.Errorf("unable to show Feed item: %v", err)
	}

	defer serverutils.RecordGraphqlResolverMetrics(ctx, startTime, "showFeedItem", err)

	return item, nil
}

// HideNudge is the resolver for the hideNudge field.
func (r *mutationResolver) HideNudge(ctx context.Context, flavour feedlib.Flavour, nudgeID string) (*feedlib.Nudge, error) {
	startTime := time.Now()

	uid, err := r.getLoggedInUserUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't get logged in user UID")
	}
	nudge, err := r.usecases.HideNudge(ctx, uid, flavour, nudgeID)
	if err != nil {
		return nil, fmt.Errorf("unable to hide nudge: %v", err)
	}

	defer serverutils.RecordGraphqlResolverMetrics(ctx, startTime, "hideNudge", err)

	return nudge, nil
}

// ShowNudge is the resolver for the showNudge field.
func (r *mutationResolver) ShowNudge(ctx context.Context, flavour feedlib.Flavour, nudgeID string) (*feedlib.Nudge, error) {
	startTime := time.Now()

	uid, err := r.getLoggedInUserUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't get logged in user UID")
	}
	nudge, err := r.usecases.ShowNudge(ctx, uid, flavour, nudgeID)
	if err != nil {
		return nil, fmt.Errorf("unable to show nudge: %v", err)
	}

	defer serverutils.RecordGraphqlResolverMetrics(ctx, startTime, "showNudge", err)

	return nudge, nil
}

// PostMessage is the resolver for the postMessage field.
func (r *mutationResolver) PostMessage(ctx context.Context, flavour feedlib.Flavour, itemID string, message feedlib.Message) (*feedlib.Message, error) {
	startTime := time.Now()

	uid, err := r.getLoggedInUserUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't get logged in user UID")
	}
	msg, err := r.usecases.PostMessage(ctx, uid, flavour, itemID, &message)
	if err != nil {
		return nil, fmt.Errorf("unable to post a message: %v", err)
	}

	defer serverutils.RecordGraphqlResolverMetrics(ctx, startTime, "postMessage", err)

	return msg, nil
}

// DeleteMessage is the resolver for the deleteMessage field.
func (r *mutationResolver) DeleteMessage(ctx context.Context, flavour feedlib.Flavour, itemID string, messageID string) (bool, error) {
	startTime := time.Now()

	uid, err := r.getLoggedInUserUID(ctx)
	if err != nil {
		return false, fmt.Errorf("can't get logged in user UID")
	}
	err = r.usecases.DeleteMessage(ctx, uid, flavour, itemID, messageID)
	if err != nil {
		return false, fmt.Errorf("can't delete message: %w", err)
	}

	defer serverutils.RecordGraphqlResolverMetrics(ctx, startTime, "deleteMessage", err)

	return true, nil
}

// ProcessEvent is the resolver for the processEvent field.
func (r *mutationResolver) ProcessEvent(ctx context.Context, flavour feedlib.Flavour, event feedlib.Event) (bool, error) {
	startTime := time.Now()

	uid, err := r.getLoggedInUserUID(ctx)
	if err != nil {
		return false, fmt.Errorf("can't get logged in user UID")
	}
	err = r.usecases.ProcessEvent(ctx, uid, flavour, &event)
	if err != nil {
		return false, fmt.Errorf("can't process event: %w", err)
	}

	defer serverutils.RecordGraphqlResolverMetrics(ctx, startTime, "processEvent", err)

	return true, nil
}

// GetFeed is the resolver for the getFeed field.
func (r *queryResolver) GetFeed(ctx context.Context, flavour feedlib.Flavour, playMp4 *bool, isAnonymous bool, persistent feedlib.BooleanFilter, status *feedlib.Status, visibility *feedlib.Visibility, expired *feedlib.BooleanFilter, filterParams *helpers.FilterParams) (*domain.Feed, error) {
	startTime := time.Now()
	uid, err := r.getLoggedInUserUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't get logged in user UID")
	}
	var playmp4FeedVideos bool
	if playMp4 != nil {
		playmp4FeedVideos = true
	}
	feed, err := r.usecases.GetFeed(
		ctx,
		&uid,
		&isAnonymous,
		flavour,
		playmp4FeedVideos,
		persistent,
		status,
		visibility,
		expired,
		filterParams,
	)
	if err != nil {
		return nil, fmt.Errorf("can't get Feed: %w", err)
	}

	defer serverutils.RecordGraphqlResolverMetrics(ctx, startTime, "getFeed", err)

	return feed, nil
}

// Labels is the resolver for the labels field.
func (r *queryResolver) Labels(ctx context.Context, flavour feedlib.Flavour) ([]string, error) {
	startTime := time.Now()

	uid, err := r.getLoggedInUserUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't get logged in user UID")
	}
	labels, err := r.usecases.Labels(ctx, uid, flavour)
	if err != nil {
		return nil, fmt.Errorf("unable to get Labels count: %v", err)
	}

	defer serverutils.RecordGraphqlResolverMetrics(ctx, startTime, "labels", err)

	return labels, nil
}

// UnreadPersistentItems is the resolver for the unreadPersistentItems field.
func (r *queryResolver) UnreadPersistentItems(ctx context.Context, flavour feedlib.Flavour) (int, error) {
	startTime := time.Now()

	uid, err := r.getLoggedInUserUID(ctx)
	if err != nil {
		return -1, fmt.Errorf("can't get logged in user UID")
	}
	count, err := r.usecases.UnreadPersistentItems(ctx, uid, flavour)
	if err != nil {
		return -1, fmt.Errorf("unable to count unread persistent items: %v", err)
	}

	defer serverutils.RecordGraphqlResolverMetrics(ctx, startTime, "unreadPersistentItems", err)

	return count, nil
}
