package channel

import (
	"log"
	"strings"
	"time"

	"gocnnews/model"

	"github.com/PuerkitoBio/goquery"
)

func Golangnews() []*model.Article {
	// log.Println("[info]\t golangnews : start")
	url := "https://golangnews.com/"
	reader, err := getPage(url)
	if err != nil {
		log.Println("[error]\t golangnews : request page fail,", err)
		return nil
	}
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Println("[error]\t golangnews : read page fail,", err)
		return nil
	}
	list := doc.Find("li.story")
	if list.Length() == 0 {
		log.Println("[warn]\t golangnews : find nothing")
		return nil
	}
	log.Printf("[info]\t golangnews : find %d", list.Length())
	now := time.Now().Unix()
	var articles []*model.Article
	list.Each(func(_ int, sec *goquery.Selection) {
		title := sec.Find("h3 a.name").Text()
		url, _ := sec.Find("h3 a.name").Attr("href")
		if title == "" || url == "" {
			log.Println("[info]\t golangnews : no title or url")
			return
		}
		a := &model.Article{
			URL:        url,
			Title:      title,
			Content:    "",
			AccessTime: now,
			From:       "golangnews",
		}
		if !strings.HasPrefix(a.URL, "http") {
			a.URL = "https://golangnews.com" + a.URL
		}
		articles = append(articles, a)
	})
	return articles
}
