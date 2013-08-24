package main

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"time"
)

type ClientHit struct {
	ClientID string
	UserID   string
	Date     time.Time
	Referer  string
	Page     string
}

func (c *ClientHit) Save() error {
	session, err := connectToMongo()
	if err != nil {
		return err
	}
	defer session.Close()
	c.Date = time.Now()
	collection := session.DB("").C("ClientHits")
	err = collection.Insert(c)
	if err != nil {
		logError("mongo", err)
	}
	return nil
}

func LatestClientHits(session *mgo.Session, clientId string) *mgo.Query {
	collection := session.DB("").C("ClientHits")
	after := time.Now().Add(-5 * time.Minute)
	query := collection.Find(bson.M{"clientid": clientId, "date": bson.M{"$gte": after}})
	return query
}
