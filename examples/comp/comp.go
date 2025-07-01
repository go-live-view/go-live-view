package comp

import (
	"fmt"

	"github.com/go-live-view/go-live-view/html"
	"github.com/go-live-view/go-live-view/rend"
)

func Unpkg(pkg, version string) rend.Node {
	return html.Script(
		html.Attrs(
			html.DeferAttr("true"),
			html.TypeAttr("text/javascript"),
			html.SrcAttr(
				fmt.Sprintf("https://unpkg.com/%s@%s", pkg, version),
			),
		),
	)
}

func RootLayout(children ...rend.Node) rend.Node {
	return html.Html(
		html.Head(
			Unpkg("phoenix", "1.7.14"),
			Unpkg("phoenix_live_view", "1.0.17"),
			Unpkg("topbar", "2.0.2"),
			Unpkg("apexcharts", "3.26.0"),
		),
		html.Body(
			html.Div(
				children...,
			),
			html.Script(
				html.Attrs(
					html.DeferAttr("true"),
					html.Attr("phx-track-static"),
					html.TypeAttr("text/javascript"),
					html.SrcAttr("/assets/app.js"),
				),
			),
		),
	)
}
