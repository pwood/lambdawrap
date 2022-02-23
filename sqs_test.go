package lambdawrap

import (
	"context"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func TestSQS(t *testing.T) {
	t.Run("each event.SQSMessage is processed, calls the next function that takes a []byte and the result is aggregated", func(t *testing.T) {
		in := events.SQSEvent{
			Records: []events.SQSMessage{
				{
					Body: "1",
				},
				{
					Body: "2",
				},
				{
					Body: "3",
				},
			},
		}

		next := func(_ context.Context, d []byte) ([]byte, error) {
			return d, nil
		}

		d, err := SQS(next)(context.TODO(), in)
		assert.NoError(t, err)
		assert.Equal(t, "123", string(d))
	})

	t.Run("each event.SQSMessage is processed, calls the next function that takes a string and the result is aggregated", func(t *testing.T) {
		in := events.SQSEvent{
			Records: []events.SQSMessage{
				{
					Body: "1",
				},
				{
					Body: "2",
				},
				{
					Body: "3",
				},
			},
		}

		next := func(_ context.Context, d string) ([]byte, error) {
			return []byte(d), nil
		}

		d, err := SQS(next)(context.TODO(), in)
		assert.NoError(t, err)
		assert.Equal(t, "123", string(d))
	})

	t.Run("each event.SQSMessage is processed, calls the next function that takes a structure, the result is aggregated and topic ARN is available on context", func(t *testing.T) {
		type myStruct struct {
			Val string
		}

		in := events.SQSEvent{
			Records: []events.SQSMessage{
				{
					EventSourceARN: "sqsARN",
					Body:           "{\"val\": \"1\"}",
				},
				{
					EventSourceARN: "sqsARN",
					Body:           "{\"val\": \"2\"}",
				},
				{
					EventSourceARN: "sqsARN",
					Body:           "{\"val\": \"3\"}",
				},
			},
		}

		next := func(ctx context.Context, d myStruct) ([]byte, error) {
			topic, ok := SQSTopicARNFromContext(ctx)
			assert.True(t, ok)
			assert.Equal(t, "sqsARN", topic)

			return []byte(d.Val), nil
		}

		d, err := SQS(next)(context.TODO(), in)
		assert.NoError(t, err)
		assert.Equal(t, "123", string(d))
	})

	t.Run("unmarshalling a Message with a JSON error will result in an error", func(t *testing.T) {
		type myStruct struct {
			Val string
		}

		in := events.SQSEvent{
			Records: []events.SQSMessage{
				{
					Body: "{\"val\": \"1\"",
				},
			},
		}

		next := func(_ context.Context, d myStruct) ([]byte, error) {
			return []byte(d.Val), nil
		}

		d, err := SQS(next)(context.TODO(), in)
		assert.Error(t, err)
		assert.Nil(t, d)
	})

	t.Run("an error from next will result in an error", func(t *testing.T) {
		in := events.SQSEvent{
			Records: []events.SQSMessage{
				{
					Body: "1",
				},
			},
		}

		next := func(_ context.Context, d string) ([]byte, error) {
			return nil, io.ErrUnexpectedEOF
		}

		d, err := SQS(next)(context.TODO(), in)
		assert.True(t, errors.Is(err, io.ErrUnexpectedEOF))
		assert.Nil(t, d)
	})
}
