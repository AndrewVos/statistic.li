package main

import (
	"strconv"
	"testing"
)

func TestUniqueUsers(t *testing.T) {
	DeleteAllClientHits()
	for userId := 0; userId < 10; userId++ {
		hit1 := &ClientHit{ClientID: "site.com", UserID: strconv.Itoa(userId)}
		hit2 := &ClientHit{ClientID: "site.com", UserID: strconv.Itoa(userId)}
		hit1.Save()
		hit2.Save()
	}
	uniques := GetUniques("site.com")
	if uniques.Count != 10 {
		t.Errorf("Expected unique users to be %v, but was %v", 10, uniques.Count)
	}
}
