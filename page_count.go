package main

type PageCount struct {
	Page  string
	Count int
}

type PageCounts []*PageCount

func (s PageCounts) Len() int           { return len(s) }
func (s PageCounts) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s PageCounts) Less(i, j int) bool { return s[i].Count > s[j].Count }
