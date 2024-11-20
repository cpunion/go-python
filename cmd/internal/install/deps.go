package install

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/cpunion/go-python/internal/env"
)

// Dependencies installs all required dependencies for the project
func Dependencies(projectPath string, goVersion, tinyPkgConfigVersion, pyVersion, pyBuildDate string, freeThreaded, debug bool, verbose bool) error {
	if err := installTinyPkgConfig(projectPath, tinyPkgConfigVersion, verbose); err != nil {
		return err
	}
	// Only install MSYS2 on Windows
	if runtime.GOOS == "windows" {
		if err := installMingw(projectPath, verbose); err != nil {
			return err
		}
	}

	if err := installGo(projectPath, goVersion, verbose); err != nil {
		return err
	}
	env.SetBuildEnv(projectPath)

	// Install Go dependencies
	if err := installGoDeps(projectPath); err != nil {
		return err
	}

	// Install Python environment and dependencies
	if err := installPythonEnv(projectPath, pyVersion, pyBuildDate, freeThreaded, debug, verbose); err != nil {
		return err
	}

	if runtime.GOOS == "windows" {
		pythonPath := env.GetPythonRoot(projectPath)
		pkgConfigDir := env.GetPythonPkgConfigDir(projectPath)
		if err := generatePkgConfig(pythonPath, pkgConfigDir); err != nil {
			return err
		}
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
