package install

import (
	"fmt"
	"os"
	"strings"
)

const (
	msys2Dir   = "msys2"
	releaseTag = "2024-11-16"
)

func installMsys2(projectPath string, verbose bool) error {
	msys2Root := GetMsys2Dir(projectPath)
	fmt.Printf("Installing msys2 in %v\n", msys2Root)

	msys2Version := strings.ReplaceAll(releaseTag, "-", "")
	msys2URL := fmt.Sprintf("https://github.com/msys2/msys2-installer/releases/download/%s/msys2-base-x86_64-%s.tar.zst", releaseTag, msys2Version)

	if verbose {
		fmt.Printf("Downloading MSYS2 from %s\n", msys2URL)
	}

	path, err := downloadFileWithCache(msys2URL)
	if err != nil {
		return fmt.Errorf("error downloading MSYS2: %v", err)
	}

	if verbose {
		fmt.Println("Extracting MSYS2...")
	}

	// Create MSYS2 directory
	if err := os.MkdirAll(msys2Root, 0755); err != nil {
		return fmt.Errorf("error creating MSYS2 directory: %v", err)
	}

	// Extract archive
	if err := extractTarZst(path, msys2Root, "", verbose); err != nil {
		return fmt.Errorf("error extracting MSYS2: %v", err)
	}

	if verbose {
		fmt.Println("MSYS2 installation completed")
	}

	return nil
}
