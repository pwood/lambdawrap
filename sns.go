package lambdawrap

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
)

// SNS provides a wrapper to iterate through multiple SNS records included in an events.SNSEvent. Default behaviour is
// to concatenate the []byte output from each message, returning to the caller.
//
// SNS will attempt to unmarshal any destination structure with JSON, this is implemented for chaining wraps
// (e.g. SNS(SQS(S3Notification()))). It is recommended you use DomainObject instead, a Codec can be provided to support
// encodings other than JSON.
//
// Do not use this directly for domain objects, use DomainObject.
func SNS[O any](n func(context.Context, O) ([]byte, error)) func(context.Context, events.SNSEvent) ([]byte, error) {
	return func(ctx context.Context, e events.SNSEvent) ([]byte, error) {
		var ret []byte

		for _, r := range e.Records {
			if p, err := sliceStringOrUnmarshal[O]([]byte(r.SNS.Message)); err != nil {
				return nil, fmt.Errorf("SNS unmarshal: %w", err)
			} else {
				if d, err := n(ctx, p); err != nil {
					return nil, fmt.Errorf("SNS next: %w", err)
				} else {
					ret = append(ret, d...)
				}
			}
		}

		return ret, nil
	}
}

func sliceStringOrUnmarshal[O any](data []byte) (O, error) {
	v := new(O)

	switch p := any(v).(type) {
	case *[]byte:
		*p = data
		return *v, nil
	case *string:
		*p = string(data)
		return *v, nil
	}

	err := json.Unmarshal(data, v)
	if err != nil {
		return *new(O), fmt.Errorf("couldn't unmarshal: %w", err)
	}

	return *v, nil
}
