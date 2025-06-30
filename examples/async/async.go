package async

import (
	"time"

	"github.com/go-live-view/go-live-view/async"
	"github.com/go-live-view/go-live-view/dynamic"
	"github.com/go-live-view/go-live-view/html"
	lv "github.com/go-live-view/go-live-view/liveview"
	"github.com/go-live-view/go-live-view/params"
	"github.com/go-live-view/go-live-view/rend"
)

type User struct {
	Name string
}

type Live struct {
	User *async.Async[*User]
}

func (l *Live) Mount(s lv.Socket, _ params.Params) error {
	l.User = async.New(s, func() (*User, error) {
		time.Sleep(2 * time.Second)
		return &User{Name: "John"}, nil
	})

	return nil
}

func (l *Live) Render(_ rend.Node) (rend.Node, error) {
	return html.Div(
		html.H1(
			dynamic.GoEmbed(func() rend.Node {
				switch l.User.State() {
				case async.Loading:
					loadingMessage := "Loading..."
					return dynamic.Text(loadingMessage)
				case async.Failed:
					err := l.User.Error().Error()
					return dynamic.Textf("failed to load user: %s", err)
				default:
					return dynamic.Text(l.User.Value().Name)
				}
			}),
		),
	), nil
}
