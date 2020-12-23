package cuahang

//
//func SaveNews(newitem models.News) string {
//
//	col := db.C("addons_news")
//
//	// if prod.Code {
//
//	// 	err := col.Insert(prod)
//	// 	c3mcommon.CheckError("product Insert", err)
//	// } else {
//	if len(newitem.Langs) > 0 {
//		if newitem.ID == "" {
//			newitem.ID = bson.NewObjectId()
//		}
//		_, err := col.UpsertId(newitem.ID, &newitem)
//		c3mcommon.CheckError("news Update", err)
//	} else {
//		col.RemoveId(newitem.ID)
//	}
//	//}
//	for lang, _ := range newitem.Langs {
//		newitem.Langs[lang].Content = ""
//	}
//	langinfo, _ := json.Marshal(newitem.Langs)
//	return "{\"Code\":\"" + newitem.Code + "\",\"Langs\":" + string(langinfo) + "}"
//}
//
//func GetAllNews(userid, shopid string) []models.News {
//	col := db.C("addons_news")
//	var rs []models.News
//	shop := GetShopById(userid, shopid)
//	err := col.Find(bson.M{"shopid": shop.ID.Hex()}).All(&rs)
//	c3mcommon.CheckError("get all news", err)
//	return rs
//}
//func GetDemoNews() []models.News {
//	col := db.C("addons_news")
//	var rs []models.News
//	shop := GetDemoShop()
//	err := col.Find(bson.M{"shopid": shop.ID.Hex()}).All(&rs)
//	c3mcommon.CheckError("get all demo news", err)
//	return rs
//}
//
//func GetNewsByCode(userid, shopid, code string) models.News {
//	col := db.C("addons_news")
//	var rs models.News
//	cond := bson.M{"shopid": shopid, "code": code}
//	if userid != "594f665df54c58a2udfl54d3er" {
//		cond["userid"] = userid
//	}
//	err := col.Find(cond).One(&rs)
//	c3mcommon.CheckError("getnewbycode", err)
//
//	return rs
//
//}
//func GetNewsByCatId(userid, shopid, catcode string) []models.Product {
//	col := db.C("addons_news")
//	var rs []models.Product
//
//	err := col.Find(bson.M{"userid": userid, "shopid": shopid, "catid": catcode}).All(&rs)
//	c3mcommon.CheckError("getprod", err)
//
//	return rs
//
//}
//
////=========================cat function==================
//func SaveNewsCat(cat models.NewsCat) string {
//	col := db.C("addons_newscats")
//	if len(cat.Langs) > 0 {
//		if cat.ID == "" {
//			cat.ID = bson.NewObjectId()
//		}
//		col.UpsertId(cat.ID, cat)
//	} else {
//		col.RemoveId(cat.ID)
//	}
//	langinfo, _ := json.Marshal(cat.Langs)
//	return "{\"Code\":\"" + cat.Code + "\",\"Langs\":" + string(langinfo) + "}"
//}
//func GetDemoNewsCats() []models.NewsCat {
//	col := db.C("addons_newscats")
//	shop := GetDemoShop()
//	var rs []models.NewsCat
//	err := col.Find(bson.M{"shopid": shop.ID.Hex()}).All(&rs)
//	c3mcommon.CheckError("getcatprod", err)
//	return rs
//}
//func GetAllNewsCats(userid, shopid string) []models.NewsCat {
//	col := db.C("addons_newscats")
//	var rs []models.NewsCat
//	cond := bson.M{"shopid": shopid}
//	if userid != "594f665df54c58a2udfl54d3er" {
//		cond["userid"] = userid
//	}
//	err := col.Find(cond).Sort("-created").All(&rs)
//	c3mcommon.CheckError("getcat ", err)
//	return rs
//}
//
//func GetNewsCatByCode(userid, shopid, code string) models.NewsCat {
//	col := db.C("addons_newscats")
//	var rs models.NewsCat
//	err := col.Find(bson.M{"shopid": shopid, "code": code}).One(&rs)
//	c3mcommon.CheckError("getcatbycode", err)
//	return rs
//}
