package dynamic

import (
	"testing"

	"github.com/go-live-view/go-live-view/html"
	"github.com/go-live-view/go-live-view/testutils"
)

func TestGroup(t *testing.T) {
	tt := []testutils.TestCase{
		{
			Name: "simple group",
			Node: html.Div(
				Group(
					testutils.TestString(`<div>Hello World 1</div>`),
					testutils.TestString(`<div>Hello World 2</div>`),
					testutils.TestString(`<div>Hello World 3</div>`),
				),
			),
			Expected: "<div><div>Hello World 1</div><div>Hello World 2</div><div>Hello World 3</div></div>",
		},
	}

	testutils.RunTestCases(t, tt, "group")
}
