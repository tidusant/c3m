package models

import (
	"gopkg.in/mgo.v2/bson"
)

type City struct {
	ID        bson.ObjectId       `bson:"_id,omitempty"`
	GhtkCode  string              `bson:"ghtkcode"`
	Source    string              `bson:"source"`
	Name      string              `bson:"name"`
	Pid       string              `bson:"pid"`
	Districts map[string]District `bson:districts`
	Code      string
}

type District struct {
	GhtkCode    string          `bson:"ghtkcode"`
	Source      string          `bson:"source"`
	Name        string          `bson:"name"`
	Pid         string          `bson:"pid"`
	IsPicked    bool            `bson:"is_picked"`
	IsDelivered bool            `bson:"is_delivered"`
	Wards       map[string]Ward `bson:wards`
	Code        string
}

type Ward struct {
	GhtkCode    string `bson:"ghtkcode"`
	Source      string `bson:"source"`
	Name        string `bson:"name"`
	Pid         string `bson:"pid"`
	IsPicked    bool   `bson:"is_picked"`
	IsDelivered bool   `bson:"is_delivered"`
	Code        string
}
