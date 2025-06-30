package dynamic

import (
	"testing"

	"github.com/go-live-view/go-live-view/rend"
	"github.com/go-live-view/go-live-view/testutils"
)

func TestNotNil(t *testing.T) {
	tt := []testutils.TestCase{
		{
			Name: "notnil true",
			Node: NotNil(1, func() rend.Node {
				return testutils.TestString(`<div>Hello World</div>`)
			}),
			Expected: "<div>Hello World</div>",
		},
		{
			Name: "notnil false",
			Node: testutils.TestNode(NotNil(nil, func() rend.Node {
				return testutils.TestString(`<div>Hello World</div>`)
			})),
			Expected: "",
		},
		{
			Name: "notnil zero int",
			Node: NotNil(0, func() rend.Node {
				return testutils.TestString(`<div>Zero</div>`)
			}),
			Expected: "<div>Zero</div>",
		},
		{
			Name: "notnil false bool",
			Node: NotNil(false, func() rend.Node {
				return testutils.TestString(`<div>False</div>`)
			}),
			Expected: "<div>False</div>",
		},
		{
			Name: "notnil empty string",
			Node: NotNil("", func() rend.Node {
				return testutils.TestString(`<div>Empty</div>`)
			}),
			Expected: "<div>Empty</div>",
		},
		{
			Name: "notnil empty slice",
			Node: NotNil([]int{}, func() rend.Node {
				return testutils.TestString(`<div>EmptySlice</div>`)
			}),
			Expected: "<div>EmptySlice</div>",
		},
	}

	testutils.RunTestCases(t, tt, "notnil")
}
