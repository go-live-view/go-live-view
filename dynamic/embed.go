package dynamic

import (
	"github.com/go-live-view/go-live-view/rend"
)

// GoEmbed creates a dynamic node from a callback function
// The callback is evaluated each time the node is rendered
func GoEmbed(cb func() rend.Node) rend.Node {
	return Wrap(cb())
}
