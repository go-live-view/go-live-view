package html

import (
	"testing"

	"github.com/go-live-view/go-live-view/testutils"
)

func TestAttribute(t *testing.T) {
	tt := []testutils.TestCase{
		{
			Name:     "non attribute children are ignored for void elements",
			Node:     Void("div", Div()),
			Expected: "<div/>",
		},
		{
			Name:     "simple void",
			Node:     Void("div", Attr("phx-click", "click")),
			Expected: "<div phx-click=\"click\"/>",
		},
		{
			Name: "static attributes",
			Node: Div(
				Attr("attr", "hello"),
				Attr("attr", "123"),
			),
			Expected: "<div attr=\"hello\" attr=\"123\"></div>",
		},
		{
			Name:     "valueless attribute",
			Node:     Div(Attr("attr")),
			Expected: "<div attr></div>",
		},
	}

	testutils.RunTestCases(t, tt, "attributes")
}

func TestAttrs(t *testing.T) {
	tt := []testutils.TestCase{
		{
			Name: "AttributesNode with multiple attributes",
			Node: Div(
				Attrs(Attr("class", "container"), Attr("id", "main")),
				Text("Content"),
			),
			Expected: "<div class=\"container\" id=\"main\">Content</div>",
		},
		{
			Name: "AttributesNode with valueless attributes",
			Node: Div(
				Attrs(Attr("disabled"), Attr("readonly")),
				Text("Content"),
			),
			Expected: "<div disabled readonly>Content</div>",
		},
		{
			Name: "mixed individual attributes and AttributesNode",
			Node: Div(
				Attr("data-test", "value"),
				Attrs(Attr("class", "container"), Attr("id", "main")),
				Text("Content"),
			),
			Expected: "<div data-test=\"value\" class=\"container\" id=\"main\">Content</div>",
		},
	}

	testutils.RunTestCases(t, tt, "attribute")
}
