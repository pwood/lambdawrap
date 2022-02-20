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

func TestSwitch(t *testing.T) {
	t.Run("returns an error if no match is found", func(t *testing.T) {
		_, err := Switch[string, string](func(s string) string {
			return s
		}, map[string]func(context.Context, string) ([]byte, error){
			"a": func(ctx context.Context, s string) ([]byte, error) {
				return nil, nil
			},
		})(context.TODO(), "b")

		assert.Error(t, err)
	})

	t.Run("returns data from function in map", func(t *testing.T) {
		d, err := Switch[string, string](func(s string) string {
			return s
		}, map[string]func(context.Context, string) ([]byte, error){
			"a": func(ctx context.Context, s string) ([]byte, error) {
				return []byte("data"), nil
			},
		})(context.Background(), "a")

		assert.Equal(t, []byte("data"), d)
		assert.NoError(t, err)
	})

	t.Run("returns error from function in map", func(t *testing.T) {
		d, err := Switch[string, string](func(s string) string {
			return s
		}, map[string]func(context.Context, string) ([]byte, error){
			"a": func(ctx context.Context, s string) ([]byte, error) {
				return nil, io.ErrUnexpectedEOF
			},
		})(context.Background(), "a")

		assert.Nil(t, d)
		assert.Error(t, err)
		assert.Equal(t, io.ErrUnexpectedEOF, err)
	})
}
