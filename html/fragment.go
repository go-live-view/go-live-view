package html

import (
	"strings"

	"github.com/go-live-view/go-live-view/rend"
)

type FragmentNode struct {
	Children []rend.Node
}

func Fragment(children ...rend.Node) *FragmentNode {
	return &FragmentNode{Children: children}
}

func (f *FragmentNode) Render(diff bool, root *rend.Root, t *rend.Rend, b *strings.Builder) error {
	for _, child := range f.Children {
		if child == nil {
			continue
		}

		if err := child.Render(diff, root, t, b); err != nil {
			return err
		}
	}
	return nil
}
