package main

import (
	"os"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

type TidyIface interface {
	Read() error
	Dup() error
	Re() error
	Old() error
	Write() error
}

type FishHistoryItem struct {
	Cmd   string   `json:"cmd"`
	When  int      `json:"when"`
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
func (h *FishHistory) Re() error {
	return nil
}
func (h *FishHistory) Old() error {
	return nil
}
func (h *FishHistory) Write() error {
	bytes, err := yaml.Marshal(h)
	if err != nil {
		return err
	}
	return os.WriteFile(relativePath(FISH_HISTORY_RELATIVE_PATH), bytes, 0644)
}

type ZshHistoryItem struct {
	Cmd  string
	When int
}
type ZshHistory []ZshHistoryItem
func (a ZshHistory) Len() int           { return len(a) }
func (a ZshHistory) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ZshHistory) Less(i, j int) bool { return a[i].When < a[j].When }
func (h *ZshHistory) Read() error {
	bytes, err := os.ReadFile(relativePath(ZSH_HISTORY_RELATIVE_PATH))
	if err != nil {
		return err
	}
	lines := strings.Split(string(bytes), "\n")
	for idx := range lines {
		if lines[idx] == "" {
			continue
		}
		*h = append(*h, ZshHistoryItem{
			Cmd:  lines[idx],
			When: idx,
		})
	}
	return nil
}
func (h *ZshHistory) Dup() error {
	return nil
}
func (h *ZshHistory) Re() error {
	return nil
}
func (h *ZshHistory) Old() error {
	return nil
}
func (h *ZshHistory) Write() error {
	return nil
}
