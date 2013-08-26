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

func TestTopReferers(t *testing.T) {
	DeleteAllClientHits()
	for userId := 0; userId < 20; userId++ {
		hit1 := &ClientHit{ClientID: "client.com", Referer: "site.com", UserID: strconv.Itoa(userId)}
		hit1.Save()
	}
	for userId := 0; userId < 10; userId++ {
		hit2 := &ClientHit{ClientID: "client.com", Referer: "othersite.com", UserID: strconv.Itoa(userId)}
		hit2.Save()
	}
	referers := TopReferers("client.com")
	if referers[0].Referer != "site.com" || referers[0].Count != 20 {
		t.Errorf("Expected top referer to be site.com, with count 20, but was %q, with count %v", referers[0].Referer, referers[0].Count)
	}
}

func TestTopPages(t *testing.T) {
	DeleteAllClientHits()

	for i := 0; i < 20; i++ {
		hit := &ClientHit{ClientID: "client.com", Page: "client.com/page1.html", UserID: strconv.Itoa(i)}
		hit.Save()
	}

	for i := 0; i < 10; i++ {
		hit := &ClientHit{ClientID: "client.com", Page: "client.com/page2.html", UserID: strconv.Itoa(i)}
		hit.Save()
	}

	topPages := TopPages("client.com")
	if len(topPages) != 2 {
		t.Error("Expected there to be two top pages")
	}

	if expected := "client.com/page1.html"; topPages[0].Page != expected {
		t.Errorf("Expected first page to be %q, but was %q", expected, topPages[0].Page)
	}

	if expected := "client.com/page2.html"; topPages[1].Page != expected {
		t.Errorf("Expected first page to be %q, but was %q", expected, topPages[1].Page)
	}
}
