package dynamic

import (
	"fmt"
	"strings"

	"github.com/go-live-view/go-live-view/rend"
)

type TextNode struct {
	Value string
}

func Text(value string) rend.Node {
	return Wrap(&TextNode{value})
}

func Textf(format string, a ...any) rend.Node {
	return Wrap(&TextNode{fmt.Sprintf(format, a...)})
}

func (txt *TextNode) Render(diff bool, root *rend.Root, t *rend.Rend, b *strings.Builder) error {
	_, err := b.WriteString(txt.Value)
	return err
}
