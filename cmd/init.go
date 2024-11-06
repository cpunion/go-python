/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/cpunion/go-python/cmd/internal/create"
	"github.com/cpunion/go-python/cmd/internal/install"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var bold = color.New(color.Bold).SprintFunc()

// isDirEmpty checks if a directory is empty
func isDirEmpty(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err
}

// promptYesNo asks user for confirmation
func promptYesNo(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s [y/N]: ", prompt)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false
	}

	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes"
}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init [path]",
	Short: "Initialize a new go-python project",
	Long: `Initialize a new go-python project in the specified directory.
If no path is provided, it will initialize in the current directory.

Example:
  gopy init
  gopy init my-project
  gopy init --debug my-project
  gopy init -v my-project`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get project path
		projectPath := "."
		if len(args) > 0 {
			projectPath = args[0]
		}

		// Get flags
		debug, _ := cmd.Flags().GetBool("debug")
		verbose, _ := cmd.Flags().GetBool("verbose")
		version, _ := cmd.Flags().GetString("version")
		buildDate, _ := cmd.Flags().GetString("build-date")
		freeThreaded, _ := cmd.Flags().GetBool("free-threaded")

		// Check if directory exists
		if _, err := os.Stat(projectPath); err == nil {
			// Directory exists, check if it's empty
			empty, err := isDirEmpty(projectPath)
			if err != nil {
				fmt.Printf("Error checking directory: %v\n", err)
				return
			}

			if !empty {
				if !promptYesNo(fmt.Sprintf("Directory %s is not empty. Do you want to continue?", projectPath)) {
					fmt.Println("Operation cancelled")
					return
				}
			}
		} else if !os.IsNotExist(err) {
			fmt.Printf("Error checking directory: %v\n", err)
			return
		}

		// Create project using the create package
		fmt.Printf("\n%s\n", bold("Creating project..."))
		if err := create.Project(projectPath, verbose); err != nil {
			fmt.Printf("Error creating project: %v\n", err)
			return
		}

		// Install dependencies
		fmt.Printf("\n%s\n", bold("Installing dependencies..."))
		if err := install.Dependencies(projectPath, version, buildDate, freeThreaded, debug, verbose); err != nil {
			fmt.Printf("Error installing dependencies: %v\n", err)
			return
		}

		fmt.Printf("\n%s\n", bold("Successfully initialized go-python project in "+projectPath))
		fmt.Println("\nNext steps:")
		fmt.Println("1. cd", projectPath)
		fmt.Println("2. gopy run .")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().Bool("debug", false, "Install debug version of Python (not available on Windows)")
	initCmd.Flags().BoolP("verbose", "v", false, "Enable verbose output")
	initCmd.Flags().String("version", "3.13.0", "Python version to install")
	initCmd.Flags().String("build-date", "20241016", "Python build date")
	initCmd.Flags().Bool("free-threaded", false, "Install free-threaded version of Python")
}
