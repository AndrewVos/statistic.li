package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

var server *httptest.Server

func setup() {
	DeleteAllClientHits()
	if server != nil {
		server.Close()
	}
	server = httptest.NewServer(http.HandlerFunc(clientHandler))
}

func get(url string, cookies []*http.Cookie) (*http.Response, string) {
	client := &http.Client{}
	request, _ := http.NewRequest("GET", url, nil)
	if cookies != nil {
		for _, cookie := range cookies {
			request.AddCookie(cookie)
		}
	}

	response, _ := client.Do(request)
	body, _ := ioutil.ReadAll(response.Body)
	response.Body.Close()
	return response, string(body)
}

func expectContentType(t *testing.T, path string, contentType string) {
	response, _ := get(server.URL+path, nil)
	if response.Header["Content-Type"][0] != contentType {
		t.Errorf("Expected %q to return a Content-Type of %q, not %q\n", path, contentType, response.Header["Content-Type"][0])
	}
}

func TestContentTypes(t *testing.T) {
	setup()
	expectContentType(t, "/client/CLIENT_ID/uniques", "application/json")
	expectContentType(t, "/client/CLIENT_ID/referers", "application/json")
	expectContentType(t, "/client/CLIENT_ID/pages", "application/json")
	expectContentType(t, "/client/CLIENT_ID/tracker.gif", "image/gif")
}

func TestTrackerRespondsWithGif(t *testing.T) {
	setup()
	_, gif := get(server.URL+"/client/MY_CLIENT_ID/tracker.gif", nil)
	if gif != string(tracker_gif()) {
		t.Errorf("Expected:\n%q\nGot:\n%q\n", string(tracker_gif()), gif)
	}
}

func TestTopPagesRoute(t *testing.T) {
	setup()
	get(server.URL+"/client/site.com/tracker.gif?page=page1&referer=referer1", nil)
	_, body := get(server.URL+"/client/site.com/pages", nil)
	expected, _ := json.Marshal(TopPages("site.com"))
	if body != string(expected) {
		t.Errorf("Expected:\n%q\nGot:\n%q\n", string(expected), body)
		t.Fail()
	}
}

func TestTopReferersRoute(t *testing.T) {
	setup()
	get(server.URL+"/client/site.com/tracker.gif?page=page1&referer=referer1", nil)
	_, body := get(server.URL+"/client/site.com/referers", nil)
	expected, _ := json.Marshal(TopReferers("site.com"))
	if body != string(expected) {
		t.Errorf("Expected:\n%q\nGot:\n%q\n", string(expected), body)
		t.Fail()
	}
}

func TestStoresClientHit(t *testing.T) {
	setup()
	cookies := []*http.Cookie{
		{Name: "sts", Value: "theUserID1213"},
	}
	get(server.URL+"/client/site.com/tracker.gif?page=page1&referer=referer1", cookies)
	hits := LatestClientHits("site.com")
	if len(hits) != 1 {
		t.Errorf("Expected there to be %v hits, but there was %v", 2, len(hits))
	}
	if hits[0].Referer != "referer1" {
		t.Errorf("Expected referer to be set")
	}
	if hits[0].ClientID != "site.com" {
		t.Errorf("Expected client id to be set")
	}
	if hits[0].UserID != "theUserID1213" {
		t.Errorf("Expected user id to be set")
	}
}

func BenchmarkStoreClientHit(b *testing.B) {
	for i := 0; i < b.N; i++ {
		referer := "referer.com"
		page := "http://www.client.com/page1.html"
		get(server.URL+"/client/CLIENT_ID/tracker.gif?referer="+url.QueryEscape(referer)+"&page="+url.QueryEscape(page), nil)
	}
}
