package main

type UniquesCount struct {
	Count int
}

func Uniques(clientId string) UniquesCount {
	session, err := connectToMongo()
	defer session.Close()
	if err != nil {
		logError("mongo", err)
		return UniquesCount{}
	}

	query := getLatestClientHitsQuery(session, clientId)
	var distinctUserIds []string
	query.Distinct("userid", &distinctUserIds)
	return UniquesCount{len(distinctUserIds)}
}
