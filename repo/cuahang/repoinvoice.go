package cuahang

import (
	"github.com/tidusant/c3m-common/c3mcommon"

	"context"
	"github.com/tidusant/chadmin-repo/models"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"

	//	"c3m/log"

	//"strings"

	"gopkg.in/mgo.v2/bson"
)

func SaveInvoice(invc models.Invoice) models.Invoice {

	col := db.Collection("addons_invoice")

	if invc.ID == "" {
		invc.ID = bson.NewObjectId()
	}
	opts := options.Update().SetUpsert(true)
	filter := bson.D{{"_id", invc.ID}}
	update := bson.D{{"$set", invc}}
	col.UpdateOne(context.TODO(), filter, update, opts)
	return invc
}

func GetInvoices(shopid string, imp bool) []models.Invoice {

	col := db.Collection("addons_invoice")

	var rs []models.Invoice
	opts := options.Find().SetSort(bson.M{"Created": "desc"})
	cursor, err := col.Find(context.TODO(), bson.M{"shopid": shopid, "import": imp}, opts)
	if err = cursor.All(context.TODO(), &rs); err != nil {
		log.Fatal(err)
	}
	c3mcommon.CheckError("GetInvoices", err)
	return rs
}
func GetInvcById(shopid, invcid string) models.Invoice {

	col := db.Collection("addons_invoice")
	var rs models.Invoice
	err := col.FindOne(context.TODO(), bson.M{"shopid": shopid, "_id": bson.ObjectIdHex(invcid)}).Decode(&rs)
	c3mcommon.CheckError("GetInvcById", err)

	return rs
}
func RemoveInvcById(shopid, invcid string) bool {

	col := db.Collection("addons_invoice")

	_, err := col.DeleteOne(context.TODO(), bson.M{"_id": invcid})
	c3mcommon.CheckError("RemoveInvcById", err)
	return true
}
