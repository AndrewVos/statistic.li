package main

import (
	"strconv"
	"testing"
)

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
	referers := GetTopReferers("client.com")
	if referers[0].Referer != "site.com" || referers[0].Count != 20 {
		t.Errorf("Expected top referer to be site.com, with count 20, but was %q, with count %v", referers[0].Referer, referers[0].Count)
	}
}
