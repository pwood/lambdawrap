package lambdawrap

import (
	"context"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func TestS3Notification(t *testing.T) {
	t.Run("each event.S3EventRecord is processed, calls the next function that takes a the events.S3EventRecord and the result is aggregated", func(t *testing.T) {
		in := events.S3Event{
			Records: []events.S3EventRecord{
				{
					EventName: "1",
				},
				{
					EventName: "2",
				},
				{
					EventName: "3",
				},
			},
		}

		next := func(_ context.Context, d events.S3EventRecord) ([]byte, error) {
			return []byte(d.EventName), nil
		}

		d, err := S3Notification(next)(context.TODO(), in)
		assert.NoError(t, err)
		assert.Equal(t, "123", string(d))
	})

	t.Run("an error from next will result in an error", func(t *testing.T) {
		in := events.S3Event{
			Records: []events.S3EventRecord{
				{
					EventName: "1",
				},
			},
		}

		next := func(_ context.Context, d events.S3EventRecord) ([]byte, error) {
			return nil, io.ErrUnexpectedEOF
		}

		d, err := S3Notification(next)(context.TODO(), in)
		assert.True(t, errors.Is(err, io.ErrUnexpectedEOF))
		assert.Nil(t, d)
	})
}
