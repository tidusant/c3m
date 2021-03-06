package models

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type LPTemplate struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	UserID      primitive.ObjectID `bson:"userid"`
	Status      int                `bson:"status"` //-2: delete, -1: reject,0:local, 1: approved and publish, 2: waiting for approve
	Description string             `bson:"description"`
	Path        string             `bson:"path"`
	ScreenShot  string             `bson:"screenshot"`
	Name        string             `bson:"name"`
	Viewed      int                `bson:"viewed"`
	Installed   int                `bson:"installed"`
	Created     time.Time          `bson:"created"`
	User        string             `bson:"user"`
}
type LPBuildData struct {
	Session        string
	LPPath         string
	Favicon        string
	CampID         string
	OrgID          string
	SuccessTitle   string
	SuccessMessage string
	ErrorTitle     string
	ErrorMessage   string
}
type LandingPage struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	Path           string             `bson:"path"`
	Favicon        string             `bson:"favicon"`
	UserID         primitive.ObjectID `bson:"userid"`
	LPTemplateID   primitive.ObjectID `bson:"lptemplateid"`
	CustomHost     bool               `bson:"customhost"`
	DomainName     string             `bson:"domainname"`
	FTPHost        string             `bson:"ftphost"`
	FTPUser        string             `bson:"ftpuser"`
	FTPPass        string             `bson:"ftppass"`
	Content        string             `bson:"content"`
	CampaignID     string             `bson:"campaignid"`
	OrgID          string             `bson:"orgid"`
	SuccessTitle   string             `bson:"successtitle"`
	SuccessMessage string             `bson:"successmessage"`
	ErrorTitle     string             `bson:"errortitle"`
	ErrorMessage   string             `bson:"errormessage"`
	SFUserID       string             `bson:sfuserid`
	Viewed         int                `bson:"viewed"`
	Submitted      int                `bson:'submitted'`
	Created        time.Time          `bson:"created"`
	LastBuild      time.Time          `bson:"lastbuild"`
	Modified       time.Time          `bson:"modified"`
}
type LPTPL4User struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Name       string             `bson:"name"`
	Viewed     int                `bson:"viewed"`
	Installed  int                `bson:"installed"`
	Designer   string             `bson:"user"`
	ScreenShot string             `bson:"screenshot"`
	Created    time.Time          `bson:"created"`
}
type Template struct {
	ID           bson.ObjectId `bson:"_id,omitempty"`
	Code         string        `bson:"code"`
	UserID       string        `bson:"userid"`
	Status       int           `bson:"status"` //-2: delete, -1: reject, 1: approved and publish, 2: pending, 3: approved but not publish
	Title        string        `bson:"title"`
	Description  string        `bson:"description"`
	Viewed       int           `bson:"viewed"`
	InstalledIDs []string      `bson:"installedid"`
	ActiveIDs    []string      `bson:"activedid"`
	Configs      string        `bson:"configs"`
	Resources    string        `bson:"resources"`
	Pages        string        `bson:"pages"`
	Avatar       string        `bson:"avatar"`
	Created      time.Time     `bson:"created"`
	Modified     time.Time     `bson:"modified"`
}

type TemplateSubmit struct {
	Code  string `bson:"code"`
	Title string `bson:"title"`
}

//News ...
type TemplateConfig struct {
	ID           bson.ObjectId `bson:"_id,omitempty"`
	TemplateCode string        `bson:"templatecode"`
	ShopID       string        `bson:"shopid"`
	Key          string        `bson:"key"`
	Type         string        `bson:"type"`
	Value        string        `bson:"value"`
}

type TemplateLang struct {
	ID           bson.ObjectId `bson:"_id,omitempty"`
	TemplateCode string        `bson:"templatecode"`
	Lang         string        `bson:"lang"`
	ShopID       string        `bson:"shopid"`
	Key          string        `bson:"key"`
	Value        string        `bson:"value"`
}

type TemplateViewData struct {
	PageName     string
	Siteurl      string
	Data         map[string]json.RawMessage
	TemplatePath string
	Templateurl  string
	Imageurl     string
	Pages        map[string]string
	Resources    map[string]string
	Configs      map[string]string
}
