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
	"gopkg.in/russross/blackfriday.v2"

	"github.com/nikonov1101/kb/templates"
)

const (
	published = "published"
	private   = "private"
	hidden    = "none"

	dateFormat = "02 Jan 2006"
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

	// note title, from headers
	title string
	// publication date described by dateFormat, from headers
	date time.Time
	// visibility of rendered page:
	// "published" - rendered, listed on index page
	// "private" - rendered, not listed on index page
	// "none" - not rendered
	visibility string
	// raw markdown content
	markdown []byte
	// cached generated html
	html []byte
}

func (s *Source) pageURI() string {
	return s.baseName + ".html"
}

// parseFile turns the sourceBytes into
// ready-to-render *Source
func parseFile(filePath string, sourceBytes []byte) (*Source, error) {
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
		return nil, fmt.Errorf("missing `---` separator in a file")
	}

	// split content into parts
	body := sourceBytes[delim+3:]
	rawHeaders := sourceBytes[:delim]
	lines := bytes.Split(rawHeaders, []byte("\n"))

	source := &Source{
		markdown: body,
		path:     filePath,
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
			source.title = value
		case "date":
			v, err := time.Parse(dateFormat, value)
			if err != nil {
				return nil, fmt.Errorf("failed to parse %q as %q: %v", value, dateFormat, err)
			}
			source.date = v
		case "visibility":
			switch value {
			case published, private:
				source.visibility = value
			default:
				source.visibility = hidden
			}
		default:
			return nil, fmt.Errorf("unknown tag %s", parts[0])
		}
	}

	source.html = markdownToHTML(source.markdown)

	return source, nil
}

// loadSourceFile from disk and parse its content
func loadSourceFile(filePath string) (*Source, error) {
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
func generatePage(src *Source) []byte {
	tmpl := bytes.ReplaceAll(templates.Page, []byte("${TITLE}"), []byte(src.title))
	tmpl = bytes.ReplaceAll(tmpl, []byte("${DATE}"), []byte(src.date.Format(dateFormat)))
	tmpl = bytes.ReplaceAll(tmpl, []byte("${CONTENT}"), src.html)

	return tmpl
}

// generateIndex generate index page with links to notes given as `fs`
func generateIndex(sources []*Source) []byte {
	const template = `<div class="post-link"><a href="%s">%04d: %s</a><div class="post-date">%s</div></div>`

	linksHTML := ""
	for i := len(sources) - 1; i >= 0; i-- {
		src := sources[i]
		if src.visibility == private {
			// do not list private notes
			continue
		}

		linksHTML += fmt.Sprintf(template, src.pageURI(), src.num, src.title, src.date.Format(dateFormat))
	}

	index := Source{
		title: siteDescription,
		html:  []byte(linksHTML),
	}

	return generatePage(&index)
}

func Generate(srcDir, dstDir string) error {
	start := time.Now()

	if err := os.RemoveAll(dstDir); err != nil {
		return fmt.Errorf("failed to clean the dst dir: %v", err)
	}

	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return fmt.Errorf("failed to create the dst dir: %v", err)
	}

	list, err := listSources(srcDir)
	if err != nil {
		return err
	}

	fmt.Printf("Found %s source file(s)\n", green(strconv.Itoa(len(list))))
	for _, src := range list {
		if src.visibility == hidden {
			fmt.Printf("  %s hidden %s...\n", yellow("skipping"), src.path)
			continue
		}

		fmt.Printf("  %s %s...\n", green("processing"), src.path)

		// add content html to the rest of the page
		page := generatePage(src)
		// where to store HTML result
		out := src.pageURI()

		if err := os.WriteFile(path.Join(dstDir, out), page, 0644); err != nil {
			return fmt.Errorf("failed to write to %s: %v", out, err)
		}
	}

	// TODO: generate rss feed as well, for compatibility?
	atomFeed := generateFeeds(list)
	atomFeedPath := path.Join(dstDir, "atom.xml")
	fmt.Printf("%s %s...\n", green("processing"), atomFeedPath)
	if err := os.WriteFile(atomFeedPath, atomFeed, 0644); err != nil {
		return fmt.Errorf("failed to write atom.xml: %v", err)
	}

	index := generateIndex(list)
	indexPath := path.Join(dstDir, "index.html")
	fmt.Printf("%s %s...\n", green("processing"), indexPath)
	if err := os.WriteFile(indexPath, index, 0644); err != nil {
		return fmt.Errorf("failed to write index.html: %v", err)
	}

	fmt.Printf("Done in %s.\n", green(time.Since(start).String()))
	return nil
}
