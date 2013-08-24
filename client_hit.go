package main

import (
	"time"
)

type ClientHit struct {
	ClientID string
	UserID   string
	Date     time.Time
	Referer  string
	Page     string
}

var launchedClientHitSaver = false
var storedClientHits = map[string][]*ClientHit{}
var clientHits = make(chan *ClientHit)

func init() {
	go saveClientHits()
}

func saveClientHits() {
	for {
		clientHit := <-clientHits
		clientHit.Date = time.Now()
		if _, ok := storedClientHits[clientHit.ClientID]; ok {
			storedClientHits[clientHit.ClientID] = append(storedClientHits[clientHit.ClientID], clientHit)
		} else {
			storedClientHits[clientHit.ClientID] = []*ClientHit{clientHit}
		}
		trimOldClientHits(clientHit.ClientID)
	}
}

func trimOldClientHits(clientId string) {
	after := time.Now().Add(-5 * time.Minute)
	newerClientHits := []*ClientHit{}
	for _, hit := range storedClientHits[clientId] {
		if hit.Date.After(after) {
			newerClientHits = append(newerClientHits, hit)
		}
	}
	storedClientHits[clientId] = newerClientHits
}

func (c *ClientHit) Save() {
	clientHits <- c
}

func DeleteAllClientHits() {
	storedClientHits = map[string][]*ClientHit{}
}

func LatestClientHits(clientId string) []*ClientHit {
	if c, ok := storedClientHits[clientId]; ok {
		return c
	} else {
		return []*ClientHit{}
	}
}
