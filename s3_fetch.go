package lambdawrap

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"io"
	"io/ioutil"
)

type S3Fetcher func(context.Context, events.S3Entity) (io.ReadCloser, error)

// S3Fetch consumes an events.S3EventRecord and retrieves the object from S3, providing an io.Reader to the next
// function. Default behaviour is to concatenate the []byte output from each message, returning to the caller.
//
// The actual network code to fetch an S3 object is provided to the function via i, a function that implements the
// S3Fetcher interface. This may make this function seem very light, however it is done so due to the complexities
// surrounding constructing S3 clients and the infinite combinations of configuration that might be needed.
//
// A copy of the event.S3Entity is added to the context, and can be extracted with S3EntityFromContext.
//
// A very basic S3 Fetcher implementation is included as impl.Fetcher, it is a submodule so will need to be imported
// separately, this is to prevent the dependency of the AWS SDK leaking into lambdawrap.
func S3Fetch(n func(context.Context, io.Reader) ([]byte, error), i S3Fetcher) func(context.Context, events.S3EventRecord) ([]byte, error) {
	return func(ctx context.Context, e events.S3EventRecord) ([]byte, error) {
		if r, err := i(ctx, e.S3); err != nil {
			return nil, fmt.Errorf("s3 fetch: %w", err)
		} else {
			ctx = context.WithValue(ctx, contextKeyS3Entity, e.S3)
			d, err := n(ctx, r)

			var closeErr error

			if r != nil {
				closeErr = r.Close()
			}

			if err != nil {
				return nil, fmt.Errorf("s3 fetch next: %w", err)
			} else if closeErr != nil {
				return nil, fmt.Errorf("s3 fetch close: %w", err)
			}

			return d, nil
		}
	}
}

// S3ReadAll consumes an io.Reader and provides a []byte to the next function. Default behaviour is to concatenate the
// []byte output from each message, returning to the caller.
func S3ReadAll(n func(context.Context, []byte) ([]byte, error)) func(context.Context, io.Reader) ([]byte, error) {
	return func(ctx context.Context, r io.Reader) ([]byte, error) {
		if rd, err := ioutil.ReadAll(r); err != nil {
			return nil, fmt.Errorf("read all error: %w", err)
		} else {
			if d, err := n(ctx, rd); err != nil {
				return nil, fmt.Errorf("read all next: %w", err)
			} else {
				return d, nil
			}
		}
	}
}

// S3EntityFromContext retrieves an events.S3Entity from the context, for use after an S3Fetch wrap has been used if
// the application needs the events.S3Entity that was downloaded.
func S3EntityFromContext(ctx context.Context) (events.S3Entity, bool) {
	if val := ctx.Value(contextKeyS3Entity); val != nil {
		return val.(events.S3Entity), true
	} else {
		return events.S3Entity{}, false
	}
}
