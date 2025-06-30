package html

import (
	"strings"

	"github.com/go-live-view/go-live-view/rend"
)

type VoidNode struct {
	Tag   string
	Attrs []rend.Node
}

func Void(tag string, children ...rend.Node) *VoidNode {
	v := &VoidNode{
		Tag: tag,
	}

	for _, attr := range children {
		if attr == nil {
			continue
		}

		switch attr.(type) {
		case *AttributesNode, *AttributeNode:
			v.Attrs = append(v.Attrs, attr)
		}
	}

	return v
}

func (v *VoidNode) Render(diff bool, root *rend.Root, t *rend.Rend, b *strings.Builder) error {
	_, err := b.WriteString("<" + v.Tag)
	if err != nil {
		return err
	}

	for _, child := range v.Attrs {
		if child == nil {
			continue
		}

		if err := child.Render(diff, root, t, b); err != nil {
			return err
		}
	}

	_, err = b.WriteString("/>")
	if err != nil {
		return err
	}

	return nil
}
