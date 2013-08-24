package main

type UniquesCount struct {
	Count int
}

func Uniques(clientId string) UniquesCount {
	hits := LatestClientHits(clientId)
	count := map[string]bool{}
	for _, hit := range hits {
		count[hit.UserID] = true
	}
	return UniquesCount{len(count)}
}
