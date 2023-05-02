package cmd

import (
	"os"
	"strings"
)

type ShellType string

const (
	Fish ShellType = "fish"
	Zsh  ShellType = "zsh"
)

func getShellType() ShellType {
	shell := os.Getenv("SHELL")
	switch true {
	case strings.HasSuffix(shell, "zsh"):
		return Zsh
	case strings.HasSuffix(shell, "fish"):
		return Fish
	default:
		panic("Unsupport shell: " + shell)
	}
}
