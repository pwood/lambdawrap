package lambdawrap

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
)

// SQS provides a wrapper to iterate through multiple SQS records included in an events.SQSEvent. Default behaviour is
// to concatenate the []byte output from each message, returning to the caller.
//
// SQS will attempt to unmarshal any destination structure with JSON, this is implemented for chaining wraps
// (e.g. SNS(SQS(S3()))). It is recommended you use DomainObject instead, a Codec can be provided to support encodings
// other than JSON.
//
// Do not use this directly for domain objects, use DomainObject.
func SQS[O any](n func(context.Context, O) ([]byte, error)) func(context.Context, events.SQSEvent) ([]byte, error) {
	return func(ctx context.Context, e events.SQSEvent) ([]byte, error) {
		var ret []byte

		for _, r := range e.Records {
			if p, err := sliceStringOrUnmarshal[O]([]byte(r.Body)); err != nil {
				return nil, fmt.Errorf("SQS unmarshal: %w", err)
			} else {
				if d, err := n(ctx, p); err != nil {
					return nil, fmt.Errorf("SQS next: %w", err)
				} else {
					ret = append(ret, d...)
				}
			}
		}

		return ret, nil
	}
}
