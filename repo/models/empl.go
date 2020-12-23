package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Employee struct {
	ID       bson.ObjectId `bson:"_id,omitempty"`
	ShopId   string        `bson:"shopid"`
	Password string        `bson:"password"`
	Name     string        `bson:"name"`
	Sex      string        `bson:"sex"`
	Phone    string        `bson:"phone"`
	Address  string        `bson:"address"`
	Avatar   string        `bson:"avatar"`
	Note     string        `bson:"note"`
	Email    string        `bson:"email"`
	IsSync   bool          `bson:"issync"`
	Created  time.Time     `bson:"created"`
	Modified time.Time     `bson:"modified"`
}
