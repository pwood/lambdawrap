package lambdawrap

import (
	"context"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func TestDynamoDBStream(t *testing.T) {
	t.Run("each event.DynamoDBEventRecord is processed, calls the next function that takes a the events.DynamoDBEventRecord and the result is aggregated", func(t *testing.T) {
		in := events.DynamoDBEvent{
			Records: []events.DynamoDBEventRecord{
				{
					EventID: "1",
				},
				{
					EventID: "2",
				},
				{
					EventID: "3",
				},
			},
		}

		next := func(_ context.Context, d events.DynamoDBEventRecord) ([]byte, error) {
			return []byte(d.EventID), nil
		}

		d, err := DynamoDBStream(next)(context.TODO(), in)
		assert.NoError(t, err)
		assert.Equal(t, "123", string(d))
	})

	t.Run("an error from next will result in an error", func(t *testing.T) {
		in := events.DynamoDBEvent{
			Records: []events.DynamoDBEventRecord{
				{
					EventID: "1",
				},
			},
		}

		next := func(_ context.Context, d events.DynamoDBEventRecord) ([]byte, error) {
			return nil, io.ErrUnexpectedEOF
		}

		d, err := DynamoDBStream(next)(context.TODO(), in)
		assert.True(t, errors.Is(err, io.ErrUnexpectedEOF))
		assert.Nil(t, d)
	})
}
