package renderer

import (
	"fmt"
	"github.com/a-h/templ"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"templBlog/internal/models"
	"templBlog/internal/templ/views"
)

// postMetadataList is used to generate the blog index.
var postMetadataList []models.PostMetadata

func GetPostMetadataList() []models.PostMetadata {
	return postMetadataList
}

// PostData is the data for a blog post.
type PostData struct {
	Metadata models.PostMetadata // Metadata for the post.
	Content  string              // Content of the post in HTML.
	TocHTML  string              // Table of contents for the post in HTML.
}

// ServeBlogPosts loads blog posts from a posts directory and serves them.
func ServeBlogPosts(mux *http.ServeMux, postsDir string) {
	renderer := &CustomHTMLRenderer{
		html.NewRenderer(html.RendererOptions{
			Flags:          html.CommonFlags,
			RenderNodeHook: CustomRenderNodeHook,
		}),
	}

	slugMap := make(map[string]bool)
	var unsortedPosts []PostData

	// Load the posts from the embedded file system.
	if err := filepath.Walk(postsDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(info.Name(), ".md") {
			return nil // Skip directories and non-Markdown files
		}

		log.Debug().Msgf("Processing post file: %s", path)

		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read post file: %w", err)
		}

		// Extract metadata
		parts := strings.SplitN(string(content), "---", 3)
		if len(parts) <= 2 {
			return fmt.Errorf("post %s does not have metadata", path)
		}

		var metadata models.PostMetadata
		if err := yaml.Unmarshal([]byte(parts[1]), &metadata); err != nil {
			return fmt.Errorf("failed to unmarshal metadata for post %s: %w", path, err)
		}

		// Process Markdown content to extract and remove the title
		markdownContent := parts[2]
		strippedContent, extractedTitle := removeMarkdownTitle(markdownContent)
		if metadata.Title == "" {
			metadata.Title = extractedTitle
		}

		// Generate slug from title
		if metadata.Slug == "" {
			metadata.Slug = slugify(metadata.Title)
		}

		// Check for duplicate slug
		if _, exists := slugMap[metadata.Slug]; exists {
			return fmt.Errorf("duplicate slug found: %s", metadata.Slug)
		}
		slugMap[metadata.Slug] = true

		// Normalize tags to lowercase
		if metadata.Tags != nil {
			for i, tag := range metadata.Tags {
				metadata.Tags[i] = slugify(tag)
			}
		}

		log.Debug().Msgf("Post metadata: %+v", metadata)

		// Parse and render Markdown content to HTML
		p := parser.NewWithExtensions(parser.CommonExtensions | parser.Footnotes)
		markdownDoc := markdown.Parse([]byte(strippedContent), p)
		htmlContent := markdown.Render(markdownDoc, renderer)

		// Generate the TOC
		tocHTML := renderer.writeTOC(markdownDoc)

		// Temporarily store this post's data for sorting
		unsortedPosts = append(unsortedPosts, PostData{
			metadata,
			string(htmlContent),
			tocHTML,
		})

		return nil
	}); err != nil {
		log.Error().Err(err).Msg("Failed to walk posts")
		return
	}
	log.Debug().Msgf("Fetched %d posts", len(unsortedPosts))

	// Sort posts by Date (newest first)
	sort.Slice(unsortedPosts, func(i, j int) bool {
		return unsortedPosts[i].Metadata.Date.After(unsortedPosts[j].Metadata.Date)
	})

	// Serve the sorted posts and update postMetadataList
	for i, post := range unsortedPosts {
		var prevSlug, nextSlug string
		if i > 0 {
			nextSlug = unsortedPosts[i-1].Metadata.Slug
		}
		if i < len(unsortedPosts)-1 {
			prevSlug = unsortedPosts[i+1].Metadata.Slug
		}

		mux.Handle(fmt.Sprintf("/posts/%s", post.Metadata.Slug), templ.Handler(
			views.BlogPost(post.Metadata.Title, post.Metadata.Author, post.Metadata.Date.Format("Jan 02, 2006"), post.Content, post.TocHTML, prevSlug, nextSlug, post.Metadata.Tags),
		))

		postMetadataList = append(postMetadataList, post.Metadata)
		log.Debug().Msgf("Served post `%s` at page `/posts/%s`", post.Metadata.Title, post.Metadata.Slug)
	}
}

// removeMarkdownTitle removes the first Markdown H1 header (# Title) from the content
// and returns the modified content and the title.
func removeMarkdownTitle(content string) (modifiedContent string, title string) {
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, "# ") {
			title = strings.TrimSpace(line[2:])
			modifiedContent = strings.Join(lines[i+1:], "\n")
			return
		}
	}
	// If no title is found, return the original content and empty title.
	modifiedContent = content
	return
}
