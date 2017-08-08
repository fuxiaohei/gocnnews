package model

import (
	"database/sql"
	"log"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var (
	DB        *sqlx.DB
	createSQL = `CREATE TABLE news(id INTEGER PRIMARY KEY AUTOINCREMENT, title VARCHAR(100), url VARCHAR(100), content TEXT, hash VARCHAR(50),fromsite VARCHAR(50), access INT(15) default '0', marked INT(4) default '0')`
)

func Init() {
	db, err := sqlx.Open("sqlite3", "./gocnnews.db")
	if err != nil {
		log.Fatalln("[fatal]\tdb: open fail", err)
		return
	}
	if _, err = db.Exec(createSQL); err != nil {
		if err.Error() != "table news already exists" {
			log.Fatalln("[fatal]\tdb: create table fail,", err)
			return
		}
	}
	DB = db
}

var (
	selectHashSQL = "SELECT id FROM news WHERE hash = ?"
	insertSQL     = " INSERT INTO news(title,url,fromsite,content,hash,access) VALUES(?,?,?,?,?,?)"
	selectSQL     = "SELECT title,url,fromsite,access,marked,id,hash FROM news WHERE access >= ? and marked >= ? ORDER BY id DESC"
	updateMarkSQL = "UPDATE news SET marked = ? WHERE id = ? AND hash = ?"
)

func SaveArticles(articles []*Article) {
	for _, article := range articles {
		saveArticle(article)
	}
}

func saveArticle(article *Article) {
	row := DB.QueryRow(selectHashSQL, article.Hash())
	var id int
	if err := row.Scan(&id); err != nil {
		if err != sql.ErrNoRows {
			log.Println("[warn]\t db:query fail,", selectHashSQL, article.Hash(), err)
			return
		}
	}
	if id > 0 {
		// log.Println("[info]\tdb:saved,", id, article.Title)
		return
	}
	res, err := DB.Exec(insertSQL, strings.TrimSpace(article.Title), article.URL, article.From, article.Content, article.Hash(), time.Now().Unix())
	if err != nil {
		log.Println("[warn]\t db:insert fail,", insertSQL, article.Title, article.From, err)
		return
	}
	lastID, _ := res.LastInsertId()
	log.Println("[info]\t db:save new,", lastID, article.Title, article.From)
	return
}

func ListArticles(access, status int64) []*Article {
	rows, err := DB.Queryx(selectSQL, access, status)
	if err != nil {
		log.Println("[warn]\tdb:query fail,", selectSQL, err)
		return nil
	}
	var articles []*Article
	for rows.Next() {
		article := new(Article)
		if err = rows.StructScan(article); err != nil {
			log.Println("[warn]\t db:scan fail,", err)
			continue
		}
		articles = append(articles, article)
	}
	return articles
}

func MarkArticle(id, hash string, mark int) {
	_, err := DB.Exec(updateMarkSQL, mark, id, hash)
	if err != nil {
		log.Println("[warn]\t db:update fail,", updateMarkSQL, id, mark, err)
		return
	}
	log.Println("[info]\t db: mark,", id, hash, mark)
}
