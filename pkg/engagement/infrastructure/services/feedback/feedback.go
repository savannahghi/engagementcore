package feedback

import (
	"context"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/google/uuid"
	"github.com/savannahghi/engagementcore/pkg/engagement/application/common/dto"
	"github.com/savannahghi/engagementcore/pkg/engagement/application/common/helpers"
	"github.com/savannahghi/engagementcore/pkg/engagement/infrastructure/database"
	"github.com/savannahghi/firebasetools"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("github.com/savannahghi/engagementcore/pkg/engagement/services/surveys")

const FeedbackCollectionName = "patient_feedback"

type ServiceFeedback interface {
	RecordPatientFeedback(ctx context.Context, input dto.PatientFeedbackInput) (bool, error)
}

type ServiceFeedbackImpl struct {
	firestoreClient *firestore.Client
	Repository      database.Repository
}

func (s ServiceFeedbackImpl) checkPreconditions() {
	if s.firestoreClient == nil {
		log.Panicf("surveys service has a nil firestore client")
	}
}

func NewService(repository database.Repository) *ServiceFeedbackImpl {
	fc := &firebasetools.FirebaseClient{}
	firebaseApp, err := fc.InitFirebase()
	if err != nil {
		log.Panicf("unable to initialize Firebase app for Surveys service: %s", err)
	}

	ctx := context.Background()

	firebaseClient, err := firebaseApp.Firestore(ctx)
	if err != nil {
		log.Panicf("unable to initialize Firestore client: %s", err)
	}

	srv := &ServiceFeedbackImpl{
		firestoreClient: firebaseClient,
		Repository:      repository,
	}
	srv.checkPreconditions()
	return srv
}

func (s ServiceFeedbackImpl) RecordPatientFeedback(ctx context.Context, input dto.PatientFeedbackInput) (bool, error) {
	ctx, span := tracer.Start(ctx, "RecordPatientFeedback")
	defer span.End()
	s.checkPreconditions()

	response := &dto.PatientFeedbackResponse{}
	response.SetID(uuid.New().String())

	feedbacks := []dto.Feedback{}
	if input.Feedback != nil {
		for _, input := range input.Feedback {
			feedback := dto.Feedback{
				Question: input.Question,
				Answer:   input.Answer,
			}
			feedbacks = append(feedbacks, feedback)
		}
	}
	response.Feedback = feedbacks
	response.ExtraFeedback = input.ExtraFeedback
	response.Timestamp = time.Now()

	err := s.Repository.SavePatientFeedback(ctx, response)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return false, fmt.Errorf("cannot patient's feedback to firestore: %w", err)
	}

	return true, nil
}
