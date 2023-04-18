package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

func run() {
	app := cli.App{
		Name:        "shtg",
		Usage:       "Shell History Tool written in Go",
		Suggest:     true,
		Commands: []*cli.Command{
			{
				Name:    "dup",
				Aliases: []string{"d"},
				Action: func(ctx *cli.Context) error {
					return tidy(ctx, ModeDup)
				},
				Usage:     "Remove duplicate history",
				UsageText: "shtg dup",
			},
			{
				Name:    "re",
				Aliases: []string{"r"},
				Action: func(ctx *cli.Context) error {
					return tidy(ctx, ModeRe)
				},
				Usage:     "Remove history which match regex",
				UsageText: "shtg re 'scp xx x:/xxx'",
			},
			{
				Name:    "recent",
				Aliases: []string{"o"},
				Action: func(ctx *cli.Context) error {
					return tidy(ctx, ModeRecent)
				},
				Usage:     "Remove history in duration",
				UsageText: "shtg recent 12h",
			},
			{
				Name: "rmlast",
				Aliases: []string{"rl"},
				Action: func(ctx *cli.Context) error {
					return tidy(ctx, ModeRmLast)
				},
				Usage: "Remove last cmd",
				UsageText: "shtg rmlast",
			},
			{
				Name:    "sync",
				Aliases: []string{"s"},
				Action: func(ctx *cli.Context) error {
					return sync(ctx)
				},
				Usage:     "Sync history between zsh / fish",
				UsageText: "shtg sync",
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
				Usage:   "without write to file",
			},
			&cli.StringFlag{
				Name:    "path",
				Aliases: []string{"p"},
				Usage:   "history file path",
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		panic(err)
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
		println("Usage: " + c.Command.UsageText)
		return nil
	}
	beforeLen := iface.Len()
	err = mode.Do(iface, c)
	if err != nil {
		return err
	}
	afterLen := iface.Len()
	printChanges(typ, beforeLen, afterLen)

	dryRun := c.Bool("dry-run")
	if dryRun {
		println("output: " + DRY_RUN_OUTPUT_PATH)
	}
	return iface.Write(dryRun)
}

func sync(c *cli.Context) error {
	zsh := &ZshHistory{}
	err := zsh.Read()
	if err != nil {
		return err
	}
	fish := &FishHistory{}
	err = fish.Read()
	if err != nil {
		return err
	}

	fBeforeLen := fish.Len()
	zBeforeLen := zsh.Len()
	fish.Combine(zsh)
	zsh.Combine(fish)
	fAfterLen := fish.Len()
	zAfterLen := zsh.Len()
	printChanges(Fish, fBeforeLen, fAfterLen)
	printChanges(Zsh, zBeforeLen, zAfterLen)

	dryRun := c.Bool("dry-run")
	if dryRun {
		println("output: " + DRY_RUN_OUTPUT_PATH)
	}
	err = fish.Write(dryRun)
	if err != nil {
		return err
	}
	return zsh.Write(dryRun)
}

func printChanges(typ ShellType, beforeLen, afterLen int) {
	if beforeLen > afterLen {
		fmt.Printf(
			"[%s] Origin %d, Removed %d, Now %d",
			typ,
			beforeLen,
			beforeLen-afterLen,
			afterLen,
		)
	} else if beforeLen < afterLen {
		fmt.Printf(
			"[%s] Origin %d, Added %d, Now %d",
			typ,
			beforeLen,
			afterLen-beforeLen,
			afterLen,
		)
	} else {
		println("No history changed")
	}
}
