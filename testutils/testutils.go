package testutils

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/go-live-view/go-live-view/rend"
	"github.com/stretchr/testify/assert"
)

var (
	update = flag.Bool("update", false, "update .json files")
)

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
}

// GoldenFile handles reading and writing golden files for tests
type GoldenFile struct {
	path string
}

// NewGoldenFile creates a new golden file handler
func NewGoldenFile(path string) *GoldenFile {
	return &GoldenFile{path: path}
}

// Read reads the golden file content
func (g *GoldenFile) Read(t *testing.T) string {
	t.Helper()

	content, err := os.ReadFile(g.path)
	if err != nil {
		t.Fatalf("Error opening file %s: %s", g.path, err)
	}

	return string(content)
}

// Write writes content to the golden file
func (g *GoldenFile) Write(t *testing.T, content string) {
	t.Helper()

	// Ensure the directory exists
	dir := filepath.Dir(g.path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("Error creating directory %s: %s", dir, err)
	}

	err := os.WriteFile(g.path, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Error writing to file %s: %s", g.path, err)
	}
}

// Update reads or writes the golden file based on the update flag
func (g *GoldenFile) Update(t *testing.T, actual string) string {
	t.Helper()

	if *update {
		g.Write(t, actual)
		return actual
	}

	return g.Read(t)
}

// Stringify converts a name to a filename-safe string
func Stringify(name string) string {
	return strings.ReplaceAll(name, " ", "-")
}

// TestCase represents a single test case with node and expected HTML
type TestCase struct {
	Name     string
	Node     rend.Node
	Expected string
}

// RunTestCase runs a single test case with golden file validation
func RunTestCase(t *testing.T, tc TestCase, jsonPath string) {
	t.Run(tc.Name, func(t *testing.T) {
		assert.NotNil(t, tc.Node)

		jsonOut := rend.RenderJSON(tc.Node)
		htmlOut := rend.RenderString(tc.Node)

		goldenFile := NewGoldenFile(jsonPath)
		actualJSON := goldenFile.Update(t, jsonOut)

		assert.JSONEq(t, actualJSON, jsonOut)
		assert.Equal(t, tc.Expected, htmlOut)
	})
}

// RunTestCases runs multiple test cases with golden file validation
func RunTestCases(t *testing.T, tt []TestCase, prefix string) {
	for _, tc := range tt {
		jsonPath := "testdata/" + prefix + "/" + Stringify(tc.Name) + ".json"
		RunTestCase(t, tc, jsonPath)
	}
}

func TestNode(node rend.Node) rend.Node {
	if node == nil {
		return TestString("") // Return empty string node instead of nil
	}
	return node
}

func TestString(s string) rend.Node {
	return &testStringNode{s}
}

type testStringNode struct {
	content string
}

func (n *testStringNode) Render(diff bool, root *rend.Root, t *rend.Rend, b *strings.Builder) error {
	_, err := b.WriteString(n.content)
	return err
}
