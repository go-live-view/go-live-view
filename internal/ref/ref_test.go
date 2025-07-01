package ref

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tt := []struct {
		name  string
		start int64
	}{
		{"zero start", 0},
		{"positive start", 100},
		{"negative start", -50},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ref := New(tc.start)
			assert.NotNil(t, ref)

			// Test that the initial value is correct by calling NextRef
			// Since NextRef increments first, we expect start + 1
			next := ref.NextRef()
			assert.Equal(t, tc.start+1, next)
		})
	}
}

func TestNextRef(t *testing.T) {
	tt := []struct {
		name     string
		start    int64
		calls    int
		expected []int64
	}{
		{
			name:     "start from zero",
			start:    0,
			calls:    3,
			expected: []int64{1, 2, 3},
		},
		{
			name:     "start from 100",
			start:    100,
			calls:    3,
			expected: []int64{101, 102, 103},
		},
		{
			name:     "start from -5",
			start:    -5,
			calls:    3,
			expected: []int64{-4, -3, -2},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ref := New(tc.start)

			var results []int64
			for i := 0; i < tc.calls; i++ {
				results = append(results, ref.NextRef())
			}

			assert.Equal(t, tc.expected, results)
		})
	}
}

func TestNextStringRef(t *testing.T) {
	tt := []struct {
		name     string
		start    int64
		calls    int
		expected []string
	}{
		{
			name:     "start from zero",
			start:    0,
			calls:    3,
			expected: []string{"1", "2", "3"},
		},
		{
			name:     "start from negative",
			start:    -2,
			calls:    3,
			expected: []string{"-1", "0", "1"},
		},
		{
			name:     "large numbers",
			start:    999,
			calls:    2,
			expected: []string{"1000", "1001"},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ref := New(tc.start)

			var results []string
			for i := 0; i < tc.calls; i++ {
				results = append(results, ref.NextStringRef())
			}

			assert.Equal(t, tc.expected, results)
		})
	}
}
