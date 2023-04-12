package main

import (
	"time"

	"github.com/urfave/cli/v2"
)

type ShellType string

const (
	Fish ShellType = "fish"
	Zsh  ShellType = "zsh"
)

type Mode string

const (
	ModeDup    Mode = "dup"
	ModeRe     Mode = "re"
	ModeRecent Mode = "recent"
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
	case ModeRecent:
		// shtg old 1d
		return ctx.NArg() >= 1
	default:
		panic("Unknown mode" + string(m))
	}
}
