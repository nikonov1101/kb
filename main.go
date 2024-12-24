package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/nikonov1101/colors.go"
	"github.com/nikonov1101/kb/tool"
	"github.com/nikonov1101/kb/version"
)

func main() {
	rootCmd := &cobra.Command{
		Use: "kb",
	}

	absp, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	srcDir := rootCmd.PersistentFlags().String("src", absp, "directory with markdown files")
	dstDir := rootCmd.PersistentFlags().String("www", filepath.Join(absp, "/www"), "directory for generated html files")
	baseURL := rootCmd.PersistentFlags().String("base-url", "http://localhost:8000", "base site URL")
	siteName := rootCmd.PersistentFlags().String("site-name", "Make computers fun again", "site name")

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "Show list of notes in the source dir",
		RunE: func(cmd *cobra.Command, args []string) error {
			list, err := tool.ListSources(*srcDir, true)
			if err != nil {
				return err
			}

			for _, file := range list {
				pre := fmt.Sprintf("%04d", file.Num)
				switch file.Visibility {
				case tool.Published:
					pre = colors.BGreen(pre)
				case tool.Private:
					pre = colors.BYellow(pre)
				default:
					pre = colors.BRed(pre)
				}

				fmt.Printf("%s :: %s :: %s\n", pre, colors.BWhite(file.Title), file.URL(*baseURL))
			}

			return nil
		},
	}

	generateCmd := &cobra.Command{
		Use:     "gen",
		Aliases: []string{"generate", "build"},
		Short:   "Generate site content",
		RunE: func(cmd *cobra.Command, args []string) error {
			isPrivate := cmd.Flag("private").Value.String() == "true"
			list, err := tool.ListSources(*srcDir, isPrivate)
			if err != nil {
				return err
			}

			fmt.Printf("source: %s, files %s\n", colors.Green(*srcDir), colors.Yellow(strconv.Itoa(len(list))))

			if err := tool.GeneratePages(list, *dstDir, *siteName, *baseURL); err != nil {
				return errors.Wrap(err, "generate pages")
			}

			if err := tool.GenerateIndex(list, *dstDir, *siteName); err != nil {
				return errors.Wrap(err, "generate index")
			}

			if err := tool.GenerateRSSFeed(list, *dstDir, *siteName, *baseURL); err != nil {
				return errors.Wrap(err, "generate RSS")
			}

			return nil
		},
	}
	generateCmd.PersistentFlags().Bool("private", false, "render private notes")

	newCmd := &cobra.Command{
		Use:   "new <name>",
		Short: fmt.Sprintf("Create new empty note and open in %q", tool.EDITOR),
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			isPrivate := cmd.Flag("private").Value.String() == "true"
			isRandomName := cmd.Flag("rand").Value.String() == "true"
			isEdit := cmd.Flag("edit").Value.String() == "true"
			isOpenBrowser := cmd.Flag("web").Value.String() == "true"

			var name string
			if isRandomName {
				name = uuid.New().String()
			} else {
				if len(args) != 1 {
					return errors.New("accepts 1 arg(s), received 0")
				}
				name = args[0]
			}

			diskPath, webPath, err := tool.New(*srcDir, name, isPrivate)
			if err != nil {
				return err
			}

			if isEdit {
				if !isOpenBrowser {
					// edit but not open browser: block by a editor process
					if err := openEditor(diskPath); err != nil {
						fmt.Printf("failed to open editor: %v\n", err)
					}
				} else {
					go func() {
						// edit and view in browser: detach editor, block by a web server process
						time.Sleep(200 * time.Millisecond)
						if err := openEditor(diskPath); err != nil {
							fmt.Printf("failed to open browser: %v\n", err)
						}
					}()
				}
			}

			if isOpenBrowser {
				listen := cmd.Flag("addr").Value.String()
				go func() {
					time.Sleep(50 * time.Millisecond)
					openURL := "http://" + listen + "/" + webPath
					if err := openBrowser(openURL); err != nil {
						fmt.Printf("failed to invoke `open` command: %v\n", err)
					}
				}()

				fmt.Printf("starting web-server on %s ...\n", colors.BGreen("http://"+listen))
				return tool.Serve(*srcDir, *dstDir, listen, *siteName, *baseURL, isPrivate)
			}
			return nil
		},
	}
	newCmd.PersistentFlags().String("addr", "127.0.0.1:8000", "address to listen to")
	newCmd.PersistentFlags().Bool("private", true, "create as private or public")
	newCmd.PersistentFlags().Bool("rand", false, "use random temporary id")
	newCmd.PersistentFlags().Bool("edit", true, "open new note in editor")
	newCmd.PersistentFlags().Bool("web", false, "open a browser with new note")

	serveCmd := &cobra.Command{
		Use:   "serve",
		Short: "Generate pages and start http server in result dir",
		RunE: func(cmd *cobra.Command, args []string) error {
			isPrivate := cmd.Flag("private").Value.String() == "true"
			isOpenBrowser := cmd.Flag("web").Value.String() == "true"
			listen := cmd.Flag("addr").Value.String()

			if isOpenBrowser {
				go func() {
					time.Sleep(50 * time.Millisecond)
					openURL := "http://" + listen
					if err := openBrowser(openURL); err != nil {
						fmt.Printf("failed to open browser: %v\n", err)
					}
				}()
			}

			fmt.Printf("starting web-server on %s ...\n", colors.BGreen("http://"+listen))
			return tool.Serve(*srcDir, *dstDir, listen, *siteName, *baseURL, isPrivate)
		},
	}
	serveCmd.PersistentFlags().String("addr", "127.0.0.1:8000", "address to listen to")
	serveCmd.PersistentFlags().Bool("private", false, "render private notes")
	serveCmd.PersistentFlags().Bool("web", false, "open a browser with new note")

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "show version, build info, and current configuration parameters",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Printf("Built with: %s\n", version.CompillerVersion())
			cmd.Printf("Built at:   %s\n", version.BuildTime())
			cmd.Printf("Version:    %s\n", version.BuildCommit())
			cmd.Printf("Source dir: %s\n", *srcDir)
			return nil
		},
	}

	rootCmd.AddCommand(
		listCmd,
		generateCmd,
		newCmd,
		serveCmd,
		versionCmd,
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func openBrowser(url string) error {
	if err := exec.Command("/usr/bin/open", url).Run(); err != nil {
		return errors.Wrapf(err, "open %q in browser", url)
	}
	return nil
}

func openEditor(filePath string) error {
	if err := exec.Command(tool.EDITOR, "-a", filePath).Run(); err != nil {
		return errors.Wrapf(err, "open %q in editor", filePath)
	}
	return nil
}
