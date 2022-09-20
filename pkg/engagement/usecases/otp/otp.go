package otp

import (
	"context"

	"github.com/savannahghi/engagementcore/pkg/engagement/application/common/dto"
	"github.com/savannahghi/engagementcore/pkg/engagement/infrastructure"
)

// UsecaseOTP defines otp service usecases interface
type UsecaseOTP interface {
	GenerateAndSendOTP(
		ctx context.Context,
		msisdn string,
		appID *string,
	) (string, error)
	SendOTPToEmail(
		ctx context.Context,
		msisdn,
		email *string,
		appID *string,
	) (string, error)
	SaveOTPToFirestore(
		otp dto.OTP,
	) error
	VerifyOtp(
		ctx context.Context,
		msisdn *string,
		verificationCode *string,
	) (bool, error)
	VerifyEmailOtp(
		ctx context.Context,
		email *string,
		verificationCode *string,
	) (bool, error)
	GenerateRetryOTP(
		ctx context.Context,
		msisdn *string,
		retryStep int,
		appID *string,
	) (string, error)
	EmailVerificationOtp(
		ctx context.Context,
		email *string,
	) (string, error)
	GenerateOTP(
		ctx context.Context,
	) (string, error)
	SendTemporaryPIN(
		ctx context.Context,
		input dto.TemporaryPIN,
	) error
}

// ImplOTP is the OTP service implementation
type ImplOTP struct {
	infrastructure infrastructure.Interactor
}

// NewOTP initializes a otp service instance
func NewOTP(infrastructure infrastructure.Interactor) *ImplOTP {
	return &ImplOTP{
		infrastructure: infrastructure,
	}
}

// GenerateAndSendOTP creates an OTP and sends it to the
// supplied phone number as a text message
func (f *ImplOTP) GenerateAndSendOTP(
	ctx context.Context,
	msisdn string,
	appID *string,
) (string, error) {
	i := f.infrastructure.ServiceOTPImpl
	return i.GenerateAndSendOTP(
		ctx,
		msisdn,
		appID,
	)
}

// SendOTPToEmail is a companion to GenerateAndSendOTP function
// It will send the generated OTP to the provided email address
func (f *ImplOTP) SendOTPToEmail(
	ctx context.Context,
	msisdn string,
	email *string,
	appID *string) (string, error) {
	i := f.infrastructure.ServiceOTPImpl
	return i.SendOTPToEmail(
		ctx,
		&msisdn,
		email,
		appID,
	)
}

// SaveOTPToFirestore persists the supplied OTP
func (f *ImplOTP) SaveOTPToFirestore(
	otp dto.OTP,
) error {
	i := f.infrastructure.ServiceOTPImpl
	return i.SaveOTPToFirestore(
		otp,
	)
}

// VerifyOtp checks for the validity of the supplied OTP but does not invalidate it
func (f *ImplOTP) VerifyOtp(
	ctx context.Context,
	msisdn string,
	verificationCode *string,
) (bool, error) {
	i := f.infrastructure.ServiceOTPImpl
	return i.VerifyOtp(
		ctx,
		&msisdn,
		verificationCode,
	)
}

// VerifyEmailOtp checks for the validity of the supplied OTP but does not invalidate it
func (f *ImplOTP) VerifyEmailOtp(
	ctx context.Context,
	email,
	verificationCode *string,
) (bool, error) {
	i := f.infrastructure.ServiceOTPImpl
	return i.VerifyEmailOtp(
		ctx,
		email,
		verificationCode,
	)
}

// GenerateRetryOTP generates fallback OTPs when Africa is talking sms fails
func (f *ImplOTP) GenerateRetryOTP(
	ctx context.Context,
	msisdn *string,
	retryStep int,
	appID *string,
) (string, error) {
	i := f.infrastructure.ServiceOTPImpl
	return i.GenerateRetryOTP(
		ctx,
		msisdn,
		retryStep,
		appID,
	)
}

// EmailVerificationOtp generates an OTP to the supplied email for verification
func (f *ImplOTP) EmailVerificationOtp(
	ctx context.Context,
	email *string,
) (string, error) {
	i := f.infrastructure.ServiceOTPImpl
	return i.EmailVerificationOtp(
		ctx,
		email,
	)
}

// GenerateOTP generates an OTP
func (f *ImplOTP) GenerateOTP(
	ctx context.Context,
) (string, error) {
	i := f.infrastructure.ServiceOTPImpl
	return i.GenerateOTP(
		ctx,
	)
}

// SendTemporaryPIN sends a temporary PIN message to user via whatsapp and SMS
func (f *ImplOTP) SendTemporaryPIN(
	ctx context.Context,
	input dto.TemporaryPIN,
) error {
	i := f.infrastructure.ServiceOTPImpl
	return i.SendTemporaryPIN(
		ctx,
		input,
	)
}
