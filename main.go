package main

import (
  "io"
  "os"
  "fmt"
  "time"
  "strconv"
  "strings"
  "github.com/garyburd/redigo/redis"
  "github.com/hoisie/web"
  "github.com/soveran/redisurl"
)

func main() {
  web.Get("/client/(.*)/tracker.gif", tracker)
  web.Get("/client/(.*)/views", clientViews)

  if os.Getenv("PORT") == "" {
    web.Run(":8080")
  } else {
    web.Run(":" + os.Getenv("PORT"))
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
  connection := getConnection()
  defer connection.Close()
  if connection != nil {
    now := time.Now().Unix()
    connection.Do("ZREM", clientId, userId)
    connection.Do("ZADD", clientId, now, userId)
  }
}

func tracker(ctx *web.Context, clientId string) {
  separator := " | "
  host := strings.Join(ctx.Request.Header["X-Forwarded-For"], ",")
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
  connection := getConnection()
  defer connection.Close()
  if connection != nil {
    now := time.Now().Unix()
    connection.Do("ZREMRANGEBYSCORE", clientId, 0, now - 60)
    result,err := redis.Int(connection.Do("ZCOUNT", clientId, "-inf", "+inf"))

    if err == nil {
      return `{"views": "` + strconv.Itoa(result) + `"}`
    } else {
      fmt.Println(err)
      return `{"views": "0"}`
    }
  } else {
    return `{"error": true}`
  }
}
