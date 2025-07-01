package router

import (
	"github.com/go-live-view/go-live-view/html"
	"github.com/go-live-view/go-live-view/rend"
)

type notFound struct{}

func (n *notFound) Render(_ rend.Node) (rend.Node, error) {
	return html.Div(
		html.Text("404 Not Found"),
	), nil
}
