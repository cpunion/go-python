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

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install [flags] [packages]",
	Short: "Install Go packages with Python environment configured",
	Long: `Install compiles and installs Go packages with the Python environment properly configured.
If package is a directory, it will be used directly. Otherwise, the directory
containing the package will be determined.

Example:
  gopy install .
  gopy install -v ./cmd/myapp
  gopy install -v -race .`,
	DisableFlagParsing: true, // Let go install handle all flags
	Run: func(cmd *cobra.Command, args []string) {
		if err := rungo.RunGoCommand("install", args); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
