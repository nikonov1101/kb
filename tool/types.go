package tool

import (
	"bytes"
	"html/template"
	"time"

	"github.com/nikonov1101/kb/tool/templates"
)

const (
	Published = "published"
	Private   = "private"
	Hidden    = "none"

	dateFormat = "02 Jan 2006"
)

// Source describes the Note source in all possible states:
// a path on disk, metadata parsed from markdown,
// the markdown itself, and generated HTML.
type Source struct {
	// Path on disk, relative to source_dir
	Path string
	// BaseName is a file name without source dir and .md extension
	BaseName string
	// sequential file number
	Num int64
	// Title of the note, from headers
	Title string
	// Date of publication, described by dateFormat, from headers
	Date time.Time
	// Visibility of rendered page:
	// "published" - rendered, listed on index page
	// "private" - rendered, not listed on index page
	// "none" - not rendered
	Visibility string

	// raw markdown content
	markdown []byte
	// cached generated html
	html []byte
}

func (s *Source) HTMLFileName() string {
	return s.BaseName + ".html"
}

func (s Source) URL(root string) string {
	return root + "/" + s.HTMLFileName()
}

func (s Source) Render(siteName string, posts []Source) ([]byte, error) {
	args := templates.Args{
		Title:     s.Title,
		OrContent: template.HTML(s.html),
	}

	if args.Title == "" {
		args.Title = siteName
	}

	if !s.Date.IsZero() {
		args.Date = s.Date.Format(dateFormat)
	}

	// index page
	if len(posts) > 0 {
		args.Title = siteName
		for _, post := range posts {
			args.Posts = append(args.Posts, templates.Post{
				Title: post.Title,
				Href:  post.HTMLFileName(),
				Date:  post.Date.Format(dateFormat),
			})
		}
	}

	tmpl := templates.Load()
	outb := bytes.NewBuffer(make([]byte, 1024))
	if err := tmpl.Execute(outb, args); err != nil {
		// TODO(nikonov): RETURN ERROR
		panic(err)
	}
	return outb.Bytes(), nil
}
