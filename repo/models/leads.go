package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Lead struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	OrgID      string             `bson:"orgid"`
	CampaignID string             `bson:"campaignid"`
	Status     int                `bson:"status"` //0:new, 1: synced
	Name       string             `bson:"name"`
	Email      string             `bson:"email"`
	Phone      string             `bson:"phone"`
	Message    string             `bson:"message"`
	Created    time.Time          `bson:"created"`
	Lastsync   time.Time          `bson:"lastsync"`
}
