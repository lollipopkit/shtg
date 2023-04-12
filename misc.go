package main

import (
	"github.com/urfave/cli/v2"
)

type ShellType string

const (
	Fish ShellType = "fish"
	Zsh  ShellType = "zsh"
)

type Mode string

const (
	ModeDup Mode = "dup"
	ModeRe  Mode = "re"
	ModeOld Mode = "old"
)

func (m Mode) Do(iface TidyIface) error {
	switch m {
	case ModeDup:
		return iface.Dup()
	case ModeRe:
		return iface.Re()
	case ModeOld:
		return iface.Old()
	default:
		panic("Unknown mode" + string(m))
	}
}
func (m Mode) Check(ctx *cli.Context) bool {
	switch m {
	case ModeDup:
		// shtg dup
		return true
	case ModeRe:
		// shtg re xxx
		return ctx.NArg() >= 1
	case ModeOld:
		// shtg old 1d
		return ctx.NArg() >= 1
	default:
		panic("Unknown mode" + string(m))
	}
}
