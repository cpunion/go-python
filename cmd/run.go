/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/cpunion/go-python/cmd/internal/install"
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
		if len(args) < 1 {
			fmt.Println("Error: package argument is required")
			os.Exit(1)
		}

		// Find the package argument by skipping flags
		pkgIndex := 0
		for i, arg := range args {
			if !strings.HasPrefix(arg, "-") {
				pkgIndex = i
				break
			}
		}

		if pkgIndex >= len(args) {
			fmt.Println("Error: package argument is required")
			os.Exit(1)
		}

		// Get the package path
		pkgPath := args[pkgIndex]

		// Get the absolute path
		absPath, err := filepath.Abs(pkgPath)
		if err != nil {
			fmt.Printf("Error resolving path: %v\n", err)
			os.Exit(1)
		}

		// If it's not a directory, get its parent directory
		fi, err := os.Stat(absPath)
		if err != nil {
			fmt.Printf("Error checking path: %v\n", err)
			os.Exit(1)
		}

		var dir string
		if !fi.IsDir() {
			dir = filepath.Dir(absPath)
		} else {
			dir = absPath
		}

		// Set up environment variables
		env := os.Environ()

		// Load additional environment variables from env.txt
		if additionalEnv, err := install.LoadEnvFile(dir); err == nil {
			env = append(env, additionalEnv...)
		} else {
			fmt.Fprintf(os.Stderr, "Warning: could not load environment variables: %v\n", err)
		}

		// Prepare go run command with all arguments
		goArgs := append([]string{"run"}, args...)
		run := exec.Command("go", goArgs...)
		run.Env = env
		run.Stdout = os.Stdout
		run.Stderr = os.Stderr
		run.Stdin = os.Stdin

		// Execute the command
		if err := run.Run(); err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				os.Exit(exitErr.ExitCode())
			}
			fmt.Printf("Error running command: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
