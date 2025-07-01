package dynamic

import (
	"testing"

	"github.com/go-live-view/go-live-view/rend"
	"github.com/go-live-view/go-live-view/testutils"
)

func TestIf(t *testing.T) {
	tt := []testutils.TestCase{
		{
			Name:     "if true",
			Node:     If(true, testutils.TestString(`<div>true</div>`)),
			Expected: "<div>true</div>",
		},
		{
			Name:     "if false",
			Node:     testutils.TestNode(If(false, testutils.TestString(`<div>true</div>`))),
			Expected: "",
		},
		{
			Name:     "if with string",
			Node:     If(true, Text("hello world")),
			Expected: "hello world",
		},
	}

	testutils.RunTestCases(t, tt, "if")
}

func TestIfElse(t *testing.T) {
	tt := []testutils.TestCase{
		{
			Name:     "ifelse true",
			Node:     IfElse(true, testutils.TestString(`<div>true</div>`), testutils.TestString(`<div>false</div>`)),
			Expected: "<div>true</div>",
		},
		{
			Name:     "ifelse false",
			Node:     IfElse(false, testutils.TestString(`<div>true</div>`), testutils.TestString(`<div>false</div>`)),
			Expected: "<div>false</div>",
		},
	}

	testutils.RunTestCases(t, tt, "ifelse")
}

func TestUnless(t *testing.T) {
	tt := []testutils.TestCase{
		{
			Name:     "unless true",
			Node:     testutils.TestNode(Unless(true, testutils.TestString(`<div>unless</div>`))),
			Expected: "",
		},
		{
			Name:     "unless false",
			Node:     Unless(false, testutils.TestString(`<div>unless</div>`)),
			Expected: "<div>unless</div>",
		},
	}

	testutils.RunTestCases(t, tt, "unless")
}

func TestSwitch(t *testing.T) {
	tt := []testutils.TestCase{
		{
			Name: "switch match",
			Node: Switch("a",
				Case("a", testutils.TestString(`<div>case a</div>`)),
				Case("b", testutils.TestString(`<div>case b</div>`)),
			),
			Expected: "<div>case a</div>",
		},
		{
			Name: "switch no match",
			Node: testutils.TestNode(Switch("c",
				Case("a", testutils.TestString(`<div>case a</div>`)),
				Case("b", testutils.TestString(`<div>case b</div>`)),
			)),
			Expected: "",
		},
	}

	testutils.RunTestCases(t, tt, "switch")
}

func TestSwitchDefault(t *testing.T) {
	tt := []testutils.TestCase{
		{
			Name: "switch default match",
			Node: SwitchDefault("a",
				testutils.TestString(`<div>default</div>`),
				Case("a", testutils.TestString(`<div>case a</div>`)),
				Case("b", testutils.TestString(`<div>case b</div>`)),
			),
			Expected: "<div>case a</div>",
		},
		{
			Name: "switch default no match",
			Node: SwitchDefault("c",
				testutils.TestString(`<div>default</div>`),
				Case("a", testutils.TestString(`<div>case a</div>`)),
				Case("b", testutils.TestString(`<div>case b</div>`)),
			),
			Expected: "<div>default</div>",
		},
	}

	testutils.RunTestCases(t, tt, "switch-default")
}

func TestWhen(t *testing.T) {
	tt := []testutils.TestCase{
		{
			Name: "when first true",
			Node: When(
				WhenIs(true, testutils.TestString(`<div>first</div>`)),
				WhenIs(true, testutils.TestString(`<div>second</div>`)),
			),
			Expected: "<div>first</div>",
		},
		{
			Name: "when second true",
			Node: When(
				WhenIs(false, testutils.TestString(`<div>first</div>`)),
				WhenIs(true, testutils.TestString(`<div>second</div>`)),
			),
			Expected: "<div>second</div>",
		},
	}

	testutils.RunTestCases(t, tt, "when")
}

func TestIfFunc(t *testing.T) {
	tt := []testutils.TestCase{
		{
			Name: "if func true",
			Node: IfFunc(true, func() rend.Node {
				return testutils.TestString(`<div>lazy true</div>`)
			}),
			Expected: "<div>lazy true</div>",
		},
		{
			Name: "if func false",
			Node: testutils.TestNode(IfFunc(false, func() rend.Node {
				return testutils.TestString(`<div>lazy false</div>`)
			})),
			Expected: "",
		},
		{
			Name: "if func returns nil",
			Node: testutils.TestNode(IfFunc(true, func() rend.Node {
				return nil
			})),
			Expected: "",
		},
	}

	testutils.RunTestCases(t, tt, "if-func")
}

func TestIfElseFunc(t *testing.T) {
	tt := []testutils.TestCase{
		{
			Name: "ifelse func true",
			Node: IfElseFunc(true,
				func() rend.Node { return testutils.TestString(`<div>true branch</div>`) },
				func() rend.Node { return testutils.TestString(`<div>false branch</div>`) },
			),
			Expected: "<div>true branch</div>",
		},
		{
			Name: "ifelse func false",
			Node: IfElseFunc(false,
				func() rend.Node { return testutils.TestString(`<div>true branch</div>`) },
				func() rend.Node { return testutils.TestString(`<div>false branch</div>`) },
			),
			Expected: "<div>false branch</div>",
		},
	}

	testutils.RunTestCases(t, tt, "ifelse-func")
}

func TestUnlessFunc(t *testing.T) {
	tt := []testutils.TestCase{
		{
			Name: "unless func true",
			Node: testutils.TestNode(UnlessFunc(true, func() rend.Node {
				return testutils.TestString(`<div>unless true</div>`)
			})),
			Expected: "",
		},
		{
			Name: "unless func false",
			Node: UnlessFunc(false, func() rend.Node {
				return testutils.TestString(`<div>unless false</div>`)
			}),
			Expected: "<div>unless false</div>",
		},
	}

	testutils.RunTestCases(t, tt, "unless-func")
}
