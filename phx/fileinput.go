package phx

import (
	"github.com/go-live-view/go-live-view/dynamic"
	"github.com/go-live-view/go-live-view/html"
	"github.com/go-live-view/go-live-view/rend"
	"github.com/go-live-view/go-live-view/uploads"
)

// FileInput creates a file input for LiveView uploads
func FileInput(u *uploads.Config, attrs ...rend.Node) rend.Node {
	if u == nil {
		return html.Input(
			html.TypeAttr("file"),
			html.Attr("disabled", "true"),
		)
	}

	return dynamic.Component(html.Input(
		html.Attrs(
			html.TypeAttr("file"),
			html.AcceptAttr(u.Accepts()),
			html.DataAttr("phx-hook", "Phoenix.LiveFileUpload"),
			html.DataAttr("phx-update", "ignore"),
			dynamic.Group(
				html.IdAttr(u.Ref),
				html.NameAttr(u.Name),
				html.DataAttr("phx-upload-ref", u.Ref),
				html.DataAttr("phx-active-refs", u.ActiveRefs()),
				html.DataAttr("phx-done-refs", u.DoneRefs()),
				html.DataAttr("phx-preflighted-refs", u.PreflightRefs()),
			),
			dynamic.If(u.MaxEntries > 1, html.MultipleAttr("true")),
			dynamic.If(u.AutoUpload, html.DataAttr("phx-auto-upload")),
		),
		html.Attrs(attrs...),
	))
}
