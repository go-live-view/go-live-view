package phx

import (
	"github.com/go-live-view/go-live-view/html"
	"github.com/go-live-view/go-live-view/rend"
)

type LinkOption func(*link)

type link struct {
	href     string
	method   string
	replace  bool
	attrs    []rend.Node
	children []rend.Node
}

// LinkNavigate creates a link that navigates to a LiveView
func LinkNavigate(to string, options ...LinkOption) rend.Node {
	l := &link{href: to}

	for _, option := range options {
		option(l)
	}

	return html.A(
		html.HrefAttr(l.href),
		html.DataAttr("phx-link", "redirect"),
		html.DataAttr("phx-link-state", linkState(l.replace)),
		html.Attrs(l.attrs...),
		html.Fragment(l.children...),
	)
}

// LinkPatch creates a link that patches the current LiveView
func LinkPatch(to string, options ...LinkOption) rend.Node {
	l := &link{href: to}

	for _, option := range options {
		option(l)
	}

	return html.A(
		html.HrefAttr(l.href),
		html.DataAttr("phx-link", "patch"),
		html.DataAttr("phx-link-state", linkState(l.replace)),
		html.Attrs(l.attrs...),
		html.Fragment(l.children...),
	)
}

// LinkHref creates a traditional link with browser navigation
func LinkHref(href string, options ...LinkOption) rend.Node {
	l := &link{href: href}

	for _, option := range options {
		option(l)
	}

	var dataAttrs []rend.Node
	if l.method != "" && l.method != "get" {
		dataAttrs = append(dataAttrs,
			html.DataAttr("method", l.method),
			html.DataAttr("to", l.href),
		)
	}

	return html.A(
		html.HrefAttr(l.href),
		html.Attrs(dataAttrs...),
		html.Attrs(l.attrs...),
		html.Fragment(l.children...),
	)
}

func linkState(replace bool) string {
	if replace {
		return "replace"
	}
	return "push"
}

func LinkMethod(method string) LinkOption {
	return func(l *link) {
		l.method = method
	}
}

func LinkReplace(replace bool) LinkOption {
	return func(l *link) {
		l.replace = replace
	}
}

func LinkAttr(attr rend.Node) LinkOption {
	return func(l *link) {
		l.attrs = append(l.attrs, attr)
	}
}

func LinkSlot(children ...rend.Node) LinkOption {
	return func(l *link) {
		l.children = children
	}
}
