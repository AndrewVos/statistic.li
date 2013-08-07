package main

type PageHitCount struct {
  Referer string
  Count int
}

type PageHitCounts []*PageHitCount
func (s PageHitCounts) Len() int           { return len(s) }
func (s PageHitCounts) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s PageHitCounts) Less(i, j int) bool { return s[i].Count > s[j].Count }
