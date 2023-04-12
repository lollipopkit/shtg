package main

import (
	"os"

	"github.com/lollipopkit/gommon/term"
	"github.com/urfave/cli/v2"
)

func run() {
	app := cli.App{
		Name:        "shtg",
		Usage:       "Shell History Tool written in Go",
		Description: "Shell history tool for zsh / fish",
		Suggest:     true,
		Copyright:   "2023 lollipopkit",
		Commands: []*cli.Command{
			{
				Name:    "dup",
				Aliases: []string{"d"},
				Action: func(ctx *cli.Context) error {
					return tidy(ctx, ModeDup)
				},
				Description: "remove duplicate history",
				Usage:       "shtg dup",
			},
			{
				Name:    "re",
				Aliases: []string{"r"},
				Action: func(ctx *cli.Context) error {
					return tidy(ctx, ModeRe)
				},
				Description: "remove history which match regex",
				Usage:       "shtg re 'scp xx x:/xxx'",
			},
			{
				Name:    "recent",
				Aliases: []string{"o"},
				Action: func(ctx *cli.Context) error {
					return tidy(ctx, ModeRecent)
				},
				Description: "remove history in duration",
				Usage:       "shtg recent 12h",
			},
			{
				Name:    "sync",
				Aliases: []string{"s"},
				Action: func(ctx *cli.Context) error {
					return sync()
				},
				Description: "sync history between zsh / fish",
				Usage:       "shtg sync",
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "type",
				Aliases: []string{"t"},
				Usage:   "fish / zsh",
			},
			&cli.BoolFlag{
				Name:    "dry-run",
				Aliases: []string{"d"},
				Value:   false,
			},
			&cli.StringFlag{
				Name:    "path",
				Aliases: []string{"p"},
				Usage:   "history file path",
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		term.Err(err.Error())
	}
}

func tidy(c *cli.Context, mode Mode) error {
	_typ := c.String("type")
	var typ ShellType
	if _typ == "" {
		typ = getShellType()
	} else {
		typ = ShellType(_typ)
	}

	var iface TidyIface
	switch typ {
	case Fish:
		iface = &FishHistory{}
	case Zsh:
		iface = &ZshHistory{}
	}
	err := iface.Read()
	if err != nil {
		return err
	}

	if !mode.Check(c) {
		term.Warn("Usage: " + c.Command.Usage)
		return nil
	}
	err = mode.Do(iface, c)
	if err != nil {
		return err
	}

	dryRun := c.Bool("dry-run")
	if dryRun {
		term.Info("output: " + DRY_RUN_OUTPUT)
	}
	return iface.Write(dryRun)
}

func sync() error {
	panic("not implemented")
}
