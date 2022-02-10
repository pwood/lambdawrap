package lambdawrap

import (
	"context"
)

// Nop is a construct to do nothing further to an event, and return an empty non error response.
func Nop[O any]() func(context.Context, O) ([]byte, error) {
	return func(ctx context.Context, o O) ([]byte, error) {
		return nil, nil
	}
}

// Err is a construct to return an empty error response.
func Err[O any](e error) func(context.Context, O) ([]byte, error) {
	return func(ctx context.Context, o O) ([]byte, error) {
		return nil, e
	}
}
