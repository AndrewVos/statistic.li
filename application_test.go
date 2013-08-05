package main

import (
  "testing"
  "strconv"
  "net/http"
  "io/ioutil"
  "net/http/httptest"
)

func TestUniqueViews(t *testing.T) {
  server := httptest.NewServer(http.HandlerFunc(clientHandler))
  defer server.Close()

  ipAddress := func(lastPart int) string { return "0.0.0." + strconv.Itoa(lastPart) }
  numberOfUsers := 5
  clientId := "andrewvos.com"

  for userHitCount := 0; userHitCount < numberOfUsers; userHitCount++ {
    client := &http.Client {}
    request, _ := http.NewRequest("GET", server.URL + "/client/" + clientId + "/tracker.gif", nil)
    request.Header.Set("X-Forwarded-For", ipAddress(userHitCount%numberOfUsers))

    for requestCount := 0; requestCount < 5; requestCount++ {
      response, _ := client.Do(request)
      gif,_ := ioutil.ReadAll(response.Body)
      if string(gif) != string(tracker_gif()) {
        t.Errorf("Expected:\n%q\nGot:\n%q\n", string(tracker_gif()), string(gif))
      }
      response.Body.Close()
    }
  }

  response,_ := http.Get(server.URL + "/client/" + clientId + "/views")
  responseBody,_ := ioutil.ReadAll(response.Body)
  views := string(responseBody)

  expectedContentType := "application/json"
  if response.Header["Content-Type"][0] != expectedContentType {
    t.Error("Expected a Content-Type of %q, not %q\n", expectedContentType, response.Header["Content-Type"][0])
  }

  expectedViews := `{"views":5}`
  if views != expectedViews {
    t.Error(`Views was wrong, expected `, expectedViews, ` got `, views)
  }
}
