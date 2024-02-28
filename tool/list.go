package tool

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// returns sorted list of files in the source's root directory
func listSources(root string) ([]*Source, error) {
	var fs []*Source
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !(info.IsDir() || info.Name()[0] == '.') {
			ff, err := loadSourceFile(path)
			if err != nil {
				return err
			}

			fs = append(fs, ff)
		}
		return err
	})

	sort.SliceStable(fs, func(i, j int) bool {
		return fs[i].num < fs[j].num
	})

	return fs, err
}

func ListSources(root string) error {
	fs, err := listSources(root)
	if err != nil {
		return err
	}

	for _, ff := range fs {
		fmt.Printf("%s %q\n", yellow(ff.path), ff.title)
	}

	return nil
}
