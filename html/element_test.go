package html

import (
	"testing"

	"github.com/go-live-view/go-live-view/testutils"
)

func TestElement(t *testing.T) {
	tt := []testutils.TestCase{
		{
			Name:     "simple element",
			Node:     Element("div", Text("Hello")),
			Expected: "<div>Hello</div>",
		},
		{
			Name:     "element with attributes",
			Node:     Element("div", Attr("class", "container"), Text("Content")),
			Expected: "<div class=\"container\">Content</div>",
		},
		{
			Name:     "element with AttributesNode",
			Node:     Element("div", Attrs(Attr("class", "container"), Attr("id", "main")), Text("Content")),
			Expected: "<div class=\"container\" id=\"main\">Content</div>",
		},
		{
			Name:     "element with mixed attributes and AttributesNode",
			Node:     Element("div", Attr("data-test", "value"), Attrs(Attr("class", "container"), Attr("id", "main")), Text("Content")),
			Expected: "<div data-test=\"value\" class=\"container\" id=\"main\">Content</div>",
		},
		{
			Name: "nested elements",
			Node: Element("div",
				Element("span", Text("Nested")),
			),
			Expected: "<div><span>Nested</span></div>",
		},
		{
			Name:     "empty element",
			Node:     Element("div"),
			Expected: "<div></div>",
		},
		{
			Name:     "nil children",
			Node:     Element("div", nil, Text("Hello"), nil),
			Expected: "<div>Hello</div>",
		},
	}

	testutils.RunTestCases(t, tt, "element")
}
