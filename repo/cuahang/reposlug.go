package cuahang

//
//import (
//	c3mcommon "colis/common/common"
//	"colis/common/log"
//	"colis/models"
//	"gopkg.in/mgo.v2/bson"
//	"strconv"
//)
//
//
//func SaveSlug(newslug models.Slug) string {
//	col := db.C("addons_slugs")
//
//	//check exist slug
//	cond := bson.M{"shopid": newslug.ShopId,
//		"object":   newslug.Object,
//		"objectid": newslug.ObjectId,
//		"lang":     newslug.Lang,
//		"slug":     newslug.Slug,
//	}
//	var oldslug models.Slug
//	col.Find(cond).One(&oldslug)
//	if oldslug.ID.Hex() != "" {
//		//just rebuild
//		newslug.ID = oldslug.ID
//	} else {
//		//no need to check slug if page ==> override slug
//		if newslug.Object != "page" {
//			//get all slug
//			allslugs := GetAllSlugs(newslug.ShopId)
//			mapslugs := make(map[string]models.Slug)
//			for _, slug := range allslugs {
//				mapslugs[slug.Slug] = slug
//			}
//			//check slug duplicate
//			i := 1
//			slugname := newslug.Slug
//			log.Debugf("checkslug %s", newslug)
//			for {
//				if oldslug, ok := mapslugs[newslug.Slug]; ok {
//					//check old slug
//					log.Debugf("oldslug %s", oldslug)
//					if oldslug.Object == newslug.Object && oldslug.ObjectId == newslug.ObjectId && oldslug.ShopId == newslug.ShopId && oldslug.Lang == newslug.Lang {
//						break
//					}
//					newslug.Slug = slugname + strconv.Itoa(i)
//					i++
//				} else {
//					log.Debugf("newslug %s", oldslug)
//					break
//				}
//			}
//		}
//		if newslug.ID.Hex() == "" {
//			newslug.ID = bson.NewObjectId()
//		}
//	}
//
//	_, err := col.UpsertId(newslug.ID, &newslug)
//	c3mcommon.CheckError("SaveSlug "+newslug.Slug, err)
//	return newslug.Slug
//}
//
//func SaveSlugNoBuild(newslug models.Slug) string {
//	col := db.C("addons_slugs")
//
//	//check exist slug
//	cond := bson.M{"shopid": newslug.ShopId,
//		"object":   newslug.Object,
//		"objectid": newslug.ObjectId,
//		"lang":     newslug.Lang,
//		"slug":     newslug.Slug,
//	}
//	var oldslug models.Slug
//	col.Find(cond).One(&oldslug)
//	if oldslug.ID.Hex() != "" {
//		//just rebuild
//		newslug.ID = oldslug.ID
//	} else {
//		//no need to check slug if page ==> override slug
//		if newslug.Object != "home" {
//			//get all slug
//			allslugs := GetAllSlugs(newslug.ShopId)
//			mapslugs := make(map[string]models.Slug)
//			for _, slug := range allslugs {
//				mapslugs[slug.Slug] = slug
//			}
//			//check slug duplicate
//			i := 1
//			slugname := newslug.Slug
//			log.Debugf("checkslug %s", newslug)
//			for {
//				if oldslug, ok := mapslugs[newslug.Slug]; ok {
//					//check old slug
//					log.Debugf("oldslug %s", oldslug)
//					if oldslug.Object == newslug.Object && oldslug.ObjectId == newslug.ObjectId && oldslug.ShopId == newslug.ShopId && oldslug.Lang == newslug.Lang {
//						break
//					}
//					newslug.Slug = slugname + strconv.Itoa(i)
//					i++
//				} else {
//					log.Debugf("newslug %s", oldslug)
//					break
//				}
//			}
//		}
//		if newslug.ID.Hex() == "" {
//			newslug.ID = bson.NewObjectId()
//		}
//	}
//
//	_, err := col.UpsertId(newslug.ID, &newslug)
//	c3mcommon.CheckError("SaveSlug "+newslug.Slug, err)
//	return newslug.Slug
//}
//
//func RemoveSlug(slug string, shopid string) bool {
//	col := db.C("addons_slugs")
//	cond := bson.M{"slug": slug, "shopid": shopid}
//	err := col.Remove(cond)
//	return c3mcommon.CheckError("RemoveSlug "+slug, err)
//}
//
//func GetAllSlugs(shopid string) []models.Slug {
//	col := db.C("addons_slugs")
//	var rs []models.Slug
//	cond := bson.M{"shopid": shopid}
//	err := col.Find(cond).All(&rs)
//	c3mcommon.CheckError("GetAllSlugs", err)
//	return rs
//}
//
//func CreateSlug(slug, shopid, object string) bool {
//	col := db.C("addons_slugs")
//	col.Insert(bson.M{"shopid": shopid, "slug": slug, "object": object})
//	return true
//}
