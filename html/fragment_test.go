package html

import (
	"testing"

	"github.com/go-live-view/go-live-view/testutils"
)

func TestFragment(t *testing.T) {
	tt := []testutils.TestCase{
		{
			Name: "simple fragment",
			Node: Fragment(
				Div(Text("Hello")),
				Span(Text("World")),
			),
			Expected: "<div>Hello</div><span>World</span>",
		},
		{
			Name: "fragment with attributes",
			Node: Fragment(
				Div(Attr("class", "container"), Text("Content")),
				Span(Attr("id", "label"), Text("Label")),
			),
			Expected: "<div class=\"container\">Content</div><span id=\"label\">Label</span>",
		},
	}

	testutils.RunTestCases(t, tt, "fragment")
}
