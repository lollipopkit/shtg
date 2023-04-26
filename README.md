## Shell History Tool

### Features
- [x] Support zsh / fish
- [x] Sync history between multiple shells
- [x] Tidy history using `duplicate` / `regexp` / `duration` / `last`

### Usage
Common usage:
```bash
# remove duplicate history
shtg dup
# rm previous cmd
shtg previous
shtg p
# remove history with regexp pattern
shtg re 'pattern'
# remove history with duration
shtg recent 24h
shtg r 12h
# remove last N history
shtg last 10
shtg l 3
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

COMMANDS:
   dup, d       Remove duplicate history
   re           Remove history which match regex
   recent, r    Remove history in duration
   previous, p  Remove previous cmd
   last, l      Remove last N cmd
   sync, s      Sync history between zsh / fish
   help, h      Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --type value, -t value  fish / zsh
   --dry-run, -d           without write to file (default: false)
   --path value, -p value  history file path
   --help, -h              show help
```

### Issues
- [ ] Will ignore fish history attr `paths`