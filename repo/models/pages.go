package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

//Page ...
type Page struct {
	ID           bson.ObjectId       `bson:"_id,omitempty"`
	Code         string              `bson:"code"`
	UserID       string              `bson:"userid"`
	ShopID       string              `bson:"shopid"`
	TemplateCode string              `bson:"templatecode"`
	Langs        map[string]PageLang `bson:"langs"`
	Created      time.Time           `bson:"created"`
	Modified     time.Time           `bson:"modified"`
	LangLinks    []LangLink          `bson:"langlinks"`
	Blocks       []PageBlock         `bson:blocks`
	Seo          string              `bson:"seo"`
}

//NewsLang ...
type PageBlock struct {
	Name  string          `bson:name`
	Items []PageBlockItem `bson:"items"`
}

type PageBlockItem struct {
	Key   string            `bson:key`
	Type  string            `bson:type`
	Value map[string]string `bson:value`
}

//NewsLang ...
type PageLang struct {
	Title       string `bson:"title"`
	Name        string `bson:"name"`
	Content     string `bson:"content"`
	Description string `bson:"description"`
	Slug        string `bson:"slug"`
	Catname     string
}

type PageView struct {
	Title       string
	Slug        string
	Code        string
	Description string
	Content     string
	PageType    string
	Pagename    string
	Templ       string
	AltPagename string
	CatName     string
	CatSlug     string
	LangLinks   []LangLink
	Lang        string
}

type LangLink struct {
	Href string `bson:"href"`
	Code string `bson:"code"`
	Flag string `bson:"flag"`
	Name string `bson:"name"`
}
