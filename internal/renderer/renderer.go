package renderer

import (
	"bytes"
	"fmt"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"io"
	"regexp"
	"strings"
)

// CustomHTMLRenderer is a custom html.Renderer for HTML output.
type CustomHTMLRenderer struct {
	r *html.Renderer
}

func (c *CustomHTMLRenderer) RenderHeader(w io.Writer, ast ast.Node) {
	c.r.RenderHeader(w, ast)
}

func (c *CustomHTMLRenderer) RenderFooter(w io.Writer, ast ast.Node) {
	c.r.RenderFooter(w, ast)
}

func (c *CustomHTMLRenderer) RenderNode(w io.Writer, node ast.Node, entering bool) ast.WalkStatus {
	return c.r.RenderNode(w, node, entering)
}

// writeTOC generates a Table of Contents from the headings in the document
func (c *CustomHTMLRenderer) writeTOC(doc ast.Node) string {
	buf := bytes.Buffer{}
	tocLevel := 1
	headingCount := 0

	ast.WalkFunc(doc, func(node ast.Node, entering bool) ast.WalkStatus {
		if nodeData, ok := node.(*ast.Heading); ok && !nodeData.IsTitleblock {
			// Skip H1 headings for the TOC
			if nodeData.Level == 1 {
				return ast.GoToNext
			}
			if !entering {
				buf.WriteString("</a>")
				return ast.GoToNext
			}
			if nodeData.HeadingID == "" {
				// Extract heading text and generate slug
				nodeData.HeadingID = slugify(extractText(nodeData))
			}
			// Adjust TOC nesting level
			if nodeData.Level == tocLevel {
				buf.WriteString("</li>\n\n<li>")
			} else if nodeData.Level < tocLevel {
				for nodeData.Level < tocLevel {
					tocLevel--
					buf.WriteString("</li>\n</ul>")
				}
				buf.WriteString("</li>\n\n<li>")
			} else {
				for nodeData.Level > tocLevel {
					tocLevel++
					buf.WriteString("\n<ul>\n<li>")
				}
			}
			// Add the heading link to the TOC
			fmt.Fprintf(&buf, `<a href="#%s">`, nodeData.HeadingID)
			buf.WriteString(extractText(nodeData))
			headingCount++
			return ast.GoToNext
		}

		return ast.GoToNext
	})

	// Close any open tags based on the current TOC level
	for ; tocLevel > 1; tocLevel-- {
		buf.WriteString("</li>\n</ul>")
	}

	// If the TOC has content, wrap it in a <nav> tag and return as a string
	if buf.Len() > 0 {
		return fmt.Sprintf("<nav>\n%s\n\n</nav>\n", buf.String())
	}
	return ""
}

// Helper function to extract text from a heading node
func extractText(node *ast.Heading) string {
	var buf bytes.Buffer
	for _, child := range node.Children {
		if text, ok := child.(*ast.Text); ok {
			buf.Write(text.Literal)
		}
	}
	return buf.String()
}

// Helper function to generate a slug from a string
func slugify(s string) (slug string) {
	slug = regexp.MustCompile(`[^a-z0-9\-]+`).ReplaceAllString(strings.ToLower(s), "-")
	slug = strings.Trim(slug, "-")
	slug = regexp.MustCompile(`-+`).ReplaceAllString(slug, "-")
	return
}
