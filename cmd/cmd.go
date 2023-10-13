package cmd

import (
	"fmt"
	"os"

	"github.com/lollipopkit/shtg/consts"
	"github.com/lollipopkit/shtg/iface"
	"github.com/urfave/cli/v2"
)

func Run() {
	app := cli.App{
		Name:    "shtg",
		Usage:   "Shell History Tool written in Go",
		Suggest: true,
		Commands: []*cli.Command{
			{
				Name:    "dup",
				Aliases: []string{"d"},
				Action: func(ctx *cli.Context) error {
					return tidy(ctx, iface.ModeDup)
				},
				Usage:     "Remove duplicate history",
				UsageText: "shtg dup",
			},
			{
				Name: "re",
				Action: func(ctx *cli.Context) error {
					return tidy(ctx, iface.ModeRe)
				},
				Usage:     "Remove history which match regex",
				UsageText: "shtg re 'scp xx x:/xxx'",
			},
			{
				Name:    "recent",
				Aliases: []string{"r"},
				Action: func(ctx *cli.Context) error {
					return tidy(ctx, iface.ModeRecent)
				},
				Usage:     "Remove history in duration",
				UsageText: "shtg recent 12h",
			},
			{
				Name:    "previous",
				Aliases: []string{"p"},
				Action: func(ctx *cli.Context) error {
					return tidy(ctx, iface.ModePre)
				},
				Usage:     "Remove previous cmd",
				UsageText: "shtg previous",
			},
			{
				Name:    "last",
				Aliases: []string{"l"},
				Action: func(ctx *cli.Context) error {
					return tidy(ctx, iface.ModeRmLastN)
				},
				Usage:     "Remove last N cmd",
				UsageText: "shtg last",
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
			{
				Name:    "restore",
				Aliases: []string{"rs"},
				Action: func(ctx *cli.Context) error {
					return restore()
				},
				Usage:     "Restore history from previous backup",
				UsageText: "shtg restore",
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "shell",
				Aliases: []string{"s"},
				Usage:   "fish / zsh",
			},
			&cli.BoolFlag{
				Name:    "dry",
				Aliases: []string{"d"},
				Value:   false,
				Usage:   "without write to file",
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}

func tidy(c *cli.Context, mode iface.Mode) error {
	_typ := c.String("shell")
	var typ ShellType
	if _typ == "" {
		typ = getShellType()
	} else {
		typ = ShellType(_typ)
	}

	// Init instance of correct type
	var instance iface.HistoryIface
	switch typ {
	case Fish:
		instance = &iface.FishHistory{}
	case Zsh:
		instance = &iface.ZshHistory{}
	}

	// Read data from file
	err := instance.Read()
	if err != nil {
		return err
	}

	// Check args of current mode
	if !mode.Check(c) {
		println("Usage: " + c.Command.UsageText)
		return nil
	}

	// Do the task
	beforeLen := instance.Len()
	err = mode.Do(instance, c)
	if err != nil {
		return err
	}
	afterLen := instance.Len()
	printChanges(typ, beforeLen, afterLen)

	// If dra-run, then return
	dryRun := c.Bool("dry")
	if dryRun {
		println("\noutput: " + consts.DRY_RUN_OUTPUT_PATH)
		return nil
	}

	// Backup
	err = instance.Backup()
	if err != nil {
		return err
	}

	// Write task result to disk
	return instance.Write(dryRun)
}

func sync(c *cli.Context) error {
	zsh := &iface.ZshHistory{}
	err := zsh.Read()
	if err != nil {
		return err
	}
	fish := &iface.FishHistory{}
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
		println("output: " + consts.DRY_RUN_OUTPUT_PATH)
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

func restore() error {
	var iface_ iface.HistoryIface
	typ := getShellType()
	switch typ {
	case Fish:
		iface_ = &iface.FishHistory{}
	case Zsh:
		iface_ = &iface.ZshHistory{}
	}
	return iface_.Restore()
}
