package main

import (
	"testing"
)

func TestShowsGoogleSearchResults(t *testing.T) {
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
	DeleteAllClientHits()
}
