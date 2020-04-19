package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"

	"github.com/sshaman1101/kb/tool"
)

//go:generate go run templates/gen.go

func main() {
	var rootCmd = &cobra.Command{
		Use: "kb",
	}

	srcDir := rootCmd.PersistentFlags().String("src", "./src", "path to source directory")
	dstDir := rootCmd.PersistentFlags().String("www", "./www", "path to results directory")

	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "Show list of notes in the source dir",
		RunE: func(cmd *cobra.Command, args []string) error {
			return tool.ListSources(*srcDir)
		},
	}

	var generateCmd = &cobra.Command{
		Use:   "gen",
		Short: "Generate site content",
		RunE: func(cmd *cobra.Command, args []string) error {
			return tool.Generate(*srcDir, *dstDir)
		},
	}

	var newCmd = &cobra.Command{
		Use:   "new <name>",
		Short: "Create new empty note",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			file, err := tool.New(*srcDir, args[0])
			if err != nil {
				return err
			}

			if err := exec.Command("/usr/local/bin/subl", "-a", file).Run(); err != nil {
				fmt.Printf("failed to open editor: %v\n", err)
			}

			return nil
		},
	}

	var serveCmd = &cobra.Command{
		Use:   "serve",
		Short: "Generate pages and start http server in result dir",
		RunE: func(cmd *cobra.Command, args []string) error {
			return tool.Serve(
				*srcDir,
				*dstDir,
				cmd.Flag("addr").Value.String(),
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
