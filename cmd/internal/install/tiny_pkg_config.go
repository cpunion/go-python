package install

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	tinyPkgDownloadURL = "https://github.com/cpunion/tiny-pkg-config/releases/download/%s/%s"
)

func installTinyPkgConfig(projectPath, version string, verbose bool) error {
	dir := GetTinyPkgConfigDir(projectPath)
	// Determine OS and architecture
	goos := runtime.GOOS
	arch := runtime.GOARCH

	// Convert OS/arch to match release file naming
	osName := strings.ToUpper(goos[:1]) + goos[1:] // "darwin" -> "Darwin", "linux" -> "Linux"
	archName := arch
	if arch == "amd64" {
		archName = "x86_64"
	}

	// Construct filename and URL
	ext := ".tar.gz"
	if osName == "Windows" {
		ext = ".zip"
	}

	filename := fmt.Sprintf("tiny-pkg-config_%s_%s%s", osName, archName, ext)
	downloadURL := fmt.Sprintf(tinyPkgDownloadURL, version, filename)

	if err := downloadAndExtract("tiny-pkg-config", version, downloadURL, dir, "", verbose); err != nil {
		return fmt.Errorf("download and extract tiny-pkg-config failed: %w", err)
	}

	// After extraction, rename the executable
	oldName := "tiny-pkg-config"
	newName := "pkg-config"
	if runtime.GOOS == "windows" {
		oldName += ".exe"
		newName += ".exe"
	}

	oldPath := filepath.Join(dir, oldName)
	newPath := filepath.Join(dir, newName)

	// Rename the file
	if err := os.Rename(oldPath, newPath); err != nil {
		return fmt.Errorf("failed to rename executable: %w", err)
	}

	if verbose {
		fmt.Printf("Renamed %s to %s\n", oldName, newName)
	}

	return nil
}
