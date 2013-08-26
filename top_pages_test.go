package main

import (
	"strconv"
	"testing"
)

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

	topPages := GetTopPages("client.com")
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
