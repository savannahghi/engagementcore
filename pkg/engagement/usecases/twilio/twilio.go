package twilio

import (
	"context"

	"github.com/savannahghi/engagementcore/pkg/engagement/application/common/dto"
	"github.com/savannahghi/engagementcore/pkg/engagement/infrastructure"
)

// UsecaseTwilio defines twilio service usecases interface
type UsecaseTwilio interface {
	Room(
		ctx context.Context,
	) (*dto.Room, error)

	TwilioAccessToken(
		ctx context.Context,
	) (*dto.AccessToken, error)

	SendSMS(
		ctx context.Context,
		to string,
		msg string,
	) error

	SaveTwilioVideoCallbackStatus(
		ctx context.Context,
		data dto.CallbackData,
	) error

	PhoneNumberVerificationCode(
		ctx context.Context,
		to string,
		code string,
		marketingMessage string,
	) (bool, error)

	SaveTwilioCallbackResponse(
		ctx context.Context,
		data dto.Message,
	) error

	TemporaryPIN(
		ctx context.Context,
		to string,
		message string,
	) (bool, error)
}

// ImplTwilio is the twilio service implementation
type ImplTwilio struct {
	infrastructure infrastructure.Interactor
}

// NewImplTwilio initializes a twilio service instance
func NewImplTwilio(
	infrastructure infrastructure.Interactor,
) *ImplTwilio {
	return &ImplTwilio{
		infrastructure: infrastructure,
	}
}

// Room represents a real-time audio, data, video, and/or screen-share session,
// and is the basic building block for a Programmable Video application.
func (t *ImplTwilio) Room(
	ctx context.Context,
) (*dto.Room, error) {
	i := t.infrastructure.ServiceTwilioImpl
	return i.Room(
		ctx,
	)
}

// TwilioAccessToken is used to generate short-lived credentials used to authenticate
// the client-side application to Twilio.
func (t *ImplTwilio) TwilioAccessToken(
	ctx context.Context,
) (*dto.AccessToken, error) {
	i := t.infrastructure.ServiceTwilioImpl
	return i.TwilioAccessToken(
		ctx,
	)
}

// SendSMS sends a text message through Twilio's programmable SMS
func (t *ImplTwilio) SendSMS(
	ctx context.Context,
	to string,
	msg string,
) error {
	i := t.infrastructure.ServiceTwilioImpl
	return i.SendSMS(
		ctx,
		to,
		msg,
	)
}

// SaveTwilioVideoCallbackStatus saves status callback data
func (t *ImplTwilio) SaveTwilioVideoCallbackStatus(
	ctx context.Context,
	data dto.CallbackData,
) error {
	i := t.infrastructure.ServiceTwilioImpl
	return i.SaveTwilioVideoCallbackStatus(
		ctx,
		data,
	)
}

// PhoneNumberVerificationCode sends Phone Number verification codes via WhatsApp
func (t *ImplTwilio) PhoneNumberVerificationCode(
	ctx context.Context,
	to string,
	code string,
	marketingMessage string,
) (bool, error) {
	i := t.infrastructure.ServiceTwilioImpl
	return i.PhoneNumberVerificationCode(
		ctx,
		to,
		code,
		marketingMessage,
	)
}

// SaveTwilioCallbackResponse saves the twilio callback response for future
// analysis
func (t *ImplTwilio) SaveTwilioCallbackResponse(
	ctx context.Context,
	data dto.Message,
) error {
	i := t.infrastructure.ServiceTwilioImpl
	return i.SaveTwilioCallbackResponse(
		ctx,
		data,
	)
}

//TemporaryPIN send PIN via whatsapp to user
func (t *ImplTwilio) TemporaryPIN(
	ctx context.Context,
	to string,
	message string,
) (bool, error) {
	i := t.infrastructure.ServiceTwilioImpl
	return i.TemporaryPIN(
		ctx,
		to,
		message,
	)
}
