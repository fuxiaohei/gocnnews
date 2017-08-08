package channel

import (
	"log"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocnnews/model"
)

func Oschina() []*model.Article {
	// log.Println("[info]\t oschina : start")
	url := "https://www.oschina.net/blog"
	reader, err := getPage(url)
	if err != nil {
		log.Println("[error]\t oschina : request page fail,", err)
		return nil
	}
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Println("[error]\t oschina : read page fail,", err)
		return nil
	}
	list := doc.Find("div#topsOfNew .item")
	if list.Length() == 0 {
		log.Println("[warn]\t oschina : find nothing")
		return nil
	}
	now := time.Now().Unix()
	var articles []*model.Article
	list.Each(func(_ int, sec *goquery.Selection) {
		title := sec.Find("a.blog-title-link h2").Text()
		url, _ := sec.Find("a.blog-title-link").Attr("href")
		if title == "" || url == "" {
			log.Println("[info]\t oschina : no title or url")
			return
		}
		content := sec.Find("section.blog-brief").Text()
		a := &model.Article{
			URL:        url,
			Title:      title,
			Content:    content,
			AccessTime: now,
			From:       "oschina",
		}
		articles = append(articles, a)
	})
	return articles
}
