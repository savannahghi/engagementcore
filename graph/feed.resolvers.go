package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"gitlab.slade360emr.com/go/feed/graph/feed"
	"gitlab.slade360emr.com/go/feed/graph/generated"
)

func (r *linkResolver) LinkType(ctx context.Context, obj *feed.Link) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ResolveFeedItem(ctx context.Context, flavour feed.Flavour, itemID string) (*feed.Item, error) {
	r.checkPreconditions()

	thinFeed, err := r.getThinFeed(ctx, flavour)
	if err != nil {
		return nil, fmt.Errorf("can't instantiate new feed: %w", err)
	}

	return thinFeed.ResolveFeedItem(ctx, itemID)
}

func (r *mutationResolver) UnresolveFeedItem(ctx context.Context, flavour feed.Flavour, itemID string) (*feed.Item, error) {
	r.checkPreconditions()

	thinFeed, err := r.getThinFeed(ctx, flavour)
	if err != nil {
		return nil, fmt.Errorf("can't instantiate new feed: %w", err)
	}

	return thinFeed.UnresolveFeedItem(ctx, itemID)
}

func (r *mutationResolver) PinFeedItem(ctx context.Context, flavour feed.Flavour, itemID string) (*feed.Item, error) {
	r.checkPreconditions()

	thinFeed, err := r.getThinFeed(ctx, flavour)
	if err != nil {
		return nil, fmt.Errorf("can't instantiate new feed: %w", err)
	}

	return thinFeed.PinFeedItem(ctx, itemID)
}

func (r *mutationResolver) UnpinFeedItem(ctx context.Context, flavour feed.Flavour, itemID string) (*feed.Item, error) {
	r.checkPreconditions()

	thinFeed, err := r.getThinFeed(ctx, flavour)
	if err != nil {
		return nil, fmt.Errorf("can't instantiate new feed: %w", err)
	}

	return thinFeed.UnpinFeedItem(ctx, itemID)
}

func (r *mutationResolver) HideFeedItem(ctx context.Context, flavour feed.Flavour, itemID string) (*feed.Item, error) {
	r.checkPreconditions()

	thinFeed, err := r.getThinFeed(ctx, flavour)
	if err != nil {
		return nil, fmt.Errorf("can't instantiate new feed: %w", err)
	}

	return thinFeed.HideFeedItem(ctx, itemID)
}

func (r *mutationResolver) ShowFeedItem(ctx context.Context, flavour feed.Flavour, itemID string) (*feed.Item, error) {
	r.checkPreconditions()

	thinFeed, err := r.getThinFeed(ctx, flavour)
	if err != nil {
		return nil, fmt.Errorf("can't instantiate new feed: %w", err)
	}

	return thinFeed.ShowFeedItem(ctx, itemID)
}

func (r *mutationResolver) HideNudge(ctx context.Context, flavour feed.Flavour, nudgeID string) (*feed.Nudge, error) {
	r.checkPreconditions()

	thinFeed, err := r.getThinFeed(ctx, flavour)
	if err != nil {
		return nil, fmt.Errorf("can't instantiate new feed: %w", err)
	}

	return thinFeed.HideNudge(ctx, nudgeID)
}

func (r *mutationResolver) ShowNudge(ctx context.Context, flavour feed.Flavour, nudgeID string) (*feed.Nudge, error) {
	r.checkPreconditions()

	thinFeed, err := r.getThinFeed(ctx, flavour)
	if err != nil {
		return nil, fmt.Errorf("can't instantiate new feed: %w", err)
	}

	return thinFeed.ShowNudge(ctx, nudgeID)
}

func (r *mutationResolver) PostMessage(ctx context.Context, flavour feed.Flavour, itemID string, message feed.Message) (*feed.Message, error) {
	r.checkPreconditions()

	thinFeed, err := r.getThinFeed(ctx, flavour)
	if err != nil {
		return nil, fmt.Errorf("can't instantiate new feed: %w", err)
	}

	return thinFeed.PostMessage(ctx, itemID, &message)
}

func (r *mutationResolver) DeleteMessage(ctx context.Context, flavour feed.Flavour, itemID string, messageID string) (bool, error) {
	r.checkPreconditions()

	thinFeed, err := r.getThinFeed(ctx, flavour)
	if err != nil {
		return false, fmt.Errorf("can't instantiate new feed: %w", err)
	}

	err = thinFeed.DeleteMessage(ctx, itemID, messageID)
	if err != nil {
		return false, fmt.Errorf("can't delete message: %w", err)
	}

	return true, nil
}

func (r *mutationResolver) ProcessEvent(ctx context.Context, flavour feed.Flavour, event feed.Event) (bool, error) {
	r.checkPreconditions()

	thinFeed, err := r.getThinFeed(ctx, flavour)
	if err != nil {
		return false, fmt.Errorf("can't instantiate new feed: %w", err)
	}

	err = thinFeed.ProcessEvent(ctx, &event)
	if err != nil {
		return false, fmt.Errorf("can't process event: %w", err)
	}

	return true, nil
}

func (r *queryResolver) GetFeed(ctx context.Context, flavour feed.Flavour, persistent feed.BooleanFilter, status *feed.Status, visibility *feed.Visibility, expired *feed.BooleanFilter, filterParams *feed.FilterParams) (*feed.Feed, error) {
	r.checkPreconditions()

	uid, err := r.getLoggedInUserUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't get logged in user UID")
	}

	agg, err := feed.NewCollection(r.repository, r.notificationService)
	if err != nil {
		return nil, fmt.Errorf("can't initialize feed aggregate")
	}

	feed, err := agg.GetFeed(
		ctx,
		uid,
		flavour,
		persistent,
		status,
		visibility,
		expired,
		filterParams,
	)
	if err != nil {
		return nil, fmt.Errorf("can't instantiate new feed: %w", err)
	}

	return feed, nil
}

// Link returns generated.LinkResolver implementation.
func (r *Resolver) Link() generated.LinkResolver { return &linkResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type linkResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
