package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"

	"github.com/nikonov1101/kb/tool"
)

var rootDir = "/Users/alex/src/kb"

func init() {
	if root := os.Getenv("KB_ROOT"); len(root) > 0 {
		rootDir = root
	}
}

func main() {
	var rootCmd = &cobra.Command{
		Use: "kb",
	}

	srcDir := rootCmd.PersistentFlags().String("src", filepath.Join(rootDir, "/src"), "directory with markdown files")
	dstDir := rootCmd.PersistentFlags().String("www", filepath.Join(rootDir, "/www"), "directory with generated html files")

	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "Show list of notes in the source dir",
		RunE: func(cmd *cobra.Command, args []string) error {
			return tool.ListSources(*srcDir)
		},
	}

	var generateCmd = &cobra.Command{
		Use:     "gen",
		Aliases: []string{"generate", "build"},
		Short:   "Generate site content",
		RunE: func(cmd *cobra.Command, args []string) error {
			return tool.Generate(*srcDir, *dstDir)
		},
	}

	var newCmd = &cobra.Command{
		Use:   "new <name>",
		Short: fmt.Sprintf("Create new empty note and open in %q", tool.EDITOR),
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			diskPath, webPath, err := tool.New(*srcDir, args[0])
			if err != nil {
				return err
			}

			go func() {
				time.Sleep(200 * time.Millisecond)
				if err := exec.Command(tool.EDITOR, "-a", diskPath).Run(); err != nil {
					fmt.Printf("failed to open editor: %v\n", err)
				}
			}()

			listen := cmd.Flag("addr").Value.String()
			openURL := "http://" + listen + "/" + webPath
			return tool.Serve(
				*srcDir,
				*dstDir,
				listen,
				openURL,
			)
		},
	}
	newCmd.PersistentFlags().String("addr", "127.0.0.1:8000", "address to listen to")

	var serveCmd = &cobra.Command{
		Use:   "serve",
		Short: "Generate pages and start http server in result dir",
		RunE: func(cmd *cobra.Command, args []string) error {
			listen := cmd.Flag("addr").Value.String()
			openURL := "http://" + listen
			return tool.Serve(
				*srcDir,
				*dstDir,
				listen,
				openURL,
			)
		},
	}
	serveCmd.PersistentFlags().String("addr", "127.0.0.1:8000", "address to listen to")

	rootCmd.AddCommand(
		listCmd,
		generateCmd,
		newCmd,
		serveCmd,
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
