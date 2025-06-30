package phx

import (
	"github.com/go-live-view/go-live-view/html"
	"github.com/go-live-view/go-live-view/rend"
)

// ImgPreview creates an image preview for LiveView uploads
func ImgPreview(entry string, attrs ...rend.Node) rend.Node {
	return html.Img(
		html.IdAttr("phx-preview-"+entry),
		html.DataAttr("phx-upload-ref", entry),
		html.DataAttr("phx-entry-ref", entry),
		html.DataAttr("phx-hook", "Phoenix.LiveImgPreview"),
		html.DataAttr("phx-update", "ignore"),
		html.Attrs(attrs...),
	)
}
