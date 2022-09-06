package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"time"

	"github.com/savannahghi/engagementcore/pkg/engagement/domain"
	"github.com/savannahghi/engagementcore/pkg/engagement/presentation/graph/generated"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/serverutils"
)

// GetLibraryContent is the resolver for the getLibraryContent field.
func (r *queryResolver) GetLibraryContent(ctx context.Context) ([]*domain.GhostCMSPost, error) {
	startTime := time.Now()

	ghostCMSPost, err := r.infra.GetLibraryContent(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get library content: %v", err)
	}

	defer serverutils.RecordGraphqlResolverMetrics(ctx, startTime, "getLibraryContent", err)

	return ghostCMSPost, nil
}

// GetFaqsContent is the resolver for the getFaqsContent field.
func (r *queryResolver) GetFaqsContent(ctx context.Context, flavour feedlib.Flavour) ([]*domain.GhostCMSPost, error) {
	startTime := time.Now()

	faqs, err := r.infra.GetFaqsContent(ctx, flavour)
	if err != nil {
		return nil, fmt.Errorf("unable to get FAQs content: %v", err)
	}

	defer serverutils.RecordGraphqlResolverMetrics(ctx, startTime, "getFaqsContent", err)

	return faqs, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
