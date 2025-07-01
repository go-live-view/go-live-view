package counter

import (
	"strconv"

	"github.com/go-live-view/go-live-view/dynamic"
	"github.com/go-live-view/go-live-view/html"
	lv "github.com/go-live-view/go-live-view/liveview"
	"github.com/go-live-view/go-live-view/params"
	"github.com/go-live-view/go-live-view/rend"
)

type Live struct {
	Count int
}

func (l *Live) Event(s lv.Socket, event string, _ params.Params) error {

	if event == "inc" {
		l.Count++
	}

	if event == "dec" {
		l.Count--
	}

	return nil
}

func (l *Live) Render(_ rend.Node) (rend.Node, error) {
	return html.Div(
		html.H1(
			dynamic.Text(strconv.Itoa(l.Count)),
		),
		html.Button(
			html.Text("inc"),
			html.Attr("phx-click", "inc"),
		),
		html.Button(
			html.Text("dec"),
			html.Attr("phx-click", "dec"),
		),
	), nil
}
