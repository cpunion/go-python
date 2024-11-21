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
	Long: func() string {
		intro := "Install compiles and installs Go packages with the Python environment properly configured.\n\n"
		help, err := rungo.GetGoCommandHelp("install")
		if err != nil {
			return intro + "Failed to get go help: " + err.Error()
		}
		return intro + help
	}(),
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		if err := rungo.RunCommand("go", append([]string{"install"}, args...)); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
