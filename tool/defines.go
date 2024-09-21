package tool

import (
	"os"
)

var (
	rssFeedURL      = "http://localhost:8000/"
	siteDescription = "Making computers fun again"
)

func init() {
	if v := os.Getenv("KB_URL"); v != "" {
		rssFeedURL = v
	}
	if v := os.Getenv("KB_NAME"); v != "" {
		siteDescription = v
	}
}
