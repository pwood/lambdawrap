package lambdawrap

import (
	"context"
)

// Filter is a generic component that can be added to most wrapper chains, it calls the user provided function for each
// object passed in the chain. That function should return true if the message should be processed further, false if
// processing should go no further, or raise an error if a problem occurs.
//
// An example might be S3 notification processing:
//
//   myS3Filter := func(_ context.Context, s3Not events.S3EventRecord) (bool, error) {
//     return (s3Not.EventName == "s3:ObjectCreated:Put"), nil
//   }
//
//   myS3Consumer := func(_ context.Context, fileContents []byte) ([]byte, error) {
//     return fileContents, nil
//   }
//
//   wrap := SQS(S3Notification(Filter(myS3Filter, S3Fetch(RealAll(myS3Consumer)))))
func Filter[O any](f func(context.Context, O) (bool, error), n func(context.Context, O) ([]byte, error)) func(context.Context, O) ([]byte, error) {
	return Match[O](f, n, Nop[O]())
}

// Match is similar to Filter, however it permits handling the failing case with a different function, n. This permits
// the code path to diverge based upon the match result. This can be useful with guardrail filters that should error
// if the filter fails a match.
func Match[O any](f func(context.Context, O) (bool, error), m func(context.Context, O) ([]byte, error), n func(context.Context, O) ([]byte, error)) func(context.Context, O) ([]byte, error) {
	return func(ctx context.Context, o O) ([]byte, error) {
		if r, err := f(ctx, o); err != nil {
			return nil, err
		} else {
			if r {
				return m(ctx, o)
			} else {
				return n(ctx, o)
			}
		}
	}
}
