package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Product struct {
	ID         primitive.ObjectID      `bson:"_id,omitempty"`
	Code       string                  `bson:"code"`
	UserId     string                  `bson:"userid"`
	ShopId     string                  `bson:"shopid"`
	CatId      string                  `bson:"catid"`
	Langs      map[string]*ProductLang `bson:"langs"`
	Properties []ProductProperty       `bson:"properties"`
	Options    []ProductOption         `bson:"options"`
	Status     string                  `bson:"status"`
	Publish    bool                    `bson:"publish"`
	Main       bool                    `bson:"main"`
	IsSync     bool                    `bson:"issync"`
	Created    time.Time               `bson:"created"`
	Modified   time.Time               `bson:"modified"`
	LangLinks  map[string]string       `bson:"langlinks"`
}
type ProductOption struct {
	Name  string   `bson:"name"`
	Items []string `bson:"items"`
}
type ProductProperty struct {
	Name      string `bson:"name"`
	Code      string `bson:"code"`
	Price     int    `bson:"price"`
	BasePrice int    `bson:"baseprice"`
	Stock     int    `bson:"stock"`
}
type ProductLang struct {
	Avatar          string   `bson:"avatar"`
	Title           string   `bson:"title"`
	Name            string   `bson:"name"`
	Content         string   `bson:"content"`
	Description     string   `bson:"description"`
	Slug            string   `bson:"slug"`
	Price           int      `bson:"price"`
	BasePrice       int      `bson:"baseprice"`
	DiscountPrice   int      `bson:"discountprice"`
	PercentDiscount int      `bson:"percentdiscount"`
	Unit            string   `bson:"unit"`
	Images          []string `bson:"images"`
}

type ProdCat struct {
	ID        primitive.ObjectID   `bson:"_id,omitempty"`
	Code      string               `bson:"code"`
	UserId    string               `bson:"userid"`
	ShopId    string               `bson:"shopid"`
	Created   time.Time            `bson:"created"`
	Langs     map[string]*PageLang `bson:"langs"`
	IsSync    bool                 `bson:"issync"`
	Main      bool                 `bson:"main"`
	LangLinks map[string]string    `bson:"langlinks"`
}
