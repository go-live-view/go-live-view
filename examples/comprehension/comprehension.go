package comprehension

import (
	"fmt"
	"strconv"

	"github.com/go-live-view/go-live-view/dynamic"
	"github.com/go-live-view/go-live-view/html"
	lv "github.com/go-live-view/go-live-view/liveview"
	"github.com/go-live-view/go-live-view/params"
	"github.com/go-live-view/go-live-view/rend"
)

type User struct {
	ID   int
	Name string
}

type Live struct {
	users []*User
}

func (l *Live) Event(s lv.Socket, event string, p params.Params) error {
	if event == "add-user" {
		l.users = append(l.users,
			&User{
				ID:   len(l.users) + 1,
				Name: fmt.Sprintf("User %d", len(l.users)+1),
			},
		)
	}

	if event == "delete-user" {
		id := p.Map("value").Int("id")

		for i, u := range l.users {
			if u.ID == id {
				l.users = append(l.users[:i], l.users[i+1:]...)
				break
			}
		}
	}

	return nil
}

func (l *Live) Render(_ rend.Node) (rend.Node, error) {
	return html.Div(
		html.Button(
			html.Text("Add User"),
			html.Attrs(
				html.Attr("phx-click", "add-user"),
			),
		),
		html.Table(
			html.Tbody(
				dynamic.Range(l.users, func(u *User) rend.Node {
					return html.Tr(
						html.Td(
							dynamic.Text(u.Name),
							html.Button(
								html.Attrs(
									html.Attr("phx-click", "delete-user"),
									dynamic.Wrap(html.Attr("phx-value-id", strconv.Itoa(u.ID))),
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
