package dynamic

import (
	"strings"
	"testing"
)

func TestText(t *testing.T) {
	tt := []struct {
		name     string
		input    string
		expected string
	}{
		{"static text", "hello", "hello"},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			node := Text(tc.input)
			var b strings.Builder
			err := node.Render(false, nil, nil, &b)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if b.String() != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, b.String())
			}
		})
	}
}
