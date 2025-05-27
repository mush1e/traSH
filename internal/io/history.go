package io

import "strings"

type History struct {
	entries []string
	index   int
}

func NewHistory() *History {
	return &History{
		entries: []string{},
		index:   -1,
	}
}

func (h *History) Add(cmd string) {
	if strings.TrimSpace(cmd) == "" {
		return
	}
	h.entries = append(h.entries, cmd)
	h.index = len(h.entries)
}

func (h *History) GetPrevious() string {
	if len(h.entries) == 0 || h.index <= 0 {
		return ""
	}
	h.index--
	return h.entries[h.index]
}

func (h *History) GetNext() string {
	if len(h.entries) == 0 || h.index >= len(h.entries)-1 {
		return ""
	}
	h.index++
	return h.entries[h.index]
}
