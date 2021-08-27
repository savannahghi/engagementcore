package whatsapp

import (
	"context"

	"github.com/savannahghi/engagement/pkg/engagement/infrastructure"
)

// UseCases ...
type UseCases interface {
	PhoneNumberVerificationCode(
		ctx context.Context,
		to string,
		code string,
		marketingMessage string,
	) (bool, error)

	WellnessCardActivationDependant(
		ctx context.Context,
		to string,
		memberName string,
		cardName string,
		marketingMessage string,
	) (bool, error)

	WellnessCardActivationPrincipal(
		ctx context.Context,
		to string,
		memberName string,
		cardName string,
		minorAgeThreshold string,
		marketingMessage string,
	) (bool, error)

	BillNotification(
		ctx context.Context,
		to string,
		productName string,
		billingPeriod string,
		billAmount string,
		paymentInstruction string,
		marketingMessage string,
	) (bool, error)

	VirtualCards(
		ctx context.Context,
		to string,
		wellnessCardFamily string,
		virtualCardLink string,
		marketingMessage string,
	) (bool, error)

	VisitStart(
		ctx context.Context,
		to string,
		memberName string,
		benefitName string,
		locationName string,
		startTime string,
		balance string,
		marketingMessage string,
	) (bool, error)

	ClaimNotification(
		ctx context.Context,
		to string,
		claimReference string,
		claimTypeParenthesized string,
		provider string,
		visitType string,
		claimTime string,
		marketingMessage string,
	) (bool, error)

	PreauthApproval(
		ctx context.Context,
		to string,
		currency string,
		amount string,
		benefit string,
		provider string,
		member string,
		careContact string,
		marketingMessage string,
	) (bool, error)

	PreauthRequest(
		ctx context.Context,
		to string,
		currency string,
		amount string,
		benefit string,
		provider string,
		requestTime string,
		member string,
		careContact string,
		marketingMessage string,
	) (bool, error)

	SladeOTP(
		ctx context.Context,
		to string,
		name string,
		otp string,
		marketingMessage string,
	) (bool, error)
}

// UseCasesImpl ...
type UseCasesImpl struct {
	infrastructure infrastructure.Infrastructure
}

// NewWhatsappUsecasesImpl ...
func NewWhatsappUsecasesImpl(infrastructure infrastructure.Infrastructure) *UseCasesImpl {
	return &UseCasesImpl{infrastructure: infrastructure}
}

// PhoneNumberVerificationCode ...
func (w *UseCasesImpl) PhoneNumberVerificationCode(
	ctx context.Context,
	to string,
	code string,
	marketingMessage string,
) (bool, error) {
	return w.infrastructure.PhoneNumberVerificationCode(ctx, to, code, marketingMessage)
}

// WellnessCardActivationDependant ...
func (w *UseCasesImpl) WellnessCardActivationDependant(
	ctx context.Context,
	to string,
	memberName string,
	cardName string,
	marketingMessage string,
) (bool, error) {
	return w.infrastructure.WellnessCardActivationDependant(ctx, to, memberName, cardName, marketingMessage)
}

// WellnessCardActivationPrincipal ...
func (w *UseCasesImpl) WellnessCardActivationPrincipal(
	ctx context.Context,
	to string,
	memberName string,
	cardName string,
	minorAgeThreshold string,
	marketingMessage string,
) (bool, error) {
	return w.infrastructure.WellnessCardActivationPrincipal(ctx, to, memberName, cardName, minorAgeThreshold, marketingMessage)
}

// BillNotification ...
func (w *UseCasesImpl) BillNotification(
	ctx context.Context,
	to string,
	productName string,
	billingPeriod string,
	billAmount string,
	paymentInstruction string,
	marketingMessage string,
) (bool, error) {
	return w.infrastructure.BillNotification(ctx, to, productName, billingPeriod, billAmount, paymentInstruction, marketingMessage)
}

// VirtualCards ...
func (w *UseCasesImpl) VirtualCards(
	ctx context.Context,
	to string,
	wellnessCardFamily string,
	virtualCardLink string,
	marketingMessage string,
) (bool, error) {
	return w.infrastructure.VirtualCards(ctx, to, wellnessCardFamily, virtualCardLink, marketingMessage)
}

// VisitStart ...
func (w *UseCasesImpl) VisitStart(
	ctx context.Context,
	to string,
	memberName string,
	benefitName string,
	locationName string,
	startTime string,
	balance string,
	marketingMessage string,
) (bool, error) {
	return w.infrastructure.VisitStart(ctx, to, memberName, benefitName, locationName, startTime, balance, marketingMessage)
}

// ClaimNotification ...
func (w *UseCasesImpl) ClaimNotification(
	ctx context.Context,
	to string,
	claimReference string,
	claimTypeParenthesized string,
	provider string,
	visitType string,
	claimTime string,
	marketingMessage string,
) (bool, error) {
	return w.infrastructure.ClaimNotification(ctx, to, claimReference, claimTypeParenthesized, provider, visitType, claimTime, marketingMessage)
}

// PreauthApproval ...
func (w *UseCasesImpl) PreauthApproval(
	ctx context.Context,
	to string,
	currency string,
	amount string,
	benefit string,
	provider string,
	member string,
	careContact string,
	marketingMessage string,
) (bool, error) {
	return w.infrastructure.PreauthApproval(ctx, to, currency, amount, benefit, provider, member, careContact, marketingMessage)
}

// PreauthRequest ...
func (w *UseCasesImpl) PreauthRequest(
	ctx context.Context,
	to string,
	currency string,
	amount string,
	benefit string,
	provider string,
	requestTime string,
	member string,
	careContact string,
	marketingMessage string,
) (bool, error) {
	return w.infrastructure.PreauthRequest(ctx, to, currency, amount, benefit, provider, requestTime, member, careContact, marketingMessage)
}

// SladeOTP ...
func (w *UseCasesImpl) SladeOTP(
	ctx context.Context,
	to string,
	name string,
	otp string,
	marketingMessage string,
) (bool, error) {
	return w.infrastructure.SladeOTP(ctx, to, name, otp, marketingMessage)
}
