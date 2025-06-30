package html

import (
	"strings"

	"github.com/go-live-view/go-live-view/rend"
)

type ElementNode struct {
	Tag      string
	Children []rend.Node
	Attrs    []rend.Node
}

func Element(tag string, children ...rend.Node) *ElementNode {
	e := &ElementNode{
		Tag: tag,
	}

	for _, child := range children {
		switch child.(type) {
		case *AttributesNode, *AttributeNode:
			e.Attrs = append(e.Attrs, child)
		default:
			e.Children = append(e.Children, child)
		}
	}

	return e
}

func (el *ElementNode) Render(diff bool, root *rend.Root, t *rend.Rend, b *strings.Builder) error {
	_, err := b.WriteString("<" + el.Tag)
	if err != nil {
		return err
	}

	for _, child := range el.Attrs {
		if child == nil {
			continue
		}

		if err := child.Render(diff, root, t, b); err != nil {
			return err
		}
	}

	_, err = b.WriteString(">")
	if err != nil {
		return err
	}

	for _, child := range el.Children {
		if child == nil {
			continue
		}

		if err := child.Render(diff, root, t, b); err != nil {
			return err
		}
	}

	_, err = b.WriteString("</" + el.Tag + ">")
	if err != nil {
		return err
	}

	return nil
}
