package examiner

import (
	"io"
	"os"
	"path/filepath"
)

// SaveUploadedFile saves the multipart file content into destDir with given filename.
// It returns the absolute saved path.
func SaveUploadedFile(destDir, filename string, r io.Reader) (string, error) {
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return "", err
	}
	dst := filepath.Join(destDir, filename)
	f, err := os.Create(dst)
	if err != nil {
		return "", err
	}
	defer f.Close()
	if _, err := io.Copy(f, r); err != nil {
		return "", err
	}
	abs, err := filepath.Abs(dst)
	if err != nil {
		return dst, nil
	}
	return abs, nil
}

// RemoveFile best-effort cleanup
func RemoveFile(path string) {
	_ = os.Remove(path)
}
