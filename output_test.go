package lambdawrap

import (
	"context"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func TestOutput(t *testing.T) {
	t.Run("returns error if next function returns an error", func(t *testing.T) {
		next := func(_ context.Context, _ string) ([]byte, error) {
			return nil, io.ErrUnexpectedEOF
		}

		d, err := Output(next, nil)(context.TODO(), "")
		assert.Error(t, err)
		assert.Nil(t, d)
	})

	t.Run("returns error if output function returns an error", func(t *testing.T) {
		next := func(_ context.Context, _ string) ([]byte, error) {
			return []byte("data"), nil
		}

		outFn := func(_ context.Context, _ []byte) error {
			return io.ErrUnexpectedEOF
		}

		d, err := Output(next, outFn)(context.TODO(), "")
		assert.Error(t, err)
		assert.Nil(t, d)
	})

	t.Run("calls output function with data from next, and returns nothing upwards", func(t *testing.T) {
		next := func(_ context.Context, _ string) ([]byte, error) {
			return []byte("data"), nil
		}

		outFn := func(_ context.Context, d []byte) error {
			assert.Equal(t, []byte("data"), d)
			return nil
		}

		d, err := Output(next, outFn)(context.TODO(), "")
		assert.NoError(t, err)
		assert.Nil(t, d)
	})
}
