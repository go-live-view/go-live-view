package router

import (
	"testing"

	"github.com/go-live-view/go-live-view/dynamic"
	"github.com/go-live-view/go-live-view/html"
	lv "github.com/go-live-view/go-live-view/liveview"
	"github.com/go-live-view/go-live-view/params"
	"github.com/go-live-view/go-live-view/rend"

	"github.com/stretchr/testify/assert"
)

func testLayout(_ ...rend.Node) rend.Node {
	return nil
}

type testLive struct {
	name string
}

func (t *testLive) Render(n rend.Node) (rend.Node, error) {
	return html.Div(
		dynamic.Text(t.name),
		n,
	), nil
}

type routes struct {
	path     string
	lv       lv.View
	children []routes
	opts     []routeOption
}

type handler interface {
	Handle(path string, view lv.View, opts ...routeOption) *route
	Group(path string, view lv.View, opts ...routeOption) *routeGroup
}

func TestGetRoute(t *testing.T) {
	tt := []struct {
		name           string
		path           string
		routes         []routes
		expected       string
		expectedParams params.Params
		expectedError  error
	}{
		{
			name: "not found",
			path: "/nonexistent",
			routes: []routes{
				{path: "/", lv: &testLive{name: "home"}},
			},
			expected:      "<div>404 Not Found</div>",
			expectedError: lv.NotFoundError,
		},
		{
			name: "partial match",
			path: "/test/nonexistent",
			routes: []routes{
				{path: "/test", lv: &testLive{name: "test"}},
			},
			expected:       "<div>test</div>",
			expectedParams: params.Params{},
		},
		{
			name: "simple routes",
			path: "/",
			routes: []routes{
				{path: "/", lv: &testLive{
					name: "test",
				}},
			},
			expected:       "<div>test</div>",
			expectedParams: params.Params{},
		},
		{
			name: "simple routes with route params",
			path: "/123",
			routes: []routes{
				{path: "/:id", lv: &testLive{
					name: "test",
				}},
			},
			expected: "<div>test</div>",
			expectedParams: params.Params{
				"id": "123",
			},
		},
		{
			name: "simple routes with wildcard",
			path: "/test/123",
			routes: []routes{
				{path: "/test/*", lv: &testLive{
					name: "test",
				}},
			},
			expected: "<div>test</div>",
			expectedParams: params.Params{
				"*": "123",
			},
		},
		{
			name: "simple routes with query params",
			path: "/test?a=1&b=2",
			routes: []routes{
				{path: "/test", lv: &testLive{
					name: "test",
				}},
			},
			expected: "<div>test</div>",
			expectedParams: params.Params{
				"a": "1",
				"b": "2",
			},
		},
		{
			name: "simple routes with extra params",
			path: "/",
			routes: []routes{
				{path: "/", lv: &testLive{
					name: "test",
				},
					opts: []routeOption{
						WithParams(map[string]interface{}{
							"extra": "extra",
						}),
					},
				},
			},
			expected: "<div>test</div>",
			expectedParams: map[string]any{
				"extra": "extra",
			},
		},
		{
			name: "nested routes with route params",
			path: "/test/123",
			routes: []routes{
				{path: "/test", lv: &testLive{
					name: "test",
				},
					children: []routes{
						{path: "/:id", lv: &testLive{
							name: "child",
						}},
					},
				},
			},
			expected: "<div>test<div>child</div></div>",
			expectedParams: params.Params{
				"id": "123",
			},
		},
		{
			name: "nested routes with extra params",
			path: "/test/child/deep",
			routes: []routes{
				{
					path: "/test",
					lv: &testLive{
						name: "test",
					},
					opts: []routeOption{
						WithParams(map[string]interface{}{
							"extra1": "extra 1",
						}),
					},
					children: []routes{
						{
							path: "/child",
							lv: &testLive{
								name: "child",
							},
							opts: []routeOption{
								WithParams(map[string]interface{}{
									"extra2": "extra 2",
								}),
							},
							children: []routes{
								{
									path: "/deep",
									lv: &testLive{
										name: "deep",
									},
									opts: []routeOption{
										WithParams(map[string]interface{}{
											"extra3": "extra 3",
										}),
									},
								},
							},
						},
					},
				},
			},
			expected: "<div>test<div>child<div>deep</div></div></div>",
			expectedParams: map[string]any{
				"extra1": "extra 1",
				"extra2": "extra 2",
				"extra3": "extra 3",
			},
		},
		{
			name: "nested routes",
			path: "/test/child",
			routes: []routes{
				{path: "/test", lv: &testLive{
					name: "test",
				},
					children: []routes{
						{path: "/child", lv: &testLive{
							name: "child",
						}},
					},
				},
			},
			expected:       "<div>test<div>child</div></div>",
			expectedParams: params.Params{},
		},
		{
			name: "deep nested routes",
			path: "/test/child/deep",
			routes: []routes{
				{path: "/test", lv: &testLive{
					name: "test",
				},
					children: []routes{
						{path: "/child", lv: &testLive{
							name: "child",
						},
							children: []routes{
								{path: "/deep", lv: &testLive{
									name: "deep",
								}},
							},
						},
					},
				},
			},
			expected:       "<div>test<div>child<div>deep</div></div></div>",
			expectedParams: params.Params{},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			rt := NewRouter(
				testLayout,
			)
			createRoutes(rt, tc.routes)
			route, err := rt.GetRoute(tc.path)
			assert.Equal(t, tc.expectedError, err)

			assert.Equal(t, tc.expectedParams, route.GetParams())

			node, err := route.GetView().Render(nil)
			if err != nil {
				t.Errorf("error rendering route: %v", err)
			}

			result := rend.RenderString(node)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestRoutable(t *testing.T) {
	tt := []struct {
		name     string
		from     string
		to       string
		routes   []routes
		expected bool
	}{
		{
			name: "same session",
			from: "/test",
			to:   "/test",
			routes: []routes{
				{path: "/test", lv: &testLive{
					name: "test",
				}},
			},
			expected: true,
		},
		{
			name: "same session, different path",
			from: "/test",
			to:   "/foo",
			routes: []routes{
				{path: "/test", lv: &testLive{
					name: "test",
				}},
				{path: "/foo", lv: &testLive{
					name: "foo",
				}},
			},
			expected: true,
		},
		{
			name: "different session",
			from: "/test",
			to:   "/foo",
			routes: []routes{
				{path: "/test", lv: &testLive{
					name: "test",
				},
					opts: []routeOption{
						WithSession("123"),
					},
				},
				{path: "/foo", lv: &testLive{
					name: "foo",
				},
					opts: []routeOption{
						WithSession("456"),
					},
				},
			},
			expected: false,
		},
		{
			name: "same parent, different session",
			from: "/test/foo",
			to:   "/test/bar",
			routes: []routes{
				{path: "/test", lv: &testLive{
					name: "test",
				},
					children: []routes{
						{path: "/foo", lv: &testLive{
							name: "foo",
						},
							opts: []routeOption{
								WithSession("123"),
							},
						},
						{path: "/bar", lv: &testLive{
							name: "bar",
						},
							opts: []routeOption{
								WithSession("456"),
							},
						},
					},
				},
			},
			expected: false,
		},
		{
			name: "group session",
			from: "/test/foo",
			to:   "/test/bar",
			routes: []routes{
				{path: "/test", lv: &testLive{
					name: "test",
				},
					opts: []routeOption{
						WithSession("123"),
					},
					children: []routes{
						{path: "/foo", lv: &testLive{
							name: "foo",
						}},
						{path: "/bar", lv: &testLive{
							name: "bar",
						}},
					},
				},
			},
			expected: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			rt := NewRouter(testLayout)
			createRoutes(rt, tc.routes)
			from, err := rt.GetRoute(tc.from)
			assert.NoError(t, err)
			to, err := rt.GetRoute(tc.to)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, rt.Routable(from, to))
		})
	}
}

func createRoutes(rt handler, routes []routes) {
	for _, route := range routes {
		if len(route.children) > 0 {
			createRoutes(rt.Group(route.path, route.lv, route.opts...), route.children)
		} else {
			rt.Handle(route.path, route.lv, route.opts...)
		}
	}
}
