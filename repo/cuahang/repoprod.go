package cuahang

import (
	"github.com/tidusant/c3m-common/c3mcommon"

	"github.com/tidusant/chadmin-repo/models"
	//	"c3m/log"

	//"strings"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
)

func SaveProd(prod models.Product) string {

	col := db.Collection("addons_products")

	// if prod.Code {

	// 	err := col.Insert(prod)
	// 	c3mcommon.CheckError("product Insert", err)
	// } else {
	if len(prod.Langs) > 0 {
		if prod.ID.Hex() == "" {
			prod.ID = primitive.NewObjectID()
		}
		//_, err := col.UpsertId(prod.ID, &prod)
		//c3mcommon.CheckError("SaveProd", err)
	} else {
		col.DeleteOne(ctx, bson.M{"_id": prod.ID})
	}
	//}
	langinfo, _ := json.Marshal(prod.Langs)
	propinfo, _ := json.Marshal(prod.Properties)
	return "{\"Code\":\"" + prod.Code + "\",\"Langs\":" + string(langinfo) + ",\"Properties\":" + string(propinfo) + "}"
}
func SaveProperties(shopid primitive.ObjectID, code string, props []models.ProductProperty) bool {
	col := db.Collection("addons_products")

	cond := bson.M{"shopid": shopid.Hex(), "code": code}
	change := bson.M{"properties": props}
	_, err := col.UpdateOne(ctx, cond, bson.M{"$set": change})

	return c3mcommon.CheckError("SaveProperties", err)

}
func GetProds(userid, shopid primitive.ObjectID, isMain bool) []models.Product {
	col := db.Collection("addons_products")
	var rs []models.Product

	cursor, err := col.Find(ctx, bson.M{"shopid": shopid, "main": isMain})
	if err = cursor.All(ctx, &rs); err != nil {
		c3mcommon.CheckError("GetProds", err)
	}

	return rs
}
func GetAllProds(userid, shopid primitive.ObjectID) []models.Product {
	col := db.Collection("addons_products")
	var rs []models.Product

	cursor, err := col.Find(ctx, bson.M{"shopid": shopid})
	if err = cursor.All(ctx, &rs); err != nil {
		c3mcommon.CheckError("GetAllProds", err)
	}

	return rs
}
func GetDemoProds() []models.Product {
	col := db.Collection("addons_products")
	var rs []models.Product
	shop := GetDemoShop()
	cursor, err := col.Find(ctx, bson.M{"shopid": shop.ID.Hex()})
	if err = cursor.All(ctx, &rs); err != nil {
		c3mcommon.CheckError("get demo prod", err)
	}

	return rs
}
func GetProdBySlug(shopid primitive.ObjectID, slug string) models.Product {
	col := db.Collection("addons_products")
	var rs models.Product
	cond := bson.M{"shopid": shopid, "slug": slug}

	err := col.FindOne(ctx, cond).Decode(&rs)
	c3mcommon.CheckError("GetProdBySlug", err)
	return rs
}
func GetProdByCode(shopid primitive.ObjectID, code string) models.Product {
	col := db.Collection("addons_products")
	var rs models.Product
	cond := bson.M{"shopid": shopid, "code": code}

	err := col.FindOne(ctx, cond).Decode(&rs)
	c3mcommon.CheckError("GetProdByCode", err)
	return rs
}

func GetProdsByCatId(shopid primitive.ObjectID, catcode string) []models.Product {
	col := db.Collection("addons_products")
	var rs []models.Product
	cond := bson.M{"shopid": shopid, "catid": catcode}

	cursor, err := col.Find(ctx, cond)
	if err = cursor.All(ctx, &rs); err != nil {
		c3mcommon.CheckError("GetProdsByCatId", err)
	}

	return rs

}

func ExportItem(exportitems []models.ExportItem) bool {
	col := db.Collection("addons_products")
	var rs models.Product

	//subcond := bson.M{"$elemMatch": bson.M{"code": itemcode}}
	for _, item := range exportitems {
		cond := bson.M{"shopid": item.ShopId, "code": item.Code}
		err := col.FindOne(ctx, cond).Decode(&rs)
		for k, v := range rs.Properties {
			if v.Code == item.ItemCode {
				rs.Properties[k].Stock -= item.Num
				SaveProd(rs)
				break
			}
		}
		c3mcommon.CheckError("ExportItem", err)

	}

	return true

}

//=========================cat function==================
func SaveCat(cat models.ProdCat) string {
	col := db.Collection("addons_prodcats")
	if len(cat.Langs) > 0 {
		if cat.ID.Hex() == "" {
			cat.ID = primitive.NewObjectID()
			//save slug
		} else {
			//update slug
		}

		//col.UpsertId(cat.ID, cat)
	} else {
		col.DeleteOne(ctx, bson.M{"_id": cat.ID})
	}
	langinfo, _ := json.Marshal(cat.Langs)
	return "{\"Code\":\"" + cat.Code + "\",\"Langs\":" + string(langinfo) + "}"
}
func GetAllCats(userid, shopid primitive.ObjectID) []models.ProdCat {
	col := db.Collection("addons_prodcats")
	var rs []models.ProdCat
	cond := bson.M{"shopid": shopid}

	cursor, err := col.Find(ctx, cond)
	if err = cursor.All(ctx, &rs); err != nil {
		c3mcommon.CheckError("GetAllCats", err)
	}
	return rs
}

func GetCats(userid, shopid primitive.ObjectID, ismain bool) []models.ProdCat {
	col := db.Collection("addons_prodcats")
	var rs []models.ProdCat
	cond := bson.M{"shopid": shopid, "main": ismain}

	cursor, err := col.Find(ctx, cond)
	if err = cursor.All(ctx, &rs); err != nil {
		c3mcommon.CheckError("GetCats", err)
	}
	return rs
}

func GetDemoProdCats() []models.ProdCat {
	col := db.Collection("addons_prodcats")
	shop := GetDemoShop()
	var rs []models.ProdCat
	cursor, err := col.Find(ctx, bson.M{"shopid": shop.ID.Hex()})
	if err = cursor.All(ctx, &rs); err != nil {
		c3mcommon.CheckError("GetCats", err)
	}
	return rs
}
func GetCatByCode(shopid primitive.ObjectID, code string) models.ProdCat {
	col := db.Collection("addons_prodcats")
	var rs models.ProdCat
	cond := bson.M{"shopid": shopid, "code": code}

	err := col.FindOne(ctx, cond).Decode(&rs)
	c3mcommon.CheckError("GetCatByCode", err)
	return rs
}
