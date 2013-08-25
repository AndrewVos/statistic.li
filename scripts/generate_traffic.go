package main

import (
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

func hitTracker(clientId string, page string, referer string, cookies []*http.Cookie) []*http.Cookie {
	client := &http.Client{}

	page = url.QueryEscape(page)
	referer = url.QueryEscape(referer)
	request, _ := http.NewRequest("GET", "http://localhost:8080/client/"+clientId+"/tracker.gif?page="+page+"&referer="+referer, nil)
	if cookies != nil {
		for _, cookie := range cookies {
			request.AddCookie(cookie)
		}
	}
	response, _ := client.Do(request)
	if response != nil {
		response.Body.Close()
		return response.Cookies()
	}
	return nil
}

func randomPage() string {
	number := rand.Intn(100)
	return "http://client.com/pages/page" + strconv.Itoa(number) + ".html"
}

func randomlyHitTracker(clientId string) {
	cookies := hitTracker("client.com", randomPage(), "google.co.uk", nil)
	for i := 0; i < rand.Intn(10000); i++ {
		hitTracker("client.com", randomPage(), "google.co.uk", cookies)
	}
}

func main() {
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			randomlyHitTracker("client.com")
			defer wg.Done()
		}()
	}
	wg.Wait()

	for i := 0; i < 100000; i++ {
		randomlyHitTracker("client.com")
		time.Sleep(100 * time.Millisecond)
	}
}
