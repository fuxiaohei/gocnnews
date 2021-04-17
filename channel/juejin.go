package channel

import (
	"bytes"
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"gocnnews/model"

	"github.com/tidwall/gjson"
)

func Juejin() []*model.Article {
	// log.Println("[info]\t juejin : start")
	data, err := getJuejinAPI("https://api.juejin.cn/recommend_api/v1/article/recommend_all_feed")
	if err != nil {
		log.Println("[error]\t juejin : request page fail,", err)
		return nil
	}

	list := gjson.ParseBytes(data).Get("data").Array()
	if len(list) == 0 {
		log.Println("[warn]\t juejin : find nothing")
		return nil
	}
	log.Printf("[info]\t juejin : find %d", len(list))
	now := time.Now().Unix()
	var articles []*model.Article
	for _, item := range list {
		title := item.Get("item_info.article_info.title").String()
		content := item.Get("item_info.article_info.brief_content").String()
		articleid := item.Get("item_info.article_id").String()
		a := &model.Article{
			URL:        "https://juejin.im/post/" + articleid,
			Title:      title,
			Content:    content,
			AccessTime: now,
			From:       "juejin",
		}
		articles = append(articles, a)
	}
	return articles
	/*
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
		return articles*/
}

func getJuejinAPI(url string) ([]byte, error) {
	payload := strings.NewReader(`{
		"id_type": 2,
		"client_type": 2608,
		"sort_type": 300,
		"cursor": "0",
		"limit": 20
	}`)
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.25 Safari/537.36")
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: false,
			},
		},
		Timeout: 20 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	buf := bytes.NewBuffer(nil)
	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
