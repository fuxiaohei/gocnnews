package model

import (
	"crypto/md5"
	"encoding/hex"
)

type Article struct {
	ID         int    `db:"id"`
	URL        string `db:"url"`
	From       string `db:"fromsite"`
	AccessTime int64  `db:"access"`
	Title      string `db:"title"`
	Content    string `db:"content"`
	Marked     int    `db:"marked"`
	HashStr    string `db:"hash"`
}

func (a *Article) Hash() string {
	key := []byte(a.From + a.URL)
	m := md5.New()
	m.Write(key)
	return hex.EncodeToString(m.Sum(nil))
}
