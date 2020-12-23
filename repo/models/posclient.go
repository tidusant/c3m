package models

import (
	"gopkg.in/mgo.v2/bson"
)

type PosClient struct {
	ID     bson.ObjectId `bson:"_id,omitempty"`
	Name   string        `bson:"name"`
	ShopId string        `bson:"shopid"`
	UserId string        `bson:"userid"`
	IsSync bool          `bson:"issync"`
}
