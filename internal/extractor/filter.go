package extractor

import (
	"slices"
)

type Filter interface {
	isDuplicate(candidate string) bool
}

type MemoryFilter struct {
	Words []string
}

func (filter *MemoryFilter) isDuplicate(candidate string) bool {
	return slices.Contains(filter.Words, candidate)
}
