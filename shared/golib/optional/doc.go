// Package optional provides a type-safe way to handle nullable values using Go generics.
//
// The Option[T] type eliminates nil pointer panics and makes nullability explicit in the type system.
// It follows functional programming patterns inspired by Rust's Option type and other languages.
//
// Basic Usage:
//
//	// Creating optional values
//	some := optional.Some("hello")
//	none := optional.None[string]()
//
//	// Checking for values
//	if some.IsSome() {
//		value := some.Unwrap() // Safe to unwrap
//	}
//
//	// Safe access with defaults
//	value := none.UnwrapOr("default")
//
// Functional Operations:
//
//	// Transform values if present
//	lengths := optional.Map(some, func(s string) int { return len(s) })
//
//	// Chain operations
//	result := optional.FlatMap(some, func(s string) optional.Option[int] {
//		if len(s) > 0 {
//			return optional.Some(len(s))
//		}
//		return optional.None[int]()
//	})
//
//	// Filter values
//	filtered := some.Filter(func(s string) bool { return len(s) > 3 })
//
// JSON Serialization:
//
//	The Option[T] type automatically handles JSON marshaling and unmarshaling:
//	- Some(value) serializes to the value itself
//	- None serializes to null
//
// Pointer Interoperability:
//
//	// Convert to/from pointers for legacy APIs
//	ptr := some.Ptr()        // Returns *string
//	opt := optional.FromPtr(ptr) // Returns Option[string]
//
// Performance:
//
//	Option[T] stores values inline, avoiding heap allocations for primitive types.
//	This provides better cache locality and performance compared to pointer-based approaches.
package optional
