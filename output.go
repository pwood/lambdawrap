package lambdawrap

import (
	"context"
	"fmt"
)

// Output is used to capture the result of the next, pass it to another function and then discard the data. This is
// useful for the output of DomainObject when the source has multiple records (and thus multiple next invocations), such
// as SQS, SNS, DynamoDBStream or S3Notification.
func Output[O any](n func(context.Context, O) ([]byte, error), ofn func(context.Context, []byte) error) func(context.Context, O) ([]byte, error) {
	return func(ctx context.Context, o O) ([]byte, error) {
		d, err := n(ctx, o)
		if err != nil {
			return nil, fmt.Errorf("output next: %w", err)
		}

		err = ofn(ctx, d)
		if err != nil {
			return nil, fmt.Errorf("output fn: %w", err)
		} else {
			return nil, nil
		}
	}
}
