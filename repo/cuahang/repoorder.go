package cuahang

import (
	"github.com/tidusant/c3m/common/c3mcommon"
	"github.com/tidusant/c3m/common/log"
	"github.com/tidusant/c3m/common/mystring"
	"github.com/tidusant/c3m/repo/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func (r *Repo) GetOrderByID(orderid, shopid primitive.ObjectID) models.Order {
	start := time.Now()
	col := db.Collection("addons_orders")
	var rs models.Order
	cond := bson.M{"shopid": shopid, "_id": orderid}
	err := col.FindOne(ctx, cond).Decode(&rs)
	r.QueryCount++
	c3mcommon.CheckError("GetOrdersByID", err)
	r.QueryTime += time.Since(start)
	return rs
}

func (r *Repo) GetOrdersByStatus(shopid primitive.ObjectID, statusid primitive.ObjectID, page int, pagesize int, searchterm string) []models.Order {
	start := time.Now()
	col := db.Collection("addons_orders")
	var rs []models.Order

	var cond bson.D
	cond = append(cond, bson.E{"shopid", shopid})
	if statusid != primitive.NilObjectID {
		cond = append(cond, bson.E{"status", statusid})
	}
	if searchterm != "" {
		//searchtermslug := strings.Replace(searchterm, " ", "-", -1)
		searchtermslug := mystring.ParameterizeJoin(searchterm, " ")
		cond = append(cond, bson.E{"$or", bson.D{
			{"searchindex", bson.M{"$regex": searchtermslug}},
			// bson.M{"email": bson.M{"$regex": bson.RegEx{searchterm, "si"}}},
			// bson.M{"name": bson.M{"$regex": bson.RegEx{searchterm, "si"}}},
			// bson.M{"name": bson.M{"$regex": bson.RegEx{searchtermslug, "si"}}},
			// bson.M{"address": bson.M{"$regex": bson.RegEx{searchterm, "si"}}},
			// bson.M{"address": bson.M{"$regex": bson.RegEx{searchtermslug, "si"}}},
			// bson.M{"note": bson.M{"$regex": bson.RegEx{searchterm, "si"}}},
			// bson.M{"note": bson.M{"$regex": bson.RegEx{searchtermslug, "si"}}},
		}})
	}

	//groupStage := bson.D{{"$group", bson.D{{"_id", "$status"}, {"order_count", bson.D{{"$sum", 1}}}}}}
	//lookupStage := bson.D{{"$lookup", bson.D{{"from", "addons_order_status"}, {"localField", "_id"}, {"foreignField", "_id"}, {"as", "status"}}}}
	//unwindStage := bson.D{{"$unwind", bson.D{{"path", "$status"}, {"preserveNullAndEmptyArrays", false}}}}
	//addfield:= bson.D{{"$addFields", bson.D{
	//	{"title", "$status.title"},
	//	{"color", "$status.color"},
	//	{"default", "$status.default"},
	//	{"finish", "$status.finish"},
	//}}}

	if page < 1 {
		page = 1
	}

	//sort:=bson.D{{"$sort", bson.D{
	//	{"_id", -1},
	//}}}
	//skip:=bson.D{{"$skip", (page-1) * pagesize}}
	//limit:=bson.D{{"$limit",  pagesize}}
	//
	//pipe:=mongo.Pipeline{
	//	bson.D{{"$match", cond}},
	//	sort,
	//	skip,
	//	limit,
	//}

	//cursor, err := col.Aggregate(ctx, pipe)

	opts := options.Find()
	opts.SetSkip(int64((page - 1) * pagesize))
	opts.SetLimit(int64(pagesize))
	opts.SetSort(bson.D{{"_id", -1}})
	cursor, err := col.Find(ctx, cond, opts)
	r.QueryCount++

	if c3mcommon.CheckError("GetOrdersByStatus", err) {
		//var rss []bson.M
		if err = cursor.All(ctx, &rs); err != nil {
			c3mcommon.CheckError("GetOrdersByStatus cursor.All", err)
		}
	}

	r.QueryTime += time.Since(start)
	return rs
}
func (r *Repo) CountOrdersByStatus(shopid, status primitive.ObjectID, searchterm string) int64 {
	start := time.Now()
	col := db.Collection("addons_orders")

	cond := bson.M{"shopid": shopid}
	if status != primitive.NilObjectID {
		cond["status"] = status
	}
	if searchterm != "" {
		searchtermslug := mystring.ParameterizeJoin(searchterm, " ")
		log.Debugf("searchteram slug: $s", searchtermslug)
		//searchtermslug = strings.Replace(searchtermslug, "-", " ", -1)
		//log.Debugf("searchteram slug: $s", searchtermslug)
		cond["$or"] = []bson.M{
			bson.M{"searchindex": bson.M{"$regex": searchtermslug}},
		}
	}

	count, err := col.CountDocuments(ctx, cond)
	r.QueryCount++
	log.Debugf("count search: %v", count)
	c3mcommon.CheckError("CountOrdersByStatus", err)
	r.QueryTime += time.Since(start)
	return count
}
func (r *Repo) GetOrdersByCampaign(camp models.Campaign) []models.Order {
	start := time.Now()
	col := db.Collection("addons_orders")
	var rs []models.Order
	cond := bson.M{"shopid": camp.ShopId}
	cursor, err := col.Find(ctx, cond)
	r.QueryCount++
	if err = cursor.All(ctx, &rs); err != nil {
		c3mcommon.CheckError("GetOrdersByCampaign", err)
	}

	c3mcommon.CheckError("GetOrdersByCamp", err)
	r.QueryTime += time.Since(start)
	return rs
}

func (r *Repo) GetDefaultOrderStatus(shopid primitive.ObjectID) models.OrderStatus {
	start := time.Now()
	col := db.Collection("addons_order_status")
	var rs models.OrderStatus

	cond := bson.M{"shopid": shopid.Hex(), "default": true}

	err := col.FindOne(ctx, cond).Decode(&rs)
	r.QueryCount++
	c3mcommon.CheckError("GetDefaultOrderStatus", err)
	r.QueryTime += time.Since(start)
	return rs
}

func (r *Repo) UpdateOrderStatus(shopid, orderid primitive.ObjectID, status string) {
	start := time.Now()
	col := db.Collection("addons_orders")
	//var arrIdObj []bson.ObjectId
	// for _, v := range orderid {
	// 	arrIdObj = append(arrIdObj, bson.ObjectIdHex(v))
	// }
	// cond := bson.M{"_id": bson.M{"$in": arrIdObj}, "shopid": shopid}
	cond := bson.M{"_id": orderid, "shopid": shopid.Hex()}
	change := bson.M{"status": status}
	stats := r.GetAllOrderStatus(shopid)
	mapstat := make(map[string]models.OrderStatus)
	for _, v := range stats {
		mapstat[v.ID.Hex()] = v
	}

	if mapstat[status].Finish {
		change["whookupdate"] = time.Now().Unix()
	}
	log.Debugf("udpate order cond:%v", cond)
	_, err := col.UpdateOne(ctx, cond, bson.M{"$set": change})
	r.QueryCount++
	c3mcommon.CheckError("Update order status", err)
	r.QueryTime += time.Since(start)
	return
}

func (r *Repo) SaveOrder(order models.Order) models.Order {
	start := time.Now()
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
	r.QueryTime += time.Since(start)
	return order
}

func (r *Repo) GetCountOrderByStatus(stat models.OrderStatus) int64 {
	start := time.Now()
	col := db.Collection("addons_orders")

	cond := bson.M{"shopid": stat.ShopId, "status": stat.ID.Hex()}
	n, err := col.CountDocuments(ctx, cond)
	r.QueryCount++
	c3mcommon.CheckError("GetCountOrderByStatus", err)
	r.QueryTime += time.Since(start)
	return n
}

//====================== whook

//===============status function
func (r *Repo) UpdateOrderStatusByShipmentCode(shipmentCode string, statusid, shopid primitive.ObjectID) {
	start := time.Now()
	col := db.Collection("addons_orders")
	cond := bson.M{"shopid": shopid, "shipmentcode": shipmentCode}
	change := bson.M{"$set": bson.M{"status": statusid, "whookupdate": time.Now().Unix()}}

	_, err := col.UpdateOne(ctx, cond, change)
	r.QueryCount++
	c3mcommon.CheckError("UpdateOrderStatusByShipmentCode", err)
	r.QueryTime += time.Since(start)
}
func (r *Repo) GetOrderByShipmentCode(shipmentCode string, shopid primitive.ObjectID) models.Order {
	start := time.Now()
	col := db.Collection("addons_orders")
	var rs models.Order
	cond := bson.M{"shopid": shopid.Hex(), "shipmentcode": shipmentCode}
	err := col.FindOne(ctx, cond).Decode(&rs)
	r.QueryCount++
	c3mcommon.CheckError("GetOrderByShipmentCode", err)
	r.QueryTime += time.Since(start)
	return rs
}
func (r *Repo) GetStatusByPartnerStatus(shopid primitive.ObjectID, partnercode, partnerstatus string) models.OrderStatus {
	start := time.Now()
	col := db.Collection("addons_order_status")
	var rs models.OrderStatus

	cond := bson.M{"shopid": shopid.Hex(), "partnerstatus." + partnercode: partnerstatus}
	err := col.FindOne(ctx, cond).Decode(&rs)
	r.QueryCount++
	c3mcommon.CheckError("GetStatusByPartnerStatus", err)
	r.QueryTime += time.Since(start)
	return rs
}
func (r *Repo) CoundOrderByPhones(phones []string, shopid primitive.ObjectID) map[string]int {
	start := time.Now()
	col := db.Collection("addons_orders")

	pipe := mongo.Pipeline{
		bson.D{{"$match", bson.D{
			{"shopid", shopid},
			{"phone", bson.D{{"$in", phones}}},
		}}},
		//bson.D{{"$lookup", bson.D{{"from", "addons_orders"}, {"localField", "phone"}, {"foreignField", "phone"}, {"as", "order"}}}},
		//bson.D{{"$unwind", bson.D{{"path", "$order"}, {"preserveNullAndEmptyArrays", false}}}},
		bson.D{{"$group", bson.D{
			{"_id", "$phone"},
			{"count", bson.D{{"$sum", 1}}},
		}}},
		//bson.D{{"$addFields", bson.D{
		//	{"name", "$_id.name"},
		//	{"_id", "$_id._id"},
		//	{"phone", "$_id.phone"},
		//}}},
	}
	cursor, err := col.Aggregate(ctx, pipe)
	r.QueryCount++
	type RT struct {
		ID    string `bson:"_id""`
		Count int
	}
	var rs []RT
	if err = cursor.All(ctx, &rs); err != nil {
		c3mcommon.CheckError("CoundOrderByPhones", err)
	}

	//convert data to map
	cuscount := make(map[string]int)
	for _, v := range rs {
		cuscount[v.ID] = v.Count
	}

	r.QueryTime += time.Since(start)
	return cuscount
}
func (r *Repo) GetAllOrderStatus(shopid primitive.ObjectID) []models.OrderStatus {
	start := time.Now()
	col := db.Collection("addons_orders")
	var rs []models.OrderStatus
	//cond := bson.M{"shopid": shopid.Hex()}
	matchStage := bson.D{{"$match", bson.D{{"shopid", shopid}}}}
	groupStage := bson.D{{"$group", bson.D{{"_id", "$status"}, {"orderCount", bson.D{{"$sum", 1}}}}}}
	lookupStage := bson.D{{"$lookup", bson.D{{"from", "addons_order_status"}, {"localField", "_id"}, {"foreignField", "_id"}, {"as", "status"}}}}
	unwindStage := bson.D{{"$unwind", bson.D{{"path", "$status"}, {"preserveNullAndEmptyArrays", false}}}}
	addfield := bson.D{{"$addFields", bson.D{
		{"title", "$status.title"},
		{"color", "$status.color"},
		{"default", "$status.default"},
		{"finish", "$status.finish"},
	}}}
	cursor, err := col.Aggregate(ctx, mongo.Pipeline{
		matchStage, groupStage, lookupStage, unwindStage, addfield,
		bson.D{{"$sort", bson.D{
			{"_id", -1},
		}}},
	})
	r.QueryCount++

	if err = cursor.All(ctx, &rs); err != nil {
		c3mcommon.CheckError("GetAllOrderStatus", err)
	}
	//log.Debugf("aggregate rs:%+v",rs)
	r.QueryTime += time.Since(start)
	return rs
}
func (r *Repo) SaveOrderStatus(status models.OrderStatus) models.OrderStatus {
	start := time.Now()
	//col := db.Collection("addons_order_status")
	//if status.ID.Hex() == "" {
	//	status.ID = primitive.NewObjectID()
	//	status.Created = time.Now().UTC()
	//}
	//
	//status.Modified = status.Created
	//col.UpsertId(status.ID, status)
	r.QueryTime += time.Since(start)
	return status
}

func (r *Repo) GetStatusByID(statusid, shopid primitive.ObjectID) models.OrderStatus {
	start := time.Now()
	col := db.Collection("addons_order_status")
	var rs models.OrderStatus
	cond := bson.M{"shopid": shopid.Hex(), "_id": statusid}
	err := col.FindOne(ctx, cond).Decode(&rs)
	r.QueryCount++
	c3mcommon.CheckError("GetstatusByID", err)
	r.QueryTime += time.Since(start)
	return rs
}

func (r *Repo) DeleteOrderStatus(stat models.OrderStatus) bool {
	start := time.Now()
	col := db.Collection("addons_order_status")

	cond := bson.M{"shopid": stat.ShopId, "_id": stat.ID}
	_, err := col.DeleteOne(ctx, cond)
	r.QueryCount++
	r.QueryTime += time.Since(start)
	return c3mcommon.CheckError("GetstatusByID", err)

}

func (r *Repo) UnSetStatusDefault(shopid primitive.ObjectID) {
	start := time.Now()
	col := db.Collection("addons_order_status")

	cond := bson.M{"shopid": shopid, "default": true}
	change := bson.M{"$set": bson.M{"default": false}}
	_, err := col.UpdateOne(ctx, cond, change)
	r.QueryCount++
	c3mcommon.CheckError("UnSetStatusDefault", err)
	r.QueryTime += time.Since(start)

}
