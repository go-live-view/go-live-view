package index

import (
	"fmt"

	"github.com/go-live-view/go-live-view/dynamic"
	"github.com/go-live-view/go-live-view/html"
	lv "github.com/go-live-view/go-live-view/liveview"
	"github.com/go-live-view/go-live-view/params"
	"github.com/go-live-view/go-live-view/rend"
)

type Live struct {
	Links []string

	flashes []string
}

func (l *Live) fectchFlashes() []string {
	flashes := l.flashes
	l.flashes = []string{}
	return flashes
}

func (l *Live) Mount(s lv.Socket, p params.Params) error {
	flash := p.Map("flash")
	if flash != nil {
		for k, v := range flash {
			l.flashes = append(l.flashes, fmt.Sprintf("%s: %s", k, v))
		}
	}

	delete(p, "flash")

	return nil
}

func (l *Live) Params(s lv.Socket, p params.Params) error {
	flash := p.Map("flash")
	if flash != nil {
		for k, v := range flash {
			l.flashes = append(l.flashes, fmt.Sprintf("%s: %s", k, v))
		}
	}

	return nil
}

func (l *Live) Event(_ lv.Socket, event string, p params.Params) error {
	if event == "lv:flash" {
		value := p.Map("value")
		key := value.String("key")
		msg := value.String("msg")
		l.flashes = append(l.flashes, fmt.Sprintf("%s: %s", key, msg))
	}
	return nil
}

func (i *Live) Render(child rend.Node) (rend.Node, error) {
	return html.Div(
		html.Div(
			html.Ol(
				dynamic.Range(i.fectchFlashes(), func(flash string) rend.Node {
					return html.Li(
						dynamic.Text(flash),
					)
				}),
			),
		),
		html.Ol(
			dynamic.Range(i.Links, func(link string) rend.Node {
				return html.Li(
					html.A(
						html.HrefAttr(link),
						dynamic.Text(link),
						html.DataAttr("phx-link", "patch"),
						html.DataAttr("phx-link-state", "push"),
					),
				)
			}),
		),
		html.Div(
			dynamic.Wrap(child),
		),
	), nil
}
