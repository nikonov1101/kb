package tool

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/nikonov1101/kb/tool/templates"
	"github.com/pkg/errors"
	"gopkg.in/russross/blackfriday.v2"
)

// GeneratePages site content using notes in srcDir, saving html files in dstDir
func GeneratePages(notes []Source, destDir string, siteName string, baseURL string) error {
	if err := os.RemoveAll(destDir); err != nil {
		return fmt.Errorf("failed to clean the dst dir: %v", err)
	}

	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return fmt.Errorf("failed to create the dst dir: %v", err)
	}

	for _, src := range notes {
		// add content html to the rest of the page
		page := generatePage(siteName, src)
		// where to store HTML result
		out := src.HTMLFileName()

		if err := os.WriteFile(path.Join(destDir, out), page, 0o644); err != nil {
			return fmt.Errorf("failed to write to %s: %v", out, err)
		}
	}

	return nil
}

// GenerateIndex generate index page with links to notes given as `fs`
func GenerateIndex(sources []Source, destDir string, siteName string) error {
	const template = `<div class="post"><a href="%s">%s</a><span class="post-date">%s</span></div>`

	var linksHTML string
	for i := len(sources) - 1; i >= 0; i-- {
		src := sources[i]
		linksHTML += fmt.Sprintf(
			template,
			src.HTMLFileName(),
			src.Title,
			displayDate(src.Date),
		) + "\n"
	}

	index := Source{
		html:    []byte(linksHTML),
		isIndex: true,
	}

	page := generatePage(siteName, index)
	indexPath := path.Join(destDir, "index.html")
	if err := os.WriteFile(indexPath, page, 0o644); err != nil {
		return errors.Wrapf(err, "write %s", indexPath)
	}

	return nil
}

// parseFile turns the sourceBytes into
// ready-to-render *Source
func parseFile(filePath string, sourceBytes []byte) (Source, error) {
	// parse path and parts
	baseName := strings.Split(path.Base(filePath), ".")[0]
	// parse file number,
	parts := strings.SplitN(baseName, "-", 2)
	num, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		panic("malformed file name: " + filePath)
	}

	// look for headers
	delim := bytes.Index(sourceBytes, []byte("---"))
	if delim == -1 {
		return Source{}, fmt.Errorf("missing `---` separator in a file")
	}

	// split content into parts
	body := sourceBytes[delim+3:]
	rawHeaders := sourceBytes[:delim]
	lines := bytes.Split(rawHeaders, []byte("\n"))

	source := Source{
		markdown: body,
		Path:     filePath,
		BaseName: baseName,
		Num:      num,
	}

	// parse headers
	for i, ln := range lines {
		if len(ln) == 0 {
			break
		}
		s := string(ln)
		parts := strings.SplitN(s, ":", 2)
		if len(parts) != 2 {
			return Source{}, fmt.Errorf("malformed tag line at %d (%s): no separator", i+1, s)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "title":
			source.Title = value
		case "date":
			v, err := time.Parse(dateFormat, value)
			if err != nil {
				return Source{}, fmt.Errorf("failed to parse %q as %q: %v", value, dateFormat, err)
			}
			source.Date = v
		case "visibility":
			switch value {
			case Published, Private:
				source.Visibility = value
			default:
				source.Visibility = Hidden
			}
		default:
			return Source{}, fmt.Errorf("unknown tag %s", parts[0])
		}
	}

	source.html = markdownToHTML(source.markdown)
	return source, nil
}

// loadSourceFile from disk and parse its content
func loadSourceFile(filePath string) (Source, error) {
	text, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	source, err := parseFile(filePath, text)
	if err != nil {
		panic(err)
	}

	return source, nil
}

// markdownToHTML turns raw markdown bytes into the HTML markup
func markdownToHTML(data []byte) []byte {
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

// generatePage makes a complete web-page from given source
func generatePage(siteName string, src Source) []byte {
	title := src.Title
	pageTitle := fmt.Sprintf(`<div class="heading"><h1>%s</h1></div>`, title)
	about := ""

	if src.isIndex {
		title = siteName
		pageTitle = ""
		about = templates.Intro
	}

	timeStr := ""
	if !src.Date.IsZero() {
		// do not render date string for index page
		timeStr = src.Date.Format(dateFormat)
	}

	tmpl := bytes.ReplaceAll(templates.Page, []byte("${TITLE}"), []byte(title))
	tmpl = bytes.ReplaceAll(tmpl, []byte("${PAGE_TITLE}"), []byte(pageTitle))
	tmpl = bytes.ReplaceAll(tmpl, []byte("${DATE}"), []byte(displayDate(src.Date)))
	tmpl = bytes.ReplaceAll(tmpl, []byte("${CONTENT}"), src.html)
	tmpl = bytes.ReplaceAll(tmpl, []byte("${TIMESTAMP}"), []byte(timeStr))
	tmpl = bytes.ReplaceAll(tmpl, []byte("${INDEX_ABOUT}"), []byte(about))

	return tmpl
}

// displayDate returns localized date string for displaying in templates
func displayDate(t time.Time) string {
	return t.Format(dateFormat)
}
