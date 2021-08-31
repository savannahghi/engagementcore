package otp

import (
	"context"

	"github.com/savannahghi/engagement/pkg/engagement/infrastructure"
)

// UseCases ...
type UseCases interface {
	VerifyOtp(ctx context.Context, msisdn string, otp string) (bool, error)
	VerifyEmailOtp(ctx context.Context, email string, otp string) (bool, error)
	GenerateAndSendOTP(ctx context.Context, msisdn string, appID *string) (string, error)
	SendOTPToEmail(ctx context.Context, msisdn string, email *string, appID *string) (string, error)
	GenerateRetryOTP(ctx context.Context, msisdn string, retryStep int, appID *string) (string, error)
	EmailVerificationOtp(ctx context.Context, email string) (string, error)
}

// UseCasesImpl ...
type UseCasesImpl struct {
	infrastructure infrastructure.Infrastructure
}

// NewOTPUsecasesImpl ...
func NewOTPUsecasesImpl(infrastructure infrastructure.Infrastructure) *UseCasesImpl {
	return &UseCasesImpl{infrastructure: infrastructure}
}

// GenerateAndSendOTP ...
func (o *UseCasesImpl) GenerateAndSendOTP(ctx context.Context, msisdn string, appID *string) (string, error) {
	return o.infrastructure.GenerateAndSendOTP(ctx, msisdn, appID)
}

// SendOTPToEmail ...
func (o *UseCasesImpl) SendOTPToEmail(ctx context.Context, msisdn string, email *string, appID *string) (string, error) {
	return o.infrastructure.SendOTPToEmail(ctx, &msisdn, email, appID)
}

// VerifyOtp ...
func (o *UseCasesImpl) VerifyOtp(ctx context.Context, msisdn, verificationCode string) (bool, error) {
	return o.infrastructure.VerifyOtp(ctx, &msisdn, &verificationCode)
}

// VerifyEmailOtp ...
func (o *UseCasesImpl) VerifyEmailOtp(ctx context.Context, email, verificationCode string) (bool, error) {
	return o.infrastructure.VerifyEmailOtp(ctx, &email, &verificationCode)
}

// GenerateRetryOTP ...
func (o *UseCasesImpl) GenerateRetryOTP(ctx context.Context, msisdn string, retryStep int, appID *string) (string, error) {
	return o.infrastructure.GenerateRetryOTP(ctx, &msisdn, retryStep, appID)
}

// EmailVerificationOtp ...
func (o *UseCasesImpl) EmailVerificationOtp(ctx context.Context, email string) (string, error) {
	return o.infrastructure.EmailVerificationOtp(ctx, &email)
}
