package renderer

import (
	"bytes"
	"fmt"
	"github.com/alecthomas/chroma/v2"
	chromaHtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/gomarkdown/markdown/ast"
	"github.com/rs/zerolog/log"
	"io"
)

var (
	htmlFormatter  *chromaHtml.Formatter
	highlightStyle *chroma.Style
)

func init() {
	htmlFormatter = chromaHtml.New(chromaHtml.WithClasses(true), chromaHtml.TabWidth(4))
	if htmlFormatter == nil {
		log.Panic().Msg("couldn't create html formatter")
	}
	styleName := "onedark"
	highlightStyle = styles.Get(styleName)
	if highlightStyle == nil {
		log.Panic().Msgf("didn't find style '%s'", styleName)
	}
}

// CustomRenderNodeHook is a custom hook for rendering nodes.
// Use this to modify the rendering of specific nodes.
func CustomRenderNodeHook(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
	if heading, ok := node.(*ast.Heading); ok && entering && heading.Level > 1 {
		renderHeadingLink(w, heading, entering)
	} else if codeBlock, ok := node.(*ast.CodeBlock); ok && entering {
		renderCodeColor(w, codeBlock, entering)
		return ast.GoToNext, true
	}

	return ast.GoToNext, false
}

// renderHeadingLink renders an anchor link before a heading.
func renderHeadingLink(w io.Writer, heading *ast.Heading, entering bool) {
	if !entering {
		return
	}

	// If we are a Header of level 2 or higher, add an anchor link before the header.
	var buf bytes.Buffer
	for _, child := range heading.GetChildren() {
		if text, ok := child.(*ast.Text); ok {
			buf.WriteString(string(text.Literal))
		}
	}
	headerText := buf.String()

	// Create a kebab-case anchor link.
	kebabCase := slugify(headerText)

	heading.Attribute = &ast.Attribute{
		ID: []byte(kebabCase),
	}

	anchorLink := fmt.Sprintf(`<a href="#%s">#</a>`, kebabCase)

	anchorNode := &ast.HTMLBlock{
		Leaf: ast.Leaf{
			Literal: []byte(anchorLink),
		},
	}

	// Prepend the anchor node as the first child of the heading.
	heading.SetChildren(append([]ast.Node{anchorNode}, heading.GetChildren()...))
}

// renderCodeColor renders a code block with syntax highlighting.
// The code syntax is determined by the language specified in the code block.
// If no language is specified, the code block is rendered as plain text.
// Chroma is used to provide syntax highlighting.
func renderCodeColor(w io.Writer, codeBlock *ast.CodeBlock, entering bool) {
	if !entering {
		return
	}

	// Get the language of the code block.
	lang := string(codeBlock.Info)
	code := string(codeBlock.Literal)

	l := lexers.Get(lang)
	if l == nil {
		l = lexers.Analyse(code)
	}
	if l == nil {
		l = lexers.Fallback
	}
	l = chroma.Coalesce(l)

	it, err := l.Tokenise(nil, code)
	if err != nil {
		log.Error().Err(err).Msg("failed to tokenize code block")
		return
	}

	if err = htmlFormatter.Format(w, highlightStyle, it); err != nil {
		log.Error().Err(err).Msg("failed to format code block")
		return
	}
}
