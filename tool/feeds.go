package tool

import (
	"os"
	"path"
	"time"

	"github.com/gopherlibs/feedhub/feedhub"
	"github.com/pkg/errors"
)

func GenerateRSSFeed(items []Source, destDir string, siteName string, baseURL string) error {
	author := &feedhub.Author{
		Name:  "Alex Nikonov",
		Email: "alex@nikonov.tech",
	}

	feed := feedhub.Feed{
		// seems like it's better to sign your blog with your name,
		// at least in people' RSS feeds.
		Title:       "alex nikonov",
		Link:        &feedhub.Link{Href: baseURL},
		Description: siteName,
		Author:      author,
		Created:     time.Now().UTC(),
	}

	for _, note := range items {
		if note.Visibility == Published {
			feed.Items = append(feed.Items, &feedhub.Item{
				Title:       note.Title,
				Description: note.Title,
				Link:        &feedhub.Link{Href: note.URL(baseURL)},
				Author:      author,
				Created:     note.Date,
				Content:     string(note.html),
			})
		}
	}

	atom, err := feed.ToAtom()
	if err != nil {
		return errors.Wrap(err, "generate atom feed")
	}

	atomFeedPath := path.Join(destDir, "atom.xml")
	if err := os.WriteFile(atomFeedPath, []byte(atom), 0o644); err != nil {
		return errors.Wrapf(err, "write %s", atomFeedPath)
	}

	return nil
}
