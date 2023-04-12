package main

import "regexp"

const (
	FISH_HISTORY_RELATIVE_PATH = ".local/share/fish/fish_history"
	ZSH_HISTORY_RELATIVE_PATH  = ".zsh_history"

	DRY_RUN_OUTPUT_PATH = "shtg_result.txt"
)

var (
	zshRegExp = regexp.MustCompile(`: (\d+):0;(.*)`)
)
