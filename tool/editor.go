package tool

import (
	"log"
	"os"
)

const defaultEditor = "/usr/bin/vi"

func Editor() string {
	ed := os.Getenv("EDITOR")
	if len(ed) == 0 {
		log.Printf("WARN: no $EDITOR, fallback to %q", defaultEditor)
		ed = defaultEditor
	}

	return ed
}
