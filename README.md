## Shell Histories Tools

### Features
- [x] Support zsh / fish
- [x] Sync history between multiple shells
- [x] Tidy history using `duplicate` / `regexp` / `duration`

### Usage
Common usage:
```bash
# remove duplicate history
shtg dup
# remove history with regexp pattern
shtg re 'pattern'
# remove history with duration
shtg recent 12h
# you can specify the shell
shtg -t zsh dup
# or dry-run
shtg -d dup
```

Details:
```
> shtg
NAME:
   shtg - Shell History Tool written in Go

USAGE:
   shtg [global options] command [command options] [arguments...]

DESCRIPTION:
   Shell history tool for zsh / fish

COMMANDS:
   dup, d     remove duplicate history
   re, r      remove history which match regex
   recent, o  remove history in duration
   sync, s    sync history between zsh / fish
   help, h    Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --type value, -t value  fish / zsh
   --dry-run, -d           (default: false)
   --path value, -p value  history file path
   --help, -h              show help

COPYRIGHT:
   2023 lollipopkit
```

### Issues
- [ ] Will ignore fish history attr `paths`