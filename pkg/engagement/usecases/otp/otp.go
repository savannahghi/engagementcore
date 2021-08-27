package otp

import (
	"context"

	"github.com/savannahghi/engagement/pkg/engagement/infrastructure"
)

type OTPUsecases interface {
	VerifyOtp(ctx context.Context, msisdn string, otp string) (bool, error)
	VerifyEmailOtp(ctx context.Context, email string, otp string) (bool, error)
	GenerateAndSendOTP(ctx context.Context, msisdn string, appID *string) (string, error)
	SendOTPToEmail(ctx context.Context, msisdn string, email *string, appID *string) (string, error)
	GenerateRetryOTP(ctx context.Context, msisdn string, retryStep int, appID *string) (string, error)
	EmailVerificationOtp(ctx context.Context, email string) (string, error)
}

type OTPUsecasesImpl struct {
	infrastructure infrastructure.Infrastructure
}

func NewOTPUsecasesImpl(infrastructure infrastructure.Infrastructure) *OTPUsecasesImpl {
	return &OTPUsecasesImpl{infrastructure: infrastructure}
}

func (o *OTPUsecasesImpl) GenerateAndSendOTP(ctx context.Context, msisdn string, appID *string) (string, error) {
	return o.infrastructure.GenerateAndSendOTP(ctx, msisdn, appID)
}

func (o *OTPUsecasesImpl) SendOTPToEmail(ctx context.Context, msisdn string, email *string, appID *string) (string, error) {
	return o.infrastructure.SendOTPToEmail(ctx, &msisdn, email, appID)
}

func (o *OTPUsecasesImpl) VerifyOtp(ctx context.Context, msisdn, verificationCode string) (bool, error) {
	return o.infrastructure.VerifyOtp(ctx, &msisdn, &verificationCode)
}

func (o *OTPUsecasesImpl) VerifyEmailOtp(ctx context.Context, email, verificationCode string) (bool, error) {
	return o.infrastructure.VerifyEmailOtp(ctx, &email, &verificationCode)
}

func (o *OTPUsecasesImpl) GenerateRetryOTP(ctx context.Context, msisdn string, retryStep int, appID *string) (string, error) {
	return o.infrastructure.GenerateRetryOTP(ctx, &msisdn, retryStep, appID)
}

func (o *OTPUsecasesImpl) EmailVerificationOtp(ctx context.Context, email string) (string, error) {
	return o.infrastructure.EmailVerificationOtp(ctx, &email)
}
