package cuahang

//
//import (
//	c3mcommon "colis/common/common"
//	"time"
//
//	"colis/models"
//
//	"gopkg.in/mgo.v2/bson"
//)
//
//func GetAllShipper(shopid string) []models.Shipper {
//	col := db.C("addons_shippers")
//	var rs []models.Shipper
//	cond := bson.M{"shopid": shopid}
//	err := col.Find(cond).All(&rs)
//	c3mcommon.CheckError("GetAllShipper", err)
//	return rs
//}
//func GetShipperByID(itemid, shopid string) models.Shipper {
//	col := db.C("addons_shippers")
//	var rs models.Shipper
//	cond := bson.M{"shopid": shopid, "_id": bson.ObjectIdHex(itemid)}
//	err := col.Find(cond).One(&rs)
//	c3mcommon.CheckError("GetShipperByID", err)
//	return rs
//}
//
//func GetDefaultShipper(shopid string) models.Shipper {
//	col := db.C("addons_shippers")
//	var rs models.Shipper
//	cond := bson.M{"shopid": shopid, "default": true}
//	err := col.Find(cond).One(&rs)
//	c3mcommon.CheckError("GetDefaultShipper", err)
//	return rs
//}
//
//func SaveShipper(shipper models.Shipper) models.Shipper {
//	col := db.C("addons_shippers")
//	if shipper.ID == "" {
//		shipper.ID = bson.NewObjectId()
//		shipper.Created = time.Now().UTC()
//	}
//
//	shipper.Modified = shipper.Created
//	_, err := col.UpsertId(shipper.ID, shipper)
//	c3mcommon.CheckError("SaveShipper", err)
//	return shipper
//}
//
//func GetCountOrderByShipper(shipper models.Shipper) int {
//	col := db.C("addons_orders")
//
//	cond := bson.M{"shopid": shipper.ShopId, "shipperid": shipper.ID.Hex()}
//	n, err := col.Find(cond).Count()
//	c3mcommon.CheckError("GetCountOrderByShipper", err)
//	return n
//}
//
//func DeleteShipper(shipper models.Shipper) bool {
//	col := db.C("addons_shippers")
//
//	cond := bson.M{"shopid": shipper.ShopId, "_id": shipper.ID}
//	err := col.Remove(cond)
//	return c3mcommon.CheckError("DeleteShipper", err)
//
//}
//
//func UnSetShipperDefault(shopid string) {
//	col := db.C("addons_shippers")
//
//	cond := bson.M{"shopid": shopid, "default": true}
//	change := bson.M{"$set": bson.M{"default": false}}
//	err := col.Update(cond, change)
//	c3mcommon.CheckError("UnSetShipperDefault", err)
//
//}
