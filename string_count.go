package main

type StringCount struct {
	String string
	Count  int
}

type StringCounts []*StringCount

func (s StringCounts) Len() int           { return len(s) }
func (s StringCounts) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s StringCounts) Less(i, j int) bool { return s[i].Count > s[j].Count }
