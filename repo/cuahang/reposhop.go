package cuahang

import (
	//"github.com/spf13/viper"
	"github.com/tidusant/c3m-common/c3mcommon"
	"github.com/tidusant/c3m-common/log"
	"github.com/tidusant/chadmin-repo/models"
	"go.mongodb.org/mongo-driver/bson/primitive" // for BSON ObjectID
	"gopkg.in/mgo.v2/bson"
)

/*for dashboard
=============================================================================
*/

func UpdateTheme(shopid, code string) string {
	col := db.Collection("addons_shops")

	change := bson.M{"$set": bson.M{"theme": code}}
	_, err := col.UpdateOne(ctx, bson.M{"_id": shopid}, change)
	c3mcommon.CheckError("update theme", err)

	return code
}
func LoadShopById(session string, userid, shopid primitive.ObjectID) models.Shop {
	col := db.Collection("addons_userlogin")
	if shopid == primitive.NilObjectID {
		//get first shop
		shopid = GetShopDefault(userid)
	}

	shop := GetShopById(userid, shopid)
	if shop.Name != "" {
		log.Debugf("update login info:shopid %s", shop.ID)
		cond := bson.M{"session": session, "userid": userid}
		change := bson.M{"$set": bson.M{"shopid": shop.ID}}
		col.UpdateOne(ctx, cond, change)
	}
	return shop
}
func GetShopDefault(userid primitive.ObjectID) primitive.ObjectID {
	col := db.Collection("addons_shops")
	var result models.Shop

	col.FindOne(ctx, bson.M{"users": userid}).Decode(&result)
	if result.Name != "" {
		return result.ID
	}

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
	return primitive.NilObjectID
}
func GetShopById(userid, shopid primitive.ObjectID) models.Shop {
	coluser := db.Collection("addons_shops")
	var shop models.Shop
	if shopid == primitive.NilObjectID {
		return shop
	}
	// Create a BSON ObjectID by passing string to ObjectIDFromHex() method
	//shopidObj,_ := primitive.ObjectIDFromHex(shopid)
	//useridObj:= bson.ObjectIdHex(userid)

	cond := bson.M{"_id": shopid}
	//cond := bson.M{"users": userid}
	cond["users"] = userid
	log.Debugf("condition query user %v,%s, %v", cond, userid, shopid)

	err := coluser.FindOne(ctx, cond).Decode(&shop)
	c3mcommon.CheckError("Error query shop in GetShopById", err)
	return shop
}
func GetShopLimitbyKey(shopid primitive.ObjectID, key string) int {

	coluser := db.Collection("shoplimits")

	cond := bson.M{"shopid": shopid, "key": key}
	var rs models.ShopLimit
	err := coluser.FindOne(ctx, cond).Decode(&rs)
	c3mcommon.CheckError("GetShopLimitbyKey :", err)
	return rs.Value
}
func GetShopLimits(shopid primitive.ObjectID) []models.ShopLimit {

	coluser := db.Collection("shoplimits")

	cond := bson.M{"shopid": shopid}
	var rs []models.ShopLimit
	cursor, err := coluser.Find(ctx, cond)
	if err = cursor.All(ctx, &rs); err != nil {
		c3mcommon.CheckError("Update Error:", err)
	}

	return rs
}
func GetOtherShopById(userid, shopid primitive.ObjectID) []models.Shop {
	coluser := db.Collection("addons_shops")
	var shops []models.Shop
	if shopid == primitive.NilObjectID {
		return shops
	}

	cond := bson.M{"_id": bson.M{"$ne": shopid}}

	//if userid != "594f665df54c58a2udfl54d3er" && userid != viper.GetString("config.webuserapi") {
	cond["users"] = userid

	log.Debugf("GetOtherShopById %v %v", cond)
	//}
	cursor, err := coluser.Find(ctx, cond)
	if err = cursor.All(ctx, &shops); err != nil {
		c3mcommon.CheckError("GetOtherShopById", err)
	}
	log.Debugf("GetOtherShopById %v ", shops)
	return shops
}
func GetDemoShop() models.Shop {
	coluser := db.Collection("addons_shops")
	var shop models.Shop
	coluser.FindOne(ctx, bson.M{"name": "demo"}).Decode(&shop)
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
func LoadAllShopAlbums(shopid string) []models.ShopAlbum {
	col := db.Collection("shopalbums")
	var rs []models.ShopAlbum
	cursor, err := col.Find(ctx, bson.M{"shopid": shopid})
	if err = cursor.All(ctx, &rs); err != nil {
		c3mcommon.CheckError("Update Error:", err)
	}

	c3mcommon.CheckError("get ShopAlbum", err)
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
