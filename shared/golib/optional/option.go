// Package optional provides type-safe nullable values using Go generics
package optional

import (
	"encoding/json"
	"fmt"
)

// Option represents an optional value that may or may not be present
type Option[T any] struct {
	value T
	some  bool
}

// Some creates an Option with a value present
func Some[T any](value T) Option[T] {
	return Option[T]{value: value, some: true}
}

// None creates an Option with no value present
func None[T any]() Option[T] {
	return Option[T]{}
}

// IsSome returns true if the Option contains a value
func (o Option[T]) IsSome() bool {
	return o.some
}

// IsNone returns true if the Option contains no value
func (o Option[T]) IsNone() bool {
	return !o.some
}

// Unwrap returns the contained value
// Panics if the Option is None
func (o Option[T]) Unwrap() T {
	if !o.some {
		panic("called Unwrap on None Option")
	}
	return o.value
}

// UnwrapOr returns the contained value or a default value if None
func (o Option[T]) UnwrapOr(defaultValue T) T {
	if o.some {
		return o.value
	}
	return defaultValue
}

// UnwrapOrElse returns the contained value or calls a function to get a default
func (o Option[T]) UnwrapOrElse(f func() T) T {
	if o.some {
		return o.value
	}
	return f()
}

// String implements the Stringer interface
func (o Option[T]) String() string {
	if o.some {
		return fmt.Sprintf("Some(%v)", o.value)
	}
	return "None"
}

// Equal checks if two Options are equal
func (o Option[T]) Equal(other Option[T], equals func(T, T) bool) bool {
	if o.some != other.some {
		return false
	}
	if !o.some {
		return true // Both are None
	}
	return equals(o.value, other.value)
}

// MarshalJSON implements json.Marshaler
func (o Option[T]) MarshalJSON() ([]byte, error) {
	if o.some {
		return json.Marshal(o.value)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON implements json.Unmarshaler
func (o *Option[T]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*o = None[T]()
		return nil
	}

	var value T
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	*o = Some(value)
	return nil
}

// Ptr returns a pointer to the value if present, nil otherwise
// Useful for interacting with APIs that expect pointers for optional fields
func (o Option[T]) Ptr() *T {
	if o.some {
		return &o.value
	}
	return nil
}

// FromPtr creates an Option from a pointer
// Returns Some if the pointer is non-nil, None otherwise
func FromPtr[T any](ptr *T) Option[T] {
	if ptr != nil {
		return Some(*ptr)
	}
	return None[T]()
}
