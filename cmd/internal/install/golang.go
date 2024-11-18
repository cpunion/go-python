package install

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
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
	goDir := filepath.Join(projectPath, DepsDir, GoDir)
	fmt.Printf("Installing Go %s in %s\n", version, goDir)

	// Create Go directory if it doesn't exist
	if err := os.MkdirAll(goDir, 0755); err != nil {
		return fmt.Errorf("error creating go directory: %v", err)
	}

	// Get download URL
	url := getGoURL(version)
	if url == "" {
		return fmt.Errorf("unsupported platform")
	}

	if verbose {
		fmt.Printf("Downloading Go %s from %s\n", version, url)
	}

	path, err := downloadFileWithCache(url)
	if err != nil {
		return fmt.Errorf("error downloading Go: %v", err)
	}

	if verbose {
		fmt.Println("Extracting Go...")
	}

	// Extract based on file extension
	if strings.HasSuffix(path, ".zip") {
		if err := extractZip(path, goDir); err != nil {
			return fmt.Errorf("error extracting Go: %v", err)
		}
	} else if strings.HasSuffix(path, ".tar.gz") {
		if err := extractTarGz(path, goDir); err != nil {
			return fmt.Errorf("error extracting Go: %v", err)
		}
	}

	return nil
}

// extractZip extracts a zip file to the specified directory
func extractZip(zipFile, destDir string) error {
	r, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		// Skip the root "go" directory
		if f.Name == "go/" || f.Name == "go" {
			continue
		}

		// Remove "go/" prefix from paths
		destPath := filepath.Join(destDir, strings.TrimPrefix(f.Name, "go/"))

		if f.FileInfo().IsDir() {
			os.MkdirAll(destPath, f.Mode())
			continue
		}

		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			return err
		}

		destFile, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		srcFile, err := f.Open()
		if err != nil {
			destFile.Close()
			return err
		}

		_, err = io.Copy(destFile, srcFile)
		srcFile.Close()
		destFile.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

// extractTarGz extracts a tar.gz file to the specified directory
func extractTarGz(tarFile, destDir string) error {
	file, err := os.Open(tarFile)
	if err != nil {
		return err
	}
	defer file.Close()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// Skip the root "go" directory
		if header.Name == "go/" || header.Name == "go" {
			continue
		}

		// Remove "go/" prefix from paths
		destPath := filepath.Join(destDir, strings.TrimPrefix(header.Name, "go/"))

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(destPath, os.FileMode(header.Mode)); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
				return err
			}
			outFile, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			if _, err := io.Copy(outFile, tr); err != nil {
				outFile.Close()
				return err
			}
			outFile.Close()
		}
	}
	return nil
}
