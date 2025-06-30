package router

import (
	"github.com/go-live-view/go-live-view/dynamic"
	"github.com/go-live-view/go-live-view/html"
	"github.com/go-live-view/go-live-view/rend"
)

type notFound struct{}

func (n *notFound) Render(_ rend.Node) (rend.Node, error) {
	return html.Div(
		dynamic.Text("404 Not Found"),
	), nil
}
