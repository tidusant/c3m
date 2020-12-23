package models

import (
	"gopkg.in/mgo.v2/bson"
)

type BuildScript struct {
	ID           bson.ObjectId `bson:"_id,omitempty"`
	Data         string        `bson:"data"`
	IsRemove     bool          `bson:"isremove"`
	Status       int           `bson:"status"` //0: news, 1: building, 2: finish
	Created      int64         `bson:"created"`
	Modified     int64         `bson:"modified"`
	Retry        int           `bson:"retry"`
	ShopConfigs  ShopConfigs   `bson:"shopconfig"`
	ObjectId     string        `bson:"objectid"`
	Object       string        `bson:"object"`
	ShopId       string        `bson:"shopid"`
	TemplateCode string        `bson:"templatecode`
	Domain       string        `bson:"domain"`
}
type CommonData struct {
	Pages       []Page    `bson:"pages"`
	News        []News    `bson:"news"`
	NewsCats    []NewsCat `bson:"newscat"`
	Products    []Product `bson:"products"`
	ProductCats []ProdCat `bson:"productcats"`
	//Images Page `bson:""`
	//Albums Album `bson:"albums"`
}
type BuildConfig struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	ShopId      string        `bson:"shopid"`
	Domain      string        `bson:"domain"`
	Host        string        `bson:"host"` //0: news, 1: building, 2: finish
	Path        string        `bson:"path"`
	ApiUrl      string        `bson:"apiurl"`
	FTPUsername string        `bson:"ftpusername"`
	FTPPassword string        `bson:"ftppassword"`
}
