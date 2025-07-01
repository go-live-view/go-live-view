package html

import (
	"testing"

	"github.com/go-live-view/go-live-view/rend"
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
			assert.Equal(t, tc.expected, rend.RenderString(Text(tc.input)))
		})
	}
}
