package channel

import (
	"log"
	"strings"
	"time"

	"gocnnews/model"

	"github.com/PuerkitoBio/goquery"
)

func Reddit() []*model.Article {
	// log.Println("[info]\t reddit : start")
	url := "https://www.reddit.com/r/golang/new/"
	reader, err := getPage(url)
	if err != nil {
		log.Println("[error]\t reddit : request page fail,", err)
		return nil
	}
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Println("[error]\t reddit : read page fail,", err)
		return nil
	}
	list := doc.Find("div.link div.entry")
	if list.Length() == 0 {
		log.Println("[warn]\t reddit : find nothing")
		return nil
	}
	now := time.Now().Unix()
	var articles []*model.Article
	list.Each(func(_ int, sec *goquery.Selection) {
		title := sec.Find("p.title a.title").Text()
		url, _ := sec.Find("p.title a.title").Attr("href")
		if title == "" || url == "" {
			log.Println("[info]\t reddit : no title or url")
			return
		}
		a := &model.Article{
			URL:        url,
			Title:      title,
			Content:    "",
			AccessTime: now,
			From:       "reddit",
		}
		if !strings.HasPrefix(a.URL, "http") {
			a.URL = "https://reddit.com" + a.URL
		}
		articles = append(articles, a)
	})
	return articles
}
