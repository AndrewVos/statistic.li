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
	hits := LatestClientHits(clientId)

	countedPages := make(map[string]bool)
	pageCounts := make(map[string]int)

	for _, clientHit := range hits {
		if _, ok := countedPages[clientHit.UserID+clientHit.Referer]; ok == false {
			countedPages[clientHit.UserID+clientHit.Referer] = true
			if _, ok := pageCounts[clientHit.Referer]; ok != true {
				pageCounts[clientHit.Referer] = 0
			}
			pageCounts[clientHit.Referer] += 1
		}
	}
	var pageHitCounts TopReferers
	for referer, count := range pageCounts {
		pageHitCounts = append(pageHitCounts, &TopReferer{Referer: referer, Count: count})
	}
	sort.Sort(pageHitCounts)
	if pageHitCounts == nil {
		return TopReferers{}
	}
	return pageHitCounts
}
