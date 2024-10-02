package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/google/uuid"
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
	rootCmd := &cobra.Command{
		Use: "kb",
	}

	srcDir := rootCmd.PersistentFlags().String("src", filepath.Join(rootDir, "/src"), "directory with markdown files")
	dstDir := rootCmd.PersistentFlags().String("www", filepath.Join(rootDir, "/www"), "directory with generated html files")

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "Show list of notes in the source dir",
		RunE: func(cmd *cobra.Command, args []string) error {
			return tool.ListSources(*srcDir)
		},
	}

	generateCmd := &cobra.Command{
		Use:     "gen",
		Aliases: []string{"generate", "build"},
		Short:   "Generate site content",
		RunE: func(cmd *cobra.Command, args []string) error {
			isPrivate := cmd.Flag("private").Value.String() == "true"
			return tool.Generate(*srcDir, *dstDir, isPrivate)
		},
	}
	generateCmd.PersistentFlags().Bool("private", false, "render private notes")

	newCmd := &cobra.Command{
		Use:   "new <name>",
		Short: fmt.Sprintf("Create new empty note and open in %q", tool.EDITOR),
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			isPrivate := cmd.Flag("private").Value.String() == "true"
			name := uuid.New().String()
			if !isPrivate {
				if len(args) != 1 {
					return errors.New("accepts 1 arg(s), received 0")
				}
			}

			diskPath, webPath, err := tool.New(*srcDir, name, isPrivate)
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
				isPrivate,
			)
		},
	}
	newCmd.PersistentFlags().String("addr", "127.0.0.1:8000", "address to listen to")
	newCmd.PersistentFlags().Bool("private", false, "render private notes")

	serveCmd := &cobra.Command{
		Use:   "serve",
		Short: "Generate pages and start http server in result dir",
		RunE: func(cmd *cobra.Command, args []string) error {
			isPrivate := cmd.Flag("private").Value.String() == "true"
			listen := cmd.Flag("addr").Value.String()
			openURL := "http://" + listen
			return tool.Serve(
				*srcDir,
				*dstDir,
				listen,
				openURL,
				isPrivate,
			)
		},
	}
	serveCmd.PersistentFlags().String("addr", "127.0.0.1:8000", "address to listen to")
	serveCmd.PersistentFlags().Bool("private", false, "render private notes")

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
