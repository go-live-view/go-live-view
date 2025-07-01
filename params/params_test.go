package params

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFromAny(t *testing.T) {
	tt := []struct {
		name     string
		input    any
		expected Params
	}{
		{
			name:     "valid map[string]any",
			input:    map[string]any{"key": "value", "num": 42},
			expected: Params{"key": "value", "num": 42},
		},
		{
			name:     "empty map",
			input:    map[string]any{},
			expected: Params{},
		},
		{
			name:     "invalid type string",
			input:    "not a map",
			expected: Params{},
		},
		{
			name:     "invalid type int",
			input:    123,
			expected: Params{},
		},
		{
			name:     "nil input",
			input:    nil,
			expected: Params{},
		},
		{
			name:     "invalid type slice",
			input:    []string{"a", "b"},
			expected: Params{},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := FromAny(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestMerge(t *testing.T) {
	tt := []struct {
		name     string
		params   []Params
		expected Params
	}{
		{
			name:     "empty merge",
			params:   []Params{},
			expected: Params{},
		},
		{
			name:     "single param",
			params:   []Params{{"key": "value"}},
			expected: Params{"key": "value"},
		},
		{
			name: "merge two params",
			params: []Params{
				{"key1": "value1"},
				{"key2": "value2"},
			},
			expected: Params{"key1": "value1", "key2": "value2"},
		},
		{
			name: "merge with overwrite",
			params: []Params{
				{"key": "value1", "other": "keep"},
				{"key": "value2"},
			},
			expected: Params{"key": "value2", "other": "keep"},
		},
		{
			name: "merge multiple params",
			params: []Params{
				{"a": "1", "b": "2"},
				{"c": "3", "b": "overwrite"},
				{"d": "4"},
			},
			expected: Params{"a": "1", "b": "overwrite", "c": "3", "d": "4"},
		},
		{
			name:     "merge nil params",
			params:   []Params{nil, {"key": "value"}, nil},
			expected: Params{"key": "value"},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := Merge(tc.params...)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestSet(t *testing.T) {
	p := Params{}
	p.Set("key1", "value1")
	p.Set("key2", 42)
	p.Set("key1", "updated")

	expected := Params{
		"key1": "updated",
		"key2": 42,
	}
	assert.Equal(t, expected, p)

	// Test setting nil value
	p.Set("nil_key", nil)
	assert.Nil(t, p["nil_key"])
}

func TestMap(t *testing.T) {
	tt := []struct {
		name     string
		params   Params
		keys     []string
		expected Params
	}{
		{
			name:     "nested map[string]any",
			params:   Params{"nested": map[string]any{"inner": "value"}},
			keys:     []string{"nested"},
			expected: Params{"inner": "value"},
		},
		{
			name:     "nested map[string]string",
			params:   Params{"nested": map[string]string{"inner": "value"}},
			keys:     []string{"nested"},
			expected: Params{"inner": "value"},
		},
		{
			name: "nested map[any]any",
			params: Params{"nested": map[any]any{
				"str_key": "value1",
				123:       "value2",
			}},
			keys:     []string{"nested"},
			expected: Params{"str_key": "value1", "123": "value2"},
		},
		{
			name:     "missing key",
			params:   Params{"other": "value"},
			keys:     []string{"missing"},
			expected: Params{},
		},
		{
			name:     "invalid type",
			params:   Params{"key": "not a map"},
			keys:     []string{"key"},
			expected: Params{},
		},
		{
			name:     "multiple keys first found",
			params:   Params{"key2": map[string]any{"found": "value"}},
			keys:     []string{"key1", "key2"},
			expected: Params{"found": "value"},
		},
		{
			name:     "multiple keys none found",
			params:   Params{"other": "value"},
			keys:     []string{"key1", "key2"},
			expected: Params{},
		},
		{
			name:     "empty keys",
			params:   Params{"key": map[string]any{"inner": "value"}},
			keys:     []string{},
			expected: Params{},
		},
		{
			name:     "nil map value",
			params:   Params{"key": nil},
			keys:     []string{"key"},
			expected: Params{},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.params.Map(tc.keys...)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestSlice(t *testing.T) {
	tt := []struct {
		name     string
		params   Params
		keys     []string
		expected []Params
	}{
		{
			name: "slice of any",
			params: Params{"list": []any{
				map[string]any{"id": 1},
				map[string]any{"id": 2},
			}},
			keys:     []string{"list"},
			expected: []Params{{"id": 1}, {"id": 2}},
		},
		{
			name: "slice of maps",
			params: Params{"list": []map[string]any{
				{"id": 1},
				{"id": 2},
			}},
			keys:     []string{"list"},
			expected: []Params{{"id": 1}, {"id": 2}},
		},
		{
			name:     "missing key",
			params:   Params{"other": "value"},
			keys:     []string{"missing"},
			expected: []Params{},
		},
		{
			name:     "invalid type",
			params:   Params{"key": "not a slice"},
			keys:     []string{"key"},
			expected: []Params{},
		},
		{
			name:     "empty slice",
			params:   Params{"list": []any{}},
			keys:     []string{"list"},
			expected: []Params{},
		},
		{
			name: "mixed slice with invalid items",
			params: Params{"list": []any{
				map[string]any{"valid": 1},
				"invalid",
				map[string]any{"valid": 2},
			}},
			keys:     []string{"list"},
			expected: []Params{{"valid": 1}, {}, {"valid": 2}},
		},
		{
			name:     "multiple keys first found",
			params:   Params{"key2": []any{map[string]any{"found": "value"}}},
			keys:     []string{"key1", "key2"},
			expected: []Params{{"found": "value"}},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.params.Slice(tc.keys...)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestInt(t *testing.T) {
	tt := []struct {
		name     string
		params   Params
		keys     []string
		expected int
	}{
		{"string to int", Params{"key": "123"}, []string{"key"}, 123},
		{"int value", Params{"key": 456}, []string{"key"}, 456},
		{"int64 value", Params{"key": int64(789)}, []string{"key"}, 789},
		{"float64 value", Params{"key": 123.7}, []string{"key"}, 123},
		{"invalid string", Params{"key": "abc"}, []string{"key"}, 0},
		{"missing key", Params{"other": 123}, []string{"key"}, 0},
		{"invalid type", Params{"key": true}, []string{"key"}, 0},
		{"multiple keys first found", Params{"key2": 42}, []string{"key1", "key2"}, 42},
		{"negative string", Params{"key": "-123"}, []string{"key"}, -123},
		{"zero value", Params{"key": 0}, []string{"key"}, 0},
		{"empty string", Params{"key": ""}, []string{"key"}, 0},
		{"float32 value", Params{"key": float32(99.9)}, []string{"key"}, 0}, // should default
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.params.Int(tc.keys...)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestFloat32(t *testing.T) {
	tt := []struct {
		name     string
		params   Params
		keys     []string
		expected float32
	}{
		{"string to float32", Params{"key": "123.45"}, []string{"key"}, float32(123.45)},
		{"int value", Params{"key": 123}, []string{"key"}, float32(123)},
		{"int64 value", Params{"key": int64(456)}, []string{"key"}, float32(456)},
		{"float32 value", Params{"key": float32(789.1)}, []string{"key"}, float32(789.1)},
		{"float64 value", Params{"key": 123.7}, []string{"key"}, float32(123.7)},
		{"invalid string", Params{"key": "abc"}, []string{"key"}, float32(0)},
		{"missing key", Params{"other": 123.0}, []string{"key"}, float32(0)},
		{"negative value", Params{"key": "-456.78"}, []string{"key"}, float32(-456.78)},
		{"zero value", Params{"key": float32(0)}, []string{"key"}, float32(0)},
		{"invalid type", Params{"key": true}, []string{"key"}, float32(0)},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.params.Float32(tc.keys...)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestFloat64(t *testing.T) {
	tt := []struct {
		name     string
		params   Params
		keys     []string
		expected float64
	}{
		{"string to float64", Params{"key": "123.45"}, []string{"key"}, 123.45},
		{"int value", Params{"key": 123}, []string{"key"}, float64(123)},
		{"int64 value", Params{"key": int64(456)}, []string{"key"}, float64(456)},
		{"float32 value", Params{"key": float32(789.1)}, []string{"key"}, float64(float32(789.1))},
		{"float64 value", Params{"key": 123.7}, []string{"key"}, 123.7},
		{"invalid string", Params{"key": "abc"}, []string{"key"}, 0.0},
		{"missing key", Params{"other": 123.0}, []string{"key"}, 0.0},
		{"negative value", Params{"key": "-456.78"}, []string{"key"}, -456.78},
		{"zero value", Params{"key": 0.0}, []string{"key"}, 0.0},
		{"scientific notation", Params{"key": "1.23e2"}, []string{"key"}, 123.0},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.params.Float64(tc.keys...)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestString(t *testing.T) {
	tt := []struct {
		name     string
		params   Params
		keys     []string
		expected string
	}{
		{"string value", Params{"key": "hello"}, []string{"key"}, "hello"},
		{"int to string", Params{"key": 123}, []string{"key"}, "123"},
		{"int64 to string", Params{"key": int64(456)}, []string{"key"}, "456"},
		{"float64 to string", Params{"key": 123.45}, []string{"key"}, "123.45"},
		{"bool to string", Params{"key": true}, []string{"key"}, "true"},
		{"bool false to string", Params{"key": false}, []string{"key"}, "false"},
		{"missing key", Params{"other": "value"}, []string{"key"}, ""},
		{"invalid type", Params{"key": []int{1, 2, 3}}, []string{"key"}, ""},
		{"multiple keys first found", Params{"key2": "found"}, []string{"key1", "key2"}, "found"},
		{"empty string", Params{"key": ""}, []string{"key"}, ""},
		{"zero int", Params{"key": 0}, []string{"key"}, "0"},
		{"negative int", Params{"key": -123}, []string{"key"}, "-123"},
		{"nil value", Params{"key": nil}, []string{"key"}, ""},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.params.String(tc.keys...)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestBool(t *testing.T) {
	tt := []struct {
		name     string
		params   Params
		keys     []string
		expected bool
	}{
		{"bool true", Params{"key": true}, []string{"key"}, true},
		{"bool false", Params{"key": false}, []string{"key"}, false},
		{"string true", Params{"key": "true"}, []string{"key"}, true},
		{"string false", Params{"key": "false"}, []string{"key"}, false},
		{"non-false string", Params{"key": "anything"}, []string{"key"}, true},
		{"empty string", Params{"key": ""}, []string{"key"}, true}, // empty string != "false"
		{"int non-zero", Params{"key": 1}, []string{"key"}, true},
		{"int zero", Params{"key": 0}, []string{"key"}, false},
		{"int64 non-zero", Params{"key": int64(5)}, []string{"key"}, true},
		{"int64 zero", Params{"key": int64(0)}, []string{"key"}, false},
		{"float64 non-zero", Params{"key": 1.5}, []string{"key"}, true},
		{"float64 zero", Params{"key": 0.0}, []string{"key"}, false},
		{"missing key", Params{"other": true}, []string{"key"}, false},
		{"invalid type", Params{"key": []string{"test"}}, []string{"key"}, false},
		{"multiple keys first found", Params{"key2": true}, []string{"key1", "key2"}, true},
		{"nil value", Params{"key": nil}, []string{"key"}, false},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.params.Bool(tc.keys...)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestIntSlice(t *testing.T) {
	tt := []struct {
		name     string
		params   Params
		keys     []string
		expected []int
	}{
		{
			name:     "valid int slice",
			params:   Params{"key": []any{1, 2, 3}},
			keys:     []string{"key"},
			expected: []int{1, 2, 3},
		},
		{
			name:     "mixed types in slice",
			params:   Params{"key": []any{1, "not int", 3}},
			keys:     []string{"key"},
			expected: []int{1, 3},
		},
		{
			name:     "empty slice",
			params:   Params{"key": []any{}},
			keys:     []string{"key"},
			expected: []int{},
		},
		{
			name:     "missing key",
			params:   Params{"other": []any{1, 2, 3}},
			keys:     []string{"key"},
			expected: []int{},
		},
		{
			name:     "invalid type",
			params:   Params{"key": "not a slice"},
			keys:     []string{"key"},
			expected: []int{},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.params.IntSlice(tc.keys...)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestFloatSlice(t *testing.T) {
	tt := []struct {
		name     string
		params   Params
		keys     []string
		expected []float64
	}{
		{
			name:     "valid float slice",
			params:   Params{"key": []any{1.1, 2.2, 3.3}},
			keys:     []string{"key"},
			expected: []float64{1.1, 2.2, 3.3},
		},
		{
			name:     "mixed types in slice",
			params:   Params{"key": []any{1.1, "not float", 3.3}},
			keys:     []string{"key"},
			expected: []float64{1.1, 3.3},
		},
		{
			name:     "empty slice",
			params:   Params{"key": []any{}},
			keys:     []string{"key"},
			expected: []float64{},
		},
		{
			name:     "missing key",
			params:   Params{"other": []any{1.1, 2.2}},
			keys:     []string{"key"},
			expected: []float64{},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.params.FloatSlice(tc.keys...)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestStringSlice(t *testing.T) {
	tt := []struct {
		name     string
		params   Params
		keys     []string
		expected []string
	}{
		{
			name:     "valid string slice",
			params:   Params{"key": []any{"a", "b", "c"}},
			keys:     []string{"key"},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "mixed types in slice",
			params:   Params{"key": []any{"a", 123, "c"}},
			keys:     []string{"key"},
			expected: []string{"a", "c"},
		},
		{
			name:     "empty slice",
			params:   Params{"key": []any{}},
			keys:     []string{"key"},
			expected: []string{},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.params.StringSlice(tc.keys...)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestBoolSlice(t *testing.T) {
	tt := []struct {
		name     string
		params   Params
		keys     []string
		expected []bool
	}{
		{
			name:     "valid bool slice",
			params:   Params{"key": []any{true, false, true}},
			keys:     []string{"key"},
			expected: []bool{true, false, true},
		},
		{
			name:     "mixed types in slice",
			params:   Params{"key": []any{true, "not bool", false}},
			keys:     []string{"key"},
			expected: []bool{true, false},
		},
		{
			name:     "empty slice",
			params:   Params{"key": []any{}},
			keys:     []string{"key"},
			expected: []bool{},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.params.BoolSlice(tc.keys...)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestByteSlice(t *testing.T) {
	tt := []struct {
		name     string
		params   Params
		keys     []string
		expected []byte
	}{
		{
			name:     "string to bytes",
			params:   Params{"key": "hello"},
			keys:     []string{"key"},
			expected: []byte("hello"),
		},
		{
			name:     "byte slice",
			params:   Params{"key": []byte{1, 2, 3}},
			keys:     []string{"key"},
			expected: []byte{1, 2, 3},
		},
		{
			name:     "empty string",
			params:   Params{"key": ""},
			keys:     []string{"key"},
			expected: []byte{},
		},
		{
			name:     "missing key",
			params:   Params{"other": "value"},
			keys:     []string{"key"},
			expected: []byte{},
		},
		{
			name:     "invalid type",
			params:   Params{"key": 123},
			keys:     []string{"key"},
			expected: []byte{},
		},
		{
			name:     "multiple keys first found",
			params:   Params{"key2": "found"},
			keys:     []string{"key1", "key2"},
			expected: []byte("found"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.params.ByteSlice(tc.keys...)
			assert.Equal(t, tc.expected, result)
		})
	}
}

// Test edge cases and integration scenarios
func TestParamsIntegration(t *testing.T) {
	t.Run("complex nested structure", func(t *testing.T) {
		p := Params{
			"user": map[string]any{
				"id":   123,
				"name": "John",
				"profile": map[string]any{
					"age":    30,
					"active": true,
				},
			},
			"tags": []any{"go", "testing", "params"},
		}

		// Test nested map access
		user := p.Map("user")
		assert.Equal(t, 123, user.Int("id"))
		assert.Equal(t, "John", user.String("name"))

		// Test double nesting
		profile := user.Map("profile")
		assert.Equal(t, 30, profile.Int("age"))
		assert.True(t, profile.Bool("active"))

		// Test string slice
		tags := p.StringSlice("tags")
		assert.Equal(t, []string{"go", "testing", "params"}, tags)
	})

	t.Run("merge and access", func(t *testing.T) {
		p1 := Params{"a": 1, "b": "old"}
		p2 := Params{"b": "new", "c": true}

		merged := Merge(p1, p2)

		assert.Equal(t, 1, merged.Int("a"))
		assert.Equal(t, "new", merged.String("b"))
		assert.True(t, merged.Bool("c"))
	})

	t.Run("fallback key behavior", func(t *testing.T) {
		p := Params{"backup": 42}

		// Test multiple keys fallback
		result := p.Int("primary", "secondary", "backup")
		assert.Equal(t, 42, result)

		// Test all keys missing
		result = p.Int("missing1", "missing2", "missing3")
		assert.Equal(t, 0, result)
	})
}
