package html

import (
	"fmt"
	"strings"

	"github.com/go-live-view/go-live-view/rend"
)

type AttributeNode struct {
	Tag    string
	Values []string
}

func Attr(tag string, values ...string) *AttributeNode {
	return &AttributeNode{
		Tag:    tag,
		Values: values,
	}
}

func (attr *AttributeNode) Render(diff bool, root *rend.Root, t *rend.Rend, b *strings.Builder) error {
	if len(attr.Values) == 0 {
		_, err := fmt.Fprintf(b, " %s", attr.Tag)
		return err
	}

	_, err := fmt.Fprintf(b, " %s=\"", attr.Tag)
	if err != nil {
		return err
	}

	for _, value := range attr.Values {
		_, err = b.WriteString(value)
		if err != nil {
			return err
		}
	}

	_, err = b.WriteString("\"")
	return err
}

type AttributesNode struct {
	Attrs []rend.Node
}

func Attrs(children ...rend.Node) *AttributesNode {
	return &AttributesNode{Attrs: children}
}

func (g *AttributesNode) Render(diff bool, root *rend.Root, t *rend.Rend, b *strings.Builder) error {
	for _, child := range g.Attrs {
		if child == nil {
			continue
		}

		if err := child.Render(diff, root, t, b); err != nil {
			return err
		}
	}
	return nil
}
