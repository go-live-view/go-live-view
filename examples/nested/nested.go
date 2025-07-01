package nested

import (
	"github.com/go-live-view/go-live-view/dynamic"
	"github.com/go-live-view/go-live-view/html"
	lv "github.com/go-live-view/go-live-view/liveview"
	"github.com/go-live-view/go-live-view/params"
	"github.com/go-live-view/go-live-view/rend"
)

type Live struct {
}

func (u *Live) Render(child rend.Node) (rend.Node, error) {
	return html.Div(
		html.H1(
			html.Text("Nested"),
		),
		html.Button(
			html.A(
				html.Text("Show"),
				html.HrefAttr("/nested/1"),
				html.DataAttr("phx-link", "patch"),
				html.DataAttr("phx-link-state", "push"),
			),
		),
		html.Button(
			html.A(
				html.Text("Edit"),
				html.HrefAttr("/nested/1/edit"),
				html.DataAttr("phx-link", "patch"),
				html.DataAttr("phx-link-state", "push"),
			),
		),
		child,
	), nil
}

type ShowLive struct {
	params params.Params
}

func (l *ShowLive) Params(s lv.Socket, p params.Params) error {
	l.params = p
	return nil
}

func (l *ShowLive) Render(_ rend.Node) (rend.Node, error) {
	id := l.params.String("id")
	return html.Div(
		html.H1(
			dynamic.Textf("Show %s", id),
		),
	), nil
}

type EditLive struct {
	params params.Params
}

func (l *EditLive) Params(s lv.Socket, p params.Params) error {
	l.params = p
	return nil
}

func (l *EditLive) Render(_ rend.Node) (rend.Node, error) {
	id := l.params.String("id")
	return html.Div(
		html.H1(
			dynamic.Textf("Edit %s", id),
		),
	), nil
}
