// Package docs serves generated API documentation as HTML.
package docs

import (
	"bytes"
	"html/template"
	"io/fs"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

type docPage struct {
	Title   string
	Slug    string
	Content template.HTML
}

var (
	pages []docPage
	md    goldmark.Markdown
)

func init() {
	md = goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithRendererOptions(html.WithUnsafe()),
	)
}

func loadPages(docsFS fs.FS) {
	pages = nil
	entries, _ := fs.ReadDir(docsFS, ".")
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		data, _ := fs.ReadFile(docsFS, e.Name())
		var buf bytes.Buffer
		md.Convert(data, &buf)

		slug := strings.TrimSuffix(e.Name(), ".md")
		title := strings.ReplaceAll(slug, "-", " ")
		pages = append(pages, docPage{
			Title:   title,
			Slug:    slug,
			Content: template.HTML(buf.String()),
		})
	}
}

var indexTmpl = template.Must(template.New("index").Parse(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>GesiTr API Docs</title>
<style>
  body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; max-width: 900px; margin: 0 auto; padding: 2rem; line-height: 1.6; color: #1a1a1a; }
  a { color: #0366d6; text-decoration: none; }
  a:hover { text-decoration: underline; }
  ul { padding-left: 1.5rem; }
  li { margin: 0.5rem 0; }
  h1 { border-bottom: 1px solid #eee; padding-bottom: 0.5rem; }
</style>
</head>
<body>
<h1>GesiTr API Documentation</h1>
<ul>
{{range .}}<li><a href="/docs/{{.Slug}}">{{.Title}}</a></li>
{{end}}</ul>
</body>
</html>`))

var pageTmpl = template.Must(template.New("page").Parse(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>{{.Title}} — GesiTr API Docs</title>
<style>
  body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; max-width: 900px; margin: 0 auto; padding: 2rem; line-height: 1.6; color: #1a1a1a; }
  a { color: #0366d6; text-decoration: none; }
  a:hover { text-decoration: underline; }
  pre { background: #f6f8fa; padding: 1rem; border-radius: 6px; overflow-x: auto; }
  code { font-family: "SFMono-Regular", Consolas, monospace; font-size: 0.9em; }
  :not(pre) > code { background: #f0f0f0; padding: 0.2em 0.4em; border-radius: 3px; }
  details { margin: 0.5rem 0; }
  summary { cursor: pointer; font-weight: 600; }
  h1, h2, h3 { margin-top: 1.5rem; }
  h1 { border-bottom: 1px solid #eee; padding-bottom: 0.5rem; }
  h2 { border-bottom: 1px solid #f0f0f0; padding-bottom: 0.3rem; }
  table { border-collapse: collapse; width: 100%; }
  th, td { border: 1px solid #ddd; padding: 0.5rem; text-align: left; }
  th { background: #f6f8fa; }
  .nav { margin-bottom: 1.5rem; }
</style>
</head>
<body>
<div class="nav"><a href="/docs">← All docs</a></div>
{{.Content}}
</body>
</html>`))

// SetupRoutes registers the /docs routes on the given engine.
// docsFS should contain the generated markdown files at its root.
func SetupRoutes(r *gin.Engine, docsFS fs.FS) {
	loadPages(docsFS)

	r.GET("/docs", func(c *gin.Context) {
		var buf bytes.Buffer
		indexTmpl.Execute(&buf, pages)
		c.Data(http.StatusOK, "text/html; charset=utf-8", buf.Bytes())
	})

	for i := range pages {
		p := pages[i]
		r.GET("/docs/"+p.Slug, func(c *gin.Context) {
			var buf bytes.Buffer
			pageTmpl.Execute(&buf, p)
			c.Data(http.StatusOK, "text/html; charset=utf-8", buf.Bytes())
		})
	}

	// Redirect /docs/ with trailing slash to /docs
	r.GET("/docs/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/docs")
	})

}
