package cuahang

import (
	//"github.com/spf13/viper"
	"github.com/tidusant/c3m/common/c3mcommon"
	"github.com/tidusant/c3m/common/log"

	"github.com/tidusant/c3m/repo/models"
	"go.mongodb.org/mongo-driver/bson/primitive" // for BSON ObjectID
	"gopkg.in/mgo.v2/bson"
	"time"
)

/*for dashboard
=============================================================================
*/

func (r *Repo) UpdateTheme(shopid, code string) string {
	start := time.Now()
	col := db.Collection("addons_shops")

	change := bson.M{"$set": bson.M{"theme": code}}
	_, err := col.UpdateOne(ctx, bson.M{"_id": shopid}, change)
	r.QueryCount++
	if !c3mcommon.CheckError("update theme", err) {
		code = ""
	}

	r.QueryTime += time.Since(start)
	return code
}
func (r *Repo) LoadShopById(session string, userid, shopid primitive.ObjectID) models.Shop {
	start := time.Now()
	col := db.Collection("addons_userlogin")
	var shop models.Shop
	if shopid == primitive.NilObjectID {
		//get first shop
		shop = r.GetShopDefault(userid)
	} else {
		shop = r.GetShopById(userid, shopid)
	}
	if shop.ID != primitive.NilObjectID {
		cond := bson.M{"session": session, "userid": userid}
		change := bson.M{"$set": bson.M{"shopid": shop.ID}}
		_, err := col.UpdateOne(ctx, cond, change)
		r.QueryCount++
		c3mcommon.CheckError("Update shop to UserLogin", err)
	}
	r.QueryTime += time.Since(start)
	return shop
}

// just get the first query shop for user
func (r *Repo) GetShopDefault(userid primitive.ObjectID) models.Shop {
	start := time.Now()
	col := db.Collection("addons_shops")

	var result models.Shop

	col.FindOne(ctx, bson.M{"users": userid}).Decode(&result)
	r.QueryCount++

	//pipeline := []bson.M{{"$match": bson.M{"name": "abc"}}}
	//col.Pipe(pipeline).All(&result)
	//	for {
	//		if iter.Next(&result) {
	//			log.Printf("result %v", result)
	//			return result.ID.Hex()
	//		} else {
	//			return ""
	//		}
	//	}

	//	if len(result) > 0 {
	//		return result[0].ID.Hex()
	//	}
	r.QueryTime += time.Since(start)
	return result
}
func (r *Repo) GetShopById(userid, shopid primitive.ObjectID) models.Shop {
	start := time.Now()
	coluser := db.Collection("addons_shops")
	var shop models.Shop

	// Create a BSON ObjectID by passing string to ObjectIDFromHex() method
	//shopidObj,_ := primitive.ObjectIDFromHex(shopid)
	//useridObj:= bson.ObjectIdHex(userid)

	cond := bson.M{"_id": shopid}
	//cond := bson.M{"users": userid}
	cond["users"] = userid
	log.Debugf("GetShopById:%s,%s", userid, shopid)

	err := coluser.FindOne(ctx, cond).Decode(&shop)
	r.QueryCount++
	c3mcommon.CheckError("Error query shop in GetShopById", err)
	r.QueryTime += time.Since(start)
	return shop
}
func (r *Repo) GetShopLimitbyKey(shopid primitive.ObjectID, key string) int {
	start := time.Now()
	coluser := db.Collection("shoplimits")

	cond := bson.M{"shopid": shopid, "key": key}
	var rs models.ShopLimit
	err := coluser.FindOne(ctx, cond).Decode(&rs)
	r.QueryCount++
	c3mcommon.CheckError("GetShopLimitbyKey :", err)
	r.QueryTime += time.Since(start)
	return rs.Value
}
func (r *Repo) GetShopLimits(shopid primitive.ObjectID) []models.ShopLimit {
	start := time.Now()
	coluser := db.Collection("shoplimits")

	cond := bson.M{"shopid": shopid}
	var rs []models.ShopLimit
	cursor, err := coluser.Find(ctx, cond)
	r.QueryCount++
	if err = cursor.All(ctx, &rs); err != nil {
		c3mcommon.CheckError("Update Error:", err)
	}
	r.QueryTime += time.Since(start)
	return rs
}
func (r *Repo) GetShopsByUserId(userid primitive.ObjectID) []models.Shop {
	start := time.Now()
	coluser := db.Collection("addons_shops")
	var shops []models.Shop

	cond := bson.M{"users": userid}
	//if userid != "594f665df54c58a2udfl54d3er" && userid != viper.GetString("config.webuserapi") {

	//}
	cursor, err := coluser.Find(ctx, cond)
	r.QueryCount++
	if err = cursor.All(ctx, &shops); err != nil {
		c3mcommon.CheckError("GetShopsByUserId", err)
	}

	r.QueryTime += time.Since(start)
	return shops
}
func (r *Repo) GetDemoShop() models.Shop {
	start := time.Now()
	coluser := db.Collection("addons_shops")
	var shop models.Shop
	coluser.FindOne(ctx, bson.M{"name": "demo"}).Decode(&shop)
	r.QueryCount++
	r.QueryTime += time.Since(start)
	return shop
}

// func SaveCat(userid, shopid string, cat models.ProdCat) string {

// 	shop := GetShopById(userid, shopid)
// 	newcat := false
// 	if cat.Code == "" {
// 		newcat = true
// 	}
// 	//get all cats
// 	cats := GetAllCats(userid, shopid)
// 	var oldcat models.ProdCat
// 	//check max cat limited
// 	if shop.Config.MaxCat <= len(cats) && newcat {
// 		return "-1"
// 	}
// 	//get array of album slug
// 	catslugs := map[string]string{}
// 	catcodes := map[string]string{}
// 	for _, c := range cats {
// 		catcodes[c.Code] = c.Code
// 		for _, ci := range c.Langs {
// 			catslugs[ci.Slug] = ci.Slug
// 		}
// 		if newcat && c.Code == cat.Code {
// 			oldcat = c
// 		}
// 	}

// 	for lang, _ := range cat.Langs {
// 		if cat.Langs[lang].Name == "" {
// 			delete(cat.Langs, lang)
// 			continue
// 		}
// 		//newslug
// 		i := 1
// 		newslug := inflect.Parameterize(cat.Langs[lang].Name)
// 		cat.Langs[lang].Slug = newslug
// 		//check slug duplicate
// 		for {
// 			if _, ok := catslugs[cat.Langs[lang].Slug]; ok {
// 				cat.Langs[lang].Slug = newslug + strconv.Itoa(i)
// 				i++
// 			} else {
// 				catslugs[cat.Langs[lang].Slug] = cat.Langs[lang].Slug
// 				break
// 			}
// 		}
// 	}

// 	//check code duplicate
// 	if newcat {
// 		//insert new
// 		newcode := ""
// 		for {
// 			newcode = mystring.RandString(3)
// 			if _, ok := catcodes[newcode]; !ok {
// 				break
// 			}
// 		}
// 		cat.Code = newcode
// 		cat.UserId = userid
// 		cat.Created = time.Now().UTC().Add(time.Hour + 7)
// 	} else {
// 		//update
// 		oldcat.Langs = cat.Langs
// 		cat = oldcat
// 	}

// 	UpdateCat(shop)
// 	return cat.Code
// }

//func SaveCat(userid, shopid, code string, catinfo models.ShopCatInfo) string {

//	//slug
//	rt := "0"
//	i := 1
//	slug := inflect.Parameterize(catinfo.Name)
//	catinfo.Slug = slug
//	shop := GetShopById(userid, shopid)

//	//get array of album slug
//	catslugs := map[string]string{}
//	for _, c := range shop.ShopCats {
//		for _, ci := range c.Langs {
//			if ci.Slug != catinfo.Slug {
//				catslugs[ci.Slug] = ci.Slug
//			}
//		}

//	}

//	for {
//		if _, ok := catslugs[catinfo.Slug]; ok {
//			catinfo.Slug = slug + strconv.Itoa(i)
//			i++
//		} else {
//			break
//		}
//	}

//	for i, _ := range shop.ShopCats {
//		if shop.ShopCats[i].Code == code && shop.ShopCats[i].UserId == userid {
//			isnewlang := true
//			for j, _ := range shop.ShopCats[i].Langs {
//				if shop.ShopCats[i].Langs[j].Lang == catinfo.Lang {
//					//shop.ShopCats[i].Langs[j] = catinfo
//					isnewlang = false
//					break
//				}
//			}
//			if isnewlang {
//				//shop.ShopCats[i].Langs = append(shop.ShopCats[i].Langs, catinfo)

//			}
//			rt = "1"
//			break
//		}
//	}
//	UpdateCat(shop)
//	return rt

//}
// func SaveShopConfig(shop models.Shop) models.Shop {
// 	coluser := db.C("addons_shops")

// 	cond := bson.M{"_id": shop.ID}
// 	change := bson.M{"$set": bson.M{"config": shop.Config}}

// 	coluser.Update(cond, change)
// 	return shop
// }
func (r *Repo) LoadAllShopAlbums(shopid string) []models.ShopAlbum {
	start := time.Now()
	col := db.Collection("shopalbums")
	var rs []models.ShopAlbum
	cursor, err := col.Find(ctx, bson.M{"shopid": shopid})
	r.QueryCount++
	if err = cursor.All(ctx, &rs); err != nil {
		c3mcommon.CheckError("Update Error:", err)
	}

	c3mcommon.CheckError("get ShopAlbum", err)
	r.QueryTime += time.Since(start)
	return rs
}

//
//func SaveAlbum(album models.ShopAlbum) models.ShopAlbum {
//	coluser := db.Collection("shopalbums")
//	if album.ID.Hex() == "" {
//		album.ID = bson.NewObjectId()
//		album.Created = time.Now()
//	}
//
//	_, err := coluser.UpsertId(album.ID, album)
//	c3mcommon.CheckError("SaveAlbum", err)
//	return album
//}
//func UpdateAlbum(shop models.Shop) models.Shop {
//	coluser := db.C("addons_shops")
//
//	cond := bson.M{"_id": shop.ID}
//	change := bson.M{"$set": bson.M{"albums": shop.Albums}}
//
//	coluser.Update(cond, change)
//	return shop
//}
//func SaveShopConfig(shop models.Shop) {
//	col := db.C("addons_shops")
//	//check  exist:
//	cond := bson.M{"_id": shop.ID}
//	change := bson.M{"$set": bson.M{"config": shop.Config}}
//	err := col.Update(cond, change)
//
//	c3mcommon.CheckError("SaveShopConfig :", err)
//
//}
