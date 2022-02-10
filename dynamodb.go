package lambdawrap

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
)

// DynamoDBStream provides a wrapper to iterate through multiple events.DynamoDBEvent. Default behaviour is to
// concatenate the []byte output from each message, returning to the caller.
func DynamoDBStream(n func(context.Context, events.DynamoDBEventRecord) ([]byte, error)) func(context.Context, events.DynamoDBEvent) ([]byte, error) {
	return func(ctx context.Context, e events.DynamoDBEvent) ([]byte, error) {
		var ret []byte

		for _, r := range e.Records {
			if d, err := n(ctx, r); err != nil {
				return nil, fmt.Errorf("DynamoDBStream next: %w", err)
			} else {
				ret = append(ret, d...)
			}
		}

		return ret, nil
	}
}
