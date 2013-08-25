package main

import (
	"strconv"
	"testing"
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

func TestUniqueUsers(t *testing.T) {
	DeleteAllClientHits()
	for userId := 0; userId < 10; userId++ {
		hit1 := &ClientHit{ClientID: "site.com", UserID: strconv.Itoa(userId)}
		hit2 := &ClientHit{ClientID: "site.com", UserID: strconv.Itoa(userId)}
		hit1.Save()
		hit2.Save()
	}
	uniques := Uniques("site.com")
	if uniques.Count != 10 {
		t.Errorf("Expected unique users to be %v, but was %v", 10, uniques.Count)
	}
}
