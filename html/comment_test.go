package html

import (
	"testing"

	"github.com/go-live-view/go-live-view/testutils"
)

func TestComment(t *testing.T) {
	tt := []testutils.TestCase{
		{
			Name:     "simple comment",
			Node:     Comment("This is a comment"),
			Expected: "<!--This is a comment-->",
		},
		{
			Name: "comment in element",
			Node: Div(
				Comment("Start of content"),
				Text("Hello World"),
				Comment("End of content"),
			),
			Expected: "<div><!--Start of content-->Hello World<!--End of content--></div>",
		},
	}

	testutils.RunTestCases(t, tt, "comment")
}
