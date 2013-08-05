package main

import (
  "io"
  "os"
  "fmt"
  "time"
  "strings"
  "net/http"
  "encoding/json"
  "github.com/garyburd/redigo/redis"
  "github.com/soveran/redisurl"
  "github.com/hoisie/mustache"
)

func createHandler(pattern string, handler func(http.ResponseWriter, *http.Request)) {
  http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
    before := time.Now()
    handler(w, r)
    duration :=  time.Now().Sub(before)
    fmt.Printf("%v %v - %v %v - %v\n", time.Now().Format("2006/01/02 15:04:05"), r.RemoteAddr, r.Method, r.URL.Path, duration)
  })
}

func serveFile(pattern string, filename string) {
  http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, filename)
  })
}

func Start() {
  createHandler("/client/", clientHandler)
  serveFile("/scripts/dash-updater.js", "./public/scripts/dash-updater.js")

  if os.Getenv("PORT") == "" {
    http.ListenAndServe(":8080", nil)
  } else {
    http.ListenAndServe(":" + os.Getenv("PORT"), nil)
  }
}

func getConnection() redis.Conn {
  if os.Getenv("REDISCLOUD_URL") != "" {
    os.Setenv("REDIS_URL", os.Getenv("REDISCLOUD_URL"))
  }

  if redisUrl := os.Getenv("REDIS_URL"); redisUrl != "" {
    connection, err := redisurl.Connect()
    if err != nil {
      fmt.Println(err)
      return nil
    }
    return connection
  } else {
    connection, err := redis.Dial("tcp", ":6379")
    if err != nil {
      fmt.Println(err)
      return nil
    }
    return connection
  }
}

func storeClientHit(clientId string, userId string) {
  fmt.Println("client hit: ", clientId, userId)
  connection := getConnection()
  defer connection.Close()
  if connection != nil {
    now := time.Now().Unix()
    connection.Do("ZADD", clientId, now, userId)
  }
}

func clientHandler(w http.ResponseWriter, r *http.Request) {
  pathParts := strings.Split(r.URL.Path[1:], "/")
  clientId := pathParts[1]
  if pathParts[2] == "dash" {
    dash(clientId, w, r)
  } else if pathParts[2] == "tracker.gif" {
    tracker(clientId, w, r)
  } else if pathParts[2] == "views" {
    views(clientId, w, r)
  } else {
    io.WriteString(w, "Not Found")
  }
}

func dash(clientId string, w http.ResponseWriter, r *http.Request) {
  context := map[string]interface{}{"clientId": clientId}
  io.WriteString(w, mustache.RenderFile("./views/dash.mustache", context))
}

func tracker(clientId string, w http.ResponseWriter, r *http.Request) {
  separator := " | "
  host := strings.Join(r.Header["X-Forwarded-For"], ",")
  fmt.Printf(host)
  userId := host + separator + r.UserAgent()
  fmt.Printf(userId)
  storeClientHit(clientId, userId)
  w.Header().Set("Content-Type", "image/gif")
  w.Write(tracker_gif())
}

func tracker_gif() []byte {
	return []byte{
    0x47,0x49,0x46,0x38,0x39,0x61,0x01,0x00,0x01,0x00,0x80,0x00,
    0x00,0x00,0x00,0x00,0xff,0xff,0xff,0x21,0xf9,0x04,0x01,0x00,
    0x00,0x00,0x00,0x2c,0x00,0x00,0x00,0x00,0x01,0x00,0x01,0x00,
    0x00,0x02,0x01,0x44,0x00,0x3b,
	}
}

func views(clientId string, w http.ResponseWriter, r *http.Request) {
  connection := getConnection()
  defer connection.Close()
  w.Header().Set("Content-Type", "application/json")
  if connection != nil {
    now := time.Now().Unix()
    connection.Do("ZREMRANGEBYSCORE", clientId, 0, now - 300)
    result,err := redis.Int(connection.Do("ZCOUNT", clientId, "-inf", "+inf"))

    if err != nil { fmt.Println(err) }
    response,_ := json.Marshal(map[string] int {
      "views": result,
    })
    io.WriteString(w, string(response))
  } else {
    io.WriteString(w, `{"error": true}`)
  }
}
