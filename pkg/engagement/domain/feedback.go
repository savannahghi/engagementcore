package domain

import (
	"time"
)

// SurveyFeedback describes the set of feedback response from the user
// the question and the answer
type SurveyFeedback struct {
	Question string `json:"question" firestore:"question"`
	Answer   string `json:"answer" firestore:"answer"`
}

// SurveyFeedbackInput refers to the Question and Answer Input
type SurveyFeedbackInput struct {
	Question string `json:"question" firestore:"question"`
	Answer   string `json:"answer" firestore:"answer"`
}

// SurveyInput refers to the whole response collected from the survey:
// the set of questions and their answers,
// and any extraFeedback from the user
type SurveyInput struct {
	Feedback      []*SurveyFeedbackInput `json:"feedback"`
	ExtraFeedback string                 `json:"extraFeedback" firestore:"extraFeedback"`
}

// SurveyFeedbackResponse shows the response that will be saved to firestore
type SurveyFeedbackResponse struct {
	Feedback      []SurveyFeedback `json:"feedback" firestore:"feedback"`
	ExtraFeedback string           `json:"extraFeedback" firestore:"extraFeedback"`
	Timestamp     time.Time        `json:"timestamp,omitempty" firestore:"timestamp,omitempty"`
}
