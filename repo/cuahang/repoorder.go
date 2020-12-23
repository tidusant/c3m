package cuahang

import (
	"github.com/tidusant/c3m-common/c3mcommon"
	"github.com/tidusant/c3m-common/log"
	"github.com/tidusant/c3m-common/mystring"
	"github.com/tidusant/chadmin-repo/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
	"time"
)

func GetOrderByID(orderid, shopid primitive.ObjectID) models.Order {
	col := db.Collection("addons_orders")
	var rs models.Order
	cond := bson.M{"shopid": shopid, "_id": orderid}
	err := col.FindOne(ctx, cond).Decode(&rs)
	c3mcommon.CheckError("GetOrdersByID", err)
	return rs
}

func GetOrdersByStatus(shopid primitive.ObjectID, status string, page int, pagesize int64, searchterm string) []models.Order {
	col := db.Collection("addons_orders")
	var rs []models.Order

	cond := bson.M{"shopid": shopid.Hex()}
	if status != "all" {
		cond["status"] = status
	}
	if searchterm != "" {
		//searchtermslug := strings.Replace(searchterm, " ", "-", -1)
		searchtermslug := mystring.ParameterizeJoin(searchterm, " ")
		log.Debugf("searchteram slug: $s", searchtermslug)
		//searchtermslug = strings.Replace(searchtermslug, "-", " ", -1)
		//log.Debugf("searchteram slug: $s", searchtermslug)
		cond["$or"] = []bson.M{
			bson.M{"searchindex": bson.M{"$regex": bson.RegEx{searchtermslug, "si"}}},
			// bson.M{"email": bson.M{"$regex": bson.RegEx{searchterm, "si"}}},
			// bson.M{"name": bson.M{"$regex": bson.RegEx{searchterm, "si"}}},
			// bson.M{"name": bson.M{"$regex": bson.RegEx{searchtermslug, "si"}}},
			// bson.M{"address": bson.M{"$regex": bson.RegEx{searchterm, "si"}}},
			// bson.M{"address": bson.M{"$regex": bson.RegEx{searchtermslug, "si"}}},
			// bson.M{"note": bson.M{"$regex": bson.RegEx{searchterm, "si"}}},
			// bson.M{"note": bson.M{"$regex": bson.RegEx{searchtermslug, "si"}}},
		}
	}

	var err error
	var cursor *mongo.Cursor
	if page == 0 {
		cursor, err = col.Find(ctx, cond)
	} else {
		opts := options.Find()
		opts.SetSort(bson.D{{"_id", -1}})
		opts.SetLimit(pagesize)
		opts.SetSkip(int64(page-1) * pagesize)
		cursor, err = col.Find(ctx, cond, opts) //.Sort("-_id").Skip((page - 1) * pagesize).Limit(pagesize).All(&rs)
	}
	log.Debugf("GetOrdersByStatus %+v", cond)
	var rss []bson.M
	if err = cursor.All(ctx, &rss); err != nil {
		c3mcommon.CheckError("GetOrdersByStatus", err)
	}

	return rs
}
func CountOrdersByStatus(shopid primitive.ObjectID, status, searchterm string) int64 {
	col := db.Collection("addons_orders")

	cond := bson.M{"shopid": shopid.Hex()}
	if status != "all" {
		cond["status"] = status
	}
	if searchterm != "" {
		searchtermslug := mystring.ParameterizeJoin(searchterm, " ")
		log.Debugf("searchteram slug: $s", searchtermslug)
		//searchtermslug = strings.Replace(searchtermslug, "-", " ", -1)
		//log.Debugf("searchteram slug: $s", searchtermslug)
		cond["$or"] = []bson.M{
			bson.M{"searchindex": bson.M{"$regex": bson.RegEx{searchtermslug, "si"}}},
		}
	}

	count, err := col.CountDocuments(ctx, cond)
	log.Debugf("count search: %v", count)
	c3mcommon.CheckError("CountOrdersByStatus", err)
	return count
}
func GetOrdersByCampaign(camp models.Campaign) []models.Order {
	col := db.Collection("addons_orders")
	var rs []models.Order
	cond := bson.M{"shopid": camp.ShopId}
	cursor, err := col.Find(ctx, cond)
	if err = cursor.All(ctx, &rs); err != nil {
		c3mcommon.CheckError("GetOrdersByCampaign", err)
	}

	c3mcommon.CheckError("GetOrdersByCamp", err)
	return rs
}

func GetDefaultOrderStatus(shopid primitive.ObjectID) models.OrderStatus {
	col := db.Collection("addons_order_status")
	var rs models.OrderStatus

	cond := bson.M{"shopid": shopid.Hex(), "default": true}

	err := col.FindOne(ctx, cond).Decode(&rs)

	c3mcommon.CheckError("GetDefaultOrderStatus", err)
	return rs
}

func UpdateOrderStatus(shopid, orderid primitive.ObjectID, status string) {
	col := db.Collection("addons_orders")
	//var arrIdObj []bson.ObjectId
	// for _, v := range orderid {
	// 	arrIdObj = append(arrIdObj, bson.ObjectIdHex(v))
	// }
	// cond := bson.M{"_id": bson.M{"$in": arrIdObj}, "shopid": shopid}
	cond := bson.M{"_id": orderid, "shopid": shopid.Hex()}
	change := bson.M{"status": status}
	stats := GetAllOrderStatus(shopid)
	mapstat := make(map[string]models.OrderStatus)
	for _, v := range stats {
		mapstat[v.ID.Hex()] = v
	}

	if mapstat[status].Finish {
		change["whookupdate"] = time.Now().Unix()
	}
	log.Debugf("udpate order cond:%v", cond)
	_, err := col.UpdateOne(ctx, cond, bson.M{"$set": change})
	c3mcommon.CheckError("Update order status", err)
	return
}

func SaveOrder(order models.Order) models.Order {
	//col := db.Collection("addons_orders")
	//if order.ID == primitive.NilObjectID {
	//	order.ID = primitive.NewObjectID()
	//	order.Created = time.Now().Unix()
	//
	//}
	//
	//order.Modified = time.Now().Unix()
	//opts := options.Update().SetUpsert(true)
	//col.UpsertId(order.ID, order)
	//_, err := coluserlogin.UpdateOne(ctx, filter, update, opts)
	//c3mcommon.CheckError("SaveOrder", err)
	return order
}

func GetCountOrderByStatus(stat models.OrderStatus) int64 {
	col := db.Collection("addons_orders")

	cond := bson.M{"shopid": stat.ShopId, "status": stat.ID.Hex()}
	n, err := col.CountDocuments(ctx, cond)
	c3mcommon.CheckError("GetCountOrderByStatus", err)
	return n
}

//====================== whook

//===============status function
func UpdateOrderStatusByShipmentCode(shipmentCode string, statusid, shopid primitive.ObjectID) {
	col := db.Collection("addons_orders")
	cond := bson.M{"shopid": shopid, "shipmentcode": shipmentCode}
	change := bson.M{"$set": bson.M{"status": statusid, "whookupdate": time.Now().Unix()}}

	_, err := col.UpdateOne(ctx, cond, change)
	c3mcommon.CheckError("UpdateOrderStatusByShipmentCode", err)
}
func GetOrderByShipmentCode(shipmentCode string, shopid primitive.ObjectID) models.Order {
	col := db.Collection("addons_orders")
	var rs models.Order
	cond := bson.M{"shopid": shopid.Hex(), "shipmentcode": shipmentCode}
	err := col.FindOne(ctx, cond).Decode(&rs)

	c3mcommon.CheckError("GetOrderByShipmentCode", err)
	return rs
}
func GetStatusByPartnerStatus(shopid primitive.ObjectID, partnercode, partnerstatus string) models.OrderStatus {
	col := db.Collection("addons_order_status")
	var rs models.OrderStatus

	cond := bson.M{"shopid": shopid.Hex(), "partnerstatus." + partnercode: partnerstatus}
	err := col.FindOne(ctx, cond).Decode(&rs)
	c3mcommon.CheckError("GetStatusByPartnerStatus", err)
	return rs
}
func GetAllOrderStatus(shopid primitive.ObjectID) []models.OrderStatus {
	col := db.Collection("addons_order_status")
	var rs []models.OrderStatus
	cond := bson.M{"shopid": shopid.Hex()}
	cursor, err := col.Find(ctx, cond)
	if err = cursor.All(ctx, &rs); err != nil {
		c3mcommon.CheckError("GetAllOrderStatus", err)
	}

	return rs
}
func SaveOrderStatus(status models.OrderStatus) models.OrderStatus {
	//col := db.Collection("addons_order_status")
	//if status.ID.Hex() == "" {
	//	status.ID = primitive.NewObjectID()
	//	status.Created = time.Now().UTC()
	//}
	//
	//status.Modified = status.Created
	//col.UpsertId(status.ID, status)
	return status
}

func GetStatusByID(statusid, shopid primitive.ObjectID) models.OrderStatus {
	col := db.Collection("addons_order_status")
	var rs models.OrderStatus
	cond := bson.M{"shopid": shopid.Hex(), "_id": statusid}
	err := col.FindOne(ctx, cond).Decode(&rs)
	c3mcommon.CheckError("GetstatusByID", err)
	return rs
}

func DeleteOrderStatus(stat models.OrderStatus) bool {
	col := db.Collection("addons_order_status")

	cond := bson.M{"shopid": stat.ShopId, "_id": stat.ID}
	_, err := col.DeleteOne(ctx, cond)
	return c3mcommon.CheckError("GetstatusByID", err)

}

func UnSetStatusDefault(shopid primitive.ObjectID) {
	col := db.Collection("addons_order_status")

	cond := bson.M{"shopid": shopid, "default": true}
	change := bson.M{"$set": bson.M{"default": false}}
	_, err := col.UpdateOne(ctx, cond, change)
	c3mcommon.CheckError("UnSetStatusDefault", err)

}
