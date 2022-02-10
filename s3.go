package lambdawrap

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
)

// S3Notification provides a wrapper to iterate through multiple S3 records included in an events.S3Event. Default behaviour is
// to concatenate the []byte output from each message, returning to the caller.
func S3Notification(n func(context.Context, events.S3EventRecord) ([]byte, error)) func(context.Context, events.S3Event) ([]byte, error) {
	return func(ctx context.Context, e events.S3Event) ([]byte, error) {
		var ret []byte

		for _, r := range e.Records {
			if d, err := n(ctx, r); err != nil {
				return nil, fmt.Errorf("S3Notification next: %w", err)
			} else {
				ret = append(ret, d...)
			}
		}

		return ret, nil
	}
}
