package main

import (
	"net/url"
	"time"
)

type ClientHit struct {
	ClientID string
	UserID   string
	Date     time.Time
	Referer  string
	Page     string
}

var storedClientHits = map[string][]*ClientHit{}
var clientHitsToAdd = make(chan *ClientHit)
var requestsForAllClientHits = make(chan *clientHitRequest)
var requestsToDeleteAllClientHits = make(chan bool)

type clientHitRequest struct {
	ClientId       string
	ResponseReader chan []*ClientHit
}

func init() {
	go saveClientHits()
}

func saveClientHits() {
	for {
		select {
		case clientHit := <-clientHitsToAdd:
			if _, ok := storedClientHits[clientHit.ClientID]; ok {
				storedClientHits[clientHit.ClientID] = append(storedClientHits[clientHit.ClientID], clientHit)
			} else {
				storedClientHits[clientHit.ClientID] = []*ClientHit{clientHit}
			}
		case clientHitRequest := <-requestsForAllClientHits:
			trimOldClientHits(clientHitRequest.ClientId)
			clientHits := storedClientHits[clientHitRequest.ClientId]
			if clientHits != nil {
				clientHitRequest.ResponseReader <- clientHits
			} else {
				clientHitRequest.ResponseReader <- []*ClientHit{}
			}
		case <-requestsToDeleteAllClientHits:
			storedClientHits = map[string][]*ClientHit{}
		}
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
	c.Date = time.Now()
	url, err := url.Parse(c.Referer)
	if err == nil {
		search := url.Query().Get("q")
		if search != "" {
			c.Referer = "Search: " + search
		}
	}
	clientHitsToAdd <- c
}

func DeleteAllClientHits() {
	requestsToDeleteAllClientHits <- true
}

func LatestClientHits(clientId string) []*ClientHit {
	responseReader := make(chan []*ClientHit)
	request := &clientHitRequest{
		ClientId:       clientId,
		ResponseReader: responseReader,
	}
	requestsForAllClientHits <- request
	return <-responseReader
}
