package phx

import (
	"github.com/go-live-view/go-live-view/html"
	"github.com/go-live-view/go-live-view/rend"
)

type TitleOption func(*title)

type title struct {
	prefix   string
	suffix   string
	def      string
	attrs    []rend.Node
	children []rend.Node
}

// Title creates a title with automatic prefix/suffix on page_title updates
func Title(options ...TitleOption) rend.Node {
	t := &title{}

	for _, option := range options {
		option(t)
	}

	return html.Title(
		html.DataAttr("prefix", t.prefix),
		html.DataAttr("default", t.def),
		html.DataAttr("suffix", t.suffix),
		html.Attr("phx-no-format"),
		html.Attrs(t.attrs...),
		html.Fragment(t.children...),
	)
}

func TitlePrefix(prefix string) TitleOption {
	return func(t *title) {
		t.prefix = prefix
	}
}

func TitleSuffix(suffix string) TitleOption {
	return func(t *title) {
		t.suffix = suffix
	}
}

func TitleDefault(def string) TitleOption {
	return func(t *title) {
		t.def = def
	}
}

func TitleAttr(attr rend.Node) TitleOption {
	return func(t *title) {
		t.attrs = append(t.attrs, attr)
	}
}

func TitleSlot(children ...rend.Node) TitleOption {
	return func(t *title) {
		t.children = children
	}
}
