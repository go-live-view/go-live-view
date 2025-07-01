package dynamic

import (
	"strings"

	"github.com/go-live-view/go-live-view/rend"
	s "github.com/go-live-view/go-live-view/stream"
)

type stream struct {
	stream *s.Stream
	f      func(s.Item) rend.Node
}

func Stream(s *s.Stream, f func(s.Item) rend.Node) *stream {
	return &stream{
		stream: s,
		f:      f,
	}
}

func (s *stream) Render(diff bool, root *rend.Root, t *rend.Rend, b *strings.Builder) error {
	if s.stream == nil {
		return nil
	}

	if len(s.stream.Deletions) == 0 && len(s.stream.Additions) == 0 && !s.stream.Reset {
		return nil
	}

	if !diff {
		for _, d := range s.stream.Additions {
			if err := s.f(d).Render(diff, root, t, b); err != nil {
				return err
			}
		}
		return nil
	}

	if len(s.stream.Additions) == 0 {
		stream := []any{
			root.NextStreamID(),
			[]any{},
			s.stream.Deletions,
		}

		if s.stream.Reset {
			stream = append(stream, true)
		}

		t.AddDynamic(&rend.Comprehension{
			Stream: stream,
		})
		t.AddStatic(b.String())
		b.Reset()
		return nil
	}

	rends := []*rend.Rend{}

	for _, d := range s.stream.Additions {
		rend := rend.Render(root, s.f(d))
		rends = append(rends, rend)
	}

	staticsMatch := true
	for i := 1; i < len(rends); i++ {
		if !compareStatics(rends[i].Static, rends[i-1].Static) {
			staticsMatch = false
			break
		}
	}

	inserts := []any{}
	for _, r := range s.stream.Additions {
		inserts = append(inserts, []any{
			r.DomID,
			r.StreamAt,
			r.Limit,
		})
	}

	if staticsMatch {
		stream := []any{
			root.NextStreamID(),
			inserts,
			s.stream.Deletions,
		}

		if s.stream.Reset {
			stream = append(stream, true)
		}

		t.AddDynamic(&rend.Comprehension{
			Static:      rends[0].Static,
			Fingerprint: rends[0].Fingerprint,
			Dynamics:    copyDynamics(rends),
			Stream:      stream,
		})
		t.AddStatic(b.String())
		b.Reset()
	} else {
		for _, r := range rends {
			t.AddDynamic(r)
			t.AddStatic(b.String())
			b.Reset()
		}
	}

	return nil
}
