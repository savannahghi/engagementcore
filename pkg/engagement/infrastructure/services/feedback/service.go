package feedback

import (
	"context"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/savannahghi/engagementcore/pkg/engagement/application/common/helpers"
	"github.com/savannahghi/engagementcore/pkg/engagement/domain"
	"github.com/savannahghi/engagementcore/pkg/engagement/infrastructure/database"
	"github.com/savannahghi/firebasetools"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("github.com/savannahghi/engagementcore/pkg/engagement/services/feedback")

// ServiceFeedback defines the interactions with the feedback service
type ServiceFeedback interface {
	RecordSurveyFeedbackResponse(
		ctx context.Context,
		input *domain.SurveyInput,
	)
}

// ServiceFeedbackImpl is an surveys service
type ServiceFeedbackImpl struct {
	firestoreClient *firestore.Client
	Repository      database.Repository
}

// NewService initializes a feedback service
func NewService(repository database.Repository) *ServiceFeedbackImpl {
	fc := &firebasetools.FirebaseClient{}
	firebaseApp, err := fc.InitFirebase()
	if err != nil {
		log.Panicf("unable to initialize Firebase app for Feedback service: %s", err)
	}
	ctx := context.Background()

	firestoreClient, err := firebaseApp.Firestore(ctx)
	if err != nil {
		log.Panicf("unable to initialize Firestore client: %s", err)
	}

	srv := &ServiceFeedbackImpl{
		firestoreClient: firestoreClient,
		Repository:      repository,
	}

	srv.checkPreconditions()
	return srv
}

func (s ServiceFeedbackImpl) checkPreconditions() {
	if s.firestoreClient == nil {
		log.Panicf("feedback service has a nil firestore client")
	}
}

// RecordSurveyFeedbackResponse ...
func (s ServiceFeedbackImpl) RecordSurveyFeedbackResponse(ctx context.Context, input domain.SurveyInput) (bool, error) {
	ctx, span := tracer.Start(ctx, "RecordSurveyFeedbackResponse")
	defer span.End()
	s.checkPreconditions()

	response := &domain.SurveyFeedbackResponse{
		ExtraFeedback: input.ExtraFeedback,
	}

	feedbacks := []domain.SurveyFeedback{}
	if input.Feedback != nil {
		for _, input := range input.Feedback {
			feedback := domain.SurveyFeedback{
				Question: input.Question,
				Answer:   input.Answer,
			}
			feedbacks = append(feedbacks, feedback)
		}
		response.Feedback = feedbacks
	}

	response.Timestamp = time.Now()

	err := s.Repository.RecordSurveyFeedbackResponse(ctx, response)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return false, fmt.Errorf("cannot save nps response to firestore: %w", err)
	}

	return true, nil
}
