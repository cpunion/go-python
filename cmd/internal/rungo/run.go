package rungo

import (
	"bytes"
	"encoding/json"
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

// RunGoCommand executes a Go command with Python environment properly configured
func RunGoCommand(command string, args []string) error {
	// Find the package argument
	pkgIndex := FindPackageIndex(args)

	// TODO: don't depend on external go command
	listArgs := []string{"list", "-find", "-json"}

	if pkgIndex != -1 {
		pkgPath := args[pkgIndex]
		listArgs = append(listArgs, pkgPath)
	}
	cmd := exec.Command("go", listArgs...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to get module info: %v", err)
	}
	var listInfo ListInfo
	if err := json.NewDecoder(&out).Decode(&listInfo); err != nil {
		return fmt.Errorf("failed to parse module info: %v", err)
	}
	projectRoot := listInfo.Root
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

	// Process args to inject Python paths via ldflags
	processedArgs := ProcessArgsWithLDFlags(args, projectRoot, pythonPath, pythonHome)

	// Prepare go command with processed arguments
	goArgs := append([]string{"go", command}, processedArgs...)
	cmd = exec.Command(goArgs[0], goArgs[1:]...)
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
	xFlags := fmt.Sprintf("-X 'github.com/cpunion/go-python.ProjectRoot=%s'", projectRoot)

	// Prepare rpath flag if needed
	var rpathFlag string
	if pythonHome != "" {
		pythonLibDir := filepath.Join(pythonHome, "lib")
		switch runtime.GOOS {
		case "darwin", "linux":
			rpathFlag = fmt.Sprintf("-extldflags '-Wl,-rpath,%s'", pythonLibDir)
		case "windows":
			// Windows doesn't use rpath
			rpathFlag = ""
		default:
			// Use Linux format for other Unix-like systems
			rpathFlag = fmt.Sprintf("-extldflags '-Wl,-rpath=%s'", pythonLibDir)
		}
	}

	// Find existing -ldflags if any
	foundLDFlags := false
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if strings.HasPrefix(arg, "-ldflags=") || arg == "-ldflags" {
			foundLDFlags = true
			// Copy everything before this arg
			result = append(result, args[:i]...)

			// Get existing flags
			var existingFlags string
			if strings.HasPrefix(arg, "-ldflags=") {
				existingFlags = strings.TrimPrefix(arg, "-ldflags=")
			} else if i+1 < len(args) {
				existingFlags = args[i+1]
				i++ // Skip the next arg since we've consumed it
			}

			// Combine all flags
			allFlags := []string{xFlags}
			if strings.TrimSpace(existingFlags) != "" {
				allFlags = append(allFlags, existingFlags)
			}
			if rpathFlag != "" {
				allFlags = append(allFlags, rpathFlag)
			}

			// Add combined ldflags
			result = append(result, "-ldflags")
			result = append(result, strings.Join(allFlags, " "))

			// Add remaining args
			result = append(result, args[i+1:]...)
			break
		}
	}

	// If no existing -ldflags found, add new ones at the beginning if we have any flags to add
	if !foundLDFlags {
		if len(xFlags) > 0 || rpathFlag != "" {
			allFlags := []string{xFlags}
			if rpathFlag != "" {
				allFlags = append(allFlags, rpathFlag)
			}
			result = append(result, "-ldflags")
			result = append(result, strings.Join(allFlags, " "))
		}
		result = append(result, args...)
	}

	return result
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
