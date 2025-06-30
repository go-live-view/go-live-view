package dynamic

import (
	"strings"

	"github.com/go-live-view/go-live-view/rend"
)

type group struct {
	Children []rend.Node
}

func Group(children ...rend.Node) rend.Node {
	dynamicChildren := make([]rend.Node, len(children))

	for i, child := range children {
		if child == nil {
			continue
		}

		dynamicChildren[i] = Wrap(child)
	}

	return &group{Children: dynamicChildren}
}

func (group *group) Render(diff bool, root *rend.Root, t *rend.Rend, b *strings.Builder) error {
	for _, child := range group.Children {
		if child == nil {
			continue
		}

		if err := child.Render(diff, root, t, b); err != nil {
			return err
		}
	}
	return nil
}
