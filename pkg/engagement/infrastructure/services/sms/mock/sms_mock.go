package mock

import (
	"context"

	"github.com/savannahghi/engagement/pkg/engagement/application/common/dto"
	"github.com/savannahghi/enumutils"
)

// FakeServiceSMS defines the interactions with the mock sms service
type FakeServiceSMS struct {
	SendToManyFn func(
		ctx context.Context,
		message string,
		to []string,
		from enumutils.SenderID,
	) (*dto.SendMessageResponse, error)

	SendFn func(
		ctx context.Context,
		to, message string,
		from enumutils.SenderID,
	) (*dto.SendMessageResponse, error)
}

// SendToMany ...
func (f *FakeServiceSMS) SendToMany(
	ctx context.Context,
	message string,
	to []string,
	from enumutils.SenderID,
) (*dto.SendMessageResponse, error) {
	return f.SendToManyFn(ctx, message, to, from)
}

// Send ...
func (f *FakeServiceSMS) Send(
	ctx context.Context,
	to, message string,
	from enumutils.SenderID,
) (*dto.SendMessageResponse, error) {
	return f.SendFn(ctx, to, message, from)
}
