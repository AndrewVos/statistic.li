package main

type Uniques struct {
	Count int
}

func GetUniques(clientId string) Uniques {
	hits := LatestClientHits(clientId)
	count := map[string]bool{}
	for _, hit := range hits {
		count[hit.UserID] = true
	}
	return Uniques{len(count)}
}
