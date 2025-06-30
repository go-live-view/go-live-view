package ssnav

import (
	"github.com/go-live-view/go-live-view/dynamic"
	"github.com/go-live-view/go-live-view/html"
	lv "github.com/go-live-view/go-live-view/liveview"
	"github.com/go-live-view/go-live-view/params"
	"github.com/go-live-view/go-live-view/rend"
)

type Live struct {
}

func (u *Live) Event(s lv.Socket, e string, p params.Params) error {
	if e == "navigate" {
		return s.PushPatch(p.Map("value").String("href"))
	}
	return nil
}

func (u *Live) Render(child rend.Node) (rend.Node, error) {
	return html.Div(
		html.H1(
			html.Text("Server Navigation"),
		),
		html.Button(
			html.A(
				html.Text("Show"),
				html.Attr("phx-click", "navigate"),
				html.Attr("phx-value-href", "/ssnav/1"),
			),
		),
		html.Button(
			html.A(
				html.Text("Edit"),
				html.Attr("phx-click", "navigate"),
				html.Attr("phx-value-href", "/ssnav/1/edit"),
			),
		),
		child,
	), nil
}

type ShowLive struct {
	id string
}

func (l *ShowLive) Params(s lv.Socket, p params.Params) error {
	l.id = p.String("id")
	return nil
}

func (l *ShowLive) Render(_ rend.Node) (rend.Node, error) {
	return html.Div(
		html.H1(
			dynamic.Textf("Show %s", l.id),
		),
	), nil
}

type EditLive struct {
	id string
}

func (l *EditLive) Params(s lv.Socket, p params.Params) error {
	l.id = p.String("id")
	return nil
}

func (l *EditLive) Render(_ rend.Node) (rend.Node, error) {
	return html.Div(
		html.H1(
			dynamic.Textf("Edit %s", l.id),
		),
	), nil
}
