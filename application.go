package main

import (
  "io"
  "os"
  "fmt"
  "time"
  "strings"
  "sort"
  "net/http"
  "encoding/json"
  "labix.org/v2/mgo"
  "labix.org/v2/mgo/bson"
  "github.com/hoisie/mustache"
  "github.com/nu7hatch/gouuid"
)

type ClientHit struct {
  ClientID string
  UserID string
  Date time.Time
  Referer string
}

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

func getConnection() (*mgo.Session, error) {
  uri := os.Getenv("MONGOHQ_URL")
  if uri == "" {
    uri = ":27017"
  }
  session, err := mgo.Dial(uri)
  if err != nil { return nil, err }
  return session, nil
}

func storeClientHit(clientId string, userId string, referer string) {
  session, err := getConnection()
  defer session.Close()
  if err != nil {
    logError("mongo", err)
    return
  }

  collection := session.DB("").C("ClientHits")
  err = collection.Insert(&ClientHit {
    ClientID: clientId,
    UserID: userId,
    Date: time.Now(),
    Referer: referer,
  })
  if err != nil { logError("mongo", err) }
}

func getUniqueViews(clientId string) (int, error) {
  session, err := getConnection()
  defer session.Close()
  if err != nil {
    logError("mongo", err)
    return 0, err
  }

  collection := session.DB("").C("ClientHits")
  after := time.Now().Add(-5 * time.Minute)

  query := collection.Find(bson.M{"clientid": clientId, "date": bson.M{"$gte": after}})
  var distinctUserIds []string
  query.Distinct("userid", &distinctUserIds)
  return len(distinctUserIds), nil
}

func getTopReferers(clientId string) (PageHitCounts, error) {
  session, err := getConnection()
  defer session.Close()
  if err != nil {
    logError("mongo", err)
    return nil, err
  }

  collection := session.DB("").C("ClientHits")
  after := time.Now().Add(-5 * time.Minute)
  query := collection.Find(bson.M{"clientid": clientId, "date": bson.M{"$gte": after}})
  var hits []ClientHit
  query.All(&hits)

  countedPages := make(map[string] bool)
  pageCounts := make(map[string] int)

  for _,clientHit := range hits {
    r := clientHit.Referer
    if r == "" { r = "(direct)" }
    if countedPages[clientHit.UserID + r] != true {
      countedPages[clientHit.UserID + r] = true
      if _,ok := pageCounts[r]; ok != true {
        pageCounts[r] = 0
      }
      pageCounts[r] += 1
    }
  }
  var pageHitCounts PageHitCounts
  for referer, count := range pageCounts {
    pageHitCounts = append(pageHitCounts, &PageHitCount{Referer: referer, Count: count})
  }
  sort.Sort(pageHitCounts)
  return pageHitCounts, nil
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
  } else if pathParts[2] == "referers" {
    referers(clientId, w, r)
  } else {
    io.WriteString(w, "Not Found")
  }
}

func dash(clientId string, w http.ResponseWriter, r *http.Request) {
  context := map[string]interface{}{"clientId": clientId}
  io.WriteString(w, mustache.RenderFile("./views/dash.mustache", context))
}

func tracker(clientId string, w http.ResponseWriter, r *http.Request) {
  referer := "(direct)"
  for h,v := range r.Header {
    if strings.ToLower(h) == "http_referer" {
      referer = v[0]
    }
  }

  cookie, err := r.Cookie("sts")
  if err == nil {
    storeClientHit(clientId, cookie.Value, referer)
  } else {
    userId := generateNewUUID()
    http.SetCookie(w, &http.Cookie{
      Name: "sts",
      Value: userId,
      Path: "/",
      Expires: time.Date(3000, 1, 1, 1, 0, 0, 0, time.UTC),
    })
    storeClientHit(clientId, userId, referer)
  }

  w.Header().Set("Content-Type", "image/gif")
  w.Write(tracker_gif())
}

func generateNewUUID() string {
  u4, _ := uuid.NewV4()
  return u4.String()
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
  w.Header().Set("Content-Type", "application/json")

  result, err := getUniqueViews(clientId)

  if err != nil {
    logError("mongo", err)
    io.WriteString(w, `{"error": true}`)
    return
  }

  response,_ := json.Marshal(map[string] int {
    "views": result,
  })
  io.WriteString(w, string(response))
}

func referers(clientId string, w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")

  topReferers, err := getTopReferers(clientId)
  if err != nil { logError("mongo", err) }
  if err != nil || topReferers == nil {
    io.WriteString(w, `[]`)
    return
  }

  if len(topReferers) > 10 {
    topReferers = topReferers[:10]
  }
  b,_ := json.Marshal(topReferers)
  w.Write(b)
}

func logError(part string, err error) {
  fmt.Println("[" + part + "] ", err)
}
