package install

import (
	"archive/tar"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/cpunion/go-python/cmd/internal/python"
)

const (
	baseURL       = "https://github.com/indygreg/python-build-standalone/releases/download/%s"
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
			if debug {
				build.variant += "+debug"
			} else {
				build.variant += "+pgo"
			}
		} else {
			if debug {
				build.variant = "debug"
			} else {
				build.variant = "pgo"
			}
		}
	case "linux":
		build.os = "unknown-linux-gnu"
		if freeThreaded {
			build.variant = "freethreaded"
			if debug {
				build.variant += "+debug"
			} else {
				build.variant += "+pgo"
			}
		} else {
			if debug {
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

// getCacheDir returns the cache directory for downloaded files
func getCacheDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %v", err)
	}
	cacheDir := filepath.Join(homeDir, ".gopy", "cache")
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create cache directory: %v", err)
	}
	return cacheDir, nil
}

// downloadFileWithCache downloads a file from url and returns the path to the cached file
func downloadFileWithCache(url string) (string, error) {
	cacheDir, err := getCacheDir()
	if err != nil {
		return "", err
	}

	// Use URL's last path segment as filename
	urlPath := strings.Split(url, "/")
	filename := urlPath[len(urlPath)-1]
	cachedFile := filepath.Join(cacheDir, filename)

	// Check if file exists in cache
	if _, err := os.Stat(cachedFile); err == nil {
		fmt.Printf("Using cached Python from %s\n", cachedFile)
		return cachedFile, nil
	}

	fmt.Printf("Downloading Python from %s\n", url)

	// Create temporary file
	tmpFile, err := os.CreateTemp(cacheDir, "download-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary file: %v", err)
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)
	defer tmpFile.Close()

	// Download to temporary file
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status: %s", resp.Status)
	}

	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to write file: %v", err)
	}

	// Close the file before renaming
	tmpFile.Close()

	// Rename temporary file to cached file
	if err := os.Rename(tmpPath, cachedFile); err != nil {
		return "", fmt.Errorf("failed to move file to cache: %v", err)
	}

	return cachedFile, nil
}

// findPipExecutable finds the correct pip executable in the bin directory
func findPipExecutable(binDir string) (string, error) {
	if runtime.GOOS == "windows" {
		// Check for pip3.exe
		pipPath := filepath.Join(binDir, "Scripts", "pip3.exe")
		if _, err := os.Stat(pipPath); err == nil {
			return pipPath, nil
		}
		// Check for pip.exe
		pipPath = filepath.Join(binDir, "Scripts", "pip.exe")
		if _, err := os.Stat(pipPath); err == nil {
			return pipPath, nil
		}
	} else {
		// Try different pip names
		pipNames := []string{"pip3", "pip"}
		for _, name := range pipNames {
			pipPath := filepath.Join(binDir, "bin", name)
			if info, err := os.Stat(pipPath); err == nil {
				// Check if the file is executable
				if info.Mode()&0111 != 0 {
					return pipPath, nil
				}
			}
		}
	}
	return "", fmt.Errorf("pip executable not found in %s", binDir)
}

// extractTarZst extracts a tar.zst file to a destination directory
func extractTarZst(src, dst string, verbose bool) error {
	if verbose {
		fmt.Printf("Extracting from %s to %s\n", src, dst)
	}

	// First decompress with zstd
	tarFile := src + ".tar"
	cmd := exec.Command("zstd", "-d", src, "-o", tarFile)
	if verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error decompressing with zstd: %v", err)
	}
	defer os.Remove(tarFile)

	// Then extract tar
	file, err := os.Open(tarFile)
	if err != nil {
		return err
	}
	defer file.Close()

	tr := tar.NewReader(file)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// Only extract files from the install directory
		if !strings.HasPrefix(header.Name, "python/install/") {
			continue
		}

		// Remove the "python/install/" prefix
		name := strings.TrimPrefix(header.Name, "python/install/")
		if name == "" {
			continue
		}

		path := filepath.Join(dst, name)
		if verbose {
			fmt.Printf("Extracting: %s\n", path)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(path, 0755); err != nil {
				return fmt.Errorf("error creating directory %s: %v", path, err)
			}
		case tar.TypeReg:
			dir := filepath.Dir(path)
			if err := os.MkdirAll(dir, 0755); err != nil {
				return fmt.Errorf("error creating directory %s: %v", dir, err)
			}

			file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return fmt.Errorf("error creating file %s: %v", path, err)
			}

			if _, err := io.Copy(file, tr); err != nil {
				file.Close()
				return fmt.Errorf("error writing to file %s: %v", path, err)
			}
			file.Close()
		}
	}

	// Don't verify pip here anymore
	return nil
}

// updatePkgConfig updates the prefix in pkg-config files to use absolute path
func updatePkgConfig(projectPath string) error {
	pkgConfigDir := filepath.Join(projectPath, ".python/lib/pkgconfig")
	entries, err := os.ReadDir(pkgConfigDir)
	if err != nil {
		return fmt.Errorf("failed to read pkgconfig directory: %v", err)
	}

	pythonPath := filepath.Join(projectPath, ".python")
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

// installPythonEnv downloads and installs Python standalone build
func installPythonEnv(projectPath string, version, buildDate string, freeThreaded, debug bool, verbose bool) error {
	pythonDir := filepath.Join(projectPath, ".python")

	// Remove existing Python directory if it exists
	if err := os.RemoveAll(pythonDir); err != nil {
		return fmt.Errorf("error removing existing Python directory: %v", err)
	}

	// Get Python URL
	url := getPythonURL(version, buildDate, runtime.GOARCH, runtime.GOOS, freeThreaded, debug)
	if url == "" {
		return fmt.Errorf("unsupported platform")
	}

	// Download Python
	archivePath, err := downloadFileWithCache(url)
	if err != nil {
		return fmt.Errorf("error downloading Python: %v", err)
	}

	if err := os.MkdirAll(pythonDir, 0755); err != nil {
		return fmt.Errorf("error creating python directory: %v", err)
	}

	if verbose {
		fmt.Println("Extracting Python...")
	}
	// Extract to .python directory
	if err := extractTarZst(archivePath, pythonDir, verbose); err != nil {
		return fmt.Errorf("error extracting Python: %v", err)
	}

	// Create Python environment
	env := python.New(projectPath)

	// Make sure pip is executable
	pipPath, err := env.Pip()
	if err != nil {
		return fmt.Errorf("error finding pip: %v", err)
	}

	if runtime.GOOS != "windows" {
		if err := os.Chmod(pipPath, 0755); err != nil {
			return fmt.Errorf("error making pip executable: %v", err)
		}
	}

	if verbose {
		fmt.Printf("Using pip at: %s\n", pipPath)
		fmt.Println("Installing Python dependencies...")
	}

	if err := env.RunPip("install", "--upgrade", "pip", "setuptools", "wheel"); err != nil {
		return fmt.Errorf("error upgrading pip: %v", err)
	}

	if err := updatePkgConfig(projectPath); err != nil {
		return fmt.Errorf("error updating pkg-config: %v", err)
	}

	return nil
}
