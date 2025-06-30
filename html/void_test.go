package html

import (
	"testing"

	"github.com/go-live-view/go-live-view/testutils"
)

func TestVoid(t *testing.T) {
	tt := []testutils.TestCase{
		{
			Name:     "void element with children (ignored)",
			Node:     Void("img", Text("This should be ignored")),
			Expected: "<img/>",
		},
	}

	testutils.RunTestCases(t, tt, "void")
}
