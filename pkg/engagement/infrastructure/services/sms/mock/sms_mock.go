package mock

import (
	"context"

	"github.com/savannahghi/silcomms"
)

// FakeServiceSMS defines the interactions with the mock sms service
type FakeServiceSMS struct {
	SendToManyFn func(
		ctx context.Context,
		to []string,
		message string,
	) (*silcomms.BulkSMSResponse, error)

	SendFn func(
		ctx context.Context,
		to, message string,
	) (*silcomms.BulkSMSResponse, error)
}

// SendToMany ...
func (f *FakeServiceSMS) SendToMany(
	ctx context.Context,
	to []string,
	message string,
) (*silcomms.BulkSMSResponse, error) {
	return f.SendToManyFn(ctx, to, message)
}

// Send ...
func (f *FakeServiceSMS) Send(
	ctx context.Context,
	to, message string,
) (*silcomms.BulkSMSResponse, error) {
	return f.SendFn(ctx, to, message)
}
