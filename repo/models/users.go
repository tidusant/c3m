package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type UserSession struct {
	Session  string
	UserID   primitive.ObjectID
	Action   string
	AppName  string
	Params   string
	ShopID   primitive.ObjectID
	UserIP   string
	Group    string
	Username string
	Modules  map[string]bool
}
type User struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	User    string             `bson:"user"`
	Name    string             `bson:"name"`
	Email   string             `bson:"email"`
	Active  int32              `bson:"active"`
	Group   string             `bson:"group"`
	Modules string             `bson:"modules"`
	ShopId  primitive.ObjectID `bson:"shopid"`
}
type SFUser struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	SFUserID  string             `bson:"sfuserid"`
	Email     string             `bson:"email"`
	OrgID     string             `bson:"orgid"`
	LastLogin time.Time          `bson:"lastlogin"`
}

type UserConfig struct {
	Level     int `bson:"level"`
	MaxUpload int `bson:"maxupload"`
}

type UserLogin struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserId    primitive.ObjectID `bson:"userid"`
	ShopId    primitive.ObjectID `bson:"shopid"`
	Session   string             `bson:"session"`
	LastLogin time.Time          `bson:"last"`
	LoginIP   string             `bson:"ip"`
}
