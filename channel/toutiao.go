package channel

import (
	"log"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocnnews/model"
)

func Toutiao() []*model.Article {
	// log.Println("[info]\t toutiao : start")
	url := "https://toutiao.io/tags/golang?type=latest"
	reader, err := getPage(url)
	if err != nil {
		log.Println("[error]\t toutiao : request page fail,", err)
		return nil
	}
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Println("[error]\t toutiao : read page fail,", err)
		return nil
	}
	list := doc.Find("div.post")
	if list.Length() == 0 {
		log.Println("[warn]\t toutiao : find nothing")
		return nil
	}
	now := time.Now().Unix()
	var articles []*model.Article
	list.Each(func(_ int, sec *goquery.Selection) {
		title := sec.Find("h3.title a").Text()
		url, _ := sec.Find("h3.title a").Attr("href")
		if title == "" || url == "" {
			log.Println("[info]\t toutiao : no title or url")
			return
		}
		content := sec.Find("p.summary").Text()
		a := &model.Article{
			URL:        "https://toutiao.io" + url,
			Title:      title,
			Content:    content,
			AccessTime: now,
			From:       "toutiao",
		}
		articles = append(articles, a)
	})
	return articles
}
