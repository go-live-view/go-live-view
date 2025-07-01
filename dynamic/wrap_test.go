package dynamic

import (
	"testing"

	"github.com/go-live-view/go-live-view/testutils"
)

func TestDynamic(t *testing.T) {
	tt := []testutils.TestCase{
		{
			Name:     "wrap static node",
			Node:     Wrap(Text("hello")),
			Expected: "hello",
		},
		{
			Name:     "wrap nil",
			Node:     Wrap(nil),
			Expected: "",
		},
		{
			Name:     "wrap dynamic node",
			Node:     Wrap(Wrap(Text("nested"))),
			Expected: "nested",
		},
	}

	testutils.RunTestCases(t, tt, "wrap")
}
