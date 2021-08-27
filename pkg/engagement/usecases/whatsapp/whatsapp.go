package whatsapp

import (
	"context"

	"github.com/savannahghi/engagement/pkg/engagement/infrastructure"
)

type WhatsappUsecases interface {
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

type WhatsappUsecasesImpl struct {
	infrastructure infrastructure.Infrastructure
}

func NewWhatsappUsecasesImpl(infrastructure infrastructure.Infrastructure) *WhatsappUsecasesImpl {
	return &WhatsappUsecasesImpl{infrastructure: infrastructure}
}

func (w *WhatsappUsecasesImpl) PhoneNumberVerificationCode(
	ctx context.Context,
	to string,
	code string,
	marketingMessage string,
) (bool, error) {
	return w.infrastructure.PhoneNumberVerificationCode(ctx, to, code, marketingMessage)
}

func (w *WhatsappUsecasesImpl) WellnessCardActivationDependant(
	ctx context.Context,
	to string,
	memberName string,
	cardName string,
	marketingMessage string,
) (bool, error) {
	return w.infrastructure.WellnessCardActivationDependant(ctx, to, memberName, cardName, marketingMessage)
}

func (w *WhatsappUsecasesImpl) WellnessCardActivationPrincipal(
	ctx context.Context,
	to string,
	memberName string,
	cardName string,
	minorAgeThreshold string,
	marketingMessage string,
) (bool, error) {
	return w.infrastructure.WellnessCardActivationPrincipal(ctx, to, memberName, cardName, minorAgeThreshold, marketingMessage)
}

func (w *WhatsappUsecasesImpl) BillNotification(
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

func (w *WhatsappUsecasesImpl) VirtualCards(
	ctx context.Context,
	to string,
	wellnessCardFamily string,
	virtualCardLink string,
	marketingMessage string,
) (bool, error) {
	return w.infrastructure.VirtualCards(ctx, to, wellnessCardFamily, virtualCardLink, marketingMessage)
}

func (w *WhatsappUsecasesImpl) VisitStart(
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

func (w *WhatsappUsecasesImpl) ClaimNotification(
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

func (w *WhatsappUsecasesImpl) PreauthApproval(
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

func (w *WhatsappUsecasesImpl) PreauthRequest(
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

func (w *WhatsappUsecasesImpl) SladeOTP(
	ctx context.Context,
	to string,
	name string,
	otp string,
	marketingMessage string,
) (bool, error) {
	return w.infrastructure.SladeOTP(ctx, to, name, otp, marketingMessage)
}
