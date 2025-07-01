package dynamic

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, b.String())
		})
	}
}
