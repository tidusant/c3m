package cuahang

import (
	"github.com/tidusant/c3m/common/c3mcommon"
	"github.com/tidusant/c3m/common/log"
	"github.com/tidusant/c3m/repo/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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
func (r *Repo) GetAllPage(shopid, templateid primitive.ObjectID) []models.Page {
	col := db.Collection("pages")
	var rs []models.Page
	cond := bson.M{"shopid": shopid, "templateid": templateid}
	log.Debugf("Getallpage cond %+v", cond)
	cursor, err := col.Find(ctx, cond)
	r.QueryCount++
	if err = cursor.All(ctx, &rs); err != nil {
		c3mcommon.CheckError("GetAllPage Error", err)
	}
	return rs
}
func (r *Repo) GetPageById(pageid primitive.ObjectID) models.Page {
	col := db.Collection("addons_pages")
	var rs models.Page
	cond := bson.D{{"_id", pageid}}

	err := col.FindOne(ctx, cond).Decode(&rs)
	c3mcommon.CheckError("GetPageById", err)
	return rs
}
