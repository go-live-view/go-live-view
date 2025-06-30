package uploads

import (
	"fmt"

	"github.com/go-live-view/go-live-view/dynamic"
	"github.com/go-live-view/go-live-view/html"
	lv "github.com/go-live-view/go-live-view/liveview"
	"github.com/go-live-view/go-live-view/params"
	"github.com/go-live-view/go-live-view/phx"
	"github.com/go-live-view/go-live-view/rend"
	"github.com/go-live-view/go-live-view/uploads"
)

type Live struct {
	uploads *uploads.Uploads
}

func New() *Live {
	u := uploads.New()
	u.AllowUpload("mydoc",
		uploads.WithAccept(".pdf"),
		uploads.WithAutoUpload(false),
		uploads.WithMaxEntries(1),
	)

	return &Live{
		uploads: u,
	}
}

func (l *Live) Mount(s lv.Socket, p params.Params) error {
	// Upload is already configured in New()
	if s != nil {
		// Any socket-specific setup can go here
	}
	return nil
}

func (l *Live) Event(s lv.Socket, event string, p params.Params) error {
	if event == "validate" {
		l.uploads.OnValidate(p)
	}

	if event == "save" {
		l.uploads.Consume("mydoc", func(path string, entry *uploads.Entry) {
			fmt.Printf("Consuming %s", entry.Meta.Name)
		})
	}

	return nil
}

func (l *Live) Uploads() *uploads.Uploads {
	return l.uploads
}

func (l *Live) Render(_ rend.Node) (rend.Node, error) {
	return dynamic.Component(
		html.Div(
			html.Form(
				html.Attr("id", "upload-form"),
				html.Attr("phx-submit", "save"),
				html.Attr("phx-change", "validate"),
				phx.FileInput(l.uploads.GetByName("mydoc")),
				html.Button(
					html.Attr("type", "submit"),
					html.Text("Upload"),
				),
			),
		)), nil
}
