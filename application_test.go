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

func flushDatabase() {
	session, _ := getConnection()
	defer session.Close()
	session.DB("").C("ClientHits").DropCollection()
}

func TestTrackerRespondsWithGif(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(clientHandler))
	defer server.Close()
	client := &http.Client{}
	request, _ := http.NewRequest("GET", server.URL+"/client/MY_CLIENT_ID/tracker.gif", nil)

	response, _ := client.Do(request)
	gif, _ := ioutil.ReadAll(response.Body)
	if string(gif) != string(tracker_gif()) {
		t.Errorf("Expected:\n%q\nGot:\n%q\n", string(tracker_gif()), string(gif))
	}
	response.Body.Close()

	expectedContentType := "image/gif"
	if response.Header["Content-Type"][0] != expectedContentType {
		t.Errorf("Expected a Content-Type of %q, not %q\n", expectedContentType, response.Header["Content-Type"][0])
	}
}

func TestUniqueViewsAreExposed(t *testing.T) {
	flushDatabase()
	server := httptest.NewServer(http.HandlerFunc(clientHandler))
	defer server.Close()

	ipAddress := func(lastPart int) string { return "0.0.0." + strconv.Itoa(lastPart) }
	numberOfUsers := 5
	clientId := "andrewvos.com"

	for userHitCount := 0; userHitCount < numberOfUsers; userHitCount++ {
		client := &http.Client{}
		request, _ := http.NewRequest("GET", server.URL+"/client/"+clientId+"/tracker.gif", nil)
		request.Header.Set("X-Forwarded-For", ipAddress(userHitCount%numberOfUsers))

		response, _ := client.Do(request)
		client.Do(request)
		for _, cookie := range response.Cookies() {
			request.AddCookie(cookie)
		}

		for requestCount := 0; requestCount < 5; requestCount++ {
			client.Do(request)
		}
	}

	response, _ := http.Get(server.URL + "/client/" + clientId + "/views")
	responseBody, _ := ioutil.ReadAll(response.Body)
	views := string(responseBody)

	expectedContentType := "application/json"
	if response.Header["Content-Type"][0] != expectedContentType {
		t.Errorf("Expected a Content-Type of %q, not %q\n", expectedContentType, response.Header["Content-Type"][0])
	}

	expectedViews := `{"views":10}`
	if views != expectedViews {
		t.Error(`Views was wrong, expected `, expectedViews, ` got `, views)
	}
}

func TestReferersContentType(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(clientHandler))
	defer server.Close()
	response, _ := http.Get(server.URL + "/client/CLIENT_ID/referers")
	expectedContentType := "application/json"
	if response.Header["Content-Type"][0] != expectedContentType {
		t.Errorf("Expected a Content-Type of %q, not %q\n", expectedContentType, response.Header["Content-Type"][0])
	}
}

func TestEmptyReferers(t *testing.T) {
	flushDatabase()
	server := httptest.NewServer(http.HandlerFunc(clientHandler))
	defer server.Close()

	response, _ := http.Get(server.URL + "/client/CLIENT_ID/referers")
	bytes, _ := ioutil.ReadAll(response.Body)
	response.Body.Close()
	if string(bytes) != "[]" {
		t.Errorf("Expected an empty json array, but got this:\n%q", string(bytes))
	}
}

func TestReferers(t *testing.T) {
	flushDatabase()
	server := httptest.NewServer(http.HandlerFunc(clientHandler))
	defer server.Close()

	hits := PageHitCounts{
		&PageHitCount{"(direct)", 11},
		&PageHitCount{"http://blablabla.com/page1.html", 10},
		&PageHitCount{"http://blablabla.com/page2.html", 9},
		&PageHitCount{"http://blablabla.com/page3.html", 8},
		&PageHitCount{"http://blablabla.com/page4.html", 7},
		&PageHitCount{"http://blablabla.com/page5.html", 6},
		&PageHitCount{"http://blablabla.com/page6.html", 5},
		&PageHitCount{"http://blablabla.com/page7.html", 4},
		&PageHitCount{"http://blablabla.com/page8.html", 3},
		&PageHitCount{"http://blablabla.com/page9.html", 2},
		&PageHitCount{"http://blablabla.com/page10.html", 1},
	}

	userNumber := 0
	for _, hit := range hits {
		for i := 0; i < hit.Count; i++ {
			userId := "user" + strconv.Itoa(userNumber)
			r := hit.Referer
			if hit.Referer == "(direct)" {
				r = ""
			}
			hitTracker(server, "CLIENT_ID", userId, r)
			hitTracker(server, "CLIENT_ID", userId, r)
			userNumber += 1
		}
	}

	response, _ := http.Get(server.URL + "/client/CLIENT_ID/referers")
	bytes, _ := ioutil.ReadAll(response.Body)
	response.Body.Close()

	expected, _ := json.Marshal(hits[:10])

	if string(bytes) != string(expected) {
		t.Errorf("Expected:\n%q\nGot:\n%q\n", string(expected), string(bytes))
		t.Fail()
	}
}

func hitTracker(server *httptest.Server, clientId string, userId string, referer string) (*http.Response, error) {
	client := &http.Client{}
	request, _ := http.NewRequest("GET", server.URL+"/client/"+clientId+"/tracker.gif?referer="+url.QueryEscape(referer), nil)
	request.Header.Set("X-Forwarded-For", "192.134.123.23")
	request.Header.Set("HTTP_REFERER", referer)
	request.AddCookie(&http.Cookie{Name: "sts", Value: userId})

	response, err := client.Do(request)
	response.Body.Close()
	return response, err
}
