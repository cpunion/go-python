package install

import (
	"fmt"
	"runtime"

	"github.com/cpunion/go-python/internal/env"
)

const (
	// Go download URL format
	goDownloadURL = "https://go.dev/dl/go%s.%s-%s.%s"
)

// getGoURL returns the appropriate Go download URL for the current platform
func getGoURL(version string) string {
	var os, arch, ext string

	switch runtime.GOOS {
	case "windows":
		os = "windows"
		ext = "zip"
	case "darwin":
		os = "darwin"
		ext = "tar.gz"
	case "linux":
		os = "linux"
		ext = "tar.gz"
	default:
		return ""
	}

	switch runtime.GOARCH {
	case "amd64":
		arch = "amd64"
	case "386":
		arch = "386"
	case "arm64":
		arch = "arm64"
	default:
		return ""
	}

	return fmt.Sprintf(goDownloadURL, version, os, arch, ext)
}

// installGo downloads and installs Go in the project directory
func installGo(projectPath, version string, verbose bool) error {
	goDir := env.GetGoDir(projectPath)
	fmt.Printf("Installing Go %s in %s\n", version, goDir)
	// Get download URL
	url := getGoURL(version)
	if url == "" {
		return fmt.Errorf("unsupported platform")
	}

	return downloadAndExtract("Go", version, url, goDir, "", verbose)
}
