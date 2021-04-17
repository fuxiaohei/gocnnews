package channel

import (
	"log"
	"time"

	"gocnnews/model"

	"github.com/PuerkitoBio/goquery"
)

func Juejin() []*model.Article {
	// log.Println("[info]\t juejin : start")
	reader, err := getPage("https://juejin.im/zhuanlan/all")
	if err != nil {
		log.Println("[error]\t juejin : request page fail,", err)
		return nil
	}
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Println("[error]\t juejin : read page fail,", err)
		return nil
	}
	list := doc.Find("div.column-entry")
	if list.Length() == 0 {
		log.Println("[warn]\t juejin : find nothing")
		return nil
	}
	now := time.Now().Unix()
	var articles []*model.Article
	list.Each(func(_ int, sec *goquery.Selection) {
		title := sec.Find("div.abstract-row a.title").Text()
		url, _ := sec.Find("div.abstract-row a.title").Attr("href")
		if title == "" || url == "" {
			log.Println("[info]\t juejin : no title or url")
			return
		}
		content := sec.Find("div.abstract-row a.abstract").Text()
		a := &model.Article{
			URL:        "https://juejin.im" + url,
			Title:      title,
			Content:    content,
			AccessTime: now,
			From:       "juejin",
		}
		articles = append(articles, a)
	})
	return articles
}
