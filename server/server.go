package server

import (
	"bytes"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gocnnews/model"
)

var majorWords = []string{"go", "golang", "docker", "grpc", "微服务", "k8s", "kubernetes"}

func init() {
	http.HandleFunc("/news", showArticles)
	http.HandleFunc("/mark", markArticle)
	http.HandleFunc("/unmark", unmarkArticle)
	http.HandleFunc("/generate", generateArticleNews)
}

func formatTime(t int64) string {
	ti := time.Unix(t, 0)
	return ti.Format("01/02 15:04")
}

func containsMajorWord(title string) bool {
	title = strings.ToLower(title)
	for _, w := range majorWords {
		if strings.Contains(title, w) {
			return true
		}
	}
	return false
}

func showArticles(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(405)
		return
	}
	var (
		t            = time.Now().Unix() - 86400
		status int64 = -99
	)
	if r.FormValue("mark") == "1" {
		status = 10
	}
	if r.FormValue("all") == "1" {
		t = 0
	}
	buf := bytes.NewBufferString(`<html><head>
	<meta charset="UTF-8" />
	<title>最近文章</title>
	</head><body>`)
	buf.WriteString("<h3>最近文章</h3><ul>")
	articles := model.ListArticles(t, status)
	for _, a := range articles {
		buf.WriteString("<li><p>")
		if a.Marked >= 10 {
			buf.WriteString(`<strong style="color:orange">` + a.Title + "</strong>")
		} else {
			if containsMajorWord(a.Title) {
				buf.WriteString(`<strong style="color:blue">` + a.Title + "</strong>")
			} else {
				buf.WriteString("<strong>" + a.Title + "</strong>")
			}
		}
		buf.WriteString(`&nbsp;&nbsp;&nbsp;<a href="` + a.URL + `" target="_blank"><small>链接</small></a>`)
		buf.WriteString("&nbsp;&nbsp;&nbsp;(" + a.From + ")")
		buf.WriteString("&nbsp;&nbsp;&nbsp;- " + formatTime(a.AccessTime))
		if a.Marked < 10 {
			buf.WriteString(`&nbsp;&nbsp;&nbsp;<a href="/mark?id=` + strconv.Itoa(a.ID) + "&hash=" + a.HashStr + `"><small>标记</small></a>`)
		} else {
			buf.WriteString(`&nbsp;&nbsp;&nbsp;<a href="/unmark?id=` + strconv.Itoa(a.ID) + "&hash=" + a.HashStr + `"><small>取消标记</small></a>`)
		}
		buf.WriteString("</p></li>")
	}
	buf.WriteString("</ul></body></html>")
	w.Header().Set("Content-Type", "text/html")
	w.Write(buf.Bytes())
}

func markArticle(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	hash := r.FormValue("hash")
	if id != "" && hash != "" {
		model.MarkArticle(id, hash, 10)
	}
	w.Header().Set("Location", "/news")
	w.WriteHeader(302)
}

func unmarkArticle(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	hash := r.FormValue("hash")
	if id != "" && hash != "" {
		model.MarkArticle(id, hash, 1)
	}
	w.Header().Set("Location", "/news")
	w.WriteHeader(302)
}

func generateArticleNews(w http.ResponseWriter, r *http.Request) {
	var (
		t               = time.Now().Unix() - 86400
		status    int64 = 10
		newsTitle       = time.Now().Format("GOCN每日新闻(2006-01-02)")
	)
	if r.FormValue("all") == "1" {
		t = 0
	}
	buf := bytes.NewBufferString(`<html><head>
	<meta charset="UTF-8" />
	<title>` + newsTitle + `</title>
	</head><body>`)
	buf.WriteString("<h3>" + newsTitle + "</h3>")
	articles := model.ListArticles(t, status)
	for i, a := range articles {
		buf.WriteString("<p>" + strconv.Itoa(i+1) + ".&nbsp;")
		buf.WriteString(a.Title)
		buf.WriteString(`&nbsp;&nbsp;<a href="` + a.URL + `" target="_blank">` + a.URL + `</a>`)
		buf.WriteString("</p>")
	}
	buf.WriteString("<p>&nbsp;&nbsp;</p>")
	buf.WriteString("<p>编辑：????</p>")
	buf.WriteString("<p>订阅新闻：http://tinyletter.com/gocn</p>")
	buf.WriteString("<p>GoCN归档：https://gocn.io/question/???</p>")
	buf.WriteString("</body></html>")
	w.Header().Set("Content-Type", "text/html")
	w.Write(buf.Bytes())
}

func Start() {
	err := http.ListenAndServe("0.0.0.0:49999", nil)
	if err != nil {
		log.Fatalln("[fatal]\t server : listen fail,", err)
	}
}
