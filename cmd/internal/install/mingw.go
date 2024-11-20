package install

import (
	"fmt"

	"github.com/cpunion/go-python/internal/env"
)

const (
	mingwVersion = "14.2.0"
	mingwURL     = "https://github.com/brechtsanders/winlibs_mingw/releases/download/14.2.0posix-19.1.1-12.0.0-ucrt-r2/winlibs-x86_64-posix-seh-gcc-14.2.0-llvm-19.1.1-mingw-w64ucrt-12.0.0-r2.zip"
)

func installMingw(projectPath string, verbose bool) error {
	root := env.GetMingwDir(projectPath)
	fmt.Printf("Installing mingw in %v\n", root)
	return downloadAndExtract("mingw", mingwVersion, mingwURL, root, "", verbose)
}
