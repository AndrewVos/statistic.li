package application

import (
  "testing"
  "fmt"
  "strconv"
  "net/http"
  "io/ioutil"
)

func get(url string) (string) {
  res, _ := http.Get(url)
  body, _ := ioutil.ReadAll(res.Body)
  return string(body)
}

func TestApplicationTrackUniqueUsers(t *testing.T) {
  go Start()

  userIds := []string {
    "0.0.0.1",
    "0.0.0.2",
    "0.0.0.3",
    "0.0.0.4",
    "0.0.0.5",
  }
  numberOfUsers := len(userIds)
  clientId := "andrewvos.com"

  for i := 0; i < numberOfUsers*2; i++ {
    client := &http.Client {}
    request, _ := http.NewRequest("GET", "http://localhost:8080/client/" + clientId + "/tracker.gif", nil)
    request.Header.Set("X-Forwarded-For", userIds[i%numberOfUsers])
    client.Do(request)
  }

  expectedViews := `{"views": "` + strconv.Itoa(numberOfUsers) + `"}`
  views := get("http://localhost:8080/client/" + clientId + "/views")

  if views != expectedViews {
    fmt.Println(`Views was wrong, expected `, expectedViews, ` got `, views)
    t.Fail()
  }

}
