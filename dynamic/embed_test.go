package dynamic

import (
	"testing"

	"github.com/go-live-view/go-live-view/rend"
	"github.com/go-live-view/go-live-view/testutils"
)

func TestEmbed(t *testing.T) {
	tt := []testutils.TestCase{
		{
			Name: "simple embed",
			Node: GoEmbed(func() rend.Node {
				return testutils.TestString(`<div>Hello World</div>`)
			}),
			Expected: "<div>Hello World</div>",
		},
	}

	testutils.RunTestCases(t, tt, "embed")
}
