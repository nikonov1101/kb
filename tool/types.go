package tool

import "time"

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
	// isIndex tells renderer that this page is index page,
	// so extra rules might be applied.
	isIndex bool
}

func (s *Source) HTMLFileName() string {
	return s.BaseName + ".html"
}

func (s Source) URL(root string) string {
	return root + "/" + s.HTMLFileName()
}
