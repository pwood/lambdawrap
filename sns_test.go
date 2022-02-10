package lambdawrap

import (
	"context"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func TestSNS(t *testing.T) {
	t.Run("each event.SNSEventRecord is processed, calls the next function that takes a []byte and the result is aggregated", func(t *testing.T) {
		in := events.SNSEvent{
			Records: []events.SNSEventRecord{
				{
					SNS: events.SNSEntity{
						Message: "1",
					},
				},
				{
					SNS: events.SNSEntity{
						Message: "2",
					},
				},
				{
					SNS: events.SNSEntity{
						Message: "3",
					},
				},
			},
		}

		next := func(_ context.Context, d []byte) ([]byte, error) {
			return d, nil
		}

		d, err := SNS(next)(context.TODO(), in)
		assert.NoError(t, err)
		assert.Equal(t, "123", string(d))
	})

	t.Run("each event.SNSEventRecord is processed, calls the next function that takes a string and the result is aggregated", func(t *testing.T) {
		in := events.SNSEvent{
			Records: []events.SNSEventRecord{
				{
					SNS: events.SNSEntity{
						Message: "1",
					},
				},
				{
					SNS: events.SNSEntity{
						Message: "2",
					},
				},
				{
					SNS: events.SNSEntity{
						Message: "3",
					},
				},
			},
		}

		next := func(_ context.Context, d string) ([]byte, error) {
			return []byte(d), nil
		}

		d, err := SNS(next)(context.TODO(), in)
		assert.NoError(t, err)
		assert.Equal(t, "123", string(d))
	})

	t.Run("each event.SNSEventRecord is processed, calls the next function that takes a structure and the result is aggregated", func(t *testing.T) {
		type myStruct struct {
			Val string
		}

		in := events.SNSEvent{
			Records: []events.SNSEventRecord{
				{
					SNS: events.SNSEntity{
						Message: "{\"val\": \"1\"}",
					},
				},
				{
					SNS: events.SNSEntity{
						Message: "{\"val\": \"2\"}",
					},
				},
				{
					SNS: events.SNSEntity{
						Message: "{\"val\": \"3\"}",
					},
				},
			},
		}

		next := func(_ context.Context, d myStruct) ([]byte, error) {
			return []byte(d.Val), nil
		}

		d, err := SNS(next)(context.TODO(), in)
		assert.NoError(t, err)
		assert.Equal(t, "123", string(d))
	})

	t.Run("unmarshalling a Message with a JSON error will result in an error", func(t *testing.T) {
		type myStruct struct {
			Val string
		}

		in := events.SNSEvent{
			Records: []events.SNSEventRecord{
				{
					SNS: events.SNSEntity{
						Message: "{\"val\": \"1\"",
					},
				},
			},
		}

		next := func(_ context.Context, d myStruct) ([]byte, error) {
			return []byte(d.Val), nil
		}

		d, err := SNS(next)(context.TODO(), in)
		assert.Error(t, err)
		assert.Nil(t, d)
	})

	t.Run("an error from next will result in an error", func(t *testing.T) {
		in := events.SNSEvent{
			Records: []events.SNSEventRecord{
				{
					SNS: events.SNSEntity{
						Message: "1",
					},
				},
			},
		}

		next := func(_ context.Context, d string) ([]byte, error) {
			return nil, io.ErrUnexpectedEOF
		}

		d, err := SNS(next)(context.TODO(), in)
		assert.True(t, errors.Is(err, io.ErrUnexpectedEOF))
		assert.Nil(t, d)
	})
}
