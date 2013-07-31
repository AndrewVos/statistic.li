package main

import (
  "io"
  "os"
  "fmt"
  "time"
  "strconv"
  "net"
  "github.com/fzzy/radix/redis"
  "github.com/hoisie/web"
)

func main() {
  web.Get("/client/(.*)/tracker.gif", tracker)
  web.Get("/client/(.*)/views", clientViews)
  web.Run("0.0.0.0:8080")
}

func storeClientHit(clientId string, userId string) {
  fmt.Println(clientId, userId)
  connection, err := redis.DialTimeout("tcp", "127.0.0.1:6379", time.Duration(10)*time.Second)
  if err == nil {
    now := time.Now().Unix()
    connection.Cmd("zrem", clientId, userId)
    connection.Cmd("zadd", clientId, now, userId)
  } else {
    fmt.Println(err)
  }
  defer connection.Close()
}

func tracker(ctx *web.Context, clientId string) {
  separator := " | "
  host, _, _ :=  net.SplitHostPort(ctx.Request.RemoteAddr)
  fmt.Printf(host)
  userId := host + separator + ctx.Request.UserAgent()
  fmt.Printf(userId)
  storeClientHit(clientId, userId)
  ctx.ContentType("gif")
  reader,_ := os.Open("./tracker.gif")
  defer reader.Close()
  io.Copy(ctx, reader)
}

func clientViews(clientId string) string {
  connection, err := redis.DialTimeout("tcp", "127.0.0.1:6379", time.Duration(10)*time.Second)
  if err == nil {
    now := time.Now().Unix()
    connection.Cmd("zremrangebyscore", clientId, 0, now - 60)
    result := connection.Cmd("zcount", clientId, "-inf", "+inf")
    currentCount,_ := result.Int()
    return `{"views": "` + strconv.Itoa(currentCount) + `"}`
  } else {
    fmt.Println(err)
    return `{"error": true}`
  }
}
