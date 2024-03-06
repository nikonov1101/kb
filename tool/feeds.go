package tool

import (
	"time"

	"github.com/gopherlibs/feedhub/feedhub"
)

func generateFeeds(items []*Source) []byte {
	author := &feedhub.Author{
		Name:  "Alex Nikonov",
		Email: "alex@nikonov.tech",
	}

	feed := feedhub.Feed{
		Title:       "nikonov blog",
		Link:        &feedhub.Link{Href: baseURL},
		Description: siteDescription,
		Author:      author,
		Created:     time.Now().UTC(),
	}

	for _, it := range items {
		feed.Items = append(feed.Items, &feedhub.Item{
			Title:       it.title,
			Description: it.title,
			Link:        &feedhub.Link{Href: baseURL + it.pageURI()},
			Author:      author,
			Created:     it.date,
			Content:     string(it.html),
		})
	}

	atom, err := feed.ToAtom()
	if err != nil {
		panic("failed to marshal feed: " + err.Error())
	}
	return []byte(atom)
}
