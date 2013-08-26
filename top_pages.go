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
	hits := LatestClientHits(clientId)

	topPagesMap := map[string]*TopPage{}
	for _, hit := range hits {
		if _, ok := topPagesMap[hit.Page]; ok {
			count := topPagesMap[hit.Page]
			count.Count += 1
		} else {
			topPagesMap[hit.Page] = &TopPage{Page: hit.Page, Count: 1}
		}
	}

	var topPages TopPages
	for _, pageImpressionCount := range topPagesMap {
		topPages = append(topPages, pageImpressionCount)
	}
	sort.Sort(topPages)
	if topPages == nil {
		return TopPages{}
	}
	return topPages
}
