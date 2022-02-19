package lambdawrap

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"io"
	"strings"
	"testing"
	"testing/iotest"
)

func TestS3Fetch(t *testing.T) {
	t.Run("errors returned by the S3 fetcher are propogated up", func(t *testing.T) {
		s3Fetcher := func(_ context.Context, _ events.S3Entity) (io.ReadCloser, error) {
			return nil, io.ErrUnexpectedEOF
		}

		d, err := S3Fetch(nil, s3Fetcher)(context.TODO(), events.S3EventRecord{})
		assert.Nil(t, d)
		assert.Error(t, err)
	})

	t.Run("errors returned by next are propagated back", func(t *testing.T) {
		s3Fetcher := func(_ context.Context, _ events.S3Entity) (io.ReadCloser, error) {
			return nil, nil
		}

		next := func(_ context.Context, r io.Reader) ([]byte, error) {
			return nil, io.ErrUnexpectedEOF
		}

		d, err := S3Fetch(next, s3Fetcher)(context.TODO(), events.S3EventRecord{})
		assert.Nil(t, d)
		assert.Error(t, err)
	})

	t.Run("reader returned by S3Fetcher propagates back", func(t *testing.T) {
		expectedS3Entity := events.S3Entity{
			ConfigurationID: "cfg",
		}

		s3Fetcher := func(_ context.Context, e events.S3Entity) (io.ReadCloser, error) {
			assert.Equal(t, expectedS3Entity, e)
			return io.NopCloser(strings.NewReader("data")), nil
		}

		next := func(ctx context.Context, r io.Reader) ([]byte, error) {
			d, err := io.ReadAll(r)
			assert.NoError(t, err)

			actualS3Entity, ok := S3EntityFromContext(ctx)
			assert.True(t, ok)
			assert.Equal(t, expectedS3Entity, actualS3Entity)
			return d, nil
		}

		d, err := S3Fetch(next, s3Fetcher)(context.TODO(), events.S3EventRecord{S3: expectedS3Entity})
		assert.NoError(t, err)
		assert.Equal(t, []byte("data"), d)
	})
}

func TestReadAll(t *testing.T) {
	t.Run("errors from reading all are propagated back", func(t *testing.T) {
		d, err := S3ReadAll(nil)(context.TODO(), iotest.ErrReader(io.ErrUnexpectedEOF))
		assert.Nil(t, d)
		assert.Error(t, err)
	})

	t.Run("errors from next are propagated back", func(t *testing.T) {
		next := func(_ context.Context, _ []byte) ([]byte, error) {
			return nil, io.ErrUnexpectedEOF
		}

		d, err := S3ReadAll(next)(context.TODO(), strings.NewReader("data"))
		assert.Nil(t, d)
		assert.Error(t, err)
	})

	t.Run("successful read all with next sends data back", func(t *testing.T) {
		next := func(_ context.Context, d []byte) ([]byte, error) {
			return d, nil
		}

		d, err := S3ReadAll(next)(context.TODO(), strings.NewReader("data"))
		assert.NoError(t, err)
		assert.Equal(t, []byte("data"), d)
	})
}
