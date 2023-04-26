package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"time"

	"gopkg.in/yaml.v3"
)

var (
	ErrEmptyHistory = errors.New("empty history")
)

type TidyIface interface {
	Read() error
	Dup() error
	Re(exp string) error
	Recent(d time.Duration) error
	Write(dryRun bool) error
	Len() int
	Combine(other TidyIface) error
	// rm index -2
	RmPre() error
	RmLastN(n int) error
}

type FishHistoryItem struct {
	Cmd  string `json:"cmd"`
	When int64  `json:"when"`
}
type FishHistory []FishHistoryItem

func (a FishHistory) Len() int           { return len(a) }
func (a FishHistory) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a FishHistory) Less(i, j int) bool { return a[i].When < a[j].When }
func (h *FishHistory) Read() error {
	bytes, err := os.ReadFile(hoem2AbsPath(FISH_HISTORY_RELATIVE_PATH))
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
		return os.WriteFile(hoem2AbsPath(FISH_HISTORY_RELATIVE_PATH), []byte(""), 0644)
	}
	bytes, err := yaml.Marshal(h)
	if err != nil {
		return err
	}
	if dryRun {
		return os.WriteFile("fish_"+DRY_RUN_OUTPUT_PATH, bytes, 0644)
	}
	return os.WriteFile(hoem2AbsPath(FISH_HISTORY_RELATIVE_PATH), bytes, 0644)
}
func (h *FishHistory) Combine(other TidyIface) error {
	switch other.(type) {
	case *FishHistory:
		*h = append(*h, *other.(*FishHistory)...)
		return nil
	case *ZshHistory:
		zshHistory := *other.(*ZshHistory)
		for idx := range zshHistory {
			*h = append(*h, FishHistoryItem{
				Cmd:  zshHistory[idx].Cmd,
				When: zshHistory[idx].When,
			})
		}
		return nil
	default:
		return fmt.Errorf("unsupported type %T", other)
	}
}
func (h *FishHistory) RmPre() error {
	len_ := len(*h)
	if len_ > 0 {
		rmLastCmd := (*h)[len_-1]
		*h = (*h)[:len_-2]
		*h = append(*h, rmLastCmd)
		return nil
	}
	return ErrEmptyHistory
}
func (h *FishHistory) RmLastN(n int) error {
	for i := 0; i < n; i++ {
		if err := h.RmPre(); err != nil {
			return err
		}
	}
	return nil
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
	bytes, err := os.ReadFile(hoem2AbsPath(ZSH_HISTORY_RELATIVE_PATH))
	if err != nil {
		return err
	}
	macthes := zshRegExp.FindAllString(string(bytes), -1)
	for _, match := range macthes {
		groups := zshRegExp.FindStringSubmatch(match)
		if len(groups) != 3 {
			continue
		}
		t := groups[1]
		time_, err := strconv.ParseInt(t, 10, 64)
		if err != nil {
			continue
		}
		*h = append(*h, ZshHistoryItem{
			Cmd:  groups[2],
			When: time_,
		})
	}
	return nil
}
func (h *ZshHistory) Dup() error {
	sort.Sort(*h)
	historyMap := make(map[string]*ZshHistoryItem)
	for idx := range *h {
		if _, ok := historyMap[(*h)[idx].Cmd]; ok {
			continue
		}
		historyMap[(*h)[idx].Cmd] = &(*h)[idx]
	}
	*h = make([]ZshHistoryItem, 0)
	for k := range historyMap {
		*h = append(*h, *historyMap[k])
	}
	sort.Sort(*h)
	return nil
}
func (h *ZshHistory) Re(exp string) error {
	re, err := regexp.Compile(exp)
	if err != nil {
		return err
	}
	newHistory := make([]ZshHistoryItem, 0)
	for idx := range *h {
		if !re.MatchString((*h)[idx].Cmd) {
			newHistory = append(newHistory, (*h)[idx])
		}
	}
	*h = newHistory
	return nil
}
func (h *ZshHistory) Recent(d time.Duration) error {
	t := time.Now().Add(-d).Unix()
	newHistory := make([]ZshHistoryItem, 0)
	for idx := range *h {
		if (*h)[idx].When < t {
			newHistory = append(newHistory, (*h)[idx])
		}
	}
	*h = newHistory
	return nil
}
func (h *ZshHistory) Write(dryRun bool) error {
	if len(*h) == 0 {
		return os.WriteFile(hoem2AbsPath(ZSH_HISTORY_RELATIVE_PATH), []byte(""), 0644)
	}
	var buffer bytes.Buffer
	for idx := range *h {
		buffer.WriteString(fmt.Sprintf(": %d:0;%s\n", (*h)[idx].When, (*h)[idx].Cmd))
	}
	if dryRun {
		return os.WriteFile("zsh_"+DRY_RUN_OUTPUT_PATH, buffer.Bytes(), 0644)
	}
	return os.WriteFile(hoem2AbsPath(ZSH_HISTORY_RELATIVE_PATH), buffer.Bytes(), 0644)
}
func (h *ZshHistory) Combine(other TidyIface) error {
	switch other.(type) {
	case *ZshHistory:
		*h = append(*h, *other.(*ZshHistory)...)
		return nil
	case *FishHistory:
		fishHistory := *other.(*FishHistory)
		for idx := range fishHistory {
			*h = append(*h, ZshHistoryItem{
				Cmd:  fishHistory[idx].Cmd,
				When: fishHistory[idx].When,
			})
		}
		return nil
	default:
		return fmt.Errorf("unsupported type %T", other)
	}
}
func (h *ZshHistory) RmPre() error {
	len_ := len(*h)
	if len_ >= 2 {
		rmLastCmd := (*h)[len_-1]
		*h = (*h)[:len_-2]
		*h = append(*h, rmLastCmd)
		return nil
	}
	return ErrEmptyHistory
}
func (h *ZshHistory) RmLastN(n int) error {
	for i := 0; i < n; i++ {
		if err := h.RmPre(); err != nil {
			return err
		}
	}
	return nil
}
