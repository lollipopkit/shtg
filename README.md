## Shell History Tools

### Features
- [x] Support zsh / fish
- [x] Sync history between multiple shells
- [x] Tidy history using `duplicate` / `regexp` / `duration` / `last`

### Usage
Common usage:
```bash
# remove duplicate history
shtg dup
# rm last cmd
# (rm `shtg rmlast` hiistory & `YOUR LAST CMD` history)
shtg rmlast
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

COMMANDS:
   dup, d      Remove duplicate history
   re, r       Remove history which match regex
   recent, o   Remove history in duration
   rmlast, rl  Remove last cmd
   sync, s     Sync history between zsh / fish
   help, h     Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --type value, -t value  fish / zsh
   --dry-run, -d           without write to file (default: false)
   --path value, -p value  history file path
   --help, -h              show help
```

### Issues
- [ ] Will ignore fish history attr `paths`