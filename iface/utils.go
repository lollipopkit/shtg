package iface

import (
	"os"
	"path/filepath"
)

// Add HOME path to relative path
func home2AbsPath(relativePath string) string {
	return filepath.Join(os.Getenv("HOME"), relativePath)
}
