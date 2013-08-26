package main

import (
	"sort"
)

type TopPage struct {
	Page  string
	Count int
}

type TopPages []*TopPage

func (s TopPages) Len() int           { return len(s) }
func (s TopPages) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s TopPages) Less(i, j int) bool { return s[i].Count > s[j].Count }

func GetTopPages(clientId string) TopPages {
	counts := make(map[string]int)
	for _, hit := range LatestClientHits(clientId) {
		if _, ok := counts[hit.Page]; !ok {
			counts[hit.Page] = 0
		}
		counts[hit.Page] += 1
	}

	var pages TopPages
	for page, count := range counts {
		pages = append(pages, &TopPage{Page: page, Count: count})
	}
	sort.Sort(pages)
	if pages == nil {
		return TopPages{}
	}
	return pages
}
