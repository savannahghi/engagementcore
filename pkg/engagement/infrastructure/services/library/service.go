package library

import (
	"context"
	"encoding/json"
	"fmt"

	"net/http"
	"time"

	"github.com/savannahghi/engagementcore/pkg/engagement/application/common/helpers"
	"github.com/savannahghi/engagementcore/pkg/engagement/domain"
	"github.com/savannahghi/engagementcore/pkg/engagement/infrastructure/services/onboarding"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/profileutils"
	"github.com/savannahghi/serverutils"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("github.com/savannahghi/engagementcore/pkg/engagement/services/library")

// Library service constants
const (
	ghostCMSAPIEndpoint          = "GHOST_CMS_API_ENDPOINT"
	ghostCMSAPIKey               = "GHOST_CMS_API_KEY"
	apiRoot                      = "/ghost/api/v3/content/posts/?"
	includeTags                  = "&include=tags"
	includeAuthors               = "&include=authors"
	formats                      = "&formats=html,plaintext"
	allowedConsumerFeedTagFilter = "&filter=tag:feed-consumer&filter=tag:welcome&filter=tag:what-is&filter=tag:getting-started"
	allowedProFeedTagFilter      = "&filter=tag:feed-pro&filter=tag:welcome&filter=tag:what-is&filter=tag:getting-started"
	allowedPROFAQsTagFilter      = "&filter=tag:faqs-pro"
	allowedConsumerFAQsTagFilter = "&filter=tag:faqs-consumer&filter=tag:how-to"
	allowedLibraryTagFilter      = "&filter=tag:diet&filter=tag:health-tips"
	allowedAgentTagFilter        = "&filter=tag:agent-faqs"
	allowedEmployeeTagFilter     = "&filter=tag:employee-faqs"
)

// ServiceLibrary ...
type ServiceLibrary interface {
	GetFeedContent(ctx context.Context, flavour feedlib.Flavour) ([]*domain.GhostCMSPost, error)
	GetFaqsContent(ctx context.Context, flavour feedlib.Flavour) ([]*domain.GhostCMSPost, error)
	GetLibraryContent(ctx context.Context) ([]*domain.GhostCMSPost, error)
}

type requestType int

const (
	feedRequestConsumer requestType = iota + 1
	feedRequestPro
	faqsRequestConsumer
	faqsRequestPro
	libraryRequest
	employeeHelpRequest
	agentHelpRequest
)

// NewLibraryService creates a new library ServiceLibraryImpl
func NewLibraryService(
	onboarding onboarding.ProfileService,
) *ServiceLibraryImpl {
	e := serverutils.MustGetEnvVar(ghostCMSAPIEndpoint)
	a := serverutils.MustGetEnvVar(ghostCMSAPIKey)

	srv := &ServiceLibraryImpl{
		APIEndpoint:  e,
		APIKey:       a,
		PostsAPIRoot: fmt.Sprintf("%v%vkey=%v", e, apiRoot, a),
		onboarding:   onboarding,
	}
	srv.checkPreconditions()
	return srv
}

// ServiceLibraryImpl organizes library functionality
// APIEndpoint should be of the form https://<name>.ghost.io
type ServiceLibraryImpl struct {
	APIEndpoint  string
	APIKey       string
	PostsAPIRoot string
	onboarding   onboarding.ProfileService
}

func (s ServiceLibraryImpl) checkPreconditions() {
	if s.APIEndpoint == "" {
		log.Panicf("Ghost API endpoint must be set")
	}

	if s.APIKey == "" {
		log.Panicf("Ghost API key must be set")
	}

	if s.PostsAPIRoot == "" {
		log.Panicf("Ghost Post API root must be set")
	}
}

func (s ServiceLibraryImpl) composeRequest(reqType requestType) string {
	var urlRequest string
	switch reqType {
	case feedRequestConsumer:
		urlRequest = fmt.Sprintf(
			"%v%v%v%v%v",
			s.PostsAPIRoot,
			includeTags,
			allowedConsumerFeedTagFilter,
			includeAuthors,
			formats,
		)
	case feedRequestPro:
		urlRequest = fmt.Sprintf(
			"%v%v%v%v%v",
			s.PostsAPIRoot,
			includeTags,
			allowedProFeedTagFilter,
			includeAuthors,
			formats,
		)
	case faqsRequestConsumer:
		urlRequest = fmt.Sprintf(
			"%v%v%v%v%v",
			s.PostsAPIRoot,
			includeTags,
			allowedConsumerFAQsTagFilter,
			includeAuthors,
			formats,
		)
	case faqsRequestPro:
		urlRequest = fmt.Sprintf(
			"%v%v%v%v%v",
			s.PostsAPIRoot,
			includeTags,
			allowedPROFAQsTagFilter,
			includeAuthors,
			formats,
		)
	case libraryRequest:
		urlRequest = fmt.Sprintf(
			"%v%v%v%v%v",
			s.PostsAPIRoot,
			includeTags,
			allowedLibraryTagFilter,
			includeAuthors,
			formats,
		)

	case employeeHelpRequest:
		urlRequest = fmt.Sprintf(
			"%v%v%v%v%v",
			s.PostsAPIRoot,
			includeTags,
			allowedEmployeeTagFilter,
			includeAuthors,
			formats,
		)
	case agentHelpRequest:
		urlRequest = fmt.Sprintf(
			"%v%v%v%v%v",
			s.PostsAPIRoot,
			includeTags,
			allowedAgentTagFilter,
			includeAuthors,
			formats,
		)
	}
	return urlRequest
}

func (s ServiceLibraryImpl) getCMSPosts(ctx context.Context, requestType requestType) ([]*domain.GhostCMSPost, error) {
	_, span := tracer.Start(ctx, "getCMSPosts")
	defer span.End()
	s.checkPreconditions()
	url := s.composeRequest(requestType)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf("failed to create action request with error; %v", err)
	}
	c := &http.Client{Timeout: time.Second * 300}
	resp, err := c.Do(req)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf("error occurred when posting to %v with err %v", url, err)
	}
	defer resp.Body.Close()

	var rr domain.GhostCMSServerResponse
	if err := json.NewDecoder(resp.Body).Decode(&rr); err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf("failed to decoder response with err %v", err)
	}
	return rr.Posts, nil
}

// GetFeedContent fetches posts that should be added to the feed.
func (s ServiceLibraryImpl) GetFeedContent(ctx context.Context, flavour feedlib.Flavour) ([]*domain.GhostCMSPost, error) {
	var request requestType

	switch flavour {
	case feedlib.FlavourConsumer:
		request = feedRequestConsumer
	case feedlib.FlavourPro:
		request = feedRequestConsumer

	}

	return s.getCMSPosts(ctx, request)
}

// GetFaqsContent fetches posts tagged as FAQs.
func (s ServiceLibraryImpl) GetFaqsContent(ctx context.Context, flavour feedlib.Flavour) ([]*domain.GhostCMSPost, error) {
	ctx, span := tracer.Start(ctx, "GetFaqsContent")
	defer span.End()
	if !feedlib.FlavourConsumer.IsValid() {
		return nil, fmt.Errorf("flavour `%s` is invalid", flavour)
	}
	if flavour == feedlib.FlavourConsumer {
		return s.getCMSPosts(ctx, faqsRequestConsumer)
	}

	// get profile from onboarding service
	user, err := profileutils.GetLoggedInUser(ctx)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf("unable to get user: %w", err)
	}

	profile, err := s.onboarding.GetUserProfile(ctx, user.UID)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return nil, fmt.Errorf("unable to get user profile: %w", err)
	}

	switch profile.Role {
	case profileutils.RoleTypeEmployee:
		return s.getCMSPosts(ctx, employeeHelpRequest)
	case profileutils.RoleTypeAgent:
		return s.getCMSPosts(ctx, agentHelpRequest)
	default:
		return s.getCMSPosts(ctx, faqsRequestPro)

	}
}

// GetLibraryContent gets library content to be show under library section of the app.
func (s ServiceLibraryImpl) GetLibraryContent(ctx context.Context) ([]*domain.GhostCMSPost, error) {
	return s.getCMSPosts(ctx, libraryRequest)
}
