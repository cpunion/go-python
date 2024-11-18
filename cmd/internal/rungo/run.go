package rungo

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/cpunion/go-python/cmd/internal/install"
)

// FindPackageIndex finds the package argument index by skipping flags and their values
func FindPackageIndex(args []string) int {
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
		return i
	}
	return -1
}

// GetPackageDir returns the directory containing the package
func GetPackageDir(pkgPath string) (string, error) {
	// Get the absolute path
	absPath, err := filepath.Abs(pkgPath)
	if err != nil {
		return "", fmt.Errorf("error resolving path: %v", err)
	}

	// If it's not a directory, get its parent directory
	fi, err := os.Stat(absPath)
	if err != nil {
		if os.IsNotExist(err) && pkgPath == "." {
			// Special case: if "." doesn't exist, use current directory
			dir, err := os.Getwd()
			if err != nil {
				return "", fmt.Errorf("error getting working directory: %v", err)
			}
			absPath = dir
			fi, err = os.Stat(absPath)
			if err != nil {
				return "", fmt.Errorf("error checking path: %v", err)
			}
		} else {
			return "", fmt.Errorf("error checking path: %v", err)
		}
	}

	if !fi.IsDir() {
		return filepath.Dir(absPath), nil
	}
	return absPath, nil
}

// RunGoCommand executes a Go command with Python environment properly configured
func RunGoCommand(command string, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("package argument is required")
	}

	// Find the package argument
	pkgIndex := FindPackageIndex(args)
	if pkgIndex == -1 {
		return fmt.Errorf("package argument is required")
	}

	// Get the package path
	pkgPath := args[pkgIndex]

	// Get package directory
	dir, err := GetPackageDir(pkgPath)
	if err != nil {
		return err
	}

	// Set up environment variables
	env := os.Environ()

	// Load additional environment variables from env.txt
	if additionalEnv, err := install.LoadEnvFile(dir); err == nil {
		env = append(env, additionalEnv...)
	} else {
		fmt.Fprintf(os.Stderr, "Warning: could not load environment variables: %v\n", err)
	}

	// Get PYTHONPATH and PYTHONHOME from environment
	pythonPath := os.Getenv("PYTHONPATH")
	pythonHome := os.Getenv("PYTHONHOME")

	// Process args to inject Python paths via ldflags
	processedArgs := ProcessArgsWithLDFlags(args, pythonPath, pythonHome)

	// Prepare go command with processed arguments
	goArgs := append([]string{command}, processedArgs...)
	cmd := exec.Command("go", goArgs...)
	cmd.Dir = dir
	cmd.Env = env
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if command == "run" {
		cmd.Stdin = os.Stdin
	}

	// Execute the command
	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		return fmt.Errorf("error executing command: %v", err)
	}

	return nil
}

// ProcessArgsWithLDFlags processes command line arguments to inject Python paths via ldflags
func ProcessArgsWithLDFlags(args []string, pythonPath, pythonHome string) []string {
	result := make([]string, 0, len(args)+4) // reserve space for potential new flags
	result = append(result, args...)

	if pythonPath != "" {
		result = append(result, "-ldflags", fmt.Sprintf("-X 'github.com/cpunion/go-python.PythonPath=%s'", pythonPath))
	}
	if pythonHome != "" {
		result = append(result, "-ldflags", fmt.Sprintf("-X 'github.com/cpunion/go-python.PythonHome=%s'", pythonHome))
	}

	return result
}
