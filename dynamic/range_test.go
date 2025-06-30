package dynamic

import (
	"testing"

	"github.com/go-live-view/go-live-view/html"
	"github.com/go-live-view/go-live-view/rend"
	"github.com/go-live-view/go-live-view/testutils"
)

func TestRange(t *testing.T) {
	tt := []testutils.TestCase{
		{
			Name: "simple range",
			Node: Range([]string{"a", "b", "c"}, func(s string) rend.Node {
				return testutils.TestString(`<div>` + s + `</div>`)
			}),
			Expected: "<div>a</div><div>b</div><div>c</div>",
		},
		{
			Name: "range with html",
			Node: html.Div(
				Range([]string{"a", "b", "c"}, func(s string) rend.Node {
					return testutils.TestString(`<div>Hello ` + s + `</div>`)
				}),
			),
			Expected: "<div><div>Hello a</div><div>Hello b</div><div>Hello c</div></div>",
		},
		{
			Name: "nested range",
			Node: Range([]string{"a", "b"}, func(s string) rend.Node {
				return Range([]string{"1", "2"}, func(n string) rend.Node {
					return testutils.TestString(`<div>` + s + n + `</div>`)
				})
			}),
			Expected: "<div>a1</div><div>a2</div><div>b1</div><div>b2</div>",
		},
		{
			Name: "range with component",
			Node: Range([]string{"a", "b", "c"}, func(s string) rend.Node {
				return Component(testutils.TestString(`<div>Hello ` + s + `</div>`))
			}),
			Expected: "<div>Hello a</div><div>Hello b</div><div>Hello c</div>",
		},
	}

	testutils.RunTestCases(t, tt, "range")
}
