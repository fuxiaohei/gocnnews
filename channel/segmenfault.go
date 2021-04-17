package channel

import (
	"log"
	"time"

	"gocnnews/model"

	"github.com/PuerkitoBio/goquery"
)

func Segmentfault() []*model.Article {
	// log.Println("[info]\t segmentfault : start")
	url := "https://segmentfault.com/blogs/newest"
	reader, err := getPage(url)
	if err != nil {
		log.Println("[error]\t segmentfault : request page fail,", err)
		return nil
	}
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Println("[error]\t segmentfault : read page fail,", err)
		return nil
	}
	list := doc.Find(".content-list-wrap .item-wrap div.content")
	if list.Length() == 0 {
		log.Println("[warn]\t segmentfault : find nothing")
		return nil
	}
	log.Printf("[info]\t segmentfault : find %d", list.Length())
	now := time.Now().Unix()
	var articles []*model.Article
	list.Each(func(_ int, sec *goquery.Selection) {
		title := sec.Find("h5").Text()
		url, _ := sec.Find("h5 a").Attr("href")
		if title == "" || url == "" {
			log.Println("[info]\t segmentfault : no title or url")
			return
		}
		content := sec.Find("p.excerpt").Text()
		a := &model.Article{
			URL:        "https://segmentfault.com" + url,
			Title:      title,
			Content:    content,
			AccessTime: now,
			From:       "segmentfault",
		}
		articles = append(articles, a)
	})
	return articles
}
