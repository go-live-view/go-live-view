package dynamic

import (
	"strings"

	"github.com/go-live-view/go-live-view/rend"
)

type WrapNode struct {
	Node rend.Node
}

func Wrap(node rend.Node) *WrapNode {
	return &WrapNode{Node: node}
}

func (c *WrapNode) Render(diff bool, root *rend.Root, t *rend.Rend, b *strings.Builder) error {
	if c.Node == nil {
		return nil
	}

	if diff {
		t.AddDynamic(rend.Render(root, c.Node))
		t.AddStatic(b.String())
		b.Reset()

		return nil
	}

	return c.Node.Render(diff, root, t, b)
}
