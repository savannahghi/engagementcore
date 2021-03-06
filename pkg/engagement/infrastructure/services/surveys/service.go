package surveys

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

// NPSResponseCollectionName firestore collection name where nps responses are stored
const NPSResponseCollectionName = "nps_response"

// ServiceSurveys defines the interactions with the surveys service
type ServiceSurveys interface {
	RecordNPSResponse(ctx context.Context, input dto.NPSInput) (bool, error)
}

// NewService initializes a surveys service
func NewService(repository database.Repository) *ServiceSurveyImpl {
	fc := &firebasetools.FirebaseClient{}
	firebaseApp, err := fc.InitFirebase()
	if err != nil {
		log.Panicf("unable to initialize Firebase app for Surveys service: %s", err)
	}
	ctx := context.Background()

	firestoreClient, err := firebaseApp.Firestore(ctx)
	if err != nil {
		log.Panicf("unable to initialize Firestore client: %s", err)
	}

	srv := &ServiceSurveyImpl{
		firestoreClient: firestoreClient,
		Repository:      repository,
	}
	srv.checkPreconditions()
	return srv
}

func (s ServiceSurveyImpl) checkPreconditions() {
	if s.firestoreClient == nil {
		log.Panicf("surveys service has a nil firestore client")
	}
}

// ServiceSurveyImpl is an surveys service
type ServiceSurveyImpl struct {
	firestoreClient *firestore.Client
	Repository      database.Repository
}

// RecordNPSResponse ...
func (s ServiceSurveyImpl) RecordNPSResponse(ctx context.Context, input dto.NPSInput) (bool, error) {
	ctx, span := tracer.Start(ctx, "RecordNPSResponse")
	defer span.End()
	s.checkPreconditions()
	response := &dto.NPSResponse{
		Name:      input.Name,
		Score:     input.Score,
		SladeCode: input.SladeCode,
	}

	response.SetID(uuid.New().String())

	if input.Email != nil {
		response.Email = input.Email
	}

	if input.PhoneNumber != nil {
		response.MSISDN = input.PhoneNumber
	}

	feedbacks := []dto.Feedback{}
	if input.Feedback != nil {

		for _, input := range input.Feedback {
			feedback := dto.Feedback{
				Question: input.Question,
				Answer:   input.Answer,
			}
			feedbacks = append(feedbacks, feedback)
		}

		response.Feedback = feedbacks
	}

	response.Timestamp = time.Now()

	err := s.Repository.SaveNPSResponse(ctx, response)
	if err != nil {
		helpers.RecordSpanError(span, err)
		return false, fmt.Errorf("cannot save nps response to firestore: %w", err)
	}

	return true, nil
}
