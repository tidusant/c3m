package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Shipper struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	ShopId      string        `bson:"shopid"`
	UserId      string        `bson:"userid"`
	Name        string        `bson:"name"`
	Color       string        `bson:"color"`
	PartnerCode string        `bson:"partnercode"`
	Default     bool          `bson:"default"`
	Created     time.Time     `bson:"created"`
	Modified    time.Time     `bson:"modified"`
}
