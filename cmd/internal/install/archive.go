package install

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/klauspost/compress/zstd"
)

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

// getFullExtension returns the full extension for a filename (e.g., ".tar.gz" for "file.tar.gz")
func getFullExtension(filename string) string {
	// Handle common multi-level extensions
	for _, ext := range []string{".tar.gz", ".tar.zst"} {
		if strings.HasSuffix(filename, ext) {
			return ext
		}
	}
	return filepath.Ext(filename)
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

	// Calculate SHA1 hash of the URL
	hasher := sha1.New()
	hasher.Write([]byte(url))
	urlHash := hex.EncodeToString(hasher.Sum(nil))[:8] // Use first 8 characters of hash

	// Insert hash before the file extension, handling multi-level extensions
	ext := getFullExtension(filename)
	baseFilename := filename[:len(filename)-len(ext)]
	cachedFilename := fmt.Sprintf("%s-%s%s", baseFilename, urlHash, ext)
	cachedFile := filepath.Join(cacheDir, cachedFilename)

	// Check if file exists in cache
	if _, err := os.Stat(cachedFile); err == nil {
		fmt.Printf("Using cached file from %s\n", cachedFile)
		return cachedFile, nil
	}

	fmt.Printf("Downloading from %s\n", url)

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

func downloadAndExtract(name, version, url, dir, trimPrefix string, verbose bool) error {
	if verbose {
		fmt.Printf("Downloading %s %s from %s\n", name, version, url)
	}

	path, err := downloadFileWithCache(url)
	if err != nil {
		return fmt.Errorf("error downloading %s %s: %v", name, version, err)
	}

	if verbose {
		fmt.Printf("Extracting %s %s into %s...\n", name, version, dir)
	}

	if err = os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("error creating directory %s: %v", dir, err)
	}

	// Extract based on file extension
	if strings.HasSuffix(path, ".zip") {
		return extractZip(path, dir)
	} else if strings.HasSuffix(path, ".tar.gz") {
		return extractTarGz(path, dir)
	} else if strings.HasSuffix(path, ".tar.zst") {
		return extractTarZst(path, dir, trimPrefix, verbose)
	} else {
		return fmt.Errorf("unsupported file extension for %s %s", name, version)
	}
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

// extractTarZst extracts a tar.zst file to a destination directory
func extractTarZst(src, dst, trimPrefix string, verbose bool) error {
	if verbose {
		fmt.Printf("Extracting from %s to %s\n", src, dst)
	}

	// Open the zstd compressed file
	file, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	// Create zstd decoder
	decoder, err := zstd.NewReader(file)
	if err != nil {
		return fmt.Errorf("error creating zstd decoder: %v", err)
	}
	defer decoder.Close()

	// Create tar reader from the decompressed stream
	tr := tar.NewReader(decoder)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		name := header.Name

		if trimPrefix != "" {
			if !strings.HasPrefix(header.Name, trimPrefix) {
				continue
			}

			// Remove the trimPrefix prefix
			name = strings.TrimPrefix(header.Name, trimPrefix)
			if name == "" {
				continue
			}
		}

		path := filepath.Join(dst, name)
		if verbose {
			fmt.Printf("Extracting: %s\n", path)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(path, os.FileMode(header.Mode)); err != nil {
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
		case tar.TypeSymlink:
			dir := filepath.Dir(path)
			if err := os.MkdirAll(dir, 0755); err != nil {
				return fmt.Errorf("error creating directory %s: %v", dir, err)
			}

			// Remove existing symlink if it exists
			if err := os.RemoveAll(path); err != nil {
				return fmt.Errorf("error removing existing symlink %s: %v", path, err)
			}

			// Create new symlink
			if err := os.Symlink(header.Linkname, path); err != nil {
				return fmt.Errorf("error creating symlink %s -> %s: %v", path, header.Linkname, err)
			}
		case tar.TypeLink:
			dir := filepath.Dir(path)
			if err := os.MkdirAll(dir, 0755); err != nil {
				return fmt.Errorf("error creating directory %s: %v", dir, err)
			}

			// Remove existing file if it exists
			if err := os.RemoveAll(path); err != nil {
				return fmt.Errorf("error removing existing file %s: %v", path, err)
			}

			// Create hard link relative to the destination directory
			targetPath := filepath.Join(dst, strings.TrimPrefix(header.Linkname, "python/install/"))
			if err := os.Link(targetPath, path); err != nil {
				return fmt.Errorf("error creating hard link %s -> %s: %v", path, targetPath, err)
			}
		}
	}

	return nil
}
