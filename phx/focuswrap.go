package phx

import (
	"github.com/go-live-view/go-live-view/html"
	"github.com/go-live-view/go-live-view/rend"
)

type FocusWrapOption func(*focusWrap)

type focusWrap struct {
	attrs    []rend.Node
	children []rend.Node
}

// FocusWrap creates a new focus wrap component with the given ID and options
func FocusWrap(id string, options ...FocusWrapOption) rend.Node {
	fw := &focusWrap{}

	for _, option := range options {
		option(fw)
	}

	return html.Div(
		html.IdAttr(id),
		html.Attr("phx-hook", "Phoenix.FocusWrap"),
		html.Attrs(fw.attrs...),
		html.Span(
			html.IdAttr(id+"-start"),
			html.TabindexAttr("0"),
			html.AriaAttr("hidden", "true"),
		),
		html.Fragment(fw.children...),
		html.Span(
			html.IdAttr(id+"-end"),
			html.TabindexAttr("0"),
			html.AriaAttr("hidden", "true"),
		),
	)
}

func FocusWrapAttr(attr rend.Node) FocusWrapOption {
	return func(fw *focusWrap) {
		fw.attrs = append(fw.attrs, attr)
	}
}

func FocusWrapSlot(children ...rend.Node) FocusWrapOption {
	return func(fw *focusWrap) {
		fw.children = children
	}
}
