package tool

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"gopkg.in/russross/blackfriday.v2"
)

type Source struct {
	// relative path on disk,
	// including source directory
	path string
	// file name without source dir
	// and the .md extension
	baseName string
	// sequential file number
	num int64

	// title from headers
	title string
	// tags list from headers
	tags []string
	// private flag from headers
	private bool
	// raw markdown content
	data []byte
}

// parseFile turns the sourceBytes into
// ready-to-render *Source
func parseFile(p string, sourceBytes []byte) (*Source, error) {
	// parse path and parts
	baseName := strings.Split(path.Base(p), ".")[0]
	// parse file number,
	parts := strings.SplitN(baseName, "-", 2)
	num, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		panic("malformed file name in " + p)
	}

	// look for headers
	delim := bytes.Index(sourceBytes, []byte("---"))
	if delim == -1 {
		return nil, fmt.Errorf("missing `---` separator in a file")
	}

	// split content into parts
	body := sourceBytes[delim+3:]
	rawHeaders := sourceBytes[:delim]
	lines := bytes.Split(rawHeaders, []byte("\n"))

	m := &Source{
		data:     body,
		path:     p,
		baseName: baseName,
		num:      num,
	}

	// parse headers
	for i, ln := range lines {
		if len(ln) == 0 {
			break
		}
		s := string(ln)
		parts := strings.SplitN(s, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("malformed tag line at %d (%s): no separator", i+1, s)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "title":
			m.title = value
		case "tags":
			tags := strings.Split(value, " ")
			for _, t := range tags {
				m.tags = append(m.tags, strings.TrimSpace(t))
			}
		case "private":
			m.private = value == "true"
		default:
			return nil, fmt.Errorf("unknown tag %s", parts[0])
		}
	}

	return m, nil
}

// openSourceFile loads source file from a disk and parse its headers
func openSourceFile(p string) (*Source, error) {
	text, err := ioutil.ReadFile(p)
	if err != nil {
		panic(err)
	}

	meta, err := parseFile(p, text)
	if err != nil {
		panic(err)
	}

	meta.path = p
	return meta, nil
}

// renderBody turns raw markdown bytes into the HTML markup
func renderBody(data []byte) []byte {
	renderer := blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{
		// Flags: blackfriday.HrefTargetBlank | blackfriday.NofollowLinks | blackfriday.NoreferrerLinks | blackfriday.TOC,
		Flags: blackfriday.HrefTargetBlank | blackfriday.NofollowLinks | blackfriday.NoreferrerLinks,
	})

	unsafe := blackfriday.Run(
		data,
		blackfriday.WithExtensions(blackfriday.CommonExtensions|blackfriday.AutoHeadingIDs),
		blackfriday.WithRenderer(renderer))

	// allow applying various CSS classes to colorize
	// content in <code> blocks
	p := bluemonday.UGCPolicy()
	p.AllowAttrs("class").Matching(regexp.MustCompile("^language-[a-zA-Z0-9]+$")).OnElements("code")
	p.AllowStandardURLs()
	p.AddTargetBlankToFullyQualifiedLinks(true)

	return p.SanitizeBytes(unsafe)
}

// renderTemplate adds Source's headers and HTML body
// to the template, returns fill page data
func renderTemplate(src *Source, html, tmpl []byte) []byte {
	title := []byte("${TITLE}")
	tags := []byte("${TAGS}")
	content := []byte("${CONTENT}")
	t := strings.Join(src.tags, ", ")

	tmpl = bytes.Replace(tmpl, title, []byte(src.title), -1)
	tmpl = bytes.Replace(tmpl, tags, []byte(t), -1)
	tmpl = bytes.Replace(tmpl, content, html, -1)

	return tmpl
}

func Generate(srcDir, dstDir string) error {
	start := time.Now()

	pageTemplate, err := ioutil.ReadFile("templates/page.html")
	if err != nil {
		return fmt.Errorf("failed to read page template: %v", err)
	}

	indexTemplate, err := ioutil.ReadFile("templates/index.html")
	if err != nil {
		return fmt.Errorf("failed to read the index template: %v", err)
	}

	if err := os.RemoveAll(dstDir); err != nil {
		return fmt.Errorf("failed to clean the dst dir: %v", err)
	}

	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return fmt.Errorf("failed to create the dst dir: %v", err)
	}

	fs, err := list(srcDir)
	if err != nil {
		return err
	}

	fmt.Printf("Found %d source files\n", len(fs))

	notesListHTML := "<ul>\n"
	for i := len(fs) - 1; i >= 0; i-- {
		f := fs[i]
		fmt.Printf("  processing %q (%s) ...\n", f.title, f.path)

		// generate HTML from source's markdown
		html := renderBody(f.data)
		html = renderTemplate(f, html, pageTemplate)

		base := f.baseName + ".html"
		out := path.Join(dstDir, base)
		if err := ioutil.WriteFile(out, html, 0644); err != nil {
			return fmt.Errorf("failed to write to %s: %v", out, err)
		}

		notesListHTML += fmt.Sprintf(`<li><a href="%s">%04d: %s %s</a></li>`+"\n", base, f.num, f.title, f.tags)
		fmt.Printf("  output in %s\n", path.Clean(out))
	}

	// append notes list to the template
	notesListHTML += "</ul>\n"
	ix := strings.Replace(string(indexTemplate), "${CONTENT}", notesListHTML, 1)
	if err := ioutil.WriteFile(path.Join(dstDir, "index.html"), []byte(ix), 0644); err != nil {
		return fmt.Errorf("failed to write index.html: %v", err)
	}

	fmt.Printf("Done in %s.\n", time.Since(start))
	return nil
}
