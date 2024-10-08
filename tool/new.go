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

const dateFormat = "02 Jan 2006"

func New(root, name string, isPrivate bool) (string, string, error) {
	if strings.Contains(name, " ") {
		return "", "", fmt.Errorf("name must not contain spaces")
	}

	sources, err := listSources(root)
	if err != nil {
		return "", "", err
	}

	var next int64
	if len(sources) > 0 {
		next = sources[len(sources)-1].num + 1
	}

	visibility := published
	if isPrivate {
		visibility = private
	}

	fileName := fmt.Sprintf("%04d-%s.md", next, name)
	webPath := fmt.Sprintf("%04d-%s.html", next, name)

	diskPath := path.Join(root, fileName)
	fmt.Println("new file: ", diskPath)

	date := time.Now().Format(dateFormat)
	content := strings.ReplaceAll(newFileContent, "$date$", date)
	content = strings.ReplaceAll(content, "$name$", name)
	content = strings.ReplaceAll(content, "$visibility$", visibility)
	if err := os.WriteFile(diskPath, []byte(content), 0o644); err != nil {
		return "", "", fmt.Errorf("failed to write to %s: %v", diskPath, err)
	}

	return diskPath, webPath, nil
}
