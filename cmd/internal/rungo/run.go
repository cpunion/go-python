package rungo

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/cpunion/go-python/internal/env"
)

type ListInfo struct {
	Dir  string `json:"Dir"`
	Root string `json:"Root"`
}

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

func FindProjectRoot(dir string) (string, error) {
	env := env.NewPythonEnv(env.GetPythonRoot(dir))
	_, err := env.Python()
	if err == nil {
		return dir, nil
	}
	parentDir := filepath.Dir(dir)
	if parentDir == dir {
		return "", fmt.Errorf("failed to find Gopy project")
	}
	return FindProjectRoot(parentDir)
}

// RunGoCommand executes a Go command with Python environment properly configured
func RunCommand(command string, args []string) error {
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %v", err)
	}
	projectRoot, err := FindProjectRoot(wd)
	if err != nil {
		return fmt.Errorf("should run this command in a Gopy project: %v", err)
	}
	env.SetBuildEnv(projectRoot)

	// Set up environment variables
	goEnv := []string{}
	// Get PYTHONPATH and PYTHONHOME from env.txt
	var pythonPath, pythonHome string
	if additionalEnv, err := env.ReadEnv(projectRoot); err == nil {
		for key, value := range additionalEnv {
			goEnv = append(goEnv, key+"="+value)
		}
		pythonPath = additionalEnv["PYTHONPATH"]
		pythonHome = additionalEnv["PYTHONHOME"]
	} else {
		fmt.Fprintf(os.Stderr, "Warning: could not load environment variables: %v\n", err)
	}

	cmdArgs := args
	if command == "go" {
		goCmd := args[0]
		args = args[1:]
		// Process args to inject Python paths via ldflags
		cmdArgs = append([]string{goCmd}, ProcessArgsWithLDFlags(args, projectRoot, pythonPath, pythonHome)...)
	}

	cmd := exec.Command(command, cmdArgs...)
	cmd.Env = append(goEnv, os.Environ()...)
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
func ProcessArgsWithLDFlags(args []string, projectRoot, pythonPath, pythonHome string) []string {
	result := make([]string, 0, len(args))

	// Prepare the -X flags we want to add
	ldflags := fmt.Sprintf("-X 'github.com/cpunion/go-python.ProjectRoot=%s'", projectRoot)

	// Prepare rpath flag if needed
	pythonLibDir := env.GetPythonLibDir(projectRoot)
	switch runtime.GOOS {
	case "darwin", "linux":
		ldflags += fmt.Sprintf(" -extldflags '-Wl,-rpath,%s'", pythonLibDir)
	case "windows":
		// Windows doesn't use rpath
	default:
		// Use Linux format for other Unix-like systems
		ldflags += fmt.Sprintf(" -extldflags '-Wl,-rpath=%s'", pythonLibDir)
	}

	for i := 0; i < len(args); i++ {
		arg := args[i]
		if strings.HasPrefix(arg, "-ldflags=") || arg == "-ldflags" {
			// Get existing flags
			var existingFlags string
			if strings.HasPrefix(arg, "-ldflags=") {
				existingFlags = strings.TrimPrefix(arg, "-ldflags=")
			} else if i+1 < len(args) {
				existingFlags = args[i+1]
				i++ // Skip the next arg since we've consumed it
			}
			existingFlags = strings.TrimSpace(existingFlags)
			if ldflags != "" {
				ldflags += " " + existingFlags
			}
		} else {
			result = append(result, arg)
		}
	}
	return append([]string{"-ldflags", ldflags}, result...)
}

// GetGoCommandHelp returns the formatted help text for the specified go command
func GetGoCommandHelp(command string) (string, error) {
	cmd := exec.Command("go", "help", command)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	intro := fmt.Sprintf(`The command arguments and flags are fully compatible with 'go %s'.

Following is the help message from 'go %s':
-------------------------------------------------------------------------------

`, command, command)

	return intro + out.String() + "\n-------------------------------------------------------------------------------", nil
}
