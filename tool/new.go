package tool

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"
)

const newFileContent = `title: name me
tags: foo bar
private: false
---
`

func New(root, name string) (string, error) {
	if strings.Contains(name, " ") {
		return "", fmt.Errorf("name must not contain spaces")
	}

	fs, err := list(root)
	if err != nil {
		return "", err
	}

	next := fs[len(fs)-1].num + 1

	p := path.Join(root, fmt.Sprintf("%04d-%s.md", next, name))
	fmt.Println("new file at", p)

	if err := ioutil.WriteFile(p, []byte(newFileContent), 0644); err != nil {
		return "", fmt.Errorf("failed to write to %s: %v", p, err)
	}

	return p, nil
}
