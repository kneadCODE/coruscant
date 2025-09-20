package optional

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSome(t *testing.T) {
	tests := []struct {
		name  string
		value interface{}
	}{
		{"string value", "hello"},
		{"int value", 42},
		{"bool value", true},
		{"struct value", struct{ Name string }{Name: "test"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := Some(tt.value)
			assert.True(t, opt.IsSome())
			assert.False(t, opt.IsNone())
			assert.Equal(t, tt.value, opt.Unwrap())
		})
	}
}

func TestNone(t *testing.T) {
	tests := []struct {
		name string
		opt  func() interface{}
	}{
		{"string None", func() interface{} { return None[string]() }},
		{"int None", func() interface{} { return None[int]() }},
		{"bool None", func() interface{} { return None[bool]() }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := tt.opt()
			switch o := opt.(type) {
			case Option[string]:
				assert.False(t, o.IsSome())
				assert.True(t, o.IsNone())
			case Option[int]:
				assert.False(t, o.IsSome())
				assert.True(t, o.IsNone())
			case Option[bool]:
				assert.False(t, o.IsSome())
				assert.True(t, o.IsNone())
			}
		})
	}
}

func TestUnwrap(t *testing.T) {
	t.Run("Some value unwraps successfully", func(t *testing.T) {
		opt := Some("hello")
		assert.Equal(t, "hello", opt.Unwrap())
	})

	t.Run("None panics on unwrap", func(t *testing.T) {
		opt := None[string]()
		assert.Panics(t, func() {
			opt.Unwrap()
		})
	})
}

func TestUnwrapOr(t *testing.T) {
	tests := []struct {
		name         string
		opt          Option[string]
		defaultValue string
		expected     string
	}{
		{
			name:         "Some returns value",
			opt:          Some("hello"),
			defaultValue: "default",
			expected:     "hello",
		},
		{
			name:         "None returns default",
			opt:          None[string](),
			defaultValue: "default",
			expected:     "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.opt.UnwrapOr(tt.defaultValue)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestUnwrapOrElse(t *testing.T) {
	t.Run("Some returns value", func(t *testing.T) {
		opt := Some("hello")
		result := opt.UnwrapOrElse(func() string { return "fallback" })
		assert.Equal(t, "hello", result)
	})

	t.Run("None calls function", func(t *testing.T) {
		opt := None[string]()
		called := false
		result := opt.UnwrapOrElse(func() string {
			called = true
			return "fallback"
		})
		assert.True(t, called)
		assert.Equal(t, "fallback", result)
	})
}

func TestString(t *testing.T) {
	tests := []struct {
		name     string
		opt      Option[interface{}]
		expected string
	}{
		{
			name:     "Some string",
			opt:      Some[interface{}]("hello"),
			expected: "Some(hello)",
		},
		{
			name:     "Some int",
			opt:      Some[interface{}](42),
			expected: "Some(42)",
		},
		{
			name:     "None",
			opt:      None[interface{}](),
			expected: "None",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.opt.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestEqual(t *testing.T) {
	equalString := func(a, b string) bool { return a == b }

	tests := []struct {
		name     string
		first    Option[string]
		second   Option[string]
		expected bool
	}{
		{
			name:     "Some equals Some with same value",
			first:    Some("hello"),
			second:   Some("hello"),
			expected: true,
		},
		{
			name:     "Some not equals Some with different value",
			first:    Some("hello"),
			second:   Some("world"),
			expected: false,
		},
		{
			name:     "None equals None",
			first:    None[string](),
			second:   None[string](),
			expected: true,
		},
		{
			name:     "Some not equals None",
			first:    Some("hello"),
			second:   None[string](),
			expected: false,
		},
		{
			name:     "None not equals Some",
			first:    None[string](),
			second:   Some("hello"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.first.Equal(tt.second, equalString)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestJSON(t *testing.T) {
	t.Run("Marshal Some", func(t *testing.T) {
		opt := Some("hello")
		data, err := json.Marshal(opt)
		require.NoError(t, err)
		assert.Equal(t, `"hello"`, string(data))
	})

	t.Run("Marshal None", func(t *testing.T) {
		opt := None[string]()
		data, err := json.Marshal(opt)
		require.NoError(t, err)
		assert.Equal(t, "null", string(data))
	})

	t.Run("Unmarshal to Some", func(t *testing.T) {
		var opt Option[string]
		err := json.Unmarshal([]byte(`"hello"`), &opt)
		require.NoError(t, err)
		assert.True(t, opt.IsSome())
		assert.Equal(t, "hello", opt.Unwrap())
	})

	t.Run("Unmarshal to None", func(t *testing.T) {
		var opt Option[string]
		err := json.Unmarshal([]byte("null"), &opt)
		require.NoError(t, err)
		assert.True(t, opt.IsNone())
	})

	t.Run("Round trip", func(t *testing.T) {
		original := Some(42)
		data, err := json.Marshal(original)
		require.NoError(t, err)

		var result Option[int]
		err = json.Unmarshal(data, &result)
		require.NoError(t, err)

		assert.True(t, result.IsSome())
		assert.Equal(t, 42, result.Unwrap())
	})
}

func TestPtr(t *testing.T) {
	t.Run("Some returns pointer", func(t *testing.T) {
		opt := Some("hello")
		ptr := opt.Ptr()
		require.NotNil(t, ptr)
		assert.Equal(t, "hello", *ptr)
	})

	t.Run("None returns nil", func(t *testing.T) {
		opt := None[string]()
		ptr := opt.Ptr()
		assert.Nil(t, ptr)
	})
}

func TestFromPtr(t *testing.T) {
	t.Run("Non-nil pointer", func(t *testing.T) {
		value := "hello"
		opt := FromPtr(&value)
		assert.True(t, opt.IsSome())
		assert.Equal(t, "hello", opt.Unwrap())
	})

	t.Run("Nil pointer", func(t *testing.T) {
		var ptr *string
		opt := FromPtr(ptr)
		assert.True(t, opt.IsNone())
	})
}

// Benchmark tests
func BenchmarkOption_Some(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Some(42)
	}
}

func BenchmarkOption_Unwrap(b *testing.B) {
	opt := Some(42)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = opt.Unwrap()
	}
}

func BenchmarkPointer_vs_Option(b *testing.B) {
	b.Run("Pointer access", func(b *testing.B) {
		value := 42
		ptr := &value
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if ptr != nil {
				_ = *ptr
			}
		}
	})

	b.Run("Option access", func(b *testing.B) {
		opt := Some(42)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if opt.IsSome() {
				_ = opt.Unwrap()
			}
		}
	})
}
