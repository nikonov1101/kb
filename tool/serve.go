package tool

import (
	"fmt"
	"net/http"
	"os/exec"

	"github.com/fsnotify/fsnotify"
)

func Serve(src, dst string, listenAddr string) error {
	if err := Generate(src, dst); err != nil {
		return err
	}

	fmt.Printf("Starting http server on %s...\n", green(listenAddr))
	http.Handle("/", http.FileServer(http.Dir(dst)))

	watchFs(src, func(s string) {
		fmt.Printf("\nfsnotify: write event on %s, generating...\n", yellow(s))
		if err := Generate(src, dst); err != nil {
			fmt.Printf("  Generate() failed: %v\n", err)
		}
	})

	// note: macos only
	if err := exec.Command("/usr/bin/open", "http://"+listenAddr).Run(); err != nil {
		fmt.Printf("failed to invoke `open` command: %v\n", err)
	}

	return http.ListenAndServe(listenAddr, nil)
}

func watchFs(dir string, onWrite func(string)) {
	fmt.Printf("Starting fs watcher on %s\n", green(dir))
	w, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Printf("fsnotify: failed to create watcher: %v", err)
		return
	}

	go func() {
		for {
			select {
			case evt, ok := <-w.Events:
				if !ok {
					return
				}
				if evt.Op&fsnotify.Write == fsnotify.Write {
					onWrite(evt.Name)
				}
			case err, ok := <-w.Errors:
				if !ok {
					return
				}
				fmt.Printf("fsnotify err: %v\n", err)
			}
		}
	}()

	if err := w.Add(dir); err != nil {
		fmt.Printf("fsnotify: failed to add event watcher: %v", err)
		return
	}
}
