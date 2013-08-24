package main

type RefererCount struct {
	Referer string
	Count   int
}

type RefererCounts []*RefererCount

func (s RefererCounts) Len() int           { return len(s) }
func (s RefererCounts) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s RefererCounts) Less(i, j int) bool { return s[i].Count > s[j].Count }
