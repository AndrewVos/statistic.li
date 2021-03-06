package main

import (
	"strconv"
	"testing"
	"time"
)

func TestShowsGoogleSearchResults(t *testing.T) {
	DeleteAllClientHits()

	referer := "http://www.google.co.uk/url?q=andrew%20vos"
	hit := &ClientHit{
		ClientID: "site.com",
		Referer:  referer,
	}
	hit.Save()
	hits := LatestClientHits("site.com")
	expected := "Search: andrew vos"
	if hits[0].Referer != expected {
		t.Errorf("Expected referer to be converted to a google search result, but got this instead\n%q\n", hits[0].Referer)
	}
}

func TestReferersShowDirectHits(t *testing.T) {
	DeleteAllClientHits()
	hit := &ClientHit{
		ClientID: "site.com",
	}
	hit.Save()
	hit = LatestClientHits("site.com")[0]
	expected := "(direct)"
	if hit.Referer != expected {
		t.Errorf("Expected referer to be (direct), but got this instead\n%q\n", hit.Referer)
	}
}

func TestLatestHits(t *testing.T) {
	DeleteAllClientHits()
	for userId := 0; userId < 10; userId++ {
		hit1 := &ClientHit{ClientID: "site.com", UserID: strconv.Itoa(userId)}
		hit2 := &ClientHit{ClientID: "site.com", UserID: strconv.Itoa(userId)}
		hit1.Save()
		hit2.Save()
		hit2.Date = time.Now().Add(-5 * time.Minute)
	}
	if len(LatestClientHits("site.com")) != 10 {
		t.Errorf("Expected to only see the latest client hits")
	}
}
