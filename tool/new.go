package tool

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"
)

const newFileContent = `title: $name$
date: $date$
visibility: $visibility$ 
---
`

// New creates new note in srcDir with a given name,
// returns path on dist and page URL.
func New(srcDir, name string, isPrivate bool) (string, string, error) {
	name = strings.ReplaceAll(name, " ", "-")
	sources, err := listSources(srcDir, true)
	if err != nil {
		return "", "", err
	}

	var next int64
	if len(sources) > 0 {
		next = sources[len(sources)-1].Num + 1
	}

	visibility := Published
	if isPrivate {
		visibility = Private
	}

	fileName := fmt.Sprintf("%04d-%s.md", next, name)
	webPath := fmt.Sprintf("%04d-%s.html", next, name)
	diskPath := path.Join(srcDir, fileName)

	date := time.Now().Format(dateFormat)
	content := strings.ReplaceAll(newFileContent, "$date$", date)
	content = strings.ReplaceAll(content, "$name$", name)
	content = strings.ReplaceAll(content, "$visibility$", visibility)

	if err := os.WriteFile(diskPath, []byte(content), 0o644); err != nil {
		return "", "", fmt.Errorf("failed to write to %s: %v", diskPath, err)
	}

	return diskPath, webPath, nil
}
