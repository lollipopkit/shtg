package main

import (
	"os"
	"path/filepath"
	"strings"
)

func relativePath(relativePath string) string {
	return filepath.Join(os.Getenv("HOME"), relativePath)
}

func getShellType() ShellType {
	shell := os.Getenv("SHELL")
	switch true {
	case strings.HasSuffix(shell, "zsh"):
		return Zsh
	case strings.HasSuffix(shell, "fish"):
		return Fish
	default:
		panic("Unknown shell type")
	}
}
