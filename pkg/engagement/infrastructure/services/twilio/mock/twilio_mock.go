package mock

import (
	"context"
	"net/url"

	"github.com/savannahghi/engagementcore/pkg/engagement/application/common/dto"
)

// FakeServiceTwilio defines the interaction with the twilio service
type FakeServiceTwilio struct {
	MakeTwilioRequestFn func(
		method string,
		urlPath string,
		content url.Values,
		target interface{},
	) error

	RoomFn func(ctx context.Context) (*dto.Room, error)

	TwilioAccessTokenFn func(ctx context.Context) (*dto.AccessToken, error)

	SendSMSFn func(ctx context.Context, to string, msg string) error

	SaveTwilioVideoCallbackStatusFn func(
		ctx context.Context,
		data dto.CallbackData,
	) error

	PhoneNumberVerificationCodeFn func(
		ctx context.Context,
		to string,
		code string,
		marketingMessage string,
	) (bool, error)

	SaveTwilioCallbackResponseFn func(
		ctx context.Context,
		data dto.Message,
	) error

	TemporaryPINFn func(
		ctx context.Context,
		to string,
		message string,
	) (bool, error)
}

// Room is a mock of the Room method
func (f *FakeServiceTwilio) Room(ctx context.Context) (*dto.Room, error) {
	return f.RoomFn(ctx)
}

// TwilioAccessToken is a mock of the TwilioAccessToken method
func (f *FakeServiceTwilio) TwilioAccessToken(ctx context.Context) (*dto.AccessToken, error) {
	return f.TwilioAccessTokenFn(ctx)
}

// SendSMS is a mock of the SendSMS method
func (f *FakeServiceTwilio) SendSMS(ctx context.Context, to string, msg string) error {
	return f.SendSMSFn(ctx, to, msg)
}

// SaveTwilioVideoCallbackStatus is a mock of the SendSMS method
func (f *FakeServiceTwilio) SaveTwilioVideoCallbackStatus(
	ctx context.Context,
	data dto.CallbackData,
) error {
	return f.SaveTwilioVideoCallbackStatusFn(ctx, data)
}

// PhoneNumberVerificationCode is a mock of the PhoneNumberVerificationCode method
func (f *FakeServiceTwilio) PhoneNumberVerificationCode(
	ctx context.Context,
	to string,
	code string,
	marketingMessage string,
) (bool, error) {
	return f.PhoneNumberVerificationCodeFn(ctx, to, code, marketingMessage)
}

// SaveTwilioCallbackResponse is a mock of the SaveTwilioCallbackResponse method
func (f *FakeServiceTwilio) SaveTwilioCallbackResponse(
	ctx context.Context,
	data dto.Message,
) error {
	return f.SaveTwilioCallbackResponseFn(ctx, data)
}

// MakeTwilioRequest is a mock of the MakeTwilioRequest method
func (f *FakeServiceTwilio) MakeTwilioRequest(
	method string,
	urlPath string,
	content url.Values,
	target interface{},
) error {
	return f.MakeTwilioRequestFn(method, urlPath, content, target)
}

// TemporaryPIN is a mock of the TemporaryPIN method
func (f *FakeServiceTwilio) TemporaryPIN(
	ctx context.Context,
	to string,
	message string,
) (bool, error) {
	return f.TemporaryPINFn(ctx, to, message)
}
