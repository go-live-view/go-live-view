# Go LiveView

[![Go Reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/go-live-view/go-live-view)
[![Go Report Card](https://img.shields.io/badge/go%20report-A+-brightgreen?logo=go&style=flat-square)](https://goreportcard.com/report/github.com/go-live-view/go-live-view)
[![Coverage](https://img.shields.io/codecov/c/github/go-live-view/go-live-view?logo=codecov&style=flat-square)](https://codecov.io/gh/go-live-view/go-live-view)
[![Test Status](https://img.shields.io/github/actions/workflow/status/go-live-view/go-live-view/test.yml?logo=github&style=flat-square)](https://github.com/go-live-view/go-live-view/actions)
[![Release](https://img.shields.io/github/v/release/go-live-view/go-live-view?logo=github&style=flat-square)](https://github.com/go-live-view/go-live-view/releases/latest)
[![License](https://img.shields.io/github/license/go-live-view/go-live-view?logo=opensourceinitiative&logoColor=white&style=flat-square)](https://github.com/go-live-view/go-live-view/blob/main/LICENSE)


Go backend library for the [Phoenix LiveView](httpss://github.com/phoenixframework/live_view) JS client. It enables rich, real-time user experiences with server-rendered HTML, written in Go.

> [!IMPORTANT]
> This is still very much a work in progress. The API is not stable and is subject to change.

## Table of Contents

- [What is LiveView?](#what-is-liveview)
- [Getting Started](#getting-started)
- [HTML DSL](#html-dsl)
  - [Core Building Blocks](#core-building-blocks)
  - [Element & Attribute Helpers](#element--attribute-helpers)
- [Components](#components)
- [Dynamic Content & Performance (`dynamic` package)](#dynamic-content--performance-dynamic-package)
  - [How LiveView Optimizes Updates](#how-liveview-optimizes-updates)
  - [Dynamic Helpers](#dynamic-helpers)
  - [Loop Optimization: Comprehensions](#loop-optimization-comprehensions)
  - [Stream Optimization: `dynamic.Stream`](#stream-optimization-dynamicstream)
  - [Component Optimization: `dynamic.Component`](#component-optimization-dynamiccomponent)
  - [Tips](#tips)
- [Lifecycle](#lifecycle)
  - [HttpMount](#httpmount)
  - [Mount](#mount)
  - [Params](#params)
  - [Event](#event)
  - [Render](#render)
  - [Unmount](#unmount)
  - [Lifecycle Flow](#lifecycle-flow)
- [Parameters](#parameters)
- [Events](#events)
- [Uploads](#uploads)
- [JavaScript Integration](#javascript-integration)
  - [Client-Side Commands](#client-side-commands)
  - [JavaScript Hooks](#javascript-hooks)
  - [Optimistic Updates](#optimistic-updates)
- [Routing](#routing)
  - [Basic Setup](#basic-setup)
  - [Route Options](#route-options)
  - [Middleware](#middleware)
  - [Route Groups](#route-groups)
  - [Layout Functions](#layout-functions)
  - [Navigation](#navigation)
- [Examples](#examples)


## What is LiveView?

Imagine building rich, interactive web applications where clicking a button instantly updates the page, forms validate in real-time, and data changes are pushed to all connected users—all without writing a single line of JavaScript.

That's LiveView. It's a programming model that moves interactivity back to the server while delivering the smooth user experience people expect from modern web apps.

**How it works:**

- Your Go code runs on the server and maintains state
- A persistent WebSocket connects the browser to your server
- User interactions (clicks, form submissions) are sent to the server instantly
- **Optimistic updates** show immediate feedback while waiting for the server response
- The server updates its state and sends back only the HTML that changed
- The browser updates the page seamlessly, preserving user input and loading states

**The result?** You get the responsiveness of a single-page application with the simplicity of server-rendered HTML. No complex frontend state management, no REST APIs to build and maintain, no JavaScript frameworks to learn.

This Go LiveView implementation is compatible with the [Phoenix LiveView JavaScript client](https://hexdocs.pm/phoenix_live_view/js-interop.html), allowing you to leverage the battle-tested frontend while implementing your application logic in Go. The backend architecture is designed specifically for Go's strengths - using stateful structs for views, standard HTTP routing patterns, and Go's concurrency model.


## Getting Started

Here's a simple counter example that demonstrates the core concepts:

```go
// CounterLive represents our LiveView with its state
type CounterLive struct {
	Count int
}

// Event handles user interactions from the frontend
func (l *CounterLive) Event(s lv.Socket, event string, _ params.Params) error {
	switch event {
	case "inc":
		l.Count++
	case "dec":
		l.Count--
	}
	return nil
}

// Render generates the HTML for this LiveView
func (l *CounterLive) Render(_ rend.Node) (rend.Node, error) {
	return html.Div(
		html.H1(
			dynamic.Text(strconv.Itoa(l.Count)),
		),
		html.Button(
			html.Text("Increment"),
			html.Attr("phx-click", "inc"),
		),
		html.Button(
			html.Text("Decrement"),
			html.Attr("phx-click", "dec"),
		),
	), nil
}
```

When a user clicks a button, the `phx-click` attribute sends an event to the server, the `Event` method updates the state, and the page automatically re-renders with the new count.

## HTML DSL

go-live-view provides an HTML DSL for building your user interfaces. The core functionality is in the `html` package, with performance optimizations in the `dynamic` package.

### Core Building Blocks

The HTML DSL is built on a few core functions that handle all HTML generation:

| Function | Purpose | Example | Result |
|----------|---------|---------|--------|
| `html.Element(tag, ...)` | Standard HTML elements with opening/closing tags | `html.Element("div", html.Text("Hi"))` | `<div>Hi</div>` |
| `html.Void(tag, ...)` | Self-closing HTML elements | `html.Void("input", html.TypeAttr("text"))` | `<input type="text"/>` |
| `html.Attr(name, values...)` | Individual attributes | `html.Attr("class", "btn primary")` | `class="btn primary"` |
| `html.Attrs(...)` | Groups multiple attributes | `html.Attrs(html.IdAttr("main"), ...)` | Multiple attributes |
| `html.Text(content)` | Text nodes (raw output) | `html.Text("Hello!")` | `Hello!` |
| `html.Textf(format, ...)` | Formatted text nodes | `html.Textf("Hello, %s!", name)` | `Hello, John!` |
| `html.Comment(text)` | HTML comments | `html.Comment("Note")` | `<!--Note-->` |
| `html.Fragment(...)` | Groups elements without wrapper element | `html.Fragment(h1, p)` | `<h1>...</h1><p>...</p>` |

```go
import "github.com/go-live-view/go-live-view/html"

// Element - creates standard HTML elements with opening/closing tags
html.Element("div",
    html.Attr("class", "container"),
    html.Text("Content"),
) // <div class="container">Content</div>

// Void - creates self-closing HTML elements
html.Void("input",
    html.Attr("type", "text"),
    html.Attr("name", "email"),
) // <input type="text" name="email"/>

// Attr - creates individual attributes
html.Attr("class", "container")          // class="container"
html.Attr("class", "btn btn-primary")    // class="btn btn-primary"

// Attrs - groups multiple attributes
html.Attrs(
    html.Attr("id", "main"),
    html.Attr("class", "container"),
)

// Text - creates text nodes (no escaping - raw output)
html.Text("Hello, World!")

// Textf - creates formatted text nodes (no escaping - raw output)
html.Textf("Hello, %s!", name)

// Comment - creates HTML comments
html.Comment("This is a comment") // <!--This is a comment-->

// Fragment - groups multiple elements without a wrapper
html.Fragment(
    html.H1(html.Text("Title")),
    html.P(html.Text("Content")),
) // <h1>Title</h1><p>Content</p>
```

### Element & Attribute Helpers

The `html` package provides convenient functions for all standard HTML elements and attributes:

**Element Helpers:**
```go
// Container elements (use Element internally)
html.Div(...)      // html.Element("div", ...)
html.Span(...)     // html.Element("span", ...)
html.P(...)        // html.Element("p", ...)
html.H1(...)       // html.Element("h1", ...)
html.Button(...)   // html.Element("button", ...)
html.Form(...)     // html.Element("form", ...)
html.A(...)        // html.Element("a", ...)
...

// Self-closing elements (use Void internally)
html.Input(...)    // html.Void("input", ...)
html.Br()          // html.Void("br")
html.Hr()          // html.Void("hr")
html.Img(...)      // html.Void("img", ...)
...
```

**Attribute Helpers:**
```go
// Instead of html.Attr("class", "container")
html.ClassAttr("container")
html.IdAttr("main")
html.TypeAttr("text")
html.NameAttr("email")
html.ValueAttr("default")
html.PlaceholderAttr("Enter text...")

// Data attributes
html.DataAttr("user-id", "123") // data-user-id="123"

// ARIA attributes
html.AriaAttr("label", "Close button")
```

## Components

Create reusable components as regular Go functions:

```go
func UserCard(user User) html.ElementNode {
    return html.Div(
        html.ClassAttr("user-card"),
        html.H3(html.Text(user.Name)),
        html.P(html.Text(user.Email)),
    )
}

// Use in render
func (lv *MyLiveView)Render(_ rend.Node) (rend.Node, error) {
    return html.Div(
        UserCard(lv.CurrentUser),
        html.H2(html.Text("Other Users")),
        // ... more content
    ), nil
}
```

## Dynamic Content & Performance (`dynamic` package)

The `dynamic` package provides helpers that optimize how dynamic content is sent to the client. LiveView separates static HTML from dynamic values for efficient updates.

### How LiveView Optimizes Updates

LiveView splits templates into static and dynamic parts. When content changes, only the dynamic parts are sent to the client:

```go
// This template...
html.Div(
    html.H1(dynamic.Text(user.Name)),        // Dynamic
    html.P(html.Text("Welcome back!")),  // Static
)

// Initial render generates:
{
  "s": ["<div><h1>", "</h1><p>Welcome back!</p></div>"],
  "0": "John Doe"
}
```

When `user.Name` changes from "John Doe" to "Jane Smith", only the dynamic value is sent:

```json
{"0": "Jane Smith"}
```

The client merges this with the existing static parts to rebuild the HTML. This approach dramatically reduces bandwidth usage compared to sending entire HTML templates on every update.

### Dynamic Helpers

The `dynamic` package provides optimized helpers for dynamic content:

| Helper | Purpose | Example |
|--------|---------|---------|
| `dynamic.Text(value)` | Dynamic text content | `dynamic.Text(user.Name)` |
| `dynamic.Textf(format, args...)` | Formatted dynamic text | `dynamic.Textf("Hello, %s!", user.Name)` |
| `dynamic.If(condition, node)` | Conditional rendering | `dynamic.If(user.IsAdmin, adminButton)` |
| `dynamic.Unless(condition, node)` | Inverse conditional | `dynamic.Unless(user.IsGuest, navbar)` |
| `dynamic.IfElse(condition, ifNode, elseNode)` | If-else branching | `dynamic.IfElse(loggedIn, welcome, login)` |
| `dynamic.Switch(value, cases...)` | Multi-way branching | `dynamic.Switch(status, cases...)` |
| `dynamic.Range(slice, func)` | Loop over collections | `dynamic.Range(users, renderUser)` |
| `dynamic.Component(node)` | Component optimization | `dynamic.Component(userCard)` |
| `dynamic.Group(nodes...)` | Groups nodes without wrapper + wraps each with `dynamic.Wrap()` | `dynamic.Group(header, content)` |
| `dynamic.Wrap(node)` | Mark as dynamic | `dynamic.Wrap(html.Attr("class", cls))` |
| `dynamic.Stream(items, func)` | Render stream items | `dynamic.Stream(stream.Get(), renderItem)` |

```go
import "github.com/go-live-view/go-live-view/dynamic"

// Dynamic text (used for variables that change)
dynamic.Text(user.Name)
dynamic.Textf("Hello, %s!", user.Name)

// Conditional rendering
dynamic.If(user.IsAdmin, html.Button(html.Text("Admin Panel")))
dynamic.Unless(user.IsGuest, html.Nav(...))

// If-else
dynamic.IfElse(user.IsLoggedIn,
    html.Text("Welcome back!"),
    html.A(html.Attr("href", "/login"), html.Text("Login")),
)

// Switch statements
dynamic.Switch(status,
    dynamic.Case("active", html.Span(html.Text("Active"))),
    dynamic.Case("inactive", html.Span(html.Text("Inactive"))),
)

// Loops
dynamic.Range(items, func(item Item) rend.Node {
    return html.Li(dynamic.Text(item.Name))
})

// Groups multiple elements (like Fragment/Attrs) + makes each dynamic
dynamic.Group(
    html.H1(html.Text("Title")),        // Becomes dynamic.Wrap(html.H1(...))
    html.P(html.Text("Description")),   // Becomes dynamic.Wrap(html.P(...))
    html.Button(html.Text("Action")),   // Becomes dynamic.Wrap(html.Button(...))
)
```

**Understanding `dynamic.Group`:**

`dynamic.Group` combines two functions: **grouping** (like `html.Fragment`) and **dynamic wrapping**. It renders multiple nodes without a wrapper element, but also wraps each child with `dynamic.Wrap()`:

```go
// These are equivalent...

// Option 1: Manual wrapping
html.Fragment(
    dynamic.Wrap(html.H1(html.Text("Title"))),
    dynamic.Wrap(html.P(html.Text("Content"))),
    dynamic.Wrap(html.Button(html.Text("Action"))),
)

// Option 2: Using dynamic.Group (cleaner)
dynamic.Group(
    html.H1(html.Text("Title")),
    html.P(html.Text("Content")),
    html.Button(html.Text("Action")),
)

// Both render: <h1>Title</h1><p>Content</p><button>Action</button>
// But each element participates in LiveView's diff system
```

**When to use `dynamic.Group`:**
- You need multiple elements without a wrapper (like `Fragment` or `Attrs`)
- You want those elements to be dynamic (participate in diffing)
- You want cleaner code than manual `Fragment` + `Wrap` combinations

**Comparison:**
- `html.Fragment()` - Groups elements without wrapper
- `html.Attrs()` - Groups attributes + preserves attribute identity when nested in dynamics
- `dynamic.Group()` - Groups elements without wrapper + makes each dynamic

### Loop Optimization: Comprehensions

When using `dynamic.Range` for loops, LiveView automatically optimizes rendering using **comprehensions**. If all items in a loop have the same HTML structure (same "static" parts), LiveView can optimize the payload:

```go
// This loop...
html.Ul(
    dynamic.Range(users, func(user User) rend.Node {
        return html.Li(dynamic.Text(user.Name))
    }),
)

// Initial render generates a comprehension:
{
  "0": {
    "s": ["<li>", "</li>"],
    "d": [["John"], ["Jane"], ["Bob"]],
    "f": "abc123"
  }
}

// When appending new users, only the new dynamic values are sent:
{
  "0": {
    "d": [["John"], ["Jane"], ["Bob"], ["Alice"], ["Charlie"]],
  }
}

// When sorting alphabetically, the entire dynamic values array is updated:
{
  "0": {
    "d": [["Alice"], ["Bob"], ["Charlie"], ["Jane"], ["John"]],
  }
}
```

**Combining with components**: For complex list items with multiple changing fields, you can combine both optimizations by wrapping each item in `dynamic.Component`. This gives you comprehension benefits for the list structure plus component isolation for individual item updates (see the Component section below for details).

### Stream Optimization: `dynamic.Stream`

`dynamic.Stream` provides optimized rendering for dynamic lists that change over time. Unlike `dynamic.Range` which re-renders the entire list, streams only send changes (additions, updates, deletions) to the client:

```go
// Using dynamic.Stream for efficient list updates
html.Div(
    html.IdAttr("messages"),
    html.Attr("phx-update", "stream"), // Enable stream updates
    dynamic.Stream(l.messageStream.Get(), func(item stream.Item) rend.Node {
        message := item.Item.(*Message)
        return html.Div(
            html.IdAttr(item.DomID), // Required for stream tracking
            html.Span(dynamic.Text(message.User)),
            html.Text(": "),
            html.Span(dynamic.Text(message.Content)),
        )
    }),
)
```

**How it works**:
1. Each stream item has a unique DOM ID for client-side tracking
2. When items are added/updated/deleted, only those specific changes are sent
3. The client applies changes without re-rendering the entire list
4. Perfect for real-time applications like chat, feeds, or live data

**Stream vs Range comparison**:

| Scenario | `dynamic.Range` | `dynamic.Stream` |
|----------|-------------|--------------|
| **Static lists** | ✅ Optimal (comprehensions) | ❌ Unnecessary overhead |
| **Adding items** | ❌ Re-renders entire list | ✅ Sends only new items |
| **Updating items** | ❌ Re-renders entire list | ✅ Sends only changed items |
| **Real-time updates** | ❌ Inefficient for large lists | ✅ Designed for this |
| **Memory usage** | ❌ Grows with list size | ✅ Constant overhead |

**Use cases**:
- Chat messages or comments
- Live data feeds (stocks, metrics)
- Real-time notifications
- Infinite scroll implementations
- Any list where items are frequently added/updated

**Performance benefits**: For a 1000-item list where 3 items change, `dynamic.Stream` sends ~3 DOM operations vs `dynamic.Range` sending the entire list structure.

> **💡 Example:** See the [stream example](examples/stream) for complete stream implementation patterns.

### Component Optimization: `dynamic.Component`

`dynamic.Component` is a separate client-side rendering optimization. It isolates content updates and enables DOM skipping for unchanged components:

```go
// With dynamic.Component - isolates content updates for DOM skipping
dynamic.Component(
    html.Div(
        html.H3(dynamic.Text(user.Name)),
        html.P(dynamic.Text(user.Email)),
        html.Button(html.Text("Edit")),
    ),
)
```

**How it works**:
1. When a component doesn't change, the client skips parsing/updating the entire subtree
2. The morphdom library renders unchanged components with a skip marker
3. This provides significant performance gains for complex or frequently changing content

**Performance gains**: 3-30x improvements in client-side rendering, especially for lists where only some items change.

**Use cases**:
- List items that update independently
- Complex components that change infrequently
- Any content block with expensive DOM structure

**Combining with comprehensions**: For complex list items, you can combine both optimizations:

```go
// Complex list items with components + comprehensions
html.Ul(
    dynamic.Range(users, func(user User) rend.Node {
        return dynamic.Component(
            html.Li(
                html.H3(dynamic.Text(user.Name)),
                html.P(dynamic.Text(user.Email)),
                html.Span(dynamic.Text(user.Status)),
            ),
        )
    }),
)

// Payload structure with components in comprehensions:
{
  "0": {
    "s": ["<ul>", "</ul>"],
    "d": [[0], [1], [2]]
  },
  "c": {
    "0": {
      "s": ["<li><h3>", "</h3><p>", "</p><span>", "</span></li>"],
      "d": ["Alice", "alice@example.com", "Active"]
    },
    "1": {
      "s": ["<li><h3>", "</h3><p>", "</p><span>", "</span></li>"],
      "d": ["Bob", "bob@example.com", "Inactive"]
    },
    "2": {
      "s": ["<li><h3>", "</h3><p>", "</p><span>", "</span></li>"],
      "d": ["Charlie", "charlie@example.com", "Active"]
    }
  }
}

// When sorting, only the component order changes:
{
  "0": {
    "d": [[0], [2], [1]]
  }
}
```

> [!WARNING]
> `dynamic.Component` requires a single root element.

**Reference**: For technical details on LiveView's rendering optimizations, see [Dashbit's comprehensive explanation](https://dashbit.co/blog/latency-rendering-liveview).

> **💡 Example:** See the [comprehension example](examples/comprehension) for optimized list rendering patterns.

### Tips

1. **Use `html.Text()` for static content**, `dynamic.Text()` for dynamic content
2. **Wrap dynamic attributes** - When attributes are inside dynamic helpers (like `dynamic.If`), wrap them in `html.Attrs()` so they're properly identified as attributes:

```go
// ❌ Wrong - attribute won't be recognized inside dynamic.If
html.Div(
    dynamic.If(isActive, html.ClassAttr("active")),
    html.Text("Content"),
)

// ✅ Correct - wrap dynamic attributes in html.Attrs()
html.Div(
    html.Attrs(
        dynamic.If(isActive, html.ClassAttr("active")),
    ),
    html.Text("Content"),
)
```


## Lifecycle

A LiveView follows a predictable lifecycle with several key methods that you can implement:

| Method | When Called | Purpose |
|--------|-------------|---------|
| `HttpMount` | Initial HTTP request only | HTTP-specific setup (cookies, redirects) |
| `Mount` | HTTP request + WebSocket connect | Initialize state, auth checks |
| `Params` | After Mount + URL changes | Handle URL parameters |
| `Event` | User interactions | Handle phx-* events |
| `Render` | After Mount/Params/Event | Generate HTML |
| `Unmount` | Connection closes | Clean up resources |

### HttpMount

```go
func (l *MyLiveView) HttpMount(w http.ResponseWriter, r *http.Request, p params.Params) error {
    // Handle HTTP-specific setup (cookies, headers, etc.)
    // Only called during the initial HTTP request, not WebSocket
    if !isAuthenticated(r) {
        http.Redirect(w, r, "/login", http.StatusFound)
        return nil
    }
    return nil
}
```

`HttpMount` is called only during the initial HTTP request (before the WebSocket connection). Use this for HTTP-specific operations like setting cookies, headers, or redirects. It has access to the raw `http.ResponseWriter` and `*http.Request`.

### Mount

```go
func (l *MyLiveView) Mount(s lv.Socket, p params.Params) error {
    // Initialize your LiveView state
    l.Count = 0
    return nil
}
```

`Mount` is called when the LiveView is first created, both during the initial HTTP request and when the WebSocket connection is established. This is where you set up initial state, load data, or perform authentication checks.

### Params

```go
func (l *MyLiveView) Params(s lv.Socket, p params.Params) error {
    // Handle URL parameters and query strings
    id := p.String("id")
    l.loadData(id)
    return nil
}
```

`Params` is called after `Mount` and whenever the URL changes (via `phx-link` navigation). It receives URL parameters and query strings, making it ideal for loading data based on the current route.

### Event

```go
func (l *MyLiveView) Event(s lv.Socket, event string, p params.Params) error {
    // Handle user interactions
    switch event {
    case "save":
        return l.save(p)
    case "delete":
        return l.delete(p.String("id"))
    }
    return nil
}
```

`Event` handles user interactions triggered by `phx-*` attributes (like `phx-click`, `phx-submit`, etc.). This is where you modify state in response to user actions.

### Render

```go
func (l *MyLiveView) Render(_ rend.Node) (rend.Node, error) {
    // Generate HTML based on current state
    return html.Div(
        html.H1(dynamic.Text("My LiveView")),
        // ... more HTML
    ), nil
}
```

`Render` generates the HTML for your LiveView. It's called after `Mount`, `Params`, and `Event`, and whenever the LiveView needs to update the page.

### Unmount

```go
func (l *MyLiveView) Unmount() error {
    // Clean up resources
    l.closeConnections()
    l.cancelTimers()
    return nil
}
```

`Unmount` is called when the LiveView is being destroyed. Use this to clean up resources, close connections, or cancel background processes.

**When is Unmount called?**
- **Client-initiated**: WebSocket connection closes (user closes tab, navigates away, network disconnects)
- **Server-initiated**: When navigating between routes and the LiveView is no longer needed

### Lifecycle Flow

| Scenario | Flow | Result |
|----------|------|--------|
| **Initial Request** | `HttpMount` → `Mount` → `Params` → `Render` | Static HTML sent to browser |
| **WebSocket Connection** | `Mount` → `Params` → `Render` | LiveView becomes interactive |
| **URL Changes** | `Params` → `Render` | Page updates (via `phx-link`) |
| **User Interactions** | `Event` → `Render` | Page updates (via `phx-click`, etc.) |
| **Cleanup** | `Unmount` | Resources cleaned up |

## Parameters

Parameters come from URL routes, query strings, forms, and route defaults. Access them through the `params.Params` type:

```go
// Route: /users/:id?active=true
func (l *UserShowLive) Params(s lv.Socket, p params.Params) error {
    userID := p.String("id")        // Route parameter
    active := p.Bool("active")      // Query parameter
    return nil
}

// Event handling with form data
func (l *UserFormLive) Event(s lv.Socket, event string, p params.Params) error {
    if event == "save" {
        user := p.Map("user")       // Nested form data
        name := user.String("name")
        email := user.String("email")
    }
    return nil
}
```

**Available helpers:** `String()`, `Int()`, `Bool()`, `Float64()`, `Map()`, `Slice()`, `StringSlice()`, `IntSlice()`

**Route patterns:**
- `/users/:id` - Named parameter
- `/files/*` - Wildcard (captures remaining path)
- Query parameters automatically parsed from `?key=value`

> **💡 Example:** See the [counter example](examples/counter) for basic parameter handling and state management.

## Events

Events are how users interact with your LiveView. They're triggered by `phx-*` attributes in your HTML and handled by your `Event` method.

```go
func (l *MyLive) Event(s lv.Socket, event string, p params.Params) error {
    switch event {
    case "increment":
        l.Count++
    case "save_user":
        return l.saveUser(p.Map("user"))
    case "delete":
        return l.deleteItem(p.String("id"))
    }
    return nil
}

func (l *MyLive) Render(_ rend.Node) (rend.Node, error) {
    return html.Div(
        // Click events
        html.Button(
            html.Text("Click me"),
            html.Attr("phx-click", "increment"),
        ),

        // Form submission
        html.Form(
            html.Attr("phx-submit", "save_user"),
            html.Input(html.NameAttr("user[name]"), html.ValueAttr(l.User.Name)),
            html.Button(html.Text("Save"), html.TypeAttr("submit")),
        ),

        // Click with data
        html.Button(
            html.Text("Delete"),
            html.Attr("phx-click", "delete"),
            html.Attr("phx-value-id", "123"),
        ),
    ), nil
}
```

**Common event patterns:**

| Event | Usage | Purpose |
|-------|-------|---------|
| `phx-click` | Button/link clicks | Handle user clicks |
| `phx-submit` | Form submissions | Process form data |
| `phx-change` | Input changes | Real-time validation |
| `phx-value-*` | Pass data with events | Send additional data |
| `phx-target` | Target specific components | Scope event handling |

> **📖 Complete Reference:** See the [Phoenix LiveView Events documentation](https://hexdocs.pm/phoenix_live_view/bindings.html) for comprehensive coverage of all event types, modifiers, and advanced patterns. The JavaScript client behavior is identical.

## Uploads

File uploads in LiveView are handled through the `uploads` package, providing secure, chunked uploads with real-time progress.

```go
import "github.com/go-live-view/go-live-view/uploads"

type FileUploadLive struct {
    uploads *uploads.Uploads
}

func (l *FileUploadLive) Mount(s lv.Socket, p params.Params) error {
    l.uploads = uploads.New()
    l.uploads.AllowUpload("documents",
        uploads.WithAccept(".pdf,.doc,.docx"),
        uploads.WithMaxEntries(3),
        uploads.WithMaxFileSize(10*1024*1024), // 10MB
        uploads.WithAutoUpload(false),
    )
    return nil
}

func (l *FileUploadLive) Event(s lv.Socket, event string, p params.Params) error {
    switch event {
    case "validate":
        l.uploads.OnValidate(p) // Validate files on selection
    case "save":
        return l.uploads.Consume("documents", func(path string, entry *uploads.Entry) {
            // Process uploaded file
            fmt.Printf("Uploaded: %s to %s\n", entry.Meta.Name, path)
        })
    }
    return nil
}

func (l *FileUploadLive) Uploads() *uploads.Uploads {
    return l.uploads // Required for upload integration
}

func (l *FileUploadLive) Render(_ rend.Node) (rend.Node, error) {
    return html.Form(
        html.Attr("phx-submit", "save"),
        html.Attr("phx-change", "validate"),
        phx.FileInput(l.uploads.GetByName("documents")),
        html.Button(html.Text("Upload"), html.TypeAttr("submit")),
    ), nil
}
```

**External Uploads:** External uploads to cloud storage (S3, Google Cloud, etc.) are not yet implemented in go-live-view. The `WithExternal` option exists but is not functional. For cloud storage uploads, you'll need to implement your own solution or wait for this feature to be completed.

> **💡 Example:** See the [uploads example](examples/uploads) for a complete file upload implementation.

> **📖 Complete Reference:** See the [Phoenix LiveView Uploads documentation](https://hexdocs.pm/phoenix_live_view/uploads.html) for comprehensive coverage of upload configuration, progress tracking, and error handling.

## JavaScript Integration

While LiveView minimizes JavaScript needs, you can integrate custom JavaScript through hooks and the `js` package for client-side commands.

### JavaScript Hooks

Register custom JavaScript behavior that runs when elements mount/update:

```javascript
// app.js
let Hooks = {}

Hooks.Chart = {
  mounted() {
    const options = JSON.parse(this.el.dataset.options)
    this.chart = new ApexCharts(this.el, options)
    this.chart.render()
  },
  updated() {
    const options = JSON.parse(this.el.dataset.options)
    this.chart.updateSeries(options.series)
  }
}

const lv = new LiveView.LiveSocket("/live", Phoenix.Socket, {
  hooks: Hooks
})
```

Use hooks in your LiveView:

```go
func (l *ChartLive) Render(_ rend.Node) (rend.Node, error) {
    chartData := `{"series": [{"data": [1,2,3]}]}`
    return html.Div(
        html.Attr("phx-hook", "Chart"),
        html.Attr("data-options", chartData),
    ), nil
}
```

> **💡 Examples:** See the [charts example](examples/charts) for JavaScript hooks integration and the [js example](examples/js) for client-side commands.

### Client-side Commands

Client-side commands let you execute JavaScript operations directly from your Go code without writing custom JavaScript. They're executed immediately on the client when events trigger, perfect for UI interactions like animations, focus management, and DOM manipulation.

**When to use:** For immediate visual feedback (showing/hiding elements, animations) without requiring a server round-trip.

```go
import "github.com/go-live-view/go-live-view/js"

// Hide elements with transitions
func hideModal() string {
    return js.JS(
        js.Hide(&js.HideArgs{
            To:         "#modal",
            Transition: [3]string{"fade-out", "duration-300", "ease-out"},
        }),
        js.Dispatch("modal:closed", &js.DispatchArgs{
            Detail: map[string]any{"modal_id": "user-modal"},
        }),
    )
}

// Show elements
func showToast() string {
    return js.JS(
        js.Show(&js.ShowArgs{
            To:         "#toast",
            Transition: [3]string{"fade-in", "duration-200"},
        }),
    )
}

// Navigate without server round-trip
func quickNavigation() string {
    return js.JS(
        js.Push("analytics", &js.PushArgs{
            Value: map[string]any{"action": "quick_nav"},
        }),
        js.Navigate("/dashboard"),
    )
}

func (l *ModalLive) Render(_ rend.Node) (rend.Node, error) {
    return html.Div(
        // Immediate UI response
        html.Button(
            html.Text("Close Modal"),
            html.Attr("phx-click", hideModal()),
        ),
        html.Button(
            html.Text("Show Toast"),
            html.Attr("phx-click", showToast()),
        ),

        // Modal that closes on outside click
        html.Div(
            html.IdAttr("modal"),
            html.Attr("phx-click-away", hideModal()),
            html.Attr("phx-window-keydown", hideModal()),
            html.Attr("phx-key", "escape"),
            html.Text("Modal content"),
        ),

        // Toast notification
        html.Div(
            html.IdAttr("toast"),
            html.ClassAttr("hidden"),
            html.Text("Success!"),
        ),
    ), nil
}
```

**Common commands:**

| Command | Purpose | Example Use |
|---------|---------|-------------|
| `Hide/Show/Toggle` | Control element visibility with CSS transitions | Modals, dropdowns, alerts |
| `Push` | Send events to server | Analytics, tracking |
| `Dispatch` | Trigger custom DOM events | Third-party libraries |
| `Navigate/Patch` | Client-side navigation | Faster than server routing |

**Key benefits:**
- **Immediate feedback** - No server round-trip delay
- **Smooth animations** - CSS transitions work perfectly
- **Composable** - Chain multiple commands together
- **No custom JavaScript** - All defined in your Go code

### Optimistic Updates

LiveView automatically handles the tricky problem of keeping client and server state in sync. When users interact with forms or buttons, the client needs to show immediate feedback while waiting for the server response, without losing user input or creating jarring state changes.

**The Problem:** User types in a form, presses submit (button shows "SUBMITTING"), but a delayed server update from an earlier event could roll back the form to its previous state, overwriting the user's current input and button state.

**LiveView's Solution:**
- **Input values** - The client is always the source of truth for current input values
- **Event tracking** - LiveView tracks in-flight events and only applies server updates after all pending events resolve
- **Loading states** - Automatic CSS classes provide immediate visual feedback

#### Automatic Loading Classes

LiveView automatically adds CSS classes to elements while events are processing:

```go
func (l *FormLive) Render(_ rend.Node) (rend.Node, error) {
    return html.Form(
        html.Attr("phx-submit", "save"),
        html.Input(
            html.NameAttr("email"),
            html.Attr("phx-change", "validate"),
            html.ClassAttr("phx-change-loading:opacity-50"), // Tailwind variant
        ),
        html.Button(
            html.Text("Submit"),
            html.TypeAttr("submit"),
            html.Attr("phx-disable-with", "Submitting..."), // Changes button text
            html.ClassAttr("phx-submit-loading:opacity-75"),
        ),
    ), nil
}
```

**Available loading classes:**
- `phx-click-loading` - Applied during click events
- `phx-submit-loading` - Applied during form submission
- `phx-change-loading` - Applied during input changes
- `phx-focus-loading`, `phx-blur-loading` - Applied during focus events

#### Advanced Optimistic UI with JS Commands

For more complex scenarios, combine server events with immediate client-side updates:

```go
func deleteWithOptimism() string {
    return js.JS(
        // Immediately fade out the row
        js.AddClass("opacity-50 transition-opacity", &js.AddClassArgs{
            To: "#user-row-123",
        }),
        // Then send the delete event to server
        js.Push("delete_user", &js.PushArgs{
            Value: map[string]any{"user_id": "123"},
            Loading: "#user-row-123", // Show loading state on the row
        }),
    )
}

func (l *UserListLive) Render(_ rend.Node) (rend.Node, error) {
    return html.Div(
        html.Div(
            html.IdAttr("user-row-123"),
            html.Text("John Doe"),
            html.Button(
                html.Text("Delete"),
                html.Attr("phx-click", deleteWithOptimism()),
            ),
        ),
    ), nil
}
```

**Benefits of LiveView's approach:**
- **Automatic sync** - No manual state management needed
- **No race conditions** - Event ordering is handled automatically
- **Preserved user input** - Form values never get lost during updates
- **Smooth UX** - Immediate feedback with eventual server consistency

> **📖 Complete Reference:**
>
> See the [Phoenix LiveView Syncing Changes documentation](https://hexdocs.pm/phoenix_live_view/syncing-changes.html) for comprehensive coverage of optimistic UIs, loading states, and client-server synchronization patterns.
>
> And [Phoenix LiveView JavaScript Interoperability documentation](https://hexdocs.pm/phoenix_live_view/js-interop.html) for comprehensive coverage of hooks, client-side commands, and advanced JavaScript integration patterns.

## Routing

The router defines which URLs map to which LiveViews. A key principle is that each connection gets its own LiveView instance to prevent users from sharing state.

### Basic Setup

```go
func main() {
    ctx := context.Background()

    // Setup HTTP handler with inline router callback
    mux := http.NewServeMux()
    mux.Handle("/", handler.NewHandler(ctx, func() lv.Router {
        rt := router.NewRouter(Layout)
        rt.Handle("/", &CounterLive{})
        rt.Handle("/users", &UserListLive{})
        return rt
    }))

    http.ListenAndServe(":8080", mux)
}
```

> [!NOTE]
> The `handler.NewHandler` takes a callback that returns a router. This ensures that each WebSocket connection calls the callback to create fresh LiveView instances (`&CounterLive{}`, `&UserListLive{}`), keeping user state completely separate.

### Middleware

Routes support middleware functions that run before your LiveView methods. There are two types:

**HTTP Middleware** - Runs only during the initial HTTP request:
```go
func authMiddleware(w http.ResponseWriter, r *http.Request, p params.Params) error {
    if !isAuthenticated(r) {
        http.Redirect(w, r, "/login", http.StatusFound)
        return nil
    }
    return nil
}

// Apply to specific routes
rt.Handle("/admin", &AdminLive{}, router.WithHttpMount(authMiddleware))

// Apply to route groups (inherited by all child routes)
admin := rt.Group("/admin", &AdminLayoutLive{}, router.WithHttpMount(authMiddleware))
admin.Handle("/users", &AdminUsersLive{})
admin.Handle("/settings", &AdminSettingsLive{})
```

**LiveView Middleware** - Runs for both HTTP and WebSocket connections:
```go
func loadUserMiddleware(s lv.Socket, p params.Params) error {
    userID := getCurrentUser(s)
    p["current_user"] = loadUser(userID)
    return nil
}

rt.Handle("/dashboard", &DashboardLive{}, router.WithMount(loadUserMiddleware))
```

### Route Options

**WithParams** - Add default parameters available to all routes:
```go
rt.Handle("/api", &APILive{},
    router.WithParams(params.Params{"api_version": "v1", "feature_flags": true}),
)
```

**WithSession** - Group routes by session to control navigation boundaries:
```go
// Public routes (default session)
rt.Handle("/", &HomeLive{})
rt.Handle("/login", &LoginLive{})

// User routes (user session)
user := rt.Group("/app", &AppLayoutLive{}, router.WithSession("user"))
user.Handle("/dashboard", &DashboardLive{})
user.Handle("/profile", &ProfileLive{})

// Admin routes (admin session)
admin := rt.Group("/admin", &AdminLayoutLive{}, router.WithSession("admin"))
admin.Handle("/users", &AdminUsersLive{})
admin.Handle("/settings", &AdminSettingsLive{})
```

**How sessions work:**
- Routes within the same session can navigate between each other using `phx-link` (seamless updates)
- Routes with different sessions cannot patch to each other - they trigger a full page redirect instead
- This breaks the WebSocket connection and forces a new HTTP request
- Use sessions to create logical boundaries in your application (public vs authenticated vs admin areas)

Sessions don't provide authentication by themselves - they control navigation behavior. Combine them with middleware for actual security.

### Route Groups

Route groups allow you to nest LiveViews and share middleware. The `view` parameter in `Group()` creates a parent LiveView that wraps child routes:

```go
rt := router.NewRouter(RootLayout)

// App section with shared navigation
app := rt.Group("/app", &AppSectionLive{})
app.Handle("/dashboard", &DashboardLive{})
app.Handle("/profile", &ProfileLive{})

// Admin section with different navigation
admin := rt.Group("/admin", &AdminSectionLive{})
admin.Handle("/users", &AdminUsersLive{})
admin.Handle("/settings", &AdminSettingsLive{})

// Nested groups create layered view composition
users := admin.Group("/users", &UsersSectionLive{})
users.Handle("/:id", &UserShowLive{})     // Nested: AdminSectionLive > UsersSectionLive > UserShowLive
users.Handle("/:id/edit", &UserEditLive{}) // Nested: AdminSectionLive > UsersSectionLive > UserEditLive
```

Parent LiveViews receive child content through their `Render` method:

```go
type AppSectionLive struct{}

func (l *AppSectionLive) Render(children rend.Node) (rend.Node, error) {
    return html.Div(
        html.ClassAttr("app-section"),
        html.Nav(
            html.A(html.HrefAttr("/app/dashboard"), html.Text("Dashboard")),
            html.A(html.HrefAttr("/app/profile"), html.Text("Profile")),
        ),
        html.Main(children), // Child routes render here
    ), nil
}
```

**Note:** All routes share the same layout function (set via `NewRouter(layout)`). Route groups provide view nesting and middleware inheritance, not separate layouts.

> **💡 Example:** See the [nested example](examples/nested) for a complete implementation of nested routing and view composition.

### Layout Functions

Layouts wrap your LiveView content and typically include HTML structure, CSS, and JavaScript:

```go
func RootLayout(children ...rend.Node) rend.Node {
    return html.Html(
        html.Head(
            html.Title(html.Text("My App")),
            html.Script(html.SrcAttr("https://unpkg.com/phoenix_live_view@1.0.17")),
        ),
        html.Body(
            html.Nav(
                // Navigation menu
            ),
            html.Main(children...),
            html.Script(html.Text(`
                const lv = new LiveView.LiveSocket("/live", Phoenix.Socket);
                lv.connect();
            `)),
        ),
    )
}
```

### Navigation

Use the `phx` package helpers for client-side navigation between LiveViews:

```go
import "github.com/go-live-view/go-live-view/phx"

// Patch navigation - updates URL without full page reload (within same session)
phx.LinkPatch("/users",
    phx.LinkSlot(html.Text("View Users")),
    phx.LinkAttr(html.ClassAttr("nav-link")),
)

// Navigate - full LiveView navigation (can cross sessions)
phx.LinkNavigate("/admin/dashboard",
    phx.LinkSlot(html.Text("Admin Dashboard")),
    phx.LinkReplace(true), // Replace in browser history instead of push
)

// Traditional HTTP link - full page reload
phx.LinkHref("/logout",
    phx.LinkMethod("post"),
    phx.LinkSlot(html.Text("Logout")),
)
```

**Navigation Types:**

| Type | Speed | Connection | Use Case |
|------|-------|------------|----------|
| `LinkPatch` | Fastest | Maintains WebSocket | Same LiveView updates |
| `LinkNavigate` | Fast | May break WebSocket | Cross-LiveView navigation |
| `LinkHref` | Slowest | Full page reload | External links, logout |

## Examples

The repository includes comprehensive examples demonstrating various LiveView features. Run the examples with:

```bash
cd examples
go run main.go
# Visit http://localhost:8080
```

| Example | Description | Key Features |
|---------|-------------|--------------|
| **[Counter](examples/counter)** | Basic increment/decrement counter | Events, state management, simple interactions |
| **[Uploads](examples/uploads)** | File upload with validation | File uploads, form validation, `phx.FileInput` |
| **[Streams](examples/stream)** | Dynamic list with add/delete | Streams, real-time updates, efficient list rendering |
| **[Async](examples/async)** | Async data loading with states | Async operations, loading states, error handling |
| **[Charts](examples/charts)** | Interactive ApexCharts integration | JavaScript hooks, chart updates, data visualization |
| **[JavaScript](examples/js)** | Modal with JS commands | Client-side commands, animations, DOM manipulation |
| **[Flash](examples/flash)** | Flash messages | User feedback, temporary messages |
| **[Broadcast](examples/broadcast)** | Multi-user real-time updates | PubSub, broadcasting, multi-client sync |
| **[Comprehension](examples/comprehension)** | List rendering optimization | Comprehensions, performance optimization |
| **[Nested](examples/nested)** | Nested routing example | Route groups, nested views, navigation |
| **[Server-Side Navigation](examples/ssnav)** | Server-controlled navigation | Navigation, route handling, URL updates |
| **[Scroll](examples/scroll)** | Infinite scroll implementation | Pagination, scroll events, dynamic loading |
| **[Components](examples/comp)** | Reusable component patterns | Component composition, layouts, reusability |

Each example includes complete source code and demonstrates best practices for that particular feature.

