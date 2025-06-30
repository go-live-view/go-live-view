package broadcast

import (
	"time"

	"github.com/go-live-view/go-live-view/dynamic"
	"github.com/go-live-view/go-live-view/html"
	lv "github.com/go-live-view/go-live-view/liveview"
	"github.com/go-live-view/go-live-view/params"
	"github.com/go-live-view/go-live-view/rend"
)

type Live struct {
	Time time.Time
}

func New() *Live {
	return &Live{
		Time: time.Now(),
	}
}

func (l *Live) Mount(s lv.Socket, _ params.Params) error {
	if s != nil {
		go func() {
			time.Sleep(1 * time.Second)
			s.PushSelf("update", nil)
		}()
	}

	l.Time = time.Now()

	return nil
}

func (l *Live) Event(s lv.Socket, event string, _ params.Params) error {
	if event == "update" {
		go func() {
			time.Sleep(1 * time.Second)
			s.PushSelf("update", nil)
		}()

		l.Time = time.Now()
	}

	return nil
}

func (l *Live) Render(_ rend.Node) (rend.Node, error) {
	time := l.Time.Format("3:4:5 pm")
	return html.Div(
		html.H1(
			dynamic.Text(time),
		),
	), nil
}
