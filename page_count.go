package main

import (
	"sort"
)

type PageCount struct {
	Page  string
	Count int
}

type PageCounts []*PageCount

func (s PageCounts) Len() int           { return len(s) }
func (s PageCounts) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s PageCounts) Less(i, j int) bool { return s[i].Count > s[j].Count }

func TopPages(clientId string) PageCounts {
	session, err := connectToMongo()
	defer session.Close()
	if err != nil {
		logError("mongo", err)
		return nil
	}

	query := LatestClientHits(session, clientId)

	var hits []ClientHit
	query.All(&hits)

	topPagesMap := map[string]*PageCount{}
	for _, hit := range hits {
		if _, ok := topPagesMap[hit.Page]; ok {
			count := topPagesMap[hit.Page]
			count.Count += 1
		} else {
			topPagesMap[hit.Page] = &PageCount{Page: hit.Page, Count: 1}
		}
	}

	var topPages PageCounts
	for _, pageImpressionCount := range topPagesMap {
		topPages = append(topPages, pageImpressionCount)
	}
	sort.Sort(topPages)
	if topPages == nil {
		return PageCounts{}
	}
	return topPages
}
