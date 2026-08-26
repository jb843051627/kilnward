package policy

import (
	"github.com/jb843051627/kilnward/internal/model"
	"sort"
)

func Sort(entries []model.TimelineEntry) []model.TimelineEntry {
	out := append([]model.TimelineEntry(nil), entries...)
	sort.SliceStable(out, func(i, j int) bool { return out[i].At.Before(out[j].At) })
	return out
}

func Merge(a, b []model.TimelineEntry) []model.TimelineEntry {
	merged := make([]model.TimelineEntry, 0, len(a)+len(b))
	merged = append(merged, a...)
	merged = append(merged, b...)
	return Sort(merged)
}

func Limit(entries []model.TimelineEntry, limit int) []model.TimelineEntry {
	if limit <= 0 || limit >= len(entries) {
		return entries
	}
	return entries[len(entries)-limit:]
}
