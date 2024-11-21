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

// execCmd represents the run command
var execCmd = &cobra.Command{
	Use:                "exec [flags] [arguments...]",
	Short:              "Exec command with the Go and Python environment properly configured",
	Long:               "Exec executes a command with the Go and Python environment properly configured.",
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}
		if err := rungo.RunCommand(args[0], args[1:]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(execCmd)
}
