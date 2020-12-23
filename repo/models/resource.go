package models

import (
	"gopkg.in/mgo.v2/bson"
)

//Page ...
type Resource struct {
	ID           bson.ObjectId     `bson:"_id,omitempty"`
	ShopID       string            `bson:"shopid"`
	Key          string            `bson:"key"`
	Type         string            `bson:"type"`
	TemplateCode string            `bson:"templatecode"`
	Value        map[string]string `bson:"value"`
}
