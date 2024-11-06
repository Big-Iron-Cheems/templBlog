# Example Blog Post website

This project showcases a minimalist blog website created entirely in Golang.

## Prerequisites

Ensure you have the following installed on your system:

- [Go](https://golang.org/dl/)
- [Node.js](https://nodejs.org/en/download/)
- [Air](https://github.com/air-verse/air) for live reloading
- [Chroma](https://github.com/alecthomas/chroma) for syntax highlighting

## Setup

1. Clone the repository.
2. Install Go dependencies via `go mod tidy`
3. Install Node.js dependencies via `npm install`
4. Run `npm watch` to generate the `tailwind.css` file and update it while working on the project.
5. Run `chroma --style=STYLE_NAME --html-styles > static/css/chroma.css` to generate the chroma css file.
6. Run `air` to start the live reloading service.

## Usage

Markdown files added to the `/posts` folder will be used to generate posts.  
These files are expected to begin with a metadata section as follows:

```yaml
---
author: Author name
date: YYYY-MM-DD
brief: A brief description of the article.
---
```

You may specify a custom `title`, `slug` and `tags`.  
`title` and `slug` are automatically generated based on the topmost header in the file.

## Project structure

Some important directories and files in the project are:

- `configs/config.json`: configuration file for the project.
- `internal/templ/`: reusable components and views.
- `posts/`: markdown files for blog posts.

## Used technologies

- [Templ](https://github.com/a-h/templ) for Golang based templating of HTML files.
- [Tailwind CSS](https://github.com/tailwindlabs/tailwindcss) for styling.
- [Chroma](https://github.com/alecthomas/chroma) for syntax highlighting in Markdown codeblocks.
- [Air](https://github.com/air-verse/air) for live reloading of the website.
- [Heroicons](https://heroicons.com/) for icons.
- [Simple Icons](https://simpleicons.org/) for popular brand logos.

## Contributing

Feel free to contribute to this project by creating a pull request.