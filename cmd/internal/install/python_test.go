package install

import (
	"strings"
	"testing"
)

func TestGetPythonURL(t *testing.T) {
	  files := `cpython-3.10.15+20241016-aarch64-apple-darwin-debug-full.tar.zst
cpython-3.10.15+20241016-aarch64-apple-darwin-install_only.tar.gz
cpython-3.10.15+20241016-aarch64-apple-darwin-install_only_stripped.tar.gz
cpython-3.10.15+20241016-aarch64-apple-darwin-pgo+lto-full.tar.zst
cpython-3.10.15+20241016-aarch64-apple-darwin-pgo-full.tar.zst
cpython-3.10.15+20241016-aarch64-unknown-linux-gnu-debug-full.tar.zst
cpython-3.10.15+20241016-aarch64-unknown-linux-gnu-install_only.tar.gz
cpython-3.10.15+20241016-aarch64-unknown-linux-gnu-install_only_stripped.tar.gz
cpython-3.10.15+20241016-aarch64-unknown-linux-gnu-lto-full.tar.zst
cpython-3.10.15+20241016-aarch64-unknown-linux-gnu-noopt-full.tar.zst
cpython-3.10.15+20241016-armv7-unknown-linux-gnueabi-debug-full.tar.zst
cpython-3.10.15+20241016-armv7-unknown-linux-gnueabi-install_only.tar.gz
cpython-3.10.15+20241016-armv7-unknown-linux-gnueabi-install_only_stripped.tar.gz
cpython-3.10.15+20241016-armv7-unknown-linux-gnueabi-lto-full.tar.zst
cpython-3.10.15+20241016-armv7-unknown-linux-gnueabi-noopt-full.tar.zst
cpython-3.10.15+20241016-armv7-unknown-linux-gnueabihf-debug-full.tar.zst
cpython-3.10.15+20241016-armv7-unknown-linux-gnueabihf-install_only.tar.gz
cpython-3.10.15+20241016-armv7-unknown-linux-gnueabihf-install_only_stripped.tar.gz
cpython-3.10.15+20241016-armv7-unknown-linux-gnueabihf-lto-full.tar.zst
cpython-3.10.15+20241016-armv7-unknown-linux-gnueabihf-noopt-full.tar.zst
cpython-3.10.15+20241016-i686-pc-windows-msvc-install_only.tar.gz
cpython-3.10.15+20241016-i686-pc-windows-msvc-install_only_stripped.tar.gz
cpython-3.10.15+20241016-i686-pc-windows-msvc-pgo-full.tar.zst
cpython-3.10.15+20241016-i686-pc-windows-msvc-shared-install_only.tar.gz
cpython-3.10.15+20241016-i686-pc-windows-msvc-shared-install_only_stripped.tar.gz
cpython-3.10.15+20241016-i686-pc-windows-msvc-shared-pgo-full.tar.zst
cpython-3.10.15+20241016-ppc64le-unknown-linux-gnu-debug-full.tar.zst
cpython-3.10.15+20241016-ppc64le-unknown-linux-gnu-install_only.tar.gz
cpython-3.10.15+20241016-ppc64le-unknown-linux-gnu-install_only_stripped.tar.gz
cpython-3.10.15+20241016-ppc64le-unknown-linux-gnu-lto-full.tar.zst
cpython-3.10.15+20241016-ppc64le-unknown-linux-gnu-noopt-full.tar.zst
cpython-3.10.15+20241016-s390x-unknown-linux-gnu-debug-full.tar.zst
cpython-3.10.15+20241016-s390x-unknown-linux-gnu-install_only.tar.gz
cpython-3.10.15+20241016-s390x-unknown-linux-gnu-install_only_stripped.tar.gz
cpython-3.10.15+20241016-s390x-unknown-linux-gnu-lto-full.tar.zst
cpython-3.10.15+20241016-s390x-unknown-linux-gnu-noopt-full.tar.zst
cpython-3.10.15+20241016-x86_64-apple-darwin-debug-full.tar.zst
cpython-3.10.15+20241016-x86_64-apple-darwin-install_only.tar.gz
cpython-3.10.15+20241016-x86_64-apple-darwin-install_only_stripped.tar.gz
cpython-3.10.15+20241016-x86_64-apple-darwin-pgo+lto-full.tar.zst
cpython-3.10.15+20241016-x86_64-apple-darwin-pgo-full.tar.zst
cpython-3.10.15+20241016-x86_64-pc-windows-msvc-install_only.tar.gz
cpython-3.10.15+20241016-x86_64-pc-windows-msvc-install_only_stripped.tar.gz
cpython-3.10.15+20241016-x86_64-pc-windows-msvc-pgo-full.tar.zst
cpython-3.10.15+20241016-x86_64-pc-windows-msvc-shared-install_only.tar.gz
cpython-3.10.15+20241016-x86_64-pc-windows-msvc-shared-install_only_stripped.tar.gz
cpython-3.10.15+20241016-x86_64-pc-windows-msvc-shared-pgo-full.tar.zst
cpython-3.10.15+20241016-x86_64-unknown-linux-gnu-debug-full.tar.zst
cpython-3.10.15+20241016-x86_64-unknown-linux-gnu-install_only.tar.gz
cpython-3.10.15+20241016-x86_64-unknown-linux-gnu-install_only_stripped.tar.gz
cpython-3.10.15+20241016-x86_64-unknown-linux-gnu-pgo+lto-full.tar.zst
cpython-3.10.15+20241016-x86_64-unknown-linux-gnu-pgo-full.tar.zst
cpython-3.10.15+20241016-x86_64-unknown-linux-musl-debug-full.tar.zst
cpython-3.10.15+20241016-x86_64-unknown-linux-musl-install_only.tar.gz
cpython-3.10.15+20241016-x86_64-unknown-linux-musl-install_only_stripped.tar.gz
cpython-3.10.15+20241016-x86_64-unknown-linux-musl-lto-full.tar.zst
cpython-3.10.15+20241016-x86_64-unknown-linux-musl-noopt-full.tar.zst
cpython-3.10.15+20241016-x86_64_v2-unknown-linux-gnu-debug-full.tar.zst
cpython-3.10.15+20241016-x86_64_v2-unknown-linux-gnu-install_only.tar.gz
cpython-3.10.15+20241016-x86_64_v2-unknown-linux-gnu-install_only_stripped.tar.gz
cpython-3.10.15+20241016-x86_64_v2-unknown-linux-gnu-pgo+lto-full.tar.zst
cpython-3.10.15+20241016-x86_64_v2-unknown-linux-gnu-pgo-full.tar.zst
cpython-3.10.15+20241016-x86_64_v2-unknown-linux-musl-debug-full.tar.zst
cpython-3.10.15+20241016-x86_64_v2-unknown-linux-musl-install_only.tar.gz
cpython-3.10.15+20241016-x86_64_v2-unknown-linux-musl-install_only_stripped.tar.gz
cpython-3.10.15+20241016-x86_64_v2-unknown-linux-musl-lto-full.tar.zst
cpython-3.10.15+20241016-x86_64_v2-unknown-linux-musl-noopt-full.tar.zst
cpython-3.10.15+20241016-x86_64_v3-unknown-linux-gnu-debug-full.tar.zst
cpython-3.10.15+20241016-x86_64_v3-unknown-linux-gnu-install_only.tar.gz
cpython-3.10.15+20241016-x86_64_v3-unknown-linux-gnu-install_only_stripped.tar.gz
cpython-3.10.15+20241016-x86_64_v3-unknown-linux-gnu-pgo+lto-full.tar.zst
cpython-3.10.15+20241016-x86_64_v3-unknown-linux-gnu-pgo-full.tar.zst
cpython-3.10.15+20241016-x86_64_v3-unknown-linux-musl-debug-full.tar.zst
cpython-3.10.15+20241016-x86_64_v3-unknown-linux-musl-install_only.tar.gz
cpython-3.10.15+20241016-x86_64_v3-unknown-linux-musl-install_only_stripped.tar.gz
cpython-3.10.15+20241016-x86_64_v3-unknown-linux-musl-lto-full.tar.zst
cpython-3.10.15+20241016-x86_64_v3-unknown-linux-musl-noopt-full.tar.zst
cpython-3.10.15+20241016-x86_64_v4-unknown-linux-gnu-debug-full.tar.zst
cpython-3.10.15+20241016-x86_64_v4-unknown-linux-gnu-install_only.tar.gz
cpython-3.10.15+20241016-x86_64_v4-unknown-linux-gnu-install_only_stripped.tar.gz
cpython-3.10.15+20241016-x86_64_v4-unknown-linux-gnu-lto-full.tar.zst
cpython-3.10.15+20241016-x86_64_v4-unknown-linux-gnu-noopt-full.tar.zst
cpython-3.10.15+20241016-x86_64_v4-unknown-linux-musl-debug-full.tar.zst
cpython-3.10.15+20241016-x86_64_v4-unknown-linux-musl-install_only.tar.gz
cpython-3.10.15+20241016-x86_64_v4-unknown-linux-musl-install_only_stripped.tar.gz
cpython-3.10.15+20241016-x86_64_v4-unknown-linux-musl-lto-full.tar.zst
cpython-3.10.15+20241016-x86_64_v4-unknown-linux-musl-noopt-full.tar.zst
cpython-3.11.10+20241016-aarch64-apple-darwin-debug-full.tar.zst
cpython-3.11.10+20241016-aarch64-apple-darwin-install_only.tar.gz
cpython-3.11.10+20241016-aarch64-apple-darwin-install_only_stripped.tar.gz
cpython-3.11.10+20241016-aarch64-apple-darwin-pgo+lto-full.tar.zst
cpython-3.11.10+20241016-aarch64-apple-darwin-pgo-full.tar.zst
cpython-3.11.10+20241016-aarch64-unknown-linux-gnu-debug-full.tar.zst
cpython-3.11.10+20241016-aarch64-unknown-linux-gnu-install_only.tar.gz
cpython-3.11.10+20241016-aarch64-unknown-linux-gnu-install_only_stripped.tar.gz
cpython-3.11.10+20241016-aarch64-unknown-linux-gnu-lto-full.tar.zst
cpython-3.11.10+20241016-aarch64-unknown-linux-gnu-noopt-full.tar.zst
cpython-3.11.10+20241016-armv7-unknown-linux-gnueabi-debug-full.tar.zst
cpython-3.11.10+20241016-armv7-unknown-linux-gnueabi-install_only.tar.gz
cpython-3.11.10+20241016-armv7-unknown-linux-gnueabi-install_only_stripped.tar.gz
cpython-3.11.10+20241016-armv7-unknown-linux-gnueabi-lto-full.tar.zst
cpython-3.11.10+20241016-armv7-unknown-linux-gnueabi-noopt-full.tar.zst
cpython-3.11.10+20241016-armv7-unknown-linux-gnueabihf-debug-full.tar.zst
cpython-3.11.10+20241016-armv7-unknown-linux-gnueabihf-install_only.tar.gz
cpython-3.11.10+20241016-armv7-unknown-linux-gnueabihf-install_only_stripped.tar.gz
cpython-3.11.10+20241016-armv7-unknown-linux-gnueabihf-lto-full.tar.zst
cpython-3.11.10+20241016-armv7-unknown-linux-gnueabihf-noopt-full.tar.zst
cpython-3.11.10+20241016-i686-pc-windows-msvc-install_only.tar.gz
cpython-3.11.10+20241016-i686-pc-windows-msvc-install_only_stripped.tar.gz
cpython-3.11.10+20241016-i686-pc-windows-msvc-pgo-full.tar.zst
cpython-3.11.10+20241016-i686-pc-windows-msvc-shared-install_only.tar.gz
cpython-3.11.10+20241016-i686-pc-windows-msvc-shared-install_only_stripped.tar.gz
cpython-3.11.10+20241016-i686-pc-windows-msvc-shared-pgo-full.tar.zst
cpython-3.11.10+20241016-ppc64le-unknown-linux-gnu-debug-full.tar.zst
cpython-3.11.10+20241016-ppc64le-unknown-linux-gnu-install_only.tar.gz
cpython-3.11.10+20241016-ppc64le-unknown-linux-gnu-install_only_stripped.tar.gz
cpython-3.11.10+20241016-ppc64le-unknown-linux-gnu-lto-full.tar.zst
cpython-3.11.10+20241016-ppc64le-unknown-linux-gnu-noopt-full.tar.zst
cpython-3.11.10+20241016-s390x-unknown-linux-gnu-debug-full.tar.zst
cpython-3.11.10+20241016-s390x-unknown-linux-gnu-install_only.tar.gz
cpython-3.11.10+20241016-s390x-unknown-linux-gnu-install_only_stripped.tar.gz
cpython-3.11.10+20241016-s390x-unknown-linux-gnu-lto-full.tar.zst
cpython-3.11.10+20241016-s390x-unknown-linux-gnu-noopt-full.tar.zst
cpython-3.11.10+20241016-x86_64-apple-darwin-debug-full.tar.zst
cpython-3.11.10+20241016-x86_64-apple-darwin-install_only.tar.gz
cpython-3.11.10+20241016-x86_64-apple-darwin-install_only_stripped.tar.gz
cpython-3.11.10+20241016-x86_64-apple-darwin-pgo+lto-full.tar.zst
cpython-3.11.10+20241016-x86_64-apple-darwin-pgo-full.tar.zst
cpython-3.11.10+20241016-x86_64-pc-windows-msvc-install_only.tar.gz
cpython-3.11.10+20241016-x86_64-pc-windows-msvc-install_only_stripped.tar.gz
cpython-3.11.10+20241016-x86_64-pc-windows-msvc-pgo-full.tar.zst
cpython-3.11.10+20241016-x86_64-pc-windows-msvc-shared-install_only.tar.gz
cpython-3.11.10+20241016-x86_64-pc-windows-msvc-shared-install_only_stripped.tar.gz
cpython-3.11.10+20241016-x86_64-pc-windows-msvc-shared-pgo-full.tar.zst
cpython-3.11.10+20241016-x86_64-unknown-linux-gnu-debug-full.tar.zst
cpython-3.11.10+20241016-x86_64-unknown-linux-gnu-install_only.tar.gz
cpython-3.11.10+20241016-x86_64-unknown-linux-gnu-install_only_stripped.tar.gz
cpython-3.11.10+20241016-x86_64-unknown-linux-gnu-pgo+lto-full.tar.zst
cpython-3.11.10+20241016-x86_64-unknown-linux-gnu-pgo-full.tar.zst
cpython-3.11.10+20241016-x86_64-unknown-linux-musl-debug-full.tar.zst
cpython-3.11.10+20241016-x86_64-unknown-linux-musl-install_only.tar.gz
cpython-3.11.10+20241016-x86_64-unknown-linux-musl-install_only_stripped.tar.gz
cpython-3.11.10+20241016-x86_64-unknown-linux-musl-lto-full.tar.zst
cpython-3.11.10+20241016-x86_64-unknown-linux-musl-noopt-full.tar.zst
cpython-3.11.10+20241016-x86_64_v2-unknown-linux-gnu-debug-full.tar.zst
cpython-3.11.10+20241016-x86_64_v2-unknown-linux-gnu-install_only.tar.gz
cpython-3.11.10+20241016-x86_64_v2-unknown-linux-gnu-install_only_stripped.tar.gz
cpython-3.11.10+20241016-x86_64_v2-unknown-linux-gnu-pgo+lto-full.tar.zst
cpython-3.11.10+20241016-x86_64_v2-unknown-linux-gnu-pgo-full.tar.zst
cpython-3.11.10+20241016-x86_64_v2-unknown-linux-musl-debug-full.tar.zst
cpython-3.11.10+20241016-x86_64_v2-unknown-linux-musl-install_only.tar.gz
cpython-3.11.10+20241016-x86_64_v2-unknown-linux-musl-install_only_stripped.tar.gz
cpython-3.11.10+20241016-x86_64_v2-unknown-linux-musl-lto-full.tar.zst
cpython-3.11.10+20241016-x86_64_v2-unknown-linux-musl-noopt-full.tar.zst
cpython-3.11.10+20241016-x86_64_v3-unknown-linux-gnu-debug-full.tar.zst
cpython-3.11.10+20241016-x86_64_v3-unknown-linux-gnu-install_only.tar.gz
cpython-3.11.10+20241016-x86_64_v3-unknown-linux-gnu-install_only_stripped.tar.gz
cpython-3.11.10+20241016-x86_64_v3-unknown-linux-gnu-pgo+lto-full.tar.zst
cpython-3.11.10+20241016-x86_64_v3-unknown-linux-gnu-pgo-full.tar.zst
cpython-3.11.10+20241016-x86_64_v3-unknown-linux-musl-debug-full.tar.zst
cpython-3.11.10+20241016-x86_64_v3-unknown-linux-musl-install_only.tar.gz
cpython-3.11.10+20241016-x86_64_v3-unknown-linux-musl-install_only_stripped.tar.gz
cpython-3.11.10+20241016-x86_64_v3-unknown-linux-musl-lto-full.tar.zst
cpython-3.11.10+20241016-x86_64_v3-unknown-linux-musl-noopt-full.tar.zst
cpython-3.11.10+20241016-x86_64_v4-unknown-linux-gnu-debug-full.tar.zst
cpython-3.11.10+20241016-x86_64_v4-unknown-linux-gnu-install_only.tar.gz
cpython-3.11.10+20241016-x86_64_v4-unknown-linux-gnu-install_only_stripped.tar.gz
cpython-3.11.10+20241016-x86_64_v4-unknown-linux-gnu-lto-full.tar.zst
cpython-3.11.10+20241016-x86_64_v4-unknown-linux-gnu-noopt-full.tar.zst
cpython-3.11.10+20241016-x86_64_v4-unknown-linux-musl-debug-full.tar.zst
cpython-3.11.10+20241016-x86_64_v4-unknown-linux-musl-install_only.tar.gz
cpython-3.11.10+20241016-x86_64_v4-unknown-linux-musl-install_only_stripped.tar.gz
cpython-3.11.10+20241016-x86_64_v4-unknown-linux-musl-lto-full.tar.zst
cpython-3.11.10+20241016-x86_64_v4-unknown-linux-musl-noopt-full.tar.zst
cpython-3.12.7+20241016-aarch64-apple-darwin-debug-full.tar.zst
cpython-3.12.7+20241016-aarch64-apple-darwin-install_only.tar.gz
cpython-3.12.7+20241016-aarch64-apple-darwin-install_only_stripped.tar.gz
cpython-3.12.7+20241016-aarch64-apple-darwin-pgo+lto-full.tar.zst
cpython-3.12.7+20241016-aarch64-apple-darwin-pgo-full.tar.zst
cpython-3.12.7+20241016-aarch64-unknown-linux-gnu-debug-full.tar.zst
cpython-3.12.7+20241016-aarch64-unknown-linux-gnu-install_only.tar.gz
cpython-3.12.7+20241016-aarch64-unknown-linux-gnu-install_only_stripped.tar.gz
cpython-3.12.7+20241016-aarch64-unknown-linux-gnu-lto-full.tar.zst
cpython-3.12.7+20241016-aarch64-unknown-linux-gnu-noopt-full.tar.zst
cpython-3.12.7+20241016-armv7-unknown-linux-gnueabi-debug-full.tar.zst
cpython-3.12.7+20241016-armv7-unknown-linux-gnueabi-install_only.tar.gz
cpython-3.12.7+20241016-armv7-unknown-linux-gnueabi-install_only_stripped.tar.gz
cpython-3.12.7+20241016-armv7-unknown-linux-gnueabi-lto-full.tar.zst
cpython-3.12.7+20241016-armv7-unknown-linux-gnueabi-noopt-full.tar.zst
cpython-3.12.7+20241016-armv7-unknown-linux-gnueabihf-debug-full.tar.zst
cpython-3.12.7+20241016-armv7-unknown-linux-gnueabihf-install_only.tar.gz
cpython-3.12.7+20241016-armv7-unknown-linux-gnueabihf-install_only_stripped.tar.gz
cpython-3.12.7+20241016-armv7-unknown-linux-gnueabihf-lto-full.tar.zst
cpython-3.12.7+20241016-armv7-unknown-linux-gnueabihf-noopt-full.tar.zst
cpython-3.12.7+20241016-i686-pc-windows-msvc-install_only.tar.gz
cpython-3.12.7+20241016-i686-pc-windows-msvc-install_only_stripped.tar.gz
cpython-3.12.7+20241016-i686-pc-windows-msvc-pgo-full.tar.zst
cpython-3.12.7+20241016-i686-pc-windows-msvc-shared-install_only.tar.gz
cpython-3.12.7+20241016-i686-pc-windows-msvc-shared-install_only_stripped.tar.gz
cpython-3.12.7+20241016-i686-pc-windows-msvc-shared-pgo-full.tar.zst
cpython-3.12.7+20241016-ppc64le-unknown-linux-gnu-debug-full.tar.zst
cpython-3.12.7+20241016-ppc64le-unknown-linux-gnu-install_only.tar.gz
cpython-3.12.7+20241016-ppc64le-unknown-linux-gnu-install_only_stripped.tar.gz
cpython-3.12.7+20241016-ppc64le-unknown-linux-gnu-lto-full.tar.zst
cpython-3.12.7+20241016-ppc64le-unknown-linux-gnu-noopt-full.tar.zst
cpython-3.12.7+20241016-s390x-unknown-linux-gnu-debug-full.tar.zst
cpython-3.12.7+20241016-s390x-unknown-linux-gnu-install_only.tar.gz
cpython-3.12.7+20241016-s390x-unknown-linux-gnu-install_only_stripped.tar.gz
cpython-3.12.7+20241016-s390x-unknown-linux-gnu-lto-full.tar.zst
cpython-3.12.7+20241016-s390x-unknown-linux-gnu-noopt-full.tar.zst
cpython-3.12.7+20241016-x86_64-apple-darwin-debug-full.tar.zst
cpython-3.12.7+20241016-x86_64-apple-darwin-install_only.tar.gz
cpython-3.12.7+20241016-x86_64-apple-darwin-install_only_stripped.tar.gz
cpython-3.12.7+20241016-x86_64-apple-darwin-pgo+lto-full.tar.zst
cpython-3.12.7+20241016-x86_64-apple-darwin-pgo-full.tar.zst
cpython-3.12.7+20241016-x86_64-pc-windows-msvc-install_only.tar.gz
cpython-3.12.7+20241016-x86_64-pc-windows-msvc-install_only_stripped.tar.gz
cpython-3.12.7+20241016-x86_64-pc-windows-msvc-pgo-full.tar.zst
cpython-3.12.7+20241016-x86_64-pc-windows-msvc-shared-install_only.tar.gz
cpython-3.12.7+20241016-x86_64-pc-windows-msvc-shared-install_only_stripped.tar.gz
cpython-3.12.7+20241016-x86_64-pc-windows-msvc-shared-pgo-full.tar.zst
cpython-3.12.7+20241016-x86_64-unknown-linux-gnu-debug-full.tar.zst
cpython-3.12.7+20241016-x86_64-unknown-linux-gnu-install_only.tar.gz
cpython-3.12.7+20241016-x86_64-unknown-linux-gnu-install_only_stripped.tar.gz
cpython-3.12.7+20241016-x86_64-unknown-linux-gnu-pgo+lto-full.tar.zst
cpython-3.12.7+20241016-x86_64-unknown-linux-gnu-pgo-full.tar.zst
cpython-3.12.7+20241016-x86_64-unknown-linux-musl-debug-full.tar.zst
cpython-3.12.7+20241016-x86_64-unknown-linux-musl-install_only.tar.gz
cpython-3.12.7+20241016-x86_64-unknown-linux-musl-install_only_stripped.tar.gz
cpython-3.12.7+20241016-x86_64-unknown-linux-musl-lto-full.tar.zst
cpython-3.12.7+20241016-x86_64-unknown-linux-musl-noopt-full.tar.zst
cpython-3.12.7+20241016-x86_64_v2-unknown-linux-gnu-debug-full.tar.zst
cpython-3.12.7+20241016-x86_64_v2-unknown-linux-gnu-install_only.tar.gz
cpython-3.12.7+20241016-x86_64_v2-unknown-linux-gnu-install_only_stripped.tar.gz
cpython-3.12.7+20241016-x86_64_v2-unknown-linux-gnu-pgo+lto-full.tar.zst
cpython-3.12.7+20241016-x86_64_v2-unknown-linux-gnu-pgo-full.tar.zst
cpython-3.12.7+20241016-x86_64_v2-unknown-linux-musl-debug-full.tar.zst
cpython-3.12.7+20241016-x86_64_v2-unknown-linux-musl-install_only.tar.gz
cpython-3.12.7+20241016-x86_64_v2-unknown-linux-musl-install_only_stripped.tar.gz
cpython-3.12.7+20241016-x86_64_v2-unknown-linux-musl-lto-full.tar.zst
cpython-3.12.7+20241016-x86_64_v2-unknown-linux-musl-noopt-full.tar.zst
cpython-3.12.7+20241016-x86_64_v3-unknown-linux-gnu-debug-full.tar.zst
cpython-3.12.7+20241016-x86_64_v3-unknown-linux-gnu-install_only.tar.gz
cpython-3.12.7+20241016-x86_64_v3-unknown-linux-gnu-install_only_stripped.tar.gz
cpython-3.12.7+20241016-x86_64_v3-unknown-linux-gnu-pgo+lto-full.tar.zst
cpython-3.12.7+20241016-x86_64_v3-unknown-linux-gnu-pgo-full.tar.zst
cpython-3.12.7+20241016-x86_64_v3-unknown-linux-musl-debug-full.tar.zst
cpython-3.12.7+20241016-x86_64_v3-unknown-linux-musl-install_only.tar.gz
cpython-3.12.7+20241016-x86_64_v3-unknown-linux-musl-install_only_stripped.tar.gz
cpython-3.12.7+20241016-x86_64_v3-unknown-linux-musl-lto-full.tar.zst
cpython-3.12.7+20241016-x86_64_v3-unknown-linux-musl-noopt-full.tar.zst
cpython-3.12.7+20241016-x86_64_v4-unknown-linux-gnu-debug-full.tar.zst
cpython-3.12.7+20241016-x86_64_v4-unknown-linux-gnu-install_only.tar.gz
cpython-3.12.7+20241016-x86_64_v4-unknown-linux-gnu-install_only_stripped.tar.gz
cpython-3.12.7+20241016-x86_64_v4-unknown-linux-gnu-lto-full.tar.zst
cpython-3.12.7+20241016-x86_64_v4-unknown-linux-gnu-noopt-full.tar.zst
cpython-3.12.7+20241016-x86_64_v4-unknown-linux-musl-debug-full.tar.zst
cpython-3.12.7+20241016-x86_64_v4-unknown-linux-musl-install_only.tar.gz
cpython-3.12.7+20241016-x86_64_v4-unknown-linux-musl-install_only_stripped.tar.gz
cpython-3.12.7+20241016-x86_64_v4-unknown-linux-musl-lto-full.tar.zst
cpython-3.12.7+20241016-x86_64_v4-unknown-linux-musl-noopt-full.tar.zst
cpython-3.13.0+20241016-aarch64-apple-darwin-debug-full.tar.zst
cpython-3.13.0+20241016-aarch64-apple-darwin-freethreaded+debug-full.tar.zst
cpython-3.13.0+20241016-aarch64-apple-darwin-freethreaded+pgo+lto-full.tar.zst
cpython-3.13.0+20241016-aarch64-apple-darwin-freethreaded+pgo-full.tar.zst
cpython-3.13.0+20241016-aarch64-apple-darwin-install_only.tar.gz
cpython-3.13.0+20241016-aarch64-apple-darwin-install_only_stripped.tar.gz
cpython-3.13.0+20241016-aarch64-apple-darwin-pgo+lto-full.tar.zst
cpython-3.13.0+20241016-aarch64-apple-darwin-pgo-full.tar.zst
cpython-3.13.0+20241016-aarch64-unknown-linux-gnu-debug-full.tar.zst
cpython-3.13.0+20241016-aarch64-unknown-linux-gnu-freethreaded+debug-full.tar.zst
cpython-3.13.0+20241016-aarch64-unknown-linux-gnu-freethreaded+lto-full.tar.zst
cpython-3.13.0+20241016-aarch64-unknown-linux-gnu-freethreaded+noopt-full.tar.zst
cpython-3.13.0+20241016-aarch64-unknown-linux-gnu-install_only.tar.gz
cpython-3.13.0+20241016-aarch64-unknown-linux-gnu-install_only_stripped.tar.gz
cpython-3.13.0+20241016-aarch64-unknown-linux-gnu-lto-full.tar.zst
cpython-3.13.0+20241016-aarch64-unknown-linux-gnu-noopt-full.tar.zst
cpython-3.13.0+20241016-armv7-unknown-linux-gnueabi-debug-full.tar.zst
cpython-3.13.0+20241016-armv7-unknown-linux-gnueabi-freethreaded+debug-full.tar.zst
cpython-3.13.0+20241016-armv7-unknown-linux-gnueabi-freethreaded+lto-full.tar.zst
cpython-3.13.0+20241016-armv7-unknown-linux-gnueabi-freethreaded+noopt-full.tar.zst
cpython-3.13.0+20241016-armv7-unknown-linux-gnueabi-install_only.tar.gz
cpython-3.13.0+20241016-armv7-unknown-linux-gnueabi-install_only_stripped.tar.gz
cpython-3.13.0+20241016-armv7-unknown-linux-gnueabi-lto-full.tar.zst
cpython-3.13.0+20241016-armv7-unknown-linux-gnueabi-noopt-full.tar.zst
cpython-3.13.0+20241016-armv7-unknown-linux-gnueabihf-debug-full.tar.zst
cpython-3.13.0+20241016-armv7-unknown-linux-gnueabihf-freethreaded+debug-full.tar.zst
cpython-3.13.0+20241016-armv7-unknown-linux-gnueabihf-freethreaded+lto-full.tar.zst
cpython-3.13.0+20241016-armv7-unknown-linux-gnueabihf-freethreaded+noopt-full.tar.zst
cpython-3.13.0+20241016-armv7-unknown-linux-gnueabihf-install_only.tar.gz
cpython-3.13.0+20241016-armv7-unknown-linux-gnueabihf-install_only_stripped.tar.gz
cpython-3.13.0+20241016-armv7-unknown-linux-gnueabihf-lto-full.tar.zst
cpython-3.13.0+20241016-armv7-unknown-linux-gnueabihf-noopt-full.tar.zst
cpython-3.13.0+20241016-i686-pc-windows-msvc-freethreaded+pgo-full.tar.zst
cpython-3.13.0+20241016-i686-pc-windows-msvc-install_only.tar.gz
cpython-3.13.0+20241016-i686-pc-windows-msvc-install_only_stripped.tar.gz
cpython-3.13.0+20241016-i686-pc-windows-msvc-pgo-full.tar.zst
cpython-3.13.0+20241016-i686-pc-windows-msvc-shared-freethreaded+pgo-full.tar.zst
cpython-3.13.0+20241016-i686-pc-windows-msvc-shared-install_only.tar.gz
cpython-3.13.0+20241016-i686-pc-windows-msvc-shared-install_only_stripped.tar.gz
cpython-3.13.0+20241016-i686-pc-windows-msvc-shared-pgo-full.tar.zst
cpython-3.13.0+20241016-ppc64le-unknown-linux-gnu-debug-full.tar.zst
cpython-3.13.0+20241016-ppc64le-unknown-linux-gnu-freethreaded+debug-full.tar.zst
cpython-3.13.0+20241016-ppc64le-unknown-linux-gnu-freethreaded+lto-full.tar.zst
cpython-3.13.0+20241016-ppc64le-unknown-linux-gnu-freethreaded+noopt-full.tar.zst
cpython-3.13.0+20241016-ppc64le-unknown-linux-gnu-install_only.tar.gz
cpython-3.13.0+20241016-ppc64le-unknown-linux-gnu-install_only_stripped.tar.gz
cpython-3.13.0+20241016-ppc64le-unknown-linux-gnu-lto-full.tar.zst
cpython-3.13.0+20241016-ppc64le-unknown-linux-gnu-noopt-full.tar.zst
cpython-3.13.0+20241016-s390x-unknown-linux-gnu-debug-full.tar.zst
cpython-3.13.0+20241016-s390x-unknown-linux-gnu-freethreaded+debug-full.tar.zst
cpython-3.13.0+20241016-s390x-unknown-linux-gnu-freethreaded+lto-full.tar.zst
cpython-3.13.0+20241016-s390x-unknown-linux-gnu-freethreaded+noopt-full.tar.zst
cpython-3.13.0+20241016-s390x-unknown-linux-gnu-install_only.tar.gz
cpython-3.13.0+20241016-s390x-unknown-linux-gnu-install_only_stripped.tar.gz
cpython-3.13.0+20241016-s390x-unknown-linux-gnu-lto-full.tar.zst
cpython-3.13.0+20241016-s390x-unknown-linux-gnu-noopt-full.tar.zst
cpython-3.13.0+20241016-x86_64-apple-darwin-debug-full.tar.zst
cpython-3.13.0+20241016-x86_64-apple-darwin-freethreaded+debug-full.tar.zst
cpython-3.13.0+20241016-x86_64-apple-darwin-freethreaded+pgo+lto-full.tar.zst
cpython-3.13.0+20241016-x86_64-apple-darwin-freethreaded+pgo-full.tar.zst
cpython-3.13.0+20241016-x86_64-apple-darwin-install_only.tar.gz
cpython-3.13.0+20241016-x86_64-apple-darwin-install_only_stripped.tar.gz
cpython-3.13.0+20241016-x86_64-apple-darwin-pgo+lto-full.tar.zst
cpython-3.13.0+20241016-x86_64-apple-darwin-pgo-full.tar.zst
cpython-3.13.0+20241016-x86_64-pc-windows-msvc-freethreaded+pgo-full.tar.zst
cpython-3.13.0+20241016-x86_64-pc-windows-msvc-install_only.tar.gz
cpython-3.13.0+20241016-x86_64-pc-windows-msvc-install_only_stripped.tar.gz
cpython-3.13.0+20241016-x86_64-pc-windows-msvc-pgo-full.tar.zst
cpython-3.13.0+20241016-x86_64-pc-windows-msvc-shared-freethreaded+pgo-full.tar.zst
cpython-3.13.0+20241016-x86_64-pc-windows-msvc-shared-install_only.tar.gz
cpython-3.13.0+20241016-x86_64-pc-windows-msvc-shared-install_only_stripped.tar.gz
cpython-3.13.0+20241016-x86_64-pc-windows-msvc-shared-pgo-full.tar.zst
cpython-3.13.0+20241016-x86_64-unknown-linux-gnu-debug-full.tar.zst
cpython-3.13.0+20241016-x86_64-unknown-linux-gnu-freethreaded+debug-full.tar.zst
cpython-3.13.0+20241016-x86_64-unknown-linux-gnu-freethreaded+pgo+lto-full.tar.zst
cpython-3.13.0+20241016-x86_64-unknown-linux-gnu-freethreaded+pgo-full.tar.zst
cpython-3.13.0+20241016-x86_64-unknown-linux-gnu-install_only.tar.gz
cpython-3.13.0+20241016-x86_64-unknown-linux-gnu-install_only_stripped.tar.gz
cpython-3.13.0+20241016-x86_64-unknown-linux-gnu-pgo+lto-full.tar.zst
cpython-3.13.0+20241016-x86_64-unknown-linux-gnu-pgo-full.tar.zst
cpython-3.13.0+20241016-x86_64-unknown-linux-musl-debug-full.tar.zst
cpython-3.13.0+20241016-x86_64-unknown-linux-musl-install_only.tar.gz
cpython-3.13.0+20241016-x86_64-unknown-linux-musl-install_only_stripped.tar.gz
cpython-3.13.0+20241016-x86_64-unknown-linux-musl-lto-full.tar.zst
cpython-3.13.0+20241016-x86_64-unknown-linux-musl-noopt-full.tar.zst
cpython-3.13.0+20241016-x86_64_v2-unknown-linux-gnu-debug-full.tar.zst
cpython-3.13.0+20241016-x86_64_v2-unknown-linux-gnu-freethreaded+debug-full.tar.zst
cpython-3.13.0+20241016-x86_64_v2-unknown-linux-gnu-freethreaded+pgo+lto-full.tar.zst
cpython-3.13.0+20241016-x86_64_v2-unknown-linux-gnu-freethreaded+pgo-full.tar.zst
cpython-3.13.0+20241016-x86_64_v2-unknown-linux-gnu-install_only.tar.gz
cpython-3.13.0+20241016-x86_64_v2-unknown-linux-gnu-install_only_stripped.tar.gz
cpython-3.13.0+20241016-x86_64_v2-unknown-linux-gnu-pgo+lto-full.tar.zst
cpython-3.13.0+20241016-x86_64_v2-unknown-linux-gnu-pgo-full.tar.zst
cpython-3.13.0+20241016-x86_64_v2-unknown-linux-musl-debug-full.tar.zst
cpython-3.13.0+20241016-x86_64_v2-unknown-linux-musl-install_only.tar.gz
cpython-3.13.0+20241016-x86_64_v2-unknown-linux-musl-install_only_stripped.tar.gz
cpython-3.13.0+20241016-x86_64_v2-unknown-linux-musl-lto-full.tar.zst
cpython-3.13.0+20241016-x86_64_v2-unknown-linux-musl-noopt-full.tar.zst
cpython-3.13.0+20241016-x86_64_v3-unknown-linux-gnu-debug-full.tar.zst
cpython-3.13.0+20241016-x86_64_v3-unknown-linux-gnu-freethreaded+debug-full.tar.zst
cpython-3.13.0+20241016-x86_64_v3-unknown-linux-gnu-freethreaded+pgo+lto-full.tar.zst
cpython-3.13.0+20241016-x86_64_v3-unknown-linux-gnu-freethreaded+pgo-full.tar.zst
cpython-3.13.0+20241016-x86_64_v3-unknown-linux-gnu-install_only.tar.gz
cpython-3.13.0+20241016-x86_64_v3-unknown-linux-gnu-install_only_stripped.tar.gz
cpython-3.13.0+20241016-x86_64_v3-unknown-linux-gnu-pgo+lto-full.tar.zst
cpython-3.13.0+20241016-x86_64_v3-unknown-linux-gnu-pgo-full.tar.zst
cpython-3.13.0+20241016-x86_64_v3-unknown-linux-musl-debug-full.tar.zst
cpython-3.13.0+20241016-x86_64_v3-unknown-linux-musl-install_only.tar.gz
cpython-3.13.0+20241016-x86_64_v3-unknown-linux-musl-install_only_stripped.tar.gz
cpython-3.13.0+20241016-x86_64_v3-unknown-linux-musl-lto-full.tar.zst
cpython-3.13.0+20241016-x86_64_v3-unknown-linux-musl-noopt-full.tar.zst
cpython-3.13.0+20241016-x86_64_v4-unknown-linux-gnu-debug-full.tar.zst
cpython-3.13.0+20241016-x86_64_v4-unknown-linux-gnu-freethreaded+debug-full.tar.zst
cpython-3.13.0+20241016-x86_64_v4-unknown-linux-gnu-freethreaded+lto-full.tar.zst
cpython-3.13.0+20241016-x86_64_v4-unknown-linux-gnu-freethreaded+noopt-full.tar.zst
cpython-3.13.0+20241016-x86_64_v4-unknown-linux-gnu-install_only.tar.gz
cpython-3.13.0+20241016-x86_64_v4-unknown-linux-gnu-install_only_stripped.tar.gz
cpython-3.13.0+20241016-x86_64_v4-unknown-linux-gnu-lto-full.tar.zst
cpython-3.13.0+20241016-x86_64_v4-unknown-linux-gnu-noopt-full.tar.zst
cpython-3.13.0+20241016-x86_64_v4-unknown-linux-musl-debug-full.tar.zst
cpython-3.13.0+20241016-x86_64_v4-unknown-linux-musl-install_only.tar.gz
cpython-3.13.0+20241016-x86_64_v4-unknown-linux-musl-install_only_stripped.tar.gz
cpython-3.13.0+20241016-x86_64_v4-unknown-linux-musl-lto-full.tar.zst
cpython-3.13.0+20241016-x86_64_v4-unknown-linux-musl-noopt-full.tar.zst
cpython-3.9.20+20241016-aarch64-apple-darwin-debug-full.tar.zst
cpython-3.9.20+20241016-aarch64-apple-darwin-install_only.tar.gz
cpython-3.9.20+20241016-aarch64-apple-darwin-install_only_stripped.tar.gz
cpython-3.9.20+20241016-aarch64-apple-darwin-pgo+lto-full.tar.zst
cpython-3.9.20+20241016-aarch64-apple-darwin-pgo-full.tar.zst
cpython-3.9.20+20241016-aarch64-unknown-linux-gnu-debug-full.tar.zst
cpython-3.9.20+20241016-aarch64-unknown-linux-gnu-install_only.tar.gz
cpython-3.9.20+20241016-aarch64-unknown-linux-gnu-install_only_stripped.tar.gz
cpython-3.9.20+20241016-aarch64-unknown-linux-gnu-lto-full.tar.zst
cpython-3.9.20+20241016-aarch64-unknown-linux-gnu-noopt-full.tar.zst
cpython-3.9.20+20241016-armv7-unknown-linux-gnueabi-debug-full.tar.zst
cpython-3.9.20+20241016-armv7-unknown-linux-gnueabi-install_only.tar.gz
cpython-3.9.20+20241016-armv7-unknown-linux-gnueabi-install_only_stripped.tar.gz
cpython-3.9.20+20241016-armv7-unknown-linux-gnueabi-lto-full.tar.zst
cpython-3.9.20+20241016-armv7-unknown-linux-gnueabi-noopt-full.tar.zst
cpython-3.9.20+20241016-armv7-unknown-linux-gnueabihf-debug-full.tar.zst
cpython-3.9.20+20241016-armv7-unknown-linux-gnueabihf-install_only.tar.gz
cpython-3.9.20+20241016-armv7-unknown-linux-gnueabihf-install_only_stripped.tar.gz
cpython-3.9.20+20241016-armv7-unknown-linux-gnueabihf-lto-full.tar.zst
cpython-3.9.20+20241016-armv7-unknown-linux-gnueabihf-noopt-full.tar.zst
cpython-3.9.20+20241016-i686-pc-windows-msvc-install_only.tar.gz
cpython-3.9.20+20241016-i686-pc-windows-msvc-install_only_stripped.tar.gz
cpython-3.9.20+20241016-i686-pc-windows-msvc-pgo-full.tar.zst
cpython-3.9.20+20241016-i686-pc-windows-msvc-shared-install_only.tar.gz
cpython-3.9.20+20241016-i686-pc-windows-msvc-shared-install_only_stripped.tar.gz
cpython-3.9.20+20241016-i686-pc-windows-msvc-shared-pgo-full.tar.zst
cpython-3.9.20+20241016-ppc64le-unknown-linux-gnu-debug-full.tar.zst
cpython-3.9.20+20241016-ppc64le-unknown-linux-gnu-install_only.tar.gz
cpython-3.9.20+20241016-ppc64le-unknown-linux-gnu-install_only_stripped.tar.gz
cpython-3.9.20+20241016-ppc64le-unknown-linux-gnu-lto-full.tar.zst
cpython-3.9.20+20241016-ppc64le-unknown-linux-gnu-noopt-full.tar.zst
cpython-3.9.20+20241016-s390x-unknown-linux-gnu-debug-full.tar.zst
cpython-3.9.20+20241016-s390x-unknown-linux-gnu-install_only.tar.gz
cpython-3.9.20+20241016-s390x-unknown-linux-gnu-install_only_stripped.tar.gz
cpython-3.9.20+20241016-s390x-unknown-linux-gnu-lto-full.tar.zst
cpython-3.9.20+20241016-s390x-unknown-linux-gnu-noopt-full.tar.zst
cpython-3.9.20+20241016-x86_64-apple-darwin-debug-full.tar.zst
cpython-3.9.20+20241016-x86_64-apple-darwin-install_only.tar.gz
cpython-3.9.20+20241016-x86_64-apple-darwin-install_only_stripped.tar.gz
cpython-3.9.20+20241016-x86_64-apple-darwin-pgo+lto-full.tar.zst
cpython-3.9.20+20241016-x86_64-apple-darwin-pgo-full.tar.zst
cpython-3.9.20+20241016-x86_64-pc-windows-msvc-install_only.tar.gz
cpython-3.9.20+20241016-x86_64-pc-windows-msvc-install_only_stripped.tar.gz
cpython-3.9.20+20241016-x86_64-pc-windows-msvc-pgo-full.tar.zst
cpython-3.9.20+20241016-x86_64-pc-windows-msvc-shared-install_only.tar.gz
cpython-3.9.20+20241016-x86_64-pc-windows-msvc-shared-install_only_stripped.tar.gz
cpython-3.9.20+20241016-x86_64-pc-windows-msvc-shared-pgo-full.tar.zst
cpython-3.9.20+20241016-x86_64-unknown-linux-gnu-debug-full.tar.zst
cpython-3.9.20+20241016-x86_64-unknown-linux-gnu-install_only.tar.gz
cpython-3.9.20+20241016-x86_64-unknown-linux-gnu-install_only_stripped.tar.gz
cpython-3.9.20+20241016-x86_64-unknown-linux-gnu-pgo+lto-full.tar.zst
cpython-3.9.20+20241016-x86_64-unknown-linux-gnu-pgo-full.tar.zst
cpython-3.9.20+20241016-x86_64-unknown-linux-musl-debug-full.tar.zst
cpython-3.9.20+20241016-x86_64-unknown-linux-musl-install_only.tar.gz
cpython-3.9.20+20241016-x86_64-unknown-linux-musl-install_only_stripped.tar.gz
cpython-3.9.20+20241016-x86_64-unknown-linux-musl-lto-full.tar.zst
cpython-3.9.20+20241016-x86_64-unknown-linux-musl-noopt-full.tar.zst
cpython-3.9.20+20241016-x86_64_v2-unknown-linux-gnu-debug-full.tar.zst
cpython-3.9.20+20241016-x86_64_v2-unknown-linux-gnu-install_only.tar.gz
cpython-3.9.20+20241016-x86_64_v2-unknown-linux-gnu-install_only_stripped.tar.gz
cpython-3.9.20+20241016-x86_64_v2-unknown-linux-gnu-pgo+lto-full.tar.zst
cpython-3.9.20+20241016-x86_64_v2-unknown-linux-gnu-pgo-full.tar.zst
cpython-3.9.20+20241016-x86_64_v2-unknown-linux-musl-debug-full.tar.zst
cpython-3.9.20+20241016-x86_64_v2-unknown-linux-musl-install_only.tar.gz
cpython-3.9.20+20241016-x86_64_v2-unknown-linux-musl-install_only_stripped.tar.gz
cpython-3.9.20+20241016-x86_64_v2-unknown-linux-musl-lto-full.tar.zst
cpython-3.9.20+20241016-x86_64_v2-unknown-linux-musl-noopt-full.tar.zst
cpython-3.9.20+20241016-x86_64_v3-unknown-linux-gnu-debug-full.tar.zst
cpython-3.9.20+20241016-x86_64_v3-unknown-linux-gnu-install_only.tar.gz
cpython-3.9.20+20241016-x86_64_v3-unknown-linux-gnu-install_only_stripped.tar.gz
cpython-3.9.20+20241016-x86_64_v3-unknown-linux-gnu-pgo+lto-full.tar.zst
cpython-3.9.20+20241016-x86_64_v3-unknown-linux-gnu-pgo-full.tar.zst
cpython-3.9.20+20241016-x86_64_v3-unknown-linux-musl-debug-full.tar.zst
cpython-3.9.20+20241016-x86_64_v3-unknown-linux-musl-install_only.tar.gz
cpython-3.9.20+20241016-x86_64_v3-unknown-linux-musl-install_only_stripped.tar.gz
cpython-3.9.20+20241016-x86_64_v3-unknown-linux-musl-lto-full.tar.zst
cpython-3.9.20+20241016-x86_64_v3-unknown-linux-musl-noopt-full.tar.zst
cpython-3.9.20+20241016-x86_64_v4-unknown-linux-gnu-debug-full.tar.zst
cpython-3.9.20+20241016-x86_64_v4-unknown-linux-gnu-install_only.tar.gz
cpython-3.9.20+20241016-x86_64_v4-unknown-linux-gnu-install_only_stripped.tar.gz
cpython-3.9.20+20241016-x86_64_v4-unknown-linux-gnu-lto-full.tar.zst
cpython-3.9.20+20241016-x86_64_v4-unknown-linux-gnu-noopt-full.tar.zst
cpython-3.9.20+20241016-x86_64_v4-unknown-linux-musl-debug-full.tar.zst
cpython-3.9.20+20241016-x86_64_v4-unknown-linux-musl-install_only.tar.gz
cpython-3.9.20+20241016-x86_64_v4-unknown-linux-musl-install_only_stripped.tar.gz
cpython-3.9.20+20241016-x86_64_v4-unknown-linux-musl-lto-full.tar.zst
cpython-3.9.20+20241016-x86_64_v4-unknown-linux-musl-noopt-full.tar.zst
  `
	tests := []struct {
		name        string
		arch        string
		os          string
		freeThreaded bool
		debug       bool
		want        string
		wantErr     bool
	}{
		{
			name:        "darwin-arm64-freethreaded-debug",
			arch:        "arm64",
			os:          "darwin",
			freeThreaded: true,
			debug:       true,
			want:        "cpython-3.13.0+20241016-aarch64-apple-darwin-freethreaded+debug-full.tar.zst",
		},
		{
			name:        "darwin-amd64-freethreaded-pgo",
			arch:        "amd64",
			os:          "darwin",
			freeThreaded: true,
			debug:       false,
			want:        "cpython-3.13.0+20241016-x86_64-apple-darwin-freethreaded+pgo-full.tar.zst",
		},
		{
			name:        "darwin-amd64-debug",
			arch:        "amd64",
			os:          "darwin",
			freeThreaded: false,
			debug:       true,
			want:        "cpython-3.13.0+20241016-x86_64-apple-darwin-debug-full.tar.zst",
		},
		{
			name:        "darwin-amd64-pgo",
			arch:        "amd64",
			os:          "darwin",
			freeThreaded: false,
			debug:       false,
			want:        "cpython-3.13.0+20241016-x86_64-apple-darwin-pgo-full.tar.zst",
		},
		{
			name:        "linux-amd64-freethreaded-debug",
			arch:        "amd64",
			os:          "linux",
			freeThreaded: true,
			debug:       true,
			want:        "cpython-3.13.0+20241016-x86_64-unknown-linux-gnu-freethreaded+debug-full.tar.zst",
		},
		{
			name:        "windows-amd64-freethreaded-pgo",
			arch:        "amd64",
			os:          "windows",
			freeThreaded: true,
			debug:       false,
			want:        "cpython-3.13.0+20241016-x86_64-pc-windows-msvc-shared-freethreaded+pgo-full.tar.zst",
		},
		{
			name:        "windows-386-freethreaded-pgo",
			arch:        "386",
			os:          "windows",
			freeThreaded: true,
			debug:       false,
			want:        "cpython-3.13.0+20241016-i686-pc-windows-msvc-shared-freethreaded+pgo-full.tar.zst",
		},
		{
			name:        "unsupported-arch",
			arch:        "mips",
			os:          "linux",
			freeThreaded: false,
			debug:       false,
			want:        "",
			wantErr:     true,
		},
		{
			name:        "unsupported-os",
			arch:        "amd64",
			os:          "freebsd",
			freeThreaded: false,
			debug:       false,
			want:        "",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getPythonURL("3.13.0", "20241016", tt.arch, tt.os, tt.freeThreaded, tt.debug)
			
			if tt.wantErr {
				if got != "" {
					t.Errorf("getPythonURL() = %v, want empty string for error case", got)
				}
				return
			}

			if got == "" {
				t.Errorf("getPythonURL() returned empty string, want %v", tt.want)
				return
			}

			// Extract filename from URL
			parts := strings.Split(got, "/")
			filename := parts[len(parts)-1]

			if filename != tt.want {
				t.Errorf("getPythonURL() = %v, want %v", filename, tt.want)
			}

			// Verify the file exists in the provided file list
			if !strings.Contains(files, filename) {
				t.Errorf("getPythonURL() generated filename %v that doesn't exist in available files", filename)
			}
		})
	}
}
