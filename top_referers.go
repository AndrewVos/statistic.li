package main

import (
	"sort"
)

type TopReferer struct {
	Referer string
	Count   int
}

type TopReferers []*TopReferer

func (s TopReferers) Len() int           { return len(s) }
func (s TopReferers) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s TopReferers) Less(i, j int) bool { return s[i].Count > s[j].Count }

func GetTopReferers(clientId string) TopReferers {
	counts := make(map[string]int)
	for _, hit := range LatestClientHits(clientId) {
		if _, ok := counts[hit.Referer]; !ok {
			counts[hit.Referer] = 0
		}
		counts[hit.Referer] += 1
	}
	var referers TopReferers
	for referer, count := range counts {
		referers = append(referers, &TopReferer{Referer: referer, Count: count})
	}
	sort.Sort(referers)
	if referers == nil {
		return TopReferers{}
	}
	return referers
}
