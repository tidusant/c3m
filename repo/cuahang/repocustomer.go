package cuahang

import (
	"github.com/tidusant/c3m-common/c3mcommon"

	"github.com/tidusant/chadmin-repo/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"

	"gopkg.in/mgo.v2/bson"
)

func CountOrderByCus(phone, shopid string) int {
	col := db.Collection("addons_orders")
	cond := bson.M{"shopid": shopid, "phone": phone}
	rs, err := col.CountDocuments(ctx, cond)
	c3mcommon.CheckError("count order cus by phone", err)
	return int(rs)
}
func GetAllCustomers(shopid string) []models.Customer {
	col := db.Collection("addons_customers")
	var rs []models.Customer
	cond := bson.M{"shopid": shopid}
	cursor, err := col.Find(ctx, cond)
	if err = cursor.All(ctx, &rs); err != nil {
		log.Fatal(err)
	}
	c3mcommon.CheckError("GetAllCustomers", err)
	return rs
}
func GetCusByPhone(phone, shopid string) models.Customer {
	col := db.Collection("addons_customers")
	var rs models.Customer
	cond := bson.M{"shopid": shopid, "phone": phone}
	err := col.FindOne(ctx, cond).Decode(&rs)
	c3mcommon.CheckError("get cus by phone", err)
	return rs
}
func GetCusByEmail(email, shopid string) models.Customer {
	col := db.Collection("addons_customers")
	var rs models.Customer
	cond := bson.M{"shopid": shopid, "email": email}
	err := col.FindOne(ctx, cond).Decode(&rs)
	c3mcommon.CheckError("get cus by email", err)
	return rs
}
func SaveCus(cus models.Customer) bool {
	col := db.Collection("addons_customers")

	if cus.ID == primitive.NilObjectID {
		cus.ID = primitive.NewObjectID()

	}
	opts := options.Update().SetUpsert(true)
	filter := bson.D{{"_id", cus.ID}}
	update := bson.D{{"$set", cus}}
	_, err := col.UpdateOne(ctx, filter, update, opts)
	if c3mcommon.CheckError("save cus ", err) {
		return true
	}
	return false
}
