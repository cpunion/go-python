package install

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/cpunion/go-python/cmd/internal/python"
)

const (
	baseURL = "https://github.com/indygreg/python-build-standalone/releases/download/%s"
)

type pythonBuild struct {
	arch     string
	os       string
	variant  string
	debug    bool
	shared   bool
	fullPack bool
}

// getPythonURL returns the appropriate Python standalone URL for the current platform
func getPythonURL(version, buildDate, arch, os string, freeThreaded, debug bool) string {
	// Map GOARCH to Python build architecture
	archMap := map[string]string{
		"amd64": "x86_64",
		"arm64": "aarch64",
		"386":   "i686",
	}

	pythonArch, ok := archMap[arch]
	if !ok {
		return ""
	}

	build := pythonBuild{
		arch:     pythonArch,
		fullPack: true,
		debug:    debug,
	}

	switch os {
	case "darwin":
		build.os = "apple-darwin"
		if freeThreaded {
			build.variant = "freethreaded"
			if build.debug {
				build.variant += "+debug"
			} else {
				build.variant += "+pgo"
			}
		} else {
			if build.debug {
				build.variant = "debug"
			} else {
				build.variant = "pgo"
			}
		}
	case "linux":
		build.os = "unknown-linux-gnu"
		if freeThreaded {
			build.variant = "freethreaded"
			if build.debug {
				build.variant += "+debug"
			} else {
				build.variant += "+pgo"
			}
		} else {
			if build.debug {
				build.variant = "debug"
			} else {
				build.variant = "pgo"
			}
		}
	case "windows":
		build.os = "pc-windows-msvc"
		build.shared = true
		if freeThreaded {
			build.variant = "freethreaded+pgo"
		} else {
			build.variant = "pgo"
		}
	default:
		return ""
	}

	// Construct filename
	filename := fmt.Sprintf("cpython-%s+%s-%s-%s", version, buildDate, build.arch, build.os)
	if build.shared {
		filename += "-shared"
	}
	filename += "-" + build.variant
	if build.fullPack {
		filename += "-full"
	}
	filename += ".tar.zst"

	return fmt.Sprintf(baseURL, buildDate) + "/" + filename
}

// updateMacOSDylibs updates the install names of dylib files on macOS
func updateMacOSDylibs(pythonDir string, verbose bool) error {
	if runtime.GOOS != "darwin" {
		return nil
	}

	libDir := filepath.Join(pythonDir, "lib")
	entries, err := os.ReadDir(libDir)
	if err != nil {
		return fmt.Errorf("failed to read lib directory: %v", err)
	}

	absLibDir, err := filepath.Abs(libDir)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %v", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".dylib") {
			dylibPath := filepath.Join(libDir, entry.Name())
			if verbose {
				fmt.Printf("Updating install name for: %s\n", dylibPath)
			}

			// Get the current install name
			cmd := exec.Command("otool", "-D", dylibPath)
			output, err := cmd.Output()
			if err != nil {
				return fmt.Errorf("failed to get install name for %s: %v", dylibPath, err)
			}

			// Parse the output to get the current install name
			lines := strings.Split(string(output), "\n")
			if len(lines) < 2 {
				continue
			}
			currentName := strings.TrimSpace(lines[1])
			if currentName == "" {
				continue
			}

			// Calculate new install name using absolute path
			newName := filepath.Join(absLibDir, filepath.Base(currentName))

			fmt.Printf("Updating install name for %s to %s\n", dylibPath, newName)
			// Update the install name
			cmd = exec.Command("install_name_tool", "-id", newName, dylibPath)
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("failed to update install name for %s: %v", dylibPath, err)
			}
		}
	}
	return nil
}

// generatePkgConfig generates pkg-config files for Windows
func generatePkgConfig(pythonPath, pkgConfigDir string) error {
	if err := os.MkdirAll(pkgConfigDir, 0755); err != nil {
		return fmt.Errorf("failed to create pkgconfig directory: %v", err)
	}

	// Get Python version from the environment
	env := python.New(pythonPath)
	pythonBin, err := env.Python()
	if err != nil {
		return fmt.Errorf("failed to get Python executable: %v", err)
	}

	// Get Python version
	cmd := exec.Command(pythonBin, "-c", "import sys; print(f'{sys.version_info.major}.{sys.version_info.minor}')")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to get Python version: %v", err)
	}
	version := strings.TrimSpace(string(output))

	// Template for the pkg-config file
	pcTemplate := `prefix=${pcfiledir}/../..
exec_prefix=${prefix}
libdir=${exec_prefix}
includedir=${prefix}/include

Name: Python
Description: Embed Python into an application
Requires:
Version: %s
Libs.private: 
Libs: -L${libdir} -lpython313
Cflags: -I${includedir}
`
	// TODO: need update libs

	// Create the main pkg-config files
	files := []struct {
		name    string
		content string
	}{
		{
			fmt.Sprintf("python-%s.pc", version),
			fmt.Sprintf(pcTemplate, version),
		},
		{
			fmt.Sprintf("python-%s-embed.pc", version),
			fmt.Sprintf(pcTemplate, version),
		},
		{
			"python3.pc",
			fmt.Sprintf(pcTemplate, version),
		},
		{
			"python3-embed.pc",
			fmt.Sprintf(pcTemplate, version),
		},
	}

	// Write all pkg-config files
	for _, file := range files {
		pcPath := filepath.Join(pkgConfigDir, file.name)
		if err := os.WriteFile(pcPath, []byte(file.content), 0644); err != nil {
			return fmt.Errorf("failed to write %s: %v", file.name, err)
		}
	}

	return nil
}

// updatePkgConfig updates the prefix in pkg-config files to use absolute path
func updatePkgConfig(projectPath string) error {
	pythonPath := GetPythonRoot(projectPath)
	pkgConfigDir := GetPythonPkgConfigDir(projectPath)

	if runtime.GOOS == "windows" {
		if err := generatePkgConfig(pythonPath, pkgConfigDir); err != nil {
			return err
		}
	}

	entries, err := os.ReadDir(pkgConfigDir)
	if err != nil {
		return fmt.Errorf("failed to read pkgconfig directory: %v", err)
	}

	absPath, err := filepath.Abs(pythonPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %v", err)
	}

	// Helper function to write a .pc file with the correct prefix
	writePC := func(path string, content []byte) error {
		newContent := strings.ReplaceAll(string(content), "prefix=/install", "prefix="+absPath)
		return os.WriteFile(path, []byte(newContent), 0644)
	}

	// Regular expressions for matching file patterns
	normalPattern := regexp.MustCompile(`^python-(\d+\.\d+)t?\.pc$`)
	embedPattern := regexp.MustCompile(`^python-(\d+\.\d+)t?-embed\.pc$`)

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".pc") {
			pcFile := filepath.Join(pkgConfigDir, entry.Name())

			// Read file content
			content, err := os.ReadFile(pcFile)
			if err != nil {
				return fmt.Errorf("failed to read %s: %v", pcFile, err)
			}

			// Update original file
			if err := writePC(pcFile, content); err != nil {
				return fmt.Errorf("failed to update %s: %v", pcFile, err)
			}

			name := entry.Name()
			// Create additional copies based on patterns
			copies := make(map[string]bool)

			// Handle python-X.YZt.pc and python-X.YZ.pc patterns
			if matches := normalPattern.FindStringSubmatch(name); matches != nil {
				if strings.Contains(name, "t.pc") {
					// python-3.13t.pc -> python3.pc and python3t.pc
					copies["python3t.pc"] = true
					copies["python3.pc"] = true
					// Also create non-t version
					noT := fmt.Sprintf("python-%s.pc", matches[1])
					if err := writePC(filepath.Join(pkgConfigDir, noT), content); err != nil {
						return fmt.Errorf("failed to write %s: %v", noT, err)
					}
				} else {
					// python-3.13.pc -> python3.pc
					copies["python3.pc"] = true
				}
			}

			// Handle python-X.YZt-embed.pc and python-X.YZ-embed.pc patterns
			if matches := embedPattern.FindStringSubmatch(name); matches != nil {
				if strings.Contains(name, "t-embed.pc") {
					// python-3.13t-embed.pc -> python3-embed.pc and python3t-embed.pc
					copies["python3t-embed.pc"] = true
					copies["python3-embed.pc"] = true
					// Also create non-t version
					noT := fmt.Sprintf("python-%s-embed.pc", matches[1])
					if err := writePC(filepath.Join(pkgConfigDir, noT), content); err != nil {
						return fmt.Errorf("failed to write %s: %v", noT, err)
					}
				} else {
					// python-3.13-embed.pc -> python3-embed.pc
					copies["python3-embed.pc"] = true
				}
			}

			// Write all unique copies
			for copyName := range copies {
				copyPath := filepath.Join(pkgConfigDir, copyName)
				if err := writePC(copyPath, content); err != nil {
					return fmt.Errorf("failed to write %s: %v", copyPath, err)
				}
			}
		}
	}
	return nil
}

// writeEnvFile writes environment variables to .python/env.txt
func writeEnvFile(projectPath string) error {
	pythonDir := GetPythonRoot(projectPath)
	absPath, err := filepath.Abs(pythonDir)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %v", err)
	}

	// Get Python path using python executable
	env := python.New(absPath)
	pythonBin, err := env.Python()
	if err != nil {
		return fmt.Errorf("failed to get Python executable: %v", err)
	}

	// Execute Python to get sys.path
	cmd := exec.Command(pythonBin, "-c", "import sys; print(':'.join(sys.path))")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to get Python path: %v", err)
	}

	// Prepare environment variables
	envVars := []string{
		fmt.Sprintf("PKG_CONFIG_PATH=%s", filepath.Join(absPath, "lib", "pkgconfig")),
		fmt.Sprintf("PYTHONPATH=%s", strings.TrimSpace(string(output))),
		fmt.Sprintf("PYTHONHOME=%s", absPath),
	}

	// Write to env.txt
	envFile := filepath.Join(pythonDir, "env.txt")
	if err := os.WriteFile(envFile, []byte(strings.Join(envVars, "\n")), 0644); err != nil {
		return fmt.Errorf("failed to write env file: %v", err)
	}

	return nil
}

// LoadEnvFile loads environment variables from .python/env.txt in the given directory
func LoadEnvFile(dir string) ([]string, error) {
	envFile := filepath.Join(GetPythonRoot(dir), "env.txt")
	content, err := os.ReadFile(envFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read env file: %v", err)
	}

	return strings.Split(strings.TrimSpace(string(content)), "\n"), nil
}

// installPythonEnv downloads and installs Python standalone build
func installPythonEnv(projectPath string, version, buildDate string, freeThreaded, debug bool, verbose bool) error {
	fmt.Printf("Installing Python %s in %s\n", version, projectPath)
	pythonDir := GetPythonRoot(projectPath)

	// Remove existing Python directory if it exists
	if err := os.RemoveAll(pythonDir); err != nil {
		return fmt.Errorf("error removing existing Python directory: %v", err)
	}

	// Get Python URL
	url := getPythonURL(version, buildDate, runtime.GOARCH, runtime.GOOS, freeThreaded, debug)
	if url == "" {
		return fmt.Errorf("unsupported platform")
	}

	if err := downloadAndExtract("Python", version, url, pythonDir, "python/install", verbose); err != nil {
		return fmt.Errorf("error downloading and extracting Python: %v", err)
	}

	// After extraction, update dylib install names on macOS
	if err := updateMacOSDylibs(pythonDir, verbose); err != nil {
		return fmt.Errorf("error updating dylib install names: %v", err)
	}

	// Create Python environment
	env := python.New(pythonDir)

	if verbose {
		fmt.Println("Installing Python dependencies...")
	}

	if err := env.RunPip("install", "--upgrade", "pip", "setuptools", "wheel"); err != nil {
		return fmt.Errorf("error upgrading pip, setuptools, whell")
	}

	if err := updatePkgConfig(projectPath); err != nil {
		return fmt.Errorf("error updating pkg-config: %v", err)
	}

	// Write environment variables to env.txt
	if err := writeEnvFile(projectPath); err != nil {
		return fmt.Errorf("error writing environment file: %v", err)
	}

	return nil
}
