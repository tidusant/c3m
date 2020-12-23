package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type FBWhook struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	Method      string        `bson:"method"`
	ContentType string        `bson:"contenttype"`
	URL         string        `bson:"url"`
	Created     time.Time     `bson:"created"`
	Data        string        `bson:"data"`
}
type FBComment struct {
	ID           bson.ObjectId `bson:"_id,omitempty"`
	Object       string        `bson:"object"`
	ObjectID     string        `bson:"id"`
	Field        string        `bson:"field"`
	FromID       string        `bson:"fromid"`
	FromName     string        `bson:"fromname"`
	CommentID    string        `bson:"comment_id"`
	PageID       string        `bson:"page_id"`
	PostID       string        `bson:"post_id"`
	ParentID     string        `bson:"parent_id"`
	PermalinkURL string        `bson:"permalink_url"`
	Message      string        `bson:"message"`
	Created      time.Time     `bson:"updated_time"`
}
type FBConversation struct {
	ID              bson.ObjectId `bson:"_id,omitempty"`
	Object          string        `bson:"object"`
	ObjectID        string        `bson:"id"`
	PageID          string        `bson:"page_id"`
	ThreadID        string        `bson:"thread_id"`
	ScopedThreadKey string        `bson:"scoped_thread_key"`
	Field           string        `bson:"field"`
	Created         time.Time     `bson:"created_time"`
}
