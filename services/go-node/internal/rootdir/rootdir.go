// Package rootdir provides helpers to locate the project root directory
// at runtime, so that all file paths can be expressed as relative paths
// and resolved portably on any machine.
package rootdir

import (
	"os"
	"path/filepath"
	"sync"
)

var (
	once sync.Once
	root string
)

// Root returns the absolute path of the project root (digital-eval-system).
// It walks up from the current working directory looking for a marker
// directory ("services") that indicates the repo root.
// The result is cached after the first call.
func Root() string {
	once.Do(func() {
		// Start from the current working directory (the binary is always
		// launched from a known location, e.g. services/go-node or the repo root).
		dir, err := os.Getwd()
		if err != nil {
			// last resort: use executable location
			exe, _ := os.Executable()
			dir = filepath.Dir(exe)
		}
		dir, _ = filepath.Abs(dir)

		// Walk up until we find a directory that contains both "services" and "infra"
		// (unique markers for the digital-eval-system repo root).
		for {
			if hasMarkers(dir) {
				root = filepath.ToSlash(dir)
				return
			}
			parent := filepath.Dir(dir)
			if parent == dir {
				// reached filesystem root without finding markers — fall back to cwd
				cwd, _ := os.Getwd()
				root = filepath.ToSlash(cwd)
				return
			}
			dir = parent
		}
	})
	return root
}

// Resolve joins the project root with relPath and returns the absolute,
// slash-normalised result.  relPath should use forward slashes.
func Resolve(relPath string) string {
	return filepath.ToSlash(filepath.Join(Root(), relPath))
}

// hasMarkers checks whether dir looks like the project root.
func hasMarkers(dir string) bool {
	for _, marker := range []string{"services", "infra"} {
		info, err := os.Stat(filepath.Join(dir, marker))
		if err != nil || !info.IsDir() {
			return false
		}
	}
	return true
}
