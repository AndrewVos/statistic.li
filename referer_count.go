package main

import (
	"sort"
)

type RefererCount struct {
	Referer string
	Count   int
}

type RefererCounts []*RefererCount

func (s RefererCounts) Len() int           { return len(s) }
func (s RefererCounts) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s RefererCounts) Less(i, j int) bool { return s[i].Count > s[j].Count }

func TopReferers(clientId string) RefererCounts {
	session, err := connectToMongo()
	defer session.Close()
	if err != nil {
		logError("mongo", err)
		return nil
	}

	query := LatestClientHits(session, clientId)
	var hits []ClientHit
	query.All(&hits)

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
	var pageHitCounts RefererCounts
	for referer, count := range pageCounts {
		pageHitCounts = append(pageHitCounts, &RefererCount{Referer: referer, Count: count})
	}
	sort.Sort(pageHitCounts)
	if pageHitCounts == nil {
		return RefererCounts{}
	}
	return pageHitCounts
}
