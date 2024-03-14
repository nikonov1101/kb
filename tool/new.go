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
visibility: published
---
`

const headerDateFormat = "02 Jan 2006"

func New(root, name string) (string, error) {
	if strings.Contains(name, " ") {
		return "", fmt.Errorf("name must not contain spaces")
	}

	sources, err := listSources(root)
	if err != nil {
		return "", err
	}

	var next int64
	if len(sources) > 0 {
		next = sources[len(sources)-1].num + 1
	}

	p := path.Join(root, fmt.Sprintf("%04d-%s.md", next, name))
	fmt.Println("new file at", p)

	date := time.Now().Format(headerDateFormat)
	content := strings.ReplaceAll(newFileContent, "$date$", date)
	content = strings.ReplaceAll(content, "$name$", name)
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("failed to write to %s: %v", p, err)
	}

	return p, nil
}
