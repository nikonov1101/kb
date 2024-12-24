package tool

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// ListSources return sorted list of notes in the source directory
func ListSources(root string, withPrivate bool) ([]Source, error) {
	fs, err := listSources(root, withPrivate)
	if err != nil {
		return nil, err
	}

	return fs, nil
}

func listSources(root string, withPrivate bool) ([]Source, error) {
	var sourceFiles []Source
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".md") {
			source, err := loadSourceFile(path)
			if err != nil {
				return err
			}

			if source.Visibility != Published && !withPrivate {
				return nil
			}

			sourceFiles = append(sourceFiles, source)
		}
		return err
	})

	sort.SliceStable(sourceFiles, func(i, j int) bool {
		return sourceFiles[i].Num < sourceFiles[j].Num
	})

	return sourceFiles, err
}
