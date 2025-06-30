package dynamic

import (
	"strings"

	"github.com/go-live-view/go-live-view/rend"
)

type ComponentNode struct {
	Node rend.Node
}

func Component(root rend.Node) *ComponentNode {
	return &ComponentNode{Node: root}
}

func (c *ComponentNode) Render(diff bool, root *rend.Root, t *rend.Rend, b *strings.Builder) error {
	if c.Node == nil {
		return nil
	}

	if diff {
		t.AddComponent(root, rend.Render(root, c.Node))
		t.AddStatic(b.String())
		b.Reset()

		return nil
	}

	return c.Node.Render(diff, root, t, b)
}
