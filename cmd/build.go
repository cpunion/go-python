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
		if len(args) < 1 {
			fmt.Println("Error: package argument is required")
			os.Exit(1)
		}

		// Find the package argument by skipping flags and their values
		pkgIndex := -1
		for i := 0; i < len(args); i++ {
			arg := args[i]
			if strings.HasPrefix(arg, "-") {
				// Skip known flags that take values
				switch arg {
				case "-o", "-p", "-asmflags", "-buildmode", "-compiler", "-gccgoflags", "-gcflags",
					"-installsuffix", "-ldflags", "-mod", "-modfile", "-pkgdir", "-tags", "-toolexec":
					i++ // Skip the next argument as it's the flag's value
				}
				continue
			}
			pkgIndex = i
			break
		}

		if pkgIndex == -1 {
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
			if os.IsNotExist(err) && pkgPath == "." {
				// Special case: if "." doesn't exist, use current directory
				dir, err := os.Getwd()
				if err != nil {
					fmt.Printf("Error getting working directory: %v\n", err)
					os.Exit(1)
				}
				absPath = dir
				fi, err = os.Stat(absPath)
				if err != nil {
					fmt.Printf("Error checking path: %v\n", err)
					os.Exit(1)
				}
			} else {
				fmt.Printf("Error checking path: %v\n", err)
				os.Exit(1)
			}
		}

		// Rest of the function remains the same...
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

		// Prepare go build command with all arguments
		goArgs := append([]string{"build"}, args...)
		build := exec.Command("go", goArgs...)
		build.Env = env
		build.Stdout = os.Stdout
		build.Stderr = os.Stderr

		// Execute the command
		if err := build.Run(); err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				os.Exit(exitErr.ExitCode())
			}
			fmt.Printf("Error building command: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
}
