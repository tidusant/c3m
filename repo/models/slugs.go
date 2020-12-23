package models

import (
	"gopkg.in/mgo.v2/bson"
)

type Slug struct {
	ID           bson.ObjectId `bson:"_id,omitempty"`
	Slug         string        `bson:"slug"`
	Lang         string        `bson:"lang"`
	ObjectId     string        `bson:"objectid"`
	Object       string        `bson:"object"`
	View         string        `bson:"view"`
	ShopId       string        `bson:"shopid"`
	TemplateCode string        `bson:"templatecode`
	Domain       string        `bson:"domain"`
}
