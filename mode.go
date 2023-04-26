package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/urfave/cli/v2"
)

type Mode string

const (
	ModeDup     Mode = "dup"
	ModeRe           = "re"
	ModeRecent       = "recent"
	ModeRmLast       = "rmlast"
	ModeRmLastN      = "rmlastn"
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
		return iface.RmLast()
	case ModeRmLastN:
		n := ctx.Args().Get(0)
		nn, err := strconv.ParseInt(n, 10, 64)
		if err != nil {
			return fmt.Errorf("Parse %s to int failed: %w", n, err)
		}
		return iface.RmLastN(int(nn))
	default:
		panic("Unknown mode" + string(m))
	}
}
func (m Mode) Check(ctx *cli.Context) bool {
	switch m {
	case ModeDup, ModeRmLast:
		// shtg dup, shtg rmlast
		return true
	case ModeRe, ModeRecent, ModeRmLastN:
		// shtg re xxx
		// shtg old 1d
		// shtg last 3
		return ctx.NArg() == 1
	default:
		panic("Unknown mode" + string(m))
	}
}
