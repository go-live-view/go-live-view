package scroll

import (
	"fmt"

	"github.com/go-live-view/go-live-view/dynamic"
	"github.com/go-live-view/go-live-view/html"
	"github.com/go-live-view/go-live-view/internal/ref"
	lv "github.com/go-live-view/go-live-view/liveview"
	"github.com/go-live-view/go-live-view/params"
	"github.com/go-live-view/go-live-view/rend"
	"github.com/go-live-view/go-live-view/stream"
)

type Live struct {
	ref *ref.Ref

	userStream *stream.StreamGetter
}

type User struct {
	id   int
	Name string
}

func NewUser(id int) *User {
	return &User{id: id, Name: fmt.Sprintf("User %d", id)}
}

func (l *Live) Mount(s lv.Socket, _ params.Params) error {
	l.ref = ref.New(0)

	l.userStream = stream.New("users",
		stream.IDFunc(func(item any) string {
			user := item.(*User)
			return fmt.Sprintf("user-%d", user.id)
		}),
	)

	err := l.userStream.Add(
		NewUser(int(l.ref.NextRef())),
		NewUser(int(l.ref.NextRef())),
		NewUser(int(l.ref.NextRef())),
		NewUser(int(l.ref.NextRef())),
		NewUser(int(l.ref.NextRef())),
		NewUser(int(l.ref.NextRef())),
		NewUser(int(l.ref.NextRef())),
		NewUser(int(l.ref.NextRef())),
		NewUser(int(l.ref.NextRef())),
		NewUser(int(l.ref.NextRef())),
		NewUser(int(l.ref.NextRef())),
		NewUser(int(l.ref.NextRef())),
		NewUser(int(l.ref.NextRef())),
		NewUser(int(l.ref.NextRef())),
		NewUser(int(l.ref.NextRef())),
		NewUser(int(l.ref.NextRef())),
		NewUser(int(l.ref.NextRef())),
		NewUser(int(l.ref.NextRef())),
		NewUser(int(l.ref.NextRef())),
		NewUser(int(l.ref.NextRef())),
		NewUser(int(l.ref.NextRef())),
		NewUser(int(l.ref.NextRef())),
		NewUser(int(l.ref.NextRef())),
		NewUser(int(l.ref.NextRef())),
		NewUser(int(l.ref.NextRef())),
		NewUser(int(l.ref.NextRef())),
		NewUser(int(l.ref.NextRef())),
	)
	if err != nil {
		return fmt.Errorf("adding user in event: %w", err)
	}

	return nil
}

func (l *Live) Event(s lv.Socket, event string, p params.Params) error {
	if event == "add-user" {
		err := l.userStream.Add(
			NewUser(int(l.ref.NextRef())),
			NewUser(int(l.ref.NextRef())),
			NewUser(int(l.ref.NextRef())),
			NewUser(int(l.ref.NextRef())),
			NewUser(int(l.ref.NextRef())),
			NewUser(int(l.ref.NextRef())),
			NewUser(int(l.ref.NextRef())),
			NewUser(int(l.ref.NextRef())),
			NewUser(int(l.ref.NextRef())),
			NewUser(int(l.ref.NextRef())),
			NewUser(int(l.ref.NextRef())),
			NewUser(int(l.ref.NextRef())),
			NewUser(int(l.ref.NextRef())),
		)
		if err != nil {
			return fmt.Errorf("adding user in event: %w", err)
		}
	}

	if event == "delete-user" {
		err := l.userStream.Delete(p.Map("value").String("id"))
		if err != nil {
			return fmt.Errorf("deleting user in event: %w", err)
		}
	}

	return nil
}

func (l *Live) Render(_ rend.Node) (rend.Node, error) {
	return html.Div(
		html.Table(
			html.Tbody(
				html.Attrs(
					html.IdAttr("scroll-users"),
					html.Attr("phx-update", "stream"),
					html.Attr("phx-viewport-bottom", "add-user"),
				),
				dynamic.Stream(l.userStream.Get(), func(item stream.Item) rend.Node {
					u := item.Item.(*User)
					return html.Tr(
						html.Attrs(
							dynamic.Wrap(html.IdAttr(item.DomID)),
						),
						html.Td(
							dynamic.Text(u.Name),
							html.Button(
								html.Attrs(
									html.Attr("phx-click", "delete-user"),
									dynamic.Wrap(html.Attr("phx-value-id", item.DomID)),
								),
								html.Text("Delete"),
							),
						),
					)
				}),
			),
		),
	), nil
}
