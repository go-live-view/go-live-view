package html

import (
	"strings"

	"github.com/go-live-view/go-live-view/rend"
)

type CommentNode struct {
	Comment string
}

func Comment(s string) *CommentNode {
	return &CommentNode{s}
}

func (c *CommentNode) Render(diff bool, root *rend.Root, t *rend.Rend, b *strings.Builder) error {
	b.WriteString("<!--")
	b.WriteString(c.Comment)
	b.WriteString("-->")
	return nil
}
