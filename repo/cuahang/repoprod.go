package cuahang

import (
	"github.com/tidusant/c3m/common/c3mcommon"
	"time"

	"github.com/tidusant/c3m/repo/models"
	//	"c3m/log"

	//"strings"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
)

func (r *Repo) SaveProd(prod models.Product) string {
	start := time.Now()
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
		r.QueryCount++
	}
	//}
	langinfo, _ := json.Marshal(prod.Langs)
	propinfo, _ := json.Marshal(prod.Properties)
	r.QueryTime += time.Since(start)
	return "{\"Code\":\"" + prod.Code + "\",\"Langs\":" + string(langinfo) + ",\"Properties\":" + string(propinfo) + "}"
}
func (r *Repo) SaveProperties(shopid primitive.ObjectID, code string, props []models.ProductProperty) bool {
	start := time.Now()
	col := db.Collection("addons_products")

	cond := bson.M{"shopid": shopid.Hex(), "code": code}
	change := bson.M{"properties": props}
	_, err := col.UpdateOne(ctx, cond, bson.M{"$set": change})
	r.QueryCount++
	r.QueryTime += time.Since(start)
	return c3mcommon.CheckError("SaveProperties", err)

}
func (r *Repo) GetProds(userid, shopid primitive.ObjectID, isMain bool) []models.Product {
	start := time.Now()
	col := db.Collection("addons_products")
	var rs []models.Product

	cursor, err := col.Find(ctx, bson.M{"shopid": shopid, "main": isMain})
	r.QueryCount++
	if err = cursor.All(ctx, &rs); err != nil {
		c3mcommon.CheckError("GetProds", err)
	}
	r.QueryTime += time.Since(start)
	return rs
}
func (r *Repo) GetAllProds(userid, shopid primitive.ObjectID) []models.Product {
	start := time.Now()
	col := db.Collection("addons_products")
	var rs []models.Product

	cursor, err := col.Find(ctx, bson.M{"shopid": shopid})
	r.QueryCount++
	if err = cursor.All(ctx, &rs); err != nil {
		c3mcommon.CheckError("GetAllProds", err)
	}
	r.QueryTime += time.Since(start)
	return rs
}
func (r *Repo) GetDemoProds() []models.Product {
	start := time.Now()
	col := db.Collection("addons_products")
	var rs []models.Product
	shop := r.GetDemoShop()
	cursor, err := col.Find(ctx, bson.M{"shopid": shop.ID.Hex()})
	r.QueryCount++
	if err = cursor.All(ctx, &rs); err != nil {
		c3mcommon.CheckError("get demo prod", err)
	}
	r.QueryTime += time.Since(start)
	return rs
}
func (r *Repo) GetProdBySlug(shopid primitive.ObjectID, slug string) models.Product {
	start := time.Now()
	col := db.Collection("addons_products")
	var rs models.Product
	cond := bson.M{"shopid": shopid, "slug": slug}

	err := col.FindOne(ctx, cond).Decode(&rs)
	r.QueryCount++
	c3mcommon.CheckError("GetProdBySlug", err)
	r.QueryTime += time.Since(start)
	return rs
}
func (r *Repo) GetProdByCode(shopid primitive.ObjectID, code string) models.Product {
	start := time.Now()
	col := db.Collection("addons_products")
	var rs models.Product
	cond := bson.M{"shopid": shopid, "code": code}

	err := col.FindOne(ctx, cond).Decode(&rs)
	r.QueryCount++
	c3mcommon.CheckError("GetProdByCode", err)
	r.QueryTime += time.Since(start)
	return rs
}

func (r *Repo) GetProdsByCatId(shopid primitive.ObjectID, catcode string) []models.Product {
	start := time.Now()
	col := db.Collection("addons_products")
	var rs []models.Product
	cond := bson.M{"shopid": shopid, "catid": catcode}

	cursor, err := col.Find(ctx, cond)
	r.QueryCount++
	if err = cursor.All(ctx, &rs); err != nil {
		c3mcommon.CheckError("GetProdsByCatId", err)
	}
	r.QueryTime += time.Since(start)
	return rs

}

func (r *Repo) ExportItem(exportitems []models.ExportItem) bool {
	start := time.Now()
	col := db.Collection("addons_products")
	var rs models.Product

	//subcond := bson.M{"$elemMatch": bson.M{"code": itemcode}}
	for _, item := range exportitems {
		cond := bson.M{"shopid": item.ShopId, "code": item.Code}
		err := col.FindOne(ctx, cond).Decode(&rs)
		r.QueryCount++
		for k, v := range rs.Properties {
			if v.Code == item.ItemCode {
				rs.Properties[k].Stock -= item.Num
				r.SaveProd(rs)
				break
			}
		}
		c3mcommon.CheckError("ExportItem", err)

	}
	r.QueryTime += time.Since(start)
	return true

}

//=========================cat function==================
func (r *Repo) SaveCat(cat models.ProdCat) string {
	start := time.Now()
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
		r.QueryCount++
	}
	langinfo, _ := json.Marshal(cat.Langs)
	r.QueryTime += time.Since(start)
	return "{\"Code\":\"" + cat.Code + "\",\"Langs\":" + string(langinfo) + "}"
}
func (r *Repo) GetAllCats(userid, shopid primitive.ObjectID) []models.ProdCat {
	start := time.Now()
	col := db.Collection("addons_prodcats")
	var rs []models.ProdCat
	cond := bson.M{"shopid": shopid}

	cursor, err := col.Find(ctx, cond)
	r.QueryCount++
	if err = cursor.All(ctx, &rs); err != nil {
		c3mcommon.CheckError("GetAllCats", err)
	}
	r.QueryTime += time.Since(start)
	return rs
}

func (r *Repo) GetCats(userid, shopid primitive.ObjectID, ismain bool) []models.ProdCat {
	start := time.Now()
	col := db.Collection("addons_prodcats")
	var rs []models.ProdCat
	cond := bson.M{"shopid": shopid, "main": ismain}

	cursor, err := col.Find(ctx, cond)
	r.QueryCount++
	if err = cursor.All(ctx, &rs); err != nil {
		c3mcommon.CheckError("GetCats", err)
	}
	r.QueryTime += time.Since(start)
	return rs
}

func (r *Repo) GetDemoProdCats() []models.ProdCat {
	start := time.Now()
	col := db.Collection("addons_prodcats")
	shop := r.GetDemoShop()
	var rs []models.ProdCat
	cursor, err := col.Find(ctx, bson.M{"shopid": shop.ID.Hex()})
	r.QueryCount++
	if err = cursor.All(ctx, &rs); err != nil {
		c3mcommon.CheckError("GetCats", err)
	}
	r.QueryTime += time.Since(start)
	return rs
}
func (r *Repo) GetCatByCode(shopid primitive.ObjectID, code string) models.ProdCat {
	start := time.Now()
	col := db.Collection("addons_prodcats")
	var rs models.ProdCat
	cond := bson.M{"shopid": shopid, "code": code}

	err := col.FindOne(ctx, cond).Decode(&rs)
	r.QueryCount++
	c3mcommon.CheckError("GetCatByCode", err)
	r.QueryTime += time.Since(start)
	return rs
}
