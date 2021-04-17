package channel

import (
	"log"
	"time"

	"gocnnews/model"

	"github.com/PuerkitoBio/goquery"
)

func Hacknews() []*model.Article {
	// log.Println("[info]\t hacknews : start")
	url := "https://news.ycombinator.com/newest"
	reader, err := getPage(url)
	if err != nil {
		log.Println("[error]\t hacknews : request page fail,", err)
		return nil
	}
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Println("[error]\t hacknews : read page fail,", err)
		return nil
	}
	list := doc.Find("tr.athing")
	if list.Length() == 0 {
		log.Println("[warn]\t hacknews : find nothing")
		return nil
	}
	log.Printf("[info]\t hacknews : find %d", list.Length())
	now := time.Now().Unix()
	var articles []*model.Article
	list.Each(func(_ int, sec *goquery.Selection) {
		title := sec.Find("td.title a.storylink").Text()
		url, _ := sec.Find("td.title a.storylink").Attr("href")
		if title == "" || url == "" {
			log.Println("[info]\t hacknews : no title or url")
			return
		}
		a := &model.Article{
			URL:        url,
			Title:      title,
			Content:    "",
			AccessTime: now,
			From:       "hacknews",
		}
		articles = append(articles, a)
	})
	return articles
}
