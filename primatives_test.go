package lambdawrap

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func TestNop(t *testing.T) {
	t.Run("constructing and calling a Nop will return a nil byte array and nil error", func(t *testing.T) {
		d, err := Nop[string]()(context.TODO(), "")
		assert.Nil(t, d)
		assert.NoError(t, err)
	})
}

func TestErr(t *testing.T) {
	t.Run("constructing and calling a Err will return a nil byte array and the provided error", func(t *testing.T) {
		d, err := Err[string](io.ErrUnexpectedEOF)(context.TODO(), "")
		assert.Nil(t, d)
		assert.True(t, errors.Is(err, io.ErrUnexpectedEOF))
	})
}
