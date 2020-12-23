package models

import (
	"gopkg.in/mgo.v2/bson"
)

type Whook struct {
	ID bson.ObjectId `bson:"_id,omitempty"`

	Name    string `bson:"name"`
	Action  string `bson:"action"`
	ShopID  string `bson:"shopid"`
	Created int    `bson:"created"`
	Data    string `bson:"data"`
	Status  int32  `bson:"status"`
}
