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

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run [flags] [package] [arguments...]",
	Short: "Run a Go package with Python environment configured",
	Long: `Run executes a Go package with the Python environment properly configured.
If package is a directory, it will be used directly. Otherwise, the directory
containing the package will be determined.

Example:
  gopy run . arg1 arg2
  gopy run -race ./cmd/myapp arg1 arg2
  gopy run -v -race . arg1 arg2`,
	DisableFlagParsing: true, // Let go run handle all flags
	Run: func(cmd *cobra.Command, args []string) {
		if err := rungo.RunGoCommand("run", args); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
