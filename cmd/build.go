/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/cpunion/go-python/cmd/internal/rungo"
	"github.com/spf13/cobra"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build [flags] [package]",
	Short: "Build a Go package with Python environment configured",
	Long: `Build compiles a Go package with the Python environment properly configured.
If package is a directory, it will be used directly. Otherwise, the directory
containing the package will be determined.

Example:
  gopy build .
  gopy build -o myapp ./cmd/myapp
  gopy build -v -race .`,
	DisableFlagParsing: true, // Let go build handle all flags
	Run: func(cmd *cobra.Command, args []string) {
		if err := rungo.RunGoCommand("build", args); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
}
