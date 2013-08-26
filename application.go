package main

import (
	"encoding/json"
	"fmt"
	"github.com/hoisie/mustache"
	"github.com/nu7hatch/gouuid"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func createHandler(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		before := time.Now()
		handler(w, r)
		duration := time.Now().Sub(before)
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
	createHandler("/example/", exampleHandler)
	createHandler("/", homeHandler)
	serveFile("/scripts/jquery.sparkline.min.js", "./public/scripts/jquery.sparkline.min.js")
	serveFile("/scripts/tracker.js", "./public/scripts/tracker.js")
	serveFile("/styles/bootstrap.min.css", "./public/styles/bootstrap.min.css")
	serveFile("/images/ipad.jpg", "./public/images/ipad.jpg")

	if os.Getenv("PORT") == "" {
		http.ListenAndServe(":8080", nil)
	} else {
		http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	code := `<script type="text/javascript">
//<![CDATA[
  var sts = document.createElement('script'); sts.async = true;
  sts.src = '//statistic.li/scripts/tracker.js';
  var s = document.getElementsByTagName('script')[0]; s.parentNode.insertBefore(sts, s);
//]]>
</script>`

	io.WriteString(w, mustache.RenderFile("./views/home.mustache", map[string]string{"code": code}))
}

func exampleHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, mustache.RenderFile("./views/example.mustache", nil))
}

func clientHandler(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path[1:], "/")
	clientId := pathParts[1]

	if len(pathParts) == 2 {
		dash(clientId, w, r)
		return
	}

	if pathParts[2] == "dash" {
		dash(clientId, w, r)
	} else if pathParts[2] == "tracker.gif" {
		tracker(clientId, w, r)
	} else if pathParts[2] == "uniques" {
		uniques(clientId, w, r)
	} else if pathParts[2] == "referers" {
		referers(clientId, w, r)
	} else if pathParts[2] == "pages" {
		pages(clientId, w, r)
	} else if pathParts[2] == "generate" {
		generate(clientId)
	} else {
		io.WriteString(w, "Not Found")
	}
}

func generate(clientId string) {
	randomPage := func() string {
		i := rand.Intn(100)
		return "http://" + clientId + "/page" + strconv.Itoa(i) + ".html"
	}
	randomReferer := func() string {
		referers := []string{
			"http://google.co.uk/?q=search text",
			"http://www.bbc.co.uk/news/",
		}
		return referers[rand.Intn(len(referers))]
	}
	c := &ClientHit{
		ClientID: clientId,
		UserID:   generateNewUUID(),
		Referer:  randomReferer(),
		Page:     randomPage(),
	}
	c.Save()
}

func dash(clientId string, w http.ResponseWriter, r *http.Request) {
	context := map[string]interface{}{"clientId": clientId}
	io.WriteString(w, mustache.RenderFile("./views/dash.mustache", context))
}

func tracker(clientId string, w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	referer := r.URL.Query().Get("referer")
	clientHit := &ClientHit{
		ClientID: clientId,
		Page:     page,
		Referer:  referer,
	}

	cookie, err := r.Cookie("sts")
	if err == nil {
		clientHit.UserID = cookie.Value
	} else {
		userId := generateNewUUID()
		http.SetCookie(w, &http.Cookie{
			Name:    "sts",
			Value:   userId,
			Path:    "/",
			Expires: time.Date(3000, 1, 1, 1, 0, 0, 0, time.UTC),
		})
		clientHit.UserID = userId
	}

	clientHit.Save()
	w.Header().Set("Content-Type", "image/gif")
	w.Write(tracker_gif())
}

func generateNewUUID() string {
	u4, _ := uuid.NewV4()
	return u4.String()
}

func tracker_gif() []byte {
	return []byte{
		0x47, 0x49, 0x46, 0x38, 0x39, 0x61, 0x01, 0x00, 0x01, 0x00, 0x80, 0x00,
		0x00, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0x21, 0xf9, 0x04, 0x01, 0x00,
		0x00, 0x00, 0x00, 0x2c, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x01, 0x00,
		0x00, 0x02, 0x01, 0x44, 0x00, 0x3b,
	}
}

func uniques(clientId string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	uniques := GetUniques(clientId)
	b, _ := json.Marshal(uniques)
	w.Write(b)
}

func referers(clientId string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	topReferers := GetTopReferers(clientId)
	if len(topReferers) > 10 {
		topReferers = topReferers[:10]
	}
	b, _ := json.Marshal(topReferers)
	w.Write(b)
}

func pages(clientId string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	topPages := GetTopPages(clientId)
	if len(topPages) > 10 {
		topPages = topPages[:10]
	}
	b, _ := json.Marshal(topPages)
	w.Write(b)
}

func logError(part string, err error) {
	fmt.Println("["+part+"] ", err)
}
