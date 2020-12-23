package cuahang

import (
	"github.com/tidusant/c3m-common/c3mcommon"

	"github.com/tidusant/chadmin-repo/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
	"log"
	"time"
)

func GetAllCampaigns(shopid primitive.ObjectID) []models.Campaign {
	col := db.Collection("addons_campaigns")
	var rs []models.Campaign
	cond := bson.M{"shopid": shopid.Hex()}
	cursor, err := col.Find(ctx, cond)
	if err = cursor.All(ctx, &rs); err != nil {
		log.Fatal(err)
	}
	c3mcommon.CheckError("Getall campaign", err)
	return rs
}

func GetCampaignsByRange(shopid primitive.ObjectID, start time.Time, end time.Time) []models.Campaign {
	col := db.Collection("addons_campaigns")
	var rs []models.Campaign
	cond := bson.M{"shopid": shopid.Hex(), "start": bson.M{"$lt": end}, "end": bson.M{"$gt": start}}
	cursor, err := col.Find(ctx, cond)
	if err = cursor.All(ctx, &rs); err != nil {
		log.Fatal(err)
	}
	c3mcommon.CheckError("Getall campaign by range", err)
	return rs
}

//
//func GetCampaignByID(shopid, ID string) models.Campaign {
//	col := db.Collection("addons_campaigns")
//	var rs models.Campaign
//	cond := bson.M{"shopid": shopid, "_id": bson.ObjectIdHex(ID)}
//	cursor,err := col.Find(ctx,cond)
//	if err = cursor.All(ctx, &rs); err != nil {     log.Fatal(err) }
//	c3mcommon.CheckError("Get campaign by id", err)
//	return rs
//}
//func GetOrderStatusMap(shopid string) map[string]models.OrderStatus {
//	statsmap := make(map[string]models.OrderStatus)
//	stats := GetAllOrderStatus(shopid)
//	for _, stat := range stats {
//		statsmap[stat.ID.Hex()] = stat
//	}
//	return statsmap
//}
//func GetCampaignDetailByID(shopid string, camp models.Campaign) models.Campaign {
//	col := db.Collection("addons_orders")
//
//	rs := make(map[string]models.CampaignStatusDetail)
//	var ords []models.Order
//	cond := bson.M{"shopid": shopid, "campaignid": camp.ID.Hex()}
//
//	cursor,err := col.Find(ctx,cond)
//	if err = cursor.All(ctx, &rs); err != nil {     log.Fatal(err) }
//
//	stats := GetAllOrderStatus(shopid)
//	statsmap := make(map[string]models.OrderStatus)
//	for _, stat := range stats {
//		statsmap[stat.ID.Hex()] = stat
//	}
//	if c3mcommon.CheckError("Get detail campaign by id", err) {
//
//		for _, ord := range ords {
//			var detail models.CampaignStatusDetail
//			if val, ok := rs[ord.Status]; ok {
//				detail = val
//			}
//
//			detail.Base += ord.BaseTotal
//			detail.Total += ord.Total
//			detail.PartnerShipFee += ord.PartnerShipFee
//			detail.ShipFee += ord.ShipFee
//			detail.Noo++
//			rs[ord.Status] = detail
//
//			if statsmap[ord.Status].Finish {
//				camp.Base += ord.BaseTotal
//				camp.Total += ord.Total
//				camp.PartnerShipFee += ord.PartnerShipFee
//				camp.ShipFee += ord.ShipFee
//				camp.Noo++
//			}
//		}
//		camp.StatusDetail = rs
//	}
//	return camp
//}
//func SaveCampaign(camp models.Campaign) models.Campaign {
//	col := db.Collection("addons_campaigns")
//	if camp.ID == "" {
//		camp.ID = bson.NewObjectId()
//		camp.Created = time.Now().UTC()
//	}
//
//	camp.Modified = camp.Created
//	opts := options.Update().SetUpsert(true)
//	filter := bson.D{{"_id", camp.ID}}
//	update := bson.D{{"$set", camp}}
//
//	col.UpdateOne(ctx,filter, update,opts)
//	return camp
//}
//func DeleteCampaign(camp models.Campaign) bool {
//	col := db.Collection("addons_campaigns")
//	_,err := col.DeleteOne(ctx,camp)
//	return c3mcommon.CheckError("Delete campaign", err)
//
//}
