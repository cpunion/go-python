package install

import (
	"fmt"
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

	return downloadAndExtract("msys2", msys2Version, msys2URL, msys2Root, "", verbose)
}
