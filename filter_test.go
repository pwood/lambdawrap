package lambdawrap

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func TestMatch(t *testing.T) {
	t.Run("returns an error if the matching function returns an error", func(t *testing.T) {
		matchFn := func(_ context.Context, _ string) (bool, error) {
			return false, io.ErrUnexpectedEOF
		}

		uncalledFn := func(_ context.Context, _ string) ([]byte, error) {
			t.Fatal("match decision function unexpectedly called")
			return nil, nil
		}

		d, err := Match[string](matchFn, uncalledFn, uncalledFn)(context.Background(), "")
		assert.Nil(t, d)
		assert.True(t, errors.Is(err, io.ErrUnexpectedEOF))
	})

	t.Run("calls the true match if the filter matches and returns the data from the match function", func(t *testing.T) {
		matchFn := func(_ context.Context, _ string) (bool, error) {
			return true, nil
		}

		uncalledFn := func(_ context.Context, _ string) ([]byte, error) {
			t.Fatal("match decision function unexpectedly called")
			return nil, nil
		}

		calledFn := func(_ context.Context, _ string) ([]byte, error) {
			return []byte("data"), nil
		}

		d, err := Match[string](matchFn, calledFn, uncalledFn)(context.Background(), "")
		assert.Equal(t, []byte("data"), d)
		assert.NoError(t, err)
	})

	t.Run("calls the false match if the filter does not match and returns the data from the no match function", func(t *testing.T) {
		matchFn := func(_ context.Context, _ string) (bool, error) {
			return false, nil
		}

		uncalledFn := func(_ context.Context, _ string) ([]byte, error) {
			t.Fatal("match decision function unexpectedly called")
			return nil, nil
		}

		calledFn := func(_ context.Context, _ string) ([]byte, error) {
			return []byte("data"), nil
		}

		d, err := Match[string](matchFn, uncalledFn, calledFn)(context.Background(), "")
		assert.Equal(t, []byte("data"), d)
		assert.NoError(t, err)
	})
}

func TestFilter(t *testing.T) {
	t.Run("returns an error if the matching function returns an error", func(t *testing.T) {
		matchFn := func(_ context.Context, _ string) (bool, error) {
			return false, io.ErrUnexpectedEOF
		}

		uncalledFn := func(_ context.Context, _ string) ([]byte, error) {
			t.Fatal("match decision function unexpectedly called")
			return nil, nil
		}

		d, err := Filter[string](matchFn, uncalledFn)(context.Background(), "")
		assert.Nil(t, d)
		assert.True(t, errors.Is(err, io.ErrUnexpectedEOF))
	})

	t.Run("calls the true match if the filter matches and returns the data from the match function", func(t *testing.T) {
		matchFn := func(_ context.Context, _ string) (bool, error) {
			return true, nil
		}

		calledFn := func(_ context.Context, _ string) ([]byte, error) {
			return []byte("data"), nil
		}

		d, err := Filter[string](matchFn, calledFn)(context.Background(), "")
		assert.Equal(t, []byte("data"), d)
		assert.NoError(t, err)
	})

	t.Run("calls the false match if the filter does not match and returns the data from the no match function", func(t *testing.T) {
		matchFn := func(_ context.Context, _ string) (bool, error) {
			return false, nil
		}

		uncalledFn := func(_ context.Context, _ string) ([]byte, error) {
			t.Fatal("match decision function unexpectedly called")
			return nil, nil
		}

		d, err := Filter[string](matchFn, uncalledFn)(context.Background(), "")
		assert.Nil(t, d)
		assert.NoError(t, err)
	})
}
