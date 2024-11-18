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
	Long: func() string {
		intro := "Build compiles a Go package with the Python environment properly configured.\n\n"
		help, err := rungo.GetGoCommandHelp("build")
		if err != nil {
			return intro + "Failed to get go help: " + err.Error()
		}
		return intro + help
	}(),
	DisableFlagParsing: true,
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
