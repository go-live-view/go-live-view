package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// URLs for MDN documentation
const (
	baseUrl    = "https://developer.mozilla.org/en-US/docs/Web"
	elUrl      = baseUrl + "/HTML/Element"
	attrUrl    = baseUrl + "/HTML/Global_attributes"
	ariaUrl    = baseUrl + "/Accessibility/ARIA/Attributes"
	svgElUrl   = baseUrl + "/SVG/Element"
	svgAttrUrl = baseUrl + "/SVG/Attribute"
)

// Cache configuration
const (
	maxWorkers = 100
)

// Attributes to skip
var skipAttrs = map[string]bool{
	"data-*": true,
	"data":   true,
}

// Data structures
type MDNData struct {
	Elements   []Element
	Attributes map[string]AttributeInfo
	Void       map[string]bool
	Deprecated map[string]bool
}

type AttributeInfo struct {
	Name     string
	Elements []string
}

type Element struct {
	Tag      string
	FuncName string
	Void     bool
}

// HTTP client helper
func getPage(client *http.Client, url string) (*goquery.Document, error) {
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch %s: status code %d", url, resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to create document from %s: %w", url, err)
	}

	return doc, nil
}

// Element extractor
type ElementExtractor struct {
	client *http.Client
}

func NewElementExtractor(client *http.Client) *ElementExtractor {
	return &ElementExtractor{client: client}
}

func (ee *ElementExtractor) getElements() ([]Element, map[string]bool, map[string]bool) {
	doc, err := getPage(ee.client, elUrl)
	if err != nil {
		log.Fatalf("Failed to fetch elements: %v", err)
	}

	var elements []Element
	found := make(map[string]bool)
	void := make(map[string]bool)
	deprecated := make(map[string]bool)

	// Try multiple selectors for better coverage
	selectors := []string{
		"tr td:first-child a code",
		"table tbody tr td:first-child a code",
		".standard-table tbody tr td:first-child a code",
	}

	for _, selector := range selectors {
		doc.Find(selector).Each(func(i int, s *goquery.Selection) {
			name := strings.TrimSpace(cleanName(s.Text()))
			if name == "" || found[name] {
				return
			}

			found[name] = true
			link := fmt.Sprintf("%s/%s", elUrl, name)

			// Get description from the same row
			var docs string
			if row := s.ParentsUntil("tr").Parent(); row.Length() > 0 {
				docs = strings.TrimSpace(row.Find("td").Eq(1).Text())
			}
			if docs == "" {
				docs = fmt.Sprintf("HTML element: %s", name)
			}

			// Check if deprecated by looking for deprecation indicators
			if strings.Contains(strings.ToLower(docs), "deprecated") ||
				strings.Contains(strings.ToLower(docs), "obsolete") {
				deprecated[name] = true
				return // Skip deprecated elements
			}

			// Check if void by looking for void indicators or visiting the element page
			isVoid := ee.checkIfVoid(link)
			if isVoid {
				void[name] = true
			}

			elements = append(elements, Element{
				Tag:      name,
				FuncName: toCamelCase(name),
				Void:     isVoid,
			})
		})
	}

	return elements, void, deprecated
}

func (ee *ElementExtractor) checkIfVoid(elementLink string) bool {
	// Fetch the element page to check if it's void
	doc, err := getPage(ee.client, elementLink)
	if err != nil {
		return false
	}

	// Look for void indicators in the page
	content := doc.Text()
	return strings.Contains(strings.ToLower(content), "void element") ||
		strings.Contains(strings.ToLower(content), "self-closing") ||
		strings.Contains(strings.ToLower(content), "no closing tag")
}

func (ee *ElementExtractor) getElementSpecificAttributes(elements []Element) map[string]AttributeInfo {
	attrs := make(map[string]AttributeInfo)
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Create a worker pool with more workers for better performance
	semaphore := make(chan struct{}, maxWorkers*2)

	// Fetch attributes for all elements, not just a hardcoded subset
	for _, elem := range elements {
		wg.Add(1)
		go func(element Element) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			url := fmt.Sprintf("%s/%s", elUrl, element.Tag)
			if element.Tag == "svg" {
				url = fmt.Sprintf("%s/%s", svgElUrl, element.Tag)
			}

			doc, err := getPage(ee.client, url)
			if err != nil {
				log.Printf("Warning: failed to fetch attributes for %s: %v", element.Tag, err)
				return
			}

			elementAttrs := ee.extractElementAttributes(doc, element.Tag)

			mu.Lock()
			for attrName, attrInfo := range elementAttrs {
				if existing, exists := attrs[attrName]; exists {
					// Combine with existing
					existing.Elements = append(existing.Elements, attrInfo.Elements...)
					attrs[attrName] = existing
				} else {
					attrs[attrName] = attrInfo
				}
			}
			mu.Unlock()
		}(elem)
	}

	wg.Wait()
	return attrs
}

func (ee *ElementExtractor) extractElementAttributes(doc *goquery.Document, elementName string) map[string]AttributeInfo {
	attrs := make(map[string]AttributeInfo)

	// Try multiple selectors for attribute tables and lists
	selectors := []string{
		"table thead th:contains('Attribute')",
		"table thead th:contains('Attributes')",
		"h2:contains('Attributes') + table",
		"h3:contains('Attributes') + table",
		"h2:contains('Attributes') + ul li code",
		"h3:contains('Attributes') + ul li code",
		"section:contains('Attributes') table",
		"section:contains('Attributes') ul li code",
		"div:contains('Attributes') table",
		"div:contains('Attributes') ul li code",
		"dl dt code", // Definition lists
		"ul li code", // Unordered lists
		"ol li code", // Ordered lists
	}

	for _, selector := range selectors {
		doc.Find(selector).Each(func(i int, table *goquery.Selection) {
			// Handle tables
			if table.Is("table") {
				table.Find("tbody tr").Each(func(j int, row *goquery.Selection) {
					cells := row.Find("td")
					if cells.Length() >= 2 {
						rawAttrName := strings.TrimSpace(cells.Eq(0).Text())
						attrName := strings.TrimSpace(cleanName(rawAttrName))

						if attrName == "" || skipAttrs[attrName] || strings.Contains(attrName, "deprecated") {
							return
						}

						attrs[attrName] = AttributeInfo{
							Name:     attrName,
							Elements: []string{elementName},
						}
					}
				})
			} else {
				// Handle lists and other elements
				table.Each(func(j int, item *goquery.Selection) {
					attrName := strings.TrimSpace(cleanName(item.Text()))
					if attrName == "" || skipAttrs[attrName] || strings.Contains(attrName, "deprecated") {
						return
					}

					attrs[attrName] = AttributeInfo{
						Name:     attrName,
						Elements: []string{elementName},
					}
				})
			}
		})
	}

	return attrs
}

// Attribute extractors
type AttributeExtractor struct {
	client *http.Client
}

func NewAttributeExtractor(client *http.Client) *AttributeExtractor {
	return &AttributeExtractor{client: client}
}

func (ae *AttributeExtractor) getGlobalAttributes() map[string]AttributeInfo {
	doc, err := getPage(ae.client, attrUrl)
	if err != nil {
		log.Printf("Warning: failed to fetch global attributes: %v", err)
		return make(map[string]AttributeInfo)
	}

	attrs := make(map[string]AttributeInfo)

	// Try multiple selectors for better coverage
	selectors := []string{
		"dl dt",
		"table tbody tr td:first-child code",
		"ul li code",
		".standard-table tbody tr td:first-child code",
		"h2:contains('Global attributes') + dl dt",
		"h3:contains('Global attributes') + dl dt",
		"section:contains('Global attributes') dl dt",
		"div:contains('Global attributes') dl dt",
	}

	for _, selector := range selectors {
		doc.Find(selector).Each(func(i int, s *goquery.Selection) {
			name := strings.TrimSpace(cleanTag(s.Text()))
			if name == "" || skipAttrs[name] {
				return
			}

			attrs[name] = AttributeInfo{
				Name:     name,
				Elements: []string{"*"},
			}
		})
	}

	return attrs
}

func (ae *AttributeExtractor) getAriaAttributes() map[string]AttributeInfo {
	doc, err := getPage(ae.client, ariaUrl)
	if err != nil {
		log.Printf("Warning: failed to fetch ARIA attributes: %v", err)
		return make(map[string]AttributeInfo)
	}

	attrs := make(map[string]AttributeInfo)

	// Try multiple selectors for ARIA attributes
	selectors := []string{
		"ol ul li a code",
		"table tbody tr td:first-child code",
		"ul li code",
		".standard-table tbody tr td:first-child code",
		"h2:contains('ARIA attributes') + ul li code",
		"h3:contains('ARIA attributes') + ul li code",
		"section:contains('ARIA attributes') ul li code",
		"div:contains('ARIA attributes') ul li code",
	}

	for _, selector := range selectors {
		doc.Find(selector).Each(func(i int, s *goquery.Selection) {
			name := strings.TrimSpace(s.Text())
			if name == "" || !strings.HasPrefix(name, "aria-") {
				return
			}

			attrs[name] = AttributeInfo{
				Name:     name,
				Elements: []string{"*"},
			}
		})
	}

	return attrs
}

func (ae *AttributeExtractor) getSVGAttributes() map[string]AttributeInfo {
	doc, err := getPage(ae.client, svgAttrUrl)
	if err != nil {
		log.Printf("Warning: failed to fetch SVG attributes: %v", err)
		return make(map[string]AttributeInfo)
	}

	attrs := make(map[string]AttributeInfo)

	// Extract from tables
	doc.Find("table tbody tr").Each(func(i int, s *goquery.Selection) {
		cells := s.Find("td")
		if cells.Length() < 2 {
			return
		}

		name := strings.TrimSpace(cells.Eq(0).Find("code").Text())
		if name == "" {
			return
		}

		attrs[name] = AttributeInfo{
			Name:     name,
			Elements: []string{"svg"},
		}
	})

	return attrs
}

// Code generator
type CodeGenerator struct{}

func NewCodeGenerator() *CodeGenerator {
	return &CodeGenerator{}
}

func (cg *CodeGenerator) generateElements(data MDNData, buf *strings.Builder) {
	fmt.Fprintf(buf, "// Code generated by go generate; DO NOT EDIT.\n")
	fmt.Fprintf(buf, "package html\n\n")
	fmt.Fprintf(buf, "import \"github.com/go-live-view/go-live-view/rend\"\n\n")

	// Generate elements
	fmt.Fprintf(buf, "// Elements\n\n")
	for _, elem := range data.Elements {
		if elem.Void {
			fmt.Fprintf(buf, "func %s(children ...rend.Node) *VoidNode {\n", elem.FuncName)
			fmt.Fprintf(buf, "\treturn Void(\"%s\", children...)\n", elem.Tag)
		} else {
			fmt.Fprintf(buf, "func %s(children ...rend.Node) *ElementNode {\n", elem.FuncName)
			fmt.Fprintf(buf, "\treturn Element(\"%s\", children...)\n", elem.Tag)
		}
		fmt.Fprintf(buf, "}\n\n")
	}

	// Generate attributes
	fmt.Fprintf(buf, "// Attributes\n\n")
	skippedAttrs := cg.generateAttributes(data.Attributes, buf)

	// Generate data attribute helper
	fmt.Fprintf(buf, "// Data Attributes\n\n")
	fmt.Fprintf(buf, "/*\nData attribute helper for creating custom data-* attributes\n*/\n")
	fmt.Fprintf(buf, "func DataAttr(name string, value ...string) *AttributeNode {\n")
	fmt.Fprintf(buf, "\treturn Attr(\"data-\"+name, value...)\n")
	fmt.Fprintf(buf, "}\n")

	// Log skipped attributes
	if len(skippedAttrs) > 0 {
		fmt.Printf("Skipped %d invalid attributes\n", len(skippedAttrs))
	}
}

func (cg *CodeGenerator) generateAttributes(attrs map[string]AttributeInfo, buf *strings.Builder) []string {
	// Sort attribute names for consistent output
	names := make([]string, 0, len(attrs))
	for name := range attrs {
		names = append(names, name)
	}
	sort.Strings(names)

	var skippedAttrs []string
	generatedFuncs := make(map[string]bool) // Track generated function names

	for _, name := range names {
		cleaned := cleanName(name)
		funcName := toCamelCase(cleaned)

		// Skip if cleaned name or original name is empty, starts with a digit, or contains invalid characters.
		if cleaned == "" || funcName == "" ||
			(len(cleaned) > 0 && cleaned[0] >= '0' && cleaned[0] <= '9') ||
			strings.ContainsAny(cleaned, " \n\r\t:") || strings.ContainsAny(name, "\n\r\t:") {
			skippedAttrs = append(skippedAttrs, name)
			continue
		}

		// Always suffix with 'Attr'
		funcName += "Attr"

		// Skip if we've already generated this function name
		if generatedFuncs[funcName] {
			continue
		}

		generatedFuncs[funcName] = true
		fmt.Fprintf(buf, "func %s(value ...string) *AttributeNode {\n", funcName)
		fmt.Fprintf(buf, "\treturn Attr(\"%s\", value...)\n", name)
		fmt.Fprintf(buf, "}\n\n")
	}

	return skippedAttrs
}

// Utility functions
func toCamelCase(s string) string {
	// Clean the string first to remove invalid characters
	s = cleanName(s)

	// Handle special cases
	if s == "" {
		return "Attr"
	}

	// Handle numeric-only names
	if isNumeric(s) {
		return "Attr" + s
	}

	// Remove leading digits (Go identifiers can't start with a digit)
	for len(s) > 0 && s[0] >= '0' && s[0] <= '9' {
		s = s[1:]
	}
	if s == "" {
		return "Attr"
	}

	words := strings.Fields(cases.Title(language.English).String(s))
	return strings.Join(words, "")
}

var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9]+`)

func cleanName(s string) string {
	s = nonAlphanumericRegex.ReplaceAllString(s, " ")
	return strings.Join(strings.Fields(s), " ")
}

func cleanTag(s string) string {
	replacer := strings.NewReplacer("\"", " ")
	return replacer.Replace(s)
}

func isNumeric(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return len(s) > 0
}

// Main function
func main() {
	// Initialize components
	client := &http.Client{Timeout: 30 * time.Second}
	attrExtractor := NewAttributeExtractor(client)
	elementExtractor := NewElementExtractor(client)
	codeGenerator := NewCodeGenerator()

	// Extract data concurrently
	var wg sync.WaitGroup
	var globalAttrs, ariaAttrs, svgAttrs map[string]AttributeInfo
	var elements []Element
	var elementAttrs map[string]AttributeInfo
	var void, deprecated map[string]bool

	wg.Add(4)

	go func() {
		defer wg.Done()
		globalAttrs = attrExtractor.getGlobalAttributes()
	}()

	go func() {
		defer wg.Done()
		ariaAttrs = attrExtractor.getAriaAttributes()
	}()

	go func() {
		defer wg.Done()
		svgAttrs = attrExtractor.getSVGAttributes()
	}()

	go func() {
		defer wg.Done()
		elements, void, deprecated = elementExtractor.getElements()
		elementAttrs = elementExtractor.getElementSpecificAttributes(elements)
	}()

	wg.Wait()

	// Combine all attributes
	allAttrs := make(map[string]AttributeInfo)
	for name, attr := range globalAttrs {
		allAttrs[name] = attr
	}
	for name, attr := range ariaAttrs {
		allAttrs[name] = attr
	}
	for name, attr := range svgAttrs {
		allAttrs[name] = attr
	}
	for name, attr := range elementAttrs {
		if existing, exists := allAttrs[name]; exists {
			// Combine with existing
			existing.Elements = append(existing.Elements, attr.Elements...)
			allAttrs[name] = existing
		} else {
			allAttrs[name] = attr
		}
	}

	data := MDNData{
		Elements:   elements,
		Attributes: allAttrs,
		Void:       void,
		Deprecated: deprecated,
	}

	// Generate code
	var buf strings.Builder
	codeGenerator.generateElements(data, &buf)

	// Write to current directory (since go generate runs from html directory)
	outputPath := "elements.go"
	err := os.WriteFile(outputPath, []byte(buf.String()), 0644)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Generated elements.go: %d elements, %d attributes", len(elements), len(allAttrs))
	if len(deprecated) > 0 {
		fmt.Printf(", %d deprecated elements skipped", len(deprecated))
	}
	fmt.Println()
}
