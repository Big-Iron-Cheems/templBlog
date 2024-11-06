package models

import (
	"fmt"
	"time"
)

/*
PostMetadata is the metadata for a blog post.
It is expected to be at the top of a Markdown file.

Example metadata:

	---
	author: Author name
	date: YYYY-MM-DD
	brief: A brief description of the article.
	---
*/
type PostMetadata struct {
	Author string    `yaml:"author"` // Author of the post.
	Date   time.Time `yaml:"date"`   // Date format YYYY-MM-DD in MD file, Jan 01, 2000 in rendered HTML.
	Brief  string    `yaml:"brief"`  // Brief a short description of the post.
	Tags   []string  `yaml:"tags"`   // Tags are used to categorize the post.
	Title  string    `yaml:"title"`  // Title of the post, ideally the same as the main header.
	Slug   string    `yaml:"slug"`   // Slug is the URL path for the post, generated from the title.
}

func (p PostMetadata) String() string {
	return fmt.Sprintf("Author: `%s`, Date: `%s`, Title: `%s`, Brief: `%s`, Slug: `%s`",
		p.Author, p.Date.Format("2006-01-02"), p.Title, p.Brief, p.Slug)
}
