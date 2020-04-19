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

	"github.com/sshaman1101/kb/templates"
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

func (s *Source) pageURI() string {
	return s.baseName + ".html"
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

// generatePage adds Source's headers and HTML body
// to the template, returns fill page data
func generatePage(src *Source, html []byte) []byte {
	title := []byte("${TITLE}")
	tags := []byte("${TAGS}")
	content := []byte("${CONTENT}")
	t := strings.Join(src.tags, ", ")

	tmpl := templates.Page

	tmpl = bytes.Replace(tmpl, title, []byte(src.title), -1)
	tmpl = bytes.Replace(tmpl, tags, []byte(t), -1)
	tmpl = bytes.Replace(tmpl, content, html, -1)

	return tmpl
}

// generateIndex generate index page with links to notes given as `fs`
func generateIndex(fs []*Source) []byte {
	buf := bytes.NewBufferString("<ul>\n")
	for i := len(fs) - 1; i >= 0; i-- {
		ff := fs[i]
		buf.WriteString(
			fmt.Sprintf("<li><a href=\"%s\">%04d: %s %s</a></li>\n",
				ff.pageURI(), ff.num, ff.title, ff.tags),
		)
	}
	buf.WriteString("</ul>\n")

	tmpl := templates.Index
	return bytes.Replace(tmpl, []byte("${CONTENT}"), buf.Bytes(), 1)
}

func Generate(srcDir, dstDir string) error {
	start := time.Now()

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

	fmt.Printf("Found %s source file(s)\n", green(strconv.Itoa(len(fs))))
	for _, f := range fs {
		fmt.Printf("  processing %s...\n", yellow(f.path))

		// generate HTML from source's markdown
		html := renderBody(f.data)
		// add content html to the rest of the page
		html = generatePage(f, html)
		// where to store HTML result
		out := f.pageURI()

		if err := ioutil.WriteFile(path.Join(dstDir, out), html, 0644); err != nil {
			return fmt.Errorf("failed to write to %s: %v", out, err)
		}
	}

	index := generateIndex(fs)
	indexPath := path.Join(dstDir, "index.html")
	fmt.Printf("Writing %s\n", yellow(indexPath))
	if err := ioutil.WriteFile(indexPath, index, 0644); err != nil {
		return fmt.Errorf("failed to write index.html: %v", err)
	}

	fmt.Printf("Done in %s.\n", green(time.Since(start).String()))
	return nil
}
