package lambdawrap

import (
	"context"
	"fmt"
)

// Codec defines an interface that can be used to Marshal and Unmarshal different encodings back into Go structures.
//
// See codec.JSON and codec.YAML.
type Codec interface {
	// Marshal processes v and encodes into a []byte
	Marshal(v any) ([]byte, error)
	// Unmarshal processes the data, and decodes back into v
	Unmarshal(data []byte, v any) error
}

// DomainObject provides an automated approach to unmarshalling an input domain object, and then automatically
// marshalling the output domain object. For processes that are side effect only (i.e. no output type), see SideEffect
// to mask the return type, otherwise ensure the that first return value of n is nil.
func DomainObject[I any, O any](n func(context.Context, I) (O, error), c Codec) func(context.Context, []byte) ([]byte, error) {
	return func(ctx context.Context, d []byte) ([]byte, error) {
		in := new(I)
		err := c.Unmarshal(d, in)
		if err != nil {
			return nil, fmt.Errorf("DomainObject codec unmarshal failure: %w", err)
		}

		ret, err := n(ctx, *in)
		if err != nil {
			return nil, fmt.Errorf("DomainObject next: %w", err)
		}

		data, err := c.Marshal(ret)
		if err != nil {
			return nil, fmt.Errorf("DomainObject codec marshal failure: %w", err)
		}

		return data, nil
	}
}

// SideEffect provides a way to ignore the requirement of DomainObject to return an output domain object. It is unlikely
// another function in lambdawrap would then be passed into SideEffect.
//
// Example:
//
//   type domainInput struct {}
//
//   myFunc := func(ctx context.Context, in domainInput) error
//
//   DomainObject(SideEffect(myFunc))
//
func SideEffect[I any](n func(context.Context, I) error) func(context.Context, I) ([]byte, error) {
	return func(ctx context.Context, i I) ([]byte, error) {
		return nil, n(ctx, i)
	}
}
