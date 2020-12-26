package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	ShopId         string             `bson:"shopid"`
	EmployeeId     string             `bson:"employeeid"`
	CampaignId     string             `bson:"campaignid"`
	CampaignName   string
	ShipperId      string `bson:"shipperid"`
	ShipmentCode   string `bson:"shipmentcode"`
	Name           string
	Phone          string `bson:"phone"`
	OrderCount     int
	C              string `bson:"c"`
	City           string
	District       string
	Ward           string
	Address        string
	Note           string `bson:"note"`
	CusNote        string
	Email          string
	Status         string `bson:"status"`
	L              string `bson:"l"`
	Total          int    `bson:"total"`
	BaseTotal      int    `bson:"basetotal"`
	PartnerShipFee int    `bson:"partnershipfee"`
	ShipFee        int    `bson:"shipfee"`
	IsPaid         bool   `bson:"ispaid"`

	Items       []OrderItem `bson:"items"`
	Created     int64       `bson:"created"`
	Modified    int64       `bson:"modified"`
	Whookupdate int64       `bson:"whookupdate"`
	SearchIndex string      `bson:"searchindex"`
}

type OrderItem struct {
	ProdCode  string `bson:"prodcode"`
	Code      string `bson:"code"`
	CatName   string `bson:"catname"`
	Title     string `bson:"title"`
	Avatar    string `bson:"avatar"`
	BasePrice int    `bson:"baseprice"`
	Price     int    `bson:"price"`
	Num       int    `bson:"num"`
}

type OrderStatus struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	Title   string             `bson:"title"`
	Default bool               `bson:"default"`
	Finish  bool               `bson:"finish"`
	Export  bool               `bson:"export"`
	UserId  string             `bson:"userid"`
	ShopId  string             `bson:"shopid"`
	//Created       time.Time     `bson:"created"`
	//Modified      time.Time     `bson:"modified"`
	Color         string
	OrderCount    int
	PartnerStatus map[string][]string `bson:partnerstatus`
}

type ExportItem struct {
	ShopId   string `json:"prodcode"`
	Code     string `json:"code"`
	ItemCode string `json:"title"`
	Num      int    `json:"num"`
}
