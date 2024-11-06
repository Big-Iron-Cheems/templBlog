package main

import (
	"embed"
	"fmt"
	"github.com/a-h/templ"
	"github.com/rs/zerolog/log"
	"net/http"
	. "templBlog/internal/config"
	"templBlog/internal/renderer"
	"templBlog/internal/templ/views"
)

//go:embed internal/static
var static embed.FS

// setupRoutes sets up the routes for the server, returning the pages handler.
func setupRoutes() (pagesHandler *http.ServeMux) {
	pagesHandler = http.NewServeMux()

	// Serve blog posts at /posts
	renderer.ServeBlogPosts(pagesHandler, GlobalConfig.PostsDirectory)

	// Serve static files at /internal/static
	pagesHandler.Handle("/internal/static/", http.FileServer(http.FS(static)))

	// Catch-all handler for unmatched routes
	pagesHandler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			// Serve the homepage
			templ.Handler(views.Index(renderer.GetPostMetadataList())).ServeHTTP(w, r)
		} else {
			// No match found, show 404 page
			w.WriteHeader(http.StatusNotFound)
			templ.Handler(views.NotFound()).ServeHTTP(w, r)
		}
	})

	// Serve the other pages
	pagesHandler.Handle("/about", templ.Handler(views.About()))
	pagesHandler.Handle("/contacts", templ.Handler(views.Contacts()))

	return
}

func main() {
	// Init
	InitConfig()
	InitLogger()
	defer CloseLogger()

	// Set up routes
	pagesHandler := setupRoutes()

	// Start the server
	log.Debug().Msgf("Starting server on port %d", GlobalConfig.Port)
	log.Fatal().Err(http.ListenAndServe(fmt.Sprintf(":%d", GlobalConfig.Port), pagesHandler)).Msg("Failed to start server")
}
