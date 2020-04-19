package main

import (
	"fmt"
	"net/http"
	"os"

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
			return tool.New(*srcDir, args[0])
		},
	}

	var serveCmd = &cobra.Command{
		Use:   "serve",
		Short: "Serve files from the result dir",
		RunE: func(cmd *cobra.Command, args []string) error {
			addr := cmd.Flag("addr").Value.String()
			fmt.Printf("Starting http server on %s...\n", addr)
			http.Handle("/", http.FileServer(http.Dir(*dstDir)))
			return http.ListenAndServe(addr, nil)
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
