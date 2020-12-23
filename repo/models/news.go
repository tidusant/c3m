package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

//News ...
type News struct {
	ID        bson.ObjectId        `bson:"_id,omitempty"`
	Avatar    string               `bson:"avatar"`
	UserID    string               `bson:"userid"`
	ShopID    string               `bson:"shopid"`
	CatIDs    []string             `bson:"catids"`
	Langs     map[string]*PageLang `bson:"langs"`
	Created   time.Time            `bson:"created"`
	Modified  time.Time            `bson:"modified"`
	Publish   bool                 `bson:"publish"`
	Home      bool                 `bson:"home"`
	Feature   bool                 `bson:"feature"`
	Code      string               `bson:"code"`
	LangLinks []LangLink           `bson:"langlinks"`
}

//NewsCat ...
type NewsCat struct {
	ID       bson.ObjectId        `bson:"_id,omitempty"`
	Avatar   string               `bson:"avatar"`
	UserId   string               `bson:"userid"`
	ShopId   string               `bson:"shopid"`
	Created  time.Time            `bson:"created"`
	Langs    map[string]*PageLang `bson:"langs"`
	Code     string               `bson:"code"`
	ParentId string               `bson:"parentid"`
	Publish  bool                 `bson:"publish"`
	Home     bool                 `bson:"home"`
	Feature  bool                 `bson:"feature"`

	LangLinks []LangLink `bson:"langlinks"`
}
