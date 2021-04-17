package channel

import (
	"log"
	"time"

	"gocnnews/model"

	"github.com/PuerkitoBio/goquery"
)

func Hexacosa() []*model.Article {
	// log.Println("[info]\t hexacosa : start")
	url := "http://goz.hexacosa.net/"
	reader, err := getPage(url)
	if err != nil {
		log.Println("[error]\t hexacosa : request page fail,", err)
		return nil
	}
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Println("[error]\t hexacosa : read page fail,", err)
		return nil
	}
	list := doc.Find("div.article")
	if list.Length() == 0 {
		log.Println("[warn]\t hexacosa : find nothing")
		return nil
	}
	now := time.Now().Unix()
	var articles []*model.Article
	list.Each(func(_ int, sec *goquery.Selection) {
		title := sec.Find("a.picks").Text()
		url, _ := sec.Find("a.picks").Attr("href")
		if title == "" || url == "" {
			log.Println("[info]\t hexacosa : no title or url")
			return
		}
		a := &model.Article{
			URL:        url,
			Title:      title,
			Content:    "",
			AccessTime: now,
			From:       "hexacosa",
		}
		articles = append(articles, a)
	})
	return articles
}
