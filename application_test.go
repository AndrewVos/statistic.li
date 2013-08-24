package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
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

func hitTracker(server *httptest.Server, clientId string, userId string, page string, referer string) {
	cookies := []*http.Cookie{
		&http.Cookie{Name: "sts", Value: userId},
	}
	get(server.URL+"/client/"+clientId+"/tracker.gif?referer="+url.QueryEscape(referer)+"&page="+url.QueryEscape(page), cookies)
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

func TestUniques(t *testing.T) {
	setup()
	numberOfUsers := 5
	clientId := "andrewvos.com"

	for userHitCount := 0; userHitCount < numberOfUsers; userHitCount++ {
		response, _ := get(server.URL+"/client/"+clientId+"/tracker.gif", nil)
		get(server.URL+"/client/"+clientId+"/tracker.gif", response.Cookies())
	}

	response, _ := http.Get(server.URL + "/client/" + clientId + "/uniques")
	responseBody, _ := ioutil.ReadAll(response.Body)
	uniques := string(responseBody)
	response.Body.Close()

	expected, _ := json.Marshal(UniquesCount{Count: numberOfUsers})

	if uniques != string(expected) {
		t.Errorf("Expected:\n%q\nGot:\n%q\n", string(expected), uniques)
		t.Fail()
	}
}

func TestEmptyReferers(t *testing.T) {
	setup()
	_, body := get(server.URL+"/client/CLIENT_ID/referers", nil)
	if body != "[]" {
		t.Errorf("Expected an empty json array, but got this:\n%q", body)
	}
}

func TestEmptyTopPages(t *testing.T) {
	setup()
	_, body := get(server.URL+"/client/CLIENT_ID/pages", nil)
	if body != "[]" {
		t.Errorf("Expected an empty json array, but got this:\n%q", body)
	}
}

func TestReferers(t *testing.T) {
	setup()
	hits := RefererCounts{
		&RefererCount{"(direct)", 11},
		&RefererCount{"http://referer.com/page1.html", 10},
		&RefererCount{"http://referer.com/page2.html", 9},
		&RefererCount{"http://referer.com/page3.html", 8},
		&RefererCount{"http://referer.com/page4.html", 7},
		&RefererCount{"http://referer.com/page5.html", 6},
		&RefererCount{"http://referer.com/page6.html", 5},
		&RefererCount{"http://referer.com/page7.html", 4},
		&RefererCount{"http://referer.com/page8.html", 3},
		&RefererCount{"http://referer.com/page9.html", 2},
		&RefererCount{"http://referer.com/page10.html", 1},
	}

	userNumber := 0
	for _, hit := range hits {
		for i := 0; i < hit.Count; i++ {
			userId := "user" + strconv.Itoa(userNumber)
			r := hit.Referer
			if hit.Referer == "(direct)" {
				r = ""
			}
			hitTracker(server, "CLIENT_ID", userId, "", r)
			hitTracker(server, "CLIENT_ID", userId, "", r)
			userNumber += 1
		}
	}

	_, body := get(server.URL+"/client/CLIENT_ID/referers", nil)
	expected, _ := json.Marshal(hits[:10])

	if body != string(expected) {
		t.Errorf("Expected:\n%q\nGot:\n%q\n", string(expected), body)
		t.Fail()
	}
}

func TestTopPages(t *testing.T) {
	setup()

	hits := PageCounts{
		&PageCount{"http://client.com/page1.html", 10},
		&PageCount{"http://client.com/page2.html", 9},
		&PageCount{"http://client.com/page3.html", 8},
		&PageCount{"http://client.com/page4.html", 7},
		&PageCount{"http://client.com/page5.html", 6},
		&PageCount{"http://client.com/page6.html", 5},
		&PageCount{"http://client.com/page7.html", 4},
		&PageCount{"http://client.com/page8.html", 3},
		&PageCount{"http://client.com/page9.html", 2},
		&PageCount{"http://client.com/page10.html", 1},
	}

	userNumber := 0
	for _, hit := range hits {
		for i := 0; i < hit.Count; i++ {
			userId := "user" + strconv.Itoa(userNumber)
			hitTracker(server, "CLIENT_ID", userId, hit.Page, "")
			userNumber += 1
		}
	}

	_, body := get(server.URL+"/client/CLIENT_ID/pages", nil)
	expected, _ := json.Marshal(hits[:10])

	if body != string(expected) {
		t.Errorf("Expected:\n%q\nGot:\n%q\n", string(expected), body)
		t.Fail()
	}
}

func BenchmarkStoreClientHit(b *testing.B) {
	for i := 0; i < b.N; i++ {
		hitTracker(server, "CLIENT_ID", "www.client.com", "http://www.client.com/page1.html", "")
	}
}
