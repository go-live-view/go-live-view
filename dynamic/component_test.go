package dynamic

import (
	"testing"

	"github.com/go-live-view/go-live-view/testutils"
)

func TestComponent(t *testing.T) {
	testCases := []testutils.TestCase{
		{
			Name:     "simple component",
			Node:     Component(testutils.TestString(`<div>Hello World</div>`)),
			Expected: "<div>Hello World</div>",
		},
		{
			Name: "nested component",
			Node: Component(
				Component(testutils.TestString(`<span>Nested</span>`)),
			),
			Expected: "<span>Nested</span>",
		},
	}

	testutils.RunTestCases(t, testCases, "component")
}
