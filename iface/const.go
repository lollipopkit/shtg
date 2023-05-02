package iface

import "regexp"

const (
	FISH_HISTORY_RELATIVE_PATH = ".local/share/fish/fish_history"
	ZSH_HISTORY_RELATIVE_PATH  = ".zsh_history"
)

var (
	zshRegExp = regexp.MustCompile(`: (\d+):0;(.*)`)
)
