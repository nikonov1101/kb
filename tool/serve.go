package tool

import (
	"fmt"
	"log"
	"net/http"

	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
)

// Serve site content on given listenAddr, re-generate each time content of src directory changed.
func Serve(src, dst, listenAddr, siteName, baseURL string, withPrivate bool) error {
	generateAll := func() error {
		list, err := ListSources(src, withPrivate)
		if err != nil {
			return errors.Wrapf(err, "list sources in %s", src)
		}
		if err := GeneratePages(list, dst, siteName, baseURL); err != nil {
			return errors.Wrap(err, "generate pages")
		}
		if err := GenerateIndex(list, dst, siteName); err != nil {
			return errors.Wrap(err, "generate index")
		}
		return nil
	}

	if err := generateAll(); err != nil {
		return err
	}

	watchFs(src, func(s string) {
		if err := generateAll(); err != nil {
			log.Printf("generate all: error: %v", err)
		}
	})

	http.Handle("/", http.FileServer(http.Dir(dst)))
	return http.ListenAndServe(listenAddr, nil)
}

func watchFs(dir string, onWrite func(string)) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Printf("fsnotify: failed to create watcher: %v", err)
		return
	}

	go func() {
		log.Printf("watching %q for changes...", dir)
		for {
			select {
			case evt, ok := <-w.Events:
				switch {
				case !ok:
					return
				// handle normal save, like sublime-text does
				case evt.Op&fsnotify.Write == fsnotify.Write:
					onWrite(evt.Name)
					return
				// handle vim-like save when temporary file with .ext~ is used
				case evt.Op&fsnotify.Create == fsnotify.Create:
					if evt.Name[len(evt.Name)-1] != '~' {
						onWrite(evt.Name)
						return
					}
				}

			case err, ok := <-w.Errors:
				if !ok {
					return
				}
				log.Printf("fsnotify err: %v\n", err)
			}
		}
	}()

	if err := w.Add(dir); err != nil {
		log.Printf("fsnotify: failed to add event watcher: %v", err)
		return
	}
}
