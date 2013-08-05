package main

import (
  "testing"
  "fmt"
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

  for i := 0; i < numberOfUsers; i++ {
    client := &http.Client {}
    request, _ := http.NewRequest("GET", server.URL + "/client/" + clientId + "/tracker.gif", nil)
    request.Header.Set("X-Forwarded-For", ipAddress(i%numberOfUsers))

    for t := 0; t < 5; t++ {
      response, _ := client.Do(request)
      response.Body.Read(nil)
      response.Body.Close()
    }
  }

  response,_ := http.Get(server.URL + "/client/" + clientId + "/views")
  responseBody,_ := ioutil.ReadAll(response.Body)
  views := string(responseBody)

  expectedContentType := "application/json"
  if response.Header["Content-Type"][0] != expectedContentType {
    fmt.Printf("Expected a Content-Type of %q, not %q\n", expectedContentType, response.Header["Content-Type"][0])
    t.Fail()
  }

  expectedViews := fmt.Sprintf(`{"views":%d}`, numberOfUsers)
  if views != expectedViews {
    fmt.Println(`Views was wrong, expected `, expectedViews, ` got `, views)
    t.Fail()
  }
}
