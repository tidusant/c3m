package cuahang

//
//import (
//	"encoding/json"
//
//	c3mcommon "colis/common/common"
//	"colis/common/log"
//	"colis/models"
//	//	"c3m/log"
//
//	//"strings"
//
//	"gopkg.in/mgo.v2/bson"
//)
//
////=========================cat function==================
//func SavePage(newitem models.Page) string {
//
//	col := db.C("addons_pages")
//
//	// if prod.Code {
//
//	// 	err := col.Insert(prod)
//	// 	c3mcommon.CheckError("product Insert", err)
//	// } else {
//
//	if len(newitem.Langs) > 0 {
//		if newitem.ID == "" {
//			newitem.ID = bson.NewObjectId()
//		}
//
//		//slug
//		//get all slug
//		slugs := GetAllSlugs( newitem.ShopID)
//		mapslugs := make(map[string]string)
//		for i := 0; i < len(slugs); i++ {
//			mapslugs[slugs[i].Slug] = slugs[i].Slug
//		}
//		for lang, _ := range newitem.Langs {
//			if newitem.Langs[lang].Title != "" {
//				//newslug
//				// tb, _ := lzjs.DecompressFromBase64(newitem.Langs[lang].Title)
//				// newslug := inflect.Parameterize(string(tb))
//				// log.Debugf("title: %s", string(tb))
//				// log.Debugf("newslug: %s", newslug)
//				// newitem.Langs[lang].Slug = newslug
//
//				// //check slug duplicate
//				// i := 1
//				// for {
//				// 	if _, ok := mapslugs[newitem.Langs[lang].Slug]; ok {
//				// 		newitem.Langs[lang].Slug = newslug + strconv.Itoa(i)
//				// 		i++
//				// 	} else {
//				// 		mapslugs[newitem.Langs[lang].Slug] = newitem.Langs[lang].Slug
//				// 		break
//				// 	}
//				// }
//				//remove oldslug
//				log.Debugf("page slug for lang %s,%v", lang, newitem.Langs[lang])
//				tmp := newitem.Langs[lang]
//				tmp.Slug = newitem.Code
//				newitem.Langs[lang] = tmp
//				CreateSlug(newitem.Langs[lang].Slug, newitem.ShopID, "page")
//			} else {
//				delete(newitem.Langs, lang)
//			}
//		}
//
//		_, err := col.UpsertId(newitem.ID, &newitem)
//		c3mcommon.CheckError("news Update", err)
//	} else {
//		col.RemoveId(newitem.ID)
//	}
//
//	//}
//	for lang, _ := range newitem.Langs {
//		tmp := newitem.Langs[lang]
//		tmp.Content = ""
//		newitem.Langs[lang] = tmp
//	}
//	langinfo, _ := json.Marshal(newitem.Langs)
//	return "{\"Code\":\"" + newitem.Code + "\",\"Langs\":" + string(langinfo) + "}"
//}
//func GetAllPage(userid, shopid string) []models.Page {
//	col := db.C("addons_pages")
//	var rs []models.Page
//	shop := GetShopById(userid, shopid)
//	err := col.Find(bson.M{"shopid": shop.ID.Hex(), "publish": true}).All(&rs)
//	c3mcommon.CheckError("get all page", err)
//	return rs
//}
//func GetPageByCode(userid, shopid, code string) models.Page {
//	col := db.C("addons_pages")
//	var rs models.Page
//	cond := bson.M{"shopid": shopid, "code": code}
//	if userid != "594f665df54c58a2udfl54d3er" {
//		cond["userid"] = userid
//	}
//	err := col.Find(cond).One(&rs)
//	c3mcommon.CheckError("getcatbycode", err)
//	return rs
//}
