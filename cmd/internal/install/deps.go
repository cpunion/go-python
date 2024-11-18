package install

import (
	"fmt"
	"os"
	"os/exec"
)

// Dependencies installs all required dependencies for the project
func Dependencies(projectPath string, goVersion, pyVersion, pyBuildDate string, freeThreaded, debug bool, verbose bool) error {
	if err := installGo(projectPath, goVersion, verbose); err != nil {
		return err
	}
	SetEnv(projectPath)

	// Install Go dependencies
	if err := installGoDeps(projectPath); err != nil {
		return err
	}

	// Install Python environment and dependencies
	if err := installPythonEnv(projectPath, pyVersion, pyBuildDate, freeThreaded, debug, verbose); err != nil {
		return err
	}

	// Update pkg-config files
	if err := updatePkgConfig(projectPath); err != nil {
		return err
	}

	return nil
}

// installGoDeps installs Go dependencies
func installGoDeps(projectPath string) error {
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %v", err)
	}

	if err := os.Chdir(projectPath); err != nil {
		return fmt.Errorf("error changing to project directory: %v", err)
	}
	defer func() {
		_ = os.Chdir(currentDir)
	}()

	fmt.Println("Installing Go dependencies...")
	getCmd := exec.Command("go", "get", "-u", "github.com/cpunion/go-python")
	getCmd.Stdout = os.Stdout
	getCmd.Stderr = os.Stderr
	if err := getCmd.Run(); err != nil {
		return fmt.Errorf("error installing dependencies: %v", err)
	}

	return nil
}
