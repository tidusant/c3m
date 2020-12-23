package models

import (
	"gopkg.in/mgo.v2/bson"
)

type Session struct {
	ID      bson.ObjectId `bson:"_id,omitempty"`
	Uid     string        `bson:"uid"`
	Expired int64         `bson:"expired"`
}
