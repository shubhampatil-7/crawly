package main

import (
	"net/url"
	"reflect"
	"testing"
)

func TestGetHeadingFromHTML(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		expected string
	}{
		{
			name:     "returns h1 text content",
			html:     "<html><body><h1>Main Title</h1></body></html>",
			expected: "Main Title",
		},
		{
			name:     "returns h2 when no h1 present",
			html:     "<html><body><h2>Sub Title</h2></body></html>",
			expected: "Sub Title",
		}, {
			name:     "prefers h1 over h2 when both present",
			html:     "<html><body><h1>Main Title</h1><h2>Sub Title</h2></body></html>",
			expected: "Main Title",
		},
		{
			name:     "prefers h1 even when h2 appears first",
			html:     "<html><body><h2>Sub Title</h2><h1>Main Title</h1></body></html>",
			expected: "Main Title",
		},
		{
			name:     "returns empty string when neither h1 nor h2 present",
			html:     "<html><body><p>Some paragraph</p></body></html>",
			expected: "",
		},
		{
			name:     "returns empty string for empty HTML",
			html:     "",
			expected: "",
		},
		{
			name:     "strips inner HTML tags from h1",
			html:     "<html><body><h1><span>Styled</span> Title</h1></body></html>",
			expected: "Styled Title",
		},
		{
			name:     "strips inner HTML tags from h2",
			html:     "<html><body><h2><em>Italic</em> Heading</h2></body></html>",
			expected: "Italic Heading",
		},
		{
			name:     "handles h1 with extra whitespace",
			html:     "<html><body><h1>  Padded Title  </h1></body></html>",
			expected: "Padded Title",
		},
		{
			name:     "handles h2 with extra whitespace",
			html:     "<html><body><h2>  Padded Sub Title  </h2></body></html>",
			expected: "Padded Sub Title",
		},
		{
			name:     "returns empty string for empty h1 tag",
			html:     "<html><body><h1></h1></body></html>",
			expected: "",
		},
		{
			name:     "returns empty string for empty h2 tag with no h1",
			html:     "<html><body><h2></h2></body></html>",
			expected: "",
		},
	}

	for _, value := range tests {
		t.Run(value.name, func(t *testing.T) {
			result := getHeadingFromHTML(value.html)
			if result != value.expected {
				t.Errorf("getHeadingFromHTML(%q)=%q, expected: %q", value.html, result, value.expected)
			}
		})
	}

}

func TestGetHeadingFromHTMLBasic(t *testing.T) {
	inputBody := "<html><body><h1>Test Title</h1></body></html>"
	actual := getHeadingFromHTML(inputBody)
	expected := "Test Title"

	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}

func TestGetFirstParagraphFromHTML(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		expected string
	}{
		{
			name: "prefers paragraph inside main",
			html: `<html><body>
				<p>Outside paragraph.</p>
				<main>
					<p>Main paragraph.</p>
				</main>
			</body></html>`,
			expected: "Main paragraph.",
		},
		{
			name: "falls back to body paragraph when no main",
			html: `<html><body>
				<p>Only paragraph.</p>
			</body></html>`,
			expected: "Only paragraph.",
		},
		{
			name: "falls back to body paragraph when main has no paragraphs",
			html: `<html><body>
				<p>Outside paragraph.</p>
				<main>
				</main>
			</body></html>`,
			expected: "Outside paragraph.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := getFirstParagraphFromHTML(tt.html)
			if actual != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, actual)
			}
		})
	}
}

func TestGetURLsFromHTML(t *testing.T) {
	tests := []struct {
		name        string
		inputURL    string
		inputBody   string
		expected    []string
		expectError bool
	}{
		{
			name:      "absolute URL",
			inputURL:  "https://crawler-test.com",
			inputBody: `<html><body><a href="https://crawler-test.com"><span>Boot.dev</span></a></body></html>`,
			expected:  []string{"https://crawler-test.com"},
		},
		{
			name:      "relative URL converted to absolute",
			inputURL:  "https://crawler-test.com",
			inputBody: `<html><body><a href="/about">About</a></body></html>`,
			expected:  []string{"https://crawler-test.com/about"},
		},
		{
			name:     "multiple URLs",
			inputURL: "https://crawler-test.com",
			inputBody: `<html><body>
				<a href="https://crawler-test.com">Home</a>
				<a href="/about">About</a>
				<a href="/contact">Contact</a>
			</body></html>`,
			expected: []string{
				"https://crawler-test.com",
				"https://crawler-test.com/about",
				"https://crawler-test.com/contact",
			},
		},
		{
			name:      "no anchor tags returns empty slice",
			inputURL:  "https://crawler-test.com",
			inputBody: `<html><body><p>No links here</p></body></html>`,
			expected:  []string{},
		},
		{
			name:      "anchor tag with no href is skipped",
			inputURL:  "https://crawler-test.com",
			inputBody: `<html><body><a>No href</a></body></html>`,
			expected:  []string{},
		},
		{
			name:      "relative URL with nested path",
			inputURL:  "https://crawler-test.com/blog",
			inputBody: `<html><body><a href="/blog/post-1">Post 1</a></body></html>`,
			expected:  []string{"https://crawler-test.com/blog/post-1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baseURL, err := url.Parse(tt.inputURL)
			if err != nil {
				t.Errorf("couldn't parse input URL: %v", err)
				return
			}

			actual, err := getURLsFromHTML(tt.inputBody, baseURL)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected an error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, actual)
			}
		})
	}
}

func TestGetImagesFromHTML(t *testing.T) {
	tests := []struct {
		name        string
		inputURL    string
		inputBody   string
		expected    []string
		expectError bool
	}{
		{
			name:      "relative URL converted to absolute",
			inputURL:  "https://crawler-test.com",
			inputBody: `<html><body><img src="/logo.png" alt="Logo"></body></html>`,
			expected:  []string{"https://crawler-test.com/logo.png"},
		},
		{
			name:      "absolute image URL returned as-is",
			inputURL:  "https://crawler-test.com",
			inputBody: `<html><body><img src="https://cdn.example.com/image.jpg" alt="Image"></body></html>`,
			expected:  []string{"https://cdn.example.com/image.jpg"},
		},
		{
			name:     "multiple images returned in order",
			inputURL: "https://crawler-test.com",
			inputBody: `<html><body>
				<img src="/img/one.png" alt="One">
				<img src="/img/two.png" alt="Two">
				<img src="https://cdn.example.com/three.png" alt="Three">
			</body></html>`,
			expected: []string{
				"https://crawler-test.com/img/one.png",
				"https://crawler-test.com/img/two.png",
				"https://cdn.example.com/three.png",
			},
		},
		{
			name:      "img tag with no src is skipped",
			inputURL:  "https://crawler-test.com",
			inputBody: `<html><body><img alt="No src"></body></html>`,
			expected:  []string{},
		},
		{
			name:      "no img tags returns empty slice",
			inputURL:  "https://crawler-test.com",
			inputBody: `<html><body><p>No images here</p></body></html>`,
			expected:  []string{},
		},
		{
			name:      "relative URL with nested base path",
			inputURL:  "https://crawler-test.com/blog",
			inputBody: `<html><body><img src="/assets/banner.png" alt="Banner"></body></html>`,
			expected:  []string{"https://crawler-test.com/assets/banner.png"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baseURL, err := url.Parse(tt.inputURL)
			if err != nil {
				t.Errorf("couldn't parse input URL: %v", err)
				return
			}

			actual, err := getImagesFromHTML(tt.inputBody, baseURL)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected an error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, actual)
			}
		})
	}
}
func TestExtractPageData(t *testing.T) {
	inputURL := "https://crawler-test.com"
	inputBody := `<html><body>
        <h1>Test Title</h1>
        <p>This is the first paragraph.</p>
        <a href="/link1">Link 1</a>
        <img src="/image1.jpg" alt="Image 1">
    </body></html>`

	actual := extractPageData(inputBody, inputURL)

	expected := PageData{
		URL:            "https://crawler-test.com",
		Heading:        "Test Title",
		FirstParagraph: "This is the first paragraph.",
		OutgoingLinks:  []string{"https://crawler-test.com/link1"},
		ImageURLs:      []string{"https://crawler-test.com/image1.jpg"},
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %+v, got %+v", expected, actual)
	}
}
