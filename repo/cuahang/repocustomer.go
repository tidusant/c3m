package cuahang

import (
	"github.com/tidusant/c3m/common/c3mcommon"
	"go.mongodb.org/mongo-driver/bson"
	"time"

	"github.com/tidusant/c3m/repo/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (r *Repo) CountOrderByCus(phone, shopid string) int {
	start := time.Now()
	col := db.Collection("addons_orders")
	cond := bson.M{"shopid": shopid, "phone": phone}
	rs, err := col.CountDocuments(ctx, cond)
	r.QueryCount++
	c3mcommon.CheckError("count order cus by phone", err)
	r.QueryTime += time.Since(start)
	return int(rs)
}
func (r *Repo) GetAllCustomers(shopid string) []models.Customer {
	start := time.Now()
	col := db.Collection("addons_customers")
	var rs []models.Customer
	cond := bson.M{"shopid": shopid}
	cursor, err := col.Find(ctx, cond)
	r.QueryCount++
	if err = cursor.All(ctx, &rs); err != nil {
		c3mcommon.CheckError("GetAllCustomers", err)
	}

	r.QueryTime += time.Since(start)
	return rs
}
func (r *Repo) GetCustomerByPhones(phones []string, shopid primitive.ObjectID) []models.Customer {
	start := time.Now()
	col := db.Collection("addons_customers")
	var rs []models.Customer
	cond := bson.D{
		{"shopid", shopid},
		{"phone", bson.D{{"$in", phones}}},
	}
	cursor, err := col.Find(ctx, cond)
	r.QueryCount++

	if err = cursor.All(ctx, &rs); err != nil {
		c3mcommon.CheckError("GetCustomerByPhones", err)
	}

	r.QueryTime += time.Since(start)
	return rs
}
func (r *Repo) GetCusByPhone(phone, shopid string) models.Customer {
	start := time.Now()
	col := db.Collection("addons_customers")
	var rs models.Customer
	cond := bson.M{"shopid": shopid, "phone": phone}
	err := col.FindOne(ctx, cond).Decode(&rs)
	r.QueryCount++
	c3mcommon.CheckError("get cus by phone", err)
	r.QueryTime += time.Since(start)
	return rs
}
func (r *Repo) GetCusByEmail(email, shopid string) models.Customer {
	start := time.Now()
	col := db.Collection("addons_customers")
	var rs models.Customer
	cond := bson.M{"shopid": shopid, "email": email}
	err := col.FindOne(ctx, cond).Decode(&rs)
	r.QueryCount++
	c3mcommon.CheckError("get cus by email", err)
	r.QueryTime += time.Since(start)
	return rs
}
func (r *Repo) SaveCus(cus models.Customer) bool {
	start := time.Now()
	col := db.Collection("addons_customers")

	if cus.ID == primitive.NilObjectID {
		cus.ID = primitive.NewObjectID()

	}
	opts := options.Update().SetUpsert(true)
	filter := bson.D{{"_id", cus.ID}}
	update := bson.D{{"$set", cus}}
	_, err := col.UpdateOne(ctx, filter, update, opts)
	r.QueryCount++
	if c3mcommon.CheckError("save cus ", err) {
		return true
	}
	r.QueryTime += time.Since(start)
	return false
}
