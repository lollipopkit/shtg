package main

import (
	"os"
	"regexp"
	"sort"
	"time"

	"gopkg.in/yaml.v3"
)

type TidyIface interface {
	Read() error
	Dup() error
	Re(exp string) error
	Recent(d time.Duration) error
	Write(dryRun bool) error
}

type FishHistoryItem struct {
	Cmd   string   `json:"cmd"`
	When  int64    `json:"when"`
	Paths []string `json:"paths"`
}
type FishHistory []FishHistoryItem

func (a FishHistory) Len() int           { return len(a) }
func (a FishHistory) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a FishHistory) Less(i, j int) bool { return a[i].When < a[j].When }
func (h *FishHistory) Read() error {
	bytes, err := os.ReadFile(relativePath(FISH_HISTORY_RELATIVE_PATH))
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(bytes, h)
	if err != nil {
		return err
	}
	return nil
}
func (h *FishHistory) Dup() error {
	sort.Sort(*h)
	historyMap := make(map[string]*FishHistoryItem)
	for idx := range *h {
		if _, ok := historyMap[(*h)[idx].Cmd]; ok {
			continue
		}
		historyMap[(*h)[idx].Cmd] = &(*h)[idx]
	}
	*h = make([]FishHistoryItem, 0)
	for k := range historyMap {
		*h = append(*h, *historyMap[k])
	}
	sort.Sort(*h)
	return nil
}
func (h *FishHistory) Re(exp string) error {
	re, err := regexp.Compile(exp)
	if err != nil {
		return err
	}
	newHistory := make([]FishHistoryItem, 0)
	for idx := range *h {
		if !re.MatchString((*h)[idx].Cmd) {
			newHistory = append(newHistory, (*h)[idx])
		}
	}
	*h = newHistory
	return nil
}
func (h *FishHistory) Recent(d time.Duration) error {
	t := time.Now().Add(-d).Unix()
	newHistory := make([]FishHistoryItem, 0)
	for idx := range *h {
		if (*h)[idx].When < t {
			newHistory = append(newHistory, (*h)[idx])
		}
	}
	*h = newHistory
	return nil
}
func (h *FishHistory) Write(dryRun bool) error {
	if len(*h) == 0 {
		return os.WriteFile(relativePath(FISH_HISTORY_RELATIVE_PATH), []byte(""), 0644)
	}
	bytes, err := yaml.Marshal(h)
	if err != nil {
		return err
	}
	if dryRun {
		return os.WriteFile(DRY_RUN_OUTPUT, bytes, 0644)
	}
	return os.WriteFile(relativePath(FISH_HISTORY_RELATIVE_PATH), bytes, 0644)
}

type ZshHistoryItem struct {
	Cmd  string
	When int64
}
type ZshHistory []ZshHistoryItem

func (a ZshHistory) Len() int           { return len(a) }
func (a ZshHistory) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ZshHistory) Less(i, j int) bool { return a[i].When < a[j].When }
func (h *ZshHistory) Read() error {
	_, err := os.ReadFile(relativePath(ZSH_HISTORY_RELATIVE_PATH))
	if err != nil {
		return err
	}

	return nil
}
func (h *ZshHistory) Dup() error {
	return nil
}
func (h *ZshHistory) Re(exp string) error {
	return nil
}
func (h *ZshHistory) Recent(d time.Duration) error {
	return nil
}
func (h *ZshHistory) Write(dryRun bool) error {
	return nil
}
