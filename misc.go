package main

import (
	"os"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
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

type Mode string

const (
	ModeDup    Mode = "dup"
	ModeRe     Mode = "re"
	ModeRecent Mode = "recent"
	ModeRmLast Mode = "rmlast"
)

func (m Mode) Do(iface TidyIface, ctx *cli.Context) error {
	switch m {
	case ModeDup:
		return iface.Dup()
	case ModeRe:
		exp := ctx.Args().Get(0)
		return iface.Re(exp)
	case ModeRecent:
		d := ctx.Args().Get(0)
		dd, err := time.ParseDuration(d)
		if err != nil {
			return err
		}
		return iface.Recent(dd)
	case ModeRmLast:
		iface.RmLast()
		return nil
	default:
		panic("Unknown mode" + string(m))
	}
}
func (m Mode) Check(ctx *cli.Context) bool {
	switch m {
	case ModeDup, ModeRmLast:
		// shtg dup, shtg rmlast
		return true
	case ModeRe:
		// shtg re xxx
		return ctx.NArg() >= 1
	case ModeRecent:
		// shtg old 1d
		return ctx.NArg() >= 1
	default:
		panic("Unknown mode" + string(m))
	}
}
