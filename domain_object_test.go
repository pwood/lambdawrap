package lambdawrap

import (
	"context"
	"github.com/pwood/lambdawrap/codec"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func TestDomainObject(t *testing.T) {
	type in struct {
		In string
	}

	type out struct {
		Out string
	}

	t.Run("an error during unmarshal is propagated", func(t *testing.T) {
		next := func(_ context.Context, i in) (out, error) {
			t.Fatal("next called unexpectedly")
			return out{}, nil
		}

		inputBytes := []byte(`{"In":"message"`)

		_, err := DomainObject(next, codec.JSON)(context.TODO(), inputBytes)

		assert.Error(t, err)
	})

	t.Run("an error during marshal is propagated", func(t *testing.T) {
		type out struct {
			Out string
			Ch  chan string // JSON codec will error when it encounters a Channel
		}

		next := func(_ context.Context, i in) (out, error) {
			return out{Out: i.In}, nil
		}

		inputBytes := []byte(`{"In":"message"}`)

		_, err := DomainObject(next, codec.JSON)(context.TODO(), inputBytes)

		assert.Error(t, err)
	})

	t.Run("calling DomainObject correctly unmarshal the input value, calls next and marshal the output value", func(t *testing.T) {
		next := func(_ context.Context, i in) (out, error) {
			return out{Out: i.In}, nil
		}

		inputBytes := []byte(`{"In":"message"}`)
		expectedData := []byte(`{"Out":"message"}`)

		actualData, err := DomainObject(next, codec.JSON)(context.TODO(), inputBytes)

		assert.NoError(t, err)
		assert.Equal(t, expectedData, actualData)
	})
}

func TestSideEffect(t *testing.T) {
	type in struct {
		In string
	}

	t.Run("returns an error if next returns an error", func(t *testing.T) {
		next := func(_ context.Context, i in) error {
			return io.ErrUnexpectedEOF
		}

		d, err := SideEffect(next)(context.TODO(), in{})
		assert.Nil(t, d)
		assert.Error(t, err)
	})

	t.Run("returns a nil []byte after calling next", func(t *testing.T) {
		wasCalled := false

		next := func(_ context.Context, i in) error {
			wasCalled = true
			return nil
		}

		d, err := SideEffect(next)(context.TODO(), in{})
		assert.Nil(t, d)
		assert.NoError(t, err)
		assert.True(t, wasCalled)
	})
}
